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

package gf_crawl_lib

import (
	"io"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
func dashboard__render_template(p_tmpl *template.Template,
	p_subtemplates_names_lst []string,
	p_resp                   io.Writer,
	p_runtime_sys            *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_view.dashboard__render_template()")

	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)

	type tmpl_data struct {
		Sys_release_info gf_core.Sys_release_info
		Is_subtmpl_def   func(string) bool //used inside the main_template to check if the subtemplate is defined
	}

	err := p_tmpl.Execute(p_resp,tmpl_data{
		Sys_release_info: sys_release_info,
		//-------------------------------------------------
		//IS_SUBTEMPLATE_DEFINED
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
		gf_err := gf_core.Error__create("failed to render the domains_browser template",
            "template_render_error",
            &map[string]interface{}{},
            err, "gf_crawl_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}