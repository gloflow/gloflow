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
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func main() {

	fmt.Println("GF_LANG >>")

	externAPI := gf_lang.GFexternAPI{
		SetStateFun: func(pNewStateMap map[string]interface{}) []interface{} {
			fmt.Println("set state")
			return nil
		},
		CreateCubeFun: func(pXf float64, pYf float64, pZf float64,
			pRotationXf float64, pRotationYf  float64, pRotationZf float64,
			pScaleXf    float64, ScaleYf      float64, ScaleZf     float64,
			pColorRedF  float64, pColorGreenF float64, pColorBlueF float64) {

			fmt.Println("create cube")
		},
		CreateSphereFun: func(pXf float64, pYf float64, pZf float64,
			pRotationXf float64, pRotationYf  float64, pRotationZf float64,
			pScaleXf    float64, ScaleYf      float64, ScaleZf     float64,
			pColorRedF  float64, pColorGreenF float64, pColorBlueF float64) {

			fmt.Println("create sphere")
		},
		CreateLineFun: func(pXf float64, pYf float64, pZf float64,
			pRotationXf float64, pRotationYf  float64, pRotationZf float64,
			pScaleXf    float64, ScaleYf      float64, ScaleZf     float64,
			pColorRedF  float64, pColorGreenF float64, pColorBlueF float64) {
			
			fmt.Println("create line")
		},
		AnimateFun: func(pPropsToAnimateLst []map[string]interface{},
			pDurationSecF float64,
			pRepeatBool   bool) {

			fmt.Println("animate")
		},
	}
	
	//----------------------
	// JS_API

	jsLangRunFun := func(pThis js.Value, p []js.Value) interface{} {


		fmt.Println("LANG RUN")

		argLst := js.ValueOf(p[0])
		
		
		
		loadedProgramASTlst := loadProgramAST(argLst)
		spew.Dump(loadedProgramASTlst)
		

		err := gf_lang.Run(loadedProgramASTlst,
			externAPI)
		
		if err != nil {
			panic(err)
		}

		outputStr := "gf_output" // p[0].Int() + p[1].Int()
		return js.ValueOf(outputStr)
	}

	js.Global().Set("gf_lang_run", js.FuncOf(jsLangRunFun))

	//----------------------

	// have to block this program to make the JS_API functions exist and callable
	c := make(chan bool)
	<-c
}

//-------------------------------------------------

func loadProgramAST(pLst js.Value) []interface{} {

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