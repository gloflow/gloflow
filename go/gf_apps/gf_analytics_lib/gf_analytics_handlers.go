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
	"time"
	"strings"
	"net/http"
	"github.com/ianoshen/uaparser"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func init_handlers(p_templates_paths_map map[string]string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_analytics_handlers.init_handlers()")

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//--------------
	// USER_EVENT
	http.HandleFunc("/a/ue", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST --- /a/ue")

		// CORS - preflight request
		gf_rpc_lib.Http_CORS_preflight_handle(p_req, p_resp)
		// if p_req.Method == "OPTIONS" {
		// 	p_resp.Header().Set("Access-Control-Allow-Origin", "*")
		// 	p_resp.Header().Set("Access-Control-Allow-Origin", "Origin, X-Requested-With, Content-Type, Accept")
		// }
		
		

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano()) / 1000000000.0

			ip_str       := p_req.RemoteAddr
			clean_ip_str := strings.Split(ip_str,":")[0]

			cookies_lst := p_req.Cookies()
			cookies_str := gf_core.HTTP__serialize_cookies(cookies_lst,p_runtime_sys)
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
			input, session_id_str, gf_err := user_event__parse_input(p_req, p_resp, p_runtime_sys)
			if gf_err != nil {
				//IMPORTANT!! - this is a special case handler, we dont want it to return any standard JSON responses,
				//              this handler should be fire-and-forget from the users/clients perspective.
				return
			}
			
			//-----------------
						
			gf_req_ctx := &Gf_user_event_req_ctx {
				User_ip_str:      clean_ip_str,
				User_agent_str:   user_agent_str,
				Browser_name_str: browser_name_str,
				Browser_ver_str:  browser_ver_str,
				Os_name_str:      os_name_str,
				Os_ver_str:       os_version_str,
				Cookies_str:      cookies_str,
			}

			gf_err = user_event__create(input, session_id_str, gf_req_ctx, p_runtime_sys)
			if gf_err != nil {
				// IMPORTANT!! - this is a special case handler, we dont want it to return any standard JSON responses,
				//               this handler should be fire-and-forget from the users/clients perspective.
				return
			}
			//-----------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
		
			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/a/ue", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//--------------
	http.HandleFunc("/a/analytics_dashboard__ff0099__ooo", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST - /a/analytics_dashboard__ff0099__ooo ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------
			// RENDER TEMPLATE
			gf_err := dashboard__render_template(gf_templates.dashboard__tmpl,
				gf_templates.dashboard__subtemplates_names_lst,
				p_resp,
				p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/analytics_dashboard__ff0099__ooo", "failed to render analytics dashboard page", gf_err, p_resp, p_runtime_sys)
				return
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/a/analytics_dashboard__ff0099__ooo", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	//--------------
	return nil
}