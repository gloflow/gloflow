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

func post__render_template(p_post *gf_publisher_core.GFpost,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_resp                   io.Writer,
	pRuntimeSys              *gf_core.RuntimeSys) *gf_core.GFerror {
	
	template_post_elements_lst, gfErr := package_post_elements_infos(p_post, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	image_post_elements_og_info_lst, gfErr := getImagePostElementsFBOpenGraphInfo(p_post, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	post_tags_lst := []string{}
	for _, tag_str := range p_post.TagsLst {
		post_tags_lst = append(post_tags_lst, tag_str)
	}

	type tmpl_data struct {
		Post_title_str                  string
		Post_tags_lst                   []string
		Post_description_str            string
		Post_poster_user_name_str       string
		Post_thumbnail_url_str          string
		Post_elements_lst               []map[string]interface{}
		Image_post_elements_og_info_lst []map[string]string
		Sys_release_info                gf_core.SysReleaseInfo
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

	sys_release_info := gf_core.GetSysReleseInfo(pRuntimeSys)

	err := p_tmpl.Execute(p_resp, tmpl_data{
		Post_title_str:                  p_post.TitleStr,
		Post_tags_lst:                   post_tags_lst,
		Post_description_str:            p_post.DescriptionStr,
		Post_poster_user_name_str:       p_post.PosterUserNameStr,
		Post_thumbnail_url_str:          p_post.ThumbnailURLstr,
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
		gfErr := gf_core.ErrorCreate("failed to render the post template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_publisher_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//--------------------------------------------------

func package_post_elements_infos(p_post *gf_publisher_core.GFpost,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	template_post_elements_lst := []map[string]interface{}{}

	for _, postElement := range p_post.PostElementsLst {

		gfErr := gf_publisher_core.Verify_post_element_type(postElement.TypeStr, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		post_element_tags_lst := []string{}
		for _, tag_str := range postElement.TagsLst {
			post_element_tags_lst = append(post_element_tags_lst, tag_str)
		}

		switch postElement.TypeStr {
			case "link":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool": false, //for mustache template conditionals
					"post_element_type__image_bool": false,
					"post_element_type__link_bool":  true,
					"post_element_description_str":  postElement.DescriptionStr,
					"post_element_extern_url_str":   postElement.ExternURLstr,
				}
				template_post_elements_lst = append(template_post_elements_lst, post_element_map)
				continue
			case "image":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool":             false, //for mustache template conditionals
					"post_element_type__image_bool":             true,
					"post_element_type__link_bool":              false,
					"post_element_img_thumbnail_medium_url_str": postElement.ImgThumbnailMediumURLstr,
					"post_element_img_thumbnail_large_url_str":  postElement.ImgThumbnailLargeURLstr,
					"tags_lst":                                  post_element_tags_lst,
				}
				template_post_elements_lst = append(template_post_elements_lst, post_element_map)
				continue
			case "video":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool": true, //for mustache template conditionals
					"post_element_type__image_bool": false,
					"post_element_type__link_bool":  false,
					"post_element_extern_url_str":   postElement.ExternURLstr,
					"tags_lst":                      post_element_tags_lst,
				}
				template_post_elements_lst = append(template_post_elements_lst, post_element_map)
				continue
		}
	}
	return template_post_elements_lst, nil
}

//--------------------------------------------------

func getImagePostElementsFBOpenGraphInfo(pPost *gf_publisher_core.GFpost,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]string, *gf_core.GFerror) {

	imagePostElementsLst, gfErr := gf_publisher_core.Get_post_elements_of_type(pPost, "image", pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	var topImagePostElementsLst []*gf_publisher_core.GFpostElement
	if len(imagePostElementsLst) > 5 {
		topImagePostElementsLst = imagePostElementsLst[:5]
	} else { 
		topImagePostElementsLst = imagePostElementsLst
	}

	//---------------------
	ogInfoLst := []map[string]string{}
	for _, postElement := range topImagePostElementsLst {
		d := map[string]string{
			"img_thumbnail_medium_absolute_url_str": postElement.ImgThumbnailMediumURLstr,
		}
		ogInfoLst = append(ogInfoLst, d)
	}
	
	//---------------------

	return ogInfoLst,nil
}