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
	"os"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func TestMain(m *testing.M) {

	os.Setenv("GF_LOG_LEVEL", "debug")
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------

func TestVariables(pTest *testing.T) {
	
	_, logNewFun := gf_core.LogsInit()
	externAPI := getTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},

		// testing setting of a user variable, and referencing
		// and returning it.
		GFexpr{
			"return", GFexpr{
					GFexpr{"$test_var", 10},
					GFexpr{"return", "$test_var"},
				},
		},
	}

	// run program
	resultsLst, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program in test_basic", map[string]interface{}{"err": err,})
		pTest.Fail()
	}

	spew.Dump(resultsLst)

	if len(resultsLst) != 1 {
		logNewFun("ERROR", "results list should be 1 elements", nil)
		pTest.Fail()
	}
	if resultsLst[0] != 10 {
		logNewFun("ERROR", "first results should be equal to 10", map[string]interface{}{"result": resultsLst[0],})
		pTest.Fail()
	}
}

//---------------------------------------------------

func TestReturnExpression(pTest *testing.T) {
	
	_, logNewFun := gf_core.LogsInit()
	externAPI := getTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},
		GFexpr{"return", 10},
		GFexpr{"return", GFexpr{"*", 3, 12}},
	}

	// run program
	resultsLst, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program in test_basic", map[string]interface{}{"err": err,})
		pTest.Fail()
	}

	spew.Dump(resultsLst)

	if len(resultsLst) != 2 {
		logNewFun("ERROR", "results list should be 2 elements", nil)
		pTest.Fail()
	}

	if resultsLst[0] != 10 {
		logNewFun("ERROR", "first results should be equal to 10", map[string]interface{}{"result": resultsLst[0],})
		pTest.Fail()
	}

	if int(resultsLst[1].(float64)) != 36 {
		logNewFun("ERROR", "first results should be equal to 36", map[string]interface{}{"result": resultsLst[1],})
		pTest.Fail()
	}
}

//---------------------------------------------------

func getTestExternAPI() GFexternAPI {
	externAPI := GFexternAPI{

		InitEngineFun: func(pShaderDefsMap map[string]interface{}) {
			fmt.Println("init_engine")
		},
		SetStateFun: func(pStateChange GFstateChange) []interface{} {
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

		//---------------------------------------------
		// RPC_CALL
		RPCcall: func(pNodeStr string, // node
			pModuleStr   string,       // module
			pFunctionStr string,       // function
			pArgsLst     []interface{}) map[string]interface{} { // args list
			

			return nil


		},

		//---------------------------------------------
		// RPC_SERVE
		RPCserve: func(pNodeNameStr string,
			pHandlersLst []*GFrpcServerHandler,
			pExternAPI   GFexternAPI) {

			
		},

		//---------------------------------------------
	}
	return externAPI
}