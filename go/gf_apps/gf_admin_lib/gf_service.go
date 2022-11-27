/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_admin_lib

import (
	"net/http"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
)

//-------------------------------------------------

type GFserviceInfo struct {

	NameStr string

	// ADMIN_EMAIL - what the default admin email is (for auth)
	AdminEmailStr string

	// EVENTS_APP - enable sending of app events from various functions
	EnableEventsAppBool bool

	// enable storage of user_creds in a secret store
	EnableUserCredsInSecretsStoreBool bool

	// enable sending of emails for any function that needs it
	EnableEmailBool bool
}

//-------------------------------------------------

func InitNewService(pTemplatesPathsMap map[string]string,
	pServiceInfo         *GFserviceInfo,
	pIdentityServiceInfo *gf_identity_lib.GFserviceInfo,
	pHTTPmux             *http.ServeMux,
	pLocalHub            *sentry.Hub,
	pRuntimeSys          *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// STATIC FILES SERVING
	staticFilesURLbaseStr := "/v1/admin"
	localDirPathStr       := "./static"

	gf_core.HTTPinitStaticServingWithMux(staticFilesURLbaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)
		
	//------------------------
	// IDENTITY_HANDLERS

	gfErr := gf_identity_lib.InitService(pHTTPmux,
		pIdentityServiceInfo,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}

	//------------------------
	// ADMIN_HANDLERS
	
	gfErr = initHandlers(pTemplatesPathsMap,
		pHTTPmux,
		pServiceInfo,
		pIdentityServiceInfo,
		pLocalHub,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	gfErr = initHandlersUsers(pHTTPmux,
		pServiceInfo,
		pIdentityServiceInfo,
		pLocalHub,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	return nil
}