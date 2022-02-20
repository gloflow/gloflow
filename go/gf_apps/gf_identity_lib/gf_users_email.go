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
	"context"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
)

//---------------------------------------------------
func users_email__verify__pipeline(p_email_address_str string,
	p_user_name_str   GF_user_name,
	p_user_id_str     gf_core.GF_ID,
	p_domain_base_str string,
	p_ctx             context.Context,
	p_runtime_sys     *gf_core.Runtime_sys) *gf_core.GF_error {
	
	//------------------------
	// EMAIL_VERIFY_ADDRESS
	// FIX!! - have a way of checking that an email has already been verified,
	//         instead of trying to verify each time.
	gf_err := gf_aws.AWS_SES__verify_address(p_email_address_str,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------------
	// EMAIL_CONFIRM

	confirm_code_str := users_email__generate_confirmation_code()

	// DB
	gf_err = db__user_email_confirm__create(p_user_name_str,
		p_user_id_str,
		confirm_code_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	msg_subject_str, msg_body_html_str, msg_body_text_str := users_email__get_confirm_msg_info(p_user_name_str,
		confirm_code_str,
		p_domain_base_str)

	// sender address
	sender_address_str := fmt.Sprintf("gf-email-confirm@%s", p_domain_base_str)

	gf_err = gf_aws.AWS_SES__send_message(p_email_address_str,
		sender_address_str,
		msg_subject_str,
		msg_body_html_str,
		msg_body_text_str,
		p_runtime_sys)
	
	if gf_err != nil {
		return gf_err
	}

	//------------------------


	return nil
}

//---------------------------------------------------
func users_email__confirm__pipeline(p_input *GF_user__http_input_email_confirm,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, string, *gf_core.GF_error) {

	db_confirm_code_str, gf_err := users_email__get_confirmation_code(p_input.User_name_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return false, "", gf_err
	}
	
	if db_confirm_code_str == "" {
		return false, "email confirmation code has expired", nil
	}

	// confirm_code is correct
	if p_input.Confirm_code_str == db_confirm_code_str {
		
		// GET_USER_ID
		user_id_str, gf_err := db__user__get_basic_info_by_username(p_input.User_name_str,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return false, "", gf_err
		}

		//------------------------
		// initial user email confirmation. only for new users.
		// user confirmed their email as valid.
		user_email_confirmed_bool, gf_err := db__user__email_is_confirmed(p_input.User_name_str, p_ctx, p_runtime_sys)
		if gf_err != nil {
			return false, "", gf_err
		}

		if user_email_confirmed_bool {
			update_op := &GF_user__update_op{
				Email_confirmed_bool: true,
			}
	
			// UPDATE_USER - mark user as email_confirmed
			gf_err = db__user__update(user_id_str,
				update_op,
				p_ctx,
				p_runtime_sys)
			if gf_err != nil {
				return false, "", gf_err
			}
		}

		//------------------------

		//------------------------
		// UPDATE_LOGIN_ATTEMPT
		// if email is confirmed then update the login_attempt

		// get a preexisting login_attempt if one exists and hasnt expired for this user.
		// if it has then a new one will have to be created.
		var login_attempt *GF_login_attempt
		login_attempt, gf_err = login_attempt__get_if_valid(GF_user_name(p_input.User_name_str),
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return false, "", gf_err
		}

		
		login_email_confirmed_bool := true
		update_op := &GF_login_attempt__update_op{Email_confirmed_bool: &login_email_confirmed_bool}
		gf_err = db__login_attempt__update(&login_attempt.Id_str,
			update_op,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return false, "", gf_err
		}

		//------------------------

		return true, "", nil

	} else {
		return false, "received confirm code and DB confirm code are not the same", nil
	}
	return false, "", nil
}

//---------------------------------------------------
func users_email__get_confirmation_code(p_user_name_str GF_user_name,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (string, *gf_core.GF_error) {

	confirm_code_str, confirm_code_creation_time_f, gf_err := db__user_email_confirm__get_code(p_user_name_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	//------------------------
	// check confirm_code didnt expire
	current_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	confirm_code_age_time_f := current_unix_time_f - confirm_code_creation_time_f

	// check if older than 5min
	if (5.0 < confirm_code_age_time_f/60) {
		return "", nil
	}

	//------------------------

	return confirm_code_str, nil
}

//---------------------------------------------------
func users_email__generate_confirmation_code() string {
	c_str := fmt.Sprintf("%s:%s", gf_core.Str_random(), gf_core.Str_random())
	return c_str
}

//---------------------------------------------------
func users_email__get_confirm_msg_info(p_user_name_str GF_user_name,
	p_confirm_code_str string,
	p_domain_str       string) (string, string, string) {

	subject_str := fmt.Sprintf("%s - confirm your email", p_domain_str)

	html_str := fmt.Sprintf(`
		<div>
			<div id='gf_logo'>
				<img src="https://gloflow.com/images/d/gf_logo_0.3.png"></img>
			</div>
			<div>
				<div>Welcome to %s!</div>
				<div>
				<div>
					"There is no spoon ...it is only yourself."
				</div>

				
			</div>
			<div>
				<div style="font-size:'14px';">Please click on the bellow link to confirm your email address.</div>
				<a href="https://%s/v1/identity/email_confirm?u=%s&c=%s">confirm email</a>
			</div>
			<div style="font-size:'8px';">
				don't reply to this email
			</div>
		</div>`,
		p_domain_str,
		p_domain_str,
		p_user_name_str,
		p_confirm_code_str)

	text_str := fmt.Sprintf(`
		Welcome to %s!
		There is no spoon. ...it is only yourself.

		Please open the following link in your browser to confirm your email address.
		
		https://%s/v1/identity/email_confirm?c=%s`,
		p_domain_str,
		p_domain_str,
		p_confirm_code_str)

	return subject_str, html_str, text_str
}