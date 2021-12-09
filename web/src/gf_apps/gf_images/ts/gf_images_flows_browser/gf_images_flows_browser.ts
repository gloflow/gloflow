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

import * as gf_viz_group_paged from "./../../../../gf_controls/gf_viz_group_paged/ts/gf_viz_group_paged";
import * as gf_gifs_viewer  from "./../../../../gf_core/ts/gf_gifs_viewer";
import * as gf_image_viewer from "./../../../../gf_core/ts/gf_image_viewer";
import * as gf_sys_panel    from "./../../../../gf_core/ts/gf_sys_panel";

//-------------------------------------------------
declare var URLSearchParams;

// FIX!! - remove this from global scope!!
var image_view_type_str = "small_view";

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

	//-----------------
	// GET FLOW_NAME
	const url_params       = new URLSearchParams(window.location.search);
	const qs_flow_name_str = url_params.get('fname');
	var   flow_name_str;

	if (qs_flow_name_str == null) {
		flow_name_str = 'general'; // default value
	} else {
		flow_name_str = qs_flow_name_str;
	}
	
	//-----------------
	gf_sys_panel.init(p_log_fun);

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
		init_image_date(image_element, p_log_fun);

		const img_thumb_medium_url_str = $(image_element).find('img').data('img_thumb_medium_url');
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
			gf_image_viewer.init(image_element, img_thumb_medium_url_str, flow_name_str, p_log_fun);
		}

		//----------------
	});

	const current_pages_display = init__current_pages_display(p_log_fun);
	$('body').append(current_pages_display);


	init__view_type_picker(flow_name_str, p_log_fun);

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
				p_log_fun("INFO", "current_page_int - "+current_page_int);

				var current_image_view_type_str = image_view_type_str;
				load_new_page(flow_name_str,
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
// view_type_picker - picks the type of view that is used to display images in a flow.
//                    default is masonry with 6 columns.
function init__view_type_picker(p_flow_name_str :string,
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
		const current_page_int = $("#gf_images_flow_container").data("current_page");
		init__viz_group_view(p_flow_name_str,
			current_page_int,
			p_log_fun);

		//------------------
	});
}

//---------------------------------------------------
function init__viz_group_view(p_flow_name_str :string,
	p_initial_page_int :number,
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
			http__load_new_page(p_flow_name_str,
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
    gf_viz_group_paged.init(id_str,
        parent_id_str,
        initial_elements_lst,
		p_initial_page_int,
        element_create_fun,
        elements_page_get_fun,

		// the container already contains elements that are created and attached
		// to the container, so we dont want to create any initially (only if paging is done).
		false); // p_create_initial_elements_bool
}

//---------------------------------------------------
function init__current_pages_display(p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_images_flows_browser.init__current_pages_display()');

	const container = $(`
		<div id="current_pages_display"'>
			<div id="title">pages:</div>
			<div id="start_page">1</div>
			<div id="to">to</div>
			<div id="end_page">6</div>
		</div>`);

	return container;
}

//---------------------------------------------------
function load_new_page(p_flow_name_str :string,
	p_current_page_int :number,
	p_current_image_view_type_str :string,
	p_on_complete_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_images_flows_browser.load_new_page()');

	http__load_new_page(p_flow_name_str,
		p_current_page_int,
		(p_page_lst)=>{
			view_page(p_page_lst);
		},
		(p_error)=>{
			p_on_complete_fun();
		},
		p_log_fun);

	//---------------------------------------------------
	function view_page(p_page_lst) {
		p_log_fun('FUN_ENTER', 'gf_images_flows_browser.load_new_page().view_page()');

		var img_i_int = 0;
		$.each(p_page_lst, (p_i, p_e)=>{

			const img__id_str                   = p_e['id_str'];
			const img__format_str               = p_e['format_str'];
			const img__creation_unix_time_f     = p_e['creation_unix_time_f'];
			const img__origin_url_str           = p_e['origin_url_str'];
			const img__thumbnail_small_url_str  = p_e['thumbnail_small_url_str'];
			const img__thumbnail_medium_url_str = p_e['thumbnail_medium_url_str'];
			const img__tags_lst                 = p_e['tags_lst'];
			const img__origin_page_url_str      = p_e['origin_page_url_str'];


			var img_url_str;
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

			$(image_container).find('img').on('load', function() {

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
					gf_image_viewer.init(image_container, img__thumbnail_medium_url_str, p_flow_name_str, p_log_fun);
				}

				//------------------

				img_i_int++;

				// IMPORTANT!! - only declare load_new_page() as complete after all its
				//               images complete loading
				if (p_page_lst.length-1 == img_i_int) {
					p_on_complete_fun();
				}
			});

			// IMAGE_FAILED_TO_LOAD
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
			init_image_date(image_container, p_log_fun);

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

			//------------------
		});
	}

	//---------------------------------------------------
}

//-------------------------------------------------
function init_image_date(p_image_element, p_log_fun) {
	// p_log_fun('FUN_ENTER', 'gf_images_flows_browser.init_image_date()');

	const creation_time_element = $(p_image_element).find('.creation_time');
	const creation_time_f       = parseFloat($(creation_time_element).text());
	const creation_date         = new Date(creation_time_f*1000);

	const date_msg_str = $.timeago(creation_date);
	$(creation_time_element).text(date_msg_str);

	const creation_date__readable_str = creation_date.toDateString();
	const creation_date__readble      = $(`<div class="full_creation_date">${creation_date__readable_str}</div>`);

	$(creation_time_element).mouseover((p_e)=>{
		$(creation_time_element).append(creation_date__readble);

		// IMPORTANT!! - image size changed, so recalculate the Masonry layout.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container').masonry();
		$('#gf_images_flow_container').masonry(<any>'reloadItems');
	});

	$(creation_time_element).mouseout((p_e)=>{
		$(creation_date__readble).remove();

		// IMPORTANT!! - image size changed, so recalculate the Masonry layout.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container').masonry();
		$('#gf_images_flow_container').masonry(<any>'reloadItems');
		
	});
}

//---------------------------------------------------
function http__load_new_page(p_flow_name_str :string,
	p_current_page_int :number,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_images_flows_browser.http__load_new_page()');

	const page_size_int = 10;
	const url_str       = `/images/flows/browser_page?fname=${p_flow_name_str}&pg_index=${p_current_page_int}&pg_size=${page_size_int}`;
	p_log_fun('INFO', 'url_str - '+url_str);

	//-------------------------
	// HTTP AJAX
	$.get(url_str,
		function(p_data_map) {
			console.log('response received');
			// const data_map = JSON.parse(p_data);

			console.log('data_map["status"] - '+p_data_map["status"]);
			
			if (p_data_map["status"] == 'OK') {

				const pages_lst = p_data_map['data']['pages_lst'];
				p_on_complete_fun(pages_lst);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});

	//-------------------------	
}