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

//-------------------------------------------------
interface GF_viz_props {
	readonly seeker_container_height_px :number;
	readonly seeker_container_width_px  :number;
	readonly seeker_range_bar_height    :number;
	readonly bar_canvas_ctx :CanvasRenderingContext2D;
}

//-------------------------------------------------
export function init(p_parent_container :HTMLElement,
	p_first_page_int      :number,
	p_last_page_int       :number,
	p_viz_group_reset_fun :Function) {

	const asset_uri__gf_bar_handle_btn_str = "https://gf-phoenix.s3.amazonaws.com/assets/gf_images_flows_browser/gf_bar_handle_btn.svg"

	const container_element = $(`
		<div id='seeker_container'>
			<div id='seeker_bar'>

				<div id='seek_range_bar'>
					<canvas id='seek_range_bar_background'></canvas>
				</div>
				
				<div id='seek_range_bar_button'>
					<div id='button' class="gf_center">
						<div id='button_symbol'>
							<img src="${asset_uri__gf_bar_handle_btn_str}"></img>
						</div>
					</div>

					<div id='button_seek_info' class='gf_center' draggable='false'>
						<div id='seek_page_index'>
							${p_first_page_int}
						</div>
					</div>

					<div id='page_preview' draggable='false'>

					</div>
				</div>
			</div>
			<div id='conn'></div>
		</div>`)[0];
	$(p_parent_container).append(container_element);
	
	const seeker_bar_container_element = $(container_element).find("#seek_range_bar")[0];
	const seeker_bar_button_element    = $(container_element).find("#seek_range_bar_button")[0];

	const bar_canvas_ctx = init_range_bar_canvas(seeker_bar_container_element);

	const viz_props :GF_viz_props = {
		seeker_container_height_px: container_element.getBoundingClientRect().height,
		seeker_container_width_px:  container_element.getBoundingClientRect().width,
		seeker_range_bar_height:    seeker_bar_container_element.getBoundingClientRect().height,         
		bar_canvas_ctx:             bar_canvas_ctx,
	}

	init_seeking(seeker_bar_button_element,
		container_element,
		p_first_page_int,
		p_last_page_int,
		p_viz_group_reset_fun,
		viz_props);

	return container_element;
}

