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

package gf_images_lib

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------

type ClientJobImageOutput struct {
	Image_id_str                      gf_images_core.GFimageID
	Image_source_url_str              string
	Thumbnail_small_relative_url_str  string
	Thumbnail_medium_relative_url_str string
	Thumbnail_large_relative_url_str  string
	// Fetch_ok_bool                     bool
	// Transform_ok_bool                 bool
}

//-------------------------------------------------
// p_input_images_origin_pages_urls_str - urls of pages (html or some other resource) where the image image_url
//                                        was found. this is valid for gf_chrome_ext image sources.
//                                        its not relevant for direct image uploads from clients.

func ClientDispatchProcessExternImages(pInputImagesURLsLst []string,
	pInputImagesOriginPagesURLsLst []string,
	pClientTypeStr                 string,
	pTargetImageServiceHostPortStr string,
	pRuntimeSys                    *gf_core.RuntimeSys) (string, []*ClientJobImageOutput, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_images_http_client.ClientDispatchProcessExternImages()")

	runningJobIDstr, imagesOutputsLst, gfErr := clientStartJob(pInputImagesURLsLst,
		pInputImagesOriginPagesURLsLst,
		pClientTypeStr,
		pTargetImageServiceHostPortStr,
		pRuntimeSys)

	if gfErr != nil {
		return "", nil, gfErr
	}

	return runningJobIDstr, imagesOutputsLst, nil
}

//-------------------------------------------------

func clientStartJob(pInputImagesURLsLst []string,
	pInputImagesOriginPagesURLsLst []string,
	pClientTypeStr                 string,
	pTargetImageServiceHostPortStr string,
	pRuntimeSys                    *gf_core.RuntimeSys) (string, []*ClientJobImageOutput, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_images_http_client.clientStartJob()")

	//------------------
	// HTTP REQUEST

	pRuntimeSys.LogFun("INFO", "target_image_service_host_port - "+pTargetImageServiceHostPortStr)

	urlStr  := fmt.Sprintf("http://%s/images/jobs/start", pTargetImageServiceHostPortStr)
	dataMap := map[string]string{
		"job_type_str":    "process_extern_image",
		"client_type_str": pClientTypeStr,
		"imgs_urls_str":   strings.Join(pInputImagesURLsLst, ","),

		// imgs_origin_pages_urls_str - urls of pages (html or some other resource) where the image image_url
		//                              was found. this is valid for gf_chrome_ext image sources.
		//                              its not relevant for direct image uploads from clients.
		"imgs_origin_pages_urls_str": strings.Join(pInputImagesOriginPagesURLsLst, ","),
	}

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	pRuntimeSys.LogFun("INFO", "")
	pRuntimeSys.LogFun("INFO", cyan("       --- CLIENT__START_JOB")+yellow(" ----------->>>>>>"))
	pRuntimeSys.LogFun("INFO", yellow(pInputImagesURLsLst))
    pRuntimeSys.LogFun("INFO", "")

    fmt.Println("")
	spew.Dump(dataMap)

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

	// FIX!! - REMOVE THIS!! - instead use gf_rpc_client function
	dataLst, _  := json.Marshal(dataMap)
	_, body, errs := gorequest.New().
		Post(urlStr).
		Set("accept", "application/json").
		Send(string(dataLst)).
		End()

	if len(errs) > 0 {
		err   := errs[0] // FIX!! - use all errors in some way, just in case
		gfErr := gf_core.ErrorCreate("gf_images_client start_job HTTP REST API request failed",
			"http_client_req_error",
			map[string]interface{}{
				"url_str": urlStr,
			},
			err, "gf_images_lib", pRuntimeSys)
		return "", nil, gfErr
	}

	fmt.Println("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ------------------- ++++++++++++++++")
	fmt.Println(fmt.Sprintf("images_service %s RESPONSE", urlStr))
	pRuntimeSys.LogFun("INFO", fmt.Sprint(body))
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ------------------- ++++++++++++++++")

	//------------------

	rMap := map[string]interface{}{}
	jErr := json.Unmarshal([]byte(body), &rMap)
	if jErr != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from gf_images_client start_job HTTP REST API - %s", urlStr), 
			"json_decode_error",
			map[string]interface{}{
				"url_str": urlStr,
				"body":    body,
			},
			jErr, "gf_images_lib", pRuntimeSys)
		return "", nil, gfErr
	}

	r_status_str := rMap["status_str"].(string)
	if r_status_str != "OK" {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("received a non-OK response from gf_images_client start_job HTTP REST API - %s", urlStr),
			"http_client_gf_status_error",
			map[string]interface{}{
				"url_str": urlStr,
				"body":    body,
			},
			nil, "gf_images_lib", pRuntimeSys)
		return "", nil, gfErr
	}

	r_data_map := rMap["data"].(map[string]interface{})

	//-----------------
	// RUNNING_JOB_ID

	if _, ok := r_data_map["running_job_id_str"]; !ok {
		err_usr_msg := fmt.Sprintf("%s response didnt return 'running_job_id_str'", urlStr)
		gfErr := gf_core.ErrorCreate(err_usr_msg,
			"verify__missing_key_error",
			map[string]interface{}{"r_map": rMap,},
			nil, "gf_images_lib", pRuntimeSys)
		return "", nil, gfErr
	}

	runningJobIDstr := r_data_map["running_job_id_str"].(string)

	//-----------------
	// JUB_RESULT_IMAGES
	jobExpectedOutputsUntypedLst := r_data_map["job_expected_outputs_lst"].([]interface{})
	imagesOutputsLst             := []*ClientJobImageOutput{}

	for _, o := range jobExpectedOutputsUntypedLst {
		imageOutput := &ClientJobImageOutput{
			Image_id_str:                      gf_images_core.Gf_image_id(o.(map[string]interface{})["image_id_str"].(string)),
			Image_source_url_str:              o.(map[string]interface{})["image_source_url_str"].(string),
			Thumbnail_small_relative_url_str:  o.(map[string]interface{})["thumbnail_small_relative_url_str"].(string),
			Thumbnail_medium_relative_url_str: o.(map[string]interface{})["thumbnail_medium_relative_url_str"].(string),
			Thumbnail_large_relative_url_str:  o.(map[string]interface{})["thumbnail_large_relative_url_str"].(string),
			// Fetch_ok_bool                    :fetch_ok_bool,
			// Transform_ok_bool                :transform_ok_bool,
		}
		imagesOutputsLst = append(imagesOutputsLst, imageOutput)
	}
	
	//-----------------

	return runningJobIDstr, imagesOutputsLst, nil
}

