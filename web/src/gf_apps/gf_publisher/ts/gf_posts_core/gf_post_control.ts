/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

import * as gf_image_colors from "./../../../../gf_core/ts/gf_image_colors";

//---------------------------------------------------
export function create(p_log_fun) {

}

//---------------------------------------------------
// INIT_EXISTING_DOM
/*
used for templates usually, where the image element DOM structure is already
created server side when loaded into the browser, and just needs to be initialized
(no creation of the DOM tree for the image control)
*/

export function init_existing_dom(p_post_element, p_log_fun) {


    init_posts_img_num(p_post_element);

    //----------------------
    // IMAGE_PALLETE
    const img = $(p_post_element).find("img")[0];

    const assets_paths_map = {
        "copy_to_clipboard_btn": "/images/static/assets/gf_copy_to_clipboard_btn.svg",
    }
    gf_image_colors.init_pallete(img,
        assets_paths_map,
        (p_color_dominant_hex_str,
        p_colors_hexes_lst)=>{

            // set the background color of the post to its dominant color
            $(p_post_element).css("background-color", `#${p_color_dominant_hex_str}`);

        });

    //----------------------
}





//---------------------------------------------------
function init_posts_img_num(p_post_element) {

    
    const post_images_number = $(p_post_element).find(".post_images_number")[0];
    const label_element      = $(post_images_number).find(".label");

    // HACK!! - "-1" was visually inferred
    $(post_images_number).css("right", `-${$(post_images_number).outerWidth()-1}px`);
    $(label_element).css("left", `${$(post_images_number).outerWidth()}px`);

    $(p_post_element).mouseover((p_e)=>{
        $(post_images_number).css("visibility", "visible");
    });
    $(p_post_element).mouseout((p_e)=>{
        $(post_images_number).css("visibility", "hidden");
    });
}