//-------------------------------------------------
function init_seeking(p_seeker_bar_button_element :HTMLElement,
	p_container_root      :HTMLElement,
	p_seek_start_page_int :number,
	p_seek_end_page_int   :number,
	p_viz_group_reset_fun :Function,
	p_viz_props           :GF_viz_props) {
	
	const button = $(p_seeker_bar_button_element).find("#button")[0];
	const conn   = $(p_container_root).find("#conn")[0];

	//------------------------------------------------------------
	function init_button() {

		var seek_percentage_int :number;

		const button_element           = $(p_seeker_bar_button_element).find("#button")[0];
		const button_seek_info_element = $(p_seeker_bar_button_element).find('#button_seek_info')[0];
		const seek_page_index_element  = $(p_seeker_bar_button_element).find("#seek_page_index")[0];
		
		var old_button_y_int :number; 

		//-------------------------------------------------
		function button__mousemove_handler_fun(p_initial_click_coord_int :number,
			p_move_event :any) {

			// button relative coordinate of the move
			// IMPORTANT!! - using pageY property of event instead of $(...).offset().top because pageY
			//               seems to have a much higher update frequency by the browser than the offset().top property
			//               and we need that frequency to be high for smooth animation of the button.
			const movement_delta = p_move_event.pageY - p_initial_click_coord_int;

			
			const button_new_y_int = old_button_y_int + movement_delta;

			seek_percentage_int = handle_user_seek_event(p_seek_start_page_int,
				p_seek_end_page_int,
				button_new_y_int,

				p_seeker_bar_button_element,
				button_element,
				button_seek_info_element,
				conn,
				seek_page_index_element,
				p_viz_props);
		}

		//-------------------------------------------------
		function get_btn_y_distance_to_bar(p_bar :HTMLElement, p_event :MouseEvent) {
			const bar_rect         = p_bar.getBoundingClientRect();
			const bar_global_y_int = bar_rect.top;

			const mouse_y_relative_to_btn_int = p_event.offsetY;
			const btn_global_y_int = p_event.clientY - mouse_y_relative_to_btn_int;
			const y_distance_int   = btn_global_y_int - bar_global_y_int;
			return y_distance_int;
		}

		//-------------------------------------------------
		// MOUSEDOWN

		const bar = $(p_container_root).find("#seek_range_bar")[0];

		$(button).on("mousedown", (p_event :JQueryEventObject)=>{
			const mouse_event = p_event.originalEvent as MouseEvent;

			
			// get the Y-axis distance of the button top to the bar top, and memorize this as the old button position
			const btn_y_distance_to_bar = get_btn_y_distance_to_bar(bar, mouse_event);
			old_button_y_int = btn_y_distance_to_bar;
			

			// initial button relative coordinate where the user clicked
			const initial_click_coord_int = p_event.pageY;

			// when user clicks and holds on the scroll button, a mouse move event handler is added to react to movement
			$(button).on("mousemove", (p_move_event)=>{
				
				button__mousemove_handler_fun(initial_click_coord_int,
					p_move_event);
			});
		});

		// MOUSEUP
		$(button).on("mouseup", async (p_event)=>{

			// stop handling button movement
			$(button).off("mousemove");

			//--------------
			// IMPORTANT!! - since random seeking is potentially an expensive operation,
			//               it is only initialized when the user lets go of the seek button.
			
			await seek_to_new_value(seek_percentage_int,
				p_seek_start_page_int,
				p_seek_end_page_int,
				p_viz_group_reset_fun);

			//--------------
			
			const conn_y = get_conn_y(conn, seek_percentage_int, p_viz_props);
			mark_seeked_page_on_bar(conn_y,
				p_viz_props);
		});
		$(button).on("mouseleave", (p_event)=>{

			// cancel movement of the button if mouse leaves it
			$(button).off("mousemove");
		});
	}

	//------------------------------------------------------------
	init_button();
}

//-------------------------------------------------
// USER_SEEK_EVENT

function handle_user_seek_event(p_seek_start_page_int :number,
	p_seek_end_page_int        :number,
	p_button_new_y_int         :number,
	p_seek_range_bar_button    :HTMLElement,
	p_button_element           :HTMLElement,
	p_button_seek_info_element :HTMLElement,
	p_conn_element             :HTMLElement,
	p_seek_page_index_element  :HTMLElement,
	p_viz_props                :GF_viz_props) {
	
	const button_height_int = $(p_button_element).height();
	
	const seeker_range_bar_height_px_int = p_viz_props.seeker_range_bar_height;

	// BOTTOM/TOP NOT REACHED
	if (p_button_new_y_int + button_height_int <= seeker_range_bar_height_px_int && 
		p_button_new_y_int >= 0) {

		$(p_seek_range_bar_button).css("top", `${p_button_new_y_int}px`);

		const seek_percentage_int = get_seek_percentage(p_button_element,
			seeker_range_bar_height_px_int,
			p_button_new_y_int);

		//-----------
		// CONNECTION POSITIONING
		
		const conn_y = get_conn_y(p_conn_element, seek_percentage_int, p_viz_props);
		$(p_conn_element).css("top", `${Math.floor(conn_y)}px`);

		//-----------
		// GET_SEEK_PAGE_INDEX
		const seek_page_index_int = get_seek_page_index(seek_percentage_int,
			p_seek_start_page_int,
			p_seek_end_page_int);
		
		//-----------

		if (p_button_seek_info_element != null && p_seek_page_index_element != null) {
			display_button_seek_info(p_seek_page_index_element,
				seek_page_index_int);
		}
		
		return seek_percentage_int;
	}
	// BOTTOM OR TOP REACHED
	else {
		// BOTTOM
		if (p_button_new_y_int + button_height_int <= seeker_range_bar_height_px_int) {
			return 100;
		} 
		// TOPs
		else if (p_button_new_y_int >= 0) {
			return 0;
		}
	}

	return 0;
}

