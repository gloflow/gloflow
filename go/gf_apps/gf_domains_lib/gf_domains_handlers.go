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

package gf_domains_lib

import (
	"fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func Init_handlers(p_templates_paths_map map[string]string,
	p_mux         *http.ServeMux,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GF_error {

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------------

	//---------------------
	// DOMAIN_BROWSER
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/domains/browser",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {
				
				//--------------------
				//response_format_str - "json"|"html"

				qs_map := p_req.URL.Query()
				fmt.Println(qs_map)

				/*//response_format_str - "j"(for json)|"h"(for html)
				response_format_str := gf_rpc_lib.Get_response_format(qs_map,
																p_log_fun)*/
				//--------------------
				// GET DOMAINS FROM DB
				domains_lst, gf_err := db__get_domains(p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//--------------------
				// RENDER TEMPLATE
				gf_err = domains_browser__render_template(domains_lst,
					gf_templates.domains_browser__tmpl,
					gf_templates.domains_browser__subtemplates_names_lst,
					p_resp,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				return nil, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)
	
	//---------------------

	return nil
}