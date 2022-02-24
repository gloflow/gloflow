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

//-------------------------------------------------
// USER_PREFLIGHT__HTTP
export function user_preflight(p_user_name_str,
    p_user_address_eth_str) {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_name_str":        p_user_name_str,
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
// USER_ETH_LOGIN__HTTP
export function user_eth_login(p_user_address_eth_str :string,
    p_auth_signature_str :string) {
    
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_address_eth_str": p_user_address_eth_str,
            "auth_signature_str":   p_auth_signature_str,
        };

        const url_str = '/v1/identity/users/login';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                p_resolve_fun(p_response_map);
            },
            'error':(jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//-------------------------------------------------
// USER_USERPASS_LOGIN__HTTP
export function user_userpass_login(p_user_name_str :string,
    p_pass_str  :string,
    p_email_str :string,
    p_url_str   :string) {
    
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_name_str": p_user_name_str,
            "pass_str":      p_pass_str,
            "email_str":     p_email_str,
        };

        $.ajax({
            'url':         p_url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                p_resolve_fun(p_response_map);
            },
            'error':(jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//-------------------------------------------------
export function user_mfa_confirm(p_user_name_str :string,
    p_mfa_val_str :string) {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_name_str": p_user_name_str,
            "mfa_val_str":   p_mfa_val_str,
        };

        const url_str = '/v1/identity/mfa_confirm';
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
// USER_CREATE__HTTP
export function user_eth_create(p_user_address_eth_str :string,
    p_auth_signature_str :string) {

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
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

// USER_CREATE__HTTP
export function user_userpass_create(p_user_name_str :string,
    p_pass_str :string) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_name_str": p_user_name_str,
            "pass_str":      p_pass_str,
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
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//-------------------------------------------------
// USER_UPDATE__HTTP
export function user_update(p_user_data_map) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_username_str":    p_user_data_map["username_str"],
            "user_email_str":       p_user_data_map["email_str"],
            "user_description_str": p_user_data_map["description_str"],
        };

        const url_str = '/v1/identity/users/update';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{
                
                p_resolve_fun(data_map);
            },
            'error': (jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}