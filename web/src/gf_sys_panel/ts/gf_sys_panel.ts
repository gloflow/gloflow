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

///<reference path="../../d/jquery.d.ts" />

import * as gf_identity from "./../../gf_apps/gf_identity/ts/gf_identity";

//-----------------------------------------------------
export async function init_with_auth(p_log_fun) {
	
	// STANDARD - non-admin urls
	const urls_map = gf_identity.get_standard_http_urls();
	
	const auth_http_api_map = gf_identity.get_http_api(urls_map);

	await init(auth_http_api_map, p_log_fun);
}

//-----------------------------------------------------
export async function init(p_auth_http_api_map, p_log_fun) {

	const sys_panel_element = $(
		`<div id="sys_panel">
			<div id="view_handle">
				<img src="/images/d/gf_sys_panel_view_handle.png"></img>
			</div>
			
			<div id="background">
				<div id="controls">
					<div id="home_btn">
						<img src="/images/d/gf_logo_icon.png"></img>
					</div>

					<div class="apps">
						<div id="images_app_btn"><a href="/images/flows/browser">Images</a></div>
						<div id="publisher_app_btn"><a href="/posts/browser">Posts</a></div>
						<div id="domains_app_btn"><a href="/a/domains/browser">Domains</a></div>
						<div id="bookmarks_app_btn"><a href="/v1/bookmarks/get">B</a></div>
					</div>

					<div id="auth">
						<div id="login_btn">login</div>
						<div id="current_user">
							
						</div>
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

	// HOME_BTN
	$(sys_panel_element).find("#home_btn").on("click", ()=>{
		window.location.href = "/"
	});

	//--------------------------
	// AUTH
	if (p_auth_http_api_map != null) {
		const me_user_map = await p_auth_http_api_map["general"]["get_me"]();


		const user_profile_img_url_str = me_user_map["profile_image_url_str"];
		const user_name_str            = me_user_map["user_name_str"];

		// IMG
		if (user_profile_img_url_str != "") { 
			$(sys_panel_element).find("#auth #current_user").append(`
				<img>${user_profile_img_url_str}</img>
			`);
		}

		// TEXT_SHORTHAND
		else {

			const shorthand_str = user_name_str[0];
			$(sys_panel_element).find("#auth #current_user").append(`
				<div id="shorthand_username">${shorthand_str}</div>
			`);
		}
	}

	//--------------------------


}