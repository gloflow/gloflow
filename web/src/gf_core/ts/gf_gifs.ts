/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

//---------------------------------------------------
export function http__gif_get_info(p_gf_img_id_str,
	p_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER','gf_gifs.http__gif_get_info()');

	const url_str = 'http://'+p_host_str+'/images/gif/get_info?gfimg_id='+p_gf_img_id_str;
	p_log_fun('INFO','url_str - '+url_str);

	//-------------------------
	//HTTP AJAX
	$.get(url_str,
		(p_data_map) => {
			console.log('response received');
			//const data_map = JSON.parse(p_data);

			if (p_data_map["status"] == 'OK') {
				const gif_map = p_data_map['data']['gif_map'];
				p_on_complete_fun(gif_map);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});
	//-------------------------	
}