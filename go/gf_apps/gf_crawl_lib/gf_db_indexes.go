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

package gf_crawl_lib

import (
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)

//--------------------------------------------------
func db_index__init(p_runtime *gf_crawl_core.GFcrawlerRuntime,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	indexes_keys_lst := [][]string{
		[]string{"t", }, // all stat queries first match on "t"
		[]string{"t", "hash_str"},
	}

	_, gf_err := gf_core.Mongo__ensure_index(indexes_keys_lst, "gf_crawl", p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//-------------
	// LINK_INDEXES
	gf_err = gf_crawl_core.Link__db_index__init(p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	//-------------
	return nil
}