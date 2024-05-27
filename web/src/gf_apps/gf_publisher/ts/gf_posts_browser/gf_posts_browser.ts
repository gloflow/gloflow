/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

///<reference path="../../../../d/jquery.d.ts" />

import * as gf_core_utils         from "./../../../../gf_core/ts/gf_utils";
import * as gf_user_events        from "./../../../../gf_events/ts/gf_user_events";
import * as gf_sys_panel          from "./../../../../gf_sys_panel/ts/gf_sys_panel";
import * as gf_identity           from "./../../../../gf_identity/ts/gf_identity";
import * as gf_identity_http      from "./../../../../gf_identity/ts/gf_identity_http";

import * as gf_events             from "./../gf_posts_core/gf_events";
import * as gf_posts_browser_view from "./gf_posts_browser_view";

declare var gf_tagger__init_ui_v2;
declare var gf_tagger__http_add_tags_to_obj;

//-----------------------------------------------------
$(document).ready(()=>{
    //-------------------------------------------------
    function log_fun(p_g :string, p_m :string) {
        var msg_str = p_g+':'+p_m
        //chrome.extension.getBackgroundPage().console.log(msg_str);

        switch (p_g) {
            case "INFO":
                console.log("%cINFO"+":"+"%c"+p_m,"color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
                break;
            case "FUN_ENTER":
                console.log("%cFUN_ENTER"+":"+"%c"+p_m,"color:yellow; background-color:lightgray","background-color:lightgray");
                break;
        }
    }
    //-------------------------------------------------

    init(log_fun);
});

//-----------------------------------------------------
export async function init(p_log_fun :any) {

    const current_host_str = gf_core_utils.get_current_host();

    //---------------------
	// META
	const events_enabled_bool = true;
	const notifications_meta_map = {
		"login_first_stage_success": "login success"
	};

    //---------------------
	// IDENTITY
	// first complete main initialization and only then initialize gf_identity
	const urls_map          = gf_identity_http.get_standard_http_urls(current_host_str);
	const auth_http_api_map = gf_identity_http.get_http_api(urls_map, current_host_str);
	gf_identity.init_with_http(notifications_meta_map, urls_map, current_host_str);
	

	
	const parent_node = $("#right_section");
	const home_url_str = urls_map["home"];

	gf_identity.init_me_control(parent_node,
		auth_http_api_map,
		home_url_str);
	
	// inspect if user is logged-in or not
	const logged_in_bool = await auth_http_api_map["general"]["logged_in"]();

    //---------------------
	// EVENTS
	if (events_enabled_bool && logged_in_bool) {
		
		const event_meta_map = {

		};
		gf_user_events.send_event_http(gf_events.GF_POSTS_BROWSER_PAGE_LOAD,
			"browser",
			event_meta_map,
			current_host_str)
	}

    //---------------------


    const http_api_map = {

		// GF_TAGGER
		"gf_tagger": {
			"add_tags_to_obj": async (p_new_tags_lst :string[],
				p_obj_id_str    :string,
				p_obj_type_str  :string,
				p_tags_meta_map :any,
				p_log_fun       :any)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

					await gf_tagger__http_add_tags_to_obj(p_new_tags_lst,
						p_obj_id_str,
						p_obj_type_str,
						{}, // meta_map
						current_host_str,
						p_log_fun);

					p_resolve_fun({
						"added_tags_lst": p_new_tags_lst,
					});
				});
				return p;
			}
		}
	};


    // this app assumes that the first page of the posts is present 
    // in the dom on app startup... subsequent page loads happen from the server
    const static_posts_infos_lst :Object[] = load_data_from_dom(p_log_fun);

    gf_sys_panel.init_with_auth(p_log_fun);

    gf_posts_browser_view.init(static_posts_infos_lst,
        http_api_map,
        p_log_fun);
}

//-----------------------------------------------------
// DATA LOADING
//-----------------------------------------------------
function load_data_from_dom(p_log_fun) {
        
    const page_posts_infos_lst :Object[] = [];

    $('body #gf_posts_container').find('.gf_post').each((p_i, p_post)=>{

        const post_id_str :string = $(p_post).data('sys_id');
        const post_title_str :string = $(p_post).find('.post_title').text().trim();
        const post_url_str   :string = '/posts/'+post_title_str;

        //---------------------
        // TAGS
        const tags_infos_lst :Object[] = [];
        $(p_post).find('.gf_post_tag').each((p_i,p_tag_element)=>{
            const tag_str     :string = $(p_tag_element).text();
            const tag_url_str :string = $(p_tag_element).attr('href');

            const tag_info_map = {
                'tag_str':     tag_str,
                'tag_url_str': tag_url_str
            };

            tags_infos_lst.push(tag_info_map);
        });

        //--------------------
        // THUMBNAIL URL's
        var thumbnail_url_str :string = $(p_post).find('img').attr('src');
        if (thumbnail_url_str == '' || thumbnail_url_str == 'error') thumbnail_url_str = null;

        //--------------------

        const post_info_map = {
            'post':              p_post,
            'post_id_str':       post_id_str,
            'post_title_str':    post_title_str, 
            'post_url_str':      post_url_str,
            'tags_infos_lst':    tags_infos_lst,
            'thumbnail_url_str': thumbnail_url_str
        };

        page_posts_infos_lst.push(post_info_map);
    });

    return page_posts_infos_lst;
}