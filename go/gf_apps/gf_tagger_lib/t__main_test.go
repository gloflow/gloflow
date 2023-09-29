/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

package gf_tagger_lib

import (
	"os"
	// "fmt"
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

func TestCreateDiscoveredTags(pTest *testing.T) {

	ctx := context.Background()

	serviceNameStr := "gf_tagger_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys   := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)
	
	userID := gf_core.GF_ID("test")

	//--------------------
	// INIT
	
	gfErr := dbSQLcreateTables(runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	testImage := createTestImages(userID, pTest, ctx, runtimeSys)

	//--------------------
	// ADD_TAGS_TO_OBJECT

	tagsStr := "tag1 tag2 tag3"
	objectTypeStr := "image"

	metaMap := map[string]interface{}{}

	gfErr = addTagsToObject(tagsStr,
		objectTypeStr,
		string(testImage.IDstr), // objectExternIDstr,
		metaMap,
		userID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//--------------------



	gfErr = pipelineCreateDiscoveredTags(ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
}

//---------------------------------------------------

func TestCreate(pTest *testing.T) {

	ctx := context.Background()

	serviceNameStr := "gf_images_flows_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys   := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)
	
	userID := gf_core.GF_ID("test")

	//--------------------
	// INIT
	
	gfErr := dbSQLcreateTables(runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	testImage := createTestImages(userID, pTest, ctx, runtimeSys)

	//--------------------
	// ADD_TAGS_TO_OBJECT

	tagsStr := "tag1 tag2 tag3"
	objectTypeStr := "image"


	metaMap := map[string]interface{}{}

	gfErr = addTagsToObject(tagsStr,
		objectTypeStr,
		string(testImage.IDstr), // objectExternIDstr,
		metaMap,
		userID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//--------------------
	// DB_GET_IMAGE

	image, gfErr := gf_images_core.DBmongoGetImage(testImage.IDstr, ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	spew.Dump(image)
	assert.True(pTest, len(image.TagsLst) == 3, "image should have 3 tags added to it")

	//--------------------
	// DB_GET_OBJECTS_WITH_TAG

	pageIndexInt := 0
	pageSizeInt := 5

	imagesLst := []*gf_images_core.GFimage{}
	gfErr = dbMongoGetObjectsWithTag("tag1",
		"img",
		&imagesLst,
		pageIndexInt,
		pageSizeInt,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	spew.Dump(imagesLst)
	assert.True(pTest, len(imagesLst) == 1, "there should be 1 image with tag1 tag in the DB")

	//--------------------
}

//---------------------------------------------------

func createTestImages(pUserID gf_core.GF_ID,
	pTest       *testing.T,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_images_core.GFimage {

	pRuntimeSys.LogNewFun("DEBUG", "creating test images...", nil)

	testImg0 := &gf_images_core.GFimage{
		IDstr: "test_img_0",
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0"},
		Origin_url_str: "https://gloflow.com/some_url0",
		TagsLst:        []string{},
	}
	gfErr := gf_images_core.DBmongoPutImage(testImg0, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	return testImg0
}