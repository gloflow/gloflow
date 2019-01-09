///<reference path="../d/jquery.d.ts" />
var gf_post;
(function (gf_post) {
    $(document).ready(function () {
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
        gf_post.init(log_fun);
    });
    //-----------------------------------------------------
    function init(p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_post.init()');
        gf_sys_panel.init(p_log_fun);
        var post_title_str = $('#post_title').text();
        var post_tags_container_element = $('#post_tags_container');
        //------------------------------
        //INIT IMAGE TAGGING
        $('.post_element_image').each(function (p_i, p_post_element) {
            var image_element = $(p_post_element).find('img');
            var img_url_str = $(image_element).attr('src');
            var path_lst = img_url_str.split('/'); //Uri.parse(img_url_str).pathSegments;
            var img_file_str = path_lst[path_lst.length - 1];
            var tags_num_int = $(p_post_element).find('.tags_container .gf_post_element_tag').length;
            //img_file_str example - "6c4a667457f05939af6a5f68690d0f55_thumb_medium.jpeg"
            var img_id_str = img_file_str.split('_')[0];
            p_log_fun('INFO', 'img_id_str - ' + img_id_str);
            var tag_ui_added_bool = false;
            gf_tagger_input_ui.init_tag_input(img_id_str, 'image', p_post_element, 
            //p_onTagsCreated_fun
            //p_onTagsCreated_fun
            function (p_added_tags_lst) {
                view_added_tags(p_post_element, p_added_tags_lst, p_log_fun);
            }, 
            //p_onTagUIAdd_fun
            //p_onTagUIAdd_fun
            function () {
                tag_ui_added_bool = true;
            }, 
            //p_onTagUIRemove_fun
            //p_onTagUIRemove_fun
            function () {
                tag_ui_added_bool = false;
            }, p_log_fun);
            gf_post_image_view.init(p_post_element, p_log_fun);
            $(p_post_element).on('mouseenter', function (p_event) {
                //IMPORTANT!! - only show tags_container if there are tags attached to this post_element
                if (tags_num_int > 0) {
                    $(p_post_element).find('.tags_container').css('visibility', 'visible');
                }
            });
            $(p_post_element).on('mouseleave', function (p_event) {
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
        $('.post_element_video').each(function (p_i, p_post_element) {
            //ADD!! - extract video ID properly
            gf_tagger_input_ui.init_tag_input('fix', 'video', p_post_element, 
            //p_onTagsCreated_fun
            //p_onTagsCreated_fun
            function (p_added_tags_lst) {
                view_added_tags(p_post_element, p_added_tags_lst, p_log_fun);
            }, function () { }, function () { }, p_log_fun);
        });
        //------------------------------
        //final List<String> tags_lst = queryAll('.post_tag').map((p_element) => p_element.text);
        //gf_post_tag_mini_view.init_tags_mini_view(tags_lst,
        //                                          p_log_fun);
    }
    gf_post.init = init;
    //-----------------------------------------------------
    function view_added_tags(p_post_element, p_added_tags_lst, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_post.view_added_tags()');
        var tags_container_element = $(p_post_element).find('.tags_container');
        for (var _i = 0; _i < p_added_tags_lst.length; _i++) {
            var tag_str = p_added_tags_lst[_i];
            var tag_url_str = '/tags/view_objects?tag=' + tag_str + '&otype=image';
            var new_tag_ui_element = $('<a class="gf_post_element_tag">' + tag_str + '</a>');
            $(new_tag_ui_element).attr('href', tag_url_str);
            //IMPORTANT!! - add the new tag link to the DOM
            $(tags_container_element).append(new_tag_ui_element);
        }
    }
    //-----------------------------------------------------
    function get_post_element_tags_num(p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_post.get_post_element_tags_num()');
        //final DivElement tags_container_element = p_post_element.query('.tags_container');
    }
})(gf_post || (gf_post = {}));
///<reference path="../d/jquery.d.ts" />
var gf_post_image_view;
(function (gf_post_image_view) {
    //------------------------------------------------
    function init(p_image_post_element, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_post_image_view.init()');
        $(p_image_post_element).find('img').on('click', function (p_event) {
            var img_medium_url_str = $(p_event.target).attr('src');
            var img_large_url_str = img_medium_url_str.replace('medium', 'large');
            view_image(img_large_url_str, p_log_fun);
        });
    }
    gf_post_image_view.init = init;
    //------------------------------------------------
    function view_image(p_img_url_str, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_post_image_view.view_image()');
        var image_view_element = $("\n\t\t<div id='image_view'>\n\t\t\t<div id='background'></div>\n\t\t\t<img></img>\n\t\t\t<div id=\"close_button\">&#10006;</div>\n\t\t</div>");
        //--------------------------------------------------------
        function load_image() {
            p_log_fun('FUN_ENTER', 'gf_post_image_view.view_image().load_image()');
            var image = document.createElement('img');
            image.src = p_img_url_str;
            $(image).on('load', function (p_e) {
                console.log('img-------');
                var image_x_int = ($(window).innerWidth() - $(image).width()) / 2;
                var image_y_int = ($(window).innerHeight() - $(image).height()) / 2;
                $(image).css('left', image_x_int + 'px');
                $(image).css('top', image_y_int + 'px');
                $(image_view_element).append(image);
                var close_btn = $(image_view_element).find('#close_button');
                $(close_btn).css('left', (image_x_int + $(image).width()) + 'px');
                $(close_btn).css('top', image_y_int + 'px');
            });
        }
        //--------------------------------------------------------
        //offset the top of the image_viewer in case the user scrolled
        $(image_view_element).css('top', document.body.scrollTop + 'px');
        $('body').append(image_view_element);
        //prevent scrolling while in image_view
        $('body').css('overflow', 'hidden');
        //'ESCAPE' key
        $(document).on('keyup', function (p_event) {
            if (p_event.which == 27) {
                $(image_view_element).remove();
                $('body').css('overflow', 'auto');
            }
        });
        $(image_view_element).find('#close_button').on('click', function (p_event) {
            $(image_view_element).remove();
            $('body').css('overflow', 'auto');
        });
    }
})(gf_post_image_view || (gf_post_image_view = {}));
var gf_tagger_client;
(function (gf_tagger_client) {
    //-----------------------------------------------------
    //SNIPPETS
    //-----------------------------------------------------
    function get_notes(p_object_id_str, p_object_type_str, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_client.get_notes()');
        var data_map = {
            'otype': p_object_type_str,
            'o_id': p_object_id_str
        };
        var url_str = '/tags/get_notes';
        $.ajax({
            'url': url_str,
            'type': 'GET',
            'data': data_map,
            'contentType': 'application/json',
            'success': function (p_response_str) {
                var data_map = JSON.parse(p_response_str);
                var notes_lst = data_map['notes_lst'];
                if (notes_lst == null) {
                    p_onComplete_fun('success', []);
                }
                else {
                    p_onComplete_fun('success', notes_lst);
                }
                //p_onComplete_fun('error',
                //            data_str);
            },
            'error': function (jqXHR, p_text_status_str) {
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
        var data_map = {
            'otype': p_object_type_str,
            'o_id': p_object_id_str,
            'body': p_body_str
        };
        var url_str = '/tags/add_note';
        $.ajax({
            'url': url_str,
            'type': 'POST',
            'data': JSON.stringify(data_map),
            'contentType': 'application/json',
            'success': function (p_response_str) {
                var data_map = JSON.parse(p_response_str);
                p_onComplete_fun('success', data_map);
                //p_onComplete_fun('error',
                //            data_str);
            },
            'error': function (jqXHR, p_text_status_str) {
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
        var tags_str = p_tags_lst.join(' ');
        var data_map = {
            'otype': p_object_type_str,
            'o_id': p_object_id_str,
            'tags': tags_str
        };
        var url_str = '/tags/add_tags';
        $.ajax({
            'url': url_str,
            'type': 'POST',
            'data': JSON.stringify(data_map),
            'contentType': 'application/json',
            'success': function (p_response_str) {
                var data_map = JSON.parse(p_response_str);
                p_onComplete_fun('success', data_map);
            },
            'error': function (jqXHR, p_text_status_str) {
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
        var url_str = '/tags/get_objects_with_tags?tags=' + p_tag_str + '&otype=' + p_object_type_str;
        $.ajax({
            'url': url_str,
            'type': 'GET',
            //'data'       :data_args_map,
            'contentType': 'application/json',
            'success': function (p_response_str) {
                var data_map = JSON.parse(p_response_str);
                var objects_with_tags_map = data_map['objects_with_tags_dict'];
                p_onComplete_fun('success', objects_with_tags_map);
            },
            'error': function (jqXHR, p_text_status_str) {
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
        var tagging_input_ui_element = init_tagging_input_ui_element(p_obj_id_str, p_obj_type_str, p_onTagsCreated_fun, p_onTagUIRemove_fun, p_log_fun);
        var tagging_ui_element = $("\n\t\t<div class=\"post_element_controls\">\n\t\t\t<div class=\"add_tags_button\">add tags</div>\n\t\t</div>");
        //OPEN TAG INPUT UI
        $(tagging_ui_element).find('.add_tags_button').on('click', function (p_event) {
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
        $(p_obj_element).on('mouseenter', function (p_event) {
            $(p_obj_element).append(tagging_ui_element);
        });
        $(p_obj_element).on('mouseleave', function (p_event) {
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
        var tagging_input_ui_element = $("\n\t\t<div id=\"tagging_input_container\">\n\t\t\t<div id=\"background\"></div>\n\t\t\t<input type=\"text\" id=\"tags_input\" placeholder=\"(space) separated tags\">\n\t\t\t<div id=\"submit_tags_button\">add</div>\n\t\t\t<div id=\"close_tagging_input_container_button\">&#10006;</div>\n\t\t</div>");
        var tags_input_element = $(tagging_input_ui_element).find('#tags_input');
        //'ESCAPE' key
        $(document).on('keyup', function (p_event) {
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
        $(tags_input_element).on('keyup', function (p_event) {
            //'ENTER' key
            if (p_event.which == 13) {
                p_event.preventDefault();
                add_tags_to_obj(p_obj_id_str, p_obj_type_str, tagging_input_ui_element, 
                //p_onComplete_fun
                //p_onComplete_fun
                function (p_tags_lst) {
                    $(tags_input_element).val('');
                    p_onTagsCreated_fun(p_tags_lst);
                }, 
                //p_onError_fun
                //p_onError_fun
                function () {
                }, p_log_fun);
            }
        });
        $(tagging_input_ui_element).find('#submit_tags_button').on('onmouseup', function (p_event) {
            add_tags_to_obj(p_obj_id_str, p_obj_type_str, tagging_input_ui_element, 
            //p_onComplete_fun
            //p_onComplete_fun
            function (p_tags_lst) {
                $(tags_input_element).val('');
                p_onTagsCreated_fun(p_tags_lst);
            }, 
            //p_onError_fun
            //p_onError_fun
            function () {
            }, p_log_fun);
        });
        //TAG INPUT CLOSE BUTTON
        $(tagging_input_ui_element).find('#close_tagging_input_container_button').on('click', function (p_event) {
            var tagging_input_container_element = $(p_event.target).parent();
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
        var relative_element__width_int = $(p_relative_to_element).width();
        var input_ui_element__width_int = $(p_tagging_input_ui_element).width();
        //p_tagging_input_ui_element.query('input').focus();
        //------------------------
        //Y_COORDINATE
        //document.body.scrollTop - is added to get the 'y' coord relative to the whole doc, regardless of amount of scrolling done
        //const relative_to_element_y_int :number = $(p_relative_to_element).offset().top + $('body').scrollTop(); //p_relative_to_element.getClientRects()[0].top.toInt() +	
        var relative_to_element_y_int = $(p_relative_to_element).offset().top;
        //------------------------
        //X_COORDINATE
        var relative_to_element_x_int = $(p_relative_to_element).offset().left;
        var input_ui_horizontal_overflow_int = (input_ui_element__width_int - relative_element__width_int) / 2;
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
        var tagging_input_y = relative_to_element_y_int - $(p_tagging_input_ui_element).height() / 2;
        $(p_tagging_input_ui_element).css('position', 'absolute');
        $(p_tagging_input_ui_element).css('left', tagging_input_x + 'px');
        $(p_tagging_input_ui_element).css('top', tagging_input_y + 'px');
    }
    //-----------------------------------------------------
    //TAGS SENDING TO SERVER
    //-----------------------------------------------------
    function add_tags_to_obj(p_obj_id_str, p_obj_type_str, p_tagging_ui_element, p_onComplete_fun, p_onError_fun, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_tagger_input_ui.add_tags_to_obj()');
        var tags_str = $(p_tagging_ui_element).find('#tags_input').val();
        var tags_lst = tags_str.split(' ');
        p_log_fun('INFO', 'tags_lst - ' + tags_lst.toString());
        var existing_tags_lst = [];
        $(p_tagging_ui_element).parent().find('.tags_container').find('a').each(function (p_i, p_tag) {
            var tag_str = $(p_tag).text().trim();
            existing_tags_lst.push(tag_str);
        });
        //filter out only tags that are currently not existing/attached to this object
        var new_tags_lst = [];
        for (var _i = 0; _i < tags_lst.length; _i++) {
            var tag_str = tags_lst[_i];
            if (tag_str in existing_tags_lst) {
                new_tags_lst.push(tag_str);
            }
        }
        console.log('>>>>>>>>>>>>>>>>');
        console.log(existing_tags_lst);
        console.log(new_tags_lst);
        //ADD!! - some visual success/failure indicator
        gf_tagger_client.add_tags_to_obj(new_tags_lst, p_obj_id_str, p_obj_type_str, function (p_data_map) {
            var added_tags_lst = p_data_map['added_tags_lst'];
            p_log_fun('INFO', 'added_tags_lst:' + added_tags_lst);
            p_onComplete_fun(added_tags_lst);
        }, function () { }, p_log_fun);
    }
})(gf_tagger_input_ui || (gf_tagger_input_ui = {}));
var gf_sys_panel;
(function (gf_sys_panel) {
    //-----------------------------------------------------
    function init(p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_sys_panel.init()');
        var sys_panel_element = $("<div id=\"sys_panel\">\n\t\t\t<div id=\"view_handle\"></div>\n\n\t\t\t<div id=\"home_btn\">\n\t\t\t\t'<img src=\"/images/d/gf_header_logo.png\"></img>\n\t\t\t</div>\n\n\t\t\t<div id=\"images_app_btn\"><a href=\"/images/flows/browser\">Images</a></div>\n\t\t\t<div id=\"publisher_app_btn\"><a href=\"/posts/browser\">Posts</a></div>\n\t\t\t\n\t\t\t<div id=\"get_invited_btn\">get invited</div>\n\t\t\t<div id=\"login_btn\">login</div>\n\t\t</div>");
        $('body').append(sys_panel_element);
        $(sys_panel_element).find('#view_handle').on('mouseover', function (p_e) {
            $(sys_panel_element).animate({
                top: 0 //move it
            }, 200, function () {
                $(sys_panel_element).find('#view_handle').css('visibility', 'hidden');
            });
        });
    }
    gf_sys_panel.init = init;
})(gf_sys_panel || (gf_sys_panel = {}));
