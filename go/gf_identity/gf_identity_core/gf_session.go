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
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func Validate(pReq *http.Request,
	pKeyServerInfo         *GFkeyServerInfo,
	pAuthSubsystemTypeStr  string,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {
	
	cookieNameStr := "gf_sess"
	cookieFoundBool, sessionDataStr := GetFromReq(cookieNameStr, pReq)
	
	if !cookieFoundBool {

		// gf_sess cookie was never found
		return false, "", nil
	}

	var userIdentifierStr string
	switch pAuthSubsystemTypeStr {
	
	//---------------------
	// USERPASS
	case GF_AUTH_SUBSYSTEM_TYPE__USERPASS:
		
		JWTtokenValStr := sessionDataStr

		// JWT_VALIDATE
		JWTuserIdentifierStr, gfErr := JWTpipelineValidate(GFjwtTokenVal(JWTtokenValStr),
			pAuthSubsystemTypeStr,
			pKeyServerInfo,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		userIdentifierStr = JWTuserIdentifierStr

	//---------------------
	// AUTH0
	case GF_AUTH_SUBSYSTEM_TYPE__AUTH0:
		sessionIDstr := gf_core.GF_ID(sessionDataStr)

		validBool, gfErr := Auth0validateSession(sessionIDstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		if !validBool {
			return false, "", nil
		}
		
	//---------------------
	}

	return true, userIdentifierStr, nil
}

//---------------------------------------------------
// GET/SET FROM COOKIES
//---------------------------------------------------

func GetFromReq(pCookieNameStr string,
	pReq *http.Request) (bool, string) {

	for _, cookie := range pReq.Cookies() {
		if (cookie.Name == pCookieNameStr) {
			sessionDataStr := cookie.Value
			return true, sessionDataStr
		}
	}
	return false, ""
}

//---------------------------------------------------

func SetOnReq(pSessionCookieNameStr string,
	pSessionDataStr string,
	pResp           http.ResponseWriter,
	pTTLhoursInt    int) {

	ttl    := time.Duration(pTTLhoursInt) * time.Hour
	expire := time.Now().Add(ttl)
	
	cookie := http.Cookie{
		Name:    pSessionCookieNameStr,
		Value:   pSessionDataStr,
		Expires: expire,

		// IMPORTANT!! - session cookie should be set for all paths
		//               on the same domain, not just the /v1/identity/...
		//               paths, because session is verified on all of them
		Path: "/", 
		
		// ADD!! - ability to specify multiple domains that the session is
		//         set for in case the GF services and API endpoints are spread
		//         across multiple domains.
		// Domain: "", 
		
		// IMPORTANT!! - make cookie http_only, disabling browser js context
		//               from being able to read its value
		HttpOnly: true,

		// SameSite allows a server to define a cookie attribute making it impossible for
		// the browser to send this cookie along with cross-site requests. The main
		// goal is to mitigate the risk of cross-origin information leakage, and provide
		// some protection against cross-site request forgery attacks.
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(pResp, &cookie)
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
	sessionTTLhoursInt   := 24 * 30 // 1 month
	sessionTTLsecondsInt := int64(60*60*24*7)
	return sessionTTLhoursInt, sessionTTLsecondsInt
}