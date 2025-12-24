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

package gf_identity

import (
	"context"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_policy"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
)

//-------------------------------------------------

func InitService(pTemplatesPathsMap map[string]string,
	pHTTPmux          *http.ServeMux,
	pRPCglobalMetrics *gf_rpc_lib.GFglobalMetrics,
	pServiceInfo      *gf_identity_core.GFserviceInfo,
	pRuntimeSys       *gf_core.RuntimeSys) (*gf_identity_core.GFkeyServerInfo, *gf_core.GFerror) {
	
	pRuntimeSys.LogNewFun("INFO", "initializing gf_identity service...", map[string]interface{}{
		"auth_subsystem_type_str": pServiceInfo.AuthSubsystemTypeStr,
	})
	
	//------------------------
	// DB
	ctx := context.Background()
	gfErr := gf_identity_core.DBsqlCreateTables(ctx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// KEYS_SERVER

	auth0initBool := false
	if pServiceInfo.AuthSubsystemTypeStr == gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__AUTH0 {
		auth0initBool = true
	}

	keyServerInfo, gfErr := gf_identity_core.KSinit(auth0initBool, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//------------------------
	// POLICIES

	gfErr = gf_policy.DBsqlCreateTables(pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//------------------------
	// STATIC FILES SERVING
	staticFilesURLbaseStr := "/v1/identity"
	localDirPathStr       := "./static"

	gf_core.HTTPinitStaticServingWithMux(staticFilesURLbaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)

	//------------------------
	// HANDLERS
	gfErr = initHandlers(
		pServiceInfo.AuthLoginURLstr,
		pTemplatesPathsMap,
		keyServerInfo,
		pHTTPmux,
		pRPCglobalMetrics,
		pServiceInfo,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// ETH - these handlers are always enabled, whether userpass|auth0 auth subsystem is activated
	gfErr = initHandlersEth(keyServerInfo, pHTTPmux, pServiceInfo, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	switch pServiceInfo.AuthSubsystemTypeStr {

	//------------------------
	// USERPASS
	case gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__USERPASS:
		
		gfErr = initHandlersUserpass(keyServerInfo,
			pHTTPmux,
			pServiceInfo,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	
	//------------------------
	// AUTH0
	case gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__AUTH0:
		
		auth0authenticator, auth0config, gfErr := gf_auth0.Init(pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		initHandlersAuth0(keyServerInfo,
			pHTTPmux,
			auth0authenticator,
			auth0config,
			pServiceInfo,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}

	//------------------------

	return keyServerInfo, nil
}

//-------------------------------------------------