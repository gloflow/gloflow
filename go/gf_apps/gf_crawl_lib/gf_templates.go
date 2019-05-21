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
	"fmt"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type gf_templates struct {
	dashboard__tmpl                   *template.Template
	dashboard__subtemplates_names_lst []string
}

//-------------------------------------------------
func tmpl__load(p_templates_dir_path_str string, p_runtime_sys *gf_core.Runtime_sys) (*gf_templates, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_templates.tmpl__load()")

	main_template_filename_str := "gf_crawl_dashboard.html"
	templates_dir_path_str     := fmt.Sprintf("%s/gf_crawl_dashboard", p_templates_dir_path_str)

	dashboard__tmpl, subtemplates_names_lst, gf_err := gf_core.Templates__load(main_template_filename_str, templates_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	gf_templates := &gf_templates{
		dashboard__tmpl:                   dashboard__tmpl,
		dashboard__subtemplates_names_lst: subtemplates_names_lst,
	}
	return gf_templates, nil
}