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

package gf_lang

import (
    // "fmt"
    "errors"
    "reflect"
    "github.com/gloflow/gloflow/go/gf_core"
    // "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// FAMILY - state family is a new root state and all its children
//          that are merged into each other.
//          they have their own coordinate system origin.
//-------------------------------------------------

func statePushFamily(pNewState *GFstate, pStateFamilyStackLst []*GFstate) {
    pStateFamilyStackLst = append(pStateFamilyStackLst, pNewState)
}

func statePopFamily(pStateFamilyStackLst []*GFstate) (*GFstate, error) {

    // this is the last state of the last state_family
    lastFamilyState, newStateFamilyStackLst := gf_core.ListPop(pStateFamilyStackLst)

    if lastFamilyState == nil {
        return nil, errors.New("family_state_stack has no more states to pop")
    }

    pStateFamilyStackLst = newStateFamilyStackLst
    return lastFamilyState, nil
}

//-------------------------------------------------
// create new state family, which is state like all others except
// it begins a new geometric space (a new coord system) but inherits some of the
// iteration counters and sub-stacks.

func stateCreateNewFamily(pStateParent *GFstate) *GFstate {

    state := stateGetEmpty()

    state.RulesNamesStackLst    = pStateParent.RulesNamesStackLst
    state.VarsMap               = cloneVars(pStateParent.VarsMap) // clone
    state.ItersNumGlobalInt     = pStateParent.ItersNumGlobalInt
    state.RulesItersNumStackLst = cloneItersNumStack(pStateParent.RulesItersNumStackLst) // clone

    return state
}

//-------------------------------------------------
// VAR
//-------------------------------------------------

func stateMergeChild(pState *GFstate, pChildState *GFstate) (*GFstate, error) {
    
    if pChildState == nil {
        return nil, errors.New("supplied child_state_map is nil")
    }

    pState.Xf  = pChildState.Xf
    pState.Yf  = pChildState.Yf
    pState.Zf  = pChildState.Zf
    pState.RotationXf = pChildState.RotationXf
    pState.RotationYf = pChildState.RotationYf
    pState.RotationZf = pChildState.RotationZf
    pState.ScaleXf = pChildState.ScaleXf
    pState.ScaleYf = pChildState.ScaleYf
    pState.ScaleZf = pChildState.ScaleZf
    pState.ColorRedF   = pChildState.ColorRedF
    pState.ColorGreenF = pChildState.ColorGreenF
    pState.ColorBlueF  = pChildState.ColorBlueF
    pState.ItersMaxInt        = pChildState.ItersMaxInt
    pState.RulesNamesStackLst = pChildState.RulesNamesStackLst

    // rule iteration count ("$i") has to propagate up the expression tree,
    // from child expressions to parent expressions, as a part of the state.
    // however $i only travels up to the root of a particular rule.
    varVal, err := varEval("$i", pChildState)
    if err != nil {
        return nil, err
    }
    pState.VarsMap["$i"] = varVal

    // what is the global number of iteratios executed relative to the root state
    pState.ItersNumGlobalInt = pChildState.ItersNumGlobalInt

    pState.RulesItersNumStackLst = pChildState.RulesItersNumStackLst

    //----------------------
    // ANIMATIONS - are not merged from children, they can only be propagated
    //              down the execution tree, not up.
    // 
    //----------------------

    return pState, nil
}

//-------------------------------------------------

func stateCreateNew(pStateParent *GFstate) *GFstate {

    stateNew := stateGetEmpty()

    if pStateParent != nil {
        stateNew.Xf = pStateParent.Xf
        stateNew.Yf = pStateParent.Yf
        stateNew.Zf = pStateParent.Zf
        stateNew.RotationXf  = pStateParent.RotationXf
        stateNew.RotationYf  = pStateParent.RotationYf
        stateNew.RotationZf  = pStateParent.RotationZf
        stateNew.ScaleXf     = pStateParent.ScaleXf
        stateNew.ScaleYf     = pStateParent.ScaleYf
        stateNew.ScaleZf     = pStateParent.ScaleZf
        stateNew.ColorRedF   = pStateParent.ColorRedF
        stateNew.ColorGreenF = pStateParent.ColorGreenF
        stateNew.ColorBlueF  = pStateParent.ColorBlueF
        stateNew.ItersMaxInt        = pStateParent.ItersMaxInt
        stateNew.RulesNamesStackLst = cloneRulesNamesStack(pStateParent.RulesNamesStackLst) // clone
        stateNew.VarsMap            = cloneVars(pStateParent.VarsMap)                       // clone

        stateNew.ItersNumGlobalInt     = pStateParent.ItersNumGlobalInt
        stateNew.RulesItersNumStackLst = cloneItersNumStack(pStateParent.RulesItersNumStackLst) // clone
        stateNew.AnimationsActiveMap   = cloneAnimations(pStateParent.AnimationsActiveMap)      // clone
    }

    return stateNew
}

//-------------------------------------------------

func stateGetEmpty() *GFstate {
    
    /*
    stateMap := {
        "x":  0.0,
        "y":  0.0,
        "z":  0.0,
        "rx": 0.0,
        "ry": 0.0,
        "rz": 0.0,
        "sx": 1.0,
        "sy": 1.0,
        "sz": 1.0,
        "cr": 0.0,
        "cg": 0.0,
        "cb": 0.0,

        // global max number of iterations for any rule
        "iters_max": 250,

        // list of all rules that are executing
        "rules_names_stack_lst": ["root"],

        "vars_map": {
            "$i": 0,
        },

        // global iterations number for a particular root expression
        "iters_num_global_int": 0,

        // stack of iteration numbers for each rule as its entered
        "rules_iters_num_stack_lst": [0],

        // ANIMATIONS - map of animations that are currently active
        //              in a subexpression or its children.
        "animations_active_map": {}
    }
    */

    state := &GFstate{
        ScaleXf: 1.0,
        ScaleYf: 1.0,
        ScaleZf: 1.0,
        ItersMaxInt: 250,
        RulesNamesStackLst: []string{"root",},
        VarsMap: map[string]interface{}{
            "$i": 0,
        },

        RulesItersNumStackLst: []int{0,},
    }
    return state
}

//-------------------------------------------------

// mapping of state property names that are exposed to language programs to the
// state names in the compiler. used when reflecting.
func stateGetPropertyNamesInLangMap() map[string]string {
    return map[string]string{
        "x": "Xf",
        "y": "Yf",
        "z": "Zf",
        "rx": "RotationXf",
        "ry": "RotationYf",
        "rz": "RotationZf",
        "sx": "ScaleXf",
        "sy": "ScaleYf",
        "sz": "ScaleZf",
        "cr": "ColorRedF",
        "cg": "ColorGreenF",
        "cb": "ColorBlueF",
    }
}

//-------------------------------------------------

func statePropFloatIncrement(pState *GFstate,
    pPropertyInProgramNameStr string,
    pIncrementByF             float64) {

    propInternalNameStr := stateGetPropertyNamesInLangMap()[pPropertyInProgramNameStr]
    field        := reflect.ValueOf(pState).Elem().FieldByName(propInternalNameStr)
    fieldValF    := field.Float()
    fieldNewValF := fieldValF + pIncrementByF

    field.SetFloat(fieldNewValF)
}

func statePropGet(pState *GFstate,
    pPropertyInProgramNameStr string) float64 {

    propInternalNameStr := stateGetPropertyNamesInLangMap()[pPropertyInProgramNameStr]
    field     := reflect.ValueOf(pState).Elem().FieldByName(propInternalNameStr)
    fieldValF := field.Float()
    return fieldValF
}