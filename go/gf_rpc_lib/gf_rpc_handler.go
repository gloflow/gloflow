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
	// "github.com/auth0/go-jwt-middleware/v2"
	// "github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFrpcHandlerRuntime struct {
	Mux             *http.ServeMux
	Metrics         *GF_metrics
	StoreRunBool    bool
	SentryHub       *sentry.Hub

	AuthSubsystemTypeStr string
	AuthLoginURLstr      string // url redirected too if user not logged in and tries to access auth handler
	AuthKeyServer        *gf_identity_core.GFkeyServerInfo
}

type GFrpcHandlerRun struct {
	ClassStr       string  `bson:"class_str"`
	HandlerURLstr  string  `bson:"handler_url_str"`
	StartTimeUNIXf float64 `bson:"start_time__unix_f"`
	EndTimeUNIXf   float64 `bson:"endTimeUNIXf"`
}

type handlerHTTP func(context.Context, http.ResponseWriter, *http.Request) (map[string]interface{}, *gf_core.GFerror)

//-------------------------------------------------
// HTTP

func CreateHandlerHTTP(pPathStr string,
	pHandlerFun handlerHTTP,
	pRuntimeSys *gf_core.RuntimeSys) {

	CreateHandlerHTTPwithMetrics(pPathStr,
		pHandlerFun,
		nil,
		false,
		pRuntimeSys)
}

//-------------------------------------------------
// HTTP_WITH_AUTH

func CreateHandlerHTTPwithAuth(pAuthBool bool, // if handler uses authentication or not
	pPathStr        string,
	pHandlerFun     handlerHTTP,
	pHandlerRuntime *GFrpcHandlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) {

	// check auth key_server has been initialized and passed to the handler runtime
	if pHandlerRuntime.AuthKeyServer == nil {
		panic("Auth key_server has to be defined!")
	}

	// AUTH0
	if pHandlerRuntime.AuthSubsystemTypeStr == gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__AUTH0 {
		
	} else {
		// BUILTIN_AUTH
		// set the builtin auth_subsystem type as the default value if another value is not set
		pHandlerRuntime.AuthSubsystemTypeStr = gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__BUILTIN
	}

	// HANDLER_FUN
	appHandlerFun := getHandler(pAuthBool,
		pPathStr,
		pHandlerFun,
		pHandlerRuntime.Metrics,
		pHandlerRuntime.StoreRunBool,
		pHandlerRuntime.SentryHub,
		&pHandlerRuntime.AuthLoginURLstr,
		pRuntimeSys)

	
	switch pHandlerRuntime.AuthSubsystemTypeStr {

	//------------------
	// BUILTIN
	case gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__BUILTIN:

		if pAuthBool {

			//-------------------------------------------------
			builtinHandlerFun := func(pResp http.ResponseWriter, pReq *http.Request) {
				
				ctx := pReq.Context()
				pathStr := pReq.URL.Path

				// SESSION_VALIDATE
				validBool, userIdentifierStr, gfErr := gf_session.ValidateOrRedirectToLogin(pReq,
					pResp,
					pHandlerRuntime.AuthKeyServer,
					pHandlerRuntime.AuthSubsystemTypeStr,
					&pHandlerRuntime.AuthLoginURLstr,
					ctx,
					pRuntimeSys)

				if gfErr != nil {
					ErrorInHandler(pathStr,
						fmt.Sprintf("handler %s failed to execute/validate a auth session", pathStr),
						nil, pResp, pRuntimeSys)
					return
				}

				// SESSION_NOT_VALID
				if !validBool {

					// METRICS
					if pHandlerRuntime.Metrics != nil {
						pHandlerRuntime.Metrics.HandlersAuthSessionInvalidCounter.Inc()
					}

					// if a login_url is not defined then return error, otherwise redirect to this login_url   
					if pHandlerRuntime.AuthLoginURLstr == "" {
						ErrorInHandler(pathStr,
							fmt.Sprintf("user not authenticated to access handler %s", pathStr),
							nil, pResp, pRuntimeSys)
					}
					return
				}

				ctxAuth := context.WithValue(ctx, "gf_user_id", userIdentifierStr)


				appHandlerFun(pResp, pReq.WithContext(ctxAuth))
			}

			//-------------------------------------------------

			pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(builtinHandlerFun))

		} else {
			pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(appHandlerFun))
		}

	//------------------
	// AUTH0
	case gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__AUTH0:

		if pAuthBool {

			//-------------------------------------------------
			auth0handlerFun := func(pResp http.ResponseWriter, pReq *http.Request) {

				
			
				appHandlerFun(pResp, pReq)
			}

			//-------------------------------------------------

			pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(auth0handlerFun))

		} else {
			pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(appHandlerFun))
		}
	}

	//------------------
}

//-------------------------------------------------
// HTTP_WITH_MUX

