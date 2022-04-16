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

import * as gf_dragndrop          from "./../../../gf_core/ts/gf_dragndrop";
import * as gf_home_eth_addresses from "gf_home_eth_addresses";
import * as gf_utils              from "gf_utils";

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



	// GET_PERSISTED_COMPONENTS
	const home_viz_map = await p_http_api_map["home"]["viz_get_fun"]();
	const home_viz_components_map = home_viz_map["components_map"];


	const home_container          = $("#home_container");
	const names_container         = init_names_view(home_container, p_http_api_map, p_assets_paths_map);
	const profile_image_container = await init_profile_image(p_http_api_map, p_assets_paths_map);
	const my_eth_addresses_container        = await gf_home_eth_addresses.init_my(home_container, p_http_api_map, p_assets_paths_map);
	const observed_eth_addresses_container  = await gf_home_eth_addresses.init_observed(home_container, p_http_api_map, p_assets_paths_map);
	const background_color_picker_container = await init_color_picker(home_container, p_http_api_map, p_assets_paths_map);

	//-----------------------------
	// PROFILE_IMAGE
	var profile_img_x_int;
	var profile_img_y_int;
	if ("profile_image" in home_viz_components_map) {
		profile_img_x_int = home_viz_components_map["profile_image"]["x_int"];
		profile_img_y_int = home_viz_components_map["profile_image"]["y_int"];
	}
	else {
		// default positioning
		profile_img_x_int = ($(window).width() - $(profile_image_container).width())/2;
		profile_img_y_int = 100;
	}

	$(profile_image_container).css("left", `${profile_img_x_int}px`);
	$(profile_image_container).css("top",  `${profile_img_y_int}px`);


	const profile_image_pos = $(profile_image_container).position();
	
	//-----------------------------
	// NAMES
	var names_x_int;
	var names_y_int;
	if ("names" in home_viz_components_map) {
		names_x_int = home_viz_components_map["names"]["x_int"];
		names_y_int = home_viz_components_map["names"]["y_int"];
	}
	else {
		// default positioning
		names_x_int = profile_img_x_int;
		names_y_int = profile_img_y_int - 38;
	}

	$(names_container).css("left", `${names_x_int}px`);
	$(names_container).css("top", `${names_y_int}px`);

	//-----------------------------
	// MY_ETH_ADDRESSES
	var eth_addr_x_int;
	var eth_addr_y_int;
	if ("web3_addresses_my" in home_viz_components_map) {
		eth_addr_x_int = home_viz_components_map["web3_addresses_my"]["x_int"];
		eth_addr_y_int = home_viz_components_map["web3_addresses_my"]["y_int"];
	}
	else {
		// default positioning
		eth_addr_x_int = profile_image_pos.left;
		eth_addr_y_int = profile_image_pos.top + $(profile_image_container).outerHeight() + 10;
	}

	$(my_eth_addresses_container).css("left", `${eth_addr_x_int}px`);
	$(my_eth_addresses_container).css("top",  `${eth_addr_y_int}px`);
	
	//-----------------------------
	// OBSERVED_ETH_ADDRESSES
	var obs_eth_addr_x_int;
	var obs_eth_addr_y_int;
	if ("web3_addresses_observed" in home_viz_components_map) {
		obs_eth_addr_x_int = home_viz_components_map["web3_addresses_observed"]["x_int"];
		obs_eth_addr_y_int = home_viz_components_map["web3_addresses_observed"]["y_int"];
	}
	else {
		// default positioning
		obs_eth_addr_x_int = eth_addr_x_int;
		obs_eth_addr_y_int = eth_addr_y_int + $(my_eth_addresses_container)[0].getBoundingClientRect().height + 10 ;
	}

	$(observed_eth_addresses_container).css("left", `${obs_eth_addr_x_int}px`);
	$(observed_eth_addresses_container).css("top",  `${obs_eth_addr_y_int}px`);
	
	//-----------------------------
	// COLOR_PICKER
	var background_color_picker_x_int;
	var background_color_picker_y_int;
	if ("background_color_picker" in home_viz_components_map) {
		background_color_picker_x_int = home_viz_components_map["background_color_picker"]["x_int"];
		background_color_picker_y_int = home_viz_components_map["background_color_picker"]["y_int"];

		// persist background_coilor
		const background_color_str = home_viz_components_map["background_color_picker"]["background_color_str"];
		$("body").css("background-color", background_color_str);
	}
	else {
		// default positioning
		background_color_picker_x_int = profile_img_x_int + $(profile_image_container)[0].getBoundingClientRect().width + 10;
		background_color_picker_y_int = profile_img_y_int;
	}

	$(background_color_picker_container).css("left", `${background_color_picker_x_int}px`);
	$(background_color_picker_container).css("top",  `${background_color_picker_y_int}px`);
	
    //---------------------------
}

