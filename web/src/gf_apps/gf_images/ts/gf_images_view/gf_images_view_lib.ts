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
///<reference path="../../../../d/masonry.layout.d.ts" />
///<reference path="../../../../d/jquery.timeago.d.ts" />

import * as gf_sys_panel from "./../../../../gf_sys_panel/ts/gf_sys_panel";

//-------------------------------------------------
declare var URLSearchParams;

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;

//-------------------------------------------------
export function init(p_log_fun) {
	
	gf_sys_panel.init_with_auth(p_log_fun);

	//---------------------
	// UPLOAD__INIT

    const target_flow_name_str = "general"
	init_upload(target_flow_name_str, p_log_fun);

	//------------------
}

//---------------------------------------------------
function init_upload(p_flow_name_str :string,
	p_log_fun) {

	// use "" so that no host is set in URL's for issued requests
	// (forces usage of origin host that the page came from)
	const target_full_host_str = "";
	gf_upload__init(p_flow_name_str,
		target_full_host_str,
		
		//-------------------------------------------------
		// p_on_upload_fun
		async (p_upload_gf_image_id_str)=>{






		});

		//-------------------------------------------------
}