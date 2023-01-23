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
    "encoding/json"
    "github.com/gloflow/gloflow/go/gf_core"
    // "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GFexpr     []interface{}
type GFruleDefs map[string][]*GFruleDef

type GFruleDef struct {
    NameStr        string
    ModifiersMap   map[string]interface{}
    ExpressionsLst GFexpr
}

type GFsymbols struct {
    RuleLevelMaxInt         int
    SystemRulesLst          []string
    PredefinedPropertiesLst []string
    LogicOperatorsMap       map[string]func(float64, float64) bool
    ArithmeticOpsMap        map[string]func(float64, float64) float64
    SystemVarsLst           []string
    SystemFunctionsLst      []string
}

type GFshaderDef struct {
    NameStr         string               
    UniformsDefsLst []GFshaderUniformDef
    VertexCodeStr   string
    FragmentCodeStr string
}

type GFshaderUniformDef struct {
    NameStr    string
    TypeStr    string
    DefaultVal interface{}
}

//-------------------------------------------------

func Run(pProgramASTlst GFexpr,
    pExternAPI GFexternAPI) ([]interface{}, []*GFprogramDebug, error) {

    //------------------------------------
    // AST_EXPANSION
    expandedProgramASTlst := GFexpr{}
    for _, rootExpression := range pProgramASTlst {

        rootExpressionLst := rootExpression.(GFexpr)
        expandedRootExpressionLst, err := expandTree(rootExpressionLst, 0)
        if err != nil {
            return nil, nil, err
        }

        // only include expressions which are not "expanded" to expression of 0 length
        // (expressions which are not marked for deletion)
        if len(expandedRootExpressionLst) > 0 {

            // fmt.Println("expanded expr", expandedRootExpressionLst)
            expandedProgramASTlst = append(expandedProgramASTlst, interface{}(expandedRootExpressionLst))
        }
    }

    // DEBUG - view new program AST
    // viewExpandedProgram(expandedProgramASTlst)    

    //------------------------------------
    // only load rules after AST tree expansion is complete,
    // and the rule is ready for execution.
    ruleDefsMap, programASTlnoRuleDefsLst := loadRuleDefs(expandedProgramASTlst)
    
    shaderDefsMap, programASTnoShaderDefsLst, err := loadShaderDefs(programASTlnoRuleDefsLst)
    if err != nil {
        return nil, nil, err
    }

    // INIT_ENGINE
    pExternAPI.InitEngineFun(shaderDefsMap)
    
    //------------------------------------
    // AST_EXECUTION

    fmt.Println("executing AST...")
    // spew.Dump(programASTnoShaderDefsLst)

    i:=0
    exprResultsLst := []interface{}{}

    // IMPORTANT!! - all states are accumulated here when running each root expression,
    //               and returned to the runner of the program, for debugging purposes, so that
    //               the runner can inspect the produced state and make checks.
    programsDebugLst := []*GFprogramDebug{}
    
    for _, rootExpression := range programASTnoShaderDefsLst {

        //------------------------------------
        // DEBUG
        programDebug := debugInit()
        programsDebugLst = append(programsDebugLst, programDebug)

        //------------------------------------
        rootState := stateCreateNew(nil, programDebug)
        
        //------------------------------------
        // STATE_FAMILY_STACK
        // STATE_FAMILY - is a group of stacks that are related and are treated independently,
        //                without copying/merging back into the states from which it came from.
        //
        // allowing for multiple state objects to exist durring expression/program execution.
        // initially there was only one blank state that the program started with, and all future
        // operations worked with that one state.
        // going forward there are state_setters that allow for pushing new states onto that stack,
        // and popping them.
        
        stateFamilyStackLst := []*GFstate{}
        
        //------------------------------------
        rootExpressionLst := rootExpression.(GFexpr)

        _, exprResult, err := executeTree(rootExpressionLst,
            rootState,
            ruleDefsMap,
            shaderDefsMap,
            stateFamilyStackLst,
            pExternAPI,
            programDebug)
        if err != nil {
            return nil, programsDebugLst, err
        }

        // accumulate as many results as there are top-level expressions
        exprResultsLst = append(exprResultsLst, exprResult)
        
        i+=1
    }
    
    //------------------------------------

    return exprResultsLst, programsDebugLst, nil
}

//-------------------------------------------------

