///<reference path="../d/jquery.d.ts" />
var gf_posts_browser;
(function (gf_posts_browser) {
    //-----------------------------------------------------
    $(document).ready(() => {
        //-------------------------------------------------
        function log_fun(p_g, p_m) {
            var msg_str = p_g + ':' + p_m;
            //chrome.extension.getBackgroundPage().console.log(msg_str);
            switch (p_g) {
                case "INFO":
                    console.log("%cINFO" + ":" + "%c" + p_m, "color:green; background-color:#ACCFAC;", "background-color:#ACCFAC;");
                    break;
                case "FUN_ENTER":
                    console.log("%cFUN_ENTER" + ":" + "%c" + p_m, "color:yellow; background-color:lightgray", "background-color:lightgray");
                    break;
            }
        }
        //-------------------------------------------------
        gf_posts_browser.init(log_fun);
    });
    //-----------------------------------------------------
    function init(p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_posts_browser.init()');
        //this app assumes that the first page of the posts is present 
        //in the dom on app startup... subsequent page loads happen from the server
        const static_posts_infos_lst = load_data_from_dom(p_log_fun);
        gf_sys_panel.init(p_log_fun);
        gf_posts_browser_view.init(static_posts_infos_lst, p_log_fun);
    }
    gf_posts_browser.init = init;
    //-----------------------------------------------------
    //DATA LOADING
    //-----------------------------------------------------
    function load_data_from_dom(p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_posts_browser.load_data_from_dom()');
        const page_posts_infos_lst = [];
        $('body #gf_posts_container').find('.gf_post').each((p_i, p_post) => {
            const post_title_str = $(p_post).find('.post_title').text().trim();
            const post_url_str = '/posts/' + post_title_str;
            //---------------------
            //TAGS
            const tags_infos_lst = [];
            $(p_post).find('.gf_post_tag').each((p_i, p_tag_element) => {
                const tag_str = $(p_tag_element).text();
                const tag_url_str = $(p_tag_element).attr('href');
                const tag_info_map = {
                    'tag_str': tag_str,
                    'tag_url_str': tag_url_str
                };
                tags_infos_lst.push(tag_info_map);
            });
            //--------------------
            //THUMBNAIL URL's
            var thumbnail_url_str = $(p_post).find('img').attr('src');
            if (thumbnail_url_str == '' || thumbnail_url_str == 'error')
                thumbnail_url_str = null;
            //--------------------
            const post_info_map = {
                'post': p_post,
                'post_title_str': post_title_str,
                'post_url_str': post_url_str,
                'tags_infos_lst': tags_infos_lst,
                'thumbnail_url_str': thumbnail_url_str
            };
            page_posts_infos_lst.push(post_info_map);
        });
        return page_posts_infos_lst;
    }
})(gf_posts_browser || (gf_posts_browser = {}));
///<reference path="../d/jquery.d.ts" />
///<reference path="../d/masonry.layout.d.ts" />
///<reference path="../d/jquery.timeago.d.ts" />
var gf_posts_browser_view;
(function (gf_posts_browser_view) {
    //-----------------------------------------------------
    function init(p_initial_posts_infos_lst, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_posts_browser_view.init()');
        const image_view_container_element = $('<div id="image_view_posts_container"></div>');
        //----------------
        //JS - MASONRY
        $('#gf_posts_container').masonry({
            'columnWidth': 10,
            'itemSelector': '.item'
        });
        //----------------
        init_posts_images(p_initial_posts_infos_lst, () => {
            //---------------------
            //IMPORTANT!! - masonry() is a layout call. without calling this every time a new
            //              item is added to the layout, all the items will initially overlap 
            //              (one over the other)
            $('#gf_posts_container').masonry(); //'reloadItems');
            //---------------------
        }, p_log_fun);
        //-----------------------------------------------------
        function init_page_loading() {
            //p_log_fun('FUN_ENTER','gf_posts_browser_view.init().init_page_loading()');
            var loading_page_bool = false;
            var current_page_int = 6; //the few initial pages are already statically embedded in the document
            $(window).on('scroll', (e) => {
                //print('SCROLL >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>');
                //print(document.documentElement.clientHeight);
                //print(window.scrollY);
                //print(window.innerHeight);
                //only test for possible need to page_load if we're not in the middle of a page_load
                //if (!loading_page_bool) {
                const document_height_int = document.documentElement.clientHeight;
                //DETECT BOTTOM OF WINDOW REACHED
                //window.scrollY - Y coordinate of the upper screen line relative to 
                //                 the body of the document. This Y coord is never larger then
                //                 body_height - window_height
                //if ($(window).scrollTop() >= ($(document).height() - $(window).height() - 170)) {
                if (window.scrollY >= (document_height_int - (window.innerHeight + 20))) {
                    //p_log_fun('INFO','START PAGE LOADING --- >>>');
                    //loading_page_bool = true;
                    //IMPORTANT!! - iterate first, so that in case more pages are to be loaded, before 
                    //              the following load_new_page() is done, we dont call the same 
                    //              page index (current_page_int)
                    current_page_int += 1;
                    load_new_page(current_page_int, 5, () => {
                    }, p_log_fun);
                }
                //}
            });
        }
        //-----------------------------------------------------
        init_page_loading();
        return image_view_container_element;
    }
    gf_posts_browser_view.init = init;
    //--------------------------------------------------------
    function load_new_page(p_page_index_int, p_page_elements_num_int, p_onComplete_fun, p_log_fun) {
        //p_log_fun('FUN_ENTER','gf_posts_browser_view.load_new_page()');
        gf_posts_browser_client.get_page(p_page_index_int, p_page_elements_num_int, (p_page_lst) => {
            const posts_infos_lst = create_posts_from_page(p_page_lst, p_page_index_int, p_log_fun);
            init_posts_images(posts_infos_lst, () => {
                //---------------------
                //IMPORTANT!! - masonry() is a layout call. without calling this every time a new
                //              item is added to the layout, all the items will initially overlap 
                //              (one over the other)
                $('#gf_posts_container').masonry('reloadItems');
                //---------------------
                p_onComplete_fun();
            }, p_log_fun); //load_new_page() only runs with server_comm
        }, () => { }, p_log_fun);
        return;
    }
    //--------------------------------------------------------
    function create_posts_from_page(p_page_lst, p_page_index_int, p_log_fun) {
        //p_log_fun('FUN_ENTER','gf_posts_browser_view.create_posts_from_page()');
        //--------------------------------------------------------
        function create_post(p_post_map) {
            //p_log_fun('FUN_ENTER','gf_posts_browser_view.create_posts_from_page().create_post()');
            const title_str = p_post_map['title_str'];
            const image_thumbnail_url_str = p_post_map['thumbnail_url_str'];
            const images_number_str = p_post_map['images_number_str'];
            const creation_date_str = p_post_map['creation_datetime_str'];
            const tags_lst = p_post_map['tags_lst'];
            //IMPORTANT!! - "item" class is used by Masonry
            const post = $(`
            <div class='gf_post item gf_post_image_view'>
                <div class='post_title'>` + title_str + `</div>

                <div class='post_images_number'>
                    <div class='num'>` + images_number_str + `</div>
                    <div class='label'>images #</div>
                </div>

                <div>
                    <a class="post_image" target="_blank" href="/posts/` + title_str + `">
                        <img class='thumb_small_url' src="` + image_thumbnail_url_str + `"></img>
                    </a>
                </div>

                <div class='gf_post_creation_date'>` + creation_date_str + `</div>
                <div class='tags_container'></div>
            </div>`)[0];
            const tags_container = $(post).find('.tags_container');
            for (var tag_str of tags_lst) {
                const a = $('<a class="gf_post_tag" href="/tags/objects?tag=' + tag_str + '&otype=post">#' + tag_str + '</a>')[0];
                $(tags_container).append(a);
            }
            return post;
        }
        //--------------------------------------------------------
        const posts_infos_lst = [];
        for (var p_post_map of p_page_lst) {
            const post = create_post(p_post_map);
            $('#gf_posts_container').append(post);
            const post_title_str = p_post_map['title_str'];
            const post_url_str = '/posts/' + post_title_str;
            const thumbnail_url_str = p_post_map['thumbnail_url_str'];
            const images_number_str = p_post_map['images_number_str'];
            const post_info_map = {
                'post': post,
                'post_url_str': post_url_str,
                'thumbnail_url_str': thumbnail_url_str,
            };
            posts_infos_lst.push(post_info_map);
        }
        return posts_infos_lst;
    }
    //--------------------------------------------------------
    function init_posts_images(p_posts_infos_lst, p_onComplete_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_posts_browser_view.init_posts_images()');
        const error_img_url_str = 'http://gloflow.com/images/d/gf_landing_page_logo.png';
        const video_thumb_img_url_str = 'http://gloflow.com/images/d/gf_video_thumb.png';
        //--------------------------------------------------------
        var processed_images_int = 0;
        $(p_posts_infos_lst).each((p_i, p_post_info_map) => {
            const post = p_post_info_map['post'];
            const thumbnail_image_src = p_post_info_map['thumbnail_url_str'];
            //images may be loaded out of initial load-issue order, 
            //so the ordering of displayed images is non-deterministic
            init_post_image(thumbnail_image_src, error_img_url_str, video_thumb_img_url_str, post, 
            //--------------------------------------------------------
            //--------------------------------------------------------
                (p_image_element) => {
                console.log('image loaded');
                const post_url_str = p_post_info_map['post_url_str'];
                init_post(post, post_url_str, p_log_fun);
                //post has finished loading, so make it visible
                $(post).css('visibility', 'visible');
                processed_images_int += 1;
                if (processed_images_int == p_posts_infos_lst.length) {
                    p_onComplete_fun();
                }
            }, 
            //--------------------------------------------------------
            //--------------------------------------------------------
                (p_error_str) => {
                p_log_fun("ERROR", p_error_str);
            }, 
            //--------------------------------------------------------
            p_log_fun);
        });
    }
    //--------------------------------------------------------
    function init_post_image(p_thumbnail_image_src, p_error_img_url_str, p_video_thumb_img_url_str, p_post, p_onComplete_fun, p_onError_fun, p_log_fun) {
        //p_log_fun('FUN_ENTER','gf_posts_browser_view.init_post_image()');
        const image = $(p_post).find('img')[0];
        //ADD!! - for some reason this post does not have a thumbnail image, 
        //        so use some generic post image
        if (p_thumbnail_image_src == null ||
            p_thumbnail_image_src == 'error') {
            $(image).attr('src', p_error_img_url_str);
        }
        $(image).on('load', (p_e) => {
            //---------------------
            //IMPORTANT!! - masonry() is a layout call. without calling this every time a new
            //              item is added to the layout, all the items will initially overlap 
            //              (one over the other)
            $('#gf_posts_container').masonry();
            //---------------------
            p_onComplete_fun(image);
        });
        $(image).on('error', (p_e) => {
            $(image).attr('src', p_error_img_url_str);
            p_onError_fun('image with url failed to load - ' + p_thumbnail_image_src);
        });
    }
    //--------------------------------------------------------
    function init_post(p_post, p_post_url_str, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_posts_browser_view.init_post()');
        init_post_date(p_post, p_log_fun);
        //---------------------
        //IMAGES_NUMBER
        const post_images_number = $(p_post).find('.post_images_number');
        $(p_post).on('mouseover', (p_e) => {
            $(post_images_number).css('visibility', 'visible');
        });
        $(p_post).on('mouseout', (p_e) => {
            $(post_images_number).css('visibility', 'hidden');
        });
        $(post_images_number).css('right', -$(post_images_number).width() + 'px');
        //---------------------
        //TAGGING
        const post_title_str = $(p_post).find('.post_title').text();
        gf_tagger_input_ui.init_tag_input(post_title_str, 'post', p_post, 
        //--------------------------------------------------------
        //p_onTagsCreated_fun
        //--------------------------------------------------------
        //p_onTagsCreated_fun
            (p_added_tags_lst) => {
            const tags_container_element = $(p_post).find('.tags_container');
            //-----------------------------------------------------
            //FIX!! - when adding the <a> tag of the newly added tag, to the tags_container_element,
            //        detect first if that tag already exists in the list of displayed tags
            //        (on the server this elimination of duplication is already achieved, via set 
            //        data structures, but on the client there is no duplication detection)
            for (var tag_str of p_added_tags_lst) {
                const tag_url_str = '/tags/view_objects?tag=' + tag_str + '&otype=post';
                const new_tag_ui_element = $('<a class="gf_post_tag">#' + tag_str + '</a>')[0];
                $(new_tag_ui_element).attr('href', tag_url_str);
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
        //--------------------------------------------------------
            () => { }, () => { }, p_log_fun);
        //---------------------
        //SNIPPET
        gf_tagger_notes_ui.init(post_title_str, 'post', p_post, p_log_fun);
        //---------------------
    }
    //--------------------------------------------------------
    function init_post_date(p_post, p_log_fun) {
        const creation_time_element = $(p_post).find('.creation_time')[0];
        const creation_time_utc_str = $(creation_time_element).text();
        const creation_date = new Date(creation_time_utc_str);
        const date_msg_str = $.timeago(creation_date);
        $(creation_time_element).text(date_msg_str);
        const creation_date__readable_str = creation_date.toDateString();
        const creation_date__readble = $('<div class="full_creation_date">' + creation_date__readable_str + '</div>');
        $(creation_time_element).mouseover((p_e) => {
            $(creation_time_element).append(creation_date__readble);
        });
        $(creation_time_element).mouseout((p_e) => {
            $(creation_date__readble).remove();
        });
    }
})(gf_posts_browser_view || (gf_posts_browser_view = {}));
var gf_posts_browser_client;
(function (gf_posts_browser_client) {
    //-----------------------------------------------------
    function get_page(p_page_index_int, p_page_elements_num_int, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_posts_browser_client.get_page()');
        const url_str = '/posts/browser_page';
        const data_map = {
            'pg_index': p_page_index_int,
            'pg_size': p_page_elements_num_int
        };
        $.ajax({
            'url': url_str,
            'type': 'GET',
            'data': data_map,
            'contentType': 'application/json',
            'success': (p_response_str) => {
                const page_lst = JSON.parse(p_response_str);
                p_onComplete_fun(page_lst);
            },
            'error': (jqXHR, p_text_status_str) => {
                p_onError_fun(p_text_status_str);
            }
        });
    }
    gf_posts_browser_client.get_page = get_page;
})(gf_posts_browser_client || (gf_posts_browser_client = {}));
var gf_tagger_client;
(function (gf_tagger_client) {
    //-----------------------------------------------------
    //SNIPPETS
    //-----------------------------------------------------
    function get_notes(p_object_id_str, p_object_type_str, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_client.get_notes()');
        const data_map = {
            'otype': p_object_type_str,
            'o_id': p_object_id_str
        };
        const url_str = '/tags/get_notes';
        $.ajax({
            'url': url_str,
            'type': 'GET',
            'data': data_map,
            'contentType': 'application/json',
            'success': (p_response_str) => {
                const data_map = JSON.parse(p_response_str);
                const notes_lst = data_map['notes_lst'];
                if (notes_lst == null) {
                    p_onComplete_fun('success', []);
                }
                else {
                    p_onComplete_fun('success', notes_lst);
                }
                //p_onComplete_fun('error',
                //            data_str);
            },
            'error': (jqXHR, p_text_status_str) => {
                p_onError_fun(p_text_status_str);
            }
        });
    }
    gf_tagger_client.get_notes = get_notes;
    //-----------------------------------------------------
    function add_note_to_obj(p_body_str, p_object_id_str, p_object_type_str, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_client.add_note_to_obj()');
        /*assert(p_object_type_str == 'image' ||
            p_object_type_str == 'video' ||
            p_object_type_str == 'post');*/
        const data_map = {
            'otype': p_object_type_str,
            'o_id': p_object_id_str,
            'body': p_body_str,
        };
        const url_str = '/tags/add_note';
        $.ajax({
            'url': url_str,
            'type': 'POST',
            'data': JSON.stringify(data_map),
            'contentType': 'application/json',
            'success': (p_response_str) => {
                const data_map = JSON.parse(p_response_str);
                p_onComplete_fun('success', data_map);
                //p_onComplete_fun('error',
                //            data_str);
            },
            'error': (jqXHR, p_text_status_str) => {
                p_onError_fun(p_text_status_str);
            }
        });
    }
    gf_tagger_client.add_note_to_obj = add_note_to_obj;
    //-----------------------------------------------------
    //TAGS
    //-----------------------------------------------------
    function add_tags_to_obj(p_tags_lst, p_object_id_str, p_object_type_str, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_client.add_tags_to_obj()');
        /*assert(p_object_type_str == 'image' ||
            p_object_type_str == 'video' ||
            p_object_type_str == 'post');*/
        p_log_fun('INFO', 'p_tags_lst:$p_tags_lst');
        const tags_str = p_tags_lst.join(' ');
        const data_map = {
            'otype': p_object_type_str,
            'o_id': p_object_id_str,
            'tags': tags_str,
        };
        const url_str = '/tags/add_tags';
        $.ajax({
            'url': url_str,
            'type': 'POST',
            'data': JSON.stringify(data_map),
            'contentType': 'application/json',
            'success': (p_response_str) => {
                const data_map = JSON.parse(p_response_str);
                p_onComplete_fun('success', data_map);
            },
            'error': (jqXHR, p_text_status_str) => {
                p_onError_fun(p_text_status_str);
            }
        });
    }
    gf_tagger_client.add_tags_to_obj = add_tags_to_obj;
    //-----------------------------------------------------
    function get_objs_with_tag(p_tag_str, p_object_type_str, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_client.get_objs_with_tag()');
        //this REST api supports supplying multiple tags to the backend, and it will return all of them
        //but Im doing loading from server per tag click, to make initial 
        //load times fast due to minimum network transfers
        const url_str = '/tags/get_objects_with_tags?tags=' + p_tag_str + '&otype=' + p_object_type_str;
        $.ajax({
            'url': url_str,
            'type': 'GET',
            //'data'       :data_args_map,
            'contentType': 'application/json',
            'success': (p_response_str) => {
                const data_map = JSON.parse(p_response_str);
                const objects_with_tags_map = data_map['objects_with_tags_dict'];
                p_onComplete_fun('success', objects_with_tags_map);
            },
            'error': (jqXHR, p_text_status_str) => {
                p_onError_fun(p_text_status_str);
            }
        });
    }
    gf_tagger_client.get_objs_with_tag = get_objs_with_tag;
})(gf_tagger_client || (gf_tagger_client = {}));
var gf_tagger_input_ui;
(function (gf_tagger_input_ui) {
    //-----------------------------------------------------
    //in gf_post view
    function init_tag_input(p_obj_id_str, p_obj_type_str, p_obj_element, p_onTagsCreated_fun, p_onTagUIAdd_fun, p_onTagUIRemove_fun, p_log_fun) {
        //p_log_fun('FUN_ENTER','gf_tagger_input_ui.init_tag_input()');
        const tagging_input_ui_element = init_tagging_input_ui_element(p_obj_id_str, p_obj_type_str, p_onTagsCreated_fun, p_onTagUIRemove_fun, p_log_fun);
        const tagging_ui_element = $(`
		<div class="post_element_controls">
			<div class="add_tags_button">add tags</div>
		</div>`);
        //OPEN TAG INPUT UI
        $(tagging_ui_element).find('.add_tags_button').on('click', (p_event) => {
            console.log('zzzzz');
            //remove the tagging_input_container if its already displayed
            //for tagging another post_element
            if ($('#tagging_input_container') != null) {
                $('#tagging_input_container').remove();
            }
            //post_element_element - as in part of a post. post_element_element because its a 
            //                       html element of the post_element
            //final DivElement post_element_element = p_event.target.parent.parent;
            place_tagging_input_ui_element(tagging_input_ui_element, p_obj_element, p_log_fun);
            if (p_onTagUIAdd_fun != null)
                p_onTagUIAdd_fun();
        });
        //------------------------
        ////---------------------
        //js.context
        //    .callMethod(r'$', ['document'])
        //    .callMethod('bind',['DOMNodeRemoved',(e) {
        //    		print('zzzzzzzzzzzzzz');
        //    	}]);
        ////---------------------
        /*//'T' key - open tagging UI to the element that has the cursor
        //          hovering over it
        final subscription = document.onKeyUp.listen((p_event) {
            if (p_event.keyCode == 84) {
    
                //remove the tagging_input_container if its already displayed
                //for tagging another post_element
                if (query('#tagging_input_container') != null) {
                    query('#tagging_input_container').remove();
                }
    
                place_tagging_input_ui_element(tagging_input_ui_element,
                                               p_obj_element, //post_element,
                                               p_log_fun);
    
                //prevent this handler being invoked while the user
                //is typing in tags into the input field
                //subscription.pause();
            }
        });*/
        //------------------------
        //IMPORTANT!! - onMouseEnter/onMouseLeave fire when the target element is entered/left, 
        //              but unline mouseon/mouseout it will not fire if its children are entered/left.
        $(p_obj_element).on('mouseenter', (p_event) => {
            $(p_obj_element).append(tagging_ui_element);
        });
        $(p_obj_element).on('mouseleave', (p_event) => {
            $(tagging_ui_element).remove();
            ////relatedTarget - The relatedTargert property can be used with the mouseover 
            ////                event to indicate the element the cursor just exited, 
            ////                or with the mouseout event to indicate the element the 
            ////                cursor just entered.
            //if (p_event.relatedTarget != null && 
            //	!p_event.relatedTarget.classes.contains('add_tags_button')) {
            //	tagging_ui_element.remove();
            //}
        });
        //------------------------
    }
    gf_tagger_input_ui.init_tag_input = init_tag_input;
    //-----------------------------------------------------
    //TAGS UI UTILS
    //-----------------------------------------------------
    function init_tagging_input_ui_element(p_obj_id_str, p_obj_type_str, p_onTagsCreated_fun, p_onTagUIRemove_fun, p_log_fun) {
        //p_log_fun('FUN_ENTER','gf_tagger_input_ui.init_tagging_input_ui_element()');
        const tagging_input_ui_element = $(`
		<div id="tagging_input_container">
			<div id="background"></div>
			<input type="text" id="tags_input" placeholder="(space) separated tags">
			<div id="submit_tags_button">add</div>
			<div id="close_tagging_input_container_button">&#10006;</div>
		</div>`);
        const tags_input_element = $(tagging_input_ui_element).find('#tags_input');
        //'ESCAPE' key
        $(document).on('keyup', (p_event) => {
            if (p_event.which == 27) {
                //remove any previously present tagging_input_container's
                $(tagging_input_ui_element).remove();
                if (p_onTagUIRemove_fun != null) {
                    p_onTagUIRemove_fun();
                }
            }
        });
        //to handlers for the same thing, one for the user clicking on the button,
        //the other for the user pressing 'enter'  
        $(tags_input_element).on('keyup', (p_event) => {
            //'ENTER' key
            if (p_event.which == 13) {
                p_event.preventDefault();
                add_tags_to_obj(p_obj_id_str, p_obj_type_str, tagging_input_ui_element, 
                //p_onComplete_fun
                //p_onComplete_fun
                    (p_tags_lst) => {
                    $(tags_input_element).val('');
                    p_onTagsCreated_fun(p_tags_lst);
                }, 
                //p_onError_fun
                //p_onError_fun
                    () => {
                }, p_log_fun);
            }
        });
        $(tagging_input_ui_element).find('#submit_tags_button').on('onmouseup', (p_event) => {
            add_tags_to_obj(p_obj_id_str, p_obj_type_str, tagging_input_ui_element, 
            //p_onComplete_fun
            //p_onComplete_fun
                (p_tags_lst) => {
                $(tags_input_element).val('');
                p_onTagsCreated_fun(p_tags_lst);
            }, 
            //p_onError_fun
            //p_onError_fun
                () => {
            }, p_log_fun);
        });
        //TAG INPUT CLOSE BUTTON
        $(tagging_input_ui_element).find('#close_tagging_input_container_button').on('click', (p_event) => {
            const tagging_input_container_element = $(p_event.target).parent();
            $(tagging_input_container_element).remove();
            if (p_onTagUIRemove_fun != null) {
                p_onTagUIRemove_fun();
            }
        });
        return tagging_input_ui_element;
    }
    //-----------------------------------------------------
    function place_tagging_input_ui_element(p_tagging_input_ui_element, p_relative_to_element, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_input_ui.place_tagging_input_ui_element()');
        $('body').append(p_tagging_input_ui_element);
        const relative_element__width_int = $(p_relative_to_element).width();
        const input_ui_element__width_int = $(p_tagging_input_ui_element).width();
        //p_tagging_input_ui_element.query('input').focus();
        //------------------------
        //Y_COORDINATE
        //document.body.scrollTop - is added to get the 'y' coord relative to the whole doc, regardless of amount of scrolling done
        //const relative_to_element_y_int :number = $(p_relative_to_element).offset().top + $('body').scrollTop(); //p_relative_to_element.getClientRects()[0].top.toInt() +	
        const relative_to_element_y_int = $(p_relative_to_element).offset().top;
        //------------------------
        //X_COORDINATE
        const relative_to_element_x_int = $(p_relative_to_element).offset().left;
        const input_ui_horizontal_overflow_int = (input_ui_element__width_int - relative_element__width_int) / 2;
        var tagging_input_x;
        //input_ui is wider then target element
        if (input_ui_horizontal_overflow_int > 0) {
            //input_ui is cutoff on the left side
            if ((relative_to_element_x_int - input_ui_horizontal_overflow_int) < 0) {
                //position input_ui with its left side aligned with left edge of element to be tagged
                tagging_input_x = relative_to_element_x_int;
            }
            else if (((relative_to_element_x_int + relative_element__width_int) + input_ui_horizontal_overflow_int) >
                $(window).innerWidth()) {
                //position inpout_ui with its right edge aligned with the right edge of element to be tagged
                tagging_input_x = (relative_to_element_x_int + relative_element__width_int) -
                    input_ui_element__width_int;
            }
            else {
                //positions that tag input container in the middle, and above, of the post_element
                tagging_input_x = relative_to_element_x_int - (input_ui_element__width_int - relative_element__width_int) / 2;
            }
        }
        else {
            //positions that tag input container in the middle, and above, of the post_element
            tagging_input_x = relative_to_element_x_int - (input_ui_element__width_int - relative_element__width_int) / 2;
        }
        const tagging_input_y = relative_to_element_y_int - $(p_tagging_input_ui_element).height() / 2;
        $(p_tagging_input_ui_element).css('position', 'absolute');
        $(p_tagging_input_ui_element).css('left', tagging_input_x + 'px');
        $(p_tagging_input_ui_element).css('top', tagging_input_y + 'px');
    }
    //-----------------------------------------------------
    //TAGS SENDING TO SERVER
    //-----------------------------------------------------
    function add_tags_to_obj(p_obj_id_str, p_obj_type_str, p_tagging_ui_element, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_input_ui.add_tags_to_obj()');
        const tags_str = $(p_tagging_ui_element).find('#tags_input').val();
        const tags_lst = tags_str.split(' ');
        p_log_fun('INFO', 'tags_lst - ' + tags_lst.toString());
        const existing_tags_lst = [];
        $(p_tagging_ui_element).parent().find('.tags_container').find('a').each((p_i, p_tag) => {
            const tag_str = $(p_tag).text().trim();
            existing_tags_lst.push(tag_str);
        });
        //filter out only tags that are currently not existing/attached to this object
        const new_tags_lst = [];
        for (var tag_str of tags_lst) {
            if (tag_str in existing_tags_lst) {
                new_tags_lst.push(tag_str);
            }
        }
        console.log('>>>>>>>>>>>>>>>>');
        console.log(existing_tags_lst);
        console.log(new_tags_lst);
        //ADD!! - some visual success/failure indicator
        gf_tagger_client.add_tags_to_obj(new_tags_lst, p_obj_id_str, p_obj_type_str, (p_data_map) => {
            const added_tags_lst = p_data_map['added_tags_lst'];
            p_log_fun('INFO', 'added_tags_lst:' + added_tags_lst);
            p_onComplete_fun(added_tags_lst);
        }, () => { }, p_log_fun);
    }
})(gf_tagger_input_ui || (gf_tagger_input_ui = {}));
var gf_tagger_notes_ui;
(function (gf_tagger_notes_ui) {
    //-----------------------------------------------------
    function init(p_obj_id_str, p_obj_type_str, p_obj_element, p_log_fun) {
        //p_log_fun('FUN_ENTER','gf_tagger_notes_ui.init()');
        const notes_panel_btn = $(`
			<div id='notes_panel_btn'>
				<div class='icon'>notes</div>
			</div>`);
        $(p_obj_element).append(notes_panel_btn);
        const notes_panel = $(`
		<div id='notes_panel'>
			<div id='background'></div>

			<div id='container'>
				<div id='add_note_btn'>
					<div class='icon'>+</div>
				</div>
				<div id='notes'>
				</div>
			</div>
		</div>`);
        const background = $(notes_panel).find('#background');
        const add_note_btn = $(notes_panel).find('#add_note_btn');
        var notes_open_bool = false;
        var notes_init_bool = false;
        $(notes_panel_btn).on('click', (p_event) => {
            if (notes_open_bool) {
            }
            else {
                $(p_obj_element).append(notes_panel);
                //------------------------
                //GET
                if (!notes_init_bool) {
                    //------------
                    //NOTE_INPUT_PANEL
                    const note_input_panel = $(`
					<div class='note_input_panel'>
						<textarea name="note_input" cols="30" rows="3"></textarea>
					</div>`);
                    notes_panel.append(note_input_panel);
                    //------------
                    get_notes(p_obj_id_str, p_obj_type_str, notes_panel, () => {
                        notes_init_bool = true;
                    }, p_log_fun);
                }
                //------------------------
                $(notes_panel).css('visibility', "visible");
            }
        });
        //------------------------
        add_note_btn.on('click', (p_event) => {
            console.log('>>>>> ENTER');
            run__remote_add_note(p_obj_id_str, p_obj_type_str, notes_panel, 
            //p_onComplete_fun,
            //p_onComplete_fun,
                () => {
                //print('>>>>>>>');
                //print(notes_panel.query('#notes').offsetTop);
                //print(notes_panel.query('#notes').getComputedStyle().top);
                //----------------------
                //GROW BACKGROUND
                const background_padding_size_int = 30;
                const notes_height_int = $(notes_panel).find('#notes').height();
                const notes_y_int = $(notes_panel).find('#notes').offset().top;
                const new_height_int = notes_y_int + notes_height_int + 2 * background_padding_size_int;
                //print('aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
                //print(notes_height_int);
                //print(notes_y_int);
                //print(new_height_int);
                $(background).css('height', new_height_int + 'px');
                //----------------------
            }, p_log_fun);
        });
        //------------------------
        //IMPORTANT!! - onMouseEnter/onMouseLeave fire when the target element is entered/left, 
        //              but unline mouseon/mouseout it will not fire if its children are entered/left.
        $(p_obj_element).on('mouseenter', (p_event) => {
            $(notes_panel_btn).css('visibility', 'visible');
        });
        $(p_obj_element).on('mouseleave', (p_event) => {
            $(notes_panel_btn).css('visibility', 'hidden');
        });
        //------------------------
        //'ESCAPE' key
        $(document).on('keyup', (p_event) => {
            if (p_event.which == 27) {
                //remove any previously present note_input_container's
                $(notes_panel).remove();
            }
        });
        //------------------------
    }
    gf_tagger_notes_ui.init = init;
    //-----------------------------------------------------
    function get_notes(p_obj_id_str, p_obj_type_str, p_notes_panel, p_onComplete_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_notes_ui.get_notes()');
        //------------------------
        //IMPORTANT!! - get notes via HTTP from backend gf_tagger_service
        gf_tagger_client.get_notes(p_obj_id_str, p_obj_type_str, 
        //p_onComplete_fun
        //p_onComplete_fun
            (p_notes_lst) => {
            for (var note_map of p_notes_lst) {
                const user_id_str = note_map['user_id_str'];
                const body_str = note_map['body_str'];
                add_note_view(body_str, user_id_str, p_notes_panel, p_log_fun);
            }
            p_onComplete_fun();
        }, () => { }, p_log_fun);
        //------------------------	
    }
    //-----------------------------------------------------
    function run__remote_add_note(p_obj_id_str, p_obj_type_str, p_notes_panel, p_onComplete_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_notes_ui.run__remote_add_note()');
        const text_element = $(p_notes_panel).find('.note_input_panel textarea');
        const note_body_str = $(text_element).val();
        p_log_fun('INFO', 'note_body_str        - $note_body_str');
        p_log_fun('INFO', 'note_body_str.length - ${note_body_str.length}');
        if (note_body_str.length > 0) {
            //ADD!! - some visual success/failure indicator
            gf_tagger_client.add_note_to_obj(note_body_str, p_obj_id_str, p_obj_type_str, () => {
                add_note_view(note_body_str, 'anonymouse', p_notes_panel, p_log_fun);
                $(text_element).val(''); //reset text
            }, () => { }, p_log_fun);
        }
    }
    //-----------------------------------------------------
    function add_note_view(p_body_str, p_user_id_str, p_notes_panel, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_notes_ui.add_note_view()');
        if (p_body_str.length > 20)
            console.log(p_body_str.substring(0, 20) + '...');
        var short_body_str;
        if (p_body_str.length > 20)
            short_body_str = p_body_str.substring(0, 20) + '...';
        else
            short_body_str = p_body_str;
        const new_note_element = $(`
		<div class='note'>
			<div class='icon'>n</div>
			<div class='details'>
				<div class='user'>` + p_user_id_str + `</div>
				<div class='body'>` + short_body_str + `</div>
			</div>
		</div>`);
        //other notes already exist
        if ($(p_notes_panel).find('#notes').children().length > 0) {
            const latest_note = $(p_notes_panel).find('#notes').children()[0];
            //insertBefore() - makes the new_note the first element in the list,
            //                 because the newest notes are at the top.
            /*$(p_notes_panel).find('#notes').insertBefore(new_note_element, //new_child
                                                    latest_note);        //ref_child*/
            $(new_note_element).insertBefore(latest_note);
        }
        else {
            $(p_notes_panel).find('#notes').append(new_note_element);
        }
        $(new_note_element).css('opacity', '0.0');
        $(new_note_element).animate({ 'opacity': 1.0 }, 300, () => { });
    }
})(gf_tagger_notes_ui || (gf_tagger_notes_ui = {}));
var gf_sys_panel;
(function (gf_sys_panel) {
    //-----------------------------------------------------
    function init(p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_sys_panel.init()');
        const sys_panel_element = $(`<div id="sys_panel">
			<div id="view_handle"></div>
			<div id="home_btn">
				'<img src="/images/d/gf_header_logo.png"></img>
			</div>
			<div id="images_app_btn"><a href="/images/flows/browser">Images</a></div>
			<div id="publisher_app_btn"><a href="/posts/browser">Posts</a></div>
			<div id="get_invited_btn">get invited</div>
			<div id="login_btn">login</div>
		</div>`);
        $('body').append(sys_panel_element);
        $(sys_panel_element).find('#view_handle').on('mouseover', (p_e) => {
            $(sys_panel_element).animate({
                top: 0 //move it
            }, 200, () => {
                $(sys_panel_element).find('#view_handle').css('visibility', 'hidden');
            });
        });
    }
    gf_sys_panel.init = init;
})(gf_sys_panel || (gf_sys_panel = {}));
