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

package gf_publisher_core

import (
	"fmt"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// external post_info is the one that comes from outside the system
// (it does not have an id assigned to it)

func Verify_external_post_info(p_post_info_map map[string]interface{},
	p_max_title_chars_int       int, // 100
	p_max_description_chars_int int, // 1000
	p_post_element_tag_max_int  int, // 20
	p_runtime_sys               *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_post_verify.Verify_external_post_info()")

	//-------------------
	// CLIENT_TYPE
	if _,ok := p_post_info_map["client_type_str"]; !ok {
		gf_err := gf_core.ErrorCreate("post client_type_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"post_info_map":p_post_info_map,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}
	client_type_str := p_post_info_map["client_type_str"].(string)

	//-------------------
	// TITLE
	if _,ok := p_post_info_map["title_str"]; !ok {
		gf_err := gf_core.ErrorCreate("post title_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"post_info_map":p_post_info_map,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}
	title_str := p_post_info_map["title_str"].(string)

	if len(title_str) > p_max_title_chars_int {
		gf_err := gf_core.ErrorCreate(fmt.Sprintf("title_str is longer (%d) then the max allowed number of chars (%d)", len(title_str), p_max_title_chars_int),
			"verify__string_too_long_error",
			map[string]interface{}{
				"title_str":           title_str,
				"max_title_chars_int": p_max_title_chars_int,
			},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	// ATTENTION!!
	// FB is removing/having problems with these symbols in url endings, and since the url to posts is composed of 
	// the post title, FB breaks these links
	// so striping them off right here avoids that

	clean_title_str   := title_str
	replace_chars_lst := []string{"[",",",":","#","%","&","!","]","$"}
	for _,c := range replace_chars_lst {
		strings.Replace(clean_title_str,c,"",-1)
	}

	//-------------------
	// DESCRIPTION
	if _,ok := p_post_info_map["description_str"]; !ok {
		gf_err := gf_core.ErrorCreate("post description_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"post_info_map":p_post_info_map,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}
	description_str := p_post_info_map["description_str"].(string)

	if len(description_str) > p_max_description_chars_int {
		gf_err := gf_core.ErrorCreate(fmt.Sprintf("description_str is longer (%d) then the max allowed number of chars (%d)", len(description_str), p_max_description_chars_int),
			"verify__string_too_long_error",
			map[string]interface{}{
				"description_str":          description_str,
				"max_description_chars_int":p_max_description_chars_int,
			},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	//-------------------	
	// TAGS
	tags_lst, gf_err := verify_tags(p_post_info_map, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//-------------------
	// POSTER_USER_NAME
	if _,ok := p_post_info_map["poster_user_name_str"]; !ok {
		gf_err := gf_core.ErrorCreate("post poster_user_name_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"post_info_map":p_post_info_map,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	//-------------------
	// POST ELEMENTS
	gf_err = verify_post_elements(p_post_info_map, p_post_element_tag_max_int, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//-------------------

	// "id_str" - not included here since p_post_info_map comes from outside the system
	//            and the internal id"s are for now not passed outside (or coming in from outside)
	verified_post_info_map := map[string]interface{}{
		"client_type_str":      client_type_str,
		"title_str":            clean_title_str,
		"description_str":      description_str,
		"tags_lst":             tags_lst,
		"poster_user_name_str": p_post_info_map["poster_user_name_str"].(string),
		"post_elements_lst":    p_post_info_map["post_elements_lst"],
	}
	
	return verified_post_info_map, nil
}

//---------------------------------------------------
func verify_tags(p_post_info_map map[string]interface{}, p_runtime_sys *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) { 
	p_runtime_sys.LogFun("FUN_ENTER","gf_post_verify.verify_tags()")
		
	if _, ok := p_post_info_map["tags_str"]; !ok {
		gf_err := gf_core.ErrorCreate("p_post_info_map doesnt contain the tags_str key",
			"verify__missing_key_error",
			map[string]interface{}{"post_info_map":p_post_info_map,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	input_tags_str := p_post_info_map["tags_str"].(string)
	tags_lst       := strings.Split(input_tags_str," ")

	p_runtime_sys.LogFun("INFO","input_tags_str - "+fmt.Sprint(input_tags_str))
	p_runtime_sys.LogFun("INFO","tags_lst       - "+fmt.Sprint(tags_lst))

	return tags_lst, nil
}

//---------------------------------------------------
func verify_post_elements(p_post_info_map map[string]interface{},
	p_post_element_tag_max_int int,
	p_runtime_sys              *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER","gf_post_verify.verify_post_elements()")
	
	if _, ok := p_post_info_map["post_elements_lst"]; !ok {
		gf_err := gf_core.ErrorCreate("p_post_info_map doesnt contain the post_elements_lst key",
			"verify__missing_key_error",
			map[string]interface{}{"post_info_map":p_post_info_map,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}
	post_elements_lst := p_post_info_map["post_elements_lst"].([]interface{})

	// verify each individiaul post_element
	for _,post_element := range post_elements_lst {
		post_element_map := post_element.(map[string]interface{})
		gf_err           := verify_post_element(post_element_map, p_post_element_tag_max_int, p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		//------------------------
		// SECURITY
		// ADD!! - have a external-url checking routines/whitelists/blacklists
		//         and other url sanitization routines,
		//         to prevent various XSS attacks
		//------------------------
	}

	return nil
}
//---------------------------------------------------
func verify_post_element(p_post_element_info_map map[string]interface{},
	p_post_element_tag_max_int int, //20
	p_runtime_sys              *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER","gf_post_verify.verify_post_element()")
	p_runtime_sys.LogFun("INFO"     ,"p_post_element_info_map - "+fmt.Sprint(p_post_element_info_map))

	//--------------
	// POST_ELEMENT_TYPE
	post_element_type_str := p_post_element_info_map["type_str"].(string)

	gf_err := Verify_post_element_type(post_element_type_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	if (post_element_type_str == "link"  ||
		post_element_type_str == "image" ||
		post_element_type_str == "video" ||
		post_element_type_str == "iframe") {	 

		// FIX!! - new versions of post_element_info_map format use extern_url_str
		//         instead of url_str. so when all post"s in the DB are updated to this format
		//         remove p_post_element_info_map has_key "url_str" check from this assert
		if !(gf_core.MapHasKey(p_post_element_info_map,"url_str") || gf_core.MapHasKey(p_post_element_info_map,"extern_url_str")) {
		
			gf_err := gf_core.ErrorCreate("p_post_element_info_map doesnt contain url_str|extern_url_str",
				"verify__missing_key_error",
				map[string]interface{}{"post_element_info_map":p_post_element_info_map,},
				nil, "gf_publisher_lib", p_runtime_sys)
			return gf_err
		}
	}

	//--------------
	// TAGS - OPTIONAL 
	if pe_tags_lst, ok := p_post_element_info_map["tags_lst"]; ok {
		for _,tag_str := range pe_tags_lst.([]string) {

			if len(tag_str) >= p_post_element_tag_max_int {
				gf_err := gf_core.ErrorCreate(fmt.Sprintf("tag (%s) is longer then max chars per tag (%d)", tag_str, p_post_element_tag_max_int),
					"verify__string_too_long_error",
					map[string]interface{}{
						"tag_str":                  tag_str,
						"post_element_tag_max_int": p_post_element_tag_max_int,
					},
					nil, "gf_publisher_lib", p_runtime_sys)
				return gf_err	
			}
		}
	}
	
	//--------------
	return nil
}
//---------------------------------------------------
func Verify_post_element_type(p_type_str string, p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {

	if !(p_type_str == "link"  ||
		p_type_str == "image"  ||
		p_type_str == "video"  ||
		p_type_str == "iframe" ||
		p_type_str == "text") {
		
		gf_err := gf_core.ErrorCreate(fmt.Sprintf("post_element type_str not of value image|link|video|iframe|text - instead its - %s", p_type_str),
			"verify__invalid_value_error",
			map[string]interface{}{"post_element_type_str": p_type_str,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}