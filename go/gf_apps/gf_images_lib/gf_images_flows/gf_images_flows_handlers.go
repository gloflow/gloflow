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
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------

func InitHandlers(pAuthSubsystemTypeStr string,
	pAuthLoginURLstr   string,
	pKeyServer         *gf_identity_core.GFkeyServerInfo,
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
	handlersEndpointsLst := []string{
		"/v1/images/flows/all",
		"/v1/images/flows/add_img",
		"/v1/images/flows/imgs_exist",
		"/images/flows/browser",
		"/images/flows/browser_page",
	}
	metricsGroupNameStr := "flows"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_images", handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:             pHTTPmux,
		Metrics:         metrics,
		StoreRunBool:    true,
		SentryHub:       nil,

		// AUTH
		AuthSubsystemTypeStr: pAuthSubsystemTypeStr,
		AuthLoginURLstr:      pAuthLoginURLstr,
		AuthKeyServer:        pKeyServer,
	}

	//---------------------

	//-------------------------------------------------
	// GET_ALL_FLOWS
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/images/flows/all",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {
				allFlowsNamesLst, gfErr := pipelineGetAll(pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"all_flows_lst": allFlowsNamesLst,
				}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // pStoreRunBool
		nil, 
		pRuntimeSys)

	//-------------------------------------------------
	// ADD_IMAGE
	// AUTH
	
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/flows/add_img",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//--------------------------
				// INPUT

				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				imageExternURLstr     := iMap["image_extern_url_str"].(string)
				imageOriginPageURLstr := iMap["image_origin_page_url_str"].(string) // if image is from a page, the url of the page
				clientTypeStr         := iMap["client_type_str"].(string)

				flowsNamesLst := []string{}
				for _, s := range iMap["flows_names_lst"].([]interface{}) {
					flowsNamesLst = append(flowsNamesLst, s.(string))
				}

				//--------------------------

				runningJobIDstr, thumbnailSmallRelativeURLstr, imageIDstr, gfErr := AddExternImageWithPolicy(imageExternURLstr,
					imageOriginPageURLstr,
					flowsNamesLst,
					clientTypeStr,
					userID,
					pJobsMngrCh,
					pCtx,
					pRuntimeSys)

				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"images_job_id_str":                runningJobIDstr,
					"thumbnail_small_relative_url_str": thumbnailSmallRelativeURLstr,
					"image_id_str":                     imageIDstr,
				}
				return dataMap, nil

				//------------------
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//-------------------------------------------------
	// IMAGE_EXISTS_IN_SYSTEM - check if extern image url's exist in the system,
	//                          if the image url has already been fetched/transformed and gf_image exists for it
	// AUTH

	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/flows/imgs_exist",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {
				
				//--------------------------
				// INPUT

				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				imagesExternURLsUntypedLst := iMap["images_extern_urls_lst"].([]interface{})
				imagesExternURLsLst        := []string{}
				for _, u := range imagesExternURLsUntypedLst {
					u_str                 := u.(string)
					imagesExternURLsLst = append(imagesExternURLsLst, u_str)
				}

				flowNameStr   := iMap["flow_name_str"].(string)
				clientTypeStr := iMap["client_type_str"].(string)

				//--------------------------
					
				existingImagesLst, gfErr := imagesExistCheck(imagesExternURLsLst,
					flowNameStr,
					clientTypeStr,
					userID,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"existing_images_lst": existingImagesLst,
				}
				
				return dataMap, nil

				//------------------
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//-------------------------------------------------
	// FLOWS_BROWSER
	//-------------------------------------------------
	// BROWSER
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/flows/browser",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//------------------
				// INPUT

				userID := gf_core.GF_ID("")

				qsMap := pReq.URL.Query()

				flowNameStr := "general"
				if aLst, ok := qsMap["fname"]; ok {
					flowNameStr = aLst[0]
				}

				//------------------
				// RENDER_TEMPLATE
				templateRenderedStr, gfErr := renderInitialPage(flowNameStr,
					6,  // p_initial_pages_num_int int,
					10, // p_page_size_int int,
					templates.flows_browser__tmpl,
					templates.flows_browser__subtemplates_names_lst,
					userID,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				//------------------

				pResp.Write([]byte(templateRenderedStr))
			}

			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//-------------------------------------------------
	// GET_BROWSER_PAGE (slice of posts data series)

	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/flows/browser_page",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				pagesLst, pagesUserNamesLst, gfErr := pipelineGetPage(pReq, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//--------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"pages_lst":            pagesLst,
					"pages_user_names_lst": pagesUserNamesLst,
				}
				return dataMap, nil
				
				//------------------
			}
			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//-------------------------------------------------

	return nil
}