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
	"time"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
)

//---------------------------------------------------
// io_login
type GF_admin__input_login struct {

	User_name_str string `validate:"required,min=3,max=50"`

	// pass is not provided if email-login is used
	Pass_str string `validate:"omitempty,min=8,max=50"`

	// admin email
	Email_str string `validate:"omitempty,email"`
}
type GF_admin__output_login struct {
	Email_confirmed_bool bool
	MFA_confirmed_bool   bool
	Pass_valid_bool      bool
	JWT_token_val        gf_session.GF_jwt_token_val
	User_id_str          gf_core.GF_ID 
}
type GF_admin__output_create_admin struct {
	General *GF_user_auth_userpass__output_create
}


type GF_admin__input_add_to_invite_list struct {
	User_name_str GF_user_name `validate:"required,min=3,max=50"`
	Email_str     string       `validate:"required,email"`
}

//------------------------------------------------
func Admin__pipeline__get_all_invite_list(p_ctx context.Context,
	p_service_info *GF_service_info,
	p_runtime_sys  *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.GF_error) {

	// DB
	db_invite_list_lst, gf_err := db__user__get_all_in_invite_list(p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	invite_list_lst := []map[string]interface{}{}
	for _, invite_map := range db_invite_list_lst {

		invite_list_lst = append(invite_list_lst, map[string]interface{}{
			"user_email_str":       invite_map["user_email_str"],
			"creation_unix_time_f": invite_map["creation_unix_time_f"],
		})
	}

	return invite_list_lst, nil
}

//------------------------------------------------
func Admin__pipeline__user_add_to_invite_list(p_input *GF_admin__input_add_to_invite_list,
	p_ctx          context.Context,
	p_service_info *GF_service_info,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.GF_error {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------------

	admin_user_name_str := p_input.User_name_str

	gf_err = db__user__add_to_invite_list(p_input.Email_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	// EVENT
	if p_service_info.Enable_events_app_bool {
		admin_user_id_str, gf_err := db__user__get_basic_info_by_username(admin_user_name_str,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		event_meta := map[string]interface{}{
			"user_id_str":                admin_user_id_str,
			"user_name_str":              admin_user_name_str,
			"email_added_to_invite_list": p_input.Email_str,
		}
		gf_events.Emit_app(GF_EVENT_APP__ADMIN_ADDED_USER_TO_INVITE_LIST,
			event_meta,
			p_runtime_sys)
	}

	return nil
}

//---------------------------------------------------
// PIPELINE__LOGIN

// this function is entered mutliple times for complex logins where not only pass/eth_signature
// are verified, but where email/mfa have to be confirmed as well.
// for each of the login stages this function is entered, and the login_attempt record
// is used to keep track of which stages have completed.

func Admin__pipeline__login(p_input *GF_admin__input_login,
	p_ctx          context.Context,
	p_local_hub    *sentry.Hub,
	p_service_info *GF_service_info,
	p_runtime_sys  *gf_core.Runtime_sys) (*GF_admin__output_login, *gf_core.GF_error) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_admin__output_login{}

	//------------------------
	// VERIFY

	user_exists_bool, gf_err := db__user__exists_by_username(GF_user_name(p_input.User_name_str),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	
	// BREADCRUMB
	gf_core.Breadcrumbs__add("auth", "admin user checked for existence",
		map[string]interface{}{"user_exists_bool": user_exists_bool, "user_name_str": p_input.User_name_str},
		p_local_hub)

	var user_id_str gf_core.GF_ID
	
	// user doesnt exist
	if !user_exists_bool {
		
		//------------------------	
		// PIPELINE__CREATE_ADMIN
		// if the admin user doesnt exist in the DB (most likely on first run of gloflow server),
		// create one in the DB

		input_create := &GF_user_auth_userpass__input_create{
			User_name_str: GF_user_name(p_input.User_name_str),
			Pass_str:      p_input.Pass_str,
			Email_str:     p_input.Email_str,
		}

		// BREADCRUMB
		gf_core.Breadcrumbs__add("auth", "creating new admin user",
			map[string]interface{}{"email_str": p_input.Email_str, "user_name_str": p_input.User_name_str},
			p_local_hub)
		
		output, gf_err := admin__pipeline__create_admin(input_create,
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

	// BREADCRUMB
	gf_core.Breadcrumbs__add("auth", "got user_id for admin user",
		map[string]interface{}{"user_id_str": user_id_str, "user_name_str": p_input.User_name_str},
		p_local_hub)

	//------------------------
	// LOGIN_ATTEMPT

	// get a preexisting login_attempt if one exists and hasnt expired for this user.
	// if it has then a new one will have to be created.
	var login_attempt *GF_login_attempt
	login_attempt, gf_err = login_attempt__get_if_valid(GF_user_name(p_input.User_name_str),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	if login_attempt == nil {

		//------------------------
		// CREATE_LOGIN_ATTEMPT

		user_identifier_str  := p_input.User_name_str
		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		login_attempt_id_str := users__create_id(user_identifier_str, creation_unix_time_f)

		login_attempt = &GF_login_attempt{
			V_str:                "0",
			Id_str:               login_attempt_id_str,
			Creation_unix_time_f: creation_unix_time_f,
			User_type_str:        "admin",
			User_name_str:        GF_user_name(p_input.User_name_str),
		}
		gf_err := db__login_attempt__create(login_attempt,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		//------------------------
	}
	
	//------------------------
	// VERIFY_PASSWORD

	// only verify password if the login_attempt didnt mark it yet as complete
	if !login_attempt.Pass_confirmed_bool {

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

			//------------------------
			// UPDATE_LOGIN_ATTEMPT
			// if password is valid then update the login_attempt 
			// to indicate that the password has been confirmed
			update_op := &GF_login_attempt__update_op{Pass_confirmed_bool: &pass_valid_bool}
			gf_err = db__login_attempt__update(&login_attempt.Id_str,
				update_op,
				p_ctx,
				p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			//------------------------

			// EVENT
			if p_service_info.Enable_events_app_bool {
				event_meta := map[string]interface{}{
					"user_id_str":     user_id_str,
					"user_name_str":   p_input.User_name_str,
					"domain_base_str": p_service_info.Domain_base_str,
				}
				gf_events.Emit_app(GF_EVENT_APP__ADMIN_LOGIN_PASS_CONFIRMED,
					event_meta,
					p_runtime_sys)
			}
		}
	}

	//------------------------
	// EMAIL
	if p_service_info.Enable_email_bool {

		// go through the email verification pipeline if the email
		// has not yet been confirmed
		if !login_attempt.Email_confirmed_bool {

			gf_err = users_email__verify__pipeline(p_input.Email_str,
				GF_user_name(p_input.User_name_str),
				user_id_str,
				p_service_info.Domain_base_str,
				p_ctx,
				p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			// EVENT
			if p_service_info.Enable_events_app_bool {
				event_meta := map[string]interface{}{
					"user_id_str":     user_id_str,
					"user_name_str":   p_input.User_name_str,
					"domain_base_str": p_service_info.Domain_base_str,
				}
				gf_events.Emit_app(GF_EVENT_APP__ADMIN_LOGIN_EMAIL_VERIFICATION_SENT,
					event_meta,
					p_runtime_sys)
			}

			//------------------------
		}
	}

	//------------------------
	
	return output, nil
}

//---------------------------------------------------
func admin__pipeline__create_admin(p_input *GF_user_auth_userpass__input_create,
	p_service_info *GF_service_info,
	p_ctx          context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) (*GF_admin__output_create_admin, *gf_core.GF_error) {

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
		gf_events.Emit_app(GF_EVENT_APP__ADMIN_CREATE,
			event_meta,
			p_runtime_sys)
	}

	//------------------------

	output_admin := &GF_admin__output_create_admin{
		General: output,
	}
	return output_admin, nil
}

//---------------------------------------------------
func Admin__is(p_user_name_str string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {
	if p_user_name_str != "admin" {
		gf_err := gf_core.Error__create("username thats not 'admin' is trying to login as admin",
			"verify__invalid_value_error",
			map[string]interface{}{
				"user_name_str": p_user_name_str,
			},
			nil, "gf_identity_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}