/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	"fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//------------------------------------------------
func init_handlers(p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_identity_lib.init_handlers()")

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/identity/users/preflight",
		"/v1/identity/users/login",
		"/v1/identity/users/create",
		"/v1/identity/users/update",
		"/v1/identity/users/get",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(handlers_endpoints_lst)

	//---------------------
	// USERS_PREFLIGHT
	// NO_AUTH
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/preflight",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				input :=&GF_user__input_preflight{
					User_address_eth_str: GF_user_address_eth(input_map["user_address_eth_str"].(string)),
				}

				//---------------------

				output, gf_err := users__pipeline__preflight(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"user_exists_bool": output.User_exists_bool,
					"nonce_val_str":    output.Nonce_val_str,
				}
				return output_map, nil
			}

			return nil, nil
		},
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	// USERS_LOGIN
	// NO_AUTH
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/login",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input :=&GF_user__input_login{
					User_address_eth_str: GF_user_address_eth(input_map["user_address_eth_str"].(string)),
					Auth_signature_str:   GF_auth_signature(input_map["auth_signature_str"].(string)),
				}

				//---------------------
				// LOGIN
				output, gf_err := users__pipeline__login(input, p_ctx, p_runtime_sys)
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
					"auth_signature_valid_bool": output.Auth_signature_valid_bool,
					"nonce_exists_bool":         output.Nonce_exists_bool,
					"user_id_str":               output.User_id_str,
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
	// NO_AUTH
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input :=&GF_user__input_create{
					User_address_eth_str: GF_user_address_eth(input_map["user_address_eth_str"].(string)),
					Auth_signature_str:   GF_auth_signature(input_map["auth_signature_str"].(string)),
				}

				//---------------------
				output, gf_err := users__pipeline__create(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"auth_signature_valid_bool": output.Auth_signature_valid_bool,
					"nonce_exists_bool":         output.Nonce_exists_bool,
				}
				return output_map, nil
			}

			return nil, nil
		},
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	// USERS_UPDATE
	// AUTH
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/update",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input := &GF_user__input_update{
					User_address_eth_str: GF_user_address_eth(input_map["user_address_eth_str"].(string)),
					Username_str:         input_map["user_username_str"].(string),
					Email_str:            input_map["user_email_str"].(string),
					Description_str:      input_map["user_description_str"].(string),
				}

				

				//---------------------

				_, gf_err = users__pipeline__update(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
			}

			return nil, nil
		},
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	// USERS_GET
	// AUTH/NO_AUTH
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/get",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// SESSION_VALIDATE
				gf_err := session__validate(p_req, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------

				input :=&GF_user__input_get{}

				_, gf_err = users__pipeline__get(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
			}

			return nil, nil
		},
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------

	return nil
}