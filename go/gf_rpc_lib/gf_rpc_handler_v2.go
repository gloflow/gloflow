package gf_rpc_lib

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/sentry-go"

	gf_core "github.com/gloflow/gloflow/go/gf_core"
	gf_identity_core "github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//-------------------------------------------------
// CREATE_HANDLERS_V2_HTTP - creates a batch of handlers

func CreateHandlersV2http(pHandlersLst []gf_core.HTTPhandlerV2info,
	pHTTPmux              *http.ServeMux,
	pAuthSubsystemTypeStr string,
	pAuthLoginURLstr      string,
	pKeyServer            *gf_identity_core.GFkeyServerInfo,
	pMetricsGroupNameStr  string,
	pRuntimeSys           *gf_core.RuntimeSys) {

	//---------------------
	// METRICS
	metrics := MetricsCreateForHandlers(pMetricsGroupNameStr, "gf_solo", pHandlersLst, pRuntimeSys)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &GFrpcHandlerRuntime {
		Mux:          pHTTPmux,
		Metrics:      metrics,
		StoreRunBool: true,
		SentryHub:    nil,

		// AUTH
		AuthSubsystemTypeStr: pAuthSubsystemTypeStr,
		AuthLoginURLstr:      pAuthLoginURLstr,
		AuthKeyServer:        pKeyServer,
	}

	pRuntimeSys.LogNewFun("DEBUG", "creating v2 handlers...",
		map[string]interface{}{"auth_subsystem_type": rpcHandlerRuntime.AuthSubsystemTypeStr})

	for _, handlerInfo := range pHandlersLst {

		// CREATE_HANDLER
		CreateHandlerV2http(handlerInfo.NameStr,
			handlerInfo.AuthStr,
			handlerInfo.PathStr,
			handlerInfo.DomainsLst,
			handlerInfo.HandlerFun,
			rpcHandlerRuntime,
			pRuntimeSys)
	}
}

//-------------------------------------------------
// HTTP_WITH_AUTH

