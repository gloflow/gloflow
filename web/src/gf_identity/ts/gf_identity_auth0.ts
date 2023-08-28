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

//-------------------------------------------------
export function user_auth_pipeline() {

    // IMPORTAN!! - adding a unique param to this request to disable browser cache,
    //              since it can cause inconsistent behavior.
    const unique_param = new Date().getTime();
    const url_str = "/v1/identity/auth0/login?"+unique_param;

    // redirect the user to the GF auth0 login page, which will in turn
    // redirec the auth0 domain.
    window.location.href = url_str;
}

//-------------------------------------------------
export function logout() {

    // IMPORTAN!! - adding a unique param to this request to disable browser cache,
    //              since it can cause inconsistent behavior.
    const unique_param = new Date().getTime();
    const url_str = "/v1/identity/auth0/logout?"+unique_param;

    // redirect user to logout endpoint
    window.location.href = url_str;
}