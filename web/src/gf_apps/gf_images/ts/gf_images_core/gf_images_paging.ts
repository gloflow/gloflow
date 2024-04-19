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

// import * as gf_image_viewer from "./../../../../gf_core/ts/gf_image_viewer";
// import * as gf_gifs_viewer  from "./../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_core_utils    from "../../../../gf_core/ts/gf_utils";
import * as gf_images_http   from "./gf_images_http";
import * as gf_utils         from "./gf_utils";
import * as gf_image_control from "./gf_image_control";
import * as gf_images_share  from "./gf_images_share";


//---------------------------------------------------
export async function load_new_page(p_flow_name_str :string,
	p_current_page_int            :number,
	p_current_image_view_type_str :string,
	p_logged_in_bool              :boolean,
	p_plugin_callbacks_map,
	p_on_complete_fun,
	p_log_fun) {

	const response_map = await gf_images_http.get_page(p_flow_name_str,
		p_current_page_int,
		p_log_fun);
	
	const pages_lst            = response_map["pages_lst"];
	const pages_user_names_lst = response_map["pages_user_names_lst"];
	
	view_page(pages_lst, pages_user_names_lst);


	const gf_host_str = gf_core_utils.get_gf_host();
	
	//---------------------------------------------------
	function view_page(p_pages_lst, p_pages_user_names_lst) {

		var img_i_int = 0;
		$.each(p_pages_lst, (p_i, p_page_lst)=>{
			$.each(p_page_lst, (p_j, p_e)=>{

				const img__id_str                   = p_e['id_str'];
				const img__format_str               = p_e['format_str'];
				const img__creation_unix_time_f     = p_e['creation_unix_time_f'];
				const img__flows_names_lst          = p_e["flows_names_lst"];
				const img__origin_url_str           = p_e['origin_url_str'];
				const img__thumbnail_small_url_str  = p_e['thumbnail_small_url_str'];
				const img__thumbnail_medium_url_str = p_e['thumbnail_medium_url_str'];
				const img__thumbnail_large_url_str  = p_e['thumbnail_large_url_str'];
				const img__title_str                = p_e['title_str'];
				const img__tags_lst                 = p_e['tags_lst'];
				const img__origin_page_url_str      = p_e['origin_page_url_str'];
				const img__owner_user_name_str      = p_pages_user_names_lst[p_i][p_j];

				// IMAGE_CONTROL
				gf_image_control.create(img__id_str,
					img__format_str,
					img__creation_unix_time_f,
					img__origin_url_str,
					img__thumbnail_small_url_str,
					img__thumbnail_medium_url_str,
					img__thumbnail_large_url_str,
					img__title_str,
					img__tags_lst,
					img__owner_user_name_str,
					img__flows_names_lst,
					p_current_image_view_type_str,

					//---------------------------------------------------
					// p_on_img_load_fun
					(p_image_container)=>{
						
						// IMPORTANT!! - add ".gf_image" to the DOM after the image is fully loaded.
						$("#gf_images_flow_container #items").append(p_image_container);

						// MASONRY_LAYOUT
						gf_utils.masonry_layout_after_img_load(p_image_container);

						img_i_int++;

						//----------------
						// LOGGED_IN - only initialize this part if the user is authenticated
						
						if (p_logged_in_bool) {
								
							// TAGGING
							gf_utils.init_tagging(img__id_str,
								p_image_container,
								gf_host_str,
								p_log_fun);

							// SHARE
							gf_images_share.init(img__id_str,
								p_image_container,
								p_plugin_callbacks_map,
								p_log_fun);
						}

						//----------------


						// IMPORTANT!! - only declare load_new_page() as complete after all its
						//               images complete loading
						if (p_page_lst.length-1 == img_i_int) {
							p_on_complete_fun();
						}
					},

					//---------------------------------------------------
					// p_on_img_load_error_fun
					()=>{
						// if image failed to load it still needs to be counted so that when all images
						// are done (either failed or succeeded) call p_on_complete_fun()
						img_i_int++;

						if (p_page_lst.length-1 == img_i_int) {
							p_on_complete_fun();
						}
					},

					//---------------------------------------------------
					p_log_fun);
			});
		});
	}

	//---------------------------------------------------
}

//---------------------------------------------------
export function init__current_pages_display(p_log_fun) {
	// p_log_fun('FUN_ENTER', 'gf_paging.init__current_pages_display()');

	const container = $(`
		<div id="current_pages_display"'>
			<div id="title">pages:</div>
			<div id="start_page">1</div>
			<div id="to">to</div>
			<div id="end_page">6</div>
		</div>`);

	return container;
}