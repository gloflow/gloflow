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

package gf_tagger_lib

import (
	"fmt"
	"strings"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_address"
)

//---------------------------------------------------
// pTagsStr           - "," separated list of strings
// pObjectExternIDstr - this is an external identifier for an object, not necessarily its internal. 
//                      for posts - their p_object_extern_id_str is their Title, but internally they have
//                      another ID.

func addTagsToObject(pTagsStr string,
	pObjectTypeStr     string,
	pObjectExternIDstr string,
	pMetaMap           map[string]interface{},
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	if pObjectTypeStr != "post" &&
		pObjectTypeStr != "image" &&
		pObjectTypeStr != "event" &&
		pObjectTypeStr != "address" {

		gfErr := gf_core.ErrorCreate(fmt.Sprintf("object_type (%s) is not of supported type (post|image|event)",
			pObjectTypeStr),
			"verify__invalid_value_error",
			map[string]interface{}{
				"tags_str":        pTagsStr,
				"object_type_str": pObjectTypeStr,
			},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}
	
	tagsLst, gfErr := parse_tags(pTagsStr,
		500, // p_max_tags_bulk_size_int        int, // 500
		20,  // p_max_tag_characters_number_int int, // 20	
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	pRuntimeSys.Log_fun("INFO", fmt.Sprintf("tags_lst - %s", tagsLst))

	//---------------
	// POST
	
	switch pObjectTypeStr {
		//---------------
		// POST
		case "post":
			postTitleStr      := pObjectExternIDstr
			existsBool, gfErr := gf_publisher_core.DB__check_post_exists(postTitleStr,
				pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
			
			if existsBool {
				pRuntimeSys.Log_fun("INFO", "POST EXISTS")
				gfErr := db__add_tags_to_post(postTitleStr, tagsLst, pRuntimeSys)
				return gfErr

			} else {
				gfErr := gf_core.ErrorCreate(fmt.Sprintf("post with title (%s) doesnt exist, while adding a tags - %s", 
					postTitleStr,
					tagsLst),
					"verify__invalid_value_error",
					map[string]interface{}{
						"post_title_str": postTitleStr,
						"tags_lst":       tagsLst,
					},
					nil, "gf_tagger", pRuntimeSys)
				return gfErr
			}

		//---------------
		// IMAGE
		case "image":
			imageIDstr := pObjectExternIDstr
			image_id   := gf_images_core.GF_image_id(imageIDstr)
			exists_bool, gfErr := gf_images_core.DB__image_exists(image_id, pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
			if exists_bool {
				gfErr := db__add_tags_to_image(imageIDstr, tagsLst, pRuntimeSys)
				if gfErr != nil {
					return gfErr
				}
			}

		//---------------
		// WEB3
		case "address":

			chainStr := pMetaMap["chain_str"].(string)



			addressStr := pObjectExternIDstr
			existsBool, gfErr := gf_address.DBexists(addressStr,
				chainStr,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}

			
			if existsBool {
				gfErr := gf_address.DBaddTag(tagsLst,
					addressStr,
					chainStr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return gfErr
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
	p_runtime_sys     *gf_core.RuntimeSys) (map[string][]map[string]interface{}, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger.get_objects_with_tags()")
		
	objects_with_tags_map := map[string][]map[string]interface{}{}
	for _, tag_str := range p_tags_lst {
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
	p_runtime_sys     *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger.get_objects_with_tag()")
	p_runtime_sys.Log_fun("INFO",      fmt.Sprintf("p_object_type_str - %s", p_object_type_str))

	//ADD!! - add support for tagging "image" p_object_type_str's
	if p_object_type_str != "post" {
		gf_err := gf_core.ErrorCreate(fmt.Sprintf("trying to get objects with a tag (%s) for objects type thats not supported - %s", p_tag_str, p_object_type_str),
			"verify__invalid_value_error",
			map[string]interface{}{
				"tag_str":         p_tag_str,
				"object_type_str": p_object_type_str,
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
	for _, post := range posts_with_tag_lst {
		post_info_map := map[string]interface{}{
			"title_str":               post.Title_str,
			"tags_lst":                post.Tags_lst,
			"url_str":                 fmt.Sprintf("/posts/%s", post.Title_str),
			"object_type_str":         p_object_type_str,
			"thumbnail_small_url_str": post.Thumbnail_url_str,
		}
		min_posts_infos_lst = append(min_posts_infos_lst,post_info_map)
	}

	objects_infos_lst := min_posts_infos_lst
	return objects_infos_lst, nil
}

//---------------------------------------------------
func parse_tags(pTagsStr string,
	p_max_tags_bulk_size_int        int, // 500
	p_max_tag_characters_number_int int, // 20
	p_runtime_sys                   *gf_core.RuntimeSys) ([]string, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger.parse_tags()")
	
	tags_lst := strings.Split(pTagsStr," ")
	//---------------------
	if len(tags_lst) > p_max_tags_bulk_size_int {
		gf_err := gf_core.ErrorCreate(fmt.Sprintf("too many tags supplied - max is %d", p_max_tags_bulk_size_int),
			"verify__value_too_many_error",
			map[string]interface{}{
				"tags_lst":               tags_lst,
				"max_tags_bulk_size_int": p_max_tags_bulk_size_int,
			},
			nil, "gf_tagger_lib", p_runtime_sys)
		return nil, gf_err
	}

	//---------------------
	for _, tag_str := range tags_lst {
		if len(tag_str) > p_max_tag_characters_number_int {
			gf_err := gf_core.ErrorCreate(fmt.Sprintf("tag (%s) is too long - max is (%d)", tag_str, p_max_tag_characters_number_int),
				"verify__string_too_long_error",
				map[string]interface{}{
					"tag_str":                       tag_str,
					"max_tag_characters_number_int": p_max_tag_characters_number_int,
				},
				nil, "gf_tagger_lib", p_runtime_sys)
			return nil, gf_err
		}
	}
	
	//---------------------
	return tags_lst, nil
}