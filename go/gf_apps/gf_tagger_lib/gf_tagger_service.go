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

package gf_tagger_lib

import (
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------

func InitService(pAuthSubsystemTypeStr string,
	pAuthLoginURLstr   string,
	pKeyServer         *gf_identity_core.GFkeyServerInfo,
	pHTTPmux           *http.ServeMux,
	pTemplatesPathsMap map[string]string,
	pImagesJobsMngr    gf_images_jobs_core.JobsMngr,
	pRuntimeSys        *gf_core.RuntimeSys) {
	
	// DB
	gfErr := dbSQLcreateTables(pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//------------------------
	// DB_INDEXES
	gfErr = DBmongoIndexInit(pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}
	
	//------------------------
	// STATIC FILES SERVING
	urlBaseStr      := "/tags"
	localDirPathStr := "./static"
	gf_core.HTTPinitStaticServingWithMux(urlBaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)

	//------------------------
	
	gfErr = initHandlers(pAuthSubsystemTypeStr,
		pAuthLoginURLstr,
		pKeyServer,
		pHTTPmux,
		pTemplatesPathsMap,
		pImagesJobsMngr,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//------------------------
}