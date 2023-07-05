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

func SessionValidate(pReq *http.Request,
	pKeyServerInfo         *GFkeyServerInfo,
	pAuthSubsystemTypeStr  string,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {
	
	cookieNameStr := "gf_sess"
	cookieFoundBool, sessionDataStr := gf_core.HTTPgetCookieFromReq(cookieNameStr, pReq)
	
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