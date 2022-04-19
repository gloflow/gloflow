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

package gf_rl

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GFqTable map[string]float64
type GFrlRuntime struct {
	QtableMap GFqTable
}

type GFrlActionDef struct {
	NameStr      string  // name of a particular action
	ProbabilityF float64 // probability of an action occuring
}

//-------------------------------------------------
func Init(pActionsDefsLst []GFrlActionDef,
	pRuntimeSys *gf_core.Runtime_sys) {









	env := EnvInit(pRuntimeSys)
	fmt.Println(env)


}

//-------------------------------------------------
func Train() {
	


	episodesNumInt := 100

	// Q_TABLE - quality table.
	//           stores quality of action by state.
	qTableMap := QtableCreate()

	qRuntime := &GFrlRuntime{
		QtableMap: qTableMap,
	}

	for i := 1; i <= episodesNumInt; i++ {

		QtableGetVal(qRuntime)

	}


}

//-------------------------------------------------
func epsilonGreedy() {



}

//-------------------------------------------------
// Q_TABLE
//-------------------------------------------------
func QtableCreate() map[string]float64 {
	qTableMap := GFqTable(map[string]float64{})
	return qTableMap
}
func QtableGetVal(pRuntime *GFrlRuntime) float64 {


	return 0.0
}
func QtableGetMaxVal(pRuntime *GFrlRuntime) float64 {
	return 0.0
}
func QtableSetVal(pVal float64,
	pRuntime *GFrlRuntime) {

}