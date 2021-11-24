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

///<reference path="../../../d/jquery.d.ts" />

declare const window: any;
declare var Web3;

//-------------------------------------------------
export async function init() {

    $("#identity #login").on('click', async function(p_e) {


        await wallet_pick();
        
        const user_address_eth_str = await wallet_connect();

        await user_auth_pipeline(user_address_eth_str);
        
    });
}

//-------------------------------------------------
async function wallet_pick() {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const wallet_pick_dialog = $(`
            <div id="wallet_pick_dialog">
                <div id="metamask">
                    <div id="icon">
                        <img src="/images/static/assets/gf_metamask_icon.svg"></img>
                    </div>
                    <div id="descr">metamask browser wallet</div>
                </div>
            </div>`);

        $("#identity").append(wallet_pick_dialog);

        $(wallet_pick_dialog).find("#metamask").on('click', ()=>{

            p_resolve_fun(null);
        })

    });
    return p;
}

//-------------------------------------------------
function wallet_connect() {
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

                // wallet is connected, so remove the wallet_pick_dialog
                $("#wallet_pick_dialog").remove();

                p_resolve_fun(user_address_eth_str);
            });
        } else {
            $("#identity").append(`<div id="wallet_connect_failed">failed to connect to wallet</div>`);
            p_reject_fun();
        }
    });
    return p;
}

//-------------------------------------------------
async function user_auth_pipeline(p_user_address_eth_str) {

    // HTTP_REQUEST
    const data_map = await user_preflight__http(p_user_address_eth_str);

    const user_exists_bool = data_map["user_exists_bool"];
    const nonce_val_str    = data_map["nonce_val_str"];
    

    // user exists, log them in
    if (user_exists_bool) {
        
    }
    // no-user in the system, offer to create new
    else {
        
        console.log("NO USER");

        // user is created
        await user_create(p_user_address_eth_str, nonce_val_str);


        await user_update();
    }   
}

//-------------------------------------------------
async function user_create(p_user_address_eth_str,
    p_nonce_val_str) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        const create_user_dialog = $(`
            <div id='create_user_dialog'>
                <div id='descr'>create new user?</div>
                <div id='confirm_btn'>ok</div>
            </div>`);

        $("#identity").append(create_user_dialog);

        $(create_user_dialog).find("#confirm_btn").on('click', async ()=>{

            const auth_signature_str = await window.web3.eth.personal.sign(p_nonce_val_str, p_user_address_eth_str);
            
            try {
                // user was created successfuly, remove the create_user_dialog and return
                const user_create_data_map = await user_create__http(p_user_address_eth_str, auth_signature_str);
                $(create_user_dialog).remove();
                p_resolve_fun(null);

            } catch (p_err) {            
                $(create_user_dialog).css("background-color", "red");
                p_reject_fun();
            }
        });
    });
    return p;
}

//-------------------------------------------------
function user_update() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {


        const update_user_dialog = $(`
            <div id='update_user_dialog'>
                <div id='descr'>set your user details</div>
                <input id='username'></input>
                <input id='email'></input>
                <input id='description'></input>
            </div>`);
        $("#identity").append(update_user_dialog);

    });
    return p;
}

//-------------------------------------------------
// USER_PREFLIGHT__HTTP
function user_preflight__http(p_user_address_eth_str) {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_address_eth_str": p_user_address_eth_str,
        };

        const url_str = '/v1/identity/users/preflight';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                const status_str = p_response_map["status"];
                const data_map   = p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//-------------------------------------------------
// USER_LOGIN__HTTP
function user_login__http() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {

        };

        const url_str = '/v1/identity/users/login';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                p_resolve_fun(data_map);
            },
            'error':(jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//-------------------------------------------------
// USER_CREATE__HTTP
function user_create__http(p_user_address_eth_str,
    p_auth_signature_str) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_address_eth_str": p_user_address_eth_str,
            "auth_signature_str":   p_auth_signature_str,
        };

        const url_str = '/v1/identity/users/create';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                p_resolve_fun(data_map);
            },
            'error':(jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}