func CreateHandlerHTTPwithMux(pPathStr string,
	pHandlerFun   handlerHTTP,
	pMux          *http.ServeMux,
	pMetrics      *GF_metrics,
	pStoreRunBool bool,
	pSentryHub    *sentry.Hub,
	pRuntimeSys   *gf_core.RuntimeSys) {

	handlerFun := getHandler(false, // pAuthBool
		pPathStr,
		pHandlerFun,
		pMetrics,
		pStoreRunBool,
		pSentryHub,
		nil,
		pRuntimeSys)

	pMux.HandleFunc(pPathStr, handlerFun)
}

//-------------------------------------------------
// HTTP_WITH_METRICS

func CreateHandlerHTTPwithMetrics(pPathStr string,
	pHandlerFun   handlerHTTP,
	pMetrics      *GF_metrics,
	pStoreRunBool bool,
	pRuntimeSys   *gf_core.RuntimeSys) {

	handlerFun := getHandler(false, // pAuthBool
		pPathStr,
		pHandlerFun,
		pMetrics,
		pStoreRunBool,
		nil,
		nil,
		pRuntimeSys)

	http.HandleFunc(pPathStr, handlerFun)
}

//-------------------------------------------------

func getHandler(pAuthBool bool,
	pPathStr         string,
	pHandlerFun      handlerHTTP,
	pMetrics         *GF_metrics,
	pStoreRunBool    bool,
	pSentryHub       *sentry.Hub,
	pAuthLoginURLstr *string,
	pRuntimeSys      *gf_core.RuntimeSys) func(pResp http.ResponseWriter, pReq *http.Request) {

	handlerFun := func(pResp http.ResponseWriter, pReq *http.Request) {

		startTimeUNIXf := float64(time.Now().UnixNano())/1000000000.0
		pathStr := pReq.URL.Path

		pRuntimeSys.LogNewFun("INFO", "------------------> HTTP REQ", map[string]interface{}{"path_str": pPathStr})

		//------------------
		// PANIC_HANDLING

		// IMPORTANT!! - only defered functions are run when a panic initiates in a goroutine
		//               as execution unwinds up the call-stack. in PanicCheckAndHandle() 
		//               recover() is executed for check for panic conditions. if panic exists
		//               it is treated as an error that gets processed, and the go routine exits.

		userMsgInternalStr := "gf_rpc handler panicked"
		defer gf_core.PanicCheckAndHandle(userMsgInternalStr,
			map[string]interface{}{"handler_path_str": pathStr},
			// oncomplete_fn
			func() {
				
				// IMPORTANT!! - if a panic occured, send a HTTP response to the client,
				//               and then proceed to process the panic as an error 
				//               with gf_core.PanicCheckAndHandle()
				ErrorInHandler(pathStr,
					fmt.Sprintf("handler %s failed unexpectedly", pathStr),
					nil, pResp, pRuntimeSys)
			},
			"gf_rpc_lib", pRuntimeSys)

		//------------------
		// METRICS

		if pMetrics != nil {
			if counter, ok := pMetrics.HandlersCountersMap[pPathStr]; ok {
				counter.Inc()
			}
		}

		//------------------
		ctx := pReq.Context()

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

		/*//------------------
		// AUTH

		var ctxAuth context.Context
		if pAuthBool {

			// SESSION_VALIDATE
			validBool, userIdentifierStr, gfErr := gf_session.ValidateOrRedirectToLogin(pReq,
				pResp,
				pAuthLoginURLstr,
				ctx,
				pRuntimeSys)
			if gfErr != nil {
				ErrorInHandler(pathStr,
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
					ErrorInHandler(pathStr,
						fmt.Sprintf("user not authenticated to access handler %s", pathStr),
						nil, pResp, pRuntimeSys)
				}
				return
			}

			ctxAuth = context.WithValue(ctxRoot, "gf_user_id", userIdentifierStr)
		} else {
			ctxAuth = ctxRoot
		}*/

		//------------------
		// HANDLER
		outputDataMap, gfErr := pHandlerFun(ctxRoot, pResp, pReq)

		//------------------
		// TRACE
		spanRoot.Finish()

		//------------------
		// ERROR
		if gfErr != nil {
			ErrorInHandler(pPathStr,
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

		endTimeUNIXf := float64(time.Now().UnixNano())/1000000000.0

		if pStoreRunBool {
			go func() {
				StoreRPChandlerRun(pPathStr, startTimeUNIXf, endTimeUNIXf, pRuntimeSys)
			}()
		}
	}
	return handlerFun
}

//-------------------------------------------------

func StoreRPChandlerRun(pHandlerURLstr string,
	pStartTimeUNIXf float64,
	pEndTimeUNIXf   float64,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	// dont store a run if there is no DB initialized
	if pRuntimeSys.Mongo_db == nil {
		return nil
	}

	run := &GFrpcHandlerRun{
		ClassStr:       "rpc_handler_run",
		HandlerURLstr:  pHandlerURLstr,
		StartTimeUNIXf: pStartTimeUNIXf,
		EndTimeUNIXf:   pEndTimeUNIXf,
	}

	ctx         := context.Background()
	collNameStr := "gf_rpc_handler_run"

	gfErr := gf_core.MongoInsert(run,
		collNameStr,
		map[string]interface{}{
			"handler_url_str":    pHandlerURLstr,
			"caller_err_msg_str": "failed to insert rpc_handler_run",
		},
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}