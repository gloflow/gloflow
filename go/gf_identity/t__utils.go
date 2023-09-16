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
	"time"
	"net/http"
	"context"
	"testing"
	"github.com/parnurzeal/gorequest"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

var logFun func(p_g string, p_m string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

type GFtestUserInfo struct {
	NameStr  string
	PassStr  string
	EmailStr string
}

//-------------------------------------------------

func TestCreateUserInDB(pTest *testing.T,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, gf_identity_core.GFuserName) {
	
	// DB
	gfErr := gf_identity_core.DBsqlCreateTables(pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	userID := gf_core.GF_ID(fmt.Sprintf("test_user_id_%s", gf_core.StrRandom()))
	userNameStr := gf_identity_core.GFuserName("test_user")
	screenNameStr := "test_user_screenname"

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

	user := &gf_identity_core.GFuser{
		Vstr:              "0",
		ID:                userID,
		CreationUNIXtimeF: creationUNIXtimeF,
		UserTypeStr:       "standard",
		UserNameStr:       userNameStr,
		ScreenNameStr:     screenNameStr,
	}
	
	// DB
	gfErr = gf_identity_core.DBsqlUserCreate(user, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	return userID, userNameStr
}

//-------------------------------------------------

func TestUserpassCreateAndLoginNewUser(pUserInfo *GFtestUserInfo,
	pTest                   *testing.T,
	pHTTPagent              *gorequest.SuperAgent,
	pIdentityServicePortInt int,
	pCtx                    context.Context,
	pRuntimeSys             *gf_core.RuntimeSys) []*http.Cookie {

	//---------------------------------
	// CLEANUP
	TestDBcleanup(pCtx, pRuntimeSys)
	
	//---------------------------------
	// ADD_TO_INVITE_LIST
	gfErr := gf_identity_core.DBsqlUserAddToInviteList(pUserInfo.EmailStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.FailNow()
	}

	//---------------------------------
	// GF_IDENTITY_INIT
	TestUserHTTPcreate(pUserInfo.NameStr,
		pUserInfo.PassStr,
		pUserInfo.EmailStr,
		pHTTPagent,
		pIdentityServicePortInt,
		pTest)

	cookiesInRespLst := TestUserHTTPuserpassLogin(pUserInfo.NameStr,
		pUserInfo.PassStr,
		pHTTPagent,
		pIdentityServicePortInt,
		pTest)
		
	//---------------------------------

	return cookiesInRespLst
}

//-------------------------------------------------

func TestStartService(pAuthSubsystemTypeStr string,
	pTemplatesPathsMap map[string]string,
	pPortInt    int,
	pRuntimeSys *gf_core.RuntimeSys) *gf_identity_core.GFkeyServerInfo {

	serviceInitCh := make(chan bool)
	var keyServer *gf_identity_core.GFkeyServerInfo
	go func() {

		HTTPmux := http.NewServeMux()

		
		serviceInfo := &gf_identity_core.GFserviceInfo{

			AuthSubsystemTypeStr: pAuthSubsystemTypeStr, // "userpass",
			
			// IMPORTANT!! - durring testing dont send emails
			EnableEmailBool: false,
		}

		// spew.Dump(serviceInfo)
		// spew.Dump(pRuntimeSys)

		var gfErr *gf_core.GFerror
		keyServer, gfErr = InitService(pTemplatesPathsMap, HTTPmux, serviceInfo, pRuntimeSys)
		if gfErr != nil {
			panic("failed to initialize gf_identity service")
		}

		serviceInitCh <- true
		gf_rpc_lib.ServerInitWithMux("gf_identity_test", pPortInt, HTTPmux)
	}()
	time.Sleep(2*time.Second) // let server startup

	<-serviceInitCh
	return keyServer
}

//-------------------------------------------------

func Tinit(pServiceNameStr string,
	pMongoHostStr string,
	pSQLhostStr   string,
	pLogNewFun    gf_core.GFlogFun,
	pLogFun       func(string, string)) *gf_core.RuntimeSys {

	testMongodbHostStr   := pMongoHostStr
	testMongodbDBnameStr := "gf_tests"
	testMongodbURLstr    := fmt.Sprintf("mongodb://%s", testMongodbHostStr)

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: pServiceNameStr,
		LogFun:         pLogFun,
		LogNewFun:      pLogNewFun,
		Validator:      gf_core.ValidateInit(),
	}

	//--------------------
	// MONGODB

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

	//--------------------
	// SQL

	dbNameStr := "gf_tests"
	dbUserStr := "gf"

	dbHostStr := pSQLhostStr

	sqlDB, gfErr := gf_core.DBsqlConnect(dbNameStr,
		dbUserStr,
		"", // config.SQLpassStr,
		dbHostStr,
		runtimeSys)
	if gfErr != nil {
		panic(-1)
	}

	runtimeSys.SQLdb = sqlDB

	//--------------------

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