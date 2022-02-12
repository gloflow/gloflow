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
	"bytes"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
func view__render_template_login(p_mfa_confirm_bool bool,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_runtime_sys            *gf_core.Runtime_sys) (string, *gf_core.GF_error) {
	
	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)
	
	type tmpl_data struct {
		MFA_confirm_bool bool
		Sys_release_info gf_core.Sys_release_info
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := p_tmpl.Execute(buff, tmpl_data{
		MFA_confirm_bool: p_mfa_confirm_bool,
		Sys_release_info: sys_release_info,
		
		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
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
		gf_err := gf_core.Error__create("failed to render the admin login template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_admin", p_runtime_sys)
		return "", gf_err
	}

	template_rendered_str := buff.String()
	return template_rendered_str, nil
}

//------------------------------------------------
func view__render_template_dashboard(p_tmpl *template.Template,
	p_subtemplates_names_lst []string,
	p_runtime_sys            *gf_core.Runtime_sys) (string, *gf_core.GF_error) {
	
	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)
	
	type tmpl_data struct {
		Sys_release_info gf_core.Sys_release_info
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := p_tmpl.Execute(buff, tmpl_data{
		Sys_release_info: sys_release_info,
		
		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
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
		gf_err := gf_core.Error__create("failed to render the admin dashboard template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_admin", p_runtime_sys)
		return "", gf_err
	}

	template_rendered_str := buff.String()
	return template_rendered_str, nil
}