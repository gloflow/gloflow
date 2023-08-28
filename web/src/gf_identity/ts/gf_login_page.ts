/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

///<reference path="../../d/jquery.d.ts" />

import * as gf_identity from "./gf_identity";
import * as gf_identity_http  from "./gf_identity_http";

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

//--------------------------------------------------------
function init(p_log_fun) {

	console.log("gf login UI...")
	
	let url    = new URL(window.location.href);
	let params = new URLSearchParams(url.search);

	/*
	IMPORTANT!! - currently in the Auth0 login_callback, once all login initialization is complete,
				the client is redirected to the login_page with a QS arg "login_success" set to 1.
				this indicates to this page it should not initialize into its standard login state,
				but instead should indicate to the user that login has succeeded and then via JS
				(with a time delay) redirect to Home.
				this is needed to handle a race condition where if user was redirected to home via a
				server redirect (HTTP 3xx response) the browser wouldnt have time to set the needed
				auth cookies that are necessary for the Home handler to authenticate.
				the recommended solution for this is to do the redirection via client/JS code with a slight
				time delay, giving the browser time to set the cookies.
	*/
	if (params.has('login_success')) {
		let login_success_str = params.get('login_success');
		
		if (login_success_str == "1") {

			$("#welcome .label").append(`
				<div class='label' style='font-weight: bold;'>
					redirecting to Home...
				</div>`);

			setTimeout(function() {

				const url_str = "/v1/home/view";				
				window.location.href = url_str;
			}, 3000);
		}
		else {

			console.log("login failed...");

			/*
			standard login state initialized, where the user is redirected to this page
			exclusively to login (on auth failure when trying to access some auth-ed endpoint).
			*/
			init_standard();
		}
	}
	else {

		console.log("standard login page init...");

		/*
		standard login state initialized, where the user is redirected to this page to login.
		*/
		init_standard();
	}

	//--------------------------------------------------------
	function init_standard() {

		$("#welcome .label").text("Welcome to GF Login");

		const urls_map = gf_identity_http.get_standard_http_urls();
		const auth_http_api_map = gf_identity_http.get_http_api(urls_map);
		
		// META
		const notifications_meta_map = {
			"login_first_stage_success": "login success"
		};

		gf_identity.init_with_http(notifications_meta_map, urls_map);

		


		const parent_node = $("#welcome")[0];
		const home_url_str = urls_map["home"];

		gf_identity.init_me_control(parent_node,
			auth_http_api_map,
			home_url_str);
	
	
	}

	//--------------------------------------------------------


}