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

package gf_identity

import (
	"fmt"
	"testing"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func TesUsersHTTPuserpass(pTest *testing.T) {

	fmt.Println(" TEST_USERS_HTTP_USERPASS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	ctx := context.Background()

	serviceNameStr := "gf_identity_test"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"

	runtimeSys  := Tinit(serviceNameStr, mongoHostStr)

	//---------------------------------
	// START_SERVICE
	authSubsystemTypeStr := "userpass"
	portInt := 200

	templatesPathsMap := map[string]string {
		"gf_login": "./../../web/src/gf_identity/templates/gf_login/gf_login.html",
	}

	TestStartService(authSubsystemTypeStr,
		templatesPathsMap,
		portInt,
		runtimeSys)

	testPortInt := portInt

	//---------------------------------

	HTTPagent   := gorequest.New()

	testUserNameStr := "ivan_t"
	testUserPassStr := "pass_lksjds;lkdj"
	testEmailStr    := "ivan_t@gloflow.com"

	//---------------------------------
	// CLEANUP
	TestDBcleanup(ctx, runtimeSys)

	//---------------------------------
	// ADD_TO_INVITE_LIST
	gfErr := gf_identity_core.DBuserAddToInviteList(testEmailStr,
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

	// FIX!! - use testUserHTTPupdate() for this instead of calling endpoint explicitly.
	urlStr := fmt.Sprintf("http://localhost:%d/v1/identity/update?auth_r=0", testPortInt)
	data_map := map[string]string{
		"user_name_str":   "ivan_t_new",
		"email_str":       "ivan_t_new@gloflow.com",
		"description_str": "some new description",
	}
	dataBytesLst, _ := json.Marshal(data_map)
	_, bodyStr, errs := HTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}
	
	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
		pTest.FailNow()
	}

	spew.Dump(bodyMap)

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user updating http request failed")

	//---------------------------------
	// TEST_USER_GET_ME

	fmt.Println("====================================")
	fmt.Println("test user GET ME")
	
	// FIX!! - use testUserHTTPgetMe()
	urlStr = fmt.Sprintf("http://localhost:%d/v1/identity/me?auth_r=0", testPortInt)
	_, bodyStr, errs = HTTPagent.Get(urlStr).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}
	
	bodyMap = map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
        pTest.FailNow()
    }

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user get me http request failed")

	userNameStr        := bodyMap["data"].(map[string]interface{})["user_name_str"].(string)
	emailStr           := bodyMap["data"].(map[string]interface{})["email_str"].(string)
	descriptionStr     := bodyMap["data"].(map[string]interface{})["description_str"].(string)
	profileImageURLstr := bodyMap["data"].(map[string]interface{})["profile_image_url_str"].(string)
	bannerImageURLstr  := bodyMap["data"].(map[string]interface{})["banner_image_url_str"].(string)

	fmt.Println("====================================")
	fmt.Println("user login response:")
	fmt.Println("user_name_str",         userNameStr)
	fmt.Println("email_str",             emailStr)
	fmt.Println("description_str",       descriptionStr)
	fmt.Println("profile_image_url_str", profileImageURLstr)
	fmt.Println("banner_image_url_str",  bannerImageURLstr)

	//---------------------------------
}