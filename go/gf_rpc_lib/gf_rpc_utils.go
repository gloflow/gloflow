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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
type Gf_rpc_handler_run struct {
	Class_str          string  `bson:"class_str"` //Rpc_Handler_Run
	Handler_url_str    string  `bson:"handler_url_str"`
	Start_time__unix_f float64 `bson:"start_time__unix_f"`
	End_time__unix_f   float64 `bson:"end_time__unix_f"`
}
//-------------------------------------------------
func Get_http_input(p_handler_url_path_str string,
	p_resp        http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_rpc_utils.Get_http_input()")

	var i map[string]interface{}
	body_bytes_lst,_ := ioutil.ReadAll(p_req.Body)
	err              := json.Unmarshal(body_bytes_lst,&i)

	if err != nil {
		gf_err := gf_core.Error__create("failed to parse json http input",
			"json_unmarshal_error",
			&map[string]interface{}{"handler_url_path_str":p_handler_url_path_str,},
			err,"gf_rpc_lib",p_runtime_sys)

		Error__in_handler(p_handler_url_path_str,
			"failed parsing http-request input JSON in - "+p_handler_url_path_str, //p_user_msg_str
			gf_err,
			p_resp,p_runtime_sys)
		return nil,gf_err
	}

	return i, nil
}
//-------------------------------------------------
func Get_response_format(p_qs_map map[string][]string,
	p_runtime_sys *gf_core.Runtime_sys) string {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_rpc_utils.Get_response_format()")

	response_format_str := "html" //default - "h" - HTML
	if f_lst,ok := p_qs_map["f"]; ok {
		response_format_str = f_lst[0] //user supplied value
	}

	return response_format_str
}
//-------------------------------------------------
func Http_Respond(p_data interface{},
	p_status_str  string,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_rpc_utils.Http_Respond()")

	r_lst,_ := json.Marshal(map[string]interface{}{
			"status_str":p_status_str,
			"data":      p_data,
		})
	
	p_resp.Header().Set("Content-Type","application/json")
	p_resp.Write(r_lst)
}
//-------------------------------------------------
func Store_rpc_handler_run(p_handler_url_str string,
	p_start_time__unix_f float64,
	p_end_time__unix_f   float64,
	p_runtime_sys        *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_rpc_utils.Store_rpc_handler_run()")

	run := &Gf_rpc_handler_run{
		Class_str:         "Rpc_Handler_Run", //FIX!! - thi should be "rpc_handler_run"
		Handler_url_str:   p_handler_url_str,
		Start_time__unix_f:p_start_time__unix_f,
		End_time__unix_f:  p_end_time__unix_f,
	}

	err := p_runtime_sys.Mongodb_coll.Insert(run)
	if err != nil {
		gf_err := gf_core.Error__create("failed to insert rpc_handler_run",
            "mongodb_insert_error",
            &map[string]interface{}{"handler_url_str":p_handler_url_str,},
            err,"gf_rpc_lib",p_runtime_sys)
		return gf_err
	}

	return nil
}
//-------------------------------------------------
func Error__in_handler(p_handler_url_path_str string,
	p_user_msg_str string,
	p_gf_err       *gf_core.Gf_error,
	p_resp         http.ResponseWriter,
	p_runtime_sys  *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_rpc_utils.Error__in_handler()")

	status_str := "ERROR"
	data_map   := map[string]interface{}{
		"handler_error_user_msg_str":p_user_msg_str,
		"gf_error_type_str":         p_gf_err.Type_str,
		"gf_error_user_msg_str":     p_gf_err.User_msg_str,
	}
	Http_Respond(data_map,status_str,p_resp,p_runtime_sys)

	/*http.Error(p_resp,
	p_usr_msg_str,
	http.StatusInternalServerError)*/
}