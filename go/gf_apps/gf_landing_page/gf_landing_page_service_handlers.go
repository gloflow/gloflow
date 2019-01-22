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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//------------------------------------------------
func init_handlers(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_landing_page_service_handlers.init_handlers()")

	template_path_str := "./templates/gf_landing_page.html"
	tmpl, err         := template.New("gf_landing_page.html").ParseFiles(template_path_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse a template",
			"template_create_error",
			&map[string]interface{}{"template_path_str":template_path_str,},
			err, "gf_landing_page", p_runtime_sys)
		return gf_err
	}

	//---------------------
	http.HandleFunc("/landing/main/", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /landing/main/ ----------")

		if p_req.Method == "GET" {
			gf_err := Pipeline__get_landing_page(2000, //p_max_random_cursor_position_int
				5,  //p_featured_posts_to_get_int
				10, //p_featured_imgs_to_get_int
				tmpl,
				p_resp,
				p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/landing/main", "get landing_page failed", gf_err, p_resp, p_runtime_sys)
				return
			}
		}
	})
	//---------------------
	http.HandleFunc("/landing/register_invite_email", func(p_resp http.ResponseWriter, p_req *http.Request) {

	})
	//---------------------
	return nil
}