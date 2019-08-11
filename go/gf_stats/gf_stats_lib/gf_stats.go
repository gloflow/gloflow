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

package gf_stats_lib

import (
	"fmt"
	"time"
	"net/http"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
type Stat_query_run__extern_result struct {
	Query_run_id_str   string                 `json:"query_run_id_str"`
	Stat_name_str      string                 `json:"stat_name_str"`
	Start_time__unix_f float64                `json:"start_time__unix_f"`
	End_time__unix_f   float64                `json:"end_time__unix_f"`
	Result_data_map    map[string]interface{} `json:"result_data_map"`
}

type Stat_query_run struct {
	Id                 bson.ObjectId          `bson:"_id,omitempty"`
	Id_str             string                 `bson:"id_str"`
	T_str              string                 `bson:"t_str"` //"stat_query_run",
	Stat_name_str      string                 `bson:"stat_name_str"`
	Start_time__unix_f float64                `bson:"start_time__unix_f"`
	End_time__unix_f   float64                `bson:"end_time__unix_f"`
	Result_data_map    map[string]interface{} `bson:"result_data_map"`
}

//-------------------------------------------------
func Init(p_stats_url_base_str string,
	p_py_stats_dir_path_str       string,
	p_stats_query_funs_groups_lst []map[string]func(*gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error),
	p_runtime_sys                 *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_stats.Init()")

	//----------------
	//BATCH__HANDLERS
	gf_err := batch__init_handlers(p_stats_url_base_str, p_py_stats_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//----------------
	//QUERY__HANDLERS

	//collect all query funs into a single map
	query_funs_map := map[string]func(*gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error){}

	for _, stats_query_funs_map := range p_stats_query_funs_groups_lst {

		for stat_name_str,query_fun := range stats_query_funs_map {
			if _, ok := query_funs_map[stat_name_str]; !ok {
				query_funs_map[stat_name_str] = query_fun
			} else {
				//panicking here since this is only run on code initialization, and is a 
				//development time error (not an expected error)
				panic("there is a duplicate stat name in several query_funs_groups")
			}
		} 
	}
	query__init_handlers(p_stats_url_base_str, query_funs_map, p_runtime_sys)
	//----------------

	return nil
}

//-------------------------------------------------
func query__init_handlers(p_stats_url_base_str string,
	p_stats_query_funs_map map[string]func(*gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error),
	p_runtime_sys          *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_stats.query__init_handlers()")

	url_str := p_stats_url_base_str+"/query"
	http.HandleFunc(url_str, func(p_resp http.ResponseWriter, p_req *http.Request) {

		p_runtime_sys.Log_fun("INFO", fmt.Sprintf("INCOMING HTTP REQUEST -- %s ----------",p_stats_url_base_str))
		if p_req.Method == "POST" {
			
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			//INPUT
			i, gf_err := gf_rpc_lib.Get_http_input(url_str, p_resp, p_req, p_runtime_sys)
			if gf_err != nil {
				return
			}

			stat_name_str := i["stat_name_str"].(string)
			
			//--------------------------
			//RUN_QUERY_FUNCTION
			
			query_fun_result, gf_err := query__run_fun(stat_name_str, p_stats_query_funs_map, p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler(url_str, "stat run failed", gf_err, p_resp, p_runtime_sys)
				return
			}
			gf_rpc_lib.Http_Respond(query_fun_result, "OK", p_resp, p_runtime_sys)
			//--------------------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run(url_str, start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	});
}

//-------------------------------------------------
func query__run_fun(p_stat_name_str string,
	p_stats_query_funs_map map[string](func(*gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error)),
	p_runtime_sys          *gf_core.Runtime_sys) (*Stat_query_run__extern_result, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_stats.query__run_fun()")

	if stat_fun, ok := p_stats_query_funs_map[p_stat_name_str]; ok {

		start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
		run_id_str         := fmt.Sprintf("%f_%s", start_time__unix_f, p_stat_name_str)

		result_data_map, gf_err := stat_fun(p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0


		gf_err = Stat_run__create(p_stat_name_str, result_data_map, start_time__unix_f, end_time__unix_f, p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		stat_result := &Stat_query_run__extern_result{
			Query_run_id_str:   run_id_str,
			Stat_name_str:      p_stat_name_str,
			Start_time__unix_f: start_time__unix_f,
			End_time__unix_f:   end_time__unix_f,
			Result_data_map:    result_data_map,
		}

		return stat_result,nil	
	} else {
		gf_err := gf_core.Error__create("failed to get random img range from the DB",
			"verify__invalid_key_value_error",
			map[string]interface{}{"stat_name_str": p_stat_name_str,},
			nil, "gf_stats_lib", p_runtime_sys)
		return nil, gf_err
	}

	return nil, nil
}

//-------------------------------------------------
func Stat_run__create(p_stat_name_str string,
	p_results_data_lst   map[string]interface{},
	p_start_time__unix_f float64,
	p_end_time__unix_f   float64,
	p_runtime_sys        *gf_core.Runtime_sys) *gf_core.Gf_error {

	id_str := fmt.Sprintf("stat_query_run:%f", float64(time.Now().UnixNano())/1000000000.0)
	run    := &Stat_query_run{
		Id_str:             id_str,
		T_str:              "stat_query_run",
		Stat_name_str:      p_stat_name_str,
		Start_time__unix_f: p_start_time__unix_f,
		End_time__unix_f:   p_end_time__unix_f,
		Result_data_map:    p_results_data_lst,
	}

	err := p_runtime_sys.Mongodb_coll.Insert(run)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to persist a stat_run",
			"mongodb_insert_error",
			map[string]interface{}{"stat_name_str": p_stat_name_str,},
			err, "gf_stats_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}