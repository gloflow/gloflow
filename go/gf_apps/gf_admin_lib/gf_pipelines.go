/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_admin_lib

import (
	"text/template"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
func Pipeline__render(p_tmpl *template.Template,
	p_subtemplates_names_lst []string,
	p_ctx                    context.Context,
	p_runtime_sys            *gf_core.Runtime_sys) (string, *gf_core.GF_error) {




	template_rendered_str, gf_err := admin__render_template(p_tmpl,
		p_subtemplates_names_lst,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}


	return template_rendered_str, nil
}