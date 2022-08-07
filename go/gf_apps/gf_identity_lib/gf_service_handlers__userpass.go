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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func initHandlersUserpass(pHTTPmux *http.ServeMux,
	pServiceInfo *GF_service_info,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/identity/userpass/login",
		"/v1/identity/userpass/create",
	}
	metricsGroupNameStr := "userpass"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, pServiceInfo.Name_str, handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                pHTTPmux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: "/landing/main",
	}

	//---------------------
	// USERS_LOGIN
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/userpass/login",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				inputMap, gfErr := gf_core.HTTPgetInput(pResp, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var userNameStr gf_identity_core.GFuserName
				if valStr, ok := inputMap["user_name_str"]; ok {
					userNameStr = gf_identity_core.GFuserName(valStr.(string))
				}

				var passStr string
				if valStr, ok := inputMap["pass_str"]; ok {
					passStr = valStr.(string)
				}

				var emailStr string
				if valStr, ok := inputMap["email_str"]; ok {
					emailStr = valStr.(string)
				}

				input :=&GF_user_auth_userpass__input_login{
					User_name_str: userNameStr,
					Pass_str:      passStr,
					Email_str:     emailStr,
				}

				//---------------------
				// LOGIN
				output, gfErr := usersAuthUserpassPipelineLogin(input, 
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------
				// SET_SESSION_ID - sets gf_sid cookie on all future requests
				sessionDataStr        := string(output.JWT_token_val)
				sessionTTLhoursInt, _ := gf_identity_core.GetSessionTTL()
				gf_session.SetOnReq(sessionDataStr, pResp, sessionTTLhoursInt)

				//---------------------

				outputMap := map[string]interface{}{
					"user_exists_bool": output.User_exists_bool,
					"pass_valid_bool":  output.Pass_valid_bool,
					"user_id_str":      output.User_id_str,
				}
				return outputMap, nil
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// USERS_CREATE
	// NO_AUTH - unauthenticated users are able to create new users, and do not get logged in automatically on success

	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/userpass/create",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				inputMap, gfErr := gf_core.HTTPgetInput(pResp, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input :=&GF_user_auth_userpass__input_create{
					User_name_str: gf_identity_core.GFuserName(inputMap["user_name_str"].(string)),
					Pass_str:      inputMap["pass_str"].(string),
					Email_str:     inputMap["email_str"].(string),
					UserTypeStr:   "standard",
				}

				

				//---------------------
				output, gfErr := users_auth_userpass__pipeline__create_regular(input,
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"user_exists_bool":         output.User_exists_bool,
					"user_in_invite_list_bool": output.User_in_invite_list_bool,
				}
				return outputMap, nil
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	return nil
}