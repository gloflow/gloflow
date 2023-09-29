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

import * as gf_image_viewer from "./../../../../gf_core/ts/gf_image_viewer";
import * as gf_gifs_viewer  from "./../../../../gf_core/ts/gf_gifs_viewer";

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
export function init_image_element(p_image_id_str :string,
	p_img__format_str               :string,
	p_img__creation_unix_time_f     :string,
	p_img__origin_page_url_str      :string,
	p_img__thumbnail_small_url_str  :string,
	p_img__thumbnail_medium_url_str :string,
	p_img__thumbnail_large_url_str  :string,
	p_img__tags_lst                 :string[],
	p_img__owner_user_name_str      :string,
	p_flows_names_lst               :string[],
	p_current_image_view_type_str   :string,

	p_on_img_load_fun,
	p_on_img_load_error_fun,
	p_log_fun) {

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
				p_log_fun);
		}

		//------------------

		p_on_img_load_fun(image_container);
	});

	// IMAGE_FAILED_TO_LOAD
	$(image_container).find('img').on('error', function() {

		p_log_fun("ERROR", "IMAGE_FAILED_TO_LOAD ----------");
		p_on_img_load_error_fun();
	});

	//------------------
	init_image_date(image_container, p_log_fun);

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

	return image_container;
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