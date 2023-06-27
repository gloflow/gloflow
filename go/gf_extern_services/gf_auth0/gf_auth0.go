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

package gf_auth0

import (
	"os"
	"fmt"
	"context"
	"crypto/rsa"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	jwks "github.com/MicahParks/keyfunc"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------------------

type GFconfig struct {
	Auth0domainStr       string
	Auth0clientIDstr     string
	Auth0clientSecretStr string
	Auth0apiAudienceStr  string
	Auth0callbackURLstr  string
}

type GFonLoginSuccessProfileInfo struct {
	SubStr      string
	NicknameStr string
	NameStr     string
	PictureURLstr string
}

type GFauthenticator struct {
	*oidc.Provider
	oauth2.Config
}

//-------------------------------------------------------------

func GetJWTpublicKeyForTenant(pConfig *GFconfig,
	pRuntimeSys *gf_core.RuntimeSys) (string, *rsa.PublicKey, *gf_core.GFerror) {

	jwksURLstr := fmt.Sprintf("https://%s/.well-known/jwks.json", pConfig.Auth0domainStr)
	
	jwks, err := jwks.Get(jwksURLstr, jwks.Options{}) // See recommended options in the examples directory.
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get the JWKS from the given Auth0 URL",
			"library_error",
			map[string]interface{}{
				"jwks_url_str": jwksURLstr,
			},
			err, "gf_auth0", pRuntimeSys)
		return "", nil, gfErr
	}

	var auth0keyIDstr  string
	var auth0publicKey *rsa.PublicKey

	// auth0 always returns a list of keys in the JWKS, where the first one is the active key
	// used for signing, and all subsequent keys are pending (next in queue) keys.
	for k, v := range jwks.ReadOnlyKeys() {
		
		auth0keyIDstr  = k
		auth0publicKey = v.(*rsa.PublicKey)

		break
	}

	return auth0keyIDstr, auth0publicKey, nil
}

//-------------------------------------------------------------

func Init(pRuntimeSys *gf_core.RuntimeSys) (*GFauthenticator, *GFconfig, *gf_core.GFerror) {

	pRuntimeSys.LogNewFun("INFO", "initializing Auth0...", nil)

	//-------------------
	// CONFIG
	config := LoadConfig(pRuntimeSys)
	pRuntimeSys.LogNewFun("INFO", "Auth0 config loaded...",
		map[string]interface{}{
			"auth0_app_domain": config.Auth0domainStr,
		})

	//-------------------
	
	provider, err := oidc.NewProvider(
		context.Background(),
		fmt.Sprintf("https://%s/", config.Auth0domainStr),
	)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to initialize a auth0 openID_connect provider",
			"library_error",
			map[string]interface{}{},
			err, "gf_auth0", pRuntimeSys)
		return nil, nil, gfErr
	}

	conf := oauth2.Config{
		ClientID:     config.Auth0clientIDstr,
		ClientSecret: config.Auth0clientSecretStr,
		RedirectURL:  config.Auth0callbackURLstr,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	authenticator := &GFauthenticator{
		Provider: provider,
		Config:   conf,
	}


	return authenticator, config, nil
}

//-------------------------------------------------------------

// verifies that an Oauth2 token is a valid *oidc.IDToken.
func VerifyIDtoken(pOauth2bearerToken *oauth2.Token,
	pAuthenticator *GFauthenticator,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (*oidc.IDToken, *gf_core.GFerror) {

	pRuntimeSys.LogNewFun("DEBUG", "verifying OpenID id_token...", nil)

	// not making a network request, gets id_token from pOauth2bearerToken.raw.id_token
	IDtokenEncodedStr, ok := pOauth2bearerToken.Extra("id_token").(string)
	if !ok {
		gfErr := gf_core.ErrorCreate("failed to get an id_token from oauth2 Token in auth0 flow",
			"library_error",
			map[string]interface{}{},
			nil, "gf_auth0", pRuntimeSys)
		return nil, gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "encoded id token", map[string]interface{}{
		"id_token_encoded_str": IDtokenEncodedStr,
	})

	oidcConfig := &oidc.Config{
		ClientID: pAuthenticator.ClientID,
	}

	// https://auth0.com/docs/authenticate/protocols/openid-connect-protocol
	// https://openid.net/connect/
	idToken, err := pAuthenticator.Verifier(oidcConfig).Verify(pCtx, IDtokenEncodedStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to verify raw ID token in auth0 flow",
			"library_error",
			map[string]interface{}{},
			err, "gf_auth0", pRuntimeSys)
		return nil, gfErr
	}

	return idToken, nil
}

//-------------------------------------------------------------

// load Auth0 config, mostly from ENV vars
func LoadConfig(pRuntimeSys *gf_core.RuntimeSys) *GFconfig {

	auth0domainStr       := os.Getenv("AUTH0_DOMAIN")
	auth0clientIDstr     := os.Getenv("AUTH0_CLIENT_ID")
	auth0clientSecretStr := os.Getenv("AUTH0_CLIENT_SECRET")
	auth0apiAudienceStr  := os.Getenv("AUTH0_AUDIENCE")
	auth0callbackURLstr  := os.Getenv("AUTH0_CALLBACK_URL")

	pRuntimeSys.LogNewFun("INFO", "auth0 config loaded", map[string]interface{}{
		"auth0_domain_str":             auth0domainStr,
		"auth0_login_callback_url_str": auth0callbackURLstr,
	})

	config := &GFconfig{
		Auth0domainStr:       auth0domainStr,
		Auth0clientIDstr:     auth0clientIDstr,
		Auth0clientSecretStr: auth0clientSecretStr,
		Auth0apiAudienceStr:  auth0apiAudienceStr,
		Auth0callbackURLstr:  auth0callbackURLstr,
	}
	return config
}

//-------------------------------------------------------------