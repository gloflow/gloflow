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

package gf_address

import (
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func httpIputForAdd(pUserIDstr gf_core.GF_ID,
	pReq        *http.Request,
	pResp       http.ResponseWriter,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFaddInput, *gf_core.GFerror) {

	inputMap, gfErr := gf_core.HTTPgetInput(pResp, pReq, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	var addressStr string 
	if valStr, ok := inputMap["address_str"]; ok {
		addressStr = valStr.(string)
	}

	var typeStr string // "my"|"observed"
	if valStr, ok := inputMap["type_str"]; ok {
		typeStr = valStr.(string)
	}

	var chainStr string // "eth"|"tezos"
	if valStr, ok := inputMap["chain_str"]; ok {
		chainStr = valStr.(string)
	}
	
	input := &GFaddInput{
		UserIDstr:  pUserIDstr,
		AddressStr: addressStr,
		TypeStr:    typeStr, 
		ChainStr:   chainStr,	
	}

	return input, nil
}

//---------------------------------------------------
func httpIputForGetAll(pUserIDstr gf_core.GF_ID,
	pReq        *http.Request,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFgetAllInput, *gf_core.GFerror) {

	queryArgsMap := pReq.URL.Query()

	// TYPE
	var typeStr string // "my"|"observed"
	if valuesLst, ok := queryArgsMap["type"]; ok {
		typeStr = valuesLst[0]
	} else {
		gfErr := gf_core.Error__create("incoming http request is missing the 'type' query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_address", pRuntimeSys)
		return nil, gfErr
	}

	// CHAIN
	var chainStr string // "eth"|"tezos"
	if valuesLst, ok := queryArgsMap["chain"]; ok {
		chainStr = valuesLst[0]
	} else {
		gfErr := gf_core.Error__create("incoming http request is missing the 'chain' query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_address", pRuntimeSys)
		return nil, gfErr
	}

	input := &GFgetAllInput {
		UserIDstr: pUserIDstr,
		TypeStr:   typeStr, 
		ChainStr:  chainStr,	
	}

	return input, nil
}