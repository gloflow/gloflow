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

package gf_identity

import (
	"net/http"
	"context"
	"encoding/base64"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_session"
)

//------------------------------------------------

func initHandlersAuth0(pKeyServer *gf_identity_core.GFkeyServerInfo,
	pHTTPmux       *http.ServeMux,
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
		AuthLoginURLstr: "/v1/identity/login_ui",
		AuthKeyServer:   pKeyServer,
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
					nil, "gf_identity", pRuntimeSys)
				return nil, gfErr
			}

			//------------------
			// CODE
			// 45 char length string
			codeStr := qsMap["code"][0]

			//------------------
			// STATE
			// "state" arg - this is the session ID set by GF handler initially via a cookie,
			//               and Auth0 enforces the "state" symbol name.
			if _, ok := qsMap["state"]; !ok {
				gfErr := gf_core.ErrorCreate("auth0 login callback request is missing the 'state' qs argument",
					"verify__input_data_missing_in_req_error",
					map[string]interface{}{},
					nil, "gf_identity", pRuntimeSys)
				return nil, gfErr
			}

			// state is base64 encoded session_id, so it needs to be decoded and casted into an GF_ID 
			auth0providedStateBase64str := qsMap["state"][0]
			auth0providedStateStr, _    := base64.StdEncoding.DecodeString(auth0providedStateBase64str)
			sessionID                := gf_core.GF_ID(auth0providedStateStr)
			
			//------------------
			input := &gf_identity_core.GFauth0inputLoginCallback{
				CodeStr:           codeStr,
				SessionID:         sessionID,
				Auth0appDomainStr: pConfig.Auth0domainStr,
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

			//---------------------
			// COOKIES 
			
			// SESSION_ID - sets gf_sess cookie
			gf_session.CreateSessionIDcookie(string(sessionID), pResp)

			// JWT - sets "Authorization" cookie
			gf_session.CreateAuthCookie(output.JWTtokenStr, pResp)
			
			//------------------
			// HTTP_REDIRECT - redirect user to logged in page
			
			/*
			IMPORTANT!! - currently in the Auth0 login_callback, once all login initialization is complete,
						the client is redirected to the login_page with a QS arg "login_success" set to 1.
						this indicates to the login page that it should not initialize into its standard login state,
						but instead should indicate to the user that login has succeeded and then via JS
						(with a time delay) redirect to Home.
						this is needed to handle a race condition where if user was redirected to home via a
						server redirect (HTTP 3xx response) the browser wouldnt have time to set the needed
						auth cookies that are necessary for the Home handler to authenticate.
						the recommended solution for this is to do the redirection via client/JS code with a slight
						time delay, giving the browser time to set the cookies.
			*/
			homeUrlStr := "/v1/identity/login_ui?login_success=1"
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