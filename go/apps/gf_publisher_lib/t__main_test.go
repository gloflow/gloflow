/*
GloFlow media management/publishing system
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

package gf_publisher_lib

import (
	"fmt"
	"testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_jobs"
)
//-------------------------------------------------
func Test__main(p_test *testing.T) {

	test__http_server_host_str      := "localhost:8000"
	test__mongodb_host_str          := "127.0.0.1"
	test__mongodb_db_name_str       := "gf_tests"
	test__images_local_dir_path_str := "./tests_data"
	test__images_thumbs_local_dir_path_str := "./tests_data/thumbnails"
	test__s3_bucket_name_str               := "gf--test--img"

	test_post_info_map := map[string]interface{}{
		"client_type_str":     "test_run",
		"title_str":           "test title",
		"description_str":     "some test description",
		"tags_str":            "tag1,tag2,tag3",
		"poster_user_name_str":"test_user",
		"post_elements_lst":   []interface{}{
			map[string]interface{}{
				"type_str":           "link",
				"extern_url_str":     fmt.Sprintf("http://%s/test_image_01.jpeg",test__http_server_host_str),
				"origin_page_url_str":"http://origin.com/page/url", 
				"tags_lst":           []string{"tag1","tag2"},
			},
			map[string]interface{}{
				"type_str":           "image",
				"extern_url_str":     fmt.Sprintf("http://%s/test_image_02.jpeg",test__http_server_host_str),
				"origin_page_url_str":"http://origin.com/page/url",
				"tags_lst":           []string{"tag1","tag2"},
			},
			map[string]interface{}{
				"type_str":           "video",
				"extern_url_str":     fmt.Sprintf("http://%s/test_image_03.jpeg",test__http_server_host_str),
				"origin_page_url_str":"http://origin.com/page/url",
				"tags_lst":           []string{"tag1","tag2"},
			},
		},
	}

	log_fun      := gf_core.Init_log_fun()
	mongodb_db   := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str:"gf_publisher_tests",
		Log_fun:         log_fun,
		Mongodb_coll:    mongodb_coll,
	}
	//-------------
	//S3
	s3_info := gf_core.T__get_s3_info(runtime_sys)
	//-------------
	//GF_IMAGES_LIB JOBS_MNGR
	jobs_mngr := gf_images_jobs.Jobs_mngr__init(test__images_local_dir_path_str,
		test__images_thumbs_local_dir_path_str,
		test__s3_bucket_name_str,
		s3_info,
		runtime_sys)

	gf_images_runtime_info := &Gf_images_extern_runtime_info{
		Jobs_mngr:            jobs_mngr, //use jobs_mngr thats running in the same process
		Service_host_port_str:"",        //setting this to "" causes jobs_mngr to not issue job requests over HTTP
	}
	//-------------

	test_posts_creation(test_post_info_map, gf_images_runtime_info, runtime_sys)
}
//-------------------------------------------------
func test_posts_creation(p_test_post_info_map map[string]interface{},
	p_gf_images_runtime_info *Gf_images_extern_runtime_info,
	p_runtime_sys            *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","t__main_test.test_posts_creation()")

	
	


	gf_post, images_job_id_str, gf_err := Pipeline__create_post(p_test_post_info_map, p_gf_images_runtime_info, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}


	fmt.Printf("images_job_id_str - %s\n", images_job_id_str)
	spew.Dump(gf_post)



	job_updates_ch := gf_images_jobs.Job__get_update_ch(images_job_id_str, p_gf_images_runtime_info.Jobs_mngr, p_runtime_sys)

	for ;; {

		fmt.Println("\n\n------------------------- TESTING - GET_JOB_UPDATE -----")
		update := <-job_updates_ch
		fmt.Println(update)
	}

}