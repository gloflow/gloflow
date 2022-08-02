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
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//---------------------------------------------------
func inputForVizPropsUpdate(pReq *http.Request,
	pResp       http.ResponseWriter,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (*GFvizPropsUpdateInput, *gf_core.GFerror) {

	userIDstr, _ := gf_identity_core.GetUserIDfromCtx(pCtx)

	inputMap, gfErr := gf_core.HTTPgetInput(pResp, pReq, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	var componentNameStr string
	if valStr, ok := inputMap["component_name_str"]; ok {
		componentNameStr = valStr.(string)
	}

	var screenXint int64
	if valStr, ok := inputMap["props_change_map"].(map[string]interface{})["screen_x_int"]; ok {
		screenXint = int64(valStr.(float64))
	}

	var screenYint int64
	if valStr, ok := inputMap["props_change_map"].(map[string]interface{})["screen_y_int"]; ok {
		screenYint = int64(valStr.(float64))
	}

	input := &GFvizPropsUpdateInput{
		UserIDstr:        userIDstr,
		ComponentNameStr: componentNameStr,
		ScreenXint:       screenXint,
		ScreenYint:       screenYint,
	}


	return input, nil
}