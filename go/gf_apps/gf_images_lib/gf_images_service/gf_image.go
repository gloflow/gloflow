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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------
func ImgGet(pImageIDstr gf_images_core.GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_images_core.GF_image_export, bool, *gf_core.GFerror) {

	// DB_EXISTS
	exists_bool, gfErr := gf_images_core.DBimageExists(pImageIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, false, gfErr
	}

	if exists_bool {

		// DB_GET
		gf_img, gfErr := gf_images_core.DBgetImage(pImageIDstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, false, gfErr
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
func TagsAddToImage(p_image *gf_images_core.GF_image,
	p_tags_lst    []string,
	pRuntimeSys *gf_core.RuntimeSys) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_image.TagsAddToImage()")
	
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