func CreateHandlerV2http(pNameStr string, // optional name for handler (used for metrics)
	pAuthStr string, // if handler uses authentication or not
	pPathStr        string,
	pDomainsLst     []string,
	pHandlerFun     gf_core.HTTPhandlerV2,
	pHandlerRuntime *GFrpcHandlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) {

	// AUTH0
	if pHandlerRuntime.AuthSubsystemTypeStr == gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__AUTH0 {

		// check auth key_server has been initialized and passed to the handler runtime
		if (pAuthStr == gf_core.AUTH_REQUIRED || pAuthStr == gf_core.AUTH_OPTIONAL) && pHandlerRuntime.AuthKeyServer == nil {
			panic("Auth key_server has to be defined!")
		}

	}

	// HANDLER_FUN
	appHandlerFun := getHandlerV2(pNameStr,
		pAuthStr,
		pPathStr,
		pHandlerFun,
		pHandlerRuntime.Metrics,
		pHandlerRuntime.StoreRunBool,
		pHandlerRuntime.SentryHub,
		&pHandlerRuntime.AuthLoginURLstr,
		pRuntimeSys)

	//-------------------------------------------------
	// VALIDATE_SESSION
	validateSessionFun := func(pResp http.ResponseWriter, pReq *http.Request) (exitReqBool bool, ctx context.Context) {

		ctx = pReq.Context()
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
		validBool, userID, sessionID, gfErr := gf_identity_core.SessionValidateOrRedirectToLogin(pReq,
			pResp,
			pHandlerRuntime.AuthKeyServer,
			pHandlerRuntime.AuthSubsystemTypeStr,
			&pHandlerRuntime.AuthLoginURLstr,
			authRedirectOnFailBool,
			ctx,
			pRuntimeSys)

		if gfErr != nil {
			ErrorInHandler(pathStr,
				fmt.Sprintf("handler %s failed to execute/validate auth session", pathStr),
				nil, pResp, pRuntimeSys)

			//-----------------------
			// EXIT_REQ
			exitReqBool = true
			return exitReqBool, nil

			//-----------------------
		}

		pRuntimeSys.LogNewFun("DEBUG", `>>>>>>>>>>>>>>>>> v2 session validation...`,
			map[string]interface{}{
				"path":                     pathStr,
				"valid":                    validBool,
				"user_id":                  userID,
				"session_id":               sessionID,
				"auth_redirect_on_failure": authRedirectOnFailBool,
				"auth_subsystem_type":      pHandlerRuntime.AuthSubsystemTypeStr,
			})

		//-----------------------
		// REQUIRED - SESSION_NOT_VALID
		if !validBool && pAuthStr == gf_core.AUTH_REQUIRED {

			// METRICS
			if pHandlerRuntime.Metrics != nil {
				pHandlerRuntime.Metrics.HandlersAuthSessionInvalidCounter.Inc()
			}

			// if no redirection of auth failure is specified (which happens in SessionValidateOrRedirectToLogin())
			// return an error
			if !authRedirectOnFailBool {
				msgStr := "unauthorized access"
				ErrorInHandler(pathStr,
					msgStr,
					nil, pResp, pRuntimeSys)
			}

			//-----------------------
			// EXIT_REQ
			exitReqBool = true
			return exitReqBool, nil

			//-----------------------


		//-----------------------
		// OPTIONAL - SESSION_NOT_VALID
		} else if !validBool && pAuthStr == gf_core.AUTH_OPTIONAL {

			// DO NOTHING - session is not valid, but auth is optional, so handler should still
			// continue to process the request.

		//-----------------------
		// SESSION_VALID
		} else {

			//-----------------------
			// AUTH_CONTEXT - attach user_id and session_id to a handler context
			ctxUserID := context.WithValue(ctx, "gf_user_id", *userID)
			ctxAuth   := context.WithValue(ctxUserID, "gf_session_id", string(*sessionID))

			//-----------------------

			// session is valid, no need to interupt the handler from further execution
			exitReqBool = false
			return exitReqBool, ctxAuth
		}

		//-----------------------

		// session is valid, no need to interupt the handler from further execution
		exitReqBool = false
		return exitReqBool, ctx
	}

	//-------------------------------------------------
	// CORS

	CORSfun := func(pResp http.ResponseWriter, pReq *http.Request) {

		//-----------------------
		// METRICS
		if pHandlerRuntime.Metrics != nil {
			pHandlerRuntime.Metrics.HandlersAuthSessionCORScounter.Inc()
		}

		//-----------------------

		if pRuntimeSys.ExternalHooks != nil &&
			pRuntimeSys.ExternalHooks.CORSoriginDomainsLst != nil {

			originStr := pReq.Header.Get("Origin")

			// check if the origin domain is in the list of allowed domains
			if gf_core.StringInList(originStr, pRuntimeSys.ExternalHooks.CORSoriginDomainsLst) {

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

	//-------------------------------------------------

	if pAuthStr == gf_core.AUTH_REQUIRED || pAuthStr == gf_core.AUTH_OPTIONAL {

		//-------------------------------------------------
		authHandlerFun := func(pResp http.ResponseWriter, pReq *http.Request) {

			pathStr := pReq.URL.Path

			pRuntimeSys.LogNewFun("INFO", "------------------> HTTP REQ", map[string]interface{}{"path_str": pathStr})

			//-----------------------
			// ALLOWED_DOMAINS
			// if the handler has specified allowed domains, check if the request domain is in the list

			if !isFromAllowedDomain(pReq, pResp, pDomainsLst, pRuntimeSys) {
				return
			}

			//-----------------------
			// METRICS
			if pHandlerRuntime.MetricsGlobal != nil {
				pHandlerRuntime.MetricsGlobal.HandlersAuthCounter.Inc()
			}

			//-----------------------
			// VALIDATE_SESSION
			exitReqBool, ctxAuth := validateSessionFun(pResp, pReq)
			if exitReqBool {
				return
			}

			// REQ_ID
			reqIDstr := genRequestID()
			ctxWithReqID := context.WithValue(ctxAuth, "gf_req_id", reqIDstr)

			//-----------------------
			// CORS
			// if the user has supplied CORS domains, check if the request origin domain is in the list

			// get the origin domain of the request
			originStr := pReq.Header.Get("Origin")

			pRuntimeSys.LogNewFun("DEBUG", `>>>>>>>>>>>>>>>>> req headers...`,
				map[string]interface{}{
					"headers_str": spew.Sdump(pReq.Header),
				})

			if originStr != "" {

				//-----------------------
				// METRICS
				if pHandlerRuntime.MetricsGlobal != nil {
					pHandlerRuntime.MetricsGlobal.HandlersCORScounter.Inc()
				}

				//-----------------------

				CORSfun(pResp, pReq)
			}

			//-----------------------
			// handle OPTIONS preflight request
			if pReq.Method == http.MethodOptions {

				CORSfun(pResp, pReq)

				// pResp.WriteHeader(http.StatusNoContent)
				return
			}

			//-----------------------
			// APP_HANDLER - external app request handler function, executed with an
			//               authenticated context.

			appHandlerFun(pResp, pReq.WithContext(ctxWithReqID))

			//-----------------------
		}

		//-------------------------------------------------

		pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(authHandlerFun))

	} else {

		//-------------------------------------------------
		wrappedHandlerFun := func(pResp http.ResponseWriter, pReq *http.Request) {

			domainStr := pReq.Host
			pathStr := pReq.URL.Path

			pRuntimeSys.LogNewFun("INFO", "------------------> HTTP REQ",
				map[string]interface{}{"path": pathStr, "domain": domainStr})

			//-----------------------
			// ALLOWED_DOMAINS
			// if the handler has specified allowed domains, check if the request domain is in the list

			if !isFromAllowedDomain(pReq, pResp, pDomainsLst, pRuntimeSys) {
				return
			}

			//-----------------------
			// CORS
			originStr := pReq.Header.Get("Origin")
			if originStr != "" {

				//-----------------------
				// METRICS
				if pHandlerRuntime.MetricsGlobal != nil {
					pHandlerRuntime.MetricsGlobal.HandlersCORScounter.Inc()
				}

				//-----------------------
			}

			appHandlerFun(pResp, pReq)
		}

		//-------------------------------------------------
		pHandlerRuntime.Mux.Handle(pPathStr, http.HandlerFunc(wrappedHandlerFun))
	}
}

//-------------------------------------------------

func isFromAllowedDomain(pReq *http.Request,
	pResp       http.ResponseWriter,
	pDomainsLst []string,
	pRuntimeSys *gf_core.RuntimeSys) bool {

	// if no domains are specified, allow all
	if len(pDomainsLst) == 0 {
		return true
	}


	pathStr := pReq.URL.Path
	domainStr := pReq.Host


	spew.Dump(pDomainsLst)
	fmt.Println(domainStr)

	forAllowedDomainBool := false
	for _, allowedDomainStr := range pDomainsLst {
		if domainStr == allowedDomainStr {
			forAllowedDomainBool = true
			break
		}
	}

	if !forAllowedDomainBool {

		pRuntimeSys.LogNewFun("DEBUG", "req not for allowed domain...",
			map[string]interface{}{"path": pReq.URL.Path, "domain": domainStr, "allowed_domains": pDomainsLst})

		msgStr := "unauthorized domain access"
		ErrorInHandler(pathStr,
			msgStr,
			nil, pResp, pRuntimeSys)
	}

	return forAllowedDomainBool
}

//-------------------------------------------------

func getHandlerV2(pNameStr string,
	pAuthStr string,
	pPathStr         string,
	pHandlerFun      gf_core.HTTPhandlerV2,
	pMetrics         *GFmetrics,
	pStoreRunBool    bool,
	pSentryHub       *sentry.Hub,
	pAuthLoginURLstr *string,
	pRuntimeSys      *gf_core.RuntimeSys) func(pResp http.ResponseWriter, pReq *http.Request) {

	handlerFun := func(pResp http.ResponseWriter, pReq *http.Request) {

		startTimeUNIXf := float64(time.Now().UnixNano())/1000000000.0
		pathStr := pReq.URL.Path

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
		// HOOKS - run request pre-process callback
		if pRuntimeSys.ExternalHooks != nil && pRuntimeSys.ExternalHooks.RPCreqPreProcessCallback != nil {

			continueBool, gfErr := pRuntimeSys.ExternalHooks.RPCreqPreProcessCallback(pReq, pResp, ctxRoot, pRuntimeSys)
			if gfErr != nil {
				ErrorInHandler(pPathStr,
					fmt.Sprintf("handler %s failed", pPathStr),
					gfErr, pResp, pRuntimeSys)
				return
			}

			if !continueBool {

				// FINALIZE
				endTimeUNIXf := float64(time.Now().UnixNano())/1000000000.0
				finalizeHandler(pPathStr, spanRoot, startTimeUNIXf, endTimeUNIXf, pStoreRunBool, pRuntimeSys)
				return
			}
		}

		//------------------
		// HANDLER
		outputDataMap, gfErr := pHandlerFun(ctxRoot, pResp, pReq)

		//------------------
		// FINALIZE
		endTimeUNIXf := float64(time.Now().UnixNano())/1000000000.0
		finalizeHandler(pPathStr, spanRoot, startTimeUNIXf, endTimeUNIXf, pStoreRunBool, pRuntimeSys)

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
	}
	return handlerFun
}
