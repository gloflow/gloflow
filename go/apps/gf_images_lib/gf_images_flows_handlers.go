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
	"fmt"
	"time"
	"net/http"
	"text/template"
	"gf_core"
	"gf_rpc_lib"
	"apps/gf_images_lib/gf_images_jobs"
)
//-------------------------------------------------
func Flows__init_handlers(p_templates_dir_path_str string,
					p_jobs_mngr_ch chan gf_images_jobs.Job_msg,
					p_runtime_sys       *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_flows_handlers.Flows__init_handlers()")

	//---------------------
	//TEMPLATES
	template_name_str       := "gf_images_flows_browser.html"
	template_path_str       := fmt.Sprintf("%s/gf_images_flows_browser.html",p_templates_dir_path_str)
	flows_browser__tmpl,err := template.New(template_name_str).ParseFiles(template_path_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse a template",
			"template_create_error",
			&map[string]interface{}{"template_path_str":template_path_str,},
			err,"gf_images_lib",p_runtime_sys)
		return gf_err
	}
	//---------------------
	
	//-------------------------------------------------
	http.HandleFunc("/images/flows/add_img",func(p_resp http.ResponseWriter,
												p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/flows/add_img ----------")
		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			//INPUT
			i,gf_err := gf_rpc_lib.Get_http_input("/images/flows/add_img",
											p_resp,
											p_req,
											p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/add_img",
								"failed parse input for adding an image to a flow", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}

			image_extern_url_str      := i["image_extern_url_str"].(string)
			image_origin_page_url_str := i["image_origin_page_url_str"].(string) //if image is from a page, the url of the page
			client_type_str           := i["client_type_str"].(string)

			//flow_name_str := "general" //i["flow_name_str"].(string) //DEPRECATED
			flows_names_lst := []string{}
			for _,s := range i["flows_names_lst"].([]interface{}) {
				flows_names_lst = append(flows_names_lst,s.(string))
			}
			//--------------------------

			running_job_id_str,thumbnail_small_relative_url_str,image_id_str,n_gf_err := Flows__add_extern_image(image_extern_url_str,
																							image_origin_page_url_str,
																							flows_names_lst,
																							client_type_str,
																							p_jobs_mngr_ch,
																							p_runtime_sys)
			if n_gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/add_img",
								"failed to add image to flow", //p_user_msg_str
								n_gf_err,p_resp,p_runtime_sys)
				return
			}
			//------------------
			//OUTPUT
			
			data_map := map[string]interface{}{
				"images_job_id_str":               running_job_id_str,
				"thumbnail_small_relative_url_str":thumbnail_small_relative_url_str,
				"image_id_str":                    image_id_str,
			}
			gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/add_img",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//-------------------------------------------------
	//IMAGE_EXISTS_IN_SYSTEM - check if extern image url's exist in the system,
	//                         if the image url has already been fetched/transformed and gf_image exists for it

	http.HandleFunc("/images/flows/imgs_exist",func(p_resp http.ResponseWriter,
												p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/flows/imgs_exist ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			//INPUT
			i,gf_err := gf_rpc_lib.Get_http_input("/images/flows/imgs_exist",
											p_resp,
											p_req,
											p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/imgs_exist",
								"failed to parse input to check if images exist in a flow", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}

			images_extern_urls__untyped_lst := i["images_extern_urls_lst"].([]interface{})
			images_extern_urls_lst          := []string{}
			for _,u := range images_extern_urls__untyped_lst {
				u_str                 := u.(string)
				images_extern_urls_lst = append(images_extern_urls_lst,u_str)
			}

			flow_name_str   := i["flow_name_str"].(string)
			client_type_str := i["client_type_str"].(string)
			//--------------------------
				
			existing_images_lst,gf_err := flows__images_exist_check(images_extern_urls_lst,
													flow_name_str,
													client_type_str,
													p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/imgs_exist",
								"failed to check if extern image exists in the system", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}
			//------------------
			//OUTPUT
			
			data_map := map[string]interface{}{
				"existing_images_lst":existing_images_lst,
			}
			gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/imgs_exist",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//-------------------------------------------------
	//FLOWS_BROWSER
	//-------------------------------------------------
	http.HandleFunc("/images/flows/browser",func(p_resp http.ResponseWriter,
												p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/flows/browser ----------")
		if p_req.Method == "GET" {

			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			qs_map := p_req.URL.Query()

			flow_name_str := "general" //default
			if a_lst,ok := qs_map["fname"]; ok {
				flow_name_str = a_lst[0]
			}

			//------------------
			//RENDER_TEMPLATE
			gf_err := flows__render_initial_page(flow_name_str,
										3,  //p_initial_pages_num_int int, //6
										10, //p_page_size_int int, //5
										flows_browser__tmpl,
										p_resp,
										p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/browser",
										"failed to render posts_browsers initial page",
										gf_err,p_resp,p_runtime_sys)
				return
			}
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/browser",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//-------------------------------------------------
	//GET_BROWSER_PAGE (slice of posts data series)
	http.HandleFunc("/images/flows/browser_page",func(p_resp http.ResponseWriter,
												p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/flows/browser_page ----------")
		if p_req.Method == "GET" {

			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			pages_lst,gf_err := flows__get_page__pipeline(p_req,p_resp,p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/flows/browser_page",
									"failed to get images flow browser_page",
									gf_err,p_resp,p_runtime_sys)
				return
			}

			//--------------------
			//OUTPUT
			
			data_map := map[string]interface{}{
				"pages_lst":pages_lst,
			}
			gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/flows/browser_page",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//-------------------------------------------------

	return nil
}