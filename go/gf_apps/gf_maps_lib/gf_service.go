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
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//-------------------------------------------------

func InitService(pAuthSubsystemTypeStr string,
	pAuthLoginURLstr   string,
	pKeyServer         *gf_identity_core.GFkeyServerInfo,
	pTemplatesPathsMap map[string]string,
	pHTTPmux           *http.ServeMux,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// STATIC FILES SERVING
	staticFilesURLbaseStr := "/v1/maps"
	localDirPathStr       := "./static"
	gf_core.HTTPinitStaticServingWithMux(staticFilesURLbaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)

	//------------------------
	// HANDLERS
	gfErr := initHandlers(pAuthSubsystemTypeStr,
		pAuthLoginURLstr,
		pKeyServer,
		pHTTPmux,
		pTemplatesPathsMap,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	return nil
}

//-------------------------------------------------