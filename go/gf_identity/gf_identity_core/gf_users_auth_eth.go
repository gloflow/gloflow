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

package gf_identity_core

import (
	"time"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

// io_preflight
type GFethInputPreflight struct {
	UserAddressETH GFuserAddressETH `validate:"omitempty,eth_addr"`
}
type GFethOutputPreflight struct {
	UserExistsBool bool             
	NonceValStr    GFuserNonceVal
}

// io_login
type GFethInputLogin struct {
	UserAddressETH   GFuserAddressETH `validate:"required,eth_addr"`
	AuthSignatureStr GFauthSignature  `validate:"required,len=132"` // singature length with "0x"
	AudienceStr      string           `validate:"len=100"`
}
type GFethOutputLogin struct {
	NonceExistsBool        bool
	AuthSignatureValidBool bool
	JWTtokenVal            GFjwtTokenVal
	UserID                 gf_core.GF_ID
	SessionID              gf_core.GF_ID
}

// io_create
type GFethInputCreate struct {
	UserTypeStr      string           `validate:"required"` // "admin" | "standard"
	UserAddressETH   GFuserAddressETH `validate:"required,eth_addr"`
	AuthSignatureStr GFauthSignature  `validate:"required,len=132"` // singature length with "0x"
}
type GFethOutputCreate struct {
	NonceExistsBool        bool
	AuthSignatureValidBool bool
}

//---------------------------------------------------

func ETHpipelinePreflight(pInput *GFethInputPreflight,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFethOutputPreflight, *gf_core.GFerror) {

	//------------------------
	// VALIDATE
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	output := &GFethOutputPreflight{}

	existsBool, gfErr := dbSQLuserExistsByETHaddr(pInput.UserAddressETH,
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
			pInput.UserAddressETH,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		output.UserExistsBool = false
		output.NonceValStr    = nonce.ValStr

	// user exists
	} else {

		nonceValStr, nonceExistsBool, gfErr := dbNonceGet(pInput.UserAddressETH, pCtx, pRuntimeSys)
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
// PIPELINE_LOGIN

func ETHpipelineLogin(pInput *GFethInputLogin,
	pKeyServerInfo *GFkeyServerInfo,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (*GFethOutputLogin, *gf_core.GFerror) {
	
	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	output := &GFethOutputLogin{}

	//------------------------
	// SESSION_ID
	sessionID := generateSessionID()
	output.SessionID = sessionID

	//------------------------
	userNonceVal, userNonceExistsBool, gfErr := dbNonceGet(pInput.UserAddressETH,
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
		pInput.UserAddressETH,
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

	userID, gfErr := DBsqlGetBasicInfoByETHaddr(pInput.UserAddressETH,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.UserID = userID

	//------------------------
	// JWT
	userIdentifierStr    := string(userID)
	authSubsystemTypeStr := GF_AUTH_SUBSYSTEM_TYPE__ETH
	jwtTokenVal, gfErr   := JWTpipelineGenerate(userIdentifierStr,
		authSubsystemTypeStr,
		pInput.AudienceStr,
		pKeyServerInfo,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.JWTtokenVal = *jwtTokenVal

	//------------------------

	return output, nil
}

//---------------------------------------------------
// PIPELINE_CREATE

func ETHpipelineCreate(pInput *GFethInputCreate,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFethOutputCreate, *gf_core.GFerror) {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	output := &GFethOutputCreate{}
	
	//------------------------
	// DB_NONCE_GET - get a nonce already generated in preflight for this user address,
	//                for validating the recevied auth_signature
	userNonceValStr, userNonceExistsBool, gfErr := dbNonceGet(pInput.UserAddressETH,
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
		pInput.UserAddressETH,
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
	UserAddressETH      := pInput.UserAddressETH
	userAddressesETHlst := []GFuserAddressETH{UserAddressETH, }

	userIdentifierStr := string(UserAddressETH)
	userID := usersCreateID(userIdentifierStr, creationUNIXtimeF)

	user := &GFuser{
		Vstr:              "0",
		ID:                userID,
		CreationUNIXtimeF: creationUNIXtimeF,
		UserTypeStr:       pInput.UserTypeStr,
		AddressesETHlst:   userAddressesETHlst,
	}

	//------------------------
	// DB
	gfErr = DBsqlUserCreate(user, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	return output, nil
}