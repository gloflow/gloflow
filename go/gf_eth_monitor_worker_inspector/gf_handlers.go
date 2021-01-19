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
	"context"
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
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			

			span__root := sentry.StartSpan(p_ctx, "http__worker_inspector__get_block", sentry.ContinueFromRequest(p_req))
			defer span__root.Finish()
			 
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

			span__pipeline := sentry.StartSpan(span__root.Context(), "eth_rpc__get_block__pipeline")
			defer span__pipeline.Finish() // in case a panic happens before the main .Finish() for this span

			gf_block, gf_err := gf_eth_monitor_lib.Eth_rpc__get_block__pipeline(block_num_int,
				p_runtime.eth_rpc_client,
				span__pipeline.Context(),
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

			span__root.Finish()

			return data_map, nil
		},
		p_runtime.runtime_sys)

	//---------------------
	// HEALTH
	http.HandleFunc("/gfethm_worker_inspect/v1/health", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_resp.Write([]byte("ok"))

	})

	//---------------------
}