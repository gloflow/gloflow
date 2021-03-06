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
/*BosaC.Jan30.2020. <3 volim te zauvek*/

package gf_crawl_lib

import (
	"time"
	"math/rand"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/olivere/elastic"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_utils"
)

//--------------------------------------------------
type GF_crawler_config struct {
	Crawled_images_s3_bucket_name_str string
	Images_s3_bucket_name_str         string
	Images_local_dir_path_str         string
	Cluster_node_type_str             string
	Crawl_config_file_path_str        string
}

type Gf_crawler struct {
	Name_str      string
	Start_url_str string
	// Domains_lst   []string //some sites have multiple domains
}

type Gf_crawler_cycle_run struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               string        `bson:"id_str"`
	T_str                string        `bson:"t"` //"crawler_cycle_run"
	Creation_unix_time_f float64       `bson:"creation_unix_time_f"`
	Crawler_name_str     string        `bson:"crawler_name_str"`
	Target_domain_str    string        `bson:"targit_domain_str"`
	Target_url_str       string        `bson:"target_url_str"`
	Start_time_f         float64       `bson:"start_time_f"`
	End_time_f           float64       `bson:"end_time_f"`
}

//--------------------------------------------------
func Init(p_config *GF_crawler_config,
	// p_images_local_dir_path_str  string,
	// p_cluster_node_type_str      string,
	// p_crawl_config_file_path_str string,
	p_media_domain_str           string,
	p_templates_paths_map        map[string]string,
	p_aws_access_key_id_str      string,
	p_aws_secret_access_key_str  string,
	p_aws_token_str              string,
	p_esearch_client             *elastic.Client,
	p_runtime_sys                *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl.Init()")

	//--------------
	events_ctx := gf_core.Events__init("/a/crawl/events", p_runtime_sys)

	// crawled_images_s3_bucket_name_str := "gf--discovered--img"
	// gf_images_s3_bucket_name_str      := "gf--img"

	gf_s3_info, gf_err := gf_core.S3__init(p_aws_access_key_id_str,
		p_aws_secret_access_key_str,
		p_aws_token_str,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	runtime := &gf_crawl_core.Gf_crawler_runtime{
		Events_ctx:            events_ctx,
		Esearch_client:        p_esearch_client,
		S3_info:               gf_s3_info,
		Cluster_node_type_str: p_config.Cluster_node_type_str,
	}

	//--------------
	// IMPORTANT!! - make sure mongo has indexes build for relevant queries
	db_index__init(runtime, p_runtime_sys)
	
	/*crawlers_map, gf_err := gf_crawl_core.Get_all_crawlers(p_config.Crawl_config_file_path_str,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}*/

	/*start_crawlers_cycles(crawlers_map,
		p_config.Images_local_dir_path_str,
		crawled_images_s3_bucket_name_str,
		runtime,
		p_runtime_sys)*/

	//--------------
	// HTTP_HANDLERS
	gf_err = init_handlers(p_media_domain_str,
		p_config.Crawled_images_s3_bucket_name_str,
		p_config.Images_s3_bucket_name_str,
		p_templates_paths_map,
		runtime,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
		
	// HTTP_HANDLERS__CLUSTER
	gf_err = cluster__init_handlers(p_config.Crawl_config_file_path_str,
		runtime,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	//--------------

	return nil
}

//--------------------------------------------------
func start_crawlers_cycles(p_crawlers_map map[string]gf_crawl_core.Gf_crawler_def,
	p_images_local_dir_path_str string,
	p_media_domain_str          string,
	p_images_s3_bucket_name_str string,
	p_runtime                   *gf_crawl_core.Gf_crawler_runtime,
	p_runtime_sys               *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl.start_crawlers_cycles()")

	events_id_str := "crawler_events"
	
	gf_core.Events__register_producer(events_id_str, p_runtime.Events_ctx, p_runtime_sys)

	for _, crawler := range p_crawlers_map {

		// IMPORTANT!! - each crawler runs in its own goroutine, and continuously
		//               crawls the target domains
		go func(p_crawler gf_crawl_core.Gf_crawler_def) {
			start_crawler(p_crawler,
				p_images_local_dir_path_str,

				p_media_domain_str,
				p_images_s3_bucket_name_str,
				p_runtime,
				p_runtime_sys)

		}(crawler)
	}
}

//--------------------------------------------------
func start_crawler(p_crawler gf_crawl_core.Gf_crawler_def,
	p_images_local_dir_path_str string,
	p_media_domain_str          string,
	p_images_s3_bucket_name_str string,
	p_runtime                   *gf_crawl_core.Gf_crawler_runtime,
	p_runtime_sys               *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl.start_crawler()")

	yellow := color.New(color.FgYellow).SprintFunc()
	black  := color.New(color.FgBlack).Add(color.BgGreen).SprintFunc()

	p_runtime_sys.Log_fun("INFO", black("------------------------------------"))
	p_runtime_sys.Log_fun("INFO", black(">>>    STARTING CRAWLER >>> ")+yellow(p_crawler.Name_str))
	p_runtime_sys.Log_fun("INFO", black("------------------------------------"))

	//-----------------
	// LINK_ALLOCATOR
	gf_crawl_core.Link_alloc__init(p_crawler.Name_str, p_runtime_sys)
	
	//-----------------

	// randomize r.Intn() usage, otherwise its determanistic 
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	i := 0

	for ;; {
		gf_crawl_utils.Crawler_sleep(p_crawler.Name_str, i, r, p_runtime_sys)
		//-----------------
		// RUN CRAWLER
		gf_err := Run_crawler_cycle(p_crawler,
			p_images_local_dir_path_str,

			p_media_domain_str,
			p_images_s3_bucket_name_str,
			p_runtime,
			p_runtime_sys)
		if gf_err != nil {
			// ADD!! - do something useful with this error, although its persisted to DB since its a gf_err
		}

		//-----------------
		i=i+1
	}
}
