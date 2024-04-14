/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

///<reference path="../../../../d/jquery.d.ts" />

import * as gf_image_viewer from "../../../../gf_core/ts/gf_image_viewer";
import * as gf_gifs_viewer  from "../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_images_http  from "./gf_images_http";

declare var gf_tagger__init_ui_v2;
declare var gf_tagger__http_add_tags_to_obj;

//---------------------------------------------------
// TAGGING_UI

export function init_tagging(p_image_id_str,
	p_image_container_element,
	p_gf_host_str,
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
		},

		// GF_IMAGES
		"gf_images": {
			"classify_image": async (p_image_id_str)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

					const client_type_str = "web";

					await gf_images_http.classify(p_image_id_str,
						client_type_str,
						p_log_fun);
				});
				return p;
			}
		}
	};

	const obj_type_str = "image";

	const callbacks_map = {

		//---------------------------------------------------
		// TAGS
		//---------------------------------------------------
		"tags_pre_create_fun": async (p_tags_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				// passing the image_id to the gf_tagger control via this callback allows for
				// customization of the image_id fetching mechanism (whether its in the template,
				// or fetched via rest api, etc., or pulled from some internal browser/web DB).
				p_resolve_fun(p_image_id_str);
			});
			return p;
		},
		
		//---------------------------------------------------
		"tags_created_fun": (p_tags_lst)=>{

			console.log("added tags >>>>>>>>>>>", p_tags_lst);

			p_tags_lst.forEach(p_tag_str=>{

				tag_display(p_tag_str);
			});
		},

		//---------------------------------------------------
		// NOTES
		//---------------------------------------------------
		"notes_pre_create_fun": (p_notes_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				// passing the image_id to the gf_tagger control via this callback allows for
				// customization of the image_id fetching mechanism (whether its in the template,
				// or fetched via rest api, etc., or pulled from some internal browser/web DB).
				p_resolve_fun(p_image_id_str);
			});
			return p;
		},

		//---------------------------------------------------
		"notes_created_fun": (p_notes_lst)=>{

			console.log("added notes >>>>>>>>>>>", p_notes_lst)
		}

		//---------------------------------------------------
	}

	gf_tagger__init_ui_v2(p_image_id_str,
		obj_type_str,
		p_image_container_element,
		$("body"),
		callbacks_map,
		http_api_map,
		p_log_fun);

	//-------------------------------------------------
	function tag_display(p_tag_str) {

		/*
		check if the tags_container div exists, if not create it.
		the backend template has a div with class "tags_container" in the image container only if the image
		has tags. if it does not, the .tags_container div is not created
		*/
		if ($(p_image_container_element).find(".tags_container").length == 0) {
			$(p_image_container_element).append("<div class='tags_container'></div>");
		}

		$(p_image_container_element)
			.find(".tags_container")
			.append(`<a class='gf_image_tag' href='/v1/tags/objects?tag=${p_tag_str}&otype=image'>#${p_tag_str}</a>`)
	}

	//-------------------------------------------------
}

//-------------------------------------------------
export function masonry_layout_after_img_load(p_image_container) {
	
	const masonry = $('#gf_images_flow_container #items').data('masonry');
	masonry.once('layoutComplete', (p_event, p_laid_out_items)=>{
		$(p_image_container).css('visibility', 'visible');
	});
	
	// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
	$('#gf_images_flow_container #items').masonry();
	$('#gf_images_flow_container #items').masonry(<any>"reloadItems");
}

//-------------------------------------------------
export function init_image_date(p_image_element, p_log_fun) {

	const creation_time_element = $(p_image_element).find('.creation_time');
	const creation_time_f       = parseFloat($(creation_time_element).text());
	const creation_date         = new Date(creation_time_f*1000);

	const date_msg_str = $.timeago(creation_date);
	$(creation_time_element).text(date_msg_str);

	const creation_date__readable_str = creation_date.toDateString();
	const creation_date__readble      = $(`<div class="full_creation_date">${creation_date__readable_str}</div>`);

	$(creation_time_element).mouseover((p_e)=>{
		$(creation_time_element).append(creation_date__readble);

		// IMPORTANT!! - image size changed, so recalculate the Masonry layout.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container #items').masonry();
		$('#gf_images_flow_container #items').masonry(<any>'reloadItems');
	});

	$(creation_time_element).mouseout((p_e)=>{
		$(creation_date__readble).remove();

		// IMPORTANT!! - image size changed, so recalculate the Masonry layout.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container #items').masonry();
		$('#gf_images_flow_container #items').masonry(<any>'reloadItems');
		
	});
}