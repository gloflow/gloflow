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
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	"github.com/davecgh/go-spew/spew"
	// "github.com/auth0/go-jwt-middleware/v2"
	// "github.com/auth0/go-jwt-middleware/v2/validator"
)

//---------------------------------------------------

type GFauth0session struct {
	ID                gf_core.GF_ID          `bson:"id_str"`
	DeletedBool       bool                   `bson:"deleted_bool"`
	CreationUNIXtimeF float64                `bson:"creation_unix_time_f"`
	UserID            gf_core.GF_ID          `bson:"user_id"`
	
	// marked as true once the login completes (once Auth0 initial auth returns the user to the GF system).
	// if the login_callback handler is called and this login_complete is already marked as true,
	// the http transaction will be immediatelly aborted.
	LoginCompleteBool bool `bson:"login_complete_bool"`
	
	// user can specify which page then want to be redirect to after login,
	// if they dont want to use the default GF successful-login url.
	LoginSuccessRedirectURLstr string

	// user can specify which page then want to be redirect to after logout
	LogoutSuccessRedirectURLstr string

	AccessTokenStr string                 `bson:"access_token_str"`
	ProfileMap     map[string]interface{} `bson:"profile_map"`
}

type GFauth0inputLoginCallback struct {
	CodeStr           string
	SessionID         gf_core.GF_ID
	Auth0appDomainStr string
}
type GFauth0outputLoginCallback struct {
	JWTtokenStr string    
	LoginSuccessRedirectURLstr string 
}

type GFauth0inputAPItokenGenerate struct {
	UserID             gf_core.GF_ID
	AppClientIDstr     string
	AppClientSecretStr string
	AudienceStr        string
	Auth0appDomainStr  string
}

type GFauth0basicUserInfo struct {
	SubjectTypeStr string
	UserID         gf_core.GF_ID
	UserNameStr    GFuserName
	EmailStr       string
	ProfileMap     map[string]interface{}
}

//---------------------------------------------------
// TOKEN_GENERATE_PIPELINE

