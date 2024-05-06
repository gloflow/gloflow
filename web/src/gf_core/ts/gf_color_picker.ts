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

declare var iro :any;

//--------------------------------------------------------
export function init(p_target_selector_str :string,
    p_on_color_change_fun :any) {

    const container = $(`
		<div id="color_picker">
			<div id="control"></div>
			<div id="picked_color">#000000</div>
		</div>`);

    const color_picker_color_element = $(container).find("#picked_color");

    const color_picker = new iro.ColorPicker("#color_picker #control", {
        width: 100,   // size of the picker
        color: "#f00" // initial color
    });
    color_picker.on('color:change', async (p_color :any)=>{
        
        const color_hex_str = p_color.hexString;

        $(color_picker_color_element).text(color_hex_str);

        p_on_color_change_fun(color_hex_str);

        /*
		$("body").css("background-color", `${background_color_hex_str}`);

		// update component remotely
		const component_name_str = "background_color_picker";
		await gf_utils.update_viz_background_color(component_name_str,
			background_color_hex_str,
			p_http_api_map);
        */
    });

    return container;
}