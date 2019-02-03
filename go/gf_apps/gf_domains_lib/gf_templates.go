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
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type gf_templates struct {
	domains_browser__tmpl                   *template.Template
	domains_browser__subtemplates_names_lst []string
}
//-------------------------------------------------
func tmpl__load(p_runtime_sys *gf_core.Runtime_sys) (*gf_templates, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_templates.tmpl__load()")

	main_template_filename_str := "gf_domains_browser.html"
	templates_dir_path_str     := "./web/gf_apps/gf_domains_lib/templates"

	domains_browser__tmpl, subtemplates_names_lst, gf_err := gf_core.Templates__load(main_template_filename_str, templates_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	gf_templates := &gf_templates{
		domains_browser__tmpl:                   domains_browser__tmpl,
		domains_browser__subtemplates_names_lst: subtemplates_names_lst,
	}
	return gf_templates, nil
}