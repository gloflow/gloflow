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
	p_runtime_sys            *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	sys_release_info := gf_core.GetSysReleseInfo(p_runtime_sys)
	
	type tmpl_data struct {
		MFA_confirm_bool bool
		SysReleaseInfo   gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := p_tmpl.Execute(buff, tmpl_data{
		MFA_confirm_bool: p_mfa_confirm_bool,
		SysReleaseInfo:   sys_release_info,
		
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
		gfErr := gf_core.ErrorCreate("failed to render the admin login template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_admin", p_runtime_sys)
		return "", gfErr
	}

	template_rendered_str := buff.String()
	return template_rendered_str, nil
}

//------------------------------------------------
func view__render_template_dashboard(p_tmpl *template.Template,
	p_subtemplates_names_lst []string,
	p_runtime_sys            *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(p_runtime_sys)
	
	type tmpl_data struct {
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := p_tmpl.Execute(buff, tmpl_data{
		Sys_release_info: sysReleaseInfo,
		
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
		gfErr := gf_core.ErrorCreate("failed to render the admin dashboard template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_admin", p_runtime_sys)
		return "", gfErr
	}

	template_rendered_str := buff.String()
	return template_rendered_str, nil
}

//-------------------------------------------------
func templatesLoad(p_templates_paths_map map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_templates, *gf_core.GFerror) {

	loginTemplateFilepathStr     := p_templates_paths_map["gf_admin_login"]
	dashboardTemplateFilepathStr := p_templates_paths_map["gf_admin_dashboard"]

	l_tmpl, l_subtemplates_names_lst, gfErr := gf_core.TemplatesLoad(loginTemplateFilepathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	d_tmpl, d_subtemplates_names_lst, gfErr := gf_core.TemplatesLoad(dashboardTemplateFilepathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	templates := &gf_templates{
		login__tmpl:                   l_tmpl,
		login__subtemplates_names_lst: l_subtemplates_names_lst,
		dashboard__tmpl:                   d_tmpl,
		dashboard__subtemplates_names_lst: d_subtemplates_names_lst,
	}
	return templates, nil
}