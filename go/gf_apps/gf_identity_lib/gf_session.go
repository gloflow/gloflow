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

package gf_identity_lib

import (
	"time"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// p_user_identifier_str - user ID or some other unique user identifier to be used

func session__set_on_req(p_session_data_str string,
	p_resp          http.ResponseWriter,
	p_ttl_hours_int int) {

	ttl    := time.Duration(p_ttl_hours_int) * time.Hour
	expire := time.Now().Add(ttl)
	cookie_name_str := "gf_sess_data"
	
	cookie := http.Cookie{
		Name:    cookie_name_str,
		Value:   p_session_data_str,
		Expires: expire,

		// IMPORTANT!! - make cookie http_only, disabling browser js context
		//               from reading its values
		HttpOnly: true,

		// SameSite allows a server to define a cookie attribute making it impossible for
		// the browser to send this cookie along with cross-site requests. The main
		// goal is to mitigate the risk of cross-origin information leakage, and provide
		// some protection against cross-site request forgery attacks.
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(p_resp, &cookie)
}

//---------------------------------------------------
func session__validate(p_req *http.Request,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	for _, cookie := range p_req.Cookies() {
		if (cookie.Name == "gf_sess_data") {
			session_data_str  := cookie.Value
			jwt_token_val_str := session_data_str

			//---------------------
			// JWT_VALIDATE
			gf_err := jwt__pipeline__validate(GF_jwt_token_val(jwt_token_val_str),
				p_ctx,
				p_runtime_sys)
			if gf_err != nil {
				return gf_err
			}

			//---------------------
		}
	}

	return nil
}