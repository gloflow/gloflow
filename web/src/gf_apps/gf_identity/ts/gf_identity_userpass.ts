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
export async function user_auth_pipeline() {

    const user_name_str = "";
    //--------------------------
    // PREFLIGHT_HTTP
    const data_map = await gf_identity_http.user_preflight(user_name_str, null);

    const user_exists_bool = data_map["user_exists_bool"];
    const nonce_val_str    = data_map["nonce_val_str"];

    //--------------------------

    // user exists, log them in
    if (user_exists_bool) {
        console.log("USER_EXISTS");

        const pass_str = "";
        const pass_hash_str = await hash_pass(pass_str); 

        const login_data_map = await gf_identity_http.user_userpass_login(user_name_str,
            pass_hash_str as string)

    }
    // no-user in the system, offer to create new
    else {
        
        console.log("NO USER");

        user_create();

    }
}

//-------------------------------------------------
function user_create() {


    const container = `
    <div id="user_and_pass>
        <div id="username_input>
            <input id='username_input'></input>
        </div>
        <div id="pass_input>
            <input id='pass_input'></input>
        </div>
    </div>`;
}


//-------------------------------------------------
async function hash_pass(p_pass_str) {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
        const encoder = new TextEncoder();
        const data    = encoder.encode(p_pass_str);
        
        // CRYPTO_HASH
        const pass_hash_buff = await crypto.subtle.digest("SHA-256", data);

        const hash_arr     = Array.from(new Uint8Array(pass_hash_buff));
        const hash_hex_str = hash_arr.map(b => b.toString(16).padStart(2, "0")).join("");

        p_resolve_fun(hash_hex_str);
    });
    return p;
}