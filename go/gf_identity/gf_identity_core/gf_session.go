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

func SessionValidate(pReq *http.Request,
	pKeyServerInfo         *GFkeyServerInfo,
	pAuthSubsystemTypeStr  string,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (bool, string, gf_core.GF_ID, *gf_core.GFerror) {

	//---------------------
	// JWT
	jwtTokenStr, foundBool, gfErr := JWTgetTokenFromRequest(pReq, pRuntimeSys)
	if gfErr != nil {
		return false, "", gf_core.GF_ID(""), gfErr
	}
	
	if !foundBool {
		
		/*
		IMPORTANT!! - return a false validity and not an error, since missing
		    JWT in request is not an abnormal situation (an error), and 
		    it means that the user is not authenticated yet.
		*/
		return false, "", gf_core.GF_ID(""), nil
	}

	//---------------------
	// SESSION_ID

	sessionID, sessionIDfoundBool := GetSessionID(pReq, pRuntimeSys)
	if !sessionIDfoundBool {

		/*
		IMPORTANT!! - return a false validity and not an error, since missing
		    session_id in request is not an abnormal situation (an error), and 
		    it means that the user is not authenticated yet.
		*/
		return false, "", gf_core.GF_ID(""), nil
	}

	//---------------------
	
	var userIdentifierStr string

	switch pAuthSubsystemTypeStr {
	
	//---------------------
	// AUTH0
	case GF_AUTH_SUBSYSTEM_TYPE__AUTH0:

		// KEY_SERVER
		publicKey, gfErr := KSclientJWTgetValidationKey(GF_AUTH_SUBSYSTEM_TYPE__AUTH0,
			pKeyServerInfo,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gf_core.GF_ID(""), gfErr
		}
		
		userIdentifierFromJWTstr, gfErr := gf_auth0.JWTvalidateToken(jwtTokenStr, publicKey, pRuntimeSys)
		if gfErr != nil {
			return false, "", gf_core.GF_ID(""), gfErr
		}

		userIdentifierStr = userIdentifierFromJWTstr

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

		userIdentifierStr = userIdentifierFromJWTstr

	//---------------------
	}

	return true, userIdentifierStr, sessionID, nil
}

//---------------------------------------------------
// COOKIES
//---------------------------------------------------

func CreateAuthCookie(pJWTtokenStr string,
	pDomainStr *string,
	pResp      http.ResponseWriter) {

	sessionTTLhoursInt, _ := GetSessionTTL()

	cookieNameStr := "Authorization"
	cookieDataStr := fmt.Sprintf("Bearer %s", pJWTtokenStr)
	gf_core.HTTPsetCookieOnResp(cookieNameStr,
		cookieDataStr,
		pDomainStr,
		pResp,
		sessionTTLhoursInt)
}

func CreateAuthCookieOnReq(pJWTtokenStr string,
	pDomainStr *string,
	pReq       *http.Request) {

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
	pDomainStr *string,
	pResp      http.ResponseWriter) {
	
	sessionTTLhoursInt, _ := GetSessionTTL()

	cookieNameStr := "gf_sess"
	cookieDataStr := pSessionIDstr
	gf_core.HTTPsetCookieOnResp(cookieNameStr,
		cookieDataStr,
		pDomainStr,
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

func GetSessionID(pReq *http.Request,
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