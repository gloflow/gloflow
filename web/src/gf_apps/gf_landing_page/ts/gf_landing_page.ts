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
import * as gf_images         from "./gf_images";
import * as gf_procedural_art from "./procedural_art/gf_procedural_art";

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
	// init_remote(log_fun);
});

//--------------------------------------------------------
/*export function init_remote(p_log_fun) {

	init(remote_register_user_email, p_log_fun);
	//--------------------------------------------------------
	function remote_register_user_email(p_inputed_email_str :string,
		p_on_complete_fun,
		p_log_fun) {
		
		const url_str       = "/landing/register_invite_email";
		const data_args_map = {
			"email_str": p_inputed_email_str
		};
		
		$.ajax({
			"url":         url_str,
			"type":        "POST",
			"data":        JSON.stringify(data_args_map),
			"contentType": "application/json",
			"success":     (p_data_map)=>{
	     		p_on_complete_fun("success", p_data_map);
			}
		});
	}

	//--------------------------------------------------------
}*/

//--------------------------------------------------------
// INIT
function init(p_log_fun) {
	
	//---------------------
	// IDENTITY
	gf_identity.init();
	
	//---------------------
	
	const featured_elements_infos_lst = load_static_data(p_log_fun);
	
	gf_procedural_art.init(p_log_fun);
	// gf_email_registration.init(p_register_user_email_fun, p_log_fun);

	init_posts_img_num();
	gf_images.init(p_log_fun);

	// draw a new canvas when the view is resized, and delete the old one (with the old dimensions)
	$(window).resize(()=>{

		// small screen widths dont display procedural_art
		if ($(window).innerWidth() > 660) {

			// ABOUT_SECTION - when screen is small (for mobile) dont display it all.
			//                 has to removed here directly because its included in the template.
			$("#about_section").remove();

			gf_procedural_art.remove();
			gf_procedural_art.init(p_log_fun);
		}
	});

	// UPLOAD__INIT
	// use "" so that no host is set in URL's for issued requests
	// (forces usage of origin host that the page came from)
	const target_full_host_str = "";
	gf_upload__init(target_full_host_str);

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

	$("#about_section").on('click', function() {
		$("#about_section #desc").css("visibility", "visible");
	});
}

//--------------------------------------------------------
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