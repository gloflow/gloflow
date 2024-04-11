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
export interface GF_random_access_viz_props {
    readonly seeker_container_height_px :number;
    readonly seeker_container_width_px  :number;
    readonly seeker_bar_width_px        :number;
    readonly seeker_range_bar_width     :number;
    readonly seeker_range_bar_height    :number;
    
    readonly seeker_range_bar_color_str :string;
    readonly assets_uris_map;
}

//-------------------------------------------------
export function init(p_first_page_int :number,
    p_last_page_int :number,
    p_viz_props     :GF_random_access_viz_props,
    p_viz_group_reset_fun) {



    const container_element = init_seeker_bar(p_first_page_int,
        p_last_page_int,
        p_viz_props,
        p_viz_group_reset_fun);
    return container_element;


}

//-------------------------------------------------
function init_seeker_bar(p_first_page_int :number,
    p_last_page_int :number,
    p_viz_props     :GF_random_access_viz_props,
    p_viz_group_reset_fun) {

    const asset_uri__gf_bar_handle_btn_str = p_viz_props.assets_uris_map["gf_bar_handle_btn"];

    // #seeker_container user-select: none;
    // IMPORTANT!! - so that the user doesnt accidently select the entire random_access ui element
    //               as if its contents was text to be selected/highlighted
    const container_element = $(`
		<div id='seeker_container' style='
            position: absolute;
            top: 0px;
            user-select: none;'>

            <div id='seeker_bar'>

				<div id='seek_range_bar'>
					<canvas id='seek_range_bar_background'></canvas>
				</div>
				
				<div id='seek_range_bar_button'>
					
					
					<div id='button' style='user-select: none;'>
						<div id='button_symbol' style='user-select: none; user-drag: none;'>
                            <img style='
                                width:100%;
                                user-select: none;
                                user-drag: none;'
                                src="${asset_uri__gf_bar_handle_btn_str}"></img>
                        </div>
					</div>

					<div id='button_seek_info'
                        style='
                            position:   absolute;
                            text-align: center;
                            font-size:  32px;
                            background-color: orange;'>

						<div id='seek_page_index'
                            style='
                                position: relative;
                                top:      12px;'>
                            ${p_first_page_int}
                        </div>
					</div>
				</div>
			</div>
            <div id='conn'></div>
		</div>`);

    const seeker_bar_element              = $(container_element).find('#seeker_bar');
    const seeker_range_bar_element        = $(container_element).find('#seek_range_bar');
    const seeker_range_bar_button_element = $(container_element).find('#seek_range_bar_button');
    const seek_page_index                 = $(container_element).find('#seek_page_index');


    $(seeker_bar_element).on('mouseenter', (p_event)=>{
        // $(seeker_range_bar_button_element).css('visibility', 'visible');
    });
    $(seeker_bar_element).on('mouseleave', (p_event)=>{
        // $(seeker_range_bar_button_element).css('visibility', 'hidden');
    });


    //------------
    // CSS
    
    // SEEKER CONTAINER
    $(container_element).css('height',   `${p_viz_props.seeker_container_height_px}px`);
    $(container_element).css('width',    `${p_viz_props.seeker_container_width_px}px`);
    
    // SEEKER BAR
    $(seeker_bar_element).css("position", 'absolute');
    $(seeker_bar_element).css("right",    '0px');
    $(seeker_bar_element).css("top",      '0px');
    $(seeker_bar_element).css("width",    `${p_viz_props.seeker_bar_width_px}px`);
    $(seeker_bar_element).css("height",   `${p_viz_props.seeker_container_height_px}px`);
    
    // SEEKER RANGER BAR
    $(seeker_range_bar_element).css("width",    `${p_viz_props.seeker_range_bar_width}px`);
    $(seeker_range_bar_element).css("height",   `${p_viz_props.seeker_range_bar_height}px`);
    $(seeker_range_bar_element).css("position", 'absolute');
    $(seeker_range_bar_element).css("right",    '0px');
    $(seeker_range_bar_element).css("background-color", "green"); // p_viz_props.seeker_range_bar_color_str);
    
    // seek_page_index
    $(seek_page_index).css("left", `${($(seek_page_index).width() - $(seek_page_index).width())/2}px`);

    //------------

    const seeker_bar_container_element = $(container_element).find("#seek_range_bar");
    init_range_bar_background_canvas(p_first_page_int,
        p_last_page_int,
        seeker_bar_container_element);

    const seeker_bar_button_element = $(container_element).find("#seek_range_bar_button");
    init_seeking(seeker_bar_button_element,
        container_element,
        p_first_page_int,
        p_last_page_int,
        p_viz_group_reset_fun,
        p_viz_props);

    //------------

    return container_element;
}

