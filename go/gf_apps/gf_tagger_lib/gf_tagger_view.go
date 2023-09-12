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
	"io"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------

func renderObjectsWithTag(pTagStr string,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pPageIndexInt         int,
	pPageSizeInt          int,
	pResp                 io.Writer,
	pRuntimeSys           *gf_core.RuntimeSys) *gf_core.GFerror {

	//-----------------------------
	/*
	FIX!! - SCALABILITY!! - get tag info on "image" and "post" types is a very long
		operation, and should be done in some more efficient way,
		as in with a mongodb aggregation pipeline.
	*/
	objectsInfosLst, gfErr := getObjectsWithTag(pTagStr,
		"post", // p_objectTypeStr
		pPageIndexInt,
		pPageSizeInt,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	postsWithTagLst := []map[string]interface{}{}
	for _, objectInfoMap := range objectsInfosLst {

		//----------------
		var postThumbnailURLstr string
		thumbSmallStr := objectInfoMap["thumbnail_small_url_str"].(string)

		if thumbSmallStr == "" {
			
			// FIX!! - use some user-configurable value that is configured at startup
			// IMPORTANT!! - some "thumbnail_small_url_str" are blank strings (""),
			errorImageURLstr := "https://gloflow.com/images/d/gf_landing_page_logo.png"

			postThumbnailURLstr = errorImageURLstr
		} else {
			postThumbnailURLstr = thumbSmallStr
		}

		//----------------
		postInfoMap := map[string]interface{}{
			"post_title_str":         objectInfoMap["title_str"].(string),
			"post_tags_lst":          objectInfoMap["tags_lst"].([]string),
			"post_url_str":           objectInfoMap["url_str"].(string),
			"post_thumbnail_url_str": postThumbnailURLstr,
		}

		postsWithTagLst = append(postsWithTagLst, postInfoMap)
	}
	//-----------------------------


	type templatesData struct {
		Tag_str                string
		Posts_with_tag_num_int int64
		Images_with_tag_int    int64
		Posts_with_tag_lst     []map[string]interface{}
		Sys_release_info       gf_core.SysReleaseInfo
		Is_subtmpl_def         func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	objectTypeStr := "post"
	postsWithTagCountInt, gfErr := dbMongoGetObjectsWithTagCount(pTagStr, objectTypeStr, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)

	err := pTemplate.Execute(pResp,
		templatesData{
			Tag_str:                pTagStr,
			Posts_with_tag_num_int: postsWithTagCountInt,
			Images_with_tag_int:    0, // FIX!! - image tagging is now implemented, and so counting images with tag occurance should be done ASAP. 
			Posts_with_tag_lst:     postsWithTagLst,
			Sys_release_info:       sysReleaseInfo,
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
		return gfErr
	}

	return nil
}	