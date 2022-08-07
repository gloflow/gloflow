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
	"math/rand"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GFaction string
type GFqTable map[string]float64
type GFrewardFun func(pState *GFstate, pActionStr GFaction) int

type GFrlRuntime struct {
	
	Hyperparams GFrlHyperparams
	QtableMap   GFqTable

	RuntimeSys  *gf_core.RuntimeSys
}

type GFstate struct {
	
	// values of all the state-space dimensions for this particular state
	DimensionValsLst []interface{}

	// app-level data that can be associated with each state
	DataMap map[string]interface{}
}

type GFrlActionDef struct {
	NameStr      string  // name of a particular action
	ProbabilityF float64 // probability of an action occuring
}

type GFappInfo struct {
	NameStr        string
	ActionsDefsLst []GFrlActionDef
	RewardFun      GFrewardFun
	Hyperparams    GFrlHyperparams
}

type GFrlHyperparams struct {
	// probability (0.0<e<1.0) for epsilon-greedy strategy
	// (exploration/exploatation action probability)
	EpsilonF float64
	
	// learning rate (0<a<1)
	// extent to which Q-values are being updated in every iteration.
	AlphaLearningRateF float64
}

//-------------------------------------------------
func Init(pAppInfo *GFappInfo,
	pRuntimeSys *gf_core.RuntimeSys) {

	

	// Q_TABLE - quality table.
	//           stores quality of action by state.
	qTableMap := QtableCreate()

	rlRuntime := &GFrlRuntime{
		QtableMap:  qTableMap,
		RuntimeSys: pRuntimeSys,
	}





	env := EnvInit(rlRuntime)
	fmt.Println(env)
}

//-------------------------------------------------
func Train(pRuntime *GFrlRuntime) {
	


	episodesNumInt := 100

	

	for i := 1; i <= episodesNumInt; i++ {

		QtableGetVal(pRuntime)

	}


}

//-------------------------------------------------
// get action that should be taken given the current state
func GetAction(pState *GFstate,
	pRuntime *GFrlRuntime) {


	// rand.Float64() - gives rand number between 0.0 and 1.0
	if rand.Float64() < pRuntime.Hyperparams.EpsilonF {

		// explore, get a random action
	} else {

		// exploit

	}



}

//-------------------------------------------------
func AgentLearn(pRuntime *GFrlRuntime) {

	alphaF := pRuntime.Hyperparams.AlphaLearningRateF

	d := (1.0 - alphaF) 
	fmt.Println(d)	
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