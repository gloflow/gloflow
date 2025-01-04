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

package gf_publisher_lib

import (
	"io"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//---------------------------------------------------

func postsBrowserRenderTemplate(pPostsPagesLst [][]*gf_publisher_core.GFpost,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_posts_page_size_int    int, // 5
	p_resp                   io.Writer,
	pRuntimeSys              *gf_core.RuntimeSys) *gf_core.GFerror {

	pagesLst := [][]map[string]interface{}{}
	for _, postsPageLst := range pPostsPagesLst {

		pagePostsLst := []map[string]interface{}{}
		for _, post := range postsPageLst {

			// POST_INFO
			postInfoMap := map[string]interface{}{
				"post_id_str":                post.IDstr,
				"post_title_str":             post.TitleStr,
				"post_creation_datetime_str": post.CreationDatetimeStr,
				"post_thumbnail_url_str":     post.ThumbnailURLstr,
				"images_number_int":          len(post.ImagesIDsLst),
			}

			//---------------
			// TAGS
			if len(post.TagsLst) > 0 {
				post_tags_lst := []string{}
				for _, tag_str := range post.TagsLst {

					// IMPORTANT!! - some tags attached to posts are emtpy strings ""
					if tag_str != "" {
						post_tags_lst = append(post_tags_lst, tag_str)
					}
				}

				postInfoMap["post_has_tags_bool"] = true
				postInfoMap["post_tags_lst"]      = post_tags_lst
			} else {
				postInfoMap["post_has_tags_bool"] = false
			}

			//---------------

			pagePostsLst = append(pagePostsLst, postInfoMap)
		}
		pagesLst = append(pagesLst, pagePostsLst)
	}

	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)
	
	type tmpl_data struct {
		Posts_pages_lst  [][]map[string]interface{}
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	err := p_tmpl.Execute(p_resp, tmpl_data{
		Posts_pages_lst:  pagesLst,
		Sys_release_info: sysReleaseInfo,
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
		gfErr := gf_core.ErrorCreate("failed to render the posts browser template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_publisher_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}