//------------------------------------------------
async function seek_to_new_value(p_seek_percentage_int :number,
	p_seek_start_page_int :number,
	p_seek_end_page_int   :number,
	p_viz_group_reset_fun :Function) {
	
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const page_index_to_seek_to_int = get_seek_page_index(p_seek_percentage_int,
			p_seek_start_page_int,
			p_seek_end_page_int);

		// RESET
		p_viz_group_reset_fun(page_index_to_seek_to_int,
			// p_on_complete_fun
			()=>{
				console.log("viz_group reset done...")
				p_resolve_fun({});
			});
	});
	return p;
}

//------------------------------------------------
function mark_seeked_page_on_bar(p_conn_y_int :number,
	p_viz_props :GF_viz_props) {
	
	const ctx = p_viz_props.bar_canvas_ctx;
	ctx.rect(0, p_conn_y_int, 50, 1);
	ctx.stroke();
}

//------------------------------------------------
// UTILS
//-------------------------------------------------
function init_range_bar_canvas(p_container :HTMLElement) :CanvasRenderingContext2D {
		
	const bar_width_int  = $(p_container).width();
	const bar_height_int = $(p_container).height();

	const canvas  = <HTMLCanvasElement> $(p_container).find('#seek_range_bar_background')[0];
	canvas.width  = bar_width_int;
	canvas.height = bar_height_int;
	const ctx = canvas.getContext('2d');

	// first check if the canvas 2d context was returned by the browser
	if (ctx) {
		ctx.fillStyle   = "yellow";
		ctx.strokeStyle = "#7c7c7c";
	}
	if (!ctx) {
		throw new Error("Failed to get 2D context");
	}
	return ctx;
}

//------------------------------------------------
// this is called repeatedly as the user is seeking around
// p_seek_page_index - the index of the page that the user would seek to if he/she
//                     initiated the seek (released the seek button)

function display_button_seek_info(p_seek_page_index_element :HTMLElement,
	p_seek_page_index_int :number) {
	
	$(p_seek_page_index_element).text(`${p_seek_page_index_int}`);
}

//------------------------------------------------
function get_seek_page_index(p_seek_percentage_int :number,
	p_seek_start_page_int :number,
	p_seek_end_page_int   :number) :number {
	
	const total_range = p_seek_end_page_int - p_seek_start_page_int;
	
	// MATH EXPLANATION:
	// 1. 100 : total_range = p_seek_percentage_int : x
	// 2. 100 * x           = total_range * p_seek_percentage_int
	
	const page_index_delta_int = (total_range * p_seek_percentage_int) / 100;
	const page_index_to_seek_to_int = Math.floor(p_seek_start_page_int + page_index_delta_int);
	
	return page_index_to_seek_to_int;
}

//------------------------------------------------------------
function get_seek_percentage(p_button_element :HTMLElement,
	p_seeker_range_bar_height_px :number,
	p_btn_global_page_y_int      :number) :number {
	
	// MATH EXPLANATION:
	// 1. (p_seeker_range_bar_height_px - p_button_element.height) : 100 = p_btn_global_page_y_int : x
	// 2. x = (100 * p_btn_global_page_y_int) / (p_seeker_range_bar_height_px - p_button_element.height)
	
	const seek_percentage_int = (100 * p_btn_global_page_y_int) / 
		(p_seeker_range_bar_height_px - $(p_button_element).height());

	return seek_percentage_int;
}

//------------------------------------------------------------
function get_conn_y(p_conn_element :HTMLElement,
	p_seek_percentage_int :number,
	p_viz_props :GF_viz_props) :number {
	
	const seeker_range_bar_height_px_int = p_viz_props.seeker_range_bar_height;

	// 'conn' has negligable height and so its positioning is different from the buttom itself
	// 1. p_seeker_range_bar_height_px_int : 100 = x : seek_percentage
	// 2. p_seeker_range_bar_height_px_int * seek_percentage = 100 * x
	const conn_y = (seeker_range_bar_height_px_int * p_seek_percentage_int) / 100;
	$(p_conn_element).css("top", `${Math.floor(conn_y)}px`);

	return conn_y;

}