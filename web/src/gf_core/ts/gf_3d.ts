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


//-------------------------------------------------
// rotates the target div when the mouse moves, so that it always faces the cursor
// p_transform_constraint_coefficient_int - limits/dampens the rotation. the larger the number to smaller the rotation

export function div_follow_mouse(p_target_div_element, p_container_div_element,
    p_transform_constraint_coefficient_int) {

    function transforms(p_x, p_y) {

        let box          = p_target_div_element.getBoundingClientRect();
        let x_calculated = -(p_y - box.y - (box.height / 2)) / p_transform_constraint_coefficient_int;
        let y_calculated = (p_x - box.x - (box.width / 2)) / p_transform_constraint_coefficient_int;
        
        return `perspective(200px) rotateX(${x_calculated}deg) rotateY(${y_calculated}deg) `;
    };

    $(p_container_div_element).on("mousemove", (p_e)=>{

        const x = p_e.clientX;
        const y = p_e.clientY;

        window.requestAnimationFrame(()=>{
            p_target_div_element.style.transform = transforms(x, y);
        });
    });
}