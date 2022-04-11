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
	"encoding/json"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/parnurzeal/gorequest"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// TEST_USER_HTTP_CREATE
func TestUserHTTPcreate(pTestUserNameStr string,
	pTestUserPassStr string,
	pTestEmailStr    string,
	pHTTPagent       *gorequest.SuperAgent,
	pTestPortInt     int,
	pTest            *testing.T) {

	fmt.Println("====================================")
	fmt.Println("test user CREATE USERPASS")
	fmt.Println("user_name_str", pTestUserNameStr)
	fmt.Println("pass_str",      pTestUserPassStr)
	fmt.Println("email_str",     pTestEmailStr)

	urlStr := fmt.Sprintf("http://localhost:%d/v1/identity/userpass/create", pTestPortInt)
	dataMap := map[string]string{
		"user_name_str": pTestUserNameStr,
		"pass_str":      pTestUserPassStr,
		"email_str":     pTestEmailStr,
	}
	dataBytesLst, _ := json.Marshal(dataMap)
	_, bodyStr, errs := pHTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	spew.Dump(bodyStr)

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}

	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
        pTest.FailNow()
    }

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user create http request failed")

	user_exists_bool         := bodyMap["data"].(map[string]interface{})["user_exists_bool"].(bool)
	user_in_invite_list_bool := bodyMap["data"].(map[string]interface{})["user_in_invite_list_bool"].(bool)

	if (user_exists_bool) {
		fmt.Println("supplied user already exists and cant be created")
		pTest.FailNow()
	}
	if (!user_in_invite_list_bool) {
		fmt.Println("supplied user is not in the invite list")
		pTest.FailNow()
	}
}

//-------------------------------------------------
func TestUserHTTPlogin(pTestUserNameStr string,
	pTestUserPassStr string,
	pHTTPagent       *gorequest.SuperAgent,
	pTestPortInt     int,
	pTest            *testing.T) {


	fmt.Println("====================================")
	fmt.Println("test user LOGIN USERPASS")

	urlStr  := fmt.Sprintf("http://localhost:%d/v1/identity/userpass/login", pTestPortInt)
	dataMap := map[string]string{
		"user_name_str": pTestUserNameStr,
		"pass_str":      pTestUserPassStr,
	}
	dataBytesLst, _ := json.Marshal(dataMap)
	resp, bodyStr, errs := pHTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
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
	assert.True(pTest, auth_cookie_present_bool,
		"login response does not contain the expected 'gf_sess_data' cookie")

	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
		pTest.FailNow()
	}

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user login http request failed")

	user_exists_bool := bodyMap["data"].(map[string]interface{})["user_exists_bool"].(bool)
	pass_valid_bool  := bodyMap["data"].(map[string]interface{})["pass_valid_bool"].(bool)
	user_id_str      := bodyMap["data"].(map[string]interface{})["user_id_str"].(string)

	assert.True(pTest, user_id_str != "", "user_id not set in the response")

	fmt.Println("user login response:")
	fmt.Println("user_exists_bool", user_exists_bool)
	fmt.Println("pass_valid_bool",  pass_valid_bool)
	fmt.Println("user_id_str",      user_id_str)
}

//-------------------------------------------------
func test_user_http_update(p_test *testing.T,
	p_http_agent    *gorequest.SuperAgent,
	p_test_port_int int) {

	fmt.Println("====================================")
	fmt.Println("test user UPDATE")

	url_str := fmt.Sprintf("http://localhost:%d/v1/identity/update", p_test_port_int)
	data_map := map[string]string{
		"user_name_str":   "new username",
		"email_str":       "ivan@gloflow.com",
		"description_str": "some new description",
	}
	data_bytes_lst, _ := json.Marshal(data_map)
	_, body_str, errs := p_http_agent.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	if len(errs) > 0 {
		spew.Dump(errs)
		p_test.Fail()
	}

	body_map := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
		p_test.Fail()
	}

	spew.Dump(body_map)

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user updating http request failed")
}

//-------------------------------------------------
// TEST_USER_GET_ME
func test_user_http_get_me(p_test *testing.T,
	p_http_agent    *gorequest.SuperAgent,
	p_test_port_int int) {

	url_str := fmt.Sprintf("http://localhost:%d/v1/identity/me", p_test_port_int)
	_, body_str, errs := p_http_agent.Get(url_str).
		End()

	if len(errs) > 0 {
		spew.Dump(errs)
		p_test.Fail()
	}
	
	body_map := map[string]interface{}{}
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

	fmt.Println("RESPONSE >>>>")
	spew.Dump(body_map["data"])
	
	fmt.Println("====================================")
	fmt.Println("user me response:")
	fmt.Println("user_name_str",         user_name_str)
	fmt.Println("email_str",             email_str)
	fmt.Println("description_str",       description_str)
	fmt.Println("profile_image_url_str", profile_image_url_str)
	fmt.Println("banner_image_url_str",  banner_image_url_str)
}