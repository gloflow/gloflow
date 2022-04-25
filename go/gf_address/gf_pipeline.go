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

import(
	"time"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GFgetAllInput struct {
	UserIDstr gf_core.GF_ID
	TypeStr   string 
	ChainStr  string	
}
type GFgetAllOutput struct {
	AddressesLst []string
}

type GFaddInput struct {
	UserIDstr  gf_core.GF_ID
	AddressStr string
	TypeStr    string 
	ChainStr   string	
}

//-------------------------------------------------
// PIPELINE_GET_ALL
func pipelineGetAll(pInput *GFgetAllInput,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (*GFgetAllOutput, *gf_core.GFerror) {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.Validate_struct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	// DB
	addressesLst, gfErr := DBgetAll(pInput.TypeStr,
		pInput.ChainStr,
		pInput.UserIDstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	addressesExportLst := []string{}
	for _, a := range addressesLst {
		addressesExportLst = append(addressesExportLst, a.AddressStr)
	}

	output := &GFgetAllOutput{
		AddressesLst: addressesExportLst,
	}

	return output, nil
}

//-------------------------------------------------
// PIPELINE_ADD
func pipelineAdd(pInput *GFaddInput,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GFerror {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.Validate_struct(pInput, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------
	

	creationTimeUNIXf := float64(time.Now().UnixNano()) / 1_000_000_000.0
	idStr   := createID(string(pInput.UserIDstr), creationTimeUNIXf)
	address := &GFchainAddress {
		Vstr:              "0",
		IDstr:             idStr,
		CreationUNIXtimeF: creationTimeUNIXf,
	
		OwnerUserIDstr: pInput.UserIDstr,
		AddressStr:     pInput.AddressStr,
		TypeStr:        pInput.TypeStr,
		ChainNameStr:   pInput.ChainStr,
	}

	// DB
	gfErr = DBadd(address,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	return nil
}

//-------------------------------------------------