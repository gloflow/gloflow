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

import * as gf_home from "gf_home";
import * as gf_home_http from "gf_home_http";
import * as gf_identity_http from "./../../../gf_identity/ts/gf_identity_http";
import * as gf_core_utils    from "./../../../gf_core/ts/gf_utils";

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
	// $("time.timeago").timeago();

	const assets_paths_map = {
		"gf_add_btn":        "/images/static/assets/gf_add_btn.svg",
		"gf_confirm_btn":    "/images/static/assets/gf_confirm_btn.svg",
		"gf_bar_handle_btn": "/images/static/assets/gf_bar_handle_btn.svg"
	}


	const current_host_str = gf_core_utils.get_current_host();
	const http_api_map          = gf_home_http.get_http_api();
	const identity_urls_map     = gf_identity_http.get_standard_http_urls(current_host_str);
	const identity_http_api_map = gf_identity_http.get_http_api(identity_urls_map, current_host_str);

	gf_home.init(http_api_map,
		identity_http_api_map,
		assets_paths_map,
		log_fun);
});