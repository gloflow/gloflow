/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

//-------------------------------------------------
export function get_gf_host() {
	const domain_str   = window.location.hostname;
	const protocol_str = window.location.protocol;
	const gf_host_str = `${protocol_str}//${domain_str}`;
	console.log("gf_host", gf_host_str);
	return gf_host_str;
}

//-------------------------------------------------
export function click_outside(p_element, p_on_click_fun) {
    $(document).on("click", (p_e)=>{
            
        // check if a click is not on the element or a child of the element
        if (!p_element.is(p_e.target) && $(p_element).has(p_e.target).length === 0) {
            p_on_click_fun()
        }
    });
}