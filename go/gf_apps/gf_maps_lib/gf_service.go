/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

package gf_maps_lib

import (
	"fmt"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

func InitService(pTemplates_paths_map map[string]string,
	p_http_mux    *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) {

	//------------------------
	// STATIC FILES SERVING
	static_files__url_base_str := "/v1/maps"
	local_dir_path_str         := "./static"
	gf_core.HTTPinitStaticServingWithMux(static_files__url_base_str,
		local_dir_path_str,
		p_http_mux,
		pRuntimeSys)

	//------------------------
	// HANDLERS
	gf_err := init_handlers(pTemplates_paths_map, p_http_mux, pRuntimeSys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//------------------------
}

//-------------------------------------------------