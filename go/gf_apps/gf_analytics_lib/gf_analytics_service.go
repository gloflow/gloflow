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
	// "fmt"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_domains_lib"

	// "github.com/olivere/elastic"
	// "github.com/gloflow/gloflow/go/gf_stats/gf_stats_apps"
	// "github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFserviceInfo struct {
	Port_str string

	Crawl__config_file_path_str      string
	Crawl__cluster_node_type_str     string
	Crawl__images_local_dir_path_str string

	Media_domain_str       string 
	Py_stats_dirs_lst      []string
	Run_indexer_bool       bool
	// Elasticsearch_host_str string
	Templates_paths_map map[string]string

	// IMAGES_STORAGE
	ImagesUseNewStorageEngineBool bool
}

//-------------------------------------------------

func InitService(pServiceInfo *GFserviceInfo,
	pHTTPmux  *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) {

	//-----------------
	// ELASTICSEARCH
	/*
	var esearch_client *elastic.Client
	var gfErr          *gf_core.GFerror
	if pServiceInfo.Run_indexer_bool {
		esearch_client, gfErr = gf_core.Elastic__get_client(pServiceInfo.Elasticsearch_host_str, pRuntimeSys)
		if gfErr != nil {
			panic(gfErr.Error)
		}
	}
	fmt.Println("ELASTIC_SEARCH_CLIENT >>> OK")
	*/

	//-----------------
	initHandlers(pServiceInfo.Templates_paths_map, pHTTPmux, pRuntimeSys)

	//------------------------
	// GF_DOMAINS
	gf_domains_lib.DBmongoIndexInit(pRuntimeSys)
	gf_domains_lib.InitDomainsAggregation(pRuntimeSys)
	gfErr := gf_domains_lib.InitHandlers(pServiceInfo.Templates_paths_map,
		pHTTPmux,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//------------------------
	// GF_CRAWL

	/*
	crawl_config := &gf_crawl_lib.GFcrawlerConfig{
		Crawled_images_s3_bucket_name_str: "gf--discovered--img",
		Images_s3_bucket_name_str:         "gf--img",
		Images_local_dir_path_str:         pServiceInfo.Crawl__images_local_dir_path_str,
		Cluster_node_type_str:             pServiceInfo.Crawl__cluster_node_type_str,
		Crawl_config_file_path_str:        pServiceInfo.Crawl__config_file_path_str,
		ImagesUseNewStorageEngineBool:     pServiceInfo.ImagesUseNewStorageEngineBool,
	}
	gf_crawl_lib.Init(crawl_config, // pServiceInfo.Crawl__images_local_dir_path_str,
		// pServiceInfo.Crawl__cluster_node_type_str,
		// pServiceInfo.Crawl__config_file_path_str,
		pServiceInfo.Media_domain_str,
		pServiceInfo.Templates_paths_map,
		
		esearch_client,
		pHTTPmux,
		pRuntimeSys)
	*/

	//------------------------
	/*
	// GF_STATS

	stats_url_base_str    := "/a/stats"
	py_stats_dir_path_str := pServiceInfo.Py_stats_dirs_lst[0]

	gfErr = gf_stats_apps.Init(stats_url_base_str, py_stats_dir_path_str, pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}
	*/

	//------------------------
	// STATIC FILES SERVING
	staticFilesURLbaseStr := "/a"
	localDirPathStr       := "./static"
	gf_core.HTTPinitStaticServingWithMux(staticFilesURLbaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)

	//------------------------

}

//-------------------------------------------------

func RunService(pServiceInfo *GFserviceInfo,
	pRuntimeSys *gf_core.RuntimeSys) {
	
	//------------------------
	// INIT
	http_mux := http.NewServeMux()

	InitService(pServiceInfo,
		http_mux,
		pRuntimeSys)
	
	//------------------------

	pRuntimeSys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	pRuntimeSys.LogFun("INFO", "STARTING HTTP SERVER - PORT - "+pServiceInfo.Port_str)
	pRuntimeSys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	
	err := http.ListenAndServe(":"+pServiceInfo.Port_str, nil)
	if err != nil {
		msg_str := "cant start listening on port - "+pServiceInfo.Port_str
		pRuntimeSys.LogFun("ERROR", msg_str)
		panic(err)
	}
}