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


//------------------------------------------------
func ImageGetFilepathFromID(pImageIDstr GFimageID,
	pImageFormatStr string) string {
	imageFileNameStr := fmt.Sprintf("%s.%s", pImageIDstr, pImageFormatStr)
	return imageFileNameStr
}

//------------------------------------------------
// returns the URL of a particular image path,
// this URL is where the image can be fetched from directly.

func ImageGetPublicURL(pImageFilePathStr string,
	pMediaDomainStr string,
	pRuntimeSys      *gf_core.Runtime_sys) string {

	// // IMPORTANT!! - amazon URL escapes image file names when it makes them public in a bucket
	// //               escaped_str := url.QueryEscape(*p_image_s3_file_path_str)
	// url_str := fmt.Sprintf("http://%s.s3-website-us-east-1.amazonaws.com/%s", p_s3_bucket_name_str, p_image_s3_file_path_str)

	urlStr := fmt.Sprintf("https://%s/%s", pMediaDomainStr, pImageFilePathStr)
	return urlStr
}

//---------------------------------------------------
// LOAD_FILE
func ImageLoadFile(p_image_local_file_path_str string,
	p_normalized_ext_str string,
	pRuntimeSys        *gf_core.RuntimeSys) (image.Image, *gf_core.GFerror) {

	file, fs_err := os.Open(p_image_local_file_path_str)
	if fs_err != nil {
		gf_err := gf_core.Error__create("failed to open a local file to load the image",
			"file_open_error",
			map[string]interface{}{
				"local_image_file_path_str": p_image_local_file_path_str,
			},
			fs_err, "gf_images_core", pRuntimeSys)
		return nil, gf_err
	}
	defer file.Close()

	var img     image.Image
	var img_err error
	
	if p_normalized_ext_str == "png" {
		// PNG
		img, img_err = png.Decode(file)
		if img_err != nil {
			gf_err := gf_core.Error__create("failed to decode PNG file while transforming image",
				"png_decoding_error",
				map[string]interface{}{
					"local_image_file_path_str": p_image_local_file_path_str,
				},
				img_err, "gf_images_core", pRuntimeSys)
			return nil, gf_err
		}
	} else {
		// JPEG,etc.
		img, _, img_err = image.Decode(file)
		if img_err != nil {
			gf_err := gf_core.Error__create("failed to decode image file while transforming image",
				"image_decoding_error",
				map[string]interface{}{
					"local_image_file_path_str": p_image_local_file_path_str,
				},
				img_err, "gf_images_core", pRuntimeSys)
			return nil, gf_err
		}
	}

	return img, nil
}

//---------------------------------------------------
// VAR
//---------------------------------------------------
func Get_image_original_filename_from_url(p_image_url_str string,
	pRuntimeSys *gf_core.Runtime_sys) (string, *gf_core.Gf_error) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_utils.Get_image_original_filename_from_url()")

	url,err := url.Parse(p_image_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse image_url to get image filename",
			"url_parse_error",
			map[string]interface{}{"image_url_str": p_image_url_str,},
			err, "gf_images_core", pRuntimeSys)
		return "", gf_err
	}

	image_path_str      := url.Path
	image_file_name_str := path.Base(image_path_str)
	return image_file_name_str, nil
}

//---------------------------------------------------
func Create_gf_image_file_path_from_url(p_gf_image_id_str Gf_image_id,
	p_image_url_str                   string,
	p_images_store_local_dir_path_str string,
	pRuntimeSys                       *gf_core.Runtime_sys) (string, Gf_image_id, *gf_core.Gf_error) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_utils.Create_gf_image_file_path_from_url()")

	// IMPORTANT!! - gf_image_id can be supplied externally. if its not supplied then a new gf_image_id is generated
	var gf_image_id_str Gf_image_id
	if p_gf_image_id_str == "" {
		new_gf_image_id_str, _ := Image_ID__create_from_url(p_image_url_str, pRuntimeSys)
		gf_image_id_str = new_gf_image_id_str
	} else {
		gf_image_id_str = p_gf_image_id_str
	}

	ext_str, gf_err := GetImageExtFromURL(p_image_url_str, pRuntimeSys)
	if gf_err != nil {
		return "", "", gf_err
	}

	// IMPORTANT!! - 0.4 system, image naming, new scheme containing image_id,
	//               instead of the old original_image naming scheme.
	local_image_file_name_str := fmt.Sprintf("%s.%s", gf_image_id_str, ext_str)
	local_image_file_path_str := fmt.Sprintf("%s/%s", p_images_store_local_dir_path_str, local_image_file_name_str)

	pRuntimeSys.Log_fun("INFO", fmt.Sprintf("local_image_file_path_str - %s", local_image_file_path_str))
	
	return local_image_file_path_str, gf_image_id_str, nil
}

//---------------------------------------------------
func Get_image_title_from_url(p_image_url_str string,
	pRuntimeSys *gf_core.Runtime_sys) (string,*gf_core.Gf_error) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_utils.Get_image_title_from_url()")
	
	url, err := url.Parse(p_image_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse image_url to get image title",
			"url_parse_error",
			map[string]interface{}{"image_url_str":p_image_url_str,},
			err, "gf_images_core", pRuntimeSys)
		return "", gf_err
	}
	image_path_str      := url.Path
	image_file_name_str := path.Base(image_path_str)
	image_title_str     := strings.Split(image_file_name_str, ".")[0]
	
	return image_title_str, nil
}

