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
	"context"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func Gif__init_handlers(p_mux *http.ServeMux,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.Gf_error {
	p_runtime_sys.LogFun("FUN_ENTER","gf_gif.Flows__init_handlers()")

	//-------------------------------------------------
	// GIF_GET_INFO
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/gif/get_info",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {
		
			if p_req.Method == "GET" {

				//--------------------------
				// INPUT

				qs_map := p_req.URL.Query()

				var origin_url_str string
				if a_lst, ok := qs_map["orig_url"]; ok {
					origin_url_str = a_lst[0]
				}

				var gf_img_id_str string
				if a_lst, ok := qs_map["gfimg_id"]; ok {
					gf_img_id_str = a_lst[0]
				}
				
				//--------------------------
				var gfGIF *GFgif
				var gfErr *gf_core.GFerror

				// BY_ORIGIN_URL
				if origin_url_str != "" {
					p_runtime_sys.LogFun("INFO","origin_url_str - "+origin_url_str)

					gfGIF, gfErr = gif_db__get_by_origin_url(origin_url_str, p_runtime_sys)

					if gfErr != nil {
						return nil, gfErr
					}

				// BY_GF_IMG_ID
				} else if gf_img_id_str != "" {
					p_runtime_sys.LogFun("INFO","gf_img_id_str - "+gf_img_id_str)

					gfGIF, gfErr = gif_db__get_by_img_id(gf_img_id_str,p_runtime_sys)

					if gfErr != nil {
						return nil, gfErr
					}
				}

				//------------------
				// OUTPUT
				output_map := map[string]interface{}{
					"gif_map":gfGIF,
				}
				return output_map, nil
				
				//------------------
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)
	
	//-------------------------------------------------
	
	return nil
}