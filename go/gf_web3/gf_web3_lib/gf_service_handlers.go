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
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_indexer"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_blocks"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_contract"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_tx"
)

//-------------------------------------------------
func InitHandlers(pGetHostsFun func(context.Context, *gf_eth_core.GF_runtime) []string,
	pIndexerCmdsCh                  gf_eth_indexer.GF_indexer_ch,
	pIndexerJobUpdatesNewConsumerCh gf_eth_indexer.GF_job_update_new_consumer_ch,
	pMetrics                        *gf_eth_core.GF_metrics,
	pRuntime                        *gf_eth_core.GF_runtime) *gf_core.GFerror {

	//---------------------
	// INDEXER

	gf_eth_indexer.Init_handlers(pIndexerCmdsCh,
		pIndexerJobUpdatesNewConsumerCh,
		pMetrics,
		pRuntime)

	//---------------------
	// GET__FAVORITES_TX_ADD

	gf_rpc_lib.CreateHandlerHTTP("/gfethm/v1/favorites/tx/add",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			spanRoot := sentry.StartSpan(pCtx, "http__master__favorites_tx_add", sentry.ContinueFromRequest(pReq))
			defer spanRoot.Finish()

			//------------------
			// INPUT

			txHexStr, gfErr := gf_eth_core.Http__get_arg__tx_id_hex(pResp, pReq, pRuntime.RuntimeSys)

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
	
	gf_rpc_lib.CreateHandlerHTTP("/gfethm/v1/tx/trace/plot",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			spanRoot := sentry.StartSpan(pCtx, "http__master__get_tx_trace_plot", sentry.ContinueFromRequest(pReq))
			defer spanRoot.Finish()

			//------------------
			// INPUT

			txHexStr, gfErr := gf_eth_core.Http__get_arg__tx_id_hex(pResp, pReq, pRuntime.RuntimeSys)

			if gfErr != nil {
				return nil, gfErr
			}

			//------------------
			spanPipeline := sentry.StartSpan(spanRoot.Context(), "tx_trace_plot")

			plot_svg_str, gfErr := gf_eth_tx.Trace__plot(txHexStr,
				pGetHostsFun,
				spanPipeline.Context(),
				pRuntime.Py_plugins,
				pMetrics,
				pRuntime)
			
			spanPipeline.Finish()
			
			if gfErr != nil {
				return nil, gfErr
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

	gf_rpc_lib.CreateHandlerHTTP("/gfethm/v1/block",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			spanRoot := sentry.StartSpan(pCtx, "http__master__get_block")
			ctx        := spanRoot.Context()
			defer spanRoot.Finish()

			//------------------
			// INPUT

			spanInput := sentry.StartSpan(ctx, "get_input")

			blockNumInt, gfErr := gf_eth_core.Http__get_arg__block_num(pResp,
				pReq,
				pRuntime.RuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			spanInput.Finish()

			//------------------
			// PIPELINE

			// ABI_DEFS
			abisDefsMap, gfErr := gf_eth_contract.Eth_abi__get_defs(ctx, pMetrics, pRuntime)
			if gfErr != nil {
				return nil, gfErr
			}

			spanPipeline := sentry.StartSpan(ctx, "blocks_get_from_workers")

			blockFromWorkersMap, gfErr := gf_eth_blocks.Get_from_workers__pipeline(blockNumInt,
				pGetHostsFun,
				abisDefsMap,
				spanPipeline.Context(),
				pMetrics,
				pRuntime)
			
			spanPipeline.Finish()

			if gfErr != nil {
				return nil, gfErr
			}
			
			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"block_from_workers_map": blockFromWorkersMap,
			}

			//------------------
			spanRoot.Finish()

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//---------------------
	// GET__MINER
	gf_rpc_lib.CreateHandlerHTTP("/gfethm/v1/miner",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			// INPUT
			miner_addr_str, gfErr := gf_eth_core.Http__get_arg__miner_addr(pResp, pReq, pRuntime.RuntimeSys)
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
	gf_rpc_lib.CreateHandlerHTTP("/gfethm/v1/peers",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			// METRICS
			if pMetrics != nil {
				pMetrics.Peers__http_req_num__get_peers__counter.Inc()
			}

			// PEERS__GET
			peer_names_groups_lst, gfErr := gf_eth_core.Eth_peers__db__get_pipeline(pMetrics, pRuntime)
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
	http.HandleFunc("/gfethm/v1/health", func(pResp http.ResponseWriter, pReq *http.Request) {
		pResp.Write([]byte("ok"))
	})

	//---------------------

	fs := http.FileServer(http.Dir("../static"))
  	http.Handle("/", fs)

	return nil
}