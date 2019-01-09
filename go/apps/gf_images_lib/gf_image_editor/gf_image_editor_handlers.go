/*
GloFlow media management/publishing system
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
	"time"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func Init_handlers(p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_image_editor_handlers.Init_handlers()")

	//---------------------
	http.HandleFunc("/images/editor/save",func(p_resp http.ResponseWriter,
											p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/editor/save ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0


			//-------------------

			gf_err := save_edited_image__pipeline("/images/editor/save",p_req,p_resp,p_runtime_sys)


			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/editor/save",
								"failed to save modified image", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}
	
 			//------------------
			//OUTPUT
			data_map := map[string]interface{}{}
			gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/editor/save",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//---------------------
}