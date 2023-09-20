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

///<reference path="../../../d/jquery.d.ts" />

import * as gf_identity       from "./../../../gf_identity/ts/gf_identity";
import * as gf_identity_http  from "./../../../gf_identity/ts/gf_identity_http";
import * as gf_images         from "./gf_images";
import * as gf_procedural_art from "./procedural_art/gf_procedural_art";
import * as gf_image_colors   from "./../../../gf_core/ts/gf_image_colors";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;

//--------------------------------------------------------
// INIT
export async function init(p_log_fun) {
	
	const gf_host_str = window.location.href;

    $("time.timeago").timeago();
    
	//---------------------
	// META
	const notifications_meta_map = {
		"login_first_stage_success": "login success"
	};

	//---------------------
	// POSTS_INIT
	posts_init();

	//---------------------
	// GF_IMAGES_INIT
	gf_images.init(gf_host_str, p_log_fun);

	//---------------------
	// WINDOW_RESIZE - draw a new canvas when the view is resized, and delete the old one (with the old dimensions)
	$(window).on("resize", ()=>{

		// small screen widths dont display procedural_art
		if ($(window).innerWidth() > 660) {

			// ABOUT_SECTION - when screen is small (for mobile) dont display it all.
			//                 has to removed here directly because its included in the template.
			$("#about_section").remove();

			gf_procedural_art.remove();
			gf_procedural_art.init(p_log_fun);
		}
	});

	//---------------------
	// GF_PROCEDURAL_ART
	gf_procedural_art.init(p_log_fun);

	// regenerate new piece on click 
	$("#randomized_art").on("click", ()=>{
		gf_procedural_art.remove();
		gf_procedural_art.init(p_log_fun);
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
	// IDENTITY
	// first complete main initialization and only then initialize gf_identity
	const urls_map = gf_identity_http.get_standard_http_urls();
	const auth_http_api_map = gf_identity_http.get_http_api(urls_map);
	gf_identity.init_with_http(notifications_meta_map, urls_map);
	

	
	const parent_node = $("#right_section");
	const home_url_str = urls_map["home"];

	gf_identity.init_me_control(parent_node,
		auth_http_api_map,
		home_url_str);

	//---------------------
	// ABOUT_SECTION
	$("#about_section").on('click', function() {
		$("#about_section #desc").css("visibility", "visible");
	});

	//---------------------
}

//--------------------------------------------------------
function posts_init() {

	init_posts_img_num();

	$('#featured_posts').find('.post_info').each((p_i, p_post_info_element)=>{
		
		//----------------------
		// IMAGE_PALLETE
		const img = $(p_post_info_element).find("img")[0];

		const assets_paths_map = {
			"copy_to_clipboard_btn": "/images/static/assets/gf_copy_to_clipboard_btn.svg",
		}
		gf_image_colors.init_pallete(img,
			assets_paths_map,
			(p_color_dominant_hex_str,
			p_colors_hexes_lst)=>{

				// set the background color of the post to its dominant color
				$(p_post_info_element).css("background-color", `#${p_color_dominant_hex_str}`);

			});

		//----------------------
	});

	//--------------------------------------------------------
	function init_posts_img_num() {

		$("#featured_posts .post_info").each((p_i, p_post)=>{

			const post_images_number = $(p_post).find(".post_images_number")[0];
			const label_element      = $(post_images_number).find(".label");

			// HACK!! - "-1" was visually inferred
			$(post_images_number).css("right", `-${$(post_images_number).outerWidth()-1}px`);
			$(label_element).css("left", `${$(post_images_number).outerWidth()}px`);

			$(p_post).mouseover((p_e)=>{
				$(post_images_number).css("visibility", "visible");
			});
			$(p_post).mouseout((p_e)=>{
				$(post_images_number).css("visibility", "hidden");
			});
		});
	}

	//--------------------------------------------------------
}