func Auth0apiTokenGeneratePipeline(pInput *GFauth0inputAPItokenGenerate,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	
	
	tokenStr, gfErr := gf_auth0.GenerateAccessToken(pInput.AppClientIDstr,
		pInput.AppClientSecretStr,
		pInput.AudienceStr,
		pInput.Auth0appDomainStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return tokenStr, nil
}

//---------------------------------------------------
// LOGOUT_PIPELINE

func Auth0logoutPipeline(pGFsessionID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	auth0session, gfErr := DBsqlAuth0getSession(pGFsessionID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	logoutSuccessRedirectURLstr := auth0session.LogoutSuccessRedirectURLstr

	// DELETE_SESSION
	gfErr = dbSQLauth0deleteSession(pGFsessionID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return logoutSuccessRedirectURLstr, nil
}

//---------------------------------------------------
// LOGIN

func Auth0loginPipeline(pLoginSuccessRedirectURLstr string,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	//---------------------
	// SESSION_ID
	sessionID := generateSessionID()
	
	//---------------------

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	auth0session := &GFauth0session{
		ID:                sessionID,
		CreationUNIXtimeF: creationUNIXtimeF,

		// indicate if the user already passed the initial login process,
		// and is now logged in.
		// this is a new Auth0 session, so the login is marked as not-complete.
		LoginCompleteBool: false,

		LoginSuccessRedirectURLstr: pLoginSuccessRedirectURLstr,
	}

	//---------------------
	// DB
	gfErr := dbSQLauth0createNewSession(auth0session,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gf_core.GF_ID(""), gfErr
	}

	/*
	with auth0 auth method only login_attept is created initially with session_id only.
	after auth0 logs the user in only then is the login_attempt updated with user info. 
	*/
	userTypeStr := "standard"
	_, gfErr = loginAttempCreateWithSession(sessionID, userTypeStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gf_core.GF_ID(""), gfErr
	}

	//------------------------

	return sessionID, nil
}

//---------------------------------------------------
// LOGIN_CALLBACK

func Auth0loginCallbackPipeline(pInput *GFauth0inputLoginCallback,
	pAuthenticator *gf_auth0.GFauthenticator,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (*GFauth0outputLoginCallback, *gf_core.GFerror) {
	
	sessionID := pInput.SessionID

	//---------------------
	// DB_GET_SESSION
	// verify that the sessionID (auth0 "state") corresponds to an registered session
	// created in the previously called login handler, and that a login with that session 
	// has not already been completed

	auth0session, gfErr := DBsqlAuth0getSession(sessionID,
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

	// auth0accessTokenStr := oauth2bearerToken.AccessToken
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

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("profile_map", profileMap)
	})

	subjectStr := profileMap["sub"].(string)
	subjectTypeStr := strings.Split(subjectStr, "|")[0]
	var userID      gf_core.GF_ID
	var userNameStr GFuserName
	var emailStr    string

	switch subjectTypeStr {

	//---------------------
	// DATABASE - username/password
	case "auth0":
		userID = gf_core.GF_ID(subjectStr)
		userNameStr = GFuserName(profileMap["nickname"].(string))
		emailStr = profileMap["email"].(string)

	//---------------------
	// EMAIL - passwordless login
	case "email":
		userID = gf_core.GF_ID(subjectStr)
		userNameStr = GFuserName(profileMap["nickname"].(string))
		emailStr = profileMap["email"].(string)

	//---------------------
	// GOOGLE
	case "google-oauth2":

		//---------------------------------------------------
		loadGoogleProfileFun := func() {

			googleUserIDstr   := profileMap["sub"].(string)
			googleNicknameStr := profileMap["nickname"].(string)
			userID      = gf_core.GF_ID(googleUserIDstr)
			userNameStr = GFuserName(googleNicknameStr)
			emailStr    = profileMap["email"].(string)

			googleProfile := &GFgoogleUserProfile {
				NameStr:       profileMap["name"].(string),
				GivenNameStr:  profileMap["given_name"].(string),
				FamilyNameStr: profileMap["family_name"].(string),
				NicknameStr:   googleNicknameStr,
				// LocaleStr:     profileMap["locale"].(string),
				UpdatedAtStr:  profileMap["updated_at"].(string),
				PictureURLstr: profileMap["picture"].(string),
			}



			// LOCAL
			// check if locale is in profile map, or if its nil, and if its an actual string, assign it to the property of googleProfile
			if localeVal, exists := profileMap["locale"]; exists && localeVal != nil {
				if locale, ok := localeVal.(string); ok {
					googleProfile.LocaleStr = locale
				}
			}


			pRuntimeSys.LogNewFun("DEBUG", "google user profile loaded...", nil)
			if gf_core.LogsIsDebugEnabled() {
				spew.Dump(googleProfile)
			}
		}

		//---------------------------------------------------

		/*
		ERROR_HANDLING
		using WithScope() to attach the profileMap to the Sentry error report, in case loading the google profile fails.
		this way the profileMap is attached to the error report and can be used to debug the issue, and once this scope is
		closed the profileMap is removed from the scope.
		extracting values out of google profile is a frequent source of errors, since some fields might be missing.
		*/
		sentry.WithScope(func(pScope *sentry.Scope) {
		
			gf_core.AttachMapToSentryScope(profileMap, pScope)

			loadGoogleProfileFun()
		})

	//---------------------
	}
	
	basicUserInfo := &GFauth0basicUserInfo{
		SubjectTypeStr: subjectTypeStr,
		UserID:         userID,
		UserNameStr:    userNameStr,
		EmailStr:       emailStr,
		ProfileMap:     profileMap,
	}

	pRuntimeSys.LogNewFun("DEBUG", "basic user info", map[string]interface{}{
		"subject_type": basicUserInfo.SubjectTypeStr,
		"user_id":      basicUserInfo.UserID,
		"user_name":    basicUserInfo.UserNameStr,
		"email":        basicUserInfo.EmailStr,
		"profile_map":  basicUserInfo.ProfileMap,
	})

	gfErr = Auth0createGFuserIfNone(basicUserInfo,
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

	gfErr = dbSQLauth0updateSession(sessionID,
		userID,
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
	/*
	LOGIN_ATTEMPT
	on user login success, when all the info on the user is present, update login_attempt in the DB.
	unlike with GF native userpass/eth auth methods where the login attempt can be created from username
	right away before creds are checked (even if login ultimately fails), with Auth0 auth method this has to be done
	at the end since the username is not known right away as the user is navigated to Auth0 systems and only
	returned to GF on login success.
	*/ 
	updateOp := &GFloginAttemptUpdateOp{
		UserID:      &userID,
		UserNameStr: &userNameStr,
	}
	gfErr = DBsqlLoginAttemptUpdateBySessionID(sessionID, updateOp, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//---------------------

	output := &GFauth0outputLoginCallback{
		JWTtokenStr: JWTtokenStr,
		LoginSuccessRedirectURLstr: auth0session.LoginSuccessRedirectURLstr,
	}
	return output, nil
}

//---------------------------------------------------
// CREATE_GF_USER_IF_NONE

// check if the Auth0 user exists in the DB, and if not create it.
// a user would not exist in the DB if it signed-up/logged-in for the first time.
func Auth0createGFuserIfNone(pBasicUserInfo *GFauth0basicUserInfo,
	pAuth0appDomainStr string,
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// USER_INFO - from Auth0
	/*
	map[
		family_name: Trajkovic
		given_name:  Ivan
		locale:      en
		name:        Ivan Trajkovic
		nickname:    ivan.ebiz
		picture:     https://...user_image...
		sub:         ...user_id...
		updated_at:  2023-08-16T22:13:10.084Z
	]
	*/
	/*
	auth0userInfoMap, gfErr := gf_auth0.GetUserInfo(pAuth0accessTokenStr,
		pAuth0appDomainStr,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}
	*/

	pRuntimeSys.LogNewFun("DEBUG", `>>>>>>>>>>>>>>>>> Auth0 /userinfo response recieved...`,
		map[string]interface{}{
			"auth0_user_info_map": pBasicUserInfo.ProfileMap,
		})

	//---------------------

	// the user_info returned by Auth0 contains the "sub" claim
	// which is the user_id assigned to the user in the Auth0 system
	auth0userID := pBasicUserInfo.UserID // gf_core.GF_ID(pUserProfileMap["sub"].(string))

	//---------------------
	// DB
	existsBool, gfErr := DBsqlUserExistsByID(auth0userID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//---------------------

	profileImageURLstr := pBasicUserInfo.ProfileMap["picture"].(string)
	screenNameStr := string(pBasicUserInfo.UserNameStr)
	// screenNameStr := pBasicUserInfo.ProfileMap["nickname"].(string)
	// userNameStr := GFuserName(pUserProfileMap["name"].(string))

	// user doesnt exist in the GF DB
	if !existsBool {

		creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

		//---------------------
		// EMAIL
		// only read the email param if a new user is created in the GF DB
		emailStr := pBasicUserInfo.EmailStr // pUserProfileMap["email"]

		// check if the email field is present in the user_info map. sometimes
		// its not returned, and the user would be created without an email.
		if emailStr == "" {
			gfErr := gf_core.ErrorCreate("Auth0 user_info does not contain the expected/required 'email' field",
				"library_error",
				map[string]interface{}{},
				nil, "gf_identity_core", pRuntimeSys)
			return gfErr
		}

		//---------------------


		user := &GFuser{
			Vstr:               "0",
			ID:                 auth0userID,
			CreationUNIXtimeF:  creationUNIXtimeF,
			UserTypeStr:        "standard",
			UserNameStr:        pBasicUserInfo.UserNameStr, // userNameStr,
			ScreenNameStr:      screenNameStr,              // screenNameStr,
			EmailStr:           emailStr,
			ProfileImageURLstr: profileImageURLstr,
		}
	
		//------------------------
		// DB
		gfErr = DBsqlUserCreate(user, pCtx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	
		//------------------------
	
	} else {

		update := &GFuserUpdateOp {
			UserNameStr:        &pBasicUserInfo.UserNameStr, // &userNameStr,
			ScreenNameStr:      &screenNameStr,              // &screenNameStr,
			ProfileImageURLstr: &profileImageURLstr,
		}

		gfErr = DBsqlUserUpdate(auth0userID,
			update,
			pCtx,
			pRuntimeSys)
	}

	return nil
}