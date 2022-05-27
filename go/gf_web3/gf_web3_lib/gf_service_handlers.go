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

package gf_web3_lib

import (
	"fmt"
	"net/http"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_indexer"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_blocks"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_contract"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_tx"
)

//-------------------------------------------------
func InitHandlers(p_get_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string,
	p_indexer_cmds_ch                     gf_eth_indexer.GF_indexer_ch,
	p_indexer_job_updates_new_consumer_ch gf_eth_indexer.GF_job_update_new_consumer_ch,
	p_metrics                             *gf_eth_core.GF_metrics,
	pRuntime                              *gf_eth_core.GF_runtime) *gf_core.GFerror {

	//---------------------
	// INDEXER

	gf_eth_indexer.Init_handlers(p_indexer_cmds_ch,
		p_indexer_job_updates_new_consumer_ch,
		p_metrics,
		pRuntime)

	//---------------------
	// GET__FAVORITES_TX_ADD

	gf_rpc_lib.Create_handler__http("/gfethm/v1/favorites/tx/add",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			spanRoot := sentry.StartSpan(p_ctx, "http__master__favorites_tx_add", sentry.ContinueFromRequest(p_req))
			defer spanRoot.Finish()

			//------------------
			// INPUT

			txHexStr, gfErr := gf_eth_core.Http__get_arg__tx_id_hex(p_resp, p_req, pRuntime.RuntimeSys)

			if gfErr != nil {
				return nil, gfErr
			}

			//------------------

			spanPipeline := sentry.StartSpan(spanRoot.Context(), "favorites_tx_add")

			gfErr = gf_eth_core.Eth_favorites__tx_add(txHexStr,
				spanPipeline.Context(),
				pRuntime)

			spanPipeline.Finish()

			if gfErr != nil {
				return nil, gfErr
			}

			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{}

			//------------------
			spanRoot.Finish()

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//---------------------
	// GET__TX_TRACE_PLOT
	
	gf_rpc_lib.Create_handler__http("/gfethm/v1/tx/trace/plot",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			spanRoot := sentry.StartSpan(p_ctx, "http__master__get_tx_trace_plot", sentry.ContinueFromRequest(p_req))
			defer spanRoot.Finish()

			//------------------
			// INPUT

			txHexStr, gf_err := gf_eth_core.Http__get_arg__tx_id_hex(p_resp, p_req, pRuntime.RuntimeSys)

			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			span__pipeline := sentry.StartSpan(spanRoot.Context(), "tx_trace_plot")

			plot_svg_str, gf_err := gf_eth_tx.Trace__plot(txHexStr,
				p_get_hosts_fn,
				span__pipeline.Context(),
				pRuntime.Py_plugins,
				p_metrics,
				pRuntime)
			
			span__pipeline.Finish()
			
			if gf_err != nil {
				return nil, gf_err
			}

			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"plot_svg_str": plot_svg_str,
			}

			//------------------
			spanRoot.Finish()

			return dataMap, nil
		},
		pRuntime.RuntimeSys)
		
	//---------------------
	// GET__BLOCK

	gf_rpc_lib.Create_handler__http("/gfethm/v1/block",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			spanRoot := sentry.StartSpan(p_ctx, "http__master__get_block")
			ctx        := spanRoot.Context()
			defer spanRoot.Finish()

			//------------------
			// INPUT

			span__input := sentry.StartSpan(ctx, "get_input")

			block_num_int, gfErr := gf_eth_core.Http__get_arg__block_num(p_resp,
				p_req,
				pRuntime.RuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			span__input.Finish()

			//------------------
			// PIPELINE

			// ABI_DEFS
			abisDefsMap, gfErr := gf_eth_contract.Eth_abi__get_defs(ctx, p_metrics, pRuntime)
			if gfErr != nil {
				return nil, gfErr
			}

			span__pipeline := sentry.StartSpan(ctx, "blocks_get_from_workers")

			block_from_workers_map, miners_map, gfErr := gf_eth_blocks.Get_from_workers__pipeline(block_num_int,
				p_get_hosts_fn,
				abisDefsMap,
				span__pipeline.Context(),
				p_metrics,
				pRuntime)
			
			span__pipeline.Finish()

			if gfErr != nil {
				return nil, gfErr
			}
			
			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"block_from_workers_map": block_from_workers_map,
				"miners_map":             miners_map,
			}

			//------------------
			spanRoot.Finish()

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//---------------------
	// GET__MINER
	gf_rpc_lib.Create_handler__http("/gfethm/v1/miner",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			// INPUT
			miner_addr_str, gfErr := gf_eth_core.Http__get_arg__miner_addr(p_resp, p_req, pRuntime.RuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}
			
			fmt.Println(miner_addr_str)

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{}

			//------------------
			return data_map, nil
		},
		pRuntime.RuntimeSys)

	//---------------------
	// GET__PEERS
	gf_rpc_lib.Create_handler__http("/gfethm/v1/peers",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			// METRICS
			if p_metrics != nil {
				p_metrics.Peers__http_req_num__get_peers__counter.Inc()
			}

			// PEERS__GET
			peer_names_groups_lst, gfErr := gf_eth_core.Eth_peers__db__get_pipeline(p_metrics, pRuntime)
			if gfErr != nil {
				return nil, gfErr
			}
			
			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"peer_names_groups_lst": peer_names_groups_lst,
			}

			//------------------
			return data_map, nil
		},
		pRuntime.RuntimeSys)

	//---------------------
	// GET__HEALTH
	http.HandleFunc("/gfethm/v1/health", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_resp.Write([]byte("ok"))
	})

	//---------------------

	fs := http.FileServer(http.Dir("../static"))
  	http.Handle("/", fs)

	return nil
}