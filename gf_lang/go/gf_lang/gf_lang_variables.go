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
	"strings"
	"errors"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFvariableVal struct {
    NameStr string
    Val     interface{}
}

//-------------------------------------------------

func execVarAssignExpr(pExpr GFexpr,
	pState     *GFstate,
	pExternAPI GFexternAPI) error {


	varNameStr := pExpr[0].(string)
	valUnevaluated := pExpr[1]

	fmt.Printf("variable %s\n", varNameStr)
	spew.Dump(pState.VarsMap)

	// EVALUATE
	// variable assignments can for now be relativelly simple expressions,
	// numbers, other variable references,
	// and simple arithmetic and system_function expressions. 
	varNewVal, _, err := exprEvalSimple(valUnevaluated, pState, pExternAPI)
	if err != nil {
		return err
	}

	// check if the variable is already created. in that case
	// assign it the new value
	if varVal, ok := pState.VarsMap[varNameStr]; ok {
		varVal.Val = varNewVal
	} else {

		// if var has not already been created, then create it
		// and assign it the newly evaluated value.
		_ = createVariable(varNameStr, varNewVal, pState)
	}

	return nil
}

//-------------------------------------------------

func varEval(pVarStr string, pState *GFstate) (*GFvariableVal, error) {

    // VERIFY
    err := varVerify(pVarStr)
    if err != nil {
        return nil, err
    }

    // read value
    varValue := pState.VarsMap[pVarStr]

    return varValue, nil
}

//-------------------------------------------------
// CREATE

func createVariable(pVariableNameStr string,
	pInitVal interface{},
    pState   *GFstate) *GFvariableVal {

    variable := &GFvariableVal{
        NameStr: pVariableNameStr,
        Val:     pInitVal,
    }
    pState.VarsMap[pVariableNameStr] = variable
    return nil
}

//-------------------------------------------------

func varVerify(pVarStr string) error {
    if !strings.HasPrefix(pVarStr, "$") {
        return errors.New(fmt.Sprintf("variable string %s has no '$' prefixed", pVarStr))
    }
    return nil
}

//-------------------------------------------------

func isVar(pVarStr string) bool {
    if strings.HasPrefix(pVarStr, "$") {
        return true
    }
    return false
}

//-------------------------------------------------
// CLONE

// clone program vars. this still assumes simple variables with primitive values.
func cloneVars(pVarsMap map[string]*GFvariableVal) map[string]*GFvariableVal {
    cloneMap := map[string]*GFvariableVal{}
    for k, v := range pVarsMap {
        cloneMap[k] = v
    }
    return cloneMap
}