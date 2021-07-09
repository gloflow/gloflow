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

declare var SVG;

//--------------------------------------------------------
export function run(p_width_int :number,
	p_height_int :number) {

    
    const canvas = SVG().addTo('#randomized_art')
        .size(p_width_int, p_height_int)
        

    const background_color = get_random_color();
    const background = canvas.rect(p_width_int, p_height_int)
        .fill({color: background_color, opacity: 1.0});


    const max_entities_num_int = 10;
    const entities_num_int = Math.floor(Math.random()*max_entities_num_int);
    for (var i=0;i<entities_num_int;i++) {

        const color_str = get_random_color();
        draw_entity(Math.random()*p_width_int, Math.random()*p_height_int, color_str, canvas);
    }
}

function get_random_color() :string {
	const random_r_int    :number = Math.floor(Math.random()*255);
	const random_g_int    :number = Math.floor(Math.random()*255);
	const random_b_int    :number = Math.floor(Math.random()*255);
	const random_rgba_str :string = `rgba(${random_r_int},${random_g_int},${random_b_int},${Math.floor(Math.random()*255)}`;
	return random_rgba_str;
}

//--------------------------------------------------------
function draw_entity(p_entity_x_int, p_entity_y_int, p_color, p_canvas) {

    const entity_scale_int    = 50;
    const entity_arms_num_int = 10;


    

    var rect = p_canvas.rect(entity_scale_int, entity_scale_int)
        .attr("x", p_entity_x_int)
        .attr("y", p_entity_y_int)
        .attr({ fill: p_color });


    var arm_previous_x_int = p_entity_x_int;
    var arm_previous_y_int = p_entity_y_int;
    var arms_coords_lst    = [];
    for (var i=0; i<entity_arms_num_int; i++) {

        const opacity_f = 1.0 - 0.5/entity_arms_num_int*i;
        console.log("=======", opacity_f)
        const [arm_x_int, arm_y_int] = draw_arm(arm_previous_x_int, arm_previous_y_int, opacity_f);
        arm_previous_x_int = arm_x_int;
        arm_previous_y_int = arm_y_int;


        arms_coords_lst.push([arm_x_int, arm_y_int])
    }
    
    //--------------------------------------------------------
    function draw_arm(p_previous_arm_x_int :number,
        p_previous_arm_y_int :number,
        p_arm_opacity_f      :number) {

        const arm_move__max_delta_int = 150;
        const arm_scale_int           = 5 + Math.floor(Math.random() * 5);

        // Math.random()-0.5 - produces a random number in the -.5-.5 range. 
        //                     essentially taking x left-right or y up-down randomly.

        var arm_x_int;
        var arm_y_int;

        while (true) {

            // reinitialize
            arm_x_int = p_previous_arm_x_int;
            arm_y_int = p_previous_arm_y_int;

            var move_along_x_bool;

            // randomly move arm on X/Y axis relative to previous arm
            if (Math.random() > 0.5) {
                move_along_x_bool = true;
                arm_x_int = p_previous_arm_x_int + Math.floor((Math.random()-0.5) * arm_move__max_delta_int);
            } else {
                move_along_x_bool = false;
                arm_y_int = p_previous_arm_y_int + Math.floor((Math.random()-0.5) * arm_move__max_delta_int);
            }

            // iterate until an arm is randomly placed outside the entity, to avoid cases where the arm
            // is randomly placed in the entity
            if (arm_outside_of_entity()) {
                break;
            }

            // console.log(arm_x_int, arm_y_int)
        }

        //--------------------------------------------------------
        function arm_outside_of_entity() {
            if (move_along_x_bool) {
                return (arm_x_int<p_entity_x_int || arm_x_int > p_entity_x_int+entity_scale_int ? true : false);
            } else {
                return (arm_y_int<p_entity_y_int || arm_y_int > p_entity_y_int+entity_scale_int ? true : false);
            }   
        }

        //--------------------------------------------------------
        
        // ARM_CONNECTION
        var line = p_canvas.line(p_previous_arm_x_int+arm_scale_int/2, p_previous_arm_y_int+arm_scale_int/2, arm_x_int+arm_scale_int/2, arm_y_int+arm_scale_int/2)
            .stroke({ color: `#000000`, width: 0.5 })
            .opacity(p_arm_opacity_f)

        // ARM
        var rect = p_canvas.rect(arm_scale_int, arm_scale_int)
            .attr("x", arm_x_int)
            .attr("y", arm_y_int)
            .attr({ fill: p_color })
            .opacity(p_arm_opacity_f);


        
        // PROXIMITY_REACTION
        if (arm_near_other_arms(arm_x_int, arm_y_int)) {

            const reactions_num_int = Math.random()*20;
            const reactions_max_distance_int = 50;
            const reaction_scale_int = 2;

            for (var i=0;i<reactions_num_int;i++) {
                const reaction_x_int = arm_x_int + Math.floor((Math.random()-0.5) * reactions_max_distance_int);
                const reaction_y_int = arm_y_int + Math.floor((Math.random()-0.5) * reactions_max_distance_int);
                
                var rect = p_canvas.rect(reaction_scale_int, reaction_scale_int)
                    .attr("x", reaction_x_int)
                    .attr("y", reaction_y_int)
                    .attr({ fill: '#000000' });

                var line = p_canvas.line(arm_x_int+arm_scale_int/2, arm_y_int+arm_scale_int/2, reaction_x_int+reaction_scale_int/2, reaction_y_int+reaction_scale_int/2)
                    .stroke({ color: '#00000022', width: 1 })
            }

        }

        //--------------------------------------------------------
        function arm_near_other_arms(p_arm_current_x_int, p_arm_current_y_int) {
        
            const distance_threshold_int = 5;
            for (var i=0; i<arms_coords_lst.length; i++) {

                const [arm_i_x_int, arm_i_y_int] = arms_coords_lst[i];

                if ((Math.abs(p_arm_current_x_int - arm_i_x_int)<distance_threshold_int) && (Math.abs(p_arm_current_y_int - arm_i_y_int)<distance_threshold_int)) {

                    console.log("---", p_arm_current_x_int - arm_i_x_int, p_arm_current_y_int - arm_i_y_int)
                    return true
                }

            }
            return false;
        }

        //--------------------------------------------------------

        
        return [arm_x_int, arm_y_int];
    }

    //--------------------------------------------------------
}