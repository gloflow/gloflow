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
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
)

//---------------------------------------------------

func SessionValidate(pReq *http.Request,
	pKeyServerInfo         *GFkeyServerInfo,
	pAuthSubsystemTypeStr  string,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {
	
	

	var userIdentifierStr string
	switch pAuthSubsystemTypeStr {
	
	//---------------------
	// AUTH0
	case GF_AUTH_SUBSYSTEM_TYPE__AUTH0:

		//---------------------
		// JWT
		auth0JWTtokenStr, foundBool, gfErr := JWTgetTokenFromRequest(pReq, pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}
		
		if !foundBool {
			
			// IMPORTANT!! - return a false validity and not an error, since missing
			//               JWT in request is not an ubnormal situation (an error), and 
			//               it means that the user is not authenticated yet.
			return false, "", nil
		}

		//---------------------

		// KEY_SERVER
		publicKey, gfErr := KSclientJWTgetValidationKey(GF_AUTH_SUBSYSTEM_TYPE__AUTH0,
			pKeyServerInfo,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}
		
		userIdentifierFromJWTstr, gfErr := gf_auth0.JWTvalidateToken(auth0JWTtokenStr, publicKey, pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		userIdentifierStr = userIdentifierFromJWTstr

	//---------------------
	// USERPASS
	case GF_AUTH_SUBSYSTEM_TYPE__USERPASS:
		
		//---------------------
		// JWT
		JWTtokenStr, foundBool, gfErr := JWTgetTokenFromRequest(pReq, pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}
		
		if !foundBool {
			
			// IMPORTANT!! - return a false validity and not an error, since missing
			//               JWT in request is not an ubnormal situation (an error), and 
			//               it means that the user is not authenticated yet.
			return false, "", nil
		}

		//---------------------

		// JWT_VALIDATE
		userIdentifierFromJWTstr, gfErr := JWTpipelineValidate(GFjwtTokenVal(JWTtokenStr),
			pAuthSubsystemTypeStr,
			pKeyServerInfo,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		userIdentifierStr = userIdentifierFromJWTstr

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

	//---------------------
	// FIX!! - this should be configurable
	sessionTTLhoursInt := 24 * 30 // 1 month

	//---------------------
	
	sessionTTLsecondsInt := int64(60*60*sessionTTLhoursInt)
	return sessionTTLhoursInt, sessionTTLsecondsInt
}