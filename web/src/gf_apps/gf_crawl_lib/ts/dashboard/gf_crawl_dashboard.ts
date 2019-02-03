/*
GloFlow media management/publishing system
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

///<reference path="../../../../d/jquery.d.ts" />

import "./gf_crawl_images_browser";
import "./gf_crawl_events";

namespace gf_crawl_dashboard {

declare var EventSource;
//-------------------------------------------------
$(document).ready(()=>{
	//-------------------------------------------------
	function log_fun(p_g,p_m) {
		var msg_str = p_g+':'+p_m
		//chrome.extension.getBackgroundPage().console.log(msg_str);

		switch (p_g) {
			case "INFO":
				console.log("%cINFO"+":"+"%c"+p_m, "color:green; background-color:#ACCFAC;", "background-color:#ACCFAC;");
				break;
			case "FUN_ENTER":
				console.log("%cFUN_ENTER"+":"+"%c"+p_m, "color:yellow; background-color:lightgray", "background-color:lightgray");
				break;
		}
	}
	//-------------------------------------------------
	gf_crawl_dashboard.init(log_fun);
});
//-------------------------------------------------
export function init(p_log_fun) {
	gf_crawl_events.init_SSE(p_log_fun);

	//---------------------
	//IMAGES
	$('#get_recent_images_btn').on('click', ()=>{
		gf_crawl_images_browser.init__recent_images(p_log_fun);
	});
	//---------------------
}
//---------------------------------------------------
}