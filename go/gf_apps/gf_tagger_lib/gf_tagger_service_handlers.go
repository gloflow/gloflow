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

package gf_tagger_lib

import (
	"time"
	"context"
	"net/http"
	"github.com/mitchellh/mapstructure"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func init_handlers(p_templates_paths_map map[string]string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_service_handlers.init_handlers()")

	validator := gf_core.Validate__init()

	// TEMPLATES
	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------------
	// BOOKMARKS
	//---------------------
	// CREATE
	gf_rpc_lib.Create_handler__http("/v1/tags/bookmark/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			if p_req.Method == "POST" {

				//------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				var input GF_bookmark__input_create
				err := mapstructure.Decode(input_map, &input)
				if err != nil {
					gf_err := gf_core.Error__create("failed to load http input into GF_bookmark__input_create struct",
						"mapstruct__decode",
						map[string]interface{}{},
						err, "gf_tagger_lib", p_runtime_sys)
					return nil, gf_err
				}

				input.User_id_str = "anonymous"

				//------------------

				gf_err = bookmarks__pipeline__create(&input,
					validator,
					p_ctx,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
			}

			return nil, nil
		},
		p_runtime_sys)


	// CREATE
	gf_rpc_lib.Create_handler__http("/v1/tags/bookmark/get_all",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			//------------------
			// INPUT
			input := &GF_bookmark__input_get_all{}
			input.User_id_str = "anonymous"

			//------------------


			output, gf_err := bookmarks__pipeline__get_all(input, p_ctx, p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			data_map := map[string]interface{}{
				"bookmarks_lst": output.Bookmarks_lst,
			}
			return data_map, nil

		},
		p_runtime_sys)

	//---------------------
	// NOTES
	//---------------------
	// CREATE
	gf_rpc_lib.Create_handler__http("/v1/tags/notes/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {
			if p_req.Method == "POST" {
				start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
	
				//------------
				// INPUT
				i_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//------------
	
				gf_err = notes__pipeline__add(i_map, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
	
				end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
	
				go func() {
					gf_rpc_lib.Store_rpc_handler_run("/v1/tags/add_note", start_time__unix_f, end_time__unix_f, p_runtime_sys)
				}()
				
				data_map := map[string]interface{}{}
				return data_map, nil
			}

			return nil, nil
		},
		p_runtime_sys)

	//---------------------
	// GET_NOTES

	gf_rpc_lib.Create_handler__http("/v1/tags/notes/get",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			if p_req.Method == "GET" {
				start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

				notes_lst, gf_err := notes__pipeline__get(p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err 
				}

				end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

				go func() {
					gf_rpc_lib.Store_rpc_handler_run("/v1/tags/get_notes", start_time__unix_f, end_time__unix_f, p_runtime_sys)
				}()

				data_map := map[string]interface{}{"notes_lst": notes_lst,}
				return data_map, nil
			}

			return nil, nil
		},
		p_runtime_sys)

	//---------------------
	// TAGS
	//---------------------
	// ADD_TAGS

	gf_rpc_lib.Create_handler__http("/v1/tags/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {
		

			if p_req.Method == "POST" {
				start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

				//------------
				// INPUT
				i_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//------------

				gf_err = tags__pipeline__add(i_map, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

				go func() {
					gf_rpc_lib.Store_rpc_handler_run("/v1/tags/add_tags", start_time__unix_f, end_time__unix_f, p_runtime_sys)
				}()

				data_map := map[string]interface{}{}
				return data_map, nil
			}

			return nil, nil
		},
		p_runtime_sys)
	
	//---------------------
	// GET_OBJECTS_WITH_TAG

	gf_rpc_lib.Create_handler__http("/v1/tags/objects",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			if p_req.Method == "GET" {
				start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

				objects_with_tag_lst, gf_err := tags__pipeline__get_objects(p_req, p_resp, 
					gf_templates.tag_objects__tmpl,
					gf_templates.tag_objects__subtemplates_names_lst,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

				go func() {
					gf_rpc_lib.Store_rpc_handler_run("/v1/tags/objects", start_time__unix_f, end_time__unix_f, p_runtime_sys)
				}()

				// if the response_format was HTML then objects_with_tag_lst is nil,
				// in which case there is no json to send back
				if objects_with_tag_lst != nil {

					data_map := map[string]interface{}{"objects_with_tag_lst": objects_with_tag_lst,}
					return data_map, nil
				} else {
					return nil, nil
				}
			}

			return nil, nil
		},
		p_runtime_sys)

	//---------------------

	return nil
}