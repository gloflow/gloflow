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
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__img_add_to_flow(p_test *testing.T) {


	//IMPORTANT!! - in this test there is no downloading of a file. a gf_page_img__pipeline_info reference is created manually
	//              with a local image file path set manually. this local path is the path of the test image (test__local_image_file_path_str).
	//              crawler image ADT's are manually created first 

	//-------------------
	test__crawler_name_str                  := "test-crawler"
	test__cycle_run_id_str                  := "test__cycle_run_id"
	test__image_flows_names_lst             := []string{"test_flow",}
	test__img_src_url_str                   := "/some/origin/test_image_wasp.jpeg"
	test__origin_page_url_str               := "/some/origin/url.html"
	test__local_image_file_path_str         := "../test_data/test_image_wasp.jpeg"
	test__images_store_local_dir_path_str   := "../test_data/processed_images" //image tmp thumbnails, or downloaded gif's and their frames
	test__crawled_images_s3_bucket_name_str := "gf--test--discovered--img"
	test__gf_images_s3_bucket_name_str      := "gf--test--img"
	runtime_sys, crawler_runtime            := T__init()

	t__cleanup__test_page_imgs(test__crawler_name_str, runtime_sys)

	//---------------------------------------------------
	create_image_ADTs := func() (*Gf_crawler_page_img, *Gf_crawler_page_img_ref) {
		//-------------------
		//CRAWLED_IMAGE_CREATE
		test__crawled_image, gf_err := images_adt__prepare_and_create(test__crawler_name_str,
			test__cycle_run_id_str,
			test__img_src_url_str,
			test__origin_page_url_str,
			crawler_runtime,
			runtime_sys)
		if gf_err != nil { 
			panic(gf_err.Error)
		}

		//DB - CRAWLED_IMAGE_PERSIST
		exists_bool, gf_err := Image__db_create(test__crawled_image, crawler_runtime, runtime_sys)
		if gf_err != nil {
			panic(gf_err.Error)
		}

		assert.Equal(p_test, exists_bool, false, "test page_image exists in the DB already, test cleanup hasnt been done")
		//-------------------
		//CRAWLED_IMAGE_REF_CREATE
		test__crawled_image_ref := images_adt__ref_create(test__crawler_name_str,
			test__cycle_run_id_str,
			test__crawled_image.Url_str,                    //p_image_url_str
			test__crawled_image.Domain_str,                 //p_image_url_domain_str
			test__crawled_image.Origin_page_url_str,        //p_origin_page_url_str
			test__crawled_image.Origin_page_url_domain_str, //p_origin_page_url_domain_str
			runtime_sys)

		//DB - CRAWLED_IMAGE_REF_PERSIST
		gf_err = Image__db_create_ref(test__crawled_image_ref, crawler_runtime, runtime_sys)
		if gf_err != nil {
			panic(gf_err.Error)
		}
		//-------------------

		return test__crawled_image, test__crawled_image_ref
	}
	//---------------------------------------------------
	test__crawled_image, test__crawled_image_ref := create_image_ADTs()

	//-------------------
	//PIPELINE_STAGE__PROCESS_IMAGES - apply image transformations, create thumbnails, etc.

	//GF_PAGE_IMAGE_LINK
	page_img_link := &gf_page_img_link{
		img_src_str:         test__img_src_url_str,
		origin_page_url_str: test__origin_page_url_str,
	}

	//GF_PAGE_IMAGE__PIPELINE_INFO - this is the struct thats passed through the crawler image processing pipeline, 
	//                               from stage to stage. here we're createing it manually and populating with test values. 
	page_img__pipeline_info := &gf_page_img__pipeline_info{
		link:                page_img_link,
		page_img:            test__crawled_image,
		page_img_ref:        test__crawled_image_ref,
		exists_bool:         false,               //artificially set test image to be declared as not existing already, in order to be fully processed
		local_file_path_str: test__local_image_file_path_str,
		nsfv_bool:           false,
		thumbs:              nil,
	}

	page_imgs__pinfos_lst := []*gf_page_img__pipeline_info{
		page_img__pipeline_info,
	}

	page_imgs__pinfos_with_thumbs_lst := images__stage__process_images(test__crawler_name_str,
		page_imgs__pinfos_lst,
		test__images_store_local_dir_path_str,
		test__origin_page_url_str,
		test__crawled_images_s3_bucket_name_str,
		crawler_runtime,
		runtime_sys)

	fmt.Println("   STAGE_COMPLETE --------------")

	assert.Equal(p_test, len(page_imgs__pinfos_lst), len(page_imgs__pinfos_with_thumbs_lst), "more page_imgs pipeline_info's returned from images__stage__process_images() then inputed")
	assert.Equal(p_test, len(page_imgs__pinfos_lst), len(page_imgs__pinfos_with_thumbs_lst), "more page_imgs pipeline_info's returned from images__stage__process_images() then inputed")

	for _, page_img__pinfo := range page_imgs__pinfos_with_thumbs_lst {

		fmt.Printf("  ------- page_img__pinfo")
		spew.Dump(page_img__pinfo)

		assert.Equal(p_test, page_img__pinfo.page_img.S3_stored_bool, false)

		if page_img__pinfo.thumbs == nil {
			panic("page_img.thumbs has not been set to a gf_images_utils.Gf_image_thumbs instance pointer")
		}
	}

	//------------------
	//PIPELINE_STAGE__S3_STORE_IMAGES

	page_imgs__pinfos_with_s3_lst := images_s3__stage__store_images(test__crawler_name_str,
		page_imgs__pinfos_with_thumbs_lst,
		test__origin_page_url_str,
		test__crawled_images_s3_bucket_name_str,
		crawler_runtime,
		runtime_sys)

	fmt.Println("   STAGE_COMPLETE --------------")

	spew.Dump(page_imgs__pinfos_with_s3_lst)

	for _, page_img__pinfo := range page_imgs__pinfos_with_s3_lst {

		spew.Dump(page_img__pinfo)

		assert.Equal(p_test, page_img__pinfo.page_img.S3_stored_bool, true)
	}
	//-------------------
	//FLOWS__ADD_EXTERN_IMAGE - copying files from one FS location to another (S3 bucket to another)



	fmt.Printf("+++++++++++++++++++++++++++++++++++++")
	spew.Dump(test__crawled_image)
	

	gf_err := Flows__add_extern_image(test__crawled_image.Id_str,
		test__image_flows_names_lst,
		test__crawled_images_s3_bucket_name_str,
		test__gf_images_s3_bucket_name_str,
		crawler_runtime,
		runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	//-------------------
}