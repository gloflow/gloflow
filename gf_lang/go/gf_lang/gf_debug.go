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
	"github.com/davecgh/go-spew/spew"
)

type GFprogramDebug struct {
	RulesCallsCounterMap map[string]int
	StateHistoryLst      []*GFstate
}

//-------------------------------------------------

func debugInit() *GFprogramDebug {
	programDebug := &GFprogramDebug{
		RulesCallsCounterMap: map[string]int{},
		StateHistoryLst:      []*GFstate{},
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