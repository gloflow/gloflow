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

// ///<reference path="./../../d/jquery.d.ts" />

import * as gf_core_utils    from "./../../gf_core/ts/gf_utils";
import * as gf_identity      from "./../../gf_identity/ts/gf_identity";
import * as gf_identity_http from "./../../gf_identity/ts/gf_identity_http";

//-----------------------------------------------------
export async function init_with_auth(p_log_fun) {
	
	// STANDARD - non-admin urls
	const current_host_str = gf_core_utils.get_current_host();
	const urls_map = gf_identity_http.get_standard_http_urls(current_host_str);
	
	const auth_http_api_map = gf_identity_http.get_http_api(urls_map, current_host_str);

	await init(auth_http_api_map, urls_map, p_log_fun);
}

//-----------------------------------------------------
export async function init(p_auth_http_api_map,
	p_urls_map,
	p_log_fun) {

	const handle_img_url_str = "https://assetspub.gloflow.com/assets/gf_sys/gf_sys_panel_view_handle.svg";
	const sys_panel_element = $(
		`<div id="sys_panel">
			<div id="view_handle" class="gf_center">
				<img src="${handle_img_url_str}"></img>
			</div>
			
			<div id="background">
				<div id="controls">
					<div id="landing_page_btn">
						<img src="/images/d/gf_logo_icon.png"></img>
					</div>

					<div id="apps">
						<div id="images_app_btn"    class="gf_center"><a href="/images/flows/browser">Images</a></div>
						<div id="bookmarks_app_btn" class="gf_center"><a href="/v1/bookmarks/get">Bookmarks</a></div>
						<div id="domains_app_btn"   class="gf_center"><a href="/a/domains/browser">D</a></div>
					</div>
				</div>
			</div>
			
		</div>`);

	$('body').append(sys_panel_element);

	// VIEW_HANDLE
	$(sys_panel_element).find('#view_handle').on('mouseover', (p_e)=>{
		$(sys_panel_element).animate({
			top: 0 // move it
		},
		200,
		()=>{
			$(sys_panel_element).find('#view_handle').css('visibility', 'hidden');
		});
	});

	// LANDING_PAGE_BTN
	$(sys_panel_element).find("#landing_page_btn").on("click", ()=>{
		window.location.href = "/"
	});

	//--------------------------
	// AUTH
	if (p_auth_http_api_map != null) {

		const parent_node  = $(sys_panel_element).find("#controls")[0];
		const home_url_str = p_urls_map["home"];

		gf_identity.init_me_control(parent_node,
			p_auth_http_api_map,
			home_url_str);
	}

	//--------------------------
}