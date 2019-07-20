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
	"os"
	"strings"
	"path"
	"path/filepath"
	"net/url"
	"io"
	"image"
	"image/jpeg"
	"image/png"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func Image__load_file(p_image_local_file_path_str string,
	p_normalized_ext_str string,
	p_runtime_sys        *gf_core.Runtime_sys) (image.Image,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_utils.Image__load_file()")

	file, fs_err := os.Open(p_image_local_file_path_str)
	if fs_err != nil {
		gf_err := gf_core.Error__create("failed to open a local file to load the image",
			"file_open_error",
			map[string]interface{}{
				"local_image_file_path_str": p_image_local_file_path_str,
			},
			fs_err, "gf_images_utils", p_runtime_sys)
		return nil,gf_err
	}
	defer file.Close()

	var img     image.Image
	var img_err error
	
	if p_normalized_ext_str == "png" {
		//PNG
		img,img_err = png.Decode(file)
		if img_err != nil {
			gf_err := gf_core.Error__create("failed to decode PNG file while transforming image",
				"png_decoding_error",
				map[string]interface{}{
					"local_image_file_path_str":p_image_local_file_path_str,
				},
				img_err, "gf_images_utils", p_runtime_sys)
			return nil,gf_err
		}
	} else {
		//JPEG,etc.
		img,_,img_err = image.Decode(file)
		if img_err != nil {
			gf_err := gf_core.Error__create("failed to decode image file while transforming image",
				"image_decoding_error",
				map[string]interface{}{
					"local_image_file_path_str":p_image_local_file_path_str,
				},
				img_err, "gf_images_utils", p_runtime_sys)
			return nil, gf_err
		}
	}

	return img, nil
}

//---------------------------------------------------
//VAR
//---------------------------------------------------
func Get_image_original_filename_from_url(p_image_url_str string, p_runtime_sys *gf_core.Runtime_sys) (string,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.Get_image_original_filename_from_url()")

	url,err := url.Parse(p_image_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse image_url to get image filename",
			"url_parse_error",
			map[string]interface{}{"image_url_str": p_image_url_str,},
			err, "gf_images_utils", p_runtime_sys)
		return "", gf_err
	}

	image_path_str      := url.Path
	image_file_name_str := path.Base(image_path_str)
	return image_file_name_str, nil
}

//---------------------------------------------------
func Get_image_title_from_url(p_image_url_str string,
	p_runtime_sys *gf_core.Runtime_sys) (string,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.Get_image_title_from_url()")
	
	url,err := url.Parse(p_image_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse image_url to get image title",
			"url_parse_error",
			map[string]interface{}{"image_url_str":p_image_url_str,},
			err, "gf_images_utils", p_runtime_sys)
		return "", gf_err
	}
	image_path_str      := url.Path
	image_file_name_str := path.Base(image_path_str)
	image_title_str     := strings.Split(image_file_name_str, ".")[0]
	
	return image_title_str, nil
}

//---------------------------------------------------
func Get_image_dimensions__from_image(p_img image.Image, p_runtime_sys *gf_core.Runtime_sys) (int,int) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.Get_image_dimensions__from_image()")

	p          := p_img.Bounds()
	width_int  := p.Max.X - p.Min.X
	height_int := p.Max.Y - p.Min.Y
	return width_int, height_int
}

//---------------------------------------------------
func Get_image_dimensions__from_filepath(p_image_local_file_path_str string, p_runtime_sys *gf_core.Runtime_sys) (int,int,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_utils.Get_image_dimensions__from_filepath()")

	//-------------------
	file, fs_err := os.Open(p_image_local_file_path_str)
	if fs_err != nil {
		gf_err := gf_core.Error__create("failed to open a local image file to get its dimensions",
			"file_open_error",
			map[string]interface{}{"image_local_file_path_str":p_image_local_file_path_str,},
			fs_err, "gf_images_utils", p_runtime_sys)
		return 0, 0, gf_err
	}
	defer file.Close()
	//-------------------
	format,gf_err := Get_image_ext_from_url(p_image_local_file_path_str,p_runtime_sys)
	if gf_err != nil {
		return 0, 0, gf_err
	}
	//-------------------
	image_width_int,image_height_int,gf_err := Get_image_dimensions__from_file(file,format,p_runtime_sys)
	if gf_err != nil {
		return 0, 0, gf_err
	}
	//-------------------
	return image_width_int, image_height_int, nil
}

