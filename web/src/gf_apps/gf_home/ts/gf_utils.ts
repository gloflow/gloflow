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


//--------------------------------------------------------
export async function get_user_info(p_identity_http_api_map) {

	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
		var screen_name_str;
		var profile_img_url_str;
		try {
			const me_user_map = await p_identity_http_api_map["general"]["get_me"]();
			screen_name_str     = me_user_map["screen_name_str"];
			profile_img_url_str = me_user_map["user_profile_img_url_str"];

		} catch (error_map) {

			// DEFAULT
			screen_name_str     = "user_screenname";
			profile_img_url_str = "https://media.gloflow.com/thumbnails/786f79c0c85c08c7b1c0b3e11d6cae1e_thumb_small.png";
		}

		const user_info_map = {
			"screen_name_str":     screen_name_str,
			"profile_img_url_str": profile_img_url_str
		};

		p_resolve_fun(user_info_map);

	});
	return p;
}

//--------------------------------------------------------
export function update_viz_component_remote(p_component_name_str :string,
	p_drag_data_map,
	p_http_api_map) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const x_new_int = p_drag_data_map["x_int"];
		const y_new_int = p_drag_data_map["y_int"];
		const props_change_map = {
			"screen_x_int": x_new_int,
			"screen_y_int": y_new_int,
		};

		const output_map = await p_http_api_map["home"]["viz_update_fun"](p_component_name_str,
			props_change_map);

		p_resolve_fun(output_map);
	});
	return p;
}

//--------------------------------------------------------
export function update_viz_background_color(p_component_name_str :string,
	p_background_color_str :string,
	p_http_api_map) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
		const props_change_map = {
			"background_color_str": p_background_color_str,
		};

		const output_map = await p_http_api_map["home"]["viz_update_fun"](p_component_name_str,
			props_change_map);

		p_resolve_fun(output_map);
	});
	return p;
}