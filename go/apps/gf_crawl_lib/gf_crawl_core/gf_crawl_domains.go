/*
GloFlow media management/publishing system
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