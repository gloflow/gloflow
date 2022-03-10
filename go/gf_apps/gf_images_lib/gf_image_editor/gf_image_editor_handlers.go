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

package gf_image_editor

import (
	"context"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func Init_handlers(p_mux *http.ServeMux,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_image_editor_handlers.Init_handlers()")

	//---------------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/editor/save",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {
		
			if p_req.Method == "POST" {

				//-------------------
				gf_err := save_edited_image__pipeline("/images/editor/save", p_req, p_resp, p_req.Context(), p_runtime_sys)

				if gf_err != nil {
					return nil, gf_err
				}
				
				//------------------
				// OUTPUT
				data_map := map[string]interface{}{}
				return data_map, nil

				//------------------
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)
	
	//---------------------
}