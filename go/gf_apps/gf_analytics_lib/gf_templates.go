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

package gf_analytics_lib

import (
	// "fmt"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type gf_templates struct {
	dashboard__tmpl                   *template.Template
	dashboard__subtemplates_names_lst []string
}

//-------------------------------------------------

func tmplLoad(p_templates_paths_map map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_templates, *gf_core.GFerror) {

	main_template_filepath_str := p_templates_paths_map["gf_analytics_dashboard"]

	dashboard__tmpl, subtemplates_names_lst, gfErr := gf_core.TemplatesLoad(main_template_filepath_str, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	gf_templates := &gf_templates{
		dashboard__tmpl:                   dashboard__tmpl,
		dashboard__subtemplates_names_lst: subtemplates_names_lst,
	}
	return gf_templates, nil
}