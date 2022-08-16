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
type Gf_stat__crawled_images_domain struct {
	Domain_str              string    `bson:"_id"                     json:"domain_str"`
	Imgs_count_int          int       `bson:"imgs_count_int"          json:"imgs_count_int"`
	Creation_unix_times_lst []float64 `bson:"creation_unix_times_lst" json:"creation_unix_times_lst"`
	Urls_lst                []string  `bson:"urls_lst"                json:"urls_lst"`
	Origin_urls_lst         []string  `bson:"origin_urls_lst"         json:"origin_urls_lst"`
	Downloaded_lst          []bool    `bson:"downloaded_lst"          json:"downloaded_lst"`
	Valid_for_usage_lst     []bool    `bson:"valid_for_usage_lst"     json:"valid_for_usage_lst"`
	S3_stored_lst           []bool    `bson:"s3_stored_lst"           json:"s3_stored_lst"`
}

type Gf_stat__crawled_gifs struct {
	Domain_str             string                   `bson:"_id"                    json:"domain_str"`
	Imgs_count_int         int                      `bson:"imgs_count_int"         json:"imgs_count_int"`
	Urls_by_origin_url_lst []map[string]interface{} `bson:"urls_by_origin_url_lst" json:"urls_by_origin_url_lst"`
}

//-------------------------------------------------
func stats__gifs_by_days(p_runtime_sys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_crawl_stats__images.stats__gifs_by_days()")


	stats__gifs_by_days, gf_err := stats__objs_by_days(map[string]interface{}{"img_ext_str": "gif",}, "crawler_page_img", p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	data_map := map[string]interface{}{
		"gifs_by_days_map": stats__gifs_by_days,
	}
	return data_map, nil
}

//-------------------------------------------------
func stats__gifs(p_runtime_sys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.LogFun("FUN_ENTER","gf_crawl_stats__images.stats__gifs()")


	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{
				"t":           "crawler_page_img",
				"img_ext_str": "gif",
			}},
		},
		{
			{"$project", bson.M{
				"domain_str":           true,
				"creation_unix_time_f": true,
				"origin_url_str":       true,
				"url_str":              true,
				"nsfv_bool":            true,
			}},
		},
		{
			{"$group", bson.M{
				"_id":                bson.M{"origin_url_str": "$origin_url_str", "domain_str": "$domain_str",},
				"creation_times_lst": bson.M{"$push": "$creation_unix_time_f"},
				"urls_lst":           bson.M{"$push": "$url_str"},
				"nsfv_lst":           bson.M{"$push": "$nsfv_bool"},
				"count_int":          bson.M{"$sum":  1},
			}},
		},
		{
			{"$group", bson.M{
				"_id":                        "$_id.domain_str",
				"imgs_count_int":             bson.M{"$sum":  "$count_int",}, // add up counts from the previous grouping operation
				"urls_by_origin_url_lst":     bson.M{"$push": bson.M{
						"origin_url_str":     "$_id.origin_url_str",
						"creation_times_lst": "$creation_times_lst",
						"urls_lst":           "$urls_lst",
						"nsfv_lst":           "$nsfv_lst",
					},
				},
			}},
		},
		{
			{"$sort", bson.M{
				"imgs_count_int": -1,
			}},
		},
	}


	/*pipe := p_runtime_sys.Mongo_db.Collection("gf_crawl").Pipe([]bson.M{
		bson.M{"$match"  :bson.M{
				"t":          "crawler_page_img",
				"img_ext_str":"gif",
			},
		},

		// bson.M{"$project":bson.M{
		// 		"id_str"              :true,
		// 		"creation_unix_time_f":true,
		// 		"cycle_run_id_str"    :true,
		// 		"domain_str"          :true,
		// 		"url_str"             :true,
		// 		"origin_url_str"      :true, //page url from whos html this element was extracted
		// 		"downloaded_bool"     :true,
		// 		"valid_for_usage_bool":true,
		// 		"s3_stored_bool"      :true,
		// 
		// 		//"errors_num_i"        :bson.M{"$size":"$errors_lst",},
		// 	},
		// },

		bson.M{"$project":bson.M{
				"domain_str":          true,
				"creation_unix_time_f":true,
				"origin_url_str":      true,
				"url_str":             true,
				"nsfv_bool":           true,
			},
		},

		bson.M{"$group":bson.M{
				"_id":                bson.M{"origin_url_str": "$origin_url_str", "domain_str": "$domain_str",},
				"creation_times_lst": bson.M{"$push": "$creation_unix_time_f"},
				"urls_lst":           bson.M{"$push": "$url_str"},
				"nsfv_lst":           bson.M{"$push": "$nsfv_bool"},
				"count_int":          bson.M{"$sum":  1},
			},
		},

		bson.M{"$group":bson.M{
				"_id":                        "$_id.domain_str",
				"imgs_count_int":             bson.M{"$sum":  "$count_int",}, // add up counts from the previous grouping operation
				"urls_by_origin_url_lst":     bson.M{"$push": bson.M{
						"origin_url_str":     "$_id.origin_url_str",
						"creation_times_lst": "$creation_times_lst",
						"urls_lst":           "$urls_lst",
						"nsfv_lst":           "$nsfv_lst",
					},
				},
			},
		},

		bson.M{"$sort":bson.M{
				"imgs_count_int": -1,
			},
		},
	})*/

	cursor, err := p_runtime_sys.Mongo_db.Collection("gf_crawl").Aggregate(ctx, pipeline)
	if err != nil {

		gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to get GIF's (crawler_page_img) by domain",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_crawl_stats", p_runtime_sys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)

	/*results_lst := []Gf_stat__crawled_gifs{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to get GIF's (crawler_page_img) by domain",
			"mongodb_aggregation_error",
			nil, err, "gf_crawl_stats", p_runtime_sys)
		return nil, gf_err
	}*/

	results_lst := []Gf_stat__crawled_gifs{}
	for cursor.Next(ctx) {

		var r Gf_stat__crawled_gifs
		err := cursor.Decode(&r)
		if err != nil {
			gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to get GIF's (crawler_page_img) by domain",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_crawl_stats", p_runtime_sys)
			return nil, gf_err
		}
	
		results_lst = append(results_lst, r)
	}


	data_map := map[string]interface{}{
		"crawled_gifs_lst": results_lst,
	}
	return data_map, nil
}

