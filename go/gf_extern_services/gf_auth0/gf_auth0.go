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
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"github.com/gloflow/gloflow/go/gf_core"
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

func Init(pRuntimeSys *gf_core.RuntimeSys) (*GFauthenticator, *GFconfig, *gf_core.GFerror) {

	fmt.Println("INITIALIZING AUTH0")

	config := loadConfig()
	
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
func VerifyIDtoken(pOauthToken *oauth2.Token,
	pAuthenticator *GFauthenticator,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (*oidc.IDToken, *gf_core.GFerror) {

	rawIDtoken, ok := pOauthToken.Extra("id_token").(string)
	if !ok {
		gfErr := gf_core.ErrorCreate("failed to get an id_token from oauth2 Token in auth0 flow",
			"library_error",
			map[string]interface{}{},
			nil, "gf_auth0", pRuntimeSys)
		return nil, gfErr
	}

	oidcConfig := &oidc.Config{
		ClientID: pAuthenticator.ClientID,
	}

	// https://auth0.com/docs/authenticate/protocols/openid-connect-protocol
	// https://openid.net/connect/
	idToken, err := pAuthenticator.Verifier(oidcConfig).Verify(pCtx, rawIDtoken)
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
func loadConfig() *GFconfig {

	auth0domainStr       := os.Getenv("AUTH0_DOMAIN")
	auth0clientIDstr     := os.Getenv("AUTH0_CLIENT_ID")
	auth0clientSecretStr := os.Getenv("AUTH0_CLIENT_SECRET")
	auth0apiAudienceStr  := os.Getenv("AUTH0_AUDIENCE")
	auth0callbackURLstr  := os.Getenv("AUTH0_CALLBACK_URL")

	fmt.Println("auth0 domain       - ", auth0domainStr)
	fmt.Println("auth0 callback URL - ", auth0callbackURLstr)

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