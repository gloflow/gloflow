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

// ///<reference path="../../../../d/jquery.d.ts" />

import * as gf_user_events  from "./../../../../gf_events/ts/gf_user_events";
import * as gf_viz_group    from "./../../../../gf_controls/gf_viz_group/ts/gf_viz_group";
import * as gf_gifs_viewer  from "./../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_image_viewer from "./gf_image_viewer";
import * as gf_images_http  from "./gf_images_http";
import * as gf_paging       from "./gf_images_paging";
import * as gf_events       from "./gf_events";

// FIX!! - remove this from global scope!!
export var image_view_type_str = "small_view";

//---------------------------------------------------
export function get_current_view_type() {
	return image_view_type_str;
}

//---------------------------------------------------
// view_type_picker - picks the type of view that is used to display images in a flow.
//                    default is masonry with 6 columns.
export function init(p_flow_name_str :string,
	p_logged_in_bool       :boolean,
	p_events_enabled_bool  :boolean,
	p_host_str             :string,
	p_plugin_callbacks_map :any,
	p_log_fun :any) {

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


		//---------------------
		// EVENTS
		if (p_events_enabled_bool && p_logged_in_bool) {
			
			const event_meta_map = {
				"view_type": "masonry_small"
			};
			gf_user_events.send_event_http(gf_events.GF_IMAGES_VIEW_TYPE_PICKER_ACTIVATE_VIEW,
				"browser",
				event_meta_map,
				p_host_str)
		}
		
		//---------------------
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

		//---------------------
		// EVENTS
		if (p_events_enabled_bool && p_logged_in_bool) {
			
			const event_meta_map = {
				"view_type": "masonry_medium"
			};
			gf_user_events.send_event_http(gf_events.GF_IMAGES_VIEW_TYPE_PICKER_ACTIVATE_VIEW,
				"browser",
				event_meta_map,
				p_host_str)
		}
		
		//---------------------
	});

	// VIZ_GROUP
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

			p_events_enabled_bool,
			p_host_str,
			p_plugin_callbacks_map,
			p_log_fun);

		//---------------------

		//---------------------
		// EVENTS
		if (p_events_enabled_bool && p_logged_in_bool) {
			
			const event_meta_map = {
				"view_type": "viz_group_medium"
			};
			gf_user_events.send_event_http(gf_events.GF_IMAGES_VIEW_TYPE_PICKER_ACTIVATE_VIEW,
				"browser",
				event_meta_map,
				p_host_str)
		}
		
		//---------------------
	});
}

//---------------------------------------------------
function init__viz_group_view(p_flow_name_str :string,
	p_flow_pages_num_int :number,
	p_initial_page_int   :number,

	p_events_enabled_bool  :boolean,
	p_host_str             :string,
	p_plugin_callbacks_map :any,
	p_log_fun              :any) {

	const current_image_view_type_str = "viz_group_medium_view";
	const initial_elements_lst :any = [];

	//-------------------------------------------------
	// ELEMENT_CREATE
    function element_create_fun(p_element_map :any) {

		const image_id_str                  = p_element_map['id_str'];
		const img__format_str               = p_element_map['format_str'];
		const img__creation_unix_time_f     = p_element_map['creation_unix_time_f'];
		const img__flows_names_lst          = p_element_map["flows_names_lst"];
		const img__thumbnail_medium_url_str = p_element_map['thumbnail_medium_url_str'];
		const img__thumbnail_large_url_str  = p_element_map['thumbnail_large_url_str'];
		const img__origin_page_url_str      = p_element_map['origin_page_url_str'];

		const img_url_str = img__thumbnail_medium_url_str;

		// IMPORTANT!! - '.gf_image' is initially invisible, and is faded into view when its image is fully loaded
		//               and its positioned appropriatelly in the Masonry grid
		const image_container = $(`
			<div class="gf_image item ${current_image_view_type_str}"

				data-img_id="${image_id_str}"
				data-img_format="${img__format_str}"
				data-img_flows_names="${img__flows_names_lst.join(' ')}"
				
				style='visibility:hidden;'>
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
			gf_gifs_viewer.init(image_container, image_id_str, img__flows_names_lst, p_log_fun);

		} else {

			gf_image_viewer.init(image_container,
				image_id_str,
				img__thumbnail_medium_url_str,
				img__thumbnail_large_url_str,
				img__flows_names_lst,
				p_events_enabled_bool,
				p_host_str,
				p_plugin_callbacks_map,
				p_log_fun);
		}

		//------------------
		
        return image_container;
    }

    //-------------------------------------------------
	// ELEMENTS_PAGE_GET
    function elements_page_get_fun(p_new_page_number_int: number) {
        return new Promise(async function(p_resolve_fun, p_reject_fun) {

            

			// HTTP_LOAD_NEW_PAGE

			const resp_map = await gf_images_http.get_page(p_flow_name_str,
				p_new_page_number_int,
				p_log_fun);
			const page_elements_lst = resp_map["pages_lst"][0];

			p_resolve_fun(page_elements_lst);
        });
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

	const viz_props :gf_viz_group.GF_viz_props = {
        seeker_container_height_px: $(window).height(), // 500,
        seeker_container_width_px:  100,
        seeker_bar_width_px:        50, 
        seeker_range_bar_width:     30,
        seeker_range_bar_height:    500,
        seeker_range_bar_color_str: "red",
        assets_uris_map:            assets_uris_map,
    }


    const props :gf_viz_group.GF_props = {

		// IDs
		container_id_str:        id_str, 
		parent_container_id_str: parent_id_str, 

        start_page_int:   0,
        end_page_int:     p_flow_pages_num_int,
        initial_page_int: p_initial_page_int,
        assets_uris_map:  assets_uris_map,
        viz_props:        viz_props,
    };

	const seeker__container_element = gf_viz_group.init(initial_elements_lst,
		props,
        element_create_fun,
        elements_page_get_fun,

		// the container already contains elements that are created and attached
		// to the container, so we dont want to create any initially (only if paging is done).
		false); // p_create_initial_elements_bool

	// seeker should be in fixed position, as the user scrolls the images themselves
	$(seeker__container_element).css("position", "fixed");
}