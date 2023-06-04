/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

import * as gf_identity       from "./../../gf_identity/ts/gf_identity";
import * as gf_identity_http  from "./../../gf_identity/ts/gf_identity_http";
import * as gf_images         from "./gf_images";
import * as gf_procedural_art from "./procedural_art/gf_procedural_art";
import * as gf_image_colors   from "./../../../gf_core/ts/gf_image_colors";

// import * as gf_calc from "./gf_calc";
// import * as gf_email_registration from "./gf_email_registration";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;

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
	$("time.timeago").timeago();

	init(log_fun);
});

//--------------------------------------------------------
// INIT
export function init(p_log_fun) {
	
	//---------------------
	// META
	const notifications_meta_map = {
		"login_first_stage_success": "login success"
	};

	//---------------------
	// IDENTITY
	const urls_map = gf_identity_http.get_standard_http_urls();
	gf_identity.init_with_http(notifications_meta_map, urls_map);
	
	//---------------------
	
	// const featured_elements_infos_lst = load_static_data(p_log_fun);
	
	gf_procedural_art.init(p_log_fun);
	// gf_email_registration.init(p_register_user_email_fun, p_log_fun);

	posts_init();
	gf_images.init(p_log_fun);

	// draw a new canvas when the view is resized, and delete the old one (with the old dimensions)
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

	$("#about_section").on('click', function() {
		$("#about_section #desc").css("visibility", "visible");
	});
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

//--------------------------------------------------------
/*
function load_static_data(p_log_fun) :Object[] {
	
	const featured_elements_infos_lst :Object[] = []; 

	$("#posts .post_info").each((p_i)=>{
		const element = this;
		const featured_element_image_url_str  :string = $(element).find("img").attr("src");
		const featured_element_images_num_str :string = $(element).find(".post_images_number").find(".num").text();
		const featured_element_title_str      :string = $(element).find(".post_title").text();

		const featured_element_info_map :Object = {
			"element":    $(element),
			"image_src":  featured_element_image_url_str,
			"images_num": featured_element_images_num_str,
			"title_str":  featured_element_title_str
		};

		featured_elements_infos_lst.push(featured_element_info_map);
	});

	return featured_elements_infos_lst;
}
*/