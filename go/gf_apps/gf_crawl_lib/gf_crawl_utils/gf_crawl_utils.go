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

package gf_crawl_utils

import (
	"fmt"
	"time"
	"bytes"
	"math/rand"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
func Get__html_doc_over_http(pURLstr string,
	pRuntimeSys *gf_core.RuntimeSys) (*goquery.Document, *gf_core.GFerror) {

	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_utils.Get__html_doc_over_http()")

	//-----------------------
	userAgentStr := "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1"
	headersMap   := map[string]string{}
	
	HTTPfetch, gfErr := gf_core.HTTPfetchURL(pURLstr,
		headersMap,
		userAgentStr,
		context.Background(),
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	defer HTTPfetch.Resp.Body.Close()
	
	//-----------------------
	
	if !(HTTPfetch.Status_code_int >= 200 && HTTPfetch.Status_code_int < 400) {

		// GET_RESPONSE_BODY
		buff := new(bytes.Buffer)
		buff.ReadFrom(HTTPfetch.Resp.Body)
		body_str := buff.String() 

		gfErr := gf_core.Error__create("crawler fetch failed with HTTP status error",
			"http_client_req_status_error",
			map[string]interface{}{
				"url_str":         pURLstr,
				"status_code_int": HTTPfetch.Status_code_int,
				"body_str":        body_str,
			},
			nil, "gf_crawl_utils", pRuntimeSys)
		return nil, gfErr
	}

	// doc, err := goquery.NewDocument(pURLstr)
	doc, err := goquery.NewDocumentFromResponse(HTTPfetch.Resp)
	if err != nil {
		gfErr := gf_core.Error__create("failed to parse a fetched HTML page from a crawled url",
			"html_parse_error",
			map[string]interface{}{"url_str": pURLstr,},
			err, "gf_crawl_utils", pRuntimeSys)
		return nil, gfErr
	}

	return doc, nil
}

//--------------------------------------------------
func Crawler_sleep(p_crawler_name_str string,
	p_cycle_index_int int,
	p_rand            *rand.Rand,
	pRuntimeSys       *gf_core.RuntimeSys) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_utils.Crawler_sleep()")

	black  := color.New(color.FgBlack).Add(color.BgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	var wait_random_sec int

	//--------------
	// EVERY 100'th CYCLE
	// "p_cycle_index_int != 0" - since on the first cycle_index ("0") module divide "%"
	//                            is always == 0
	if p_cycle_index_int != 0 && p_cycle_index_int%100 == 0 {
		wait_random_sec = 60*60*2 // run every 2h

	//--------------
	// EVERY 200'th CYCLE
	} else if p_cycle_index_int != 0 && p_cycle_index_int % 200 == 0 {
		wait_random_sec = 60*60*5 // run every 5h

	//--------------
	// EVERY OTHER CYCLE
	} else {
		wait_random_sec = p_rand.Intn(60*10) //runs at least every 10min (or less)
	}


	wait_random_min := float32(wait_random_sec)/float32(60)
	sleep_length    := time.Second * time.Duration(wait_random_sec)

	fmt.Println("INFO", black(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))
	fmt.Println("INFO", black(">>>    SLEEPING CRAWLER >>> ")+yellow(p_crawler_name_str)+black(" - for min - "+fmt.Sprint(wait_random_min)))
	fmt.Println("INFO", black(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))

	time.Sleep(sleep_length)
}
