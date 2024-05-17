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

//-------------------------------------------------

func CreateHandlersHTTP(pMetricsGroupNameStr string,
	pHandlersLst          []gf_core.HTTPhandlerInfo,
	pHTTPmux              *http.ServeMux,
	pAuthSubsystemTypeStr string,
	pAuthLoginURLstr      string,
	pKeyServer            *gf_identity_core.GFkeyServerInfo,
	pRuntimeSys           *gf_core.RuntimeSys) {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{}
	for _, handlerInfo := range pHandlersLst {
		pathStr := handlerInfo.PathStr
		handlersEndpointsLst = append(handlersEndpointsLst, pathStr)
	}

	metrics := MetricsCreateForHandlers(pMetricsGroupNameStr, "gf_solo", handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &GFrpcHandlerRuntime {
		Mux:             pHTTPmux,
		Metrics:         metrics,
		StoreRunBool:    true,
		SentryHub:       nil,

		// AUTH
		AuthSubsystemTypeStr: pAuthSubsystemTypeStr,
		AuthLoginURLstr:      pAuthLoginURLstr,
		AuthKeyServer:        pKeyServer,
	}

	
	for _, handlerDescr := range pHandlersLst {

		// CREATE_HANDLER
		CreateHandlerHTTPwithAuth(handlerDescr.AuthBool,
			handlerDescr.PathStr,
			handlerDescr.HandlerFun,
			rpcHandlerRuntime,
			pRuntimeSys)
	}

}

//-------------------------------------------------
// CREATE
//-------------------------------------------------
// HTTP

func CreateHandlerHTTP(pPathStr string,
	pHandlerFun gf_core.HTTPhandler,
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
	pHandlerFun     gf_core.HTTPhandler,
	pHandlerRuntime *GFrpcHandlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) {

	// check auth key_server has been initialized and passed to the handler runtime
	if pAuthBool && pHandlerRuntime.AuthKeyServer == nil {
		panic("Auth key_server has to be defined!")
	}

	// AUTH0
	if pHandlerRuntime.AuthSubsystemTypeStr == gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__AUTH0 {
		
	} else {
		// DEFAULT - USERPASS_AUTH
		// set the userpass auth_subsystem type as the default value if another value is not set
		pHandlerRuntime.AuthSubsystemTypeStr = gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__USERPASS
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

	

	//-------------------------------------------------
	// VALIDATE_SESSION
	validateSessionFun := func(pResp http.ResponseWriter, pReq *http.Request) *context.Context {
				
		ctx := pReq.Context()
		pathStr := pReq.URL.Path

		//-----------------------
		// AUTH_REDIRECT_ON_FAIL - QS that can toggle if the user should be redirected to the auth login url
		//                         on failure to validate auth. by default we always redirect, and the user
		//                         has the ability to turn that behavior off.
		authRedirectOnFailBool := true
		valuesMap := pReq.URL.Query()
		if vLst, ok := valuesMap["auth_r"]; ok {
			if vLst[0] == "0" {
				authRedirectOnFailBool = false
			}
		}

		//-----------------------
		
		// SESSION_VALIDATE
		validBool, userIdentifierStr, sessionID, gfErr := gf_session.ValidateOrRedirectToLogin(pReq,
			pResp,
			pHandlerRuntime.AuthKeyServer,
			pHandlerRuntime.AuthSubsystemTypeStr,
			&pHandlerRuntime.AuthLoginURLstr,
			authRedirectOnFailBool,
			ctx,
			pRuntimeSys)
		
		pRuntimeSys.LogNewFun("DEBUG", `>>>>>>>>>>>>>>>>> session validation...`,
			map[string]interface{}{
				"path_str":            pathStr,
				"valid_bool":          validBool,
				"user_identifier_str": userIdentifierStr,
				"session_id_str":      sessionID,
				"auth_redirect_on_failure_bool": authRedirectOnFailBool,
				"auth_method_str":               pHandlerRuntime.AuthSubsystemTypeStr,
			})
			
		if gfErr != nil {
			ErrorInHandler(pathStr,
				fmt.Sprintf("handler %s failed to execute/validate a auth session", pathStr),
				nil, pResp, pRuntimeSys)
			return nil
		}

		// SESSION_NOT_VALID
		if !validBool {

			// METRICS
			if pHandlerRuntime.Metrics != nil {
				pHandlerRuntime.Metrics.HandlersAuthSessionInvalidCounter.Inc()
			}

			// if no redirection of auth failure is specified (which happens in ValidateOrRedirectToLogin())
			// return an error
			if !authRedirectOnFailBool {
				msgStr := "unauthorized access"
				ErrorInHandler(pathStr,
					msgStr,
					nil, pResp, pRuntimeSys)
			}
			return nil
		}

		//-----------------------
		// AUTH_CONTEXT - attach user_id and session_id to a handler context
		ctxUserID := context.WithValue(ctx, "gf_user_id", userIdentifierStr)
		ctxAuth   := context.WithValue(ctxUserID, "gf_session_id", string(sessionID))

		//-----------------------

		return &ctxAuth
	}

	//-------------------------------------------------

	if pAuthBool {

		//-------------------------------------------------
		authHandlerFun := func(pResp http.ResponseWriter, pReq *http.Request) {
			
			pRuntimeSys.LogNewFun("DEBUG", `>>>>>>>>>>>>>>>>> auth http handler...`,
				map[string]interface{}{
					"path_str":                pReq.URL.Path,
					"auth_subsystem_type_str": pHandlerRuntime.AuthSubsystemTypeStr,
				})

			//-----------------------
			// VALIDATE_SESSION
			ctxAuth := validateSessionFun(pResp, pReq)
			if ctxAuth == nil {
				return
			}

			//-----------------------
			// CORS
			// if the user has supplied CORS domains, check if the request origin domain is in the list
			if pRuntimeSys.ExternalPlugins != nil &&
				pRuntimeSys.ExternalPlugins.CORSoriginDomainsLst != nil {
				
				// get the origin domain of the request
				originStr := pReq.Header.Get("Origin")

				if originStr != "" {

					// check if the origin domain is in the list of allowed domains
					if gf_core.StringInList(originStr, pRuntimeSys.ExternalPlugins.CORSoriginDomainsLst) {
						pResp.Header().Set("Access-Control-Allow-Origin", originStr)
						pResp.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

						/*
						specify which headers are allowed to be received by CORS requests.
						if the request includes other (non-simple) headers (Authorization, Content-Type with application/json),
						its necesary to explicitly allow these headers using the Access-Control-Allow-Headers header.
						*/
						pResp.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

						/*
						The XMLHttpRequest.withCredentials property is a Boolean that indicates
						whether or not cross-site Access-Control requests should be made using
						credentials such as cookies, authorization headers or TLS client certificates.
						Setting withCredentials has no effect on same-site requests
						
						js jquery $.ajax() param:
							xhrFields: {
								withCredentials: true
							}
						*/
						pResp.Header().Set("Access-Control-Allow-Credentials", "true")
					}
				}
			}

			//-----------------------
			// APP_HANDLER - external app request handler function, executed with an
			//               authenticated context.

			appHandlerFun(pResp, pReq.WithContext(*ctxAuth))

			//-----------------------
		}

		//-------------------------------------------------

		pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(authHandlerFun))

	} else {

		pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(appHandlerFun))
	}
}

//-------------------------------------------------
// HTTP_WITH_MUX

func CreateHandlerHTTPwithMux(pPathStr string,
	pHandlerFun   gf_core.HTTPhandler,
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
	pHandlerFun   gf_core.HTTPhandler,
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
	pHandlerFun      gf_core.HTTPhandler,
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