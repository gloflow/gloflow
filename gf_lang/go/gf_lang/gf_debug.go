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

package gf_lang

import (
	"fmt"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFprogramDebug struct {
	RulesCallsCounterMap map[string]int
	StateHistoryLst      []*GFstate

	//----------------
	// OUTPUT
	
	OutputEntitiesMap map[gf_core.GF_ID]*GFentityOutput
	
	// this is the full output built up over program execution
	// that stores all extern API calls.
	// this includes both the entities and state operations
	OutputLst []interface{}

	//----------------
}

//-------------------------------------------------

type GFstateChangeOutput struct {
	IDstr   gf_core.GF_ID  `json:"id_str"`
	Change  *GFstateChange `json:"change_map"`
}

type GFentityOutput struct {
	IDstr   gf_core.GF_ID  `json:"id_str"`
	TypeStr string         `json:"type_str"`
	Props   *GFentityProps `json:"props_map"`
}

type GFentityProps struct {
	Xf float64 `json:"x_f"`
	Yf float64 `json:"y_f"`
	Zf float64 `json:"z_f"`
	RotationXf  float64 `json:"rotation_x_f"`
	RotationYf  float64 `json:"rotation_y_f"`
	RotationZf  float64 `json:"rotation_z_f"`
	ScaleXf     float64 `json:"scale_x_f"`
	ScaleYf     float64 `json:"scale_y_f"`
	ScaleZf     float64 `json:"scale_z_f"`
	ColorRedF   float64 `json:"color_red_f"`
	ColorGreenF float64 `json:"color_green_f"`
	ColorBlueF  float64 `json:"color_blue_f"`
}

//-------------------------------------------------

func addEntityToOutput(pTypeStr string,
	pProps *GFentityProps,
	pDebug *GFprogramDebug) {

	//------------------------
	// ID
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	fieldsForIDlst := []string{
	}
	IDstr := gf_core.IDcreate(fieldsForIDlst,
		creationUNIXtimeF)
		
	//------------------------

	entity := &GFentityOutput{
		IDstr:   IDstr,
		TypeStr: pTypeStr,
		Props:   pProps,
	}

	pDebug.OutputEntitiesMap[IDstr] = entity
	pDebug.OutputLst = append(pDebug.OutputLst, entity)
}

//-------------------------------------------------

func addExternStateChange(pStateChange *GFstateChange,
	pDebug *GFprogramDebug) {

	//------------------------
	// ID
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	fieldsForIDlst := []string{
	}
	IDstr := gf_core.IDcreate(fieldsForIDlst,
		creationUNIXtimeF)
		
	//------------------------
	
	stateChange := &GFstateChangeOutput{
		IDstr:   IDstr,
		Change:  pStateChange,
	}

	pDebug.OutputLst = append(pDebug.OutputLst, stateChange)
}

//-------------------------------------------------

func debugInit() *GFprogramDebug {
	programDebug := &GFprogramDebug{
		RulesCallsCounterMap: map[string]int{},
		StateHistoryLst:      []*GFstate{},
		OutputEntitiesMap:    map[gf_core.GF_ID]*GFentityOutput{},
		OutputLst:            []interface{}{},

	}
	return programDebug
}

//-------------------------------------------------

func viewDebug(pProgramASTlst GFexpr,
	pProgramsDebugLst []*GFprogramDebug) {

	for i, programDebug := range pProgramsDebugLst {

		// "1+i" - skip over the "lang_v" statement
		fmt.Println("PROGRAM ", i, ">>>>>>>>>>>>>", pProgramASTlst[1+i])
		
		for j, state := range programDebug.StateHistoryLst {
			fmt.Println("state history item ", j, "colors", state.ColorRedF, state.ColorGreenF, state.ColorBlueF)
			// spew.Dump(state)
		}

		spew.Dump(programDebug.RulesCallsCounterMap)
	}
}