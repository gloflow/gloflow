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
	"gf_rpc_lib"
)
//------------------------------------------------
func init_handlers(p_mongodb_coll *mgo.Collection, p_log_fun func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_landing_page_service_handlers.init_handlers()")

	tmpl,err := template.New("gf_landing_page.html").ParseFiles("./templates/gf_landing_page.html")
	if err != nil {
		return err
	}

	//---------------------
	http.HandleFunc("/landing/main/",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_log_fun("INFO","INCOMING HTTP REQUEST - /landing/main/ ----------")

		if p_req.Method == "GET" {
			err := Pipeline__get_landing_page(2000, //p_max_random_cursor_position_int
				5,  //p_featured_posts_to_get_int
				10, //p_featured_imgs_to_get_int
				tmpl,
				p_resp,
				p_mongodb_coll,
				p_log_fun)
			if err != nil {
				gf_rpc_lib.Error__in_handler("/landing/main", err, "get landing_page failed", p_resp, p_mongodb_coll, p_log_fun)
				return
			}
		}
	})
	//---------------------
	http.HandleFunc("/landing/register_invite_email",func(p_resp http.ResponseWriter, p_req *http.Request) {

	})
	//---------------------

	return nil
}