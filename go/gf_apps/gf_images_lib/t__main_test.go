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

package gf_images_service

import (
	"os"
	"fmt"
	"testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------
type Gf_test_image_data struct {
	image_client_type_str            string
	image_flows_names_lst            []string
	images_local_filepaths_lst       []string
	local_thumbs_target_dir_path_str string
	origin_url_str                   string
	origin_page_url_str              string
	small_thumb_max_size_px_int      int
	medium_thumb_max_size_px_int     int
	large_thumb_max_size_px_int      int
}

var log_fun func(string, string)
var cli_args_map map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {
	log_fun = gf_core.Init_log_fun()
	cli_args_map = gf_images_core.CLI__parse_args(log_fun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func Test__main(p_test *testing.T) {

	fmt.Println("TEST__MAIN ==============================================")
	
	//-----------------
	// TEST_DATA
	test__image_client_type_str      := "test_run"
	test__image_flows_names_lst      := []string{"test_flow",}
	test__images_local_filepaths_lst := []string{
		"./tests_data/test_image_01.jpeg",
		"./tests_data/test_image_02.jpeg",
		"./tests_data/test_image_03.jpeg",
	}
	test__local_thumbs_target_dir_path_str := "./tests_data/thumbnails"
	test__origin_url_str                   := "https://some_test_domain.com/page_1/test_image.jpeg"
	test__origin_page_url_str              := "https://some_test_domain.com/page_1"
	small_thumb_max_size_px_int            := 200
	medium_thumb_max_size_px_int           := 400
	large_thumb_max_size_px_int            := 600

	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) //"127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"
	

	// RUNTIME
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_images_tests",
		Log_fun:          log_fun,
	}
	
	//-----------------
	// MONGODB	


	// MONGODB
	test__mongodb_url_str     := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)
	mongodb_db, _, gf_err := gf_core.Mongo__connect_new(test__mongodb_url_str, test__mongodb_db_name_str, nil, runtime_sys)
	if gf_err != nil {
		fmt.Println(gf_err.Error)
		p_test.Fail()
	}
	mongodb_coll := mongodb_db.Collection("data_symphony")
	runtime_sys.Mongo_coll = mongodb_coll



	//-----------------

	test_image_data := &Gf_test_image_data{
		image_client_type_str:            test__image_client_type_str,
		image_flows_names_lst:            test__image_flows_names_lst,
		images_local_filepaths_lst:       test__images_local_filepaths_lst,
		local_thumbs_target_dir_path_str: test__local_thumbs_target_dir_path_str,
		origin_url_str:                   test__origin_url_str,
		origin_page_url_str:              test__origin_page_url_str,
		small_thumb_max_size_px_int:      small_thumb_max_size_px_int,
		medium_thumb_max_size_px_int:     medium_thumb_max_size_px_int,
		large_thumb_max_size_px_int:      large_thumb_max_size_px_int,
	}

	test__images_transformer(test_image_data,runtime_sys)
	test__images_ops(test_image_data,runtime_sys)
}

