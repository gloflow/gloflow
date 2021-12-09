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

///<reference path="../../../d/jquery.d.ts" />

import * as gf_viz_group_random_access from "./gf_viz_group_random_access";
import {GF_random_access_viz_props} from "./gf_viz_group_random_access";

declare var Draggabilly;

//-------------------------------------------------
export function init(p_id_str: string,
    p_parent_id_str: string,
    p_elements_lst,
    p_initial_page_int: number,
    p_element_create_fun,
    p_elements_page_get_fun,
    p_assets_uris_map,
    p_create_initial_elements_bool: boolean=true) {



    



    //------------------------
    var container;

    // check if a container with this name already exists, and if it does use that.
    // this is for cases where the DOM structure already exists (maybe from template rendering)
    // and there are some items already in that container.
    if ($(`#${p_id_str}`).length > 0) {
        container = $(`#${p_id_str}`)[0];
    }

    // otherwise create the div from scratch
    else {
        container = $(`<div id=${p_id_str}>
            <div id="items"></div>
        </div>`);
        $(`#${p_parent_id_str}`).append(container);
    }
    




    //------------------------
    // add elements
    if (p_create_initial_elements_bool) {
        for (let element_map of p_elements_lst) {


            const element = p_element_create_fun(element_map);
            $(element).addClass("item");

            $(container).find("#items").append(element);
        }
    }

    /*//------------------------
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
		itemSelector: '.item',
		columnWidth:  6
	});

    //------------------------*/


    $(container).find('img').on('load', ()=>{
		$(container).packery();
	});

    const packery_grid = $(container).packery({
        itemSelector: '.item',
        gutter: 10,
        // columnWidth: 60
    });

    // trigger initial layout
    packery_grid.packery();


    // enable draggability
    enable_draggability(container, packery_grid);



    //------------------------
    // INIT_RANDOM_ACCESS
    const viz_props :GF_random_access_viz_props = {
        seeker_container_height_px: 500,
        seeker_container_width_px:  100,
        seeker_bar_width_px:        50, 
        seeker_range_bar_width:     30,
        seeker_range_bar_height:    500,
        seeker_range_bar_color_str: "red",
        assets_uris_map: p_assets_uris_map,
    }

    const start_page_int = 0;
    const end_page_int   = 20;
    const random_access__container_element = gf_viz_group_random_access.init(start_page_int,
        end_page_int,
        viz_props,

        //-------------------------------------------------
        // p_viz_group_reset_fun
        (p_start_page_int :number,
        p_on_complete_fun)=>{
            
            reset_with_new_start_pages(container,
                p_start_page_int,
                p_element_create_fun,
                p_elements_page_get_fun);

            p_on_complete_fun();
        });
        
        //-------------------------------------------------

    // position seeker on the far right
    $(random_access__container_element).css("position", "absolute");
    $(random_access__container_element).css("right", "0px");
    $(container).append(random_access__container_element);

    //------------------------
	// LOAD_PAGES_ON_SCROLL

	var current_page_int     = p_initial_page_int; // the few initial pages are already statically embedded in the document
	var page_is_loading_bool = false;
    const pages_container = $(container).find("#items");

    window.onscroll = async ()=>{

		// $(document).height() - height of the HTML document
		// window.innerHeight   - Height (in pixels) of the browser window viewport including, if rendered, the horizontal scrollbar
		if (window.scrollY >= $(document).height() - (window.innerHeight+50)) {

            // IMPORTANT!! - only load 1 page at a time
			if (!page_is_loading_bool) {


                await load_new_pages(current_page_int,
                    pages_container,
                    p_element_create_fun,
                    p_elements_page_get_fun);


                page_is_loading_bool = false;

                current_page_int += 1;
				$(container).data("current_page", current_page_int);
            }
        }
    }

    //------------------------
}

//-------------------------------------------------
async function reset_with_new_start_pages(p_container,
    p_start_page_int :number, // this is where it was seeked to, and is different from first_page/last_page
    p_element_create_fun,
    p_elements_page_get_fun,

    // this is an initial load of viz_group, so load some >1 number of pages
    // starting from the page where to user seeked to.
    p_pages_to_get_num_int :number=6) {


    // remove all items currently displayed by viz_group, 
    // since new ones have to be shown
    $(p_container).find("#items .item").remove();

    const pages_container = $(p_container).find("#items");


    await load_new_pages(p_start_page_int,
        pages_container,
        p_element_create_fun,
        p_elements_page_get_fun,
        p_pages_to_get_num_int);

}

//-------------------------------------------------
function enable_draggability(p_container,
    p_packery_grid) {
    const draggable_items_lst = $(p_container).find(".item");
    $(draggable_items_lst).each((p_i, p_e)=>{

        const element = p_e; // $(this)[0];

        console.log(element);


        var draggie = new Draggabilly(element, {
            // options...
        });

        p_packery_grid.packery('bindDraggabillyEvents', draggie);
    });
}

//-------------------------------------------------
async function load_new_pages(p_page_index_int :number,
    p_pages_container,
    p_element_create_fun,
    p_elements_page_get_fun,
    p_pages_to_get_num_int :number=1) {



    // fetch page
    const elements_lst = await p_elements_page_get_fun(p_page_index_int, p_pages_to_get_num_int);
    
    // create elements
    for (let element_map of elements_lst) {


        const element = p_element_create_fun(element_map);
        $(element).addClass("item");
        $(element).css("visibility", "hidden"); // initially elements are not visible until they load


        $(p_pages_container).find("#items").append(element);


        console.log(element)
        $(element).find('img').on('load', function() {

            $(p_pages_container).packery();

            /*// MASONRY
            $(p_pages_container).masonry();
            $(p_pages_container).masonry(<any>"reloadItems");
            
            const masonry = $(p_pages_container).data('masonry');
            masonry.once('layoutComplete', (p_event, p_laid_out_items)=>{
                $(element).css('visibility', 'visible');
            });*/
        });
    }




}