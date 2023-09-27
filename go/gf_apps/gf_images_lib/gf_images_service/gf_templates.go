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

package gf_images_service

import (
	// "fmt"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type gfTemplates struct {
	imagesViewTmpl                 *template.Template
	imagesViewSubtemplatesNamesLst []string
}

//-------------------------------------------------

func templateLoad(pTemplatesPathsMap map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) (*gfTemplates, *gf_core.GFerror) {

	imagesViewTemplateFilepathStr := pTemplatesPathsMap["gf_images_view"]

	imagesViewTmpl, subtemplatesNamesSst, gf_err := gf_core.TemplatesLoad(imagesViewTemplateFilepathStr,
		pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	gfTemplates := &gfTemplates{
		imagesViewTmpl:                 imagesViewTmpl,
		imagesViewSubtemplatesNamesLst: subtemplatesNamesSst,
	}
	return gfTemplates, nil
}