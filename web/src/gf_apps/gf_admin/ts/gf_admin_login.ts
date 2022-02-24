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

///<reference path="../../../d/jquery.d.ts" />

import * as gf_identity     from "./../../gf_identity/ts/gf_identity";
import * as gf_identity_mfa from "./../../gf_identity/ts/gf_identity_mfa";

//--------------------------------------------------------
$(document).ready(()=>{
	//-------------------------------------------------
	function log_fun(p_g, p_m) {
		var msg_str = p_g+':'+p_m;
		switch (p_g) {
			case "INFO":
				console.log("%cINFO"+":"+"%c"+p_m, "color:green; background-color:#ACCFAC;", "background-color:#ACCFAC;");
				break;
			case "FUN_ENTER":
				console.log("%cFUN_ENTER"+":"+"%c"+p_m, "color:yellow; background-color:lightgray", "background-color:lightgray");
				break;
		}
	}

	//-------------------------------------------------
	init(log_fun);
});

//-------------------------------------------------
function init(p_log_fun) {


    console.log("admin login")


	//---------------------
	// IDENTITY
	const urls_map = gf_identity.get_admin_http_urls();
	gf_identity.init_with_http(urls_map);

	
	const http_api_map = gf_identity.get_http_api(urls_map);

	// MFA_CONFIRM - check if login is in MFA_confirm stage of login
	if ($('#identiy_mfa_confirm').length > 0) {

		// when backend redirects user to MFA confirm page, it appends
		// the user_name to query_string args.
		const url_params    = new URLSearchParams(window.location.search);
		const user_name_str = url_params.get('user_name');

		const mfa_container = gf_identity_mfa.init(user_name_str,
			http_api_map,
			//-------------------------------------------------
			(p_mfa_valid_bool)=>{

				// MFA_VALID
				if (p_mfa_valid_bool) {

					// redirect to dashboard
					window.location.pathname = "/v1/admin/dashboard";
				}
			});

			//-------------------------------------------------

		$("body").append(mfa_container);
	}

	//---------------------
}