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

///<reference path="../../../../d/jquery.d.ts" />

import * as gf_tag_objects_lib from "./gf_tag_objects_lib";

//-------------------------------------------------
// main entrypoint for the OS gloflow version.
// all the code in the the *_lib module, init(), which can be called from
// other external code that controls initialization.
$(document).ready(()=>{
	//-------------------------------------------------
	function log_fun(p_g,p_m) {
		var msg_str = p_g+':'+p_m
		switch (p_g) {
			case "INFO":
				console.log("%cINFO"+":"+"%c"+p_m,"color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
				break;
			case "FUN_ENTER":
				console.log("%cFUN_ENTER"+":"+"%c"+p_m,"color:yellow; background-color:lightgray","background-color:lightgray");
				break;
		}
	}

	//-------------------------------------------------


	// PLUGINS
	const plugin_callbacks_map = {};
	
	const tag_str = $("#tag_info #tag_name").text();

	gf_tag_objects_lib.init(tag_str, plugin_callbacks_map, log_fun);
});