



///<reference path="../../../d/jquery.d.ts" />

import * as gf_admin from "./../ts/gf_admin";

//--------------------------------------------------------
$(document).ready(()=>{
	main();
});

//--------------------------------------------------------
function main() {

    var user_exists_bool = false;
    var login_pass_not_valid__returned_bool = false;
    var creation_allowed__returned_bool     = false;

    const http_api_map = {
        // ADMIN
        "admin": {
            "get_all_invite_list": async ()=>{
                const output_map = {
                    "invite_list_lst": [
                        {
                            "user_email_str": "it@tra.com",
                            "creation_unix_time_f": "46546987465"
                        }
                    ]
                };
                return output_map;
            },
            "add_to_invite_list": async (p_email_str :string)=>{
                const output_map = {};
                return output_map;
            }
        }
    };


    gf_admin.init(http_api_map);
}