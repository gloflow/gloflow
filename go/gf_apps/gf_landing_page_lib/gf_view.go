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

package gf_landing_page_lib

import (
	"io"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------

func renderTemplate(pFeaturedPostsLst []*GFfeaturedPost,
	p_featured_imgs_0_lst []*GFfeaturedImage,
	p_featured_imgs_1_lst []*GFfeaturedImage,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pResp                 io.Writer,
	pRuntimeSys           *gf_core.RuntimeSys) *gf_core.GFerror {
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)
	
	type tmplData struct {
		Featured_posts_lst  []*GFfeaturedPost
		Featured_imgs_0_lst []*GFfeaturedImage
		Featured_imgs_1_lst []*GFfeaturedImage
		Sys_release_info    gf_core.SysReleaseInfo
		Is_subtmpl_def      func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	err := pTemplate.Execute(pResp, tmplData{
		Featured_posts_lst:  pFeaturedPostsLst,
		Featured_imgs_0_lst: p_featured_imgs_0_lst,
		Featured_imgs_1_lst: p_featured_imgs_1_lst,
		Sys_release_info:    sysReleaseInfo,
		
		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
		Is_subtmpl_def: func(p_subtemplate_name_str string) bool {
			for _, n := range pSubtemplatesNamesLst {
				if n == p_subtemplate_name_str {
					return true
				}
			}
			return false
		},
		
		//-------------------------------------------------
	})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to render the landing_page template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_landing_page", pRuntimeSys)
		return gfErr
	}

	return nil
}