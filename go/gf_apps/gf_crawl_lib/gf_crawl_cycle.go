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

package gf_crawl_lib

import (
	"time"
	"fmt"
	"context"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)

//--------------------------------------------------

func RunCrawlerCycle(pCrawler gf_crawl_core.GFcrawlerDef,
	pImagesLocalDirPathStr string,
	pMediaDomainStr        string,
	pS3bucketNameStr       string,
	pUserID                gf_core.GF_ID,
	pRuntime               *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) *gf_core.GFerror {

	yellow := color.New(color.FgYellow).SprintFunc()
	black  := color.New(color.FgBlack).Add(color.BgWhite).SprintFunc()

	start_time_f := float64(time.Now().UnixNano())/1000000000.0

	//-------------------------
	// STAGE - GET TARGET URL


	// IMPORTANT!! - get unresolved links to pages on the domain to which the crawler belongs.
	//               so if the a page contains links to domains external to the domain to which the 
	//               crawler belongs, it wont get fetched/parsed here
	unresolved_link, gfErr := gf_crawl_core.DBmongoLinkGetUnresolved(pCrawler.NameStr, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// // IMPORTANT!! - no unresolved links were found, this is a valid possible state. 
	// //               do nothing, because further bellow unresolved_link is tested for == nil and then
	// //               crawlenrs start_url is used - hence the empty IF block.
	// if gfErr != nil && fmt.Sprint(gfErr.Type_str) == "mongodb_not_found_error" {

	// } else if gfErr != nil {
	// 	return gfErr
	// }
	
	var url_str string
	if unresolved_link != nil {

		url_str = unresolved_link.A_href_str
		fmt.Println(">> "+black(">>>>>>>>>>>>>>>>> UNRESOLVED_LINK FOUND - ")+yellow(url_str))

		// HACK!! - this is for strange links that only contain the protocol specifier
		if url_str == "http://" {
			return nil
		}

		//----------------------
		// MARK_LINK IMPORT_IN_PROGRESS

		// IMPORTANT!! - because the global DB of unresolved_link's is parsed by many crawler instances running
		//               across the whole gf_crawl cluster. by marking links as "in_progress" they are flagged
		//               as being processed (which might take some unknown amount of time) and other crawler
		//               instances running on other nodes should not load this link of importing as well, 
		//               to avoid duplicate work/data
		start_time_f := float64(time.Now().UnixNano())/1000000000.0
		gfErr       := gf_crawl_core.DBmongoLinkMarkImportInProgress(true, start_time_f, unresolved_link, pRuntime, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		//----------------------
	} else {

		// IMPORTANT!! - if all links were resolved, then start at the initial Crawler StartURLstr
		//               and begin a new sweep of crawling the domain
		url_str = pCrawler.StartURLstr
		fmt.Println("INFO",black(">>>>>>>>>>>>>>>>> UNRESOLVED_LINK NOT FOUND - using start_url - ")+yellow(url_str))
	}
	
	//-------------------------
	// CYCLE_RUN
	cycle_run__creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	cycle_run__id_str               := "crawler_cycle_run:"+fmt.Sprint(cycle_run__creation_unix_time_f)

	//-------------------------
	
	//-------------------
	// STAGE - FETCH THE LINK
	url_fetch, domain_str, gfErr := gf_crawl_core.FetchURL(url_str,
		unresolved_link,
		cycle_run__id_str,
		pCrawler.NameStr,
		pRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-------------------
	// STAGE - PARSE THE FETCHED PAGE
	gfErr = gf_crawl_core.FetchParseResult(url_fetch,
		cycle_run__id_str,
		pCrawler.NameStr,
		pImagesLocalDirPathStr,

		pMediaDomainStr,
		pS3bucketNameStr,
		pUserID,
		pRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-------------------
	// STAGE - END
	
	end_time_f := float64(time.Now().UnixNano())/1000000000.0

	cycle_run := &GFcrawlerCycleRun{
		Id_str:               cycle_run__id_str,
		T_str:                "crawler_cycle_run",
		Creation_unix_time_f: cycle_run__creation_unix_time_f,
		Crawler_name_str:     pCrawler.NameStr,
		Target_domain_str:    domain_str,
		Target_url_str:       url_str,
		Start_time_f:         start_time_f,
		End_time_f:           end_time_f,
	}

	ctx           := context.Background()
	coll_name_str := "gf_crawl"
	gfErr         = gf_core.MongoInsert(cycle_run,
		coll_name_str,
		map[string]interface{}{
			"cycle_run__id_str":  cycle_run__id_str,
			"crawler_name_str":   pCrawler.NameStr,
			"domain_str":         domain_str,
			"caller_err_msg_str": "failed to insert a Crawler_cycle_run into the DB",
		},
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-------------------
	// LINK MARK AS RESOLVED

	// IMPORTANT!! - unresolved_link is nil if no links are present in DB or if all links have been resolved, 
	//              in which case the pCrawler.StartURLstr was used
	if unresolved_link != nil {
		gfErr = gf_crawl_core.DBmongoLinkMarkAsResolved(unresolved_link,
			url_fetch.Id_str,
			url_fetch.Creation_unix_time_f,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	//-------------------
	// unresolved_link - is nil if pCrawler.StartURLstr is used
	if unresolved_link != nil {

		// IMPORTANT!! - mark the link as no longer import_in_progress
		gfErr := gf_crawl_core.DBmongoLinkMarkImportInProgress(false, // p_status_bool
			end_time_f,
			unresolved_link,
			pRuntime,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}
	
	//-------------------
	return nil
}