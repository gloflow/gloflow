/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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
///<reference path="../../../../d/masonry.layout.d.ts" />
///<reference path="../../../../d/jquery.timeago.d.ts" />

import * as gf_core_utils    from "./../../../../gf_core/ts/gf_utils";
import * as gf_identity_http from "./../../../../gf_identity/ts/gf_identity_http";
import * as gf_image_control from "./../../../gf_images/ts/gf_images_core/gf_image_control";
import * as gf_images_paging from "./../../../gf_images/ts/gf_images_core/gf_images_paging";
import * as gf_sys_panel     from "./../../../../gf_sys_panel/ts/gf_sys_panel";
import * as gf_flows_picker  from "./../../../gf_images/ts/gf_images_flows_browser/gf_flows_picker";
import * as gf_tags_picker   from "./../../../gf_tagger/ts/gf_tags_picker/gf_tags_picker";
// import * as gf_image_viewer  from "./../../../../gf_core/ts/gf_image_viewer";
// import * as gf_utils         from "./../../../gf_images/ts/gf_images_core/gf_utils";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
// declare var gf_upload__init;

//-------------------------------------------------
export async function init(p_tag_str :string,
	p_plugin_callbacks_map,
	p_log_fun) {

	const events_enabled_bool = true;
	const current_host_str = gf_core_utils.get_current_host();
	
	// inspect if user is logged-in or not
	const urls_map          = gf_identity_http.get_standard_http_urls(current_host_str);
	const auth_http_api_map = gf_identity_http.get_http_api(urls_map, current_host_str);
	const logged_in_bool = await auth_http_api_map["general"]["logged_in"]();

	
	//---------------------
	// SYS_PANEL
    gf_sys_panel.init_with_auth(p_log_fun);

	//---------------------
	// FLOWS_PICKER - display it if the user is logged in
	if (logged_in_bool) {

		gf_flows_picker.init(events_enabled_bool,
			p_plugin_callbacks_map,
			current_host_str,
			p_log_fun)
	}

	// TAGS_PICKER - display it if the user is logged in
	if (logged_in_bool) {

		gf_tags_picker.init(p_log_fun)
	}

    //---------------------
	// IMAGES
    init_images(logged_in_bool,
		current_host_str,
		events_enabled_bool,
		p_plugin_callbacks_map,
		p_log_fun);

	//------------------
	// LOAD_PAGES_ON_SCROLL

	var current_page_int = 6; // the few initial pages are already statically embedded in the document
	$("#gf_images_flow_container").data("current_page", current_page_int); // used in other functions to inspect current page

	var page_is_loading_bool = false;

	window.onscroll = async ()=>{

		// $(document).height() - height of the HTML document
		// window.innerHeight   - Height (in pixels) of the browser window viewport including, if rendered, the horizontal scrollbar
		if (window.scrollY >= $(document).height() - (window.innerHeight+50)) {
			
			// IMPORTANT!! - only load 1 page at a time
			if (!page_is_loading_bool) {
				
				page_is_loading_bool = true;
				p_log_fun("INFO", `current_page_int - ${current_page_int}`);

				var image_view_type_str = "masonry_small_images";
				

				const page_source_ref_str  = p_tag_str;
				const page_source_type_str = "tag"

				await gf_images_paging.load_new_page(page_source_ref_str,
					page_source_type_str,
					current_page_int,
					image_view_type_str,
					logged_in_bool,
					events_enabled_bool,
					p_plugin_callbacks_map,
					p_log_fun);
				

				current_page_int += 1;
				$("#gf_images_flow_container").data("current_page", current_page_int);

				page_is_loading_bool = false;
			}
		}
	};

	//------------------
}

//-------------------------------------------------
function init_images(p_logged_in_bool :boolean,
	p_gf_host_str :string,
	p_events_enabled_bool :boolean,
	p_plugin_callbacks_map :any,
	p_log_fun :any) {

	$('#images_container').masonry({
		itemSelector: '.gf_image',
		columnWidth:  6
	});
	
	/*
	IMPORTANT!! - as each image loads call masonry to reconfigure the view.
		this is necessary so that initial images in the page, before
		load_new_page() starts getting called, are properly laid out
		by masonry.
	*/
	$('.gf_image img').on('load', ()=>{
		$('#images_container').masonry();
	});

    $('#images_container .gf_image').each((p_i, p_e)=>{

		const image_element = p_e;
		
		/*
		const image_id_str = $(image_element).data('img_id');
		const image_flows_names_lst = $(image_element).data('img_flows_names').split(",");
		const img_thumb_medium_url_str = $(image_element).find('img').data('img_thumb_medium_url');
		const img_thumb_large_url_str  = $(image_element).find('img').data('img_thumb_large_url');
		const img_format_str           = $(image_element).data('img_format');
		const origin_page_url_link = $(image_element).find(".origin_page_url a")[0];
		*/
		
		const image_flows_names_lst = $(image_element).data('img_flows_names').split(",");

		// IMAGE_CONTROL
		gf_image_control.init_existing_dom(image_element,
			image_flows_names_lst,

			p_gf_host_str,
			p_logged_in_bool,
			p_events_enabled_bool,
			p_plugin_callbacks_map,
			p_log_fun);

		/*	
		gf_utils.init_image_date(image_element, p_log_fun);


		//----------------
		// CLEANUP - for images that dont come from some origin page (direct uploads, or generated images)
		//           this origin_page_url is set to empty string. check for that and remove it.
		// FIX!! - potentially on the server/template-generation side this div node shouldnt get included
		//         at all for images that dont have an origin_page_url.
		if ($(origin_page_url_link).text().trim() == "") {
			$(image_element).find(".origin_page_url").remove();
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
		else {

			gf_image_viewer.init(image_element,
				image_id_str,
				img_thumb_medium_url_str,
				img_thumb_large_url_str,
				image_flows_names_lst,
				p_log_fun);
		}

		//----------------
		*/
	});
}