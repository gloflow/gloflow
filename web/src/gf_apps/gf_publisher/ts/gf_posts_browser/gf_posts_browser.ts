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

import * as gf_sys_panel          from "./../../../../gf_sys_panel/ts/gf_sys_panel";
import * as gf_posts_browser_view from "./gf_posts_browser_view";

declare var gf_tagger__init_ui_v2;
declare var gf_tagger__http_add_tags_to_obj;

//-----------------------------------------------------
$(document).ready(()=>{
    //-------------------------------------------------
    function log_fun(p_g,p_m) {
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
export function init(p_log_fun) {

    //-------------------
    // HOST
    const domain_str   = window.location.hostname;
	const protocol_str = window.location.protocol;
	const gf_host_str = `${protocol_str}//${domain_str}`;
	console.log("gf_host", gf_host_str);

    //-------------------

    const http_api_map = {

		// GF_TAGGER
		"gf_tagger": {
			"add_tags_to_obj": async (p_new_tags_lst,
				p_obj_id_str,
				p_obj_type_str,
				p_tags_meta_map,
				p_log_fun)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

					await gf_tagger__http_add_tags_to_obj(p_new_tags_lst,
						p_obj_id_str,
						p_obj_type_str,
						{}, // meta_map
						gf_host_str,
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