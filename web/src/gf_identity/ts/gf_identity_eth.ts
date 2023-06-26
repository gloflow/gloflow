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

///<reference path="../../d/jquery.d.ts" />

import * as gf_identity_http from "./gf_identity_http";
import * as gf_ops           from "./gf_ops";

declare const window: any;
declare var Web3;

//-------------------------------------------------
export async function user_auth_pipeline(p_http_api_map) {

    const user_address_eth_str = await wallet_connect() as string;

    //--------------------------
    // PREFLIGHT_HTTP
    const output_map = await p_http_api_map["eth"]["user_preflight_fun"](user_address_eth_str);
    const user_exists_bool = output_map["user_exists_bool"];
    const nonce_val_str    = output_map["nonce_val_str"];
    
    //--------------------------
    // USER_EXISTS - log them in
    if (user_exists_bool) {
        console.log("USER_EXISTS");
        
        //--------------------------
        // WEB3_SIGN
        const auth_signature_str = await sign(nonce_val_str, user_address_eth_str);
        
        //--------------------------
        
        // login this newly created user
        const login_output_map = await p_http_api_map["eth"]["user_login_fun"](user_address_eth_str, auth_signature_str);
        console.log(" ============== LOGIN_OUTPUT", login_output_map);
    }

    //--------------------------
    // NO_USER - in the system, offer to create new
    else {
        
        console.log("NO USER");

        // CREATE - create new user
        const new_user_create_data_map    = await user_create(null, user_address_eth_str, nonce_val_str, p_http_api_map);
        const new_user_auth_signature_str = new_user_create_data_map["auth_signature_str"];

        // LOGIN - login this newly created user
        const login_output_map = await p_http_api_map["eth"]["user_login_fun"](user_address_eth_str, new_user_auth_signature_str);
        console.log(" ============== LOGIN_OUTPUT", login_output_map);


        // only after that offer to the user to upload their details.
        // for update to succeed the user has to be logedin
        const user_data_map = await gf_ops.user_update_dialog();
    }

    //--------------------------
}

//-------------------------------------------------
async function user_create(p_username_str :string,
    p_user_address_eth_str :string,
    p_nonce_val_str,
    p_http_api_map) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        const create_user_dialog = $(`
            <div id='create_user_prompt_dialog'>
                <div id='descr'>create new user?</div>
                <div id='confirm_btn'>ok</div>
            </div>`);

        $("#identity").append(create_user_dialog);

        $(create_user_dialog).find("#confirm_btn").on('click', async ()=>{

            //--------------------------
            // WEB3_SIGN
            const auth_signature_str = await sign(p_nonce_val_str, p_user_address_eth_str);

            //--------------------------

            try {

                //--------------------------
                // USER_CREATE_HTTP
                // user was created successfuly, remove the create_user_dialog and return
                const http_output_map = await p_http_api_map["eth"]["user_create_fun"](p_user_address_eth_str, auth_signature_str);
                
                //--------------------------
                $(create_user_dialog).remove();


                const user_create_data_map = {
                    "http_output_map":    http_output_map,
                    "auth_signature_str": auth_signature_str,
                };
                p_resolve_fun(user_create_data_map);

            } catch (p_err) {            
                $(create_user_dialog).css("background-color", "red");
                p_reject_fun();
            }
        });
    });
    return p;
}

//-------------------------------------------------
async function sign(p_nonce_val_str :string,
    p_user_address_eth_str :string) {

    const auth_signature_str = await window.web3.eth.personal.sign(p_nonce_val_str, p_user_address_eth_str);    
    return auth_signature_str;
}

//-------------------------------------------------
export function wallet_connect() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        // providers like MetaMask and Status must continue to inject window.ethereum,
        // but now the window.ethereum object itself is a provider type that supports the 
        // methods defined in EIP-1102 and EIP-1193
        if (window.ethereum) {

            // IMPORTANT!! - if the users wallet is not connected to this page
            //               or the wallet is not unlocked, this async method will return only
            //               once all that is done. 
            const p = window.ethereum.send('eth_requestAccounts');
            p.then((p_r)=>{

                // get the first user address? or give the user choice 
                // over which address to use if multiple are returned?
                const user_address_eth_str = p_r["result"][0];

                window.web3 = new Web3(window.ethereum);

                // // wallet is connected, so remove the auth_pick_dialog
                // $("#auth_pick_dialog").remove();

                console.log("ETH_ADDRESS", user_address_eth_str);
                p_resolve_fun(user_address_eth_str);
            });
        } else {
            $("#identity").append(`<div id="wallet_connect_failed">failed to connect to wallet</div>`);
            p_reject_fun();
        }
    });
    return p;
}