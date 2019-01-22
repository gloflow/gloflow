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
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
type Gf_crawler_error struct {
	Id                   bson.ObjectId          `bson:"_id,omitempty"    json:"-"`
	Id_str               string                 `bson:"id_str"           json:"id_str"`
	T_str                string                 `bson:"t"                json:"t"` //"crawler_error"
	Creation_unix_time_f float64                `bson:"creation_unix_time_f"`
	Type_str             string                 `bson:"type_str"         json:"type_str"`
	Msg_str              string                 `bson:"msg_str"          json:"msg_str"` 
	Data_map             map[string]interface{} `bson:"data_map"         json:"data_map"` //if an error is related to a particular URL, it is noted here.
	Gf_error_id_str      string                 `bson:"gf_error_id_str"  json:"gf_error_id_str"`
	Crawler_name_str     string                 `bson:"crawler_name_str" json:"crawler_name_str"`
	Url_str              string                 `bson:"url_str"          json:"url_str"`
}
//--------------------------------------------------
func Create_error_and_event(p_error_type_str string,
	p_error_msg_str    string,
	p_error_data_map   map[string]interface{},
	p_error_url_str    string,
	p_crawler_name_str string,
	p_gf_err           *gf_core.Gf_error,
	p_runtime          *Gf_crawler_runtime,
	p_runtime_sys      *gf_core.Runtime_sys) (*Gf_crawler_error, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_error.Create_error_and_event()")

	if p_runtime.Events_ctx != nil {
		events_id_str  := "crawler_events"
		event_type_str := "error"

		gf_core.Events__send_event(events_id_str,
			event_type_str,   //p_type_str
			p_error_msg_str,  //p_msg_str
			p_error_data_map, //p_data_map
			p_runtime.Events_ctx,
			p_runtime_sys)
	}

	crawl_err, gf_err := create_error(p_error_type_str,
		p_error_msg_str,
		p_error_data_map,
		p_error_url_str,
		p_crawler_name_str,
		p_gf_err,
		p_runtime,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return crawl_err, nil
}
//--------------------------------------------------
func create_error(p_type_str string,
	p_msg_str          string,
	p_data_map         map[string]interface{},
	p_url_str          string,
	p_crawler_name_str string,
	p_gf_err           *gf_core.Gf_error,
	p_runtime          *Gf_crawler_runtime,
	p_runtime_sys      *gf_core.Runtime_sys) (*Gf_crawler_error, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_error.create_error()")

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := "crawl_error:"+fmt.Sprint(creation_unix_time_f)
	crawl_err            := &Gf_crawler_error{
		Id_str:              id_str,
		T_str:               "crawler_error",
		Creation_unix_time_f:creation_unix_time_f,
		Type_str:            p_type_str,
		Msg_str:             p_msg_str,
		Data_map:            p_data_map,
		Gf_error_id_str:     p_gf_err.Id_str,
		Crawler_name_str:    p_crawler_name_str,
		Url_str:             p_url_str,
	}

	if p_runtime.Cluster_node_type_str == "master" {
		err := p_runtime_sys.Mongodb_coll.Insert(crawl_err)
		if err != nil {
			gf_err := gf_core.Error__create("failed to persist a crawler_error",
				"mongodb_insert_error",
				&map[string]interface{}{
					"type_str":         p_type_str,
					"crawler_name_str": p_crawler_name_str,
				},
				err, "gf_crawl_core", p_runtime_sys)
			return nil, gf_err
		}
	}
	return crawl_err, nil
}