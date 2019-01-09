

namespace gf_tagger_client {
//-----------------------------------------------------
//SNIPPETS
//-----------------------------------------------------
export function get_notes(p_object_id_str :string,
                    p_object_type_str :string,
                    p_onComplete_fun,
                    p_onError_fun,
                    p_log_fun) {
    p_log_fun('FUN_ENTER','gf_tagger_client.get_notes()');

    const data_map = {
        'otype':p_object_type_str,
        'o_id' :p_object_id_str
    };
    const url_str = '/tags/get_notes';

    $.ajax({
        'url'        :url_str,
        'type'       :'GET',
        'data'       :data_map,
        'contentType':'application/json',
        'success'    :(p_response_str)=>{
            const data_map  :Object   = JSON.parse(p_response_str);
            const notes_lst :Object[] = data_map['notes_lst'];

            if (notes_lst == null) {
                p_onComplete_fun('success',
                            []);
            } else {
                p_onComplete_fun('success',
                            notes_lst);
            }

             //p_onComplete_fun('error',
            //            data_str);
        },
        'error':(jqXHR,p_text_status_str)=>{
            p_onError_fun(p_text_status_str);
        }
    });
}
//-----------------------------------------------------
export function add_note_to_obj(p_body_str :string,
                    p_object_id_str   :string,
                    p_object_type_str :string,
                    p_onComplete_fun,
                    p_onError_fun,
                    p_log_fun) {
    p_log_fun('FUN_ENTER','gf_tagger_client.add_note_to_obj()');

    /*assert(p_object_type_str == 'image' ||
        p_object_type_str == 'video' ||
        p_object_type_str == 'post');*/

    const data_map = {
        'otype':p_object_type_str,
        'o_id' :p_object_id_str,
        'body' :p_body_str,
    };
    const url_str = '/tags/add_note';

    $.ajax({
        'url'        :url_str,
        'type'       :'POST',
        'data'       :JSON.stringify(data_map),
        'contentType':'application/json',
        'success'    :(p_response_str)=>{

            const data_map :Object = JSON.parse(p_response_str);
            p_onComplete_fun('success',
                         data_map);

             //p_onComplete_fun('error',
            //            data_str);
        },
        'error':(jqXHR,p_text_status_str)=>{
            p_onError_fun(p_text_status_str);
        }
    });
}
//-----------------------------------------------------
//TAGS
//-----------------------------------------------------
export function add_tags_to_obj(p_tags_lst :string[],  
                    p_object_id_str   :string,
                    p_object_type_str :string,
                    p_onComplete_fun,
                    p_onError_fun,
                    p_log_fun) {
    p_log_fun('FUN_ENTER','gf_tagger_client.add_tags_to_obj()');

    /*assert(p_object_type_str == 'image' ||
        p_object_type_str == 'video' ||
        p_object_type_str == 'post');*/

    p_log_fun('INFO','p_tags_lst:$p_tags_lst');

    const tags_str :string = p_tags_lst.join(' ');
    const data_map = {
        'otype':p_object_type_str,
        'o_id' :p_object_id_str,
        'tags' :tags_str,
    };
    const url_str = '/tags/add_tags';

    $.ajax({
        'url'        :url_str,
        'type'       :'POST',
        'data'       :JSON.stringify(data_map),
        'contentType':'application/json',
        'success'    :(p_response_str)=>{

            const data_map :Object = JSON.parse(p_response_str);
            p_onComplete_fun('success',
                        data_map);
        },
        'error':(jqXHR,p_text_status_str)=>{
            p_onError_fun(p_text_status_str);
        }
    });
}
//-----------------------------------------------------
export function get_objs_with_tag(p_tag_str :string, 
                        p_object_type_str :string,
                        p_onComplete_fun, 
                        p_onError_fun,
                        p_log_fun) {
    p_log_fun('FUN_ENTER','gf_tagger_client.get_objs_with_tag()');
  
    //this REST api supports supplying multiple tags to the backend, and it will return all of them
    //but Im doing loading from server per tag click, to make initial 
    //load times fast due to minimum network transfers
    const url_str = '/tags/get_objects_with_tags?tags='+p_tag_str+'&otype='+p_object_type_str;

    $.ajax({
        'url'        :url_str,
        'type'       :'GET',
        //'data'       :data_args_map,
        'contentType':'application/json',
        'success'    :(p_response_str)=>{
            const data_map              :Object   = JSON.parse(p_response_str);
            const objects_with_tags_map :Object[] = data_map['objects_with_tags_dict'];

            p_onComplete_fun('success',
                        objects_with_tags_map);
        },
        'error':(jqXHR,p_text_status_str)=>{
            p_onError_fun(p_text_status_str);
        }
    });
}
//-----------------------------------------------------
}