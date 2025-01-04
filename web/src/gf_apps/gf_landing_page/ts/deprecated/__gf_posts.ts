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
///<reference path="../../../d/jquery.timeago.d.ts" />

import * as gf_post_control from "../../../gf_publisher/ts/gf_posts_core/gf_post_control";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_tagger__init_ui_v2;
declare var gf_tagger__http_add_tags_to_obj;

//--------------------------------------------------------
// POSTS_INIT

export function init(p_gf_host_str :string,
	p_log_fun) {

	const http_api_map = {

		// GF_TAGGER
		"gf_tagger": {
			"add_tags_to_obj": async (p_new_tags_lst,
				p_obj_id_str,
				p_obj_type_str,
				p_tags_meta_map,
				p_log_fun)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

					await gf_tagger__http_add_tags_to_obj(p_new_tags_lst,
						p_obj_id_str,
						p_obj_type_str,
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

	// init_posts_img_num();

	$('#featured_posts').find('.post_info').each((p_i, p_post_info_element)=>{
		
		gf_post_control.init_existing_dom(p_post_info_element, http_api_map, p_log_fun)

		/*
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
		*/
	});

	//--------------------------------------------------------
	/*
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
	*/
	//--------------------------------------------------------
}