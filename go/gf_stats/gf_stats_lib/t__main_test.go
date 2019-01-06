package gf_stats_lib

import (
	"fmt"
	"testing"
	"strings"
	"path/filepath"
	"gf_core"
)

//---------------------------------------------------
func Test__main(p_test *testing.T) {

	//-----------------
	test__mongodb_host_str    := "127.0.0.1"
	test__mongodb_db_name_str := "test_db"

	log_fun  := gf_core.Init_log_fun()
	mongo_db := gf_core.Mongo__connect(test__mongodb_host_str,
							test__mongodb_db_name_str,
							log_fun )
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