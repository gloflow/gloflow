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
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func init_handlers(p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_identity_lib.init_handlers()")

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/identity/update",
		"/v1/identity/me",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers(handlers_endpoints_lst)

	//---------------------
	// USERS_UPDATE
	// AUTH - only logged in users can update their own details

	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/update",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//---------------------
				// SESSION_VALIDATE
				valid_bool, _, gf_err := Session__validate(p_req, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				if !valid_bool {
					return nil, nil
				}

				//---------------------
				// INPUT
				http_input, gf_err := http__get_user_update(p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input := &GF_user__input_update{
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

				_, gf_err = users__pipeline__update(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{}
				return output_map, nil
			}
			return nil, nil
		},
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------
	// USERS_GET_ME
	// AUTH
	gf_rpc_lib.Create_handler__http_with_metrics("/v1/identity/me",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//---------------------
				// SESSION_VALIDATE
				valid_bool, me_user_identifier_str, gf_err := Session__validate(p_req, p_ctx, p_runtime_sys)
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
		metrics,
		true, // p_store_run_bool
		p_runtime_sys)

	//---------------------

	return nil
}