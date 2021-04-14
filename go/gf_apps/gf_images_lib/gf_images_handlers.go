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
	"time"
	"strings"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
)

//-------------------------------------------------
func init_handlers(p_jobs_mngr_ch chan gf_images_jobs.Job_msg,
	p_img_config  *gf_images_utils.GF_config,
	p_s3_info     *gf_core.Gf_s3_info,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_handlers.init_handlers()")

	//---------------------
	// ADD!! - /images/d/bulk - return resolved final url's from a list of source /images/d/image_name urls
	//---------------------
	// GET_IMAGE_URL
	
	http.HandleFunc("/images/d/", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//-----------------
			// INPUT
			path_str            := p_req.URL.Path
			image_path_name_str := strings.Replace(path_str, "/images/d/", "", 1)

			qs_map        := p_req.URL.Query()
			flow_name_str := "general" // default
			if a_lst,ok := qs_map["fname"]; ok {
				flow_name_str = a_lst[0]
			}
			
			//-----------------

			if _, ok := p_img_config.Images_flow_to_s3_bucket_map[flow_name_str]; !ok {
				gf_rpc_lib.Error__in_handler("/images/d",
					"supplied fname argument is for a non-existed flow - "+flow_name_str, //p_user_msg_str
					nil, p_resp, p_runtime_sys)
				return
			}
			s3_bucket_name_str := p_img_config.Images_flow_to_s3_bucket_map[flow_name_str]

			image_s3_url_str := gf_images_utils.S3__get_image_url(image_path_name_str,
				s3_bucket_name_str,
				p_runtime_sys)

			// redirect user to S3 image url
			http.Redirect(p_resp,
				p_req,
				image_s3_url_str,
				301)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/d", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//---------------------
	// UPLOAD_INIT - client calls this to get the presigned URL to then upload the image to directly.
	//               this is done mainly to save on bandwidth and avoid one extra hop.
	
	http.HandleFunc("/images/v1/upload_init", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------------
			// CORS - in "simple" requests a CORS PREFLIGHT request is not necessary, 
			//        and a CORS header needs to be set on the response of the GET request itself
			//        (not on the preflight OPTIONS request).
			//        "simple" requests are GET/HEAD/POST with standard form/text_plain content-types.
			p_resp.Header().Set("Access-Control-Allow-Origin", "*")

			//------------------
			// INPUT
			qs_map := p_req.URL.Query()

			// IMAGE_FORMAT
			var image_format_str string
			if a_lst, ok := qs_map["imgf"]; ok {
				image_format_str = a_lst[0]
			}

			// IMAGE_NAME - name that the user has potentially assigned to the image
			var image_name_str string
			if a_lst, ok := qs_map["imgn"]; ok {
				image_name_str = a_lst[0]
			}

			// FLOWS_NAMES - names of flows to which this image should be added
			var flows_names_lst []string
			if a_lst, ok := qs_map["f"]; ok {
				flows_names_str := a_lst[0]
				flows_names_lst = strings.Split(flows_names_str, ",")
			}

			// CLIENT_TYPE - type of client thats doing the upload
			var client_type_str string
			if a_lst, ok := qs_map["ct"]; ok {
				client_type_str = a_lst[0]
			}

			//------------------
			// UPLOAD__INIT
			upload_info, gf_err := Upload__init(image_name_str,
				image_format_str,
				flows_names_lst,
				client_type_str,
				p_s3_info,
				p_img_config,
				p_runtime_sys)

			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/v1/upload_init",
					"failed to initialize an image upload",
					nil, p_resp, p_runtime_sys)
				return
			}

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{
				"upload_info_map": upload_info,
			}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)

			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/v1/upload_init", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//---------------------
	// UPLOAD_COMPLETE - client calls this to get the presigned URL to then upload the image to directly.
	//               this is done mainly to save on bandwidth and avoid one extra hop.
	
	http.HandleFunc("/images/v1/upload_complete", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------------
			// CORS - in "simple" requests a CORS PREFLIGHT request is not necessary, 
			//        and a CORS header needs to be set on the response of the GET request itself
			//        (not on the preflight OPTIONS request).
			//        "simple" requests are GET/HEAD/POST with standard form/text_plain content-types.
			p_resp.Header().Set("Access-Control-Allow-Origin", "*")

			//------------------
			// INPUT
			qs_map := p_req.URL.Query()

			// UPLOAD_GF_IMAGE_ID - gf_image ID that was assigned to this uploaded image. it is used here
			//                      to know which ID to use for the new gf_image thats going to be constructed,
			//                      and to know by which ID to query the DB for Gf_image_upload_info.
			var upload_gf_image_id_str gf_images_utils.Gf_image_id
			if a_lst, ok := qs_map["imgid"]; ok {
				upload_gf_image_id_str = gf_images_utils.Gf_image_id(a_lst[0])
			}
			
			//------------------
			// COMPLETE
			running_job, gf_err := Upload__complete(upload_gf_image_id_str,
				p_jobs_mngr_ch,
				p_s3_info,
				p_runtime_sys)

			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/v1/upload_complete",
					"failed to complete an image upload",
					nil, p_resp, p_runtime_sys)
				return
			}

			//------------------
			// OUTPUT
			data_map := map[string]interface{}{}

			if running_job != nil {
				data_map["images_job_id_str"] = running_job.Id_str
			}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)
			
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/v1/upload_complete", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//---------------------
	// IMAGE_JOB_RESULT FROM CLIENT_BROWSER (distributed jobs run on client machines)
	
	http.HandleFunc("/images/c", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			// INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input("/images/c", p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}

			browser_jobs_runs_results_lst      := i_map["jr"].([]interface{}) //map[string]interface{})
			cast_browser_jobs_runs_results_lst := []map[string]interface{}{}
			for _, r := range browser_jobs_runs_results_lst {
				cast_browser_jobs_runs_results_lst = append(cast_browser_jobs_runs_results_lst,
					r.(map[string]interface{}))
			}
			//--------------------------
			// STORE BROWSER_IMAGE_CALC_RESULT
			gf_err = Process__browser_image_calc_result(cast_browser_jobs_runs_results_lst, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/c",
					"failed processing browser_image_calc_result", //p_user_msg_str
					gf_err, p_resp, p_runtime_sys)
			}
			//--------------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/c", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//---------------------

	return nil
}