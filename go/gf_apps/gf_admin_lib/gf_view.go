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

func viewRenderTemplateLogin(pMFAconfirmBool bool,
	pTmpl                 *template.Template,
	pSubtemplatesNamesLst []string,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)
	
	type tmplData struct {
		MFA_confirm_bool bool
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := pTmpl.Execute(buff, tmplData{
		MFA_confirm_bool: pMFAconfirmBool,
		Sys_release_info: sysReleaseInfo,
		
		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
		Is_subtmpl_def: func(pSubtemplateNameStr string) bool {
			for _, n := range pSubtemplatesNamesLst {
				if n == pSubtemplateNameStr {
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
			err, "gf_admin", pRuntimeSys)
		return "", gfErr
	}

	templateRenderedStr := buff.String()
	return templateRenderedStr, nil
}

//------------------------------------------------

func viewRenderTemplateDashboard(pTmpl *template.Template,
	pSubtemplatesNamesLst []string,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)
	
	type tmpl_data struct {
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := pTmpl.Execute(buff, tmpl_data{
		Sys_release_info: sysReleaseInfo,
		
		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
		Is_subtmpl_def: func(pSubtemplateNameStr string) bool {
			for _, n := range pSubtemplatesNamesLst {
				if n == pSubtemplateNameStr {
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
			err, "gf_admin", pRuntimeSys)
		return "", gfErr
	}

	templateRenderedStr := buff.String()
	return templateRenderedStr, nil
}

//-------------------------------------------------

func templatesLoad(pTemplatesPathsMap map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) (*gfTemplates, *gf_core.GFerror) {

	loginTemplateFilepathStr     := pTemplatesPathsMap["gf_admin_login"]
	dashboardTemplateFilepathStr := pTemplatesPathsMap["gf_admin_dashboard"]

	lTmpl, lSubtemplatesNamesLst, gfErr := gf_core.TemplatesLoad(loginTemplateFilepathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	dTmpl, dSubtemplatesNamesLst, gfErr := gf_core.TemplatesLoad(dashboardTemplateFilepathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	templates := &gfTemplates{
		loginTmpl:                     lTmpl,
		loginSubtemplatesNamesLst:     lSubtemplatesNamesLst,
		dashboardTmpl:                 dTmpl,
		dashboardSubtemplatesNamesLst: dSubtemplatesNamesLst,
	}
	return templates, nil
}