//-------------------------------------------------

func clientGetStatus(pRunningJobIDstr string,
	pTargetImageServiceHostPortStr string,
	pRuntimeSys                    *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_images_http_client.clientGetStatus()")

	url_str := fmt.Sprintf("http://%s/images/jobs/status", pTargetImageServiceHostPortStr)

	_, body, errs := gorequest.New().
		Get(url_str).
		Set("accept", "text/event-stream").
		Query(fmt.Sprintf(`running_job_id_str=%s`, pRunningJobIDstr)).
		End()

	if len(errs) > 0 {
		err := errs[0]
		gfErr := gf_core.ErrorCreate("failed make a client HTTP request to /images/jobs/status",
			"http_client_req_error",
			map[string]interface{}{
				"running_job_id_str":                  pRunningJobIDstr,
				"target__image_service_host_port_str": pTargetImageServiceHostPortStr,
			},
			err, "gf_images_lib", pRuntimeSys)
		return nil, gfErr
	}

	update_items_lst,gfErr := clientParseSSEresponse(body, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return update_items_lst, nil
}

//-------------------------------------------------
// SSE

func clientParseSSEresponse(pBodyStr string,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_images_http_client.clientParseSSEresponse()")

	data_items_lst := []map[string]interface{}{}

	for _, line_str := range strings.Split(pBodyStr, `\n`) {

		pRuntimeSys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>")
		pRuntimeSys.LogFun("INFO", line_str)

		// filter out keep-alive new lines
		if line_str != "" && strings.HasPrefix(line_str, "data: ") {


			msg_str := strings.Replace(line_str, "data: ", "", 1)
			msg_map := map[string]interface{}{}
			err     := json.Unmarshal([]byte(msg_str), &msg_map)

			if err != nil {

				gfErr := gf_core.ErrorCreate("failed to parse JSON response line of the SSE stream (of even updates from a gf_images server)",
					"json_decode_error",
					map[string]interface{}{"line_str": line_str,},
					err, "gf_images_lib", pRuntimeSys)

				return nil, gfErr
			}

			//-------------------
			// STATUS
			if _,ok := msg_map["status_str"]; !ok {
				err_usr_msg := "sse message json doesnt container key status_str"
				gfErr       := gf_core.ErrorCreate(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", pRuntimeSys)
				return nil, gfErr
			}
			status_str := msg_map["status_str"].(string)

			if !(status_str == "ok" || status_str == "error") {

				err_usr_msg := "sse message json status_str key is not of value ok|error"
				gfErr      := gf_core.ErrorCreate(err_usr_msg,
					"verify__invalid_key_value_error",
					map[string]interface{}{
						"status_str": status_str,
						"msg_map":    msg_map,
					},
					nil, "gf_images_lib", pRuntimeSys)
				return nil, gfErr
			}

			//-------------------
			// DATA
			if _,ok := msg_map["data_map"]; !ok {
				err_usr_msg := "sse message json doesnt container key data_map"
				gfErr      := gf_core.ErrorCreate(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", pRuntimeSys)
				return nil, gfErr
			}
			
			dataMap := msg_map["data_map"].(map[string]interface{})
			
			//-------------------
			data_items_lst = append(data_items_lst, dataMap)
		}
	}
	return data_items_lst, nil
}