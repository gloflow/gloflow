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
)

//---------------------------------------------------

type GFgeometryFunc func(float64, float64, float64,
    float64, float64, float64,
    float64, float64, float64,
    float64, float64, float64)

type GFexternAPI struct {
    InitEngineFun   func(map[string]*GFshaderDef)
    SetStateFun     func(GFstateChange) []interface{}
    CreateCubeFun   GFgeometryFunc
    CreateSphereFun GFgeometryFunc
    CreateLineFun   GFgeometryFunc
    AnimateFun      func([]map[string]interface{}, float64, bool)

    //------------------------------------
    // RPC
    // RPC_CALL
    RPCcall GFrpcCallFun

    // RPC_SERVE
    RPCserve GFrpcServeFun

    //------------------------------------
}

//---------------------------------------------------

func GetTestExternAPI() GFexternAPI {
	externAPI := GFexternAPI{

		InitEngineFun: func(pShaderDefsMap map[string]*GFshaderDef) {
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