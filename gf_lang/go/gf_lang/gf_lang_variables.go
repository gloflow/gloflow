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
    "reflect"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFvariableVal struct {
    NameStr string
    TypeStr string // "string"|"number"|"list"
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

        //------------------------------------
        // UPDATE_EXISTING_VARIABLE

        newValTypeStr, err := inferVarType(varNewVal)
        if err != nil {
            return err
        }

        if varVal.TypeStr != newValTypeStr {
            return errors.New(fmt.Sprintf("trying to assign to variable %s of type %s a value %s of incompatible type %s",
                varVal.TypeStr,
                varNameStr,
                varNewVal,
                newValTypeStr))
        }

        if newValTypeStr == "list" {
            // gurantee the list/[]interface{} type,
            // since all lists in the code are first cast to GFexpr type by the preprocessor.
            varNewValLst := varNewVal.([]interface{})
            varVal.Val = varNewValLst
        } else {
            varVal.Val = varNewVal
        }

        //------------------------------------
		
	} else {

        //------------------------------------
        // CREATE_NEW_VARIABLE
		// if var has not already been created, then create it
		// and assign it the newly evaluated value.
		_, err = createVariable(varNameStr, varNewVal, pState)
        if err != nil {
            return err
        }

        //------------------------------------
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
    pState   *GFstate) (*GFvariableVal, error) {
    
    fmt.Printf("create new variable - %s | %s\n", pVariableNameStr, pInitVal)

    varTypeStr, err := inferVarType(pInitVal)
    if err != nil {
        return nil, err
    }

    variable := &GFvariableVal{
        NameStr: pVariableNameStr,
        TypeStr: varTypeStr,
    }

    if varTypeStr == "list" {

        // gurantee the list/[]interface{} type,
        // since all lists in the code are first cast to GFexpr type by the preprocessor.
        varValLst := pInitVal.([]interface{})
        variable.Val = varValLst
    } else {
        variable.Val = pInitVal
    }

    pState.VarsMap[pVariableNameStr] = variable
    return variable, nil
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

//-------------------------------------------------

func inferVarType(pVal interface{}) (string, error) {
    
    switch pVal.(type) {
    case int:
        return "int", nil
    case float64:
        return "float", nil
    case string:
        return "string", nil
    case []interface{}:
        return "list", nil
    case GFexpr:
        return "list", nil
    case map[string]interface{}:
        return "map", nil
    default:
        unknownTypeStr := reflect.TypeOf(pVal)
        return "", errors.New(fmt.Sprintf("variable of unsupported type - %s", unknownTypeStr))
    }
    return "", nil
}