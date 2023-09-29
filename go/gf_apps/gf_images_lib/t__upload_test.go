/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_images_lib

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------

func TestUpload(pTest *testing.T) {

	fmt.Println("TEST__UPLOAD ==============================================")

	//-----------------
	// TEST_DATA

	// MONGODB
	test__mongodb_host_str    := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_url_str     := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)
	test__mongodb_db_name_str := "gf_tests"
	
	testImageNameStr       := "test_image"
	testImageFlowsNamesLst := []string{"test_flow",}
	testConfigFilePathStr  := "./test_config/gf_images_config.yaml"

	testImageLocalFilePathLst := "./tests_data/test_image_02.jpeg"

	test__images_local_dir_path_str        := "./tests_data"
	test__images_thumbs_local_dir_path_str := "./tests_data/thumbnails"
	testVideosLocalDirPathStr              := "./tests_data/videos"
	// test__s3_bucket_name_str               := "gf--test--img"

	userID := "test_user"

	//-------------
	
	mongodbDB, _, gfErr := gf_core.MongoConnectNew(test__mongodb_url_str, test__mongodb_db_name_str, logFun)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.FailNow()
	}
	mongodbColl := mongodbDB.Collection("data_symphony")
	
	runtime_sys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_images_tests",
		LogFun:           logFun,
		Mongodb_db:       mongodbDB,
		Mongodb_coll:     mongodbColl,
	}

	// CONFIG
	img_config, gfErr := gf_images_core.Config__get(testConfigFilePathStr, runtime_sys)
	if gfErr != nil {
		return
	}
	
	//-------------
	// S3
	s3testInfo := gf_core.T__get_s3_info(runtime_sys)
	fmt.Println(s3testInfo)

	//-------------
	// JOBS_MANAGER
	jobsMngr := gf_images_jobs.Init(test__images_local_dir_path_str,
		test__images_thumbs_local_dir_path_str,
		testVideosLocalDirPathStr,
		img_config,
		s3testInfo.Gf_s3_info,
		runtime_sys)

	//-------------
	// UPLOAD_INIT
	uploadInfo, gfErr := Upload__init(testImageNameStr,
		"jpeg", // image_format_str,
		testImageFlowsNamesLst,
		"browser", // client_type_str,
		userID,
		s3testInfo.Gf_s3_info,
		img_config,
		runtime_sys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	spew.Dump(uploadInfo)

	assert.Equal(pTest, uploadInfo.T_str, "img_upload_info", "upload_info doesnt have a proper type set")
	assert.Equal(pTest, len(uploadInfo.Upload_gf_image_id_str), 32, "S3 image upload_info gf_image_id is not of correct length - 32 chars")
	assert.NotEqual(pTest, uploadInfo.Presigned_url_str, "", "S3 image presigned_url is not set")

	//-------------
	// DB
	db_upload_info, gfErr := Upload_db__get_info(uploadInfo.Upload_gf_image_id_str, runtime_sys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	fmt.Println("done DB READING")

	assert.Equal(pTest, uploadInfo.Creation_unix_time_f, db_upload_info.Creation_unix_time_f,
		"upload_info struct that was created in memory and the one refetched from the DB (after the memory one was persisted) dont have the same creation time ")

	//-------------
	// S3_UPLOAD - HTTP PUT
	headers_map := map[string]string{
		"Content-Type": "image/jpeg",
		"x-amz-acl":    "public-read",
	}

	resp, gfErr := gf_core.HTTP__put_file(uploadInfo.Presigned_url_str,
		testImageLocalFilePathLst,
		headers_map,
		runtime_sys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)

	fmt.Printf("HTTP PUT COMPLETE - %s\n", testImageLocalFilePathLst)

	assert.Equal(pTest, resp.StatusCode, 200,
		"image upload to S3 via presigned-url failed")
		
	//-------------
	// UPLOAD_COMPLETE
	running_job, gfErr := Upload__complete(uploadInfo.Upload_gf_image_id_str,
		jobsMngr,
		userID,
		s3testInfo.Gf_s3_info,
		runtime_sys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	T__test_image_job__updates(running_job.Id_str, jobsMngr, runtime_sys)
	
	//-------------
}