func expandTree(pExpressionASTlst GFexpr,
    pTreeLevelInt int) (GFexpr, error) {

    expressionLst := cloneExpr(pExpressionASTlst) // clone in case of mutations of expression

    // fmt.Println("expand tree", expressionLst)
    
    //-------------------------------------------------
    // SUB_EXPRESSION
    // [...]

    handleSubExpressionFun := func(pSubExpressionLst GFexpr, pIndexInt int) error {

        // RECURSION
        expandedSubExpressionLst, err := expandTree(pSubExpressionLst, pTreeLevelInt+1)
        if err != nil {
            return err
        }

        // IMPORTANT!! - splice the expanded sub-expression in the place of the old unexpanded element
        expressionLst[pIndexInt] = expandedSubExpressionLst

        return nil
    }

    //-------------------------------------------------

    for i:=0; i < len(expressionLst); {

        element := expressionLst[i]
        elementIsStrBool, elementStr := gf_core.CastToStr(element)

        if elementIsStrBool && elementStr == "rule" {

            //------------------------------------
            // RULE_DEFINITION
            // [rule rule_name, [...]]
            
            //-----------
            // VERIFICATION
            if pTreeLevelInt != 0 {
                return nil, errors.New("rule definitions can only exist at the top level")
            }
            
            if i != 0 {
                return nil, errors.New("rule definition has to be at the start of the expression")
            }

            // rule definition can be of form:
            // ["rule", name_str, expressions_lst]
            // ["rule", name_str, ruleModifiersLst, expressions_lst]
            if len(expressionLst) != 3 && len(expressionLst) != 4 {
                return nil, errors.New("rule definition expression can only have 3|4 elements")
            }

            if len(expressionLst) == 3 {
                if _, ok := expressionLst[2].(GFexpr); !ok {
                    return nil, errors.New("rule definitions 3rd element has to be a list of its expressions")
                }
            }

            if len(expressionLst) == 4 { 
                if _, ok := expressionLst[3].(GFexpr); !ok {
                    return nil, errors.New("rule definitions 4rd element has to be a list of its expressions")
                }
            }

            // fast-forward to the 3rd|4th element of the expression, which represents rules expressions
            // so that tree expansion can be run on that rules expressions.
            if len(expressionLst) == 3 {
                i+=2
            }
            if len(expressionLst) == 4 {
                i+=3
            }

            //-----------

            //------------------------------------
            
            continue

        } else if elementIsStrBool && elementStr == "set" {

            //------------------------------------
            // SET
            if i != 0 {
                return nil, errors.New("'set' declaration has to be at the start of the expression");
            }

            //------------------------------------

            // ADD!! - localize "i" incrementing to each sub-block, instead at the end of 'for' loop.
            //         for better viewing/debug.
            // i+=1

        } else if elementIsStrBool && elementStr == "*" && i==0 {

            //------------------------------------
            // MULTIPLICATION__TOP_LEVEL
            // [* op1 op2]

            operand1indexInt := i+1
            operand2indexInt := i+2

            operand1 := expressionLst[operand1indexInt]
            operand2 := expressionLst[operand2indexInt]

            //-----------
            // OPERAND_1
            _, operand1isFloatBool  := operand1.(float64)
            _, operand1isIntBool    := operand1.(int)
            _, operand1isExprBool   := operand1.(GFexpr)
            _, operand1isStringBool := operand1.(string)

            var operand1isNumBool bool
            if operand1isFloatBool || operand1isIntBool {
                operand1isNumBool = true
            }

            fmt.Println(expressionLst)
            fmt.Println("* operand 1 is", operand1isNumBool, operand1isExprBool, operand1isStringBool)

            // IMPORTANT!! - if operand1 is not a static number no compile-time expansion is done.
            //
            // operand1 is not a static value so cant be expanded at compile time.
            // instead skip, and it will be handled dynamically at interpretation time.
            // do some validation of it if its a variable reference, and if subexpression then
            // handle it as such.
            //
            // first operand not being a float64 use to be an error, it was meant for static/compile-time expansion;
            // but with general multiplication "*" that can happen at runtime with operands being expressions with runtime values, it no longer is.
            
            // EXIT_LOGIC
            if !operand1isNumBool {

                //-----------
                // SUB_EXPRESSION
                subExpressionLst, _ := operand1.(GFexpr)
                if operand1isExprBool {
                    err := handleSubExpressionFun(subExpressionLst, operand1indexInt)
                    if err != nil {
                        return nil, err
                    }

                    // both operand2 and operand1 have been expanded at this point,
                    // so exit the loop for processing this multiply subexpression.
                    break
                }
                
                //-----------
                // VARIABLE
                operand1str, _ := operand1.(string)
                if operand1isStringBool {
                    varStr := operand1str

                    // verify variable reference
                    err := varVerify(varStr)
                    if err != nil {
                        return nil, err
                    }

                    // when handling variables no expansion is done.
                    // since operand2 was already expanded, this multiplication statement 
                    break
                }
                
                //-----------

                if !(operand1isNumBool && operand1isExprBool && operand1isStringBool) {
                    // if operand1 is not a number, or list, or variable reference
                    return nil, errors.New(fmt.Sprintf("first operand of multiplication expression is not a number or a sub_expression or a string (variable reference) - %s", expressionLst))
                }
            }

            //-----------
            // MAIN_LOGIC
            // operand1 is a float64 (static/compile-time) so the expansion of the operand2 can be completed.
            if _, ok := operand2.(GFexpr); ok && operand1isNumBool {

                expressionToMultiplyLst := operand2.(GFexpr)
                factorInt               := castToInt(operand1)
                expandedExpressionsLst  := cloneExprNtimes(expressionToMultiplyLst, factorInt)

                /*
                ["*", 10, [["y", -2.0], "cube"]], // 10 * {x -2} cube
                transforms to:

                [
                    [["y", -2.0], "cube"],
                    [["y", -2.0], "cube"],
                    [["y", -2.0], "cube"],
                ]
                */
                // multiplication of sub-expression has been projected,
                // and multiplication expression itself eliminated/replaced by new cloned expressions
                expressionLst = expandedExpressionsLst

                //-----------
                // rewind to start, since the expression has a new form, and should be re-processed (where operand2 is multiplied
                // at compile-time multiple times)
                i=0
                
                // go straight to new iteration, without incrementing 'i' (keeping it at 0 instead)
                continue

                //-----------
            }

            //------------------------------------

        } else if subExpressionLst, ok := element.(GFexpr); ok {

            //------------------------------------
            // SUB_EXPRESSION
            err := handleSubExpressionFun(subExpressionLst, i)
            if err != nil {
                return nil, err
            }

            //------------------------------------

            // ADD!! - localize "i" incrementing to each sub-block, instead at the end of 'for' loop.
            //         for better viewing/debug.
            // i+=1

        } else if elementIsStrBool && elementStr == "lang_v" {

            //------------------------------------
            // LANG_VERSION

            if i != 0 {
                return nil, errors.New("lang_v expression identifier can only be the first element in the expression")
            }

            // this expression is expanded to expression of 0 length, meaning it should be removed.
            return GFexpr{}, nil

            //------------------------------------
        }

        i+=1
    }
    return expressionLst, nil
}

