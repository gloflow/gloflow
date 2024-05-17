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
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/olivere/elastic"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_utils"
)

//--------------------------------------------------

type GFcrawlerConfig struct {
	Crawled_images_s3_bucket_name_str string
	Images_s3_bucket_name_str         string
	Images_local_dir_path_str         string
	Cluster_node_type_str             string
	Crawl_config_file_path_str        string
	ImagesUseNewStorageEngineBool     bool
}

type GFcrawler struct {
	Name_str      string
	StartURLstr string
	// Domains_lst   []string //some sites have multiple domains
}

type GFcrawlerCycleRun struct {
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

func Init(pConfig *GFcrawlerConfig,
	pMediaDomainStr             string,
	pTemplatesPathsMap          map[string]string,

	/*
	// DEPRECATE!!
	p_aws_access_key_id_str     string,
	p_aws_secret_access_key_str string,
	p_aws_token_str             string,
	*/
	
	pEsearchClient              *elastic.Client,
	pHTTPmux                    *http.ServeMux,
	pRuntimeSys                 *gf_core.RuntimeSys) *gf_core.GFerror {

	//--------------
	eventsCtx := gf_events.Init("/a/crawl/events", pRuntimeSys)

	// crawled_images_s3_bucket_name_str := "gf--discovered--img"
	// gf_images_s3_bucket_name_str      := "gf--img"

	gf_s3_info, gfErr := gf_aws.S3init(/*p_aws_access_key_id_str,
		p_aws_secret_access_key_str,
		p_aws_token_str,*/
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	runtime := &gf_crawl_core.GFcrawlerRuntime{
		EventsCtx:                     eventsCtx,
		// EsearchClient:                 pEsearchClient,
		S3info:                        gf_s3_info,
		ImagesUseNewStorageEngineBool: pConfig.ImagesUseNewStorageEngineBool,
	}

	//--------------
	// IMPORTANT!! - make sure mongo has indexes build for relevant queries
	dbMongoIndexInit(pRuntimeSys)
	
	/*
	crawlersMap, gfErr := gf_crawl_core.GetAllCrawlers(pConfig.Crawl_config_file_path_str,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	*/

	/*
	startCrawlersCycles(crawlersMap,
		pConfig.Images_local_dir_path_str,
		crawled_images_s3_bucket_name_str,
		runtime,
		pRuntimeSys)
	*/

	//--------------
	// HTTP_HANDLERS
	gfErr = initHandlers(pMediaDomainStr,
		pConfig.Crawled_images_s3_bucket_name_str,
		pConfig.Images_s3_bucket_name_str,
		pTemplatesPathsMap,
		pHTTPmux,
		runtime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	//--------------

	return nil
}

//--------------------------------------------------

func startCrawlersCycles(pCrawlersMap map[string]gf_crawl_core.GFcrawlerDef,
	pImagesLocalDirPathStr string,
	pMediaDomainStr        string,
	pImagesS3bucketNameStr string,
	pRuntime               *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) {

	eventsIDstr := "crawler_events"
	
	gf_events.RegisterProducer(eventsIDstr, pRuntime.EventsCtx, pRuntimeSys)

	for _, crawler := range pCrawlersMap {

		// IMPORTANT!! - each crawler runs in its own goroutine, and continuously
		//               crawls the target domains
		go func(pCrawler gf_crawl_core.GFcrawlerDef) {
			startCrawler(pCrawler,
				pImagesLocalDirPathStr,

				pMediaDomainStr,
				pImagesS3bucketNameStr,
				pRuntime,
				pRuntimeSys)

		}(crawler)
	}
}

//--------------------------------------------------

func startCrawler(pCrawler gf_crawl_core.GFcrawlerDef,
	pImagesLocalDirPathStr string,
	pMediaDomainStr        string,
	pImagesS3bucketNameStr string,
	pRuntime               *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) {

	yellow := color.New(color.FgYellow).SprintFunc()
	black  := color.New(color.FgBlack).Add(color.BgGreen).SprintFunc()

	pRuntimeSys.LogFun("INFO", black("------------------------------------"))
	pRuntimeSys.LogFun("INFO", black(">>>    STARTING CRAWLER >>> ")+yellow(pCrawler.NameStr))
	pRuntimeSys.LogFun("INFO", black("------------------------------------"))

	//-----------------
	// USER_ID
	// IMPORTANT!! - use system user for the regular crawling runs.
	//               use some specific users ID if a particular user runs the
	//               crawler in some way.
	userID := gf_core.GF_ID("gf")

	//-----------------
	// LINK_ALLOCATOR
	gf_crawl_core.LinkAllocInit(pCrawler.NameStr, pRuntimeSys)
	
	//-----------------

	// randomize r.Intn() usage, otherwise its determanistic 
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	i := 0

	for ;; {
		gf_crawl_utils.CrawlerSleep(pCrawler.NameStr, i, r, pRuntimeSys)
		//-----------------
		// RUN CRAWLER
		gfErr := RunCrawlerCycle(pCrawler,
			pImagesLocalDirPathStr,

			pMediaDomainStr,
			pImagesS3bucketNameStr,
			userID,
			pRuntime,
			pRuntimeSys)
		if gfErr != nil {
			// ADD!! - do something useful with this error, although its persisted to DB since its a gfErr
		}

		//-----------------
		i=i+1
	}
}
