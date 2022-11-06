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
	"io/ioutil"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func batch__init_handlers(p_stats_url_base_str string,
	p_py_stats_dir_path_str string,
	p_runtime_sys           *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_stats_batch.batch__init_handlers()")

	stats_list_lst, gf_err := batch__get_stats_list(p_py_stats_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	url_str := fmt.Sprintf("%s/batch/list", p_stats_url_base_str)
	http.HandleFunc(url_str, func(p_resp http.ResponseWriter, p_req *http.Request) {

		p_runtime_sys.LogFun("INFO",fmt.Sprintf("INCOMING HTTP REQUEST -- %s ----------", p_stats_url_base_str))
		if p_req.Method == "GET" {

			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			data_map := map[string]interface{}{
				"stats_list_lst": stats_list_lst,
			}

			gf_rpc_lib.HTTPrespond(data_map, "OK", p_resp, p_runtime_sys)
			//--------------------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run(url_str, start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})

	return nil
}

//-------------------------------------------------
func batch__get_stats_list(p_py_stats_dir_path_str string,
	p_runtime_sys *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_stats_batch.batch__get_stats_list()")

	files_lst, err := ioutil.ReadDir(p_py_stats_dir_path_str)
	if err != nil {
		gf_err := gf_core.ErrorCreate("failed to list py_stats dir in order to get a list of batch py_stats",
			"dir_list_error",
			map[string]interface{}{"py_stats_dir_path_str": p_py_stats_dir_path_str,},
			err, "gf_stats_lib", p_runtime_sys)
		return nil, gf_err
	}

	py_stats__names_lst := []string{}
	for _, file := range files_lst {
		
		file_basename_str := file.Name()

		if strings.HasSuffix(file_basename_str, ".py") {
			py_stat__name_str  := strings.TrimSuffix(file_basename_str, ".py")
			py_stats__names_lst = append(py_stats__names_lst, py_stat__name_str)
		}
	}

	return py_stats__names_lst, nil
}