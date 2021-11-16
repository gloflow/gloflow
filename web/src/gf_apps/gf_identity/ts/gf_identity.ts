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
export function init() {

    $("#identity #login").on('click', function(p_e) {



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

            login(()=>{

                },
                ()=>{
                    
                });
        })

        
    });
}

//-------------------------------------------------
function login(p_on_complete_fun,
    p_on_error_fun) {

    //-------------------------------------------------
    /* providers like MetaMask and Status must continue to inject window.ethereum,
    but now the window.ethereum object itself is a provider type that supports the 
    methods defined in EIP-1102 and EIP-1193*/
    function eth_is_enabled(p_on_complete_fun) {
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
                p_on_complete_fun(true, user_address_eth_str);
            });
        } else {
            p_on_complete_fun(false, null);
        }
    }

    //-------------------------------------------------
    

    eth_is_enabled((p_enabled_bool, p_user_address_eth_str)=>{

        console.log(p_enabled_bool);

        if (!p_enabled_bool) {
            $("#identity").append(`<div id="wallet_connect_failed">failed to connect to wallet</div>`);
        } else {


            console.log("user address", p_user_address_eth_str);

            user_auth_pipeline(p_user_address_eth_str,
                ()=>{},
                ()=>{});
        }
    });
}

//-------------------------------------------------
function user_auth_pipeline(p_user_address_eth_str,
    p_on_complete_fun,
    p_on_error_fun) {

    // HTTP_REQUEST
    user_preflight__http(p_user_address_eth_str,
        (p_data_map)=>{
            
            const user_exists_bool = p_data_map["user_exists_bool"];
            const nonce_val_str    = p_data_map["nonce_val_str"];
            

            // user exists, log them in
            if (user_exists_bool) {
                
            }
            // no-user in the system, offer to create new
            else {
                
                console.log("NO USER");

                create_new_user(p_user_address_eth_str,
                    nonce_val_str,
                    (p_data_map)=>{

                        // user is created, allow them to update basic profile information



                    },
                    (p_error_data_map)=>{
                        p_on_error_fun(p_error_data_map);
                    });
            }
        },
        (p_error_data_map)=>{
            p_on_error_fun(p_error_data_map);
        });
}

//-------------------------------------------------
function create_new_user(p_user_address_eth_str,
    p_nonce_val_str,
    p_on_complete_fun,
    p_on_error_fun) {



    



    const create_user_dialog = $(`
        <div id='create_user_dialog'>
            <div id='descr'>create new user?</div>
            <div id='confirm'>ok</div>
        </div>`);


    $("#identity").append(create_user_dialog);

    $(create_user_dialog).find("#confirm").on('click', ()=>{


        const s = window.ethereum.personal.sign(p_nonce_val_str, p_user_address_eth_str).then((p_auth_signature_str)=>{
            

            user_create__http(p_user_address_eth_str,
                p_auth_signature_str,
                (p_data_map)=>{
                    
                    // user was created successfuly, remove the create_user_dialog and return
                    $(create_user_dialog).remove();
                    p_on_complete_fun();
                },
                (p_error_data_map)=>{
                    
                    $(create_user_dialog).css("background-color", "red");
                });
            


        });


        

    });
}

//-------------------------------------------------
// USER_PREFLIGHT__HTTP
function user_preflight__http(p_user_address_eth_str,
    p_on_complete_fun,
    p_on_error_fun) {

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
                p_on_complete_fun(data_map);
            } else {
                p_on_error_fun(data_map);
            }
        },
        'error': (jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}

//-------------------------------------------------
// USER_LOGIN__HTTP
function user_login__http(p_on_complete_fun,
    p_on_error_fun) {

    const data_map = {

    };

    const url_str = '/v1/identity/users/login';
    $.ajax({
        'url':         url_str,
        'type':        'POST',
        'data':        JSON.stringify(data_map),
        'contentType': 'application/json',
        'success':     (p_response_map)=>{
            
            p_on_complete_fun(data_map);
        },
        'error':(jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}

//-------------------------------------------------
// USER_CREATE__HTTP
function user_create__http(p_user_address_eth_str,
    p_auth_signature_str,
    p_on_complete_fun,
    p_on_error_fun) {

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
            
            p_on_complete_fun(data_map);
        },
        'error':(jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}