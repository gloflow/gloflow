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



import * as gf_identity_eth      from "./gf_identity_eth";
import * as gf_identity_userpass from "./gf_identity_userpass";

declare const window: any;
declare var Web3;

//-------------------------------------------------
export async function init() {

    $("#identity #login").on('click', async function(p_e) {

        const method_str = await auth_method_pick();
        
        switch (method_str) {
            //--------------------------
            // ETH_METAMASK
            case "eth_metamask":
                await gf_identity_eth.user_auth_pipeline();
                break;
            
            //--------------------------
            // USER_AND_PASS
            case "userpass":
                await gf_identity_userpass.user_auth_pipeline();
                break;

            //--------------------------
        }
    });
}

//-------------------------------------------------
async function auth_method_pick() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {
        const auth_pick_dialog = $(`
        <div id="auth_pick_dialog">

            <div id="wallet_pick_dialog">
                <div id="metamask">
                    <div id="icon">
                        <img src="/images/static/assets/gf_metamask_icon.svg"></img>
                    </div>
                    <div id="descr">metamask browser wallet</div>
                </div>
            </div>

            <div id="user_and_pass>
                <div id="label">user and password</div>
            </div>
        </div>`);

        $("#identity").append(auth_pick_dialog);

        $(auth_pick_dialog).find("#metamask").on('click', ()=>{
            p_resolve_fun("eth_metamask");
        })

        $(auth_pick_dialog).find("#user_and_pass").on('click', ()=>{
            p_resolve_fun("userpass");
        });
    });
    return p;
}