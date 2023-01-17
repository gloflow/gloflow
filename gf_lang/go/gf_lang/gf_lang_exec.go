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
    "fmt"
    "errors"
    "strings"
    "github.com/gloflow/gloflow/go/gf_core"
    // "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func executeTree(pExpressionASTlst []interface{},
    pStateParent         *GFstate,
    pRulesDefsMap        GFruleDefs,
    pShaderDefsMap       map[string]interface{},
    pStateFamilyStackLst []*GFstate,
    pExternAPI           GFexternAPI) (*GFstate, interface{}, error) {
    
    symbols := getSymbolsAndConstants()
    
    //--------------------
    // STATE_NEW
    // IMPORTANT!! - on every tree descent a new independent state is constructed
    state := stateCreateNew(pStateParent)

    //--------------------

    // clone in case of mutations of expression
    expressionLst := cloneExpr(pExpressionASTlst)



    //-------------------------------------------------
    handleSubExpressionFun := func(pSubExprLst GFexpr) (interface{}, error) {
            
        // RECURSION
        childState, subExprResult, err := executeTree(pSubExprLst,
            state,
            pRulesDefsMap,
            pShaderDefsMap,
            pStateFamilyStackLst,
            pExternAPI)
        if err != nil {
            return nil, err
        }

        mergedState, err := stateMergeChild(state, childState)
        if err != nil {
            return nil, err
        }
        state = mergedState
        return subExprResult, nil
    }

    //-------------------------------------------------

    // iterate over each expression element
    for i:=0; i < len(expressionLst); i++ {

        element := expressionLst[i]
        elementIsStrBool, elementStr := gf_core.CastToStr(element)

        //------------------------------------
        // PROPERTY_MODIFIER
        if elementIsStrBool && gf_core.ListContainsStr(elementStr, symbols.PredefinedPropertiesLst) {

            propertyNameStr := elementStr
            modifier        := expressionLst[i+1]
            var modifierFinalValF float64

            // modifier is a number
            if modifierF, ok := modifier.(float64); ok {
                modifierFinalValF = modifierF

            } else if modifierLst, ok := modifier.([]interface{}); ok {

                // modifier is a sub-expression - needs to be evaluated

                subExprLst := modifierLst

                //--------------------
                // SYSTEM_FUNCTION
                if isSysFunc(subExprLst) {

                    result, err := sysFuncEval(subExprLst, state, pExternAPI)
                    if err != nil {
                        return nil, nil, err
                    }

                    modifierFinalValF = result.(float64)
                
                //--------------------

                } else {
                    resultF, err := arithmeticEval(subExprLst, state, pExternAPI)
                    if err != nil {
                        return nil, nil, err
                    }
                    modifierFinalValF = *resultF
                }
            }

            // IMPORTANT!! - incremental modification
            statePropFloatIncrement(state, propertyNameStr, modifierFinalValF)

            // if (propertyNameStr.startsWith("r"))
            //     console.log("rotation", propertyNameStr, modifierVal, state[propertyNameStr])

            i+=1 // fast-forward, modifiers can be listed sequentially in the same expression
            continue

        } else if subExprLst, ok := element.(GFexpr); ok {

            //------------------------------------
            // SUB_EXPRESSION

            subExprResult, err := handleSubExpressionFun(subExprLst)
            if err != nil {
                return nil, nil, err
            }

            // sub-expression evaluated to a value
            if subExprResult != nil {

                // RETURN - special handling if a subexpression is a return statement.
                //          in that case break out of the expression processing loop
                //          immediatelly and return to parent expression a result.
                if checkIsReturnExpr(subExprLst) {
                    return state, subExprResult, nil

                } else {
                    // REGULAR_CASE
                
                    // substitute sub-expression for its results
                    expressionLst[i] = subExprResult
                    
                    // continue looping through the expression elements,
                    // without incrementing "i". because we evaluated the expression
                    // at position "i" and substituted results of that expression at that slot.
                    continue
                }
            }

            //------------------------------------

        } else if checkIsArithmeticOp(elementStr) && i==0 {
            
            //------------------------------------
            // ARITHMETIC

            arithmeticResult, err := arithmeticEval(expressionLst, state, pExternAPI)
            if err != nil {
                return nil, nil, err
            }

            expressionResult := *arithmeticResult
            return state, expressionResult, nil

            //------------------------------------

        } else if elementIsStrBool && elementStr == "if" {

            //------------------------------------
            // CONDITIONALS

            if i != 0 {
                return nil, nil, errors.New("'if' keyword can ony be the first element in the expression")
            }

            childState, err := exprConditional(expressionLst,
                state,
                pRulesDefsMap,
                pShaderDefsMap,
                pStateFamilyStackLst,
                pExternAPI)
            if err != nil {
                return nil, nil, err
            }
            mergedState, err := stateMergeChild(state, childState)
            if err != nil {
                return nil, nil, err
            }
            state = mergedState

            break

            //------------------------------------

        } else if elementIsStrBool && (elementStr == "set" || elementStr == "push" || elementStr == "pop") {
            
            //------------------------------------
            // STATE_SETTERS - global state setters

            setterTypeStr   := elementStr
            propertyNameStr := expressionLst[1].(string)
            vals            := expressionLst[2]


            if (setterTypeStr == "push" || setterTypeStr == "pop") && propertyNameStr == "coord_origin" {
                
                // coord_origin state_setter is the only setter so far that returns
                // a new state, all other state setters dont modify the gf_lang state.
                newState, err := execStateSetterExpr(setterTypeStr,
                    propertyNameStr,
                    vals,
                    state,
                    pStateFamilyStackLst,
                    pExternAPI)
                if err != nil {
                    return nil, nil, err
                }

                state = newState

            } else {
                _, err := execStateSetterExpr(setterTypeStr,
                    propertyNameStr,
                    vals,
                    state,
                    pStateFamilyStackLst,
                    pExternAPI)
                if err != nil {
                    return nil, nil, err
                }
            }
            break

            //------------------------------------

        } else if elementIsStrBool && elementStr == "print" {

            //------------------------------------
            // PRINT

            err := exprPrint(expressionLst, state)
            if err != nil {
                return nil, nil, err
            }
            break

            //------------------------------------

        } else if elementIsStrBool && elementStr == "animate" {

            //------------------------------------
            // ANIMATION

            if i != 0 {
                return nil, nil, errors.New("'animate' keyword can ony be the first element in the expression")
            }

            exprAnimation(expressionLst, state, pExternAPI)
            break

            //------------------------------------
        
        } else if elementIsStrBool && isVar(elementStr) && i==0 && len(expressionLst) == 2 {

            //------------------------------------
            // VARIABLE_ASSIGNMENT
            // as the first element is the expression. 
            // when assigning to a variable (["$some", 10])
            // the name of the var is always expected to be the first.

            err := execVarAssignExpr(expressionLst, state, pExternAPI)
            if err != nil {
                return nil, nil, err
            }
            
            // spew.Dump(state.VarsMap)

            return state, nil, nil
            
            //------------------------------------
        
        } else if elementIsStrBool && elementStr == "return" && i==0 {

            //------------------------------------
            // RETURN
            valUnevaluated := expressionLst[1]

            // EVALUATE
            returnVal, complexSubExprBool, err := exprEvalSimple(valUnevaluated, state, pExternAPI)
            if err != nil {
                return nil, nil, err
            }

            var expressionResult interface{}

            // return statement contains a complex sub-expression, which cant be handled by exprEvalSimple(),
            // and instead it has to be handled by a full executeTree() run.
            if complexSubExprBool {
                
                subExprLst := valUnevaluated.(GFexpr)
                subExprResult, err := handleSubExpressionFun(subExprLst)
                if err != nil {
                    return nil, nil, err
                }

                expressionResult = subExprResult

            } else {
                expressionResult = returnVal
            }

            return state, expressionResult, nil

            //------------------------------------

        } else {

            //------------------------------------
            // RULE_CALL

            ruleNameStr := elementStr

            if i != len(expressionLst)-1 {
                return nil, nil, errors.New(fmt.Sprintf("rule call can only be the last element in expression; name %s, expression %s",
                    ruleNameStr,
                    expressionLst))
            }

            newState, err := exprRuleCall(ruleNameStr,
                expressionLst,
                state,
                pRulesDefsMap,
                pShaderDefsMap,
                pStateFamilyStackLst,
                pExternAPI)
            if err != nil {
                return nil, nil, err
            }

            state = newState
            break

            //------------------------------------
        }

        //------------------------------------
    }

    return state, nil, nil
}

