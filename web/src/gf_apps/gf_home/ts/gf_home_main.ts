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



	const http_api_map = {
		"home": {
			//------------------------
			// MY_ETH_ADDRESSES
			"get_my_eth_addresses_fun": async ()=>{
				gf_home_http
			},

			//------------------------
			// OBSERVED_ETH_ADDRESSES
			"get_observed_eth_addresses_fun": async ()=>{

			},

			//------------------------
			// ADD_ETH_ADDRESS
			"add_eth_address_fun": async (p_address_str :string,
				p_type_str :string)=>{




			},

			//------------------------
		}

	};
	gf_home.init(http_api_map, assets_paths_map, log_fun);
});