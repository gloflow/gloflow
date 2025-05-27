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
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/stretchr/testify/assert"
)

var logFun func(string, string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//---------------------------------------------------

func TestMain(m *testing.M) {
	logFun, logNewFun = gf_core.LogsInit()
	cliArgsMap = gf_images_core.CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------

func TestCreateDiscoveredTags(pTest *testing.T) {

	ctx := context.Background()

	serviceNameStr := "gf_tagger_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr := cliArgsMap["sql_host_str"].(string)
	runtimeSys := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

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
	sqlHostStr := cliArgsMap["sql_host_str"].(string)
	runtimeSys := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

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

	imageFromDB, gfErr := gf_images_core.DBgetImage(testImage.IDstr, ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	/*
		image, gfErr := gf_images_core.DBmongoGetImage(testImage.IDstr, ctx, runtimeSys)
		if gfErr != nil {
			pTest.FailNow()
		}
	*/

	fmt.Println("=================================================")
	spew.Dump(imageFromDB)

	//--------------------
	// CHECK TAGS COUNT USING dbSQLgetObjectsWithTag

	tagsToCheckLst := []string{"tag1", "tag2", "tag3"}
	imageIDSetMap := make(map[string]struct{})

	pageIndexInt := 0
	pageSizeInt := 100

	for _, tag := range tagsToCheckLst {

		imagesFromDBwithTagLst := []*gf_images_core.GFimage{}

		gfErr = dbSQLgetObjectsWithTag(tag,
			objectTypeStr, // "img",
			&imagesFromDBwithTagLst,
			pageIndexInt,
			pageSizeInt,
			ctx,
			runtimeSys)
		if gfErr != nil {
			pTest.FailNow()
		}

		fmt.Println("images with tag", len(imagesFromDBwithTagLst))
		for _, img := range imagesFromDBwithTagLst {
			fmt.Println("image ID:", img.IDstr)
		}

		foundBool := false
		for _, imgFromDB := range imagesFromDBwithTagLst {

			fmt.Println("---", imgFromDB.IDstr)

			imageIDSetMap[string(imgFromDB.IDstr)] = struct{}{}

			if imgFromDB.IDstr == testImage.IDstr {
				foundBool = true
			}
		}

		assert.True(pTest, foundBool, "the test image ID should be present for tag '"+tag+"'")
	}

	assert.True(pTest, len(tagsToCheckLst) == 3, "there should be 3 tags checked")

	//--------------------
}

//---------------------------------------------------

func createTestImages(pUserID gf_core.GF_ID,
	pTest *testing.T,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_images_core.GFimage {

	pRuntimeSys.LogNewFun("DEBUG", "creating test images...", nil)

	testImg0 := &gf_images_core.GFimage{
		IDstr:          gf_images_core.GFimageID(fmt.Sprintf("test_img_%d", time.Now().UnixNano())),
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0"},
		Origin_url_str: "https://gloflow.com/some_url0",
		TagsLst:        []string{},
	}
	gfErr := gf_images_core.DBsqlPutImage(testImg0, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	return testImg0
}
