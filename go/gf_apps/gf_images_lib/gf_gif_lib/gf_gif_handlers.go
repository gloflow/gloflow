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

package gf_gif_lib

import (
	"time"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func Gif__init_handlers(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_gif.Flows__init_handlers()")

	//-------------------------------------------------
	//GIF_GET_INFO
	http.HandleFunc("/images/gif/get_info", func(p_resp http.ResponseWriter, p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/gif/get_info ----------")
		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			//INPUT

			qs_map := p_req.URL.Query()

			var origin_url_str string
			if a_lst,ok := qs_map["orig_url"]; ok {
				origin_url_str = a_lst[0]
			}

			var gf_img_id_str string
			if a_lst,ok := qs_map["gfimg_id"]; ok {
				gf_img_id_str = a_lst[0]
			}
			//--------------------------
			var gf_gif *Gf_gif
			var gf_err *gf_core.Gf_error

			//BY_ORIGIN_URL
			if origin_url_str != "" {
				p_runtime_sys.Log_fun("INFO","origin_url_str - "+origin_url_str)

				gf_gif,gf_err = gif_db__get_by_origin_url(origin_url_str,p_runtime_sys)

				if gf_err != nil {
					usr_msg_str := "failed to get GIF from DB, with gf_url - "+origin_url_str
					gf_rpc_lib.Error__in_handler("/images/gif/get_info",usr_msg_str,gf_err,p_resp,p_runtime_sys)
					return
				}

			//BY_GF_IMG_ID
			} else if gf_img_id_str != "" {
				p_runtime_sys.Log_fun("INFO","gf_img_id_str - "+gf_img_id_str)

				gf_gif,gf_err = gif_db__get_by_img_id(gf_img_id_str,p_runtime_sys)

				if gf_err != nil {
					usr_msg_str := "failed to get GIF from DB, with gf_img_id - "+gf_img_id_str
					gf_rpc_lib.Error__in_handler("/images/gif/get_info",usr_msg_str,gf_err,p_resp,p_runtime_sys)
					return
				}
			}
			//------------------
			//OUTPUT
			data_map := map[string]interface{}{
				"gif_map":gf_gif,
			}
			gf_rpc_lib.Http_respond(data_map, "OK", p_resp,p_runtime_sys)
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/gif/get_info", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//-------------------------------------------------
	
	return nil
}