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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func init_handlers(p_mux *http.ServeMux,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/admin",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(handlers_endpoints_lst)

	//---------------------
	// ADMIN
	gf_rpc_lib.Create_handler__http_with_mux("/v1/admin/",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

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