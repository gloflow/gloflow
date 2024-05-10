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

import * as gf_time             from "./../../../../gf_core/ts/gf_time";

import * as gf_core_utils       from "./../../../../gf_core/ts/gf_utils";
import * as gf_sys_panel        from "./../../../../gf_sys_panel/ts/gf_sys_panel";
import * as gf_identity         from "./../../../../gf_identity/ts/gf_identity";
import * as gf_identity_http    from "./../../../../gf_identity/ts/gf_identity_http";
import * as gf_images_http      from "./../gf_images_core/gf_images_http";
import * as gf_image_control    from "./../gf_images_core/gf_image_control";
import * as gf_images_paging    from "../gf_images_core/gf_images_paging";
import * as gf_view_type_picker from "./../gf_images_core/gf_view_type_picker";
import * as gf_utils            from "../gf_images_core/gf_utils";
import * as gf_flows_picker     from "./gf_flows_picker";

//-------------------------------------------------
declare var URLSearchParams;

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;

//-------------------------------------------------
export async function init(p_plugin_callbacks_map,
	p_log_fun) {

	const current_host_str = gf_core_utils.get_current_host();
	
	//---------------------
	// META
	const notifications_meta_map = {
		"login_first_stage_success": "login success"
	};

	//-------------------------------------------------
	function get_current_flow() {
		const url_params       = new URLSearchParams(window.location.search);
		const qs_flow_name_str = url_params.get('fname');
		var   flow_name_str;

		if (qs_flow_name_str == null) {
			flow_name_str = 'general'; // default value
		} else {
			flow_name_str = qs_flow_name_str;
		}
		return flow_name_str;
	}
	
	//-------------------------------------------------

	//-----------------
	gf_sys_panel.init_with_auth(p_log_fun);
	gf_flows_picker.init(p_log_fun);
	
	//---------------------
	// IDENTITY
	// first complete main initialization and only then initialize gf_identity
	const urls_map          = gf_identity_http.get_standard_http_urls(current_host_str);
	const auth_http_api_map = gf_identity_http.get_http_api(urls_map);
	gf_identity.init_with_http(notifications_meta_map, urls_map);
	

	
	const parent_node = $("#right_section");
	const home_url_str = urls_map["home"];

	gf_identity.init_me_control(parent_node,
		auth_http_api_map,
		home_url_str);
	
	// inspect if user is logged-in or not
	const logged_in_bool = await auth_http_api_map["general"]["logged_in"]();

	const flow_name_str = get_current_flow();

	//---------------------
	// UPLOAD

	init_upload(flow_name_str, p_log_fun);

	//---------------------
	// MASONRY
	$('#gf_images_flow_container #items').masonry({
		// options...
		itemSelector: '.item',
		columnWidth:  6
	});

	/*
	IMPORTANT!! - as each image loads call masonry to reconfigure the view.
		this is necessary so that initial images in the page, before
		load_new_page() starts getting called, are properly laid out
		by masonry.
	*/
	$('.gf_image img').on('load', ()=>{
		$('#gf_images_flow_container #items').masonry();
	});

	//---------------------
	// IMAGE_CONTROLS

	$('.gf_image').each((p_i, p_e)=>{

		const image_element = p_e;

		// IMAGE_CONTROL
		gf_image_control.init_existing_dom(image_element,
			[flow_name_str],

			current_host_str,
			logged_in_bool,
			p_plugin_callbacks_map,
			p_log_fun);
	});

	//---------------------
	// CURRENT_PAGE_DISPLAY

	const current_pages_display = gf_images_paging.init__current_pages_display(p_log_fun);
	$('body').append(current_pages_display);

	//---------------------
	// VIEW_TYPE_PICKER

	gf_view_type_picker.init(flow_name_str, p_log_fun);

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

				var current_image_view_type_str = gf_view_type_picker.image_view_type_str;

				const page_source_ref_str  = flow_name_str;
				const page_source_type_str = "flow"

				await gf_images_paging.load_new_page(page_source_ref_str,
					page_source_type_str,
					current_page_int,
					current_image_view_type_str,
					logged_in_bool,
					p_plugin_callbacks_map,
					p_log_fun);

				current_page_int += 1;
				$("#gf_images_flow_container").data("current_page", current_page_int);

				page_is_loading_bool = false;

				$(current_pages_display).find('#end_page').text(current_page_int);
				
			}
		}
	};

	//------------------
}

