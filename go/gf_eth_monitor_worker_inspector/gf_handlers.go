/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package main

import (
	"net/http"
	// "github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func init_handlers(p_metrics *GF_metrics,
	p_runtime *GF_runtime) {

	//---------------------
	// GET_BLOCKS
	http.HandleFunc("/gfethm_worker_inspect/v1/blocks", func(p_resp http.ResponseWriter, p_req *http.Request) {
		
		// METRICS
		if p_metrics != nil {
			p_metrics.counter__http_req_num__get_blocks.Inc()
		}



		//------------------
		// OUTPUT
		data_map := map[string]interface{}{
			
		}
		gf_rpc_lib.Http_respond(data_map, "OK", p_resp, p_runtime.Runtime_sys)

		//------------------



	})

	//---------------------
	// HEALTH
	http.HandleFunc("/gfethm_worker_inspect/v1/health", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_resp.Write([]byte("ok"))

	})

	//---------------------
}