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

const (

	// USER
	GF_EVENT_APP__USER_CREATE_REGULAR = "GF_EVENT_APP__USER_CREATE_REGULAR"
	GF_EVENT_APP__USER_LOGIN          = "GF_EVENT_APP__USER_LOGIN"
	
	// ADMIN
	GF_EVENT_APP__ADMIN_CREATE                        = "GF_EVENT_APP__ADMIN_CREATE"
	GF_EVENT_APP__ADMIN_LOGIN                         = "GF_EVENT_APP__ADMIN_LOGIN"
	GF_EVENT_APP__ADMIN_LOGIN_PASS_CONFIRMED          = "GF_EVENT_APP__ADMIN_LOGIN_PASS_CONFIRMED"
	GF_EVENT_APP__ADMIN_LOGIN_EMAIL_VERIFICATION_SENT = "GF_EVENT_APP__ADMIN_LOGIN_EMAIL_VERIFICATION_SENT"
	GF_EVENT_APP__ADMIN_ADDED_USER_TO_INVITE_LIST     = "GF_EVENT_APP__ADMIN_ADDED_USER_TO_INVITE_LIST"
	GF_EVENT_APP__ADMIN_REMOVED_USER_FROM_INVITE_LIST = "GF_EVENT_APP__ADMIN_REMOVED_USER_FROM_INVITE_LIST"
)

