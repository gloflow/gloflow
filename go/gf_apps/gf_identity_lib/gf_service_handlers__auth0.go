/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_identity_lib

import (
	// "fmt"
	"net/http"
	"context"
	"encoding/base64"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
)

//------------------------------------------------

func initHandlersAuth0(pHTTPmux *http.ServeMux,
	pAuthenticator *gf_auth0.GFauthenticator,
	pConfig        *gf_auth0.GFconfig,
	pServiceInfo   *gf_identity_core.GFserviceInfo,
	pRuntimeSys    *gf_core.RuntimeSys) {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/identity/auth0/login",
		"/v1/identity/auth0/login_callback",
		"/v1/identity/auth0/logout_callback",
	}
	metricsGroupNameStr := "auth0"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, pServiceInfo.NameStr, handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:             pHTTPmux,
		Metrics:         metrics,
		StoreRunBool:    true,
		SentryHub:       nil,

		AuthSubsystemTypeStr: pServiceInfo.AuthSubsystemTypeStr,
		
		// url redirected too if user not logged in and tries to access auth handler
		AuthLoginURLstr: "/landing/main",
	}

	//---------------------
	// Auth0 may need to redirect back to the application's Login Initiation endpoint, using OIDC third-party initiated login
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/auth0/login",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				sessionIDstr, gfErr := gf_identity_core.Auth0loginPipeline(pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				//------------------
				// HTTP_REDIRECT - redirect user to Auth0 login url
				auth0appStateBase64str := base64.StdEncoding.EncodeToString([]byte(string(sessionIDstr)))

				http.Redirect(pResp,
					pReq,
					pAuthenticator.AuthCodeURL(auth0appStateBase64str),
					301)

				//------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// user redirected to this URL by Auth0 on successful login
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/auth0/login_callback",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			//------------------
			// INPUT
			qsMap := pReq.URL.Query()

			// "code" arg - provided by Auth0
			if _, ok := qsMap["code"]; !ok {
				gfErr := gf_core.ErrorCreate("auth0 login callback request is missing the 'code' qs argument",
					"verify__input_data_missing_in_req_error",
					map[string]interface{}{},
					nil, "gf_identity_lib", pRuntimeSys)
				return nil, gfErr
			}

			// 45 char length string
			codeStr := qsMap["code"][0]

			// "state" arg - this is the session ID set by GF handler initially via a cookie,
			//               and Auth0 enforces the "state" symbol name.
			if _, ok := qsMap["state"]; !ok {
				gfErr := gf_core.ErrorCreate("auth0 login callback request is missing the 'state' qs argument",
					"verify__input_data_missing_in_req_error",
					map[string]interface{}{},
					nil, "gf_identity_lib", pRuntimeSys)
				return nil, gfErr
			}
			auth0providedStateBase64str := qsMap["state"][0]
			auth0providedStateStr, _ := base64.StdEncoding.DecodeString(auth0providedStateBase64str)
			gfSessionIDauth0providedStr := gf_core.GF_ID(auth0providedStateStr)

			input := &gf_identity_core.GFauth0inputLoginCallback{
				CodeStr:                     codeStr,
				GFsessionIDauth0providedStr: gfSessionIDauth0providedStr,
			}

			//------------------

			// ADD!! - create login attempt record in the DB

			//------------------
			output, gfErr := gf_identity_core.Auth0loginCallbackPipeline(input,
				pAuthenticator,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			sessionIDstr := output.SessionIDstr

			//---------------------
			// SET_SESSION_ID - sets gf_sess cookie on all future requests
			sessionDataStr := string(sessionIDstr)
			gf_session.Create(sessionDataStr, pResp)

			//------------------
			// HTTP_REDIRECT - redirect user to logged in page
			
			homeUrlStr := "/v1/home/view"
			http.Redirect(pResp,
				pReq,
				homeUrlStr,
				301)
			
			//------------------

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------


	//---------------------
	//---------------------
	// FINISH!! - logout handler
	//---------------------
	//---------------------


	// user redirected to this URL by Auth0 on successful logout
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/auth0/logout_callback",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				

				// ADD!! - mark auth0 session as deleted
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------

}