/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

import * as gf_utils from "./gf_utils";
import * as gf_image_colors from "./../../../../gf_core/ts/gf_image_colors";

//---------------------------------------------------
// CREATE

export async function create(p_post_map,
    p_http_api_map,
    p_log_fun) {

    const id_str                  :string   = p_post_map['id_str'];
    const title_str               :string   = p_post_map['title_str'];
    const image_thumbnail_url_str :string   = p_post_map['thumbnail_url_str'];
    const images_number_str       :string   = p_post_map['images_number_str'];
    const creation_date_str       :string   = p_post_map['creation_datetime_str'];
    const tags_lst                :string[] = p_post_map['tags_lst'];
    const post_url_str            :string   = p_post_map["post_url_str"];

    // IMPORTANT!! - "item" class is used by Masonry
    const post :HTMLDivElement = <HTMLDivElement> $(`
        <div class='gf_post item gf_post_image_view'
            data-sys_id='${id_str}'>

            <div class='post_title'>${title_str}</div>

            <div class='post_images_number'>
                <div class='num'>${images_number_str}</div>
                <div class='label'>images #</div>
            </div>

            <div>
                <a class="post_image" target="_blank" href="${post_url_str}">
                    <img class='thumb_small_url' src="${image_thumbnail_url_str}"></img>
                </a>
            </div>

            <div class='gf_post_creation_date'>${creation_date_str}</div>
            <div class='tags_container'></div>
        </div>`)[0];

    // TAGS
    const tags_container = $(post).find('.tags_container');

    for (var tag_str of tags_lst) {

        const a :HTMLAnchorElement = <HTMLAnchorElement> $('<a class="gf_post_tag" href="/v1/tags/objects?tag='+tag_str+'&otype=post">#'+tag_str+'</a>')[0];
        $(tags_container).append(a);
    }

    //-------------------------
    // INIT_POST

    await init_post(id_str,
        post,
        p_http_api_map,
        p_log_fun);

    //-------------------------

    return post;
}

//---------------------------------------------------
// INIT_EXISTING_DOM
/*
used for templates usually, where the image element DOM structure is already
created server side when loaded into the browser, and just needs to be initialized
(no creation of the DOM tree for the image control)
*/

export async function init_existing_dom(p_post_element,
    p_http_api_map,
    p_log_fun) {

    /*
    init_posts_img_num(p_post_element);

    //----------------------
    // IMAGE_PALLETE
    const img = $(p_post_element).find("img")[0];

    const assets_paths_map = {
        "copy_to_clipboard_btn": "/images/static/assets/gf_copy_to_clipboard_btn.svg",
    }
    gf_image_colors.init_pallete(img,
        assets_paths_map,
        (p_color_dominant_hex_str,
        p_colors_hexes_lst)=>{

            // set the background color of the post to its dominant color
            $(p_post_element).css("background-color", `#${p_color_dominant_hex_str}`);

        });

    //----------------------
    */

    const post_id_str = $(p_post_element).data("sys_id");
    const post_title_str = $(p_post_element).find('.post_title').text();
    

    await init_post(post_id_str,
        p_post_element,
        p_http_api_map,
        p_log_fun);
}

