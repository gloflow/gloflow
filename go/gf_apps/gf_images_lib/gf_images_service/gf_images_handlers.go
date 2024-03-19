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
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
// INIT_HANDLERS

func InitHandlers(pAuthSubsystemTypeStr string,
	pAuthLoginURLstr   string,
	pKeyServer         *gf_identity_core.GFkeyServerInfo,
	pHTTPmux           *http.ServeMux,
	pTemplatesPathsMap map[string]string,
	pJobsMngrCh      chan gf_images_jobs_core.JobMsg,
	pImgConfig       *gf_images_core.GFconfig,
	pServiceInfo     *gf_images_core.GFserviceInfo,
	pMediaDomainStr  string,
	pStorage         *gf_images_storage.GFimageStorage,
	pS3info          *gf_aws.GFs3Info,
	pMetrics         *gf_images_core.GFmetrics,
	pRuntimeSys      *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//---------------------
	// TEMPLATES
	templates, gfErr := templateLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		
		//-------------
		/*
		PERMANENT_URLS
		IMPORTANT!! - these are url's from the GF system embedded into third-party pages and text,
			and cannot be versioned, they're permanent. so these handlers dont have /v1 in them (versioning segment). 
		*/
		"/images/d/",
		"/images/v/",

		//-------------

		"/v1/images/classify",
		"/v1/images/share",
		"/v1/images/get",
		"/v1/images/upload_init",
		"/v1/images/upload_complete",
		"/images/c",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_images", handlersEndpointsLst)

	//---------------------
	// rpcHandlerRuntime
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:             pHTTPmux,
		Metrics:         metrics,
		StoreRunBool:    true,
		SentryHub:       nil,
		AuthSubsystemTypeStr: pAuthSubsystemTypeStr,
		AuthLoginURLstr:      pAuthLoginURLstr,
		AuthKeyServer:        pKeyServer,
	}

	//---------------------
	// PERMANENT_URLS
	//---------------------
	// GET_IMAGE_URL
	// ADD!! - /images/d/bulk - return resolved final url's from a list of source /images/d/image_name urls

	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/d/",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			if pReq.Method == "GET" {

				//-----------------
				// INPUT
				pathStr          := pReq.URL.Path
				imagePathNameStr := strings.Replace(pathStr, "/images/d/", "", 1)

				qsMap       := pReq.URL.Query()
				flowNameStr := "general"
				if aLst, ok := qsMap["fname"]; ok {
					flowNameStr = aLst[0]
				}
				
				//-----------------

				// check if flow exists
				if _, ok := pImgConfig.ImagesFlowToS3bucketMap[flowNameStr]; !ok {
					gfErr := gf_core.ErrorCreate("image to resolve in unexisting flow",
						"verify__invalid_value_error",
						map[string]interface{}{
							"flow_name_str":    flowNameStr,
							"handler_path_str": "/images/d/",
						},
						nil, "gf_images_lib", pRuntimeSys)
					return nil, gfErr
				}

				imageURLstr := gf_images_core.ImageGetPublicURL(imagePathNameStr,
					pMediaDomainStr,
					pRuntimeSys)

				// redirect user to image url
				http.Redirect(pResp,
					pReq,
					imageURLstr,
					301)
			}

			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//---------------------
	// VIEW_IMAGE
	// renders a solo image
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/v/",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			if pReq.Method == "GET" {

				//-----------------
				// INPUT

				userID := gf_core.GF_ID("anon")

				pathStr    := pReq.URL.Path
				imageIDstr := strings.Replace(pathStr, "/images/v/", "", 1)
				imageID    := gf_images_core.GFimageID(imageIDstr)

				/*
				qsMap       := pReq.URL.Query()
				flowNameStr := "general"
				if aLst, ok := qsMap["fname"]; ok {
					flowNameStr = aLst[0]
				}
				*/

				//-----------------
				// RENDER_TEMPLATE
				templateRenderedStr, gfErr := renderImageViewPage(imageID,
					templates.imagesViewTmpl,
					templates.imagesViewSubtemplatesNamesLst,
					userID,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//-----------------
				pResp.Write([]byte(templateRenderedStr))
			}

			return nil, nil
		},
		pHTTPmux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)


	//---------------------
	// CLASSIFY_IMAGE
	
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/classify",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//-----------------
				// INPUT

				iMap, gfErr :=  gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				// USER_ID
				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				clientTypeStr := iMap["client_type_str"].(string)
				
				// IMAGES_IDS
				imagesIDsLst := iMap["images_ids_lst"].([]string)

				imagesIDsCastedLst := []gf_images_core.GFimageID{}
				for _, imageIDstr := range imagesIDsLst {
					imageID := gf_images_core.GFimageID(imageIDstr)
					imagesIDsCastedLst = append(imagesIDsCastedLst, imageID)
				}

				

				input := &GFimageClassifyInput{
					ClientTypeStr: clientTypeStr,
					ImagesIDsLst:  imagesIDsCastedLst,

				}

				//-----------------
				
				// SHARE
				gfErr = ImageClassify(input,
					userID,
					pJobsMngrCh,
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					
				}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// SHARE_IMAGE
	
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/share",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {
				//-----------------
				// INPUT

				iMap, gfErr :=  gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				// USER_ID
				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				// IMAGE_ID
				imageIDstr := iMap["image_id"].(string)
				imageID := gf_images_core.GFimageID(imageIDstr)

				emailAddressStr := iMap["email_address"].(string)
				emailSubjectStr := iMap["email_subject"].(string)
				emailBodyStr := iMap["email_body"].(string)

				input := &gf_images_core.GFshareInput{
					ImageID:         imageID,
					EmailAddressStr: emailAddressStr,
					EmailSubjectStr: emailSubjectStr,
					EmailBodyStr:    emailBodyStr,
				}

				//-----------------
				
				// SHARE
				gfErr = gf_images_core.SharePipeline(input,
					userID,
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					
				}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// GET_IMAGE
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/images/get",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {
				//-----------------
				// INPUT
				qsMap := pReq.URL.Query()
				
				var imgIDstr string 
				if aLst, ok := qsMap["img_id"]; ok {
					imgIDstr = aLst[0]
				} else {
					gfErr := gf_core.MongoHandleError("failed to get img_id arg from request query string",
						"verify__input_data_missing_in_req_error",
						map[string]interface{}{},
						nil, "gf_images_lib", pRuntimeSys)
					return nil, gfErr
				}

				//-----------------

				imgID := gf_images_core.GFimageID(imgIDstr)
				imageExport, existsBool, gfErr := ImageGet(imgID, pCtx, pRuntimeSys)
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
	// UPLOAD_INIT - client calls this to get the presigned URL to then upload the image to directly.
	//               this is done mainly to save on bandwidth and avoid one extra hop.
	// AUTH

	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/upload_init",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//------------------
				// CORS - in "simple" requests a CORS PREFLIGHT request is not necessary, 
				//        and a CORS header needs to be set on the response of the GET request itself
				//        (not on the preflight OPTIONS request).
				//        "simple" requests are GET/HEAD/POST with standard form/text_plain content-types.
				pResp.Header().Set("Access-Control-Allow-Origin", "*")

				//------------------
				// INPUT
				qsMap := pReq.URL.Query()
				
				// USER_ID
				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				// IMAGE_FORMAT
				var imageFormatStr string
				if aLst, ok := qsMap["imgf"]; ok {
					imageFormatStr = aLst[0]
				}

				// IMAGE_NAME - name that the user has potentially assigned to the image
				var imageNameStr string
				if aLst, ok := qsMap["imgn"]; ok {
					imageNameStr = aLst[0]
				}

				// FLOWS_NAMES - names of flows to which this image should be added
				var flowsNamesLst []string
				if aLst, ok := qsMap["f"]; ok {
					flowsNamesStr := aLst[0]
					flowsNamesLst = strings.Split(flowsNamesStr, ",")
				}

				// CLIENT_TYPE - type of client thats doing the upload
				var clientTypeStr string
				if aLst, ok := qsMap["ct"]; ok {
					clientTypeStr = aLst[0]
				}

				//------------------
				// UPLOAD__INIT
				uploadInfo, gfErr := UploadInit(imageNameStr,
					imageFormatStr,
					flowsNamesLst,
					clientTypeStr,
					userID,
					pStorage,
					pS3info,
					pImgConfig,
					pServiceInfo,
					pCtx,
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
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// UPLOAD_COMPLETE - client calls this to get the presigned URL to then upload the image to directly.
	//                   this is done mainly to save on bandwidth and avoid one extra hop.
	
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/upload_complete",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------------
				// CORS - in "simple" requests a CORS PREFLIGHT request is not necessary, 
				//        and a CORS header needs to be set on the response of the GET request itself
				//        (not on the preflight OPTIONS request).
				//        "simple" requests are GET/HEAD/POST with standard form/text_plain content-types.
				pResp.Header().Set("Access-Control-Allow-Origin", "*")

				//------------------
				// INPUT

				// USER_ID
				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				qsMap := pReq.URL.Query()

				iMap, gfErr :=  gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				// UPLOAD_GF_IMAGE_ID - gf_image ID that was assigned to this uploaded image. it is used here
				//                      to know which ID to use for the new gf_image thats going to be constructed,
				//                      and to know by which ID to query the DB for Gf_image_upload_info.
				var uploadImageIDstr gf_images_core.GFimageID
				if aLst, ok := qsMap["imgid"]; ok {
					uploadImageIDstr = gf_images_core.GFimageID(aLst[0])
				}
				
				// FLOWS_NAMES - names of flows to which this image should be added
				var flowsNamesLst []string
				if aLst, ok := qsMap["f"]; ok {
					flowsNamesStr := aLst[0]
					flowsNamesLst = strings.Split(flowsNamesStr, ",")
				}

				// image metadata (optional)
				var metaMap map[string]interface{}
				if metaMap, ok := iMap["meta_map"]; ok {
					metaMap = metaMap.(map[string]interface{})
				}

				//------------------
				// COMPLETE
				runningJob, gfErr := UploadComplete(uploadImageIDstr,
					flowsNamesLst,
					metaMap,
					userID,
					pJobsMngrCh,
					pServiceInfo,
					pCtx,
					pRuntimeSys)

				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{}

				if runningJob != nil {
					dataMap["images_job_id_str"] = runningJob.Id_str
				}
				return dataMap, nil
				
				//------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// UPLOAD_METRICS - client reports upload metrics from its own perspective
	
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/images/upload_metrics",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------------
				// CORS - in "simple" requests a CORS PREFLIGHT request is not necessary, 
				//        and a CORS header needs to be set on the response of the GET request itself
				//        (not on the preflight OPTIONS request).
				//        "simple" requests are GET/HEAD/POST with standard form/text_plain content-types.
				pResp.Header().Set("Access-Control-Allow-Origin", "*")

				//------------------
				// INPUT

				// USER_ID
				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				qsMap := pReq.URL.Query()

				iMap, gfErr :=  gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				metricsDataMap := iMap

				// UPLOAD_GF_IMAGE_ID - gf_image ID that was assigned to this uploaded image. it is used here
				//                      to know which ID to use for the new gf_image thats going to be constructed,
				//                      and to know by which ID to query the DB for Gf_image_upload_info.
				var uploadImageIDstr gf_images_core.GF_image_id
				if aLst, ok := qsMap["imgid"]; ok {
					uploadImageIDstr = gf_images_core.GFimageID(aLst[0])
				}
				
				// CLIENT_TYPE - type of client thats doing the upload
				var clientTypeStr string
				if aLst, ok := qsMap["ct"]; ok {
					clientTypeStr = aLst[0]
				}

				//------------------
				gfErr = UploadMetricsCreate(uploadImageIDstr,
					clientTypeStr,
					metricsDataMap,
					userID,
					pMetrics,
					pCtx,
					pRuntimeSys)

				if gfErr != nil {
					return nil, gfErr
				}

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

	//---------------------
	// IMAGE_JOB_RESULT FROM CLIENT_BROWSER (distributed jobs run on client machines)
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/c",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//--------------------------
				// INPUT
				iMap, gfErr :=  gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				browser_jobs_runs_results_lst      := iMap["jr"].([]interface{}) //map[string]interface{})
				cast_browser_jobs_runs_results_lst := []map[string]interface{}{}
				for _, r := range browser_jobs_runs_results_lst {
					cast_browser_jobs_runs_results_lst = append(cast_browser_jobs_runs_results_lst,
						r.(map[string]interface{}))
				}

				//--------------------------
				// STORE BROWSER_IMAGE_CALC_RESULT
				gfErr = ProcessBrowserImageCalcResult(cast_browser_jobs_runs_results_lst, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				//--------------------------
			}
			return nil, nil
		},
		pHTTPmux,
		metrics,
		false, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// HEALTH

	// FIX!! - change to "/v1/images/healthz" but have to also fix infra healthcheck path 
	//         otherwise service is going to get marked as unhealthy
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/v1/healthz",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			return nil, nil
		},
		pHTTPmux,
		nil,   // no metrics for health endpoint
		false, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//---------------------
	return nil
}