//---------------------------------------------------
function init_upload(p_flow_name_str :string,
	p_log_fun) {

	// use "" so that no host is set in URL's for issued requests
	// (forces usage of origin host that the page came from)
	const target_full_host_str = "";
	gf_upload__init(p_flow_name_str,
		target_full_host_str,
		
		//-------------------------------------------------
		// p_on_upload_fun
		async (p_upload_gf_image_id_str)=>{

			var image_exists_bool;
			var image_result_map;

			// start attempting to get uploaded image metadata, until the upload succeeds
			const attempts_num_int = 6;
			for (var i=0; i<5; i++) {
				
				//------------------
				// SLEEP - it takes time for the image to get uploaded.
				//         so dont run gf_images_http.get() until the system had time to add the image,
				//         otherwise it will return a response that the image doesnt exist yet.
				// ADD!! - some way to immediatelly display a placeholder for the image that is being uploaded.
				const wait_time_miliseconds_int = 1500; // 1s
				await gf_time.sleep(wait_time_miliseconds_int);

				//------------------
				// HTTP_GET_IMAGE
				image_result_map  = await gf_images_http.get(p_upload_gf_image_id_str, p_log_fun);
				image_exists_bool = image_result_map["image_exists_bool"];

				//------------------

				if (image_exists_bool) {

					// image now exists and we can stop attempting to fetch its metadata
					break;
				}
			}
			
			// uploaded image is not in the system even after all the retries,
			// so just display the failure and do nothign else
			if (!image_exists_bool) {

				// ERROR_DISPLAY
				$("body").append(`<div id='upload_display_failed'
					style='position:'fixed';right='20px';top='20px';background-color='red';width='10px';height='10px'>
					</div>`);
			}

			else {

				const image_export_map = image_result_map["image_export_map"];

				const img__format_str               = image_export_map["format_str"];
				const img__creation_unix_time_f     = image_export_map["creation_unix_time_f"];
				const img__owner_user_name_str      = image_export_map["user_name_str"];
				const img__flows_names_lst          = image_export_map["flows_names_lst"];
				const img__origin_page_url_str      = image_export_map["origin_page_url_str"];
				const img__thumbnail_small_url_str  = image_export_map["thumbnail_small_url_str"];
				const img__thumbnail_medium_url_str = image_export_map["thumbnail_medium_url_str"];
				const img__thumbnail_large_url_str  = image_export_map["thumbnail_large_url_str"];
				const img__title_str                = image_export_map['title_str'];
				const img__tags_lst                 = image_export_map["tags_lst"];
				

				const current_image_view_type_str = gf_view_type_picker.get_current_view_type();

				//------------------
				// IMAGE_CONTROL

				gf_image_control.create(p_upload_gf_image_id_str,
					img__format_str,
					img__creation_unix_time_f,
					img__origin_page_url_str,
					img__thumbnail_small_url_str,
					img__thumbnail_medium_url_str,
					img__thumbnail_large_url_str,
					img__title_str,
					img__tags_lst,
					img__owner_user_name_str,
					img__flows_names_lst,
					current_image_view_type_str,

					//---------------------------------------------------
					// p_on_img_load_fun
					(p_image_container)=>{
						// IMPORTANT!! - add ".gf_image" to the DOM after the image is fully loaded.
						// add it as the first element since its an uploaded image
						$("#gf_images_flow_container #items").prepend(p_image_container);

						// MASONRY_LAYOUT
						gf_utils.masonry_layout_after_img_load(p_image_container);
					},

					//---------------------------------------------------
					// p_on_img_load_error_fun
					()=>{},

					//---------------------------------------------------
					p_log_fun);

				//------------------
			}
		});

		//-------------------------------------------------
}