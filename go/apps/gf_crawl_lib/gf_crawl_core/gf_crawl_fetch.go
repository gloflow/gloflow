/*
GloFlow media management/publishing system
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
	"github.com/globalsign/mgo/bson"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_crawl_lib/gf_crawl_utils"
)
//--------------------------------------------------
//ELASTIC_SEARCH - INDEXED
type Crawler_url_fetch struct {
	Id                   bson.ObjectId     `bson:"_id,omitempty"`
	Id_str               string            `bson:"id_str"               json:"id_str"`
	T_str                string            `bson:"t"                    json:"t"` //"crawler_url_fetch"
	Creation_unix_time_f float64           `bson:"creation_unix_time_f" json:"creation_unix_time_f"`
	Cycle_run_id_str     string            `bson:"cycle_run_id_str"     json:"cycle_run_id_str"`
	Domain_str           string            `bson:"domain_str"           json:"domain_str"`
	Url_str              string            `bson:"url_str"              json:"url_str"`
	Start_time_f         float64           `bson:"start_time_f"         json:"-"`
	End_time_f           float64           `bson:"end_time_f"           json:"-"`
	Page_text_str        string            `bson:"page_text_str"        json:"page_text_str"` //full text of the page html - indexed in ES
	goquery_doc          *goquery.Document `bson:"-"                    json:"-"`

	//-------------------
	//IMPORTANT!! - last error that occured/interupted processing of this link
	Error_type_str       string            `bson:"error_type_str,omitempty"`
	Error_id_str         string            `bson:"error_id_str,omitempty"`
	//-------------------
}
//--------------------------------------------------
func Fetch__url(p_url_str string,
	p_link             *Crawler_page_outgoing_link,
	p_cycle_run_id_str string,
	p_crawler_name_str string,
	p_runtime          *Crawler_runtime,
	p_runtime_sys      *gf_core.Runtime_sys) (*Crawler_url_fetch, string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_fetch.Fetch__url()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	p_runtime_sys.Log_fun("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))
	p_runtime_sys.Log_fun("INFO","FETCHING >> - "+yellow(p_url_str))
	p_runtime_sys.Log_fun("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))

	start_time_f := float64(time.Now().UnixNano())/1000000000.0

	//-------------------
	url,err := url.Parse(p_url_str)
	if err != nil {
		t:="fetcher_parse_url__failed"
		m:=fmt.Sprintf("failed to parse url for fetch - %s",p_url_str)

		gf_err := gf_core.Error__create(m,
			"url_parse_error",
			&map[string]interface{}{"url_str":p_url_str,},
			err, "gf_crawl_core", p_runtime_sys)

		_,fe_gf_err := fetch__error(t, m, p_url_str, p_link, p_crawler_name_str, gf_err, p_runtime, p_runtime_sys)
		if fe_gf_err != nil {
			return nil,"",fe_gf_err
		}

		return nil,"",gf_err
	}

	domain_str           := url.Host
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := "crawler_fetch__"+fmt.Sprint(creation_unix_time_f)
	fetch                := &Crawler_url_fetch{
		Id_str:              id_str,
		T_str:               "crawler_url_fetch",
		Creation_unix_time_f:creation_unix_time_f,
		Cycle_run_id_str:    p_cycle_run_id_str,
		Domain_str:          domain_str,
		Url_str:             p_url_str,
		Start_time_f:        start_time_f,
		//End_time_f          :end_time_f,
		//Page_text_str       :doc.Text(),
		//goquery_doc         :doc,
	}

	err = p_runtime_sys.Mongodb_coll.Insert(fetch)
	if err != nil {
		t:="fetch_record_persist__failed"
		m:=fmt.Sprintf("failed to DB persist Crawler_url_fetch struct of fetch for url - %s",p_url_str)
		
		gf_err := gf_core.Error__create(m,
			"mongodb_insert_error",
			&map[string]interface{}{"url_str":p_url_str,},
			err,"gf_crawl_core",p_runtime_sys)

		_,fe_gf_err := fetch__error(t,m,p_url_str,p_link,p_crawler_name_str,gf_err,p_runtime,p_runtime_sys)
		if fe_gf_err != nil {
			return nil,"",fe_gf_err
		}

		return nil,"",gf_err
	}
	//-------------------
	//HTTP REQUEST

	doc, gf_err := gf_crawl_utils.Get__html_doc_over_http(p_url_str, p_runtime_sys)

	if gf_err != nil {
		t := "fetch_url__failed"
		m := fmt.Sprintf("failed to HTTP fetch url - %s - err - %s", p_url_str, fmt.Sprint(gf_err.Error))
		
		crawler_error, fe_gf_err := fetch__error(t, m, p_url_str, p_link, p_crawler_name_str, gf_err, p_runtime, p_runtime_sys)
		if fe_gf_err != nil {
			return nil, "", fe_gf_err
		}

		fetch__mark_as_failed(crawler_error,
			fetch,
			p_runtime,
			p_runtime_sys)

		return nil, "", gf_err
	}

	end_time_f := float64(time.Now().UnixNano())/1000000000.0
	//-------------
	//UPDATE FETCH
	fetch.End_time_f    = end_time_f
	fetch.Page_text_str = doc.Text()
	fetch.goquery_doc   = doc
	err = p_runtime_sys.Mongodb_coll.Update(bson.M{"id_str":fetch.Id_str,"t":"crawler_url_fetch"},
		bson.M{"$set":bson.M{
			"end_time_f":   end_time_f,
			"page_text_str":doc.Text(),
		}})
	if err != nil {
		gf_err := gf_core.Error__create("failed to to update fetch record with end_time and page_text",
			"mongodb_update_error",
			&map[string]interface{}{"fetch_id_str":fetch.Id_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return nil,"",gf_err
	}
	//-------------
	//SEND_EVENT
	if p_runtime.Events_ctx != nil {
		events_id_str  := "crawler_events"
		event_type_str := "fetch__http_request__done"
		msg_str        := "completed fetching a document over HTTP"
		data_map       := map[string]interface{}{
			"url_str":     p_url_str,
			"start_time_f":start_time_f,
			"end_time_f":  end_time_f,
		}

		gf_core.Events__send_event(events_id_str,
			event_type_str, //p_type_str
			msg_str,        //p_msg_str
			data_map,
			p_runtime.Events_ctx,
			p_runtime_sys)
	}
	//-------------
	
	return fetch,domain_str,nil
}
//--------------------------------------------------
func Fetch__parse_result(p_url_fetch *Crawler_url_fetch,
	p_cycle_run_id_str          string,
	p_crawler_name_str          string,
	p_images_local_dir_path_str string,
	p_s3_bucket_name_str        string,
	p_runtime                   *Crawler_runtime,
	p_runtime_sys               *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_fetch.Fetch__parse_result()")

	//----------------
	//GET LINKS
	Links__get_outgoing_in_page(p_url_fetch,
		p_cycle_run_id_str,
		p_crawler_name_str,
		p_runtime,
		p_runtime_sys)
	//----------------
	//GET IMAGES
	images_pipe__from_html(p_url_fetch,
		p_cycle_run_id_str,
		p_crawler_name_str,
		p_images_local_dir_path_str,
		p_s3_bucket_name_str,
		p_runtime,
		p_runtime_sys)
	//----------------
	//INDEX URL_FETCH

	//IMPORTANT!! - index only if the indexer is initialized
	if p_runtime.Esearch_client != nil {
		gf_err := index__add_to__of_url_fetch(p_url_fetch, p_runtime, p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
	}
	//----------------

	return nil
}
//--------------------------------------------------
func fetch__error(p_error_type_str string,
	p_error_msg_str    string,
	p_url_str          string,
	p_link             *Crawler_page_outgoing_link,
	p_crawler_name_str string,
	p_gf_err           *gf_core.Gf_error,
	p_runtime          *Crawler_runtime,
	p_runtime_sys      *gf_core.Runtime_sys) (*Crawler_error,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_fetch.fetch__error()")

	crawler_error,ce_err := Create_error_and_event(p_error_type_str,
		p_error_msg_str,
		map[string]interface{}{}, p_url_str, p_crawler_name_str,
		p_gf_err,
		p_runtime,
		p_runtime_sys)
	if ce_err != nil {
		return nil,ce_err
	}

	if p_link != nil {
		//IMPORTANT!! - mark link as failed, so that it is not repeatedly tried
		lm_err := link__mark_as_failed(crawler_error, p_link, p_runtime, p_runtime_sys)
		if lm_err != nil {
			return nil,lm_err
		}
	}

	return crawler_error,nil
}
//--------------------------------------------------
func fetch__mark_as_failed(p_error *Crawler_error,
	p_fetch       *Crawler_url_fetch,
	p_runtime     *Crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_fetch.fetch__mark_as_failed()")

	p_fetch.Error_id_str   = p_error.Id_str
	p_fetch.Error_type_str = p_error.Type_str

	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"id_str":p_fetch.Id_str,
			"t":     "crawler_url_fetch",
		},
		bson.M{"$set":bson.M{
				"error_id_str":  p_error.Id_str,
				"error_type_str":p_error.Type_str,
			},
		})
	if err != nil {
		gf_err := gf_core.Error__create("failed to mark a crawler_url_fetch as failed in mongodb",
			"mongodb_update_error",
			&map[string]interface{}{"fetch_id_str":p_fetch.Id_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}

	return nil
}