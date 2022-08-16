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

package gf_images_flows

import (
	"fmt"
	"context"
	"strconv"
	"bytes"
	"text/template"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------
func flows__render_initial_page(p_flow_name_str string,
	p_initial_pages_num_int  int, // 6
	p_page_size_int          int, // 5
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_ctx                    context.Context,
	p_runtime_sys            *gf_core.RuntimeSys) (string, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows_views.flows__render_initial_page()")

	//---------------------
	// GET_TEMPLATE_DATA

	pages_lst := [][]*gf_images_core.Gf_image{}

	for i:=0; i < p_initial_pages_num_int; i++ {

		start_position_int := i*p_page_size_int
		// int end_position_int = start_position_int+p_page_size_int;

		p_runtime_sys.Log_fun("INFO", fmt.Sprintf(">>>>>>> start_position_int - %d - %d", start_position_int, p_page_size_int))
		//------------
		// DB GET PAGE

		// initial page might be larger then subsequent pages, that are requested 
		// dynamically by the front-end
		page_lst, gf_err := flows_db__get_page(p_flow_name_str,
			start_position_int, // p_cursor_start_position_int
			p_page_size_int,    // p_elements_num_int
			p_ctx,
			p_runtime_sys)

		if gf_err != nil {
			return "", gf_err
		}

		//------------

		pages_lst = append(pages_lst, page_lst)
	}


	flow_pages_num_int, gf_err := flows_db__get_pages_total_num(p_flow_name_str,
		p_page_size_int,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	//---------------------
	template_rendered_str, gf_err := flows__render_template(p_flow_name_str,
		pages_lst,
		flow_pages_num_int,
		p_tmpl,
		p_subtemplates_names_lst,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	return template_rendered_str, nil
}

//-------------------------------------------------
func flows__render_template(p_flow_name_str string,
	p_images_pages_lst       [][]*gf_images_core.GF_image, // 2D list
	p_flow_pages_num_int     int64,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_runtime_sys            *gf_core.RuntimeSys) (string, *gf_core.GF_error) {

	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)
	//-------------------------
	images_pages_lst := [][]map[string]interface{}{}
	for _, images_page_lst := range p_images_pages_lst {

		page_images_lst := []map[string]interface{}{}
		for _, image := range images_page_lst {

			meta_json_bytes_lst, _ := json.Marshal(image.Meta_map)
			meta_json_str          := string(meta_json_bytes_lst)

			image_info_map := map[string]interface{}{
				"creation_unix_time_str":    strconv.FormatFloat(image.Creation_unix_time_f, 'f', 6, 64),
				"id_str":                    image.Id_str,
				"title_str":                 image.Title_str,
				"meta_json_str":             meta_json_str,
				"format_str":                image.Format_str,
				"thumbnail_small_url_str":   image.Thumbnail_small_url_str,
				"thumbnail_medium_url_str":  image.Thumbnail_medium_url_str,
				"thumbnail_large_url_str":   image.Thumbnail_large_url_str,
				"image_origin_page_url_str": image.Origin_page_url_str,
			}

			if len(image.Tags_lst) > 0 {
				image_info_map["image_has_tags_bool"] = true
				image_info_map["tags_lst"]            = image.Tags_lst
			} else {
				image_info_map["image_has_tags_bool"] = false
			}

			page_images_lst = append(page_images_lst, image_info_map)
		}
		images_pages_lst = append(images_pages_lst, page_images_lst)
	}

	//-------------------------

	type tmpl_data struct {
		Flow_name_str      string
		Images_pages_lst   [][]map[string]interface{}
		Flow_pages_num_int int64
		Sys_release_info   gf_core.Sys_release_info
		Is_subtmpl_def     func(string) bool //used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := p_tmpl.Execute(buff, tmpl_data{
		Flow_name_str:      p_flow_name_str,
		Images_pages_lst:   images_pages_lst,
		Flow_pages_num_int: p_flow_pages_num_int,
		Sys_release_info:   sys_release_info,

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
		gf_err := gf_core.ErrorCreate("failed to render the images flow template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_images_lib", p_runtime_sys)
		return "", gf_err
	}

	template_rendered_str := buff.String()
	return template_rendered_str, nil
}