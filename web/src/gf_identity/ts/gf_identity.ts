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

import * as gf_identity_auth0    from "./gf_identity_auth0";
import * as gf_identity_eth      from "./gf_identity_eth";
import * as gf_identity_userpass from "./gf_identity_userpass";
import * as gf_identity_http     from "./gf_identity_http";
import * as gf_utils             from "../../gf_core/ts/gf_utils";

declare const window: any;
declare var Web3;

//-------------------------------------------------
export async function init_me_control(p_parent_node,
    p_auth_http_api_map,
    p_home_url_str) {
    
    var me_user_map;
    try {
        me_user_map = await p_auth_http_api_map["general"]["get_me"]();

    } catch (error_map) {

        // failed to run the /me endpoint.
        // user is likely not logged in.
        return
    }

    const user_profile_img_url_str = me_user_map["profile_image_url_str"];
    const user_name_str            = me_user_map["user_name_str"];

    const auth_me_container = $(`
        <div id="auth_me">
            <div id="current_user">
                
            </div>
        </div>`);
    $(p_parent_node).append(auth_me_container);
    
    // IMG
    if (user_profile_img_url_str != "") { 
        $(auth_me_container).find("#current_user").append(`
            <img src="${user_profile_img_url_str}"></img>
        `);
    }

    // TEXT_SHORTHAND
    else {

        const shorthand_str = user_name_str[0];
        $(auth_me_container).find("#current_user").append(`
            <div id="shorthand_username">${shorthand_str}</div>
        `);
    }

    // HOME_REDIRECT
    $("#current_user").on("click", ()=>{

        // IMPORTAN!! - adding a unique param to this request to disable browser cache,
        //              since it can cause inconsistent behavior.
        const unique_param = new Date().getTime();
        const url_str = p_home_url_str+"?"+unique_param;
        
        window.location.href = url_str;
    });
}

//-------------------------------------------------
// INIT_WITH_HTTP
// p_notifications_meta_map - map of notifications metadata.
//                            allows external users of gf_identity functionality
//                            to customize text and various user messages that are displayed.
//                            example: different for admins and regular users.
export async function init_with_http(p_notifications_meta_map,
    p_urls_map) {
    
    const http_api_map = gf_identity_http.get_http_api(p_urls_map);
    init(p_notifications_meta_map, http_api_map, p_urls_map);
}

//-------------------------------------------------
export async function init(p_notifications_meta_map, p_http_api_map, p_urls_map) {

    const logged_in_bool = await p_http_api_map["general"]["logged_in"]();

    $("#identity #login").on("click", async function(p_e) {
        
        const opened_bool = $("#identity").has("#auth_pick_dialog").length > 0;

        if (!opened_bool) {
            
            const method_str = await auth_method_pick(logged_in_bool);
            switch (method_str) {
                //--------------------------
                // AUTH0
                case "auth0":
                    await gf_identity_auth0.user_auth_pipeline();
                    break;

                //--------------------------
                // ETH_METAMASK
                case "eth_metamask":
                    await gf_identity_eth.user_auth_pipeline(p_http_api_map);
                    break;
                
                //--------------------------
                // USER_AND_PASS
                case "internal_userpass":
                    await gf_identity_userpass.user_auth_pipeline(p_notifications_meta_map,
                        p_http_api_map,
                        p_urls_map);
                    break;

                //--------------------------
            }
        }
        else {
            $("#identity").find("#auth_pick_dialog").remove();
        }
    });
}

//-------------------------------------------------
async function auth_method_pick(p_logged_in_bool) {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        const container_identity = $("#identity");
        const container = $(`
        <div id="auth_pick_dialog">

            <div id="auth0_pick_dialog">
                <div id="label">classic</div>
            </div>

            <div id="wallet_pick_dialog">
                <div id="metamask">
                    <div id="icon">
                        <img src="https://gloflow.com/images/static/assets/gf_metamask_icon.svg"></img>
                    </div>
                    <div id="descr">metamask browser wallet</div>
                </div>
            </div>

            <!--
            // INTERNAL_USERNAME_PASSWORD_AUTH_METHOD
            <div id="internal_username_and_pass_pick_dialog">
                <div id="label">username and password</div>
            </div>
            -->

        </div>`);

        // if user is already logged in, the auth0 auth method should become the log-out button
        if (p_logged_in_bool) {
            $(container).find("#label").text("log out");
        }

        container_identity.append(container);

        // close_dialog
        gf_utils.click_outside(container_identity, ()=>{
            $(container).remove();
        });
        
        
        $(container).find("#auth0_pick_dialog").on('click', (p_e)=>{
            p_e.stopPropagation();

            // if user is already logged in, and logout button is pressed, redirect user
            // to logout endpoint so that the system logs them out. 
            if (p_logged_in_bool) {
                gf_identity_auth0.logout();
                return
            }

            // remove initial auth method pick dialog
            $(container).remove();

            p_resolve_fun("auth0");
        });

        $(container).find("#metamask").on('click', (p_e)=>{
            p_e.stopPropagation();

            // user picked to auth with wallet so remove initial p
            // auth method pick dialog.
            // $(container).remove();

            p_resolve_fun("eth_metamask");
        });

        /*
        
        // ADD!! - for now auth0 is the only enabled non-web3 auth method;
        //         but instead this setting (if internal username/pass method is enabled or Auth0)
        //         should be loaded from the backend (on the backend its set on startup via ENV var).

        $(container).find("#internal_username_and_pass_pick_dialog").on('click', (p_e)=>{
            p_e.stopPropagation();

            // remove initial auth method pick dialog
            $(container).remove();

            p_resolve_fun("internal_userpass");
        });
        */
    });
    return p;
}