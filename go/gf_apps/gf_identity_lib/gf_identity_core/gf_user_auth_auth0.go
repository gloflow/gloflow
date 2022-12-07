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
	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
)

//---------------------------------------------------

type GFauth0session struct {
	IDstr             gf_core.GF_ID          `bson:"id_str"`
	DeletedBool       bool                   `bson:"deleted_bool"`
	CreationUNIXtimeF float64                `bson:"creation_unix_time_f"`
	AccessTokenStr    string                 `bson:"access_token_str"`
	ProfileMap        map[string]interface{} `bson:"profile_map"`
}

type GFauth0inputLoginCallback struct {
	CodeStr                     string
	GFsessionIDauth0providedStr gf_core.GF_ID
	GFsessionIDstr              gf_core.GF_ID
}

//---------------------------------------------------
// USED AT ALL??

func Auth0middlewareInit(pRuntimeSys *gf_core.RuntimeSys) *jwtmiddleware.JWTMiddleware {

	//-------------------------------------------------
	keyGenerateFun := func(pCtx context.Context) (interface{}, error) {
		
		userIdentifierStr := ""
		jwtSecretKeyValStr, gfErr := JWTgenerateSecretSigningKey(userIdentifierStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr.Error
		}

		return []byte(string(jwtSecretKeyValStr)), nil
	}

	//-------------------------------------------------

	apiAudienceStr := ""
	jwtValidator, err := validator.New(
		keyGenerateFun,
		validator.HS256,
		"https://<issuer-url>/",
		[]string{apiAudienceStr,},
	)
	if err != nil {
		panic(err)
	}


	jwtAuth0middleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	return jwtAuth0middleware
}

//---------------------------------------------------

func Auth0loginPipeline(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	sessionIDstr := generateSessionID()

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	auth0session := &GFauth0session{
		IDstr:             sessionIDstr,
		CreationUNIXtimeF: creationUNIXtimeF,
	}

	//---------------------
	// DB
	gfErr := dbAuth0createNewSession(auth0session,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gf_core.GF_ID(""), gfErr
	}

	//---------------------

	return sessionIDstr, nil
}

//---------------------------------------------------

func Auth0loginCallbackPipeline(pInput *GFauth0inputLoginCallback,
	pAuthenticator *gf_auth0.GFauthenticator,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//---------------------
	// check if the GF session ID stored in the users browsers cookie is the same 
	// as the GF session ID (auth0 "state" arg) that was provided by Auth0.
	if pInput.GFsessionIDauth0providedStr != pInput.GFsessionIDstr {
		gfErr := gf_core.ErrorCreate("invalid 'state' input argument",
			"verify__input_data_missing_in_req_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return gfErr
	}
	
	//---------------------
	// exchange an authorization code for a token
	token, err := pAuthenticator.Exchange(pCtx, pInput.CodeStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to exchange an authorization code for a token",
			"library_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	//---------------------
	// verify token
	idToken, gfErr := gf_auth0.VerifyIDtoken(token,
		pAuthenticator,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------
	// this is provided by Auth0

	accessTokenStr := token.AccessToken

	var auth0profileMap map[string]interface{}
	if err := idToken.Claims(&auth0profileMap); err != nil {
		gfErr := gf_core.ErrorCreate("failed to verify ID Token",
			"library_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	//---------------------
	// DB
	gfErr = dbAuth0UpdateSession(pInput.GFsessionIDstr,
		accessTokenStr,
		auth0profileMap,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------

	return nil
}