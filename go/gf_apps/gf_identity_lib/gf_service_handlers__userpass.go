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

package gf_identity_lib

import (
	// "fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func init_handlers__userpass(p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/identity/userpass/login",
		"/v1/identity/userpass/create",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(handlers_endpoints_lst)

	//---------------------
	// USERS_LOGIN
	// NO_AUTH
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/userpass/login",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, user_name_str, _, gf_err := http__get_user_std_input(p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				pass_hash_str := input_map["pass_hash_str"].(string)

				input :=&GF_user_auth_userpass__input_login{
					User_name_str: user_name_str,
					Pass_hash_str: pass_hash_str,
				}

				//---------------------
				// LOGIN
				output, gf_err := users_auth_userpass__pipeline__login(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------
				// SET_SESSION_ID - sets gf_sid cookie on all future requests
				session_data_str      := string(output.JWT_token_val)
				session_ttl_hours_int := 24 // 1 day
				session__set_on_req(session_data_str, p_resp, session_ttl_hours_int)

				//---------------------

				output_map := map[string]interface{}{
					"user_exists_bool":     output.User_exists_bool,
					"pass_hash_valid_bool": output.Pass_hash_valid_bool,
					"user_id_str":          output.User_id_str,
				}
				return output_map, nil
			}

			return nil, nil
		},
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	// USERS_CREATE
	// NO_AUTH - unauthenticated users are able to create new users

	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/userpass/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input :=&GF_user_auth_userpass__input_create{
					User_name_str: input_map["user_name_str"].(string),
					Pass_hash_str: input_map["pass_hash_str"].(string),
					Email_str:     input_map["email_str"].(string),
				}

				//---------------------
				output, gf_err := users_auth_userpass__pipeline__create(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"user_exists_bool": output.User_exists_bool,
				}
				return output_map, nil
			}

			return nil, nil
		},
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	return nil
}