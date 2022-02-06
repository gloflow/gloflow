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
	"fmt"
	"time"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
)

//---------------------------------------------------
// io_login
type GF_user_auth_userpass__input_login struct {

	// username is always required, with both pass and email login
	User_name_str string `validate:"required,min=3,max=50"`

	// pass is not provided if email-login is used
	Pass_str string `validate:"omitempty,min=8,max=50"`

	// for certain emails allow email-login
	Email_str string `validate:"omitempty,email"`
}
type GF_user_auth_userpass__output_login struct {
	User_exists_bool     bool
	Email_confirmed_bool bool
	Pass_valid_bool      bool
	JWT_token_val        GF_jwt_token_val
	User_id_str          gf_core.GF_ID 
}

// io_create
type GF_user_auth_userpass__input_create struct {
	User_name_str GF_user_name `validate:"required,min=3,max=50"`
	Pass_str      string       `validate:"required,min=8,max=50"`
	Email_str     string       `validate:"required,email"`
}
type GF_user_auth_userpass__output_create struct {
	User_exists_bool         bool
	User_in_invite_list_bool bool
}

//---------------------------------------------------
// PIPELINE__LOGIN
func users_auth_userpass__pipeline__login(p_input *GF_user_auth_userpass__input_login,
	p_service_info *GF_service_info,
	p_ctx          context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_login, *gf_core.GF_error) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user_auth_userpass__output_login{}

	//------------------------
	// VERIFY

	user_exists_bool, gf_err := db__user__exists_by_username(GF_user_name(p_input.User_name_str),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// user doesnt exists, so abort login
	if !user_exists_bool {
		output.User_exists_bool = false
		return output, nil
	}

	// VERIFY_EMAIL_CONFIRMED
	// if this check is enabled, users that have not confirmed their email cant login
	if p_service_info.Enable_email_require_confirm_for_login_bool {

		email_confirmed_bool, gf_err := db__user__get_email_confirmed_by_username(GF_user_name(p_input.User_name_str),
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		if !email_confirmed_bool {
			output.Email_confirmed_bool = false
			return output, nil
		} else {
			output.Email_confirmed_bool = true
		}
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
	// USER_ID
	user_id_str, gf_err := db__user__get_basic_info_by_username(GF_user_name(p_input.User_name_str),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output.User_id_str = user_id_str

	//------------------------
	// JWT
	user_identifier_str := string(user_id_str)
	jwt_token_val, gf_err := jwt__pipeline__generate(user_identifier_str, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output.JWT_token_val = jwt_token_val

	//------------------------

	return output, nil
}

//---------------------------------------------------
// PIPELINE__CREATE
func users_auth_userpass__pipeline__create(p_input *GF_user_auth_userpass__input_create,
	p_service_info *GF_service_info,
	p_ctx          context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_create, *gf_core.GF_error) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user_auth_userpass__output_create{}

	//------------------------
	// VALIDATE

	user_exists_bool, gf_err := db__user__exists_by_username(p_input.User_name_str, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// user already exists, so abort creation
	if user_exists_bool {
		output.User_exists_bool = true
		return output, nil
	}

	// check if in invite list
	in_invite_list_bool, gf_err := db__user__check_in_invitelist_by_username(p_input.Email_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// user is not in the invite list, so abort the creation
	if !in_invite_list_bool {
		output.User_in_invite_list_bool = false
		return output, nil
	} else {
		output.User_in_invite_list_bool = true
	}

	

	//------------------------

	creation_unix_time_f  := float64(time.Now().UnixNano())/1000000000.0
	user_name_str := p_input.User_name_str
	pass_str      := p_input.Pass_str
	email_str     := p_input.Email_str

	user_identifier_str := string(user_name_str)
	user_id_str := users__create_id(user_identifier_str, creation_unix_time_f)

	user := &GF_user{
		V_str:                "0",
		Id_str:               user_id_str,
		Creation_unix_time_f: creation_unix_time_f,
		User_name_str:        user_name_str,
		Email_str:            email_str,
	}

	
	pass_salt_str := users_auth_userpass__get_pass_salt()
	pass_hash_str := users_auth_userpass__get_pass_hash(pass_str, pass_salt_str)

	creds__creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	user_creds_id_str           := users__create_id(user_identifier_str, creds__creation_unix_time_f)

	user_creds := &GF_user_creds {
		V_str:                "0",
		Id_str:               user_creds_id_str,
		Creation_unix_time_f: creds__creation_unix_time_f,
		User_id_str:          user_id_str,
		User_name_str:        user_name_str,
		Pass_salt_str:        pass_salt_str,
		Pass_hash_str:        pass_hash_str,
	}

	//------------------------
	// DB
	gf_err = db__user__create(user, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// SECRETS_STORE
	if p_service_info.Enable_user_creds_in_secrets_store_bool && 
		p_runtime_sys.External_plugins.Secret_app__create_callback != nil {

		secret_name_str := fmt.Sprintf("gf_user_creds:%s", user_name_str)
		secret_description_str := fmt.Sprintf("user creds for a particular user")



		user_creds_map := map[string]interface{}{
			"user_creds_id_str":    user_creds_id_str, 
			"creation_unix_time_f": creds__creation_unix_time_f,
			"user_id_str":          user_id_str,
			"user_name_str":        user_name_str,
			"pass_salt_str":        pass_salt_str,
			"pass_hash_str":        pass_hash_str,
		}



		gf_err := p_runtime_sys.External_plugins.Secret_app__create_callback(secret_name_str,
			user_creds_map,
			secret_description_str,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}
	} else {


		// DB - otherwise use the regular DB
		gf_err = db__user_creds__create(user_creds, p_ctx, p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}
	}
	

	//------------------------
	// EMAIL
	if p_service_info.Enable_email_bool {

		gf_err = users_email__verify__pipeline(email_str,
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
			"user_name_str":   user_name_str,
			"domain_base_str": p_service_info.Domain_base_str,
		}
		gf_events.Emit_app(GF_EVENT_APP__USER_CREATE,
			event_meta,
			p_runtime_sys)
	}

	//------------------------

	return output, nil
}

//---------------------------------------------------
// PASS
//---------------------------------------------------
func users_auth_userpass__verify_pass(p_user_name_str GF_user_name,
	p_pass_str     string,
	p_service_info *GF_service_info,
	p_ctx          context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {


	// GET_PASS_AND_SALT

	
	var pass_salt__loaded_str string
	var pass_hash__loaded_str string

	// SECRETS_STORE
	if p_service_info.Enable_user_creds_in_secrets_store_bool && 
		p_runtime_sys.External_plugins.Secret_app__get_callback != nil {

		secret_name_str := fmt.Sprintf("gf_user_creds:%s", p_user_name_str)
		secret_map, gf_err := p_runtime_sys.External_plugins.Secret_app__get_callback(secret_name_str,
			p_runtime_sys)
		if gf_err != nil {
			return false, gf_err
		}

		pass_salt__loaded_str = secret_map["pass_salt_str"].(string)
		pass_hash__loaded_str = secret_map["pass_hash_str"].(string)
		
	} else {

		// DB
		db_pass_salt_str, db_pass_hash_str, gf_err := db__user_creds__get_pass_hash(p_user_name_str,
			p_ctx, p_runtime_sys)
		if gf_err != nil {
			return false, gf_err
		}

		pass_salt__loaded_str = db_pass_salt_str
		pass_hash__loaded_str = db_pass_hash_str
	}

	// GENERATE_PASS_HASH
	pass_hash__expected_str := users_auth_userpass__get_pass_hash(p_pass_str, pass_salt__loaded_str)


	if (pass_hash__loaded_str == pass_hash__expected_str) {
		return true, nil
	} else {
		return false, nil
	}

	return false, nil
}

//---------------------------------------------------
func users_auth_userpass__get_pass_hash(p_pass_str string,
	p_pass_salt_str string) string {

	salted_pass_str := fmt.Sprintf("%s:%s", p_pass_salt_str, p_pass_str)
	pass_hash_str   := gf_core.Hash_val_sha256(salted_pass_str)
	return pass_hash_str
}

//---------------------------------------------------
func users_auth_userpass__get_pass_salt() string {
	rand_str := gf_core.Str_random()
	return rand_str
}