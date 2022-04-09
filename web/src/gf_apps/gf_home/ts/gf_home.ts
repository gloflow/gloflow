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

import * as gf_dragndrop from "./../../../gf_core/ts/gf_dragndrop";
import * as gf_home_eth_addresses from "gf_home_eth_addresses";

declare var WebFont; 
declare var iro;

//--------------------------------------------------------
// INIT
export async function init(p_http_api_map,
	p_assets_paths_map,
	p_log_fun) {

	// FONTS
	WebFont.load({
		google: {
			families: ['IBM Plex Sans']
		},
		/*loading: function() {
			console.log("Fonts are being loaded");
		},
		active: function() {
			console.log("Fonts have been rendered")
		}*/
	});
	
	const home_container          = $("#home_container");
	const names_container         = init_names_view(home_container, p_assets_paths_map);
	const profile_image_container = await init_profile_image(p_assets_paths_map);
	const my_eth_addresses_container       = await gf_home_eth_addresses.init_my(home_container, p_http_api_map, p_assets_paths_map);
	const observed_eth_addresses_container = await gf_home_eth_addresses.init_observed(home_container, p_http_api_map, p_assets_paths_map);
	
	//-----------------------------
	// PROFILE_IMAGE
	const profile_image_pos = $(profile_image_container).position();
	
	//-----------------------------
	// NAMES
	const names_x = profile_image_pos.left;
	const names_y = profile_image_pos.top - 38;

	$(names_container).css("left", `${names_x}px`);
	$(names_container).css("top", `${names_y}px`);

	//-----------------------------
	// MY_ETH_ADDRESSES

	const eth_addr_x = profile_image_pos.left;
	const eth_addr_y = profile_image_pos.top + $(profile_image_container).outerHeight() + 10 ;

	$(my_eth_addresses_container).css("left", `${eth_addr_x}px`);
	$(my_eth_addresses_container).css("top",  `${eth_addr_y}px`);
	
	//-----------------------------
	// OBSERVED_ETH_ADDRESSES

	const obs_eth_addr_x = eth_addr_x;
	const obs_eth_addr_y = eth_addr_y + $(my_eth_addresses_container)[0].getBoundingClientRect().height + 10 ;

	$(observed_eth_addresses_container).css("left", `${obs_eth_addr_x}px`);
	$(observed_eth_addresses_container).css("top",  `${obs_eth_addr_y}px`);
	
	//-----------------------------



	//---------------------------
    // COLOR_PICKER

	const color_picker_container = $(`
		<div id="background_color_picker">
			<div id="control"></div>
			<div id="picked_color">#000000</div>
		</div>`);
	$(home_container).append(color_picker_container);

    const color_picker_color_element = $(color_picker_container).find("#picked_color");
    const color_picker = new iro.ColorPicker("#background_color_picker #control", {
        width: 100,   // size of the picker
        color: "#f00" // initial color
    });
    color_picker.on('color:change', (p_color)=>{
        
        const picked_color_hex_str = p_color.hexString;

        $(color_picker_color_element).text(picked_color_hex_str);
        
		
		$("body").css("background-color", `${picked_color_hex_str}`);
    });

    //---------------------------
}

//--------------------------------------------------------
function init_names_view(p_parent_element, p_assets_paths_map) {
	const container = $(`
		<div id="names">
			<div>@ivan</div>
		</div>`);
	$(p_parent_element).append(container);


	// DRAG_N_DROP
	gf_dragndrop.init(container,
		//--------------------------------------------------------
		// p_on_dnd_event_fun
		(p_dnd_event_type_str :string)=>{

			switch (p_dnd_event_type_str) {
				case "drag_start":
					break;

				case "drag_stop":
					break;
			}
		},
		
		//--------------------------------------------------------
		p_assets_paths_map);

	return container;
}

//--------------------------------------------------------
function init_profile_image(p_assets_paths_map) {
	const p = new Promise(function(p_resolve_fun, p_reject_fun) {

		const container = $(`
			<div id="profile_image">
				<img class="image_preview"
					src="https://media.gloflow.com/thumbnails/786f79c0c85c08c7b1c0b3e11d6cae1e_thumb_small.png"
					draggable="false"></img>
			</div>`);
		
		const home_container = $("#home_container");
		$(home_container).append(container);

		
		console.log($(window).width(), $(container).width())
		const x = ($(window).width() - $(container).width())/2;
		const y = 100;

		$(container).css("left", `${x}px`);
		$(container).css("top",  `${y}px`);
		  
		$(container).find("img").one("load", function() {
			
			// DRAG_N_DROP
			gf_dragndrop.init(container,
				//--------------------------------------------------------
				// p_on_dnd_event_fun
				(p_dnd_event_type_str :string)=>{

					switch (p_dnd_event_type_str) {
						case "drag_start":
							break;

						case "drag_stop":
							break;
					}
				},
				
				//--------------------------------------------------------
				p_assets_paths_map);

			p_resolve_fun(container);
		});
	});
	return p;
}