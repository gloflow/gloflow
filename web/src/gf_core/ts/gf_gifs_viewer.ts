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

import * as gf_gifs from "./gf_gifs";

//-------------------------------------------------
export function init(p_image_element,
	p_img_id_str    :string,
	p_flow_name_str :string,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_gifs_viewer.init()');

	// const img_thumb_medium_url = $(p_image_element).find('img').data('img_thumb_medium_url');

	//----------------------
	// GIF_INDICATOR

	const gif_indicator_e = $('<div class="gif_indicator">GIF</div>');
	$(p_image_element).append(gif_indicator_e);

	const img_element = $(p_image_element).find('img')[0];

	// coordinates of <img> relative to its parent <div>, so that gif_indicator_e 
	// can take those coords into account for its own positioning.
	const img_x_int = img_element.offsetLeft;
	const img_y_int = img_element.offsetTop;

	// position in the middle of the image
	// offsetWidth/Height - The width and height of the entire element, including borders and padding, excluding margins.
	const gif_indicator_x_int = img_x_int + (img_element.offsetWidth - $(gif_indicator_e)[0].offsetWidth)/2;
	const gif_indicator_y_int = img_y_int + (img_element.offsetHeight - $(gif_indicator_e)[0].offsetHeight)/2;
	$(gif_indicator_e).css('left', gif_indicator_x_int+'px');
	$(gif_indicator_e).css('top' , gif_indicator_y_int+'px');


	$(gif_indicator_e).on('click', ()=>{

		// IMPORTANT!! - when gif_indicator is clicked activate the click event handler
		//               on the <img> tag of the GIF.
		$(p_image_element).find('img').trigger('click');
	});

	//----------------------


	$(p_image_element).find('img').on('click', ()=>{

		console.log('click');

		gf_gifs.http__gif_get_info(p_img_id_str,
			'gloflow.com',
			(p_gif_map)=>{

				console.log('GIF RECEIVED >>>>>>>>>> --------------')
				console.log(p_gif_map);

				const gif_gf_url_str = p_gif_map['gf_url_str'];

				view_gif(gif_gf_url_str);
			},
			(p_error_data_map)=>{

			},
			p_log_fun);
	});

	//-------------------------------------------------
	function view_gif(p_gif_gf_url_str) {
		
		const image_view = $(`
			<div id="gif_viewer">
				<div id="background"></div>
				<div id="gif">
					<img src="`+p_gif_gf_url_str+`"></img>
				</div>
			</div>`);

		console.log(p_gif_gf_url_str)
		$(p_image_element).append(image_view);

		//----------------------
		// BAKCGROUND
		const bg = $(image_view).find('#background');

		//----------------------

		$(image_view).find('img').on('load', ()=>{

		});

		//----------------------
		// CLOSE_VIEWER - when GIF thats playing is clicked again
		//                the image_view is removed.
	    $(bg).on('click', ()=>{
	    	$(image_view).remove();
	    });
		
	    //----------------------
	}

	//-------------------------------------------------
}