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
	// "github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)

//---------------------------------------------------
func Test__upload(p_test *testing.T) {

	fmt.Println("TEST__UPLOAD ==============================================")

	//-----------------
	// TEST_DATA

	// MONGODB
	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"
	
	test__image_name_str        := "test_image"
	test__image_flows_names_lst := []string{"test_flow",}
	test__config_file_path_str  := "./test_config/gf_images_config.yaml"

	test__image_local_file_path_lst := "./tests_data/test_image_02.jpeg"

	test__images_local_dir_path_str        := "./tests_data"
	test__images_thumbs_local_dir_path_str := "./tests_data/thumbnails"
	// test__s3_bucket_name_str               := "gf--test--img"
	//-------------
	
	mongodb_db   := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_images_tests",
		Log_fun:          log_fun,
		Mongodb_db:       mongodb_db,
		Mongodb_coll:     mongodb_coll,
	}

	// CONFIG
	img_config, gf_err := gf_images_utils.Config__get(test__config_file_path_str, runtime_sys)
	if gf_err != nil {
		return
	}
	
	//-------------
	// S3
	gf_s3_test_info := gf_core.T__get_s3_info(runtime_sys)
	fmt.Println(gf_s3_test_info)

	//-------------
	// JOBS_MANAGER
	jobs_mngr := gf_images_jobs.Jobs_mngr__init(test__images_local_dir_path_str,
		test__images_thumbs_local_dir_path_str,
		img_config,
		gf_s3_test_info.Gf_s3_info,
		runtime_sys)

	//-------------
	// UPLOAD_INIT
	upload_info, gf_err := Upload__init(test__image_name_str,
		"jpeg", // image_format_str,
		test__image_flows_names_lst,
		"browser", // client_type_str,
		gf_s3_test_info.Gf_s3_info,
		img_config,
		runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	spew.Dump(upload_info)

	assert.Equal(p_test, upload_info.T_str, "img_upload_info", "upload_info doesnt have a proper type set")
	assert.Equal(p_test, len(upload_info.Upload_gf_image_id_str), 32, "S3 image upload_info gf_image_id is not of correct length - 32 chars")
	assert.NotEqual(p_test, upload_info.Presigned_url_str, "", "S3 image presigned_url is not set")

	//-------------
	// DB
	db_upload_info, gf_err := Upload_db__get_info(upload_info.Upload_gf_image_id_str, runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	fmt.Println("done DB READING")

	assert.Equal(p_test, upload_info.Creation_unix_time_f, db_upload_info.Creation_unix_time_f,
		"upload_info struct that was created in memory and the one refetched from the DB (after the memory one was persisted) dont have the same creation time ")

	//-------------
	// S3_UPLOAD - HTTP PUT
	headers_map := map[string]string{
		"Content-Type": "image/jpeg",
		"x-amz-acl":    "public-read",
	}

	resp, gf_err := gf_core.HTTP__put_file(upload_info.Presigned_url_str,
		test__image_local_file_path_lst,
		headers_map,
		runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)

	fmt.Printf("HTTP PUT COMPLETE - %s\n", test__image_local_file_path_lst)

	assert.Equal(p_test, resp.StatusCode, 200,
		"image upload to S3 via presigned-url failed")
		
	//-------------
	// UPLOAD_COMPLETE
	running_job, gf_err := Upload__complete(upload_info.Upload_gf_image_id_str,
		jobs_mngr,
		gf_s3_test_info.Gf_s3_info,
		runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	T__test_image_job__updates(running_job.Id_str, jobs_mngr, runtime_sys)
	//-------------
}