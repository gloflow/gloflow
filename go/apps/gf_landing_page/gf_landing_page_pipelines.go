/*
GloFlow media management/publishing system
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

package main

import (
	"net/http"
	"text/template"
	"gopkg.in/mgo.v2"
)
//------------------------------------------------
func Pipeline__get_landing_page(p_max_random_cursor_position_int int, //500
				p_featured_posts_to_get_int int, //5
				p_featured_imgs_to_get_int  int, //10
				p_tmpl                      *template.Template,
				p_resp                      http.ResponseWriter,
				p_mongodb_coll              *mgo.Collection,
				p_log_fun                   func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_landing_page_pipelines.Pipeline__get_landing_page()")

	featured_posts_lst,err := get_featured_posts(p_max_random_cursor_position_int,
											p_featured_posts_to_get_int,
											p_mongodb_coll,
											p_log_fun)
	if err != nil {
		return err
	}

	featured_imgs_lst,err := get_featured_imgs(p_max_random_cursor_position_int,
										p_featured_imgs_to_get_int,
										p_mongodb_coll,
										p_log_fun)
	if err != nil {
		return err
	}

	err = render_template(featured_posts_lst,
					featured_imgs_lst,
					p_tmpl,
					p_resp,
					p_log_fun)
	if err != nil {
		return err
	}

	return nil
}