//--------------------------------------------------------
function init_post(p_post_id_str :string,
    p_post :HTMLDivElement,
    p_http_api_map,
    p_log_fun) {

    return new Promise(async function(p_resolve_fun, p_reject_fun) {

        const post_image_url_str = $(p_post).find('.thumb_small_url').attr('src');

        init_post_date(p_post, p_log_fun);

        //---------------------
        // INIT_IMAGES_NUMBER
        init_posts_img_num(p_post);

        // INIT_IMAGES
        await init_post_image(post_image_url_str, p_post);

        /*
        const post_images_number = $(p_post).find('.post_images_number');

        $(p_post).on('mouseover',(p_e)=>{
            $(post_images_number).css('visibility', 'visible');
        });
        $(p_post).on('mouseout',(p_e)=>{
            $(post_images_number).css('visibility', 'hidden');
        });
        $(post_images_number).css('right', -$(post_images_number).width()+'px');
        */

        //---------------------
        // TAGGING
        const post_sys_id_str :string = $(p_post).data('sys_id');
        const post_title_str  :string = $(p_post).find('.post_title').text();

        gf_utils.init_tagging(p_post_id_str,
            p_post,
            p_http_api_map,
            p_log_fun)

        /*
        gf_tagger_input_ui.init_tag_input(post_title_str, //p_obj_id_str
            'post',                                       //p_obj_type_str
            p_post,
            //--------------------------------------------------------
            // p_onTagsCreated_fun
            (p_added_tags_lst :string[])=>{
                const tags_container_element = $(p_post).find('.tags_container');
            
                //-----------------------------------------------------
                //FIX!! - when adding the <a> tag of the newly added tag, to the tags_container_element,
                //        detect first if that tag already exists in the list of displayed tags
                //        (on the server this elimination of duplication is already achieved, via set 
                //        data structures, but on the client there is no duplication detection)

                for (var tag_str of p_added_tags_lst) {
                    const tag_url_str        :string            = `/v1/tags/view_objects?tag=${tag_str}&otype=post`;
                    const new_tag_ui_element :HTMLAnchorElement = <HTMLAnchorElement> $(`<a class="gf_post_tag">#${tag_str}</a>`)[0];

                    $(new_tag_ui_element).attr('href',tag_url_str);
                    $(tags_container_element).append(new_tag_ui_element);
                }
                
                //---------------------
                // JS - MASONRY
                
                // IMPORTANT!! - masonry() is a layout call. without calling this every time a new
                //               item is added to the layout, all the items will initially overlap 
                //               (one over the other)
                $('#gf_posts_container').masonry();

                //---------------------
            },
            //--------------------------------------------------------
            ()=>{},
            ()=>{},
            p_http_api_map,
            p_log_fun);
        
        // NOTES
        gf_tagger_notes_ui.init(post_title_str, //p_obj_id_str
            'post', //p_obj_type_str
            p_post,
            p_log_fun);

        */
        
        //---------------------
    });
}

//--------------------------------------------------------
function init_post_image(p_thumbnail_image_src :string,
    p_post :HTMLDivElement) {
    
    const error_img_url_str       = 'https://gloflow.com/images/d/gf_landing_page_logo.png';
    const video_thumb_img_url_str = 'https://gloflow.com/images/d/gf_video_thumb.png';

    return new Promise(function(p_resolve_fun, p_reject_fun) {

        const image :HTMLImageElement = <HTMLImageElement> $(p_post).find('img')[0];

        // ADD!! - for some reason this post does not have a thumbnail image, 
        //         so use some generic post image
        if (p_thumbnail_image_src == null ||
            p_thumbnail_image_src == 'error') {
            $(image).attr('src', error_img_url_str);
        }

        $(image).on('load',(p_e)=>{

            //---------------------
            // IMPORTANT!! - masonry() is a layout call. without calling this every time a new
            //               item is added to the layout, all the items will initially overlap 
            //               (one over the other)

            $('#gf_posts_container').masonry();
            //---------------------
            p_resolve_fun(image);
        });

        $(image).on('error', (p_e)=>{
            $(image).attr('src', error_img_url_str);
            p_reject_fun(`image with url failed to load - ${p_thumbnail_image_src}`);
        });
    });
}

//--------------------------------------------------------
function init_post_date(p_post :HTMLDivElement, p_log_fun) {

    const creation_time_element :HTMLDivElement = <HTMLDivElement> $(p_post).find('.creation_time')[0];
    const creation_time_utc_str :string         = $(creation_time_element).text();
    const creation_date         :Date           = new Date(creation_time_utc_str);

    const date_msg_str = $.timeago(creation_date);
    $(creation_time_element).text(date_msg_str);

    const creation_date__readable_str = creation_date.toDateString();
    const creation_date__readble      = $('<div class="full_creation_date">'+creation_date__readable_str+'</div>');

    $(creation_time_element).mouseover((p_e)=>{
        $(creation_time_element).append(creation_date__readble);
    });

    $(creation_time_element).mouseout((p_e)=>{
        $(creation_date__readble).remove();
    });
}

//---------------------------------------------------
function init_posts_img_num(p_post_element) {

    
    const post_images_number = $(p_post_element).find(".post_images_number")[0];
    const label_element      = $(post_images_number).find(".label");

    // HACK!! - "-1" was visually inferred
    $(post_images_number).css("right", `-${$(post_images_number).outerWidth()-1}px`);
    $(label_element).css("left", `${$(post_images_number).outerWidth()}px`);

    $(p_post_element).mouseover((p_e)=>{
        $(post_images_number).css("visibility", "visible");
    });
    $(p_post_element).mouseout((p_e)=>{
        $(post_images_number).css("visibility", "hidden");
    });
}