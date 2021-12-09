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

//-------------------------------------------------
export async function init(p_draw_element_fun,
	p_get_elements_pages_info_fun,
	p_log_fun,

	p_columns_int                           = 4,
	p_pages_in_page_set_number_int          = 4,
	p_initial_pages_number_to_display_int   = 4,
	p_on_scroll_pages_number_to_dislpay_int = 2) {




    const pages_cache_map = {
        'current_page_index_int': 1,
        'pages_map': {}
    };

    const layout_context_map = {};

    const elements_container_element = $(`
			<div class='gf_vis_group'>
			</div>`);
            
    $(elements_container_element).style("position", "relative");



    // REMOVE!! - this control should be able to be appended to any div
    $("body").append(elements_container_element);




    const start_page_index_int = pages_cache_map["current_page_index_int"];
    await pages_display(start_page_index_int,
        "down",
        pages_cache_map,
        layout_context_map,
        elements_container_element,

        p_draw_element_fun,
        p_get_elements_pages_info_fun,
        p_log_fun,
        
        p_columns_int,
        p_initial_pages_number_to_display_int,
        p_pages_in_page_set_number_int,
        true);

}

//----------------------------------------------------
function get_api_functions(p_pages_cache_map,
	p_layout_context_map,
	p_elements_container_element,

	p_draw_element_fun,
	p_get_elements_pages_info_fun,
	p_log_fun,
	p_columns_int                         = 4,
	p_pages_in_page_set_number_int        = 4,
	p_initial_pages_number_to_display_int = 4) {
	p_log_fun('FUN_ENTER', 'gf_vis_group.get_api_functions()');

	//----------------------------------------------------
	// used in displaying/loading new pages on scroll (for example)
	// p_initial_layout_bool - false - because these are additional pages that are displayed
	//                                 after the initial layout
	
	function pages_display_down_fun(p_on_complete_fun) {
		p_log_fun('FUN_ENTER', 'gf_vis_group.get_api_functions().pages_display_down_fun()');
		
		const current_page_index_int = p_pages_cache_map['current_page_index_int'];
		
		_pages_display(current_page_index_int, // p_start_page_index_int
			'down',                        // p_direction_str
			p_pages_cache_map,
			p_layout_context_map,
			p_elements_container_element,

			p_draw_element_fun,
			p_get_elements_pages_info_fun,

			// p_on_complete_fun
			(p_new_layout_context_map)=>{ 
				// new all_elements_height after the additional pages have been displayed
				const all_elements_height_int = p_new_layout_context_map['all_elements_height'];

				console.log('OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO');
				console.log(all_elements_height_int);
				console.log(p_on_complete_fun);

				p_on_complete_fun(all_elements_height_int);
			}, 
			p_log_fun,
			p_columns_int,
			p_initial_pages_number_to_display_int, // p_pages_number_to_display_int
			p_pages_in_page_set_number_int,
			false);                                // p_initial_layout_bool
	}

	//----------------------------------------------------
	// used in displaying/loading new pages on scroll (for example)
	// p_initial_layout_bool - false - because these are additional pages that are displayed
	//                                 after the initial layout
	
	function pages_display_up_fun(p_on_complete_fun) {
		p_log_fun('FUN_ENTER', 'gf_vis_group.get_api_functions().pages_display_up_fun()');
		
	    const current_page_index_int = p_pages_cache_map['current_page_index_int'];
	
		_pages_display(current_page_index_int, // p_start_page_index_int
				'up',                          // p_direction_str
				p_pages_cache_map,
				layout_context_map,
				p_elements_container_element,

				p_draw_element_fun,
				p_get_elements_pages_info_fun,

				// p_on_complete_fun
				(p_new_layout_context_map) => { 
					p_on_complete_fun(); 

					// new all_elements_height after the additional pages have been displayed
				    const all_elements_height_int = p_new_layout_context_map['all_elements_height'];
					p_on_complete_fun(all_elements_height_int);
				}, 

				p_log_fun,
				p_columns_int,
				p_initial_pages_number_to_display_int, // p_pages_number_to_display_int
				p_pages_in_page_set_number_int,
				false);                                // p_initial_layout_bool
	}

	//----------------------------------------------------
	// used for seeking 
	function pages_display_fun(p_page_index_int,
		p_on_complete_fun) {
		p_log_fun('FUN_ENTER', 'gf_vis_group.get_api_functions().pages_display_fun()');
		
		_pages_display(p_page_index_int, // p_start_page_index_int
			'down',                      // p_direction_str
			p_pages_cache_map,
			p_layout_context_map,
			p_elements_container_element,

			p_draw_element_fun,
			p_get_elements_pages_info_fun,

			// p_on_complete_fun,
			(p_new_layout_context_map)=>{p_on_complete_fun();},

			p_log_fun,
			p_columns_int,
			p_initial_pages_number_to_display_int, // p_pages_number_to_display_int
			p_pages_in_page_set_number_int,

			// pages_display_fun() is for seeking, and so on each seek the layout 
			// should be done as if its initial
			true); // p_initial_layout_bool
    }

    //----------------------------------------------------
	return {
		'pages_display_down_fun': pages_display_down_fun,
		'pages_display_up_fun':   pages_display_up_fun,
		'pages_display_fun':      pages_display_fun
	};
}

//-------------------------------------------------
async function pages_display(p_start_page_index_int,
	p_direction_str,
	p_pages_cache_map,
	p_layout_context_map,
	p_container_element,

	p_draw_element_fun,
	p_get_elements_pages_info_fun,
	p_log_fun,

	p_columns_int                  = 4,
	p_pages_number_to_display_int  = 3,
	p_pages_in_page_set_number_int = 3,
	p_initial_layout_bool          = true) {

    // int next_page_index_int = p_current_page_index + 1;
    const elements_in_pages_elements_lst = await gf_vis_group_get_pages.get_pages(p_start_page_index_int, // next_page_index_int, //int p_start_page_index,
        p_direction_str,
        p_pages_cache_map,

        p_draw_element_fun,
        p_get_elements_pages_info_fun,
        p_log_fun,
        p_pages_in_page_set_number_int);
        

    const new_layout_context_map = gf_vis_group_draw.display_elements_in_pages(elements_in_pages_elements_lst,
        p_layout_context_map,
        p_container_element,
        p_log_fun,
        p_columns_int,
        p_initial_layout_bool);
    
        
    return new_layout_context_map;
}