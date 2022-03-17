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
	"time"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_rpc_handler_runtime struct {
	Mux                *http.ServeMux
	Metrics            *GF_metrics
	Store_run_bool     bool
	Sentry_hub         *sentry.Hub
	Auth_login_url_str string // url redirected too if user not logged in and tries to access auth handler
}

type GF_rpc_handler_run struct {
	Class_str          string  `bson:"class_str"` // Rpc_Handler_Run
	Handler_url_str    string  `bson:"handler_url_str"`
	Start_time__unix_f float64 `bson:"start_time__unix_f"`
	End_time__unix_f   float64 `bson:"end_time__unix_f"`
}

type handler_http func(context.Context, http.ResponseWriter, *http.Request) (map[string]interface{}, *gf_core.GF_error)

//-------------------------------------------------
// HTTP
func Create_handler__http(p_path_str string,
	p_handler_fun  handler_http,
	p_runtime_sys *gf_core.Runtime_sys) {

	Create_handler__http_with_metrics(p_path_str,
		p_handler_fun,
		nil,
		false,
		p_runtime_sys)
}

//-------------------------------------------------
// HTTP_WITH_AUTH
func CreateHandlerHTTPwithAuth(p_auth_bool bool, // if handler uses authentication or not
	p_path_str        string,
	p_handler_fun     handler_http,
	p_handler_runtime *GF_rpc_handler_runtime,
	p_runtime_sys     *gf_core.Runtime_sys) {

	handler_fun := getHandler(p_auth_bool,
		p_path_str,
		p_handler_fun,
		p_handler_runtime.Metrics,
		p_handler_runtime.Store_run_bool,
		p_handler_runtime.Sentry_hub,
		&p_handler_runtime.Auth_login_url_str,
		p_runtime_sys)

	p_handler_runtime.Mux.HandleFunc(p_path_str, handler_fun)
}

//-------------------------------------------------
// HTTP_WITH_MUX
func CreateHandlerHTTPwithMux(p_path_str string,
	p_handler_fun    handler_http,
	p_mux            *http.ServeMux,
	p_metrics        *GF_metrics,
	p_store_run_bool bool,
	p_sentry_hub     *sentry.Hub,
	p_runtime_sys    *gf_core.Runtime_sys) {

	handler_fun := getHandler(false, // p_auth_bool
		p_path_str,
		p_handler_fun,
		p_metrics,
		p_store_run_bool,
		p_sentry_hub,
		nil,
		p_runtime_sys)

	p_mux.HandleFunc(p_path_str, handler_fun)
}

//-------------------------------------------------
// HTTP_WITH_METRICS
func Create_handler__http_with_metrics(p_path_str string,
	p_handler_fun    handler_http,
	p_metrics        *GF_metrics,
	p_store_run_bool bool,
	p_runtime_sys    *gf_core.Runtime_sys) {

	handler_fun := getHandler(false, // p_auth_bool
		p_path_str,
		p_handler_fun,
		p_metrics,
		p_store_run_bool,
		nil,
		nil,
		p_runtime_sys)

	http.HandleFunc(p_path_str, handler_fun)
}

