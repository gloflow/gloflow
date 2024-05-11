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

///<reference path="./../../../d/jquery.d.ts" />

import * as gf_core_utils    from "./../../../gf_core/ts/gf_utils";
import * as gf_identity      from "./../../../gf_identity/ts/gf_identity";
import * as gf_identity_http from "./../../../gf_identity/ts/gf_identity_http";
import * as gf_identity_mfa  from "./../../../gf_identity/ts/gf_identity_mfa";

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
	// META
	const notifications_meta_map = {

		// for admin logins, success of first_stage login leads to a confirmation email being sent to the
		// user, which requires them to follow a link for the second stage to succeed
		// (which may lead to further login stages, MFA/etc.)
		"login_first_stage_success": "login success - please check your email for a confirmation message"
	};

	//---------------------
	// IDENTITY

	const current_host_str = gf_core_utils.get_current_host();
	const urls_map = gf_identity_http.get_admin_http_urls(current_host_str);
	gf_identity.init_with_http(notifications_meta_map, urls_map, current_host_str);

	
	const http_api_map = gf_identity_http.get_http_api(urls_map, current_host_str);

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