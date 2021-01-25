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

package gf_eth_monitor_lib

import (
	"fmt"
	"net/http"
	"context"
	"time"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_core"
)

//-------------------------------------------------
func init_handlers(p_get_hosts_fn func() []string,
	p_metrics *gf_eth_monitor_core.GF_metrics,
	p_runtime *gf_eth_monitor_core.GF_runtime) *gf_core.Gf_error {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_eth_monitor_handlers.init_handlers()")

	//---------------------
	// REMOVE!! - temporary, just to test Sentry spans issue.
	http.HandleFunc("/test", func(p_resp http.ResponseWriter, p_req *http.Request) {
		
		ctx := p_req.Context()

		span__root := sentry.StartSpan(ctx, "1111111111")
		time.Sleep(2 * time.Second)
		span__root.Finish()

		span__root_222 := sentry.StartSpan(span__root.Context(), "2222222222222222")
		time.Sleep(2 * time.Second)
		span__root_222.Finish()
	})

	//---------------------
	// GET_BLOCK

	gf_rpc_lib.Create_handler__http("/gfethm/v1/block",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			span__root := sentry.StartSpan(p_ctx, "http__master__get_block")
			defer span__root.Finish()

			//------------------
			// INPUT

			span__input := sentry.StartSpan(span__root.Context(), "get_input")
			defer span__input.Finish() // in case a panic happens before the main .Finish() for this span

			block_num_int, gf_err := Http__get_arg__block_num(p_resp,
				p_req,
				p_runtime.Runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			span__input.Finish()


			//------------------
			// PIPELINE

			span__pipeline := sentry.StartSpan(span__root.Context(), "get_block_pipeline")
			defer span__pipeline.Finish() // in case a panic happens before the main .Finish() for this span

			block_from_workers_map, miners_map, gf_err := gf_eth_monitor_core.Eth_blocks__get_block_pipeline(block_num_int,
				p_get_hosts_fn,
				span__pipeline.Context(),
				p_metrics,
				p_runtime)
			
			span__pipeline.Finish()

			if gf_err != nil {
				return nil, gf_err
			}
			
			//------------------
			data_map := map[string]interface{}{
				"block_from_workers_map": block_from_workers_map,
				"miners_map":             miners_map,
			}

			span__root.Finish()

			return data_map, nil
		},
		p_runtime.Runtime_sys)

	//---------------------
	// GET_MINER
	gf_rpc_lib.Create_handler__http("/gfethm/v1/miner",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			

			// INPUT
			miner_addr_str, gf_err := Http__get_arg__miner_addr(p_resp, p_req, p_runtime.Runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}
			





			




			fmt.Println(miner_addr_str)


			data_map := map[string]interface{}{}
			return data_map, nil
		},
		p_runtime.Runtime_sys)

	//---------------------
	// GET_PEERS
	gf_rpc_lib.Create_handler__http("/gfethm/v1/peers",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

		
			// METRICS
			if p_metrics != nil {
				p_metrics.Counter__http_req_num__get_peers.Inc()
			}

			// PEERS__GET
			peer_names_groups_lst, gf_err := gf_eth_monitor_core.Eth_peers__db__get_pipeline(p_metrics, p_runtime)
			if gf_err != nil {
				return nil, gf_err
			}
			

			data_map := map[string]interface{}{
				"peer_names_groups_lst": peer_names_groups_lst,
			}
			return data_map, nil
		},
		p_runtime.Runtime_sys)

	//---------------------
	// HEALTH
	http.HandleFunc("/gfethm/v1/health", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_resp.Write([]byte("ok"))
	})

	//---------------------

	fs := http.FileServer(http.Dir("../static"))
  	http.Handle("/", fs)

	return nil
}