//-------------------------------------------------

func exprRuleCall(pCalledRuleNameStr string,
    pExpressionLst       []interface{},
    pStateParent         *GFstate,
    pRulesDefsMap        GFruleDefs,
    pShaderDefsMap       map[string]interface{},
    pStateFamilyStackLst []*GFstate,
    pExternAPI           GFexternAPI) (*GFstate, error) {

    symbols := getSymbolsAndConstants()

    //------------------------------------
    // SYSTEM_RULE
    // rules predefined in the system

    if gf_core.ListContainsStr(pCalledRuleNameStr, symbols.SystemRulesLst) {

        newState := exprRuleSysCall(pCalledRuleNameStr,
            pStateParent,
            pExternAPI)
        return newState, nil

    } else if gf_core.MapHasKey(pRulesDefsMap, pCalledRuleNameStr) {

        //------------------------------------
        // USER_RULE
        // rules defined by the user in their program

        //--------------------
        // STATE_NEW
        // for each rule invocation a new state object is created, that inherits 
        // the values of its parent state, within the same state family.

        newState := stateCreateNew(pStateParent)

        //--------------------
        // get the name of the rule that is making this call to another rule
        currentRuleNameStr, err := ruleGetName(pStateParent)
        if err != nil {
            return nil, err
        }

        // all user calls (even recursive) are stored in the call stack
        newState.RulesNamesStackLst = append(newState.RulesNamesStackLst, pCalledRuleNameStr)

        // console.log(`calling rule ${currentRuleNameStr}->${pCalledRuleNameStr}`, pStateParent *GFstate["rules_names_stack_lst"]);


        // new rule getting executed
        if currentRuleNameStr != pCalledRuleNameStr {

            // start a new iterations counter since we're entering a new rule
            // (not recursively iterating within the same rule)
            addNewItersNumState(newState)
        }

        currentRuleItersNumInt := ruleGetItersNum(newState)



        // pick a random definition for a rule, which can have many definitions.
        ruleDef, ruleExpressionsLst := pickRuleRandomDef(pCalledRuleNameStr,
            pRulesDefsMap)

        //------------------------------------
        // RECURSION_STOP
        //  - prevent infinite rules execution
        //  - if global iters_max limit is reached (global for all rules)
        //  - if local rule-specific (rule modifier) iters_max limit is reached (for current rule only)

        // GLOBAL_LIMIT
        if newState.ItersNumGlobalInt > newState.ItersMaxInt-1 {

            fmt.Println("global iter limit reached")
            return newState, nil

        } else if gf_core.MapHasKey(ruleDef.ModifiersMap, "iters_max") {

            // RULE_LIMIT
            // check if rule has a iters_max rule modifier specified

            ruleItersLimitInt := int(ruleDef.ModifiersMap["iters_max"].(float64))

            if currentRuleItersNumInt > ruleItersLimitInt-1 {

                // console.log(`local iter limit ${currentRuleItersNumInt} for rule ${pCalledRuleNameStr} reached`,
                //    pStateParent *GFstate["rules_names_stack_lst"]);

                //-----------------
                // RULE_EXIT
                // rule naturally ended with iter_num limit, without entering
                // a different rule. the state has to be reset to the callers
                // state.

                // pop
                _, newRulesItersNumStackLst := gf_core.ListPop[int](newState.RulesItersNumStackLst)
                newState.RulesItersNumStackLst = newRulesItersNumStackLst

                // pop
                _, newRulesNamesStackLst := gf_core.ListPop[string](newState.RulesNamesStackLst)
                newState.RulesNamesStackLst = newRulesNamesStackLst

                oldRuleItersNumInt := ruleGetItersNum(newState)
                newState.VarsMap["$i"].Val = oldRuleItersNumInt
                
                //-----------------

                return newState, nil
            }

            //------------------------------------
        }

        // RECURSION
        // IMPORTANT!! - rules are not yet treated as expressions, and cant return
        //               results of evaluating its expressions. so ignoring execution results.
        childState, _, err := executeTree(ruleExpressionsLst,
            newState,
            pRulesDefsMap,
            pShaderDefsMap,
            pStateFamilyStackLst,
            pExternAPI)
        if err != nil {
            return nil, err
        }
        
        // remove rule_name from the stack of rules that were executed
        _, newRulesNamesStackLst := gf_core.ListPop[string](newState.RulesNamesStackLst)
        newState.RulesNamesStackLst = newRulesNamesStackLst

        if currentRuleNameStr != pCalledRuleNameStr {

            // RULE_EXIT
            // we returned from a new rule into the old rule context,
            // so the iterations count for that new rule is no longer needed (and removed from stack).
            restorePreviousRulesItersNum(newState)

            return newState, nil
        } else {

            // we're still within the same rule, in one of its iterations, so just merge state

            mergedState, err := stateMergeChild(newState, childState)
            if err != nil {
                return nil, err
            }
            return mergedState, nil
        }

        //------------------------------------

    } else {
        return nil, errors.New(fmt.Sprintf("rule call referencing an unexisting rule - %s", pCalledRuleNameStr))
    }

    //------------------------------------
}

