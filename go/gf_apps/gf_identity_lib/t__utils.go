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
	"time"
	"net/http"
	"context"
	"testing"
	"github.com/parnurzeal/gorequest"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

var logFun func(p_g string, p_m string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//-------------------------------------------------

func TestCreateAndLoginNewUser(pTest *testing.T,
	pHTTPagent              *gorequest.SuperAgent,
	pIdentityServicePortInt int,
	pCtx                    context.Context,
	pRuntimeSys             *gf_core.RuntimeSys) {

	testUserNameStr := "ivan_t"
	testUserPassStr := "pass_lksjds;lkdj"
	testEmailStr    := "ivan_t@gloflow.com"

	//---------------------------------
	// CLEANUP
	TestDBcleanup(pCtx, pRuntimeSys)
	
	//---------------------------------
	// ADD_TO_INVITE_LIST
	gfErr := gf_identity_core.DBuserAddToInviteList(testEmailStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.FailNow()
	}

	//---------------------------------
	// GF_IDENTITY_INIT
	TestUserHTTPcreate(testUserNameStr,
		testUserPassStr,
		testEmailStr,
		pHTTPagent,
		pIdentityServicePortInt,
		pTest)

	TestUserHTTPlogin(testUserNameStr,
		testUserPassStr,
		pHTTPagent,
		pIdentityServicePortInt,
		pTest)
		
	//---------------------------------
}

//-------------------------------------------------

func TestStartService(pPortInt int,
	pRuntimeSys *gf_core.RuntimeSys) {

	// testPortInt := 2000
	go func() {

		HTTPmux := http.NewServeMux()

		serviceInfo := &gf_identity_core.GFserviceInfo{

			AuthSubsystemTypeStr: "builtin",
			
			// IMPORTANT!! - durring testing dont send emails
			EnableEmailBool: false,
		}

		spew.Dump(serviceInfo)
		InitService(HTTPmux, serviceInfo, pRuntimeSys)
		gf_rpc_lib.ServerInitWithMux("gf_identity_test", pPortInt, HTTPmux)
	}()
	time.Sleep(2*time.Second) // let server startup
}

//-------------------------------------------------

func Tinit() *gf_core.RuntimeSys {

	testMongodbHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	testMongodbDBnameStr := "gf_tests"
	testMongodbURLstr    := fmt.Sprintf("mongodb://%s", testMongodbHostStr)

	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_identity_tests",
		LogFun:           logFun,
		LogNewFun:        logNewFun,
		Validator:        gf_core.ValidateInit(),
	}

	mongoDB, _, gfErr := gf_core.MongoConnectNew(testMongodbURLstr,
		testMongodbDBnameStr,
		nil,
		runtimeSys)
	if gfErr != nil {
		panic(-1)
	}

	mongoColl := mongoDB.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongoDB
	runtimeSys.Mongo_coll = mongoColl

	return runtimeSys
}

//-------------------------------------------------

func TestDBcleanup(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) {
	
	// CLEANUP
	collNameStr := "gf_users"
	gf_core.MongoDelete(bson.M{}, collNameStr, 
		map[string]interface{}{
			"caller_err_msg_str": "failed to cleanup test user DB",
		},
		pCtx, pRuntimeSys)
}