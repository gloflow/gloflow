/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func TestImageClassify(pTest *testing.T) {

	ctx := context.Background()
	serviceNameStr := "gf_images_service_tests"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr     := cliArgsMap["sql_host_str"].(string)
	runtimeSys     := gf_identity.Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	//-------------------
	// CREATE USER
	userID, _ := gf_identity.TestCreateUserInDB(pTest, ctx, runtimeSys)

	//-------------------



	fmt.Println(userID)







	testImagesStoreLocalDirPathStr           := "./../tests_data/images"
	testImagesThumbnailsStoreLocalDirPathStr := "./../tests_data/thumbnails"
	testVideosStoreLocalDirPathStr           := "./../tests_data/videos"
	testMediaDomainStr                       := ""
	pluginsPyDirPathStr        := "./../../../../py/src/gf_apps/gf_images/plugins"
	imagesClassifyPyDirPathStr := "./../../../../py/src/gf_apps/gf_images/gf_images_classify"


	//------------------------
	// IMAGES_JOBS_MNGR

	jobsMngr := gf_images_jobs_core.TgetJobsMngr(testImagesStoreLocalDirPathStr,
		testImagesThumbnailsStoreLocalDirPathStr,
		testVideosStoreLocalDirPathStr,
		testMediaDomainStr,
		pluginsPyDirPathStr,
		imagesClassifyPyDirPathStr,
		runtimeSys)

	//------------------------
		


	input := &GFimageClassifyInput{
		ClientTypeStr: "testing",
		ImagesIDsLst: []gf_images_core.GFimageID{"image_id_1", "image_id_2"},
	}

	serviceInfo := &gf_images_core.GFserviceInfo{
		EnableEventsAppBool: false,
	}

	imageClassesLst, gfErr := ImageClassify(input,
		userID,
		jobsMngr,
		serviceInfo,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
		panic(gfErr.Error)
	}



	fmt.Println(imageClassesLst)

}