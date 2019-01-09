package gf_crawl_stats

import (
	"gf_core"
)
//-------------------------------------------------
func Get_query_funs(p_runtime_sys *gf_core.Runtime_sys) map[string]func(*gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_stats.Init()")

	stats_funs_map := map[string]func(*gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
		"crawler_fetches_by_url": stats__crawler_fetches_by_url,
		"crawler_fetches_by_days":stats__crawler_fetches_by_days,
		"crawled_links_domains":  stats__crawled_links_domains,
		"crawled_images_domains": stats__crawled_images_domains,
		"gifs":                   stats__gifs,
		"gifs_by_days":           stats__gifs_by_days,
		"errors":                 stats__errors,
		"unresolved_links":       stats__unresolved_links,
		"new_links_by_day":       stats__new_links_by_day,
	}
	return stats_funs_map
}