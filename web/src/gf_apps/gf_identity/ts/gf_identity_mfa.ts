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

import * as gf_identity_http from "./gf_identity_http";

//-------------------------------------------------
export function init(p_user_name_str,
    p_http_api_map,
    p_on_mfa_validate_fun) {

    const container = $(`<div id="mfa_dialog">
        <div id="mfa_background"></div>
        <div id="mfa_confirm_code">
            <div id="label">MFA code</div>
            <input id="mfa_val" type="number"></input>
            <div id="confirm_btn">confirm</div>
        </div>
    </div>`);

    $(container).find("#confirm_btn").on("click", async ()=>{

        const mfa_val_str = $(container).find("#mfa_val input").val();
        const output_map  = await p_http_api_map["mfa"]["user_mfa_confirm"](p_user_name_str, mfa_val_str);

        

        const mfa_valid_bool = output_map["mfa_valid_bool"];
        if (mfa_valid_bool) {

            $(container).remove();
            p_on_mfa_validate_fun(true);
        } 
        else {
            // p_on_mfa_validate_fun(false);

            $(container).find("#mfa_confirm_code").append(`
                <div id="confirm_failed">
                    MFA confirm failed
                </div>
            `);
        }

    });
    return container;
}