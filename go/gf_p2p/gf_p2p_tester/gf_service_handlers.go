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

package main

import (
	"context"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_p2p"
)

//-------------------------------------------------
func initHandlers(pStatusServerCh gf_p2p.GFp2pStatusServerCh,
	pHTTPmux           *http.ServeMux,
	pTemplatesPathsMap map[string]string,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES
	templates, gfErr := tmplLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/p2p/status",
	}
	metricsGroupNameStr := "status"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_p2p", handlersEndpointsLst, pRuntimeSys)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:          pHTTPmux,
		Metrics:      metrics,
		StoreRunBool: false,
		SentryHub:    nil,
	}

	//-------------------------------------------------
	// P2P_STATUS
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/p2p/status",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {


				p2pStatus := gf_p2p.GetStatusFromServer(pStatusServerCh)

				templateRenderedStr, gfErr := viewRenderStatus(p2pStatus,
					templates.p2pStatus,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				pResp.Write([]byte(templateRenderedStr))
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//-------------------------------------------------

	return nil
}