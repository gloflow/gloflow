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
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------

func ImageGet(pImageIDstr gf_images_core.GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_images_core.GFimageExport, bool, *gf_core.GFerror) {


	// DB_EXISTS
	existsBool, gfErr := gf_images_core.DBimageExists(pImageIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, false, gfErr
	}

	/*
	existsBool, gfErr := gf_images_core.DBmongoImageExists(pImageIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, false, gfErr
	}
	*/

	if existsBool {

		// DB_GET
		image, gfErr := gf_images_core.DBgetImage(pImageIDstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, false, gfErr
		}
		
		/*
		image, gfErr := gf_images_core.DBmongoGetImage(pImageIDstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, false, gfErr
		}
		*/
		
		resolvedUserNameStr := gf_identity_core.ResolveUserName(image.UserID, pCtx, pRuntimeSys)

		imageExport := &gf_images_core.GFimageExport{
			Creation_unix_time_f:  image.Creation_unix_time_f,
			UserNameStr:           resolvedUserNameStr,
			Title_str:             image.TitleStr,
			Flows_names_lst:       image.FlowsNamesLst,
			ThumbnailSmallURLstr:  image.ThumbnailSmallURLstr,
			ThumbnailMediumURLstr: image.ThumbnailMediumURLstr,
			ThumbnailLargeURLstr:  image.ThumbnailLargeURLstr,
			Format_str:            image.Format_str,
			Tags_lst:              image.TagsLst,
		}
		return imageExport, true, nil
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