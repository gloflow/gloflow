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
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func init_handlers(p_auth_login_url_str string,
	p_http_mux     *http.ServeMux,
	p_service_info *GF_service_info,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/identity/email_confirm",
		"/v1/identity/mfa_confirm",
		"/v1/identity/update",
		"/v1/identity/me",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(p_service_info.Name_str, handlers_endpoints_lst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpc_handler_runtime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                p_http_mux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: p_auth_login_url_str,
	}

	//---------------------
	// EMAIL_CONFIRM
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/email_confirm",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//---------------------
				// INPUT
				http_input, gf_err := http__get_email_confirm_input(p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------

				confirmed_bool, fail_msg_str, gf_err := users_email__confirm__pipeline(http_input,
					p_ctx,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				if confirmed_bool {

					// redirect user to login page
					// "email_confirmed=1" - signals to the UI that email has been confirmed
					url_redirect_str := fmt.Sprintf("%s?email_confirmed=1&user_name=%s",
						rpc_handler_runtime.Auth_login_url_str,
						http_input.User_name_str)

					// REDIRECT
					http.Redirect(p_resp,
						p_req,
						url_redirect_str,
						301)

				} else {
					output_map := map[string]interface{}{
						"fail_msg_str": fail_msg_str,
					}
					return output_map, nil
				}
			}
			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// MFA_CONFIRM
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/mfa_confirm",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT

				input_map, user_name_str, _, gf_err := Http__get_user_std_input(p_ctx, p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				var extern_htop_value_str string
				if input_extern_htop_value_str, ok := input_map["mfa_val_str"].(string); ok {
					extern_htop_value_str = input_extern_htop_value_str
				}

				input := &GF_user_auth_mfa__input_confirm{
					User_name_str:         GF_user_name(user_name_str),
					Extern_htop_value_str: extern_htop_value_str,
					Secret_key_base32_str: p_service_info.Admin_mfa_secret_key_base32_str,
				}
				
				//---------------------
				
				valid_bool, gf_err := mfa__pipeline__confirm(input,
					p_ctx,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				if valid_bool {
					//---------------------
					// LOGIN_FINALIZE

					login_finalize_input := &GF_user_auth_userpass__input_login_finalize{
						User_name_str: GF_user_name(user_name_str),
					}
					login_finalize_output, gf_err := users_auth_userpass__pipeline__login_finalize(login_finalize_input,
						p_service_info,
						p_ctx,
						p_runtime_sys)
					if gf_err != nil {
						return nil, gf_err
					}

					//---------------------
					
					// SET_SESSION_ID - sets gf_sid cookie on all future requests
					session_data_str      := string(login_finalize_output.JWT_token_val)
					session_ttl_hours_int := 24 // 1 day
					gf_session.Set_on_req(session_data_str, p_resp, session_ttl_hours_int)

					//---------------------
				}

				output_map := map[string]interface{}{
					"mfa_valid_bool": valid_bool,
				}
				return output_map, nil
			}

			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// USERS_UPDATE
	// AUTH - only logged in users can update their own details

	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/update",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// SESSION_VALIDATE
				valid_bool, user_identifier_str, gf_err := gf_session.Validate(p_req, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				if !valid_bool {
					return nil, nil
				}

				user_name_str := user_identifier_str

				//---------------------
				// INPUT
				http_input, gf_err := http__get_user_update_input(p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input := &GF_user__input_update{
					User_name_str:         GF_user_name(user_name_str),
					Email_str:             http_input.Email_str,
					Description_str:       http_input.Description_str,
					Profile_image_url_str: http_input.Profile_image_url_str,
					Banner_image_url_str:  http_input.Banner_image_url_str,
				}
				
				// VALIDATE
				gf_err = gf_core.Validate_struct(input, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				
				//---------------------

				_, gf_err = users__pipeline__update(input,
					p_service_info,
					p_ctx,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{}
				return output_map, nil
			}
			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// USERS_GET_ME
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/identity/me",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//---------------------
				// SESSION_VALIDATE
				valid_bool, me_user_identifier_str, gf_err := gf_session.Validate(p_req, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				if !valid_bool {
					return nil, nil
				}

				//---------------------
				// INPUT
				user_id_str := gf_core.GF_ID(me_user_identifier_str)
				input := &GF_user__input_get{
					User_id_str: user_id_str,
				}

				//---------------------

				output, gf_err := users__pipeline__get(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"user_name_str":         output.User_name_str,
					"email_str":             output.Email_str,
					"description_str":       output.Description_str,
					"profile_image_url_str": output.Profile_image_url_str,
					"banner_image_url_str":  output.Banner_image_url_str,
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