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
	"fmt"
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Test__signing(p_test *testing.T) {
	fmt.Println(" TEST__IDENTITY_SIGNING >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")





}

//-------------------------------------------------
func Test__users(p_test *testing.T) {

	

	fmt.Println(" TEST__IDENTITY_USERS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	runtime_sys := T__init()

	test_user_address_eth_str := "0xBA47Bef4ca9e8F86149D2f109478c6bd8A642C97"
	test_user_signature_str   := "0x07c582de2c6fb11310495815c993fa978540f0c0cdc89fd51e6fe3b8db62e913168d9706f32409f949608bcfd372d41cbea6eb75869afe2f189738b7fb764ef91c"
	test_user_nonce_str       := "gf_test_message_to_sign"
	ctx := context.Background()




	//------------------
	// NONCE_CREATE

	unexisting_user_id_str := gf_core.GF_ID("")
	_, gf_err := nonce__create(GF_user_nonce_val(test_user_nonce_str),
		unexisting_user_id_str,
		GF_user_address_eth(test_user_address_eth_str),
		ctx,
		runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	//------------------
	// USER_CREATE
	
	input__create := &GF_user__input_create{
		Auth_signature_str:   GF_auth_signature(test_user_signature_str),
		User_address_eth_str: GF_user_address_eth(test_user_address_eth_str),
		// Nonce_val_str:   nonce.Val_str,
	}


	output__create, gf_err := users__pipeline__create(input__create, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}


	spew.Dump(output__create)


	assert.True(p_test, output__create.Auth_signature_valid_bool, "crypto signature supplied for user creation pipeline is invalid")


	//------------------


	input__login := &GF_user__input_login{
		Auth_signature_str:   GF_auth_signature(test_user_signature_str),
		User_address_eth_str: GF_user_address_eth(test_user_address_eth_str),
	}
	output__login, gf_err := users__pipeline__login(input__login, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	spew.Dump(output__login)
	
	//------------------

}