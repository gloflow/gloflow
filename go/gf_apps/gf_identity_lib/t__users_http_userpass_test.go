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
	"testing"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/parnurzeal/gorequest"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Test__users_http_userpass(pTest *testing.T) {

	fmt.Println(" TEST__IDENTITY_USERS_HTTP_USERPASS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	testPortInt := 2000
	ctx         := context.Background()
	runtimeSys  := T__init()
	HTTPagent   := gorequest.New()

	testUserNameStr := "ivan_t"
	testUserPassStr := "pass_lksjds;lkdj"
	testEmailStr    := "ivan_t@gloflow.com"

	//---------------------------------
	// CLEANUP
	TestDBcleanup(ctx, runtimeSys)

	//---------------------------------
	// ADD_TO_INVITE_LIST
	gfErr := DBuserAddToInviteList(testEmailStr,
		ctx,
		runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.FailNow()
	}
	
	//---------------------------------
	// TEST_USER_CREATE_HTTP

	TestUserHTTPcreate(testUserNameStr,
		testUserPassStr,
		testEmailStr,
		HTTPagent,
		testPortInt,
		pTest)

	//---------------------------------
	// TEST_USER_LOGIN

	TestUserHTTPlogin(testUserNameStr,
		testUserPassStr,
		HTTPagent,
		testPortInt,
		pTest)

	//---------------------------------
	// TEST_USER_UPDATE

	fmt.Println("====================================")
	fmt.Println("test user UPDATE")

	url_str := fmt.Sprintf("http://localhost:%d/v1/identity/update", testPortInt)
	data_map := map[string]string{
		"user_name_str":   "ivan_t_new",
		"email_str":       "ivan_t_new@gloflow.com",
		"description_str": "some new description",
	}
	data_bytes_lst, _ := json.Marshal(data_map)
	_, body_str, errs := HTTPagent.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}
	
	body_map := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
		pTest.FailNow()
	}

	spew.Dump(body_map)

	assert.True(pTest, body_map["status"].(string) != "ERROR", "user updating http request failed")

	//---------------------------------
	// TEST_USER_GET_ME

	fmt.Println("====================================")
	fmt.Println("test user GET ME")
	
	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/me", testPortInt)
	data_bytes_lst, _ = json.Marshal(data_map)
	_, body_str, errs = HTTPagent.Get(url_str).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}
	
	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        pTest.FailNow()
    }

	assert.True(pTest, body_map["status"].(string) != "ERROR", "user get me http request failed")

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