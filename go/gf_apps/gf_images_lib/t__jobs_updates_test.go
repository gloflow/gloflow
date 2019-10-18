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
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
)
//---------------------------------------------------
func Test__jobs_updates(p_test *testing.T) {

	log_fun      := gf_core.Init_log_fun()
	cli_args_map := gf_images_utils.CLI__parse_args(log_fun)
	//-----------------
	//TEST_DATA
	
	test__http_server_host_str             := "localhost:8000"
	test__gf_images_service_port_str       := "8010"
	test__images_local_dir_path_str        := "./tests_data"
	test__images_thumbs_local_dir_path_str := "./tests_data/thumbnails"
	test__s3_bucket_name_str               := "gf--test--img"
	
	test__image_flows_names_lst            := []string{"test_flow",}
	test__image_url_str                    := fmt.Sprintf("http://%s/test_image_01.jpeg", test__http_server_host_str)
	test__origin_page_url_str              := "https://some_test_domain.com/page_1"
	test__image_client_type_str            := "test_run"
	
	test__service_templates_dir_paths_map  := map[string]interface{}{
		"flows_str":"./templates",
	}

	// MONGODB
	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) //"127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"

	fmt.Println(fmt.Sprintf("test__http_server_host_str       - %s", test__http_server_host_str))
	fmt.Println(fmt.Sprintf("test__gf_images_service_port_str - %s", test__gf_images_service_port_str))
	fmt.Println("")
	//-------------
	
	mongodb_db   := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str:"gf_images_tests",
		Log_fun:         log_fun,
		Mongodb_coll:    mongodb_coll,
	}
	//-------------
	//S3
	gf_s3_test_info := gf_core.T__get_s3_info(runtime_sys)
	//-------------
	//JOBS_MNGR
	jobs_mngr := gf_images_jobs.Jobs_mngr__init(test__images_local_dir_path_str,
		test__images_thumbs_local_dir_path_str,
		test__s3_bucket_name_str,
		gf_s3_test_info.Gf_s3_info,
		runtime_sys)
	//-------------
	//START_HTTP_SERVICE
	done_ch := make(chan bool)
	go func() {
		Run_service(test__gf_images_service_port_str,
			test__mongodb_host_str,
			test__mongodb_db_name_str,
			test__images_local_dir_path_str,
			test__images_thumbs_local_dir_path_str,
			test__s3_bucket_name_str,
			gf_s3_test_info.Aws_access_key_id_str,
			gf_s3_test_info.Aws_secret_access_key_str,
			gf_s3_test_info.Aws_token_str,
			test__service_templates_dir_paths_map,
			done_ch,
			log_fun)
	}()
	<-done_ch //wait for the service to finish initializing
	//-------------
	//HTTP
	test__job_updates__via_http(test__image_url_str,
		test__origin_page_url_str,
		test__image_client_type_str,
		test__gf_images_service_port_str,
		runtime_sys)

	//IN_PROCESS
	test__job_updates__in_process(test__image_url_str,
		test__image_flows_names_lst,
		test__image_client_type_str,
		test__origin_page_url_str,
		jobs_mngr,
		runtime_sys)
	//-------------
}
//---------------------------------------------------
func test__job_updates__via_http(p_test__image_url_str string,
	p_test__origin_page_url_str   string,
	p_test__image_client_type_str string,
	p_test_image_service_port_str string,
	p_runtime_sys                 *gf_core.Runtime_sys) {
	
	test__input_images_urls_lst                    := []string{p_test__image_url_str,}
	test__input_images_origin_pages_urls_lst       := []string{p_test__origin_page_url_str,}
	test__image_service_host_port_str              := fmt.Sprintf("localhost:%s", p_test_image_service_port_str)
	running_job_id_str, images_outputs_lst, gf_err := Client__dispatch_process_extern_images(test__input_images_urls_lst,
		test__input_images_origin_pages_urls_lst,
		p_test__image_client_type_str,
		test__image_service_host_port_str,
		p_runtime_sys)

	if gf_err != nil {
		panic(gf_err.Error)
	}


	fmt.Println(running_job_id_str)
	fmt.Println(images_outputs_lst)

	
}
//---------------------------------------------------
func test__job_updates__in_process(p_test__image_url_str string,
	p_test__image_flows_names_lst []string,
	p_test__image_client_type_str string,
	p_test__origin_page_url_str   string,
	p_jobs_mngr                   gf_images_jobs.Jobs_mngr,
	p_runtime_sys                 *gf_core.Runtime_sys) {

	images_to_process_lst := []gf_images_jobs.Image_to_process{
		gf_images_jobs.Image_to_process{
			Source_url_str:     p_test__image_url_str,
			Origin_page_url_str:p_test__origin_page_url_str,
		},
	}

	running_job, output, gf_err := gf_images_jobs.Job__start(p_test__image_client_type_str,
		images_to_process_lst,
		p_test__image_flows_names_lst,
		p_jobs_mngr,
		p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	fmt.Println(running_job)
	spew.Dump(output)

	T__test_image_job__updates(running_job.Id_str, p_jobs_mngr, p_runtime_sys)
}