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

import * as gf_viz_group_paged         from "./../../../../gf_controls/gf_viz_group_paged/ts/gf_viz_group_paged";
import * as gf_viz_group_random_access from "./../../../../gf_controls/gf_viz_group_paged/ts/gf_viz_group_random_access";
import * as gf_gifs_viewer             from "./../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_image_viewer            from "./../../../../gf_core/ts/gf_image_viewer";
import * as gf_paging                  from "gf_paging";

// FIX!! - remove this from global scope!!
export var image_view_type_str = "small_view";

//---------------------------------------------------
// view_type_picker - picks the type of view that is used to display images in a flow.
//                    default is masonry with 6 columns.
export function init(p_flow_name_str :string,
	p_log_fun) {

	const container = $(`
		<div id='view_type_picker'>
			<div id='masonry_small_images'>
			</div>
			<div id='masonry_medium_images'>
			</div>

			<div id='viz_group_medium_images'>
			</div>
		</div>`);
	$('body').append(container);

	// MASONRY_SMALL_IMAGES
	$(container).find('#masonry_small_images').on('click', function() {

		// FIX!! - global var. handle this differently;
		image_view_type_str = "small_view";

		$(".gf_image").each(function(p_i, p_e) {

			const img_e                   = $(p_e).find('img');
			const img_thumb_small_url_str = $(img_e).data('img_thumb_small_url');
			$(img_e).attr("src", img_thumb_small_url_str);

			// switch gf_image class to "small_view"
			$(p_e).removeClass("medium_view");
			$(p_e).addClass("small_view");
		});

		// dimensions of items changed, re-layout masonry.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container #items').masonry();
		$('#gf_images_flow_container #items').masonry(<any>"reloadItems");
	});

	// MASONRY_MEDIUM_IMAGES
	$(container).find('#masonry_medium_images').on('click', function() {

		// FIX!! - global var. handle this differently;
		image_view_type_str = "medium_view";

		$(".gf_image").each(function(p_i, p_e) {

			const img_e                    = $(p_e).find('img');
			const img_thumb_medium_url_str = $(img_e).data('img_thumb_medium_url');
			$(img_e).attr("src", img_thumb_medium_url_str);

			// switch gf_image class to "medium_view"
			$(p_e).removeClass("small_view");
			$(p_e).addClass("medium_view");
		});

		// dimensions of items changed, re-layout masonry.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container #items').masonry();
		$('#gf_images_flow_container #items').masonry(<any>"reloadItems");
	});

	// MASONRY_MEDIUM_IMAGES
	$(container).find('#viz_group_medium_images').on('click', function() {

		// FIX!! - global var. handle this differently;
		image_view_type_str = "viz_group_medium_view";

		$(".gf_image").each(function(p_i, p_e) {

			const img_e                    = $(p_e).find('img');
			const img_thumb_medium_url_str = $(img_e).data('img_thumb_medium_url');
			$(img_e).attr("src", img_thumb_medium_url_str);

			// switch gf_image class to "medium_view"
			$(p_e).removeClass("small_view");
			$(p_e).addClass("medium_view");
		});

		// switching to viz_group view, so remove masonry
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container #items').masonry(<any>"destroy");

		//------------------
		// VIZ_GROUP
		const flow_pages_num_int = $("#gf_images_flow_container").data("flow_pages_num");
		const current_page_int   = $("#gf_images_flow_container").data("current_page");
		init__viz_group_view(p_flow_name_str,
			flow_pages_num_int,
			current_page_int,
			p_log_fun);

		//------------------
	});
}

//---------------------------------------------------
function init__viz_group_view(p_flow_name_str :string,
	p_flow_pages_num_int :number,
	p_initial_page_int   :number,
	p_log_fun) {

	const current_image_view_type_str = "viz_group_medium_view";
	const initial_elements_lst = [];
	//-------------------------------------------------
	// ELEMENT_CREATE
    function element_create_fun(p_element_map) {

		const img__id_str                   = p_element_map['id_str'];
		const img__format_str               = p_element_map['format_str'];
		const img__creation_unix_time_f     = p_element_map['creation_unix_time_f'];
		const img__thumbnail_medium_url_str = p_element_map['thumbnail_medium_url_str'];
		const img__origin_page_url_str      = p_element_map['origin_page_url_str'];

		const img_url_str = img__thumbnail_medium_url_str;

		// IMPORTANT!! - '.gf_image' is initially invisible, and is faded into view when its image is fully loaded
		//               and its positioned appropriatelly in the Masonry grid
		const image_container = $(`
			<div class="gf_image item ${current_image_view_type_str}" data-img_id="${img__id_str}" data-img_format="${img__format_str}" style='visibility:hidden;'>
				<img src="${img_url_str}" data-img_thumb_medium_url="${img__thumbnail_medium_url_str}"></img>
				<div class="tags_container"></div>
				<div class="origin_page_url">
					<a href="${img__origin_page_url_str}" target="_blank">${img__origin_page_url_str}</a>
				</div>
				<div class="creation_time">${img__creation_unix_time_f}</div>
			</div>`);
			
		//------------------
		// VIEWER_INIT

		if (img__format_str == 'gif') {
			gf_gifs_viewer.init(image_container, img__id_str, p_flow_name_str, p_log_fun);
		} else {
			gf_image_viewer.init(image_container, img__thumbnail_medium_url_str, p_flow_name_str, p_log_fun);
		}

		//------------------
		
        return image_container;
    }

    //-------------------------------------------------
	// ELEMENTS_PAGE_GET
    function elements_page_get_fun(p_new_page_number_int: number) {
        const p = new Promise(function(p_resolve_fun, p_reject_fun) {

            const page_elements_lst = [];
            p_resolve_fun(page_elements_lst);

			// HTTP_LOAD_NEW_PAGE
			gf_paging.http__load_new_page(p_flow_name_str,
				p_new_page_number_int,

				// p_on_complete_fun
				(p_page_elements_lst)=>{
					p_resolve_fun(p_page_elements_lst);
				},
				(p_error)=>{
					p_reject_fun();
				},
				p_log_fun);
        });
        return p;
    }

	//-------------------------------------------------


	// IMPORTANT!! - already existing div element
	const id_str = "gf_images_flow_container";

	// this is empty because gf_viz_group wont append to parent itself,
	// the container div is already present in the DOM
    const parent_id_str = "";
    
	const assets_uris_map = {
        "gf_bar_handle_btn": "https://gloflow.com/images/static/assets/gf_bar_handle_btn.svg",
    };

	const random_access_viz_props :gf_viz_group_random_access.GF_random_access_viz_props = {
        seeker_container_height_px: $(window).height(), // 500,
        seeker_container_width_px:  100,
        seeker_bar_width_px:        50, 
        seeker_range_bar_width:     30,
        seeker_range_bar_height:    500,
        seeker_range_bar_color_str: "red",
        assets_uris_map:            assets_uris_map,
    }


    const props :gf_viz_group_paged.GF_props = {
        start_page_int:   0,
        end_page_int:     p_flow_pages_num_int,
        initial_page_int: p_initial_page_int,
        assets_uris_map:  assets_uris_map,
        random_access_viz_props: random_access_viz_props,
    };

	const seeker__container_element = gf_viz_group_paged.init(id_str,
        parent_id_str,
        initial_elements_lst,
		props,
        element_create_fun,
        elements_page_get_fun,

		// the container already contains elements that are created and attached
		// to the container, so we dont want to create any initially (only if paging is done).
		false); // p_create_initial_elements_bool

	// seeker should be in fixed position, as the user scrolls the images themselves
	$(seeker__container_element).css("position", "fixed");
}