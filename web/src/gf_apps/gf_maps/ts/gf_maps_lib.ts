/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

// ///<reference path="../../../d/jquery.d.ts" />

import * as gf_core_utils     from "./../../../gf_core/ts/gf_utils";
import * as gf_identity       from "./../../../gf_identity/ts/gf_identity";
import * as gf_identity_http  from "./../../../gf_identity/ts/gf_identity_http";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var L;

//--------------------------------------------------------
// INIT
export async function init(p_plugin_callbacks_map,
	p_log_fun) {
	
		
	const current_host_str = gf_core_utils.get_current_host();


    $("time.timeago").timeago();
    
	//---------------------
	// META
	const notifications_meta_map = {
		"login_first_stage_success": "login success"
	};

	//---------------------
	// IDENTITY
	// first complete main initialization and only then initialize gf_identity
	const urls_map          = gf_identity_http.get_standard_http_urls(current_host_str);
	const auth_http_api_map = gf_identity_http.get_http_api(urls_map, current_host_str);
	gf_identity.init_with_http(notifications_meta_map, urls_map, current_host_str);
	

	
	const parent_node = $("#right_section");
	const home_url_str = urls_map["home"];

	gf_identity.init_me_control(parent_node,
		auth_http_api_map,
		home_url_str);
	
	// inspect if user is logged-in or not
	const logged_in_bool = await auth_http_api_map["general"]["logged_in"]();



	//---------------------
	// WINDOW_RESIZE - draw a new canvas when the view is resized, and delete the old one (with the old dimensions)
	$(window).on("resize", ()=>{

		// small screen widths dont display procedural_art
		if ($(window).innerWidth() > 660) {

        
		}
	});

	//---------------------

	init_map();
}

//--------------------------------------------------------
function init_map() {
	var map = L.map('map').setView([51.505, -0.09], 13);

	L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
		maxZoom: 19,
		attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);



	var polygon = L.polygon([
		[51.509, -0.08],
		[51.503, -0.06],
		[51.51, -0.047]
	]).addTo(map);

	var circle = L.circle([51.508, -0.11], {
		color: 'red',
		fillColor: '#f03',
		fillOpacity: 0.5,
		radius: 500
	}).addTo(map);
	
}