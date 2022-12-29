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



///<reference path="../../../../../../gloflow/web/src/d/jquery.d.ts" />



declare var _;
declare var THREE;
declare var iro;

//-------------------------------------------------
export function init(p_engine_api_map) {



    draw_helpers(p_engine_api_map);

    camera_animation(p_engine_api_map);


    
}

//-------------------------------------------------
function camera_animation(p_engine_api_map) {

    const container = $(`
        <div id="camera_animate">
            <div id="record_btn">
                start
            </div>
            <div id="play_btn">play</div>
            <div id="animations_label">current animations</div>
            <div id="animations"></div>
        </div>`);
    $("body").append(container);


    //-------------------------------------------------
    function view_animation_params(p_animation_map) {
        $(container).find("#animations_label").css("visibility", "visible");
        const animation_container = $(`
            <div class="animation">
                <input id="duration_sec_f" value="${p_animation_map["duration_sec_f"]}"></input>
                <div class="start">
                    x<span>${_.truncate(p_animation_map["start_props_map"].x, {"length": 8})}</span> |
                    y<span>${_.truncate(p_animation_map["start_props_map"].y, {"length": 8})}</span> |
                    z<span>${_.truncate(p_animation_map["start_props_map"].z, {"length": 8})}</span>
                </div>
                <div class="end">
                    x<span>${_.truncate(p_animation_map["end_props_map"].x, {"length": 8})}</span> |
                    y<span>${_.truncate(p_animation_map["end_props_map"].y, {"length": 8})}</span> |
                    z<span>${_.truncate(p_animation_map["end_props_map"].z, {"length": 8})}</span>
                </div>
            </div>
        `);
        $(container).find("#animations").append(animation_container);

        $(animation_container).find("input").on("change", (p_e)=>{

            console.log(p_e, )
            const new_animation_duration_sec_f = parseFloat($(p_e.target).val());
            p_animation_map["duration_sec_f"] = new_animation_duration_sec_f;
        });
    }

    //-------------------------------------------------

    // RECORD_CAMERA_POSITION
    var recording_bool = false;
    var camera__start_props_map;
    var camera__end_props_map;
    var animations_lst = [];

    // RECORD_BTN
    $(container).find("#record_btn").on('click', ()=>{

        // START
        if (!recording_bool) {
            camera__start_props_map = p_engine_api_map["camera__get_props_fun"]();

            // background color
            // FIX!! - setting this on camera_props for testing purposes, this should
            //         have its own global state map (not directly camera related)
            const color_background_rgb_lst = p_engine_api_map["get_state_fun"]("color_background");
            camera__start_props_map["color_background_r"] = color_background_rgb_lst[0];
            camera__start_props_map["color_background_g"] = color_background_rgb_lst[1];
            camera__start_props_map["color_background_b"] = color_background_rgb_lst[2];

            $(container).find("#record_btn").text("stop");
            recording_bool = true;
        }

        // STOP
        else {
            camera__end_props_map = p_engine_api_map["camera__get_props_fun"]();

            // background color
            // FIX!! - setting this on camera_props for testing purposes, this should
            //         have its own global state map (not directly camera related)
            const color_background_rgb_lst = p_engine_api_map["get_state_fun"]("color_background");
            camera__end_props_map["color_background_r"] = color_background_rgb_lst[0];
            camera__end_props_map["color_background_g"] = color_background_rgb_lst[1];
            camera__end_props_map["color_background_b"] = color_background_rgb_lst[2];

            //---------------------------
            // ANIMATION_RECORD
            const animation_map = {
                "duration_sec_f":  5,
                "start_props_map": camera__start_props_map,
                "end_props_map":   camera__end_props_map,
                "repeat_bool":     false,
            }
            animations_lst.push(animation_map);

            //---------------------------

            view_animation_params(animation_map);

            $(container).find("#record_btn").text("start");
            recording_bool = false;
        }
    });

    // PLAY_BTN
    $(container).find("#play_btn").on('click', ()=>{

        p_engine_api_map["camera__animate_fun"](animations_lst);
    });

    return container;
}

