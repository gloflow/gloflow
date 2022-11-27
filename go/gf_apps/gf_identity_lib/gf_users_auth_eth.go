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

package gf_identity_lib

import (
	"time"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
)

//---------------------------------------------------

// io_preflight
type GFuserAuthETHinputPreflight struct {
	UserAddressETHstr gf_identity_core.GFuserAddressETH `validate:"omitempty,eth_addr"`
}
type GFuserAuthETHoutputPreflight struct {
	UserExistsBool bool             
	NonceValStr    GFuserNonceVal
}

// io_login
type GFuserAuthETHinputLogin struct {
	UserAddressETHstr gf_identity_core.GFuserAddressETH `validate:"required,eth_addr"`
	AuthSignatureStr  gf_identity_core.GFauthSignature  `validate:"required,len=132"` // singature length with "0x"
}
type GFuserAuthETHoutputLogin struct {
	NonceExistsBool        bool
	AuthSignatureValidBool bool
	JWTtokenVal            gf_session.GFjwtTokenVal
	UserIDstr              gf_core.GF_ID 
}

// io_create
type GFuserAuthETHinputCreate struct {
	UserTypeStr       string                            `validate:"required"` // "admin" | "standard"
	UserAddressETHstr gf_identity_core.GFuserAddressETH `validate:"required,eth_addr"`
	AuthSignatureStr  gf_identity_core.GFauthSignature  `validate:"required,len=132"` // singature length with "0x"
}
type GFuserAuthETHoutputCreate struct {
	NonceExistsBool        bool
	AuthSignatureValidBool bool
}

//---------------------------------------------------

func usersAuthETHpipelinePreflight(pInput *GFuserAuthETHinputPreflight,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserAuthETHoutputPreflight, *gf_core.GFerror) {

	//------------------------
	// VALIDATE
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	output := &GFuserAuthETHoutputPreflight{}

	existsBool, gfErr := dbUserExistsByETHaddr(pInput.UserAddressETHstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// no user exists so create a new nonce
	if !existsBool {

		// user doesnt exist yet so no user_id
		userIDstr := gf_core.GF_ID("")
		nonce, gfErr := nonceCreateAndPersist(userIDstr,
			pInput.UserAddressETHstr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		output.UserExistsBool = false
		output.NonceValStr    = nonce.ValStr

	// user exists
	} else {

		nonceValStr, nonceExistsBool, gfErr := dbNonceGet(pInput.UserAddressETHstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		if !nonceExistsBool {
			// generate new nonce, because the old one has been invalidated?
		} else {
			output.UserExistsBool = true
			output.NonceValStr    = nonceValStr
		}
	}

	return output, nil
}

//---------------------------------------------------
// PIPELINE__LOGIN

func usersAuthETHpipelineLogin(pInput *GFuserAuthETHinputLogin,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserAuthETHoutputLogin, *gf_core.GFerror) {
	
	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	output := &GFuserAuthETHoutputLogin{}

	//------------------------
	userNonceVal, userNonceExistsBool, gfErr := dbNonceGet(pInput.UserAddressETHstr,
		pCtx,
		pRuntimeSys)
		
	if gfErr != nil {
		return nil, gfErr
	}
	
	if !userNonceExistsBool {
		output.NonceExistsBool = false
		return output, nil
	} else {
		output.NonceExistsBool = true
	}

	//------------------------
	// VERIFY

	signatureValidBool, gfErr := verifyAuthSignatureAllMethods(pInput.AuthSignatureStr,
		userNonceVal,
		pInput.UserAddressETHstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	if !signatureValidBool {
		output.AuthSignatureValidBool = false
		return output, nil
	} else {
		output.AuthSignatureValidBool = true
	}

	//------------------------
	// USER_ID

	userIDstr, gfErr := gf_identity_core.DBgetBasicInfoByETHaddr(pInput.UserAddressETHstr,
		pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.UserIDstr = userIDstr

	//------------------------
	// JWT
	userIdentifierStr := string(userIDstr)
	jwtTokenVal, gfErr := gf_session.JWTpipelineGenerate(userIdentifierStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.JWTtokenVal = jwtTokenVal

	//------------------------

	return output, nil
}

//---------------------------------------------------
// PIPELINE__CREATE

func usersAuthETHpipelineCreate(pInput *GFuserAuthETHinputCreate,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserAuthETHoutputCreate, *gf_core.GFerror) {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	output := &GFuserAuthETHoutputCreate{}
	
	//------------------------
	// DB_NONCE_GET - get a nonce already generated in preflight for this user address,
	//                for validating the recevied auth_signature
	userNonceValStr, userNonceExistsBool, gfErr := dbNonceGet(pInput.UserAddressETHstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	if !userNonceExistsBool {
		output.NonceExistsBool = false
		return output, nil
	} else {
		output.NonceExistsBool = true
	}

	//------------------------
	// VALIDATE

	signatureValidBool, gfErr := verifyAuthSignatureAllMethods(pInput.AuthSignatureStr,
		userNonceValStr,
		pInput.UserAddressETHstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	if signatureValidBool {
		output.AuthSignatureValidBool = true
	} else {
		output.AuthSignatureValidBool = false
		return output, nil
	}

	//------------------------

	creationUNIXtimeF   := float64(time.Now().UnixNano())/1000000000.0
	userAddressETHstr   := pInput.UserAddressETHstr
	userAddressesETHlst := []gf_identity_core.GFuserAddressETH{userAddressETHstr, }

	userIdentifierStr := string(userAddressETHstr)
	userID := usersCreateID(userIdentifierStr, creationUNIXtimeF)

	user := &GFuser{
		Vstr:              "0",
		IDstr:             userID,
		CreationUNIXtimeF: creationUNIXtimeF,
		UserTypeStr:       pInput.UserTypeStr,
		AddressesETHlst:   userAddressesETHlst,
	}

	//------------------------
	// DB
	gfErr = dbUserCreate(user, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	return output, nil
}