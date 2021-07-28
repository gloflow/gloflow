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
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_core"
)

//-------------------------------------------------
func init_handlers(p_get_hosts_fn func(context.Context, *gf_eth_monitor_core.GF_runtime) []string,
	p_metrics *gf_eth_monitor_core.GF_metrics,
	p_runtime *gf_eth_monitor_core.GF_runtime) *gf_core.GF_error {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_eth_monitor_handlers.init_handlers()")



	//---------------------
	// GET_BLOCK_PERSIST_BULK

	gf_rpc_lib.Create_handler__http("/gfethm/v1/block/persist_bulk",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			span__root := sentry.StartSpan(p_ctx, "http__master__block_persist_bulk", sentry.ContinueFromRequest(p_req))
			ctx := span__root.Context()
			defer span__root.Finish()

			//------------------
			// INPUT

			block_start_uint, block_end_uint, gf_err := Http__get_arg__block_range(p_resp, p_req, p_runtime.Runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			
			
			// ABI_DEFS
			abis_defs_map, gf_err := gf_eth_monitor_core.Eth_abi__get_defs(ctx, p_metrics, p_runtime)
			if gf_err != nil {
				return nil, gf_err
			}

			span__pipeline := sentry.StartSpan(ctx, "blocks_persist_bulk")

			gf_errs_lst := gf_eth_monitor_core.Eth_blocks__get_and_persist_bulk__pipeline(block_start_uint,
				block_end_uint,
				p_get_hosts_fn,
				abis_defs_map,
				span__pipeline.Context(),
				p_metrics,
				p_runtime)
			
			span__pipeline.Finish()

			if len(gf_errs_lst) > 0 {

				// FIX!! - see if all errors should be returned maybe.
				gf_err__first := gf_errs_lst[0]
				return nil, gf_err__first
			}

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{}

			//------------------
			span__root.Finish()

			return data_map, nil
		},
		p_runtime.Runtime_sys)

	//---------------------
	// GET_FAVORITES_TX_ADD

	gf_rpc_lib.Create_handler__http("/gfethm/v1/favorites/tx/add",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			span__root := sentry.StartSpan(p_ctx, "http__master__favorites_tx_add", sentry.ContinueFromRequest(p_req))
			defer span__root.Finish()

			//------------------
			// INPUT

			tx_hex_str, gf_err := Http__get_arg__tx_id_hex(p_resp, p_req, p_runtime.Runtime_sys)

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------

			span__pipeline := sentry.StartSpan(span__root.Context(), "favorites_tx_add")

			gf_err = gf_eth_monitor_core.Eth_favorites__tx_add(tx_hex_str,
				span__pipeline.Context(),
				p_runtime)

			span__pipeline.Finish()

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{}

			//------------------
			span__root.Finish()

			return data_map, nil
		},
		p_runtime.Runtime_sys)

	//---------------------
	// GET_TX_TRACE_PLOT
	
	gf_rpc_lib.Create_handler__http("/gfethm/v1/tx/trace/plot",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			span__root := sentry.StartSpan(p_ctx, "http__master__get_tx_trace_plot", sentry.ContinueFromRequest(p_req))
			defer span__root.Finish()

			//------------------
			// INPUT

			tx_hex_str, gf_err := Http__get_arg__tx_id_hex(p_resp, p_req, p_runtime.Runtime_sys)

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			span__pipeline := sentry.StartSpan(span__root.Context(), "tx_trace_plot")

			plot_svg_str, gf_err := gf_eth_monitor_core.Eth_tx_trace__plot(tx_hex_str,
				p_get_hosts_fn,
				span__pipeline.Context(),
				p_runtime.Py_plugins,
				p_metrics,
				p_runtime)
			
			span__pipeline.Finish()
			
			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"plot_svg_str": plot_svg_str,
			}

			//------------------
			span__root.Finish()

			return data_map, nil
		},
		p_runtime.Runtime_sys)
		
	//---------------------
	// GET_BLOCK

	gf_rpc_lib.Create_handler__http("/gfethm/v1/block",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			span__root := sentry.StartSpan(p_ctx, "http__master__get_block")
			ctx        := span__root.Context()
			defer span__root.Finish()

			//------------------
			// INPUT

			span__input := sentry.StartSpan(ctx, "get_input")

			block_num_int, gf_err := Http__get_arg__block_num(p_resp,
				p_req,
				p_runtime.Runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			span__input.Finish()

			//------------------
			// PIPELINE

			// ABI_DEFS
			abis_defs_map, gf_err := gf_eth_monitor_core.Eth_abi__get_defs(ctx, p_metrics, p_runtime)
			if gf_err != nil {
				return nil, gf_err
			}

			span__pipeline := sentry.StartSpan(ctx, "blocks_get_from_workers")

			block_from_workers_map, miners_map, gf_err := gf_eth_monitor_core.Eth_blocks__get_from_workers__pipeline(block_num_int,
				p_get_hosts_fn,
				abis_defs_map,
				span__pipeline.Context(),
				p_metrics,
				p_runtime)
			
			span__pipeline.Finish()

			if gf_err != nil {
				return nil, gf_err
			}
			
			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"block_from_workers_map": block_from_workers_map,
				"miners_map":             miners_map,
			}

			//------------------
			span__root.Finish()

			return data_map, nil
		},
		p_runtime.Runtime_sys)

	//---------------------
	// GET_MINER
	gf_rpc_lib.Create_handler__http("/gfethm/v1/miner",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			// INPUT
			miner_addr_str, gf_err := Http__get_arg__miner_addr(p_resp, p_req, p_runtime.Runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}
			
			fmt.Println(miner_addr_str)

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{}

			//------------------
			return data_map, nil
		},
		p_runtime.Runtime_sys)

	//---------------------
	// GET_PEERS
	gf_rpc_lib.Create_handler__http("/gfethm/v1/peers",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			// METRICS
			if p_metrics != nil {
				p_metrics.Peers__http_req_num__get_peers__counter.Inc()
			}

			// PEERS__GET
			peer_names_groups_lst, gf_err := gf_eth_monitor_core.Eth_peers__db__get_pipeline(p_metrics, p_runtime)
			if gf_err != nil {
				return nil, gf_err
			}
			
			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"peer_names_groups_lst": peer_names_groups_lst,
			}

			//------------------
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