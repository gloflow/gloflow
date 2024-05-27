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

///<reference path="../../../../d/jquery.d.ts" />

import * as gf_gifs_viewer  from "./../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_image_colors from "./../../../../gf_core/ts/gf_image_colors";
import * as gf_color        from "./../../../../gf_core/ts/gf_color";
import * as gf_image_viewer from "./gf_image_viewer";
import * as gf_images_share from "./gf_images_share";
import * as gf_utils from "./gf_utils";

//-------------------------------------------------
// CREATE
/*
create the image control DOM element from scratch.
initialize it. does not append it to the parent DOM tree,
instead just returns it.
*/

export function create(p_image_id_str :string,
	p_img__format_str               :string,
	p_img__creation_unix_time_f     :string,
	p_img__origin_page_url_str      :string,
	p_img__thumbnail_small_url_str  :string,
	p_img__thumbnail_medium_url_str :string,
	p_img__thumbnail_large_url_str  :string,
	p_img__title_str                :string,
	p_img__tags_lst                 :string[],
	p_img__owner_user_name_str      :string,
	p_flows_names_lst               :string[],
	p_current_image_view_type_str   :string,

	p_events_enabled_bool	   :boolean,
	p_host_str				   :string,
	p_on_img_load_fun :any,
	p_on_img_load_error_fun :any,
	p_log_fun :any) {

	var img_url_str;
	switch (p_current_image_view_type_str) {
		case "small_view":
			img_url_str = p_img__thumbnail_medium_url_str;
			break;
		case "medium_view":
			img_url_str = p_img__thumbnail_medium_url_str;
			break;
	}

	// IMPORTANT!! - '.gf_image' is initially invisible, and is faded into view when its image is fully loaded
	//               and its positioned appropriatelly in the Masonry grid
	const image_container = $(`
		<div class="gf_image item ${p_current_image_view_type_str}"
			data-img_id="${p_image_id_str}"
			data-img_format="${p_img__format_str}"
			data-img_flows_names="${p_flows_names_lst.join(',')}"
			style='visibility:hidden;'>

			<div class="image_title">${p_img__title_str}</div>

			<img src="${img_url_str}" data-img_thumb_medium_url="${p_img__thumbnail_medium_url_str}"></img>
			
			<div class="tags_container"></div>
			
			<div class="origin_page_url">
				<a href="${p_img__origin_page_url_str}" target="_blank">${p_img__origin_page_url_str}</a>
			</div>
			<div class="creation_time">${p_img__creation_unix_time_f}</div>
			<div class="owner_user_name">by <span>${p_img__owner_user_name_str}</span></div>
		</div>`);

	//------------------
	
	// FIX!! - this needs to happen after the image <div> is added to the DOM, 
	//         here reloading masonry layout doesnt have the intended effect, since 
	//         the image hasnt been added yet.
	//         move it to be after $("#gf_images_flow_container").append(image);

	$(image_container).find('img').on('load', ()=>{

		//------------------
		// MASONRY_RELOAD
		// var masonry = $('#gf_images_flow_container #items').data('masonry');
		// masonry.once('layoutComplete', (p_event, p_laid_out_items)=>{
		// 	$(image_container).css('visibility', 'visible');
		// });
		
		
		// // IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		// $('#gf_images_flow_container #items').masonry();
		// $('#gf_images_flow_container #items').masonry(<any>"reloadItems");

		//------------------

		// CLEANUP - for images that dont come from some origin page (direct uploads, or generated images)
		//           this origin_page_url is set to empty string. check for that and remove it.
		// FIX!! - potentially on the server/template-generation side this div node shouldnt get included
		//         at all for images that dont have an origin_page_url.
		if ($(image_container).find(".origin_page_url a").text().trim() == "") {
			$(image_container).find(".origin_page_url").remove();
		}
		
		//------------------
		// VIEWER_INIT

		if (p_img__format_str == 'gif') {
			gf_gifs_viewer.init(image_container, p_image_id_str, p_flows_names_lst, p_log_fun);
		} else {
			
			gf_image_viewer.init(image_container,
				p_image_id_str,
				p_img__thumbnail_medium_url_str,
				p_img__thumbnail_large_url_str,
				p_flows_names_lst,
				p_events_enabled_bool,
				p_host_str,
				p_log_fun);
		}

		//------------------
		// IMAGE_PALLETE

		init_pallete(image_container);

		//------------------
		
		p_on_img_load_fun(image_container);
	});

	// IMAGE_FAILED_TO_LOAD
	$(image_container).find('img').on('error', function() {

		p_log_fun("ERROR", "IMAGE_FAILED_TO_LOAD ----------");
		p_on_img_load_error_fun();
	});

	//------------------
	gf_utils.init_image_date(image_container, p_log_fun);

	//------------------
	// TAGS
	if (p_img__tags_lst != null && p_img__tags_lst.length > 0) {
		$.each(p_img__tags_lst, function(p_i, p_tag_str) {
			const tag = $(
				`<a class='gf_image_tag' href='/v1/tags/objects?tag=${p_tag_str}&otype=image'>
					${p_tag_str}
				</a>`);

			$(image_container).find('.tags_container').append(tag);
		});
	}

	//------------------

	//---------------------------------------------------
	function init_pallete(p_image_info_element :any) {

		console.log("init_pallete...");

		const img = $(p_image_info_element).find("img")[0];

		const assets_paths_map = {
			"copy_to_clipboard_btn": "/images/static/assets/gf_copy_to_clipboard_btn.svg",
		}
		gf_image_colors.init_pallete(img,
			assets_paths_map,
			(p_color_dominant_hex_str :string,
			p_colors_hexes_lst :string[])=>{

				console.log("init_pallete post callback...");

				console.log("p_color_dominant_hex_str: ", p_color_dominant_hex_str);
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
						// if background is light, then the text should be dark, so setting it here explicitly
						// on dominant color classification.
						
						$(p_image_info_element).find(".image_title").css("color", "black");
						$(p_image_info_element).find(".origin_page_url a").css("color", "black");
						$(p_image_info_element).find(".creation_time").css("color", "black");
						$(p_image_info_element).find(".owner_user_name").css("color", "black");

						break;

					// DARK
					case "dark":
						// css rules external to this function set the default color of
						// text to white, so dark background dominant-color works fine.
						// no need to set anything here yet.
						break;
				};
			});
	}

	//---------------------------------------------------
	return image_container;
}

