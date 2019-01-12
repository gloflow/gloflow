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

	"gf_core"
)
//------------------------------------------------
func render_template(p_featured_posts_lst []*Featured_post,
				p_featured_imgs_lst []*Featured_img,
				p_tmpl              *template.Template,
				p_resp              http.ResponseWriter,
				p_log_fun           func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_landing_page_view.render_template()")
	
	sys_release_info := gf_core.Get_sys_relese_info(p_log_fun)
	
	type tmpl_data struct {
		Featured_posts_lst []*Featured_post
		Featured_imgs_lst  []*Featured_img
		Sys_release_info   gf_core.Sys_release_info
	}

	err := p_tmpl.Execute(p_resp,tmpl_data{
		Featured_posts_lst:p_featured_posts_lst,
		Featured_imgs_lst :p_featured_imgs_lst,
		Sys_release_info  :sys_release_info,
	})

	if err != nil {
		return err
	}

	return nil
}