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
	// "fmt"
	"net/http"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func init_handlers__users(p_http_mux *http.ServeMux,
	p_service_info          *GF_service_info,
	p_identity_service_info *gf_identity_lib.GF_service_info,
	p_local_hub             *sentry.Hub,
	p_runtime_sys           *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/admin/users/get_all_invite_list",
		"/v1/admin/users/add_to_invite_list",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers("gf_admin", handlers_endpoints_lst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpc_handler_runtime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                p_http_mux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         p_local_hub,
		Auth_login_url_str: "/v1/admin/login_ui",
	}

	//---------------------
	// GET_ALL_INVITE_LIST
	// AUTH
	gf_rpc_lib.Create_handler__http_with_auth(true, "/v1/admin/users/get_all_invite_list",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				
				_, user_name_str, _, gf_err := gf_identity_lib.Http__get_user_std_input(p_ctx, p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}


				gf_err = gf_identity_lib.Admin__is(user_name_str, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------

				invite_list_lst, gf_err := gf_identity_lib.Admin__pipeline__get_all_invite_list(p_ctx,
					p_identity_service_info,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"invite_list_lst": invite_list_lst,
				}
				return output_map, nil
			}
			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// ADD_TO_INVITE_LIST
	// AUTH
	gf_rpc_lib.Create_handler__http_with_auth(true, "/v1/admin/users/add_to_invite_list",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				
				input_map, user_name_str, _, gf_err := gf_identity_lib.Http__get_user_std_input(p_ctx, p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				var email_str string
				if val_str, ok := input_map["email_str"]; ok {
					email_str = val_str.(string)
				}

				gf_err = gf_identity_lib.Admin__is(user_name_str, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input := &gf_identity_lib.GF_admin__input_add_to_invite_list{
					User_name_str: gf_identity_lib.GF_user_name(user_name_str),
					Email_str:     email_str,
				}

				//---------------------

				gf_err = gf_identity_lib.Admin__pipeline__user_add_to_invite_list(input,
					p_ctx,
					p_identity_service_info,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					
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