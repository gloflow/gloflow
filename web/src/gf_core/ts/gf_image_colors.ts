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
export function init_pallete(p_image,
    p_assets_paths_map,
    p_on_color_compute_fun) {

    var image_colors_shown_bool = false;
    $(p_image).on("mouseover", async (p_event)=>{


        if (!image_colors_shown_bool) {
            
            const image        = p_event.target;
            const image_colors = gf_image_process.get_colors(image);
            const color_dominant_hex_str = image_colors.color_hex_str;
            const color_palette_lst      = image_colors.color_palette_lst;
            
            console.log("dominant color - ", color_dominant_hex_str);

            const color_info_element = $(`<div class="colors_info">
                <div class="color_dominant">
                    <div class="color" style="background-color:#${color_dominant_hex_str};">
                    </div>
                </div>
                <div class="color_pallete">
                    <div class="colors">
                    </div>
                </div>
            </div>`);

            // // IMPORTANT!! - change to color of the whole image_info control to match the dominant color of the
            // //               image its displaying.
            // $(p_image_info_element).css("background-color", `#${color_dominant_hex_str}`);				
            
            color_info_element.insertAfter(image);

            //-------------
            // COLOR_DOMINANT
            const color_dominant_element       = $(color_info_element).find(".color_dominant");
            var   color_dominant_label_element = $(`<div class="color_dominant_label">dominant color</div>`);
            var   color_dominant__copy_to_clipboard_btn;
            $(color_dominant_element).on("mouseover", ()=>{

                // COLOR_DOMINANT_LABEL
                color_info_element.append(color_dominant_label_element);

                //-------------
                // COLOR_INSPECTOR
                var color_inspect_element = $(`<div class="color_inspect">
                    <div class='color_hex'>#${color_dominant_hex_str}</div>
                    <div class='color_large' style="background-color:#${color_dominant_hex_str};"></div>
                </div>`);
                $(color_info_element).append(color_inspect_element);

                /*$("body").on("click", ()=>{
                    $(color_inspect_element).remove();
                });*/

                /*$(color_dominant_element).on("mouseout", ()=>{
                    $(color_inspect_element).remove();
                });*/
                
                //-------------

                //-------------
                // COPY_TO_CLIPBOARD
                if (color_dominant__copy_to_clipboard_btn == null) {
                    color_dominant__copy_to_clipboard_btn = init_copy_to_clipboard_btn([color_dominant_hex_str]);
                    $(color_dominant_element).append(color_dominant__copy_to_clipboard_btn);
                }

                //-------------
            });
            $(color_dominant_element).on("mouseout", (p_e)=>{


                $(color_dominant_label_element).remove();
                
                // $(color_dominant__copy_to_clipboard_btn).remove();
                // color_dominant__copy_to_clipboard_btn = null;
            });

            //-------------
            // COLOR_PALLETE
            const color_pallete_element = $(color_info_element).find(".color_pallete");
            // const color_pallete_sub_lst = image_colors.color_palette_lst.slice(1, 6);

            const colors_hexes_lst = [];
            color_palette_lst.forEach((p_color_hex_str)=>{

                const color_element = $(`<div class="color" style="background-color:#${p_color_hex_str};"></div>`);
                $(color_pallete_element).find(".colors").append(color_element);

                //-------------
                // COLOR_INSPECTOR
                var color_inspect_element = $(`<div class="color_inspect">
                    <div class='color_hex'>#${p_color_hex_str}</div>
                    <div class='color_large' style="background-color:#${p_color_hex_str};"></div>

                </div>`);
                $(color_element).on("mouseover", ()=>{
                    
                    // remove a previously appended color_inspect element, if its there
                    $(color_info_element).find("#color_inspect").remove();

                    // color_pallete_element.append(color_inspect_element);
                    color_info_element.append(color_inspect_element);
                });

                /*$("body").on("click", ()=>{
                    $(color_inspect_element).remove();
                });*/

                /*$(color_element).on("mouseout", ()=>{
                    $(color_inspect_element).remove();
                });*/
                
                //-------------

                colors_hexes_lst.push(p_color_hex_str);
            })

            //-------------
            // COLOR_PALLETE_LABEL
            var color_pallete_label_element = $(`<div class="color_pallete_label">color pallete</div>`);
            var color_pallete__copy_to_clipboard_btn;
            $(color_pallete_element).on("mouseover", ()=>{
                color_info_element.append(color_pallete_label_element);

                //-------------
                // COPY_TO_CLIPBOARD
                if (color_pallete__copy_to_clipboard_btn == null) {
                    color_pallete__copy_to_clipboard_btn = init_copy_to_clipboard_btn(colors_hexes_lst);
                    $(color_pallete_element).append(color_pallete__copy_to_clipboard_btn);
                }

                //-------------
            });
            $(color_pallete_element).on("mouseout", ()=>{
                $(color_pallete_label_element).remove();
            });
            
            //-------------

            /*// leave color picking info on screen even when mouse moves out of the control,
            // to allow users to interact with it later.
            $("body").on("click", ()=>{

                //-------------
                // COPY_TO_CLIPBOARD
                if (color_pallete__copy_to_clipboard_btn != null) {
                    // has to be here, removed when user goes out of the entire color_info element,
                    // and not just the pallete, so that the user has a chance to click it
                    $(color_pallete__copy_to_clipboard_btn).remove();
                    color_pallete__copy_to_clipboard_btn = null;
                }

                //-------------
            });*/

            /*$(color_info_element).on("mouseout", ()=>{

                //-------------
                // COPY_TO_CLIPBOARD
                if (copy_to_clipboard_btn != null) {
                    // has to be here, removed when user goes out of the entire color_info element,
                    // and not just the pallete, so that the user has a chance to click it
                    $(copy_to_clipboard_btn).remove();
                    copy_to_clipboard_btn = null;
                }

                //-------------
            });*/

            image_colors_shown_bool = true;

            p_on_color_compute_fun(color_dominant_hex_str,
                colors_hexes_lst);
        }
    });

    //--------------------------------------------------------
    function init_copy_to_clipboard_btn(p_colors_hexes_lst) {
        const element = $(`
            <div id='copy_to_clipboard_btn'>
                <img src="${p_assets_paths_map["copy_to_clipboard_btn"]}"></img>
            </div>`);
        
        $(element).on("click", async ()=>{
            var colors_for_clipboard_str = p_colors_hexes_lst.join(",");

            // COPY_TO_CLIPBOARD
            // navigator.clipboard - only defined when served over https|localhost
            await navigator.clipboard.writeText(colors_for_clipboard_str);

            $(element).css("background-color", "green");
        });

        var label_element = $(`<div class="color_to_clipboard_label">copy</div>`);
        var copy_to_clipboard_btn;
        $(element).on("mouseover", ()=>{
            $(element).append(label_element);
        });
        $(element).on("mouseout", ()=>{
            $(label_element).remove();
        });

        return element;
    }

    //--------------------------------------------------------
}
