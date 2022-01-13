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

import * as gf_identity_http from "./gf_identity_http";

//-------------------------------------------------
export async function user_auth_pipeline(p_user_name_str :string) {


    //--------------------------
    // PREFLIGHT_HTTP
    const data_map = await gf_identity_http.user_preflight(p_user_name_str, null);

    const user_exists_bool = data_map["user_exists_bool"];
    const nonce_val_str    = data_map["nonce_val_str"];

    //--------------------------
    // user exists, log them in
    if (user_exists_bool) {
        console.log("USER_EXISTS");


        const pass_hash_str :string = hash_pass();

        const login_data_map = await gf_identity_http.user_userpass_login(p_user_name_str,
            pass_hash_str)

    }
    // no-user in the system, offer to create new
    else {
        
        console.log("NO USER");



    }
}



function hash_pass() :string {


    return "";
}