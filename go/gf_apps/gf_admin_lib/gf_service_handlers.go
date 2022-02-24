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
func init_handlers(p_templates_paths_map map[string]string,
	p_http_mux              *http.ServeMux,
	p_service_info          *GF_service_info,
	p_identity_service_info *gf_identity_lib.GF_service_info,
	p_local_hub             *sentry.Hub,
	p_runtime_sys           *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/admin/login",
		"/v1/admin/login_ui",
		"/v1/admin/dashboard",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers("gf_admin", handlers_endpoints_lst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpc_handler_runtime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                p_http_mux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         p_local_hub,
		Auth_login_url_str: "/v1/admin/login_ui",
	}

	//---------------------
	// ADMIN_LOGIN_UI
	// NO_AUTH
	gf_rpc_lib.Create_handler__http_with_auth(false, "/v1/admin/login_ui",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//---------------------
				// INPUT
				qs_map := p_req.URL.Query()

				mfa_confirm_bool := false

				// email_confirmed - signals that email has been confirmed. appended by the email_confirm handler
				//                   when the user gets redirected to admin login_ui URL.
				//                   in admin login this means that MFA_confirmation is next
				if _, ok := qs_map["email_confirmed"]; ok {
					mfa_confirm_bool = true
				}

				//---------------------

				template_rendered_str, gf_err := Pipeline__render_login(mfa_confirm_bool,
					gf_templates.login__tmpl,
					gf_templates.dashboard__subtemplates_names_lst,
					p_ctx,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				p_resp.Write([]byte(template_rendered_str))
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// ADMIN_LOGIN
	// NO_AUTH
	gf_rpc_lib.Create_handler__http_with_auth(false, "/v1/admin/login",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				var user_name_str string
				if val_str, ok := input_map["user_name_str"]; ok {
					user_name_str = val_str.(string)
				}

				var pass_str string
				if val_str, ok := input_map["pass_str"]; ok {
					pass_str = val_str.(string)
				}

				gf_err = gf_identity_lib.Admin__is(user_name_str, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input := &gf_identity_lib.GF_admin__input_login{
					User_name_str: user_name_str,
					Pass_str:      pass_str,
					Email_str:     p_service_info.Admin_email_str,
				}

				//---------------------

				output, gf_err := gf_identity_lib.Admin__pipeline__login(input,
					p_ctx,
					p_local_hub,
					p_identity_service_info,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"pass_valid_bool": output.Pass_valid_bool,
				}
				return output_map, nil
			}
			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// ADMIN_DASHBOARD
	// AUTH - only logged in admins can use the dashboard
	gf_rpc_lib.Create_handler__http_with_auth(true, "/v1/admin/dashboard",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				template_rendered_str, gf_err := Pipeline__render_dashboard(gf_templates.dashboard__tmpl,
					gf_templates.dashboard__subtemplates_names_lst,
					p_ctx,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				p_resp.Write([]byte(template_rendered_str))
			}

			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// HEALHTZ - admin has its own healthz because its started separately from other apps services
	//           on a different port, and registeres separate handlers (on a separate mux).

	rpc_handler_runtime_health := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:            p_http_mux,
		Metrics:        metrics,
		Store_run_bool: false,
		Sentry_hub:     p_local_hub,
	}

	gf_rpc_lib.Create_handler__http_with_auth(false, "/v1/admin/healthz",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			return nil, nil
		},
		rpc_handler_runtime_health,
		p_runtime_sys)

	//---------------------
	
	return nil
}

//-------------------------------------------------
func tmpl__load(p_templates_paths_map map[string]string,
	p_runtime_sys *gf_core.Runtime_sys) (*gf_templates, *gf_core.Gf_error) {

	login_template_filepath_str     := p_templates_paths_map["gf_admin_login"]
	dashboard_template_filepath_str := p_templates_paths_map["gf_admin_dashboard"]

	l_tmpl, l_subtemplates_names_lst, gf_err := gf_core.Templates__load(login_template_filepath_str,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	d_tmpl, d_subtemplates_names_lst, gf_err := gf_core.Templates__load(dashboard_template_filepath_str,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	gf_templates := &gf_templates{
		login__tmpl:                   l_tmpl,
		login__subtemplates_names_lst: l_subtemplates_names_lst,
		dashboard__tmpl:                   d_tmpl,
		dashboard__subtemplates_names_lst: d_subtemplates_names_lst,
	}
	return gf_templates, nil
}