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
	"errors"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFrpcCallFun func(string, // node
    string,                    // module
    string,                    // function
    []interface{}) map[string]interface{} // args list

type GFrpcServeFun func(string, // node_name
	[]*GFrpcServerHandler,      // handlers
	GFexternAPI)               // extern_api

type GFrpcServerHandler struct {
	URLpathStr  string        // URL path
    ModuleStr   string        // module
    FunctionStr string        // function
    ArgsSpecLst []interface{} // args_spec
    CodeASTlst  []interface{} // code
}	

//-------------------------------------------------
// RPC_CALL

func rpcCallEval(pExprLst []interface{},
	pState     *GFstate,
	pExternAPI GFexternAPI) (map[string]interface{}, error) {

	nodeNameStr     := pExprLst[0].(string)
	moduleNameStr   := pExprLst[1].(string)
	functionNameStr := pExprLst[2].(string)
	argsLst         := pExprLst[3].([]interface{})

	resultMap := pExternAPI.RPCcall(nodeNameStr, moduleNameStr, functionNameStr, argsLst)

	return resultMap, nil
}

//-------------------------------------------------
// RPC_SERVE

func rpcServeEval(pExprLst []interface{},
	pState     *GFstate,
	pExternAPI GFexternAPI) error {

	nodeNameStr      := pExprLst[0].(string)        // node_name
	handlersBlockLst := pExprLst[1].([]interface{}) // handlers

	if _, ok := handlersBlockLst[0].(string); !ok {
		return errors.New(fmt.Sprintf("node name (%s) has to be a string", nodeNameStr))
	}

	if handlersBlockLst[0].(string) != "handlers" {
		return errors.New(fmt.Sprintf("first element has to be a string 'handlers', not %s", handlersBlockLst[0]))
	}

	handlersExprsLst := handlersBlockLst[1].([]interface{})

	//---------------------
	// LOAD_HANDLERS
	handlersLst, err := loadHandlers(handlersExprsLst)
	if err != nil {
		return err
	}

	//---------------------

	pExternAPI.RPCserve(nodeNameStr, handlersLst, pExternAPI)

	return nil
}

//-------------------------------------------------
// LOAD_HANDLERS

func loadHandlers(pHandlersExprsLst []interface{}) ([]*GFrpcServerHandler, error) {
	spew.Dump(pHandlersExprsLst)

	handlersLst := []*GFrpcServerHandler{}
	for _, h := range pHandlersExprsLst {

		handlerExprLst       := h.([]interface{})
		handlerURLstr        := handlerExprLst[0].(string)
		handlerModuleNameStr := handlerExprLst[1].(string)
		handlerFunctionStr   := handlerExprLst[2].(string)

		// code
		handlerCodeBlockLst := handlerExprLst[4].([]interface{})

		if len(handlerCodeBlockLst) != 2 {
			return nil, errors.New("handler code definition block can only have length 2 - ['code', [...code_blocks...]]")
		}
		if handlerCodeBlockLst[0].(string) != "code" {
			return nil, errors.New("first element of the handler code block has to be a string 'code'")
		}
		handlerCodeASTlst := handlerCodeBlockLst[1].([]interface{})


		handler := &GFrpcServerHandler{
			URLpathStr:  handlerURLstr,
			ModuleStr:   handlerModuleNameStr,
			FunctionStr: handlerFunctionStr,
			CodeASTlst:  handlerCodeASTlst,
		}
		handlersLst = append(handlersLst, handler)
	}
	return handlersLst, nil
}

//-------------------------------------------------