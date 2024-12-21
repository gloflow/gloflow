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

func TestImagesExist(pTest *testing.T) {

	ctx := context.Background()

	serviceNameStr := "gf_images_flows_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys   := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	flowNameStr := "flow_0"
	userID := gf_core.GF_ID("test")
	clientTypeStr := "test"

	//--------------------
	// INIT
	gfErr := Init(runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//------------------
	// CREATE_TEST_IMAGES
	
	gf_images_core.CreateTestImages(userID, pTest, ctx, runtimeSys)
	//------------------

	imagesExternURLsLst := []string{
		"https://gloflow.com/some_url0",
		"https://gloflow.com/some_url1",
	}

	existingImagesLst, gfErr := imagesExistCheck(imagesExternURLsLst,
		flowNameStr,
		clientTypeStr,
		userID,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	
	spew.Dump(existingImagesLst)

	assert.True(pTest, len(existingImagesLst) == 2, "2 images should be found to exist in target flow")
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
		pTest.FailNow()
	}

	//--------------------
	// CREATE_FLOW
	flow, gfErr := Create(flowNameStr,
		userID,
		ctx,
		runtimeSys)

	if gfErr != nil {
		pTest.FailNow()
	}

	spew.Dump(flow)

	//--------------------
}

//---------------------------------------------------

func TestGetAll(pTest *testing.T) {

	ctx := context.Background()
	serviceNameStr := "gf_images_flows_test"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr     := cliArgsMap["sql_host_str"].(string)
	runtimeSys     := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	userID := gf_core.GF_ID("test_user")

	//------------------
	// INIT_TABLES
	gfErr := gf_images_core.DBsqlCreateTables(ctx, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()
	}

	//------------------

	//------------------
	// CREATE_TEST_IMAGES
	
	gf_images_core.CreateTestImages(userID, pTest, ctx, runtimeSys)
	//------------------


	allFlowsNamesLst, gfErr := pipelineGetAll(ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
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