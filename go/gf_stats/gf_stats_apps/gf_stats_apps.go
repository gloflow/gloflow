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

package gf_stats_apps

import (
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_stats/gf_stats_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_stats"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_stats"
)

//-------------------------------------------------
func Init(p_stats_url_base_str string,
	p_py_stats_dir_path_str string,
	p_runtime_sys           *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_stats_apps.Init()")

	images_stats__query_funs_map := gf_images_stats.Get_query_funs(p_runtime_sys)
	crawl_stats__query_funs_map  := gf_crawl_stats.Get_query_funs(p_runtime_sys)

	stats_query_funs_groups_lst := []map[string]func(*gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error){
		images_stats__query_funs_map,
		crawl_stats__query_funs_map,
	}

	gf_err := gf_stats_lib.Init(p_stats_url_base_str,
		p_py_stats_dir_path_str,
		stats_query_funs_groups_lst,
		p_runtime_sys)
	return gf_err
}