//-------------------------------------------------

func loadRuleDefs(pProgramASTlst GFexpr) (GFruleDefs, GFexpr) {

    fmt.Println("loading rule defs...")

    //------------------------
    // copy program_ast
    newProgramASTlst := cloneExpr(pProgramASTlst)

    //------------------------

    ruleDefsMap := map[string][]*GFruleDef{}

    for i:=0; i < len(newProgramASTlst); {

        rootExpressionLst := newProgramASTlst[i].(GFexpr)
        rootExprFirstElementIsStrBool, rootExprFirstElementStr := gf_core.CastToStr(rootExpressionLst[0])

        if rootExprFirstElementIsStrBool && rootExprFirstElementStr == "rule" {

            // rule with no modifiers
            if len(rootExpressionLst) == 3 {

                ruleNameStr        := rootExpressionLst[1].(string)
                ruleExpressionsLst := rootExpressionLst[2].(GFexpr)

                ruleDef := &GFruleDef{
                    NameStr:        ruleNameStr,
                    ExpressionsLst: ruleExpressionsLst,
                    ModifiersMap:   map[string]interface{}{},
                }


                if _, ok := ruleDefsMap[ruleNameStr]; ok {
                    ruleDefsMap[ruleNameStr] = append(ruleDefsMap[ruleNameStr], ruleDef)

                } else {
                    ruleDefsMap[ruleNameStr] = []*GFruleDef{ruleDef,}
                }
                

                // remove the rule definition element from the program_ast, 
                // as it has been expanded and loaded and ready for execution,
                // it doesnt need to be iterated over during execution.
                newProgramASTlst = gf_core.ListRemoveElementAtIndex(newProgramASTlst, i)

                // run next iteration without incrementing "i"
                continue

            } else if len(rootExpressionLst) == 4 {

                // rule with modifiers

                ruleNameStr        := rootExpressionLst[1].(string)
                ruleModifiersLst   := rootExpressionLst[2].(GFexpr)
                ruleExpressionsLst := rootExpressionLst[3].(GFexpr)

                // MODIFIERS
                ruleModifiersMap := map[string]interface{}{}
                for _, modifier := range ruleModifiersLst {

                    modifierLst     := modifier.(GFexpr)
                    modifierNameStr := modifierLst[0].(string)
                    modifierVal     := modifierLst[1]

                    ruleModifiersMap[modifierNameStr] = modifierVal
                }

                ruleDef := &GFruleDef{
                    NameStr:        ruleNameStr, 
                    ModifiersMap:   ruleModifiersMap,
                    ExpressionsLst: ruleExpressionsLst,
                }

                if _, ok := ruleDefsMap[ruleNameStr]; ok {
                    ruleDefsMap[ruleNameStr] = append(ruleDefsMap[ruleNameStr], ruleDef)
                    
                } else {
                    ruleDefsMap[ruleNameStr] = []*GFruleDef{ruleDef,}
                }

                // remove the rule definition element from the program_ast, 
                // as it has been expanded and loaded and ready for execution,
                // it doesnt need to be iterated over during execution.
                newProgramASTlst = gf_core.ListRemoveElementAtIndex(newProgramASTlst, i)

                // run next iteration without incrementing "i"
                continue
            }
        }

        i+=1
    }
    return ruleDefsMap, newProgramASTlst
}