//---------------------------------------------------
func Get_image_dimensions__from_file(p_file io.Reader,
	p_img_extension_str string,
	p_runtime_sys       *gf_core.Runtime_sys) (int,int,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.Get_image_dimensions__from_file()")

	var image_config image.Config
	var config_err   error

	//-------------------
	//JPEG
	if p_img_extension_str == "jpeg" {
		image_config, config_err = jpeg.DecodeConfig(p_file)
		if config_err != nil {
			gf_err := gf_core.Error__create("failed to decode config for JPEG image file to get image dimensions",
				"image_decoding_config_error",
				map[string]interface{}{"img_extension_str":p_img_extension_str,},
				config_err, "gf_images_utils", p_runtime_sys)
			return 0, 0, gf_err
		}
	//-------------------
	//PNG
	} else if p_img_extension_str == "png" {
		image_config, config_err = png.DecodeConfig(p_file)
		if config_err != nil {
			gf_err := gf_core.Error__create("failed to decode config for PNG image file to get image dimensions",
				"image_decoding_config_error",
				map[string]interface{}{"img_extension_str":p_img_extension_str,},
				config_err, "gf_images_utils", p_runtime_sys)
			return 0, 0, gf_err
		}
	//-------------------
	//GENERAL
	} else {
		image_config, _, config_err = image.DecodeConfig(p_file)
		if config_err != nil {
			gf_err := gf_core.Error__create("failed to decode config for image file to get image dimensions",
				"image_decoding_config_error",
				map[string]interface{}{"img_extension_str":p_img_extension_str,},
				config_err, "gf_images_utils", p_runtime_sys)
			return 0, 0, gf_err
		}
	}
	//-------------------

	image_width_int  := image_config.Width
	image_height_int := image_config.Height

	return image_width_int, image_height_int, nil
}

//---------------------------------------------------
func Get_image_ext_from_url(p_image_url_str string, p_runtime_sys *gf_core.Runtime_sys) (string,*gf_core.Gf_error) {
	//p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.Get_image_ext_from_url()")
	
	fmt.Println("p_image_url_str - "+p_image_url_str)

	//urlparse() - used so that any possible url query parameters are not used in the 
	//             os.path.basename() result
	url,err := url.Parse(p_image_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse image_url to get image extension",
			"url_parse_error",
			map[string]interface{}{"image_url_str":p_image_url_str,},
			err, "gf_images_utils", p_runtime_sys)
		return "", gf_err
	}

	image_path_str      := url.Path
	image_file_name_str := path.Base(image_path_str)
	ext_str             := filepath.Ext(image_file_name_str)
	clean_ext_str       := strings.TrimPrefix(strings.ToLower(ext_str),".")

	normalized_ext_str,ok := Image__check_image_format(clean_ext_str,p_runtime_sys)

	if !ok {
		gf_err := gf_core.Error__create(fmt.Sprintf("invalid image extension (%s) found in image url - %s",ext_str,p_image_url_str),
			"verify__invalid_image_extension_error",
			map[string]interface{}{
				"image_url_str":p_image_url_str,
				"ext_str":      ext_str,
			},
			nil, "gf_images_utils", p_runtime_sys)
		return "", gf_err
	}


	return normalized_ext_str,nil
}

//---------------------------------------------------	
func Image__check_image_format(p_format_str string, p_runtime_sys *gf_core.Runtime_sys) (string,bool) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.Image__check_image_format()")
	
	if p_format_str != "jpeg" && 
		p_format_str != "jpg" &&
		p_format_str != "gif" &&
		p_format_str != "png" {

		//IMPORTANT!! - format is not a valid image format, so 'false' is returned
		return p_format_str, false
	}

	var normalized_format_str string

	//normalize "jpg" (variation on jpeg) to "jpeg"
	if p_format_str == "jpg" {
		normalized_format_str = "jpeg"
	} else {
		normalized_format_str = p_format_str
	}

	return normalized_format_str, true
}

//---------------------------------------------------
//IMPORTANT!! - look at JS library, for content-aware image cropping
//              https://github.com/jwagner/smartcrop.js/
//---------------------------------------------------
//ADD!! - a function that will detect how many duplicates are there in the image DB/collection
//		 use "fdupes" for this (command line utility)
//		 fdupes - fdupes is a program written by Adrian Lopez to scan directories for duplicate files, 
//				  with options to list, delete or replace the files with hardlinks pointing to the 
//				  duplicate. It first compares file sizes and MD5 signatures, and then 
//				  performs a byte-by-byte check for verification.
//---------------------------------------------------
//ADD!! - calculate RGBA histogram for every image
/*var histogram [16][4]int
for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		r, g, b, a := m.At(x, y).RGBA()

		// A color's RGBA method returns values in the range [0, 65535].
		// Shifting by 12 reduces this to the range [0, 15].
		histogram[r>>12][0]++
		histogram[g>>12][1]++
		histogram[b>>12][2]++
		histogram[a>>12][3]++
	}
}*/
//---------------------------------------------------
//TAGS
//---------------------------------------------------
/*func add_tags_to_image_in_db(p_id_str string,
					p_tags_lst   []string,
					p_mongo_coll *mgo.Collection,
					p_log_fun    func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_images_utils.add_tags_to_image_in_db()")
	
	image,err := image__db_get(p_id_str,
						p_mongo_coll,
						p_log_fun)
	if err != nil {
		return err
	}

	add_tags_to_image(post_adt,
					  p_tags_lst,
					  p_log_fun)
	
	i_err := image__db_put(image,
					p_mongo_coll,
					p_log_fun)
	if i_err != nil {
		return i_err
	}

	return nil
}*/