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
	"context"
	"encoding/json"
	// "net/http/cookiejar"
	// "net/url"
	"strings"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func TestUsersHTTPeth(pTest *testing.T) {


	serviceNameStr := "gf_identity_test"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr     := cliArgsMap["sql_host_str"].(string)
	runtimeSys     := Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)
	runtimeSys.LogNewFun("INFO", "TEST_USERS_HTTP_ETH >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", nil)


	//---------------------------------
	// START_SERVICE
	authSubsystemTypeStr := "eth"
	portInt := 2001

	templatesPathsMap := map[string]string {
		"gf_login": "./../../web/src/gf_identity/templates/gf_login/gf_login.html",
	}

	TestStartService(authSubsystemTypeStr,
		templatesPathsMap,
		portInt,
		runtimeSys)

	testPortInt := portInt
	
	//---------------------------------

	HTTPagent := gorequest.New()
	
	//---------------------------------
	// GENERATE_WALLET
	privateKeyHexStr, publicKeyHexStr, addressStr, err := gf_eth_core.EthGenerateKeys()
	if err != nil {
		runtimeSys.LogNewFun("DEBUG", "wallet generation failed", map[string]interface{}{"err": err,})
		pTest.FailNow()
	}

	runtimeSys.LogNewFun("DEBUG", "wallet generation done...", nil)

	//---------------------------------
	// TEST_PREFLIGHT_HTTP
	
	dataMap := map[string]string{
		"user_address_eth_str": addressStr,
	}
	dataBytesLst, _ := json.Marshal(dataMap)
	urlStr := fmt.Sprintf("http://localhost:%d/v1/identity/eth/preflight", testPortInt)
	_, bodyStr, errs := HTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	runtimeSys.LogNewFun("DEBUG", "eth preflight HTTP request done...", map[string]interface{}{"body": bodyStr,})

	if (len(errs) > 0) {
		runtimeSys.LogNewFun("DEBUG", "eth preflight HTTP failed", map[string]interface{}{"errs": errs,})
		pTest.FailNow()
	}

	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
        pTest.FailNow()
    }

	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(bodyMap)
	}

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user preflight http request failed")

	nonceValStr    := bodyMap["data"].(map[string]interface{})["nonce_val_str"].(string)
	userExistsBool := bodyMap["data"].(map[string]interface{})["user_exists_bool"].(bool)
	
	runtimeSys.LogNewFun("DEBUG", "eth preflight response...", map[string]interface{}{
		"nonce_val_str":    nonceValStr,
		"user_exists_bool": userExistsBool,
	})


	// we're testing user creation and the rest of the flow as well, so user shouldnt exist
	if (userExistsBool) {
		pTest.FailNow()
	}
	
	//---------------------------------
	// TEST_USER_CREATE_HTTP

	signatureStr, err := gf_eth_core.EthSignData(nonceValStr, privateKeyHexStr)
	if err != nil {
		fmt.Println(err)
		pTest.FailNow()
	}

	fmt.Println("====================================")
	fmt.Println("user create inputs:")
	fmt.Println("address",   addressStr, len(addressStr))
	fmt.Println("priv key",  privateKeyHexStr, len(privateKeyHexStr))
	fmt.Println("pub key",   publicKeyHexStr, len(publicKeyHexStr))
	fmt.Println("signature", signatureStr, len(signatureStr))
	fmt.Println("nonce",     nonceValStr)

	urlStr = fmt.Sprintf("http://localhost:%d/v1/identity/eth/create", testPortInt)
	dataMap = map[string]string{
		"user_address_eth_str": addressStr,
		"auth_signature_str":   signatureStr,
	}
	dataBytesLst, _  = json.Marshal(dataMap)
	_, bodyStr, errs = HTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	spew.Dump(bodyStr)

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}

	bodyMap = map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
        pTest.FailNow()
    }

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user create http request failed")

	nonceExistsBool        := bodyMap["data"].(map[string]interface{})["nonce_exists_bool"].(bool)
	authSignatureValidBool := bodyMap["data"].(map[string]interface{})["auth_signature_valid_bool"].(bool)

	if (!nonceExistsBool) {
		fmt.Println("supplied nonce doesnt exist")
		pTest.FailNow()
	}
	if (!authSignatureValidBool) {
		fmt.Println("signature is not valid")
		pTest.FailNow()
	}

	//---------------------------------
	// TEST_USER_LOGIN

	fmt.Println("====================================")
	fmt.Println("user login inputs:")
	fmt.Println("address",   addressStr, len(addressStr))
	fmt.Println("signature", signatureStr, len(signatureStr))

	
	dataMap = map[string]string{
		"user_address_eth_str": addressStr,
		"auth_signature_str":   signatureStr,
	}

	urlStr = fmt.Sprintf("http://localhost:%d/v1/identity/eth/login", testPortInt)
	dataBytesLst, _ = json.Marshal(dataMap)
	resp, bodyStr, errs := HTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}

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
	
	bodyMap = map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
        pTest.FailNow()
    }

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "user login http request failed")

	nonceExistsBool        = bodyMap["data"].(map[string]interface{})["nonce_exists_bool"].(bool)
	authSignatureValidBool = bodyMap["data"].(map[string]interface{})["auth_signature_valid_bool"].(bool)
	userIDstr             := bodyMap["data"].(map[string]interface{})["user_id_str"].(string)

	fmt.Println("RESPONSE >>>>")
	spew.Dump(bodyMap["data"])

	fmt.Println("====================================")
	fmt.Println("user login response:")
	fmt.Println("nonce_exists_bool",         nonceExistsBool)
	fmt.Println("auth_signature_valid_bool", authSignatureValidBool)
	fmt.Println("user_id_str",               userIDstr)

	if (!nonceExistsBool) {
		fmt.Println("supplied nonce doesnt exist")
		pTest.FailNow()
	}
	if (!authSignatureValidBool) {
		fmt.Println("signature is not valid")
		pTest.FailNow()
	}

	//---------------------------------
	// TEST_USER_HTTP_UPDATE
	testUserHTTPupdate(pTest, cookiesInRespLst, HTTPagent, testPortInt)

	//---------------------------------
	// TEST_USER_HTTP_GET_ME
	testUserHTTPgetMe(pTest, cookiesInRespLst, HTTPagent, testPortInt)

	//---------------------------------
}

