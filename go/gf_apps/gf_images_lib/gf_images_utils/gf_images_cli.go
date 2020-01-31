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

package gf_images_utils

import (
	"os"
	"flag"
)

//-------------------------------------------------
func CLI__parse_args(p_log_fun func(string,string)) map[string]interface{} {
	p_log_fun("FUN_ENTER", "gf_images_cli.CLI__parse_args()")

	//-------------------
	// CLI_ARGS
	run__start_service_bool                    := flag.Bool("run__start_service",                       true,                  "run the service daemon")
	port_str                                   := flag.String("port",                                   "3050",                "port for the service to use")
	images_store_local_dir_path_str            := flag.String("images_store_local_dir_path",            "./images",            "local dir to store processed images")
	images_thumbnails_store_local_dir_path_str := flag.String("images_thumbnails_store_local_dir_path", "./images/thumbnails", "local dir to store images thumbnails")
	images_s3_bucket_name_str                  := flag.String("images_s3_bucket_name",                  "gf--img",             "AWS S3 bucket name where to store/serve images")

	//-------------------
	// MONGODB
	mongodb_host_str    := flag.String("mongodb_host",    "127.0.0.1", "host of mongodb to use")
	mongodb_db_name_str := flag.String("mongodb_db_name", "prod_db",   "DB name to use")

	// MONGODB_ENV
	mongodb_host_env_str    := os.Getenv("GF_MONGODB_HOST")
	mongodb_db_name_env_str := os.Getenv("GF_MONGODB_DB_NAME")

	if mongodb_db_name_env_str != "" {
		*mongodb_db_name_str = mongodb_db_name_env_str
	}

	if mongodb_host_env_str != "" {
		*mongodb_host_str = mongodb_host_env_str
	}
	
	//-------------------
	// AWS_CREDS
	aws_access_key_id_str     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	aws_secret_access_key_str := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	aws_token_str             := os.Getenv("GF_AWS_TOKEN")

	if len(aws_access_key_id_str) != 20 {
		panic("GF_AWS_ACCESS_KEY_ID ENV var not set/of correct length")
	}

	if len(aws_secret_access_key_str) != 40 {
		panic("GF_AWS_SECRET_ACCESS_KEY ENV var not set/of correct length")
	}

	//-------------------
	// AWS_S3

	images_s3_bucket_name_env_str := os.Getenv("GF_IMAGES_S3_BUCKET_NAME")
	if images_s3_bucket_name_env_str != "" {
		*images_s3_bucket_name_str = images_s3_bucket_name_env_str
	}
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