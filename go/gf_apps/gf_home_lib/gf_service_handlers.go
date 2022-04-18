/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_home_lib

import (
	// "fmt"
	"net/http"
	"context"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
type gfTemplates struct {
	mainTmpl                 *template.Template
	mainSubtemplatesNamesLst []string
}

//------------------------------------------------
func initHandlers(pTemplatesPathsMap map[string]string,
	pAuthLoginURLstr string,
	pHTTPmux         *http.ServeMux,
	pRuntimeSys      *gf_core.Runtime_sys) *gf_core.GF_error {
	
	//---------------------
	// TEMPLATES

	gfTemplates, gfErr := templatesLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/home/view",
		"/v1/home/viz/get",
		"/v1/home/viz/update",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_home", handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                pHTTPmux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: pAuthLoginURLstr,
	}

	//---------------------
	// VIZ_GET
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/home/viz/get",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			if pReq.Method == "GET" {

				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)
				
				homeViz, gfErr := PipelineVizPropsGet(userIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"components_lst": homeViz.ComponentsLst,
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
	// VIZ_UPDATE
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/home/viz/update",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			if pReq.Method == "POST" {

				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				gfErr := PipelineVizPropsUpdate(userIDstr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}



				
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// VIEW
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/home/view",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {


			if pReq.Method == "GET" {

				templateRenderedStr, gfErr := PipelineRenderDashboard(gfTemplates.mainTmpl,
					gfTemplates.mainSubtemplatesNamesLst,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				pResp.Write([]byte(templateRenderedStr))
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

//-------------------------------------------------
func templatesLoad(pTemplatesPathsMap map[string]string,
	pRuntimeSys *gf_core.Runtime_sys) (*gfTemplates, *gf_core.Gf_error) {

	mainTemplateFilepathStr := pTemplatesPathsMap["gf_home_main"]

	// MAIN
	tmpl, subtemplatesNamesLst, gfErr := gf_core.Templates__load(mainTemplateFilepathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	templates := &gfTemplates{
		mainTmpl:                 tmpl,
		mainSubtemplatesNamesLst: subtemplatesNamesLst,
	}
	return templates, nil
}