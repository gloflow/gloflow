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

package gf_crawl_stats

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type Gf_stat__crawled_url_fetches struct {
	Url_str         string    `bson:"_id"             json:"url_str"`
	Count_int       int       `bson:"count_int"       json:"count_int"`
	Start_times_lst []float64 `bson:"start_times_lst" json:"start_times_lst"`
}

//-------------------------------------------------

func stats__crawler_fetches_by_days(pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_stats__fetches.stats__crawler_fetches_by_days()")

	stats__fetches_by_days, gf_err := stats__objs_by_days(map[string]interface{}{}, "crawler_url_fetch", pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	data_map := map[string]interface{}{
		"fetches_by_days_map": stats__fetches_by_days,
	}
	return data_map, nil
}

//-------------------------------------------------

func stats__crawler_fetches_by_url(pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_stats__fetches.stats__crawler_fetches_by_url()")


	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{"t": "crawler_url_fetch"}},
		},
		{
			{"$project", bson.M{
				"id_str":       true,
				"start_time_f": true,
				"end_time_f":   true,
				"domain_str":   true,
				"url_str":      true, // actual link from the HTML <a> page ('href' parameter)
				// "errors_num_i":bson.M{"$size":"$errors_lst",},
			}},
		},
		{
			{"$group", bson.M{
				"_id":             "$url_str",
				"count_int":       bson.M{"$sum":  1},
				"start_times_lst": bson.M{"$push": "$start_time_f"},
			}},
		},
		{
			{"$sort", bson.M{"count_int": -1}},
		},
	}


	/*pipe := pRuntimeSys.Mongo_db.Collection("gf_crawl").Pipe([]bson.M{
		bson.M{"$match": bson.M{
				"t":     "crawler_url_fetch",
			},
		},

		bson.M{"$project": bson.M{
				"id_str":       true,
				"start_time_f": true,
				"end_time_f":   true,
				"domain_str":   true,
				"url_str":      true, // actual link from the HTML <a> page ('href' parameter)
				// "errors_num_i":bson.M{"$size":"$errors_lst",},
			},
		},

		bson.M{"$group": bson.M{
				"_id":             "$url_str",
				"count_int":       bson.M{"$sum":  1},
				"start_times_lst": bson.M{"$push": "$start_time_f"},
			},
		},

		bson.M{"$sort": bson.M{
				"count_int": -1,
			},
		},
	})*/

	cursor, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").Aggregate(ctx, pipeline)
	if err != nil {

		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to group all crawler_url_fetch's",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_crawl_stats", pRuntimeSys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)

	/*results_lst := []Gf_stat__crawled_url_fetches{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to group all crawler_url_fetch's",
			"mongodb_aggregation_error",
			nil, err, "gf_crawl_stats", pRuntimeSys)
		return nil, gf_err
	}*/
	
	results_lst := []Gf_stat__crawled_url_fetches{}
	for cursor.Next(ctx) {

		var r Gf_stat__crawled_url_fetches
		err := cursor.Decode(&r)
		if err != nil {
			gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to group all crawler_url_fetch's",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_crawl_stats", pRuntimeSys)
			return nil, gf_err
		}
		results_lst = append(results_lst, r)
	}

	data_map := map[string]interface{}{
		"crawled_url_fetches_lst": results_lst,
	}

	return data_map, nil
}