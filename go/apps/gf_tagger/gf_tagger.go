/*
GloFlow media management/publishing system
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

package main

import (
	"fmt"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/apps/gf_publisher_lib"
)
//---------------------------------------------------
//p_tags_str      - :String - "," separated list of strings
//p_object_id_str - :String - this is an external identifier for an object, not necessarily its internal. 
//                            for posts - their p_object_extern_id_str is their Title, but internally they have
//                                        another ID.

func add_tags_to_object(p_tags_str string,
	p_object_type_str      string,
	p_object_extern_id_str string,
	p_runtime_sys          *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_tagger.add_tags_to_object()")

	if p_object_type_str != "post" &&
		p_object_type_str != "image" &&
		p_object_type_str != "event" {
		gf_err := gf_core.Error__create(fmt.Sprintf("p_object_type_str (%s) is not of supported type (post|image|event)",p_object_type_str),
			"verify__invalid_value_error",
			&map[string]interface{}{
				"tags_str":       p_tags_str,
				"object_type_str":p_object_type_str,
			},
			nil, "gf_tagger", p_runtime_sys)
		return gf_err
	}
	
	tags_lst, gf_err := parse_tags(p_tags_str,
		500, //p_max_tags_bulk_size_int        int, //500
		20,  //p_max_tag_characters_number_int int, //20	
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	p_runtime_sys.Log_fun("INFO","tags_lst - "+fmt.Sprint(tags_lst))
	//---------------
	//POST
	
	switch p_object_type_str {
		//---------------
		//POST
		case "post":
			post_title_str      := p_object_extern_id_str
			exists_bool, gf_err := gf_publisher_lib.DB__check_post_exists(post_title_str, p_runtime_sys)
			if gf_err != nil {
				return gf_err
			}
			
			if exists_bool {
				p_runtime_sys.Log_fun("INFO","POST EXISTS")
				gf_err := db__add_tags_to_post(post_title_str, tags_lst, p_runtime_sys)
				return gf_err
			} else {
				gf_err := gf_core.Error__create(fmt.Sprintf("post with title (%s) doesnt exist, while adding a tags - %s", post_title_str, tags_lst),
					"verify__invalid_value_error",
					&map[string]interface{}{
						"post_title_str":post_title_str,
						"tags_lst":      tags_lst,
					},
					nil, "gf_tagger", p_runtime_sys)
				return gf_err
			}

		//---------------
		//IMAGE
		case "image":
			image_id_str        := p_object_extern_id_str
			exists_bool, gf_err := gf_images_lib.DB__image_exists(image_id_str, p_runtime_sys)
			if gf_err != nil {
				return gf_err
			}
			if exists_bool {
				gf_err := db__add_tags_to_image(image_id_str, tags_lst, p_runtime_sys)
				if gf_err != nil {
					return gf_err
				}
			}
		//---------------
	}
	return nil
}
//---------------------------------------------------
func get_objects_with_tags(p_tags_lst []string,
	p_object_type_str string,
	p_page_index_int  int,
	p_page_size_int   int,
	p_runtime_sys     *gf_core.Runtime_sys) (map[string][]map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_tagger.get_objects_with_tags()")
		
	objects_with_tags_map := map[string][]map[string]interface{}{}
	for _,tag_str := range p_tags_lst {
		objects_with_tag_lst, gf_err := get_objects_with_tag(tag_str,
			p_object_type_str,
			p_page_index_int,
			p_page_size_int,
			p_runtime_sys)

		if gf_err != nil {
			return nil, gf_err
		}
		objects_with_tags_map[tag_str] = objects_with_tag_lst
	}
	return objects_with_tags_map, nil
}
//---------------------------------------------------
func get_objects_with_tag(p_tag_str string,
	p_object_type_str string,
	p_page_index_int  int,
	p_page_size_int   int,
	p_runtime_sys     *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_tagger.get_objects_with_tag()")
	p_runtime_sys.Log_fun("INFO",     "p_object_type_str - "+p_object_type_str)

	//ADD!! - add support for tagging "image" p_object_type_str's
	if p_object_type_str != "post" {
		gf_err := gf_core.Error__create(fmt.Sprintf("trying to get objects with a tag (%s) for objects type thats not supported - %s", p_tag_str, p_object_type_str),
			"verify__invalid_value_error",
			&map[string]interface{}{
				"tag_str":        p_tag_str,
				"object_type_str":p_object_type_str,
			},
			nil, "gf_tagger", p_runtime_sys)
		return nil, gf_err
	}
	
	posts_with_tag_lst, gf_err := db__get_posts_with_tag(p_tag_str,
		p_page_index_int,
		p_page_size_int,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//package up info of each post that was found with tag 
	min_posts_infos_lst := []map[string]interface{}{}
	for _,post := range posts_with_tag_lst {
		post_info_map := map[string]interface{}{
			"title_str":              post.Title_str,
			"tags_lst":               post.Tags_lst,
			"url_str":                fmt.Sprintf("/posts/%s",post.Title_str),
			"object_type_str":        p_object_type_str,
			"thumbnail_small_url_str":post.Thumbnail_url_str,
		}
		min_posts_infos_lst = append(min_posts_infos_lst,post_info_map)
	}

	objects_infos_lst := min_posts_infos_lst
	return objects_infos_lst, nil
}
//---------------------------------------------------
func parse_tags(p_tags_str string,
	p_max_tags_bulk_size_int        int, //500
	p_max_tag_characters_number_int int, //20
	p_runtime_sys                   *gf_core.Runtime_sys) ([]string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_tagger.parse_tags()")
	
	tags_lst := strings.Split(p_tags_str," ")
	//---------------------
	if len(tags_lst) > p_max_tags_bulk_size_int {
		gf_err := gf_core.Error__create(fmt.Sprintf("too many tags supplied - max is %s",p_max_tags_bulk_size_int),
			"verify__value_too_many_error",
			&map[string]interface{}{
				"tags_lst":              tags_lst,
				"max_tags_bulk_size_int":p_max_tags_bulk_size_int,
			},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}
	//---------------------
	for _,tag_str := range tags_lst {
		if len(tag_str) > p_max_tag_characters_number_int {
			gf_err := gf_core.Error__create(fmt.Sprintf("tag (%s) is too long - max is (%s)", tag_str, fmt.Sprint(p_max_tag_characters_number_int)),
				"verify__string_too_long_error",
				&map[string]interface{}{
					"tag_str":                      tag_str,
					"max_tag_characters_number_int":p_max_tag_characters_number_int,
				},
				nil, "gf_publisher_lib", p_runtime_sys)
			return nil, gf_err
		}
	}
	//---------------------
	return tags_lst, nil
}