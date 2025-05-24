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

///<reference path="./../../../d/jquery.d.ts" />

import * as gf_time from "./../../../gf_core/ts/gf_time";

//-------------------------------------------------
export async function init(p_http_api_map, p_log_fun) {


    console.log("admin dashboard")

    init_users_list(p_http_api_map, p_log_fun)
    init_invite_list(p_http_api_map, p_log_fun);
}

//-------------------------------------------------
function init_users_list(p_http_api_map, p_log_fun) {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
        const container = $(`
            <div id="users_list">
                <div id="title">users list</div>
                <div id="current">
                </div>
            </div>`);

        $("body").append(container);

        //--------------------------
        // GET_ALL

        // HTTP
        const output_map = await p_http_api_map["admin"]["get_all_users"]();
        const users_lst  = output_map["users_lst"];

        //-------------------------------------------------
        function view_user(p_user_map) {

            const user_id_str          = p_user_map["id_str"];
            const user_name_str        = p_user_map["user_name_str"];
            const email_str            = p_user_map["email_str"];
            const creation_unix_time   = p_user_map["creation_unix_time_f"];
            const email_confirmed_bool = p_user_map["email_confirmed_bool"];

            const user_element = $(`
                <div class="user">
                    <div class="user_name">${user_name_str}</div>
                    <div class="email">${email_str}</div>                    
                    <div class="delete_btn">delete</div>
                    <div class="creation_time">${creation_unix_time}</div>
                    <div>
                </div>`)[0];

            gf_time.init_creation_date(user_element, p_log_fun);

            // IMPORTANT!! - if email is not confirmed yet for a particular user,
            //               show the button allowing for resending of email_confirmation emails
            if (!email_confirmed_bool) {

                const resend_btn = $(`<div class="resend_email_confirm_btn">email confirm resend</div>`);
                $(user_element).append(resend_btn);

                $(resend_btn).on("click", async ()=>{
                    const output_map = await p_http_api_map["admin"]["resend_email_confirm"](user_id_str,
                        user_name_str,
                        email_str);

                    $(this).css("background-color", "green");
                });
            }
            else {
                const email_confirmed = $(`<div class="email_confirmed">email confirmed</div>`);
                $(user_element).append(email_confirmed);
            }
            
            // DELETE
            $(user_element).find(".delete_btn").on("click", async ()=>{

                const confirm_deletion_dialog = $(`<div class="confirm_deletion_dialog">
                    <div class="label">really want to delete user [${user_name_str}]</div>
                    <div class="confirm_btn">confirm</div>
                    <div class="close_btn">x</div>
                </div>`)
                $(user_element).append(confirm_deletion_dialog);

                $(confirm_deletion_dialog).find(".close_btn").on("click", async ()=>{
                    $(confirm_deletion_dialog).remove();
                });

                $(confirm_deletion_dialog).find(".confirm_btn").on("click", async ()=>{

                    // HTTP
                    const output_map = await p_http_api_map["admin"]["delete_user"](user_id_str,
                        user_name_str);
    
                    $(user_element).remove();
                });
                
            });

            return user_element;
        }

        //-------------------------------------------------

        for (const user_map of users_lst) {

            const user_element = view_user(user_map);
            $(container).find("#current").append(user_element);
        }

        //--------------------------
        p_resolve_fun(container);
    });
    return p;
}

//-------------------------------------------------
async function init_invite_list(p_http_api_map :any, p_log_fun: Function) {
    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

        const container = $(`
            <div id="invite_list">

                <div id="title">invite list</div>
                <div id="add_new">
                    <div id="new_email">
                        <input></input>
                    </div>
                    <div id="confirm_btn">add</div>
                </div>
                <div id="current">
                </div>
            </div>`);

        $("body").append(container);

        //--------------------------
        // ADD_NEW
        $(container).find("#add_new #confirm_btn").on("click", async ()=>{

            const new_email_str = $(container).find("#new_email input").val();
            const output_map    = await p_http_api_map["admin"]["add_to_invite_list"](new_email_str);


            const new_invite_element = $(`
                <div class="invite">
                    <div class="email">${new_email_str}</div>
                    <div class="creation_time">now</div>
                </div>`);
            $(container).find("#current").append(new_invite_element);

        });

        //--------------------------
        // GET_ALL

        // HTTP
        const output_map      = await p_http_api_map["admin"]["get_all_invite_list"]();
        const invite_list_lst = output_map["invite_list_lst"];

        //-------------------------------------------------
        function view_invite(p_invite_map: any) {

            const email_str          = p_invite_map["user_email_str"];
            const creation_unix_time = p_invite_map["creation_unix_time_f"];
            const invite_element = $(`
                <div class="invite">
                    <div class="email">${email_str}</div>
                    <div class="creation_time">${creation_unix_time}</div>
                    <div class="remove_btn">x</div>
                </div>`)[0];

            gf_time.init_creation_date(invite_element, p_log_fun);

            // REMOVE_INVITE
            $(invite_element).find(".remove_btn").on("click", async ()=>{
                const output_map = await p_http_api_map["admin"]["remove_from_invite_list"](email_str);

                $(invite_element).remove();
            });

            return invite_element;
        }

        //-------------------------------------------------

        for (const invite_map of invite_list_lst) {

            const invite_element = view_invite(invite_map);
            $(container).find("#current").append(invite_element);
        }

        //--------------------------
        p_resolve_fun(container);
    });
    return p;
}