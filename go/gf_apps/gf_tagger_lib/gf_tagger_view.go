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
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_page_index_int         int,
	p_page_size_int          int,
	p_resp                   io.Writer,
	pRuntimeSys              *gf_core.RuntimeSys) *gf_core.GFerror {

	//-----------------------------
	// FIX!! - SCALABILITY!! - get tag info on "image" and "post" types is a very long
	//                         operation, and should be done in some more efficient way,
	//                         as in with a mongodb aggregation pipeline

	objectsInfosLst, gfErr := getObjectsWithTag(pTagStr,
		"post", // p_objectTypeStr
		p_page_index_int,
		p_page_size_int,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	posts_with_tag_lst := []map[string]interface{}{}
	for _, p_object_info_map := range objectsInfosLst {

		//----------------
		var post_thumbnail_url_str string
		thumb_small_str := p_object_info_map["thumbnail_small_url_str"].(string)

		if thumb_small_str == "" {
			
			// FIX!! - use some user-configurable value that is configured at startup
			// IMPORTANT!! - some "thumbnail_small_url_str" are blank strings (""),
			error_img_url_str     := "http://gloflow.com/images/d/gf_landing_page_logo.png"
			post_thumbnail_url_str = error_img_url_str
		} else {
			post_thumbnail_url_str = thumb_small_str
		}

		//----------------
		post_info_map := map[string]interface{}{
			"post_title_str":         p_object_info_map["title_str"].(string),
			"post_tags_lst":          p_object_info_map["tags_lst"].([]string),
			"post_url_str":           p_object_info_map["url_str"].(string),
			"post_thumbnail_url_str": post_thumbnail_url_str,
		}

		posts_with_tag_lst = append(posts_with_tag_lst, post_info_map)
	}
	//-----------------------------


	type tmpl_data struct {
		Tag_str                string
		Posts_with_tag_num_int int64
		Images_with_tag_int    int64
		Posts_with_tag_lst     []map[string]interface{}
		Sys_release_info       gf_core.SysReleaseInfo
		Is_subtmpl_def         func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	objectTypeStr := "post"
	posts_with_tag_count_int, gfErr := db__get_objects_with_tag_count(pTagStr, objectTypeStr, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)

	err := p_tmpl.Execute(p_resp,
		tmpl_data{
			Tag_str:                pTagStr,
			Posts_with_tag_num_int: posts_with_tag_count_int,
			Images_with_tag_int:    0, // FIX!! - image tagging is now implemented, and so counting images with tag occurance should be done ASAP. 
			Posts_with_tag_lst:     posts_with_tag_lst,
			Sys_release_info:       sysReleaseInfo,
			//-------------------------------------------------
			// IS_SUBTEMPLATE_DEFINED
			Is_subtmpl_def: func(p_subtemplate_name_str string) bool {
				for _, n := range p_subtemplates_names_lst {
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
			err, "gf_tagger", pRuntimeSys)
		return gfErr
	}

	return nil
}	