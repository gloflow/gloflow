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
type Client_job_image_output struct {
	Image_id_str                      gf_images_core.Gf_image_id
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

func Client__dispatch_process_extern_images(p_input_images_urls_lst []string,
	p_input_images_origin_pages_urls_lst  []string,
	p_client_type_str                     string,
	p_target__image_service_host_port_str string,
	p_runtime_sys                         *gf_core.Runtime_sys) (string, []*Client_job_image_output, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_http_client.Client__dispatch_process_extern_images()")

	running_job_id_str, images_outputs_lst, gf_err := client__start_job(p_input_images_urls_lst,
		p_input_images_origin_pages_urls_lst,
		p_client_type_str,
		p_target__image_service_host_port_str,
		p_runtime_sys)

	if gf_err != nil {
		return "", nil, gf_err
	}

	return running_job_id_str, images_outputs_lst, nil
}

//-------------------------------------------------
func client__start_job(p_input_images_urls_lst []string,
	p_input_images_origin_pages_urls_lst  []string,
	p_client_type_str                     string,
	p_target__image_service_host_port_str string,
	p_runtime_sys                         *gf_core.Runtime_sys) (string, []*Client_job_image_output, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_http_client.client__start_job()")

	//------------------
	// HTTP REQUEST

	p_runtime_sys.Log_fun("INFO","p_target__image_service_host_port_str - "+p_target__image_service_host_port_str)

	url_str  := fmt.Sprintf("http://%s/images/jobs/start", p_target__image_service_host_port_str)
	data_map := map[string]string{
		"job_type_str":    "process_extern_image",
		"client_type_str": p_client_type_str,
		"imgs_urls_str":   strings.Join(p_input_images_urls_lst, ","),

		// imgs_origin_pages_urls_str - urls of pages (html or some other resource) where the image image_url
		//                              was found. this is valid for gf_chrome_ext image sources.
		//                              its not relevant for direct image uploads from clients.
		"imgs_origin_pages_urls_str": strings.Join(p_input_images_origin_pages_urls_lst, ","),
	}

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	p_runtime_sys.Log_fun("INFO", "")
	p_runtime_sys.Log_fun("INFO", cyan("       --- CLIENT__START_JOB")+yellow(" ----------->>>>>>"))
	p_runtime_sys.Log_fun("INFO", yellow(p_input_images_urls_lst))
    p_runtime_sys.Log_fun("INFO", "")

    fmt.Println("")
	spew.Dump(data_map)

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

	// FIX!! - REMOVE THIS!! - instead use gf_rpc_client function
	data_lst, _  := json.Marshal(data_map)
	_, body,errs := gorequest.New().
		Post(url_str).
		Set("accept", "application/json").
		Send(string(data_lst)).
		End()

	if errs != nil {
		err    := errs[0] //FIX!! - use all errors in some way, just in case
		gf_err := gf_core.Error__create("gf_images_client start_job HTTP REST API request failed - "+url_str,
			"http_client_req_error",
			map[string]interface{}{"url_str":url_str,},
			err, "gf_images_lib", p_runtime_sys)
		return "", nil, gf_err
	}

	fmt.Println("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ------------------- ++++++++++++++++")
	fmt.Println(fmt.Sprintf("images_service %s RESPONSE", url_str))
	p_runtime_sys.Log_fun("INFO", fmt.Sprint(body))
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ------------------- ++++++++++++++++")

	//------------------

	r_map := map[string]interface{}{}
	j_err := json.Unmarshal([]byte(body), &r_map)
	if j_err != nil {
		gf_err := gf_core.Error__create(fmt.Sprintf("failed to parse json response from gf_images_client start_job HTTP REST API - %s", url_str), 
			"json_unmarshal_error",
			map[string]interface{}{
				"url_str": url_str,
				"body":    body,
			},
			j_err, "gf_images_lib", p_runtime_sys)
		return "", nil, gf_err
	}

	r_status_str := r_map["status_str"].(string)
	if r_status_str != "OK" {
		gf_err := gf_core.Error__create(fmt.Sprintf("received a non-OK response from gf_images_client start_job HTTP REST API - %s", url_str),
			"http_client_gf_status_error",
			map[string]interface{}{
				"url_str": url_str,
				"body":    body,
			},
			nil, "gf_images_lib", p_runtime_sys)
		return "", nil, gf_err
	}

	r_data_map := r_map["data"].(map[string]interface{})

	//-----------------
	// RUNNING_JOB_ID

	if _, ok := r_data_map["running_job_id_str"]; !ok {
		err_usr_msg := fmt.Sprintf("%s response didnt return 'running_job_id_str'", url_str)
		gf_err := gf_core.Error__create(err_usr_msg,
			"verify__missing_key_error",
			map[string]interface{}{"r_map": r_map,},
			nil, "gf_images_lib", p_runtime_sys)
		return "", nil, gf_err
	}

	running_job_id_str := r_data_map["running_job_id_str"].(string)

	//-----------------
	// JUB_RESULT_IMAGES
	job_expected_outputs_untyped_lst := r_data_map["job_expected_outputs_lst"].([]interface{})
	images_outputs_lst               := []*Client_job_image_output{}

	for _, o := range job_expected_outputs_untyped_lst {
		image_output := &Client_job_image_output{
			Image_id_str:                      gf_images_core.Gf_image_id(o.(map[string]interface{})["image_id_str"].(string)),
			Image_source_url_str:              o.(map[string]interface{})["image_source_url_str"].(string),
			Thumbnail_small_relative_url_str:  o.(map[string]interface{})["thumbnail_small_relative_url_str"].(string),
			Thumbnail_medium_relative_url_str: o.(map[string]interface{})["thumbnail_medium_relative_url_str"].(string),
			Thumbnail_large_relative_url_str:  o.(map[string]interface{})["thumbnail_large_relative_url_str"].(string),
			// Fetch_ok_bool                    :fetch_ok_bool,
			// Transform_ok_bool                :transform_ok_bool,
		}
		images_outputs_lst = append(images_outputs_lst, image_output)
	}
	
	//-----------------

	return running_job_id_str, images_outputs_lst, nil
}

//-------------------------------------------------
func client__get_status(p_running_job_id_str string,
	p_target__image_service_host_port_str string,
	p_runtime_sys                         *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_http_client.client__get_status()")

	url_str := fmt.Sprintf("http://%s/images/jobs/status", p_target__image_service_host_port_str)

	_, body, errs := gorequest.New().
		Get(url_str).
		Set("accept", "text/event-stream").
		Query(fmt.Sprintf(`running_job_id_str=%s`, p_running_job_id_str)).
		End()

	if errs != nil {
		err := errs[0]
		gf_err := gf_core.Error__create("failed make a client HTTP request to /images/jobs/status",
			"http_client_req_error",
			map[string]interface{}{
				"running_job_id_str":                  p_running_job_id_str,
				"target__image_service_host_port_str": p_target__image_service_host_port_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}

	update_items_lst,gf_err := client__parse_sse_response(body, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return update_items_lst, nil
}

//-------------------------------------------------
func client__parse_sse_response(p_body_str string, p_runtime_sys *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_http_client.client__parse_sse_response()")

	data_items_lst := []map[string]interface{}{}

	for _, line_str := range strings.Split(p_body_str, `\n`) {

		p_runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>")
		p_runtime_sys.Log_fun("INFO", line_str)

		//filter out keep-alive new lines
		if line_str != "" && strings.HasPrefix(line_str, "data: ") {


			msg_str := strings.Replace(line_str, "data: ", "", 1)
			msg_map := map[string]interface{}{}
			err     := json.Unmarshal([]byte(msg_str), &msg_map)

			if err != nil {

				gf_err := gf_core.Error__create("failed to parse JSON response line of the SSE stream (of even updates from a gf_images server)",
					"json_unmarshal_error",
					map[string]interface{}{"line_str": line_str,},
					err, "gf_images_lib", p_runtime_sys)

				return nil, gf_err
			}

			//-------------------
			//STATUS
			if _,ok := msg_map["status_str"]; !ok {
				err_usr_msg := "sse message json doesnt container key status_str"
				gf_err      := gf_core.Error__create(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", p_runtime_sys)
				return nil, gf_err
			}
			status_str := msg_map["status_str"].(string)

			if !(status_str == "ok" || status_str == "error") {

				err_usr_msg := "sse message json status_str key is not of value ok|error"
				gf_err      := gf_core.Error__create(err_usr_msg,
					"verify__invalid_key_value_error",
					map[string]interface{}{
						"status_str": status_str,
						"msg_map":    msg_map,
					},
					nil, "gf_images_lib", p_runtime_sys)
				return nil, gf_err
			}
			//-------------------
			//DATA
			if _,ok := msg_map["data_map"]; !ok {
				err_usr_msg := "sse message json doesnt container key data_map"
				gf_err      := gf_core.Error__create(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", p_runtime_sys)
				return nil, gf_err
			}
			
			data_map := msg_map["data_map"].(map[string]interface{})
			//-------------------
			data_items_lst = append(data_items_lst, data_map)
		}
	}
	return data_items_lst, nil
}