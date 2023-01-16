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
	"regexp"
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

func main() {

	fmt.Println("\n   GF_LANG SERVER >>\n")

	serverPortInt := 5000
	
	logFun, logNewFun := gf_core.LogsInitNew(true, "DEBUG")
	log.SetOutput(os.Stdout)

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_lang_server",
		EnvStr:         "dev",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
	}

	//--------------------------
	localTestProgramStr := "./tests/rpc_test.gf"
	programASTlst, gfErr := ParseProgramASTfromFile(localTestProgramStr, runtimeSys)
	if gfErr != nil {
		panic(1)
	}

	//--------------------------

	externAPI := gf_lang.GFexternAPI{

		InitEngineFun: func(pShaderDefsMap map[string]interface{}) {
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

					resultsLst, err := gf_lang.Run(handlerProgramASTlst,
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

					
				fmt.Println(handlerFun)
			}

			//-------------
			// SERVER_INIT - blocking
			gf_rpc_lib.ServerInitWithMux("gf_lang", serverPortInt, HTTPmux)

			//-------------
		},

		//---------------------------------------------
	}
	
	// RUN
	_, err := gf_lang.Run(programASTlst,
		externAPI)
	
	if err != nil {
		panic(err)
	}

}

//-------------------------------------------------

func ParseProgramASTfromFile(pLocalFilePathStr string,
	pRuntimeSys *gf_core.RuntimeSys) ([]interface{}, *gf_core.GFerror) {

	//------------------------
	// READ_FILE
	programCodeStr, gfErr := gf_core.FileRead(pLocalFilePathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// REMOVE_COMMENTS
	commentRegex := regexp.MustCompile(`(?m)(.*)//.*$`)

	cleanJSONcodeStr := ""
	for _, lineStr := range strings.Split(programCodeStr, "\n") {

		lineNoCommentsStr := commentRegex.ReplaceAllString(lineStr, "$1")

		if strings.TrimSpace(lineNoCommentsStr) != "" {
			cleanJSONcodeStr += fmt.Sprintf("%s\n", lineNoCommentsStr)
		}
	}

	fmt.Println("clean JSON code:", cleanJSONcodeStr)

	//------------------------
	// PARSE

	code, gfErr := gf_core.ParseJSONfromString(cleanJSONcodeStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	codeLst := code.([]interface{})

	//------------------------
    
	return codeLst, nil
}