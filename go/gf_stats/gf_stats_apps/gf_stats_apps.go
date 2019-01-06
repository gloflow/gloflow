package gf_stats_apps

import (
	"gf_core"
	"gf_stats/gf_stats_lib"
	"apps/gf_images_lib/gf_images_stats"
	"apps/gf_crawl_lib/gf_crawl_stats"
)
//-------------------------------------------------
func Init(p_stats_url_base_str string,
	p_py_stats_dir_path_str string,
	p_runtime_sys           *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_stats_apps.Init()")

	images_stats__query_funs_map := gf_images_stats.Get_query_funs(p_runtime_sys)
	crawl_stats__query_funs_map  := gf_crawl_stats.Get_query_funs(p_runtime_sys)

	stats_query_funs_groups_lst := []map[string]func(*gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error){
		images_stats__query_funs_map,
		crawl_stats__query_funs_map,
	}

	gf_err := gf_stats_lib.Init(p_stats_url_base_str,
					p_py_stats_dir_path_str,
					stats_query_funs_groups_lst,
					p_runtime_sys)
	return gf_err
}