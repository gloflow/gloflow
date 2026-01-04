/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_analytics_lib

import (
	"context"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//-------------------------------------------------

func initHandlers(pTemplatesPathsMap map[string]string,
	pAuthSubsystemTypeStr string,
	pAuthLoginURLstr      string,
	pKeyServer            *gf_identity_core.GFkeyServerInfo,
	pRPCglobalMetrics     *gf_rpc_lib.GFglobalMetrics,
	pHTTPmux              *http.ServeMux,
	pRuntimeSys           *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/a/ue",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_analytics", handlersEndpointsLst, pRuntimeSys)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:           pHTTPmux,
		Metrics:       metrics,
		MetricsGlobal: pRPCglobalMetrics,
		
		StoreRunBool: true,
		SentryHub:    nil,

		// AUTH
		AuthSubsystemTypeStr: pAuthSubsystemTypeStr,
		AuthLoginURLstr:      pAuthLoginURLstr,
		AuthKeyServer:        pKeyServer,
	}

	//---------------------
	
	//---------------------
	// USER_EVENT
	/*
	IMPORTANT!! - this is a special case handler, we dont want it to return any standard JSON responses,
		this handler should be fire-and-forget from the users/clients perspective.
	*/
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/a/ue",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			// CORS - preflight request
			gf_rpc_lib.HTTPcorsPreflightHandle(pReq, pResp)

			if pReq.Method == "POST" {

				//-----------------
				// INPUT
				
				input, gfErr := gf_events.UserEventParseInput(pReq, pResp, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				var userIDfinalStr gf_core.GF_ID

				userLoggedInBool, userID, gfErr := gf_identity_core.GetUserIDfromReq(pReq,
					pAuthSubsystemTypeStr,
					pCtx,
					pRuntimeSys)

				if userLoggedInBool {
					userIDfinalStr = userID
				}

				//-----------------
							
				gf_events.UserEventCreate(input,
					userIDfinalStr,
					pCtx,
					pRuntimeSys)

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//--------------
	return nil
}