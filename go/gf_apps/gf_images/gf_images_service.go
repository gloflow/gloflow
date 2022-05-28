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
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------
func main() {
	log_fun := gf_core.Init_log_fun()

	// REMOVE!! - "_str"/"_bool" postfixes from the CLI args, the names that are used on the CLI
	cli_args_map            := gf_images_core.CLI__parse_args(log_fun)
	run__start_service_bool := cli_args_map["run__start_service_bool"].(bool)
	port_str                := cli_args_map["port_str"].(string)
	mongodb_host_str        := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str     := cli_args_map["mongodb_db_name_str"].(string)
	images_store_local_dir_path_str            := cli_args_map["images_store_local_dir_path_str"].(string)
	images_thumbnails_store_local_dir_path_str := cli_args_map["images_thumbnails_store_local_dir_path_str"].(string)
	media_domain_str                           := cli_args_map["media_domain_str"].(string)
	images_main_s3_bucket_name_str             := cli_args_map["images_s3_bucket_name_str"].(string)
	aws_access_key_id_str                      := cli_args_map["aws_access_key_id_str"].(string)
	aws_secret_access_key_str                  := cli_args_map["aws_secret_access_key_str"].(string)
	aws_token_str                              := cli_args_map["aws_token_str"].(string)

	templates_dir_paths_map := map[string]string{
		"flows_str": "./templates",
		"gif_str":   "./templates",
	}

	config_file_path_str := "./config/gf_images_config.yaml"

	// fmt.Println("AWS------------------------------------------")
	// fmt.Println(aws_access_key_id_str)
	// fmt.Println(aws_secret_access_key_str)
	// fmt.Println(aws_token_str)

	// START_SERVICE
	if run__start_service_bool {

		serviceInfo := &gf_images_core.GFserviceInfo{
			Port_str:                                   port_str,
			Mongodb_host_str:                           mongodb_host_str,
			Mongodb_db_name_str:                        mongodb_db_name_str,
			Images_store_local_dir_path_str:            images_store_local_dir_path_str,
			Images_thumbnails_store_local_dir_path_str: images_thumbnails_store_local_dir_path_str,
			Media_domain_str:                           media_domain_str,
			Images_main_s3_bucket_name_str:             images_main_s3_bucket_name_str,
			AWS_access_key_id_str:                      aws_access_key_id_str,
			AWS_secret_access_key_str:                  aws_secret_access_key_str,
			AWS_token_str:                              aws_token_str,
			Templates_paths_map:                        templates_dir_paths_map,
			Config_file_path_str:                       config_file_path_str,
		}

		// HTTP_MUX
		serviceHTTPmux := http.NewServeMux()

		gf_images_lib.Run_service(serviceHTTPmux,
			serviceInfo,
			nil, // init_done_ch,
			log_fun)
	}
}