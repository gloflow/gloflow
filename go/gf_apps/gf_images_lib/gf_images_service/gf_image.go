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
	pRuntimeSys *gf_core.RuntimeSys) (*gf_images_core.GFimageExport, bool, *gf_core.GFerror) {

	// DB_EXISTS
	existsBool, gfErr := gf_images_core.DBimageExists(pImageIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, false, gfErr
	}

	if existsBool {

		// DB_GET
		gfImage, gfErr := gf_images_core.DBgetImage(pImageIDstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, false, gfErr
		}

		gfImageExport := &gf_images_core.GFimageExport{
			Creation_unix_time_f:     gfImage.Creation_unix_time_f,
			Title_str:                gfImage.TitleStr,
			Flows_names_lst:          gfImage.FlowsNamesLst,
			Thumbnail_small_url_str:  gfImage.Thumbnail_small_url_str,
			Thumbnail_medium_url_str: gfImage.Thumbnail_medium_url_str,
			Thumbnail_large_url_str:  gfImage.Thumbnail_large_url_str,
			Format_str:               gfImage.Format_str,
			Tags_lst:                 gfImage.TagsLst,
		}
		return gfImageExport, true, nil
	} else {
		return nil, false, nil
	}

	return nil, false, nil
}

//---------------------------------------------------

func TagsAddToImage(pImage *gf_images_core.GFimage,
	pTagsLst    []string,
	pRuntimeSys *gf_core.RuntimeSys) {
	
	if len(pTagsLst) > 0 {

		//add all new tags with the current tags associated with an image,
		//with possible duplicates existing
		pImage.TagsLst = append(pImage.TagsLst, pTagsLst...)

		//-----------
		set := map[string]bool{}
		for _, t_str := range pImage.TagsLst {
			set[t_str] = true
		}

		//-----------
		listNoDuplicatesLst := []string{}
		for kStr, _ := range set {
			listNoDuplicatesLst = append(listNoDuplicatesLst, kStr)
		}

		//eliminate duplicates from the list
		pImage.TagsLst = listNoDuplicatesLst
		
		//-----------
	}
}