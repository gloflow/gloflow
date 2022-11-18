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
	"os"
	"flag"
	"strings"
)

//-------------------------------------------------
func CLI__parse_args(pLogFun func(string, string)) map[string]interface{} {
	pLogFun("FUN_ENTER", "gf_analytics_cli.CLI__parse_args()")

	default_command_str := "start_service"
	//-------------------
	run_str                           := flag.String("run",                     default_command_str,          "start_service|discover_domains_in_db|run_crawler - name of the command to run")
	port_str                          := flag.String("port",                    "3060",                       "port for the service to use")
	
	crawler_name_str                  := flag.String("crawler_name",            "gloflow.com",                "name of the crawler to run")
	crawler_cycles_to_run_int         := flag.Int("crawler_cycles_to_run",      1,                            "DEBUGGING - when running 'run_crawler' command this indicates how many crawler cylcles o run")
	cluster_node_type_str             := flag.String("cluster_node_type",       "master",                     "master|worker - crawler node type")
	crawl_config_file_path_str        := flag.String("crawl_config_file_path",  "./config/crawl_config.yaml", "local image tmp dir for crawled images before upload to permanent storage")
	crawler_images_local_dir_path_str := flag.String("crawler_images_dir_path", "./data/images",              "local image tmp dir for crawled images before upload to permanent storage")
	run_indexer_bool                  := flag.Bool("run_indexer",               true,                         "DEBUG - if the indexer (elasticsearch) should be used")
	py_stats_dirs                     := flag.String("py_stats_dirs",           "./py/stats",                 "path to the dir that contains stats .py script files")

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

	// ELASTICSEARCH
	elasticsearch_host_str := flag.String("elasticsearch_host", "127.0.0.1:9200", "host of elasticsearch to use")

	// ELASTICSEARCH_ENV
	elasticsearch_host_env_str := os.Getenv("GF_ELASTICSEARCH_HOST")

	if elasticsearch_host_env_str != "" {
		*elasticsearch_host_str = elasticsearch_host_env_str
	}

	//-------------------
	// AWS_ENV_VARS
	aws_access_key_id_str     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	aws_secret_access_key_str := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	aws_token_str             := os.Getenv("GF_AWS_TOKEN")

	if aws_access_key_id_str == "" || aws_secret_access_key_str == "" {
		pLogFun("ERROR", "ENV vars not set - GF_AWS_ACCESS_KEY_ID, GF_AWS_SECRET_ACCESS_KEY")
		panic(1)
	}
	//-------------------

	flag.Parse()

	return map[string]interface{}{
		"run_str":                           *run_str,
		"port_str":                          *port_str,
		"mongodb_host_str":                  *mongodb_host_str,
		"mongodb_db_name_str":               *mongodb_db_name_str,
		"elasticsearch_host_str":            *elasticsearch_host_str,
		"crawler_name_str":                  *crawler_name_str,
		"crawler_cycles_to_run_int":         *crawler_cycles_to_run_int,
		"cluster_node_type_str":             *cluster_node_type_str,
		"crawl_config_file_path_str":        *crawl_config_file_path_str,
		"crawler_images_local_dir_path_str": *crawler_images_local_dir_path_str,
		"run_indexer_bool":                  *run_indexer_bool,
		"py_stats_dirs_lst":                 strings.Split(*py_stats_dirs,","),
		"aws_access_key_id_str":             aws_access_key_id_str,
		"aws_secret_access_key_str":         aws_secret_access_key_str,
		"aws_token_str":                     aws_token_str,
	}
}