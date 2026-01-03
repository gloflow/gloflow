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
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
)
 
//---------------------------------------------------

type GFsession struct {
	ID                gf_core.GF_ID
	DeletedBool       bool
	CreationUNIXtimeF float64
	UserID            gf_core.GF_ID
	
	// marked as true once the login completes (once Auth0 initial auth returns the user to the GF system).
	// if the login_callback handler is called and this login_complete is already marked as true,
	// the http transaction will be immediatelly aborted.
	LoginCompleteBool bool
	
	// user can specify which page then want to be redirect to after login,
	// if they dont want to use the default GF successful-login url.
	LoginSuccessRedirectURLstr string

	// user can specify which page then want to be redirect to after logout
	LogoutSuccessRedirectURLstr string

	AccessTokenStr string
	ProfileMap     map[string]interface{}

	AuthSubsystemTypeStr string
	AuthMethodStr        *string

	// USER_ID_IDP - user ID in the IDP system (google, github, etc.)
	UserIDidp *gf_core.GF_ID
}

//---------------------------------------------------

func SessionValidateOrRedirectToLogin(pReq *http.Request,
	pResp                   http.ResponseWriter,
	pKeyServerInfo          *GFkeyServerInfo,
	pAuthSubsystemTypeStr   string,
	pAuthLoginURLstr        *string,
	pAuthRedirectOnFailBool bool,
	pCtx                    context.Context,
	pRuntimeSys             *gf_core.RuntimeSys) (bool, gf_core.GF_ID, gf_core.GF_ID, *gf_core.GFerror) {

	var validBool bool
	var userID    gf_core.GF_ID
	var sessionID gf_core.GF_ID
	var gfErr     *gf_core.GFerror

	//---------------------
	// API_KEY_VALIDATE
	// PLUGIN - use custom session validation for api keys if provided

	// check if an API key is supplied; if not supplied Get() will return ""
	if apiKeyStr := pReq.Header.Get("X-Api-Key"); apiKeyStr != "" {

		if pRuntimeSys.ExternalPlugins != nil && pRuntimeSys.ExternalPlugins.IdentitySessionValidateApiKeyCallback != nil {

			validBool, userID, gfErr = pRuntimeSys.ExternalPlugins.IdentitySessionValidateApiKeyCallback(apiKeyStr,
				pReq,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return false, "", gf_core.GF_ID(""), gfErr
			}

			// calls with an API key dont have session IDs
			sessionID = gf_core.GF_ID("")
			
		} else {
			return false, gf_core.GF_ID(""), gf_core.GF_ID(""), nil
		}
	}
	
	//---------------------
	// SESSION_VALIDATE

	// PLUGIN - use custom session validation if provided
	if pRuntimeSys.ExternalPlugins != nil && pRuntimeSys.ExternalPlugins.IdentitySessionValidateCallback != nil {

		validBool, userID, sessionID, gfErr = pRuntimeSys.ExternalPlugins.IdentitySessionValidateCallback(pReq,
			pResp,
			pCtx,
			pRuntimeSys)

		if gfErr != nil {
			return false, gf_core.GF_ID(""), gf_core.GF_ID(""), gfErr
		}

	// INTERNAL - use built-in session validation
	} else {
		
		validBool, userID, sessionID, gfErr = SessionValidate(pReq,
			pKeyServerInfo,
			pAuthSubsystemTypeStr,
			pCtx,
			pRuntimeSys)
	}

	//---------------------------------------------------
	redirectFun := func() {			
		//-------------------------
		// HTTP_REDIRECT - redirect user to login url
		http.Redirect(pResp,
			pReq,
			*pAuthLoginURLstr,
			301)
		
		//-------------------------
	}

	//---------------------------------------------------

	if gfErr != nil {

		// if the JWT supplied by the user to auth is invalid,
		// redirect the user to the login page so that they can auth.
		if gfErr.Type_str == "crypto_jwt_verify_token_error" {
			if pAuthRedirectOnFailBool && pAuthLoginURLstr != nil {
				redirectFun()
			}
		}

		return false, gf_core.GF_ID(""), gf_core.GF_ID(""), gfErr
	}

	if !validBool {
		if pAuthRedirectOnFailBool && pAuthLoginURLstr != nil {
			redirectFun()
		}

		return false, gf_core.GF_ID(""), gf_core.GF_ID(""), nil
	}

	return validBool, userID, sessionID, nil
}

