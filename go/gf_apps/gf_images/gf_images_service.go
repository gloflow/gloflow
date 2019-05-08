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

package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
)

//-------------------------------------------------
func main() {
	log_fun := gf_core.Init_log_fun()

	cli_args_map            := parse__cli_args(log_fun)
	run__start_service_bool := cli_args_map["run__start_service_bool"].(bool)
	port_str                := cli_args_map["port_str"].(string)
	mongodb_host_str        := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str     := cli_args_map["mongodb_db_name_str"].(string)
	images_store_local_dir_path_str            := cli_args_map["images_store_local_dir_path_str"].(string)
	images_thumbnails_store_local_dir_path_str := cli_args_map["images_thumbnails_store_local_dir_path_str"].(string)
	images_main_s3_bucket_name_str             := cli_args_map["images_s3_bucket_name_str"].(string)
	aws_access_key_id_str                      := cli_args_map["aws_access_key_id_str"].(string)
	aws_secret_access_key_str                  := cli_args_map["aws_secret_access_key_str"].(string)
	aws_token_str                              := cli_args_map["aws_token_str"].(string)

	templates_dir_paths_map := map[string]interface{}{
		"flows_str": "./templates",
		"gif_str":   "./templates",
	}

	//fmt.Println("AWS------------------------------------------")
	//fmt.Println(aws_access_key_id_str)
	//fmt.Println(aws_secret_access_key_str)
	//fmt.Println(aws_token_str)

	//START_SERVICE
	if run__start_service_bool {

		gf_images_lib.Run_service(port_str,
			mongodb_host_str,
			mongodb_db_name_str,
			images_store_local_dir_path_str,
			images_thumbnails_store_local_dir_path_str,
			images_main_s3_bucket_name_str,
			aws_access_key_id_str,
			aws_secret_access_key_str,
			aws_token_str,
			templates_dir_paths_map,
			nil, //init_done_ch,
			log_fun)
	}
}

//-------------------------------------------------
func parse__cli_args(p_log_fun func(string,string)) map[string]interface{} {
	p_log_fun("FUN_ENTER", "gf_images_service.parse__cli_args()")

	//-------------------
	//CLI_ARGS
	run__start_service_bool                    := flag.Bool("run__start_service",                       true,                  "run the service daemon")
	port_str                                   := flag.String("port",                                   "3050",                "port for the service to use")
	images_store_local_dir_path_str            := flag.String("images_store_local_dir_path",            "./images",            "local dir to store processed images")
	images_thumbnails_store_local_dir_path_str := flag.String("images_thumbnails_store_local_dir_path", "./images/thumbnails", "local dir to store images thumbnails")
	images_s3_bucket_name_str                  := flag.String("images_s3_bucket_name",                  "gf--img",             "AWS S3 bucket name where to store/serve images")

	//MONGODB
	mongodb_host_str    := flag.String("mongodb_host",    "127.0.0.1", "host of mongodb to use")
	mongodb_db_name_str := flag.String("mongodb_db_name", "prod_db",   "DB name to use")
	//-------------------
	//ENV VARS
	aws_access_key_id_str     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	aws_secret_access_key_str := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	aws_token_str             := os.Getenv("GF_AWS_TOKEN")
	//-------------------

	flag.Parse()
	
	return map[string]interface{}{
		"run__start_service_bool":                    *run__start_service_bool,
		"port_str":                                   *port_str,
		"images_store_local_dir_path_str":            *images_store_local_dir_path_str,
		"images_thumbnails_store_local_dir_path_str": *images_thumbnails_store_local_dir_path_str,
		"images_s3_bucket_name_str":                  *images_s3_bucket_name_str,

		//MONGODB
		"mongodb_host_str":    *mongodb_host_str,
		"mongodb_db_name_str": *mongodb_db_name_str,

		//AWS
		"aws_access_key_id_str":     aws_access_key_id_str,
		"aws_secret_access_key_str": aws_secret_access_key_str,
		"aws_token_str":             aws_token_str,
	}
}