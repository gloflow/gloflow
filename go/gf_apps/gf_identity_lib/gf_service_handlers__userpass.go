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
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func init_handlers__userpass(p_http_mux *http.ServeMux,
	p_service_info *GF_service_info,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/identity/userpass/login",
		"/v1/identity/userpass/create",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(p_service_info.Name_str, handlers_endpoints_lst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpc_handler_runtime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                p_http_mux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: "/landing/main/",
	}

	//---------------------
	// USERS_LOGIN
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/userpass/login",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, _, _, gf_err := gf_identity_core.Http__get_user_std_input(p_ctx, p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				var userNameStr string
				if valStr, ok := input_map["user_name_str"]; ok {
					userNameStr = valStr.(string)
				}

				var pass_str string
				if valStr, ok := input_map["pass_str"]; ok {
					pass_str = valStr.(string)
				}

				var email_str string
				if valStr, ok := input_map["email_str"]; ok {
					email_str = valStr.(string)
				}

				input :=&GF_user_auth_userpass__input_login{
					User_name_str: gf_identity_core.GFuserName(userNameStr),
					Pass_str:      pass_str,
					Email_str:     email_str,
				}

				//---------------------
				// LOGIN
				output, gf_err := users_auth_userpass__pipeline__login(input, 
					p_service_info,
					p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------
				// SET_SESSION_ID - sets gf_sid cookie on all future requests
				sessionDataStr        := string(output.JWT_token_val)
				session_ttl_hours_int := 24 // 1 day
				gf_session.SetOnReq(sessionDataStr, p_resp, session_ttl_hours_int)

				//---------------------

				outputMap := map[string]interface{}{
					"user_exists_bool": output.User_exists_bool,
					"pass_valid_bool":  output.Pass_valid_bool,
					"user_id_str":      output.User_id_str,
				}
				return outputMap, nil
			}

			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// USERS_CREATE
	// NO_AUTH - unauthenticated users are able to create new users, and do not get logged in automatically on success

	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/userpass/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input :=&GF_user_auth_userpass__input_create{
					User_name_str: gf_identity_core.GFuserName(input_map["user_name_str"].(string)),
					Pass_str:      input_map["pass_str"].(string),
					Email_str:     input_map["email_str"].(string),
				}

				//---------------------
				output, gf_err := users_auth_userpass__pipeline__create_regular(input,
					p_service_info,
					p_ctx,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"user_exists_bool":         output.User_exists_bool,
					"user_in_invite_list_bool": output.User_in_invite_list_bool,
				}
				return output_map, nil
			}

			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	return nil
}