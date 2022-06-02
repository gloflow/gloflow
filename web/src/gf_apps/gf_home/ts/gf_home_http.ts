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

//-------------------------------------------------
export function get_http_api() {
    const http_api_map = {
		"home": {

			//------------------------
			// VIZ
			//------------------------
			"viz_get_fun": async ()=>{
                const output_map = await viz_get();
				return output_map;
			},
			"viz_update_fun": async (p_component_name_str :string,
				p_props_change_map)=>{

                const output_map = await viz_update(p_component_name_str,
					p_props_change_map);
                return output_map;
			},

			//------------------------
			// WEB3
			//------------------------
			// WEB3_ADDRESSES_GET_ALL
			"web3_addresses_get_all_fun": async (p_type_str :string,
				p_chain_str :string)=>{
				const output_map = await web3_addresses_get_all(p_type_str,
					p_chain_str);
				return output_map;
			},

			//------------------------
			// WEB3_ADD_ADDRESS
			"web3_address_add_fun": async (p_address_str :string,
				p_type_str  :string,
				p_chain_str :string)=>{

				const output_map = await web3_address_add(p_address_str,
					p_type_str,
					p_chain_str);
				return output_map;
			},
            
            // WEB3_NFT_INDEX_FOR_ADDRESS
            "web3_nft_index_for_address_fun": async (p_address_str :string,
                p_chain_str :string)=>{

                const output_map = await web3_nft_index_for_address(p_address_str,
                    p_chain_str);
                return output_map;
            },

			//------------------------
		}

	};
    return http_api_map;
}

//--------------------------------------------------------
async function viz_get() {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const url_str = "/v1/home/viz/get"
        $.ajax({
            'url':         url_str,
            'type':        'GET',
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                const status_str = p_response_map["status"];
                const data_map   = p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error':(jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//--------------------------------------------------------
async function viz_update(p_component_name_str :string,
    p_props_change_map) {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const url_str = "/v1/home/viz/update"
		const data_map = {
            "component_name_str": p_component_name_str,
            "props_change_map":   p_props_change_map,
        };

        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                const status_str = p_response_map["status"];
                const data_map   = p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//--------------------------------------------------------
export async function web3_nft_index_for_address(p_address_str :string,
    p_chain_str :string) {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const url_str = `/v1/web3/nft/index_address`;
        const data_map = {
            "address_str": p_address_str,
            "chain_str":   p_chain_str,
        };

        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                const status_str = p_response_map["status"];
                const data_map   = p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//--------------------------------------------------------
export async function web3_addresses_get_all(p_type_str  :string,
	p_chain_str :string) {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const url_str = `/v1/web3/address/get_all?type=${p_type_str}&chain=${p_chain_str}`;

        $.ajax({
            'url':         url_str,
            'type':        'GET',
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                const status_str = p_response_map["status"];
                const data_map   = p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//--------------------------------------------------------
export async function web3_address_add(p_address_str :string,
	p_type_str  :string,
	p_chain_str :string) {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const url_str = "/v1/web3/address/add"
		const data_map = {
            "address_str": p_address_str,
            "type_str":    p_type_str,
            "chain_str":   p_chain_str,
        };

        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                const status_str = p_response_map["status"];
                const data_map   = p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}