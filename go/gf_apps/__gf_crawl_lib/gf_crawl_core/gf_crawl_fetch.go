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

package gf_crawl_core

import (
	"fmt"
	"time"
	"net/url"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_utils"
)

//--------------------------------------------------
// ELASTIC_SEARCH - INDEXED

type GFcrawlerURLfetch struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               string             `bson:"id_str"               json:"id_str"`
	T_str                string             `bson:"t"                    json:"t"` // "crawler_url_fetch"
	Creation_unix_time_f float64            `bson:"creation_unix_time_f" json:"creation_unix_time_f"`
	Cycle_run_id_str     string             `bson:"cycle_run_id_str"     json:"cycle_run_id_str"`
	Domain_str           string             `bson:"domain_str"           json:"domain_str"`
	Url_str              string             `bson:"url_str"              json:"url_str"`
	Start_time_f         float64            `bson:"start_time_f"         json:"-"`
	End_time_f           float64            `bson:"end_time_f"           json:"-"`
	Page_text_str        string             `bson:"page_text_str"        json:"page_text_str"` // full text of the page html - indexed in ES
	goquery_doc          *goquery.Document  `bson:"-"                    json:"-"`

	//-------------------
	// IMPORTANT!! - last error that occured/interupted processing of this link
	Error_type_str       string            `bson:"error_type_str,omitempty"`
	Error_id_str         string            `bson:"error_id_str,omitempty"`

	//-------------------
}

//--------------------------------------------------
// FETCH_URL

func FetchURL(pURLstr string,
	pLink           *GFcrawlerPageOutgoingLink,
	pCycleRunIDstr  string,
	pCrawlerNameStr string,
	pRuntime        *GFcrawlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFcrawlerURLfetch, string, *gf_core.GFerror) {

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	pRuntimeSys.LogFun("INFO", cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))
	pRuntimeSys.LogFun("INFO", "FETCHING >> - "+yellow(pURLstr))
	pRuntimeSys.LogFun("INFO", cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))

	start_time_f := float64(time.Now().UnixNano())/1000000000.0

	//-------------------
	url, err := url.Parse(pURLstr)
	if err != nil {
		t := "fetcher_parse_url__failed"
		m := fmt.Sprintf("failed to parse url for fetch - %s", pURLstr)

		gfErr := gf_core.ErrorCreate(m,
			"url_parse_error",
			map[string]interface{}{"url_str": pURLstr,},
			err, "gf_crawl_core", pRuntimeSys)

		_, fe_gf_err := fetchError(t, m, pURLstr, pLink, pCrawlerNameStr, gfErr, pRuntime, pRuntimeSys)
		if fe_gf_err != nil {
			return nil, "", fe_gf_err
		}

		return nil, "", gfErr
	}

	domain_str           := url.Host
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := "crawler_fetch__"+fmt.Sprint(creation_unix_time_f)
	fetch                := &GFcrawlerURLfetch{
		Id_str:               id_str,
		T_str:                "crawler_url_fetch",
		Creation_unix_time_f: creation_unix_time_f,
		Cycle_run_id_str:     pCycleRunIDstr,
		Domain_str:           domain_str,
		Url_str:              pURLstr,
		Start_time_f:         start_time_f,
		// End_time_f           : end_time_f,
		// Page_text_str        : doc.Text(),
		// goquery_doc          : doc,
	}

	t := "fetch_record_persist__failed"
	m := fmt.Sprintf("failed to DB persist GFcrawlerURLfetch struct of fetch for url - %s", pURLstr)

	ctx           := context.Background()
	coll_name_str := "gf_crawl"
	gfErr         := gf_core.MongoInsert(fetch,
		coll_name_str,
		map[string]interface{}{
			"url_str":            pURLstr,
			"caller_err_msg_str": m,
		},
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		
		_, feGFerr := fetchError(t, m, pURLstr, pLink, pCrawlerNameStr, gfErr, pRuntime, pRuntimeSys)
		if feGFerr != nil {
			return nil, "", feGFerr
		}

		return nil, "", gfErr
	}
	
	//-------------------
	// HTTP REQUEST

	doc, gfErr := gf_crawl_utils.Get__html_doc_over_http(pURLstr, pRuntimeSys)

	if gfErr != nil {
		t := "fetch_url__failed"
		m := fmt.Sprintf("failed to HTTP fetch url - %s - err - %s", pURLstr, fmt.Sprint(gfErr.Error))
		
		crawler_error, feGFerr := fetchError(t, m, pURLstr, pLink, pCrawlerNameStr, gfErr, pRuntime, pRuntimeSys)
		if feGFerr != nil {
			return nil, "", feGFerr
		}

		fetchMarkAsFailed(crawler_error, fetch, pRuntime, pRuntimeSys)

		return nil, "", gfErr
	}

	end_time_f := float64(time.Now().UnixNano())/1000000000.0

	//-------------
	// UPDATE FETCH
	fetch.End_time_f    = end_time_f
	fetch.Page_text_str = doc.Text()
	fetch.goquery_doc   = doc
	_, err = pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"id_str": fetch.Id_str,
			"t":      "crawler_url_fetch",
		},
		bson.M{"$set": bson.M{
			"end_time_f":    end_time_f,
			"page_text_str": doc.Text(),
		}})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to to update fetch record with end_time and page_text",
			"mongodb_update_error",
			map[string]interface{}{"fetch_id_str":fetch.Id_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return nil, "", gfErr
	}

	//-------------
	// SEND_EVENT
	if pRuntime.EventsCtx != nil {
		events_id_str  := "crawler_events"
		event_type_str := "fetch__http_request__done"
		msg_str        := "completed fetching a document over HTTP"
		data_map       := map[string]interface{}{
			"url_str":      pURLstr,
			"start_time_f": start_time_f,
			"end_time_f":   end_time_f,
		}

		gf_events.SendEvent(events_id_str,
			event_type_str, // p_type_str
			msg_str,        // p_msg_str
			data_map,
			pRuntime.EventsCtx,
			pRuntimeSys)
	}

	//-------------
	
	return fetch, domain_str, nil
}

