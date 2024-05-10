/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

// ///<reference path="../../d/jquery.d.ts" />

//-------------------------------------------------
// GET_ALL_USERS
export function delete_user(p_user_id_str :string,
    p_user_name_str :string) {
    return new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = { 
            "user_id_str":   p_user_id_str,
            "user_name_str": p_user_name_str
        };

        const url_str = '/v1/admin/users/delete';
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
}

//-------------------------------------------------
// GET_ALL_USERS
export function resend_email_confirm(p_user_id_str :string,
    p_user_name_str :string,
    p_email_str     :string) {
    return new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "user_id_str":   p_user_id_str,
            "user_name_str": p_user_name_str,
            "email_str":     p_email_str
        };

        const url_str = '/v1/admin/users/resend_confirm_email';
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
}

//-------------------------------------------------
// GET_ALL_USERS
export function get_all_users() {
    return new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {};

        const url_str = '/v1/admin/users/get_all';
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
}

//-------------------------------------------------
// GET_ALL_INVITE_LIST
export function get_all_invite_list() {
    return new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {};

        const url_str = '/v1/admin/users/get_all_invite_list';
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
}

//-------------------------------------------------
export function add_to_invite_list(p_email_str :string) {
    return new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "email_str": p_email_str,
        };

        const url_str = '/v1/admin/users/add_to_invite_list';
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
}

//-------------------------------------------------
export function remove_from_invite_list(p_email_str :string) {
    return new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "email_str": p_email_str,
        };

        const url_str = '/v1/admin/users/remove_from_invite_list';
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
}