//-------------------------------------------------

func loadShaderDefs(pProgramASTlst GFexpr) (map[string]*GFshaderDef, GFexpr, error) {

    shaderDefsMap := map[string]*GFshaderDef{}
    
    //------------------------
    // copy program_ast
    newProgramASTlst := cloneExpr(pProgramASTlst)

    //------------------------

    for i:=0; i < len(newProgramASTlst);  {

        rootExpressionLst := newProgramASTlst[i].(GFexpr)
        rootExprFirstElementIsStrBool, rootExprFirstElementStr := gf_core.CastToStr(rootExpressionLst[0])

        if rootExprFirstElementIsStrBool && rootExprFirstElementStr == "shader" {

            if len(rootExpressionLst) != 4 && len(rootExpressionLst) != 5 {
                return nil, nil, errors.New(fmt.Sprintf("shader definition expression %s can only have 4|5 elements", rootExpressionLst))
            }

            if len(rootExpressionLst) == 4 {
                shaderNameStr     := rootExpressionLst[1].(string)
                vertexShaderLst   := rootExpressionLst[2].(GFexpr)
                fragmentShaderLst := rootExpressionLst[3].(GFexpr)

                vertexCodeStr   := vertexShaderLst[1].(string)
                fragmentCodeStr := fragmentShaderLst[1].(string)

                shaderDef := &GFshaderDef{
                    NameStr:         shaderNameStr,
                    VertexCodeStr:   vertexCodeStr,
                    FragmentCodeStr: fragmentCodeStr,
                }

                shaderDefsMap[shaderNameStr] = shaderDef
            }

            if len(rootExpressionLst) == 5 {

                shaderNameStr       := rootExpressionLst[1].(string)
                uniformsDefsExprLst := rootExpressionLst[2].(GFexpr)
                vertexShaderLst     := rootExpressionLst[3].(GFexpr)
                fragmentShaderLst   := rootExpressionLst[4].(GFexpr)

                vertexCodeStr   := vertexShaderLst[1].(string)
                fragmentCodeStr := fragmentShaderLst[1].(string)
                uniformsDefsLst := uniformsDefsExprLst[1]

                uniformDefsLst := []GFshaderUniformDef{}
                for _, uniformDefinterface := range uniformsDefsLst.(GFexpr) {
                    uniformDefExpr, _ := uniformDefinterface.(GFexpr)

                    uniformNameStr    := uniformDefExpr[0].(string)
                    uniformTypeStr    := uniformDefExpr[1].(string)
                    uniformDefaultVal := uniformDefExpr[2].(interface{})

                    uniformDef := GFshaderUniformDef{
                        NameStr:    uniformNameStr,
                        TypeStr:    uniformTypeStr,
                        DefaultVal: uniformDefaultVal,
                    }

                    uniformDefsLst = append(uniformDefsLst, uniformDef)
                }

                shaderDef := &GFshaderDef{
                    NameStr:         shaderNameStr,
                    UniformsDefsLst: uniformDefsLst,
                    VertexCodeStr:   vertexCodeStr,
                    FragmentCodeStr: fragmentCodeStr,
                }
                shaderDefsMap[shaderNameStr] = shaderDef
            }

            // remove the shader definition element from the program_ast, 
            // as it has been expanded and loaded and ready for execution,
            // it doesnt need to be iterated over during execution.
            newProgramASTlst = gf_core.ListRemoveElementAtIndex(newProgramASTlst, i)
        }

        i+=1
    }
    return shaderDefsMap, newProgramASTlst, nil
}

//-------------------------------------------------

func viewExpandedProgram(pExpandedProgramASTlst GFexpr) {

    fmt.Println("=======================>>>>>>>>>>>>>>")
    fmt.Println("EXPANDED PROGRAM")
    jsonData, _ := json.MarshalIndent(pExpandedProgramASTlst, "", "    ")
    
    fmt.Println(string(jsonData))
    // spew.Dump(pExpandedProgramASTlst)

    fmt.Println("=======================>>>>>>>>>>>>>>")
}