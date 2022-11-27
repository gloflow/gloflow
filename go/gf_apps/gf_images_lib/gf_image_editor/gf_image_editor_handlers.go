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

func InitHandlers(pMux *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) {

	//---------------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/editor/save",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {
		
			if p_req.Method == "POST" {

				//-------------------
				gfErr := saveEditedImagePipeline("/images/editor/save", p_req, p_resp,
					pCtx,
					pRuntimeSys)

				if gfErr != nil {
					return nil, gfErr
				}
				
				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		pMux,
		nil,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)
	
	//---------------------
}