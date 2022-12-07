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

package gf_session

import (
	"time"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//---------------------------------------------------

func Create(pJWTtokenValStr string,
	pResp http.ResponseWriter) {
	
	sessionDataStr        := pJWTtokenValStr
	sessionTTLhoursInt, _ := gf_identity_core.GetSessionTTL()

	sessionCookieNameStr := "gf_sess"
	SetOnReq(sessionCookieNameStr,
		sessionDataStr,
		pResp,
		sessionTTLhoursInt)
}

//---------------------------------------------------

func ValidateOrRedirectToLogin(pReq *http.Request,
	pResp            http.ResponseWriter,
	pAuthLoginURLstr *string,
	pCtx             context.Context,
	pRuntimeSys      *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {

	validBool, userIdentifierStr, gfErr := Validate(pReq, pCtx, pRuntimeSys)
	if gfErr != nil {
		return false, "", gfErr
	}

	if !validBool {
		if pAuthLoginURLstr != nil {

			// redirect user to login url
			http.Redirect(pResp,
				pReq,
				*pAuthLoginURLstr,
				301)

			return false, "", nil
		} else {
			return false, "", nil
		}
	}

	return validBool, userIdentifierStr, nil
}

//---------------------------------------------------

func Validate(pReq *http.Request,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {
	
	cookieNameStr := "gf_sess"
	cookieFoundBool, cookieValStr := GetFromReq(cookieNameStr, pReq)
	
	if !cookieFoundBool {

		// gf_sess cookie was never found
		return false, "", nil
	}
	
	sessionDataStr := cookieValStr
	JWTtokenValStr := sessionDataStr

	//---------------------
	// JWT_VALIDATE
	userIdentifierStr, gfErr := gf_identity_core.JWTpipelineValidate(gf_identity_core.GFjwtTokenVal(JWTtokenValStr),
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return false, "", gfErr
	}

	//---------------------

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