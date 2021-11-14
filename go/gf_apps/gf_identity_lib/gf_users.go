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

	Username_str      string                `bson:"username_str"`
	Description_str   string                `bson:"description_str"`
	Addresses_eth_lst []GF_user_address_eth `bson:"addresses_eth_lst"`
}

// io_preflight
type GF_user__input_preflight struct {
	Address_eth_str GF_user_address_eth `json:"address_eth_str"`
}
type GF_user__output_preflight struct {
	User_exists_bool bool              `json:"user_exists_bool"`
	Nonce_val_str    GF_user_nonce_val `json:"nonce_val_str"`
}

// io_login
type GF_user__input_login struct {
	Signature_str   GF_auth_signature   `json:"signature_str"`
	Address_eth_str GF_user_address_eth `json:"address_eth_str"`
}
type GF_user__output_login struct {
	Signature_valid_bool bool             `json:"signature_valid_bool"`
	JWT_token_val        GF_jwt_token_val `json:"jwt_token_val_str"`
	User_id_str          gf_core.GF_ID    `json:"user_id_str"`
}

// io_create
type GF_user__input_create struct {
	Signature_str   GF_auth_signature   `json:"signature_str"`
	Nonce_val_str   GF_user_nonce_val   `json:"nonce_val_str"`
	Address_eth_str GF_user_address_eth `json:"address_eth_str"`
}
type GF_user__output_create struct {
	Signature_valid_bool bool             `json:"signature_valid_bool"`
	JWT_token_val        GF_jwt_token_val `json:"jwt_token_val_str"`
}

// io_update
type GF_user__input_update struct {
}
type GF_user__output_update struct {
}

type GF_user__update struct {
	Username_str    string
	Description_str string
}

// io_get
type GF_user__input_get struct {
}
type GF_user__output_get struct {
}

//---------------------------------------------------
func users__pipeline__preflight(p_input *GF_user__input_preflight,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_preflight, *gf_core.GF_error) {

	output := &GF_user__output_preflight{}





	exists_bool, gf_err := db__user__exists(p_input.Address_eth_str,
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
			p_input.Address_eth_str,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		output.User_exists_bool = false
		output.Nonce_val_str    = nonce.Val_str

	} else {



		nonce_val_str, gf_err := db__nonce__get(p_input.Address_eth_str, p_ctx, p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		output.User_exists_bool = true
		output.Nonce_val_str    = nonce_val_str
	}

	


	return output, nil
}

//---------------------------------------------------
// PIPELINE__LOGIN
func users__pipeline__login(p_input *GF_user__input_login,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_login, *gf_core.GF_error) {

	output := &GF_user__output_login{}





	user_nonce_val, gf_err := db__nonce__get(p_input.Address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------
	// VERIFY
	valid_bool, gf_err := verify__auth_signature__all_methods(p_input.Signature_str,
		user_nonce_val,
		p_input.Address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	
	if !valid_bool {
		output.Signature_valid_bool = false
		return output, nil
	} else {
		output.Signature_valid_bool = true
	}

	//------------------------
	// JWT
	jwt_token_val, gf_err := jwt__pipeline__generate(p_input.Address_eth_str, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output.JWT_token_val = jwt_token_val

	//------------------------



	return nil, nil
}

//---------------------------------------------------
// PIPELINE__CREATE
func users__pipeline__create(p_input *GF_user__input_create,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_create, *gf_core.GF_error) {

	output := &GF_user__output_create{}

	//------------------------
	// VERIFY
	valid_bool, gf_err := verify__auth_signature__all_methods(p_input.Signature_str,
		p_input.Nonce_val_str,
		p_input.Address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	
	if !valid_bool {
		output.Signature_valid_bool = false
		return output, nil
	} else {
		output.Signature_valid_bool = true
	}

	//------------------------


	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	address_eth_str      := p_input.Address_eth_str
	addresses_eth_lst    := []GF_user_address_eth{address_eth_str, }

	user_id := users__create_id(address_eth_str, creation_unix_time_f)

	user := &GF_user{
		V_str:                "0",
		Id_str:               user_id,
		Creation_unix_time_f: creation_unix_time_f,
		Addresses_eth_lst:    addresses_eth_lst, 
	}


	//------------------------
	// DB
	gf_err = db__user__create(user, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------
	// JWT
	jwt_token_val, gf_err := jwt__pipeline__generate(address_eth_str, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	output.JWT_token_val = jwt_token_val

	//------------------------

	return output, nil
}

//---------------------------------------------------
// PIPELINE__UPDATE
func users__pipeline__update(p_input *GF_user__input_update,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_update, *gf_core.GF_error) {

	output := &GF_user__output_update{}

	return output, nil
}

//---------------------------------------------------
// PIPELINE__GET
func users__pipeline__get(p_input *GF_user__input_get,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_get, *gf_core.GF_error) {

	output := &GF_user__output_get{}

	return output, nil
}