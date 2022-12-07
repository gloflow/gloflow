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
	"net/http"
	"context"
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

				// IMPORTANT!! - Auth0 expects the "state" variable name
				gfAuth0sessionCookieNameStr := "state"
				gfAuth0sessionCookieValStr  := string(sessionIDstr)
				sessionTTLhoursInt, _ := gf_identity_core.GetSessionTTL()

				gf_session.SetOnReq(gfAuth0sessionCookieNameStr,
					gfAuth0sessionCookieValStr,
					pResp,
					sessionTTLhoursInt)
				
				//------------------
				// HTTP_REDIRECT - redirect user to Auth0 login url
				http.Redirect(pResp,
					pReq,
					pAuthenticator.AuthCodeURL(string(sessionIDstr)),
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
			gfSessionIDauth0providedStr := gf_core.GF_ID(qsMap["state"][0])


			// this is cookie state set on the users browser by the GF auth0 login handler
			// /v1/identity/auth0/login before the user got redirected to Auth0's Login screen
			gfSessionIDcookieNameStr       := "state"
			cookieExistsBool, cookieValStr := gf_session.GetFromReq(gfSessionIDcookieNameStr, pReq)
			if !cookieExistsBool {
				gfErr := gf_core.ErrorCreate("auth0 login callback request is missing the 'state' cookie argument",
					"verify__input_data_missing_in_req_error",
					map[string]interface{}{},
					nil, "gf_identity_lib", pRuntimeSys)
				return nil, gfErr
			}

			gfSessionIDstr := gf_core.GF_ID(cookieValStr)

			input := &gf_identity_core.GFauth0inputLoginCallback{
				CodeStr:                     codeStr,
				GFsessionIDauth0providedStr: gfSessionIDauth0providedStr,
				GFsessionIDstr:              gfSessionIDstr,
			}

			//------------------

			// ADD!! - create login attempt record in the DB

			//------------------
			gfErr := gf_identity_core.Auth0loginCallbackPipeline(input,
				pAuthenticator,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			//------------------
			// HTTP_REDIRECT - redirect user to logged in page
			http.Redirect(pResp,
				pReq,
				"/landing/main",
				301)
			
			//------------------

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

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