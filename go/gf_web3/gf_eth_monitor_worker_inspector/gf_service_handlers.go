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

package main

import (
	"fmt"
	"net/http"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_tx"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_blocks"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func init_handlers(p_metrics *GF_metrics,
	p_runtime *GF_runtime) {


	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/gfethm_worker_inspect/v1/account/info",
		"/gfethm_worker_inspect/v1/tx/trace",
		"/gfethm_worker_inspect/v1/blocks",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr,
		"gf_web3_monitor_worker_inspector",
		handlers_endpoints_lst)

	//---------------------
	// GET_ACCOUNT_INFO

	gf_rpc_lib.Create_handler__http_with_metrics("/gfethm_worker_inspect/v1/account/info",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {


			//------------------
			// INPUT

			account_address_hex_str, gf_err := gf_eth_core.Http__get_arg__acc_address_hex(p_req, p_runtime.runtime_sys)

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------


			fmt.Println(account_address_hex_str)

			return nil, nil

		},
		metrics,
		false, // true, // pStoreRunBool
		p_runtime.runtime_sys)

	//---------------------
	// GET_TX_TRACE

	gf_rpc_lib.Create_handler__http_with_metrics("/gfethm_worker_inspect/v1/tx/trace",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			span__root := sentry.StartSpan(p_ctx, "http__worker_inspector__get_tx_trace", sentry.ContinueFromRequest(p_req))
			defer span__root.Finish()

			//------------------
			// INPUT
			tx_hex_str, gf_err := gf_eth_core.Http__get_arg__tx_id_hex(p_resp, p_req, p_runtime.runtime_sys)

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			// GET_TRACE
			trace_map, gf_err := gf_eth_tx.Trace__get(tx_hex_str,
				p_runtime.eth_rpc_host_str,
				span__root.Context(),
				p_runtime.runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}
			
			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"trace_map": trace_map,
			}

			//------------------

			return data_map, nil
		},
		metrics,
		false, // true, // pStoreRunBool
		p_runtime.runtime_sys)

	//---------------------
	// GET_BLOCKS

	gf_rpc_lib.Create_handler__http_with_metrics("/gfethm_worker_inspect/v1/blocks",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			span__root := sentry.StartSpan(p_ctx, "http__worker_inspector__get_block", sentry.ContinueFromRequest(p_req))
			defer span__root.Finish()
			 
			// METRICS
			if p_metrics != nil {
				p_metrics.counter__http_req_num__get_blocks.Inc()
			}

			//------------------
			// INPUT

			block_num_int, gf_err := gf_eth_core.Http__get_arg__block_num(p_resp, p_req, p_runtime.runtime_sys)

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			// GET_BLOCK

			span__pipeline := sentry.StartSpan(span__root.Context(), "eth_rpc__get_block__pipeline")
			defer span__pipeline.Finish() // in case a panic happens before the main .Finish() for this span

			gf_block, gf_err := gf_eth_blocks.Get__pipeline(block_num_int,
				p_runtime.eth_rpc_client,
				span__pipeline.Context(),
				p_runtime.py_plugins,
				p_runtime.runtime_sys)

			span__pipeline.Finish()

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"block_map": gf_block, // spew.Sdump(),
			}

			//------------------

			return data_map, nil
		},
		metrics,
		false, // true, // pStoreRunBool
		p_runtime.runtime_sys)

	//---------------------
	// HEALTH
	http.HandleFunc("/gfethm_worker_inspect/v1/health", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_resp.Write([]byte("ok"))

	})

	//---------------------
}