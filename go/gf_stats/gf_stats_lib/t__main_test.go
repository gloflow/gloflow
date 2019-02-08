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
	"testing"
	"strings"
	"path/filepath"
	"github.com/gloflow/gloflow/go/gf_core"
)
//---------------------------------------------------
func Test__main(p_test *testing.T) {

	//-----------------
	test__mongodb_host_str    := "127.0.0.1"
	test__mongodb_db_name_str := "test_db"

	log_fun      := gf_core.Init_log_fun()
	mongo_db     := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongo_db.C("data_symphony")

	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str:"gf_stats_lib_test",
		Log_fun:         log_fun,
		Mongodb_coll:    mongodb_coll,
	}
	//-----------------

	test__py_stats_dir_path_str,err := filepath.Abs("../../apps/gf_crawl_lib/py/stats")
	if err != nil {
		p_test.Errorf("failed to get absolute path of test__py_stats_dir - %s",test__py_stats_dir_path_str)
	}


	fmt.Println("TEST ---")
	fmt.Println(test__py_stats_dir_path_str)


	py_stats__names_lst,gf_err := batch__get_stats_list(test__py_stats_dir_path_str,runtime_sys)
	if gf_err != nil {
		p_test.Errorf("failed to list py_stats files in py_stats_dir - %s",test__py_stats_dir_path_str)
	}

	if len(py_stats__names_lst) == 0 {
		p_test.Errorf("no py_stats found in py_stats_dir - %s",test__py_stats_dir_path_str)
	}

	for _,py_stat_name_str := range py_stats__names_lst {

		fmt.Println("py_stat_name_str - "+py_stat_name_str)
		
		if strings.HasSuffix(py_stat_name_str,".py") {
			p_test.Errorf("list py_stats file still has a .py extension - %s",py_stat_name_str)
		}
	}
}