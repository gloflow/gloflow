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
	"fmt"
	"time"
	"context"
	"strings"
	"github.com/golang-jwt/jwt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	"github.com/davecgh/go-spew/spew"
	// "github.com/auth0/go-jwt-middleware/v2"
	// "github.com/auth0/go-jwt-middleware/v2/validator"
)

//---------------------------------------------------

type GFauth0session struct {
	IDstr             gf_core.GF_ID          `bson:"id_str"`
	DeletedBool       bool                   `bson:"deleted_bool"`
	CreationUNIXtimeF float64                `bson:"creation_unix_time_f"`

	// marked as true once the login completes (once Auth0 initial auth returns the user to the GF system).
	// if the login_callback handler is called and this login_complete is already marked as true,
	// the http transaction will be immediatelly aborted.
	LoginCompleteBool bool                   `bson:"login_complete_bool"`

	AccessTokenStr    string                 `bson:"access_token_str"`
	ProfileMap        map[string]interface{} `bson:"profile_map"`
}

type GFauth0inputLoginCallback struct {
	CodeStr           string
	GFsessionIDstr    gf_core.GF_ID
	Auth0appDomainStr string
}
type GFauth0outputLoginCallback struct {
	SessionIDstr gf_core.GF_ID
	JWTtokenStr  string     
}

//---------------------------------------------------
// CREATE_GF_USER_IF_NONE

