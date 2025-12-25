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
	"fmt"
	"net/http"
	"net/url"
	"context"
	"encoding/base64"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
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
		"/v1/identity/auth0/api_token/generate",
		"/v1/identity/auth0/login",
		"/v1/identity/auth0/login_callback",
		"/v1/identity/auth0/logout",
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

		EnableEventsBool: true,
	}

	//---------------------
	// API_TOKEN_GENERATE

	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/auth0/api_token/generate",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			//---------------------
			// INPUT

			userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

			iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}
			
			audienceStr := iMap["audience_str"].(string)
			

			// ADD!! - set these to proper secret values
			appClientIDstr := ""
			appClientSecretStr := ""

			input := &gf_identity_core.GFauth0inputAPItokenGenerate{
				UserID:             userID,
				AppClientIDstr:     appClientIDstr,
				AppClientSecretStr: appClientSecretStr,
				AudienceStr:        audienceStr,
				Auth0appDomainStr:  pConfig.Auth0domainStr,
			}

			//---------------------
			
			if pReq.Method == "GET" {

				tokenStr, gfErr := gf_identity_core.Auth0apiTokenGeneratePipeline(input,
					pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				outputMap := map[string]interface{}{
					"token_str": tokenStr,
				}
				return outputMap, nil

				//------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// Auth0 may need to redirect back to the application's Login Initiation endpoint, using OIDC third-party initiated login
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/auth0/login",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//---------------------
				// INPUT

				// extract QS arg "redirect_url" - where to redirect user after successful login.
				// allows user to be redirected to the page they were on before being redirected to Auth0 login.
				loginSuccessRedirectURLstr := pReq.URL.Query().Get("redirect_url")
				if loginSuccessRedirectURLstr == "" {
					loginSuccessRedirectURLstr = "/"
				} else {
					decodedURLstr, err := url.QueryUnescape(loginSuccessRedirectURLstr)
					if err != nil {
						gfErr := gf_core.ErrorCreate("failed to decode redirect_url",
							"http__input_data_invalid_error",
							map[string]interface{}{},
							err, "gf_identity", pRuntimeSys)
						return nil, gfErr
					}

					loginSuccessRedirectURLstr = decodedURLstr
				}
				
				//---------------------

				sessionIDstr, gfErr := gf_identity_core.Auth0loginPipeline(loginSuccessRedirectURLstr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// HTTP_REDIRECT - redirect user to Auth0 login url

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)

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
			sessionID                   := gf_core.GF_ID(auth0providedStateStr)
			
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

			/*
			IMPORTANT!! - strict mode set to None, since the user is authned via Auth0 and on successful auth 
				user is redirected by Auth0 to this handler. Even though the cookies are set by this server/domain,
				because the request was a redirect from a third-party domain (Auth0), the browser will treat the
				cookies as third-party cookies and will block them if the SameSite attribute is set to "Strict".
				So initially they're set as "None" to allow the browser to set them, and then on subsequent
				login-finalize request they're set to "Strict" to prevent CSRF attacks.
			*/
			sameSiteStrictBool := false

			// SESSION_ID - sets gf_sess cookie
			gf_identity_core.CreateSessionIDcookie(string(sessionID),
				pServiceInfo.DomainForAuthCookiesStr,
				sameSiteStrictBool,
				pResp)

			// JWT - sets "Authorization" cookie
			gf_identity_core.CreateAuthCookie(output.JWTtokenStr,
				pServiceInfo.DomainForAuthCookiesStr,
				sameSiteStrictBool,
				pResp)
			
			//------------------
			// HTTP_REDIRECT - redirect user to logged in page
			
			// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
			gf_core.HTTPdisableCachingOfResponse(pResp)


			var redirectURLstr string
			if output.LoginSuccessRedirectURLstr == "" {
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


				defaultLoginSuccessURLstr := "/v1/identity/login_ui?login_success=1"
				redirectURLstr = defaultLoginSuccessURLstr

				
			} else {
				redirectURLstr = output.LoginSuccessRedirectURLstr
			}
			
			http.Redirect(pResp,
				pReq,
				redirectURLstr,
				301)

			//------------------

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// this can only be called by authenticated users.

	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/auth0/login_finalize",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {
				//------------------
				// INPUT

				jwtTokenStr, foundBool, gfErr := gf_identity_core.JWTgetTokenFromRequest(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				if !foundBool {
					gfErr := gf_core.ErrorCreate("'Authorization' is missing from auth0 login_finalize handler request cookies",
						"auth_missing_cookie",
						map[string]interface{}{},
						nil, "gf_identity", pRuntimeSys)
					return nil, gfErr
				}
				sessionID, sessionIDfoundBool := gf_identity_core.GetSessionID(pReq, pRuntimeSys)
				if !sessionIDfoundBool {
					gfErr := gf_core.ErrorCreate("'session_id' is missing from auth0 login_finalize handler request cookies",
						"auth_missing_cookie",
						map[string]interface{}{},
						nil, "gf_identity", pRuntimeSys)
					return nil, gfErr
				}

				//---------------------
				// COOKIES

				/*
				IMPORTANT!! - overwrite cookies for the user in same-site strict mode, allowing only this domain
					to read/write cookies. this is done to prevent CSRF attacks.
				*/
				sameSiteStrictBool := true

				// SESSION_ID - sets gf_sess cookie
				gf_identity_core.CreateSessionIDcookie(string(sessionID),
					pServiceInfo.DomainForAuthCookiesStr,
					sameSiteStrictBool,
					pResp)

				// JWT - sets "Authorization" cookie
				gf_identity_core.CreateAuthCookie(jwtTokenStr,
					pServiceInfo.DomainForAuthCookiesStr,
					sameSiteStrictBool,
					pResp)
				
				//------------------
				gf_core.HTTPdisableCachingOfResponse(pResp)

				//------------------
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/auth0/logout",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {
				
				//---------------------
				// INPUT

				// extract QS arg "redirect_url" - where to redirect user after successful login.
				// allows user to be redirected to the page they were on before being redirected to Auth0 login.
				logoutSuccessRedirectURLstr := pReq.URL.Query().Get("redirect_url")
				if logoutSuccessRedirectURLstr == "" {
					logoutSuccessRedirectURLstr = ""
				} else {
					decodedURLstr, err := url.QueryUnescape(logoutSuccessRedirectURLstr)
					if err != nil {
						gfErr := gf_core.ErrorCreate("failed to decode redirect_url",
							"http__input_data_invalid_error",
							map[string]interface{}{},
							err, "gf_identity", pRuntimeSys)
						return nil, gfErr
					}

					logoutSuccessRedirectURLstr = decodedURLstr
				}
				
				sessionID, existsBool := gf_identity_core.GetSessionID(pReq, pRuntimeSys)
				if !existsBool {
					gfErr := gf_core.ErrorCreate("session_id is missing from auth0 logout handler request cookies",
						"http_cookie",
						map[string]interface{}{},
						nil, "gf_identity", pRuntimeSys)
					return nil, gfErr
				}

				//---------------------
				// DB - update session with logout URL
				gfErr := gf_identity_core.DBsqlUpdateSessionLogoutRedirectURL(sessionID,
					logoutSuccessRedirectURLstr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------
				/*
				IMPORTANT!! - redirect user to Auth0 logout url, which will log them out of Auth0
					and any third-party identity providers. this does not log them out of the GF
					application, which needs to be done separatelly. 
				*/

				logoutURLstr := fmt.Sprintf("https://%s/v2/logout?client_id=%s&returnTo=%s",
					pConfig.Auth0domainStr,
					pConfig.Auth0clientIDstr,
					pConfig.Auth0logoutCallbackURLstr)

				pRuntimeSys.LogNewFun("DEBUG", "redirecting to Auth0 logout url...",
					map[string]interface{}{"logout_url_str": logoutURLstr})

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)

				// HTTP_REDIRECT
				http.Redirect(pResp,
					pReq,
					logoutURLstr,
					301)

				//---------------------
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// user redirected to this URL by Auth0 on successful logout.
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/auth0/logout_callback",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//------------------
				// INPUT

				/*
				session_id in this handler has to be fetched directly from request cookies instead from
				context because this handler is not authed, and so in gf_rpc the auth path is not chosen
				and context is not enriched with the session_id.
				*/
				sessionID, existsBool := gf_identity_core.GetSessionID(pReq, pRuntimeSys)
				if !existsBool {
					gfErr := gf_core.ErrorCreate("session_id is missing from auth0 logout_callback handler request cookies",
						"http_cookie",
						map[string]interface{}{},
						nil, "gf_identity", pRuntimeSys)
					return nil, gfErr
				}

				//------------------

				domainForAuthCookiesStr := *pServiceInfo.DomainForAuthCookiesStr

				logoutSuccessRedirectURLstr, gfErr := gf_identity_core.LogoutPipeline(sessionID,
					domainForAuthCookiesStr,
					pResp,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				

				var redirectURLstr string
				if logoutSuccessRedirectURLstr == "" {
					redirectURLstr = "/landing/main"
				} else {
					redirectURLstr = logoutSuccessRedirectURLstr
				}

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)
				
				// HTTP_REDIRECT
				http.Redirect(pResp,
					pReq,
					redirectURLstr,
					301)

			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------

}