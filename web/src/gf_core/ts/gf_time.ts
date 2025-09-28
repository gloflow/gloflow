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
export function init_creation_date(p_target_element :HTMLElement, p_log_fun :Function) {

	init_timeago(p_target_element,
		".creation_time",
		p_log_fun);
}

//-------------------------------------------------
export function init_timeago(p_target_element :HTMLElement,
	p_selector_str: string,
	p_log_fun :Function) {

	const time_element = $(p_target_element).find(p_selector_str);
	const time_f       = parseFloat($(time_element).text());
	const date         = new Date(time_f*1000);

	const date_msg_str = $.timeago(date);
	$(time_element).text(date_msg_str);

	const date__readable_str = date.toDateString();

	$(time_element).mouseover((p_e)=>{
		$(time_element).text(date__readable_str);
	});

	$(time_element).mouseout((p_e)=>{
		$(time_element).text(date_msg_str);
	});
}

