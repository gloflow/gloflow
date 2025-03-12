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

import * as gf_core_utils from "./../../gf_core/ts/gf_utils";

//-------------------------------------------------
export function login() {
    const redirect_url_str = window.location.href;

    /*
    IMPORTANT!! - "login_success" is used to indicate to the page/js that the user has successfully logged in,
        and that that page load is the first load after the login.
        gives the page a change to finalize login.
    */
    const redirect_url_with_flag_str = gf_core_utils.add_query_param(redirect_url_str, "login_success", "1");

    user_auth_pipeline(redirect_url_with_flag_str);
}

//-------------------------------------------------
export function user_auth_pipeline(p_redirect_url_str :string="") {

    var url_str = "/v1/identity/auth0/login";

    if (p_redirect_url_str != "") {
        url_str += "?redirect_url=" + encodeURIComponent(p_redirect_url_str);
    }

    // redirect the user to the GF auth0 login page, which will in turn
    // redirec the auth0 domain.
    window.location.href = url_str;
}

//-------------------------------------------------
export async function login_finalize_if_needed() {
    const params = new URLSearchParams(window.location.search);
    if (params.get("login_success") === "1") {
        login_finalize();

        params.delete("login_success");

        const new_url_str = window.location.pathname + (params.toString() ? "?" + params.toString() : "");
        window.history.replaceState({}, "", new_url_str);
    }
}

//-------------------------------------------------
export function logout(p_redirect_url_str :string="") {

    var url_str = "/v1/identity/auth0/logout";

    if (p_redirect_url_str != "") {
        url_str += "?redirect_url="+encodeURIComponent(p_redirect_url_str);
    }

    // redirect user to logout endpoint
    window.location.href = url_str;
}

//-------------------------------------------------
function login_finalize() {

    console.log("login_finalize...");

    const url_str = "/v1/identity/auth0/login_finalize";
    fetch(url_str, { method: "POST", 
        // Ensures that cookies are sent and received with the request.
        credentials: "include" });
}