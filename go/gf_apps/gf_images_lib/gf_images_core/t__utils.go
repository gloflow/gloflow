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

package gf_images_core

import (
	"context"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func CreateTestImages(pUserID gf_core.GF_ID,
	pTest       *testing.T,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *GFimage {

	pRuntimeSys.LogNewFun("DEBUG", "creating test images...", nil)

	testImg0 := &GFimage{
		IDstr: "test_img_0",
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0"},
		Origin_url_str: "https://gloflow.com/some_url0",
	}
	testImg1 := &GFimage{
		IDstr: "test_img_1",
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0"},
		Origin_url_str: "https://gloflow.com/some_url1",
	}
	testImg2 := &GFimage{
		IDstr: "test_img_2",
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0", "flow_1"},
		Origin_url_str: "https://gloflow.com/some_url2",
	}
	testImg3 := &GFimage{
		IDstr: "test_img_3",
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_1", "flow_2"},
		Origin_url_str: "https://gloflow.com/some_url3",
	}

	//----------------------------
	// DB
	
	gfErr := DBmongoPutImage(testImg0, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	gfErr = DBmongoPutImage(testImg1, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	gfErr = DBmongoPutImage(testImg2, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	gfErr = DBmongoPutImage(testImg3, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//----------------------------

	return testImg0
}