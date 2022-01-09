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
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------
func Get_img(p_image_id_str gf_images_core.GF_image_id,
	p_runtime_sys *gf_core.Runtime_sys) (*gf_images_core.GF_image_export, bool, *gf_core.GF_error) {

	// DB_EXISTS
	exists_bool, gf_err := gf_images_core.DB__image_exists(p_image_id_str, p_runtime_sys)
	if gf_err != nil {
		return nil, false, gf_err
	}

	if exists_bool {

		// DB_GET
		gf_img, gf_err := gf_images_core.DB__get_image(p_image_id_str, p_runtime_sys)
		if gf_err != nil {
			return nil, false, gf_err
		}

		gf_img_export := &gf_images_core.GF_image_export{
			Creation_unix_time_f:     gf_img.Creation_unix_time_f,
			Title_str:                gf_img.Title_str,
			Flows_names_lst:          gf_img.Flows_names_lst,
			Thumbnail_small_url_str:  gf_img.Thumbnail_small_url_str,
			Thumbnail_medium_url_str: gf_img.Thumbnail_medium_url_str,
			Thumbnail_large_url_str:  gf_img.Thumbnail_large_url_str,
			Format_str:               gf_img.Format_str,
			Tags_lst:                 gf_img.Tags_lst,
		}
		return gf_img_export, true, nil
	} else {
		return nil, false, nil
	}

	return nil, false, nil
}

//---------------------------------------------------
func Add_tags_to_image(p_image *gf_images_core.GF_image,
	p_tags_lst    []string,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_image.Add_tags_to_image()")
	
	if len(p_tags_lst) > 0 {

		//add all new tags with the current tags associated with an image,
		//with possible duplicates existing
		p_image.Tags_lst = append(p_image.Tags_lst, p_tags_lst...)

		//-----------
		set := map[string]bool{}
		for _, t_str := range p_image.Tags_lst {
			set[t_str] = true
		}

		//-----------
		list_no_duplicates_lst := []string{}
		for k_str, _ := range set {
			list_no_duplicates_lst = append(list_no_duplicates_lst, k_str)
		}

		//eliminate duplicates from the list
		p_image.Tags_lst = list_no_duplicates_lst
		
		//-----------
	}
}