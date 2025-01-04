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
	"github.com/gloflow/gloflow/go/gf_identity/gf_policy"
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

func TestCreateIfMissing(pTest *testing.T) {

	ctx := context.Background()
	serviceNameStr := "gf_images_flows_test"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr     := cliArgsMap["sql_host_str"].(string)
	runtimeSys     := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	initialFlowsNamesLst := []string{
		"flow5",
	}

	secondFlowsNamesLst := []string{
		"flow10",
	}

	//-------------------
	// CREATE USER

	userID, _ := gf_identity.TestCreateUserInDB(pTest, ctx, runtimeSys)

	//------------------
	// INIT_TABLES
	gfErr := gf_images_core.DBsqlCreateTables(ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	gfErr = DBsqlCreateTables(runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	gfErr = gf_policy.DBsqlCreateTables(runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//------------------
	// INIT FLOW AND POLICY
	testFlowsIDsLst, gfErr := CreateIfMissing(initialFlowsNamesLst,
		userID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	testFlowID := testFlowsIDsLst[0]

	// CREATE
	policyOwnerUserID := userID
	targetResourceID := testFlowID
	_, gfErr = gf_policy.CreateTestPolicy(targetResourceID, policyOwnerUserID, ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//------------------
	// TEST
	// test creating a flow and adding permission to the users policy to add images to that flow

	gfErr = CreateIfMissingWithPolicy(secondFlowsNamesLst,
		userID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//------------------
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

	imagesExternURLsLst := []string{
		"https://gloflow.com/some_url0",
		"https://gloflow.com/some_url1",
	}

	//--------------------
	// INIT
	gfErr := Init(runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//------------------
	// INIT_TABLES
	gfErr = gf_images_core.DBsqlCreateTables(ctx, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()
	}

	//------------------
	initExistingImagesLst, gfErr := gf_images_core.DBimageExistsByURLs(imagesExternURLsLst,
		flowNameStr,
		clientTypeStr,
		userID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//------------------
	// CREATE_TEST_IMAGES
	
	gf_images_core.CreateTestImages(userID, pTest, ctx, runtimeSys)
	//------------------

	existingImagesLst, gfErr := gf_images_core.DBimageExistsByURLs(imagesExternURLsLst,
		flowNameStr,
		clientTypeStr,
		userID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	
	spew.Dump(existingImagesLst)

	assert.True(pTest, len(existingImagesLst) == len(initExistingImagesLst) + 2, "2 images should be found to exist in target flow")
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

	//-------------------
	// CREATE USER

	userID, _ := gf_identity.TestCreateUserInDB(pTest, ctx, runtimeSys)

	//------------------
	// INIT_TABLES
	gfErr := gf_images_core.DBsqlCreateTables(ctx, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()
	}

	//------------------
	// INITIAL_COUNTS - to compare against

	initAllFlowsCountsLst, gfErr := pipelineGetAll(ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	initCountAint := initAllFlowsCountsLst[0]["flow_imgs_count_int"].(int)
	initCountBint := initAllFlowsCountsLst[1]["flow_imgs_count_int"].(int)
	initCountCint := initAllFlowsCountsLst[2]["flow_imgs_count_int"].(int)

	//------------------
	// CREATE_TEST_IMAGES
	
	testImageIDsLst := CreateTestImagesInFlows(userID, pTest, ctx, runtimeSys)
	spew.Dump(testImageIDsLst)
	
	//------------------

	newFlowsCountsLst, gfErr := pipelineGetAll(ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	spew.Dump(newFlowsCountsLst)

	newCountAint := newFlowsCountsLst[0]["flow_imgs_count_int"].(int)
	newCountBint := newFlowsCountsLst[1]["flow_imgs_count_int"].(int)
	newCountCint := newFlowsCountsLst[2]["flow_imgs_count_int"].(int)

	assert.True(pTest, len(newFlowsCountsLst) >= 3, "minimum of 3 flows in total should have been discovered")
	assert.True(pTest, newFlowsCountsLst[0]["flow_name_str"].(string) == "flow_0", "first flow should be flow_0")
	assert.True(pTest, newFlowsCountsLst[1]["flow_name_str"].(string) == "flow_1", "first flow should be flow_1")
	assert.True(pTest, newFlowsCountsLst[2]["flow_name_str"].(string) == "flow_2", "first flow should be flow_2")

	fmt.Println(initCountAint, initCountBint, initCountCint)
	fmt.Println(newCountAint, newCountBint, newCountCint)


	// check new counts after image creation
	assert.True(pTest, newCountAint == initCountAint+3, "first flow should have a count of 3")
	assert.True(pTest, newCountBint == initCountBint+2, "second flow should have a count of 2")
	assert.True(pTest, newCountCint == initCountCint+1, "third flow should have a count of 1")
}