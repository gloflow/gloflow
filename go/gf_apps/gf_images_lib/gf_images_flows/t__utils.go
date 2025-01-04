/*
GloFlow application and media management/publishing platform
Copyright (C) 2025 Ivan Trajkovic

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
	// "fmt"
	"testing"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func CreateTestImagesInFlows(pUserID gf_core.GF_ID,
	pTest       *testing.T,
	pCtx		context.Context,
	pRuntimeSys	*gf_core.RuntimeSys) *gf_core.GFerror {

	// CREATE_TEST_IMAGES
	testImagesIDsLst := gf_images_core.CreateTestImages(pUserID, pTest, pCtx, pRuntimeSys)

	// get flows names attached to these images
	imagesFlowsNamesMap, gfErr := gf_images_core.DBsqlGetImagesFlows(testImagesIDsLst, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	uniqueFlowsNamesLst := gf_images_core.GetUniqueFlowNames(imagesFlowsNamesMap)

	// CREATE_TEST_FLOWS
	_, gfErr = CreateIfMissing(uniqueFlowsNamesLst, pUserID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}