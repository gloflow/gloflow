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
import * as gf_images_http from "./../gf_images_core/gf_images_http";
import * as gf_utils       from "./gf_utils";


//---------------------------------------------------
export async function load_new_page(p_flow_name_str :string,
	p_current_page_int            :number,
	p_current_image_view_type_str :string,
	p_on_complete_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_paging.load_new_page()');

	const page_lst = await gf_images_http.get_page(p_flow_name_str,
		p_current_page_int,
		p_log_fun);
	
	view_page(page_lst);

	//---------------------------------------------------
	function view_page(p_page_lst) {
		p_log_fun('FUN_ENTER', 'gf_paging.load_new_page().view_page()');

		var img_i_int = 0;
		$.each(p_page_lst, (p_i, p_e)=>{

			const img__id_str                   = p_e['id_str'];
			const img__format_str               = p_e['format_str'];
			const img__creation_unix_time_f     = p_e['creation_unix_time_f'];
			const img__origin_url_str           = p_e['origin_url_str'];
			const img__thumbnail_small_url_str  = p_e['thumbnail_small_url_str'];
			const img__thumbnail_medium_url_str = p_e['thumbnail_medium_url_str'];
			const img__thumbnail_large_url_str  = p_e['thumbnail_large_url_str'];
			const img__tags_lst                 = p_e['tags_lst'];
			const img__origin_page_url_str      = p_e['origin_page_url_str'];


			gf_utils.init_image_element(img__id_str,
				img__format_str,
				img__creation_unix_time_f,
				img__origin_url_str,
				img__thumbnail_small_url_str,
				img__thumbnail_medium_url_str,
				img__thumbnail_large_url_str,
				img__tags_lst,
				p_flow_name_str,
				p_current_image_view_type_str,

				//---------------------------------------------------
				// p_on_img_load_fun
				(p_image_container)=>{
					
					// IMPORTANT!! - add ".gf_image" to the DOM after the image is fully loaded.
					$("#gf_images_flow_container #items").append(p_image_container);

					// MASONRY_LAYOUT
					gf_utils.masonry_layout_after_img_load(p_image_container);

					img_i_int++;

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

			/*var img_url_str;
			switch (p_current_image_view_type_str) {
				case "small_view":
					img_url_str = img__thumbnail_small_url_str;
					break;
				case "medium_view":
					img_url_str = img__thumbnail_medium_url_str;
					break;
			}

			// IMPORTANT!! - '.gf_image' is initially invisible, and is faded into view when its image is fully loaded
			//               and its positioned appropriatelly in the Masonry grid
			const image_container = $(`
				<div class="gf_image item ${p_current_image_view_type_str}" data-img_id="${img__id_str}" data-img_format="${img__format_str}" style='visibility:hidden;'>
					<img src="${img_url_str}" data-img_thumb_medium_url="${img__thumbnail_medium_url_str}"></img>
					<div class="tags_container"></div>
					<div class="origin_page_url">
						<a href="${img__origin_page_url_str}" target="_blank">${img__origin_page_url_str}</a>
					</div>
					<div class="creation_time">${img__creation_unix_time_f}</div>
				</div>`);

			//------------------
			
			// FIX!! - this needs to happen after the image <div> is added to the DOM, 
			//         here reloading masonry layout doesnt have the intended effect, since 
			//         the image hasnt been added yet.
			//         move it to be after $("#gf_images_flow_container").append(image);

			$(image_container).find('img').on('load', ()=>{

				// IMPORTANT!! - add ".gf_image" to the DOM after the image is fully loaded
				$("#gf_images_flow_container #items").append(image_container);
				
				//------------------
				// MASONRY_RELOAD
				var masonry = $('#gf_images_flow_container #items').data('masonry');

				masonry.once('layoutComplete', (p_event, p_laid_out_items)=>{
					$(image_container).css('visibility', 'visible');
				});
				
				
				// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
				$('#gf_images_flow_container #items').masonry();
				$('#gf_images_flow_container #items').masonry(<any>"reloadItems");

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

				if (img__format_str == 'gif') {
					gf_gifs_viewer.init(image_container, img__id_str, p_flow_name_str, p_log_fun);
				} else {
					gf_image_viewer.init(image_container,
						img__thumbnail_medium_url_str,
						img__thumbnail_large_url_str,
						p_flow_name_str,
						p_log_fun);
				}

				//------------------

				img_i_int++;

				// IMPORTANT!! - only declare load_new_page() as complete after all its
				//               images complete loading
				if (p_page_lst.length-1 == img_i_int) {
					p_on_complete_fun();
				}
			});*/

			/*// IMAGE_FAILED_TO_LOAD
			$(image_container).find('img').on('error', function() {

				p_log_fun("ERROR", "IMAGE_FAILED_TO_LOAD ----------");

				// if image failed to load it still needs to be counted so that when all images
				// are done (either failed or succeeded) call p_on_complete_fun()
				img_i_int++;
				if (p_page_lst.length-1 == img_i_int) {
					p_on_complete_fun();
				}
			});

			//------------------
			gf_utils.init_image_date(image_container, p_log_fun);

			//------------------
			// TAGS
			if (img__tags_lst != null && img__tags_lst.length > 0) {
				$.each(img__tags_lst, function(p_i, p_tag_str) {
					const tag = $(
						`<a class='gf_image_tag' href='/v1/tags/objects?tag=${p_tag_str}&otype=image'>
							${p_tag_str}
						</a>`);

					$(image_container).find('.tags_container').append(tag);
				});
			}

			//------------------*/
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