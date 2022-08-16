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

package gf_images_core

import (
	"fmt"
	"strings"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
/*map[
	origin_url_str:http://66.media.tumblr.com/b5b70dc4e994c5111c5318177b965b57/tumblr_o87f7n8XZr1vqf2mpo1_500.jpg 
	original_file_internal_uri_str:/home/gf/data/images/4615191efaae8f4b530b31d5d26bbd9f.jpeg 
	height_int:0 
	id_str:4615191efaae8f4b530b31d5d26bbd9f 
	title_str:tumblr_o87f7n8XZr1vqf2mpo1_500 
	format_str:jpeg 
	width_int:0 
	thumbnail_small_url_str:/images/d/thumbnails/4615191efaae8f4b530b31d5d26bbd9f_thumb_small.jpeg 
	thumbnail_medium_url_str:/images/d/thumbnails/4615191efaae8f4b530b31d5d26bbd9f_thumb_medium.jpeg 
	thumbnail_large_url_str:/images/d/thumbnails/4615191efaae8f4b530b31d5d26bbd9f_thumb_large.jpeg
]*/

//---------------------------------------------------
func Image__verify_image_info(pImageInfoMap map[string]interface{},
	pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	pRuntimeSys.Log_fun("FUN_ENTER","gf_image_verify.Image__verify_image_info()")
	
	spew.Dump(pImageInfoMap)

	max_title_characters_int := 100
	//-------------------
	//'id_str' - is None if image_info_dict comes from the outside of the system
	//           and the ID has not yet been assigned. 
	//           ID is determined if image_info_dict comes from the DB or someplace
	//           else within the system
	
	id_str := pImageInfoMap["id_str"].(string)

	//-------------------
	// TITLE
	
	if _, ok := pImageInfoMap["title_str"]; !ok {
		gfErr := gf_core.ErrorCreate("image title_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map": pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	title_str := pImageInfoMap["title_str"].(string)
		
	if len(title_str) > max_title_characters_int {

		usr_msg_str := fmt.Sprintf("image title_str length (%d) is longer then max_title_characters_int (%d)", len(title_str), max_title_characters_int)
		gfErr := gf_core.ErrorCreate(usr_msg_str,
			"verify__string_too_long_error",
			map[string]interface{}{
				"image_info_map":           pImageInfoMap,
				"max_title_characters_int": max_title_characters_int,
			},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	//-------------------
	// IMAGE_CLIENT_TYPE
	
	if _, ok := pImageInfoMap["image_client_type_str"]; !ok {
		gfErr := gf_core.ErrorCreate("image image_client_type_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map": pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}
	
	image_client_type_str := pImageInfoMap["image_client_type_str"].(string)

	//-------------------
	// FORMAT
	
	if _, ok := pImageInfoMap["format_str"]; !ok {
		gfErr := gf_core.ErrorCreate("image format_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map": pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}
	
	lowercaseFormatStr := strings.ToLower(pImageInfoMap["format_str"].(string))
	
	ok := CheckImageFormat(lowercaseFormatStr, pRuntimeSys)
	if !ok {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("invalid image extension (%s) found in image_info - %s", lowercaseFormatStr, title_str),
			"verify__invalid_image_extension_error",
			map[string]interface{}{
				"title_str":             title_str,
				"lower_case_format_str": lowercaseFormatStr,
			},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	normalizedFormatStr := NormalizeImageFormat(lowercaseFormatStr)

	//-------------------
	// WIDTH/HEIGHT
	
	if _, ok := pImageInfoMap["width_int"]; !ok {
		gfErr := gf_core.ErrorCreate("image width_int not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map": pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	if _, ok := pImageInfoMap["height_int"]; !ok {
		gfErr := gf_core.ErrorCreate("image height_int not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map": pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	width_int  := int(pImageInfoMap["width_int"].(int))
	height_int := int(pImageInfoMap["height_int"].(int))
	
	/*if _, err := strconv.Atoi(width_int); err != nil {
		return nil,errors.New("image width_int is not a digit")
	}
	if _, err := strconv.Atoi(height_int); err != nil {
		return nil,errors.New("image height_int is not a digit")
	}*/

	//-------------------
	// IMAGE FLOWS NAMES
	if _, ok := pImageInfoMap["flows_names_lst"]; !ok {
		gfErr := gf_core.ErrorCreate("image flows_names_lst not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map": pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	flows_names_lst := pImageInfoMap["flows_names_lst"].([]string)

	//-------------------
	// ORIGIN URL
	
	if _, ok := pImageInfoMap["origin_url_str"]; !ok {
		gfErr := gf_core.ErrorCreate("image origin_url_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map": pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	origin_url_str := pImageInfoMap["origin_url_str"].(string)

	//-------------------
	if _, ok := pImageInfoMap["origin_page_url_str"]; !ok {
		gfErr := gf_core.ErrorCreate("image origin_page_url_str not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"image_info_map":pImageInfoMap,},
			nil, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	origin_page_url_str := pImageInfoMap["origin_page_url_str"].(string)

	//-------------------
	original_file_internal_uri_str := pImageInfoMap["original_file_internal_uri_str"].(string)
	thumbnail_small_url_str        := pImageInfoMap["thumbnail_small_url_str"].(string)
	thumbnail_medium_url_str       := pImageInfoMap["thumbnail_medium_url_str"].(string)
	thumbnail_large_url_str        := pImageInfoMap["thumbnail_large_url_str"].(string)
	
	//-------------------
	pRuntimeSys.Log_fun("INFO",fmt.Sprintf("image (id - %s) verified",id_str))
	
	verifiedImageInfoMap := map[string]interface{}{
		"id_str":                         id_str,
		"title_str":                      title_str,
		"image_client_type_str":          image_client_type_str,
		"flows_names_lst":                flows_names_lst,
		"origin_url_str":                 origin_url_str,
		"origin_page_url_str":            origin_page_url_str,
		"original_file_internal_uri_str": original_file_internal_uri_str,
		"thumbnail_small_url_str":        thumbnail_small_url_str,
		"thumbnail_medium_url_str":       thumbnail_medium_url_str,
		"thumbnail_large_url_str":        thumbnail_large_url_str,
		"format_str":                     normalizedFormatStr,
		"width_int":                      width_int,
		"height_int":                     height_int,
		
		//"dominant_color_hex_str":p_image_info_map["dominant_color_hex_str"],
	}
	
	return verifiedImageInfoMap,nil
}