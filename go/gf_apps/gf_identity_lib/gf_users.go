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


type GF_user__input_login struct {
}
type GF_user__output_login struct {
}

type GF_user__input_create struct {
	Auth_proof_sig_str string              `json:"auth_proof_sig_str"`
	Nonce_str          string              `json:"nonce_str"`
	Address_eth_str    GF_user_address_eth `json:"address_eth_str"`
}
type GF_user__output_create struct {
	Auth_proof_sig_valid_bool bool       `json:"auth_proof_sig_valid_bool"`
	JWT_token_val             GF_jtw_val `json:"jwt_token_val_str"`
}



type GF_user__input_update struct {
}
type GF_user__output_update struct {
}

type GF_user__update struct {
	Username_str    string
	Description_str string
}


type GF_user__input_get struct {
}
type GF_user__output_get struct {
}

//---------------------------------------------------
// PIPELINE__CREATE
func users__pipeline__login(p_input *GF_user__input_login,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_login, *gf_core.GF_error) {




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
	valid_bool, gf_err := verify__auth_proof_signature(p_input.Auth_proof_sig_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	
	if !valid_bool {
		output.Auth_proof_sig_valid_bool = false
		return output, nil
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
func users__pipeline__update(p_input *GF_user__input_update,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_update, *gf_core.GF_error) {

	output := &GF_user__output_update{}

	return output, nil
}

//---------------------------------------------------
func users__pipeline__get(p_input *GF_user__input_get,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__output_get, *gf_core.GF_error) {

	output := &GF_user__output_get{}

	return output, nil
}