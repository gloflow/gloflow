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

package gf_rpc_lib

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)


//-------------------------------------------------
func Get_http_input(p_handler_url_path_str string,
	p_resp        http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_rpc_utils.Get_http_input()")

	var i map[string]interface{}
	body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
	err               := json.Unmarshal(body_bytes_lst, &i)

	if err != nil {
		gf_err := gf_core.Error__create("failed to parse json http input",
			"json_decode_error",
			map[string]interface{}{"handler_url_path_str": p_handler_url_path_str,},
			err, "gf_rpc_lib", p_runtime_sys)

		Error__in_handler(p_handler_url_path_str,
			fmt.Sprintf("failed parsing http-request input JSON in - %s", p_handler_url_path_str), //p_user_msg_str
			gf_err,
			p_resp,
			p_runtime_sys)
		return nil, gf_err
	}

	return i, nil
}

//-------------------------------------------------
func Get_response_format(p_qs_map map[string][]string,
	p_runtime_sys *gf_core.Runtime_sys) string {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_rpc_utils.Get_response_format()")

	response_format_str := "html" //default - "h" - HTML
	if f_lst, ok := p_qs_map["f"]; ok {
		response_format_str = f_lst[0] //user supplied value
	}

	return response_format_str
}

//-------------------------------------------------
func Http_CORS_preflight_handle(p_req *http.Request,
	p_resp http.ResponseWriter) {
	
	// CORS - preflight request
	if p_req.Method == "OPTIONS" {
		p_resp.Header().Set("Access-Control-Allow-Origin", "*")
		p_resp.Header().Set("Access-Control-Allow-Origin", "Origin, X-Requested-With, Content-Type, Accept")
	}
}





//-------------------------------------------------
func Panic__check_and_handle(p_resp http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) {

	// IMPORTANT!! - if a panic occured, send a HTTP response to the client,
	//               and then proceed to process the panic as an error 
	//               with gf_core.Panic__check_and_handle()
	if panic_info := recover(); panic_info != nil {

		path_str := p_req.URL.Path
		Error__in_handler(path_str,
			fmt.Sprintf("handler %s failed unexpectedly", path_str),
			nil, p_resp, p_runtime_sys)


		user_msg__internal_str := "gf_rpc handler panicked"
		gf_core.Panic__check_and_handle(user_msg__internal_str,
			map[string]interface{}{"handler_path_str": path_str},
			"gf_rpc", p_runtime_sys)
	}
}