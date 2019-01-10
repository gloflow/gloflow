package gf_crawl_stats

import (
	"fmt"
	"time"
	"strconv"
	"sort"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
type Stat__crawled_links_domain struct {
	Domain_str              string    `bson:"_id"                     json:"domain_str"`
	Links_count_int         int       `bson:"links_count_int"         json:"links_count_int"`
	Creation_unix_times_lst []float64 `bson:"creation_unix_times_lst" json:"creation_unix_times_lst"`
	A_href_lst              []string  `bson:"a_href_lst"              json:"a_href_lst"`
	Origin_urls_lst         []string  `bson:"origin_urls_lst"         json:"origin_urls_lst"`
	Valid_for_crawl_lst     []bool    `bson:"valid_for_crawl_lst"     json:"valid_for_crawl_lst"`  //if the link is to be crawled/followed, or should be ignored
	Fetched_lst             []bool    `bson:"fetched_lst"             json:"fetched_lst"`          //if the link's HTML was downloaded
	Images_processed_lst    []bool    `bson:"images_processed_lst"    json:"images_processed_lst"` //if the images of this links HTML page were downloaded/processed
}

type Stat__unresolved_links struct {
	Origin_domain_str             string     `bson:"_id"                           json:"origin_domain_str"`
	Origin_urls_lst               []string   `bson:"origin_urls_lst"               json:"origin_urls_lst"`
	Counts__from_origin_urls_lst  []int      `bson:"counts__from_origin_urls_lst"  json:"counts__from_origin_urls_lst"`
	A_hrefs__from_origin_urls_lst [][]string `bson:"a_hrefs__from_origin_urls_lst" json:"a_hrefs__from_origin_urls_lst"`
}

type Stat__links_in_day struct {
	Total_count_int           int `bson:"total_count_int"           json:"total_count_int"`
	Valid_for_crawl_total_int int `bson:"valid_for_crawl_total_int" json:"valid_for_crawl_total_int"`
	Fetched_total_int         int `bson:"fetched_total_int"         json:"fetched_total_int"`
}
//-------------------------------------------------
func stats__new_links_by_day(p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_stats__links.stats__new_links_by_day()")

	type Minimal_link struct {
		Creation_unix_time_f float64 `bson:"creation_unix_time_f"`
		Valid_for_crawl_bool bool    `bson:"valid_for_crawl_bool"`
		Fetched_bool         bool    `bson:"fetched_bool"`
	}
	pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
		bson.M{"$match":bson.M{
				"t":"crawler_page_outgoing_link",
			},
		},
		bson.M{"$project":bson.M{
				//"id_str"               :true,
				"creation_unix_time_f" :true,
				//"cycle_run_id_str"     :true,
				//"domain_str"           :true,
				//"a_href_str"           :true, //actual link from the html <a> page ('href' parameter)
				//"origin_url_str"       :true, //page url from whos html this element was extracted
				"valid_for_crawl_bool" :true,
				"fetched_bool"         :true,
				//"images_processed_bool":true,
			},
		},
		bson.M{"$sort":bson.M{
				"creation_unix_time_f":-1,
			},
		},
	})

	results_lst := []Minimal_link{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to run an aggregation pipeline to get new links by day",
			"mongodb_aggregation_error",
			nil,err,"gf_crawl_stats",p_runtime_sys)
		return nil,gf_err
	}

	//--------------------
	//AGGREGATE DAY COUNTS - app-layer DB join
	new_links_counts_map := map[int]*Stat__links_in_day{}
	keys_lst             := []int{}
	for _,l := range results_lst {

		tm                := time.Unix(int64(l.Creation_unix_time_f),0)
		year_day_id_int,_ := strconv.Atoi(fmt.Sprintf("%d%d",tm.Year(),tm.YearDay()))

		//--------------
		var stat_r *Stat__links_in_day
		if stat,ok := new_links_counts_map[year_day_id_int]; ok {
			stat_r = stat
		} else {
			//-----------------
			//CREATE_NEW
			stat                                 := &Stat__links_in_day{}
			new_links_counts_map[year_day_id_int] = stat
			stat_r                                = stat

			keys_lst = append(keys_lst,year_day_id_int)
			//-----------------
		}
		//--------------

		stat_r.Total_count_int = stat_r.Total_count_int+1

		if l.Valid_for_crawl_bool {
			stat_r.Valid_for_crawl_total_int = stat_r.Valid_for_crawl_total_int + 1
		}

		if l.Fetched_bool {
			stat_r.Fetched_total_int = stat_r.Fetched_total_int + 1
		}
	}
	//--------------------

	sort.Ints(keys_lst)

	new_links_counts__sorted_lst := []*Stat__links_in_day{}
	for _,k := range keys_lst {

		stat                        := new_links_counts_map[k]
		new_links_counts__sorted_lst = append(new_links_counts__sorted_lst,stat)
	}

	fmt.Println("DONE SORTING >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(len(new_links_counts__sorted_lst))
	for _,a := range new_links_counts__sorted_lst {
		fmt.Println(a)
	}
	
	data_map := map[string]interface{}{
		"new_links_per_day_lst":new_links_counts__sorted_lst,
	}
	return data_map,nil
}
//-------------------------------------------------
func stats__unresolved_links(p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_stats__links.stats__unresolved_links()")

	pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
		bson.M{"$match":bson.M{
				"t"                   :"crawler_page_outgoing_link",
				"valid_for_crawl_bool":true,
				"fetched_bool"        :false,
			},
		},

		bson.M{"$sort":bson.M{
				"creation_unix_time_f":-1,
			},
		},

		bson.M{"$group":bson.M{
				"_id"        :bson.M{"origin_domain_str":"$origin_domain_str","origin_url_str":"$origin_url_str",},
				"count_int"  :bson.M{"$sum" :1},
				"a_hrefs_lst":bson.M{"$push":"$a_href_str",},
			},
		},

		bson.M{"$group":bson.M{
				"_id"                          :"$_id.origin_domain_str",
				"origin_urls_lst"              :bson.M{"$push":"$_id.origin_url_str"},
				"counts__from_origin_urls_lst" :bson.M{"$push":"$count_int"},
				"a_hrefs__from_origin_urls_lst":bson.M{"$push":"$a_hrefs_lst"},
			},
		},
	})

	results_lst := []Stat__unresolved_links{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to run an aggregation pipeline to unresolved links",
			"mongodb_aggregation_error",
			nil,err,"gf_crawl_stats",p_runtime_sys)
		return nil,gf_err
	}

	data_map := map[string]interface{}{
		"unresolved_links_lst":results_lst,
	}
	return data_map,nil
}
//-------------------------------------------------
func stats__crawled_links_domains(p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_stats__links.stats__crawled_links_domains()")

	pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
		bson.M{"$match":bson.M{
				"t":"crawler_page_outgoing_link",
			},
		},

		bson.M{"$project":bson.M{
				"id_str"               :true,
				"creation_unix_time_f" :true,
				"cycle_run_id_str"     :true,
				"domain_str"           :true,
				"a_href_str"           :true, //actual link from the html <a> page ('href' parameter)
				"origin_url_str"       :true, //page url from whos html this element was extracted
				"valid_for_crawl_bool" :true,
				"fetched_bool"         :true,
				"images_processed_bool":true,
			},
		},

		bson.M{"$group":bson.M{
				"_id"                    :"$domain_str",
				"links_count_int"        :bson.M{"$sum"     :1},
				"creation_unix_times_lst":bson.M{"$push"    :"$creation_unix_time_f"},
				"a_href_lst"             :bson.M{"$push"    :"$a_href_str"},
				"origin_urls_lst"        :bson.M{"$addToSet":"$origin_url_str"},
				"valid_for_crawl_lst"    :bson.M{"$push"    :"$valid_for_crawl_bool"},  //if the link is to be crawled/followed, or should be ignored
				"fetched_lst"            :bson.M{"$push"    :"$fetched_bool"},          //if the link's HTML was downloaded
				"images_processed_lst"   :bson.M{"$push"    :"$images_processed_bool"}, //if the images of this links HTML page were downloaded/processed
			},
		},

		bson.M{"$sort":bson.M{
				"links_count_int":-1,
			},
		},
	})

	results_lst := []Stat__crawled_links_domain{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to run an aggregation pipeline to get crawled links domains",
			"mongodb_aggregation_error",
			nil,err,"gf_crawl_stats",p_runtime_sys)
		return nil,gf_err
	}

	data_map := map[string]interface{}{
		"crawled_links_domains_lst":results_lst,
	}
	return data_map,nil
}