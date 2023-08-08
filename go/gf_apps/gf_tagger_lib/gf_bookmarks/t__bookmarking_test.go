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

package gf_bookmarks

import (
	"fmt"
	"os"
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_tagger_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

var logFun func(string,string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//---------------------------------------------------

func TestMain(m *testing.M) {
	logFun, logNewFun  = gf_core.LogsInit()
	cliArgsMap = gf_tagger_core.CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------

func TestBookmarks(pTest *testing.T) {

	fmt.Println(" TEST__BOOKMARKS_MAIN >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	serviceNameStr := "gf_bookmarks_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys   := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	testBookmarkingFlow(pTest, runtimeSys)
}

//-------------------------------------------------

func testBookmarkingFlow(pTest *testing.T,
	pRuntimeSys *gf_core.RuntimeSys) {

	ctx := context.Background()

	testUserIDstr := gf_core.GF_ID("test_user")
	//------------------
	// CREATE
	inputCreate := &GFbookmarkInputCreate{
		User_id_str:     testUserIDstr,
		Url_str:         "https://gloflow.com",
		Description_str: "test bookmark",
		Tags_lst: []string{
			"test", "code", "art",
		},
	}
	gfErr := PipelineCreate(inputCreate,
		nil, // p_images_jobs_mngr
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//------------------
	// GET_ALL__JSON

	inputGet := &GFbookmarkInputGet{
		Response_format_str: "json",
		User_id_str:         testUserIDstr,
	}
	output, gfErr := PipelineGet(inputGet,
		nil,
		nil,
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	spew.Dump(output.Bookmarks_lst)
	assert.True(pTest, len(output.Bookmarks_lst) > 0, "no bookmarks were returned")
	assert.True(pTest, output.TemplateRenderedStr == "", "bookmarks were rendered as a template, when it should be data-only")

	//------------------



	templatesPathsMap := map[string]string{
		"gf_tag_objects": "./../../../../web/src/gf_apps/gf_tagger/templates/gf_tag_objects/gf_tag_objects.html",
		"gf_bookmarks":   "./../../../../web/src/gf_apps/gf_tagger/templates/gf_bookmarks/gf_bookmarks.html",
	}
	// TEMPLATES
	gfTemplates, gfErr := gf_tagger_core.TemplatesLoad(templatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}




	inputGetHTML := &GFbookmarkInputGet{
		Response_format_str: "html",
		User_id_str:         testUserIDstr,
	}
	outputHTML, gfErr := PipelineGet(inputGetHTML,
		gfTemplates.Bookmarks,
		gfTemplates.BookmarksSubtemplatesNamesLst,
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	fmt.Println(outputHTML.TemplateRenderedStr)
	assert.True(pTest, len(outputHTML.Bookmarks_lst) == 0, "bookmarks were returned when it should only be html template string")
	assert.True(pTest, outputHTML.TemplateRenderedStr != "", "bookmarks were not rendered as a html template,")

	//------------------
}