
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
    "errors"
    "fmt"
    "strings"
    "math"
    "math/rand"
    "github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// EVALUATION
//-------------------------------------------------

func exprEval(pExpr interface{}, pState *GFstate) (interface{}, error) {

    symbols := getSymbolsAndConstants()

    // NUMBER
    if valF, ok := pExpr.(float64); ok {
        return valF, nil

    } else if exprStr, ok := pExpr.(string); ok {

        if strings.HasPrefix(exprStr, "$") {
            // VAR_REFERENCE

            // for now only system-defined vars are available, no user-defined vars yet.
            if !gf_core.ListContainsStr(exprStr, symbols.SystemVarsLst) {
                return nil, errors.New(fmt.Sprintf("variable operand %s is not one of the system defined vars %s",
                    exprStr,
                    symbols.SystemVarsLst))
            }

            varValue, err := varEval(exprStr, pState)
            if err != nil {
                return nil, err
            }

            return varValue, nil
        }

    } else if _, ok := pExpr.([]interface{}); ok {
        
        // SUB_EXPRESSION
        
        exprLst := pExpr.([]interface{})
        firstElementIsStrBool, arithmeticOpStr := gf_core.CastToStr(exprLst[0])

        //-------------
        // ARITHMETIC_OPERATION
        if firstElementIsStrBool && gf_core.MapHasKey(symbols.ArithmeticOpsMap, arithmeticOpStr) {
            val, err := arithmeticEval(exprLst, pState)
            if err != nil {
                return nil, err
            }
            
            return val, nil
        
        //-------------

        } else if isSysFunc(exprLst) {

            // SYSTEM_FUNCTION
            val, err := sysFuncEval(exprLst)
            if err != nil {
                return nil, err
            }
            return val, nil
        }
    }
    return nil, nil
}

//-------------------------------------------------
// ARITHMETIC_EVALUATION

