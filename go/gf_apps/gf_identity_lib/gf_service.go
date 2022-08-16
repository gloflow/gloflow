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
)

//-------------------------------------------------
type GF_service_info struct {
	
	// name of this service, in case multiple are spawned
	Name_str string

	// DOMAIN - where this gf_solo instance is reachable on
	Domain_base_str string

	// ADMIN_MFA_SECRET_KEY_BASE32
	Admin_mfa_secret_key_base32_str string

	//------------------------
	// AUTH

	// AUTH_LOGIN_URL - url of the login page to which the system should
	//                  redirect users after certain operations
	AuthLoginURLstr string

	// AUTH_LOGIN_SUCCESS_REDIRECT_URL - url to redirect to when the user 
	//                                   logs in successfuly. if ommited then dont redirect.
	AuthLoginSuccessRedirectURLstr string

	//------------------------
	// FEATURE_FLAGS

	// EVENTS_APP - enable sending of app events from various functions
	Enable_events_app_bool bool

	// enable storage of user_creds in a secret store
	Enable_user_creds_in_secrets_store_bool bool

	// enable sending of emails for any function that needs it
	Enable_email_bool bool

	// enable login only for users that have confirmed their email
	Enable_email_require_confirm_for_login_bool bool

	// enable login only for users that have confirmed their MFA code
	Enable_mfa_require_confirm_for_login_bool bool

	//------------------------
}

//-------------------------------------------------
func InitService(pHTTPmux *http.ServeMux,
	pServiceInfo *GF_service_info,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GF_error {

	//------------------------
	// HANDLERS
	gfErr := initHandlers(pServiceInfo.AuthLoginURLstr,
		pHTTPmux, pServiceInfo, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	gfErr = init_handlers__eth(pHTTPmux, pServiceInfo, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	gfErr = initHandlersUserpass(pHTTPmux,
		pServiceInfo,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	return nil
}

//-------------------------------------------------
func CLI__parse_args(pLogFun func(string, string)) map[string]interface{} {

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