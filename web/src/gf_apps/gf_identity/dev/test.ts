



///<reference path="../../../d/jquery.d.ts" />

import * as gf_identity from "./../ts/gf_identity";

//--------------------------------------------------------
$(document).ready(()=>{
	main();
});

//--------------------------------------------------------
function main() {

    var user_exists_bool = false;
    const http_api_map = {
        "eth": {
            "user_preflight_fun": async (p_user_address_eth_str)=>{
                const output_map = {
                    "user_exists_bool": user_exists_bool,
                    "nonce_val_str":    "random_string",
                };
                return output_map;
            },
            "user_login_fun": async (p_user_address_eth_str, p_auth_signature_str)=>{
                const output_map = {};
                return output_map;
            },
            "user_create_fun": async (p_user_address_eth_str, p_auth_signature_str)=>{
                const output_map = {};
                user_exists_bool = true;
                return output_map;
            },
            "user_update_fun": async (p_username_str,
                p_email_str,
                p_description_str)=>{

            }
        },
        "userpass": {
            "username_exists_fun": async (p_user_name_str)=>{
                const output_map = {
                    "user_exists_bool": user_exists_bool,
                };
                return output_map;
            },
            "user_login_fun": async (p_user_name_str, p_pass_str)=>{
                const output_map = {};
                return output_map;
            },
            "user_create_fun": async (p_user_name_str, p_pass_str)=>{
                const output_map = {};
                user_exists_bool = true;
                return output_map;
            },
            "user_update_fun": async (p_username_str,
                p_email_str,
                p_description_str)=>{

            }


        }
    };


    gf_identity.init(http_api_map);
}