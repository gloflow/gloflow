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

package gf_analytics_lib

import (
	// "os"
	"fmt"
	"net/http"
	"github.com/olivere/elastic"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_stats/gf_stats_apps"
	"github.com/gloflow/gloflow/go/gf_apps/gf_domains_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_service_info struct {
	Port_str string

	Crawl__config_file_path_str      string
	Crawl__cluster_node_type_str     string
	Crawl__images_local_dir_path_str string

	Media_domain_str       string 
	Py_stats_dirs_lst      []string
	Run_indexer_bool       bool
	Elasticsearch_host_str string

	AWS_access_key_id_str     string
	AWS_secret_access_key_str string
	AWS_token_str             string
	
	Templates_paths_map map[string]string
}

//-------------------------------------------------
func Init_service(p_service_info *GF_service_info,
	p_runtime_sys *gf_core.Runtime_sys) {

	//-----------------
	// ELASTICSEARCH
	var esearch_client *elastic.Client
	var gf_err         *gf_core.Gf_error
	if p_service_info.Run_indexer_bool {
		esearch_client, gf_err = gf_core.Elastic__get_client(p_service_info.Elasticsearch_host_str, p_runtime_sys)
		if gf_err != nil {
			panic(gf_err.Error)
		}
	}
	fmt.Println("ELASTIC_SEARCH_CLIENT >>> OK")

	//-----------------
	// TEMPLATES_DIR
	/*templates_dir_path_str := "./templates"
	if _, err := os.Stat(templates_dir_path_str); os.IsNotExist(err) {
		p_runtime_sys.Log_fun("ERROR", fmt.Sprintf("templates dir doesnt exist - %s", templates_dir_path_str))
		panic(1)
	}*/

	init_handlers(p_service_info.Templates_paths_map, p_runtime_sys)

	//------------------------
	// GF_DOMAINS
	gf_domains_lib.DB_index__init(p_runtime_sys)
	gf_domains_lib.Init_domains_aggregation(p_runtime_sys)
	gf_err = gf_domains_lib.Init_handlers(p_service_info.Templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//------------------------
	// GF_CRAWL

	crawl_config := &gf_crawl_lib.GF_crawler_config{
		Crawled_images_s3_bucket_name_str: "gf--discovered--img",
		Images_s3_bucket_name_str:         "gf--img",
		Images_local_dir_path_str:         p_service_info.Crawl__images_local_dir_path_str,
		Cluster_node_type_str:             p_service_info.Crawl__cluster_node_type_str,
		Crawl_config_file_path_str:        p_service_info.Crawl__config_file_path_str,
	}
	gf_crawl_lib.Init(crawl_config, // p_service_info.Crawl__images_local_dir_path_str,
		// p_service_info.Crawl__cluster_node_type_str,
		// p_service_info.Crawl__config_file_path_str,
		p_service_info.Media_domain_str,
		p_service_info.Templates_paths_map,
		p_service_info.AWS_access_key_id_str,
		p_service_info.AWS_secret_access_key_str,
		p_service_info.AWS_token_str,
		esearch_client,
		p_runtime_sys)

	//------------------------
	// GF_STATS

	stats_url_base_str    := "/a/stats"
	py_stats_dir_path_str := p_service_info.Py_stats_dirs_lst[0]

	gf_err = gf_stats_apps.Init(stats_url_base_str, py_stats_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//------------------------
	// STATIC FILES SERVING
	static_files__url_base_str := "/a"
	gf_core.HTTP__init_static_serving(static_files__url_base_str, p_runtime_sys)

	//------------------------

}

//-------------------------------------------------
func Run_service(p_service_info *GF_service_info,
	p_runtime_sys *gf_core.Runtime_sys) {
	
	//------------------------
	// INIT
	Init_service(p_service_info,
		p_runtime_sys)
	
	//------------------------

	p_runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	p_runtime_sys.Log_fun("INFO", "STARTING HTTP SERVER - PORT - "+p_service_info.Port_str)
	p_runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	
	err := http.ListenAndServe(":"+p_service_info.Port_str, nil)
	if err != nil {
		msg_str := "cant start listening on port - "+p_service_info.Port_str
		p_runtime_sys.Log_fun("ERROR", msg_str)
		panic(err)
	}
}

//-------------------------------------------------
/*func main() {
	log_fun := gf_core.Init_log_fun()

	cli_args_map        := CLI__parse_args(log_fun)
	run_str             := cli_args_map["run_str"].(string)
	port_str            := cli_args_map["port_str"].(string)
	mongodb_host_str       := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str    := cli_args_map["mongodb_db_name_str"].(string)
	elasticsearch_host_str := cli_args_map["elasticsearch_host_str"].(string)
	crawl_config_file_path_str := cli_args_map["crawl_config_file_path_str"].(string)
	//-----------------
	// MONGODB

	mongodb_db   := gf_core.Mongo__connect(mongodb_host_str, mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")

	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_analytics",
		Log_fun:          log_fun,
		Mongodb_db:       mongodb_db,
		Mongodb_coll:     mongodb_coll,
	}
	
	//-----------------
	// ELASTICSEARCH
 	
 	// used in case we want to skip using elasticsearch, to avoid that
 	// dependency needing to be present
	run_indexer_bool := cli_args_map["run_indexer_bool"].(bool)

	var esearch_client *elastic.Client
	var gf_err         *gf_core.Gf_error
	if run_indexer_bool {
		esearch_client, gf_err = gf_core.Elastic__get_client(elasticsearch_host_str, runtime_sys)
		if gf_err != nil {
			panic(gf_err.Error)
		}
	}

	fmt.Println("ELASTIC_SEARCH_CLIENT >>> - "+fmt.Sprint(esearch_client))

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
			
			cluster_node_type_str             := cli_args_map["cluster_node_type_str"].(string)
			crawler_images_local_dir_path_str := cli_args_map["crawler_images_local_dir_path_str"].(string)
			py_stats_dirs_lst                 := cli_args_map["py_stats_dirs_lst"].([]string)

			// AWS
			aws_access_key_id_str     := cli_args_map["aws_access_key_id_str"].(string)
			aws_secret_access_key_str := cli_args_map["aws_secret_access_key_str"].(string)
			aws_token_str             := cli_args_map["aws_token_str"].(string)

			// TEMPLATES_DIR
			templates_dir_path_str := "./templates"
			if _, err := os.Stat(templates_dir_path_str); os.IsNotExist(err) {
				log_fun("ERROR", fmt.Sprintf("templates dir doesnt exist - %s", templates_dir_path_str))
				panic(1)
			}

			init_handlers(templates_dir_path_str, runtime_sys)
			//------------------------
			// GF_DOMAINS
			gf_domains_lib.DB_index__init(runtime_sys)
			gf_domains_lib.Init_domains_aggregation(runtime_sys)
			gf_err := gf_domains_lib.Init_handlers(templates_dir_path_str, runtime_sys)
			if gf_err != nil {
				panic(gf_err.Error)
			}

			//------------------------
			// GF_CRAWL

			gf_crawl_lib.Init(crawler_images_local_dir_path_str,
				cluster_node_type_str,
				crawl_config_file_path_str,
				templates_dir_path_str,
				aws_access_key_id_str,
				aws_secret_access_key_str,
				aws_token_str,
				esearch_client,
				runtime_sys)

			//------------------------
			// GF_STATS

			stats_url_base_str    := "/a/stats"
			py_stats_dir_path_str := py_stats_dirs_lst[0]

			gf_err = gf_stats_apps.Init(stats_url_base_str, py_stats_dir_path_str, runtime_sys)
			if gf_err != nil {
				panic(gf_err.Error)
			}

			//------------------------
			// STATIC FILES SERVING
			static_files__url_base_str := "/a"
			gf_core.HTTP__init_static_serving(static_files__url_base_str, runtime_sys)

			//------------------------

			log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			log_fun("INFO", "STARTING HTTP SERVER - PORT - "+port_str)
			log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			
			err := http.ListenAndServe(":"+port_str, nil)
			if err != nil {
				msg_str := "cant start listening on port - "+port_str
				log_fun("ERROR", msg_str)
				panic(err)
			}

		//-----------------------------
	}
}*/