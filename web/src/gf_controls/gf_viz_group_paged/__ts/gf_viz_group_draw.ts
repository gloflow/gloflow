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

//----------------------------------------------------
// p_elements_in_pages_lst - :List<:List<:DivElement>>

export function display_elements_in_pages(p_elements_in_pages_lst,
    p_layout_context_map,
    p_container_element,
    p_log_fun,
    
    p_layout_direction_str  = 'down',
    p_initial_layout_bool   = true,
    p_columns_int           = 4,
    p_elements_width_pixels = 200) {
    p_log_fun('FUN_ENTER', 'gf_vis_group_draw.display_elements_in_pages()');


    // flaten all elements in pages into a single list
    const elements_in_pages_flat_lst = p_elements_in_pages_lst.fold([], // initialValue, 
    	(p_accum_lst,
        p_elements_in_page_lst)=>{
    		p_accum_lst.addAll(p_elements_in_page_lst);
    		return p_accum_lst;
    	});
    
    const layout_context_map = gf_vis_layout_lib.layout_elements(elements_in_pages_flat_lst,
        p_layout_direction_str,
        p_layout_context_map,
        p_container_element,
        p_log_fun,

        p_columns_int,
        p_initial_layout_bool,
        p_elements_width_pixels);

    return layout_context_map;
}

//----------------------------------------------------
function draw_elements_in_pages(p_elements_infos_pages_lst,					 
    p_draw_element_fun,
    p_on_complete_fun,
    p_log_fun) {

    // List<:List<:DivElement>>
    const elements_in_pages_elements_2d_lst = new List(p_elements_infos_pages_lst.length);

    //----------------------------------------------------
    function draw_pages_on_complete_fun() {  
        p_on_complete_fun(elements_in_pages_elements_2d_lst);
    }

    //----------------------------------------------------  
    
    for (let p_elements_infos_in_page_map of p_elements_infos_pages_lst) {

        // assert(p_elements_infos_in_page_map.containsKey('page_index'));
        // assert(p_elements_infos_in_page_map.containsKey('elements_infos_lst'));

        var page_index                 = p_elements_infos_in_page_map['page_index'];
        const elements_infos_in_page_lst = p_elements_infos_in_page_map['elements_infos_lst'];


        _draw_elements_in_page(
            elements_infos_in_page_lst, 
            page_index,
            
            p_draw_element_fun,
            (p_elements_in_page_lst)=>{
                
                // extends the posts_elements_lst with elements from post_elements_in_page_lst
                // "page_index-1" - because list indexes start from 0, whereas page_indexes start from 1
                elements_in_pages_elements_2d_lst[page_index-1] = p_elements_in_page_lst; 
                
                // checks that all elements were created and added to the list
                if (!elements_in_pages_elements_2d_lst.contains(null)) {
                    draw_pages_on_complete_fun();
                }
            },
            p_log_fun);
    }
}

//----------------------------------------------------
// p_page_index - this is not the page_number, but a counter index

function _draw_elements_in_page(p_elements_infos_in_page_lst,
    p_page_index,
    p_draw_element_fun,
    p_on_complete_fun,
    p_log_fun) {

    const elements_in_page_lst = new List(p_elements_infos_in_page_lst.length);
    var current_post_index_int = 0;

    // by using a forEach loop, all draw_post() calls are issued in parallel, 
    // and as they complete thier results are added to the post_elements_in_page_lst 
    // results list  
    for (let p_element_info_map of p_elements_infos_in_page_lst) {
        p_draw_element_fun(
        	p_element_info_map,
        	current_post_index_int,
        	p_page_index,
            //----------------------------------------------------
        	(p_container_div_element,
            p_element_index_int)=>{
    			
    			//------------
    			// by using p_post_index it is guranteed that created Elements will
    			// be in the same order as p_post_info_dict in post_elements_in_page_lst
    			elements_in_page_lst[p_element_index_int] = p_container_div_element;
    			
                //------------
    		
    			// checks that all elements were created and added to the list
    			if (!elements_in_page_lst.contains(null)) {
    				p_on_complete_fun(elements_in_page_lst);
    			}
        	},

            //----------------------------------------------------
        	p_log_fun);

        current_post_index_int++;
    }
}