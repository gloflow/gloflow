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

func initHandlersEth(pKeyServer *gf_identity_core.GFkeyServerInfo,
	pHTTPmux     *http.ServeMux,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/identity/eth/preflight",
		"/v1/identity/eth/login",
		"/v1/identity/eth/create",
	}
	metricsGroupNameStr := "eth"
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
	// USERS_PREFLIGHT
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/eth/preflight",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				_, _, userAddressETHstr, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input :=&gf_identity_core.GFethInputPreflight{
					UserAddressETHstr: userAddressETHstr,
				}

				//---------------------

				output, gfErr := gf_identity_core.ETHpipelinePreflight(input, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"user_exists_bool": output.UserExistsBool,
					"nonce_val_str":    output.NonceValStr,
				}
				return outputMap, nil
			}

			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// USERS_LOGIN
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/eth/login",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				inputMap, _, userAddressETHstr, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				authSignatureStr := gf_identity_core.GFauthSignature(inputMap["auth_signature_str"].(string))

				input :=&gf_identity_core.GFethInputLogin{
					UserAddressETHstr: userAddressETHstr,
					AuthSignatureStr:  authSignatureStr,
				}

				//---------------------
				// LOGIN
				output, gfErr := gf_identity_core.ETHpipelineLogin(input,
					pKeyServer,
					pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------
				// SET_SESSION_ID - sets gf_sess cookie on all future requests
				jwtTokenValStr := string(output.JWTtokenVal)
				gf_session.Create(jwtTokenValStr, pResp)

				//---------------------

				outputMap := map[string]interface{}{
					"auth_signature_valid_bool": output.AuthSignatureValidBool,
					"nonce_exists_bool":         output.NonceExistsBool,
					"user_id_str":               output.UserIDstr,
				}
				return outputMap, nil
			}


			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// USERS_CREATE
	// NO_AUTH - unauthenticated users are able to create new users
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/eth/create",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input :=&gf_identity_core.GFethInputCreate{
					UserTypeStr:       "standard",
					UserAddressETHstr: gf_identity_core.GFuserAddressETH(inputMap["user_address_eth_str"].(string)),
					AuthSignatureStr:  gf_identity_core.GFauthSignature(inputMap["auth_signature_str"].(string)),
				}
				
				//---------------------
				output, gfErr := gf_identity_core.ETHpipelineCreate(input, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"auth_signature_valid_bool": output.AuthSignatureValidBool,
					"nonce_exists_bool":         output.NonceExistsBool,
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