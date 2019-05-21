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
	"fmt"
	"io"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
func Get_posts_page(p_page_index_int int,
	p_page_elements_num_int int,
	p_runtime_sys           *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_posts_browser_pipelines.Get_posts_page()")

	cursor_start_position_int := p_page_index_int*p_page_elements_num_int
	page_lst, gf_err          := DB__get_posts_page(cursor_start_position_int, p_page_elements_num_int, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	serialized_page_lst := []map[string]interface{}{}
	for _, post := range page_lst {
		post_map := map[string]interface{}{
			"title_str":             post.Title_str,
			"images_number_str":     len(post.Images_ids_lst),
			"creation_datetime_str": post.Creation_datetime_str,
			"thumbnail_url_str":     post.Thumbnail_url_str,
			"tags_lst":              post.Tags_lst,
		}
		serialized_page_lst = append(serialized_page_lst, post_map)
	}

	return serialized_page_lst, nil
}

//------------------------------------------------
//get initial pages - the pages that are rendered in the initial HTML template. 
//                    subsequent pages are loaded as AJAX requests, via HTTP API. 

func Render_initial_pages(p_response_format_str string,
	p_initial_pages_num_int  int, //6
	p_page_size_int          int, //5
	p_tmpl                   *template.Template,
	p_subtempaltes_names_lst []string,
	p_resp                   io.Writer,
	p_runtime_sys            *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_posts_browser_pipelines.Render_initial_pages()")
	
	posts_pages_lst := [][]*Gf_post{}

	for i:=0;i<p_initial_pages_num_int;i++ {

		start_position_int := i*p_page_size_int
		//int end_position_int   = start_position_int+p_page_size_int;

		p_runtime_sys.Log_fun("INFO", fmt.Sprintf(">>>>>>> start_position_int - %d - %d", start_position_int, p_page_size_int))

		//initial page might be larger then subsequent pages, that are requested 
		//dynamically by the front-end
		page_lst, gf_err := DB__get_posts_page(start_position_int, p_page_size_int, p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		posts_pages_lst = append(posts_pages_lst, page_lst)
	}
	
	gf_err := posts_browser__render_template(posts_pages_lst, p_tmpl, p_subtempaltes_names_lst, p_page_size_int, p_resp, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}