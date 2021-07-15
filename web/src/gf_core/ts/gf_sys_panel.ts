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

//-----------------------------------------------------
export function init(p_log_fun) {
	// p_log_fun('FUN_ENTER', 'gf_sys_panel.init()');

	const sys_panel_element = $(
		`<div id="sys_panel">
			<div id="view_handle">
				<img src="/images/d/gf_sys_panel_view_handle.png"></img>
			</div>
			<div id="home_btn">
				<img src="/images/d/gf_logo_icon.png"></img>
			</div>
			<div id="images_app_btn"><a href="/images/flows/browser">Images</a></div>
			<div id="publisher_app_btn"><a href="/posts/browser">Posts</a></div>
			<div id="domains_app_btn"><a href="/a/domains/browser">Domains</a></div>
			<div id="bookmarks_app_btn"><a href="/v1/bookmarks/get">B</a></div>
			<div id="get_invited_btn">get invited</div>
			<div id="login_btn">login</div>
		</div>`);

	$('body').append(sys_panel_element);

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
}