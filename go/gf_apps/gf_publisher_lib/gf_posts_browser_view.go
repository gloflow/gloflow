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
func posts_browser__render_template(p_posts_pages_lst [][]*gf_publisher_core.Gf_post, // list-of-lists
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_posts_page_size_int    int, // 5
	p_resp                   io.Writer,
	p_runtime_sys            *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_posts_browser_view.posts_browser__render_template()")

	pages_lst := [][]map[string]interface{}{}
	for _, posts_page_lst := range p_posts_pages_lst {

		page_posts_lst := []map[string]interface{}{}
		for _, post := range posts_page_lst {

			post_info_map := map[string]interface{}{
				"post_title_str":             post.Title_str,
				"post_creation_datetime_str": post.Creation_datetime_str,
				"post_thumbnail_url_str":     post.Thumbnail_url_str,
				"images_number_int":          len(post.Images_ids_lst),
			}

			//---------------
			// TAGS
			if len(post.Tags_lst) > 0 {
				post_tags_lst := []string{}
				for _, tag_str := range post.Tags_lst {

					// IMPORTANT!! - some tags attached to posts are emtpy strings ""
					if tag_str != "" {
						post_tags_lst = append(post_tags_lst, tag_str)
					}
				}

				post_info_map["post_has_tags_bool"] = true
				post_info_map["post_tags_lst"]      = post_tags_lst
			} else {
				post_info_map["post_has_tags_bool"] = false
			}

			//---------------

			page_posts_lst = append(page_posts_lst, post_info_map)
		}
		pages_lst = append(pages_lst, page_posts_lst)
	}

	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)
	
	type tmpl_data struct {
		Posts_pages_lst  [][]map[string]interface{}
		Sys_release_info gf_core.Sys_release_info
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	err := p_tmpl.Execute(p_resp, tmpl_data{
		Posts_pages_lst:  pages_lst,
		Sys_release_info: sys_release_info,
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
		gf_err := gf_core.ErrorCreate("failed to render the posts browser template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}