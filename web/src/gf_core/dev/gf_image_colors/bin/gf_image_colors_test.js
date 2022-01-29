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
System.register("ts/gf_image_process", [], function (exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    //--------------------------------------------------------
    function get_colors(p_image) {
        const color_thief = new ColorThief();
        // get a dominant color for an image
        const color_lst = color_thief.getColor(p_image); // DOMINANT COLOR
        const palette_lst = color_thief.getPalette(p_image); // COLOR PALLETE
        const color_hex_str = rgb_to_hex(color_lst[0], color_lst[1], color_lst[2]);
        const hex_pallete_lst = [];
        for (var e of palette_lst) {
            const hex_str = rgb_to_hex(e[0], e[1], e[2]);
            hex_pallete_lst.push(hex_str);
        }
        ;
        const colors = {
            color_hex_str: color_hex_str,
            color_palette_lst: hex_pallete_lst,
        };
        return colors;
    }
    exports_1("get_colors", get_colors);
    //--------------------------------------------------------
    function rgb_to_hex(r, g, b) {
        return to_hex(r) + to_hex(g) + to_hex(b);
    }
    //--------------------------------------------------------
    function to_hex(c) {
        const hex = c.toString(16);
        return hex.length == 1 ? "0" + hex : hex;
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
System.register("ts/gf_image_colors", ["ts/gf_image_process"], function (exports_2, context_2) {
    "use strict";
    var gf_image_process;
    var __moduleName = context_2 && context_2.id;
    //--------------------------------------------------------
    function init_pallete(p_image, p_assets_paths_map, p_on_color_compute_fun) {
        var image_colors_shown_bool = false;
        $(p_image).on("mouseover", async (p_event) => {
            if (!image_colors_shown_bool) {
                const image = p_event.target;
                const image_colors = gf_image_process.get_colors(image);
                const color_dominant_hex_str = image_colors.color_hex_str;
                const color_palette_lst = image_colors.color_palette_lst;
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
                //-------------
                // COLOR_DOMINANT
                const color_dominant_element = $(color_info_element).find(".color_dominant");
                var color_dominant_label_element = $(`<div class="color_dominant_label">dominant color</div>`);
                var color_dominant__copy_to_clipboard_btn;
                $(color_dominant_element).on("mouseenter", () => {
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
                $(color_dominant_element).on("mouseleave", (p_e) => {
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
                color_palette_lst.forEach((p_color_hex_str) => {
                    const color_element = $(`<div class="color" style="background-color:#${p_color_hex_str};"></div>`);
                    $(color_pallete_element).find(".colors").append(color_element);
                    //-------------
                    // COLOR_INSPECTOR
                    var color_inspect_element = $(`<div class="color_inspect">
                    <div class='color_hex'>#${p_color_hex_str}</div>
                    <div class='color_large' style="background-color:#${p_color_hex_str};"></div>
                </div>`);
                    $(color_element).on("mouseenter", () => {
                        // remove a previously appended color_inspect element, if its there
                        $(color_info_element).find(".color_inspect").remove();
                        color_info_element.append(color_inspect_element);
                    });
                    $(color_element).on("mouseleave", () => {
                        $(color_inspect_element).remove();
                    });
                    //-------------
                    colors_hexes_lst.push(p_color_hex_str);
                });
                //-------------
                // COLOR_PALLETE_LABEL
                var color_pallete_label_element = $(`<div class="color_pallete_label">color pallete</div>`);
                var color_pallete__copy_to_clipboard_btn;
                $(color_pallete_element).on("mouseenter", () => {
                    color_info_element.append(color_pallete_label_element);
                    //-------------
                    // COPY_TO_CLIPBOARD
                    if (color_pallete__copy_to_clipboard_btn == null) {
                        color_pallete__copy_to_clipboard_btn = init_copy_to_clipboard_btn(colors_hexes_lst);
                        $(color_pallete_element).append(color_pallete__copy_to_clipboard_btn);
                    }
                    //-------------
                });
                $(color_pallete_element).on("mouseleave", () => {
                    $(color_pallete_label_element).remove();
                    $(color_pallete__copy_to_clipboard_btn).remove();
                    color_pallete__copy_to_clipboard_btn = null;
                });
                //-------------
                image_colors_shown_bool = true;
                p_on_color_compute_fun(color_dominant_hex_str, colors_hexes_lst);
            }
        });
        //--------------------------------------------------------
        function init_copy_to_clipboard_btn(p_colors_hexes_lst) {
            const element = $(`
            <div id='copy_to_clipboard_btn'>
                <img src="${p_assets_paths_map["copy_to_clipboard_btn"]}"></img>
            </div>`);
            $(element).on("click", async () => {
                var colors_for_clipboard_str = p_colors_hexes_lst.join(",");
                // COPY_TO_CLIPBOARD
                // navigator.clipboard - only defined when served over https|localhost
                await navigator.clipboard.writeText(colors_for_clipboard_str);
                $(element).css("background-color", "green");
            });
            var label_element = $(`<div class="color_to_clipboard_label">copy</div>`);
            var copy_to_clipboard_btn;
            $(element).on("mouseover", () => {
                $(element).append(label_element);
            });
            $(element).on("mouseout", () => {
                $(label_element).remove();
            });
            return element;
        }
        //--------------------------------------------------------
    }
    exports_2("init_pallete", init_pallete);
    return {
        setters: [
            function (gf_image_process_1) {
                gf_image_process = gf_image_process_1;
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
///<reference path="./../../../d/jquery.d.ts" />
System.register("dev/gf_image_colors/gf_image_colors_test", ["ts/gf_image_colors"], function (exports_3, context_3) {
    "use strict";
    var gf_image_colors;
    var __moduleName = context_3 && context_3.id;
    return {
        setters: [
            function (gf_image_colors_1) {
                gf_image_colors = gf_image_colors_1;
            }
        ],
        execute: function () {
            ///<reference path="./../../../d/jquery.d.ts" />
            //-------------------------------------------------
            $(document).ready(() => {
                const img = $(".image_info").find("img")[0];
                const assets_paths_map = {
                    "copy_to_clipboard_btn": "./../../../../assets/gf_copy_to_clipboard_btn.svg"
                };
                gf_image_colors.init_pallete(img, assets_paths_map, (p_color_dominant_hex_str, p_colors_hexes_lst) => {
                    // $(".image_info").css("background-color", p_color_dominant_hex_str)
                });
            });
        }
    };
});