//--------------------------------------------------------
function init_color_picker(p_parent_element,
	p_http_api_map,
	p_assets_paths_map) {

	const container = $(`
		<div id="background_color_picker">
			<div id="control"></div>
			<div id="picked_color">#000000</div>
		</div>`);
	$(p_parent_element).append(container);

    const color_picker_color_element = $(container).find("#picked_color");
    const color_picker = new iro.ColorPicker("#background_color_picker #control", {
        width: 100,   // size of the picker
        color: "#f00" // initial color
    });
    color_picker.on('color:change', async (p_color)=>{
        
        const background_color_hex_str = p_color.hexString;

        $(color_picker_color_element).text(background_color_hex_str);
		$("body").css("background-color", `${background_color_hex_str}`);

		// update component remotely
		const component_name_str = "background_color_picker";
		await gf_utils.update_viz_background_color(component_name_str,
			background_color_hex_str,
			p_http_api_map);
    });

	// DRAG_N_DROP
	gf_dragndrop.init(container,
		//--------------------------------------------------------
		// p_on_dnd_event_fun
		async (p_dnd_event_type_str :string, p_drag_data_map)=>{

			switch (p_dnd_event_type_str) {
				case "drag_start":
					break;

				case "drag_stop":

					// update component remotely on each drag/coord change
					const component_name_str = "background_color_picker";
					await gf_utils.update_viz_component_remote(component_name_str,
						p_drag_data_map,
						p_http_api_map);
						
					break;
			}
		},
		
		//--------------------------------------------------------
		p_assets_paths_map);

	return container;
}

//--------------------------------------------------------
function init_names_view(p_parent_element,
	p_http_api_map,
	p_assets_paths_map) {
	const container = $(`
		<div id="names">
			<div>@ivan</div>
		</div>`);
	$(p_parent_element).append(container);


	// DRAG_N_DROP
	gf_dragndrop.init(container,
		//--------------------------------------------------------
		// p_on_dnd_event_fun
		async (p_dnd_event_type_str :string, p_drag_data_map)=>{

			switch (p_dnd_event_type_str) {
				case "drag_start":
					break;

				case "drag_stop":

					// update component remotely on each drag/coord change
					const component_name_str = "names";
					await gf_utils.update_viz_component_remote(component_name_str,
						p_drag_data_map,
						p_http_api_map);

					break;
			}
		},
		
		//--------------------------------------------------------
		p_assets_paths_map);

	return container;
}

//--------------------------------------------------------
function init_profile_image(p_http_api_map,
	p_assets_paths_map) {
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
		
		  
		$(container).find("img").one("load", function() {
			
			// DRAG_N_DROP
			gf_dragndrop.init(container,
				//--------------------------------------------------------
				// p_on_dnd_event_fun
				async (p_dnd_event_type_str :string, p_drag_data_map)=>{

					switch (p_dnd_event_type_str) {
						case "drag_start":
							break;

						case "drag_stop":

							// update component remotely on each drag/coord change
							const component_name_str = "profile_image";
							await gf_utils.update_viz_component_remote(component_name_str,
								p_drag_data_map,
								p_http_api_map);

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