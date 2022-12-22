/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	// "fmt"
	"time"
	"context"
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

type GFjwtTokenVal string

//---------------------------------------------------
// GENERATE
//---------------------------------------------------

func JWTpipelineGenerate(pUserIdentifierStr string,
	pAuthSubsystemTypeStr string,
	pKeyServerInfo        *GFkeyServerInfo,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.RuntimeSys) (*GFjwtTokenVal, *gf_core.GFerror) {
	
	// KEY_SERVER
	privateKey, gfErr := ksClientJWTgetSigningKey(pAuthSubsystemTypeStr, pKeyServerInfo, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	signingKey := privateKey

	//----------------------
	// JWT_GENERATE
	tokenValStr, gfErr := jwtGenerate(pUserIdentifierStr,
		signingKey,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//----------------------

	return tokenValStr, nil
}

//---------------------------------------------------

func jwtGenerate(pUserIdentifierStr string,
	pSigningKey *rsa.PrivateKey,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFjwtTokenVal, *gf_core.GFerror) {
	
	pRuntimeSys.LogNewFun("DEBUG", "JWT generated for user", map[string]interface{}{
		"user_identifier_str": pUserIdentifierStr,})
		
	issuerStr := "gf"
	_, jwtTokenTTLsecInt  := GetSessionTTL()
	creationUNIXtimeF     := float64(time.Now().UnixNano())/1000000000.0
	expirationUNIXtimeInt := int64(creationUNIXtimeF) + jwtTokenTTLsecInt

	//----------------------
	// CLAIMS

	/*
	type StandardClaims struct {
		Audience  string `json:"aud,omitempty"`
		ExpiresAt int64  `json:"exp,omitempty"`
		Id        string `json:"jti,omitempty"`
		IssuedAt  int64  `json:"iat,omitempty"`
		Issuer    string `json:"iss,omitempty"`
		NotBefore int64  `json:"nbf,omitempty"`
		Subject   string `json:"sub,omitempty"`
	}
	*/
	claimsMap := map[string]interface{}{

		//----------------------
		// standard claims
		"aud": "",                     // audience
		"exp": expirationUNIXtimeInt,  // expires_at
		"jti": "",                     // id
		"iat": int(creationUNIXtimeF), // issued_at
		"iss": issuerStr,              // issuer
		"nbf": int(creationUNIXtimeF), // not_before
		"sub": "",                     // subject

		//----------------------
		// GF claims
		"user_identifier_str": pUserIdentifierStr,
		
		//----------------------
	}

	// claims := token.Claims.(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(10 * time.Minute)
	// claims["authorized"] = true
	// claims["user"] = "username"
	
	pRuntimeSys.LogNewFun("DEBUG", "claims create for new generated JWT", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(claimsMap)
	}

	//----------------------

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims(claimsMap))

	jwtTokenSignedStr, err := jwtToken.SignedString(pSigningKey)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to sign JWT token for user",
			"crypto_jwt_sign_token_error",
			map[string]interface{}{
				"user_identifier_str": pUserIdentifierStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}


	tokenValStr := GFjwtTokenVal(jwtTokenSignedStr)
	return &tokenValStr, nil
}

//---------------------------------------------------
// VALIDATE
//---------------------------------------------------

func JWTpipelineValidate(pJWTtokenVal GFjwtTokenVal,
	pAuthSubsystemTypeStr string,
	pKeyServerInfo        *GFkeyServerInfo,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	pRuntimeSys.LogNewFun("DEBUG", "validating JWT token...", nil)

	// KEY_SERVER
	publicKey, gfErr := ksClientJWTgetValidationKey(pAuthSubsystemTypeStr, pKeyServerInfo, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}
	
	// VALIDATE
	validBool, userIdentifierStr, gfErr := JWTvalidate(pJWTtokenVal,
		publicKey,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	if !validBool {
		gfErr := gf_core.ErrorCreate("JWT token supplied for validation is invalid",
			"crypto_jwt_verify_token_invalid_error",
			map[string]interface{}{
				"jwt_token_val_str": pJWTtokenVal,
			},
			nil, "gf_identity_core", pRuntimeSys)
		return "", gfErr
	}

	return userIdentifierStr, nil
}

//---------------------------------------------------

func JWTvalidate(pJWTtokenVal GFjwtTokenVal,
	pPublicKey  *rsa.PublicKey,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {

	// token validation
	jwtToken, err := jwt.Parse(string(pJWTtokenVal), func(pToken *jwt.Token) (interface{}, error) {

		return pPublicKey, nil
	})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to verify a JWT token",
			"crypto_jwt_verify_token_error",
			map[string]interface{}{
				"jwt_token_val_str": pJWTtokenVal,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, "", gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "token validation has been executed...", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(jwtToken)
	}

	validBool := jwtToken.Valid
	var userIdentifierStr string

	if userIdentifierClaimStr, ok := jwtToken.Claims.(jwt.MapClaims)["user_identifier_str"]; ok {
		userIdentifierStr = userIdentifierClaimStr.(string)
	} else {
		gfErr := gf_core.ErrorCreate("validated JWT token is missing an expected 'user_identifier_str' claim",
			"crypto_jwt_verify_token_error",
			map[string]interface{}{
				"jwt_token_val_str": pJWTtokenVal,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, "", gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "validated JWT token", map[string]interface{}{"valid_bool": validBool,})

	return validBool, userIdentifierStr, nil
}