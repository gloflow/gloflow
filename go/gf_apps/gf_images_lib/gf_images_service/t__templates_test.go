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

package gf_images_service

import (
	"fmt"
	"testing"
	"context"
	// "github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func TestTemplates(pTest *testing.T) {

	ctx := context.Background()
	serviceNameStr := "gf_images_service_tests"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr     := cliArgsMap["sql_host_str"].(string)
	runtimeSys     := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	userID := gf_core.GF_ID("test_user")
	firstTestImage := gf_images_core.CreateTestImages(userID, pTest, ctx, runtimeSys)

	// TEMPLATES
	templatesPathsMap := map[string]string{
		"gf_images_view": "./../../../../web/src/gf_apps/gf_images/templates/gf_images_view/gf_images_view.html",
	}
	
	gfTemplates, gfErr := templateLoad(templatesPathsMap, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	imageID := firstTestImage.IDstr
	templateRenderedStr, gfErr := renderImageViewPage(imageID,
		gfTemplates.imagesViewTmpl,
		gfTemplates.imagesViewSubtemplatesNamesLst,
		userID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	fmt.Println(templateRenderedStr)
}