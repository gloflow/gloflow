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
	// "fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------

func initHandlersUserpass(pKeyServer *gf_identity_core.GFkeyServerInfo,
	pHTTPmux     *http.ServeMux,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/identity/userpass/login",
		"/v1/identity/userpass/create",
	}
	metricsGroupNameStr := "userpass"
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
		AuthKeyServer:   pKeyServer,
	}

	//---------------------
	// USERS_LOGIN
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/userpass/login",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
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

				input :=&gf_identity_core.GFuserpassInputLogin{
					UserNameStr: userNameStr,
					PassStr:     passStr,
					EmailStr:    emailStr,
				}

				//---------------------
				// LOGIN
				output, gfErr := gf_identity_core.UserpassPipelineLogin(input, 
					pKeyServer,
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------
				// COOKIES 

				if output.UserExistsBool {

					// SESSION_ID - sets gf_sess cookie
					sessionIDstr := string(output.SessionID)
					gf_session.CreateSessionIDcookie(sessionIDstr, pResp)

					


					// JWT - sets "Authorization" cookie
					jwtTokenValStr := string(output.JWTtokenVal)
					gf_session.CreateAuthCookie(jwtTokenValStr, pResp)
				}

				//---------------------

				outputMap := map[string]interface{}{
					"user_exists_bool": output.UserExistsBool,
					"pass_valid_bool":  output.PassValidBool,
					"user_id_str":      output.UserID,
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
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input :=&gf_identity_core.GFuserpassInputCreate{
					UserNameStr: gf_identity_core.GFuserName(inputMap["user_name_str"].(string)),
					PassStr:     inputMap["pass_str"].(string),
					EmailStr:    inputMap["email_str"].(string),
					UserTypeStr: "standard",
				}

				//---------------------
				output, gfErr := gf_identity_core.UserpassPipelineCreateRegular(input,
					pServiceInfo,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"user_exists_bool":         output.UserExistsBool,
					"user_in_invite_list_bool": output.UserInInviteListBool,
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