//-------------------------------------------------

func exprRuleSysCall(pRuleNameStr string,
    pState     *GFstate,
    pExternAPI GFexternAPI) *GFstate {

    //----------------------
    pState.ItersNumGlobalInt += 1

    // IMPORTANT!! - rule iterations are counted only for actual rule evaluations.
    //               important not to count "set" statements, expression tree
    //               descending, property modifiers execution, etc.
    incrementItersNum(pState)

    //----------------------

    if pRuleNameStr == "cube" {

        // CUBE
        x := pState.Xf
        y := pState.Yf
        z := pState.Zf
        rx := pState.RotationXf
        ry := pState.RotationYf
        rz := pState.RotationZf
        sx := pState.ScaleXf
        sy := pState.ScaleYf
        sz := pState.ScaleZf
        cr := pState.ColorRedF
        cg := pState.ColorGreenF
        cb := pState.ColorBlueF

        pExternAPI.CreateCubeFun(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb)

    } else if pRuleNameStr == "sphere" {

        // SPHERE
        x := pState.Xf
        y := pState.Yf
        z := pState.Zf
        rx := pState.RotationXf
        ry := pState.RotationYf
        rz := pState.RotationZf
        sx := pState.ScaleXf
        sy := pState.ScaleYf
        sz := pState.ScaleZf
        cr := pState.ColorRedF
        cg := pState.ColorGreenF
        cb := pState.ColorBlueF

        pExternAPI.CreateSphereFun(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb)

    } else if pRuleNameStr == "line" {

        // LINE
        x := pState.Xf
        y := pState.Yf
        z := pState.Zf
        rx := pState.RotationXf
        ry := pState.RotationYf
        rz := pState.RotationZf
        sx := pState.ScaleXf
        sy := pState.ScaleYf
        sz := pState.ScaleZf
        cr := pState.ColorRedF
        cg := pState.ColorGreenF
        cb := pState.ColorBlueF

        pExternAPI.CreateLineFun(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb)
    }

    return pState
}

