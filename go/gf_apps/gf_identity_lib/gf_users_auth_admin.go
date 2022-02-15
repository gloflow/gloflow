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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
)

//---------------------------------------------------
// io_login
type GF_user_auth_admin__input_login struct {

	User_name_str string `validate:"required,min=3,max=50"`

	// pass is not provided if email-login is used
	Pass_str string `validate:"omitempty,min=8,max=50"`

	// admin email
	Email_str string `validate:"omitempty,email"`
}

type GF_user_auth_admin__output_login struct {
	Email_confirmed_bool bool
	MFA_confirmed_bool   bool
	Pass_valid_bool      bool
	JWT_token_val        gf_session.GF_jwt_token_val
	User_id_str          gf_core.GF_ID 
}

type GF_user_auth_userpass__output_create_admin struct {
	General *GF_user_auth_userpass__output_create
}

//---------------------------------------------------
// PIPELINE__LOGIN
func Users_auth_admin__pipeline__login(p_input *GF_user_auth_admin__input_login,
	p_service_info *GF_service_info,
	p_ctx          context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) (*GF_user_auth_admin__output_login, *gf_core.GF_error) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------


	output := &GF_user_auth_admin__output_login{}

	//------------------------
	// VERIFY


	user_exists_bool, gf_err := db__user__exists_by_username(GF_user_name(p_input.User_name_str),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	var user_id_str gf_core.GF_ID

	// user doesnt exist
	if !user_exists_bool {
		
		input_create := &GF_user_auth_userpass__input_create{
			User_name_str: GF_user_name(p_input.User_name_str),
			Pass_str:      p_input.Pass_str,
			Email_str:     p_input.Email_str,
		}
		//------------------------	
		// PIPELINE__CREATE_ADMIN
		// if the admin user doesnt exist in the DB (most likely on first run of gloflow server),
		// create one in the DB
		output, gf_err := users_auth_admin__pipeline__create_admin(input_create,
			p_service_info,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		//------------------------

		user_id_str = output.General.User_id_str
	
	} else {
		existing_user_id_str, gf_err := db__user__get_basic_info_by_username(GF_user_name(p_input.User_name_str),
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		user_id_str = existing_user_id_str
	}



	// VERIFY_PASSWORD
	pass_valid_bool, gf_err := users_auth_userpass__verify_pass(GF_user_name(p_input.User_name_str),
		p_input.Pass_str,
		p_service_info,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	if !pass_valid_bool {
		output.Pass_valid_bool = false
		return output, nil
	} else {
		output.Pass_valid_bool = true
	}


	

	//------------------------
	// EMAIL
	if p_service_info.Enable_email_bool {

		gf_err = users_email__verify__pipeline(p_input.Email_str,
			user_id_str,
			p_service_info.Domain_base_str,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}
	}

	//------------------------
	// EVENT
	if p_service_info.Enable_events_app_bool {
		event_meta := map[string]interface{}{
			"user_id_str":     user_id_str,
			"user_name_str":   p_input.User_name_str,
			"domain_base_str": p_service_info.Domain_base_str,
		}
		gf_events.Emit_app(GF_EVENT_APP__USER_LOGIN_ADMIN,
			event_meta,
			p_runtime_sys)
	}

	//------------------------
	
	return output, nil
}

//---------------------------------------------------
func users_auth_admin__pipeline__create_admin(p_input *GF_user_auth_userpass__input_create,
	p_service_info *GF_service_info,
	p_ctx          context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_create_admin, *gf_core.GF_error) {




	//------------------------
	// PIPELINE
	output, gf_err := users_auth_userpass__pipeline__create(p_input,
		p_service_info,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}


	//------------------------
	// EVENT
	if p_service_info.Enable_events_app_bool {
		event_meta := map[string]interface{}{
			"user_id_str":     output.User_id_str,
			"user_name_str":   p_input.User_name_str,
			"domain_base_str": p_service_info.Domain_base_str,
		}
		gf_events.Emit_app(GF_EVENT_APP__USER_CREATE_ADMIN,
			event_meta,
			p_runtime_sys)
	}

	//------------------------

	output_admin := &GF_user_auth_userpass__output_create_admin{
		General: output,
	}
	return output_admin, nil

}