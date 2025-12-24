/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_policy"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------

func initHandlers(pAuthLoginURLstr string,
	pTemplatesPathsMap map[string]string,
	pKeyServer         *gf_identity_core.GFkeyServerInfo,
	pHTTPmux           *http.ServeMux,
	pRPCglobalMetrics  *gf_rpc_lib.GFglobalMetrics,
	pServiceInfo       *gf_identity_core.GFserviceInfo,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES

	gfTemplates, gfErr := templatesLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/identity/logged_in",
		"/v1/identity/me",
		"/v1/identity/policy/update",
		"/v1/identity/login_ui",
		"/v1/identity/email_confirm",
		"/v1/identity/mfa_confirm",
		"/v1/identity/update",
		"/v1/identity/register_invite_email",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, pServiceInfo.NameStr, handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:                  pHTTPmux,
		Metrics:              metrics,
		MetricsGlobal:        pRPCglobalMetrics,
		StoreRunBool:         true,
		SentryHub:            nil,

		// AUTH
		AuthSubsystemTypeStr: pServiceInfo.AuthSubsystemTypeStr,
		AuthLoginURLstr:      pAuthLoginURLstr,
		AuthKeyServer:        pKeyServer,
	}

	//---------------------
	// LOGGED_IN - used to check efficiently by the front-end if the user is logged in
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/logged_in",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {
				outputMap := map[string]interface{}{
					"logged_in_bool": true,
				}
				return outputMap, nil
			}

			// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
			gf_core.HTTPdisableCachingOfResponse(pResp)

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// USERS_GET_ME
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/me",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//---------------------
				// INPUT

				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				input := &gf_identity_core.GFuserGetInput{
					UserID: userID,
				}

				//---------------------

				output, gfErr := gf_identity_core.UsersPipelineGet(input, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)

				outputMap := map[string]interface{}{
					"user_name_str":         output.UserNameStr,
					"user_id":               userID,
					"screen_name_str":       output.ScreenNameStr,
					"email_str":             output.EmailStr,
					"description_str":       output.DescriptionStr,
					"profile_image_url_str": output.ProfileImageURLstr,
					"banner_image_url_str":  output.BannerImageURLstr,
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)
		
	//---------------------
	// POLICY
	//---------------------
	// POLICY_UPDATE
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/policy/update",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT

				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var targetResourceIDstr gf_core.GF_ID
				if targetResourceIDinputStr, ok := inputMap["target_resource_id_str"]; ok {
					targetResourceIDstr = gf_core.GF_ID(targetResourceIDinputStr.(string))
				}

				var polidyIDstr gf_core.GF_ID
				if polidyIDinputStr, ok := inputMap["policy_id_str"]; ok {
					polidyIDstr = gf_core.GF_ID(polidyIDinputStr.(string))
				}

				//---------------------

				
				output, gfErr := gf_policy.PipelineUpdate(targetResourceIDstr, polidyIDstr, userIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)

				//---------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"policy_exists_bool": output.PolicyExistsBool,
				}
				return dataMap, nil

				//---------------------
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// VAR
	//---------------------
	// LOGIN_UI
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/login_ui",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				templateRenderedStr, gfErr := viewRenderTemplateLogin(pServiceInfo.AuthSubsystemTypeStr,
					gfTemplates.loginTmpl,
					gfTemplates.loginSubtemplatesNamesLst,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				pResp.Write([]byte(templateRenderedStr))

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// EMAIL_LOGIN
	// NO_AUTH

	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/email/logn",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {
				//---------------------
				// INPUT
				httpInput, gfErr := gf_identity_core.HTTPgetEmailLoginInput(pReq)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				gfErr = gf_identity_core.UsersEmailLoginPipeline(httpInput,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	// EMAIL_LOGIN_CONFIRM
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/email/login_confirm",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//---------------------
				// INPUT
				httpInput, gfErr := gf_identity_core.HTTPgetEmailLoginConfirmInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				confirmedBool, failMsgStr, gfErr := gf_identity_core.UsersEmailLoginConfirmPipeline(httpInput,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				if confirmedBool {

					userNameStr := httpInput.UserNameStr



					// for non-admins email confirmation is only run initially on user creation
					// and if successfuly will login the user
					//---------------------
					// LOGIN_FINALIZE

					loginFinalizeInput := &gf_identity_core.GFuserpassInputLoginFinalize{
						UserNameStr: userNameStr,
					}
					loginFinalizeOutput, gfErr := gf_identity_core.UserpassPipelineLoginFinalize(loginFinalizeInput,
						pKeyServer,
						pServiceInfo,
						pCtx,
						pRuntimeSys)
					if gfErr != nil {
						return nil, gfErr
					}

					//---------------------
					// SET_COOKIES

					sameSiteStrictBool := true
					jwtTokenValStr := string(loginFinalizeOutput.JWTtokenVal)
					gf_identity_core.CreateAuthCookie(jwtTokenValStr,
						pServiceInfo.DomainForAuthCookiesStr,
						sameSiteStrictBool,
						pResp)

					//---------------------

					// now that user is logged in redirect them if a redirect URL was specified. 
					if pServiceInfo.AuthLoginSuccessRedirectURLstr != "" {
					
						http.Redirect(pResp,
							pReq,
							pServiceInfo.AuthLoginSuccessRedirectURLstr,
							301)
					}
					

				} else {
					outputMap := map[string]interface{}{
						"fail_msg_str": failMsgStr,
					}
					return outputMap, nil
				}

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// EMAIL_CONFIRM
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/email_confirm",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//---------------------
				// INPUT
				httpInput, gfErr := gf_identity_core.HTTPgetEmailConfirmInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				confirmedBool, failMsgStr, gfErr := gf_identity_core.UsersEmailPipelineConfirm(httpInput,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				if confirmedBool {

					userNameStr := httpInput.UserNameStr

					// for admins the login process has not completed yet after email confirmation
					if userNameStr == "admin" {

						// redirect user to login page
						// "email_confirmed=1" - signals to the UI that email has been confirmed
						URLredirectStr := fmt.Sprintf("%s?email_confirmed=1&user_name=%s",
							rpcHandlerRuntime.AuthLoginURLstr,
							userNameStr)

						// REDIRECT
						http.Redirect(pResp,
							pReq,
							URLredirectStr,
							301)
						
					} else {

						// for non-admins email confirmation is only run initially on user creation
						// and if successfuly will login the user
						//---------------------
						// LOGIN_FINALIZE

						loginFinalizeInput := &gf_identity_core.GFuserpassInputLoginFinalize{
							UserNameStr: userNameStr,
						}
						loginFinalizeOutput, gfErr := gf_identity_core.UserpassPipelineLoginFinalize(loginFinalizeInput,
							pKeyServer,
							pServiceInfo,
							pCtx,
							pRuntimeSys)
						if gfErr != nil {
							return nil, gfErr
						}

						//---------------------
						// SET_COOKIES
						jwtTokenValStr := string(loginFinalizeOutput.JWTtokenVal)
						sameSiteStrictBool := true
						gf_identity_core.CreateAuthCookie(jwtTokenValStr,
							pServiceInfo.DomainForAuthCookiesStr,
							sameSiteStrictBool,
							pResp)

						//---------------------

						// now that user is logged in redirect them if a redirect URL was specified. 
						if pServiceInfo.AuthLoginSuccessRedirectURLstr != "" {
						
							http.Redirect(pResp,
								pReq,
								pServiceInfo.AuthLoginSuccessRedirectURLstr,
								301)
						}
					}

				} else {
					outputMap := map[string]interface{}{
						"fail_msg_str": failMsgStr,
					}
					return outputMap, nil
				}

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// MFA_CONFIRM
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/mfa_confirm",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT

				inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var userNameStr gf_identity_core.GFuserName
				if inputUserNameStr, ok := inputMap["user_name_str"].(string); ok {
					userNameStr = gf_identity_core.GFuserName(inputUserNameStr)
				}

				var externHtopValueStr string
				if inputExternHtopValueStr, ok := inputMap["mfa_val_str"].(string); ok {
					externHtopValueStr = inputExternHtopValueStr
				}

				input := &GFuserAuthMFAinputConfirm{
					UserNameStr:        userNameStr,
					ExternHtopValueStr: externHtopValueStr,
					SecretKeyBase32str: pServiceInfo.AdminMFAsecretKeyBase32str,
				}
				
				//---------------------
				
				validBool, gfErr := mfaPipelineConfirm(input,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				if validBool {
					//---------------------
					// LOGIN_FINALIZE

					loginFinalizeInput := &gf_identity_core.GFuserpassInputLoginFinalize{
						UserNameStr: userNameStr,
					}
					loginFinalizeOutput, gfErr := gf_identity_core.UserpassPipelineLoginFinalize(loginFinalizeInput,
						pKeyServer,
						pServiceInfo,
						pCtx,
						pRuntimeSys)
					if gfErr != nil {
						return nil, gfErr
					}

					//---------------------	
					// SET_COOKIES
					sameSiteStrictBool := true
					jwtTokenValStr := string(loginFinalizeOutput.JWTtokenVal)
					gf_identity_core.CreateAuthCookie(jwtTokenValStr,
						pServiceInfo.DomainForAuthCookiesStr,
						sameSiteStrictBool,
						pResp)

					//---------------------
				}

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)

				outputMap := map[string]interface{}{
					"mfa_valid_bool": validBool,
				}
				return outputMap, nil
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// USERS_UPDATE
	// AUTH - only logged in users can update their own details

	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/update",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT

				userID, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				HTTPinput, gfErr := gf_identity_core.HTTPgetUserUpdateInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input := &gf_identity_core.GFuserUpdateInput{
					UserID:             userID,
					EmailStr:           HTTPinput.EmailStr,
					DescriptionStr:     HTTPinput.DescriptionStr,
					ProfileImageURLstr: HTTPinput.ProfileImageURLstr,
					BannerImageURLstr:  HTTPinput.BannerImageURLstr,
				}
				
				// VALIDATE
				gfErr = gf_core.ValidateStruct(input, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				//---------------------

				_, gfErr = gf_identity_core.UsersPipelineUpdate(input,
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
				gf_core.HTTPdisableCachingOfResponse(pResp)

				outputMap := map[string]interface{}{}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// REGISTER_INVITE_EMAIL
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/register_invite_email",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			// IMPORTANT!! - disable client caching for this endpoint, to avoid incosistent behavior
			gf_core.HTTPdisableCachingOfResponse(pResp)
			
			dataMap := map[string]interface{}{}
			return dataMap, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------


	return nil
}