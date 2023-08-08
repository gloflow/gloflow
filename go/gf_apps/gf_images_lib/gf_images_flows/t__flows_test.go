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

package gf_images_flows

import (
	"os"
	"fmt"
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/davecgh/go-spew/spew"
)

var logFun func(string,string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//---------------------------------------------------

func TestMain(m *testing.M) {
	logFun, logNewFun  = gf_core.LogsInit()
	cliArgsMap = gf_images_core.CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------

func TestCreate(pTest *testing.T) {

	ctx := context.Background()

	serviceNameStr := "gf_images_flows_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys   := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)
	
	
	flowNameStr := "test"
	userID := gf_core.GF_ID("test")

	//--------------------
	// INIT
	gfErr := Init(runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//--------------------
	// CREATE_FLOW
	flow, gfErr := Create(flowNameStr,
		userID,
		ctx,
		runtimeSys)

	if gfErr != nil {
		pTest.Fail()
	}


	spew.Dump(flow)

	//--------------------
}

//---------------------------------------------------

func TestGetAll(pTest *testing.T) {

	ctx := context.Background()

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_images_flows_tests",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
	}

	//------------------
	// MONGODB
	testMongodbHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	testMongodbURLstr    := fmt.Sprintf("mongodb://%s", testMongodbHostStr)
	testMongodbDBnameStr := "gf_tests"
	mongodbDB, _, gfErr  := gf_core.MongoConnectNew(testMongodbURLstr, testMongodbDBnameStr, nil, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()
	}
	mongodbColl := mongodbDB.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongodbDB
	runtimeSys.Mongo_coll = mongodbColl
	
	//------------------
	// CREATE_TEST_IMAGES
	testImg0 := &gf_images_core.GFimage{
		IDstr: "test_img_0",
		T_str: "img",
		FlowsNamesLst: []string{"flow_0"},
	}
	testImg1 := &gf_images_core.GFimage{
		IDstr: "test_img_1",
		T_str: "img",
		FlowsNamesLst: []string{"flow_0"},
	}
	testImg2 := &gf_images_core.GFimage{
		IDstr: "test_img_2",
		T_str: "img",
		FlowsNamesLst: []string{"flow_0", "flow_1"},
	}
	testImg3 := &gf_images_core.GFimage{
		IDstr: "test_img_3",
		T_str: "img",
		FlowsNamesLst: []string{"flow_1", "flow_2"},
	}
	gfErr = gf_images_core.DBputImage(testImg0, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}
	gfErr = gf_images_core.DBputImage(testImg1, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}
	gfErr = gf_images_core.DBputImage(testImg2, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}
	gfErr = gf_images_core.DBputImage(testImg3, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}
 
	//------------------


	allFlowsNamesLst, gfErr := pipelineGetAll(ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	spew.Dump(allFlowsNamesLst)

	assert.True(pTest, len(allFlowsNamesLst) == 3, "3 flows in total should have been discovered")
	assert.True(pTest, allFlowsNamesLst[0]["flow_name_str"].(string) == "flow_0", "first flow should be flow_0")
	assert.True(pTest, allFlowsNamesLst[1]["flow_name_str"].(string) == "flow_1", "first flow should be flow_1")
	assert.True(pTest, allFlowsNamesLst[2]["flow_name_str"].(string) == "flow_2", "first flow should be flow_2")

	assert.True(pTest, allFlowsNamesLst[0]["flow_imgs_count_int"].(int32) == 3, "first flow should have a count of 3")
	assert.True(pTest, allFlowsNamesLst[1]["flow_imgs_count_int"].(int32) == 2, "second flow should have a count of 2")
	assert.True(pTest, allFlowsNamesLst[2]["flow_imgs_count_int"].(int32) == 1, "third flow should have a count of 1")
}