//-------------------------------------------------
func getHandler(p_auth_bool bool,
	p_path_str       string,
	p_handler_fun    handler_http,
	p_metrics        *GF_metrics,
	p_store_run_bool bool,
	p_sentry_hub         *sentry.Hub,
	p_auth_login_url_str *string,
	pRuntimeSys          *gf_core.Runtime_sys) func(pResp http.ResponseWriter, pReq *http.Request) {

	handler_fun := func(pResp http.ResponseWriter, pReq *http.Request) {

		start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

		path_str := pReq.URL.Path

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
					nil, pResp, pRuntimeSys)
			},
			"gf_rpc_lib", pRuntimeSys)

		//------------------
		// METRICS

		if p_metrics != nil {
			if counter, ok := p_metrics.Handlers_counters_map[p_path_str]; ok {
				counter.Inc()
			}
		}

		//------------------
		ctx := pReq.Context()
		fmt.Println("CTX >>>", ctx, p_sentry_hub)

		// FIX!! - when creating additional http servers outside the default global
		//         http server and default global Sentry context, the clone sentry hub
		//         is being passed in explicitly.
		//         figure out a cleaner way to abstract all Sentry details from this handler wrapper.
		var hub *sentry.Hub
		if p_sentry_hub == nil {

			// use the default global pre-created Hub (one thats used by the main go-routine)
			hub = sentry.GetHubFromContext(ctx)
		} else {
			hub = p_sentry_hub
		}
		hub.Scope().SetTag("url", path_str)

		//------------------
		// TRACE
		span_op_str := p_path_str
		span__root  := sentry.StartSpan(ctx, span_op_str)
		defer span__root.Finish()

		ctx_root := span__root.Context()

		//------------------
		// AUTH

		var ctxAuth context.Context
		if p_auth_bool {
			
			// SESSION_VALIDATE
			validBool, userIdentifierStr, gfErr := gf_session.Validate(pReq, ctx, pRuntimeSys)
			if gfErr != nil {
				Error__in_handler(path_str,
					fmt.Sprintf("handler %s failed to execute auth validation of session", path_str),
					nil, pResp, pRuntimeSys)
				return
			}

			if !validBool {

				if p_auth_login_url_str != nil {

					// redirect user to login url
					http.Redirect(pResp,
						pReq,
						*p_auth_login_url_str,
						301)

				} else {
					Error__in_handler(path_str,
						fmt.Sprintf("user not authenticated to access handler %s", path_str),
						nil, pResp, pRuntimeSys)
				}
				return
			}


			ctxAuth = context.WithValue(ctx_root, "gf_user_name", userIdentifierStr)
		} else {
			ctxAuth = ctx_root
		}

		//------------------
		// HANDLER
		dataMap, gfErr := p_handler_fun(ctxAuth, pResp, pReq)

		//------------------
		// TRACE
		span__root.Finish()

		//------------------
		// ERROR
		if gfErr != nil {
			Error__in_handler(p_path_str,
				fmt.Sprintf("handler %s failed", p_path_str),
				gfErr, pResp, pRuntimeSys)
			return
		}
		
		//------------------
		// OUTPUT
		// IMPORTANT!! - currently testing if dataMap != nil because routes that render templates
		//               (render html into body) should not also return a JSON map
		if dataMap != nil {
			Http_respond(dataMap, "OK", pResp, pRuntimeSys)
		}

		//------------------

		end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

		if p_store_run_bool {
			go func() {
				Store_rpc_handler_run(p_path_str, start_time__unix_f, end_time__unix_f, pRuntimeSys)
			}()
		}
	}
	return handler_fun
}

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
func Store_rpc_handler_run(p_handler_url_str string,
	p_start_time__unix_f float64,
	p_end_time__unix_f   float64,
	pRuntimeSys          *gf_core.Runtime_sys) *gf_core.GF_error {

	// dont store a run if there is no DB initialized
	if pRuntimeSys.Mongo_db == nil {
		return nil
	}

	run := &GF_rpc_handler_run{
		Class_str:          "rpc_handler_run",
		Handler_url_str:    p_handler_url_str,
		Start_time__unix_f: p_start_time__unix_f,
		End_time__unix_f:   p_end_time__unix_f,
	}

	ctx           := context.Background()
	coll_name_str := "gf_rpc_handler_run"

	gf_err := gf_core.Mongo__insert(run,
		coll_name_str,
		map[string]interface{}{
			"handler_url_str":    p_handler_url_str,
			"caller_err_msg_str": "failed to insert rpc_handler_run",
		},
		ctx,
		pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}

//-------------------------------------------------
func Error__in_handler(p_handler_url_path_str string,
	p_user_msg_str string,
	p_gf_err       *gf_core.GF_error,
	p_resp         http.ResponseWriter,
	p_runtime_sys  *gf_core.Runtime_sys) {

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
}