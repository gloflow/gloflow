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

///<reference path="../../../../d/jquery.d.ts" />

//-------------------------------------------------
export function init_image_date(p_image_element, p_log_fun) {

	const creation_time_element = $(p_image_element).find('.creation_time');
	const creation_time_f       = parseFloat($(creation_time_element).text());
	const creation_date         = new Date(creation_time_f*1000);

	const date_msg_str = $.timeago(creation_date);
	$(creation_time_element).text(date_msg_str);

	const creation_date__readable_str = creation_date.toDateString();
	const creation_date__readble      = $(`<div class="full_creation_date">${creation_date__readable_str}</div>`);

	$(creation_time_element).mouseover((p_e)=>{
		$(creation_time_element).append(creation_date__readble);

		// IMPORTANT!! - image size changed, so recalculate the Masonry layout.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container #items').masonry();
		$('#gf_images_flow_container #items').masonry(<any>'reloadItems');
	});

	$(creation_time_element).mouseout((p_e)=>{
		$(creation_date__readble).remove();

		// IMPORTANT!! - image size changed, so recalculate the Masonry layout.
		// IMPORTANT!! - for some reason both masonry() and masonry("reloadItems") are needed.
		$('#gf_images_flow_container #items').masonry();
		$('#gf_images_flow_container #items').masonry(<any>'reloadItems');
		
	});
}