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
	"errors"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func exprEval(pExpr interface{},
    pState     *GFstate,
    pExternAPI GFexternAPI) (interface{}, error) {

    symbols := getSymbolsAndConstants()

    //-------------
    // NUMBER_FLOAT
    if valF, ok := pExpr.(float64); ok {
        return valF, nil

    //-------------
    // NUMBER_INT
    } else if valInt, ok := pExpr.(int); ok {
        return valInt, nil
    
    //-------------

    } else if exprStr, ok := pExpr.(string); ok {

        if isVar(exprStr) {

            //-------------
            // VAR_REFERENCE            

            // // for now only system-defined vars are available, no user-defined vars yet.
            // if !gf_core.ListContainsStr(exprStr, symbols.SystemVarsLst) {
            //     return nil, errors.New(fmt.Sprintf("variable operand %s is not one of the system defined vars %s",
            //         exprStr,
            //         symbols.SystemVarsLst))
            // }

            varValue, err := varEval(exprStr, pState)
            if err != nil {
                return nil, err
            }

            return varValue.Val, nil

            //-------------
        }

    } else if exprLst, ok := pExpr.(GFexpr); ok {
        
        // SUB_EXPRESSION        
        isStrBool, arithmeticOpStr := gf_core.CastToStr(exprLst[0])

        //-------------
        // ARITHMETIC_OPERATION
        if isStrBool && gf_core.MapHasKey(symbols.ArithmeticOpsMap, arithmeticOpStr) {
            val, err := arithmeticEval(exprLst, pState, pExternAPI)
            if err != nil {
                return nil, err
            }
            
            valF := *val
            return valF, nil
        
        //-------------

        } else if isSysFunc(exprLst) {

            //-------------
            // SYSTEM_FUNCTION
            // these are mosly functions that are in turn invoking external_api functions.

            val, err := sysFuncEval(exprLst, pState, pExternAPI)
            if err != nil {
                return nil, err
            }
            return val, nil

            //-------------
        }
    }
    return nil, nil
}

//-------------------------------------------------
// ARITHMETIC_EVALUATION

func arithmeticEval(pExprLst []interface{},
    pState     *GFstate,
    pExternAPI GFexternAPI) (*float64, error) {

    symbols := getSymbolsAndConstants()

    if len(pExprLst) != 3 {
        return nil, errors.New(fmt.Sprintf("arithmetic expression %s has to be of length 3", pExprLst))
    }

    firstElementIsStrBool, arithmeticOpStr := gf_core.CastToStr(pExprLst[0])

    if !(firstElementIsStrBool && gf_core.MapHasKey(symbols.ArithmeticOpsMap, arithmeticOpStr)) {
        return nil, errors.New(fmt.Sprintf("arithmetic op %s is not supported", arithmeticOpStr))
    }

    opStr    := pExprLst[0].(string)
    operand1 := pExprLst[1]
    operand2 := pExprLst[2]

    //-------------------------------------------------
    evalOpFunc := func(pOperand interface{}) (interface{}, error) {
        
        var operand interface{}

        // SUB_EXPRESSION
        if subExprLst, ok := pOperand.(GFexpr); ok {

            // system_function sub-expression
            if isSysFunc(subExprLst) {
                subResult, err := sysFuncEval(subExprLst, pState, pExternAPI)
                if err != nil {
                    return nil, err
                }
                operand = subResult
            } else {

                // arithmetic sub-expression

                subResult, err := arithmeticEval(subExprLst, pState, pExternAPI)
                if err != nil {
                    return nil, err
                }
                operand = *subResult
            }

        } else if operandStr, ok := pOperand.(string); ok {
            
            // VARIABLE
            if strings.HasPrefix(operandStr, "$") {
                varVal, err := varEval(operandStr, pState)
                if err != nil {
                    return nil, err
                }
                operand = varVal.Val
            } else {
                return nil, errors.New("operator is a string but not a variable reference with '$'")
            }

        } else if _, ok := pOperand.(float64); ok {
            operand = pOperand
        
        } else if opInt, ok := pOperand.(int); ok {

            // all arithmetic ops are done as floats, so cast int to float
            operand = float64(opInt)

        } else {
            
            // if operand is not a subexpression, var reference, or number, its not valid
            return nil, errors.New(fmt.Sprintf("operand %s is not a subexpression|var_reference|number", pOperand))
        }
        return operand, nil
    }

    //-------------------------------------------------

    op1, err := evalOpFunc(operand1)
    if err != nil {
        return nil, err
    }
    
    op2, err := evalOpFunc(operand2)
    if err != nil {
        return nil, err
    }

    // EVALUATE
    resultF := symbols.ArithmeticOpsMap[opStr](op1.(float64), op2.(float64))
    return &resultF, nil
}

//-------------------------------------------------
// VARIABLES

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