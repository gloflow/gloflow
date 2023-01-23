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
    "math"
    "math/rand"
    "github.com/gloflow/gloflow/go/gf_core"
)

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
    pState.VarsMap["$i"].Val = newRuleItersNumInt
    
    return newRuleItersNumInt
}

func addNewItersNumState(pState *GFstate) {
    pState.RulesItersNumStackLst = append(pState.RulesItersNumStackLst, 0)
    pState.VarsMap["$i"].Val = 0
}

//-------------------------------------------------
// called when one rule exits (finishes executing) and returns execution
// to its parent rule (not when the same rule recurses into itself).

func restorePreviousRulesItersNum(pState *GFstate) {

    _, newRulesItersNumStackLst := gf_core.ListPop(pState.RulesItersNumStackLst)
    pState.RulesItersNumStackLst = newRulesItersNumStackLst

    // reinitialize $i to the parents number of iterations
    pState.VarsMap["$i"].Val = pState.RulesItersNumStackLst[len(pState.RulesItersNumStackLst)-1]
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

func cloneExpr(pExprLst GFexpr) GFexpr {
    copyLst := make([]interface{}, len(pExprLst))
    copy(copyLst, pExprLst)
    return copyLst
}

//-------------------------------------------------

func cloneExprNtimes(pExprLst GFexpr, pNint int) GFexpr {
    
    clonesLst := []interface{}{}
    for i := 0; i < pNint; i++ {

        clonedExprLst := cloneExpr(pExprLst)
        clonesLst = append(clonesLst, clonedExprLst)
    }
    return GFexpr(clonesLst)
}

//-------------------------------------------------

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
    systemFunctionsLst := getSysFunctionNames()

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

func checkIsArithmeticOp(pOpToCheckStr string) bool {
    for opStr, _ := range getSymbolsAndConstants().ArithmeticOpsMap {
        if pOpToCheckStr == opStr {
            return true
        }
    }
    return false
}

//-------------------------------------------------

func checkIsReturnExpr(pExprLst GFexpr) bool {
    if pExprLst[0].(string) == "return" {
        return true
    }
    return false
}

//-------------------------------------------------
// CASTS
//-------------------------------------------------

func CastToExpr(pExprLst []interface{}) GFexpr {

    exprLst := []interface{}{}
    for _, el := range pExprLst {
        if eLst, ok := el.([]interface{}); ok {

            exprCastedLst := CastToExpr(eLst)
            exprLst = append(exprLst, exprCastedLst)
        } else {
            // not an expression, so no need to cast
            exprLst = append(exprLst, el)
        }
    }
    return GFexpr(exprLst)
}

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

func castToInt(pN interface{}) int {
    if nInt, ok := pN.(int); ok {
        return nInt
    }
    if nF, ok := pN.(float64); ok {
        return int(nF)
    }
    panic(fmt.Sprintf("number %s is not a int or float64", pN))
    return 0
}