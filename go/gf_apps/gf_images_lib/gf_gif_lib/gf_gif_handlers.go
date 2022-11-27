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

func InitHandlers(pMux *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//-------------------------------------------------
	// GIF_GET_INFO
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/gif/get_info",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {
		
			if p_req.Method == "GET" {

				//--------------------------
				// INPUT

				qs_map := p_req.URL.Query()

				var originURLstr string
				if a_lst, ok := qs_map["orig_url"]; ok {
					originURLstr = a_lst[0]
				}

				var imgIDstr string
				if a_lst, ok := qs_map["gfimg_id"]; ok {
					imgIDstr = a_lst[0]
				}
				
				//--------------------------
				var gif *GFgif
				var gfErr *gf_core.GFerror

				// BY_ORIGIN_URL
				if originURLstr != "" {

					gif, gfErr = dbGetByOriginURL(originURLstr, pRuntimeSys)

					if gfErr != nil {
						return nil, gfErr
					}

				// BY_GF_IMG_ID
				} else if imgIDstr != "" {

					gif, gfErr = gifDBgetByImgID(imgIDstr,pRuntimeSys)

					if gfErr != nil {
						return nil, gfErr
					}
				}

				//------------------
				// OUTPUT
				outputMap := map[string]interface{}{
					"gif_map": gif,
				}
				return outputMap, nil
				
				//------------------
			}
			return nil, nil
		},
		pMux,
		nil,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)
	
	//-------------------------------------------------
	
	return nil
}