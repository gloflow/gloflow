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

///<reference path="../../d/jquery.d.ts" />

import * as gf_image_process from "./gf_image_process";

//--------------------------------------------------------
export function init_pallete(p_image) {

    var image_colors_shown_bool = false;
    $(p_image).on("mouseover", (p_event)=>{


        if (!image_colors_shown_bool) {
            
            const image        = p_event.target;
            const image_colors = gf_image_process.get_colors(image);

            console.log("dominant color - ", image_colors.color_hex_str);

            const color_info_element = $(`<div class="colors_info">
                <div class="color_dominant" style="background-color:#${image_colors.color_hex_str};"></div>
                <div class="color_pallete"></div>
            </div>`);


            // // IMPORTANT!! - change to color of the whole image_info control to match the dominant color of the
            // //               image its displaying.
            // $(p_image_info_element).css("background-color", `#${image_colors.color_hex_str}`);				
            


            color_info_element.insertAfter(image);

            //-------------
            // COLOR_PALLETE
            const color_pallete_element = $(color_info_element).find(".color_pallete");
            // const color_pallete_sub_lst = image_colors.color_palette_lst.slice(1, 6);
            image_colors.color_palette_lst.forEach((p_color_hex_str)=>{

                // console.log("-------------")
                // console.log(p_color_hex_str);


                const color_element = $(`<div class="color" style="background-color:#${p_color_hex_str};"></div>`);
                $(color_pallete_element).append(color_element);





                //-------------
                // COLOR_INSPECTOR
                var color_inspect_element = $(`<div class="color_inspect">
                    <div class='color_hex'>#${p_color_hex_str}</div>
                    <div class='color_large' style="background-color:#${p_color_hex_str};"></div>

                </div>`);
                $(color_element).on("mouseover", ()=>{
                    color_pallete_element.append(color_inspect_element);
                });
                $(color_element).on("mouseout", ()=>{
                    $(color_inspect_element).remove();
                });
                
                //-------------
            })

            //-------------
            const color_dominant_element = $(color_info_element).find(".color_dominant");

            var color_dominant_label_element = $(`<div class="color_dominant_label">color dominant</div>`);
            $(color_dominant_element).on("mouseover", ()=>{
                color_info_element.append(color_dominant_label_element);
            });
            $(color_dominant_element).on("mouseout", ()=>{
                $(color_dominant_label_element).remove();
            });


            var color_pallete_label_element = $(`<div class="color_pallete_label">color pallete</div>`);
            $(color_pallete_element).on("mouseover", ()=>{
                color_info_element.append(color_pallete_label_element);
            });
            $(color_pallete_element).on("mouseout", ()=>{
                $(color_pallete_label_element).remove();
            });
            


            image_colors_shown_bool = true;
        }
    });



}