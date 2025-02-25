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

import * as gf_user_events   from "../../../../gf_events/ts/gf_user_events";
import * as gf_events        from "../gf_images_core/gf_events";
import * as gf_image_control from "../gf_images_core/gf_image_control";
import * as gf_viz_group     from "./gf_viz_group";
import * as gf_images_paging from "../gf_images_core/gf_images_paging";

// FIX!! - remove this from global scope!!
export var image_view_type_str = "small_view";

//---------------------------------------------------
export function get_current_view_type() {
	return image_view_type_str;
}

function init_tooltip(p_view_type_btn_element :HTMLElement) {

	$(p_view_type_btn_element).on('mouseenter', function() {
		$(p_view_type_btn_element).find('#tooltip').animate({
			opacity: 1
		}, 200);
	});
	$(p_view_type_btn_element).on('mouseleave', function() {
		$(p_view_type_btn_element).find('#tooltip').animate({
			opacity: 0
		}, 200);
	});
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

	const view_picker_element = $(`
		<div id='view_type_picker'>
			<div id='masonry_small_images'>
				<img src="https://gf-phoenix.s3.amazonaws.com/assets/gf_images_flows_browser/gf_small_view_icon.svg"></img>
				<div id="tooltip">small size view</div>
			</div>
			<div id='masonry_medium_images'>
				<img src="https://gf-phoenix.s3.amazonaws.com/assets/gf_images_flows_browser/gf_medium_view_icon.svg"></img>
				<div id="tooltip">medium size view</div>
			</div>

			<div id='viz_group_medium_images'>
				<img src="https://gf-phoenix.s3.amazonaws.com/assets/gf_images_flows_browser/gf_seek_view_icon.svg"></img>
				<div id="tooltip">seek view</div>
			</div>
		</div>`);
	$('body').append(view_picker_element);

	const small_view_btn  = $(view_picker_element).find('#masonry_small_images')[0];
	const medium_view_btn = $(view_picker_element).find('#masonry_medium_images')[0];
	const seek_view_btn   = $(view_picker_element).find('#viz_group_medium_images')[0];

	init_tooltip(small_view_btn);
	init_tooltip(medium_view_btn);
	init_tooltip(seek_view_btn);
	

	// SMALL_VIEW
	$(small_view_btn).on('click', function() {

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

	// MEDIUM_VIEW
	$(medium_view_btn).on('click', function() {

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

	// SEEK_VIEW
	$(seek_view_btn).on('click', function() {

		

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

			p_logged_in_bool,
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
	p_flow_pages_num_int   :number,
	p_initial_page_int     :number,
	p_logged_in_bool	   :boolean,
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
		const img__thumbnail_small_url_str  = p_element_map['thumbnail_small_url_str'];
		const img__thumbnail_medium_url_str = p_element_map['thumbnail_medium_url_str'];
		const img__thumbnail_large_url_str  = p_element_map['thumbnail_large_url_str'];
		const img__origin_page_url_str      = p_element_map['origin_page_url_str'];
		const img__title_str = p_element_map['title_str'];
		const img__tags_lst  = p_element_map['tags_lst'];
		const img_url_str = img__thumbnail_medium_url_str;

		const owner_user_name_str = "anon";



		const image_element = gf_image_control.create(image_id_str,
			img__format_str,
			img__creation_unix_time_f,
			img__origin_page_url_str,
			img__thumbnail_small_url_str,
			img__thumbnail_medium_url_str,
			img__thumbnail_large_url_str,
			img__title_str,
			img__tags_lst,
			owner_user_name_str,
			img__flows_names_lst,
			current_image_view_type_str,


			p_events_enabled_bool,
			p_host_str,

			// p_on_img_load_fun
			()=>{

			},

			// p_on_img_load_error_fun
			()=>{

			},
			p_plugin_callbacks_map,

			// p_on_viz_change_fun
			()=>{

			},
			p_log_fun);
			

        return image_element;
    }

	//-------------------------------------------------


	const container_element = $('#gf_images_flow_container')[0];
	const initial_page_int = 0;

    const props :gf_viz_group.GF_props = {

		flow_name_str:    p_flow_name_str,
		container:        container_element,

        start_page_int:   initial_page_int,
        end_page_int:     p_flow_pages_num_int,
        initial_page_int: p_initial_page_int,

		image_view_type_str: current_image_view_type_str,
		logged_in_bool:      p_logged_in_bool,     
    	events_enabled_bool: p_events_enabled_bool,

		plugin_callbacks_map: p_plugin_callbacks_map
    };

	const seeker__container_element = gf_viz_group.init(initial_elements_lst,
		props,
        element_create_fun,
		p_log_fun,
		
		// the container already contains elements that are created and attached
		// to the container, so we dont want to create any initially (only if paging is done).
		false); // p_create_initial_elements_bool




	// INIT_PAGING
	gf_images_paging.init(initial_page_int,
		p_flow_name_str,
		current_image_view_type_str,
		p_logged_in_bool,
		p_events_enabled_bool,
		p_plugin_callbacks_map,

		// p_on_page_load_fun
		(p_new_page_int :number)=>{

			$("#gf_images_flow_container").data("current_page", p_new_page_int);

			// $(current_pages_display).find('#end_page').text(p_new_page_int);
		},
		p_log_fun);
}