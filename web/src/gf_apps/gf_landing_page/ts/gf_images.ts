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
///<reference path="../../../d/jquery.timeago.d.ts" />

import * as gf_color        from "./../../../gf_core/ts/gf_color";
import * as gf_image_colors from "./../../../gf_core/ts/gf_image_colors";
import * as gf_time         from "./../../../gf_core/ts/gf_time";
import * as gf_image_viewer from "./../../../gf_core/ts/gf_image_viewer";
import * as gf_images_share from "./../../gf_images/ts/gf_images_core/gf_images_share";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_tagger__init_ui;
declare var gf_tagger__http_add_tags_to_obj;

//-------------------------------------------------
export function init(p_logged_in_bool,
	p_plugin_callbacks_map,
	p_gf_host_str,
	p_log_fun) {

	$('#featured_images_0').find('.image_info').each((p_i, p_image_info_element)=>{
		
		init_img(p_image_info_element);
	});
	$('#featured_images_1').find('.image_info').each((p_i, p_image_info_element)=>{
		
		init_img(p_image_info_element);
	});

	//-------------------------------------------------
	function init_img(p_image_info_element) {
		
		// CLEANUP - for images that dont come from some origin page (direct uploads, or generated images)
		//           this origin_page_url is set to empty string. check for that and remove it.
		// FIX!! - potentially on the server/template-generation side this div node shouldnt get included
		//         at all for images that dont have an origin_page_url.
		if ($(p_image_info_element).find(".origin_page_url a").text().trim() == "") {
			$(p_image_info_element).find(".origin_page_url").remove();
		}

		// some images dont have a title set, so for those remove the title element
		if ($(p_image_info_element).find(".image_title").text().trim() === "") {
			$(p_image_info_element).find(".image_title").remove();
		}

		//----------------------
		// IMAGE_VIEWER

		const image_id_str = $(p_image_info_element).data("img_system_id");
		const img_thumb_medium_url = $(p_image_info_element).find("img").data("img_thumb_medium_url");
		const img_thumb_large_url  = $(p_image_info_element).find("img").data("img_thumb_medium_url");
		const flows_names_lst = $(p_image_info_element).data("img_flows_names").split(",")

		gf_image_viewer.init(p_image_info_element,
			image_id_str,
			img_thumb_medium_url,
			img_thumb_large_url,
			flows_names_lst,
			p_log_fun);

		//----------------------
		// CREATION_DATE
		gf_time.init_creation_date(p_image_info_element, p_log_fun);

		//----------------------
		// IMAGE_PALLETE
		const img = $(p_image_info_element).find("img")[0];

		const assets_paths_map = {
			"copy_to_clipboard_btn": "/images/static/assets/gf_copy_to_clipboard_btn.svg",
		}
		gf_image_colors.init_pallete(img,
			assets_paths_map,
			(p_color_dominant_hex_str,
			p_colors_hexes_lst)=>{

				// set a few of the other needed elements to the same dominant color
				$(p_image_info_element).css("background-color", `#${p_color_dominant_hex_str}`);
				$(p_image_info_element).find(".image_title").css("background-color", `#${p_color_dominant_hex_str}`);
				$(p_image_info_element).find(".origin_page_url").css("background-color", `#${p_color_dominant_hex_str}`);


				//----------------------
				// COLOR_CLASSIFY
				const color_class_str = gf_color.classify(p_color_dominant_hex_str);

				//----------------------

				switch (color_class_str) {

					// LIGHT
					case "light":

						/*
						if background is light, then the text should be dark, so setting it here explicitly
						on dominant color classification.
						*/
						$(p_image_info_element).find(".image_title").css("color", "black");
						$(p_image_info_element).find(".origin_page_url a").css("color", "black");
						$(p_image_info_element).find(".creation_time").css("color", "black");
						$(p_image_info_element).find(".owner_user_name").css("color", "black");

						break;

					// DARK
					case "dark":

						/*
						css rules external to this function set the default color of
						text to white, so dark background dominant-color works fine.
						no need to set anything here yet.
						*/
						break;
				};
			});

		//----------------------

		// only initialize tagging UI for logged-in users
		if (p_logged_in_bool) {
			
			init_tagging(p_image_info_element,
				p_gf_host_str,
				p_log_fun);



			// SHARE
			gf_images_share.init(image_id_str,
				p_image_info_element,
				p_plugin_callbacks_map,
				p_log_fun);
		}
	}

	//-------------------------------------------------
}

//---------------------------------------------------
// TAGGING_UI

function init_tagging(p_image_container_element,
	p_gf_host_str,
	p_log_fun) {
	
	var image_system_id_str;

	const http_api_map = {
		"gf_tagger": {
			"add_tags_to_obj": async (p_new_tags_lst,
				p_obj_id_str,
				p_obj_type_str,
				p_tags_meta_map,
				p_log_fun)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
					
					const object_type_str = "img";

					await gf_tagger__http_add_tags_to_obj(p_new_tags_lst,
						image_system_id_str,
						object_type_str,
						{}, // meta_map
						p_gf_host_str,
						p_log_fun);

					p_resolve_fun({
						"added_tags_lst": p_new_tags_lst,
					});
				});
				return p;
			}
		}
	};

	const obj_type_str = "image";
	const input_element_parent_selector_str = "body";

	gf_tagger__init_ui(obj_type_str,
		p_image_container_element,
		input_element_parent_selector_str,

		//---------------------------------------------------
		// TAGS
		//---------------------------------------------------
		// tags_create_pre_fun
		async (p_tags_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				// p_tags_create_pre_fun resolves the system_id of the item being tagged
				/*
				IMPORTANT!! - img_system_id is attached as a data property to the image container element
					in the server template rendering.
				*/
				var image_system_id_str = $(p_image_container_element).attr("data-img_system_id_str");
				p_resolve_fun(image_system_id_str);
			});
			return p;
		},

		//---------------------------------------------------
		// on_tags_created_fun
		(p_tags_lst)=>{

			console.log("added tags >>>>>>>>>>>", p_tags_lst);

			p_tags_lst.forEach(p_tag_str=>{

				tag_display(p_tag_str);
			})
		},

		//---------------------------------------------------
		// NOTES
		//---------------------------------------------------
		// p_notes_create_pre_fun
		(p_notes_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				var image_system_id_str = $(p_image_container_element).attr("data-img_system_id_str");
				p_resolve_fun(image_system_id_str);
			});
			return p;
		},

		//---------------------------------------------------
		// on_notes_created_fun
		(p_notes_lst)=>{

			console.log("added notes >>>>>>>>>>>", p_notes_lst)
		},

		//---------------------------------------------------		
		()=>{}, // on_tagging_ui_add_fun
		()=>{}, // on_tagging_ui_remove_fun
		http_api_map,
		p_log_fun);

	//-------------------------------------------------
	function tag_display(p_tag_str) {

		$(p_image_container_element)
			.find(".tags_container")
			.append(`<a class='gf_image_tag' href='/v1/tags/objects?tag=${p_tag_str}&otype=image'>#${p_tag_str}</a>`)
	}

	//-------------------------------------------------
}