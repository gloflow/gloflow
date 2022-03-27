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

///<reference path="../../../d/jquery.d.ts" />

import * as gf_time         from "./../../../gf_core/ts/gf_time";

//-------------------------------------------------
export async function init(p_http_api_map, p_log_fun) {


    console.log("admin dashboard")


    init_invite_list(p_http_api_map, p_log_fun);
}

//-------------------------------------------------
async function init_invite_list(p_http_api_map, p_log_fun) {
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
        function view_invite(p_invite_map) {

            // /v1/admin/users/remove_from_invite_list

            const email_str          = p_invite_map["user_email_str"];
            const creation_unix_time = p_invite_map["creation_unix_time_f"];
            const invite_element = $(`
                <div class="invite">
                    <div class="email">${p_invite_map["user_email_str"]}</div>
                    <div class="creation_time">${creation_unix_time}</div>
                    <div>

                    </div class="remove_btn">x</div>
                </div>`);

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