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

func execVarExpr(pExpr GFexpr,
	pState     *GFstate,
	pExternAPI GFexternAPI) (interface{}, error) {


	varNameStr := pExpr[0].(string)
	valUnevaluated := pExpr[1]

	fmt.Printf("variable %s\n", varNameStr)

	// EVALUATE
	// variable assignments can for now be relativelly simple expressions,
	// numbers, other variable references,
	// and simple arithmetic and system_function expressions. 
	varVal, err := exprEval(valUnevaluated, pState, pExternAPI)
	if err != nil {
		return nil, err
	}




	spew.Dump(pState)

	expressionResult := varVal
	return expressionResult, nil
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
// CLONE_VARS

// clone program vars. this still assumes simple variables with primitive values.
func cloneVars(pVarsMap map[string]*GFvariableVal) map[string]*GFvariableVal {
    cloneMap := map[string]*GFvariableVal{}
    for k, v := range pVarsMap {
        cloneMap[k] = v
    }
    return cloneMap
}