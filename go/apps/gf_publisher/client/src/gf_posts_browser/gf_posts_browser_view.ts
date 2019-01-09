///<reference path="../d/jquery.d.ts" />
///<reference path="../d/masonry.layout.d.ts" />
///<reference path="../d/jquery.timeago.d.ts" />

namespace gf_posts_browser_view {
//-----------------------------------------------------
export function init(p_initial_posts_infos_lst :Object[],
                p_log_fun) {
    p_log_fun('FUN_ENTER','gf_posts_browser_view.init()');
  
    const image_view_container_element = $('<div id="image_view_posts_container"></div>');

    //----------------
    //JS - MASONRY

    $('#gf_posts_container').masonry(
            {
                'columnWidth' :10,
                'itemSelector':'.item'
            });
    //----------------
    init_posts_images(p_initial_posts_infos_lst,
                ()=>{

                    //---------------------
                    //IMPORTANT!! - masonry() is a layout call. without calling this every time a new
                    //              item is added to the layout, all the items will initially overlap 
                    //              (one over the other)

                    $('#gf_posts_container').masonry(); //'reloadItems');
                    //---------------------
                },
                p_log_fun);

    //-----------------------------------------------------
    function init_page_loading() {
        //p_log_fun('FUN_ENTER','gf_posts_browser_view.init().init_page_loading()');

        var loading_page_bool = false;
        var current_page_int  = 6; //the few initial pages are already statically embedded in the document
        $(window).on('scroll',(e)=>{

            //print('SCROLL >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>');
            //print(document.documentElement.clientHeight);
            //print(window.scrollY);
            //print(window.innerHeight);

            //only test for possible need to page_load if we're not in the middle of a page_load
            //if (!loading_page_bool) {

                const document_height_int :number = document.documentElement.clientHeight;

                //DETECT BOTTOM OF WINDOW REACHED
                //window.scrollY - Y coordinate of the upper screen line relative to 
                //                 the body of the document. This Y coord is never larger then
                //                 body_height - window_height
                //if ($(window).scrollTop() >= ($(document).height() - $(window).height() - 170)) {
                if (window.scrollY >= (document_height_int - (window.innerHeight+20))) {

                    //p_log_fun('INFO','START PAGE LOADING --- >>>');
                    
                    //loading_page_bool = true;

                    //IMPORTANT!! - iterate first, so that in case more pages are to be loaded, before 
                    //              the following load_new_page() is done, we dont call the same 
                    //              page index (current_page_int)
                    current_page_int += 1;

                    load_new_page(current_page_int,
                            5, //p_page_elements_num_int
                            ()=>{

                            },
                            p_log_fun);
                }
            //}
        });
    }
    //-----------------------------------------------------
    init_page_loading();

    return image_view_container_element;
}
//--------------------------------------------------------
function load_new_page(p_page_index_int :number,
                p_page_elements_num_int :number,
                p_onComplete_fun,
                p_log_fun) {
    //p_log_fun('FUN_ENTER','gf_posts_browser_view.load_new_page()');

    gf_posts_browser_client.get_page(p_page_index_int,
                                p_page_elements_num_int,
                                (p_page_lst :Object[])=>{
                                    const posts_infos_lst :Object[] = create_posts_from_page(p_page_lst,
                                                                                    p_page_index_int,
                                                                                    p_log_fun);
                                    init_posts_images(posts_infos_lst,
                                            ()=>{

                                                //---------------------
                                                //IMPORTANT!! - masonry() is a layout call. without calling this every time a new
                                                //              item is added to the layout, all the items will initially overlap 
                                                //              (one over the other)

                                                $('#gf_posts_container').masonry('reloadItems');
                                                //---------------------

                                                p_onComplete_fun();
                                            },
                                            p_log_fun); //load_new_page() only runs with server_comm

                                    

                                    
                                },
                                ()=>{},
                                p_log_fun);
    return;   
}
//--------------------------------------------------------
function create_posts_from_page(p_page_lst :Object[],
                            p_page_index_int :number,
                            p_log_fun) {
    //p_log_fun('FUN_ENTER','gf_posts_browser_view.create_posts_from_page()');

    //--------------------------------------------------------
    function create_post(p_post_map :Object) {
        //p_log_fun('FUN_ENTER','gf_posts_browser_view.create_posts_from_page().create_post()');

        const title_str               :string   = p_post_map['title_str'];
        const image_thumbnail_url_str :string   = p_post_map['thumbnail_url_str'];
        const images_number_str       :string   = p_post_map['images_number_str'];
        const creation_date_str       :string   = p_post_map['creation_datetime_str'];
        const tags_lst                :string[] = p_post_map['tags_lst'];

        //IMPORTANT!! - "item" class is used by Masonry
        const post :HTMLDivElement = <HTMLDivElement> $(`
            <div class='gf_post item gf_post_image_view'>
                <div class='post_title'>`+title_str+`</div>

                <div class='post_images_number'>
                    <div class='num'>`+images_number_str+`</div>
                    <div class='label'>images #</div>
                </div>

                <div>
                    <a class="post_image" target="_blank" href="/posts/`+title_str+`">
                        <img class='thumb_small_url' src="`+image_thumbnail_url_str+`"></img>
                    </a>
                </div>

                <div class='gf_post_creation_date'>`+creation_date_str+`</div>
                <div class='tags_container'></div>
            </div>`)[0];

        const tags_container = $(post).find('.tags_container');

        for (var tag_str of tags_lst) {

            const a :HTMLAnchorElement = <HTMLAnchorElement> $('<a class="gf_post_tag" href="/tags/objects?tag='+tag_str+'&otype=post">#'+tag_str+'</a>')[0];
            $(tags_container).append(a);
        }
        return post;
    }
    //--------------------------------------------------------
    const posts_infos_lst :Object[] = [];

    for (var p_post_map of p_page_lst) {

        const post :HTMLDivElement = create_post(p_post_map);
        $('#gf_posts_container').append(post);

        const post_title_str    :string = p_post_map['title_str'];
        const post_url_str      :string = '/posts/'+post_title_str;
        const thumbnail_url_str :string = p_post_map['thumbnail_url_str'];
        const images_number_str :string = p_post_map['images_number_str'];

        const post_info_map = {
            'post'             :post,
            'post_url_str'     :post_url_str,
            'thumbnail_url_str':thumbnail_url_str,
            //'images_number_str':images_number_str
        };
        posts_infos_lst.push(post_info_map);
    }

    return posts_infos_lst;
}
//--------------------------------------------------------
function init_posts_images(p_posts_infos_lst :Object[],
                        p_onComplete_fun,
                        p_log_fun) {
    p_log_fun('FUN_ENTER','gf_posts_browser_view.init_posts_images()');
        
    const error_img_url_str       = 'http://gloflow.com/images/d/gf_landing_page_logo.png';
    const video_thumb_img_url_str = 'http://gloflow.com/images/d/gf_video_thumb.png';

    //--------------------------------------------------------
    var processed_images_int = 0;

    $(p_posts_infos_lst).each((p_i,p_post_info_map)=>{

        const post                :HTMLDivElement = p_post_info_map['post'];
        const thumbnail_image_src :string         = p_post_info_map['thumbnail_url_str'];





        //images may be loaded out of initial load-issue order, 
        //so the ordering of displayed images is non-deterministic
        init_post_image(thumbnail_image_src,
                error_img_url_str,
                video_thumb_img_url_str,
                post,
                //--------------------------------------------------------
                (p_image_element :HTMLImageElement)=>{
                    console.log('image loaded');
                    const post_url_str :string = p_post_info_map['post_url_str'];
                    
                    init_post(post,
                            post_url_str,
                            p_log_fun);

                    //post has finished loading, so make it visible
                    $(post).css('visibility','visible');

                    processed_images_int += 1;

                    if (processed_images_int == p_posts_infos_lst.length) {
                        p_onComplete_fun();
                    }
                },
                //--------------------------------------------------------
                (p_error_str)=>{
                    p_log_fun("ERROR",p_error_str);
                },
                //--------------------------------------------------------
                p_log_fun);
    });
}
//--------------------------------------------------------
function init_post_image(p_thumbnail_image_src :string,
            p_error_img_url_str       :string,
            p_video_thumb_img_url_str :string,
            p_post                    :HTMLDivElement,
            p_onComplete_fun,
            p_onError_fun,
            p_log_fun) {
    //p_log_fun('FUN_ENTER','gf_posts_browser_view.init_post_image()');

    const image :HTMLImageElement = <HTMLImageElement> $(p_post).find('img')[0];

    //ADD!! - for some reason this post does not have a thumbnail image, 
    //        so use some generic post image
    if (p_thumbnail_image_src == null ||
        p_thumbnail_image_src == 'error') {
        $(image).attr('src',p_error_img_url_str);
    }

    $(image).on('load',(p_e)=>{

        //---------------------
        //IMPORTANT!! - masonry() is a layout call. without calling this every time a new
        //              item is added to the layout, all the items will initially overlap 
        //              (one over the other)

        $('#gf_posts_container').masonry();
        //---------------------
        p_onComplete_fun(image);
    });

    $(image).on('error',(p_e)=>{
        $(image).attr('src',p_error_img_url_str);
        p_onError_fun('image with url failed to load - '+p_thumbnail_image_src);
    });
}
//--------------------------------------------------------
function init_post(p_post :HTMLDivElement,
            p_post_url_str :string,
            p_log_fun) {
    p_log_fun('FUN_ENTER','gf_posts_browser_view.init_post()');

    init_post_date(p_post,
                p_log_fun);
    //---------------------
    //IMAGES_NUMBER

    const post_images_number = $(p_post).find('.post_images_number');

    $(p_post).on('mouseover',(p_e)=>{
        $(post_images_number).css('visibility','visible');
    });
    $(p_post).on('mouseout',(p_e)=>{
        $(post_images_number).css('visibility','hidden');
    });
    $(post_images_number).css('right',-$(post_images_number).width()+'px');
    //---------------------
    //TAGGING
    const post_title_str :string = $(p_post).find('.post_title').text();

    gf_tagger_input_ui.init_tag_input(post_title_str, //p_obj_id_str
                                'post',               //p_obj_type_str
                                p_post,
            //--------------------------------------------------------
            //p_onTagsCreated_fun
            (p_added_tags_lst :string[])=>{
                const tags_container_element = $(p_post).find('.tags_container');
            
                //-----------------------------------------------------
                //FIX!! - when adding the <a> tag of the newly added tag, to the tags_container_element,
                //        detect first if that tag already exists in the list of displayed tags
                //        (on the server this elimination of duplication is already achieved, via set 
                //        data structures, but on the client there is no duplication detection)

                for (var tag_str of p_added_tags_lst) {
                    const tag_url_str        :string            = '/tags/view_objects?tag='+tag_str+'&otype=post';
                    const new_tag_ui_element :HTMLAnchorElement = <HTMLAnchorElement> $('<a class="gf_post_tag">#'+tag_str+'</a>')[0];

                    $(new_tag_ui_element).attr('href',tag_url_str);
                    $(tags_container_element).append(new_tag_ui_element);
                }

                //---------------------
                //JS - MASONRY
                
                //IMPORTANT!! - masonry() is a layout call. without calling this every time a new
                //              item is added to the layout, all the items will initially overlap 
                //              (one over the other)
                $('#gf_posts_container').masonry();
                //---------------------
            },
            //--------------------------------------------------------
            ()=>{},
            ()=>{},
            p_log_fun);
    //---------------------
    //SNIPPET
    
    gf_tagger_notes_ui.init(post_title_str, //p_obj_id_str
                        'post',         //p_obj_type_str
                        p_post,
                        p_log_fun);
    //---------------------
}
//--------------------------------------------------------
function init_post_date(p_post :HTMLDivElement,
                    p_log_fun) {

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
//--------------------------------------------------------
}