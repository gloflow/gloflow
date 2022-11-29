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

func InitHandlers(p_templates_paths_map map[string]string,
	p_mux         *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES

	gfTemplates, gfErr := tmplLoad(p_templates_paths_map, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------

	//---------------------
	// DOMAIN_BROWSER
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/domains/browser",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if p_req.Method == "GET" {
				
				//--------------------
				//response_format_str - "json"|"html"

				qs_map := p_req.URL.Query()
				fmt.Println(qs_map)

				/*//response_format_str - "j"(for json)|"h"(for html)
				response_format_str := gf_rpc_lib.Get_response_format(qs_map,
																pLogFun)*/
				//--------------------
				// GET DOMAINS FROM DB
				domains_lst, gfErr := dbGetDomains(pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//--------------------
				// RENDER TEMPLATE
				gfErr = domainsBrowserRenderTemplate(domains_lst,
					gfTemplates.domains_browser__tmpl,
					gfTemplates.domains_browser__subtemplates_names_lst,
					p_resp,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				return nil, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//---------------------

	return nil
}