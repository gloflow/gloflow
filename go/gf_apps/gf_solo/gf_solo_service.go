/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	"net/http"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_landing_page_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_analytics_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_ml_lib"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func service__run(p_config *GF_config,
	p_runtime_sys *gf_core.Runtime_sys) {

	yellow := color.New(color.BgYellow).Add(color.FgBlack).SprintFunc()
	green  := color.New(color.BgGreen).Add(color.FgBlack).SprintFunc()

	fmt.Printf("%s%s\n", yellow("GF_SOLO"), green("==============="))

	//-------------
	// GF_IMAGES

	// CONFIG
	gf_images__config, gf_err := gf_images_utils.Config__get(p_config.Images__config_file_path_str, p_runtime_sys)
	if gf_err != nil {
		return
	}
	
	gf_images__service_info := &gf_images_lib.GF_service_info{
		Mongodb_host_str:                           p_config.Mongodb_host_str,
		Mongodb_db_name_str:                        p_config.Mongodb_db_name_str,

		Images_store_local_dir_path_str:            gf_images__config.Store_local_dir_path_str,
		Images_thumbnails_store_local_dir_path_str: gf_images__config.Thumbnails_store_local_dir_path_str,
		Images_main_s3_bucket_name_str:             gf_images__config.Main_s3_bucket_name_str,

		AWS_access_key_id_str:                      p_config.AWS_access_key_id_str,
		AWS_secret_access_key_str:                  p_config.AWS_secret_access_key_str,
		AWS_token_str:                              p_config.AWS_token_str,

		Templates_paths_map: p_config.Templates_paths_map,
	}

	gf_images_lib.Init_service(gf_images__service_info,
		gf_images__config,
		p_runtime_sys)

	//-------------
	// GF_LANDING_PAGE

	gf_landing_page_lib.Init_service(p_config.Templates_paths_map, p_runtime_sys)

	//-------------
	// GF_ANALYTICS
	
	gf_analytics__service_info := &gf_analytics_lib.GF_service_info{

		Crawl__config_file_path_str:      p_config.Crawl__config_file_path_str,
		Crawl__cluster_node_type_str:     p_config.Crawl__cluster_node_type_str,
		Crawl__images_local_dir_path_str: p_config.Crawl__images_local_dir_path_str,

		Py_stats_dirs_lst:      p_config.Analytics__py_stats_dirs_lst,
		Run_indexer_bool:       p_config.Analytics__run_indexer_bool,
		Elasticsearch_host_str: p_config.Elasticsearch_host_str,

		AWS_access_key_id_str:     p_config.AWS_access_key_id_str,
		AWS_secret_access_key_str: p_config.AWS_secret_access_key_str,
		AWS_token_str:             p_config.AWS_token_str,

		Templates_paths_map: p_config.Templates_paths_map,
	}
	gf_analytics_lib.Init_service(gf_analytics__service_info,
		p_runtime_sys)

	//-------------
	// GF_PUBLISHER

	// FIX!! - find a soloution where gf_solo gf_publisher functions can invoke
	//         gf_images functions in the same process if in non-distributed mode.
	//         specifying gf_images host
	//         is there because of the default distributed design that assumes
	//         gf_publisher and gf_images run as separate processes.
	gf_images_service_host_port_str := "127.0.0.1"
	gf_images_runtime_info := &gf_publisher_lib.GF_images_extern_runtime_info{
		Jobs_mngr:               nil, // indicates not to send in-process messages to jobs_mngr goroutine, instead use HTTP REST API of gf_images
		Service_host_port_str:   gf_images_service_host_port_str,
		Templates_dir_paths_map: p_config.Templates_paths_map,
	}
	
	gf_publisher_lib.Init_service(gf_images_runtime_info,
		p_runtime_sys)

	//-------------
	// GF_TAGGER
	gf_tagger_lib.Init_service(p_config.Templates_paths_map, p_runtime_sys)

	//-------------
	// GF_ML
	gf_ml_lib.Init_service(p_runtime_sys)

	//-------------

	p_runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	p_runtime_sys.Log_fun("INFO", "STARTING HTTP SERVER - PORT - "+p_config.Port_str)
	p_runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	
	err := http.ListenAndServe(":"+p_config.Port_str, nil)
	if err != nil {
		msg_str := "cant start listening on port - "+p_config.Port_str
		p_runtime_sys.Log_fun("ERROR", msg_str)
		panic(err)
	}
}