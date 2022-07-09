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

package gf_images_service

import (
	// "fmt"
	"strings"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
func InitHandlers(pAuthLoginURLstr string,
	pHTTPmux           *http.ServeMux,
	p_jobs_mngr_ch     chan gf_images_jobs_core.JobMsg,
	p_img_config       *gf_images_core.GFconfig,
	p_media_domain_str string,
	pS3info            *gf_core.GFs3Info,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_handlers.init_handlers()")
	
	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/images/d/",
		"/v1/images/get",
		"/v1/images/upload_init",
		"/v1/images/upload_complete",
		"/images/c",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_images", handlers_endpoints_lst)

	//---------------------
	// rpcHandlerRuntime
	rpcHandlerRuntime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                pHTTPmux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: pAuthLoginURLstr,
	}

	//---------------------
	// GET_IMAGE
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/images/get",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "GET" {
				//-----------------
				// INPUT
				qsMap := pReq.URL.Query()
				
				var imgIDstr string 
				if a_lst, ok := qsMap["img_id"]; ok {
					imgIDstr = a_lst[0]
				} else {
					gfErr := gf_core.Mongo__handle_error("failed to get img_id arg from request query string",
						"verify__input_data_missing_in_req_error",
						map[string]interface{}{},
						nil, "gf_images_lib", pRuntimeSys)
					return nil, gfErr
				}

				//-----------------

				imgID := gf_images_core.GF_image_id(imgIDstr)
				imageExport, existsBool, gfErr := ImgGet(imgID, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"image_exists_bool": existsBool,      
					"image_export_map":  imageExport,
				}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------

	//---------------------
	// ADD!! - /images/d/bulk - return resolved final url's from a list of source /images/d/image_name urls
	//---------------------
	// GET_IMAGE_URL
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/d/",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			if pReq.Method == "GET" {

				//-----------------
				// INPUT
				path_str            := pReq.URL.Path
				image_path_name_str := strings.Replace(path_str, "/images/d/", "", 1)

				qsMap       := pReq.URL.Query()
				flowNameStr := "general"
				if aLst, ok := qsMap["fname"]; ok {
					flowNameStr = aLst[0]
				}
				
				//-----------------

				if _, ok := p_img_config.Images_flow_to_s3_bucket_map[flowNameStr]; !ok {
					gfErr := gf_core.Error__create("image to resolve in unexisting flow",
						"verify__invalid_value_error",
						map[string]interface{}{
							"flow_name_str":    flowNameStr,
							"handler_path_str": "/images/d/",
						},
						nil, "gf_images_lib", pRuntimeSys)
					return nil, gfErr
				}

				image_s3_url_str := gf_images_core.Image__get_public_url(image_path_name_str,
					p_media_domain_str,
					pRuntimeSys)

				// redirect user to S3 image url
				http.Redirect(pResp,
					pReq,
					image_s3_url_str,
					301)
			}

			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------
	// UPLOAD_INIT - client calls this to get the presigned URL to then upload the image to directly.
	//               this is done mainly to save on bandwidth and avoid one extra hop.
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/images/upload_init",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//------------------
				// CORS - in "simple" requests a CORS PREFLIGHT request is not necessary, 
				//        and a CORS header needs to be set on the response of the GET request itself
				//        (not on the preflight OPTIONS request).
				//        "simple" requests are GET/HEAD/POST with standard form/text_plain content-types.
				p_resp.Header().Set("Access-Control-Allow-Origin", "*")

				//------------------
				// INPUT
				qsMap := p_req.URL.Query()

				// IMAGE_FORMAT
				var imageFormatStr string
				if a_lst, ok := qsMap["imgf"]; ok {
					imageFormatStr = a_lst[0]
				}

				// IMAGE_NAME - name that the user has potentially assigned to the image
				var imageNameStr string
				if a_lst, ok := qsMap["imgn"]; ok {
					imageNameStr = a_lst[0]
				}

				// FLOWS_NAMES - names of flows to which this image should be added
				var flowsNamesLst []string
				if aLst, ok := qsMap["f"]; ok {
					flowsNamesStr := aLst[0]
					flowsNamesLst = strings.Split(flowsNamesStr, ",")
				}

				// CLIENT_TYPE - type of client thats doing the upload
				var client_type_str string
				if aLst, ok := qsMap["ct"]; ok {
					client_type_str = aLst[0]
				}

				//------------------
				// UPLOAD__INIT
				uploadInfo, gfErr := UploadInit(imageNameStr,
					imageFormatStr,
					flowsNamesLst,
					client_type_str,
					pS3info,
					p_img_config,
					pRuntimeSys)

				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"upload_info_map": uploadInfo,
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

	//---------------------
	// UPLOAD_COMPLETE - client calls this to get the presigned URL to then upload the image to directly.
	//               this is done mainly to save on bandwidth and avoid one extra hop.
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/images/upload_complete",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//------------------
				// CORS - in "simple" requests a CORS PREFLIGHT request is not necessary, 
				//        and a CORS header needs to be set on the response of the GET request itself
				//        (not on the preflight OPTIONS request).
				//        "simple" requests are GET/HEAD/POST with standard form/text_plain content-types.
				p_resp.Header().Set("Access-Control-Allow-Origin", "*")

				//------------------
				// INPUT
				qs_map := p_req.URL.Query()

				i_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				// UPLOAD_GF_IMAGE_ID - gf_image ID that was assigned to this uploaded image. it is used here
				//                      to know which ID to use for the new gf_image thats going to be constructed,
				//                      and to know by which ID to query the DB for Gf_image_upload_info.
				var upload_gf_image_id_str gf_images_core.GF_image_id
				if a_lst, ok := qs_map["imgid"]; ok {
					upload_gf_image_id_str = gf_images_core.GF_image_id(a_lst[0])
				}
				
				// image metadata (optional)
				var meta_map map[string]interface{}
				if meta_map, ok := i_map["meta_map"]; ok {
					meta_map = meta_map
				}

				//------------------
				// COMPLETE
				running_job, gf_err := UploadComplete(upload_gf_image_id_str,
					meta_map,
					p_jobs_mngr_ch,
					pS3info,
					pRuntimeSys)

				if gf_err != nil {
					return nil, gf_err
				}

				//------------------
				// OUTPUT
				data_map := map[string]interface{}{}

				if running_job != nil {
					data_map["images_job_id_str"] = running_job.Id_str
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

	//---------------------
	// IMAGE_JOB_RESULT FROM CLIENT_BROWSER (distributed jobs run on client machines)
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/c",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//--------------------------
				// INPUT
				i_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				browser_jobs_runs_results_lst      := i_map["jr"].([]interface{}) //map[string]interface{})
				cast_browser_jobs_runs_results_lst := []map[string]interface{}{}
				for _, r := range browser_jobs_runs_results_lst {
					cast_browser_jobs_runs_results_lst = append(cast_browser_jobs_runs_results_lst,
						r.(map[string]interface{}))
				}

				//--------------------------
				// STORE BROWSER_IMAGE_CALC_RESULT
				gf_err = Process__browser_image_calc_result(cast_browser_jobs_runs_results_lst, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}
				
				//--------------------------
			}
			return nil, nil
		},
		pHTTPmux,
		metrics,
		false, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------
	// HEALTH

	// FIX!! - change to "/v1/images/healthz" but have to also fix infra healthcheck path 
	//         otherwise service is going to get marked as unhealthy
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/v1/healthz",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			return nil, nil
		},
		pHTTPmux,
		nil,   // no metrics for health endpoint
		false, // p_store_run_bool
		nil,
		pRuntimeSys)
	
	//---------------------
	return nil
}