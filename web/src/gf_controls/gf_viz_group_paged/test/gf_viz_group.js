var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
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
System.register("ts/gf_viz_group_random_access", [], function (exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    //-------------------------------------------------
    function init(p_first_page_int, p_last_page_int, p_viz_props, p_viz_group_reset_fun) {
        const container_element = init_seeker_bar(p_first_page_int, p_last_page_int, p_viz_props, p_viz_group_reset_fun);
        return container_element;
    }
    exports_1("init", init);
    //-------------------------------------------------
    function init_seeker_bar(p_first_page_int, p_last_page_int, p_viz_props, p_viz_group_reset_fun) {
        const container_element = $(`
		<div id='seeker_container'>

            <div id='seeker_bar'>

				<div id='seek_range_bar'>
					<canvas id='seek_range_bar_background'></canvas>
				</div>
				
				<div id='seek_range_bar_button'>
					<div id='conn'></div>
					
					<div id='button'>
						<div id='button_symbol'>
                            <img style="width:100%;" src="./../../../../assets/gf_bar_handle_btn.svg"></img>
                        </div>
					</div>

					<div id='button_seek_info'>
						<div id='seek_page_index'>${p_first_page_int}</div>
					</div>
				</div>
			</div>
		</div>`);
        const seeker_bar_element = $(container_element).find('#seeker_bar');
        const seeker_range_bar_element = $(container_element).find('#seek_range_bar');
        const seeker_range_bar_button_element = $(container_element).find('#seek_range_bar_button');
        const seek_page_index = $(container_element).find('#seek_page_index');
        $(seeker_bar_element).on('mouseenter', (p_event) => {
            $(seeker_range_bar_button_element).css('visibility', 'visible');
        });
        $(seeker_bar_element).on('mouseleave', (p_event) => {
            $(seeker_range_bar_button_element).css('visibility', 'hidden');
        });
        //------------
        // CSS
        // SEEKER CONTAINER
        $(container_element).css('position', 'absolute');
        $(container_element).css('overflow', 'hidden');
        $(container_element).css('height', `${p_viz_props.seeker_container_height_px}px`);
        $(container_element).css('width', `${p_viz_props.seeker_container_width_px}px`);
        $(container_element).css('top', '0px');
        // IMPORTANT!! - so that the user doesnt accidently select the entire random_access ui element
        //               as if its contents was text to be selected/highlighted
        $(container_element).css('user-select', 'none');
        // SEEKER BAR
        $(seeker_bar_element).css("position", 'absolute');
        $(seeker_bar_element).css("right", '0px');
        $(seeker_bar_element).css("top", '0px');
        $(seeker_bar_element).css("height", `${p_viz_props.seeker_container_height_px}px`);
        $(seeker_bar_element).css("width", `${p_viz_props.seeker_bar_width_px}px`);
        // SEEKER RANGER BAR
        $(seeker_range_bar_element).css("height", `${p_viz_props.seeker_container_height_px}px`);
        $(seeker_range_bar_element).css("width", `${p_viz_props.seeker_range_bar_width}px`);
        $(seeker_range_bar_element).css("position", 'absolute');
        $(seeker_range_bar_element).css("right", '0px');
        $(seeker_range_bar_element).css("background-color", "green"); // p_viz_props.seeker_range_bar_color_str);
        // seek_page_index
        $(seek_page_index).css("left", `${($(seek_page_index).width() - $(seek_page_index).width()) / 2}px`);
        //------------
        const seeker_bar_container_element = $(container_element).find("#seek_range_bar");
        init_range_bar_background_canvas(p_first_page_int, p_last_page_int, seeker_bar_container_element);
        const seeker_bar_button_element = $(container_element).find("#seek_range_bar_button");
        init_seeking(seeker_bar_button_element, p_first_page_int, p_last_page_int, p_viz_group_reset_fun, p_viz_props);
        //------------
        return container_element;
    }
    //-------------------------------------------------
    function init_seeking(p_container_root, p_seek_start_page_int, p_seek_end_page_int, p_viz_group_reset_fun, p_viz_props) {
        const button_width_px_int = 60;
        const button_height_px_int = 149;
        const button_conn_width_px_int = 40;
        const button_info_width_px = 60;
        const button_info_height_px = 60;
        const button_conn_right_px_int = 0;
        const button_seek_info_y_offset_int = 15;
        const seek_range_bar_button = $(p_container_root).find("#seek_range_bar_button");
        const button = $(p_container_root).find("#button");
        const button_symbol = $(button).find('#button_symbol');
        const button_seek_info = $(p_container_root).find('#button_seek_info');
        const conn = $(p_container_root).find("#conn");
        //------------------------------------------------------------
        function init_button() {
            $(seek_range_bar_button).css("position", 'absolute');
            $(seek_range_bar_button).css("top", '0px');
            $(seek_range_bar_button).css("right", '0px');
            $(seek_range_bar_button).css("width", `${button_width_px_int}px`);
            $(seek_range_bar_button).css("height", `${button_height_px_int}px`);
            $(button).css("position", 'absolute');
            $(button).css("top", '0px');
            $(button).css("right", '0px');
            $(button).css("width", `${button_width_px_int}px`);
            $(button).css("height", `${button_height_px_int}px`);
            $(button).css("background-color", "gray");
            $(button).css("cursor", 'pointer');
            $(button).css("opacity", '0.7');
            $(button_symbol).css("position", 'absolute');
            $(button_symbol).css("width", '37px');
            $(button_symbol).css("height", '37px');
            $(button_symbol).css("top", '54px');
            $(button_symbol).css("left", '14px');
            // $(button_symbol).css("background-image", `url('${p_button_atlas_url_str}')`);
            // $(button_symbol).css("background-position", "38px 59px"); 
            $(button_symbol).find("img").css("width");
            // the 'button' is on the left side of the 'conn'
            // (positioned relative to the right edge)
            const button_right_px = button_conn_right_px_int + button_conn_width_px_int;
            $(button).css("right", `${button_right_px}px`);
            $(button_seek_info).css("position", 'absolute');
            $(button_seek_info).css("width", `${button_info_width_px}px`);
            $(button_seek_info).css("height", `${button_info_height_px}px`);
            $(button_seek_info).css("background-color", "orange");
            $(button_seek_info).css("top", `${button_seek_info_y_offset_int}px`);
            // the 'button' is on the left side of the 'conn'
            const button_seek_info_right_px = button_width_px_int + button_right_px;
            $(button_seek_info).css("right", `${button_seek_info_right_px}px`);
            $(button).on("mouseover", (p_event) => {
                $(button).css("opacity", '1');
            });
            $(button).on("mouseout", (p_event) => {
                $(button).css("opacity", '0.7');
            });
            var seek_percentage_int;
            const button_element = $(p_container_root).find("#button");
            const button_seek_info_element = $(p_container_root).find('#button_seek_info');
            const conn_element = $(p_container_root).find("#conn");
            const seek_page_index_element = $(p_container_root).find("#seek_page_index");
            //-------------------------------------------------
            function button__mousemove_handler_fun(p_initial_click_coord, p_move_event) {
                // button relative coordinate of the move
                // IMPORTANT!! - using pageY property of event instead of $(...).offset().top because pageY
                //               seems to have a much higher update frequency by the browser than the offset().top property
                //               and we need that frequency to be high for smooth animation of the button.
                const movement_delta = p_move_event.pageY - p_initial_click_coord;
                // const old_button_y = $(button).offset().top; 
                const button_new_y_int = movement_delta; // old_button_y + movement_delta;
                console.log("**");
                seek_percentage_int = handle_user_seek_event(p_container_root, p_seek_start_page_int, p_seek_end_page_int, button_new_y_int, seek_range_bar_button, button_element, button_seek_info_element, conn_element, seek_page_index_element, p_viz_props, p_move_event);
                /*seek_percentage_int = handle_user_seek_event(button_new_y_int,
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
                    p_button_seek_info_draw_fun);*/
            }
            //-------------------------------------------------
            // MOUSEDOWN
            $(button).on("mousedown", (p_event) => {
                // initial button relative coordinate where the user clicked
                const initial_click_coord = p_event.pageY;
                // when user clicks and holds on the scroll button, a mouse move event handler is added 
                // to react to movement
                $(button).on("mousemove", (p_move_event) => {
                    console.log("=====");
                    button__mousemove_handler_fun(initial_click_coord, p_move_event);
                });
            });
            // MOUSEUP
            $(button).on("mouseup", (p_event) => {
                $(button).off("mousemove", button__mousemove_handler_fun);
                //--------------
                // IMPORTANT!! - since random seeking is potentially an expensive operation,
                //               it is only initialized when the user lets go of the seek button
                // FIX!! - how is the mousewheel seeking handled? (no mouse seeking? since it cant
                //         be done on touch devices anyway)
                seek_to_new_value(seek_percentage_int, p_seek_start_page_int, p_seek_end_page_int, p_viz_group_reset_fun, () => {
                });
                //--------------
            });
            $(button).on("mouseleave", (p_event) => {
                $(button).off("mousemove", button__mousemove_handler_fun);
            });
        }
        //------------------------------------------------------------
        function init_button_conn() {
            //------------
            // CSS
            $(conn).css('position', 'absolute');
            $(conn).css('height', '1px');
            $(conn).css('width', `${button_conn_width_px_int}px`);
            $(conn).css('backgroundColor', 'black');
            $(conn).css('top', '0px');
            $(conn).css('right', `${button_conn_right_px_int}px`); // p_seeker_range_bar_element.style.width;
            //------------
        }
        //------------------------------------------------------------
        init_button_conn();
        init_button();
    }
    //------------------------------------------------
    function seek_to_new_value(p_seek_percentage_int, p_seek_start_page_int, p_seek_end_page_int, p_viz_group_reset_fun, p_on_complete_fun) {
        const page_index_to_seek_to_int = get_seek_page_index(p_seek_percentage_int, p_seek_start_page_int, p_seek_end_page_int);
        p_viz_group_reset_fun(page_index_to_seek_to_int, 
        // p_on_complete_fun
        () => {
            p_on_complete_fun();
        });
    }
    //-------------------------------------------------
    function init_range_bar_background_canvas(p_first_page_int, p_last_page_int, p_container) {
        const bar_width_int = $(p_container).width();
        const bar_height_int = $(p_container).height();
        const pages_num_int = Math.abs(p_last_page_int - p_first_page_int);
        const page_single_height_px_int = Math.floor(bar_height_int / pages_num_int);
        const canvas = $(p_container).find('#seek_range_bar_background')[0];
        canvas.width = bar_width_int;
        canvas.height = bar_height_int;
        const ctx = canvas.getContext('2d');
        /*// background
        ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
        ctx.fillStyle = "blue";
        ctx.fillRect(0,
            0,
            bar_width_int,
            bar_height_int);*/
        ctx.fillStyle = "yellow";
        ctx.strokeStyle = "black";
        for (var i = 0; i < pages_num_int; i++) {
            ctx.rect(0, i * page_single_height_px_int, bar_width_int, page_single_height_px_int);
            ctx.fill();
            ctx.stroke();
        }
    }
    //-------------------------------------------------
    function init_button_seek_info_background_canvas() {
    }
    //-------------------------------------------------
    function handle_user_seek_event(p_container_root, p_seek_start_page_int, p_seek_end_page_int, p_button_new_y_int, p_seek_range_bar_button, p_button_element, p_button_seek_info_element, p_conn_element, p_seek_page_index_element, p_viz_props, p_move_event) {
        const button_height_int = $(p_button_element).height();
        const seeker_range_bar_height_px_int = p_viz_props.seeker_range_bar_height;
        // BOTTOM/TOP NOT REACHED
        if (p_button_new_y_int + button_height_int <= seeker_range_bar_height_px_int &&
            p_button_new_y_int >= 0) {
            console.log("---", p_button_new_y_int);
            // $(p_button_element).css("top", `${p_button_new_y_int}px`);
            $(p_seek_range_bar_button).css("top", `${p_button_new_y_int}px`);
            const seek_percentage_int = get_seek_percentage(p_button_element, seeker_range_bar_height_px_int, p_move_event.pageY);
            //-----------
            // CONN POSITIONING
            // 'conn' has negligable height and so its positioning is different from the buttom itself
            // 1. p_seeker_range_bar_height_px_int : 100 = x : seek_percentage
            // 2. p_seeker_range_bar_height_px_int * seek_percentage = 100 * x
            const conn_y = (seeker_range_bar_height_px_int * seek_percentage_int) / 100;
            $(p_conn_element).css("top", `${Math.floor(conn_y)}px`);
            //-----------
            const seek_page_index_int = get_seek_page_index(seek_percentage_int, p_seek_start_page_int, p_seek_end_page_int);
            if (p_button_seek_info_element != null && p_seek_page_index_element != null) {
                const button_seek_info_y_offset_int = $(p_button_seek_info_element).offset().top;
                display_button_seek_info(p_button_seek_info_element, p_seek_page_index_element, seek_page_index_int, p_button_new_y_int, button_seek_info_y_offset_int);
                // p_button_seek_info_draw_fun);
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
    // this is called repeatedly as the user is seeking around
    // p_seek_page_index - the index of the page that the user would seek to if he/she
    //                     initiated the seek (released the seek button)
    function display_button_seek_info(p_button_seek_info_element, p_seek_page_index_element, p_seek_page_index_int, p_button_new_y_int, p_button_seek_info_y_offset_int) {
        // p_button_seek_info_draw_fun) {
        $(p_seek_page_index_element).text(`${p_seek_page_index_int}`);
        $(p_seek_page_index_element).css("left", `${(p_button_seek_info_element.offsetWidth - p_seek_page_index_element.offsetWidth) / 2}px`);
        // $(p_button_seek_info_element).css("top", `${p_button_new_y_int+p_button_seek_info_y_offset_int}px`);
        /*//--------------
        // IMPORTANT!! - user-passed callback execution
        
        if (p_button_seek_info_draw_fun != null){
            p_button_seek_info_draw_fun();
        }
    
        //--------------*/
    }
    //------------------------------------------------
    function get_seek_page_index(p_seek_percentage_int, p_seek_start_page_int, p_seek_end_page_int) {
        const total_range = p_seek_end_page_int - p_seek_start_page_int;
        // 1. 100 : total_range = p_seek_percentage_int : x
        // 2. 100 * x           = total_range * p_seek_percentage_int
        const page_index_delta_int = (total_range * p_seek_percentage_int) / 100;
        const page_index_to_seek_to_int = Math.floor(p_seek_start_page_int + page_index_delta_int);
        console.log("----===", total_range, p_seek_percentage_int, page_index_delta_int, page_index_to_seek_to_int);
        return page_index_to_seek_to_int;
    }
    //------------------------------------------------------------
    function get_seek_percentage(p_button_element, p_seeker_range_bar_height_px, p_global_page_y_int) {
        // MATH EXPLANATION:
        // 1. (p_seeker_range_bar_height_px - p_button_element.height) : 100 = p_global_page_y_int : x
        // 2. x = (100 * p_global_page_y_int) / (p_seeker_range_bar_height_px - p_button_element.height)
        const seek_percentage_int = (100 * p_global_page_y_int) /
            (p_seeker_range_bar_height_px - $(p_button_element).height());
        console.log("###", p_global_page_y_int, p_seeker_range_bar_height_px, $(p_button_element).height(), seek_percentage_int);
        return seek_percentage_int;
    }
    return {
        setters: [],
        execute: function () {/*
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
        }
    };
});
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
System.register("ts/gf_viz_group_paged", ["ts/gf_viz_group_random_access"], function (exports_2, context_2) {
    "use strict";
    var gf_viz_group_random_access;
    var __moduleName = context_2 && context_2.id;
    //-------------------------------------------------
    function init(p_id_str, p_parent_id_str, p_elements_lst, p_initial_page_int, p_element_create_fun, p_elements_page_get_fun, p_create_initial_elements_bool = true) {
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
        $(container).find('img').on('load', () => {
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
        const viz_props = {
            seeker_container_height_px: 500,
            seeker_container_width_px: 100,
            seeker_bar_width_px: 50,
            seeker_range_bar_width: 30,
            seeker_range_bar_height: 500,
            seeker_range_bar_color_str: "red"
        };
        const start_page_int = 0;
        const end_page_int = 20;
        const random_access__container_element = gf_viz_group_random_access.init(start_page_int, end_page_int, viz_props, 
        //-------------------------------------------------
        // p_viz_group_reset_fun
        (p_start_page_int, p_on_complete_fun) => {
            reset_with_new_start_pages(container, p_start_page_int, p_element_create_fun, p_elements_page_get_fun);
            p_on_complete_fun();
        });
        //-------------------------------------------------
        // position seeker on the far right
        $(random_access__container_element).css("position", "absolute");
        $(random_access__container_element).css("right", "0px");
        $(container).append(random_access__container_element);
        //------------------------
        // LOAD_PAGES_ON_SCROLL
        var current_page_int = p_initial_page_int; // the few initial pages are already statically embedded in the document
        var page_is_loading_bool = false;
        const pages_container = $(container).find("#items");
        window.onscroll = () => __awaiter(this, void 0, void 0, function* () {
            // $(document).height() - height of the HTML document
            // window.innerHeight   - Height (in pixels) of the browser window viewport including, if rendered, the horizontal scrollbar
            if (window.scrollY >= $(document).height() - (window.innerHeight + 50)) {
                // IMPORTANT!! - only load 1 page at a time
                if (!page_is_loading_bool) {
                    yield load_new_pages(current_page_int, pages_container, p_element_create_fun, p_elements_page_get_fun);
                    page_is_loading_bool = false;
                    current_page_int += 1;
                    $(container).data("current_page", current_page_int);
                }
            }
        });
        //------------------------
    }
    exports_2("init", init);
    //-------------------------------------------------
    function reset_with_new_start_pages(p_container, p_start_page_int, // this is where it was seeked to, and is different from first_page/last_page
    p_element_create_fun, p_elements_page_get_fun, 
    // this is an initial load of viz_group, so load some >1 number of pages
    // starting from the page where to user seeked to.
    p_pages_to_get_num_int = 6) {
        return __awaiter(this, void 0, void 0, function* () {
            // remove all items currently displayed by viz_group, 
            // since new ones have to be shown
            $(p_container).find("#items .item").remove();
            const pages_container = $(p_container).find("#items");
            yield load_new_pages(p_start_page_int, pages_container, p_element_create_fun, p_elements_page_get_fun, p_pages_to_get_num_int);
        });
    }
    //-------------------------------------------------
    function enable_draggability(p_container, p_packery_grid) {
        const draggable_items_lst = $(p_container).find(".item");
        $(draggable_items_lst).each((p_i, p_e) => {
            const element = p_e; // $(this)[0];
            console.log(element);
            var draggie = new Draggabilly(element, {
            // options...
            });
            p_packery_grid.packery('bindDraggabillyEvents', draggie);
        });
    }
    //-------------------------------------------------
    function load_new_pages(p_page_index_int, p_pages_container, p_element_create_fun, p_elements_page_get_fun, p_pages_to_get_num_int = 1) {
        return __awaiter(this, void 0, void 0, function* () {
            // fetch page
            const elements_lst = yield p_elements_page_get_fun(p_page_index_int, p_pages_to_get_num_int);
            // create elements
            for (let element_map of elements_lst) {
                const element = p_element_create_fun(element_map);
                $(element).addClass("item");
                $(element).css("visibility", "hidden"); // initially elements are not visible until they load
                $(p_pages_container).find("#items").append(element);
                console.log(element);
                $(element).find('img').on('load', function () {
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
        });
    }
    return {
        setters: [
            function (gf_viz_group_random_access_1) {
                gf_viz_group_random_access = gf_viz_group_random_access_1;
            }
        ],
        execute: function () {/*
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
        }
    };
});
///<reference path="../../../d/jquery.d.ts" />
System.register("test/gf_viz_group_test", ["ts/gf_viz_group_paged"], function (exports_3, context_3) {
    "use strict";
    var gf_viz_group_paged;
    var __moduleName = context_3 && context_3.id;
    return {
        setters: [
            function (gf_viz_group_paged_1) {
                gf_viz_group_paged = gf_viz_group_paged_1;
            }
        ],
        execute: function () {
            ///<reference path="../../../d/jquery.d.ts" />
            //-------------------------------------------------
            $(document).ready(() => {
                const test_elements_lst = [
                    // page 1 (10 items)
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/b1b448df22b2767a8769f644f5f9e719_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/b1b448df22b2767a8769f644f5f9e719_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/b1b448df22b2767a8769f644f5f9e719_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    // page 2 (10 items)
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                    {
                        "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                    },
                ];
                //-------------------------------------------------
                function element_create_fun(p_element_map) {
                    const img_url_str = p_element_map["img_url_str"];
                    // console.log(img_url_str)
                    // console.log(p_element_container);
                    const element = $(`<div><img src='${img_url_str}'></img></div>`);
                    return element;
                }
                //-------------------------------------------------
                function elements_page_get_fun(p_page_index_int, p_pages_to_get_num_int) {
                    const p = new Promise(function (p_resolve_fun, p_reject_fun) {
                        const page_elements_lst = [
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                            {
                                "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                            },
                        ];
                        p_resolve_fun(page_elements_lst);
                    });
                    return p;
                }
                //-------------------------------------------------
                const id_str = "test_viz_group";
                const parent_id_str = "test_parent";
                // number of initial pages that are supplied to gf_viz_group to display
                // before it has to initiate its own page fetching logic.
                const initial_pages_num_int = 2;
                gf_viz_group_paged.init(id_str, parent_id_str, test_elements_lst, initial_pages_num_int, element_create_fun, elements_page_get_fun);
            });
        }
    };
});
