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
	"bytes"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------

type gfTemplates struct {
	template             *template.Template
	subtemplatesNamesLst []string
}

//------------------------------------------------

func renderTemplate(pFeaturedPostsLst []*GFfeaturedPost,
	pFeaturedImages0lst   []*GFfeaturedImage,
	pFeaturedImages1lst   []*GFfeaturedImage,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)
	
	type tmplData struct {
		FeaturedPostsLst   []*GFfeaturedPost
		FeaturedImages0lst []*GFfeaturedImage
		FeaturedImages1lst []*GFfeaturedImage
		SysReleaseInfo     gf_core.SysReleaseInfo
		IsSubtmplDef       func(string) bool // used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := pTemplate.Execute(buff, tmplData{
		FeaturedPostsLst:   pFeaturedPostsLst,
		FeaturedImages0lst: pFeaturedImages0lst,
		FeaturedImages1lst: pFeaturedImages1lst,
		SysReleaseInfo:     sysReleaseInfo,
		
		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
		IsSubtmplDef: func(pSubtemplateNameStr string) bool {
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
		gfErr := gf_core.ErrorCreate("failed to render the landing_page template",
			"template_render_error",
			map[string]interface{}{},
			err, "gf_landing_page", pRuntimeSys)
		return "", gfErr
	}

	templateRenderedStr := buff.String()
	return templateRenderedStr, nil	
}

//-------------------------------------------------

func templatesLoad(pTemplatesPathsMap map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) (*gfTemplates, *gf_core.GFerror) {

	mainTemplateFilepathStr := pTemplatesPathsMap["gf_landing_page"]

	template, subtemplatesNamesLst, gf_err := gf_core.TemplatesLoad(mainTemplateFilepathStr,
		pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	gfTemplates := &gfTemplates{
		template:             template,
		subtemplatesNamesLst: subtemplatesNamesLst,
	}
	return gfTemplates, nil
}