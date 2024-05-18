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
	// "fmt"
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

func getRequestID(pReq *http.Request) string {
    if reqIDstr, ok := pReq.Context().Value("gf_req_id").(string); ok {
        return reqIDstr
    }
    return ""
}

func genRequestID() string {
	reqIDstr := uuid.New().String()
	return reqIDstr
}

//-------------------------------------------------

func HTTPrespond(p_data interface{},
	p_status_str string,
	p_resp       http.ResponseWriter,
	pRuntimeSys  *gf_core.RuntimeSys) {

	r_byte_lst, _ := json.Marshal(map[string]interface{}{
		"status": p_status_str,
		"data":   p_data,
	})
	
	p_resp.Header().Set("Content-Type", "application/json")
	p_resp.Write(r_byte_lst)
}

//-------------------------------------------------

func ErrorInHandler(p_handler_url_path_str string,
	p_user_msg_str string,
	p_gf_err       *gf_core.GFerror,
	p_resp         http.ResponseWriter,
	pRuntimeSys    *gf_core.RuntimeSys) {

	statusStr := "ERROR"
	dataMap   := map[string]interface{}{
		"handler_error_user_msg": p_user_msg_str,
	}

	if p_gf_err != nil {
		dataMap["gf_error_type"]     = p_gf_err.Type_str
		dataMap["gf_error_user_msg"] = p_gf_err.User_msg_str
	}

	// DEBUG
	if pRuntimeSys.Debug_bool {
		dataMap["error"] = p_gf_err.Error
	}

	HTTPrespond(dataMap, statusStr, p_resp, pRuntimeSys)
}

//-------------------------------------------------

func GetResponseFormat(p_qs_map map[string][]string,
	pRuntimeSys *gf_core.RuntimeSys) string {

	response_format_str := "html" // default - "h" - HTML
	if f_lst, ok := p_qs_map["f"]; ok {
		response_format_str = f_lst[0] // user supplied value
	}

	return response_format_str
}

//-------------------------------------------------

func HTTPcorsPreflightHandle(p_req *http.Request,
	p_resp http.ResponseWriter) {
	
	// CORS - preflight request
	if p_req.Method == "OPTIONS" {
		p_resp.Header().Set("Access-Control-Allow-Origin", "*")
		p_resp.Header().Set("Access-Control-Allow-Origin", "Origin, X-Requested-With, Content-Type, Accept")
	}
}

//-------------------------------------------------
/*
// FIX!! - passing in structs as interface{} doesnt preserve their tags.
func Get_http_input_to_struct(p_input_struct interface{},
	p_resp        http.ResponseWriter,
	p_req         *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	handler_url_path_str := p_req.URL.Path
	body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
	err               := json.Unmarshal(body_bytes_lst, &p_input_struct)
		
	if err != nil {
		gf_err := gf_core.ErrorCreate("failed to parse json http input",
			"json_decode_error",
			map[string]interface{}{"handler_url_path_str": handler_url_path_str,},
			err, "gf_rpc_lib", pRuntimeSys)

		Error__in_handler(handler_url_path_str,
			fmt.Sprintf("failed parsing http-request input JSON in - %s", handler_url_path_str),
			gf_err,
			p_resp,
			pRuntimeSys)
		return gf_err
	}

	return nil
}*/