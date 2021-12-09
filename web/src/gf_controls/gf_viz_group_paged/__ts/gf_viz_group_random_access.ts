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

import * as gf_vis_group_random_access_canvas from "gf_vis_group_random_access_canvas";
import * as gf_vis_group_scroll               from "gf_vis_group_scroll";

//------------------------------------------------------------
// p_button_seek_info_draw_fun - :Function - the user can pass in a callback that will get invoked when
//                                           the button_seek_info is drawn, and the DivElement that this 
//                                           function returns is then attached/placed right next to the
//                                           page_index indicator.

export async function init(p_seek_start_page_int,
	p_seek_end_page_int,

	p_draw_element_fun,
	p_get_elements_pages_info_fun,
	p_onComplete_fun,
	p_log_fun,

	p_columns_int                           = 4,
	p_pages_in_page_set_number_int          = 4,
	p_initial_pages_number_to_display_int   = 4,
	p_on_scroll_pages_number_to_dislpay_int = 2,
	p_scroll_container_height_px            = 600,
	p_scroll_container_width_px             = 600,
	p_scroll_bar_width_px                   = 40,
	p_seeker_container_width_px             = 800,
	p_seeker_bar_width_px                   = 60,
	p_button_seek_info_draw_fun, 
	p_button_atlas_url_str                  = 'buttons.png'}) {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access.init()');
  
	var seeker_bar_draw_cached_pages_fun;

	//------------------------------------------------------------
	// called whenever gf_vis_group_scroll loads new pages
	function on_new_pages_load_fun() {
		p_log_fun('FUN_ENTER','gf_vis_group_random_access.init().on_new_pages_load_fun()');

		// seeker_bar_draw_cached_pages_fun();
	}

	//------------------------------------------------------------

	const result_map = await gf_vis_group_scroll.init(p_draw_element_fun,
		p_get_elements_pages_info_fun,
		on_new_pages_load_fun,
		p_log_fun,
		p_columns_int,
		p_pages_in_page_set_number_int,
		p_initial_pages_number_to_display_int,
		p_on_scroll_pages_number_to_dislpay_int,

		p_scroll_container_height_px,
		p_scroll_container_width_px);
	

	const gf_vis_group_scroll_element    = result_map['gf_vis_group_scroll_element'];
	const pages_cache_map                = result_map['pages_cache_map'];
	const gf_vis_group_scroll_height_int = result_map['gf_vis_group_scroll_height_int'];
	const visGroup_pages_display_fun     = result_map['visGroup_pages_display_fun'];

													
	const seeker_bar_info_map = init_seeker_bar(gf_vis_group_scroll_element,
		gf_vis_group_scroll_height_int,
		p_seek_start_page_int,
		p_seek_end_page_int,
		pages_cache_map,
		visGroup_pages_display_fun,
		p_log_fun,

		p_seeker_container_width_px,
		p_scroll_container_height_px, // p_seeker_container_height_px,
		p_button_seek_info_draw_fun,
		p_seeker_bar_width_px,
		p_seeker_bar_width_px, // p_seeker_range_bar_width,
		p_button_atlas_url_str);
			 
	const seeker_bar_element               = seeker_bar_info_map['container_element'];	
	const seeker_bar_draw_cached_pages_fun = seeker_bar_info_map['draw_cached_pages_fun'];


	return seeker_bar_element;
}

