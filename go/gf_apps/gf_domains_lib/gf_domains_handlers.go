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

package gf_domains_lib

import (
	"fmt"
	"time"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func Init_handlers(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_domains_handlers.Init_handlers()")

	//---------------------
	//TEMPLATES

	gf_templates, gf_err := tmpl__load(p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	/*main_template_filename_str := "gf_domains_browser.html"
	templates_dir_path_str     := "./templates"

	domains_browser__tmpl, subtemplates_names_lst, gf_err := gf_core.Templates__load(main_template_filename_str, templates_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}*/

	/*domains_browser__tmpl, err := template.New("gf_domains_browser.html").ParseFiles(template_path_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse a template",
			"template_create_error",
			&map[string]interface{}{"template_path_str":template_path_str,},
			err, "gf_images_lib", p_runtime_sys)
		return gf_err
	}*/
	//---------------------

	//---------------------
	//POSTS_ELEMENTS
	http.HandleFunc("/a/domains/browser", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST - /a/domains/browser ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------
			//response_format_str - "json"|"html"

			qs_map := p_req.URL.Query()
			fmt.Println(qs_map)

			/*//response_format_str - "j"(for json)|"h"(for html)
			response_format_str := gf_rpc_lib.Get_response_format(qs_map,
															p_log_fun)*/
			//--------------------
			//GET DOMAINS FROM DB
			domains_lst, gf_err := db__get_domains(p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/domains/browser", "rpc_handler failed getting domains", gf_err, p_resp, p_runtime_sys)
				return
			}
			//--------------------
			//RENDER TEMPLATE
			gf_err = domains_browser__render_template(domains_lst,
				gf_templates.domains_browser__tmpl,
				gf_templates.domains_browser__subtemplates_names_lst,
				p_resp,
				p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/domains/browser", "failed to render domains_browser page", gf_err, p_resp, p_runtime_sys)
				return
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/a/domains/browser", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//---------------------

	return nil
}