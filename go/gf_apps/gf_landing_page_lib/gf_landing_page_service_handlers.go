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

package gf_landing_page_lib

import (
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//------------------------------------------------
func init_handlers(p_templates_paths_map map[string]string,
	p_http_mux    *http.ServeMux,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_landing_page_service_handlers.init_handlers()")

	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/landing/main/",
		"/landing/register_invite_email",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers("gf_landing_page", handlers_endpoints_lst)

	//---------------------
	// MAIN
	gf_rpc_lib.Create_handler__http_with_mux("/landing/main/",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {


				imgs__max_random_cursor_position_int  := 10000
				posts__max_random_cursor_position_int := 2000
				gf_err := Pipeline__render_landing_page(imgs__max_random_cursor_position_int,
					posts__max_random_cursor_position_int,
					5,  // p_featured_posts_to_get_int
					10, // p_featured_imgs_to_get_int
					gf_templates.tmpl,
					gf_templates.subtemplates_names_lst,
					p_resp,
					p_runtime_sys)

				if gf_err != nil {
					return nil, gf_err
				}
			}
			
			// IMPORTANT!! - this handler renders and writes template output to HTTP response, 
			//               and should not return any JSON data, so mark data_map as nil t prevent gf_rpc_lib
			//               from returning it.
			return nil, nil
		},
		p_http_mux,
		metrics,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)

	//---------------------
	// REGISTER_INVITE_EMAIL
	gf_rpc_lib.Create_handler__http_with_mux("/landing/register_invite_email",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			
			
			data_map := map[string]interface{}{}
			return data_map, nil
		},
		p_http_mux,
		metrics,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)

	//---------------------
	return nil
}