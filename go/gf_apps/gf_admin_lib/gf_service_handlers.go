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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
type gf_templates struct {
	dashboard__tmpl                   *template.Template
	dashboard__subtemplates_names_lst []string
}

//------------------------------------------------
func init_handlers(p_templates_paths_map map[string]string,
	p_mux         *http.ServeMux,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/admin",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(handlers_endpoints_lst)

	//---------------------
	// ADMIN
	gf_rpc_lib.Create_handler__http_with_mux("/v1/admin",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				template_rendered_str, gf_err := Pipeline__render(gf_templates.dashboard__tmpl,
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
		p_mux,
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	// HEALHTZ - admin has its own healthz because its started separately from other apps services
	//           on a different port, and registeres separate handlers (on a separate mux).

	gf_rpc_lib.Create_handler__http_with_mux("/v1/admin/healthz",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			return nil, nil
		},
		p_mux,
		nil,
		false, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	
	return nil
}

//-------------------------------------------------
func tmpl__load(p_templates_paths_map map[string]string,
	p_runtime_sys *gf_core.Runtime_sys) (*gf_templates, *gf_core.Gf_error) {

	main_template_filepath_str := p_templates_paths_map["gf_admin_dashboard"]

	tmpl, subtemplates_names_lst, gf_err := gf_core.Templates__load(main_template_filepath_str,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	gf_templates := &gf_templates{
		dashboard__tmpl:                   tmpl,
		dashboard__subtemplates_names_lst: subtemplates_names_lst,
	}
	return gf_templates, nil
}