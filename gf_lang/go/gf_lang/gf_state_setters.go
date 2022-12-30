
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
    "strings"
    "github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

func execExpr(pSetterTypeStr string,
    pPropertyNameStr     string,
    pVals                interface{},
    pState               *GFstate,
    pStateFamilyStackLst []*GFstate,
    pExternAPI           GFexternAPI) (*GFstate, error) {
    
    symbols := getSymbolsAndConstants()

    if pSetterTypeStr != "set" && pSetterTypeStr != "push" && pSetterTypeStr != "pop" {
        return nil, errors.New("state setter is not of type 'set|push|pop'")
    }

    //------------------------------------
    // SCALE
    if pPropertyNameStr == "scale" {
        
        if valF, ok := pVals.(float64); ok {

            scaleValF := valF
            pState.ScaleXf = scaleValF
            pState.ScaleYf = scaleValF
            pState.ScaleZf = scaleValF

        } else if valLst, ok := pVals.([]interface{}); ok {

            valsLst := valLst
            if len(valsLst) != 3 {
                return nil, errors.New("scale values can only be of length 3 [sx, sy, sz]")
            }

            sx := valsLst[0]
            sy := valsLst[1]
            sz := valsLst[2] 

            sxResult, err := exprEval(sx, pState)
            if err != nil {
                return nil, err
            }
            syResult, err := exprEval(sy, pState)
            if err != nil {
                return nil, err
            }
            szResult, err := exprEval(sz, pState)
            if err != nil {
                return nil, err
            }

            pState.ScaleXf = sxResult.(float64)
            pState.ScaleYf = syResult.(float64)
            pState.ScaleZf = szResult.(float64)
        }

    } else if pPropertyNameStr == "color" {
        
        //------------------------------------
        // COLOR

        if valLst, ok := pVals.([]interface{}); ok {

            valsLst := valLst
            if valsLst[0] != "rgb" {
                return nil, errors.New("only rgb type is allowed")
            }

            if len(valsLst) != 4 {
                return nil, errors.New("rgb values can only be of length 4 ['rgb', r, g, b]")
            }

            rF := valsLst[1].(float64)
            gF := valsLst[2].(float64)
            bF := valsLst[3].(float64)

            stateChangeMap := map[string]interface{}{
                "setter_type_str": pSetterTypeStr,
                "color_rgb":       []float64{rF, gF, bF},
            }
            pExternAPI.SetStateFun(stateChangeMap)

            pState.ColorRedF   = rF
            pState.ColorGreenF = gF
            pState.ColorBlueF  = bF

        } else {

            isStringBool, valsStr := gf_core.CastToStr(pVals)
            if !isStringBool {
                return nil, errors.New("setting color has to be either an array or string")
            }

            colorHexStr := valsStr

            if !strings.HasPrefix(colorHexStr, "#") {
                return nil, errors.New("setting color with a string has to be done in hex format starting with #")
            }

            stateChangeMap := map[string]interface{}{
                "setter_type_str": pSetterTypeStr,
                "color_rgb":       colorHexStr,
            }
            // set_state_fun() will return parsed color hex string, with rgb channels in 0-1 range.
            resultLst := pExternAPI.SetStateFun(stateChangeMap)

            pState.ColorRedF   = resultLst[0].(float64)
            pState.ColorGreenF = resultLst[1].(float64)
            pState.ColorBlueF  = resultLst[2].(float64)
        }

        //------------------------------------

    } else if pPropertyNameStr == "color-background" {

        //------------------------------------
        // COLOR_BACKGROUND

        valsLst := pVals.([]interface{})
        if valsLst[0].(string) != "rgb" {
            return nil, errors.New("only rgb type is allowed")
        }
        if len(valsLst) != 4 {
            return nil, errors.New("rgb values can only be of length 4 ['rgb', r, g, b]")
        }

        r := valsLst[1].(float64)
        g := valsLst[2].(float64)
        b := valsLst[3].(float64)
        
        stateChangeMap := map[string]interface{}{
            "setter_type_str":  pSetterTypeStr,
            "color_background": []float64{r, g, b,},
        }
        pExternAPI.SetStateFun(stateChangeMap)

        //------------------------------------

    } else if pPropertyNameStr == "iters_max" {

        //------------------------------------
        // ITERS_MAX

        iterationsMaxInt := pVals.(int)
        pState.ItersMaxInt = iterationsMaxInt

        //------------------------------------
        
    } else if pPropertyNameStr == "material" {

        //------------------------------------
        // MATERIAL

        valsLst         := pVals.([]interface{})
        materialTypeStr := valsLst[0].(string)

        if materialTypeStr != "wireframe" &&
            materialTypeStr != "shader" {
            return nil, errors.New("only 'wireframe|shader' material types are supported")
        }

        stateChangeMap := map[string]interface{}{
            "setter_type_str":    pSetterTypeStr,
            "material_type_str":  materialTypeStr,
        }

        if materialTypeStr == "wireframe" {

            if valBool, ok := valsLst[1].(bool); ok {
                stateChangeMap["material_value_bool"] = valBool
            } else {
                return nil, errors.New("if material_type is 'wireframe', the value has to be a bool indicating if wireframe is on/off")
            }

        } else if materialTypeStr == "shader" {
            valIsStrBool, valStr := gf_core.CastToStr(valsLst[1])
            if !valIsStrBool {
                return nil, errors.New("if material_type is 'shader', the value has to be a string representing the shader name")
            }
            stateChangeMap["material_value_str"] = valStr

        } else {
            return nil, errors.New("only 'wireframe|shader' material types are supported")
        }

        pExternAPI.SetStateFun(stateChangeMap)      
        
        //------------------------------------

    } else if pPropertyNameStr == "material_prop" {

        //------------------------------------
        // MATERIAL_PROPERTY

        valsLst := pVals.([]interface{})
        materialNameStr := valsLst[0].(string)
        materialPropStr := valsLst[1].(string)
        materialPropVal := valsLst[2]

        if materialPropStr != "shader_uniform" {
            return nil, errors.New("only 'shader_uniform' material properties are supported")
        }
    
        if materialPropStr == "shader_uniform" {

            materialPropValLst := materialPropVal.([]interface{})
            uniformNameStr := materialPropValLst[0].(string)
            uniformVal     := materialPropValLst[1]
            uniformValIsStrBool, uniformValStr := gf_core.CastToStr(uniformVal)

            var loadedVal interface{}

            // VARIABLE_REFERENCE
            if uniformValIsStrBool && strings.HasPrefix(uniformValStr, "$") {

                possiblePropNameStr := strings.Trim(uniformValStr, "$") // remove "$"

                // SYSTEM_PROPERTY - x|y|z|...|cr|cg|cb
                if gf_core.ListContainsStr(possiblePropNameStr, symbols.PredefinedPropertiesLst) {
                    propertyNameStr := possiblePropNameStr
                    loadedVal = statePropGet(pState, propertyNameStr)

                } else {
                    // USER_DEFINED_VARIABLE
                    // evalue the variable reference to get its value
                    loadedVal = pState.VarsMap[uniformValStr]
                }

            } else if uniformValLst, ok := uniformVal.([]interface{}); ok {

                // ARITHMETIC_EXPRESSION

                subExprLst := uniformValLst
                mulResult, err := arithmeticEval(subExprLst, pState)
                if err != nil {
                    return nil, err
                }

                loadedVal = mulResult

            } else {

                // NUMERIC_VALUE
                loadedVal = uniformVal
            }
            
            stateChangeMap := map[string]interface{}{
                "setter_type_str": pSetterTypeStr,
                "material_prop_map": map[string]interface{}{
                    "material_shader_name_str":         materialNameStr,
                    "material_shader_uniform_name_str": uniformNameStr,
                    "material_shader_uniform_val":      loadedVal,
                },
            }
            pExternAPI.SetStateFun(stateChangeMap)
        }

        //------------------------------------

    } else if pPropertyNameStr == "line" {

        //------------------------------------
        // LINE

        // valsLst := pVals.([]interface{})
        // cmdStr  := valsLst[0].(string)

        stateChangeMap := map[string]interface{}{
            "setter_type_str": pSetterTypeStr,
            "line_cmd_str":    "start",
        }
        pExternAPI.SetStateFun(stateChangeMap)

        //------------------------------------
        
    } else if pPropertyNameStr == "rotation_pivot" {

        //------------------------------------
        // ROTATION_PIVOT

        axisTypeStr := pVals.(string)
        if axisTypeStr == "current_pos" {

            stateChangeMap := map[string]interface{}{
                "property_name_str": "rotation_pivot",
                "setter_type_str":   pSetterTypeStr,
                "axis_type_str":     "current_pos",
                "x":  pState.Xf,
                "y":  pState.Yf,
                "z":  pState.Zf,
                "rx": pState.RotationXf,
                "ry": pState.RotationYf,
                "rz": pState.RotationZf,
            }
            pExternAPI.SetStateFun(stateChangeMap)
        }

        //------------------------------------

    } else if pPropertyNameStr == "coord_origin" {

        //------------------------------------
        // COORD_ORIGIN - setting where the origin for subsequent operation should be.
        //                it can either be the current_position or world origin.

        originTypeStr := pVals
        if originTypeStr != "current_pos" {
            return nil, errors.New("'coord_origin' setter has to have an type of 'current_pos'")
        }

        var newState *GFstate
        switch pSetterTypeStr {
            
            //------------------------------------
            case "push":

                newStateChangeMap := map[string]interface{}{
                    "property_name_str": "coord_origin",
                    "setter_type_str":   "push",
                    "origin_type_str":   originTypeStr,
        
                    "x":  pState.Xf,
                    "y":  pState.Yf,
                    "z":  pState.Zf,
                    "rx": pState.RotationXf,
                    "ry": pState.RotationYf,
                    "rz": pState.RotationZf,
                }
                pExternAPI.SetStateFun(newStateChangeMap)

                //---------------------------
                // NEW_BLANK_STATE - only other place where this is being done
                //                   is at the root of the program execution.
                newFamilyState := stateCreateNewFamily(pState)

                //---------------------------

                // this is the last state of the current family, that will need to be restored
                // to as the current state when the current family is popped
                currentFamilyLastState := pState
                statePushFamily(currentFamilyLastState, pStateFamilyStackLst)
                
                newState = newFamilyState

            //------------------------------------
            case "pop":
                
                lastFamilyState, err := statePopFamily(pStateFamilyStackLst)
                if err != nil {
                    return nil, err
                }
                
                newState = lastFamilyState

                restoreStateChangeMap := map[string]interface{}{
                    "property_name_str": "coord_origin",
                    "setter_type_str":   "pop",
                    "origin_type_str":   originTypeStr,
        
                    "x":  lastFamilyState.Xf,
                    "y":  lastFamilyState.Yf,
                    "z":  lastFamilyState.Zf,
                    "rx": lastFamilyState.RotationXf,
                    "ry": lastFamilyState.RotationYf,
                    "rz": lastFamilyState.RotationZf,
                }
                pExternAPI.SetStateFun(restoreStateChangeMap)

            //------------------------------------
        }

        //------------------------------------
        return newState, nil
    }

    return nil, nil
}