//-------------------------------------------------
function init_seeking(p_seeker_bar_button_element,
    p_container_root,
    p_seek_start_page_int :number,
    p_seek_end_page_int   :number,
    p_viz_group_reset_fun,
    p_viz_props           :GF_random_access_viz_props) {

    const button_width_px_int      = 60;
    const button_height_px_int     = 149;
    const button_conn_width_px_int = 40;
    const button_info_width_px     = 60;
    const button_info_height_px    = 60;
    const button_conn_right_px_int = 0;
    // const button_seek_info_y_offset_int = 15;

    const button           = $(p_seeker_bar_button_element).find("#button");
    const button_symbol    = $(button).find('#button_symbol'); 
    const button_seek_info = $(p_seeker_bar_button_element).find('#button_seek_info');
    const conn             = $(p_container_root).find("#conn");

    //------------------------------------------------------------
    function init_button() {


        $(p_seeker_bar_button_element).css("position", 'absolute');
        $(p_seeker_bar_button_element).css("top",      '0px');
        $(p_seeker_bar_button_element).css("right",    '0px');
        $(p_seeker_bar_button_element).css("width",    `${button_width_px_int}px`);
        $(p_seeker_bar_button_element).css("height",   `${button_height_px_int}px`);
        
        $(button).css("position", 'absolute');
        $(button).css("top",      '0px');
        $(button).css("right",    '0px');
		$(button).css("width",    `${button_width_px_int}px`);
        $(button).css("height",   `${button_height_px_int}px`);
		$(button).css("background-color", "gray");
		$(button).css("cursor",           'pointer');
		$(button).css("opacity",          '0.7');
        

        $(button_symbol).css("position", 'absolute');
		$(button_symbol).css("width",    '37px');
		$(button_symbol).css("height",   '18');
		$(button_symbol).css("top",      `${($(button).height()-$(button_symbol).height())/2}px`);
		$(button_symbol).css("left",     '12px');
    

        // the 'button' is on the left side of the 'conn'
	  	// (positioned relative to the right edge)
	  	const button_right_px = button_conn_right_px_int + button_conn_width_px_int;
	  	$(button).css("right", `${button_right_px}px`);


        

        $(button_seek_info).css("width",  `${button_info_width_px}px`);
	  	$(button_seek_info).css("height", `${button_info_height_px}px`);
	  	$(button_seek_info).css("top",    `${($(button).height()-$(button_seek_info).height())/2}px`); // `${button_seek_info_y_offset_int}px`);
	  	
	  	// the 'button' is on the left side of the 'conn'
	  	const button_seek_info_right_px = button_width_px_int + button_right_px;
	  	$(button_seek_info).css("right", `${button_seek_info_right_px}px`);
        


        $(button).on("mouseover", (p_event)=>{
            $(button).css("opacity", '1');
        });
        
        $(button).on("mouseout", (p_event)=>{
            $(button).css("opacity", '0.7');
        });


        var seek_percentage_int :number;

        const button_element           = $(p_seeker_bar_button_element).find("#button");
        const button_seek_info_element = $(p_seeker_bar_button_element).find('#button_seek_info');
        const seek_page_index_element  = $(p_seeker_bar_button_element).find("#seek_page_index");
        
        //-------------------------------------------------
        function button__mousemove_handler_fun(p_initial_click_coord_int,
            p_move_event) {


            // button relative coordinate of the move
            // IMPORTANT!! - using pageY property of event instead of $(...).offset().top because pageY
            //               seems to have a much higher update frequency by the browser than the offset().top property
            //               and we need that frequency to be high for smooth animation of the button.
            const movement_delta = p_move_event.pageY - p_initial_click_coord_int;


            

            // const old_button_y = $(button).offset().top; 
            const button_new_y_int = movement_delta; // old_button_y + movement_delta;



            console.log("-----", button_new_y_int, "init click coord", p_initial_click_coord_int)


            seek_percentage_int = handle_user_seek_event(p_seek_start_page_int,
                p_seek_end_page_int,
                button_new_y_int,

                p_seeker_bar_button_element,
                button_element,
                button_seek_info_element,
                conn,
                seek_page_index_element,

                p_viz_props,
                p_move_event);
        }

        //-------------------------------------------------
        // MOUSEDOWN
        $(button).on("mousedown", (p_event)=>{
	  			
            // initial button relative coordinate where the user clicked
            const initial_click_coord_int = p_event.pageY;

            // when user clicks and holds on the scroll button, a mouse move event handler is added 
            // to react to movement
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
			//               it is only initialized when the user lets go of the seek button
			// FIX!! - how is the mousewheel seeking handled? (no mouse seeking? since it cant
			//         be done on touch devices anyway)
			
			await seek_to_new_value(seek_percentage_int,
				p_seek_start_page_int,
				p_seek_end_page_int,
				p_viz_group_reset_fun);

			//--------------
        });
        $(button).on("mouseleave", (p_event)=>{

            // cancel movement of the button if mouse leaves it
            // $(button).off("mousemove");
        });
    }

    //------------------------------------------------------------
	function init_button_conn() {
		
		//------------
		// CSS
		$(conn).css('position',        'absolute');
		$(conn).css('height',          '1px');
		$(conn).css('width',           `${button_conn_width_px_int}px`);
		$(conn).css('backgroundColor', 'black');
		$(conn).css('top',             '0px');
		
		$(conn).css('right', `${button_conn_right_px_int}px`); // p_seeker_range_bar_element.style.width;

		//------------
	}

    //------------------------------------------------------------

    init_button_conn();
    init_button();
}

