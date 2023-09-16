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

package gf_tagger_lib

import (
	"context"
	"bytes"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------

func renderObjectsWithTag(pTagStr string,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pPageIndexInt         int,
	pPageSizeInt          int,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	//-----------------------------
	// IMAGES

	/*
	FIX!! - SCALABILITY!! - get tag info on "image" and "post" types is a very long
		operation, and should be done in some more efficient way,
		as in with a mongodb aggregation pipeline.
	*/
	imagesWithTagLst, gfErr := exportObjectsWithTag(pTagStr,
		"image", // p_objectTypeStr
		pPageIndexInt,
		pPageSizeInt,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	imagesUserNamesLst := resolveUserIDStoUserNames(imagesWithTagLst, pCtx, pRuntimeSys)

	for i, imageMap := range imagesWithTagLst {
		imageMap["owner_user_name_str"] = imagesUserNamesLst[i]
	}

	//-----------------------------
	// POSTS

	/*
	FIX!! - SCALABILITY!! - get tag info on "image" and "post" types is a very long
		operation, and should be done in some more efficient way,
		as in with a mongodb aggregation pipeline.
	*/
	postsInfosLst, gfErr := exportObjectsWithTag(pTagStr,
		"post", // p_objectTypeStr
		pPageIndexInt,
		pPageSizeInt,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	postsWithTagLst := []map[string]interface{}{}
	for _, postInfoMap := range postsInfosLst {

		//----------------
		var postThumbnailURLstr string
		thumbSmallStr := postInfoMap["thumbnail_small_url_str"].(string)

		if thumbSmallStr == "" {
			
			// FIX!! - use some user-configurable value that is configured at startup
			// IMPORTANT!! - some "thumbnail_small_url_str" are blank strings (""),
			errorImageURLstr := "https://gloflow.com/images/d/gf_landing_page_logo.png"

			postThumbnailURLstr = errorImageURLstr
		} else {
			postThumbnailURLstr = thumbSmallStr
		}

		//----------------
		postWithTagMap := map[string]interface{}{
			"post_title_str":         postInfoMap["title_str"].(string),
			"post_tags_lst":          postInfoMap["tags_lst"].([]string),
			"post_url_str":           postInfoMap["url_str"].(string),
			"post_thumbnail_url_str": postThumbnailURLstr,
		}

		postsWithTagLst = append(postsWithTagLst, postWithTagMap)
	}
	
	//-----------------------------


	type templatesData struct {
		TagStr                string
		ImagesWithTagCountInt int64
		PostsWithTagCountInt  int64
		
		ImagesWithTagLst      []map[string]interface{}
		PostsWithTagLst       []map[string]interface{}
		Sys_release_info      gf_core.SysReleaseInfo
		Is_subtmpl_def        func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	imagesWithTagCountInt, gfErr := dbMongoGetObjectsWithTagCount(pTagStr, "image", pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}
	postsWithTagCountInt, gfErr := dbMongoGetObjectsWithTagCount(pTagStr, "post", pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)

	buff := new(bytes.Buffer)
	err := pTemplate.Execute(buff,
		templatesData{
			TagStr:                pTagStr,
			ImagesWithTagCountInt: imagesWithTagCountInt,
			PostsWithTagCountInt:  postsWithTagCountInt,
			ImagesWithTagLst:      imagesWithTagLst,
			PostsWithTagLst:       postsWithTagLst,
			Sys_release_info:      sysReleaseInfo,
			
			//-------------------------------------------------
			// IS_SUBTEMPLATE_DEFINED
			Is_subtmpl_def: func(p_subtemplate_name_str string) bool {
				for _, n := range pSubtemplatesNamesLst {
					if n == p_subtemplate_name_str {
						return true
					}
				}
				return false
			},

			//-------------------------------------------------
		})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to render the objects_with_tag template",
			"template_render_error",
			map[string]interface{}{"tag_str": pTagStr,},
			err, "gf_tagger_lib", pRuntimeSys)
		return "", gfErr
	}

	templateRenderedStr := buff.String()
	return templateRenderedStr, nil
}	