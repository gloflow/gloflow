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

///<reference path="../../../d/jquery.d.ts" />

import * as gf_core_utils     from "./../../../gf_core/ts/gf_utils";
import * as gf_identity       from "./../../../gf_identity/ts/gf_identity";
import * as gf_identity_http  from "./../../../gf_identity/ts/gf_identity_http";
import * as gf_flows_picker   from "./../../gf_images/ts/gf_images_flows_browser/gf_flows_picker";
import * as gf_tags_picker    from "./../../gf_tagger/ts/gf_tags_picker/gf_tags_picker";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;

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
	// FLOWS_PICKER - display it if the user is logged in
	if (logged_in_bool) {

		gf_flows_picker.init(p_log_fun)
	}

	// TAGS_PICKER - display it if the user is logged in
	if (logged_in_bool) {

		gf_tags_picker.init(p_log_fun)
	}

	//---------------------
	// WINDOW_RESIZE - draw a new canvas when the view is resized, and delete the old one (with the old dimensions)
	$(window).on("resize", ()=>{

		// small screen widths dont display procedural_art
		if ($(window).innerWidth() > 660) {

        
		}
	});

	//---------------------
	// UPLOAD__INIT

	var   upload_flow_name_str;
	const default_flow_name_str  = "general";


	// LOCAL_STORAGE
	const previous_flow_name_str = localStorage.getItem("gf:upload_flow_name_str");
	
	// get old value from localStorage if it exists, if it doesnt use the default
	if (previous_flow_name_str == null) {
		upload_flow_name_str = default_flow_name_str;
	} else {
		upload_flow_name_str = previous_flow_name_str;
	}

	// use "" so that no host is set in URL's for issued requests
	// (forces usage of origin host that the page came from)
	const target_full_host_str = "";

	gf_upload__init(upload_flow_name_str,
		target_full_host_str,
		(p_upload_gf_image_id_str)=>{

		});

	//---------------------
	// ABOUT_SECTION
	$("#about_section").on('click', function() {
		$("#about_section #desc").css("visibility", "visible");
	});

	//---------------------
	// PLUGINS
	// INIT - allows for arbitrary init code to be run when landing_page is done initializing
	if ("init" in p_plugin_callbacks_map) {

		p_plugin_callbacks_map["init"]();
	}

	//---------------------
}