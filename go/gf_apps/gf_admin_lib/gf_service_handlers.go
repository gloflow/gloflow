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
	"text/template"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
type gf_templates struct {
	login__tmpl                   *template.Template
	login__subtemplates_names_lst []string
	dashboard__tmpl                   *template.Template
	dashboard__subtemplates_names_lst []string
}

//------------------------------------------------
func initHandlers(pTemplatesPathsMap map[string]string,
	p_http_mux              *http.ServeMux,
	p_service_info          *GFserviceInfo,
	p_identity_service_info *gf_identity_lib.GFserviceInfo,
	p_local_hub             *sentry.Hub,
	pRuntimeSys             *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := templatesLoad(pTemplatesPathsMap, pRuntimeSys)
	if gf_err != nil {
		return gf_err
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
		Mux:                p_http_mux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         p_local_hub,
		Auth_login_url_str: "/v1/admin/login_ui",
	}

	//---------------------
	// ADMIN_LOGIN_UI
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/admin/login_ui",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//---------------------
				validBool, _, gfErr := gf_session.Validate(pReq, pCtx, pRuntimeSys)
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

				templateRenderedStr, gf_err := Pipeline__render_login(mfaConfirmBool,
					gf_templates.login__tmpl,
					gf_templates.dashboard__subtemplates_names_lst,
					pCtx,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
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
				
				inputMap, gf_err := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				var userNameStr gf_identity_core.GFuserName
				if val_str, ok := inputMap["user_name_str"]; ok {
					userNameStr = gf_identity_core.GFuserName(val_str.(string))
				}

				var passStr string
				if val_str, ok := inputMap["pass_str"]; ok {
					passStr = val_str.(string)
				}

				input := &gf_identity_lib.GF_admin__input_login{
					User_name_str: userNameStr,
					Pass_str:      passStr,
					Email_str:     p_service_info.Admin_email_str,
				}

				//---------------------

				output, gf_err := gf_identity_lib.Admin__pipeline__login(input,
					pCtx,
					p_local_hub,
					p_identity_service_info,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"user_exists_bool": output.User_exists_bool,
					"pass_valid_bool":  output.Pass_valid_bool,
				}
				return output_map, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// ADMIN_DASHBOARD
	// AUTH - only logged in admins can use the dashboard
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/dashboard",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if p_req.Method == "GET" {

				templateRenderedStr, gfErr := Pipeline__render_dashboard(gf_templates.dashboard__tmpl,
					gf_templates.dashboard__subtemplates_names_lst,
					p_ctx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				p_resp.Write([]byte(templateRenderedStr))
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
		Mux:            p_http_mux,
		Metrics:        metrics,
		Store_run_bool: false,
		Sentry_hub:     p_local_hub,
	}

	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/admin/healthz",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			return nil, nil
		},
		rpcHandlerRuntimeHealth,
		pRuntimeSys)

	//---------------------
	
	return nil
}