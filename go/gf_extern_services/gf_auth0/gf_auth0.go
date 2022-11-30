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
)

//-------------------------------------------------------------

type GFconfig struct {
	Auth0domainStr      string
	Auth0apiAudienceStr string
}

type GFonLoginSuccessProfileInfo struct {
	SubStr      string
	NicknameStr string
	NameStr     string
	PictureURLstr string
}

//-------------------------------------------------------------

func Init() *GFconfig {

	config := loadConfig()
	
	return config
}

//-------------------------------------------------------------

// load Auth0 config, mostly from ENV vars
func loadConfig() *GFconfig {

	auth0domainStr := os.Getenv("AUTH0_DOMAIN")
	auth0apiAudienceStr := os.Getenv("AUTH0_AUDIENCE")

	config := &GFconfig{
		Auth0domainStr:      auth0domainStr,
		Auth0apiAudienceStr: auth0apiAudienceStr,
	}
	return config
}

//-------------------------------------------------------------