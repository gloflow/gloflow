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
func Run_crawler_cycle(p_crawler gf_crawl_core.GFcrawlerDef,
	pImagesLocalDirPathStr string,
	pMediaDomainStr        string,
	pS3bucketNameStr       string,
	p_runtime              *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_cycle.Run_crawler_cycle()")
	pRuntimeSys.Log_fun("INFO",      "pS3bucketNameStr - "+pS3bucketNameStr)

	yellow := color.New(color.FgYellow).SprintFunc()
	black  := color.New(color.FgBlack).Add(color.BgWhite).SprintFunc()

	start_time_f := float64(time.Now().UnixNano())/1000000000.0

	//-------------------------
	// STAGE - GET TARGET URL


	// IMPORTANT!! - get unresolved links to pages on the domain to which the crawler belongs.
	//               so if the a page contains links to domains external to the domain to which the 
	//               crawler belongs, it wont get fetched/parsed here
	unresolved_link, gfErr := gf_crawl_core.Link__db_get_unresolved(p_crawler.Name_str, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// // IMPORTANT!! - no unresolved links were found, this is a valid possible state. 
	// //               do nothing, because further bellow unresolved_link is tested for == nil and then
	// //               crawlenrs start_url is used - hence the empty IF block.
	// if gf_err != nil && fmt.Sprint(gf_err.Type_str) == "mongodb_not_found_error" {

	// } else if gf_err != nil {
	// 	return gf_err
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
		gfErr       := gf_crawl_core.Link__db_mark_import_in_progress(true, start_time_f, unresolved_link, p_runtime, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		//----------------------
	} else {

		// IMPORTANT!! - if all links were resolved, then start at the initial Crawler Start_url_str
		//               and begin a new sweep of crawling the domain
		url_str = p_crawler.Start_url_str
		fmt.Println("INFO",black(">>>>>>>>>>>>>>>>> UNRESOLVED_LINK NOT FOUND - using start_url - ")+yellow(url_str))
	}
	
	//-------------------------
	// CYCLE_RUN
	cycle_run__creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	cycle_run__id_str               := "crawler_cycle_run:"+fmt.Sprint(cycle_run__creation_unix_time_f)

	//-------------------------
	
	//-------------------
	// STAGE - FETCH THE LINK
	url_fetch, domain_str, gfErr := gf_crawl_core.Fetch__url(url_str,
		unresolved_link,
		cycle_run__id_str,
		p_crawler.Name_str,
		p_runtime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-------------------
	// STAGE - PARSE THE FETCHED PAGE
	gfErr = gf_crawl_core.Fetch__parse_result(url_fetch,
		cycle_run__id_str,
		p_crawler.Name_str,
		pImagesLocalDirPathStr,

		pMediaDomainStr,
		pS3bucketNameStr,
		p_runtime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-------------------
	// STAGE - END
	
	end_time_f := float64(time.Now().UnixNano())/1000000000.0

	cycle_run := &Gf_crawler_cycle_run{
		Id_str:               cycle_run__id_str,
		T_str:                "crawler_cycle_run",
		Creation_unix_time_f: cycle_run__creation_unix_time_f,
		Crawler_name_str:     p_crawler.Name_str,
		Target_domain_str:    domain_str,
		Target_url_str:       url_str,
		Start_time_f:         start_time_f,
		End_time_f:           end_time_f,
	}

	ctx           := context.Background()
	coll_name_str := "gf_crawl"
	gfErr         = gf_core.Mongo__insert(cycle_run,
		coll_name_str,
		map[string]interface{}{
			"cycle_run__id_str":  cycle_run__id_str,
			"crawler_name_str":   p_crawler.Name_str,
			"domain_str":         domain_str,
			"caller_err_msg_str": "failed to insert a Crawler_cycle_run into the DB",
		},
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	/*err := pRuntimeSys.Mongodb_db.C("gf_crawl").Insert(cycle_run)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a Crawler_cycle_run in mongodb",
			"mongodb_insert_error",
			map[string]interface{}{
				"cycle_run__id_str": cycle_run__id_str,
				"crawler_name_str":  p_crawler.Name_str,
				"domain_str":        domain_str,
			},
			err, "gf_crawl_lib", pRuntimeSys)
		return gf_err
	}*/

	//-------------------
	// LINK MARK AS RESOLVED

	// IMPORTANT!! - unresolved_link is nil if no links are present in DB or if all links have been resolved, 
	//              in which case the p_crawler.Start_url_str was used
	if unresolved_link != nil {
		gfErr = gf_crawl_core.Link__db_mark_as_resolved(unresolved_link,
			url_fetch.Id_str,
			url_fetch.Creation_unix_time_f,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	//-------------------
	// unresolved_link - is nil if p_crawler.Start_url_str is used
	if unresolved_link != nil {

		// IMPORTANT!! - mark the link as no longer import_in_progress
		gfErr := gf_crawl_core.Link__db_mark_import_in_progress(false, //p_status_bool
			end_time_f,
			unresolved_link,
			p_runtime,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}
	
	//-------------------
	return nil
}