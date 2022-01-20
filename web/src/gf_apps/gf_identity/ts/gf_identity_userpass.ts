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

import * as gf_3d from "./../../../gf_core/ts/gf_3d";
import * as gf_identity_http from "./gf_identity_http";

//-------------------------------------------------
export async function user_auth_pipeline(p_http_api_map) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const container = $(`
        <div id="user_and_pass_dialog">
            <div id="username_input">
                <input id='username_input' placeholder="user name"></input>
            </div>
            <div id="pass_input">
                <input id="pass_input" placeholder="password" type="password"></input>
            </div>
            <div id="email_input">
                <input id='email_input' placeholder="email"></input>
            </div>
            <div id="login_btn">login</div>
            <div id="create_btn">create user</div>
        </div>`);
        $("#identity").append(container);

        // gf_3d.div_follow_mouse($(container)[0], document, 90);

        $(container).find("input#username_input").focus();

        //-------------------------------------------------
        async function login_activate() {
            console.log("login activate");

            const user_name_str = $(container).find("#username_input").val();
            const pass_str      = $(container).find("#pass_input").val();
            const pass_hash_str = await hash_pass(pass_str); 

            const login_output_map = await p_http_api_map["userpass"]["user_login_fun"](user_name_str,
                pass_hash_str as string);
        }

        //-------------------------------------------------

        $(container).on('keyup', function (e) {
            if (e.key === 'Enter' || e.keyCode === 13) {
                login_activate();
            }
        });

        $(container).find("#login_btn").on('click', async ()=>{
            login_activate();
        });

        $(container).find("#create_btn").on('click', async ()=>{

            // ADD!! - do frontend validation on username
            const user_name_str = $(container).find("#username_input").val();
            const pass_str      = $(container).find("#pass_input").val();
            const email_str     = $(container).find("#email_input").val();
            const pass_hash_str = await hash_pass(pass_str); 

            const create_output_map = await p_http_api_map["userpass"]["user_create_fun"](user_name_str,
                pass_hash_str as string,
                email_str as string);
        });
    });
    return p;
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