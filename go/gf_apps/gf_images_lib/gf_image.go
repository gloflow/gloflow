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

package gf_images_lib 

import (
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)

//---------------------------------------------------
func Add_tags_to_image(p_image *gf_images_utils.Gf_image,
	p_tags_lst    []string,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_image.Add_tags_to_image()")
	
	if len(p_tags_lst) > 0 {

		//add all new tags with the current tags associated with an image,
		//with possible duplicates existing
		p_image.Tags_lst = append(p_image.Tags_lst, p_tags_lst...)

		//-----------
		set := map[string]bool{}
		for _,t_str := range p_image.Tags_lst {
			set[t_str]=true
		}
		//-----------
		list_no_duplicates_lst := []string{}
		for k_str,_ := range set {
			list_no_duplicates_lst = append(list_no_duplicates_lst, k_str)
		}

		//eliminate duplicates from the list
		p_image.Tags_lst = list_no_duplicates_lst
		//-----------
	}
}