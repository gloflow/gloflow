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

package gf_images_lib

import (
	"fmt"
	"testing"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
)

//---------------------------------------------------

func TstartService(pMongoHostStr string,
	pTest   *testing.T,
	pLogFun func(string, string)) {

	//-----------------
	// TEST_DATA
	
	test__http_server_host_str             := "localhost:8000"
	testPortStr                            := "8010"
	test__images_local_dir_path_str        := "./tests_data"
	test__images_thumbs_local_dir_path_str := "./tests_data/thumbnails"
	testServiceTemplatesDirPathsMap  := map[string]string{
		"flows_str": "./../../../web/src/gf_apps/gf_images/templates",
	}

	testConfigFilePathStr := "./test_config/gf_images_config.yaml"
	
	// test__origin_page_url_str   := "https://some_test_domain.com/page_1"
	// test__image_client_type_str := "test_run"
	// test__s3_bucket_name_str    := "gf--test--img"
	// test__image_flows_names_lst := []string{"test_flow",}
	// test__image_url_str         := fmt.Sprintf("http://%s/test_image_01.jpeg", test__http_server_host_str)

	//-------------
	// MONGODB
	testMongoHostStr   := pMongoHostStr // cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	testMongoURLstr    := fmt.Sprintf("mongodb://%s", testMongoHostStr)
	testMongoDBnameStr := "gf_tests"

	fmt.Println(fmt.Sprintf("test__http_server_host_str - %s", test__http_server_host_str))
	fmt.Println(fmt.Sprintf("test_port                  - %s", testPortStr))
	fmt.Println("")
	
	//-------------
	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_images_tests",
		LogFun:         pLogFun,
	}

	mongodbDB, _, gfErr := gf_core.MongoConnectNew(testMongoURLstr,
		testMongoDBnameStr,
		nil,
		runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()
	}

	mongodbColl := mongodbDB.Collection("data_symphony")
	
	runtimeSys.Mongo_db   = mongodbDB
	runtimeSys.Mongo_coll = mongodbColl

	// CONFIG
	useNewStorageEngineBool := false
	IPFSnodeHostStr := ""
	config, gfErr := gf_images_core.ConfigGet(testConfigFilePathStr,
		useNewStorageEngineBool,
		IPFSnodeHostStr,
		runtimeSys)
	if gfErr != nil {
		return
	}

	//-------------
	// S3
	// s3testInfo := gf_aws.TgetS3info(runtimeSys)

	//-------------
	// START_HTTP_SERVICE
	
	serviceInfo := &gf_images_core.GFserviceInfo{
		Port_str:                             testPortStr,
		Mongodb_host_str:                     testMongoHostStr,
		Mongodb_db_name_str:                  testMongoDBnameStr,
		ImagesStoreLocalDirPathStr:           test__images_local_dir_path_str,
		ImagesThumbnailsStoreLocalDirPathStr: test__images_thumbs_local_dir_path_str,
		Images_main_s3_bucket_name_str:       config.Main_s3_bucket_name_str, // test__s3_bucket_name_str,
		TemplatesPathsMap:                    testServiceTemplatesDirPathsMap,
		Config_file_path_str:                 testConfigFilePathStr,
	}

	doneCh := make(chan bool)
	go func() {
		
		// HTTP_MUX
		serviceHTTPmux := http.NewServeMux()

		RunService(serviceHTTPmux,
			serviceInfo,
			doneCh,
			pLogFun)
	}()
	<-doneCh // wait for the service to finish initializing
}

//---------------------------------------------------

func T__test_image_job__updates(pJobIDstr string,
	pJobsMngr   gf_images_jobs_core.JobsMngr,
	pRuntimeSys *gf_core.RuntimeSys) {

	//-------------
	// TEST_JOB_UPDATES
	jobUpdatesCh := gf_images_jobs_client.GetJobUpdateCh(pJobIDstr, pJobsMngr, pRuntimeSys)

	for ;; {

		fmt.Println("\n\n------------------------- TESTING - GET_JOB_UPDATE -----")
		jobUpdate := <- jobUpdatesCh

		spew.Dump(jobUpdate)

		jobUpdateTypeStr := jobUpdate.Type_str
		if jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__ERROR {
			panic("job encountered an error while processing")
		}

		if !(jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__OK || jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__COMPLETED) {
			panic(fmt.Sprintf("job_update is expected to be of type 'ok' but instead is - %s", jobUpdateTypeStr))
		}
		
		// test complete
		if jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__COMPLETED {
			break
		}
	}

	//-------------
}