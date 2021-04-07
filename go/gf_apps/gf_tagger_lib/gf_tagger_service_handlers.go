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
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func init_handlers(p_templates_paths_map map[string]string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_service_handlers.init_handlers()")


	//---------------------
	// TEMPLATES
	
	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	/*main_template_filename_str := "gf_tag_objects.html"
	templates_dir_path_str     := "./web/gf_apps/gf_tagger/gf_tag_objects/templates"

	tag_objects__tmpl, subtemplates_names_lst, gf_err := gf_core.Templates__load(main_template_filename_str, templates_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	/*template_path_str      := "./templates/gf_tag_objects.html"
	tag_objects__tmpl, err := template.New("gf_tag_objects.html").ParseFiles(template_path_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse a template",
			"template_create_error",
			&map[string]interface{}{"template_path_str":template_path_str,},
			err, "gf_images_lib", p_runtime_sys)
		return gf_err
	}*/

	//---------------------
	// NOTES
	//---------------------
	// ADD_NOTE

	http.HandleFunc("/tags/add_note", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST - /tags/add_note ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------
			// INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input("/tags/add_note", p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}
			//------------

			gf_err = pipeline__add_note(i_map, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/tags/add_note", "add_note pipeline failed", gf_err, p_resp, p_runtime_sys)
				return
			}

			data_map := map[string]interface{}{}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/add_note", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//---------------------
	// GET_NOTES

	http.HandleFunc("/tags/get_notes", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST - /tags/get_notes ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			notes_lst, gf_err := pipeline__get_notes(p_req, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/tags/get_notes", "get_notes pipeline failed", gf_err, p_resp, p_runtime_sys)
				return 
			}

			data_map := map[string][]*Gf_note{"notes_lst":notes_lst,}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/get_notes", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//---------------------
	// TAGS
	//---------------------
	// ADD_TAGS

	http.HandleFunc("/tags/add_tags", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /tags/add_tags ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------
			//INPUT
			i_map, gf_err := gf_rpc_lib.Get_http_input("/tags/add_tags", p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}
			//------------

			gf_err = pipeline__add_tags(i_map, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/tags/add_tags", "add_tags pipeline failed", gf_err, p_resp, p_runtime_sys)
				return
			}

			data_map := map[string]interface{}{}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/add_tags", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	
	//---------------------
	// GET_OBJECTS_WITH_TAG

	http.HandleFunc("/tags/objects", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST - /tags/objects ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			objects_with_tag_lst, gf_err := pipeline__get_objects_with_tag(p_req, p_resp, 
				gf_templates.tag_objects__tmpl,
				gf_templates.tag_objects__subtemplates_names_lst,
				p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/tags/objects", "failed to get html/json objects with tag", gf_err, p_resp, p_runtime_sys)
				return
			}

			//if the response_format was HTML then objects_with_tag_lst is nil,
			//in which case there is no json to send back
			if objects_with_tag_lst != nil {

				data_map := map[string]interface{}{"objects_with_tag_lst":objects_with_tag_lst,}
				gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime_sys)
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/objects", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//---------------------

	return nil
}