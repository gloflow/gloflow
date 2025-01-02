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
	"fmt"
	"time"
	"context"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func CreateTestImages(pUserID gf_core.GF_ID,
	pTest       *testing.T,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) []GFimageID {

	pRuntimeSys.LogNewFun("DEBUG", "creating test images...", nil)

	testImg0 := &GFimage{
		IDstr: GFimageID(fmt.Sprintf("test_img_%d", time.Now().UnixNano())),
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0"},
		Origin_url_str: "https://gloflow.com/some_url0",
	}
	testImg1 := &GFimage{
		IDstr: GFimageID(fmt.Sprintf("test_img_%d", time.Now().UnixNano())),
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0"},
		Origin_url_str: "https://gloflow.com/some_url1",
	}
	testImg2 := &GFimage{
		IDstr: GFimageID(fmt.Sprintf("test_img_%d", time.Now().UnixNano())),
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_0", "flow_1"},
		Origin_url_str: "https://gloflow.com/some_url2",
	}
	testImg3 := &GFimage{
		IDstr: GFimageID(fmt.Sprintf("test_img_%d", time.Now().UnixNano())),
		T_str: "img",
		UserID:         pUserID,
		FlowsNamesLst:  []string{"flow_1", "flow_2"},
		Origin_url_str: "https://gloflow.com/some_url3",
	}

	//----------------------------
	// DB
	
	gfErr := DBsqlPutImage(testImg0, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	gfErr = DBsqlPutImage(testImg1, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	gfErr = DBsqlPutImage(testImg2, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}
	gfErr = DBsqlPutImage(testImg3, pCtx, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	//----------------------------

	testImageIDsLst := []GFimageID{
		testImg0.IDstr,
		testImg1.IDstr,
		testImg2.IDstr,
		testImg3.IDstr,
	}

	return testImageIDsLst
}