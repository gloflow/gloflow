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
	"os"
	"image"
	"image/jpeg"
	"github.com/nfnt/resize"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
//p_image_origin_page_url_str - urls of pages (html or some other resource) where the image image_url
//                              was found. this is valid for gf_chrome_ext image sources.
//                              its not relevant for direct image uploads from clients.

func Transform_image(p_image_id_str string,
	p_image_client_type_str                      string,
	p_images_flows_names_lst                     []string,
	p_image_origin_url_str                       string,
	p_image_origin_page_url_str                  string,
	p_image_local_file_path_str                  string,
	p_images_store_thumbnails_local_dir_path_str string,
	p_runtime_sys                                *gf_core.Runtime_sys) (*Gf_image,*Gf_image_thumbs,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_transformer.Transform_image()")

	normalized_ext_str,gf_err := Get_image_ext_from_url(p_image_origin_url_str,p_runtime_sys)
	if gf_err != nil {
		return nil,nil,gf_err
	}

	gf_image,gf_image_thumbs,gf_err := Trans__process_image(p_image_id_str,
		p_image_client_type_str,
		p_images_flows_names_lst,
		p_image_origin_url_str,
		p_image_origin_page_url_str,
		normalized_ext_str,
		p_image_local_file_path_str,
		p_images_store_thumbnails_local_dir_path_str,
		p_runtime_sys)
	if gf_err != nil {
		return nil,nil,gf_err
	}

	return gf_image,gf_image_thumbs,nil
}
//---------------------------------------------------
func Trans__process_image(p_image_id_str string,
	p_image_client_type_str                string,
	p_images_flows_names_lst               []string,
	p_image_origin_url_str                 string,
	p_image_origin_page_url_str            string,
	p_normalized_ext_str                   string,
	p_image_local_file_path_str            string,
	p_local_thumbnails_target_dir_path_str string,
	p_runtime_sys                          *gf_core.Runtime_sys) (*Gf_image,*Gf_image_thumbs,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_transformer.Trans__process_image()")
	fmt.Println("p_image_local_file_path_str - "+p_image_local_file_path_str)

	//---------------------------------
	//LOAD_IMAGE

	img,gf_err := Image__load_file(p_image_local_file_path_str,p_normalized_ext_str,p_runtime_sys)
	if gf_err != nil {
		return nil, nil, gf_err
	}
	//--------------------------
	//CREATE THUMBNAILS

	small_thumb_max_size_px_int  := 200
	medium_thumb_max_size_px_int := 400
	large_thumb_max_size_px_int  := 600

	gf_image_thumbs,gf_err := Create_thumbnails(p_image_id_str,
		p_normalized_ext_str,
		p_image_local_file_path_str,
		p_local_thumbnails_target_dir_path_str,
		small_thumb_max_size_px_int,
		medium_thumb_max_size_px_int,
		large_thumb_max_size_px_int,
		img,
		p_runtime_sys)
	if gf_err != nil {
		return nil, nil, gf_err
	}
	//--------------------------
	/*//DOMINANT COLOR DETERMINATION
	//it"s computed only for non-gif"s
	dominant_color_hex_str := gf_images_utils_graphic.get_dominant_image_color(p_image_local_file_path_str,p_log_fun)*/
	//--------------------------
	image_width_int,image_height_int := Get_image_dimensions__from_image(img,p_runtime_sys)
	//--------------------------

	//SECURITY ISSUE!!
	//When you open a file, the file header is read to determine the file 
	//format and extract things like mode, size, and other properties 
	//required to decode the file, but the rest of the file is not processed until later
	
	//someone can forge header information in an image
	image_title_str,gf_err := Get_image_title_from_url(p_image_origin_url_str,p_runtime_sys)
	if gf_err != nil {
		return nil,nil,gf_err
	}

	gf_image_info := &Gf_image_new_info{
		Id_str:                        p_image_id_str,
		Title_str:                     image_title_str,
		Flows_names_lst:               p_images_flows_names_lst,
		Image_client_type_str:         p_image_client_type_str,
		Origin_url_str:                p_image_origin_url_str,
		Origin_page_url_str:           p_image_origin_page_url_str,
		Original_file_internal_uri_str:p_image_local_file_path_str,
		Thumbnail_small_url_str:       gf_image_thumbs.Small_relative_url_str,
		Thumbnail_medium_url_str:      gf_image_thumbs.Medium_relative_url_str,
		Thumbnail_large_url_str:       gf_image_thumbs.Large_relative_url_str,
		Format_str:                    p_normalized_ext_str,
		Width_int:                     image_width_int,
		Height_int:                    image_height_int,
	}
	//--------------------------
	//IMAGE_CREATE

	//IMPORTANT!! - creates a GF_Image struct and stores it in the DB
	gf_image,gf_err := Image__create_new(gf_image_info,p_runtime_sys)
	if gf_err != nil {
		return nil,nil,gf_err
	}
	//--------------------------

	return gf_image,gf_image_thumbs,nil
}
//---------------------------------------------------
func resize_image(p_img image.Image,
	p_image_output_path_str string,
	p_image_format_str      string,
	p_size_px_int           int,
	p_runtime_sys           *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_transformer.resize_image()")
	
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	
	m := resize.Resize(uint(p_size_px_int), 0, p_img, resize.Lanczos3)

	out, err := os.Create(p_image_output_path_str)
	if err != nil {
		gf_err := gf_core.Error__create("OS failed to create a file to save a resized image to FS",
			"file_create_error",
			&map[string]interface{}{"image_output_path_str":p_image_output_path_str,},
			err,"gf_images_utils",p_runtime_sys)
		return gf_err
	}
	defer out.Close()

	//IMPORTANT!! - using JPEG instead of PNG, because JPEG compression was made for photographic images,
	//              and so for these kinds of images it comes out with much smaller file size
	// write new image to file
	//jpeg.Encode(out, m, nil)
	jpeg.Encode(out,m,nil)

	return nil
}