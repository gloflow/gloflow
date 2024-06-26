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

import * as gf_home from "./../ts/gf_home";

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

	main(log_fun);
});

//--------------------------------------------------------
function main(p_log_fun) {
	
	const my_eth_addresses_lst = [
		// "0xBA47Bef4ca9e8F86149D2f109478c6bd8A642C97",
		// "0xBA47Bef4ca9e8F86149D2f109478c6bd8A642C97",
		// "0xBA47Bef4ca9e8F86149D2f109478c6bd8A642C97"
	];

	const observed_eth_addresses_lst = [];

	const tags_by_address_map = {};

	const http_api_map = {
		//------------------------
        // TAGGING
        //------------------------
        "gf_tagger": {
            
            "add_tags_to_obj": async (p_tags_lst :string[],  
                p_object_id_str   :string,
                p_object_type_str :string,
                p_meta_map,
                p_log_fun)=>{
				
				const output_map = {
					"added_tags_lst": p_tags_lst
				};
				return output_map;
            }
        },
		
		//------------------------
		"home": {

			//------------------------
			// VIZ
			//------------------------

			"viz_get_fun": async ()=>{

				const output_map = {
					"components_map": {
						"profile_image": {
							"name_str":     "profile_image",
							"screen_x_int": 500,
							"screen_y_int": 113
						}
					},
				};
				return output_map;
			},
			"viz_update_fun": async ()=>{
				const output_map = {};
				return output_map;
			},

			//------------------------
			// WEB3
			//------------------------
			// WEB3_GET_ALL_ADDRESSES
			"web3_addresses_get_all_fun": async (p_type_str :string,
				p_chain_str :string)=>{

				var output_map;
				switch (p_type_str) {
					case "my":
						output_map = {
							"addresses_lst": my_eth_addresses_lst,
						};
						return output_map;
						break;
					
					case "observed":
						output_map = {
							"addresses_lst": observed_eth_addresses_lst
						};
						return output_map;
						break;
				}
			},

			//------------------------
			"web3_address_add_fun": async (p_address_str :string,
				p_type_str  :string,
				p_chain_str :string)=>{

				switch (p_type_str) {
					case "my":
						my_eth_addresses_lst.push(p_address_str);
						tags_by_address_map[p_address_str] = [];
						break;

					case "observed":
						observed_eth_addresses_lst.push(p_address_str);
						tags_by_address_map[p_address_str] = [];
						break;
				}
				
			},

			//------------------------
		}

	};

    const assets_paths_map = {
		"gf_add_btn":        "./../../../../assets/gf_add_btn.svg",
		"gf_confirm_btn":    "./../../../../assets/gf_confirm_btn.svg",
		"gf_bar_handle_btn": "./../../../../assets/gf_bar_handle_btn.svg"
	}

	const identity_http_api_map = {};

    gf_home.init(http_api_map,
		identity_http_api_map,
		assets_paths_map,
		p_log_fun);
}