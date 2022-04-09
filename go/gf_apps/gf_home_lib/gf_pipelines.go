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

package gf_home_lib

import (
	"text/template"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
type GFhomeViz struct {
	ColorBackgroundStr string               `bson:"color_background_str" json:"color_background_str"`
	ComponentsLst      []GFhomeVizComponent `bson:"components_lst"       json:"components_lst"`
}

type GFhomeVizComponent struct {
	ScreenXint int64 `bson:"screen_x_int" json:"screen_x_int"`
	ScreenYint int64 `bson:"screen_y_int" json:"screen_y_int"`
}

//------------------------------------------------
// VIZ_PROPS_GET
func PipelineVizPropsGet(pCtx context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (*GFhomeViz, *gf_core.GF_error) {



	return nil, nil
}

//------------------------------------------------
// VIZ_PROPS_UPDATE
func PipelineVizPropsUpdate(pCtx context.Context,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GF_error {
		

	return nil
}

//------------------------------------------------
func PipelineRenderDashboard(pTmpl *template.Template,
	pSubtemplatesNamesLst []string,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.Runtime_sys) (string, *gf_core.GF_error) {

	templateRenderedStr, gfErr := viewRenderTemplateDashboard(pTmpl,
		pSubtemplatesNamesLst,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return templateRenderedStr, nil
}