//------------------------------------------------------------
function init_seeker_bar(p_gf_vis_group_scroll_element,
	p_gf_vis_group_scroll_height_int,
	p_seek_start_page_int,
	p_seek_end_page_int,

	p_pages_cache_map,           // used to visualize the cached pages
	p_visGroup_pages_display_fun,
	p_log_fun,

	p_seeker_container_height_px    = 800,
	p_seeker_container_width_px     = 820,
	p_seeker_bar_width_px           = 200,
	p_seeker_range_bar_width        = 60,
	p_seeker_range_bar_color_str    = 'rgb(168, 199, 162)',
	p_seeker_range_bar_button_width = 40,
	
    p_button_seek_info_draw_fun,
	p_seeker_range_bar_button_info_width_px  = 60,
	p_seeker_range_bar_button_info_height_px = 60,
	p_button_atlas_url_str                   = 'buttons.png') {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access.init_seeker_bar()');

	// seek_range_bar - this is the bar that represents the full range 
	//                  that can be seeked over, so if its time then the lenght
	//                  of this bar would represent the full time period over which
	//                  the user can randomly seek new page-ranges. 
	
	const container_element = $(`
		<div class='seeker_container'>
			<div class='seeker_bar'>
				<div class='seek_range_bar'>
					<canvas class='seek_range_bar_background'></canvas>
				</div>
				
				<div class='seek_range_bar_button'>
					<div class='conn'></div>
					
					<div class='button'>
						<div class='button_symbol'></div>
					</div>
					<div class='button_seek_info'>
						<canvas class='button_seek_info_background'></canvas>
						<div class='seek_page_index'>$p_seek_start_page_int</div>
					</div>
				</div>
			</div>
		</div>`);
	
	const seeker_bar_element              = $(container_element).find('.seeker_bar');
	const seeker_range_bar_element        = $(container_element).find('.seek_range_bar');
	const seeker_range_bar_button_element = $(container_element).find('.seek_range_bar_button');
	const seek_page_index                 = $(container_element).find('.seek_page_index');
	
	$(seeker_bar_element).on('mouseenter', (p_event)=>{
		$(p_gf_vis_group_scroll_element).style('opacity', '0.3');
	    $(seeker_range_bar_button_element).style('visibility', 'visible');
	});
	$(seeker_bar_element).on('mouseleave', (p_event)=>{
		$(p_gf_vis_group_scroll_element).style('opacity', '1');
		$(seeker_range_bar_button_element).style('visibility', 'hidden');
	});

	//------------------------------------------------------------
	function set_css_styling() {
		p_log_fun('FUN_ENTER', 'gf_vis_group_random_access.init_seeker_bar().set_css_styling()');
		
		//------------
		// CSS
		
		// SEEKER CONTAINER
		$(container_element).style('position', 'relative');         
		$(container_element).style('overflow', 'hidden');
		$(container_element).style('height',   `${p_seeker_container_height_px}px`);
		$(container_element).style('width',    `${p_seeker_container_width_px}px`);
		$(container_element).style('top',      '0px');

		// IMPORTANT!! - so that the user doesnt accidently select the entire random_access ui element
		//               as if its contents was text to be selected/highlighted
		$(container_element).style('user-select', 'none'); 
		
		// SEEKER BAR
		$(seeker_bar_element).style("position", 'absolute');
		$(seeker_bar_element).style("right",    '0px');
		$(seeker_bar_element).style("top",      '0px');
		$(seeker_bar_element).style("height",   `${p_seeker_container_height_px}px`);
		$(seeker_bar_element).style("width",    `${p_seeker_bar_width_px}px`);
		
		// SEEKER RANGER BAR
		$(seeker_range_bar_element).style("height",   `${p_seeker_container_height_px}px`);
		$(seeker_range_bar_element).style("width",    `${p_seeker_range_bar_width}px`);
		$(seeker_range_bar_element).style("position", 'absolute');
		$(seeker_range_bar_element).style("right",    '0px');
		$(seeker_range_bar_element).style("background-color", p_seeker_range_bar_color_str);
		
		// seek_page_index
		$(seek_page_index).style("left", `${(seek_page_index.offsetWidth - seek_page_index.offsetWidth)/2}px`);

		//------------
	}

	//------------------------------------------------------------
		
	set_css_styling();
	
	//------------
	// CANVAS-es (html5)

	const range_bar_background_element        = $(container_element).find('.seek_range_bar_background');
	const seek_info_button_background_element = $(container_element).find('.button_seek_info_background');
	
	const draw_cached_pages_fun = gf_vis_group_random_access_canvas.init_range_bar_background_canvas(range_bar_background_element,
		p_pages_cache_map,
		p_seek_start_page_int,
		p_seek_end_page_int,
		p_log_fun,

		p_seeker_range_bar_width,      // p_width_px
		p_seeker_container_height_px); // p_height_px
	
	gf_vis_group_random_access_canvas.init_button_seek_info_background_canvas(seek_info_button_background_element,
		p_log_fun,
		p_seeker_range_bar_button_info_height_px,  // p_height_px,
		p_seeker_range_bar_button_info_width_px);  // p_width_px

	//------------
	
	_init_seek_range_bar_button_dom(seeker_range_bar_button_element,
		seeker_range_bar_element,
		p_seeker_container_height_px,

		p_seek_start_page_int,
		p_seek_end_page_int,

		draw_cached_pages_fun,
		p_visGroup_pages_display_fun,
		p_button_seek_info_draw_fun,
		p_log_fun,

		p_seeker_range_bar_button_info_height_px, // p_button_info_height_px,
		p_seeker_range_bar_button_info_width_px,  // p_button_info_width_px,
		p_seeker_range_bar_width,                 // p_button_conn_width_px,
		p_button_atlas_url_str);                  // p_button_atlas_url_str
	
	// so that the seeker_bar_element is drawn over the gf_vis_group_scroll_element
	$(container_element).insertBefore(p_gf_vis_group_scroll_element, seeker_bar_element);
	
	return {
		'container_element':     container_element,
		'draw_cached_pages_fun': draw_cached_pages_fun
	};
}