//---------------------------------------------------

func SessionValidate(pReq *http.Request,
	pKeyServerInfo         *GFkeyServerInfo,
	pAuthSubsystemTypeStr  string,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (bool, gf_core.GF_ID, gf_core.GF_ID, *gf_core.GFerror) {

	//---------------------
	// JWT
	jwtTokenStr, foundBool, gfErr := JWTgetTokenFromRequest(pReq, pRuntimeSys)
	if gfErr != nil {
		return false, gf_core.GF_ID(""), gf_core.GF_ID(""), gfErr
	}
	
	if !foundBool {
		
		/*
		IMPORTANT!! - return a false validity and not an error, since missing
		    JWT in request is not an abnormal situation (an error), and 
		    it means that the user is not authenticated yet.
		*/
		return false, gf_core.GF_ID(""), gf_core.GF_ID(""), nil
	}

	//---------------------
	// SESSION_ID

	sessionID, sessionIDfoundBool := GetSessionIDfromReq(pReq, pRuntimeSys)
	if !sessionIDfoundBool {

		/*
		IMPORTANT!! - return a false validity and not an error, since missing
		    session_id in request is not an abnormal situation (an error), and 
		    it means that the user is not authenticated yet.
		*/
		return false, gf_core.GF_ID(""), gf_core.GF_ID(""), nil
	}

	//---------------------
	
	var userID gf_core.GF_ID

	switch pAuthSubsystemTypeStr {
	
	//---------------------
	// AUTH0
	case GF_AUTH_SUBSYSTEM_TYPE__AUTH0:

		// KEY_SERVER
		publicKey, gfErr := KSclientJWTgetValidationKey(GF_AUTH_SUBSYSTEM_TYPE__AUTH0,
			pKeyServerInfo,
			pRuntimeSys)
		if gfErr != nil {
			return false, gf_core.GF_ID(""), gf_core.GF_ID(""), gfErr
		}
		
		userIdentifierFromJWTstr, gfErr := gf_auth0.JWTvalidateToken(jwtTokenStr, publicKey, pRuntimeSys)
		if gfErr != nil {
			return false, gf_core.GF_ID(""), gf_core.GF_ID(""), gfErr
		}

		userID = gf_core.GF_ID(userIdentifierFromJWTstr)

	//---------------------
	// USERPASS
	case GF_AUTH_SUBSYSTEM_TYPE__USERPASS:

		// JWT_VALIDATE
		userIdentifierFromJWTstr, gfErr := JWTpipelineValidate(GFjwtTokenVal(jwtTokenStr),
			pAuthSubsystemTypeStr,
			pKeyServerInfo,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gf_core.GF_ID(""), gfErr
		}

		userID = gf_core.GF_ID(userIdentifierFromJWTstr)

	//---------------------
	}

	return true, userID, sessionID, nil
}

//---------------------------------------------------

func SessionCreate(pUserID *gf_core.GF_ID,
	pLoginSuccessRedirectURLstr *string,
	pAuthSubsystemTypeStr       string,
	pAuthMethodStr			    *string,
	pUserIDidp                  *gf_core.GF_ID,
	pCtx                        context.Context,
	pRuntimeSys                 *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	//---------------------
	// SESSION_ID
	sessionID := generateSessionID()
	
	//---------------------

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	session := &GFsession{
		ID:                sessionID,
		CreationUNIXtimeF: creationUNIXtimeF,

		// indicate if the user already passed the initial login process,
		// and is now logged in.
		// this is a new Auth0 session, so the login is marked as not-complete.
		LoginCompleteBool: false,

		AuthSubsystemTypeStr: pAuthSubsystemTypeStr,
		AuthMethodStr:        pAuthMethodStr,
		UserIDidp:            pUserIDidp,
	}

	// with Auth0 the userID is only known after the login completes, and the session
	// record is updated then. session is created initially without userID at the login start.
	if pUserID != nil {
		session.UserID = *pUserID
	}

	if pLoginSuccessRedirectURLstr != nil {
		session.LoginSuccessRedirectURLstr = *pLoginSuccessRedirectURLstr
	}

	//---------------------
	// DB
	gfErr := dbSQLcreateNewSession(session,
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
// COOKIES
//---------------------------------------------------

func CreateAuthCookie(pJWTtokenStr string,
	pDomainStr          *string,
	pSameSiteStrictBool bool,
	pResp               http.ResponseWriter) {

	sessionTTLhoursInt, _ := GetSessionTTL()

	cookieNameStr := "Authorization"
	cookieDataStr := fmt.Sprintf("Bearer %s", pJWTtokenStr)
	gf_core.HTTPsetCookieOnResp(cookieNameStr,
		cookieDataStr,
		pDomainStr,
		pSameSiteStrictBool,
		pResp,
		sessionTTLhoursInt)
}

func CreateAuthCookieOnReq(pJWTtokenStr string,
	pDomainStr          *string,
	pReq                *http.Request) {

	sessionTTLhoursInt, _ := GetSessionTTL()

	cookieNameStr := "Authorization"
	cookieDataStr := fmt.Sprintf("Bearer %s", pJWTtokenStr)

	gf_core.HTTPsetCookieOnReq(cookieNameStr,
		cookieDataStr,
		pDomainStr,
		pReq,
		sessionTTLhoursInt)
}

//---------------------------------------------------

func CreateSessionIDcookie(pSessionIDstr string,
	pDomainStr          *string,
	pSameSiteStrictBool bool,
	pResp               http.ResponseWriter) {
	
	sessionTTLhoursInt, _ := GetSessionTTL()

	cookieNameStr := "gf_sess"
	cookieDataStr := pSessionIDstr
	gf_core.HTTPsetCookieOnResp(cookieNameStr,
		cookieDataStr,
		pDomainStr,
		pSameSiteStrictBool,
		pResp,
		sessionTTLhoursInt)
}

//---------------------------------------------------

func DeleteCookies(pDomainForAuthCookiesStr string,
	pResp http.ResponseWriter) {

	sessCookieNameStr := "gf_sess"
	gf_core.HTTPdeleteCookieOnResp(sessCookieNameStr, pDomainForAuthCookiesStr, pResp)

	jwtCookieNameStr := "Authorization"
	gf_core.HTTPdeleteCookieOnResp(jwtCookieNameStr, pDomainForAuthCookiesStr, pResp)
}

//---------------------------------------------------

func GetSessionIDfromReq(pReq *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, bool) {

	sessCookieNameStr := "gf_sess"
	existsBool, valStr := gf_core.HTTPgetCookieFromReq(sessCookieNameStr, pReq, pRuntimeSys)

	sessionID := gf_core.GF_ID(valStr)
	return sessionID, existsBool
}

//---------------------------------------------------
// VAR
//---------------------------------------------------

func generateSessionID() gf_core.GF_ID {

	creationUNIXtimeF  := float64(time.Now().UnixNano())/1000000000.0
	randomStr          := gf_core.StrRandom()
	uniqueValsForIDlst := []string{
		randomStr,
	}
	sessionIDstr := gf_core.IDcreate(uniqueValsForIDlst, creationUNIXtimeF)
	return sessionIDstr
}

//---------------------------------------------------

func GetSessionTTL() (int, int64) {

	//---------------------
	// FIX!! - this should be configurable
	sessionTTLhoursInt := 24 * 30 // 1 month

	//---------------------
	
	sessionTTLsecondsInt := int64(60*60*sessionTTLhoursInt)
	return sessionTTLhoursInt, sessionTTLsecondsInt
}