// check if the Auth0 user exists in the DB, and if not create it.
// a user would not exist in the DB if it signed-up/logged-in for the first time.
func Auth0createGFuserIfNone(pAuth0accessTokenStr string,
	pAuth0appDomainStr string,
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {


	// GET_USER_INFO - from Auth0

	auth0userInfoMap, gfErr := gf_auth0.GetUserInfo(pAuth0accessTokenStr,
		pAuth0appDomainStr,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}



	pRuntimeSys.LogNewFun("DEBUG", `>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
		Auth0 /userinfo response recieved (for fetching user info for the current user)...`,
		map[string]interface{}{
			"auth0_user_info_map": auth0userInfoMap,
		})

	// the user_info returned by Auth0 contains the "sub" claim
	// which is the user_id assigned to the user in the Auth0 system
	auth0userID := gf_core.GF_ID(auth0userInfoMap["sub"].(string))


	//---------------------
	// DB
	existsBool, gfErr := DBuserExistsByID(auth0userID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------

	// user doesnt exist in the GF DB
	if !existsBool {

	}

	return nil
}

//---------------------------------------------------
// LOGIN

func Auth0loginPipeline(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	//---------------------
	// SESSION_ID
	sessionIDstr := generateSessionID()
	
	//---------------------

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	auth0session := &GFauth0session{
		IDstr:             sessionIDstr,
		CreationUNIXtimeF: creationUNIXtimeF,

		// indicate if the user already passed the initial login process,
		// and is now logged in.
		// this is a new Auth0 session, so the login is marked as not-complete.
		LoginCompleteBool: false,
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
// LOGIN_CALLBACK

func Auth0loginCallbackPipeline(pInput *GFauth0inputLoginCallback,
	pAuthenticator *gf_auth0.GFauthenticator,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (*GFauth0outputLoginCallback, *gf_core.GFerror) {
	
	sessionIDstr := pInput.GFsessionIDstr

	//---------------------
	// DB_GET_SESSION
	// verify that the sessionID (auth0 "state") corresponds to an registered session
	// created in the previously called login handler, and that a login with that session 
	// has not already been completed

	auth0session, gfErr := dbAuth0GetSession(sessionIDstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//---------------------
	
	/*
	// user is already logged in.
	if auth0session.LoginCompleteBool {
		gfErr := gf_core.ErrorCreate("'state' input argument supplied is invalid, it has already been used by the user to login",
			"verify__invalid_value_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	if !loginAttemptCheckAgeIsValid(auth0session.CreationUNIXtimeF) {
		gfErr := gf_core.ErrorCreate("'state' input argument supplied is invalid, too long has passed since it was created and it has expired",
			"verify__invalid_value_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}
	*/
	
	//---------------------
	// USER_NOT_LOGGED_IN
	// user has not already logged in, so check some things first.
	// a user can be already logged in, and this redirect (login-callback) by the Auth0 system was done in an already logged in state.
	if !auth0session.LoginCompleteBool {

		// check max-age of a login_attempt
		if !loginAttemptCheckAgeIsValid(auth0session.CreationUNIXtimeF) {
			gfErr := gf_core.ErrorCreate("'state' input argument supplied is invalid, too long has passed since it was created and it has expired",
				"verify__invalid_value_error",
				map[string]interface{}{},
				nil, "gf_identity_core", pRuntimeSys)
			return nil, gfErr
		}
	}
	
	//---------------------
	// EXCHANGE_CODE_FOR_TOKEN
	// IMPORTANT!! - the code parameter returned by Auth0 supplied as an HTTP QS argument "code".
	//               this auth code gets exchanged with Auth0 servers for an Oauth2 bearer token.
	//
	// "...exchange the authorization code (obtained after the user authenticates 
	// and grants authorization) with the provider, which returns an OAuth2 token.
	// This token can be used to make authenticated API requests to the provider's resources."
	oauth2bearerToken, err := pAuthenticator.Exchange(pCtx, pInput.CodeStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to exchange an authorization code for a token",
			"library_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "Auth0 received Oauth2 bearer token", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(oauth2bearerToken)
	}






	// extract the ID token from OAuth2 token.
	rawIDTokenStr, ok := oauth2bearerToken.Extra("id_token").(string)
	if !ok {
		gfErr := gf_core.ErrorCreate("failed find OpenID token in returned Oauth2 Bearer token",
			"library_error",
			map[string]interface{}{
				"oauth2_bearer_token_str": fmt.Sprint(oauth2bearerToken),
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	// parse the raw ID token without validating the signature so that the JWT token can be extracted
	JWTtmpToken, _, err := new(jwt.Parser).ParseUnverified(rawIDTokenStr, &jwt.StandardClaims{})
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to parse unverified OpenID Token as JWT",
			"library_error",
			map[string]interface{}{
				"raw_id_token_str": rawIDTokenStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	JWTtokenStr := JWTtmpToken.Raw

	//---------------------
	// verify token - get ID token
	// https://auth0.com/docs/secure/tokens/id-tokens
	// ID token is used to get user profile information.
	idToken, gfErr := gf_auth0.VerifyIDtoken(oauth2bearerToken,
		pAuthenticator,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	JWTtoken := idToken

	pRuntimeSys.LogNewFun("DEBUG", "Auth0 verified openID ID_token/JWT", nil)

	//---------------------
	// EXTERNAL_USER_PROFILE
	// JWT_CLAIMS - user_profile information coming from a third-party identity provider
	//              is encoded in JWT claims section of the token.

	var profileMap map[string]interface{}
	if err := JWTtoken.Claims(&profileMap); err != nil {
		gfErr := gf_core.ErrorCreate("failed to verify ID Token",
			"library_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "parsed user profile from JWT token", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(profileMap)
	}

	// GOOGLE
	// check if the "subject" name starts with google prefix
	if strings.HasPrefix(profileMap["sub"].(string), "google-oauth2") {
		googleProfile := &GFgoogleUserProfile {
			NameStr:       profileMap["name"].(string),
			GivenNameStr:  profileMap["given_name"].(string),
			FamilyNameStr: profileMap["family_name"].(string),
			NicknameStr:   profileMap["nickname"].(string),
			LocaleStr:     profileMap["locale"].(string),
			UpdatedAtStr:  profileMap["updated_at"].(string),
			PictureURLstr: profileMap["picture"].(string),
		}

		pRuntimeSys.LogNewFun("DEBUG", "google user profile loaded...", nil)
		if gf_core.LogsIsDebugEnabled() {
			spew.Dump(googleProfile)
		}
	}

	//---------------------
	
	auth0accessTokenStr := oauth2bearerToken.AccessToken

	gfErr = Auth0createGFuserIfNone(auth0accessTokenStr,
		pInput.Auth0appDomainStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//---------------------
	// DB

	// mark the session as successfuly logged in, so that the login_callback handler
	// cant be invoked again
	loginCompleteBool := true

	gfErr = dbAuth0UpdateSession(sessionIDstr,
		loginCompleteBool,

		//---------------
		// currently active profile is stored in the DB.
		profileMap,

		//---------------
		
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//---------------------

	output := &GFauth0outputLoginCallback{
		SessionIDstr: sessionIDstr,
		JWTtokenStr:  JWTtokenStr,
	}
	return output, nil
}