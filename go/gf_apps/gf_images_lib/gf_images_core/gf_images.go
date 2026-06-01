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

package gf_images_core

import (
	"time"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func ImageCreateNew(pImageInfo *GFimageNewInfo,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	image := &GFimage{
		IDstr:                          pImageInfo.IDstr,
		T_str:                          "img",
		Creation_unix_time_f:           creationUNIXtimeF,
		UserID:                         pImageInfo.UserID,
		ClientTypeStr:                  pImageInfo.Image_client_type_str,
		TitleStr:                       pImageInfo.Title_str,
		FlowsNamesLst:                  pImageInfo.Flows_names_lst,
		Origin_url_str:                 pImageInfo.Origin_url_str,
		Origin_page_url_str:            pImageInfo.Origin_page_url_str,
		Original_file_internal_uri_str: pImageInfo.Original_file_internal_uri_str,
		ThumbnailSmallURLstr:           pImageInfo.ThumbnailSmallURLstr,
		ThumbnailMediumURLstr:          pImageInfo.ThumbnailMediumURLstr,
		ThumbnailLargeURLstr:           pImageInfo.ThumbnailLargeURLstr,
		Format_str:                     pImageInfo.Format_str,
		Width_int:                      pImageInfo.Width_int,
		Height_int:                     pImageInfo.Height_int,

		TagsLst: []string{},
		MetaMap: pImageInfo.Meta_map,
	}

	//----------------------------------
	// DB PERSIST

	/*
	// MONGO
	gfErr := DBmongoPutImage(image, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	*/

	// SQL
	gfErr := DBsqlPutImage(image, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//----------------------------------

	return image, nil
}

//---------------------------------------------------
// DEPRECATED!! - use Image__create_new and its structured input

/*
func ImageCreate(pImageInfoMap map[string]interface{},
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {
	
	newImageInfoMap, gfErr := VerifyImageInfo(pImageInfoMap, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	title_str       := newImageInfoMap["title_str"].(string)
	flows_names_lst := newImageInfoMap["flows_names_lst"].([]string)
	gf_image_id_str := GFimageID(newImageInfoMap["id_str"].(string))

	image := &GFimage{
		IDstr:                          gf_image_id_str,
		T_str:                          "img",
		Creation_unix_time_f:           float64(time.Now().UnixNano())/1000000000.0,
		ClientTypeStr:                  newImageInfoMap["image_client_type_str"].(string),
		TitleStr:                       title_str,
		FlowsNamesLst:                  flows_names_lst,
		Origin_url_str:                 newImageInfoMap["origin_url_str"].(string),
		Origin_page_url_str:            newImageInfoMap["origin_page_url_str"].(string),
		Original_file_internal_uri_str: newImageInfoMap["original_file_internal_uri_str"].(string),
		Thumbnail_small_url_str:        newImageInfoMap["thumbnail_small_url_str"].(string),
		Thumbnail_medium_url_str:       newImageInfoMap["thumbnail_medium_url_str"].(string),
		Thumbnail_large_url_str:        newImageInfoMap["thumbnail_large_url_str"].(string),
		Format_str:                     newImageInfoMap["format_str"].(string),
	}
	
	//----------------------------------
	// DB PERSIST

	gfErr = DBputImage(image, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//----------------------------------

	return image, nil
}
*/