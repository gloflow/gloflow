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

import 'package:gf_collection_scroll/gf_collection_scroll.dart' as gf_collection_scroll;
// import '../visGroup.dart' as visGroup;

//------------------------------------------------------------
export async function init(p_draw_element_fun,
	p_get_elements_pages_info_fun,
	p_on_new_pages_load_fun,
	p_log_fun,
	p_columns_int                           = 4,
	p_pages_in_page_set_number_int          = 4,
	p_initial_pages_number_to_display_int   = 4,
	p_on_scroll_pages_number_to_dislpay_int = 2,

	p_scroll_container_height_px = 600,
	p_scroll_container_width_px  = 600,
	p_scroll_bar_width_px        = 40,

	p_scroll_bar_color_str                    = 'rgb(255, 212, 178)',
	p_scroll_bar_button_color_str             = 'rgb(255, 144, 56)',
	p_scroll_bar_button_onMouseOver_color_str = 'rgb(235, 133, 52)') {
	p_log_fun('FUN_ENTER', 'gf_vis_group_scroll.init()');

	const p = new Promise(function(p_resolve_fun, p_reject_fun) {

		const result_map = visGroup.init(p_draw_element_fun,
			p_get_elements_pages_info_fun,
			p_log_fun,
			p_columns_int,
			p_pages_in_page_set_number_int,
			p_initial_pages_number_to_display_int,
			p_on_scroll_pages_number_to_dislpay_int);



		const visGroup_element       = result_map['visGroup_element'];
		const visGroup_height_int    = result_map['visGroup_height_int'];
		const pages_cache_map        = result_map['pages_cache_map'];
		const pages_display_down_fun = result_map['pages_display_down_fun'];
		const pages_display_up_fun   = result_map['pages_display_up_fun'];
		const pages_display_fun      = result_map['pages_display_fun'];

		// Map collection_scroll_info_map;
		const collection_scroll_info_map = gf_collection_scroll.init(visGroup_element, // p_target_content_element,
			visGroup_height_int,
			p_log_fun,
			p_scroll_container_height_px,
			p_scroll_container_width_px,
			p_scroll_bar_width_px,
			p_scroll_bar_color_str,
			p_scroll_bar_button_color_str,
			p_scroll_bar_button_onMouseOver_color_str,
			//------------------------------------------------------------
			// p_on_top_reached_fun
			(p_collection_scroll_info_map)=>{
				console.log('TOP REACHED+++++++++++++++++++++++++++++');
				console.log(p_collection_scroll_info_map);

				const scroll_reinitialize_fun = p_collection_scroll_info_map['reinitialize_fun'];

				p_on_new_pages_load_fun();

				pages_display_up_fun((p_all_elements_height_int)=>{
					scroll_reinitialize_fun(p_all_elements_height_int);
				});
			},

			//------------------------------------------------------------
			// p_on_bottom_reached_fun
			(p_collection_scroll_info_map)=>{
				console.log('BOTTOM REACHED+++++++++++++++++++++++++++++');
				console.log(p_collection_scroll_info_map);

				const scroll_reinitialize_fun = p_collection_scroll_info_map['reinitialize_fun'];

				p_on_new_pages_load_fun();

				pages_display_down_fun((p_all_elements_height_int)=>{
					scroll_reinitialize_fun(p_all_elements_height_int);
				});
			});

			//------------------------------------------------------------
			
		const visGroup_scroll_element = collection_scroll_info_map['scroll_container_element'];
		
		const final_result_map = {
			'visGroup_scroll_element': visGroup_scroll_element,
			'pages_cache_map':         pages_cache_map,
			'visGroup_height_int':     visGroup_height_int,
			'pages_display_fun':       pages_display_fun
		};
		p_resolve_fun(final_result_map);
	});
	return p;
}