//-------------------------------------------------

func TestUsersETHunit(pTest *testing.T) {

	serviceNameStr := "gf_identity_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys   := Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)
	
	runtimeSys.LogNewFun("INFO", "TEST_USERS_ETH_UNIT >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", nil)

	
	testUserAddressEthStr := "0xBA47Bef4ca9e8F86149D2f109478c6bd8A642C97"
	testUserSignatureStr  := "0x07c582de2c6fb11310495815c993fa978540f0c0cdc89fd51e6fe3b8db62e913168d9706f32409f949608bcfd372d41cbea6eb75869afe2f189738b7fb764ef91c"
	testUserNonceStr      := "gf_test_message_to_sign"
	ctx := context.Background()

	//------------------------
	// KEY_SERVER
	keyServerInfo, gfErr := gf_identity_core.KSinit(false, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//------------------
	// NONCE_CREATE

	unexistingUserIDstr := gf_core.GF_ID("")
	_, gfErr = gf_identity_core.NonceCreate(gf_identity_core.GFuserNonceVal(testUserNonceStr),
		unexistingUserIDstr,
		gf_identity_core.GFuserAddressETH(testUserAddressEthStr),
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//------------------
	// USER_CREATE
	
	inputCreate := &gf_identity_core.GFethInputCreate{
		UserTypeStr:      "standard",
		AuthSignatureStr: gf_identity_core.GFauthSignature(testUserSignatureStr),
		UserAddressETH:   gf_identity_core.GFuserAddressETH(testUserAddressEthStr),
	}

	outputCreate, gfErr := gf_identity_core.ETHpipelineCreate(inputCreate, ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	spew.Dump(outputCreate)

	assert.True(pTest, outputCreate.AuthSignatureValidBool, "crypto signature supplied for user creation pipeline is invalid")

	//------------------
	inputLogin := &gf_identity_core.GFethInputLogin{
		AuthSignatureStr: gf_identity_core.GFauthSignature(testUserSignatureStr),
		UserAddressETH:   gf_identity_core.GFuserAddressETH(testUserAddressEthStr),
	}
	outputLogin, gfErr := gf_identity_core.ETHpipelineLogin(inputLogin,
		keyServerInfo,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	spew.Dump(outputLogin)
	
	//------------------
}