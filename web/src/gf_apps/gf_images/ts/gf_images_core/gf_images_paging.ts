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
import * as gf_tagger_http   from "./../../../gf_tagger/ts/gf_tagger_client/gf_tagger_http";

//---------------------------------------------------
export function init(p_initial_page_int :number,
	p_flow_name_str        :string,
	p_view_type_str        :string,
	p_logged_in_bool	   :boolean,
	p_events_enabled_bool  :boolean,
	p_plugin_callbacks_map :any,
	p_on_page_load_fun     :Function,
	p_log_fun              :Function) {

	var current_page_int = p_initial_page_int;
	var page_is_loading_bool = false;
	
	//---------------------------------------------------
	const scroll_handler_fun = async ()=>{

		// $(document).height() - height of the HTML document
		// window.innerHeight   - Height (in pixels) of the browser window viewport including, if rendered, the horizontal scrollbar
		if (window.scrollY >= $(document).height() - (window.innerHeight+50)) {
			
			// IMPORTANT!! - only load 1 page at a time
			if (!page_is_loading_bool) {
				
				page_is_loading_bool = true;
				p_log_fun("INFO", `current_page_int - ${current_page_int}`);

				const page_source_ref_str  = p_flow_name_str;
				const page_source_type_str = "flow"


				await load_new_pages(page_source_ref_str,
					current_page_int,
					p_view_type_str,
					p_logged_in_bool,
					p_plugin_callbacks_map,
					p_log_fun,
					1,
					page_source_type_str,
					p_events_enabled_bool);
				
				current_page_int += 1;
				page_is_loading_bool = false;

				// callback
				p_on_page_load_fun(current_page_int);
			}
		}
	};

	//---------------------------------------------------
	window.onscroll = scroll_handler_fun;
	return scroll_handler_fun;
}

//---------------------------------------------------
export async function load_new_pages(p_page_source_ref_str :string, // p_flow_name_str :string,
	p_current_page_int            :number,
	p_current_image_view_type_str :string,
	p_logged_in_bool              :boolean,
	p_plugin_callbacks_map        :any,
	p_log_fun                     :Function,
	p_pages_num_int               :number  = 1,
	p_page_source_type_str        :string  = "flow",
	p_events_enabled_bool         :boolean = true) {

	const gf_host_str = gf_core_utils.get_current_host();

	return new Promise(async function(p_resolve_fun, p_reject_fun) {

		//---------------------------------------------------
		async function fetch_pages() :Promise<Object> {
			return new Promise(async function(p_resolve_fun, p_reject_fun) {
				switch (p_page_source_type_str) {

					// FLOWS
					case "flow":

						const flow_name_str = p_page_source_ref_str;
						const resp_pg_map = await gf_images_http.get_page(flow_name_str,
							p_current_page_int,
							p_pages_num_int,
							p_log_fun);

						const pages_lst            = resp_pg_map["pages_lst"];
						const pages_user_names_lst = resp_pg_map["pages_user_names_lst"];
						
						p_resolve_fun({pages_lst, pages_user_names_lst});
						break;

					// TAGS
					case "tag":

						const tag_name_str = p_page_source_ref_str;
						const object_type_str = "image";

						const resp_objs_with_tag_map = await gf_tagger_http.get_objs_with_tag(tag_name_str,
							object_type_str,
							p_log_fun);


						const objects_with_tag_lst = resp_objs_with_tag_map["objects_with_tag_lst"];
						p_resolve_fun({objects_with_tag_lst});

						break;
				}
			});
		}
		
		//---------------------------------------------------
		
		const pages_map = await fetch_pages() as { pages_lst: Object[][]; pages_user_names_lst: string[][] };
		const pages_lst = pages_map.pages_lst;
		const pages_user_names_lst = pages_map.pages_user_names_lst;
		
		view_page(pages_lst, pages_user_names_lst);

		//---------------------------------------------------
		function view_page(p_pages_lst :any[][], p_pages_user_names_lst: string[][]) {

			console.log("VIEW PAGE...", p_pages_lst)

			var img_i_int = 0;
			var total_images_int = p_pages_lst.length*p_pages_lst[0].length;

			console.log("total_images_int", total_images_int)

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
					const image_element = gf_image_control.create(img__id_str,
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

						p_events_enabled_bool,
						gf_host_str,
						//---------------------------------------------------
						// p_on_img_load_fun
						(p_image_element :HTMLElement)=>{
							
							console.log(">>>> image")

							// MASONRY_LAYOUT
							gf_utils.masonry_layout_after_img_load(p_image_element);

							img_i_int++;

							//----------------
							// LOGGED_IN - only initialize this part if the user is authenticated
							
							if (p_logged_in_bool) {
									
								// TAGGING
								gf_utils.init_tagging(img__id_str,
									p_image_element,
									gf_host_str,
									p_log_fun);

								// SHARE
								gf_images_share.init(img__id_str,
									p_image_element,
									p_plugin_callbacks_map,
									p_log_fun);
							}

							//----------------

							// IMPORTANT!! - only declare load_new_page() as complete after all its
							//               images complete loading
							if (total_images_int-1 == img_i_int) {
								p_resolve_fun({});
							}
						},

						//---------------------------------------------------
						// p_on_img_load_error_fun
						()=>{
							// if image failed to load it still needs to be counted so that when all images
							// are done (either failed or succeeded) call p_on_complete_fun()
							img_i_int++;

							if (total_images_int-1 == img_i_int) {
								p_resolve_fun({});
							}
						},

						//---------------------------------------------------
						p_plugin_callbacks_map,

						// p_on_viz_change_fun
						()=>$("#gf_images_flow_container #items").masonry(),
						p_log_fun);
					
					// IMPORTANT!! - add image_element to container right away, so that image stay in-order
					$("#gf_images_flow_container #items").append(image_element);
				});
			});
		}

		//---------------------------------------------------
	});
}