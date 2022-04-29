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
	"fmt"
	"time"
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

type GFvizPropsUpdateInput struct {
	UserIDstr        gf_core.GF_ID
	ComponentNameStr string
	ScreenXint       int64
	ScreenYint       int64
}

//------------------------------------------------
// VIZ_PROPS_UPDATE
func PipelineVizPropsUpdate(pInput *GFvizPropsUpdateInput,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GF_error {
	
	homeVizExisting, gfErr := DBgetHomeViz(pInput.UserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	// update component viz properties
	if pInput.ComponentNameStr != "" {
		foundBool := false
		for componentNameStr, vizComponent := range homeVizExisting.ComponentsMap {

			if componentNameStr == pInput.ComponentNameStr {
				foundBool = true
				vizComponent.ScreenXint = pInput.ScreenXint
				vizComponent.ScreenYint = pInput.ScreenYint
			}
		}

		// if the component is not already in the list of components 
		// then insert it as new.
		// it wouldnt be present if its viz properties havent been 
		// customized yet
		if !foundBool {
			newComponent := GFhomeVizComponent{
				NameStr:    pInput.ComponentNameStr,
				ScreenXint: pInput.ScreenXint,
				ScreenYint: pInput.ScreenYint,
			}
			homeVizExisting.ComponentsMap[pInput.ComponentNameStr] = newComponent
		}

		// DB
		gfErr = DBupdateHomeVizComponents(pInput.UserIDstr,
			homeVizExisting.ComponentsMap,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	return nil
}

//------------------------------------------------
// VIZ_PROPS_CREATE
func PipelineVizPropsCreate(pUserIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (*GFhomeViz, *gf_core.GF_error) {
		


	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	userIdentifierStr := string(pUserIDstr)
	IDstr := homeVizCreateID(userIdentifierStr,
		creationUNIXtimeF)

	homeViz := &GFhomeViz{
		Vstr:               "0",
		IDstr:              IDstr,
		CreationUNIXtimeF:  creationUNIXtimeF,
		OwnerUserIDstr:     pUserIDstr,
		ComponentsMap:      map[string]GFhomeVizComponent{},
	}

	// DB
	gfErr := DBcreateHomeViz(homeViz, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return homeViz, nil
}

//------------------------------------------------
// VIZ_PROPS_GET
func PipelineVizPropsGet(pUserIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (*GFhomeViz, *gf_core.GF_error) {
	

	homeVizExisting, gfErr := DBgetHomeViz(pUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// NONE_FOUND
	var homeViz *GFhomeViz
	if homeVizExisting == nil {

		fmt.Println("no home_viz found for user, creating new...")

		// CREATE
		homeVizNew, gfErr := PipelineVizPropsCreate(pUserIDstr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		homeViz = homeVizNew
	} else {
		homeViz = homeVizExisting
	}

	return homeViz, nil
}