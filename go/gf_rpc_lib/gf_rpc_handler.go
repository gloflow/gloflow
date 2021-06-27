/*
MIT License

Copyright (c) 2021 Ivan Trajkovic

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
	"net/http"
	"context"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_rpc_handler_run struct {
	Class_str          string  `bson:"class_str"` // Rpc_Handler_Run
	Handler_url_str    string  `bson:"handler_url_str"`
	Start_time__unix_f float64 `bson:"start_time__unix_f"`
	End_time__unix_f   float64 `bson:"end_time__unix_f"`
}

type handler_http func(context.Context, http.ResponseWriter, *http.Request) (map[string]interface{}, *gf_core.Gf_error)

//-------------------------------------------------
func Create_handler__http(p_path_str string,
	p_handler_fun  handler_http,
	p_runtime_sys *gf_core.Runtime_sys) {

	http.HandleFunc(p_path_str, func(p_resp http.ResponseWriter, p_req *http.Request) {

		path_str := p_req.URL.Path

		//------------------
		// PANIC_HANDLING

		// IMPORTANT!! - only defered functions are run when a panic initiates in a goroutine
		//               as execution unwinds up the call-stack. in Panic__check_and_handle() 
		//               recover() is executed for check for panic conditions. if panic exists
		//               it is treated as an error that gets processed, and the go routine exits.

		user_msg__internal_str := "gf_rpc handler panicked"
		defer gf_core.Panic__check_and_handle(user_msg__internal_str,
			map[string]interface{}{"handler_path_str": path_str},
			// oncomplete_fn
			func() {
				
				// IMPORTANT!! - if a panic occured, send a HTTP response to the client,
				//               and then proceed to process the panic as an error 
				//               with gf_core.Panic__check_and_handle()
				Error__in_handler(path_str,
					fmt.Sprintf("handler %s failed unexpectedly", path_str),
					nil, p_resp, p_runtime_sys)
			},
			"gf_rpc_lib", p_runtime_sys)

		//------------------
		ctx := p_req.Context()

		hub := sentry.GetHubFromContext(ctx)
		hub.Scope().SetTag("url", path_str)

		//------------------
		// TRACE
		span_op_str := p_path_str
		span__root := sentry.StartSpan(ctx, span_op_str)
		defer span__root.Finish()

		//------------------
		// HANDLER
		data_map, gf_err := p_handler_fun(span__root.Context(), p_resp, p_req)

		//------------------
		// TRACE
		span__root.Finish()

		//------------------
		// ERROR
		if gf_err != nil {

			Error__in_handler(p_path_str,
				fmt.Sprintf("handler %s failed", p_path_str),
				gf_err, p_resp, p_runtime_sys)

			return
		}
		
		//------------------
		// OUTPUT
		if data_map != nil {
			Http_respond(data_map, "OK", p_resp, p_runtime_sys)
		}

		//------------------
	})
}

//-------------------------------------------------
func Http_respond(p_data interface{},
	p_status_str  string,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) {
	// p_runtime_sys.Log_fun("FUN_ENTER", "gf_rpc_utils.Http_respond()")

	r_byte_lst, _ := json.Marshal(map[string]interface{}{
		"status": p_status_str,
		"data":   p_data,
	})
	
	p_resp.Header().Set("Content-Type", "application/json")
	p_resp.Write(r_byte_lst)
}

//-------------------------------------------------
func Store_rpc_handler_run(p_handler_url_str string,
	p_start_time__unix_f float64,
	p_end_time__unix_f   float64,
	p_runtime_sys        *gf_core.Runtime_sys) *gf_core.Gf_error {
	// p_runtime_sys.Log_fun("FUN_ENTER", "gf_rpc_handler.Store_rpc_handler_run()")

	run := &GF_rpc_handler_run{
		Class_str:          "rpc_handler_run", // FIX!! - thi should be "rpc_handler_run"
		Handler_url_str:    p_handler_url_str,
		Start_time__unix_f: p_start_time__unix_f,
		End_time__unix_f:   p_end_time__unix_f,
	}

	ctx           := context.Background()
	coll_name_str := p_runtime_sys.Mongo_coll.Name()

	gf_err := gf_core.Mongo__insert(run,
		coll_name_str,
		map[string]interface{}{
			"handler_url_str":    p_handler_url_str,
			"caller_err_msg_str": "failed to insert rpc_handler_run",
		},
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	/*err := p_runtime_sys.Mongo_coll.Insert(run)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert rpc_handler_run",
            "mongodb_insert_error",
            map[string]interface{}{"handler_url_str": p_handler_url_str,},
            err, "gf_rpc_lib", p_runtime_sys)
		return gf_err
	}*/

	return nil
}

//-------------------------------------------------
func Error__in_handler(p_handler_url_path_str string,
	p_user_msg_str string,
	p_gf_err       *gf_core.Gf_error,
	p_resp         http.ResponseWriter,
	p_runtime_sys  *gf_core.Runtime_sys) {
	// p_runtime_sys.Log_fun("FUN_ENTER", "gf_rpc_handler.Error__in_handler()")

	status_str := "ERROR"
	data_map   := map[string]interface{}{
		"handler_error_user_msg": p_user_msg_str,
	}

	if p_gf_err != nil {
		data_map["gf_error_type"]     = p_gf_err.Type_str
		data_map["gf_error_user_msg"] = p_gf_err.User_msg_str
	}

	// DEBUG
	if p_runtime_sys.Debug_bool {
		data_map["error"] = p_gf_err.Error
	}

	Http_respond(data_map, status_str, p_resp, p_runtime_sys)

	/*http.Error(p_resp,
	p_usr_msg_str,
	http.StatusInternalServerError)*/
}