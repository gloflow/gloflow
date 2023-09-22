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

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_tagger__init_ui;
declare var gf_tagger__http_add_tags_to_obj;

//-------------------------------------------------
export function init(p_logged_in_bool,
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


		// FIX!! - this function has been moved to gf_core/gf_images_viewer.ts, as a general viewer,
		//         to use universaly gf_images/gf_landing_page.
		//         gf_images flows_browser already uses the version from gf_core.
		init_image_viewer(p_image_info_element, p_log_fun);

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
			
			const gf_container_element = $("body");
				
			init_tagging(p_image_info_element,
				gf_container_element,
		
				// tags_create_pre_fun
				// called before a tag is about to be added to an image
				async (p_tags_lst)=>{
					const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
						
						/*
						IMPORTANT!! - img_system_id is attached as a data property to the image container element
							in the server template rendering.
						*/
						var img_system_id_str = $(p_image_info_element).attr("data-img_system_id_str");
						p_resolve_fun(img_system_id_str);

					});
					return p;
				},
				p_gf_host_str,
				p_log_fun);
		}
	}

	//-------------------------------------------------
}

//---------------------------------------------------
// TAGGING_UI

function init_tagging(p_image_container_element,
	p_gf_container,
	p_tags_create_pre_fun,
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
	const input_element_parent_selector_str = "#page_info_container";

	gf_tagger__init_ui(obj_type_str,
		p_image_container_element,
		input_element_parent_selector_str,

		//---------------------------------------------------
		// tags_create_pre_fun
		async (p_tags_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				// p_tags_create_pre_fun resolves the system_id of the item being tagged
				image_system_id_str = await p_tags_create_pre_fun(p_tags_lst);

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
		()=>{}, // on_tag_ui_add_fun
		()=>{}, // on_tag_ui_remove_fun
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

//-------------------------------------------------
// DEPRECATED!!
// REMOVE!!
function init_image_viewer(p_image_element, p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_landing_page.init_image_viewer()');

	const img_thumb_medium_url = $(p_image_element).find('img').data('img_thumb_medium_url');
	const image_view           = $(`
		<div id="image_view">
			<div id="background"></div>
			<div id="image_detail">
				<img src="${img_thumb_medium_url}"></img>
			</div>
		</div>`);

	$(p_image_element).find('img').click(()=>{
		
		console.log(img_thumb_medium_url)
		$('body').append(image_view);

		//----------------------
		// BAKCGROUND
		const bg = $(image_view).find('#background');

		// position the background vertically where the user has scrolled to
		$(bg).css('top', $(window).scrollTop()+'px');

		//----------------------
		// IMPORTANT!! - turn off vertical scrolling while viewing the image
		$("body").css("overflow-y", "hidden");

		//----------------------
		const image_detail = $(image_view).find('#image_detail');
		$(image_detail).css("position", "absolute");
	    $(image_detail).css("top",  Math.max(0, (($(window).height() - $(image_detail).outerHeight()) / 2) + $(window).scrollTop())  + "px");
	    $(image_detail).css("left", Math.max(0, (($(window).width()  - $(image_detail).outerWidth()) / 2)  + $(window).scrollLeft()) + "px");
	    
		//----------------------
	    $(bg).click(()=>{
	    	$(image_view).remove();
	    	$("body").css("overflow-y", "auto"); // turn vertical scrolling back on when done viewing the image
	    });

	    //----------------------
	});
}

//-------------------------------------------------
/*function init_image_date(p_target_element, p_log_fun) {

	const creation_time_element = $(p_target_element).find('.creation_time');
	const creation_time_f       = parseFloat($(creation_time_element).text());
	const creation_date         = new Date(creation_time_f*1000);

	const date_msg_str = $.timeago(creation_date);
	$(creation_time_element).text(date_msg_str);

	const creation_date__readable_str = creation_date.toDateString();
	const creation_date__readble      = $('<div class="full_creation_date">'+creation_date__readable_str+'</div>');

	$(creation_time_element).mouseover((p_e)=>{
		$(creation_time_element).append(creation_date__readble);
	});

	$(creation_time_element).mouseout((p_e)=>{
		$(creation_date__readble).remove();
	});
}*/