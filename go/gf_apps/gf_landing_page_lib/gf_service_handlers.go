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

package gf_landing_page_lib

import (
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//------------------------------------------------

func init_handlers(p_templates_paths_map map[string]string,
	pHTTPmux    *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}
	
	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/",
		"/landing/main",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_landing_page", handlers_endpoints_lst)

	//---------------------
	// MAIN

	//------------------------------------------------
	landingMainHandlerFun := func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

		if p_req.Method == "GET" {


			imgs__max_random_cursor_position_int  := 10000
			posts__max_random_cursor_position_int := 2000
			gfErr := pipelineRenderLandingPage(imgs__max_random_cursor_position_int,
				posts__max_random_cursor_position_int,
				5,  // p_featured_posts_to_get_int
				10, // p_featured_imgs_to_get_int
				gf_templates.tmpl,
				gf_templates.subtemplates_names_lst,
				p_resp,
				pRuntimeSys)

			if gfErr != nil {
				return nil, gfErr
			}
		}
		
		// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
		//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
		//               from returning it.
		return nil, nil
	}

	//------------------------------------------------
	// ROOT
	gf_rpc_lib.CreateHandlerHTTPwithMux("/",
		landingMainHandlerFun,
		pHTTPmux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	gf_rpc_lib.CreateHandlerHTTPwithMux("/landing/main",
		landingMainHandlerFun,
		pHTTPmux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------
	return nil
}