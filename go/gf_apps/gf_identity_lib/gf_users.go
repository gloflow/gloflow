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
	// "fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GF_auth_signature   string
type GF_user_address_eth string

type GF_user struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`

	Username_str      string                `bson:"username_str"`   // set once at the creation of the user
	Screenname_str    string                `bson:"screenname_str"` // changable durring the lifetime of the user
	Email_str         string                `bson:"email_str"`
	Description_str   string                `bson:"description_str"`
	Addresses_eth_lst []GF_user_address_eth `bson:"addresses_eth_lst"`

	// IMAGES
	Profile_image_url_str string `bson:"profile_image_url_str"`
	Banner_image_url_str  string `bson:"banner_image_url_str"`
}

// io_preflight
type GF_user__input_preflight struct {
	User_name_str        string              `validate:"omitempty,min=3,max=50"`
	User_address_eth_str GF_user_address_eth `validate:"omitempty,eth_addr"`
}
type GF_user__output_preflight struct {
	User_exists_bool bool             
	Nonce_val_str    GF_user_nonce_val
}

// io_login
type GF_user__input_login struct {
	User_address_eth_str GF_user_address_eth `validate:"required,eth_addr"`
	Auth_signature_str   GF_auth_signature   `validate:"required,len=132"` // singature length with "0x"
}
type GF_user__output_login struct {
	Nonce_exists_bool         bool
	Auth_signature_valid_bool bool
	JWT_token_val             GF_jwt_token_val
	User_id_str               gf_core.GF_ID 
}

// io_create
type GF_user__input_create struct {
	User_address_eth_str GF_user_address_eth `validate:"required,eth_addr"`
	Auth_signature_str   GF_auth_signature   `validate:"required,len=132"` // singature length with "0x"
}
type GF_user__output_create struct {
	Nonce_exists_bool         bool
	Auth_signature_valid_bool bool
}

// io_update
type GF_user__input_update struct {
	User_address_eth_str GF_user_address_eth `validate:"required,eth_addr"`
	Username_str         *string             `validate:"omitempty,min=3,max=50"`   // optional
	Screenname_str       *string             `validate:"omitempty,min=3,max=50"`   // optional
	Email_str            *string             `validate:"omitempty,min=6,max=50"`   // optional
	Description_str      *string             `validate:"omitempty,min=1,max=2000"` // optional

	Profile_image_url_str *string `validate:"omitempty,min=1,max=100"` // optional // FIX!! - validation
	Banner_image_url_str  *string `validate:"omitempty,min=1,max=100"` // optional // FIX!! - validation
}
type GF_user__output_update struct {
	
}

type GF_user__update struct {
	Username_str    string
	Description_str string
}

// io_get
type GF_user__input_get struct {
	User_address_eth_str GF_user_address_eth `validate:"required,eth_addr"`
}

type GF_user__output_get struct {
	Username_str    string
	Email_str       string
	Description_str string
	Profile_image_url_str string
	Banner_image_url_str  string
}

//---------------------------------------------------
func users__pipeline__preflight(p_input *GF_user__input_preflight,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_preflight, *gf_core.GF_error) {

	//------------------------
	// VALIDATE
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user__output_preflight{}

	exists_bool, gf_err := db__user__exists(p_input.User_address_eth_str,
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
func users__pipeline__login(p_input *GF_user__input_login,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_login, *gf_core.GF_error) {
	
	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------


	output := &GF_user__output_login{}

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
	// JWT
	user_identifier_str := string(p_input.User_address_eth_str)
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
func users__pipeline__create(p_input *GF_user__input_create,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_create, *gf_core.GF_error) {

	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user__output_create{}
	
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

	creation_unix_time_f  := float64(time.Now().UnixNano())/1000000000.0
	user_address_eth_str  := p_input.User_address_eth_str
	user_ddresses_eth_lst := []GF_user_address_eth{user_address_eth_str, }

	user_id := users__create_id(user_address_eth_str, creation_unix_time_f)

	user := &GF_user{
		V_str:                "0",
		Id_str:               user_id,
		Creation_unix_time_f: creation_unix_time_f,
		Addresses_eth_lst:    user_ddresses_eth_lst,
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

//---------------------------------------------------
// PIPELINE__UPDATE
func users__pipeline__update(p_input *GF_user__input_update,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_update, *gf_core.GF_error) {
	
	//------------------------
	// VALIDATE_INPUT
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------

	output := &GF_user__output_update{}

	return output, nil
}

//---------------------------------------------------
// PIPELINE__GET
func users__pipeline__get(p_input *GF_user__input_get,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_get, *gf_core.GF_error) {

	//------------------------
	// VALIDATE
	gf_err := gf_core.Validate_struct(p_input, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------
	
	user, gf_err := db__user__get(p_input.User_address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}


	output := &GF_user__output_get{
		Username_str:    user.Username_str,
		Email_str:       user.Email_str,
		Description_str: user.Description_str,
		Profile_image_url_str: user.Profile_image_url_str,
		Banner_image_url_str:  user.Banner_image_url_str,
	}

	return output, nil
}