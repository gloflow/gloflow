package main

import (
	"os"
	"fmt"
	"flag"
	"net/http"
	"strings"
	"github.com/olivere/elastic"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_stats/gf_stats_apps"
	"github.com/gloflow/gloflow/go/gf_apps/gf_domains_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)
//-------------------------------------------------
func main() {
	log_fun := gf_core.Init_log_fun()

	cli_args_map        := parse__cli_args(log_fun)
	run_str             := cli_args_map["run_str"].(string)
	port_str            := cli_args_map["port_str"].(string)
	mongodb_host_str    := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str := cli_args_map["mongodb_db_name_str"].(string)

	//-----------------
	//MONGODB

	mongodb_db   := gf_core.Mongo__connect(mongodb_host_str, mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")

	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_analytics",
		Log_fun:          log_fun,
		Mongodb_coll:     mongodb_coll,
	}
	//-----------------
	//ELASTICSEARCH
 	
 	//used in case we want to skip using elasticsearch, to avoid that
 	//dependency needing to be present
	run_indexer_bool := cli_args_map["run_indexer_bool"].(bool)

	var esearch_client *elastic.Client
	var gf_err         *gf_core.Gf_error
	if run_indexer_bool {
		esearch_client, gf_err = gf_core.Elastic__get_client(runtime_sys)
		if gf_err != nil {
			panic(gf_err.Error)
		}
	}

	fmt.Println("ELASTIC_SEARCH_CLIENT >>> - "+fmt.Sprint(esearch_client))
	//-----------------

	switch run_str {

		//-----------------------------
		//RUN CRAWLER

		case "run_crawler":
				
			crawler_name_str                  := cli_args_map["crawler_name_str"].(string)
			crawler_cycles_to_run_int         := cli_args_map["crawler_cycles_to_run_int"].(int)
			cluster_node_type_str             := cli_args_map["cluster_node_type_str"].(string)
			crawler_images_local_dir_path_str := cli_args_map["crawler_images_local_dir_path_str"].(string)
			all_crawlers_map                  := gf_crawl_lib.Get_all_crawlers()
			crawler                           := all_crawlers_map[crawler_name_str]

			//-------------
			//S3
			images_s3_bucket_name_str := "gf--discovered--img"
			aws_access_key_id_str     := cli_args_map["aws_access_key_id_str"].(string)
			aws_secret_access_key_str := cli_args_map["aws_secret_access_key_str"].(string)
			aws_token_str             := cli_args_map["aws_token_str"].(string)
			s3_info,gf_err            := gf_core.S3__init(aws_access_key_id_str, aws_secret_access_key_str, aws_token_str, runtime_sys)
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

			//run a certain number of crawl cycles
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
		//DISCOVER DOMAINS IN DB

		case "discover_domains_in_db":
			gf_err := gf_domains_lib.Discover_domains_in_db(runtime_sys)
			if gf_err != nil {
				panic(gf_err.Error)
			}
		//-----------------------------
		//START SERVICE
		case "start_service":
			
			cluster_node_type_str             := cli_args_map["cluster_node_type_str"].(string)
			crawler_images_local_dir_path_str := cli_args_map["crawler_images_local_dir_path_str"].(string)
			py_stats_dirs_lst                 := cli_args_map["py_stats_dirs_lst"].([]string)
			aws_access_key_id_str             := cli_args_map["aws_access_key_id_str"].(string)
			aws_secret_access_key_str         := cli_args_map["aws_secret_access_key_str"].(string)
			aws_token_str                     := cli_args_map["aws_token_str"].(string)

			init_handlers(runtime_sys)
			//------------------------
			//GF_DOMAINS

			gf_domains_lib.Init_domains_aggregation(runtime_sys)
			gf_err := gf_domains_lib.Init_handlers(runtime_sys)
			if gf_err != nil {
				panic(gf_err.Error)
			}
			//------------------------
			//GF_CRAWL

			gf_crawl_lib.Init(crawler_images_local_dir_path_str,
				cluster_node_type_str,
				aws_access_key_id_str,
				aws_secret_access_key_str,
				aws_token_str,
				esearch_client,
				runtime_sys)			
			//------------------------
			//GF_STATS

			stats_url_base_str    := "/a/stats"
			py_stats_dir_path_str := py_stats_dirs_lst[0]

			gf_err = gf_stats_apps.Init(stats_url_base_str, py_stats_dir_path_str, runtime_sys)
			if gf_err != nil {
				panic(gf_err.Error)
			}
			//------------------------
			//STATIC FILES SERVING
			static_files__url_base_str := "/a"
			gf_core.HTTP__init_static_serving(static_files__url_base_str, runtime_sys)
			//------------------------

			log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			log_fun("INFO","STARTING HTTP SERVER - PORT - "+port_str)
			log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			
			err := http.ListenAndServe(":"+port_str,nil)
			if err != nil {
				msg_str := "cant start listening on port - "+port_str
				log_fun("ERROR", msg_str)
				panic(err)
			}
		//-----------------------------
	}
}
//-------------------------------------------------
func parse__cli_args(p_log_fun func(string,string)) map[string]interface{} {
	p_log_fun("FUN_ENTER","gf_analytics_service.parse__cli_args()")

	//-------------------
	run_str                           := flag.String("run",                     "start_service", "start_service|discover_domains_in_db|run_crawler - name of the command to run")
	port_str                          := flag.String("port",                    "3060",          "port for the service to use")
	mongodb_host_str                  := flag.String("mongodb_host",            "127.0.0.1",     "host of mongodb to use")
	mongodb_db_name_str               := flag.String("mongodb_db_name",         "prod_db",       "DB name to use")
	crawler_name_str                  := flag.String("crawler_name",            "gloflow.com",   "name of the crawler to run")
	crawler_cycles_to_run_int         := flag.Int("crawler_cycles_to_run",      1,               "DEBUGGING - when running 'run_crawler' command this indicates how many crawler cylcles o run")
	cluster_node_type_str             := flag.String("cluster_node_type",       "master",        "master|worker - crawler node type")
	crawler_images_local_dir_path_str := flag.String("crawler_images_dir_path", "./image_data",  "local image tmp dir for crawled images before upload to permanent storage")
	run_indexer_bool                  := flag.Bool("run_indexer",               true,            "DEBUG - if the indexer (elasticsearch) should be used")
	py_stats_dirs                     := flag.String("py_stats_dirs",           "./py/stats",    "path to the dir that contains stats .py script files")
	//-------------------
	//ENV VARS
	aws_access_key_id_str     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	aws_secret_access_key_str := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	aws_token_str             := os.Getenv("GF_AWS_TOKEN")
	//-------------------

	flag.Parse()

	return map[string]interface{}{
		"run_str":                           *run_str,
		"port_str":                          *port_str,
		"mongodb_host_str":                  *mongodb_host_str,
		"mongodb_db_name_str":               *mongodb_db_name_str,
		"crawler_name_str":                  *crawler_name_str,
		"crawler_cycles_to_run_int":         *crawler_cycles_to_run_int,
		"cluster_node_type_str":             *cluster_node_type_str,
		"crawler_images_local_dir_path_str": *crawler_images_local_dir_path_str,
		"run_indexer_bool":                  *run_indexer_bool,
		"py_stats_dirs_lst":                 strings.Split(*py_stats_dirs,","),
		"aws_access_key_id_str":             aws_access_key_id_str,
		"aws_secret_access_key_str":         aws_secret_access_key_str,
		"aws_token_str":                     aws_token_str,
	}
}