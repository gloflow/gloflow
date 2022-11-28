/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_eth_indexer

import (
	"fmt"
	"time"
	"context"
	"strings"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_contract"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_blocks"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_worker"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_indexer_job_id string

type GF_indexer_ch chan(GF_indexer_cmd)
type GF_indexer_cmd struct {
	Block_start_uint uint64
	Block_end_uint   uint64
	Response_ch      chan(GF_indexer_job_id)
	Response_err_ch  chan(gf_core.GFerror)
}

//-------------------------------------------------
func Init(p_get_worker_hosts_fn gf_eth_worker.Get_worker_hosts_fn,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) (GF_indexer_ch, GF_job_update_new_consumer_ch, *gf_core.GFerror) {
	
		
	// AWS_CLIENT
	sqs_client, gfErr := gf_aws.SQSinit(p_runtime.RuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	// incoming commands to begin indexing jobs
	indexer_cmds_ch                     := make(GF_indexer_ch, 100)
	indexer_job_updates_new_consumer_ch := make(chan GF_job_update_new_consumer, 10)


	go func() {

		// IMPORTANT!! - sentry needs a new hub per go-routine, so cloning the toplevel hub
		hub := sentry.CurrentHub().Clone()
		
		
		
		for {
			select {

			//----------------------------
			// INDEXER_COMMANDS
			case cmd := <- indexer_cmds_ch:

				// run job in a new go-routine to be able to handle other messages
				// while completion of this job.
				go func() {

					job_id_str := job_get_id()

					
					// IMPORTANT!! - using a background context, and not a client supplied context
					//               (via cmd.Ctx) because clients just submit an index operation,
					//               and continue their work (or get response to their request). 
					//               the index op should complete independently of the client, in the future.
					ctx_bg := context.Background()

					// IMPORTANT!! - associate context for this job with the hub for this job processing goroutine
					ctx := sentry.SetHubOnContext(ctx_bg, hub)
					
					hub.Scope().SetTag("job_id", string(job_id_str))

					// TRACE
					// span has to be started and its context passed to job_run
					// so that all the subsequent nested sentry calls dont 
					// fail with nil exception.
					span__root := sentry.StartSpan(ctx, "indexer_job")
					defer span__root.Finish()

					gfErr := job_run(job_id_str,
						cmd,
						p_get_worker_hosts_fn,
						span__root.Context(),
						sqs_client,
						p_metrics,
						p_runtime)
					if gfErr != nil {
						cmd.Response_err_ch <- *gfErr
						return
					}

					span__root.Finish()

					cmd.Response_ch <- job_id_str
				}()
			
			//----------------------------

			case new_consumer_msg := <- indexer_job_updates_new_consumer_ch:

				fmt.Println(new_consumer_msg)



				job_id_str   := new_consumer_msg.Job_id_str
				ctx_consumer := new_consumer_msg.ctx
				consumer__job_updates_ch, consumer__job_err_ch, consumer__job_complete_ch := Updates__consume_stream(job_id_str,
					ctx_consumer,
					sqs_client,
					p_runtime)

				response := GF_job_update_new_consumer_response{
					Job_updates_ch:  consumer__job_updates_ch,
					Job_err_ch:      consumer__job_err_ch,
					Job_complete_ch: consumer__job_complete_ch,
				}
				new_consumer_msg.Response_ch <- response
				

			//----------------------------
			}
		}
	}()
	
	return indexer_cmds_ch, indexer_job_updates_new_consumer_ch, nil
}

//-------------------------------------------------
func job_get_id() GF_indexer_job_id {
	job_start_time_f := float64(time.Now().UnixNano())/1000000000.0
	job_id_str       := GF_indexer_job_id(fmt.Sprintf("jid_%s", strings.ReplaceAll(fmt.Sprintf("%f", job_start_time_f), ".", "_")))
	return job_id_str
}

//-------------------------------------------------
func job_run(p_job_id_str GF_indexer_job_id,
	p_cmd                 GF_indexer_cmd,
	p_get_worker_hosts_fn gf_eth_worker.Get_worker_hosts_fn,
	p_ctx                 context.Context,
	p_sqs_client          *sqs.Client,
	p_metrics             *gf_eth_core.GF_metrics,
	p_runtime             *gf_eth_core.GF_runtime) *gf_core.GFerror {

	job_updates_ch  := make(GF_job_updates_ch, 10)
	job_complete_ch := make(GF_job_complete_ch, 1)


	

	//----------------------------
	// ABI_DEFS
	abis_defs_map, gfErr := gf_eth_contract.Eth_abi__get_defs(p_ctx, p_metrics, p_runtime)
	if gfErr != nil {
		return gfErr
	}

	//----------------------------



	
	gf_sqs_queue, gfErr := Updates__init_stream(p_job_id_str,
		p_ctx,
		p_sqs_client,
		p_runtime)
	if gfErr != nil {
		return gfErr
	}

	// process indexing job in separate go-routine so that updates/completion
	// messages can be processed in parallel.
	go func() {
		// PERSIST_RANGE
		gfErrsLst := index__range(p_cmd.Block_start_uint,
			p_cmd.Block_end_uint,
			p_get_worker_hosts_fn,
			abis_defs_map,

			job_updates_ch,
			job_complete_ch,
			p_ctx,
			p_metrics,
			p_runtime)
		
		if len(gfErrsLst) > 0 {

		}
	}()

	//----------------------------
	// UPDATES - push to SQS queue
	go func() {
		for {
			select {
			case update_msg := <- job_updates_ch:
				
				gfErr := gf_aws.SQSmsgPush(interface{}(update_msg),
					gf_sqs_queue,
					p_sqs_client,
					p_ctx,
					p_runtime.RuntimeSys)
				if gfErr != nil {

				}

			case complete_bool := <- job_complete_ch:

				if complete_bool {
					gfErr := gf_aws.SQSqueueDelete(gf_sqs_queue.Name_str,
						p_sqs_client,
						p_ctx,
						p_runtime.RuntimeSys)
					if gfErr != nil {
						break
					}
				}
			}
		}
	}()

	//----------------------------
	return nil
}

//-------------------------------------------------
func index__range(p_block_start_uint uint64,
	p_block_end_uint      uint64,
	p_get_worker_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string,
	p_abis_defs_map       map[string]*gf_eth_contract.GF_eth__abi,

	p_job_updates_ch      GF_job_updates_ch,
	p_job_complete_ch     GF_job_complete_ch,
	p_ctx                 context.Context,
	p_metrics             *gf_eth_core.GF_metrics,
	p_runtime             *gf_eth_core.GF_runtime) []*gf_core.GFerror {

	gfErrsLst := []*gf_core.GFerror{}
	for b := p_block_start_uint; b <= p_block_end_uint; b++ {

		block_uint := b

		txs_num_int, gfErr := gf_eth_blocks.Index__pipeline(block_uint,
			p_get_worker_hosts_fn,
			p_abis_defs_map,
			p_ctx,
			p_metrics,
			p_runtime)
		if gfErr != nil {
			gfErrsLst = append(gfErrsLst, gfErr)
			continue // continue processing subsequent blocks
		}

		// JOB_UPDATE
		if p_job_updates_ch != nil {
			p_job_updates_ch <- GF_job_update{
				Block_num_indexed_int: block_uint,
				Txs_num_indexed_int:   txs_num_int,
			}
		}
	}

	// JOB_COMPLETE
	if p_job_complete_ch != nil {
		p_job_complete_ch <- true
	}
	
	return gfErrsLst
}