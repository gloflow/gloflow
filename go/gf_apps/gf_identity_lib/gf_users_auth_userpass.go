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
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// io_login
type GF_user_auth_userpass__input_login struct {
	User_name_str string `validate:"omitempty,min=3,max=50"`
	Pass_hash_str string `validate:"required,len=132"` // FIX!! - whats the expected pass hash length?
}
type GF_user_auth_userpass__output_login struct {
	User_exists_bool     bool
	Pass_hash_valid_bool bool
	JWT_token_val        GF_jwt_token_val
	User_id_str          gf_core.GF_ID 
}

// io_create
type GF_user_auth_userpass__input_create struct {
	User_name_str string `validate:"omitempty,min=3,max=50"`
	Pass_hash_str string
	Email_str     string
}
type GF_user_auth_userpass__output_create struct {
	User_exists_bool bool
}

//---------------------------------------------------
// PIPELINE__LOGIN
func users_auth_userpass__pipeline__login(p_input *GF_user_auth_userpass__input_login,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_login, *gf_core.GF_error) {
	
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

	

	//------------------------
	// JWT
	user_identifier_str := string(p_input.User_name_str)
	jwt_token_val, gf_err := jwt__pipeline__generate(user_identifier_str, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output.JWT_token_val = jwt_token_val

	//------------------------
	// USER_ID

	//------------------------

	return output, nil
}

//---------------------------------------------------
// PIPELINE__CREATE
func users_auth_userpass__pipeline__create(p_input *GF_user_auth_userpass__input_create,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user_auth_userpass__output_create, *gf_core.GF_error) {

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

	//------------------------

	creation_unix_time_f  := float64(time.Now().UnixNano())/1000000000.0
	user_name_str := p_input.User_name_str
	pass_hash_str := p_input.Pass_hash_str
	email_str     := p_input.Email_str

	user_identifier_str := user_name_str
	user_id := users__create_id(user_identifier_str, creation_unix_time_f)

	user := &GF_user{
		V_str:                "0",
		Id_str:               user_id,
		Creation_unix_time_f: creation_unix_time_f,
		User_name_str:        user_name_str,
		Pass_hash_str:        pass_hash_str,
		Email_str:            email_str,
	}

	//------------------------
	// DB
	gf_err = db__user__create(user, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	return output, nil
}