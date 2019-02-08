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

import * as gf_sys_panel       from "./../../../../gf_core/ts/gf_sys_panel";
import * as gf_post_image_view from "./gf_post_image_view";
import * as gf_tagger_input_ui from "./../../../gf_tagger/ts/gf_tagger_client/gf_tagger_input_ui";

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
    p_log_fun('FUN_ENTER', 'gf_post.init()');
    
    gf_sys_panel.init(p_log_fun);

    const post_title_str :string      = $('#post_title').text();
    const post_tags_container_element = $('#post_tags_container');
    
    //------------------------------
    //INIT IMAGE TAGGING
    $('.post_element_image').each((p_i, p_post_element)=>{
        
        const image_element          = $(p_post_element).find('img');
        const img_url_str  :string   = $(image_element).attr('src');
        const path_lst     :string[] = img_url_str.split('/'); //Uri.parse(img_url_str).pathSegments;
        const img_file_str :string   = path_lst[path_lst.length-1];
        const tags_num_int :number   = $(p_post_element).find('.tags_container .gf_post_element_tag').length;
        
        //img_file_str example - "6c4a667457f05939af6a5f68690d0f55_thumb_medium.jpeg"
        const img_id_str :string = img_file_str.split('_')[0];
        p_log_fun('INFO', 'img_id_str - '+img_id_str);

        var tag_ui_added_bool :boolean = false;
        gf_tagger_input_ui.init_tag_input(img_id_str, //p_obj_id_str
            'image',    //p_obj_type_str
            p_post_element,
            //p_onTagsCreated_fun
            (p_added_tags_lst :string[])=>{
                view_added_tags(p_post_element, p_added_tags_lst, p_log_fun);
            },
            //p_onTagUIAdd_fun
            ()=>{
                tag_ui_added_bool = true;
            },
            //p_onTagUIRemove_fun
            ()=>{
                tag_ui_added_bool = false;
            },
            p_log_fun);

        gf_post_image_view.init(p_post_element, p_log_fun);

        $(p_post_element).on('mouseenter', (p_event)=>{

            //IMPORTANT!! - only show tags_container if there are tags attached to this post_element
            if (tags_num_int > 0) {
                $(p_post_element).find('.tags_container').css('visibility', 'visible');
            }
        });
        $(p_post_element).on('mouseleave', (p_event)=>{

            //hide the tags_container only if the tagging UI is not open. 
            //if it is open we want the tags_container visible so that we can 
            //see the tags as they're added
            if (!tag_ui_added_bool) {
                $(p_post_element).find('.tags_container').css("visibility", 'hidden');
            }
        });
    });
    //------------------------------
    //VIDEO TAGGING

    $('.post_element_video').each((p_i, p_post_element)=>{

        //ADD!! - extract video ID properly
        gf_tagger_input_ui.init_tag_input('fix',
            'video',
            p_post_element,
                    
            //p_onTagsCreated_fun
            (p_added_tags_lst :string[])=>{
                view_added_tags(p_post_element, p_added_tags_lst, p_log_fun);
            },
            ()=>{}, //p_onTagUIAdd_fun
            ()=>{}, //p_onTagUIRemove_fun
            p_log_fun);
    });
    //------------------------------

    //final List<String> tags_lst = queryAll('.post_tag').map((p_element) => p_element.text);
    //gf_post_tag_mini_view.init_tags_mini_view(tags_lst,
    //                                          p_log_fun);
}
//-----------------------------------------------------
function view_added_tags(p_post_element,
    p_added_tags_lst :string[],
    p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_post.view_added_tags()');

    const tags_container_element = $(p_post_element).find('.tags_container');

    for (var tag_str of p_added_tags_lst) {
        const tag_url_str :string = '/tags/view_objects?tag='+tag_str+'&otype=image';
        const new_tag_ui_element  = $('<a class="gf_post_element_tag">'+tag_str+'</a>');
        $(new_tag_ui_element).attr('href',tag_url_str);

        //IMPORTANT!! - add the new tag link to the DOM
        $(tags_container_element).append(new_tag_ui_element);
    }
}
//-----------------------------------------------------
function get_post_element_tags_num(p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_post.get_post_element_tags_num()');
    
}