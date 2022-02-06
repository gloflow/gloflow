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

package gf_landing_page_lib

import (
	"io"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
func Pipeline__render_landing_page(p_imgs__max_random_cursor_position_int int, // 500
	p_posts__max_random_cursor_position_int int,
	p_featured_posts_to_get_int int, // 5
	p_featured_imgs_to_get_int  int, // 10
	p_tmpl                      *template.Template,
	p_subtemplates_names_lst    []string,
	p_resp                      io.Writer,
	p_runtime_sys               *gf_core.Runtime_sys) *gf_core.GF_error {

	//-------------------
	// FEATURED_IMAGES - two random groups of images are fetched
	featured_imgs_0_lst, gf_err := get_featured_imgs(p_imgs__max_random_cursor_position_int,
		p_featured_imgs_to_get_int,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	featured_imgs_1_lst, gf_err := get_featured_imgs(p_imgs__max_random_cursor_position_int,
		p_featured_imgs_to_get_int,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//-------------------
	featured_posts_lst, gf_err := get_featured_posts(p_posts__max_random_cursor_position_int,
		p_featured_posts_to_get_int,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	gf_err = render_template(featured_posts_lst,
		featured_imgs_0_lst,
		featured_imgs_1_lst,
		p_tmpl,
		p_subtemplates_names_lst,
		p_resp,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}