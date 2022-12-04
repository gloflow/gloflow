/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_landing_page_lib

import (
	"fmt"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_tagger_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

var logFun func(string,string)
var cliArgsMap map[string]interface{}

//---------------------------------------------------

func TestMain(m *testing.M) {
	logFun, _  = gf_core.InitLogs()
	cliArgsMap = gf_tagger_core.CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------

func TestBookmarks(pTest *testing.T) {

	fmt.Println(" TEST__MAIN >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	testMongodbHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	testMongodbDBnameStr := "gf_tests"
	testMongodbURLstr    := fmt.Sprintf("mongodb://%s", testMongodbHostStr)

	logFun, _ := gf_core.InitLogs()


	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_landing_page_test",
		LogFun:           logFun,
		Validator:        gf_core.ValidateInit(),
	}

	mongoDB, _, gfErr := gf_core.MongoConnectNew(testMongodbURLstr, testMongodbDBnameStr, nil, runtimeSys)
	if gfErr != nil {
		panic(-1)
	}

	mongoColl := mongoDB.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongoDB
	runtimeSys.Mongo_coll = mongoColl

	testBookmarkingFlow(pTest, runtimeSys)
}

//-------------------------------------------------

func testBookmarkingFlow(pTest *testing.T,
	pRuntimeSys *gf_core.RuntimeSys) {

	


	templatesPathsMap := map[string]string{
		"gf_landing_page": "./../../../web/src/gf_apps/gf_landing_page/templates/gf_landing_page/gf_landing_page.html",
	}
	// TEMPLATES
	gfTemplates, gfErr := templatesLoad(templatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}


	imagesMaxRandomCursorPositionInt := 10
	postsMaxRandomCursorPositionInt  := 10
	featuredPostsToGetInt            := 3
	featuredImagesToGetInt           := 3
	templateRenderedStr, gfErr := pipelineRenderLandingPage(imagesMaxRandomCursorPositionInt,
		postsMaxRandomCursorPositionInt,
		featuredPostsToGetInt,
		featuredImagesToGetInt,

		gfTemplates.template,
		gfTemplates.subtemplatesNamesLst,
		pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	fmt.Println(templateRenderedStr)
	assert.True(pTest, templateRenderedStr != "", "gf_landing_page was not rendered as a html template,")

	//------------------
}