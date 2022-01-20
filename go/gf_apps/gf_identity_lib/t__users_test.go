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
	"time"
	"strings"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_crypto"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Test__users_http(p_test *testing.T) {

	fmt.Println(" TEST__IDENTITY_USERS_HTTP >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	runtime_sys := T__init()
	test_port_int := 2000
	go func() {
		Init_service(runtime_sys)
		gf_rpc_lib.Server__init(test_port_int)
	}()
	time.Sleep(2*time.Second) // let server startup

	request := gorequest.New()

	//---------------------------------
	// GENERATE_WALLET
	private_key_hex_str, public_key_hex_str, address_str, err := gf_crypto.Eth_generate_keys()
	if err != nil {
		fmt.Println(err)
		p_test.Fail()
	}

	//---------------------------------
	// TEST_PREFLIGHT_HTTP
	
	data_map := map[string]string{
		"user_address_eth_str": address_str,
	}
	data_bytes_lst, _ := json.Marshal(data_map)
	url_str := fmt.Sprintf("http://localhost:%d/v1/identity/eth/preflight", test_port_int)
	_, body_str, errs := request.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	spew.Dump(body_str)

	if (len(errs) > 0) {
		fmt.Println(errs)
		p_test.Fail()
	}

	body_map := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        p_test.Fail()
    }

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user preflight http request failed")

	nonce_val_str    := body_map["data"].(map[string]interface{})["nonce_val_str"].(string)
	user_exists_bool := body_map["data"].(map[string]interface{})["user_exists_bool"].(bool)
	
	fmt.Println("====================================")
	fmt.Println("preflight response:")
	fmt.Println("nonce_val_str",    nonce_val_str)
	fmt.Println("user_exists_bool", user_exists_bool)

	// we're testing user creation and the rest of the flow as well, so user shouldnt exist
	if (user_exists_bool) {
		p_test.Fail()
	}
	
	//---------------------------------
	// TEST_USER_CREATE_HTTP

	signature_str, err := gf_crypto.Eth_sign_data(nonce_val_str, private_key_hex_str)
	if err != nil {
		fmt.Println(err)
		p_test.Fail()
	}

	fmt.Println("====================================")
	fmt.Println("user create inputs:")
	fmt.Println("address",   address_str, len(address_str))
	fmt.Println("priv key",  private_key_hex_str, len(private_key_hex_str))
	fmt.Println("pub key",   public_key_hex_str, len(public_key_hex_str))
	fmt.Println("signature", signature_str, len(signature_str))
	fmt.Println("nonce",     nonce_val_str)

	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/eth/create", test_port_int)
	data_map = map[string]string{
		"user_address_eth_str": address_str,
		"auth_signature_str":   signature_str,
	}
	data_bytes_lst, _ = json.Marshal(data_map)
	_, body_str, errs = request.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	spew.Dump(body_str)

	if (len(errs) > 0) {
		fmt.Println(errs)
		p_test.Fail()
	}

	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        p_test.Fail()
    }

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user create http request failed")

	nonce_exists_bool         := body_map["data"].(map[string]interface{})["nonce_exists_bool"].(bool)
	auth_signature_valid_bool := body_map["data"].(map[string]interface{})["auth_signature_valid_bool"].(bool)

	if (!nonce_exists_bool) {
		fmt.Println("supplied nonce doesnt exist")
		p_test.Fail()
	}
	if (!auth_signature_valid_bool) {
		fmt.Println("signature is not valid")
		p_test.Fail()
	}

	//---------------------------------
	// TEST_USER_LOGIN

	fmt.Println("====================================")
	fmt.Println("user login inputs:")
	fmt.Println("address",   address_str, len(address_str))
	fmt.Println("signature", signature_str, len(signature_str))

	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/eth/login", test_port_int)
	data_map = map[string]string{
		"user_address_eth_str": address_str,
		"auth_signature_str":   signature_str,
	}
	data_bytes_lst, _ = json.Marshal(data_map)
	resp, body_str, errs := request.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		p_test.Fail()
	}

	// check if the login response sets a cookie for all future auth requests
	auth_cookie_present_bool := false
	for k, v := range resp.Header {
		if (k == "Set-Cookie") {
			for _, vv := range v {
				o := strings.Split(vv, "=")[0]
				if o == "gf_sess_data" {
					auth_cookie_present_bool = true
				}
			}
		}
	}
	assert.True(p_test, auth_cookie_present_bool,
		"login response does not contain the expected 'gf_sess_data' cookie")

	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        p_test.Fail()
    }

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user login http request failed")

	nonce_exists_bool         = body_map["data"].(map[string]interface{})["nonce_exists_bool"].(bool)
	auth_signature_valid_bool = body_map["data"].(map[string]interface{})["auth_signature_valid_bool"].(bool)
	user_id_str              := body_map["data"].(map[string]interface{})["user_id_str"].(string)

	fmt.Println("====================================")
	fmt.Println("user login response:")
	fmt.Println("nonce_exists_bool",         nonce_exists_bool)
	fmt.Println("auth_signature_valid_bool", auth_signature_valid_bool)
	fmt.Println("user_id_str",               user_id_str)

	if (!nonce_exists_bool) {
		fmt.Println("supplied nonce doesnt exist")
		p_test.Fail()
	}
	if (!auth_signature_valid_bool) {
		fmt.Println("signature is not valid")
		p_test.Fail()
	}

	//---------------------------------
	// TEST_USER_UPDATE

	fmt.Println("====================================")
	fmt.Println("user update inputs:")
	fmt.Println("address",   address_str, len(address_str))
	fmt.Println("signature", signature_str, len(signature_str))

	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/update", test_port_int)
	data_map = map[string]string{
		"user_name_str":   "new username",
		"email_str":       "ivan@gloflow.com",
		"description_str": "some new description",
	}
	data_bytes_lst, _ = json.Marshal(data_map)
	_, body_str, errs = request.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
		p_test.Fail()
	}

	spew.Dump(body_map)

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user updating http request failed")

	//---------------------------------
	// TEST_USER_GET_ME

	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/me", test_port_int)
	data_bytes_lst, _ = json.Marshal(data_map)
	_, body_str, errs = request.Get(url_str).
		End()

	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        p_test.Fail()
    }

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user get me http request failed")

	user_name_str         := body_map["data"].(map[string]interface{})["user_name_str"].(string)
	email_str             := body_map["data"].(map[string]interface{})["email_str"].(string)
	description_str       := body_map["data"].(map[string]interface{})["description_str"].(string)
	profile_image_url_str := body_map["data"].(map[string]interface{})["profile_image_url_str"].(string)
	banner_image_url_str  := body_map["data"].(map[string]interface{})["banner_image_url_str"].(string)

	fmt.Println("====================================")
	fmt.Println("user login response:")
	fmt.Println("user_name_str",         user_name_str)
	fmt.Println("email_str",             email_str)
	fmt.Println("description_str",       description_str)
	fmt.Println("profile_image_url_str", profile_image_url_str)
	fmt.Println("banner_image_url_str",  banner_image_url_str)

	//---------------------------------
}

//-------------------------------------------------
func Test__users_unit(p_test *testing.T) {

	fmt.Println(" TEST__IDENTITY_USERS_UNIT >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

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
	
	input__create := &GF_user_auth_eth__input_create{
		Auth_signature_str:   GF_auth_signature(test_user_signature_str),
		User_address_eth_str: GF_user_address_eth(test_user_address_eth_str),
		// Nonce_val_str:   nonce.Val_str,
	}

	output__create, gf_err := users_auth_eth__pipeline__create(input__create, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	spew.Dump(output__create)

	assert.True(p_test, output__create.Auth_signature_valid_bool, "crypto signature supplied for user creation pipeline is invalid")

	//------------------
	input__login := &GF_user_auth_eth__input_login{
		Auth_signature_str:   GF_auth_signature(test_user_signature_str),
		User_address_eth_str: GF_user_address_eth(test_user_address_eth_str),
	}
	output__login, gf_err := users_auth_eth__pipeline__login(input__login, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	spew.Dump(output__login)
	
	//------------------
}