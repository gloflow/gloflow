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

// ///<reference path="../../d/jquery.d.ts" />

//--------------------------------------------------------
export function init(p_target_element :any,
    p_on_dnd_event_fun :any,
    p_assets_paths_map :any) {


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


    var control_side_str :string;
    var element_dragged_bool = false;

    // MOUSE_ENTER
    $(p_target_element).on("mouseenter", ()=>{
        
        $(p_target_element).append(control);
        init_handle();




        const taget_element_x_int = $(p_target_element).position().left;

        if (taget_element_x_int < 100) {
            $(control).css("right", "-40px");
            $(control).css("left", ""); // remove old

            control_side_str = "right";
        }
        else {
            $(control).css("left", "-40px");
            $(control).css("right", ""); // remove old

            control_side_str = "left";
        }

        handle_width_int = $(control).outerWidth();
    });

    // MOUSE_LEAVE
    $(p_target_element).on("mouseleave", ()=>{
        
        // dont remove the control if the dragging of the element
        // is in progress. otherwise its event handlers (including dragging)
        // will be removed and disturbed while user is dragging.
        if (!element_dragged_bool) {
            // $(control).remove();
        }
    });

    //--------------------------------------------------------
    function init_handle() {

        // IMPORTANT!!
        // these values indicate by how much the position of the target_element
        // have to be offset when they're moved around, to accound for the dimensions
        // and position of the movement handle and where the user clicked on that handle.
        var distance_to_target_origin_x :number;
        var distance_to_target_origin_y :number;

        //--------------------------------------------------------
        function mouse_move_fun(p_event :any) {

            // const new_x = p_event.pageX + distance_to_target_origin_x;

            // final coords of the element that was droped
            var new_x;
            switch (control_side_str) {
                case "left":
                    new_x = p_event.pageX + distance_to_target_origin_x;
                    break;
                
                case "right":
                    new_x = p_event.pageX - distance_to_target_origin_x;
                    break;
            }



            const new_y = p_event.pageY - distance_to_target_origin_y;

            $(p_target_element).css("left", `${new_x}px`);
            $(p_target_element).css("top", `${new_y}px`);
        }

        //--------------------------------------------------------
        // MOUSE_DOWN
        $(control).on("mousedown", (p_event)=>{

            switch (control_side_str) {
                case "left":
                    distance_to_target_origin_x = handle_width_int! - p_event.offsetX;
                    break;
                
                case "right":
                    const element_width_int = $(p_target_element).outerWidth()!;
                    distance_to_target_origin_x = element_width_int + handle_width_int! - p_event.offsetX;
                    break;
            }
            
            distance_to_target_origin_y = p_event.offsetY;

            $(control).css("pointer", "grab");

            // $("body").on("mousemove", mouse_move_fun);
            $(document).on("mousemove", mouse_move_fun);
            element_dragged_bool = true;
            
            // EVENT
            const data_map = {};
            p_on_dnd_event_fun("drag_start", data_map);
        });

        // MOUSE_UP
        // $(control).on("mouseup", (p_event)=>{
        $(document).on("mouseup", (p_event)=>{
            if (element_dragged_bool) {

                element_dragged_bool = false;

                $(control).css("pointer", "pointer");

                // $("body").unbind("mousemove", mouse_move_fun);
                $(document).unbind("mousemove", mouse_move_fun);


                // final coords of the element that was droped
                var final_x_int;
                switch (control_side_str) {
                    case "left":
                        final_x_int = p_event.pageX + distance_to_target_origin_x;
                        break;
                    
                    case "right":
                        final_x_int = p_event.pageX - distance_to_target_origin_x;
                        break;
                }




                const final_y_int = p_event.pageY - distance_to_target_origin_y;
                const data_map = {
                    "x_int": final_x_int,
                    "y_int": final_y_int,
                };

                // EVENT
                p_on_dnd_event_fun("drag_stop", data_map);
            }
        });
    }

    //--------------------------------------------------------
}