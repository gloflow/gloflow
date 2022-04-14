/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_analytics_lib

import (
	"strings"
	"context"
	"net/http"
	"github.com/ianoshen/uaparser"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
)

//-------------------------------------------------
func init_handlers(p_templates_paths_map map[string]string,
	p_mux         *http.ServeMux,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_analytics_handlers.init_handlers()")

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/a/ue",
		"/v1/a/dashboard",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_analytics", handlers_endpoints_lst)

	//---------------------
	// USER_EVENT
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/a/ue",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {


			// CORS - preflight request
			gf_rpc_lib.Http_CORS_preflight_handle(p_req, p_resp)
			// if p_req.Method == "OPTIONS" {
			// 	p_resp.Header().Set("Access-Control-Allow-Origin", "*")
			// 	p_resp.Header().Set("Access-Control-Allow-Origin", "Origin, X-Requested-With, Content-Type, Accept")
			// }

			if p_req.Method == "POST" {

				ip_str       := p_req.RemoteAddr
				clean_ip_str := strings.Split(ip_str,":")[0]
				
				//-----------------
				// BROWSER INFORMATION
				user_agent_str := p_req.UserAgent()
				user_agent     := uaparser.Parse(user_agent_str)

				var browser_name_str string
				var browser_ver_str  string
				if user_agent.Browser != nil {
					browser_name_str = user_agent.Browser.Name
					browser_ver_str  = user_agent.Browser.Version
				}

				os_name_str    := user_agent.OS.Name
				os_version_str := user_agent.OS.Version

				//-----------------
				// INPUT
				input, session_id_str, gf_err := gf_events.User_event__parse_input(p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					//IMPORTANT!! - this is a special case handler, we dont want it to return any standard JSON responses,
					//              this handler should be fire-and-forget from the users/clients perspective.
					return nil, gf_err
				}
				
				//-----------------
							
				gf_req_ctx := &gf_events.GF_user_event_req_ctx {
					User_ip_str:      clean_ip_str,
					User_agent_str:   user_agent_str,
					Browser_name_str: browser_name_str,
					Browser_ver_str:  browser_ver_str,
					Os_name_str:      os_name_str,
					Os_ver_str:       os_version_str,
				}

				gf_err = gf_events.User_event__create(input, session_id_str, gf_req_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//-----------------
			}
			return nil, nil
		},
		p_mux,
		metrics,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)

	//--------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/a/dashboard",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			
		if p_req.Method == "GET" {

			//---------------------
			// SESSION_VALIDATE
			valid_bool, _, gf_err := gf_session.Validate(p_req, p_ctx, p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			if !valid_bool {
				return nil, nil
			}

			//---------------------
			
			//--------------------
			// RENDER TEMPLATE
			gf_err = dashboard__render_template(gf_templates.dashboard__tmpl,
				gf_templates.dashboard__subtemplates_names_lst,
				p_resp,
				p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}
		}
		return nil, nil
	},
	p_mux,
	metrics,
	true, // p_store_run_bool
	nil,
	p_runtime_sys)

	//--------------
	return nil
}