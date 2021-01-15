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
	// "fmt"
	"net/http"
	// "context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_lib"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func init_handlers(p_metrics *GF_metrics,
	p_runtime *GF_runtime) {

	//---------------------
	// GET_BLOCKS

	gf_rpc_lib.Create_handler__http("/gfethm_worker_inspect/v1/blocks",
		func(p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			ctx := p_req.Context()
			hub := sentry.GetHubFromContext(ctx)
			hub.Scope().SetTag("url", p_req.URL.Path)
			// hub.Scope().SetTransaction("http__worker_inspector__get_block") // set custom transaction name

			// METRICS
			if p_metrics != nil {
				p_metrics.counter__http_req_num__get_blocks.Inc()
			}

			//------------------
			// INPUT

			block_num_int, gf_err := gf_eth_monitor_lib.Http__get_arg__block_num(p_resp, p_req, p_runtime.runtime_sys)

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			// GET_BLOCK

			span_pipeline := sentry.StartSpan(ctx, "eth_rpc__get_block__pipeline")

			gf_block, gf_err := gf_eth_monitor_lib.Eth_rpc__get_block__pipeline(block_num_int,
				p_runtime.eth_rpc_client,
				ctx,
				p_runtime.runtime_sys)

			span_pipeline.Finish()

			if gf_err != nil {
				
				// gf_rpc_lib.Error__in_handler("/gfethm_worker_inspect/v1/blocks",
				// 	fmt.Sprintf("failed to get block - %d", block_num_int),
				// 	gf_err, p_resp, p_runtime.runtime_sys)
				return nil, gf_err
			}

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"block": gf_block, // spew.Sdump(),
			}
			return data_map, nil
			// gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime.runtime_sys)

			//------------------
		},
		p_runtime.runtime_sys)

	//---------------------
	// HEALTH
	http.HandleFunc("/gfethm_worker_inspect/v1/health", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_resp.Write([]byte("ok"))

	})

	//---------------------
}