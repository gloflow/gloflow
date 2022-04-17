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
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_job_update struct {
	Block_num_indexed_int uint64 `json:"block_num_indexed_int"`
	Txs_num_indexed_int   uint64 `json:"txs_num_indexed_int"`
}
type GF_job_updates_ch  chan(GF_job_update)
type GF_job_err_ch      chan(gf_core.GF_error)
type GF_job_complete_ch chan(bool)


type GF_job_update_new_consumer_ch chan(GF_job_update_new_consumer)
type GF_job_update_new_consumer struct {
	Job_id_str  GF_indexer_job_id
	Response_ch chan(GF_job_update_new_consumer_response)
	ctx         context.Context
}

type GF_job_update_new_consumer_response struct {
	Job_updates_ch  GF_job_updates_ch
	Job_err_ch      GF_job_err_ch
	Job_complete_ch GF_job_complete_ch
}

//-------------------------------------------------
func Updates__consume_stream(p_job_id_str GF_indexer_job_id,
	p_ctx        context.Context,
	p_aws_client *sqs.Client,
	p_runtime    *gf_eth_core.GF_runtime) (GF_job_updates_ch, GF_job_err_ch, GF_job_complete_ch) {

	
	// CONSUMER_CHANNELS - channels returned to the consumer to consume job_update messages from.
	consumer__job_updates_ch  := make(GF_job_updates_ch, 100)
	consumer__job_err_ch      := make(GF_job_err_ch, 1)
	consumer__job_complete_ch := make(GF_job_complete_ch, 1)
	

	go func() {

		//-------------------------------------------------
		complete_job_with_error_fn := func(p_gf_err *gf_core.GF_error) {
			consumer__job_err_ch <- *p_gf_err
		}

		//-------------------------------------------------

		// GET_QUEUE
		queue_name_str := updates__get_queue_name(p_job_id_str)
		queue_info, gf_err := gf_aws.SQS_get_queue_info(queue_name_str,
			p_aws_client,
			p_ctx,
			p_runtime.Runtime_sys)
		if gf_err != nil {
			complete_job_with_error_fn(gf_err)
			return	
		}

		errs_num_int := 0
		for {

			// SQS_MSG_PULL
			msg_map, gf_err := gf_aws.SQS_msg_pull(queue_info,
				p_aws_client,
				p_ctx,
				p_runtime.Runtime_sys)

			if gf_err != nil {

				// IMPORTANT!! - fail job if a certain number of errors
				//               is encountered. to avoid cases where only a few messages
				//               happen in a queue where there are a lot of updates.
				if errs_num_int < 3 {
					complete_job_with_error_fn(gf_err)
					return
				}

				// continue processing messages and increase the error count
				errs_num_int += 1
				continue
			}

			// send the job_update message to the job_update consumer
			//
			// IMPORTANT!! - if a msg_map is nil it means that SQS_msg_pull() timed with no 
			//               message in the queue and returned nil.
			if msg_map != nil {
				job_update_msg := GF_job_update{
					Block_num_indexed_int: uint64(msg_map["block_num_indexed_int"].(float64)),
				}
				consumer__job_updates_ch <- job_update_msg
			}
		}
	}()
	
	return consumer__job_updates_ch, consumer__job_err_ch, consumer__job_complete_ch
}

//-------------------------------------------------
func Updates__init_stream(p_job_id_str GF_indexer_job_id,
	p_ctx        context.Context,
	p_sqs_client *sqs.Client,
	p_runtime    *gf_eth_core.GF_runtime) (*gf_aws.GF_SQS_queue, *gf_core.GF_error) {

	queue_name_str := updates__get_queue_name(p_job_id_str)
	
	// QUEUE_CREATE
	queue, gf_err := gf_aws.SQS_queue_create(queue_name_str,
		p_sqs_client,
		p_ctx,
		p_runtime.Runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	spew.Dump(queue)

	return queue, nil
}

//-------------------------------------------------
func updates__get_queue_name(p_job_id_str GF_indexer_job_id) string {
	queue_name_str := fmt.Sprintf(fmt.Sprintf("gf_eth_indexer__job_updates__%s", p_job_id_str))
	return queue_name_str
}