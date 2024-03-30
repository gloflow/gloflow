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
import * as gf_gifs_viewer      from "./../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_image_viewer     from "./../../../../gf_core/ts/gf_image_viewer";
import * as gf_sys_panel        from "./../../../../gf_sys_panel/ts/gf_sys_panel";
import * as gf_identity         from "./../../../../gf_identity/ts/gf_identity";
import * as gf_identity_http    from "./../../../../gf_identity/ts/gf_identity_http";
import * as gf_images_http      from "./../gf_images_core/gf_images_http";
import * as gf_images_share     from "./../gf_images_core/gf_images_share";
import * as gf_paging           from "./gf_paging";
import * as gf_view_type_picker from "./gf_view_type_picker";
import * as gf_utils            from "./gf_utils";
import * as gf_flows_picker     from "./gf_flows_picker";

//-------------------------------------------------
declare var URLSearchParams;

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;
declare var gf_tagger__init_ui_v2;
declare var gf_tagger__http_add_tags_to_obj;

//-------------------------------------------------
export async function init(p_plugin_callbacks_map,
	p_log_fun) {

	const domain_str   = window.location.hostname;
	const protocol_str = window.location.protocol;
	const gf_host_str = `${protocol_str}//${domain_str}`;
	console.log("gf_host", gf_host_str);
	
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
	const urls_map          = gf_identity_http.get_standard_http_urls();
	const auth_http_api_map = gf_identity_http.get_http_api(urls_map);
	gf_identity.init_with_http(notifications_meta_map, urls_map);
	

	
	const parent_node = $("#right_section");
	const home_url_str = urls_map["home"];

	gf_identity.init_me_control(parent_node,
		auth_http_api_map,
		home_url_str);
	
	// inspect if user is logged-in or not
	const logged_in_bool = await auth_http_api_map["general"]["logged_in"]();

	//---------------------
	// UPLOAD__INIT

	const flow_name_str = get_current_flow();

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
	$('.gf_image').each((p_i, p_e)=>{

		const image_element = p_e;
		gf_utils.init_image_date(image_element, p_log_fun);

		const image_id_str = $(image_element).data('img_id');
		const img_thumb_medium_url_str = $(image_element).find('img').data('img_thumb_medium_url');
		const img_thumb_large_url_str  = $(image_element).find('img').data('img_thumb_large_url');
		const img_format_str           = $(image_element).data('img_format');
		const flows_names_lst          = $(image_element).data('img_flows_names').split(",");

		const origin_page_url_link = $(image_element).find(".origin_page_url a")[0];

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
		// GIFS
		if (img_format_str == 'gif') {			
			gf_gifs_viewer.init(image_element, image_id_str, flow_name_str, p_log_fun);
		}

		//----------------
		else {
			gf_image_viewer.init(image_element,
				image_id_str,
				img_thumb_medium_url_str,
				img_thumb_large_url_str,
				flows_names_lst,
				p_log_fun);
		}

		//----------------
		// LOGGED_IN - only initialize this part if the user is authenticated
		
		if (logged_in_bool) {
				
			// TAGGING
			init_tagging(image_id_str,
				image_element,
				gf_host_str,
				p_log_fun);



			// SHARE
			gf_images_share.init(image_id_str,
				image_element,
				p_plugin_callbacks_map,
				p_log_fun);

			//----------------
		}
	});

	const current_pages_display = gf_paging.init__current_pages_display(p_log_fun);
	$('body').append(current_pages_display);


	gf_view_type_picker.init(flow_name_str, p_log_fun);

	//------------------
	// LOAD_PAGES_ON_SCROLL

	var current_page_int = 6; // the few initial pages are already statically embedded in the document
	$("#gf_images_flow_container").data("current_page", current_page_int); // used in other functions to inspect current page

	var page_is_loading_bool = false;

	window.onscroll = ()=>{

		// $(document).height() - height of the HTML document
		// window.innerHeight   - Height (in pixels) of the browser window viewport including, if rendered, the horizontal scrollbar
		if (window.scrollY >= $(document).height() - (window.innerHeight+50)) {
			
			// IMPORTANT!! - only load 1 page at a time
			if (!page_is_loading_bool) {
				
				page_is_loading_bool = true;
				p_log_fun("INFO", `current_page_int - ${current_page_int}`);

				var current_image_view_type_str = gf_view_type_picker.image_view_type_str;
				gf_paging.load_new_page(flow_name_str,
					current_page_int,
					current_image_view_type_str,
					()=>{

						current_page_int += 1;
						$("#gf_images_flow_container").data("current_page", current_page_int);

						page_is_loading_bool = false;

						$(current_pages_display).find('#end_page').text(current_page_int);
					},
					p_log_fun);
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
				const img__tags_lst                 = image_export_map["tags_lst"];
				

				const current_image_view_type_str = gf_view_type_picker.get_current_view_type();

				gf_utils.init_image_element(p_upload_gf_image_id_str,
					img__format_str,
					img__creation_unix_time_f,
					img__origin_page_url_str,
					img__thumbnail_small_url_str,
					img__thumbnail_medium_url_str,
					img__thumbnail_large_url_str,
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
			}
		});

		//-------------------------------------------------
}

//---------------------------------------------------
// TAGGING_UI

function init_tagging(p_image_id_str,
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

		$(p_image_container_element)
			.find(".tags_container")
			.append(`<a class='gf_image_tag' href='/v1/tags/objects?tag=${p_tag_str}&otype=image'>#${p_tag_str}</a>`)
	}

	//-------------------------------------------------
}