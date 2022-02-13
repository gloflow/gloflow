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
// import * as gf_3d from "./../../../gf_core/ts/gf_3d";

//-------------------------------------------------
export async function user_auth_pipeline(p_http_api_map) {

    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const container = $(`
        <div id="user_and_pass_dialog">
            <div id="username_input">
                <input id='username_input' placeholder="user name"></input>
            </div>
            <div id="pass_input">
                <input id="pass_input" placeholder="password" type="password"></input>
            </div>
            <div id="email_input">
                <input id='email_input' placeholder="email"></input>
            </div>
            <div id="login_btn">login</div>
            <div id="create_btn">create user</div>
        </div>`);
        $("#identity").append(container);

        // gf_3d.div_follow_mouse($(container)[0], document, 90);

        $(container).find("input#username_input").focus();

        
        //-------------------------------------------------

        $(container).on('keyup', function (e) {
            if (e.key === 'Enter' || e.keyCode === 13) {
                login_activate(container, p_http_api_map);
            }
        });

        $(container).find("#login_btn").on('click', async ()=>{
            login_activate(container, p_http_api_map);
        });

        $(container).find("#create_btn").on('click', async ()=>{

            create_activate(container, p_http_api_map);
        });
    });
    return p;
}

//-------------------------------------------------
async function create_activate(p_container,
    p_http_api_map) {
        

    // ADD!! - do frontend validation on username
    const user_name_str = $(p_container).find("#username_input input").val();
    const pass_str      = $(p_container).find("#pass_input input").val();
    const email_str     = $(p_container).find("#email_input input").val();

    // remove all previous errors that were displayed
    $(p_container).find(".error").remove();

    // ERROR
    if (user_name_str == "" || pass_str == "" || email_str == "") {
        const error = $(`<div id="error_empty_field_dialog" class="error">
            <div class="label">username or password or email field is empty</div>
        </div>`);
        $(error).css("opacity", "0.0");
        $(p_container).append(error);

        $(error).animate({
            "opacity": "1.0"
        }, 200, ()=>{});
        return;
    }

    const create_output_map = await p_http_api_map["userpass"]["user_create_fun"](user_name_str,
        pass_str as string,
        email_str as string);


    // ERROR
    // user already exist
    const user_exists_bool = create_output_map["user_exists_bool"];
    if (user_exists_bool) {
        const error = $(`<div id="error_creation_existing_user_dialog" class="error">
            <div class="label">user with this username already exists. please pick a different username</div>
        </div>`);
        $(error).css("opacity", "0.0");
        $(p_container).append(error);

        $(error).animate({
            "opacity": "1.0"
        }, 200, ()=>{});
        return;
    }

    // ERROR
    // user not on the invite-list
    const user_in_invite_list_bool = create_output_map["user_in_invite_list_bool"];
    if (!user_in_invite_list_bool) {
        const error = $(`<div id="error_create_user_not_allowed_dialog" class="error">
            <div class="label">user with this username/email has not yet been added to the invite-list</div>
        </div>`);
        $(error).css("opacity", "0.0");
        $(p_container).append(error);

        $(error).animate({
            "opacity": "1.0"
        }, 200, ()=>{});
        return;
    }
    
}
//-------------------------------------------------
async function login_activate(p_container,
    p_http_api_map) {

    console.log("login activate");

    const user_name_str = $(p_container).find("#username_input input").val();
    const pass_str      = $(p_container).find("#pass_input input").val();
    const email_str     = $(p_container).find("#email_input input").val();


    // remove all previous errors that were displayed
    $(p_container).find(".error").remove();

    // ERROR
    if (user_name_str == "" || pass_str == "") {
        const error = $(`<div id="error_empty_field_dialog" class="error">
            <div class="label">username or password field is empty</div>
        </div>`);
        $(error).css("opacity", "0.0");
        $(p_container).append(error);

        $(error).animate({
            "opacity": "1.0"
        }, 200, ()=>{});
        return;
    }

    const login_output_map = await p_http_api_map["userpass"]["user_login_fun"](user_name_str,
        pass_str as string,
        email_str as string);

    // ERROR
    // user doesnt exist
    const user_exists_bool = login_output_map["user_exists_bool"];
    if (!user_exists_bool) {
        const error = $(`<div id="error_login_no_user_dialog" class="error">
            <div class="label">no user for this username</div>
        </div>`);
        $(error).css("opacity", "0.0");
        $(p_container).append(error);

        $(error).animate({
            "opacity": "1.0"
        }, 200, ()=>{});
        return;
    }

    const pass_valid_bool = login_output_map["pass_valid_bool"];
    if (!pass_valid_bool) {
        const error = $(`<div id="error_login_pass_not_valid_dialog" class="error">
            <div class="label">password is not correct</div>
        </div>`);
        $(error).css("opacity", "0.0");
        $(p_container).append(error);

        $(error).animate({
            "opacity": "1.0"
        }, 200, ()=>{});
        return;
    }
}