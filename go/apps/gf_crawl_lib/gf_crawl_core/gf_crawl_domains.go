package gf_crawl_core

import (
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
func get_domains_blacklist(p_runtime_sys *gf_core.Runtime_sys) map[string]bool {

	domains_map := map[string]bool{
		"facebook.com"     :false,
		"l.facebook.com"   :false,
		"twitter.com"      :false,
		"apple.com"        :false,
		"apple.news"       :false,
		"tw.appstore.com"  :false,
		"tw.itunes.com"    :false,
		"itunes.apple.com" :false,
		"vimeo.com"        :false,
		"cloud.feedly.com" :false,
		"mediatemple.net"  :false,
		"pinterest.com"    :false,
		"youtube.com"      :false,
		"ffffound.com"     :false,
		
	}
	return domains_map
}