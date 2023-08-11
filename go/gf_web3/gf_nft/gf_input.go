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

package gf_nft

import (
	"net/http"
	"context"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// INDEX_ADDRESS

func httpInputForIndexAddress(pUserID gf_core.GF_ID,
	pReq        *http.Request,
	pResp       http.ResponseWriter,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFindexAddressInput, *gf_core.GFerror) {

	inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	var addressStr string
	if valStr, ok := inputMap["address_str"]; ok {
		addressStr = strings.ToLower(valStr.(string))
	}

	var chainStr string
	if valStr, ok := inputMap["chain_str"]; ok {
		chainStr = strings.ToLower(valStr.(string))
	}

	var fetcherNameStr string
	if valStr, ok := inputMap["fetcher_name_str"]; ok {
		if valStr.(string) == "" {
			// default value
			fetcherNameStr = "alchemy"
		} else {
			fetcherNameStr = valStr.(string)
		}
	}

	input := &GFindexAddressInput{
		UserID:         pUserID,
		AddressStr:     addressStr,
		ChainStr:       chainStr,
		FetcherNameStr: fetcherNameStr,
	}
	return input, nil
}

//-------------------------------------------------
// GET

func httpInputForGetByOwner(pUserID gf_core.GF_ID,
	pReq        *http.Request,
	pResp       http.ResponseWriter,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFgetByOwnerInput, *gf_core.GFerror) {

	inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	var addressStr string
	if valStr, ok := inputMap["address_str"]; ok {
		addressStr = strings.ToLower(valStr.(string))
	}

	var chainStr string
	if valStr, ok := inputMap["chain_str"]; ok {
		chainStr = strings.ToLower(valStr.(string))
	}

	input := &GFgetByOwnerInput{
		UserID:     pUserID,
		AddressStr: addressStr,
		ChainStr:   chainStr,
	}
	return input, nil
}

//-------------------------------------------------
// GET

func httpInputForGet(pUserID gf_core.GF_ID,
	pReq        *http.Request,
	pResp       http.ResponseWriter,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFgetInput, *gf_core.GFerror) {

	inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	var tokenIDstr string
	if valStr, ok := inputMap["token_id_str"]; ok {
		tokenIDstr = valStr.(string)
	}

	var collectionNameStr string
	if valStr, ok := inputMap["collection_name_str"]; ok {
		collectionNameStr = valStr.(string)
	}
	
	input := &GFgetInput{
		UserID:            pUserID,
		TokenIDstr:        tokenIDstr,
		CollectionNameStr: collectionNameStr,
	}
	return input, nil
}