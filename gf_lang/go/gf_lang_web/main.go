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

package main

import (
	"fmt"
	"syscall/js"
	"reflect"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func main() {

	fmt.Println("GF_LANG >>")

	
	
	//-------------------------------------------------
	// JS_API
	//-------------------------------------------------
	jsLangRunFun := func(pThis js.Value, p []js.Value) interface{} {

		fmt.Println("LANG RUN")

		// ARGS
		programASTlst := js.ValueOf(p[0])
		externAPImap  := js.ValueOf(p[1])
		
		loadedProgramASTlst := loadProgramAST(programASTlst)
		loadedExternAPI     := loadExternAPIfuns(externAPImap)

		// fmt.Println("LOADED PROGRAM AST >>>")
		// spew.Dump(loadedProgramASTlst)
		

		_, programsDebugLst, err := gf_lang.Run(loadedProgramASTlst,
			loadedExternAPI)
		
		if err != nil {
			panic(err)
		}


		debugView(programsDebugLst)

		outputStr := "gf_output"
		return js.ValueOf(outputStr)
	}

	//-------------------------------------------------

	//----------------------
	// register functions with the global browser object
	js.Global().Set("gf_lang_run", js.FuncOf(jsLangRunFun))

	//----------------------

	// have to block this program to make the JS_API functions exist and callable
	c := make(chan bool)
	<-c
}

//-------------------------------------------------

func loadExternAPIfuns(pExternAPImap js.Value) gf_lang.GFexternAPI {

	if pExternAPImap.Type() != js.TypeObject {
		panic("supplied extern_api is not a map")
	}

	initEngineFun := pExternAPImap.Get("init_engine_fun")
	if initEngineFun.Type() != js.TypeFunction {
		panic("supplied extern_api init_engine_fun is not a function")
	}

	setStateFun := pExternAPImap.Get("set_state_fun")
	if setStateFun.Type() != js.TypeFunction {
		panic("supplied extern_api set_state_fun is not a function")
	}

	createCubeFun := pExternAPImap.Get("create_cube_fun")
	if createCubeFun.Type() != js.TypeFunction {
		panic("supplied extern_api create_cube_fun is not a function")
	}

	createSphereFun := pExternAPImap.Get("create_sphere_fun")
	if createSphereFun.Type() != js.TypeFunction {
		panic("supplied extern_api create_sphere_fun is not a function")
	}

	createLineFun := pExternAPImap.Get("create_line_fun")
	if createLineFun.Type() != js.TypeFunction {
		panic("supplied extern_api create_line_fun is not a function")
	}

	animateFun := pExternAPImap.Get("animate_fun")
	if animateFun.Type() != js.TypeFunction {
		panic("supplied extern_api animate_fun is not a function")
	}
	
	externAPI := gf_lang.GFexternAPI{

		//-------------------------------------------------
		// INIT_ENGINE
		InitEngineFun: func(pShaderDefsMap map[string]*gf_lang.GFshaderDef) {

			shaderDefsForJSmap := transformShaderDefsForJS(pShaderDefsMap)

			initEngineFun.Invoke(shaderDefsForJSmap)
		},

		//-------------------------------------------------
		// SET_STATE
		SetStateFun: func(pStateChange gf_lang.GFstateChange) []interface{} {

			stateChangeMap := transformStateChangeForJS(pStateChange)

			// fmt.Println("JS STATE CHANGE >>>>>>>>>>>>>>>>")
			// spew.Dump(pStateChange)
			// spew.Dump(stateChangeMap)

			r := setStateFun.Invoke(stateChangeMap)

			//----------------
			// ADD!! - this only assumes that set_state_fun always returns just
			//         lists of floats.
			//         Add support for other datatypes or structs
			if !r.IsUndefined() {
				var resultsLst []interface{}
				for i:=0; i< r.Length(); i++{
					resultsLst = append(resultsLst, interface{}(r.Index(i).Float()))
				}
				return resultsLst
			}

			//----------------

			return nil
		},

		//-------------------------------------------------
		CreateCubeFun: func(pXf float64, pYf float64, pZf float64,
			pRotationXf float64, pRotationYf  float64, pRotationZf float64,
			pScaleXf    float64, ScaleYf      float64, ScaleZf     float64,
			pColorRedF  float64, pColorGreenF float64, pColorBlueF float64) {

			createCubeFun.Invoke(pXf, pYf, pZf,
				pRotationXf, pRotationYf, pRotationZf,
				pScaleXf, ScaleYf, ScaleZf,
				pColorRedF, pColorGreenF, pColorBlueF)
		},
		CreateSphereFun: func(pXf float64, pYf float64, pZf float64,
			pRotationXf float64, pRotationYf  float64, pRotationZf float64,
			pScaleXf    float64, ScaleYf      float64, ScaleZf     float64,
			pColorRedF  float64, pColorGreenF float64, pColorBlueF float64) {
				
			createSphereFun.Invoke(pXf, pYf, pZf,
				pRotationXf, pRotationYf, pRotationZf,
				pScaleXf, ScaleYf, ScaleZf,
				pColorRedF, pColorGreenF, pColorBlueF)
		},
		CreateLineFun: func(pXf float64, pYf float64, pZf float64,
			pRotationXf float64, pRotationYf  float64, pRotationZf float64,
			pScaleXf    float64, ScaleYf      float64, ScaleZf     float64,
			pColorRedF  float64, pColorGreenF float64, pColorBlueF float64) {
			
			createLineFun.Invoke(pXf, pYf, pZf,
				pRotationXf, pRotationYf, pRotationZf,
				pScaleXf, ScaleYf, ScaleZf,
				pColorRedF, pColorGreenF, pColorBlueF)
		},
		AnimateFun: func(pPropsToAnimateLst []map[string]interface{},
			pDurationSecF float64,
			pRepeatBool   bool) {

			animateFun.Invoke()
		},
	}
	return externAPI
}

//-------------------------------------------------

func loadProgramAST(pLst js.Value) gf_lang.GFexpr { // []interface{} {

	exprLst := []interface{}{}
	for i:=0; i < pLst.Length() ;i++ {

		e := pLst.Index(i)

		switch e.Type() {
		case js.TypeObject:
			subExprLst := loadProgramAST(e)
			exprLst = append(exprLst, subExprLst)

		case js.TypeString:
			eStr := e.String()
			exprLst = append(exprLst, interface{}(eStr))

		case js.TypeNumber:
			eF := e.Float()
			exprLst = append(exprLst, interface{}(eF))
		
		case js.TypeBoolean:
			eBool := e.Bool()
			exprLst = append(exprLst, interface{}(eBool))
		}
	}
	return exprLst
}


//-------------------------------------------------

func transformShaderDefsForJS(pShaderDefsMap map[string]*gf_lang.GFshaderDef) map[string]interface{} {

	shaderDefsMap := map[string]interface{}{}

	for shaderNameStr, shaderDef := range pShaderDefsMap {
		
		uniformDefsLst := []interface{}{}
		for _, u := range shaderDef.UniformsDefsLst {
			uniformDefsLst = append(uniformDefsLst, []interface{}{u.NameStr, u.TypeStr, u.DefaultVal})
		}

		shaderDefMap := map[string]interface{}{
			"name_str":          shaderDef.NameStr,
			"uniforms_defs_lst": uniformDefsLst,
			"vertex_code_str":   shaderDef.VertexCodeStr,
			"fragment_code_str": shaderDef.FragmentCodeStr,
		}

		shaderDefsMap[shaderNameStr] = shaderDefMap
	}
	return shaderDefsMap
}

//-------------------------------------------------

func transformStateChangeForJS(pStateChange gf_lang.GFstateChange) map[string]interface{} {
	
	stateChangeMap := map[string]interface{}{}

	structVal  := reflect.ValueOf(pStateChange)
	structType := reflect.TypeOf(pStateChange)

	for i := 0; i < structType.NumField(); i++ {

		fieldType  := structType.Field(i)
		fieldValue := structVal.Field(i)
		// fieldNameStr := fieldType.Name

		// get the 'json' tag from the struct member
		fieldJSONnameStr := reflect.StructTag(fieldType.Tag).Get("json")

		switch fieldType.Type.Kind() {
		case reflect.Slice:
			// golang/js bridge expects all arrays to be typed as []interface{}
			fieldLst := []interface{}{}

			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j).Interface()
				fieldLst = append(fieldLst, elem)
			}

			stateChangeMap[fieldJSONnameStr] = fieldLst
			
		case reflect.String:
			stateChangeMap[fieldJSONnameStr] = fieldValue.String()
			
		case reflect.Float64:
			stateChangeMap[fieldJSONnameStr] = fieldValue.Float()
			
		case reflect.Int:
			stateChangeMap[fieldJSONnameStr] = fieldValue.Int()
			
		case reflect.Int64:
			stateChangeMap[fieldJSONnameStr] = fieldValue.Int()
			
		case reflect.Bool:
			stateChangeMap[fieldJSONnameStr] = fieldValue.Bool()

		case reflect.Interface:
			stateChangeMap[fieldJSONnameStr] = fieldValue.Interface()
		}
	}
	return stateChangeMap	
}