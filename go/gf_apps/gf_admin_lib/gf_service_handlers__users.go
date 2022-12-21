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

package gf_admin_lib

import (
	// "fmt"
	"net/http"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------

func initHandlersUsers(pKeyServer *gf_identity_core.GFkeyServerInfo,
	pHTTPmux             *http.ServeMux,
	pServiceInfo         *GFserviceInfo,
	pIdentityServiceInfo *gf_identity_core.GFserviceInfo,
	pLocalHub            *sentry.Hub,
	pRuntimeSys          *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/admin/users/delete",
		"/v1/admin/users/get_all",
		"/v1/admin/users/get_all_invite_list",
		"/v1/admin/users/add_to_invite_list",
		"/v1/admin/users/resend_confirm_email",
	}
	metricsGroupNameStr := "users"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_admin", handlersEndpointsLst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:             pHTTPmux,
		Metrics:         metrics,
		StoreRunBool:    true,
		SentryHub:       pLocalHub,
		AuthLoginURLstr: "/v1/admin/login_ui",
		AuthKeyServer:   pKeyServer,
	}

	//---------------------
	// DELETE
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/delete",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				inputMap, adminUserIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var userIDstr string
				if valStr, ok := inputMap["user_id_str"]; ok {
					userIDstr = valStr.(string)
				}

				var userNameStr string
				if valStr, ok := inputMap["user_name_str"]; ok {
					userNameStr = valStr.(string)
				}

				gfErr = gf_identity.AdminIs(adminUserIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				//---------------------

				input := &gf_identity.GFadminUserDeleteInput{
					UserIDstr:   gf_core.GF_ID(userIDstr),
					UserNameStr: gf_identity_core.GFuserName(userNameStr),
				}
				gfErr = gf_identity.AdminPipelineDeleteUser(input,
					pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// GET_ALL
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/get_all",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				_, adminUserIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				gfErr = gf_identity.AdminIs(adminUserIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				usersLst, gfErr := gf_identity.AdminPipelineGetAllUsers(pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"users_lst": usersLst,
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// GET_ALL_INVITE_LIST
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/get_all_invite_list",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				_, adminUserIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				gfErr = gf_identity.AdminIs(adminUserIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				inviteListLst, gfErr := gf_identity.AdminPipelineGetAllInviteList(pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"invite_list_lst": inviteListLst,
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// ADD_TO_INVITE_LIST
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/add_to_invite_list",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				inputMap, adminUserIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var emailStr string
				if valStr, ok := inputMap["email_str"]; ok {
					emailStr = valStr.(string)
				}

				gfErr = gf_identity.AdminIs(adminUserIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input := &gf_identity.GFadminInputAddToInviteList{
					AdminUserIDstr: adminUserIDstr,
					EmailStr:       emailStr,
				}

				//---------------------

				gfErr = gf_identity.AdminPipelineUserAddToInviteList(input,
					pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// REMOVE_FROM_INVITE_LIST
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/remove_from_invite_list",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				inputMap, adminUserIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var emailStr string
				if valStr, ok := inputMap["email_str"]; ok {
					emailStr = valStr.(string)
				}

				gfErr = gf_identity.AdminIs(adminUserIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input := &gf_identity.GFadminRemoveFromInviteListInput{
					AdminUserIDstr: adminUserIDstr,
					EmailStr:       emailStr,
				}

				//---------------------

				gfErr = gf_identity.AdminPipelineUserRemoveFromInviteList(input,
					pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// RESEND_CONFIRM_EMAIL
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/resend_confirm_email",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				inputMap, userIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var targetUserIDstr gf_core.GF_ID
				if valStr, ok := inputMap["user_id_str"]; ok {
					targetUserIDstr = gf_core.GF_ID(valStr.(string))
				}
				
				var targetUserNameStr gf_identity_core.GFuserName
				if valStr, ok := inputMap["user_name_str"]; ok {
					targetUserNameStr = gf_identity_core.GFuserName(valStr.(string))
				}

				var emailStr string
				if valStr, ok := inputMap["email_str"]; ok {
					emailStr = valStr.(string)
				}

				gfErr = gf_identity.AdminIs(userIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input := &gf_identity.GFadminResendConfirmEmailInput{
					UserIDstr:   targetUserIDstr,
					UserNameStr: targetUserNameStr,
					EmailStr:    emailStr,
				}

				//---------------------

				gfErr = gf_identity.AdminPipelineUserResendConfirmEmail(input,
					pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					
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