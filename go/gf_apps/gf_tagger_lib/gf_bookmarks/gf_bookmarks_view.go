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

package gf_bookmarks

import (
	"bytes"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------

func renderBookmarks(pBookmarksLst []*GFbookmark,
	pTmpl                 *template.Template,
	pSubtemplatesNamesLst []string,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	sys_release_info := gf_core.GetSysReleseInfo(pRuntimeSys)

	type tmpl_data struct {
		Bookmarks_lst    []*GFbookmark
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool // used inside the main_template to check if the subtemplate is defined
	}
	

	buff := new(bytes.Buffer)
	err := pTmpl.Execute(buff,
		tmpl_data{
			Bookmarks_lst:    pBookmarksLst,
			Sys_release_info: sys_release_info,
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

		gf_err := gf_core.ErrorCreate("failed to render the gf_bookmarks template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_tagger", pRuntimeSys)
		return "", gf_err
	}


	template_rendered_str := buff.String()
	return template_rendered_str, nil	
}