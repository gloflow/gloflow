/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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
	"net/http"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//-------------------------------------------------
func Init_handlers(p_indexer_cmds_ch GF_indexer_ch,
	p_indexer_job_updates_new_consumer_ch GF_job_update_new_consumer_ch,
	p_metrics                             *gf_eth_core.GF_metrics,
	p_runtime                             *gf_eth_core.GF_runtime) {

	//---------------------
	// GET__BLOCK_INDEX__JOB_UPDATES

	gf_rpc_lib.SSE_create_handler__http("/gfethm/v1/block/index/job_updates",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (gf_rpc_lib.SSE_data_update_ch, gf_rpc_lib.SSE_data_err_ch, gf_rpc_lib.SSE_data_complete_ch, *gf_core.GF_error) {


			span__root := sentry.StartSpan(p_ctx, "http__master__block_index_job_updates", sentry.ContinueFromRequest(p_req))
			// ctx := span__root.Context()
			defer span__root.Finish()


			job_id_str := p_req.URL.Query().Get("job_id")

			job_updates_ch, job_err_ch, job_complete_ch := Client__new_consumer(GF_indexer_job_id(job_id_str),
				p_indexer_job_updates_new_consumer_ch,
				p_ctx)

			//---------------------
			// IMPORTANT!! - casting message from the indexer update format
			//               to the general SSE update format (interface{}).
			//               not very efficient because it adds an extra message
			//               relaying stage.
			data_updates_ch  := make(chan interface{}, 10)
			data_err_ch      := make(gf_rpc_lib.SSE_data_err_ch)
			data_complete_ch := make(chan bool)
			
			go func() {
				for {
					select {
					case update_msg := <- job_updates_ch:
						data_updates_ch <- interface{}(update_msg)
					
					case gf_err := <- job_err_ch:
						data_err_ch <- gf_err

					case complete_bool := <- job_complete_ch:
						data_complete_ch <- complete_bool
					}
				}
			}()

			//---------------------

			span__root.Finish()

			return data_updates_ch, data_err_ch, data_complete_ch, nil
		},
		true, // p_store_run_bool
		p_runtime.RuntimeSys)

	//---------------------
	// GET__BLOCK_INDEX

	gf_rpc_lib.Create_handler__http("/gfethm/v1/block/index",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			span__root := sentry.StartSpan(p_ctx, "http__master__block_index", sentry.ContinueFromRequest(p_req))
			// ctx := span__root.Context()
			defer span__root.Finish()

			//------------------
			// INPUT

			block_start_uint, block_end_uint, gf_err := gf_eth_core.Http__get_arg__block_range(p_resp, p_req, p_runtime.RuntimeSys)
			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			
			job_id_str, gf_err := Client__index_block_range(block_start_uint,
				block_end_uint,
				p_indexer_cmds_ch)
			if gf_err != nil {
				return nil, gf_err
			}
			
			
			// // ABI_DEFS
			// abis_defs_map, gf_err := gf_eth_core.Eth_abi__get_defs(ctx, p_metrics, p_runtime)
			// if gf_err != nil {
			// 	return nil, gf_err
			// }
			//
			// span__pipeline := sentry.StartSpan(ctx, "blocks_persist_bulk")
			//
			// gf_errs_lst := gf_eth_blocks.Eth_blocks__get_and_persist_bulk__pipeline(block_start_uint,
			// 	block_end_uint,
			// 	// p_get_hosts_fn,
			// 	// abis_defs_map,
			// 	// span__pipeline.Context(),
			// 	p_metrics,
			// 	p_runtime)
			//
			// span__pipeline.Finish()
			//
			// if len(gf_errs_lst) > 0 {
			//
			// 	// FIX!! - see if all errors should be returned maybe.
			// 	gf_err__first := gf_errs_lst[0]
			// 	return nil, gf_err__first
			// }

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"job_id_str": job_id_str,
			}

			//------------------
			span__root.Finish()

			return data_map, nil
		},
		p_runtime.RuntimeSys)

	//---------------------
}