//---------------------------------------------------
func Get_image_dimensions__from_image(p_img image.Image,
	pRuntimeSys *gf_core.RuntimeSys) (int, int) {

	p          := p_img.Bounds()
	width_int  := p.Max.X - p.Min.X
	height_int := p.Max.Y - p.Min.Y
	return width_int, height_int
}

//---------------------------------------------------
func Get_image_dimensions__from_filepath(p_image_local_file_path_str string,
	pRuntimeSys *gf_core.Runtime_sys) (int, int, *gf_core.GFerror) {

	//-------------------
	file, fs_err := os.Open(p_image_local_file_path_str)
	if fs_err != nil {
		gfErr := gf_core.Error__create("failed to open a local image file to get its dimensions",
			"file_open_error",
			map[string]interface{}{"image_local_file_path_str": p_image_local_file_path_str,},
			fs_err, "gf_images_core", pRuntimeSys)
		return 0, 0, gfErr
	}
	defer file.Close()

	//-------------------
	format, gfErr := GetImageExtFromURL(p_image_local_file_path_str, pRuntimeSys)
	if gfErr != nil {
		return 0, 0, gfErr
	}

	//-------------------
	image_width_int, image_height_int, gfErr := Get_image_dimensions__from_file(file, format, pRuntimeSys)
	if gfErr != nil {
		return 0, 0, gfErr
	}

	//-------------------
	return image_width_int, image_height_int, nil
}

//---------------------------------------------------
func Get_image_dimensions__from_file(p_file io.Reader,
	p_img_extension_str string,
	pRuntimeSys         *gf_core.Runtime_sys) (int, int, *gf_core.GFerror) {

	var image_config image.Config
	var config_err   error

	//-------------------
	// JPEG
	if p_img_extension_str == "jpeg" {
		image_config, config_err = jpeg.DecodeConfig(p_file)
		if config_err != nil {
			gf_err := gf_core.Error__create("failed to decode config for JPEG image file to get image dimensions",
				"image_decoding_config_error",
				map[string]interface{}{"img_extension_str": p_img_extension_str,},
				config_err, "gf_images_core", pRuntimeSys)
			return 0, 0, gf_err
		}

	//-------------------
	// PNG
	} else if p_img_extension_str == "png" {
		image_config, config_err = png.DecodeConfig(p_file)
		if config_err != nil {
			gf_err := gf_core.Error__create("failed to decode config for PNG image file to get image dimensions",
				"image_decoding_config_error",
				map[string]interface{}{"img_extension_str": p_img_extension_str,},
				config_err, "gf_images_core", pRuntimeSys)
			return 0, 0, gf_err
		}

	//-------------------
	// GENERAL
	} else {
		image_config, _, config_err = image.DecodeConfig(p_file)
		if config_err != nil {
			gf_err := gf_core.Error__create("failed to decode config for image file to get image dimensions",
				"image_decoding_config_error",
				map[string]interface{}{"img_extension_str":p_img_extension_str,},
				config_err, "gf_images_core", pRuntimeSys)
			return 0, 0, gf_err
		}
	}

	//-------------------

	image_width_int  := image_config.Width
	image_height_int := image_config.Height

	return image_width_int, image_height_int, nil
}

//---------------------------------------------------
func GetImageExtFromURL(pImageURLstr string,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	fmt.Printf("image url - %s\n", pImageURLstr)

	// urlparse() - used so that any possible url query parameters are not used in the 
	//              os.path.basename() result
	url, err := url.Parse(pImageURLstr)
	if err != nil {
		gfErr := gf_core.Error__create("failed to parse image_url to get image extension",
			"url_parse_error",
			map[string]interface{}{"image_url_str": pImageURLstr,},
			err, "gf_images_core", pRuntimeSys)
		return "", gfErr
	}

	imagePathStr      := url.Path
	imageFileNameStr := path.Base(imagePathStr)
	extStr             := filepath.Ext(imageFileNameStr)
	cleanExtStr       := strings.TrimPrefix(strings.ToLower(extStr),".")

	normalizedExtStr, ok := Image__check_image_format(cleanExtStr,pRuntimeSys)

	if !ok {
		gfErr := gf_core.Error__create(fmt.Sprintf("invalid image extension (%s) found in image url - %s", extStr, pImageURLstr),
			"verify__invalid_image_extension_error",
			map[string]interface{}{
				"image_url_str": pImageURLstr,
				"ext_str":       extStr,
			},
			nil, "gf_images_core", pRuntimeSys)
		return "", gfErr
	}

	return normalizedExtStr, nil
}

//---------------------------------------------------	
func Image__check_image_format(pFormatStr string, pRuntimeSys *gf_core.Runtime_sys) (string, bool) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_utils.Image__check_image_format()")
	
	if pFormatStr != "jpeg" && 
		pFormatStr != "jpg" &&
		pFormatStr != "gif" &&
		pFormatStr != "png" {

		// IMPORTANT!! - format is not a valid image format, so 'false' is returned
		return pFormatStr, false
	}

	var normalizedFormatStr string

	// normalize "jpg" (variation on jpeg) to "jpeg"
	if pFormatStr == "jpg" {
		normalizedFormatStr = "jpeg"
	} else {
		normalizedFormatStr = pFormatStr
	}

	return normalizedFormatStr, true
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