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

package gf_images_service

import (
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
func DB_index__init(p_runtime_sys *gf_core.RuntimeSys) *gf_core.Gf_error {

	indexes_keys_lst := [][]string{
		[]string{"t", }, // all stat queries first match on "t"

		// QUERIES - flows_db__images_exist() issues these queries
		// DEPRECATED!! - flow_name_str field is deprecated in favor of flows_names_lst, 
		//                but some of the old img records still use it and havent been migrated yet. 
		//                so we're creating this index until migration is complete, and then 
		//                it should be removed.
		[]string{"t", "flow_names_lst", "origin_url_str"},

		// QUERIES - flows_db__images_exist() issues these queries
		[]string{"t", "flows_names_lst", "origin_url_str"},
	}
	
	_, gf_err := gf_core.Mongo__ensure_index(indexes_keys_lst, "data_symphony", p_runtime_sys)
	return gf_err
}