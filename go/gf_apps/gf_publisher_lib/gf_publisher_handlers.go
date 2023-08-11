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
	"fmt"
	"strconv"
	"context"
	"strings"
	"net/http"
	"net/url"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//-------------------------------------------------

func initHandlers(p_gf_images_runtime_info *GF_images_extern_runtime_info,
	p_templates_paths_map map[string]string,
	p_mux                 *http.ServeMux,
	pRuntimeSys           *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES
	
	gf_templates, gfErr := tmplLoad(p_templates_paths_map, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// GET_POST
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts/",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//--------------------
				// response_format_str - "j"(for json)|"h"(for html)

				qs_map := pReq.URL.Query()

				// response_format_str - "j"(for json)|"h"(for html)
				response_format_str := gf_rpc_lib.GetResponseFormat(qs_map, pRuntimeSys)

				//--------------------
				// POST_TITLE

				url_str          := pReq.URL.Path
				url_elements_lst := strings.Split(url_str, "/")

				// IMPORTANT!! - "!=3" - because /a/b splits into {"","a","b",}
				if len(url_elements_lst) != 3 {
					usr_msg_str := fmt.Sprintf("get_post url is not of proper format - %s", url_str)
					gfErr       := gf_core.ErrorCreate(usr_msg_str,
						"verify__invalid_value_error",
						map[string]interface{}{"url_str":url_str,},
						nil, "gf_publisher_lib", pRuntimeSys)
					return nil, gfErr
				}

				raw_post_title_str := url_elements_lst[2]

				// IMPORTANT!! - replaceAll() - is used here because at the time of testing all titles were still
				//                              with their spaces (" ") encoded as "+". So for the title to be correct,
				//                              for lookups against the internal DB, this is decoded.
				// decodeComponent() - this decodes the percentage encoded symbols. it does not remove
				//                     "+" encoded spaces (" "), and the need for replaceAll()
				post_title_encoded_str := strings.Replace(raw_post_title_str, "+", " ", -1)

				// QueryUnescape() - converting each 3-byte encoded substring of the form "%AB" into the
				//                   hex-decoded byte 0xAB. It returns an error if any % is not followed by two hexadecimal digits.
				post_title_str, err := url.QueryUnescape(post_title_encoded_str)
				if err != nil {

					usr_msg_str := fmt.Sprintf("post title cant be query_unescaped - %s", post_title_encoded_str)
					gfErr      := gf_core.ErrorCreate(usr_msg_str,
						"verify__invalid_query_string_encoding_error",
						map[string]interface{}{"post_title_encoded_str": post_title_encoded_str,},
						err, "gf_publisher_lib", pRuntimeSys)

					return nil, gfErr
				}
				pRuntimeSys.LogFun("INFO", "post_title_str - "+post_title_str)

				//--------------------

				gfErr := PipelineGetPost(post_title_str,
					response_format_str,
					gf_templates.post__tmpl,
					gf_templates.post__subtemplates_names_lst,
					pesp,
					pRuntimeSys)

				if gfErr != nil {
					return nil, gfErr
				}

				return nil, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// POST_CREATE
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts/create",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------
				// INPUT

				// USER_ID
				userID := gf_core.GF_ID("gf")

				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				postInfoMap := iMap
				
				//------------

				_, imagesJobIDstr, gfErr := PipelineCreatePost(postInfoMap,
					p_gf_images_runtime_info,
					userID,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"images_job_id_str": imagesJobIDstr,
				}
				return outputMap, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//---------------------
	// POST_STATUS
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts/status",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//---------------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts/update",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts/delete",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------
				// INPUT
				i_map, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				post_title_str := i_map["title_str"].(string)

				//------------

				gfErr = gf_publisher_core.DBmarkAsDeletedPost(post_title_str, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				return nil, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// POSTS_BROWSER
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts/browser",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {
				
				//--------------------
				// response_format_str - "json"|"html"

				qs_map := pReq.URL.Query()

				// response_format_str - "j"(for json)|"h"(for html)
				response_format_str := gf_rpc_lib.GetResponseFormat(qs_map, pRuntimeSys)
				
				//--------------------

				gfErr := RenderInitialPages(response_format_str,
					6, // p_initial_pages_num_int int
					5, // p_page_size_int
					gf_templates.posts_browser__tmpl,
					gf_templates.posts_browser__subtemplates_names_lst,
					pesp,
					pRuntimeSys)

				if gfErr != nil {
					return nil, gfErr
				}
				return nil, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// GET_BROWSER_PAGE (slice of posts data series)
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts/browser_page",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//--------------------
				// INPUT

				qs_map := pReq.URL.Query()

				page_index_int := 0 // default - "h" - HTML
				var err error

				if a_lst, ok := qs_map["pg_index"]; ok {
					input_val          := a_lst[0]
					page_index_int, err = strconv.Atoi(input_val) // user supplied value
					if err != nil {
					
						usr_msg_str := "pg_index (page_index) is not an integer"
						gfErr      := gf_core.ErrorCreate(usr_msg_str,
							"verify__value_not_integer_error",
							map[string]interface{}{"input_val": input_val,},
							err, "gf_publisher_lib", pRuntimeSys)
						return nil, gfErr
					}
				}

				page_size_int := 10 //default - "h" - HTML
				if a_lst, ok := qs_map["pg_size"]; ok {
					input_val         := a_lst[0]
					page_size_int, err = strconv.Atoi(input_val) //user supplied value
					if err != nil {

						usr_msg_str := "pg_size (page_size) is not an integer"
						gfErr      := gf_core.ErrorCreate(usr_msg_str,
							"verify__value_not_integer_error",
							map[string]interface{}{"input_val": input_val,},
							err, "gf_publisher_lib", pRuntimeSys)
						return nil, gfErr
					}
				}

				//--------------------
				
				serialized_pages_lst, gfErr := Get_posts_page(page_index_int, page_size_int, pRuntimeSys)
				if err != nil {
					return nil, gfErr
				}

				//------------
				// JSON RESPONSE

				r_lst,_ := json.Marshal(serialized_pages_lst)
				r_str   := string(r_lst)
				fmt.Fprintf(pesp,r_str)

				//------------

				return nil, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// POSTS_ELEMENTS
	gf_rpc_lib.CreateHandlerHTTPwithMux("/posts_elements/create",
		func(pCtx context.Context, pesp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------

	return nil
}