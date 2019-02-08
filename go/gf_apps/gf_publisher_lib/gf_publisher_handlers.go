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

package gf_publisher_lib

import (
	"time"
	"fmt"
	"strconv"
	"strings"
	"net/http"
	"net/url"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func init_handlers(p_gf_images_runtime_info *Gf_images_extern_runtime_info,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_publisher_handlers.init_handlers()")

	//---------------------
	//TEMPLATES
	
	gf_templates, gf_err := tmpl__load(p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//---------------------
	//HIDDEN DASHBOARD

	http.HandleFunc("/posts/dash/18956180__42115/", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/dash ----------")


	})
	//---------------------
	//GET_POST
	http.HandleFunc("/posts/", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/ ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------
			//response_format_str - "j"(for json)|"h"(for html)

			qs_map := p_req.URL.Query()

			//response_format_str - "j"(for json)|"h"(for html)
			response_format_str := gf_rpc_lib.Get_response_format(qs_map, p_runtime_sys)
			//--------------------
			//POST_TITLE

			url_str          := p_req.URL.Path
			url_elements_lst := strings.Split(url_str, "/")

			//IMPORTANT!! - "!=3" - because /a/b splits into {"","a","b",}
			if len(url_elements_lst) != 3 {
				usr_msg_str := fmt.Sprintf("get_post url is not of proper format - %s", url_str)
				gf_err      := gf_core.Error__create(usr_msg_str,
					"verify__invalid_value_error",
					&map[string]interface{}{"url_str":url_str,},
					nil, "gf_publisher_lib", p_runtime_sys)
				gf_rpc_lib.Error__in_handler("/posts/", usr_msg_str, gf_err, p_resp, p_runtime_sys)
				return
			}

			raw_post_title_str := url_elements_lst[2]

			//IMPORTANT!! - replaceAll() - is used here because at the time of testing all titles were still
			//                             with their spaces (" ") encoded as "+". So for the title to be correct,
			//                             for lookups against the internal DB, this is decoded.
			//decodeComponent() - this decodes the percentage encoded symbols. it does not remove
			//                    "+" encoded spaces (" "), and the need for replaceAll()
			post_title_encoded_str := strings.Replace(raw_post_title_str, "+", " ", -1)

			//QueryUnescape() - converting each 3-byte encoded substring of the form "%AB" into the
			//                  hex-decoded byte 0xAB. It returns an error if any % is not followed by two hexadecimal digits.
			post_title_str, err := url.QueryUnescape(post_title_encoded_str)
			if err != nil {

				usr_msg_str := fmt.Sprintf("post title cant be query_unescaped - %s", post_title_encoded_str)
				gf_err      := gf_core.Error__create(usr_msg_str,
					"verify__invalid_query_string_encoding_error",
					&map[string]interface{}{"post_title_encoded_str": post_title_encoded_str,},
					err, "gf_publisher_lib", p_runtime_sys)

				gf_rpc_lib.Error__in_handler("/posts/", usr_msg_str, gf_err, p_resp, p_runtime_sys)
				return
			}
			p_runtime_sys.Log_fun("INFO","post_title_str - "+post_title_str)
			//--------------------

			gf_err := Pipeline__get_post(post_title_str,
				response_format_str,
				gf_templates.post__tmpl,
				gf_templates.post__subtemplates_names_lst,
				p_resp,
				p_runtime_sys)

			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/posts/", "get_post pipeline failed", gf_err, p_resp, p_runtime_sys)
				return
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------
	//POST_CREATE
	http.HandleFunc("/posts/create",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/create ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------
			//INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input("/posts/create", p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}
			post_info_map := i_map
			//------------

			_,images_job_id_str, gf_err := Pipeline__create_post(post_info_map, p_gf_images_runtime_info, p_runtime_sys)

			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/posts/create", "create_post pipeline failed", gf_err, p_resp, p_runtime_sys)
				return 
			}

			gf_rpc_lib.Http_Respond(map[string]interface{}{"images_job_id_str":images_job_id_str}, "OK", p_resp, p_runtime_sys)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/create", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------
	//POST_STATUS
	
	http.HandleFunc("/posts/status", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/status ----------")
	})
	//---------------------
	/*http.HandleFunc("/posts/create_with_updates",func(p_resp http.ResponseWriter,
													p_req *http.Request) {
			p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/create_with_updates ----------")
		})*/

	http.HandleFunc("/posts/update", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/update ----------")
	})

	http.HandleFunc("/posts/delete", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/delete ----------")
		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------
			//INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input("/posts/delete", p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}
			post_title_str := i_map["title_str"].(string)
			//------------

			gf_err = DB__mark_as_deleted_post(post_title_str, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/posts/delete", "delete_post pipeline failed", gf_err, p_resp, p_runtime_sys)
				return 
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/delete", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------
	//POSTS_BROWSER
	http.HandleFunc("/posts/browser", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/browser ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------
			//response_format_str - "json"|"html"

			qs_map := p_req.URL.Query()

			//response_format_str - "j"(for json)|"h"(for html)
			response_format_str := gf_rpc_lib.Get_response_format(qs_map, p_runtime_sys)
			//--------------------

			gf_err := Render_initial_pages(response_format_str,
				6, //p_initial_pages_num_int int
				5, //p_page_size_int
				gf_templates.posts_browser__tmpl,
				gf_templates.posts_browser__subtemplates_names_lst,
				p_resp,
				p_runtime_sys)

			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/posts/browser",
					"failed to render posts_browser initial page",
					gf_err, p_resp, p_runtime_sys)
				return
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/browser", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------
	//GET_BROWSER_PAGE (slice of posts data series)
	http.HandleFunc("/posts/browser_page", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /posts/browser_page ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
			
			//--------------------
			//INPUT

			qs_map := p_req.URL.Query()

			page_index_int := 0 //default - "h" - HTML
			var err error

			if a_lst,ok := qs_map["pg_index"]; ok {
				input_val          := a_lst[0]
				page_index_int, err = strconv.Atoi(input_val) //user supplied value
				if err != nil {
				
					usr_msg_str := "pg_index (page_index) is not an integer"
					gf_err      := gf_core.Error__create(usr_msg_str,
						"verify__value_not_integer_error",
						&map[string]interface{}{"input_val":input_val,},
						err, "gf_publisher_lib", p_runtime_sys)

					gf_rpc_lib.Error__in_handler("/posts/browser_page", usr_msg_str, gf_err, p_resp, p_runtime_sys)
					return
				}
			}

			page_size_int := 10 //default - "h" - HTML
			if a_lst,ok := qs_map["pg_size"]; ok {
				input_val         := a_lst[0]
				page_size_int, err = strconv.Atoi(input_val) //user supplied value
				if err != nil {

					usr_msg_str := "pg_size (page_size) is not an integer"
					gf_err      := gf_core.Error__create(usr_msg_str,
						"verify__value_not_integer_error",
						&map[string]interface{}{"input_val":input_val,},
						err, "gf_publisher_lib", p_runtime_sys)

					gf_rpc_lib.Error__in_handler("/posts/browser_page", usr_msg_str, gf_err, p_resp, p_runtime_sys)
					return
				}
			}
			//--------------------
			
			serialized_pages_lst, gf_err := Get_posts_page(page_index_int, page_size_int, p_runtime_sys)
			if err != nil {
				gf_rpc_lib.Error__in_handler("/posts/browser_page", "failed to get posts page", gf_err, p_resp, p_runtime_sys)
				return
			}

			//------------
			//JSON RESPONSE

			r_lst,_ := json.Marshal(serialized_pages_lst)
			r_str   := string(r_lst)
			fmt.Fprintf(p_resp,r_str)
			//------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/browser_page", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------
	//POSTS_ELEMENTS
	http.HandleFunc("/posts_elements/create",func(p_resp http.ResponseWriter, p_req *http.Request) {


	})
	//---------------------

	return nil
}