//---------------------------------------------------
// INIT_EXISTING_DOM
/*
used for templates usually, where the image element DOM structure is already
created server side when loaded into the browser, and just needs to be initialized
(no creation of the DOM tree for the image control)
*/

export function init_existing_dom(p_image_element :any,
	p_flows_names_lst :string[],
	p_gf_host_str     :string,
	p_logged_in_bool  :boolean,

	p_events_enabled_bool :boolean,
	p_plugin_callbacks_map :any,
	p_log_fun :any) {



    gf_utils.init_image_date(p_image_element, p_log_fun);

	const image_id_str = $(p_image_element).data('img_id');
	const img_thumb_medium_url_str = $(p_image_element).find('img').data('img_thumb_medium_url');
	const img_thumb_large_url_str  = $(p_image_element).find('img').data('img_thumb_large_url');
	const img_format_str           = $(p_image_element).data('img_format');
	const flows_names_lst          = $(p_image_element).data('img_flows_names').split(",");

	const origin_page_url_link = $(p_image_element).find(".origin_page_url a")[0];

	//----------------
	// CLEANUP - for images that dont come from some origin page (direct uploads, or generated images)
	//           this origin_page_url is set to empty string. check for that and remove it.
	// FIX!! - potentially on the server/template-generation side this div node shouldnt get included
	//         at all for images that dont have an origin_page_url.
	if ($(origin_page_url_link).text().trim() == "") {
		$(p_image_element).find(".origin_page_url").remove();
	}

	//----------------
	// LINK_TEXT_SHORTEN - if the link text (not its href) is too long, dont display it completely in the UI
	//                     because it clutters the UI too much.
	//                     instead shorten it at its cutoff length and append "..."
	const link_text_cutoff_threshold_int = 50;
	if ($(origin_page_url_link).text().length > link_text_cutoff_threshold_int) {
		const old_link_text_str = $(origin_page_url_link).text();
		const new_link_text_str = `${old_link_text_str.slice(0, link_text_cutoff_threshold_int)}...`;
		$(origin_page_url_link).text(new_link_text_str);
	}

	//----------------
	// GIFS
	if (img_format_str == 'gif') {			
		gf_gifs_viewer.init(p_image_element, image_id_str, p_flows_names_lst, p_log_fun);
	}

	//----------------
	else {
		gf_image_viewer.init(p_image_element,
			image_id_str,
			img_thumb_medium_url_str,
			img_thumb_large_url_str,
			flows_names_lst,
			p_events_enabled_bool,
			p_gf_host_str,
			p_log_fun);
	}

	//----------------
	// LOGGED_IN - only initialize this part if the user is authenticated
	
	if (p_logged_in_bool) {
			
		// TAGGING
		gf_utils.init_tagging(image_id_str,
			p_image_element,
			p_gf_host_str,
			p_log_fun);

		// SHARE
		gf_images_share.init(image_id_str,
			p_image_element,
			p_plugin_callbacks_map,
			p_log_fun);
	}

	//----------------




}