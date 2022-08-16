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
	"github.com/olivere/elastic"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_analytics_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_domains_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func main() {
	logFun := gf_core.InitLogs()

	cli_args_map        := CLI__parse_args(logFun)
	run_str             := cli_args_map["run_str"].(string)
	port_str            := cli_args_map["port_str"].(string)
	mongodb_host_str       := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str    := cli_args_map["mongodb_db_name_str"].(string)
	elasticsearch_host_str := cli_args_map["elasticsearch_host_str"].(string)
	crawl_config_file_path_str := cli_args_map["crawl_config_file_path_str"].(string)
	
	// used in case we want to skip using elasticsearch, to avoid that
 	// dependency needing to be present
	run_indexer_bool := cli_args_map["run_indexer_bool"].(bool)

	//-----------------
	// MONGODB

	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_analytics",
		LogFun:           logFun,
	}

	mongo_db, gf_err := gf_core.Mongo__connect_new(mongodb_host_str, mongodb_db_name_str, runtime_sys)
	if gf_err != nil {
		panic(-1)
	}
	runtime_sys.mongo_db   = mongo_db
	runtime_sys.Mongo_coll = mongo_db.Collection("data_symphony")

	//-----------------
 	
	switch run_str {

		//-----------------------------
		// RUN CRAWLER - run a certain number of crawler cycles.

		case "run_crawler":
				
			crawler_name_str                  := cli_args_map["crawler_name_str"].(string)
			crawler_cycles_to_run_int         := cli_args_map["crawler_cycles_to_run_int"].(int)
			cluster_node_type_str             := cli_args_map["cluster_node_type_str"].(string)
			crawler_images_local_dir_path_str := cli_args_map["crawler_images_local_dir_path_str"].(string)

			all_crawlers_map, gf_err := gf_crawl_core.Get_all_crawlers(crawl_config_file_path_str, runtime_sys)
			crawler                  := all_crawlers_map[crawler_name_str]


			spew.Dump(all_crawlers_map)


			//-------------
			// AWS_S3
			images_s3_bucket_name_str := "gf--discovered--img"
			aws_access_key_id_str     := cli_args_map["aws_access_key_id_str"].(string)
			aws_secret_access_key_str := cli_args_map["aws_secret_access_key_str"].(string)
			aws_token_str             := cli_args_map["aws_token_str"].(string)
			s3_info, gf_err           := gf_core.S3__init(aws_access_key_id_str, aws_secret_access_key_str, aws_token_str, runtime_sys)
			if gf_err != nil {
				panic(gf_err.Error)
			}

			//-------------
			// ELASTICSEARCH
			var esearch_client *elastic.Client
			if run_indexer_bool {
				esearch_client, gf_err = gf_core.Elastic__get_client(elasticsearch_host_str, runtime_sys)
				if gf_err != nil {
					panic(gf_err.Error)
				}
			}
			fmt.Println("ELASTIC_SEARCH_CLIENT >>> OK")

			//-------------
			crawler_runtime := &gf_crawl_core.Gf_crawler_runtime{
				Events_ctx:            nil,
				Esearch_client:        esearch_client,
				S3_info:               s3_info,
				Cluster_node_type_str: cluster_node_type_str,
			}

			// run a certain number of crawl cycles
			for i := 0; i < crawler_cycles_to_run_int; i++ {

				err := gf_crawl_lib.Run_crawler_cycle(crawler,
					crawler_images_local_dir_path_str,
					images_s3_bucket_name_str,
					crawler_runtime,
					runtime_sys)
				if err != nil {
					panic(err)
				}
			}

		//-----------------------------
		// DISCOVER DOMAINS IN DB

		case "discover_domains_in_db":
			gf_err := gf_domains_lib.Discover_domains_in_db(runtime_sys)
			if gf_err != nil {
				panic(gf_err.Error)
			}

		//-----------------------------
		// START SERVICE
		case "start_service":
			crawl_config_file_path_str      := cli_args_map["crawl_config_file_path_str"]
			crawl_cluster_node_type_str     := cli_args_map["cluster_node_type_str"].(string)
			crawl_images_local_dir_path_str := cli_args_map["crawler_images_local_dir_path_str"].(string)
			
			media_domain_str  := cli_args_map["media_domain_str"].(string)
			py_stats_dirs_lst := cli_args_map["py_stats_dirs_lst"].([]string)

			// AWS
			aws_access_key_id_str     := cli_args_map["aws_access_key_id_str"].(string)
			aws_secret_access_key_str := cli_args_map["aws_secret_access_key_str"].(string)
			aws_token_str             := cli_args_map["aws_token_str"].(string)

			config := &GF_analytics_config{
				Port_str: port_str,

				Crawl__config_file_path_str:      crawl_config_file_path_str,
				Crawl__cluster_node_type_str:     crawl_cluster_node_type_str,
				Crawl__images_local_dir_path_str: crawl_images_local_dir_path_str,

				Media_domain_str:       media_domain_str,
				Py_stats_dirs_lst:      py_stats_dirs_lst,
				Run_indexer_bool:       run_indexer_bool,
				Elasticsearch_host_str: elasticsearch_host_str,

				AWS_access_key_id_str:     aws_access_key_id_str,
				AWS_secret_access_key_str: aws_secret_access_key_str,
				AWS_token_str:             aws_token_str,
			}

			gf_analytics_lib.Run_service(config,
				runtime_sys)


		//-----------------------------
	}
}