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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
type GFhomeViz struct {
	Vstr               string             `bson:"v_str"` // schema_version
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	IDstr              gf_core.GF_ID      `bson:"id_str"`
	DeletedBool        bool               `bson:"deleted_bool"`
	CreationUNIXtimeF  float64            `bson:"creation_unix_time_f"`

	OwnerUserIDstr     gf_core.GF_ID                 `bson:"owner_user_id_str"`
	ComponentsMap      map[string]GFhomeVizComponent `bson:"components_map"`
}

type GFhomeVizComponent struct {
	NameStr            string `bson:"name_str"     json:"name_str"`
	ScreenXint         int64  `bson:"screen_x_int" json:"screen_x_int"`
	ScreenYint         int64  `bson:"screen_y_int" json:"screen_y_int"`
	ColorBackgroundStr string `bson:"color_background_str"`
}

//------------------------------------------------
// RENDER_DASHBOARD
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

//---------------------------------------------------
func homeVizCreateID(pUserIdentifierStr string,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		pUserIdentifierStr,
	}
	gfIDstr := gf_core.ID__create(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfIDstr
}