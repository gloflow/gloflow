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

package gf_domains_lib

import (
	"io"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------

func domainsBrowserRenderTemplate(pDomainsLst []GFdomain,
	pTmpl                 *template.Template,
	pSubtemplatesNamesLst []string,
	pResp                 io.Writer,
	pRuntimeSys           *gf_core.RuntimeSys) *gf_core.GFerror {

	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)

	type tmpl_data struct {
		Domains_lst      []GFdomain
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool //used inside the main_template to check if the subtemplate is defined
	}

	err := pTmpl.Execute(pResp, tmpl_data{
		Domains_lst:      pDomainsLst,
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
		gfErr := gf_core.ErrorCreate("failed to render the domains_browser template",
            "template_render_error",
            map[string]interface{}{"domains_lst": pDomainsLst,},
            err, "gf_domains_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}