//--------------------------------------------------
// FETCH__PARSE_RESULT

func FetchParseResult(pURLfetch *GFcrawlerURLfetch,
	pCycleRunIDstr         string,
	pCrawlerNameStr        string,
	pImagesLocalDirPathStr string,

	pMediaDomainStr        string,
	pS3bucketNameStr       string,
	pUserID                gf_core.GF_ID,
	pRuntime               *GFcrawlerRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) *gf_core.GFerror {

	//----------------
	// GET LINKS
	LinksGetOutgoingInPage(pURLfetch,
		pCycleRunIDstr,
		pCrawlerNameStr,
		pRuntime,
		pRuntimeSys)

	//----------------
	// GET IMAGES
	imagesPipeFromHTML(pURLfetch,
		pCycleRunIDstr,
		pCrawlerNameStr,
		pImagesLocalDirPathStr,

		pMediaDomainStr,
		pS3bucketNameStr,
		pUserID,
		pRuntime,
		pRuntimeSys)

	//----------------
	/*
	// INDEX URL_FETCH

	// IMPORTANT!! - index only if the indexer is initialized
	if pRuntime.EsearchClient != nil {
		gfErr := indexAddToOfURLfetch(pURLfetch, pRuntime, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}
	*/
	
	//----------------

	return nil
}

//--------------------------------------------------

func fetchError(p_error_type_str string,
	p_error_msg_str string,
	pURLstr         string,
	pLink           *GFcrawlerPageOutgoingLink,
	pCrawlerNameStr string,
	pGFerr          *gf_core.GFerror,
	pRuntime        *GFcrawlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFcrawlerError, *gf_core.GFerror) {

	crawler_error, ce_err := CreateErrorAndEvent(p_error_type_str,
		p_error_msg_str,
		map[string]interface{}{}, pURLstr, pCrawlerNameStr,
		pGFerr,
		pRuntime,
		pRuntimeSys)
	if ce_err != nil {
		return nil, ce_err
	}

	if pLink != nil {
		// IMPORTANT!! - mark link as failed, so that it is not repeatedly tried
		lm_err := dbMongoLinkMarkAsFailed(crawler_error, pLink, pRuntime, pRuntimeSys)
		if lm_err != nil {
			return nil, lm_err
		}
	}

	return crawler_error, nil
}

//--------------------------------------------------

func fetchMarkAsFailed(pError *GFcrawlerError,
	pFetch      *GFcrawlerURLfetch,
	pRuntime    *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()

	pFetch.Error_id_str   = pError.IDstr
	pFetch.Error_type_str = pError.TypeStr

	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"id_str": pFetch.Id_str,
			"t":      "crawler_url_fetch",
		},
		bson.M{"$set": bson.M{
				"error_id_str":   pError.IDstr,
				"error_type_str": pError.TypeStr,
			},
		})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to mark a crawler_url_fetch as failed in mongodb",
			"mongodb_update_error",
			map[string]interface{}{"fetch_id_str": pFetch.Id_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return gfErr
	}

	return nil
}