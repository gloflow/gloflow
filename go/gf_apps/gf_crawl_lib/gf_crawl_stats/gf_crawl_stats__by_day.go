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
	"fmt"
	"time"
	"strconv"
	"sort"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type Gf_stats__objs_by_days struct {
	Obj_type_str                     string                           `json:"obj_type_str"                     bson:"obj_type_str"`
	Counts_by_day__sorted_lst        []int                            `json:"counts_by_day__sorted_lst"        bson:"counts_by_day__sorted_lst"`        //global count of fetches per day
	Domain_counts_by_day__sorted_lst []*Gf_domain_counts_for_all_days `json:"domain_counts_by_day__sorted_lst" bson:"domain_counts_by_day__sorted_lst"` //counts of fetches per domain per day
}

type Gf_domain_counts_for_all_days struct {
	Domain_str      string `json:"domain_str"      bson:"domain_str"`
	Total_count_int int    `json:"total_count_int" bson:"total_count_int"`
	Days_counts_lst []int  `json:"days_counts_lst" bson:"days_counts_lst"`
}

//-------------------------------------------------
func stats__objs_by_days(p_match_query_map map[string]interface{},
	p_obj_type_str string,
	p_runtime_sys  *gf_core.RuntimeSys) (*Gf_stats__objs_by_days, *gf_core.Gf_error) {
	p_runtime_sys.LogFun("FUN_ENTER","gf_crawl_stats__by_day.stats__objs_by_days()")

	type Domain_objs__stat struct {
		Domain_str         string    `bson:"_id"`
		Count_int          int       `bson:"count_int"`
		Creation_times_lst []float64 `bson:"creation_times_lst"`
	}

	match_query := bson.M{
		"t": p_obj_type_str,
	}
	for k,v := range p_match_query_map {
		match_query[k] = v
	}

	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{"$match", match_query},
		},
		{
			{"$project", bson.D{
				{"creation_unix_time_f", true},
				{"domain_str",           true},
			}},
		},
		{
			{"$sort", bson.D{{"creation_unit_time_f", -1},}},
		},
		{
			{"$group", bson.D{
				{"_id",                "$domain_str"},
				{"count_int",          bson.M{"$sum":  1}},
				{"creation_times_lst", bson.M{"$push": "$creation_unix_time_f"}},
			}},
		},
		{
			{"$sort", bson.D{{"fetches_count_int", -1},}},
		},
	}

	/*pipe := p_runtime_sys.Mongo_db.Collection("gf_crawl").Pipe([]bson.M{
		// bson.M{"$match":bson.M{
		// 		"t":p_obj_type_str, //"crawler_url_fetch",
		// 	},
		// },
		bson.M{"$match": match_query},
		bson.M{"$project": bson.M{
				"creation_unix_time_f": true,
				"domain_str":           true,
			},
		},
		bson.M{"$sort":bson.M{
				"creation_unix_time_f":-1,
			},
		},
		bson.M{"$group":bson.M{
				"_id":                "$domain_str",
				"count_int":          bson.M{"$sum" :1},
				"creation_times_lst": bson.M{"$push":"$creation_unix_time_f"},
			},
		},
		bson.M{"$sort":bson.M{
				"fetches_count_int": -1,
			},
		},
	})*/

	cursor, err := p_runtime_sys.Mongo_db.Collection("gf_crawl").Aggregate(ctx, pipeline)
	if err != nil {

		gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to count objects by days",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_crawl_stats", p_runtime_sys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)

	/*results_lst := []Domain_objs__stat{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to count objects by days",
			"mongodb_aggregation_error",
			map[string]interface{}{"obj_type_str": p_obj_type_str,},
			err, "gf_crawl_stats", p_runtime_sys)
		return nil, gf_err
	}*/

	results_lst := []Domain_objs__stat{}
	for cursor.Next(ctx) {

		var r Domain_objs__stat
		err := cursor.Decode(&r)
		if err != nil {
			gf_err := gf_core.Mongo__handle_error("failed to run an aggregation pipeline to count objects by days",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_crawl_stats", p_runtime_sys)
			return nil, gf_err
		}
	
		results_lst = append(results_lst, r)
	}

	//--------------------
	// AGGREGATE DAY COUNTS

	type stat__objs_in_day struct {
		Year_int                    int            `bson:"year_int"                    json:"year_int"`
		Day_int                     int            `bson:"day_int"                     json:"day_int"`
		Total_count_int             int            `bson:"total_count_int"             json:"total_count_int"`
		Total_count__per_domain_map map[string]int `bson:"total_count__per_domain_map" json:"total_count__per_domain_map"`
	}

	days_stats_map               := map[int]*stat__objs_in_day{}
	days_keys_lst                := []int{}
	all_domains_total_counts_map := map[string]int{} // total count of all days 

	for _, domain := range results_lst {

		all_domains_total_counts_map[domain.Domain_str] = domain.Count_int


		// IMPORTANT!! - process each creation timestamp for a fetch
		for _, t := range domain.Creation_times_lst { 
			tm                 := time.Unix(int64(t), 0)
			year_day_id_int, _ := strconv.Atoi(fmt.Sprintf("%d%d", tm.Year(), tm.YearDay()))

			//--------------
			var day_stat *stat__objs_in_day

			if cached_stat, ok := days_stats_map[year_day_id_int]; ok {
				day_stat = cached_stat
			} else {
				//-----------------
				// CREATE_NEW
				new_stat := &stat__objs_in_day{
					Year_int:                    tm.Year(),
					Day_int:                     tm.YearDay(),
					Total_count__per_domain_map: map[string]int{},	
				}
				days_stats_map[year_day_id_int] = new_stat
				day_stat                        = new_stat

				days_keys_lst = append(days_keys_lst, year_day_id_int)

				//-----------------
			}
			//--------------

			day_stat.Total_count_int++
			day_stat.Total_count__per_domain_map[domain.Domain_str]++		
		}
	}

	//--------------------
	// SORT DAY STATS BY DAY
	
	sort.Ints(days_keys_lst)

	stats__sorted_by_day_lst := []*stat__objs_in_day{}
	for _,k := range days_keys_lst {
		day_stat                := days_stats_map[k]
		stats__sorted_by_day_lst = append(stats__sorted_by_day_lst, day_stat)
	}

	//------------------
	// ZERO_OUT_BLANK_VALUES

	// since in a particular day not all possible domain fetches are done, some of the domains counts
	// for some day are 0, and this needs to be set manually, by checking if counts exist for that day.
	// if no then set it to 0.
	for _, day_stat := range stats__sorted_by_day_lst {

		for domain_str, _ := range all_domains_total_counts_map {

			if _, ok := day_stat.Total_count__per_domain_map[domain_str]; !ok {
				day_stat.Total_count__per_domain_map[domain_str] = 0
			}
		}
	}
	//----------------------
	// ACCUMULATE DOMAINS COUNTS IN COLUMNS

	// gets for each domain a list of fetches number for each day
	// this is formated this way for easy feeding of columns of data (per domain)
	// to visualization routines, without need to do a bunch of app joins client side

	domains__counts_for_all_days_lst := []*Gf_domain_counts_for_all_days{} // map[string][]int{}

	for domain_str, _ := range all_domains_total_counts_map {

		domain_days_count_lst := []int{}
		for _, day_stat := range stats__sorted_by_day_lst {


			domain_day_count_int := day_stat.Total_count__per_domain_map[domain_str]
			domain_days_count_lst = append(domain_days_count_lst, domain_day_count_int)
		}

		//domains__counts_for_all_days_map[domain_str] = domain_days_count_lst

		d := &Gf_domain_counts_for_all_days{
			Domain_str:      domain_str,
			Total_count_int: all_domains_total_counts_map[domain_str],
			Days_counts_lst: domain_days_count_lst,
		}
		domains__counts_for_all_days_lst = append(domains__counts_for_all_days_lst, d)
	}

	//----------------------
	// SORT
	sort__domains_counts(domains__counts_for_all_days_lst, p_runtime_sys)

	//------------------
	// ACCUMULATE LIST OF GLOBAL FETCH COUNTS PER DAY
	total_counts_by_day__sorted_lst := []int{}
	for _, day_stat := range stats__sorted_by_day_lst {
		total_counts_by_day__sorted_lst = append(total_counts_by_day__sorted_lst, day_stat.Total_count_int)
	}

	//------------------
	stats := &Gf_stats__objs_by_days{
		Obj_type_str:                     p_obj_type_str,
		Counts_by_day__sorted_lst:        total_counts_by_day__sorted_lst,
		Domain_counts_by_day__sorted_lst: domains__counts_for_all_days_lst,
	}

	//------------------
	return stats, nil
}

//-------------------------------------------------
// SORT
type domains_counts []*Gf_domain_counts_for_all_days
func (d_lst domains_counts) Len() int {
	return len(d_lst)
}
func (d_lst domains_counts) Swap(i, j int) {
	d_lst[i],d_lst[j] = d_lst[j],d_lst[i]
}
func (d_lst domains_counts) Less(i, j int) bool {
	return d_lst[i].Total_count_int > d_lst[j].Total_count_int
}

func sort__domains_counts(p_domains__counts_for_all_days_lst []*Gf_domain_counts_for_all_days, p_runtime_sys *gf_core.RuntimeSys) {
	/*//------------------
	//SORT
	func (s domains_counts) Len() int {
		return len(s)
	}
	func (s domains_counts) Swap(i, j int) {
		s[i], s[j] = s[j], s[i]
	}
	func (s domains_counts) Less(i, j int) bool {
		//return len(s[i]) < len(s[j])

		a:=s[i].(*domain_counts_for_all_days)
		b:=s[i].(*domain_counts_for_all_days)
		return a.Domain_total_fetches_count_int < b.Domain_total_fetches_count_int
	}*/
	sort.Sort(domains_counts(p_domains__counts_for_all_days_lst))
	return
}