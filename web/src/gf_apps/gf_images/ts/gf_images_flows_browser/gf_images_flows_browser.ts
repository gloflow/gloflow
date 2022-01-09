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

///<reference path="../../../../d/jquery.d.ts" />
///<reference path="../../../../d/masonry.layout.d.ts" />
///<reference path="../../../../d/jquery.timeago.d.ts" />

import * as gf_time             from "./../../../../gf_core/ts/gf_time";
import * as gf_gifs_viewer      from "./../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_image_viewer     from "./../../../../gf_core/ts/gf_image_viewer";
import * as gf_sys_panel        from "./../../../../gf_core/ts/gf_sys_panel";
import * as gf_image_http       from "../gf_images_core/gf_images_http";
import * as gf_paging           from "./gf_paging";
import * as gf_view_type_picker from "./gf_view_type_picker";
import * as gf_utils            from "./gf_utils";
import * as gf_flows_picker     from "./gf_flows_picker";

//-------------------------------------------------
declare var URLSearchParams;

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;

//-------------------------------------------------
$(document).ready(()=>{
    //-------------------------------------------------
    function log_fun(p_g,p_m) {
        var msg_str = p_g+':'+p_m
        //chrome.extension.getBackgroundPage().console.log(msg_str);

        switch (p_g) {
            case "INFO":
                console.log("%cINFO"+":"+"%c"+p_m,"color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
                break;
            case "FUN_ENTER":
                console.log("%cFUN_ENTER"+":"+"%c"+p_m,"color:yellow; background-color:lightgray","background-color:lightgray");
                break;
        }
    }

    //-------------------------------------------------
    init(log_fun);
});

//-------------------------------------------------
export function init(p_log_fun) {

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
	gf_sys_panel.init(p_log_fun);
	gf_flows_picker.init(p_log_fun);
	
	//---------------------
	// UPLOAD__INIT

	const flow_name_str = get_current_flow();

	init_upload(flow_name_str, p_log_fun);

	//---------------------

	//-----------------
	// IMPORTANT!! - as each image loads call masonry to reconfigure the view.
	//               this is necessary so that initial images in the page, before
	//               load_new_page() starts getting called, are properly laid out
	//               by masonry.
	$('.gf_image img').on('load', ()=>{
		$('#gf_images_flow_container #items').masonry();
	});

	$('#gf_images_flow_container #items').masonry({
		// options...
		itemSelector: '.item',
		columnWidth:  6
	});

	$('.gf_image').each((p_i, p_e)=>{

		const image_element = p_e;
		gf_utils.init_image_date(image_element, p_log_fun);

		const img_thumb_medium_url_str = $(image_element).find('img').data('img_thumb_medium_url');
		const img_thumb_large_url_str  = $(image_element).find('img').data('img_thumb_large_url');
		const img_format_str           = $(image_element).attr('data-img_format');


		// CLEANUP - for images that dont come from some origin page (direct uploads, or generated images)
		//           this origin_page_url is set to empty string. check for that and remove it.
		// FIX!! - potentially on the server/template-generation side this div node shouldnt get included
		//         at all for images that dont have an origin_page_url.
		if ($(image_element).find(".origin_page_url a").text().trim() == "") {
			$(image_element).find(".origin_page_url").remove();
		}

		//----------------
		// GIFS
		if (img_format_str == 'gif') {

			const img_id_str = $(image_element).attr('data-img_id');
			gf_gifs_viewer.init(image_element, img_id_str, flow_name_str, p_log_fun);
		}

		//----------------
		else {
			gf_image_viewer.init(image_element,
				img_thumb_medium_url_str,
				img_thumb_large_url_str,
				flow_name_str,
				p_log_fun);
		}

		//----------------
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
				//         so dont run gf_image_http.get() until the system had time to add the image,
				//         otherwise it will return a response that the image doesnt exist yet.
				// ADD!! - some way to immediatelly display a placeholder for the image that is being uploaded.
				const wait_time_miliseconds_int = 1500; // 1s
				await gf_time.sleep(wait_time_miliseconds_int);

				//------------------
				// HTTP_GET_IMAGE
				image_result_map  = await gf_image_http.get(p_upload_gf_image_id_str, p_log_fun);
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
					p_flow_name_str,
					current_image_view_type_str,

					//---------------------------------------------------
					// p_on_img_load_fun
					()=>{
						
					},

					//---------------------------------------------------
					// p_on_img_load_error_fun
					()=>{
						
					},

					//---------------------------------------------------
					p_log_fun);
			}
		});

		//-------------------------------------------------
}