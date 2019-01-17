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

package gf_domains_lib

import (
	"net/http"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
func domains_browser__render_template(p_domains_lst []Gf_domain,
	p_tmpl        *template.Template,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains_view.domains_browser__render_template()")

	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)

	type tmpl_data struct {
		Domains_lst      []Gf_domain
		Sys_release_info gf_core.Sys_release_info
	}

	err := p_tmpl.Execute(p_resp,tmpl_data{
		Domains_lst:     p_domains_lst,
		Sys_release_info:sys_release_info,
	})

	if err != nil {
		gf_err := gf_core.Error__create("failed to render the domains_browser template",
            "template_render_error",
            &map[string]interface{}{"domains_lst":p_domains_lst,},
            err, "gf_domains_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}