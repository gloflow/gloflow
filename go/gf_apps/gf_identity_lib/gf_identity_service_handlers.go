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
		"/v1/identity/users/login",
		"/v1/identity/users/create",
		"/v1/identity/users/update",
		"/v1/identity/users/get",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(handlers_endpoints_lst)

	//---------------------
	// USERS_LOGIN
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/login",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {


				input :=&GF_user__input_login{}

				_, gf_err := users__pipeline__login(input, p_ctx, p_runtime_sys)
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
	// USERS_CREATE
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {


				input :=&GF_user__input_create{}

				_, gf_err := users__pipeline__create(input, p_ctx, p_runtime_sys)
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
	// USERS_UPDATE
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/update",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				// USER_ADDRESS_ETH
				user_address_eth, gf_err := http__get_user_address_eth(p_req, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------
				// JWT_VALIDATE
				gf_err = jwt__validate_from_req(user_address_eth, p_req, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------

				input :=&GF_user__input_update{}

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
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/users/get",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				// USER_ADDRESS_ETH
				user_address_eth, gf_err := http__get_user_address_eth(p_req, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------
				// JWT_VALIDATE
				gf_err = jwt__validate_from_req(user_address_eth, p_req, p_ctx, p_runtime_sys)
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