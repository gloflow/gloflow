/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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
function log_fun(p_g, p_m) {
	const msg_str = `${p_g}:${p_m}`;
	// chrome.extension.getBackgroundPage().console.log(msg_str);

	switch (p_g) {
		case "INFO":
			console.log(`%cINFO:%c${p_m}`, "color:green; background-color:#ACCFAC;", "background-color:#ACCFAC;");
			break;
		case "FUN_ENTER":
			console.log(`%cFUN_ENTER:%c${p_m}`, "color:yellow; background-color:lightgray", "background-color:lightgray");
			break;
	}
}

//---------------------------------------------------
function send_msg_to_bg_page(p_msg_map, p_on_complete_fun) {

	chrome.runtime.sendMessage(p_msg_map,
		(p_resp)=>{
			p_on_complete_fun(p_resp);
		});
}

//---------------------------------------------------
function hash_code(p_str,
	p_log_fun) {
	
	var hash = 0;
	if (p_str.length == 0) return hash;
	for (i = 0; i < p_str.length; i++) {
		char = p_str.charCodeAt(i);
		hash = ((hash<<5)-hash)+char;
		hash = hash & hash; // Convert to 32bit integer
	}
	return hash;
}

//---------------------------------------------------
function get_image_histogram(p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER', 'utils.get_image_histogram()');

	const hist_map = {};	
	const img       = new Image();
	img.onload      = () => {
		Pixastic.process(img,"colorhistogram", {
			paint:       true,
			returnValue: hist_map
		});
		hist_map.rvals; // <- array[255] red channel
		hist_map.gvals; // <- array[255] green channel
		hist_map.bvals; // <- array[255] blue channel

		p_on_complete_fun(hist_map);
	}
	document.body.appendChild(img);
	img.src = "myimage.jpg";
}