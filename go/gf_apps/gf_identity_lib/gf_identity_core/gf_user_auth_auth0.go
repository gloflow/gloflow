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
	"context"
	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func Auth0middlewareInit(pRuntimeSys *gf_core.RuntimeSys) *jwtmiddleware.JWTMiddleware {

	//-------------------------------------------------
	keyGenerateFun := func(pCtx context.Context) (interface{}, error) {
		
		userIdentifierStr := ""
		jwtSecretKeyValStr, gfErr := JWTgenerateSecretSigningKey(userIdentifierStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr.Error
		}

		return []byte(string(jwtSecretKeyValStr)), nil
	}

	//-------------------------------------------------

	apiAudienceStr := ""
	jwtValidator, err := validator.New(
		keyGenerateFun,
		validator.HS256,
		"https://<issuer-url>/",
		[]string{apiAudienceStr,},
	)
	if err != nil {
		panic(err)
	}


	jwtAuth0middleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	return jwtAuth0middleware
}