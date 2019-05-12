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
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func Get_query_funs(p_runtime_sys *gf_core.Runtime_sys) map[string]func(*gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_stats.Init()")

	stats_funs_map := map[string]func(*gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {
		"crawler_fetches_by_url":  stats__crawler_fetches_by_url,
		"crawler_fetches_by_days": stats__crawler_fetches_by_days,
		"crawled_links_domains":   stats__crawled_links_domains,
		"crawled_images_domains":  stats__crawled_images_domains,
		"gifs":                    stats__gifs,
		"gifs_by_days":            stats__gifs_by_days,
		"errors":                  stats__errors,
		"unresolved_links":        stats__unresolved_links,
		"new_links_by_day":        stats__new_links_by_day,
	}
	return stats_funs_map
}