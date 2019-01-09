package main

import (
	"flag"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib"
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

	templates_dir_paths_map := map[string]interface{}{
		"flows_str":"./templates",
		"gif_str":  "./templates",
	}

	//START_SERVICE
	if run__start_service_bool {

		gf_images_lib.Run_service(port_str,
			mongodb_host_str,
			mongodb_db_name_str,
			images_store_local_dir_path_str,
			images_thumbnails_store_local_dir_path_str,
			images_main_s3_bucket_name_str,
			templates_dir_paths_map,
			nil, //init_done_ch,
			log_fun)
	}
}
//-------------------------------------------------
func parse__cli_args(p_log_fun func(string,string)) map[string]interface{} {
	p_log_fun("FUN_ENTER","gf_images_service.parse__cli_args()")

	//-------------------
	run__start_service_bool                    := flag.Bool("run__start_service",                      true,                 "run the service daemon")
	port_str                                   := flag.String("port",                                  "3050",               "port for the service to use")
	mongodb_host_str                           := flag.String("mongodb_host",                          "127.0.0.1",          "host of mongodb to use")
	mongodb_db_name_str                        := flag.String("mongodb_db_name",                       "prod_db",            "DB name to use")
	images_store_local_dir_path_str            := flag.String("images_store_local_dir_path",           "./images",           "local dir to store processed images")
	images_thumbnails_store_local_dir_path_str := flag.String("images_thumbnails_store_local_dir_path","./images/thumbnails","local dir to store images thumbnails")
	images_s3_bucket_name_str                  := flag.String("images_s3_bucket_name",                 "gf--img",            "AWS S3 bucket name where to store/serve images")
	//-------------------

	flag.Parse()

	return map[string]interface{}{
		"run__start_service_bool":                   *run__start_service_bool,
		"port_str":                                  *port_str,
		"mongodb_host_str":                          *mongodb_host_str,
		"mongodb_db_name_str":                       *mongodb_db_name_str,
		"images_store_local_dir_path_str":           *images_store_local_dir_path_str,
		"images_thumbnails_store_local_dir_path_str":*images_thumbnails_store_local_dir_path_str,
		"images_s3_bucket_name_str":                 *images_s3_bucket_name_str,
	}
}