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
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
)

//--------------------------------------------------

type GFcrawlerError struct {
	Id                primitive.ObjectID     `bson:"_id,omitempty"    json:"-"`
	IDstr             string                 `bson:"id_str"           json:"id_str"`
	Tstr              string                 `bson:"t"                json:"t"` //"crawler_error"
	CreationUNIXtimeF float64                `bson:"creation_unix_time_f"`
	TypeStr           string                 `bson:"type_str"         json:"type_str"`
	MsgStr            string                 `bson:"msg_str"          json:"msg_str"` 
	DataMap           map[string]interface{} `bson:"data_map"         json:"data_map"` //if an error is related to a particular URL, it is noted here.
	GFerrorIDstr      string                 `bson:"gf_error_id_str"  json:"gf_error_id_str"`
	CrawlerNameStr    string                 `bson:"crawler_name_str" json:"crawler_name_str"`
	URLstr            string                 `bson:"url_str"          json:"url_str"`
}

//--------------------------------------------------

func CreateErrorAndEvent(pErrorTypeStr string,
	pErrorMsgStr    string,
	pErrorDataMap   map[string]interface{},
	pErrorURLstr    string,
	pCrawlerNameStr string,
	pGFerr          *gf_core.GFerror,
	pRuntime        *GFcrawlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFcrawlerError, *gf_core.GFerror) {

	if pRuntime.EventsCtx != nil {
		eventsIDstr  := "crawler_events"
		eventTypeStr := "error"

		gf_events.SendEvent(eventsIDstr,
			eventTypeStr,   // p_type_str
			pErrorMsgStr,   // pMsgStr
			pErrorDataMap,  // p_data_map
			pRuntime.EventsCtx,
			pRuntimeSys)
	}

	crawlErr, gfErr := createError(pErrorTypeStr,
		pErrorMsgStr,
		pErrorDataMap,
		pErrorURLstr,
		pCrawlerNameStr,
		pGFerr,
		pRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return crawlErr, nil
}

//--------------------------------------------------

func createError(pTypeStr string,
	pMsgStr         string,
	pDataMap        map[string]interface{},
	pURLstr         string,
	pCrawlerNameStr string,
	pGFerr          *gf_core.GFerror,
	pRuntime        *GFcrawlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFcrawlerError, *gf_core.GFerror) {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	IDstr             := fmt.Sprintf("crawl_error:%s", fmt.Sprint(creationUNIXtimeF))
	crawlErr          := &GFcrawlerError{
		IDstr:             IDstr,
		Tstr:              "crawler_error",
		CreationUNIXtimeF: creationUNIXtimeF,
		TypeStr:           pTypeStr,
		MsgStr:            pMsgStr,
		DataMap:           pDataMap,
		GFerrorIDstr:      pGFerr.Id_str,
		CrawlerNameStr:    pCrawlerNameStr,
		URLstr:            pURLstr,
	}

	ctx         := context.Background()
	collNameStr := "gf_crawl"
	gfErr := gf_core.MongoInsert(crawlErr,
		collNameStr,
		map[string]interface{}{
			"type_str":           pTypeStr,
			"crawler_name_str":   pCrawlerNameStr,
			"caller_err_msg_str": "failed to insert the crawler_error into the DB",
		},
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return crawlErr, nil
}