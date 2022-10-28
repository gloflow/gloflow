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

package gf_images_flows

import (
	// "fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
func InitHandlers(pAuthLoginURLstr string,
	pHTTPmux           *http.ServeMux,
	pTemplatesPathsMap map[string]string,
	pJobsMngrCh        chan gf_images_jobs_core.JobMsg,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES
	templates, gfErr := tmplLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/images/flows/all",
		"/v1/images/flows/add_img",
		"/images/flows/add_img",
		"/images/flows/imgs_exist",
		"/images/flows/browser",
		"/images/flows/browser_page",
	}
	metricsGroupNameStr := "flows"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_images", handlers_endpoints_lst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:                pHTTPmux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: pAuthLoginURLstr,
	}

	//---------------------

	//-------------------------------------------------
	// GET_ALL_FLOWS
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/images/flows/all",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {
				all_flows_names_lst, gf_err := pipelineGetAll(pCtx, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				//------------------
				// OUTPUT
				data_map := map[string]interface{}{
					"all_flows_lst": all_flows_names_lst,
				}
				return data_map, nil

				//------------------
			}
			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // p_store_run_bool
		nil, 
		pRuntimeSys)

	//-------------------------------------------------
	// ADD_IMAGE
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/flows/add_img",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//--------------------------
				// INPUT

				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				i_map, gf_err := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				imageExternURLstr         := i_map["image_extern_url_str"].(string)
				image_origin_page_url_str := i_map["image_origin_page_url_str"].(string) // if image is from a page, the url of the page
				client_type_str           := i_map["client_type_str"].(string)

				flowsNamesLst := []string{}
				for _, s := range i_map["flows_names_lst"].([]interface{}) {
					flowsNamesLst = append(flowsNamesLst, s.(string))
				}

				//--------------------------

				running_job_id_str, thumb_small_relative_url_str, image_id_str, n_gf_err := FlowsAddExternImageWithPolicy(imageExternURLstr,
					image_origin_page_url_str,
					flowsNamesLst,
					client_type_str,
					pJobsMngrCh,
					userIDstr,
					pCtx,
					pRuntimeSys)

				if n_gf_err != nil {
					return nil, n_gf_err
				}

				//------------------
				// OUTPUT
				data_map := map[string]interface{}{
					"images_job_id_str":                running_job_id_str,
					"thumbnail_small_relative_url_str": thumb_small_relative_url_str,
					"image_id_str":                     image_id_str,
				}
				return data_map, nil

				//------------------
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//-------------------------------------------------
	// ADD_IMAGE
	// DEPRECATED!! - switch to using the v1/auth based add_img handler

	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/flows/add_img",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//--------------------------
				// INPUT
				i_map, gf_err := gf_core.HTTPgetInput(p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				image_extern_url_str      := i_map["image_extern_url_str"].(string)
				image_origin_page_url_str := i_map["image_origin_page_url_str"].(string) // if image is from a page, the url of the page
				client_type_str           := i_map["client_type_str"].(string)

				// flow_name_str := "general" //i["flow_name_str"].(string) // DEPRECATED
				flows_names_lst := []string{}
				for _, s := range i_map["flows_names_lst"].([]interface{}) {
					flows_names_lst = append(flows_names_lst, s.(string))
				}

				//--------------------------

				running_job_id_str, thumb_small_relative_url_str, image_id_str, n_gf_err := FlowsAddExternImage(image_extern_url_str,
					image_origin_page_url_str,
					flows_names_lst,
					client_type_str,
					pJobsMngrCh,
					pRuntimeSys)

				if n_gf_err != nil {
					return nil, n_gf_err
				}

				//------------------
				// OUTPUT
				data_map := map[string]interface{}{
					"images_job_id_str":                running_job_id_str,
					"thumbnail_small_relative_url_str": thumb_small_relative_url_str,
					"image_id_str":                     image_id_str,
				}
				return data_map, nil

				//------------------
			}

			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//-------------------------------------------------
	// IMAGE_EXISTS_IN_SYSTEM - check if extern image url's exist in the system,
	//                          if the image url has already been fetched/transformed and gf_image exists for it

	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/flows/imgs_exist",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {
				
				//--------------------------
				// INPUT
				i_map, gf_err := gf_core.HTTPgetInput(p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				images_extern_urls__untyped_lst := i_map["images_extern_urls_lst"].([]interface{})
				images_extern_urls_lst          := []string{}
				for _, u := range images_extern_urls__untyped_lst {
					u_str                 := u.(string)
					images_extern_urls_lst = append(images_extern_urls_lst, u_str)
				}

				flow_name_str   := i_map["flow_name_str"].(string)
				client_type_str := i_map["client_type_str"].(string)

				//--------------------------
					
				existing_images_lst, gf_err := flowsImagesExistCheck(images_extern_urls_lst, flow_name_str, client_type_str, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}
				//------------------
				// OUTPUT
				data_map := map[string]interface{}{
					"existing_images_lst": existing_images_lst,
				}
				
				return data_map, nil

				//------------------
			}

			return nil, nil
		},
		pHTTPmux,
		metrics,
		false, // p_store_run_bool
		nil,
		pRuntimeSys)

	//-------------------------------------------------
	// FLOWS_BROWSER
	//-------------------------------------------------
	// BROWSER
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/flows/browser",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//------------------
				// INPUT
				qs_map := p_req.URL.Query()

				flow_name_str := "general"
				if a_lst, ok := qs_map["fname"]; ok {
					flow_name_str = a_lst[0]
				}

				//------------------
				// RENDER_TEMPLATE
				template_rendered_str, gf_err := flows__render_initial_page(flow_name_str,
					6,  // p_initial_pages_num_int int,
					10, // p_page_size_int int,
					templates.flows_browser__tmpl,
					templates.flows_browser__subtemplates_names_lst,
					p_ctx,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}
				
				//------------------

				p_resp.Write([]byte(template_rendered_str))
			}

			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//-------------------------------------------------
	// GET_BROWSER_PAGE (slice of posts data series)

	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/flows/browser_page",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "GET" {

				pagesLst, gfErr := pipelineGetPage(pReq, pResp, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//--------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"pages_lst": pagesLst,
				}
				return dataMap, nil
				
				//------------------
			}
			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//-------------------------------------------------

	return nil
}