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
	"fmt"
	"testing"
	"context"
	// "github.com/stretchr/testify/assert"
	// "github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_tagger_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func TestTemplates(pTest *testing.T) {

	ctx := context.Background()

	serviceNameStr := "gf_tagger_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys   := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)
	
	//-------------------
	// CREATE USER
	userID, _ := gf_identity.TestCreateUserInDB(pTest, ctx, runtimeSys)

	//--------------------
	// INIT
	
	gfErr := dbSQLcreateTables(runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}
	testImage := createTestImages(userID, pTest, ctx, runtimeSys)

	// TEMPLATES
	templatesPathsMap := map[string]string{
		"gf_tag_objects": "./../../../web/src/gf_apps/gf_tagger/templates/gf_tag_objects/gf_tag_objects.html",
		"gf_bookmarks": "./../../../web/src/gf_apps/gf_tagger/templates/gf_bookmarks/gf_bookmarks.html",
	}
	
	// TEMPLATES
	gfTemplates, gfErr := gf_tagger_core.TemplatesLoad(templatesPathsMap, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	

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
		pTest.Fail()
	}

	//--------------------

	pageIndexInt := 0
	pageSizeInt := 5

	templateRenderedStr, gfErr := renderObjectsWithTag("tag1",
		gfTemplates.TagObjects,
		gfTemplates.TagObjectsSubtemplatesNamesLst,
		pageIndexInt,
		pageSizeInt,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	fmt.Println(templateRenderedStr)
}