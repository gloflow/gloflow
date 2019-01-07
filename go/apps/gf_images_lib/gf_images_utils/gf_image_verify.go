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

package gf_images_utils

import (
	"fmt"
	"strings"
	"github.com/davecgh/go-spew/spew"
	"gf_core"
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
func Image__verify_image_info(p_image_info_map map[string]interface{},
				p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{},*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_image_verify.Image__verify_image_info()")
	
	spew.Dump(p_image_info_map)

	max_title_characters_int := 100
	//-------------------
	//'id_str' - is None if image_info_dict comes from the outside of the system
	//           and the ID has not yet been assigned. 
	//           ID is determined if image_info_dict comes from the DB or someplace
	//           else within the system
	
	id_str := p_image_info_map["id_str"].(string)
	//-------------------
	//TITLE
	
	if _,ok := p_image_info_map["title_str"]; !ok {
		gf_err := gf_core.Error__create("image title_str not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}

	title_str := p_image_info_map["title_str"].(string)
		
	if len(title_str) > max_title_characters_int {

		usr_msg_str := fmt.Sprintf("image title_str length (%d) is longer then max_title_characters_int (%d)",
							len(title_str),
							max_title_characters_int)
		gf_err := gf_core.Error__create(usr_msg_str,
			"verify__string_too_long_error",
			&map[string]interface{}{
				"image_info_map":          p_image_info_map,
				"max_title_characters_int":max_title_characters_int,
			},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}
	//-------------------
	//IMAGE_CLIENT_TYPE
	
	if _,ok := p_image_info_map["image_client_type_str"]; !ok {
		gf_err := gf_core.Error__create("image image_client_type_str not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}
	
	image_client_type_str := p_image_info_map["image_client_type_str"].(string)
	//-------------------
	//FORMAT
	
	if _,ok := p_image_info_map["format_str"]; !ok {
		gf_err := gf_core.Error__create("image format_str not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}
	
	lower_case_format_str := strings.ToLower(p_image_info_map["format_str"].(string))
	
	normalized_format_str,ok := Image__check_image_format(lower_case_format_str,
		                                           		p_runtime_sys)
	if !ok {
		gf_err := gf_core.Error__create(fmt.Sprintf("invalid image extension (%s) found in image_info - %s",lower_case_format_str,title_str),
			"verify__invalid_image_extension_error",
			&map[string]interface{}{
				"title_str":            title_str,
				"lower_case_format_str":lower_case_format_str,
			},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}
	//-------------------
	//WIDTH/HEIGHT
	
	if _,ok := p_image_info_map["width_int"]; !ok {
		gf_err := gf_core.Error__create("image width_int not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}

	if _,ok := p_image_info_map["height_int"]; !ok {
		gf_err := gf_core.Error__create("image height_int not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}

	width_int  := int(p_image_info_map["width_int"].(int))
	height_int := int(p_image_info_map["height_int"].(int))
	
	/*if _, err := strconv.Atoi(width_int); err != nil {
		return nil,errors.New("image width_int is not a digit")
	}
	if _, err := strconv.Atoi(height_int); err != nil {
		return nil,errors.New("image height_int is not a digit")
	}*/
	//-------------------
	//IMAGE FLOWS NAMES
	if _,ok := p_image_info_map["flows_names_lst"]; !ok {
		gf_err := gf_core.Error__create("image flows_names_lst not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}

	flows_names_lst := p_image_info_map["flows_names_lst"].([]string)
	//-------------------
	//ORIGIN URL
	
	if _,ok := p_image_info_map["origin_url_str"]; !ok {
		gf_err := gf_core.Error__create("image origin_url_str not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}

	origin_url_str := p_image_info_map["origin_url_str"].(string)
	//-------------------
	if _,ok := p_image_info_map["origin_page_url_str"]; !ok {
		gf_err := gf_core.Error__create("image origin_page_url_str not supplied",
			"verify__missing_key_error",
			&map[string]interface{}{"image_info_map":p_image_info_map,},
			nil,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}

	origin_page_url_str := p_image_info_map["origin_page_url_str"].(string)
	//-------------------
	original_file_internal_uri_str := p_image_info_map["original_file_internal_uri_str"].(string)
	thumbnail_small_url_str        := p_image_info_map["thumbnail_small_url_str"].(string)
	thumbnail_medium_url_str       := p_image_info_map["thumbnail_medium_url_str"].(string)
	thumbnail_large_url_str        := p_image_info_map["thumbnail_large_url_str"].(string)
	//-------------------
	p_runtime_sys.Log_fun("INFO",fmt.Sprintf("image (id - %s) verified",id_str))
	
	verified_image_info_map := map[string]interface{}{
		"id_str":                        id_str,
		"title_str":                     title_str,
		"image_client_type_str":         image_client_type_str,
		"flows_names_lst":               flows_names_lst,
		"origin_url_str":                origin_url_str,
		"origin_page_url_str":           origin_page_url_str,
		"original_file_internal_uri_str":original_file_internal_uri_str,
		"thumbnail_small_url_str":       thumbnail_small_url_str,
		"thumbnail_medium_url_str":      thumbnail_medium_url_str,
		"thumbnail_large_url_str":       thumbnail_large_url_str,
		"format_str":                    normalized_format_str,
		"width_int":                     width_int,
		"height_int":                    height_int,
		
		//"dominant_color_hex_str":p_image_info_map["dominant_color_hex_str"],
	}
	
	return verified_image_info_map,nil
}