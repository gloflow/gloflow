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

//--------------------------------------------------------
interface GF_image_colors {
    color_hex_str     :string
    color_palette_lst :string[]
}

declare var ColorThief;

//--------------------------------------------------------
export function get_colors(p_image) :GF_image_colors {
    
    const color_thief = new ColorThief();

    // get a dominant color for an image
	const color_lst   = color_thief.getColor(p_image);   // DOMINANT COLOR
	const palette_lst = color_thief.getPalette(p_image); // COLOR PALLETE



    const color_hex_str = rgb_to_hex(color_lst[0], color_lst[1], color_lst[2]);

    const hex_pallete_lst = [];
    for (var e of palette_lst) {

        const hex_str = rgb_to_hex(e[0],e[1],e[2]);
        hex_pallete_lst.push(hex_str);
    };

    const colors :GF_image_colors = {
        color_hex_str:     color_hex_str,
        color_palette_lst: hex_pallete_lst,
    };

    return colors;
}

//--------------------------------------------------------
function rgb_to_hex(r, g, b) {
    return to_hex(r) + to_hex(g) + to_hex(b);
}

//--------------------------------------------------------
function to_hex(c) {
    const hex = c.toString(16);
    return hex.length == 1 ? "0" + hex : hex;
}