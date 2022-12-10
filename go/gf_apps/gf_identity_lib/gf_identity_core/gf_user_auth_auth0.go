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
	// "fmt"
	"time"
	"context"
	"strings"
	// "github.com/auth0/go-jwt-middleware/v2"
	// "github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	"github.com/davecgh/go-spew/spew"
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
	CodeStr                     string
	GFsessionIDauth0providedStr gf_core.GF_ID
}
type GFauth0outputLoginCallback struct {
	SessionIDstr gf_core.GF_ID
}

//---------------------------------------------------

func Auth0validateSession(pSessionIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {


	session, gfErr := dbAuth0GetSession(pSessionIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	if !session.LoginCompleteBool {
		return false, nil
	}



	return true, nil
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
	pRuntimeSys    *gf_core.RuntimeSys) (*GFauth0outputLoginCallback, *gf_core.GFerror) {
	
	sessionIDstr := pInput.GFsessionIDauth0providedStr

	//---------------------
	// verify that the sessionID (auth0 "state") corresponds to an registered session
	// created in the previously called login handler, and that a login with that session 
	// has not already been completed

	auth0session, gfErr := dbAuth0GetSession(sessionIDstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

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
	
	//---------------------
	// exchange an authorization code for a token.
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

	//---------------------
	// verify token
	idToken, gfErr := gf_auth0.VerifyIDtoken(oauth2bearerToken,
		pAuthenticator,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "Auth0 verified openID ID token", nil)

	//---------------------

	accessTokenStr := oauth2bearerToken.AccessToken

	var profileMap map[string]interface{}
	if err := idToken.Claims(&profileMap); err != nil {
		gfErr := gf_core.ErrorCreate("failed to verify ID Token",
			"library_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "parsed user profile from openID id_token", nil)
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

	// mark the session as successfuly logged in, so that the login_callback handler
	// cant be invoked again
	loginCompleteBool := true

	//---------------------
	// DB
	gfErr = dbAuth0UpdateSession(sessionIDstr,
		loginCompleteBool,
		accessTokenStr,
		profileMap,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//---------------------

	output := &GFauth0outputLoginCallback{
		SessionIDstr: sessionIDstr,
	}
	return output, nil
}