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

func initHandlers(pTemplatesPathsMap map[string]string,
	pKeyServer           *gf_identity_core.GFkeyServerInfo,
	pHTTPmux             *http.ServeMux,
	pServiceInfo         *GFserviceInfo,
	pIdentityServiceInfo *gf_identity_core.GFserviceInfo,
	pLocalHub            *sentry.Hub,
	pRuntimeSys          *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES

	gfTemplates, gfErr := templatesLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// METRICS
	handlersEndpointsLst := []string{
		"/v1/admin/login",
		"/v1/admin/login_ui",
		"/v1/admin/dashboard",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_admin", handlersEndpointsLst)

	//---------------------
	// rpcHandlerRuntime
	rpcHandlerRuntime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:             pHTTPmux,
		Metrics:         metrics,
		StoreRunBool:    true,
		SentryHub:       pLocalHub,
		AuthSubsystemTypeStr: pServiceInfo.AuthSubsystemTypeStr,
		AuthLoginURLstr: "/v1/admin/login_ui",
		AuthKeyServer:   pKeyServer,
	}

	//---------------------
	// ADMIN_LOGIN_UI
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/admin/login_ui",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//---------------------
				validBool, _, gfErr := gf_identity_core.Validate(pReq,
					pKeyServer,
					rpcHandlerRuntime.AuthSubsystemTypeStr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				if validBool {
					// if the user is logged in and accesing login UI, then redirect to admin dashboard
					http.Redirect(pResp,
						pReq,
						"/v1/admin/dashboard",
						301)
				}

				//---------------------
				// INPUT
				qsMap := pReq.URL.Query()

				mfaConfirmBool := false

				// email_confirmed - signals that email has been confirmed. appended by the email_confirm handler
				//                   when the user gets redirected to admin login_ui URL.
				//                   in admin login this means that MFA_confirmation is next
				if _, ok := qsMap["email_confirmed"]; ok {
					mfaConfirmBool = true
				}

				//---------------------

				templateRenderedStr, gfErr := PipelineRenderLogin(mfaConfirmBool,
					gfTemplates.loginTmpl,
					gfTemplates.dashboardSubtemplatesNamesLst,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				pResp.Write([]byte(templateRenderedStr))
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// ADMIN_LOGIN
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/admin/login",
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

				input := &gf_identity.GFadminInputLogin{
					UserNameStr: userNameStr,
					PassStr:     passStr,
					EmailStr:    pServiceInfo.AdminEmailStr,
				}

				//---------------------

				output, gfErr := gf_identity.AdminPipelineLogin(input,
					pCtx,
					pLocalHub,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"user_exists_bool": output.UserExistsBool,
					"pass_valid_bool":  output.PassValidBool,
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// ADMIN_DASHBOARD
	// AUTH - only logged in admins can use the dashboard
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/dashboard",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				templateRenderedStr, gfErr := PipelineRenderDashboard(gfTemplates.dashboardTmpl,
					gfTemplates.dashboardSubtemplatesNamesLst,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				pResp.Write([]byte(templateRenderedStr))
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// HEALHTZ - admin has its own healthz because its started separately from other apps services
	//           on a different port, and registeres separate handlers (on a separate mux).

	rpcHandlerRuntimeHealth := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:          pHTTPmux,
		Metrics:      metrics,
		StoreRunBool: false,
		SentryHub:    pLocalHub,
	}

	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/admin/healthz",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			return nil, nil
		},
		rpcHandlerRuntimeHealth,
		pRuntimeSys)

	//---------------------
	
	return nil
}