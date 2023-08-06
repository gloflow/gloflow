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
	"fmt"
	"time"
	"context"
	"strings"
	"net/http"
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

		//-----------
		// ID
		// unique identifier for the token, often referred to as the "JWT ID".
		// optional and can be used to uniquely identify a specific token.
		// way to prevent token replay attacks, where an attacker tries to reuse a previously issued token
		
		// ADD!! - pass in a GF session ID here to be used as the unique JWT ID,
		//         to be able to reference sessions with JWT tokens directly. 
		"jti": "", // ID

		//-----------
		
		"iat": int(creationUNIXtimeF), // issued_at
		"iss": issuerStr,              // issuer
		"nbf": int(creationUNIXtimeF), // not_before
		"sub": pUserIdentifierStr,     // subject

		//----------------------
	}
	
	pRuntimeSys.LogNewFun("DEBUG", "claims created for new generated JWT", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(claimsMap)
	}

	//----------------------

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claimsMap))

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
	
	pRuntimeSys.LogNewFun("DEBUG", "validating JWT token...", map[string]interface{}{
		"auth_subsystem_type": pAuthSubsystemTypeStr,
	})

	// KEY_SERVER
	publicKey, gfErr := KSclientJWTgetValidationKey(pAuthSubsystemTypeStr,
		pKeyServerInfo, pRuntimeSys)
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
				"jwt_token_val_str":       pJWTtokenVal,
				"auth_subsystem_type_str": pAuthSubsystemTypeStr,
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

	//-------------------------
	// JWT_PARSE

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

	//-------------------------
	
	pRuntimeSys.LogNewFun("DEBUG", "token validation has been executed...", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(jwtToken)
	}

	validBool := jwtToken.Valid
	
	//-------------------------
	// USER_IDENTIFIER
	
	var userIdentifierStr string

	if userIdentifierClaimStr, ok := jwtToken.Claims.(jwt.MapClaims)["sub"]; ok {
		userIdentifierStr = userIdentifierClaimStr.(string)
	} else {
		gfErr := gf_core.ErrorCreate("validated JWT token is missing an expected 'sub' claim",
			"crypto_jwt_verify_token_error",
			map[string]interface{}{
				"jwt_token_val_str": pJWTtokenVal,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, "", gfErr
	}

	//-------------------------

	pRuntimeSys.LogNewFun("DEBUG", "validated JWT token", map[string]interface{}{"valid_bool": validBool,})

	return validBool, userIdentifierStr, nil
}

//-------------------------------------------------------------

// extract JWT token from a http request and return it as a string
func JWTgetTokenFromRequest(pReq *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (string, bool, *gf_core.GFerror) {
	
	// AUTHORIZATION_HEADER - set by a GF http handler /v1/identity/auth0/login_callback
	//                        on successful completion of login at the end of the handler.
	//                        this is the standard Oauth2 header symbol.

	cookieNameStr := "Authorization"
	cookieFoundBool, cookieValueStr := gf_core.HTTPgetCookieFromReq(cookieNameStr, pReq, pRuntimeSys)
	
	pRuntimeSys.LogNewFun("DEBUG", `auth0 Authorization cookie fetch attempt from incoming request...`,
		map[string]interface{}{
			"cookie_found_bool": cookieFoundBool,
		})
		
    // authCookie, err := pReq.Cookie("Authorization")
    if !cookieFoundBool {
		return "", false, nil
	}

	authHeaderStr := cookieValueStr

	// remove the "Bearer" header in the token
    authPartsLst := strings.Split(authHeaderStr, " ")
    if len(authPartsLst) != 2 || strings.ToLower(authPartsLst[0]) != "bearer" {
		gfErr := gf_core.ErrorCreate("Authorization cookie is not in a valid format (not composed of 2 components, starting with 'Bearer ...')",
			"http_cookie",
			map[string]interface{}{
				"path_str":        pReq.URL.Path,
				"auth_header_str": authHeaderStr,
			},
			nil, "gf_auth0", pRuntimeSys)
		return "", false, gfErr
    }

    return authPartsLst[1], true, nil
}