//------------------------------------------------------------
function _init_seek_range_bar_button_dom(p_button_element,
	p_seeker_range_bar_element,
	p_seeker_range_bar_height_px,

	p_seek_start_page_int,
	p_seek_end_page_int, 

	p_draw_cached_pages_fun,
	p_visGroup_pages_display_fun,
	p_button_seek_info_draw_fun,
	p_log_fun,
	
    p_button_height_px      = 149,
	p_button_width_px       = 60,
	p_button_info_height_px = 60,
	p_button_info_width_px  = 60,

	p_button_color_str = 'rgb(162, 219, 40)',

	p_button_seek_info_y_offset_int    = 15,
	p_button_seek_info_color_str       = '889977',
	p_button_seek_page_index_font_size = 25,

	p_button_conn_width_px  = 107,
	p_button_conn_color_hex = 0xF33FE4,
	p_button_atlas_url_str  = 'buttons.png') {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access._init_seek_range_bar_button_dom()');
	
	const button           = $(p_button_element).find('.button');
	const button_seek_info = $(p_button_element).find('.button_seek_info');
	const seek_page_index  = $(p_button_element).find('.seek_page_index');
	const conn             = $(p_button_element).find('.conn');

	//------------------------------------------------------------
  	function init_button() {
  		p_log_fun('FUN_ENTER', 'gf_vis_group_random_access._init_seek_range_bar_button_dom().init_button()');
  	
  		//------------
		// CSS
		
		const button_symbol = $(button).find('.button_symbol'); 


		$(button).style("position", 'absolute');
		$(button).style("height",   `${p_button_height_px}px`);
		$(button).style("width",    `${p_button_width_px}px`);
		$(button).style("background-color", p_button_color_str;
		$(button).style("top",              '0px');
		$(button).style("cursor",           'pointer');
		$(button).style("opacity",          '0.7');

		$(button_symbol).style("position", 'absolute');
		$(button_symbol).style("width",    '37px');
		$(button_symbol).style("height",   '37px');
		$(button_symbol).style("top",      '54px');
		$(button_symbol).style("left",     '14px');
		$(button_symbol).style("background-image", `url('${p_button_atlas_url_str}')`);
		$(button_symbol).style("background-position", "38px 59px"); 
		
	  	// the 'button' is on the left side of the 'conn'
	  	// (positioned relative to the right edge)
	  	const button_right_px = p_button_conn_width_px + int.parse(conn.style.right.replaceFirst('px', ''));
	  	$(button).style("right", `${button_right_px}px`);
	  	
	  	
	  	$(button_seek_info).style("position",         'absolute');
	  	$(button_seek_info).style("height",           `${p_button_info_height_px}px`);
	  	$(button_seek_info).style("width",            `${p_button_info_width_px}px`);
	  	$(button_seek_info).style("background-color", p_button_seek_info_color_str);
	  	$(button_seek_info).style("top",             `${p_button_seek_info_y_offset_int}px`);
	  	
	  	// the 'button' is on the left side of the 'conn'
	  	const button_seek_info_right_px = p_button_width_px + button_right_px;
	  	$(button_seek_info).style("right", `${button_seek_info_right_px}px`);
	  	
	  	
	  	// $(seek_page_index).style("display         = 'block';
	  	$(seek_page_index).style("position", 'absolute');
	  	$(seek_page_index).style("top",  '15px');
	  	$(seek_page_index).style("left", '10px');
	  	// $(seek_page_index).style("marginLeft"      = 'auto';
	  	// $(seek_page_index).style("marginRight"     = 'auto';
	  	$(seek_page_index).style("font-size",   `${p_button_seek_page_index_font_size}px`);
	  	$(seek_page_index).style("font-family", '"Helvetica Neue", Helvetica, Arial, sans-serif');
	  	
	  	// formula to calculate seek_page_index 'left' (x) dynamically is:
	  	// (p_button_seek_info_element.offsetWidth - p_seek_page_index_element.offsetWidth)/2
	  	// but since both p_button_seek_info_element and p_seek_page_index_element dont have dimensions
	  	// assigned to them, Im using user-passed values here...
	  	// '+6' is used to account for extra padding added to the font_size, when seek_page_index is styled
	  	$(seek_page_index).style("left", `${(p_button_info_width_px - p_button_seek_page_index_font_size+6)/2}px`);
	  	//------------
	  	
	  	button.onMouseOver.listen((p_event) {
	  		$(button).style("opacity", '1');
	  	});
	  	
	  	button.onMouseOut.listen((p_event) {
	  		$(button).style("opacity", '0.7');
	  	});
	  	
	  	var button_mouse_move_subscription;
	  	var seek_percentage_int = 0;

	  	$(button).on("mousedown", (p_event)=>{
	  			
	  		// initial button relative coordinate where the user clicked
			const initial_click_coord  = p_event.offset.y.toInt();
			const initial_button_coord = button.offsetTop;
	
			// when user clicks and holds on the scroll button, a mouse move event handler is added 
			// to react to movement
			button_mouse_move_subscription = button.onMouseMove.listen((MouseEvent p_move_event) {
	
				// button relative coordinate of the move
				const movement_delta = p_move_event.offset.y.toInt() - initial_click_coord;
				
				const old_button_y   = button.offsetTop; 
				const new_button_y   = old_button_y + movement_delta;
			
				seek_percentage_int = handle_user_seek_event(new_button_y,
					p_button_height_px,
					p_seeker_range_bar_height_px,
					p_button_seek_info_y_offset_int,
					button,
					button_seek_info,
					seek_page_index,
					conn,
					
					p_seek_start_page_int,
					p_seek_end_page_int,

					p_visGroup_pages_display_fun,
					p_button_seek_info_draw_fun,
					p_log_fun);
			});
	  	});

  		$(button).on("mouseup", (p_event)=>{

			if (button_mouse_move_subscription != null) {
				button_mouse_move_subscription.cancel(); // cancel mouse move event handler
				button_mouse_move_subscription = null;
			}
			
			//--------------
			// IMPORTANT!! - since random seeking is potentially an expensive operation,
			//               it is only initialized when the user lets go of the seek button
			// FIX!! - how is the mousewheel seeking handled? (no mouse seeking? since it cant
			//         be done on touch devices anyway)
			
			seek_to_new_value(seek_percentage_int,
				p_seek_start_page_int,
				p_seek_end_page_int,

				p_visGroup_pages_display_fun,
				()=>{
					// since this new seek might load new pages into the cache, this re-displays all
					// cached pages
					p_draw_cached_pages_fun();
				},
				p_log_fun);

			//--------------
		});
		$(button).on("mouseleave", (p_event)=>{
			if (button_mouse_move_subscription != null) {
				button_mouse_move_subscription.cancel();
				button_mouse_move_subscription = null;
			}
		});
	}

	//------------------------------------------------------------
	function init_button_conn() {
		p_log_fun('FUN_ENTER', 'gf_vis_group_random_access._init_seek_range_bar_button_dom().init_button_conn()');
		
		//------------
		// CSS
		
		$(conn).style('position',        'absolute');
		$(conn).style('height',          '1px');
		$(conn).style('width',           `${p_button_conn_width_px}px`);
		$(conn).style('backgroundColor', 'black');
		$(conn).style('top',             '0px');
		
		$(conn).style('right', '0px'); // p_seeker_range_bar_element.style.width;

		//------------
	}

	//------------------------------------------------------------

	init_button_conn();
	init_button();
}

//------------------------------------------------
function seek_to_new_value(p_seek_percentage_int,
	p_seek_start_page_int,
	p_seek_end_page_int,

	p_visGroup_pages_display_fun,
	p_on_complete_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access.seek_to_new_value()');

	const page_index_to_seek_to = get_seek_page_index(p_seek_percentage_int,
		p_seek_start_page_int,
		p_seek_end_page_int,
		p_log_fun);
	
	p_visGroup_pages_display_fun(page_index_to_seek_to,
		// p_on_complete_fun
		()=>{
			p_on_complete_fun();
		});
}

//------------------------------------------------
function handle_user_seek_event(p_new_button_y,
    p_button_height_px,
	p_seeker_range_bar_height_px,
	p_button_seek_info_y_offset_int,
	p_button_element,
	p_button_seek_info_element,
	p_seek_page_index_element,
	p_conn_element,

	p_seek_start_page_int,
	p_seek_end_page_int,

	p_visGroup_pages_display_fun,
	p_button_seek_info_draw_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access.handle_user_seek_event()');
	
	// checks if the button reached the bottom or the top
	if (p_new_button_y + p_button_height_px <= p_seeker_range_bar_height_px && 
		p_new_button_y >= 0) {
		p_button_element.style.top = '${p_new_button_y}px';
		
		const seek_percentage = get_seek_percentage(p_button_element,
			p_seeker_range_bar_height_px,
			p_log_fun);

		//-----------
		// CONN POSITIONING
		
		// 'conn' has negligable height and so its positioning is different from the buttom itself
		// 1. p_seeker_range_bar_height_px : 100 = x : seek_percentage
		// 2. p_seeker_range_bar_height_px * seek_percentage = 100 * x
		const conn_y = (p_seeker_range_bar_height_px * seek_percentage) / 100;
		p_conn_element.style.top = `${conn_y.toFixed()}px`;
		
		//-----------
		
		const seek_page_index_int = get_seek_page_index(seek_percentage,
			p_seek_start_page_int,
			p_seek_end_page_int,
			p_log_fun);
		
		if (p_button_seek_info_element != null && p_seek_page_index_element != null) {
			display_button_seek_info(p_button_seek_info_element,
				p_seek_page_index_element,
				seek_page_index_int,
				p_new_button_y,
				p_button_seek_info_y_offset_int, 

				p_button_seek_info_draw_fun,
				p_log_fun);
		}
		
		return seek_percentage;
	}
	// BOTTOM OR TOP REACHED
	else {
		// BOTTOM
		if (p_new_button_y + p_button_height_px <= p_seeker_range_bar_height_px) {
			return 100;
		} 
		// TOPs
		else if (p_new_button_y >= 0) {
			return 0;
		}
	}

	return 0;
}

//------------------------------------------------
// this is called repeatedly as the user is seeking around
// p_seek_page_index - the index of the page that the user would seek to if he/she
//                     initiated the seek (released the seek button)

function display_button_seek_info(p_button_seek_info_element,
	p_seek_page_index_element,
	p_seek_page_index,
	p_new_button_y,
	p_button_seek_info_y_offset_int, // how far from top of parent div container

	p_button_seek_info_draw_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access.display_button_seek_info()');
	
	$(p_seek_page_index_element).text(`${p_seek_page_index}`);
	$(p_seek_page_index_element).style("left", `${(p_button_seek_info_element.offsetWidth - p_seek_page_index_element.offsetWidth)/2}px`);
	
	$(p_button_seek_info_element).style("top", `${p_new_button_y+p_button_seek_info_y_offset_int}px`);
	
	//--------------
	// IMPORTANT!! - user-passed callback execution
	
	if (p_button_seek_info_draw_fun != null){
		p_button_seek_info_draw_fun();
	}

	//--------------
}

//------------------------------------------------
function get_seek_page_index(p_seek_percentage_int,
    p_seek_start_page_int,
	p_seek_end_page_int,
	p_log_fun) {
	
	const total_range = p_seek_end_page_int - p_seek_start_page_int;
	
	// 1. 100 : total_range = p_seek_percentage_int : x
	// 2. 100 * x           = total_range * p_seek_percentage_int
	
	const page_index_delta      = (total_range * p_seek_percentage_int) / 100;  
	const page_index_to_seek_to = p_seek_start_page_int + page_index_delta.toInt();
	
	return page_index_to_seek_to;
}

//------------------------------------------------------------
function get_seek_percentage(p_button_element,
	p_seeker_range_bar_height_px,
	p_log_fun) {
	
	// MATH EXPLANATION:
	// 1. (p_seeker_range_bar_height_px - p_button_element.offsetHeight):100 = p_button_element.offsetTop:x
	// 2. x = (100 * p_button_element.offsetTop)/(p_seeker_range_bar_height_px - p_button_element.offsetHeight)
	
	const seek_percentage = (100 * p_button_element.offsetTop) / 
	    (p_seeker_range_bar_height_px - p_button_element.offsetHeight);
	return seek_percentage.toInt();
}