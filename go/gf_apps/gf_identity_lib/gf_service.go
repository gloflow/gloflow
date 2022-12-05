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

package gf_identity_lib

import (
	"os"
	"flag"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
)

//-------------------------------------------------



//-------------------------------------------------

func InitService(pHTTPmux *http.ServeMux,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	
	//------------------------
	// HANDLERS
	gfErr := initHandlers(pServiceInfo.AuthLoginURLstr,
		pHTTPmux, pServiceInfo, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// ETH
	gfErr = initHandlersEth(pHTTPmux, pServiceInfo, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	switch pServiceInfo.AuthSubsystemTypeStr {

	// USERPASS
	case gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__BUILTIN:
		
		gfErr = initHandlersUserpass(pHTTPmux,
			pServiceInfo,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

	// AUTH0
	case gf_identity_core.GF_AUTH_SUBSYSTEM_TYPE__AUTH0:
		
		auth0config := gf_auth0.Init()
		initHandlersAuth0(pHTTPmux,
			auth0config,
			pServiceInfo,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	//------------------------

	return nil
}

//-------------------------------------------------

func CLIparseArgs(pLogFun func(string, string)) map[string]interface{} {

	//-------------------
	// MONGODB
	mongodb_host_str    := flag.String("mongodb_host",    "127.0.0.1", "host of mongodb to use")
	mongodb_db_name_str := flag.String("mongodb_db_name", "prod_db",   "DB name to use")

	// MONGODB_ENV
	mongodb_host_env_str    := os.Getenv("GF_MONGODB_HOST")
	mongodb_db_name_env_str := os.Getenv("GF_MONGODB_DB_NAME")

	if mongodb_db_name_env_str != "" {
		*mongodb_db_name_str = mongodb_db_name_env_str
	}

	if mongodb_host_env_str != "" {
		*mongodb_host_str = mongodb_host_env_str
	}

	//-------------------
	flag.Parse()

	return map[string]interface{}{
		"mongodb_host_str":    *mongodb_host_str,
		"mongodb_db_name_str": *mongodb_db_name_str,
	}
}