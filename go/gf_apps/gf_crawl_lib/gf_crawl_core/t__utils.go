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

package gf_crawl_core

import (
	"os"
	"testing"
	"fmt"
	"os/exec"
	"path/filepath"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)

//-------------------------------------------------
// INIT
func T__init() (*gf_core.Runtime_sys, *Gf_crawler_runtime) {

	//-------------
	// MONGODB
	test__mongodb_host_str      := "127.0.0.1"
	test__mongodb_db_name_str   := "gf_tests"

	// MONGODB_ENV
	test__mongodb_host_env_str    := os.Getenv("GF_MONGODB_HOST")
	test__mongodb_db_name_env_str := os.Getenv("GF_MONGODB_DB_NAME")

	if test__mongodb_host_env_str != "" {
		test__mongodb_host_str = test__mongodb_host_env_str
	}

	if test__mongodb_db_name_env_str != "" {
		test__mongodb_host_str = test__mongodb_db_name_env_str
	}

	// ELASTICSEARCH
	test__es_host_str := "127.0.0.1:9200"
	
	// ELASTICSEARCH_ENV
	test__es_host_env_str := os.Getenv("GF_ELASTICSEARCH_HOST")

	if test__es_host_env_str != "" {
		test__es_host_str = test__es_host_env_str
	}

	test__cluster_node_type_str := "master"
	//-------------

	log_fun      := gf_core.Init_log_fun()
	mongodb_db   := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_crawl_tests",
		Log_fun:          log_fun,
		Mongodb_db:       mongodb_db,
		Mongodb_coll:     mongodb_coll,
	}
	//-------------
	//ELASTICSEARCH
	esearch_client, gf_err := gf_core.Elastic__get_client(test__es_host_str, runtime_sys)
	if gf_err != nil {
		panic("failed to get ElasticSearch client in test initialization")
		return nil, nil
	}
	//-------------
	//S3
	s3_test_info := gf_core.T__get_s3_info(runtime_sys)
	//-------------

	crawler_runtime := &Gf_crawler_runtime{
		Events_ctx:            nil,
		Esearch_client:        esearch_client,
		S3_info:               s3_test_info.Gf_s3_info,
		Cluster_node_type_str: test__cluster_node_type_str,
	}

	return runtime_sys, crawler_runtime
}

//---------------------------------------------------
func t__create_test_image_ADTs(p_test *testing.T,
	p_test__crawler_name_str    string,
	p_test__cycle_run_id_str    string,
	p_test__img_src_url_str     string,
	p_test__origin_page_url_str string,
	p_crawler_runtime           *Gf_crawler_runtime,
	p_runtime_sys               *gf_core.Runtime_sys) (*Gf_crawler_page_image, *Gf_crawler_page_image_ref) {

	//-------------------
	//CRAWLED_IMAGE_CREATE
	test__crawled_image, gf_err := images_adt__prepare_and_create(p_test__crawler_name_str,
		p_test__cycle_run_id_str,
		p_test__img_src_url_str,
		p_test__origin_page_url_str,
		p_crawler_runtime,
		p_runtime_sys)
	if gf_err != nil { 
		p_test.Errorf("failed to prepare and create image_adt with URL [%s] and origin_page URL [%s]", p_test__img_src_url_str, p_test__origin_page_url_str)
		return nil, nil
	}

	//DB - CRAWLED_IMAGE_PERSIST
	exists_bool, gf_err := Image__db_create(test__crawled_image, p_crawler_runtime, p_runtime_sys)
	if gf_err != nil {
		p_test.Errorf("failed to DB persist image_adt with URL [%s] and origin_page URL [%s]", p_test__img_src_url_str, p_test__origin_page_url_str)
		return nil, nil
	}

	assert.Equal(p_test, exists_bool, false, "test page_image exists in the DB already, test cleanup hasnt been done")
	//-------------------
	//CRAWLED_IMAGE_REF_CREATE
	test__crawled_image_ref := images_adt__ref_create(p_test__crawler_name_str,
		p_test__cycle_run_id_str,
		test__crawled_image.Url_str,                    //p_image_url_str
		test__crawled_image.Domain_str,                 //p_image_url_domain_str
		test__crawled_image.Origin_page_url_str,        //p_origin_page_url_str
		test__crawled_image.Origin_page_url_domain_str, //p_origin_page_url_domain_str
		p_runtime_sys)

	//DB - CRAWLED_IMAGE_REF_PERSIST
	gf_err = Image__db_create_ref(test__crawled_image_ref, p_crawler_runtime, p_runtime_sys)
	if gf_err != nil {
		p_test.Errorf("failed to DB persist image_ref for image with URL [%s] and origin_page URL [%s]", p_test__img_src_url_str, p_test__origin_page_url_str)
		return nil, nil
	}
	//-------------------

	return test__crawled_image, test__crawled_image_ref
}

//-------------------------------------------------
//given some human readable (or arbitrarily named) local image file, create a new image file with the same content, that is named
//according to the gf_images image file naming scheme. here for testing this is done manually via this function, but in the gf_crawl pipeline
//this is done by calling the native gf_image functions that create this gf_images based name.

func t__create_test_gf_image_named_image_file(p_test *testing.T,
	p_test__img_src_url_str           string,
	p_test__local_image_file_path_str string,
	p_runtime_sys                     *gf_core.Runtime_sys) (string, gf_images_utils.Gf_image_id) {
	p_runtime_sys.Log_fun("FUN_ENTER", "t__utils.t__create_test_gf_image_named_image_file()")

	test__local_image_dir_path_str := filepath.Dir(p_test__local_image_file_path_str)

	//IMPORTANT!! - creates a new gf_image ID from the image URL
	test__local_gf_image_file_path_str, gf_image_id_str, gf_err := gf_images_utils.Create_gf_image_file_path_from_url("", p_test__img_src_url_str,
		test__local_image_dir_path_str,
		p_runtime_sys)
	if gf_err != nil {
		p_test.Errorf(fmt.Sprintf("failed to create a gf_image local file path from URL [%s]", p_test__img_src_url_str))
		return "", ""
	}

	source_abs_str, _ := filepath.Abs(p_test__local_image_file_path_str)
	target_abs_str, _ := filepath.Abs(test__local_gf_image_file_path_str)

	err := exec.Command("cp", source_abs_str, target_abs_str).Run()
	if err != nil {
		p_test.Errorf(fmt.Sprintf("failed to copy a image file via shell, from old path [%s] to new gf_image path [%s]", source_abs_str, target_abs_str))
		return "", ""
	}
	return test__local_gf_image_file_path_str, gf_image_id_str
}

//-------------------------------------------------
func t__cleanup__test_page_imgs(p_test__crawler_name_str string, p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "t__utils.t__cleanup__test_page_imgs()")

	_,err := p_runtime_sys.Mongodb_db.C("gf_crawl").RemoveAll(bson.M{
			"t":                bson.M{"$in": []string{"crawler_page_img", "crawler_page_img_ref",},},
			"crawler_name_str": p_test__crawler_name_str,
		})
	if err != nil {
		panic(err)
	}
}