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

	// DOMAIN - where this gf_solo instance is reachable on
	Domain_base_str string

	// EVENTS_APP - enable sending of app events from various functions
	Enable_events_app_bool bool

	// enable storage of user_creds in a secret store
	Enable_user_creds_in_secrets_store_bool bool

	// enable sending of emails for any function that needs it
	Enable_email_bool bool

	// enable login only for users that have confirmed their email
	Enable_email_require_confirm_for_login_bool bool
}

//-------------------------------------------------
func Init_service(p_mux *http.ServeMux,
	p_service_info *GF_service_info,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.GF_error {

	//------------------------
	// HANDLERS
	gf_err := init_handlers(p_mux, p_service_info, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	gf_err = init_handlers__eth(p_mux, p_service_info, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	gf_err = init_handlers__userpass(p_mux, p_service_info, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------------

	return nil
}

//-------------------------------------------------
func CLI__parse_args(p_log_fun func(string, string)) map[string]interface{} {

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