//-------------------------------------------------

func exprAnimation(pExpressionLst []interface{},
    pState     *GFstate,
    pExternAPI GFexternAPI) error {

    symbols := getSymbolsAndConstants()

    if len(pExpressionLst) != 3 && len(pExpressionLst) != 4 {
        return errors.New("animation expression can only have 3|4 elements")
    }

        
    var propsLst     []interface{}
    var durationSecF float64
    var repeatBool   bool

    if len(pExpressionLst) == 3 {

        propsLst     = pExpressionLst[1].([]interface{})
        durationSecF = pExpressionLst[2].(float64)

    } else if len(pExpressionLst) == 4 {

        propsLst     = pExpressionLst[1].([]interface{})
        durationSecF = pExpressionLst[2].(float64)
        repeatStr   := pExpressionLst[3].(string)

        if repeatStr == "repeat" {
            repeatBool = true

        } else {
            return errors.New("animation can only be enabled with the 'repeat' keyword")
        }
    }
    
    propsToAnimateLst := []map[string]interface{}{}

    for _, prop := range propsLst {

        propLst := prop.([]interface{})
        propNameStr  := propLst[0].(string)
        changeDeltaF := propLst[1].(float64)

        if !gf_core.ListContainsStr(propNameStr, symbols.PredefinedPropertiesLst) {
            return errors.New(fmt.Sprintf("cant animate property that is not predefined - %s", propNameStr))
        }
        
        startValF := statePropGet(pState, propNameStr)
        endValF   := startValF + changeDeltaF
        
        propsToAnimateLst = append(propsToAnimateLst, map[string]interface{}{
            "name_str":    propNameStr,
            "start_val_f": startValF,
            "end_val_f":   endValF,
        })
    }

    pExternAPI.AnimateFun(propsToAnimateLst, durationSecF, repeatBool)

    return nil
}

