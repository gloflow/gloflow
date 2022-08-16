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

package gf_home_lib

import (
	"os"
	"fmt"
	"time"
	"testing"
	"net/http"
	"encoding/json"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func TestMain(m *testing.M) {

	logFun     = gf_core.InitLogs()
	cliArgsMap = CLIparseArgs(logFun)

	runtimeSys := Tinit()

	templatesPathsMap := map[string]string{
		"gf_home_main": "./../../../web/src/gf_apps/gf_home/templates/gf_home_main/gf_home_main.html",
	}

	// GF_HOME_SERVICE
	testPortInt := 2000
	go func() {

		HTTPmux := http.NewServeMux()

		serviceInfo := &GFserviceInfo{}
		InitService(templatesPathsMap,
			serviceInfo,
			HTTPmux,
			runtimeSys)
		gf_rpc_lib.Server__init_with_mux(testPortInt, HTTPmux)
	}()

	// GF_IDENTITY_SERVICE
	testIdentityServicePortInt := 2001
	go func() {

		gf_identity_lib.TestStartService(testIdentityServicePortInt,
			runtimeSys)
	}()

	time.Sleep(2*time.Second) // let services startup

	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func TestHomeViz(pTest *testing.T) {

	runtimeSys := Tinit()
	fmt.Println(runtimeSys)

	HTTPagent := gorequest.New()
	ctx := context.Background()
	testPortInt := 2000
	testIdentityServicePortInt := 2001
	
	
	// CREATE_AND_LOGIN_NEW_USER
	gf_identity_lib.TestCreateAndLoginNewUser(pTest,
		HTTPagent,
		testIdentityServicePortInt,
		ctx,
		runtimeSys)

	/*//---------------------------------
	// CLEANUP
	gf_identity_lib.TestDBcleanup(ctx, runtimeSys)
	
	//---------------------------------
	// ADD_TO_INVITE_LIST
	gfErr := gf_identity_lib.DBuserAddToInviteList(testEmailStr,
		ctx,
		runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.FailNow()
	}

	//---------------------------------
	// GF_IDENTITY_INIT
	gf_identity_lib.TestUserHTTPcreate(testUserNameStr,
		testUserPassStr,
		testEmailStr,
		HTTPagent,
		testIdentityServicePortInt,
		pTest)

	gf_identity_lib.TestUserHTTPlogin(testUserNameStr,
		testUserPassStr,
		HTTPagent,
		testIdentityServicePortInt,
		pTest)
		
	//---------------------------------*/
	

	fmt.Println("======== HOME_VIZ GET HTTP")
	urlStr := fmt.Sprintf("http://localhost:%d/v1/home/viz/get", testPortInt)
	_, bodyStr, errs := HTTPagent.Get(urlStr).
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

}

//---------------------------------------------------
func TestTemplates(pTest *testing.T) {

	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_home_test",
		LogFun:           logFun,
	}

	// TEMPLATES
	templatesPathsMap := map[string]string{
		"gf_home_main": "./../../../web/src/gf_apps/gf_home/templates/gf_home_main/gf_home_main.html",
	}
	
	templates, gfErr := templatesLoad(templatesPathsMap, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	templateRenderedStr, gfErr := viewRenderTemplateDashboard(templates.mainTmpl,
		templates.mainSubtemplatesNamesLst,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	fmt.Println(templateRenderedStr)
}