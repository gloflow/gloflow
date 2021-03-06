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
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
func Flows__init_handlers(p_templates_paths_map map[string]string,
	p_jobs_mngr_ch chan gf_images_jobs_core.Job_msg,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows_handlers.Flows__init_handlers()")

	//---------------------
	// TEMPLATES
	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------------
	
	//-------------------------------------------------
	http.HandleFunc("/images/flows/add_img", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			// INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}

			image_extern_url_str      := i_map["image_extern_url_str"].(string)
			image_origin_page_url_str := i_map["image_origin_page_url_str"].(string) //if image is from a page, the url of the page
			client_type_str           := i_map["client_type_str"].(string)

			// flow_name_str := "general" //i["flow_name_str"].(string) //DEPRECATED
			flows_names_lst := []string{}
			for _, s := range i_map["flows_names_lst"].([]interface{}) {
				flows_names_lst = append(flows_names_lst, s.(string))
			}

			//--------------------------

			running_job_id_str, thumb_small_relative_url_str, image_id_str, n_gf_err := Flows__add_extern_image(image_extern_url_str,
				image_origin_page_url_str,
				flows_names_lst,
				client_type_str,
				p_jobs_mngr_ch,
				p_runtime_sys)

			if n_gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/add_img",
					"failed to add image to flow", // p_user_msg_str
					n_gf_err, p_resp, p_runtime_sys)
				return
			}
			//------------------
			// OUTPUT
			
			data_map := map[string]interface{}{
				"images_job_id_str":                running_job_id_str,
				"thumbnail_small_relative_url_str": thumb_small_relative_url_str,
				"image_id_str":                     image_id_str,
			}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)
			
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/add_img", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//-------------------------------------------------
	// IMAGE_EXISTS_IN_SYSTEM - check if extern image url's exist in the system,
	//                          if the image url has already been fetched/transformed and gf_image exists for it

	http.HandleFunc("/images/flows/imgs_exist", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			// INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
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
				
			existing_images_lst, gf_err := flows__images_exist_check(images_extern_urls_lst, flow_name_str, client_type_str, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/imgs_exist",
					"failed to check if extern image exists in the system", //p_user_msg_str
					gf_err, p_resp, p_runtime_sys)
				return
			}
			//------------------
			// OUTPUT
			
			data_map := map[string]interface{}{
				"existing_images_lst": existing_images_lst,
			}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/imgs_exist", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//-------------------------------------------------
	// FLOWS_BROWSER
	//-------------------------------------------------
	http.HandleFunc("/images/flows/browser", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "GET" {

			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------------
			// INPUT
			qs_map := p_req.URL.Query()

			flow_name_str := "general" //default
			if a_lst, ok := qs_map["fname"]; ok {
				flow_name_str = a_lst[0]
			}

			//------------------
			// RENDER_TEMPLATE
			gf_err := flows__render_initial_page(flow_name_str,
				3,  // p_initial_pages_num_int int,
				10, // p_page_size_int int,
				gf_templates.flows_browser__tmpl,
				gf_templates.flows_browser__subtemplates_names_lst,
				p_resp,
				p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/browser", "failed to render posts_browsers initial page", gf_err, p_resp, p_runtime_sys)
				return
			}
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/browser", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//-------------------------------------------------
	// GET_BROWSER_PAGE (slice of posts data series)
	http.HandleFunc("/images/flows/browser_page", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s %s ----------", p_req.Method, p_req.URL))

		if p_req.Method == "GET" {

			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			pages_lst,gf_err := flows__get_page__pipeline(p_req, p_resp, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/browser_page", "failed to get images flow browser_page", gf_err, p_resp, p_runtime_sys)
				return
			}

			//--------------------
			// OUTPUT
			
			data_map := map[string]interface{}{
				"pages_lst": pages_lst,
			}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)
			
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/browser_page", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//-------------------------------------------------

	return nil
}