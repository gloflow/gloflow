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

import * as gf_images_paging from "../gf_images_core/gf_images_paging";
import * as gf_viz_group_random_access from "./gf_viz_group_random_access";
import * as gf_flows_browser_utils from "./gf_flows_browser_utils";

//-------------------------------------------------
export interface GF_props {
	
	readonly flow_name_str           :string
	readonly container_id_str        :string
	readonly parent_container_id_str :string

	readonly start_page_int   :number
	readonly end_page_int     :number
	readonly initial_page_int :number

	readonly image_view_type_str :string
	readonly logged_in_bool :boolean
	readonly events_enabled_bool :boolean

	readonly plugin_callbacks_map :any
}

//-------------------------------------------------
export function init(p_elements_lst :Array<any>,
	p_props                 :GF_props,
	p_element_create_fun    :Function,
	p_elements_page_get_fun :Function,
	p_log_fun                      :Function,
	p_create_initial_elements_bool :boolean=true,
	p_initial_pages_num_int        :number=6) {

	//------------------------
	// CONTROL_CONTAINER
	
	const container_id_str = p_props.container_id_str;
	var container :HTMLElement;

	// check if a container with this name already exists, and if it does use that.
	// this is for cases where the DOM structure already exists (maybe from template rendering)
	// and there are some items already in that container.
	if ($(`#${container_id_str}`).length > 0) {
		container = $(`#${container_id_str}`)[0];
	}

	// otherwise create the div from scratch
	else {
		container = $(`<div id=${container_id_str}>
			<div id="items"></div>
		</div>`)[0];
		$(`#${p_props.parent_container_id_str}`).append(container);
	}
	
	const items_container = $(container).find("#items");

	//------------------------
	// CREATE_ELEMENTS - initial elements displayed in the control, before paging is done

	if (p_create_initial_elements_bool) {
		for (let element_map of p_elements_lst) {

			const element = p_element_create_fun(element_map);
			$(element).addClass("item");

			$(items_container).append(element);
		}
	}

	//------------------------
	// MASONRY

	// IMPORTANT!! - as each image loads call masonry to reconfigure the view.
	//               this is necessary so that initial images in the page, before
	//               load_new_page() starts getting called, are properly laid out
	//               by masonry.
	$('#elements img').on('load', ()=>{
		$('#elements').masonry();
		$('#elements').masonry(<any>"reloadItems");
	});

	$(container).masonry({
		// options...
		itemSelector: '.gf_image',
		columnWidth:  6,
		gutter: 10,
	});

	//------------------------
	// CURRENT_PAGE
	/*
	indicates what the current page is. needed to share that state between random access
	and regular page loading logic based on scroll.
	this is needed for regular page loading on scroll to know where to load from since random_access
	can change what the start page is.
	the initial value with p_props is to account for the few initial pages are already
	statically embedded in the document.
	*/
	var current_page_int = p_props.initial_page_int;

	//------------------------
	// RANDOM_ACCESS_INIT
	
	const seeker__container_element = gf_viz_group_random_access.init(container,
		p_props.start_page_int,
		p_props.end_page_int,

		//-------------------------------------------------
		// RESET
		(p_page_index_to_seek_to_int :number,
		p_on_complete_fun :Function)=>{
			
			const new_start_page_int = p_page_index_to_seek_to_int;

			reset_with_new_start_pages(container,
				new_start_page_int,
				p_initial_pages_num_int,
				p_element_create_fun,
				p_elements_page_get_fun);

			// user seeked to a new random page, so that should be set
			// as the current page plus the initial pages that are loaded on reset.
			const new_current_page_int = new_start_page_int + p_initial_pages_num_int;
			current_page_int = new_current_page_int;

			gf_flows_browser_utils.current_pages_display__reset(new_start_page_int,
				new_current_page_int);

			p_on_complete_fun();
		});
		
		//-------------------------------------------------

	//------------------------
	// LOAD_PAGES_ON_SCROLL
	
	var initial_page_int = 6;

	gf_images_paging.init(initial_page_int,
		p_props.flow_name_str,
		p_props.image_view_type_str,
		p_props.logged_in_bool,
		p_props.events_enabled_bool,
		p_props.plugin_callbacks_map,

		// p_on_page_load_fun
		(p_new_page_int :number)=>{

			$("#gf_images_flow_container").data("current_page", p_new_page_int);

			// $(current_pages_display).find('#end_page').text(p_new_page_int);
		},
		p_log_fun);

	/*
	const pages_container = $(container).find("#items")[0];

	var page_is_loading_bool = false;

	window.onscroll = async ()=>{

		// $(document).height() - height of the HTML document
		// window.innerHeight   - Height (in pixels) of the browser window viewport including, if rendered, the horizontal scrollbar
		if (window.scrollY >= $(document).height() - (window.innerHeight+50)) {

			// IMPORTANT!! - only load 1 page at a time
			if (!page_is_loading_bool) {


				page_is_loading_bool = true;
				p_log_fun("INFO", `current_page_int - ${current_page_int}`);

				load_new_pages(current_page_int,
					pages_container,
					p_element_create_fun,
					p_elements_page_get_fun);


				page_is_loading_bool = false;

				current_page_int += 1;
				$(container).data("current_page", current_page_int);
			}
		}
	}
	*/

	//------------------------
	return seeker__container_element;
}

//-------------------------------------------------
// RESET_WITH_NEW_START_PAGES

async function reset_with_new_start_pages(p_container :HTMLElement,
	p_start_page_int        :number, // this is where it was seeked to, and is different from first_page/last_page
	p_initial_pages_num_int :number,
	p_element_create_fun    :Function,
	p_elements_page_get_fun :Function) {

	//------------------------
	// REMOVE_ALL - items currently displayed by viz_group, 
	//              since new ones have to be shown.
	$(p_container).find("#items .item").remove();

	//------------------------

	// IMPORTANT!! - do a layout after removing all items
	// p_packery_instance.packery("layout");
	$(p_container).masonry(<any>"reloadItems");



	const pages_container = $(p_container).find("#items")[0];

	load_new_pages(p_start_page_int,
		pages_container,
		p_element_create_fun,
		p_elements_page_get_fun,
		p_initial_pages_num_int);    
}

//-------------------------------------------------
async function load_new_pages(p_page_index_int :number,
	p_pages_container       :HTMLElement,
	p_element_create_fun    :Function,
	p_elements_page_get_fun :Function,
	p_pages_to_load_int     :number=1) {

	// fetch page
	const elements_lst = await p_elements_page_get_fun(p_page_index_int, p_pages_to_load_int);

	// create elements
	for (let element_map of elements_lst) {

		const element = p_element_create_fun(element_map);
		$(element).addClass("item");
		$(element).css("visibility", "hidden"); // initially elements are not visible until they load

		$(p_pages_container).append(element);

		$(element).find('img').on('load', ()=>{

			// make element visible after its image loads
			$(element).css('visibility', 'visible');
			
			// MASONRY
			$(p_pages_container).masonry();
			$(p_pages_container).masonry(<any>"reloadItems");
			
			/*
			const masonry = $(p_pages_container).data('masonry');
			masonry.once('layoutComplete', (p_event, p_laid_out_items)=>{
				$(element).css('visibility', 'visible');
			});
			*/
		});
	}

	return elements_lst;
}

//-------------------------------------------------