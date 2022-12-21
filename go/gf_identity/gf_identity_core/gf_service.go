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

//---------------------------------------------------

const (
	GF_AUTH_SUBSYSTEM_TYPE__USERPASS = "userpass"
	GF_AUTH_SUBSYSTEM_TYPE__AUTH0    = "auth0"
)

//---------------------------------------------------

type GFserviceInfo struct {
	
	// name of this service, in case multiple are spawned
	NameStr string

	// DOMAIN - where this gf_solo instance is reachable on
	DomainBaseStr string

	//------------------------
	// AUTH_SUBSYSTEM_TYPE - userpass | auth0
	AuthSubsystemTypeStr string

	//------------------------
	// ADMIN_MFA_SECRET_KEY_BASE32
	AdminMFAsecretKeyBase32str string

	//------------------------
	// AUTH

	// AUTH_LOGIN_URL - url of the login page to which the system should
	//                  redirect users after certain operations
	AuthLoginURLstr string

	// AUTH_LOGIN_SUCCESS_REDIRECT_URL - url to redirect to when the user 
	//                                   logs in successfuly. if ommited then dont redirect.
	AuthLoginSuccessRedirectURLstr string

	//------------------------
	// FEATURE_FLAGS

	// EVENTS_APP - enable sending of app events from various functions
	EnableEventsAppBool bool

	// enable storage of user_creds in a secret store
	EnableUserCredsInSecretsStoreBool bool

	// enable sending of emails for any function that needs it
	EnableEmailBool bool

	// enable login only for users that have confirmed their email
	EnableEmailRequireConfirmForLoginBool bool

	// enable login only for users that have confirmed their MFA code
	EnableMFArequireConfirmForLoginBool bool

	//------------------------
}