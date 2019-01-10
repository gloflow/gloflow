package gf_crawl_stats

import (
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
type Stat__crawled_url_fetches struct {
	Url_str         string    `bson:"_id"             json:"url_str"`
	Count_int       int       `bson:"count_int"       json:"count_int"`
	Start_times_lst []float64 `bson:"start_times_lst" json:"start_times_lst"`
}
//-------------------------------------------------
func stats__crawler_fetches_by_days(p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_stats__fetches.stats__crawler_fetches_by_days()")

	stats__fetches_by_days,gf_err := stats__objs_by_days(map[string]interface{}{},
										"crawler_url_fetch",
										p_runtime_sys)
	if gf_err != nil {
		return nil,gf_err
	}

	data_map := map[string]interface{}{
		"fetches_by_days_map":stats__fetches_by_days,
	}
	return data_map,nil
}
//-------------------------------------------------
func stats__crawler_fetches_by_url(p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_stats__fetches.stats__crawler_fetches_by_url()")

	pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
		bson.M{"$match"  :bson.M{
				"t":"crawler_url_fetch",
			},
		},

		bson.M{"$project":bson.M{
				"id_str"      :true,
				"start_time_f":true,
				"end_time_f"  :true,
				"domain_str"  :true,
				"url_str"     :true, //actual link from the HTML <a> page ('href' parameter)
				//"errors_num_i":bson.M{"$size":"$errors_lst",},
			},
		},

		bson.M{"$group":bson.M{
				"_id"            :"$url_str",
				"count_int"      :bson.M{"$sum" :1},
				"start_times_lst":bson.M{"$push":"$start_time_f"},
			},
		},

		bson.M{"$sort":bson.M{
				"count_int":-1,
			},
		},
	})

	results_lst := []Stat__crawled_url_fetches{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to run an aggregation pipeline to group all crawler_url_fetch's",
			"mongodb_aggregation_error",
			nil,err,"gf_crawl_stats",p_runtime_sys)
		return nil,gf_err
	}
	
	data_map := map[string]interface{}{
		"crawled_url_fetches_lst":results_lst,
	}

	return data_map,nil
}