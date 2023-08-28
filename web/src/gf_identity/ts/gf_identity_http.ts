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

import * as gf_admin_http from "./gf_admin_http";

//-------------------------------------------------
export function get_standard_http_urls() {
    const login_url_str = '/v1/identity/userpass/login';
    const home_url_str  = "/v1/home/view";
    const urls_map = {
        "login": login_url_str,
        "home":  home_url_str,
    };
    return urls_map;
}

//-------------------------------------------------
export function get_admin_http_urls() {
    const login_url_str = '/v1/admin/login';
    const home_url_str  = "/v1/admin/dashboard";
    const urls_map = {
        "login": login_url_str,
        "home":  home_url_str,
    };
    return urls_map;
}

//-------------------------------------------------
export function get_http_api(p_urls_map) {
    const http_api_map = {

        // ETH
        "eth": {
            "user_preflight_fun": async (p_user_address_eth_str)=>{
                const output_map = await user_preflight(null, p_user_address_eth_str);
                return output_map;
            },
            "user_login_fun": async (p_user_address_eth_str, p_auth_signature_str)=>{
                
                const output_map = await user_eth_login(p_user_address_eth_str,
                    p_auth_signature_str);
                return output_map;
            },
            "user_create_fun": async (p_user_address_eth_str, p_auth_signature_str)=>{
                const output_map = await user_eth_create(p_user_address_eth_str, p_auth_signature_str);
                return output_map;
            }
        },

        // USERPASS
        "userpass": {
            "user_login_fun": async (p_user_name_str, p_pass_str, p_email_str)=>{
                const url_str    = p_urls_map["login"];
                const output_map = await user_userpass_login(p_user_name_str,
                    p_pass_str,
                    p_email_str,
                    url_str);
                return output_map;
            },
            "user_create_fun": async (p_user_name_str, p_pass_str, p_email_str)=>{
                const output_map = await user_userpass_create(p_user_name_str,
                    p_pass_str,
                    p_email_str);
                return output_map;
            }
        },

        // MFA
        "mfa": {
            "user_mfa_confirm": async (p_user_name_str, p_mfa_val_str)=>{
                const output_map = await user_mfa_confirm(p_user_name_str,
                    p_mfa_val_str);
                return output_map;
            }
        },

        // ADMIN
        "admin": {
            "delete_user": async (p_user_id_str :string,
                p_user_name_str :string)=>{
                const output_map = await gf_admin_http.delete_user(p_user_id_str, p_user_name_str);
                return output_map;
            },
            "get_all_users": async ()=>{
                const output_map = await gf_admin_http.get_all_users();
                return output_map;
            },
            "get_all_invite_list": async ()=>{
                const output_map = await gf_admin_http.get_all_invite_list();
                return output_map;
            },
            "add_to_invite_list": async (p_email_str :string)=>{
                const output_map = await gf_admin_http.add_to_invite_list(p_email_str);
                return output_map;
            },
            "remove_from_invite_list": async (p_email_str :string)=>{
                const output_map = await gf_admin_http.remove_from_invite_list(p_email_str);
                return output_map;
            },

            "resend_email_confirm": async (p_user_id_str :string,
                p_user_name_str :string,
                p_email_str     :string)=>{
                const output_map = await gf_admin_http.resend_email_confirm(p_user_id_str,
                    p_user_name_str,
                    p_email_str);
                return output_map;
            },
        },

        "general": {
            "get_me": async ()=>{
                const output_map = await user_get_me();
                return output_map;
            },
            "logged_in": async ()=>{
                const output_map = await logged_in();
                return output_map;
            },
        }
    };
    return http_api_map;
}

//-------------------------------------------------
// ME
//-------------------------------------------------
export function user_get_me() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        // auth_r=0 - toggle off redirecting to login url in case the validation fails.
        //            redirecting not needed since /me is called by subcomponents of pages
        //            and on failure should not redirect.
        const url_str = '/v1/identity/me?auth_r=0';
        $.ajax({
            'url':         url_str,
            'type':        'GET',
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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

export function logged_in() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        // auth_r=0 - toggle off redirecting to login url in case the validation fails.
        //            redirecting not needed since /me is called by subcomponents of pages
        //            and on failure should not redirect.
        const url_str = '/v1/identity/logged_in?auth_r=0';
        $.ajax({
            'url':         url_str,
            'type':        'GET',
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
            'contentType': 'application/json',
            'success':     (p_response_map)=>{

                const status_str = p_response_map["status"];

                if (status_str == "OK") {
                    p_resolve_fun(true);
                } else {
                    p_resolve_fun(false);
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
// LOGIN
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
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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
            'error':(jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//-------------------------------------------------
// USER_MFA_CONFIRM
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
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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
// CREATE
//-------------------------------------------------
// USER_ETH_CREATE__HTTP
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
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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
// USER_USERPASS_CREATE__HTTP
export function user_userpass_create(p_user_name_str :string,
    p_pass_str  :string,
    p_email_str :string) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_name_str": p_user_name_str,
            "pass_str":      p_pass_str,
            "email_str":     p_email_str, 
        };

        const url_str = '/v1/identity/userpass/create';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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
// UPDATE
//-------------------------------------------------
// USER_UPDATE__HTTP
export function user_update(p_user_data_map) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_username_str":    p_user_data_map["username_str"],
            "user_email_str":       p_user_data_map["email_str"],
            "user_description_str": p_user_data_map["description_str"],
        };

        const url_str = '/v1/identity/update';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'cache':       false, // IMPORTANT!! - avoids various issues. always need latest response
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