//-------------------------------------------------
function draw_helpers(p_engine_api_map) {
    
    const container = $(`
        <div id="helpers">
            <div id="origins_btn">origins</div>
            <div id="axes_btn">axes</div>
            <div id="grid_btn">grid</div>
            <div id="background_color_picker">
                <div id="control"></div>
                <div id="picked_color">#000000</div>
            </div>
        </div>`);
    $("body").append(container);

    //-------------------------------------------------
    // AXES_HELPER
    function draw_axis(p_x, p_y, p_z) {
        const axes_gr = new THREE.Group();
        var arrow_pos = new THREE.Vector3(p_x, p_y, p_z);
        axes_gr.add( new THREE.ArrowHelper( new THREE.Vector3( 1,0,0 ), arrow_pos, 6, 0x7F2020, 1, 0.5 ) );
        axes_gr.add( new THREE.ArrowHelper( new THREE.Vector3( 0,1,0 ), arrow_pos, 6, 0x207F20, 1, 0.5 ) );
        axes_gr.add( new THREE.ArrowHelper( new THREE.Vector3( 0,0,1 ), arrow_pos, 6, 0x20207F, 1, 0.5 ) );
        p_engine_api_map["scene_3d"].add(axes_gr);
        return axes_gr;
    }
    //-------------------------------------------------

    const origin_axes_gr = draw_axis(0, 0, 0);

    //---------------------------
    // GRID_HELPER
    const size = 30;
    const divisions = 30;

    const grid_helper = new THREE.GridHelper(size, divisions, 0x000000);
    p_engine_api_map["scene_3d"].add(grid_helper);




    // ORIGINS_TOGGLE
    var origins_visible_bool = true;
    var origins_lst = [];
    $(container).find("#origins_btn").on('click', ()=>{
        if (origins_visible_bool) {

            for (const origin of origins_lst) {
                p_engine_api_map["scene_3d"].remove(origin);
            }
            origins_visible_bool=false;
        } else {
            const coord_origins_stack_lst = p_engine_api_map["coord_origins_stack_lst"];

            for (const origin_v3 of coord_origins_stack_lst) {
                const origin_axis_gr = draw_axis(origin_v3.x, origin_v3.y, origin_v3.z)
                p_engine_api_map["scene_3d"].add(origin_axis_gr);
            }
            origins_visible_bool=true;
        }
    })

    // AXES_TOGGLE
    var axes_visible_bool = true;
    $(container).find("#axes_btn").on('click', ()=>{
        if (axes_visible_bool) {
            p_engine_api_map["scene_3d"].remove(origin_axes_gr);
            axes_visible_bool=false;
        } else {
            p_engine_api_map["scene_3d"].add(origin_axes_gr);
            axes_visible_bool=true;
        }
    })

    // GRID_TOGGLE
    var grid_visible_bool = true;
    $(container).find("#grid_btn").on('click', ()=>{
        if (grid_visible_bool) {
            p_engine_api_map["scene_3d"].remove(grid_helper);
            grid_visible_bool=false;
        } else {
            p_engine_api_map["scene_3d"].add(grid_helper);
            grid_visible_bool=true;
        }
    })

    

    //---------------------------
    // COLOR_PICKER
    const color_picker_color_element = $(container).find("#picked_color");
    const color_picker = new iro.ColorPicker("#background_color_picker #control", {
        width: 100,   // size of the picker
        color: "#f00" // initial color
    });
    color_picker.on('color:change', (p_color)=>{
        
        const picked_color_hex_str = p_color.hexString;

        $(color_picker_color_element).text(picked_color_hex_str);
        p_engine_api_map["set_state_fun"]({
            "color_background": picked_color_hex_str
        });
    });

    //---------------------------

}