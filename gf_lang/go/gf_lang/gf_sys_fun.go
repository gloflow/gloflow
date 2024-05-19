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
	"math/rand"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

func getSysFunctionNames() []string {
	return []string{
        "make",      // make a datastructure - list|map
		"len",       // length of a collection
        "rand",      // random number generator,
        "rpc_call",  // remot procedure call
		"rpc_serve", // rpc server definition/startup
    }
}

//-------------------------------------------------

func isSysFunc(pExprLst []interface{}) bool {
    sysFunsLst := getSymbolsAndConstants().SystemFunctionsLst
    firstElementIsStrBool, firstElementStr := gf_core.CastToStr(pExprLst[0])
    if firstElementIsStrBool && gf_core.ListContainsStr(firstElementStr, sysFunsLst) {
        return true
    }
    return false
}

//-------------------------------------------------

func sysFuncEval(pExprLst []interface{},
    pState     *GFstate,
    pExternAPI GFexternAPI) (interface{}, error) {
    
    funcNameStr := pExprLst[0].(string)
    argsExprLst := pExprLst[1].(GFexpr)
    
    var val interface{}
    switch funcNameStr {

    //---------------------
    // MAKE - create a datastructure
    case "make":

        if len(argsExprLst) != 2 {
			return nil, errors.New("'make' system function can only be called with 2 arg, 'list|map' and initial values or list length")
		}

        //---------------------
        // MAKE_TYPE
        var makeTypeStr string
        if typeStr, ok := argsExprLst[0].(string); true {

            if !ok {
                return nil, errors.New("first 'make' arg has to be a string, indicating the 'make' type")
            }
            if typeStr != "list" && typeStr != "map" {
                return nil, errors.New("'make' type can only be 'list'|'map'")
            }
            makeTypeStr = typeStr
        }

        //---------------------
        
        switch makeTypeStr {
        
        //---------------------
        // LIST
        case "list":

            var newLst []interface{}
            switch v := argsExprLst[1].(type) {
            
            // new list 'make' command specifies a length to be created, of an empty string
            case int:
                listLenInt := v
                newLst = make([]interface{}, listLenInt)

            // initial values for the new list are specified
            case GFexpr:
                initialValuesLst := []interface{}(v)
                newLst = initialValuesLst

            default:
                return nil, errors.New("unsupported type for the second argument in the 'make' sys func for creating a list")
            }
            val = newLst

        //---------------------
        // MAP
        case "map":
            newMap := map[string]interface{}{}
            initValuesLst := []interface{}(argsExprLst[1].(GFexpr))
            for _, kv := range initValuesLst {
                kvLst := []interface{}(kv.(GFexpr))
                kSstr := kvLst[0].(string)
                v     := kvLst[1].(interface{})

                newMap[kSstr] = v
            }
            val = newMap

        //---------------------
        }

	//---------------------
	// LEN - get a length of a collection
	case "len":

		if len(argsExprLst) > 1 {
			return nil, errors.New("'len' system function can only be called with 1 arg, which is a collection")
		}

		// argument should only be one, and a variable string reference
		varStr := argsExprLst[0].(string)
		varValue, err := varEval(varStr, pState)
		if err != nil {
			return nil, err
		}

        if _, ok := varValue.Val.([]interface{}); !ok {
            return nil, errors.New(fmt.Sprintf("variable %s in expression %s is not a collection, so cant call 'len' for it",
                varStr,
                pExprLst))
        }

        // the value in the variable is expected to be a collection
		collLengthInt := len(varValue.Val.([]interface{}))

		val = interface{}(collLengthInt)

    //---------------------
    // RAND
    case "rand":
        if len(argsExprLst) != 2 {
            return nil, errors.New("'rand' system function only takes 2 argument")
        }

        randomRangeMinF := argsExprLst[0].(float64)
        randomRangeMaxF := argsExprLst[1].(float64)
        valF := rand.Float64()*(randomRangeMaxF - randomRangeMinF) + randomRangeMinF
        val = interface{}(valF)

    //---------------------
    // RPC_CALL
    case "rpc_call":
        
        responseMap, err := rpcCallEval(argsExprLst, pState, pExternAPI)
        if err != nil {
            return nil, err
        }

        fmt.Println(responseMap)

    //---------------------
    // RPC_SERVE
    case "rpc_serve":

        err := rpcServeEval(argsExprLst, pState, pExternAPI)
        if err != nil {
            return nil, err
        }

    //---------------------
    }

    return val, nil
}