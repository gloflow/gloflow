/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

package gf_identity

import (
	"bytes"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------

type gfTemplates struct {
	loginTmpl                 *template.Template
	loginSubtemplatesNamesLst []string
}

//------------------------------------------------

func viewRenderTemplateLogin(pAuthSubsystemTypeStr string,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)
	

	type templateData struct {
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := pTemplate.Execute(buff, templateData{
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
		gfErr := gf_core.ErrorCreate("failed to render the login template",
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

	loginTemplateFilepathStr := pTemplatesPathsMap["gf_login"]

	lTmpl, lSubtemplatesNamesLst, gfErr := gf_core.TemplatesLoad(loginTemplateFilepathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	templates := &gfTemplates{
		loginTmpl:                 lTmpl,
		loginSubtemplatesNamesLst: lSubtemplatesNamesLst,
	}
	return templates, nil
}