func arithmeticEval(pExprLst []interface{}, pState *GFstate) (*float64, error) {

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
        if subExprLst, ok := pOperand.([]interface{}); ok {

            // system_function sub-expression
            if isSysFunc(subExprLst) {
                subResult, err := sysFuncEval(subExprLst)
                if err != nil {
                    return nil, err
                }
                operand = subResult
            } else {

                // arithmetic sub-expression

                subResult, err := arithmeticEval(subExprLst, pState)
                if err != nil {
                    return nil, err
                }
                operand = subResult
            }

        } else if operandStr, ok := pOperand.(string); ok {
            
            // VARIABLE
            if strings.HasPrefix(operandStr, "$") {
                varVal, err := varEval(operandStr, pState)
                if err != nil {
                    return nil, err
                }
                operand = varVal
            } else {
                return nil, errors.New("operator is a string but not a variable reference with '$'")
            }

        } else {

            // NUMBER
            // if operand is not a var reference, it has to be a number
            if _, ok := pOperand.(float64); !ok {
                return nil, errors.New(fmt.Sprintf("operand %s is not a number", pOperand))
            }
            operand = pOperand
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
// SYSTEM_FUNCTIONS

func isSysFunc(pExprLst []interface{}) bool {
    sysFunsLst := getSymbolsAndConstants().SystemFunctionsLst
    firstElementIsStrBool, firstElementStr := gf_core.CastToStr(pExprLst[0])
    if firstElementIsStrBool && gf_core.ListContainsStr(firstElementStr, sysFunsLst) {
        return true
    }
    return false
}

func sysFuncEval(pExprLst []interface{}) (interface{}, error) {
    
    funcNameStr := pExprLst[0].(string)
    argsLst     := pExprLst[1].([]interface{})
    
    var val interface{}
    if funcNameStr == "rand" {
        if len(argsLst) != 2 {
            return nil, errors.New("'rand' system function only takes 2 argument")
        }

        randomRangeMinF := argsLst[0].(float64)
        randomRangeMaxF := argsLst[1].(float64)
        valF := rand.Float64()*(randomRangeMaxF - randomRangeMinF) + randomRangeMinF
        val = interface{}(valF)
    }
    return val, nil
}

//-------------------------------------------------
// VARIABLES

func varEval(pVarStr string, pState *GFstate) (interface{}, error) {

    if !strings.HasPrefix(pVarStr, "$") {
        return nil, errors.New(fmt.Sprintf("variable string %s has no '$' prefixed", pVarStr))
    }
    varValue := pState.VarsMap[pVarStr]
    return varValue, nil
}

//-------------------------------------------------
// RULES
//-------------------------------------------------

func getRulesAllNames(pRulesDefsMap GFruleDefs) []string {
    rulesNamesLst := []string{}
    for k, _ := range pRulesDefsMap {
        rulesNamesLst = append(rulesNamesLst, k)
    }
    return rulesNamesLst
}

func ruleGetItersNum(pState *GFstate) int {
    return pState.RulesItersNumStackLst[len(pState.RulesItersNumStackLst)-1]
}

func ruleGetName(pState *GFstate) (string, error) {
    if len(pState.RulesNamesStackLst) == 0 {
        return "", errors.New("no rules in rules stack")
    }

    ruleNameStr := pState.RulesNamesStackLst[len(pState.RulesNamesStackLst)-1]
    return ruleNameStr, nil
}

//-------------------------------------------------

func incrementItersNum(pState *GFstate) int {
    newRuleItersNumInt := ruleGetItersNum(pState) + 1

    pState.RulesItersNumStackLst[len(pState.RulesItersNumStackLst)-1] = newRuleItersNumInt
    pState.VarsMap["$i"] = newRuleItersNumInt
    
    return newRuleItersNumInt
}

func addNewItersNumState(pState *GFstate) {
    pState.RulesItersNumStackLst = append(pState.RulesItersNumStackLst, 0)
    pState.VarsMap["$i"] = 0
}

//-------------------------------------------------
// called when one rule exits (finishes executing) and returns execution
// to its parent rule (not when the same rule recurses into itself).

func restorePreviousRulesItersNum(pState *GFstate) {

    _, newRulesItersNumStackLst := gf_core.ListPop(pState.RulesItersNumStackLst)
    pState.RulesItersNumStackLst = newRulesItersNumStackLst

    // reinitialize $i to the parents number of iterations
    pState.VarsMap["$i"] = pState.RulesItersNumStackLst[len(pState.RulesItersNumStackLst)-1]
}

//-------------------------------------------------

func pickRuleRandomDef(pRuleNameStr string,
    pRulesDefsMap GFruleDefs) (*GFruleDef, []interface{}) {

    ruleDefsLst           := pRulesDefsMap[pRuleNameStr]
    ruleDefsNumInt        := len(ruleDefsLst)
    randomRuleDefIndexInt := int(math.Floor(rand.Float64() * float64(ruleDefsNumInt)))
    ruleDef               := ruleDefsLst[randomRuleDefIndexInt]
    ruleExpressionsLst    := ruleDef.ExpressionsLst
    return ruleDef, ruleExpressionsLst
}

//-------------------------------------------------
// CLONING
//-------------------------------------------------

func cloneExpr(pExprLst []interface{}) []interface{} {
    copyLst := make([]interface{}, len(pExprLst))
    copy(copyLst, pExprLst)
    return copyLst
}

//-------------------------------------------------

func cloneExprNtimes(pExprLst []interface{}, pNint int) []interface{} {
    
    clonesLst := []interface{}{}
    for i := 0; i < pNint; i++ {

        clonedExprLst := cloneExpr(pExprLst)
        clonesLst = append(clonesLst, clonedExprLst)
    }
    return clonesLst
}

//-------------------------------------------------

// clone program vars. this still assumes simple variables with primitive values.
func cloneVars(pVarsMap map[string]interface{}) map[string]interface{} {
    cloneMap := map[string]interface{}{}
    for k, v := range pVarsMap {
        cloneMap[k] = v
    }
    return cloneMap
}

func cloneItersNumStack(pRulesItersNumStackLst []int) []int {
    copyLst := make([]int, len(pRulesItersNumStackLst))
    copy(copyLst, pRulesItersNumStackLst)
    return copyLst
}

func cloneRulesNamesStack(pRulesNamesStackLst []string) []string {
    copyLst := make([]string, len(pRulesNamesStackLst))
    copy(copyLst, pRulesNamesStackLst)
    return copyLst
}

// clone program vars. this still assumes simple variables with primitive values.
func cloneAnimations(pAnimationsMap map[string]interface{}) map[string]interface{} {
    cloneMap := map[string]interface{}{}
    for k, v := range pAnimationsMap {
        cloneMap[k] = v
    }
    return cloneMap
}

//-------------------------------------------------
// SYMBOLS
//-------------------------------------------------

func getSymbolsAndConstants() *GFsymbols {
    
    ruleLevelMaxInt := 250
    systemRulesLst := []string{
        "cube",
        "sphere",
        "line",
    }

    predefinedPropertiesLst := []string{
        "x",  // x-coordinate
        "y",  // y-coordinate
        "z",  // z-coordinate
        "rx", // x-rotation
        "ry", // y-rotation
        "rz", // z-rotation
        "sx", // x-scale
        "sy", // y-scale
        "sz", // z-scale
        "cr", // red-channel-color
        "cg", // green-channel-color
        "cb", // blue-channel-color
    }
    logicOperatorsMap := map[string]func(float64, float64) bool {
        "==": func(p1 float64, p2 float64) bool {return p1 == p2},
        "!=": func(p1 float64, p2 float64) bool {return p1 != p2},
        "<":  func(p1 float64, p2 float64) bool {return p1 < p2},
        ">":  func(p1 float64, p2 float64) bool {return p1 > p2},
        "<=": func(p1 float64, p2 float64) bool {return p1 <= p2},
        ">=": func(p1 float64, p2 float64) bool {return p1 >= p2},
    }
    arithmeticOpsMap := map[string]func(float64, float64) float64{
        "+": func(p1 float64, p2 float64) float64 {return p1 + p2},
        "-": func(p1 float64, p2 float64) float64 {return p1 - p2},
        "*": func(p1 float64, p2 float64) float64 {return p1 * p2},
        "/": func(p1 float64, p2 float64) float64 {return p1 / p2},
        "%": func(p1 float64, p2 float64) float64 {return float64(int(p1) % int(p2))},
    }
    systemVarsLst := []string{
        "$i", // current rule iteration
    }
    systemFunctionsLst := []string{
        "rand", // random number generator
    }

    symbols := &GFsymbols{
        RuleLevelMaxInt:         ruleLevelMaxInt,
        SystemRulesLst:          systemRulesLst,
        PredefinedPropertiesLst: predefinedPropertiesLst,
        LogicOperatorsMap:       logicOperatorsMap,
        ArithmeticOpsMap:        arithmeticOpsMap,
        SystemVarsLst:           systemVarsLst,
        SystemFunctionsLst:      systemFunctionsLst,
    }
    return symbols
}

//-------------------------------------------------
// CHECKS
//-------------------------------------------------

func checkArithmeticOpExists(pOpToCheckStr string) bool {
    for opStr, _ := range getSymbolsAndConstants().ArithmeticOpsMap {
        if pOpToCheckStr == opStr {
            return true
        }
    }
    return false
}

//-------------------------------------------------
// CHECKS
//-------------------------------------------------

func castToFloat(pN interface{}) float64 {
    if nInt, ok := pN.(int); ok {
        return float64(nInt)
    }
    if nF, ok := pN.(float64); ok {
        return nF
    }
    panic(fmt.Sprintf("number %s is not a int or float64", pN))
    return 0.0
}