//-------------------------------------------------
// USER_SEEK_EVENT

function handle_user_seek_event(p_seek_start_page_int :number,
    p_seek_end_page_int   :number,
    p_button_new_y_int    :number,

    p_seek_range_bar_button,
    p_button_element,
    p_button_seek_info_element,
    p_conn_element,
    p_seek_page_index_element,

    p_viz_props :GF_random_access_viz_props,
    p_move_event) {
    
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
		
		// 'conn' has negligable height and so its positioning is different from the buttom itself
		// 1. p_seeker_range_bar_height_px_int : 100 = x : seek_percentage
		// 2. p_seeker_range_bar_height_px_int * seek_percentage = 100 * x
		const conn_y = (seeker_range_bar_height_px_int * seek_percentage_int) / 100;
		$(p_conn_element).css("top", `${Math.floor(conn_y)}px`);
		
		//-----------
        // GET_SEEK_PAGE_INDEX
		const seek_page_index_int = get_seek_page_index(seek_percentage_int,
			p_seek_start_page_int,
			p_seek_end_page_int);
		
        //-----------

		if (p_button_seek_info_element != null && p_seek_page_index_element != null) {

            const button_seek_info_y_offset_int = $(p_button_seek_info_element).offset().top;

			display_button_seek_info(p_button_seek_info_element,
				p_seek_page_index_element,
				seek_page_index_int,
				p_button_new_y_int,
				button_seek_info_y_offset_int);
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

	p_viz_group_reset_fun) {
    
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

        const page_index_to_seek_to_int = get_seek_page_index(p_seek_percentage_int,
            p_seek_start_page_int,
            p_seek_end_page_int);
        

        console.log("DDDDDDDDDDDDDDDDDDDDDD", page_index_to_seek_to_int)

        // RESET
        p_viz_group_reset_fun(page_index_to_seek_to_int,
            // p_on_complete_fun
            ()=>{

                console.log("DDDDDDD2222")
                p_resolve_fun({});
            });
    });
    return p;
}

//------------------------------------------------
// UTILS
//-------------------------------------------------
function init_range_bar_background_canvas(p_first_page_int :number,
    p_last_page_int :number,
    p_container) {
        
    const bar_width_int  = $(p_container).width();
    const bar_height_int = $(p_container).height();
    const pages_num_int  = Math.abs(p_last_page_int - p_first_page_int);
    const page_single_height_px_int = Math.floor(bar_height_int/pages_num_int)

    const canvas = <HTMLCanvasElement> $(p_container).find('#seek_range_bar_background')[0];
    canvas.width  = bar_width_int;
	canvas.height = bar_height_int;
    const ctx = canvas.getContext('2d');

    // first check if the canvas 2d context was returned by the browser
    if (ctx) {

        ctx.fillStyle   = "yellow";
        ctx.strokeStyle = "black";
        for (var i=0;i<pages_num_int;i++) {

            ctx.rect(0,
                i*page_single_height_px_int,
                bar_width_int,
                page_single_height_px_int);

            ctx.fill();
            ctx.stroke();
        }
    }
}

//------------------------------------------------
// this is called repeatedly as the user is seeking around
// p_seek_page_index - the index of the page that the user would seek to if he/she
//                     initiated the seek (released the seek button)

function display_button_seek_info(p_button_seek_info_element,
	p_seek_page_index_element,
	p_seek_page_index_int           :number,
	p_button_new_y_int              :number,
	p_button_seek_info_y_offset_int :number) { // how far from top of parent div container
	// p_button_seek_info_draw_fun) {
	
	$(p_seek_page_index_element).text(`${p_seek_page_index_int}`);
	$(p_seek_page_index_element).css("left",
        `${(p_button_seek_info_element.offsetWidth - p_seek_page_index_element.offsetWidth)/2}px`);
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
function get_seek_percentage(p_button_element,
	p_seeker_range_bar_height_px :number,
    p_btn_global_page_y_int      :number) :number {
	
	// MATH EXPLANATION:
	// 1. (p_seeker_range_bar_height_px - p_button_element.height) : 100 = p_btn_global_page_y_int : x
	// 2. x = (100 * p_btn_global_page_y_int) / (p_seeker_range_bar_height_px - p_button_element.height)
	
	const seek_percentage_int = (100 * p_btn_global_page_y_int) / 
	    (p_seeker_range_bar_height_px - $(p_button_element).height());

	return seek_percentage_int;
}