//-------------------------------------------------
// EXPRESSION__CONDITIONAL

func exprConditional(pExpressionLst []interface{},
    pState               *GFstate,
    pRulesDefsMap        GFruleDefs,
    pShaderDefsMap       map[string]interface{},
    pStateFamilyStackLst []*GFstate,
    pExternAPI           GFexternAPI) (*GFstate, error) {

    // [, conditionLst, subExpressionsLst] = pExpressionLst;
    conditionLst      := pExpressionLst[1].([]interface{})
    subExpressionsLst := pExpressionLst[2].([]interface{})

    if len(conditionLst) > 3 {
        return nil, errors.New("'if' condition can only have 3 elements [logic_op, operand1, operand2]")
    }

    //-------------------------------------------------
    evaluateLogicExprFun := func(pLogicExprLst []interface{}) (bool, error)  {
        
        symbols := getSymbolsAndConstants()

        logicOpStr := conditionLst[0].(string)
        operand1 := conditionLst[1]
        operand2 := conditionLst[2]

        if !gf_core.MapHasKey(symbols.LogicOperatorsMap, logicOpStr) {
            return false, errors.New(fmt.Sprintf("specified logic operator %s is not valid", logicOpStr))
        }

        //-------------------------------------------------
        op1val, _, err := exprEvalSimple(operand1, pState, pExternAPI)
        if err != nil {
            return false, err
        }
        op2val, _, err := exprEvalSimple(operand2, pState, pExternAPI)
        if err != nil {
            return false, err
        }

        
        if symbols.LogicOperatorsMap[logicOpStr](castToFloat(op1val), castToFloat(op2val)) {
            return true, nil
        } else {
            return false, nil
        }

        return false, nil
    }
    
    //-------------------------------------------------

    logicResultBool, err := evaluateLogicExprFun(conditionLst)
    if err != nil {
        return nil, err
    }

    // if condition evaluates to true, execute subexpressions
    if logicResultBool {

        // recursion
        childState, _, err := executeTree(subExpressionsLst,
            pState,
            pRulesDefsMap,
            pShaderDefsMap,
            pStateFamilyStackLst,
            pExternAPI)
        if err != nil {
            return nil, err
        }

        mergedState, err := stateMergeChild(pState, childState)
        if err != nil {
            return nil, err
        }

        return mergedState, nil
    } else {
        return pState, nil // else returned state unchanged
    }

    return nil, nil
}

//-------------------------------------------------
// EXPRESSION__PRINT

func exprPrint(pExpressionLst []interface{},
    pState *GFstate) error {

    valsLst := pExpressionLst[1].([]interface{})

    valsStr := ""
    for _, val := range valsLst {

        valIsStringBool, valStr := gf_core.CastToStr(val)

        if valIsStringBool {
            if strings.HasPrefix(valStr, "$") {

                varRefStr := valStr

                varVal, err := varEval(varRefStr, pState)
                if err != nil {
                    return err
                }

                valFmtStr := fmt.Sprintf("%s=%s ", varRefStr, varVal.Val)
                valsStr += valFmtStr
            
            } else {
                valsStr += valStr+" "
            }
        }
    }

    fmt.Printf(`gf %s\n`, valsStr)
    return nil
}