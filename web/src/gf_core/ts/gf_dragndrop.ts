/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

///<reference path="../../d/jquery.d.ts" />

//--------------------------------------------------------
export function init(p_target_element,
    p_on_dnd_event_fun,
    p_assets_paths_map) {


    const control = $(`
        <div class="dnd_handle">
            <img class="symbol" src="${p_assets_paths_map["gf_bar_handle_btn"]}"></img>
            <div class="overlay"></div>
        </div>`);

    const handle_height_int = $(p_target_element).outerHeight();
    var   handle_width_int  = $(control).outerWidth();
    $(control).css("height", `${handle_height_int}px`);
    $(control).css("opacity", `0.6`);

    // overlay - in place to be transparent and over all other handle elements,
    //           to simplify event and click coordinate calculation over the dnd_handle.
    $(control).find(".overlay").css("height", `${handle_height_int}px`);

    // MOUSE_ENTER
    $(p_target_element).on("mouseenter", ()=>{
        
        $(p_target_element).append(control);
        init_handle();

        handle_width_int = $(control).outerWidth();
    });

    // MOUSE_LEAVE
    $(p_target_element).on("mouseleave", ()=>{
                
        $(control).remove();
    });

    //--------------------------------------------------------
    function init_handle() {

        // IMPORTANT!!
        // these values indicate by how much the position of the target_element
        // have to be offset when they're moved around, to accound for the dimensions
        // and position of the movement handle and where the user clicked on that handle.
        var distance_to_target_origin_x;
        var distance_to_target_origin_y;

        //--------------------------------------------------------
        function mouse_move_fun(p_event) {

            const new_x = p_event.pageX + distance_to_target_origin_x;
            const new_y = p_event.pageY - distance_to_target_origin_y;

            $(p_target_element).css("left", `${new_x}px`);
            $(p_target_element).css("top", `${new_y}px`);
        }

        //--------------------------------------------------------
        // MOUSE_DOWN
        $(control).on("mousedown", (p_event)=>{

            distance_to_target_origin_x = handle_width_int - p_event.offsetX;
            distance_to_target_origin_y = p_event.offsetY;

            $(control).css("pointer", "grab");

            $("body").on("mousemove", mouse_move_fun);

            // EVENT
            p_on_dnd_event_fun("drag_start");
        });

        // MOUSE_UP
        $(control).on("mouseup", ()=>{
            
            $(control).css("pointer", "pointer");

            $("body").unbind("mousemove", mouse_move_fun);

            // EVENT
            p_on_dnd_event_fun("drag_stop");
        });
    }

    //--------------------------------------------------------
}