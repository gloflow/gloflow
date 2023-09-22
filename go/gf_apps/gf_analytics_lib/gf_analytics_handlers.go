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
)

//-------------------------------------------------

func initHandlers(pTemplatesPathsMap map[string]string,
	pMux        *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/a/ue",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_analytics", handlers_endpoints_lst)

	//---------------------
	// USER_EVENT
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/a/ue",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			// CORS - preflight request
			gf_rpc_lib.HTTPcorsPreflightHandle(pReq, pResp)

			if pReq.Method == "POST" {

				ip_str       := pReq.RemoteAddr
				clean_ip_str := strings.Split(ip_str,":")[0]
				
				//-----------------
				// BROWSER INFORMATION
				user_agent_str := pReq.UserAgent()
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
				input, sessionIDstr, gfErr := gf_events.User_event__parse_input(pReq, pResp, pRuntimeSys)
				if gfErr != nil {
					//IMPORTANT!! - this is a special case handler, we dont want it to return any standard JSON responses,
					//              this handler should be fire-and-forget from the users/clients perspective.
					return nil, gfErr
				}
				
				//-----------------
							
				gfReqCtx := &gf_events.GF_user_event_req_ctx {
					User_ip_str:      clean_ip_str,
					User_agent_str:   user_agent_str,
					Browser_name_str: browser_name_str,
					Browser_ver_str:  browser_ver_str,
					Os_name_str:      os_name_str,
					Os_ver_str:       os_version_str,
				}

				gfErr = gf_events.User_event__create(input, sessionIDstr, gfReqCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//-----------------
			}
			return nil, nil
		},
		pMux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//--------------
	return nil
}