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

import * as gf_time      from "./../../../../gf_core/ts/gf_time";
import * as gf_sys_panel from "./../../../../gf_sys_panel/ts/gf_sys_panel";

//--------------------------------------------------------
$(document).ready(()=>{
	//-------------------------------------------------
	function log_fun(p_g, p_m) {
		var msg_str = p_g+':'+p_m;
		switch (p_g) {
			case "INFO":
				console.log("%cINFO"+":"+"%c"+p_m, "color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
				break;
			case "FUN_ENTER":
				console.log("%cFUN_ENTER"+":"+"%c"+p_m, "color:yellow; background-color:lightgray","background-color:lightgray");
				break;
		}
	}

	//-------------------------------------------------
    init(log_fun);
});

//-------------------------------------------------
function init(p_log_fun) {
    console.log("start")


	gf_sys_panel.init_with_auth(p_log_fun);






	$("#bookmarks .bookmark").each((p_i, p_bookmark_element)=>{

		gf_time.init_creation_date(p_bookmark_element, p_log_fun);
	});
}