//---------------------------------------------------
func test__images_transformer(p_test_image_data *Gf_test_image_data,
	p_runtime_sys *gf_core.Runtime_sys) {

	fmt.Println("")
	fmt.Println("         TEST__IMAGES_TRANSFORMER   >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("")

	for _, image_local_file_path_str := range p_test_image_data.images_local_filepaths_lst {

		//---------------
		format_str, gf_err := gf_images_core.Get_image_ext_from_url(image_local_file_path_str,p_runtime_sys)
		if gf_err != nil {
			panic(gf_err.Error)
		}

		//---------------
		test__image_id_str := gf_images_core.Image_ID__create(image_local_file_path_str,format_str,p_runtime_sys)
		fmt.Println("test__image_id_str - "+test__image_id_str)

		//---------------

		image_thumbs, gf_image, gf_err := gf_images_core.Trans__process_image(test__image_id_str,
			p_test_image_data.image_client_type_str,
			p_test_image_data.image_flows_names_lst,
			p_test_image_data.origin_url_str,
			p_test_image_data.origin_page_url_str,
			format_str,
			image_local_file_path_str,
			p_test_image_data.local_thumbs_target_dir_path_str,
			p_runtime_sys)

		if gf_err != nil {
			panic(gf_err.Error)	
		}

		spew.Dump(image_thumbs)
		fmt.Println("")

		spew.Dump(gf_image)
		fmt.Println("")
	}
}

//---------------------------------------------------
func test__images_ops(p_test_image_data *Gf_test_image_data,
	p_runtime_sys *gf_core.Runtime_sys) {

	fmt.Println("")
	fmt.Println("         TEST__IMAGES_OPS   >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("")

	for _, test__image_local_filepath_str := range p_test_image_data.images_local_filepaths_lst {

		format_str,gf_err := gf_images_core.Get_image_ext_from_url(test__image_local_filepath_str, p_runtime_sys)
		if gf_err != nil {
			panic(gf_err.Error)
		}

		test__image_ops(p_test_image_data,
			test__image_local_filepath_str,
			format_str,
			p_runtime_sys)
	}
}

//---------------------------------------------------
func test__image_ops(p_test_image_data *Gf_test_image_data,
	p_test__image_local_filepath_str string,
	p_test__image_format_str         string,
	p_runtime_sys                    *gf_core.Runtime_sys) {

	//---------------
	test__image_id_str := gf_images_core.Image_ID__create(p_test__image_local_filepath_str, p_test__image_format_str, p_runtime_sys)
	fmt.Println("test__image_id_str - "+test__image_id_str)

	//---------------
	test__image_title_str, gf_err := gf_images_core.Get_image_title_from_url(p_test__image_local_filepath_str, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	fmt.Println("test__image_title_str - "+test__image_title_str)

	//---------------
	img_width_int, img_height_int, gf_err := gf_images_core.Get_image_dimensions__from_filepath(p_test__image_local_filepath_str, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	fmt.Println(fmt.Sprintf("test__image dimensions - %d/%d",img_width_int,img_height_int))

	//---------------
	img, gf_err := gf_images_core.Image__load_file(p_test__image_local_filepath_str, p_test__image_format_str, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//---------------
	second_img_width_int, second_img_height_int := gf_images_core.Get_image_dimensions__from_image(img, p_runtime_sys)

	if img_width_int != second_img_width_int {
		err_msg_str := "gf_images_core.Get_image_dimensions__from_filepath() and gf_images_core.Get_image_dimensions__from_image() dont return the same width"
		panic(err_msg_str)
	}

	if img_height_int != second_img_height_int {
		err_msg_str := "gf_images_core.Get_image_dimensions__from_filepath() and gf_images_core.Get_image_dimensions__from_image() dont return the same width"
		panic(err_msg_str)
	}

	//---------------
	image_thumbs, gf_err := gf_images_core.Create_thumbnails(test__image_id_str,
		p_test__image_format_str,
		p_test__image_local_filepath_str,
		p_test_image_data.local_thumbs_target_dir_path_str,
		p_test_image_data.small_thumb_max_size_px_int,
		p_test_image_data.medium_thumb_max_size_px_int,
		p_test_image_data.large_thumb_max_size_px_int,
		img,
		p_runtime_sys)

	if gf_err != nil {
		panic(gf_err.Error)
	}

	fmt.Println("")
	fmt.Println("              NEW_TEST_THUMBS >>>>>>>>>>>>>>>")
	fmt.Println("")

	spew.Dump(image_thumbs)

	fmt.Println("")
	fmt.Println("")

	//---------------

	image_new_info := &gf_images_core.Gf_image_new_info{
		Id_str:                         test__image_id_str,
		Title_str:                      test__image_title_str,
		Flows_names_lst:                p_test_image_data.image_flows_names_lst,
		Image_client_type_str:          p_test_image_data.image_client_type_str,
		Origin_url_str:                 p_test_image_data.origin_url_str,
		Origin_page_url_str:            p_test_image_data.origin_page_url_str,
		Original_file_internal_uri_str: p_test__image_local_filepath_str,
		Thumbnail_small_url_str:        image_thumbs.Small_relative_url_str,
		Thumbnail_medium_url_str:       image_thumbs.Medium_relative_url_str,
		Thumbnail_large_url_str:        image_thumbs.Large_relative_url_str,
	 	Format_str:                     p_test__image_format_str,
	 	Width_int:                      img_width_int,
	 	Height_int:                     img_height_int,
	}


	gf_image, gf_err := gf_images_core.Image__create_new(image_new_info, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	fmt.Println("")
	fmt.Println("              NEW_TEST_IMAGE >>>>>>>>>>>>>>>")
	fmt.Println("")

	spew.Dump(gf_image)

	fmt.Println("")
	fmt.Println("")
}