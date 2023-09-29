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

package gf_identity

import (
	"fmt"
	"testing"
	"encoding/json"
	"strings"
	"net/http"
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

	spew.Dump(bodyMap)
	
	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user create http request failed")

	userExistsBool       := bodyMap["data"].(map[string]interface{})["user_exists_bool"].(bool)
	userInInviteListBool := bodyMap["data"].(map[string]interface{})["user_in_invite_list_bool"].(bool)

	if (userExistsBool) {
		fmt.Println("supplied user already exists and cant be created")
		pTest.FailNow()
	}
	if (!userInInviteListBool) {
		fmt.Println("supplied user is not in the invite list")
		pTest.FailNow()
	}
}

//-------------------------------------------------

func TestUserHTTPuserpassLogin(pTestUserNameStr string,
	pTestUserPassStr string,
	pHTTPagent       *gorequest.SuperAgent,
	pTestPortInt     int,
	pTest            *testing.T) []*http.Cookie {


	fmt.Println("====================================")
	fmt.Println("test user LOGIN USERPASS")

	urlStr := fmt.Sprintf("http://localhost:%d/v1/identity/userpass/login",
		pTestPortInt)

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

	//-----------------------------------
	// COOKIES
	// check if the login response sets a cookie for all future auth requests
	sessionIDcookiePresentBool := false
	authCookiePresentBool := false
	for k, v := range resp.Header {
		if (k == "Set-Cookie") {

			for _, vv := range v {
				o := strings.Split(vv, "=")[0]
				if o == "gf_sess" {
					sessionIDcookiePresentBool = true
				}
				if o == "Authorization" {
					authCookiePresentBool = true
				}
			}
		}
	}
	cookiesInRespLst := (*http.Response)(resp).Cookies()

    // Print all of the current cookies
    fmt.Println("RESPONSE COOKIES =====================================:")
    for _, cookie := range cookiesInRespLst {
        fmt.Printf("%s=%s\n", cookie.Name, cookie.Value)
    }

	assert.True(pTest, sessionIDcookiePresentBool,
		"login response does not contain the expected 'gf_sess' cookie")
	assert.True(pTest, authCookiePresentBool,
		"login response does not contain the expected 'Authorization' cookie")

	//-----------------------------------

	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
		pTest.FailNow()
	}

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user login http request failed")

	userExistsBool := bodyMap["data"].(map[string]interface{})["user_exists_bool"].(bool)
	passValidBool  := bodyMap["data"].(map[string]interface{})["pass_valid_bool"].(bool)
	userIDstr      := bodyMap["data"].(map[string]interface{})["user_id_str"].(string)

	assert.True(pTest, userIDstr != "", "user_id not set in the response")

	fmt.Println("user login response:")
	fmt.Println("user_exists_bool", userExistsBool)
	fmt.Println("pass_valid_bool",  passValidBool)
	fmt.Println("user_id_str",      userIDstr)

	return cookiesInRespLst
}

//-------------------------------------------------

func testUserHTTPupdate(pTest *testing.T,
	pCookiesLst  []*http.Cookie,
	pHTTPagent   *gorequest.SuperAgent,
	pTestPortInt int) {

	fmt.Println("====================================")
	fmt.Println("test user UPDATE")

	// auth_r=0 - this is a auth-ed call, but API only, so dont redirect on failure to auth
	urlStr := fmt.Sprintf("http://localhost:%d/v1/identity/update?auth_r=0", pTestPortInt)
	dataMap := map[string]string{
		"user_name_str":   "new username",
		"email_str":       "ivan@gloflow.com",
		"description_str": "some new description",
	}
	dataBytesLst, _ := json.Marshal(dataMap)
	_, bodyStr, errs := pHTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).

		// IMPORTANT!! - auth info is in cookies
		AddCookies(pCookiesLst).
		End()

	if len(errs) > 0 {
		spew.Dump(errs)
		pTest.FailNow()
	}

	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
		pTest.FailNow()
	}

	spew.Dump(bodyMap)

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user updating http request failed")
}

//-------------------------------------------------
// TEST_USER_GET_ME

func testUserHTTPgetMe(pTest *testing.T,
	pCookiesLst  []*http.Cookie,
	pHTTPagent   *gorequest.SuperAgent,
	pTestPortInt int) {

	// auth_r=0 - this is a auth-ed call, but API only, so dont redirect on failure to auth
	urlStr := fmt.Sprintf("http://localhost:%d/v1/identity/me?auth_r=0", pTestPortInt)
	_, bodyStr, errs := pHTTPagent.Get(urlStr).

		// IMPORTANT!! - auth info is in cookies
		AddCookies(pCookiesLst).
		End()

	if len(errs) > 0 {
		spew.Dump(errs)
		pTest.FailNow()
	}
	
	bodyMap := map[string]interface{}{}
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

	fmt.Println("RESPONSE >>>>")
	spew.Dump(bodyMap["data"])
	
	fmt.Println("====================================")
	fmt.Println("user me response:")
	fmt.Println("user_name_str",         userNameStr)
	fmt.Println("email_str",             emailStr)
	fmt.Println("description_str",       descriptionStr)
	fmt.Println("profile_image_url_str", profileImageURLstr)
	fmt.Println("banner_image_url_str",  bannerImageURLstr)
}