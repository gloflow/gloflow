/*
GloFlow media management/publishing system
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
	"time"
	"strings"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
)
//-------------------------------------------------
func init_handlers(p_jobs_mngr_ch chan gf_images_jobs.Job_msg, p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_handlers.init_handlers()")

	img_config := Config__get()
	//---------------------
	//ADD!! - /images/d/bulk - return resolved final url's from a list of source /images/d/image_name urls
	//---------------------
	//GET_IMAGE_URL
	
	http.HandleFunc("/images/d/", func(p_resp http.ResponseWriter, p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/d ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//-----------------
			//INPUT
			path_str            := p_req.URL.Path
			image_path_name_str := strings.Replace(path_str,"/images/d/","",1)

			qs_map        := p_req.URL.Query()
			flow_name_str := "general" //default
			if a_lst,ok := qs_map["fname"]; ok {
				flow_name_str = a_lst[0]
			}
			//-----------------

			if _,ok := img_config.Flow_to_s3bucket_map[flow_name_str]; !ok {
				gf_rpc_lib.Error__in_handler("/images/d",
					"supplied fname argument is for a non-existed flow - "+flow_name_str, //p_user_msg_str
					nil, p_resp, p_runtime_sys)
				return
			}
			s3_bucket_name_str := img_config.Flow_to_s3bucket_map[flow_name_str]

			image_s3_url_str := gf_images_utils.S3__get_image_url(image_path_name_str, s3_bucket_name_str, p_runtime_sys)

			//redirect user to S3 image url
			http.Redirect(p_resp, p_req, image_s3_url_str, 301)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/d", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------
	//IMAGE_JOB_RESULT FROM CLIENT_BROWSER (distributed jobs run on client machines)
	
	http.HandleFunc("/images/c", func(p_resp http.ResponseWriter, p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/c ----------")
		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			//INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input("/images/c", p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}

			browser_jobs_runs_results_lst      := i_map["jr"].([]interface{}) //map[string]interface{})
			cast_browser_jobs_runs_results_lst := []map[string]interface{}{}
			for _,r := range browser_jobs_runs_results_lst {
				cast_browser_jobs_runs_results_lst = append(cast_browser_jobs_runs_results_lst, r.(map[string]interface{}))
			}
			//--------------------------
			//STORE BROWSER_IMAGE_CALC_RESULT
			gf_err = Process__browser_image_calc_result(cast_browser_jobs_runs_results_lst, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/c",
					"failed processing browser_image_calc_result", //p_user_msg_str
					gf_err,p_resp,p_runtime_sys)
			}
			//--------------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/c", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------
}