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

package gf_crawl_stats

import (
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
type Stat__errors struct {
	Crawler_name_str string             `bson:"_id"              json:"crawler_name_str"`
	Errors_types_lst []Stat__error_type `bson:"errors_types_lst" json:"errors_types_lst"`
}
type Stat__error_type struct {
	Type_str  string   `bson:"type_str"  json:"type_str"`
	Count_int int      `bson:"count_int" json:"count_int"`
	Urls_lst  []string `bson:"urls_lst"  json:"urls_lst"`
}
//-------------------------------------------------
func stats__errors(p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_stats__errors.stats__errors()")

	pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
		bson.M{"$match":bson.M{
				"t":"crawler_error",
			},
		},

		bson.M{"$sort":bson.M{
				"creation_unix_time_f":-1,
			},
		},

		bson.M{"$group":bson.M{
				"_id"                    :bson.M{"type_str":"$type_str","crawler_name_str":"$crawler_name_str",},
				"count_int"              :bson.M{"$sum":1},
				"urls_lst"               :bson.M{"$push":"$url_str"},
				"creation_unix_times_lst":bson.M{"$push":"$creation_unix_time_f"},
			},
		},

		bson.M{"$group":bson.M{
				"_id"             :"$_id.crawler_name_str",
				"errors_types_lst":bson.M{"$push":bson.M{
											"type_str"               :"$_id.type_str",
											"count_int"              :"$count_int",
											"urls_lst"               :"$urls_lst",
											"creation_unix_times_lst":"$creation_unix_times_lst",
										},
									},
			},
		},
	})

	results_lst := []Stat__errors{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to run an aggregation pipeline to count/get_info of crawler_error's by crawler_name",
			"mongodb_aggregation_error",
			nil,err,"gf_crawl_stats",p_runtime_sys)
		return nil,gf_err
	}
	
	data_map := map[string]interface{}{
		"errors_lst":results_lst,
	}
	return data_map,nil
}