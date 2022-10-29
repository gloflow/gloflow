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
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
)

//---------------------------------------------------
// io_preflight
type GF_user_auth_eth__input_preflight struct {
	User_address_eth_str gf_identity_core.GF_user_address_eth `validate:"omitempty,eth_addr"`
}
type GF_user_auth_eth__output_preflight struct {
	User_exists_bool bool             
	Nonce_val_str    GF_user_nonce_val
}

// io_login
type GF_user_auth_eth__input_login struct {
	User_address_eth_str gf_identity_core.GF_user_address_eth `validate:"required,eth_addr"`
	Auth_signature_str   gf_identity_core.GF_auth_signature   `validate:"required,len=132"` // singature length with "0x"
}
type GF_user_auth_eth__output_login struct {
	Nonce_exists_bool         bool
	Auth_signature_valid_bool bool
	JWT_token_val             gf_session.GF_jwt_token_val
	User_id_str               gf_core.GF_ID 
}

// io_create
type GF_user_auth_eth__input_create struct {
	UserTypeStr          string                               `validate:"required"` // "admin" | "standard"
	User_address_eth_str gf_identity_core.GF_user_address_eth `validate:"required,eth_addr"`
	Auth_signature_str   gf_identity_core.GF_auth_signature   `validate:"required,len=132"` // singature length with "0x"
}
type GF_user_auth_eth__output_create struct {
	Nonce_exists_bool         bool
	Auth_signature_valid_bool bool
}

//---------------------------------------------------
func users_auth_eth__pipeline__preflight(p_input *GF_user_auth_eth__input_preflight,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) (*GF_user_auth_eth__output_preflight, *gf_core.GFerror) {

	//------------------------
	// VALIDATE
	gf_err := gf_core.ValidateStruct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user_auth_eth__output_preflight{}

	exists_bool, gf_err := db__user__exists_by_eth_addr(p_input.User_address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// no user exists so create a new nonce
	if !exists_bool {

		// user doesnt exist yet so no user_id
		user_id_str := gf_core.GF_ID("")
		nonce, gf_err := nonce__create_and_persist(user_id_str,
			p_input.User_address_eth_str,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		output.User_exists_bool = false
		output.Nonce_val_str    = nonce.Val_str

	// user exists
	} else {

		nonce_val_str, nonce_exists_bool, gf_err := db__nonce__get(p_input.User_address_eth_str, p_ctx, p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		if !nonce_exists_bool {
			// generate new nonce, because the old one has been invalidated?
		} else {
			output.User_exists_bool = true
			output.Nonce_val_str    = nonce_val_str
		}
	}

	return output, nil
}

//---------------------------------------------------
// PIPELINE__LOGIN
func users_auth_eth__pipeline__login(p_input *GF_user_auth_eth__input_login,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) (*GF_user_auth_eth__output_login, *gf_core.GFerror) {
	
	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.ValidateStruct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user_auth_eth__output_login{}

	//------------------------
	user_nonce_val, user_nonce_exists_bool, gf_err := db__nonce__get(p_input.User_address_eth_str,
		p_ctx,
		p_runtime_sys)
		
	if gf_err != nil {
		return nil, gf_err
	}
	
	if !user_nonce_exists_bool {
		output.Nonce_exists_bool = false
		return output, nil
	} else {
		output.Nonce_exists_bool = true
	}

	//------------------------
	// VERIFY

	signature_valid_bool, gf_err := verify__auth_signature__all_methods(p_input.Auth_signature_str,
		user_nonce_val,
		p_input.User_address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	if !signature_valid_bool {
		output.Auth_signature_valid_bool = false
		return output, nil
	} else {
		output.Auth_signature_valid_bool = true
	}

	//------------------------
	// USER_ID

	user_id_str, gf_err := gf_identity_core.DBgetBasicInfoByETHaddr(p_input.User_address_eth_str,
		p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output.User_id_str = user_id_str

	//------------------------
	// JWT
	user_identifier_str := string(user_id_str)
	jwt_token_val, gf_err := gf_session.JWT__pipeline__generate(user_identifier_str, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output.JWT_token_val = jwt_token_val

	//------------------------

	return output, nil
}

//---------------------------------------------------
// PIPELINE__CREATE
func users_auth_eth__pipeline__create(p_input *GF_user_auth_eth__input_create,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) (*GF_user_auth_eth__output_create, *gf_core.GFerror) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.ValidateStruct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user_auth_eth__output_create{}
	
	//------------------------
	// DB_NONCE_GET - get a nonce already generated in preflight for this user address,
	//                for validating the recevied auth_signature
	user_nonce_val_str, user_nonce_exists_bool, gf_err := db__nonce__get(p_input.User_address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	if !user_nonce_exists_bool {
		output.Nonce_exists_bool = false
		return output, nil
	} else {
		output.Nonce_exists_bool = true
	}

	//------------------------
	// VALIDATE

	signature_valid_bool, gf_err := verify__auth_signature__all_methods(p_input.Auth_signature_str,
		user_nonce_val_str,
		p_input.User_address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	
	if signature_valid_bool {
		output.Auth_signature_valid_bool = true
	} else {
		output.Auth_signature_valid_bool = false
		return output, nil
	}

	//------------------------

	creation_unix_time_f   := float64(time.Now().UnixNano())/1000000000.0
	user_address_eth_str   := p_input.User_address_eth_str
	user_addresses_eth_lst := []gf_identity_core.GF_user_address_eth{user_address_eth_str, }

	user_identifier_str := string(user_address_eth_str)
	user_id := usersCreateID(user_identifier_str, creation_unix_time_f)

	user := &GFuser{
		V_str:                "0",
		Id_str:               user_id,
		Creation_unix_time_f: creation_unix_time_f,
		UserTypeStr:          p_input.UserTypeStr,
		Addresses_eth_lst:    user_addresses_eth_lst,
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