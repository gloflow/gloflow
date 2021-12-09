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
// if page is in local cache it gets it from there, if not
// it pulls it from the remote server

export function get_pages(p_start_page_index,
    p_direction_str,
    p_pages_cache_map,

    p_draw_element_fun,
    p_get_elements_pages_info_fun,
    p_log_fun,
    
    p_page_sets_to_get_number_int  = 2,
    p_pages_in_page_set_number_int = 1,
    p_page_size_int                = 6) {
        

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {


        const pages_to_get_int = p_page_sets_to_get_number_int * p_pages_in_page_set_number_int;
        const pages_lst        = new List(pages_to_get_int);  

        var i              = 0;
        var page_index_int = p_start_page_index + i;
        //----------------------------------------------------
        function handle_next() {
        	p_log_fun('FUN_ENTER', 'gf_vis_group_get_pages.get_pages().handle_next()');
        	
        	page_index_int = p_start_page_index + i;
        	
        	_get_page(page_index_int,
                p_pages_cache_map,
                p_direction_str,

                p_draw_element_fun,
                p_get_elements_pages_info_fun,

                // on_complete_fun
                (p_page_div_elements_lst)=>{
                	
                    pages_lst[i] = p_page_div_elements_lst;
                	  
                    // since pages_to_get_int is in the counting system that starts with 1,
                	// where indicies for pages_lst are in a counting system that starts with 0
                	if (i == (pages_to_get_int-1)) {
                        
                        p_resolve_fun(pages_lst);
                    }
                    else {
                	  	i += 1;
                	  	handle_next();
                    }
                },
                p_log_fun,
                p_pages_in_page_set_number_int,
                p_page_size_int);
        }

        //----------------------------------------------------

        _get_page(page_index_int,
            p_pages_cache_map,
            p_direction_str,

            p_draw_element_fun,
            p_get_elements_pages_info_fun,

            // on_complete_fun
            (p_page_div_elements_lst)=>{
                pages_lst[i] = p_page_div_elements_lst;
                handle_next();
            },
            p_log_fun,
            p_pages_in_page_set_number_int,
            p_page_size_int);
    }); 
    return p;
}

//----------------------------------------------------
function _get_page(p_page_index_int,
    p_pages_cache_map,
    p_direction_str,

    p_draw_element_fun,
    p_get_elements_pages_info_fun,
    p_on_complete_fun,
    p_log_fun,
    
    p_pages_in_page_set_number_int = 1,
    p_page_size_int                = 6) {
	p_log_fun('FUN_ENTER', 'gf_vis_group_get_pages._get_page()');

    //-------------
    // page is in cache
    if (p_pages_cache_map['pages_map'].containsKey(`${p_page_index_int}`)) {
    	p_log_fun('INFO', 'PAGE IS IN CACHE');
    	
        const page_divs_lst = p_pages_cache_map['pages_map'][`${p_page_index_int}`];

        p_pages_cache_map['current_page_index_int'] = p_page_index_int + 1;
        p_on_complete_fun(page_divs_lst);
    }
    
    //-------------
    // page is not in cache
    else {
    	p_log_fun('INFO', 'PAGE NOT IN CACHE');
    	
        //----------------------------------------------------     
        // p_posts_in_pages_elements_lst - :List<:List<:DivElement>>

        function on_complete_fun(p_elements_in_pages_elements_lst) {
        	p_log_fun('FUN_ENTER', 'gf_publisher_posts_image_view_get_pages.get_page().onComplete_fun()');
        	
            //-------------
            // store each of the newly acquired pages into cache, 
            // appropriatelly numbered with the right page number
            var j = p_page_index_int;
            for (let p_page_div_elements_lst of p_elements_in_pages_elements_lst) {
                p_log_fun('INFO', 'CACHING PAGE [$j]');

                // place the newly loaded and drawn elements into cache
                p_pages_cache_map['pages_map'][`${j}`] = p_page_div_elements_lst;
                j += 1;
            }
            
            //-------------
            // users of _get_page() expect one page to be returned, but for performance
            // reasons multiple pages are retreived from the server by _page_load_and_draw().
            // these extra pages that are not immediatelly displayed, but instead are placed in cache,
            // so that on future invocations of _get_page() they're first found in the cache,
            // instead of having to be retreived from cache...

            // FIX!! - this however means that if a page number in a certain direction are requested
            //         even though a certain page is not in cache and is requested from the server,
            //         other pages that are adjacent to it (in that direction) will also be returned
            //         from the server but those pages might already be cached... in that case they will
            //         be unnecessarily re-rendered and re-cached...

            const first_page_in_list_lst = p_elements_in_pages_elements_lst[0];
            p_pages_cache_map['current_page_index_int'] = p_page_index_int + 1;

            p_on_complete_fun(first_page_in_list_lst);

            //-------------
        }

        //---------------------------------------------------- 
        _page_load_and_draw(p_page_index_int,
            p_direction_str,

            p_draw_element_fun,
            p_get_elements_pages_info_fun,
            on_complete_fun,
            p_log_fun,
            p_pages_in_page_set_number_int,
            p_page_size_int);
    }

    //-------------
}

//---------------------------------------------------- 
// IMPORTANT!! - does not attach the post div's, it just creates their DIV's 
//               and returns that in a list

function _page_load_and_draw(p_start_page_index_int,
    p_direction_str,

    p_draw_element_fun,
    p_get_elements_pages_info_fun,
    p_on_complete_fun,
    p_log_fun,
    p_pages_in_page_set_number_int = 1,
    p_page_size_int                = 6) {
    p_log_fun('FUN_ENTER', 'gf_vis_group_get_pages._page_load_and_draw()');

    //----------------------------------------------------
    function on_complete_fun(p, // ????
        p_element_infos_by_pages_lst) {
            
        //----------------------------------------------------
        // IMPORTANT!! - does not attach the elements, it just creates their DIV's 
        //               and returns that
        gf_vis_group_draw.draw_elements_in_pages(p_element_infos_by_pages_lst,
            p_draw_element_fun,
            (p_elements_in_pages_elements_lst)=>{
                p_on_complete_fun(p_elements_in_pages_elements_lst);
            },
            p_log_fun);
            // p_page_size_int               :p_page_size_int,
            // p_pages_in_page_set_number_int:p_pages_in_page_set_number_int);
    }

    //----------------------------------------------------

    p_get_elements_pages_info_fun(p_start_page_index_int,
        p_direction_str,
        on_complete_fun,
        p_log_fun,

        p_page_size_int,
        p_pages_in_page_set_number_int);
}