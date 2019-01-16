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
	"os"
	"fmt"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
	//"github.com/davecgh/go-spew/spew"
)
//-------------------------------------------------
func Test__main(p_test *testing.T) {

	//test__images_s3_bucket_name_str := "gf--test"


	test__mongodb_host_str    := "127.0.0.1"
	test__mongodb_db_name_str := "test_db"
	
	test_post_info_map := map[string]interface{}{
		
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
	aws_access_key_id_str     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	aws_secret_access_key_str := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	aws_token_str             := os.Getenv("GF_AWS_TOKEN")

	if aws_access_key_id_str == "" || aws_secret_access_key_str == "" {
		panic("test AWS credentials were not supplied")
	}
	
	s3_info, gf_err := gf_core.S3__init(aws_access_key_id_str, aws_secret_access_key_str, aws_token_str, runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	fmt.Println(s3_info)
}
//-------------------------------------------------
func test_posts_creation() {


	Pipeline__create_post(p_post_info_map map[string]interface{},




}