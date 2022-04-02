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

package gf_identity_lib

import (
	"fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_policy"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func initHandlers(p_auth_login_url_str string,
	pHTTPmux     *http.ServeMux,
	pServiceInfo *GF_service_info,
	pRuntimeSys  *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/identity/policy/update",
		"/v1/identity/email_confirm",
		"/v1/identity/mfa_confirm",
		"/v1/identity/update",
		"/v1/identity/me",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(pServiceInfo.Name_str, handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                pHTTPmux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: p_auth_login_url_str,
	}

	//---------------------
	// POLICY
	//---------------------
	// POLICY_UPDATE
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/policy/update",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT

				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				inputMap, gfErr := gf_rpc_lib.Get_http_input(pResp, pReq, pRuntimeSys)
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
	// EMAIL_CONFIRM
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/email_confirm",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//---------------------
				// INPUT
				http_input, gf_err := gf_identity_core.Http__get_email_confirm_input(p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------

				confirmedBool, fail_msg_str, gf_err := usersEmailPipelineConfirm(http_input,
					p_ctx,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				if confirmedBool {

					// redirect user to login page
					// "email_confirmed=1" - signals to the UI that email has been confirmed
					url_redirect_str := fmt.Sprintf("%s?email_confirmed=1&user_name=%s",
						rpcHandlerRuntime.Auth_login_url_str,
						http_input.User_name_str)

					// REDIRECT
					http.Redirect(p_resp,
						p_req,
						url_redirect_str,
						301)

				} else {
					output_map := map[string]interface{}{
						"fail_msg_str": fail_msg_str,
					}
					return output_map, nil
				}
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// MFA_CONFIRM
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/mfa_confirm",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT

				inputMap, gfErr := gf_rpc_lib.Get_http_input(pResp, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var userNameStr gf_identity_core.GFuserName
				if inputUserNameStr, ok := inputMap["user_name_str"].(string); ok {
					userNameStr = gf_identity_core.GFuserName(inputUserNameStr)
				}

				var extern_htop_value_str string
				if input_extern_htop_value_str, ok := inputMap["mfa_val_str"].(string); ok {
					extern_htop_value_str = input_extern_htop_value_str
				}

				input := &GF_user_auth_mfa__input_confirm{
					UserNameStr:           userNameStr,
					Extern_htop_value_str: extern_htop_value_str,
					Secret_key_base32_str: pServiceInfo.Admin_mfa_secret_key_base32_str,
				}
				
				//---------------------
				
				validBool, gf_err := mfaPipelineConfirm(input,
					pCtx,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				if validBool {
					//---------------------
					// LOGIN_FINALIZE

					loginFinalizeInput := &GF_user_auth_userpass__input_login_finalize{
						UserNameStr: userNameStr,
					}
					loginFinalizeOutput, gf_err := users_auth_userpass__pipeline__login_finalize(loginFinalizeInput,
						pServiceInfo,
						pCtx,
						pRuntimeSys)
					if gf_err != nil {
						return nil, gf_err
					}

					//---------------------					
					// SET_SESSION_ID - sets gf_sid cookie on all future requests
					sessionDataStr     := string(loginFinalizeOutput.JWT_token_val)
					sessionTTLhoursInt := 24 // 1 day
					gf_session.SetOnReq(sessionDataStr, pResp, sessionTTLhoursInt)

					//---------------------
				}

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
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT

				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				HTTPinput, gfErr := gf_identity_core.Http__get_user_update_input(p_req, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input := &GF_user__input_update{
					UserIDstr:             userIDstr,
					Email_str:             HTTPinput.Email_str,
					Description_str:       HTTPinput.Description_str,
					Profile_image_url_str: HTTPinput.Profile_image_url_str,
					Banner_image_url_str:  HTTPinput.Banner_image_url_str,
				}
				
				// VALIDATE
				gfErr = gf_core.Validate_struct(input, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				//---------------------

				_, gfErr = users__pipeline__update(input,
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// USERS_GET_ME
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/me",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//---------------------
				// INPUT

				userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

				input := &GF_user__input_get{
					UserIDstr: userIDstr,
				}

				//---------------------

				output, gf_err := usersPipelineGet(input, pCtx, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"user_name_str":         output.User_name_str,
					"email_str":             output.Email_str,
					"description_str":       output.Description_str,
					"profile_image_url_str": output.Profile_image_url_str,
					"banner_image_url_str":  output.Banner_image_url_str,
				}
				return output_map, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------

	return nil
}