//-------------------------------------------------
func stats__crawled_images_domains(p_runtime_sys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_crawl_stats__images.stats__crawled_images_domains()")


	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{"t": "crawler_page_img"}},
		},
		{
			{"$project", bson.M{
				"id_str":               true,
				"creation_unix_time_f": true,
				"cycle_run_id_str":     true,
				"domain_str":           true,
				"url_str":              true,
				"origin_url_str":       true, // page url from whos html this element was extracted
				"downloaded_bool":      true,
				"valid_for_usage_bool": true,
				"s3_stored_bool":       true,
			}},
		},
		{
			{"$group", bson.M{
				"_id":                     "$domain_str",
				"imgs_count_int":          bson.M{"$sum":      1},
				"creation_unix_times_lst": bson.M{"$push":     "$creation_unix_time_f"},
				"urls_lst":                bson.M{"$push":     "$url_str"},
				"origin_urls_lst":         bson.M{"$addToSet": "$origin_url_str"},
				"downloaded_lst":          bson.M{"$push":     "$downloaded_bool"},
				"valid_for_usage_lst":     bson.M{"$push":     "$valid_for_usage_bool"},
				"s3_stored_lst":           bson.M{"$push":     "$s3_stored_bool"},
			}},
		},
		{
			{"$sort", bson.M{"imgs_count_int": -1}},
		},

	}


	/*pipe := p_runtime_sys.Mongo_db.Collection("gf_crawl").Pipe([]bson.M{
		bson.M{"$match":bson.M{
				"t": "crawler_page_img",
			},
		},
		bson.M{"$project":bson.M{
				"id_str":               true,
				"creation_unix_time_f": true,
				"cycle_run_id_str":     true,
				"domain_str":           true,
				"url_str":              true,
				"origin_url_str":       true, //page url from whos html this element was extracted
				"downloaded_bool":      true,
				"valid_for_usage_bool": true,
				"s3_stored_bool":       true,
			},
		},
		bson.M{"$group":bson.M{
				"_id":                     "$domain_str",
				"imgs_count_int":          bson.M{"$sum":      1},
				"creation_unix_times_lst": bson.M{"$push":     "$creation_unix_time_f"},
				"urls_lst":                bson.M{"$push":     "$url_str"},
				"origin_urls_lst":         bson.M{"$addToSet": "$origin_url_str"},
				"downloaded_lst":          bson.M{"$push":     "$downloaded_bool"},
				"valid_for_usage_lst":     bson.M{"$push":     "$valid_for_usage_bool"},
				"s3_stored_lst":           bson.M{"$push":     "$s3_stored_bool"},
			},
		},
		// bson.M{"$group":bson.M{
		// 		"_id"           :"$_id.cycle_run_id_str",
		// 		"imgs_count_int":bson.M{"$sum"     :1},
		// 		""
		// 	},
		// },
		bson.M{"$sort": bson.M{
				"imgs_count_int": -1,
			},
		},
	})*/

	cursor, err := p_runtime_sys.Mongo_db.Collection("gf_crawl").Aggregate(ctx, pipeline)
	if err != nil {

		gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to get crawler_page_imgs by domain",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_crawl_stats", p_runtime_sys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)

	/*results_lst := []Gf_stat__crawled_images_domain{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to get crawler_page_imgs by domain",
			"mongodb_aggregation_error",
			nil, err, "gf_crawl_stats", p_runtime_sys)
		return nil, gf_err
	}*/


	results_lst := []Gf_stat__crawled_images_domain{}
	for cursor.Next(ctx) {

		var r Gf_stat__crawled_images_domain
		err := cursor.Decode(&r)
		if err != nil {
			gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to get crawler_page_imgs by domain",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_crawl_stats", p_runtime_sys)
			return nil, gf_err
		}
	
		results_lst = append(results_lst, r)
	}


	data_map := map[string]interface{}{
		"crawled_images_domains_lst": results_lst,
	}
	return data_map, nil
}