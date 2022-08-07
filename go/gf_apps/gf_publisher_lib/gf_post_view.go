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

//--------------------------------------------------
func post__render_template(p_post *gf_publisher_core.Gf_post,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_resp                   io.Writer,
	p_runtime_sys            *gf_core.RuntimeSys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_post_view.post__render_template()")
	
	template_post_elements_lst, gf_err := package_post_elements_infos(p_post, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	image_post_elements_og_info_lst, gf_err := get_image_post_elements_FBOpenGraph_info(p_post, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	post_tags_lst := []string{}
	for _, tag_str := range p_post.Tags_lst {
		post_tags_lst = append(post_tags_lst,tag_str)
	}

	type tmpl_data struct {
		Post_title_str                  string
		Post_tags_lst                   []string
		Post_description_str            string
		Post_poster_user_name_str       string
		Post_thumbnail_url_str          string
		Post_elements_lst               []map[string]interface{}
		Image_post_elements_og_info_lst []map[string]string
		Sys_release_info                gf_core.Sys_release_info
		Is_subtmpl_def                  func(string) bool //used inside the main_template to check if the subtemplate is defined
	}
	
	/*template_info_map := map[string]interface{}{
		"post_title_str"                       :p_post.title_str,
		"post_tags_lst"                        :post_tags_lst,
		"post_description_str"                 :p_post.description_str,
		"post_poster_user_name_str"   template_str         :p_post.poster_user_name_str,
		"post_elements_lst"                    :template_post_elements_lst,
		"img_thumbnail_medium_absolute_url_str":post_thumbnail_url_str,
		"image_post_elements_og_info_lst"      :image_post_elements_og_info_lst,
	}

	final String template_str = p_template.renderString(template_info_map)
	return template_str;*/

	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)

	err := p_tmpl.Execute(p_resp, tmpl_data{
		Post_title_str:                  p_post.Title_str,
		Post_tags_lst:                   post_tags_lst,
		Post_description_str:            p_post.Description_str,
		Post_poster_user_name_str:       p_post.Poster_user_name_str,
		Post_thumbnail_url_str:          p_post.Thumbnail_url_str,
		Post_elements_lst:               template_post_elements_lst,
		Image_post_elements_og_info_lst: image_post_elements_og_info_lst,
		Sys_release_info:                sys_release_info,

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
		gf_err := gf_core.Error__create("failed to render the post template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}

//--------------------------------------------------
func package_post_elements_infos(p_post *gf_publisher_core.Gf_post,
	p_runtime_sys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_post_view.package_post_elements_infos()")

	template_post_elements_lst := []map[string]interface{}{}

	for _, post_element := range p_post.Post_elements_lst {

		p_runtime_sys.Log_fun("INFO","post_element.Type_str - "+post_element.Type_str)
		gf_err := gf_publisher_core.Verify_post_element_type(post_element.Type_str, p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		post_element_tags_lst := []string{}
		for _, tag_str := range post_element.Tags_lst {
			post_element_tags_lst = append(post_element_tags_lst,tag_str)
		}

		switch post_element.Type_str {
			case "link":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool": false, //for mustache template conditionals
					"post_element_type__image_bool": false,
					"post_element_type__link_bool":  true,
					"post_element_description_str":  post_element.Description_str,
					"post_element_extern_url_str":   post_element.Extern_url_str,
				}
				template_post_elements_lst = append(template_post_elements_lst, post_element_map)
				continue
			case "image":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool":             false, //for mustache template conditionals
					"post_element_type__image_bool":             true,
					"post_element_type__link_bool":              false,
					"post_element_img_thumbnail_medium_url_str": post_element.Img_thumbnail_medium_url_str,
					"post_element_img_thumbnail_large_url_str":  post_element.Img_thumbnail_large_url_str,
					"tags_lst":                                  post_element_tags_lst,
				}
				template_post_elements_lst = append(template_post_elements_lst, post_element_map)
				continue
			case "video":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool": true, //for mustache template conditionals
					"post_element_type__image_bool": false,
					"post_element_type__link_bool":  false,
					"post_element_extern_url_str":   post_element.Extern_url_str,
					"tags_lst":                      post_element_tags_lst,
				}
				template_post_elements_lst = append(template_post_elements_lst, post_element_map)
				continue
		}
	}
	return template_post_elements_lst, nil
}

//--------------------------------------------------
func get_image_post_elements_FBOpenGraph_info(p_post *gf_publisher_core.Gf_post,
	p_runtime_sys *gf_core.RuntimeSys) ([]map[string]string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_post_view.get_image_post_elements_FBOpenGraph_info()")

	image_post_elements_lst, gf_err := gf_publisher_core.Get_post_elements_of_type(p_post, "image", p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	var top_image_post_elements_lst []*gf_publisher_core.Gf_post_element
	if len(image_post_elements_lst) > 5 {

		// getRange() - returns an Iterable<String>
		top_image_post_elements_lst = image_post_elements_lst[:5] //new List.from(image_post_elements_lst.getRange(0,5))
	} else { 
		top_image_post_elements_lst = image_post_elements_lst
	}

	//---------------------
	og_info_lst := []map[string]string{}
	for _, post_element := range top_image_post_elements_lst {
		d := map[string]string{
			"img_thumbnail_medium_absolute_url_str": post_element.Img_thumbnail_medium_url_str,
		}
		og_info_lst = append(og_info_lst, d)
	}
	
	//---------------------

	return og_info_lst,nil
}