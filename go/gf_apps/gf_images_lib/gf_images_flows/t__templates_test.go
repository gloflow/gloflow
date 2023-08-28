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
	"fmt"
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func TestTemplatesWithDB(pTest *testing.T) {

	ctx := context.Background()


	serviceNameStr := "gf_identity_test"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr     := cliArgsMap["sql_host_str"].(string)
	runtimeSys     := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	gfErr := gf_identity_core.DBsqlCreateTables(ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//-------------------
	// CREATE USER
	
	userID, userNameStr := gf_identity.TestCreateUserInDB(pTest, ctx, runtimeSys)

	//-------------------
	// CREATE_TEST_IMAGES
	createTestImages(userID, pTest, ctx, runtimeSys)



	flowNameStr        := "flow_0"
	initialPagesNumInt := 2
	pageSizeInt        := 2
	pagesLst, pagesUserNamesLst, flowPagesNumInt, gfErr := getTemplateData(flowNameStr,
		initialPagesNumInt, pageSizeInt,
		ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}




	spew.Dump(pagesLst)
	spew.Dump(pagesUserNamesLst)



	fmt.Println(flowPagesNumInt)



	for _, pLst := range pagesUserNamesLst {
		for _, resolvedUserNameStr := range pLst {

			/*
			check that the user_id that was assigned to images is resolved to the
			correct user_name of that user.
			the resolution of user_id to user_name happens via the user SQL table.
			*/
			assert.True(pTest, resolvedUserNameStr == userNameStr,
				"image user_id not resolved to the correct user_name")
		}
	}

}

//---------------------------------------------------

func TestTemplates(pTest *testing.T) {

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_images_flows_tests",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
	}

	userID := gf_core.GF_ID("test_user")

	// TEMPLATES
	templatesPathsMap := map[string]string{
		"gf_images_flows_browser": "./../../../../web/src/gf_apps/gf_images/templates/gf_images_flows_browser/gf_images_flows_browser.html",
	}
	
	gfTemplates, gfErr := tmplLoad(templatesPathsMap, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	imagesPagesLst := [][]*gf_images_core.GFimage{
		{
			&gf_images_core.GFimage{
				IDstr:      "some_test_id",
				TitleStr:   "some_test_img",
				MetaMap:    map[string]interface{}{"t_k": "val"},
				Format_str: "jpg",
				Thumbnail_small_url_str:  "url1",
				Thumbnail_medium_url_str: "url2",
				Thumbnail_large_url_str:  "url3",
				Origin_page_url_str:      "url4",
			},
		},
	}

	pagesUserNamesLst := [][]gf_identity_core.GFuserName{
		{
			gf_identity_core.GFuserName("image_owner_user_name"),
		},
	}

	flowNameStr     := "test_flow" 
	flowPagesNumInt := int64(6)
	templateRenderedStr, gfErr := renderTemplate(flowNameStr,
		imagesPagesLst,
		pagesUserNamesLst,
		flowPagesNumInt,
		gfTemplates.flows_browser__tmpl,
		gfTemplates.flows_browser__subtemplates_names_lst,
		userID,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	fmt.Println(templateRenderedStr)
}