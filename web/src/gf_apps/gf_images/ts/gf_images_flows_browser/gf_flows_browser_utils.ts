/*
GloFlow application and media management/publishing platform
Copyright (C) 2025 Ivan Trajkovic

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

//---------------------------------------------------
// CURRENT_PAGES_DISPLAY
//---------------------------------------------------

export function current_pages_display__init(p_log_fun :Function) {

	const container = $(`
		<div id="current_pages_display"'>
			<div id="title">pages</div>
			<div id="start_page">1</div>
			<div id="to">to</div>
			<div id="end_page">6</div>
		</div>`);

	return container;
}

export function current_pages_display__reset(p_start_page_int :number,
	p_end_page_int :number) {
	$("#current_pages_display #start_page").text(p_start_page_int);
	$("#current_pages_display #end_page").text(p_end_page_int);
}