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

    $("#identity #login").on("click", async function(p_e) {
        
        const opened_bool = $("#identity").has("#auth_pick_dialog").length > 0;

        if (!opened_bool) {
            
            const method_str = await auth_method_pick();
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
async function auth_method_pick() {
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

        container_identity.append(container);

        // close_dialog
        gf_utils.click_outside(container_identity, ()=>{
            $(container).remove();
        });
        
        
        $(container).find("#auth0_pick_dialog").on('click', (p_e)=>{
            p_e.stopPropagation();

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