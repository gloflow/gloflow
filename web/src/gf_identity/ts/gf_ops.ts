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

// ///<reference path="../../d/jquery.d.ts" />

import * as gf_3d from "../../gf_core/ts/gf_3d";
import * as gf_identity_http from "./gf_identity_http";

//-------------------------------------------------
// UPDATE
export function user_update_dialog() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        const update_user_dialog = $(`
            <div id='update_user_dialog'>
                <div id='dialog_label'>set your user details</div>
                <div id='username'>
                    <div class='label'>username</div>
                    <input id='username_input'></input>
                </div>
                <div id='email'>
                    <div class='label'>email</div>
                    <input id='email_input'></input>
                </div>
                <div id='description'>
                    <div class='label'>description</div>
                    <textarea id='description_input' rows="4" cols="50"></textarea>
                </div>
                <div id='confirm_btn'>ok</div>
            </div>`);

        $("#identity").append(update_user_dialog);

        // gf_3d.div_follow_mouse($(update_user_dialog)[0], document, 90);

        $(update_user_dialog).find("#confirm_btn").on('click', async ()=>{

            const username_str    = $(update_user_dialog).find("#username_input").val();
            const email_str       = $(update_user_dialog).find("#email_input").val();
            const description_str = $(update_user_dialog).find("#description_input").val();

            const data_map = {
                "username_str":    username_str,
                "email_str":       email_str,
                "description_str": description_str,
            };

            //--------------------------
            // USER_UPDATE_HTTP
            await gf_identity_http.user_update(data_map);

            //--------------------------
            $(update_user_dialog).remove();

            p_resolve_fun(data_map);
        });
    });
    return p;
}