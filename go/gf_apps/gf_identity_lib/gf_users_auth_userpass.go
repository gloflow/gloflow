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
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
)

//---------------------------------------------------
// io_login
type GF_user_auth_userpass__input_login struct {

	// username is always required, with both pass and email login
	User_name_str gf_identity_core.GFuserName `validate:"required,min=3,max=50"`

	// pass is not provided if email-login is used
	Pass_str string `validate:"omitempty,min=8,max=50"`

	// for certain emails allow email-login
	Email_str string `validate:"omitempty,email"`
}
type GF_user_auth_userpass__output_login struct {
	User_exists_bool     bool
	Email_confirmed_bool bool
	Pass_valid_bool      bool
	User_id_str          gf_core.GF_ID 
	JWT_token_val        gf_session.GF_jwt_token_val
}

// io_login_finalize
type GF_user_auth_userpass__input_login_finalize struct {
	// UserIDstr gf_core.GF_ID `validate:"required"`
	UserNameStr gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
}
type GF_user_auth_userpass__output_login_finalize struct {
	Email_confirmed_bool bool
	User_id_str          gf_core.GF_ID 
	JWT_token_val        gf_session.GF_jwt_token_val
}

// io_create
type GF_user_auth_userpass__input_create struct {
	User_name_str gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
	Pass_str      string                      `validate:"required,min=8,max=50"`
	Email_str     string                      `validate:"required,email"`
	UserTypeStr   string                      `validate:"required"` // "admin"|"standard"
}
type GF_user_auth_userpass__output_create_regular struct {
	User_exists_bool         bool
	User_in_invite_list_bool bool
	General                  *GF_user_auth_userpass__output_create
}
type GF_user_auth_userpass__output_create struct {
	User_name_str gf_identity_core.GFuserName
	User_id_str   gf_core.GF_ID
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

	user_exists_bool, gf_err := db__user__exists_by_username(gf_identity_core.GFuserName(p_input.User_name_str),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// user doesnt exists, so abort login
	if !user_exists_bool {
		output.User_exists_bool = false
		return output, nil
	} else {
		output.User_exists_bool = true
	}

	//------------------------
	// VERIFY_PASSWORD
	pass_valid_bool, gf_err := users_auth_userpass__verify_pass(gf_identity_core.GFuserName(p_input.User_name_str),
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
	// LOGIN_FINALIZE
	input := &GF_user_auth_userpass__input_login_finalize{
		UserNameStr: gf_identity_core.GFuserName(p_input.User_name_str),
	}
	login_finalize_output, gf_err := users_auth_userpass__pipeline__login_finalize(input,
		p_service_info,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------
	output.Email_confirmed_bool = login_finalize_output.Email_confirmed_bool
	output.User_id_str          = login_finalize_output.User_id_str
	output.JWT_token_val        = login_finalize_output.JWT_token_val

	return output, nil
}

//---------------------------------------------------
func users_auth_userpass__pipeline__login_finalize(pInput *GF_user_auth_userpass__input_login_finalize,
	pServiceInfo *GF_service_info,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_login_finalize, *gf_core.GF_error) {

	output := &GF_user_auth_userpass__output_login_finalize{}
	userNameStr := gf_identity_core.GFuserName(pInput.UserNameStr)

	//------------------------
	// VERIFY_EMAIL_CONFIRMED
	// if this check is enabled, users that have not confirmed their email cant login.
	// this is the initial confirmation of an email on user creation, or user email update.
	if pServiceInfo.Enable_email_require_confirm_for_login_bool {

		emailConfirmedBool, gfErr := db__user__get_email_confirmed_by_username(userNameStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		if !emailConfirmedBool {
			output.Email_confirmed_bool = false
			return output, nil
		} else {
			output.Email_confirmed_bool = true
		}
	}

	//------------------------
	// USER_ID
	
	userIDstr, gfErr := gf_identity_core.DBgetBasicInfoByUsername(userNameStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.User_id_str = userIDstr

	//------------------------
	// JWT
	userIdentifierStr := string(userIDstr)
	JWTtokenVal, gfErr := gf_session.JWT__pipeline__generate(userIdentifierStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.JWT_token_val = JWTtokenVal

	//------------------------
	return output, nil
}

//---------------------------------------------------
// PIPELINE__CREATE_REGULAR
func users_auth_userpass__pipeline__create_regular(p_input *GF_user_auth_userpass__input_create,
	p_service_info *GF_service_info,
	pCtx           context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_create_regular, *gf_core.GF_error) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output_regular := &GF_user_auth_userpass__output_create_regular{}

	//------------------------
	// VALIDATE

	user_exists_bool, gf_err := db__user__exists_by_username(p_input.User_name_str, pCtx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// user already exists, so abort creation
	if user_exists_bool {
		output_regular.User_exists_bool = true
		return output_regular, nil
	}

	// check if in invite list
	in_invite_list_bool, gf_err := db__user__check_in_invitelist_by_email(p_input.Email_str,
		pCtx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// user is not in the invite list, so abort the creation
	if in_invite_list_bool {
		output_regular.User_in_invite_list_bool = true
	} else {
		output_regular.User_in_invite_list_bool = false
		return output_regular, nil
	}

	//------------------------
	// PIPELINE
	output, gf_err := users_auth_userpass__pipeline__create(p_input,
		p_service_info,
		pCtx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output_regular.General = output

	//------------------------
	// EMAIL
	if p_service_info.Enable_email_bool {

		gf_err = usersEmailPipelineVerify(p_input.Email_str,
			p_input.User_name_str,
			output.User_id_str,
			p_service_info.Domain_base_str,
			pCtx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}
	}
	
	//------------------------
	// EVENT
	if p_service_info.Enable_events_app_bool {
		event_meta := map[string]interface{}{
			"user_id_str":     output.User_id_str,
			"user_name_str":   p_input.User_name_str,
			"domain_base_str": p_service_info.Domain_base_str,
		}
		gf_events.Emit_app(GF_EVENT_APP__USER_CREATE_REGULAR,
			event_meta,
			p_runtime_sys)
	}

	//------------------------

	return output_regular, nil
}

//---------------------------------------------------
// PIPELINE__CREATE
func users_auth_userpass__pipeline__create(pInput *GF_user_auth_userpass__input_create,
	pServiceInfo *GF_service_info,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_create, *gf_core.GF_error) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(pInput, pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	creation_unix_time_f  := float64(time.Now().UnixNano())/1000000000.0
	userTypeStr   := pInput.UserTypeStr
	userNameStr   := pInput.User_name_str
	pass_str      := pInput.Pass_str
	email_str     := pInput.Email_str

	user_identifier_str := string(userNameStr)
	user_id_str := usersCreateID(user_identifier_str, creation_unix_time_f)

	user := &GFuser{
		V_str:                "0",
		Id_str:               user_id_str,
		Creation_unix_time_f: creation_unix_time_f,
		UserTypeStr:          userTypeStr,
		User_name_str:        userNameStr,
		Email_str:            email_str,
	}

	
	pass_salt_str := users_auth_userpass__get_pass_salt()
	pass_hash_str := users_auth_userpass__get_pass_hash(pass_str, pass_salt_str)

	creds__creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	user_creds_id_str           := usersCreateID(user_identifier_str, creds__creation_unix_time_f)

	userCreds := &GF_user_creds {
		V_str:                "0",
		Id_str:               user_creds_id_str,
		Creation_unix_time_f: creds__creation_unix_time_f,
		User_id_str:          user_id_str,
		User_name_str:        userNameStr,
		Pass_salt_str:        pass_salt_str,
		Pass_hash_str:        pass_hash_str,
	}

	//------------------------
	// USER_PERSIST
	// DB__USER_CREATE
	gf_err = db__user__create(user, pCtx, pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------
	// USER_CREDS_PERSIST

	// SECRETS_STORE
	if pServiceInfo.Enable_user_creds_in_secrets_store_bool && 
		pRuntimeSys.External_plugins.Secret_app__create_callback != nil {

		secretNameStr := fmt.Sprintf("gf_user_creds@%s", userNameStr)
		secretDescriptionStr := fmt.Sprintf("user creds for a particular user")

		userCredsMap := map[string]interface{}{
			"user_creds_id_str":    user_creds_id_str, 
			"creation_unix_time_f": creds__creation_unix_time_f,
			"user_id_str":          user_id_str,
			"user_name_str":        userNameStr,
			"pass_salt_str":        pass_salt_str,
			"pass_hash_str":        pass_hash_str,
		}

		// SECRET_STORE__USER_CREDS_CREATE
		gfErr := pRuntimeSys.External_plugins.Secret_app__create_callback(secretNameStr,
			userCredsMap,
			secretDescriptionStr,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	} else {

		// DB__USER_CREDS_CREATE - otherwise use the regular DB
		gf_err = db__user_creds__create(userCreds, pCtx, pRuntimeSys)
		if gf_err != nil {
			return nil, gf_err
		}
	}
	

	//------------------------

	// LOGIN_ATTEMPT
	// on user creation initiate a login process that completes after the user
	// confirms their email.
	_, gfErr := loginAttempCreate(userNameStr, userTypeStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}


	//------------------------
	// EMAIL_VERIFY_ADDRESS
	if pServiceInfo.Enable_email_bool {

		// this SES email verification is done only once for a new email address,
		// so that SES allows sending to this email address.
		gf_err = gf_aws.AWS_SES__verify_address(email_str,
			pRuntimeSys)
		if gf_err != nil {
			return nil, gf_err
		}
	}

	//------------------------

	output := &GF_user_auth_userpass__output_create{
		User_name_str: userNameStr,
		User_id_str:   user_id_str,
	}

	return output, nil
}

//---------------------------------------------------
// PASS
//---------------------------------------------------
func users_auth_userpass__verify_pass(p_user_name_str gf_identity_core.GFuserName,
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

		secret_name_str := fmt.Sprintf("gf_user_creds@%s", p_user_name_str)
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