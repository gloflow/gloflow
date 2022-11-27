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
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//------------------------------------------------

func Get_posts_page(p_page_index_int int,
	p_page_elements_num_int int,
	pRuntimeSys             *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	cursor_start_position_int := p_page_index_int*p_page_elements_num_int
	page_lst, gfErr          := gf_publisher_core.DBgetPostsPage(cursor_start_position_int, p_page_elements_num_int, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	serialized_page_lst := []map[string]interface{}{}
	for _, post := range page_lst {
		postMap := map[string]interface{}{
			"title_str":             post.TitleStr,
			"images_number_str":     len(post.ImagesIDsLst),
			"creation_datetime_str": post.CreationDatetimeStr,
			"thumbnail_url_str":     post.ThumbnailURLstr,
			"tags_lst":              post.TagsLst,
		}
		serialized_page_lst = append(serialized_page_lst, postMap)
	}

	return serialized_page_lst, nil
}

//------------------------------------------------
// get initial pages - the pages that are rendered in the initial HTML template. 
//                     subsequent pages are loaded as AJAX requests, via HTTP API. 

func RenderInitialPages(p_response_format_str string,
	p_initial_pages_num_int  int, //6
	p_page_size_int          int, //5
	p_tmpl                   *template.Template,
	p_subtempaltes_names_lst []string,
	p_resp                   io.Writer,
	pRuntimeSys              *gf_core.RuntimeSys) *gf_core.GFerror {
	
	postsPagesLst := [][]*gf_publisher_core.GFpost{}

	for i:=0; i < p_initial_pages_num_int; i++ {

		start_position_int := i*p_page_size_int
		// int end_position_int   = start_position_int+p_page_size_int;

		pRuntimeSys.LogFun("INFO", fmt.Sprintf(">>>>>>> start_position_int - %d - %d", start_position_int, p_page_size_int))

		// initial page might be larger then subsequent pages, that are requested 
		// dynamically by the front-end
		pageLst, gfErr := gf_publisher_core.DBgetPostsPage(start_position_int, p_page_size_int, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		postsPagesLst = append(postsPagesLst, pageLst)
	}
	
	gfErr := postsBrowserRenderTemplate(postsPagesLst, p_tmpl, p_subtempaltes_names_lst, p_page_size_int, p_resp, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}