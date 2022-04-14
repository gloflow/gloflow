/*
MIT License

Copyright (c) 2019 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
func Http_respond(p_data interface{},
	p_status_str string,
	p_resp       http.ResponseWriter,
	pRuntimeSys  *gf_core.Runtime_sys) {

	r_byte_lst, _ := json.Marshal(map[string]interface{}{
		"status": p_status_str,
		"data":   p_data,
	})
	
	p_resp.Header().Set("Content-Type", "application/json")
	p_resp.Write(r_byte_lst)
}

//-------------------------------------------------
func Error__in_handler(p_handler_url_path_str string,
	p_user_msg_str string,
	p_gf_err       *gf_core.GF_error,
	p_resp         http.ResponseWriter,
	p_runtime_sys  *gf_core.Runtime_sys) {

	statusStr := "ERROR"
	dataMap   := map[string]interface{}{
		"handler_error_user_msg": p_user_msg_str,
	}

	if p_gf_err != nil {
		dataMap["gf_error_type"]     = p_gf_err.Type_str
		dataMap["gf_error_user_msg"] = p_gf_err.User_msg_str
	}

	// DEBUG
	if p_runtime_sys.Debug_bool {
		dataMap["error"] = p_gf_err.Error
	}

	Http_respond(dataMap, statusStr, p_resp, p_runtime_sys)
}

//-------------------------------------------------
/*
// FIX!! - passing in structs as interface{} doesnt preserve their tags.
func Get_http_input_to_struct(p_input_struct interface{},
	p_resp        http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {

	handler_url_path_str := p_req.URL.Path
	body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
	err               := json.Unmarshal(body_bytes_lst, &p_input_struct)
		
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse json http input",
			"json_decode_error",
			map[string]interface{}{"handler_url_path_str": handler_url_path_str,},
			err, "gf_rpc_lib", p_runtime_sys)

		Error__in_handler(handler_url_path_str,
			fmt.Sprintf("failed parsing http-request input JSON in - %s", handler_url_path_str),
			gf_err,
			p_resp,
			p_runtime_sys)
		return gf_err
	}

	return nil
}*/

//-------------------------------------------------
func Get_http_input(p_resp http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {

	handler_url_path_str := p_req.URL.Path

	var i map[string]interface{}
	body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)

	// parse body bytes only if they're larger than 0
	if len(body_bytes_lst) > 0 {
		err := json.Unmarshal(body_bytes_lst, &i)

		if err != nil {
			gf_err := gf_core.Error__create("failed to parse json http input",
				"json_decode_error",
				map[string]interface{}{"handler_url_path_str": handler_url_path_str,},
				err, "gf_rpc_lib", p_runtime_sys)

			Error__in_handler(handler_url_path_str,
				fmt.Sprintf("failed parsing http-request input JSON in - %s", handler_url_path_str), // p_user_msg_str
				gf_err,
				p_resp,
				p_runtime_sys)
			return nil, gf_err
		}
	}

	return i, nil
}

//-------------------------------------------------
func Get_response_format(p_qs_map map[string][]string,
	p_runtime_sys *gf_core.Runtime_sys) string {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_rpc_utils.Get_response_format()")

	response_format_str := "html" // default - "h" - HTML
	if f_lst, ok := p_qs_map["f"]; ok {
		response_format_str = f_lst[0] // user supplied value
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