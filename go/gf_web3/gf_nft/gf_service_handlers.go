/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_nft

import (
	// "fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func InitHandlers(pHTTPmux *http.ServeMux,
	pConfig           *gf_eth_core.GF_config,
	pImagesJobsMngrCh chan gf_images_jobs_core.JobMsg,
	pRuntimeSys       *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// METRICS

	metrics := MetricsCreate()

	handlersEndpointsLst := []string{
		"/v1/web3/nft/index_address",
	}
	metricsGroupNameStr := "main"
	metricsForHandlers := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_web3_nft", handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:                pHTTPmux,
		Metrics:            metricsForHandlers,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		// Auth_login_url_str: pAuthLoginURLstr,
	}

	//---------------------
	// INDEX_ADDRESS
	// this is potentially a long-running process
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/web3/nft/index_address",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				input, gfErr := httpInputForIndexAddress(userIDstr, pReq, pResp,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------
				// START INDEXING
				pipelineIndexAddress(input,
					pConfig,
					pImagesJobsMngrCh,
					metrics,
					pCtx,
					pRuntimeSys)
				
				//---------------------
				
				outputMap := map[string]interface{}{}
				return outputMap, nil
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// GET_BY_OWNER - all NFTs for a particular owner address
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/web3/nft/get_by_owner",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				input, gfErr := httpInputForGetByOwner(userIDstr, pReq, pResp,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				nftsExternLst, gfErr := pipelineGetByOwner(input,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"nfts_lst": nftsExternLst,
				}
				return outputMap, nil
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// GET - individual NFT information
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/web3/nft/get",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				input, gfErr := httpInputForGet(userIDstr, pReq, pResp,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				gfErr = pipelineGet(input,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{}
				return outputMap, nil
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------

	return nil
}