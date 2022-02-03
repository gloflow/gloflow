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
    $(p_image).on("mouseenter", async (p_event)=>{


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
            
            color_info_element.insertAfter(image);


            // apply shadow only when the cursor enters the color_info element
            $(color_info_element).on("mouseenter", ()=>{
                $(color_info_element).css("box-shadow", "1px 1px 10px #0000007a")
            })
            $(color_info_element).on("mouseleave", ()=>{
                $(color_info_element).css("box-shadow", "")
            })

            //-------------
            // COLOR_DOMINANT
            const color_dominant_element       = $(color_info_element).find(".color_dominant");
            var   color_dominant_label_element = $(`<div class="color_dominant_label">dominant color</div>`);
            var   color_dominant__copy_to_clipboard_btn;
            $(color_dominant_element).on("mouseenter", ()=>{

                // COLOR_DOMINANT_LABEL
                color_info_element.append(color_dominant_label_element);

                //-------------
                // COLOR_INSPECTOR
                var color_inspect_element = $(`<div class="color_inspect">
                    <div class='color_hex'>#${color_dominant_hex_str}</div>
                    <div class='color_large' style="background-color:#${color_dominant_hex_str};"></div>
                </div>`);
                $(color_info_element).append(color_inspect_element);
                
                //-------------

                //-------------
                // COPY_TO_CLIPBOARD
                if (color_dominant__copy_to_clipboard_btn == null) {
                    color_dominant__copy_to_clipboard_btn = init_copy_to_clipboard_btn([color_dominant_hex_str]);
                    $(color_dominant_element).append(color_dominant__copy_to_clipboard_btn);
                }

                //-------------
            });

            // IMPORTANT!! - using "mouseleave" so that label and clipboard btn are only removed when the whole
            //               color_dominant element is left, not just the .color element thats a child of .color_dominant
            // mouseenter and mouseleave - triggered when you enter and leave a hierarchy of nodes,
            //                             but not when you navigate that hierarchy's descendance.
            // mouseover and mouseout    - triggered when the mouse respectively enters and leaves
            //                             a node's "exclusive" space, so you get a "out" when the
            //                             mouse gets into a child node.
            $(color_dominant_element).on("mouseleave", (p_e)=>{

                // remove label
                $(color_dominant_label_element).remove();

                // remove color_inspect
                $(color_info_element).find(".color_inspect").remove();

                // remove copy_to_clipboard button
                $(color_dominant__copy_to_clipboard_btn).remove();
                color_dominant__copy_to_clipboard_btn = null;

            });

            //-------------
            // COLOR_PALLETE
            const color_pallete_element = $(color_info_element).find(".color_pallete");

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

                $(color_element).on("mouseenter", ()=>{
                    
                    // remove a previously appended color_inspect element, if its there
                    $(color_info_element).find(".color_inspect").remove();

                    color_info_element.append(color_inspect_element);
                });
                $(color_element).on("mouseleave", ()=>{
                    $(color_inspect_element).remove();
                });
                
                //-------------

                colors_hexes_lst.push(p_color_hex_str);
            })

            //-------------
            // COLOR_PALLETE_LABEL
            var color_pallete_label_element = $(`<div class="color_pallete_label">color pallete</div>`);
            var color_pallete__copy_to_clipboard_btn;
            $(color_pallete_element).on("mouseenter", ()=>{
                color_info_element.append(color_pallete_label_element);

                //-------------
                // COPY_TO_CLIPBOARD
                if (color_pallete__copy_to_clipboard_btn == null) {
                    color_pallete__copy_to_clipboard_btn = init_copy_to_clipboard_btn(colors_hexes_lst);
                    $(color_pallete_element).append(color_pallete__copy_to_clipboard_btn);
                }

                //-------------
            });
            $(color_pallete_element).on("mouseleave", ()=>{
                $(color_pallete_label_element).remove();

                $(color_pallete__copy_to_clipboard_btn).remove();
                color_pallete__copy_to_clipboard_btn = null;
            });
            
            //-------------

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
        
        var label_element = $(`<div class="color_to_clipboard_label">copy</div>`);

        $(element).on("click", async (p_e)=>{

            p_e.stopPropagation();
            
            var colors_for_clipboard_str = p_colors_hexes_lst.join(",");

            // COPY_TO_CLIPBOARD
            // navigator.clipboard - only defined when served over https|localhost
            await navigator.clipboard.writeText(colors_for_clipboard_str);

            $(element).css("background-color", "green");
            $(label_element).text("copied");
        });

        
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