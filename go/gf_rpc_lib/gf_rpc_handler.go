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
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GFrpcHandlerRuntime struct {
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

type handler_http func(context.Context, http.ResponseWriter, *http.Request) (map[string]interface{}, *gf_core.GFerror)

//-------------------------------------------------
// HTTP
func Create_handler__http(p_path_str string,
	p_handler_fun handler_http,
	pRuntimeSys   *gf_core.RuntimeSys) {

	Create_handler__http_with_metrics(p_path_str,
		p_handler_fun,
		nil,
		false,
		pRuntimeSys)
}

//-------------------------------------------------
// HTTP_WITH_AUTH
func CreateHandlerHTTPwithAuth(p_auth_bool bool, // if handler uses authentication or not
	p_path_str        string,
	p_handler_fun     handler_http,
	p_handler_runtime *GFrpcHandlerRuntime,
	pRuntimeSys       *gf_core.RuntimeSys) {

	handler_fun := getHandler(p_auth_bool,
		p_path_str,
		p_handler_fun,
		p_handler_runtime.Metrics,
		p_handler_runtime.Store_run_bool,
		p_handler_runtime.Sentry_hub,
		&p_handler_runtime.Auth_login_url_str,
		pRuntimeSys)

	p_handler_runtime.Mux.HandleFunc(p_path_str, handler_fun)
}

//-------------------------------------------------
// HTTP_WITH_MUX
func CreateHandlerHTTPwithMux(p_path_str string,
	p_handler_fun handler_http,
	p_mux         *http.ServeMux,
	pMetrics      *GF_metrics,
	pStoreRunBool bool,
	p_sentry_hub  *sentry.Hub,
	pRuntimeSys   *gf_core.RuntimeSys) {

	handler_fun := getHandler(false, // p_auth_bool
		p_path_str,
		p_handler_fun,
		pMetrics,
		pStoreRunBool,
		p_sentry_hub,
		nil,
		pRuntimeSys)

	p_mux.HandleFunc(p_path_str, handler_fun)
}

//-------------------------------------------------
// HTTP_WITH_METRICS
func Create_handler__http_with_metrics(p_path_str string,
	p_handler_fun handler_http,
	pMetrics      *GF_metrics,
	pStoreRunBool bool,
	pRuntimeSys   *gf_core.RuntimeSys) {

	handler_fun := getHandler(false, // p_auth_bool
		p_path_str,
		p_handler_fun,
		pMetrics,
		pStoreRunBool,
		nil,
		nil,
		pRuntimeSys)

	http.HandleFunc(p_path_str, handler_fun)
}

//-------------------------------------------------
func getHandler(p_auth_bool bool,
	pPathStr         string,
	p_handler_fun    handler_http,
	pMetrics         *GF_metrics,
	pStoreRunBool    bool,
	pSentryHub       *sentry.Hub,
	pAuthLoginURLstr *string,
	pRuntimeSys      *gf_core.RuntimeSys) func(pResp http.ResponseWriter, pReq *http.Request) {

	handler_fun := func(pResp http.ResponseWriter, pReq *http.Request) {

		start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

		pathStr := pReq.URL.Path

		//------------------
		// PANIC_HANDLING

		// IMPORTANT!! - only defered functions are run when a panic initiates in a goroutine
		//               as execution unwinds up the call-stack. in Panic__check_and_handle() 
		//               recover() is executed for check for panic conditions. if panic exists
		//               it is treated as an error that gets processed, and the go routine exits.

		user_msg__internal_str := "gf_rpc handler panicked"
		defer gf_core.Panic__check_and_handle(user_msg__internal_str,
			map[string]interface{}{"handler_path_str": pathStr},
			// oncomplete_fn
			func() {
				
				// IMPORTANT!! - if a panic occured, send a HTTP response to the client,
				//               and then proceed to process the panic as an error 
				//               with gf_core.Panic__check_and_handle()
				Error__in_handler(pathStr,
					fmt.Sprintf("handler %s failed unexpectedly", pathStr),
					nil, pResp, pRuntimeSys)
			},
			"gf_rpc_lib", pRuntimeSys)

		//------------------
		// METRICS

		if pMetrics != nil {
			if counter, ok := pMetrics.Handlers_counters_map[pPathStr]; ok {
				counter.Inc()
			}
		}

		//------------------
		ctx := pReq.Context()
		fmt.Println("CTX >>>", ctx, pSentryHub)

		// FIX!! - when creating additional http servers outside the default global
		//         http server and default global Sentry context, the clone sentry hub
		//         is being passed in explicitly.
		//         figure out a cleaner way to abstract all Sentry details from this handler wrapper.
		var hub *sentry.Hub
		if pSentryHub == nil {

			// use the default global pre-created Hub (one thats used by the main go-routine)
			hub = sentry.GetHubFromContext(ctx)
		} else {
			hub = pSentryHub
		}
		hub.Scope().SetTag("url", pathStr)

		//------------------
		// TRACE
		spanOpStr := pPathStr
		spanRoot  := sentry.StartSpan(ctx, spanOpStr)
		defer spanRoot.Finish()

		ctxRoot := spanRoot.Context()

		//------------------
		// AUTH

		var ctxAuth context.Context
		if p_auth_bool {

			// SESSION_VALIDATE
			validBool, userIdentifierStr, gfErr := gf_session.ValidateOrRedirectToLogin(pReq,
				pResp,
				pAuthLoginURLstr,
				ctx,
				pRuntimeSys)
			if gfErr != nil {
				Error__in_handler(pathStr,
					fmt.Sprintf("handler %s failed to execute/validate a auth session", pathStr),
					nil, pResp, pRuntimeSys)
				return
			}

			// SESSION_NOT_VALID
			if !validBool {

				// METRICS
				if pMetrics != nil {
					pMetrics.HandlersAuthSessionInvalidCounter.Inc()
				}

				// if a login_url is not defined then return error, otherwise redirect to this
				// login_url   
				if pAuthLoginURLstr == nil {
					Error__in_handler(pathStr,
						fmt.Sprintf("user not authenticated to access handler %s", pathStr),
						nil, pResp, pRuntimeSys)
				}
				return
			}

			ctxAuth = context.WithValue(ctxRoot, "gf_user_id", userIdentifierStr)
		} else {
			ctxAuth = ctxRoot
		}

		//------------------
		// HANDLER
		outputDataMap, gfErr := p_handler_fun(ctxAuth, pResp, pReq)

		//------------------
		// TRACE
		spanRoot.Finish()

		//------------------
		// ERROR
		if gfErr != nil {
			Error__in_handler(pPathStr,
				fmt.Sprintf("handler %s failed", pPathStr),
				gfErr, pResp, pRuntimeSys)
			return
		}
		
		//------------------
		// OUTPUT
		// IMPORTANT!! - currently testing if dataMap != nil because routes that render templates
		//               (render html into body) should not also return a JSON map
		if outputDataMap != nil {
			HTTPrespond(outputDataMap, "OK", pResp, pRuntimeSys)
		}

		//------------------

		end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

		if pStoreRunBool {
			go func() {
				Store_rpc_handler_run(pPathStr, start_time__unix_f, end_time__unix_f, pRuntimeSys)
			}()
		}
	}
	return handler_fun
}

//-------------------------------------------------
func Store_rpc_handler_run(p_handler_url_str string,
	p_start_time__unix_f float64,
	p_end_time__unix_f   float64,
	pRuntimeSys          *gf_core.RuntimeSys) *gf_core.GFerror {

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

	gf_err := gf_core.MongoInsert(run,
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