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

//-----------------------------------------------------
// NOTES
//-----------------------------------------------------
export function get_notes(p_object_id_str :string,
    p_object_type_str :string,
    p_on_complete_fun,
    p_on_error_fun,
    p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_tagger_client.get_notes()');

    const data_map = {
        'otype': p_object_type_str,
        'o_id':  p_object_id_str
    };

    const url_str = '/v1/tags/notes/get';
    $.ajax({
        'url':         url_str,
        'type':        'GET',
        'data':        data_map,
        'contentType': 'application/json',
        'success':     (p_response_str)=>{
            const data_map  :Object   = JSON.parse(p_response_str);
            const notes_lst :Object[] = data_map['notes_lst'];

            if (notes_lst == null) {
                p_on_complete_fun('success', []);
            } else {
                p_on_complete_fun('success', notes_lst);
            }
        },
        'error':(jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}

//-----------------------------------------------------
export function add_note_to_obj(p_body_str :string,
    p_object_id_str   :string,
    p_object_type_str :string,
    p_on_complete_fun,
    p_on_error_fun,
    p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_tagger_client.add_note_to_obj()');

    const data_map = {
        'otype': p_object_type_str,
        'o_id':  p_object_id_str,
        'body':  p_body_str,
    };

    const url_str = '/v1/tags/notes/create';
    $.ajax({
        'url':         url_str,
        'type':        'POST',
        'data':        JSON.stringify(data_map),
        'contentType': 'application/json',
        'success':     (p_response_str)=>{

            const data_map :Object = JSON.parse(p_response_str);
            p_on_complete_fun('success', data_map);
        },
        'error':(jqXHR,p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}

//-----------------------------------------------------
// TAGS
//-----------------------------------------------------
export async function add_tags_to_obj(p_tags_lst :string[],  
    p_object_id_str   :string,
    p_object_type_str :string,
    p_meta_map,
    p_log_fun) {
    p_log_fun('INFO', 'p_tags_lst:$p_tags_lst');

    const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

        const tags_str :string = p_tags_lst.join(' ');
        const data_map         = {
            "otype": p_object_type_str,
            "o_id":  p_object_id_str,
            "tags":  tags_str,
            "meta_map": p_meta_map,
        };

        const url_str = '/v1/tags/create';
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',
            'success':     (p_response_map)=>{

                const status_str = p_response_map["status"];
                const data_map   = p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error':(jqXHR, p_text_status_str)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
    return p;
}

//-----------------------------------------------------
export function get_objs_with_tag(p_tag_str :string, 
    p_object_type_str :string,
    p_on_complete_fun, 
    p_on_error_fun,
    p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_tagger_client.get_objs_with_tag()');
  
    //this REST api supports supplying multiple tags to the backend, and it will return all of them
    //but Im doing loading from server per tag click, to make initial 
    //load times fast due to minimum network transfers
    const url_str = '/v1/tags/objects?tags='+p_tag_str+'&otype='+p_object_type_str;

    $.ajax({
        'url':         url_str,
        'type':        'GET',
        'contentType': 'application/json',
        'success':     (p_response_str)=>{
            const data_map              :Object   = JSON.parse(p_response_str);
            const objects_with_tags_map :Object[] = data_map['objects_with_tags_dict'];

            p_on_complete_fun('success', objects_with_tags_map);
        },
        'error': (jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}