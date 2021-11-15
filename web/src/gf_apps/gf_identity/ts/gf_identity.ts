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


        login(()=>{

            },
            ()=>{
                
            });
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
        }

        console.log("user address", p_user_address_eth_str);
    })

    
    /*//------------------------------------
    // connect to web3 wallet
    window.ethereum.enable();


    // check if web3 lib has been properly injected into the page context
    var web3;
    if (typeof window.Web3 != "undefined") {
        web3 = new Web3(Web3.givenProvider);
    } else {
        // fallback - if there is no web3 lib injected with its own provider,
        //            than try to reach out to a local eth node as the provider.
        web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:8545"))
    }

    // check if web3 lib is connected to wallet
    if(!web3.isConnected()) {
        $("#identity").append(`<div id="wallet_connect_failed">failed to connect to wallet</div>`);
    }

    const user_address_eth_str = web3.eth.accounts[0];

    //------------------------------------

    user_preflight__http(user_address_eth_str,
        (p_data_map)=>{
            
            const user_exists_bool = p_data_map["user_exists_bool"];

            if (user_exists_bool) {
                const user_nonce_val_str = p_data_map["nonce_val_str"];
            } else {

            }
        },
        ()=>{
            p_on_error_fun();
        });
    */
}

//-------------------------------------------------
function user_preflight__http(p_user_address_eth_str,
    p_on_complete_fun,
    p_on_error_fun) {

    const data_map = {
        "address_eth_str": p_user_address_eth_str,
    };

    const url_str = '/v1/identity/users/preflight';
    $.ajax({
        'url':         url_str,
        'type':        'POST',
        'data':        data_map,
        'contentType': 'application/json',
        'success':     (p_response_str)=>{
            const data_map :Object = JSON.parse(p_response_str);

            p_on_complete_fun(data_map);
        },
        'error': (jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}

//-------------------------------------------------
function user_login__http(p_on_complete_fun,
    p_on_error_fun) {

    const data_map = {

    };

    const url_str = '/v1/identity/users/login';
    $.ajax({
        'url':         url_str,
        'type':        'POST',
        'data':        data_map,
        'contentType': 'application/json',
        'success':     (p_response_str)=>{
            const data_map :Object = JSON.parse(p_response_str);
            
            p_on_complete_fun(data_map);
        },
        'error':(jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}

//-------------------------------------------------
function user_create__http(p_on_complete_fun,
    p_on_error_fun) {

    const data_map = {

    };

    const url_str = '/v1/identity/users/create';
    $.ajax({
        'url':         url_str,
        'type':        'POST',
        'data':        data_map,
        'contentType': 'application/json',
        'success':     (p_response_str)=>{
            const data_map :Object = JSON.parse(p_response_str);
            
            p_on_complete_fun(data_map);
        },
        'error':(jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}