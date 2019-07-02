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

package gf_images_utils

import (
	"fmt"
	"image"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func Create_thumbnails(p_image_id_str Gf_image_id,
	p_image_format_str                     string,
	p_image_file_path_str                  string,
	p_local_target_thumbnails_dir_path_str string,
	p_small_thumb_max_size_px_int          int,
	p_medium_thumb_max_size_px_int         int,
	p_large_thumb_max_size_px_int          int,
	p_image                                image.Image,
	p_runtime_sys                          *gf_core.Runtime_sys) (*Gf_image_thumbs, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_thumbs.Create_thumbnails()")

	//-----------------
	//SMALL THUMBS
	new_thumb_small_file_name_str         := fmt.Sprintf("%s_thumb_small.%s", p_image_id_str, p_image_format_str)
	small__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", p_local_target_thumbnails_dir_path_str, new_thumb_small_file_name_str)

	gf_err := resize_image(p_image, //p_image_file,
		small__target_thumbnail_file_path_str,
		p_image_format_str,
		p_small_thumb_max_size_px_int,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	//-----------------
	//MEDIUM THUMBS
	new_thumb_medium_file_name_str         := fmt.Sprintf("%s_thumb_medium.%s", p_image_id_str, p_image_format_str)
	medium__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", p_local_target_thumbnails_dir_path_str, new_thumb_medium_file_name_str)

	gf_err = resize_image(p_image, //p_image_file,
		medium__target_thumbnail_file_path_str,
		p_image_format_str,
		p_medium_thumb_max_size_px_int,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	//-----------------
	//LARGE THUMBS
	new_thumb_large_file_name_str         := fmt.Sprintf("%s_thumb_large.%s",p_image_id_str,p_image_format_str)
	large__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s",p_local_target_thumbnails_dir_path_str,new_thumb_large_file_name_str)

	gf_err = resize_image(p_image, //p_image_file,
		large__target_thumbnail_file_path_str,
		p_image_format_str,
		p_large_thumb_max_size_px_int,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	//-----------------

	thumb_small_relative_url_str  := "/images/d/thumbnails/"+new_thumb_small_file_name_str
	thumb_medium_relative_url_str := "/images/d/thumbnails/"+new_thumb_medium_file_name_str
	thumb_large_relative_url_str  := "/images/d/thumbnails/"+new_thumb_large_file_name_str

	image_thumbs := &Gf_image_thumbs{
		Small_relative_url_str:     thumb_small_relative_url_str,
		Medium_relative_url_str:    thumb_medium_relative_url_str,
		Large_relative_url_str:     thumb_large_relative_url_str,
		Small_local_file_path_str:  small__target_thumbnail_file_path_str,
		Medium_local_file_path_str: medium__target_thumbnail_file_path_str,
		Large_local_file_path_str:  large__target_thumbnail_file_path_str,
	}

	return image_thumbs, nil
}