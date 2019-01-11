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

package main

import (
	"time"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func init_handlers(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_tagger_service_handlers.init_handlers()")

	//---------------------------
	//NOTES
	//---------------------------
	//ADD_NOTE

	http.HandleFunc("/tags/add_note",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /tags/add_note ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------
			//JSON INPUT
			i_map            := map[string]interface{}{}
			body_bytes_lst,_ := ioutil.ReadAll(p_req.Body)
		    err              := json.Unmarshal(body_bytes_lst,&i_map)

			if err != nil {
				gf_rpc_lib.Error__in_handler("/tags/add_note", err, "add_note pipeline received bad i_map input", p_resp, p_runtime_sys)
				return
			}
			//------------

			err = pipeline__add_note(i_map, p_runtime_sys)
			if err != nil {
				gf_rpc_lib.Error__in_handler("/tags/add_note", err, "add_note pipeline failed", p_resp, p_runtime_sys)
				return
			}

			data_map := map[string]interface{}{}
			gf_rpc_lib.Http_Respond(data_map, "OK", p_resp, p_runtime_sys)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/add_note", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------------
	//GET_notes

	http.HandleFunc("/tags/get_notes",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /tags/get_notes ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			notes_lst,err := pipeline__get_notes(p_req, p_runtime_sys)
			if err != nil {
				gf_rpc_lib.Error__in_handler("/tags/get_notes", err, "get_notes pipeline failed", p_resp, p_runtime_sys)
				return 
			}

			data_map := map[string][]*Note{"notes_lst":notes_lst,}
			gf_rpc_lib.Http_Respond(data_map, "OK", p_resp, p_runtime_sys)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/get_notes", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------------
	//TAGS
	//---------------------------
	//ADD_TAGS

	http.HandleFunc("/tags/add_tags",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /tags/add_tags ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------
			//JSON INPUT
			i_map            := map[string]interface{}{}
			body_bytes_lst,_ := ioutil.ReadAll(p_req.Body)
		    err              := json.Unmarshal(body_bytes_lst,&i_map)

			if err != nil {
				gf_rpc_lib.Error__in_handler("/tags/add_tags", err, "add_tags pipeline received bad i_map input", p_resp, p_runtime_sys)
				return
			}
			//------------

			err = pipeline__add_tags(i_map, p_runtime_sys)
			if err != nil {
				gf_rpc_lib.Error__in_handler("/tags/add_tags", err, "add_tags pipeline failed", p_resp, p_runtime_sys)
				return
			}

			data_map := map[string]interface{}{}
			gf_rpc_lib.Http_Respond(data_map, "OK", p_resp, p_log_fun)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/add_tags", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------------
	//GET_OBJECTS_WITH_TAG

	tag_objects__tmpl,err := template.New("gf_tag_objects.html").ParseFiles("./templates/gf_tag_objects.html")
	if err != nil {
		return err
	}
	http.HandleFunc("/tags/objects",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /tags/objects ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0


			objects_with_tag_lst,err := pipeline__get_objects_with_tag(p_req, p_resp, tag_objects__tmpl, p_runtime_sys)
			if err != nil {
				gf_rpc_lib.Error__in_handler("/tags/objects", err, "failed to get html/json objects with tag", p_resp, p_runtime_sys)
				return
			}

			//if the response_format was HTML then objects_with_tag_lst is nil,
			//in which case there is no json to send back
			if objects_with_tag_lst != nil {

				data_map := map[string]interface{}{"objects_with_tag_lst":objects_with_tag_lst,}
				gf_rpc_lib.Http_Respond(data_map, "OK", p_resp, p_log_fun)
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/tags/objects", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------------

	return nil
}