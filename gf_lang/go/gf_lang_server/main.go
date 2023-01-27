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
	"os"
	"fmt"
	"net/http"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
)

//-------------------------------------------------

type GFplugins struct {
	PyBaseDirStr string
}

//-------------------------------------------------

func main() {

	fmt.Printf("\n   GF_LANG SERVER >>\n\n")

	appServerPortInt := 5000
	
	logFun, logNewFun := gf_core.LogsInitNew(true, "DEBUG")
	log.SetOutput(os.Stdout)

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_lang_server",
		EnvStr:         "dev",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
	}



	plugins := &GFplugins{
		PyBaseDirStr: "./../../py/",
	}

	//--------------------------
	localTestProgramStr := "./tests/first_scene.gf"
	programASTlst, gfErr := ParseProgramASTfromFile(localTestProgramStr, runtimeSys)
	if gfErr != nil {
		panic(1)
	}

	//--------------------------

	externAPI := gf_lang.GFexternAPI{

		InitEngineFun: func(pShaderDefsMap map[string]*gf_lang.GFshaderDef) {
			fmt.Println("init_engine")
		},
		SetStateFun: func(pStateChange gf_lang.GFstateChange) []interface{} {
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
			pHandlersLst []*gf_lang.GFrpcServerHandler,
			pExternAPI   gf_lang.GFexternAPI) {

			// HTTP_MUX
			HTTPmux := http.NewServeMux()

			for _, h := range pHandlersLst {

				//-------------------------------------------------
				// handler_fun
				handlerFun := func(pCtx context.Context,
					pResp http.ResponseWriter,
					pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

					
					handlerProgramASTlst := h.CodeASTlst

					//---------------------
					// RUN_CODE
					runtimeSys.LogNewFun("DEBUG", "about to run a gf_lang program...", nil)

					resultsLst, _, err := gf_lang.Run(handlerProgramASTlst,
						pExternAPI)
					
					if err != nil {

						gfErr := gf_core.ErrorCreate("failed to execute gf_lang program in a rpc handler",
							"gf_lang_program_run_failed",
							map[string]interface{}{"program_ast_lst": handlerProgramASTlst,},
							err, "gf_lang", runtimeSys)
						return nil, gfErr
					}

					//---------------------
					
					outputMap := map[string]interface{}{
						"results_lst": resultsLst,
					}
					return outputMap, nil
				}

				//-------------------------------------------------

				gf_rpc_lib.CreateHandlerHTTPwithMux(h.URLpathStr,
					handlerFun,
					HTTPmux,
					nil,
					false, // pStoreRunBool
					nil,
					runtimeSys)
			}

			//-------------
			// SERVER_INIT - blocking
			gf_rpc_lib.ServerInitWithMux("gf_lang", appServerPortInt, HTTPmux)

			//-------------
		},

		//---------------------------------------------
	}
	
	// RUN
	_, programsDebugLst, err := gf_lang.Run(programASTlst,
		externAPI)
	
	if err != nil {
		panic(err)
	}


	// OUTPUT->FILE
	filePathStr := "serialized_output.json"

	gfErr = debugSerializeOutputToFile(filePathStr,
		programsDebugLst,
		runtimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}


	// STATE_HISTORY->FILE
	filePathStr = "state_history.json"
	debugSerializeStateHistoryToFile(filePathStr,
		programsDebugLst,
		runtimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}


	//-------------
	// DEBUG_ANALYZER
	gfErr = debugRunPyAnalyzer(programsDebugLst,
		plugins,
		runtimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	fmt.Println("debug processing done...")

	//-------------
}

//-------------------------------------------------

