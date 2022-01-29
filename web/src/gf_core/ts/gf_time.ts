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
///<reference path="../../d/jquery.timeago.d.ts" />

//-------------------------------------------------
export function sleep(p_miliseconds_int :number) {
	return new Promise(resolve => setTimeout(resolve, p_miliseconds_int));
}

//-------------------------------------------------
export function init_creation_date(p_target_element, p_log_fun) {

	const creation_time_element = $(p_target_element).find('.creation_time');
	const creation_time_f       = parseFloat($(creation_time_element).text());
	const creation_date         = new Date(creation_time_f*1000);

	const date_msg_str = $.timeago(creation_date);
	$(creation_time_element).text(date_msg_str);

	const creation_date__readable_str = creation_date.toDateString();

	$(creation_time_element).mouseover((p_e)=>{
		$(creation_time_element).text(creation_date__readable_str);
	});

	$(creation_time_element).mouseout((p_e)=>{
		$(creation_time_element).text(date_msg_str);
	});
}