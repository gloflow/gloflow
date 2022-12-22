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
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//---------------------------------------------------

func CreateCookie(pSessionDataStr string,
	pResp http.ResponseWriter) {
	
	sessionTTLhoursInt, _ := gf_identity_core.GetSessionTTL()

	cookieNameStr := "gf_sess"
	cookieDataStr := pSessionDataStr
	gf_core.HTTPsetCookieOnReq(cookieNameStr,
		cookieDataStr,
		pResp,
		sessionTTLhoursInt)
}

//---------------------------------------------------

func ValidateOrRedirectToLogin(pReq *http.Request,
	pResp                  http.ResponseWriter,
	pKeyServerInfo         *gf_identity_core.GFkeyServerInfo,
	pAuthSubsystemTypeStr  string,
	pAuthLoginURLstr       *string,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {

	validBool, userIdentifierStr, gfErr := gf_identity_core.Validate(pReq,
		pKeyServerInfo,
		pAuthSubsystemTypeStr,
		pCtx,
		pRuntimeSys)

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