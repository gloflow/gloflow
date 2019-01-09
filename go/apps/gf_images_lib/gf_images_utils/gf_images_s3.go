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
	"path"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
)
//---------------------------------------------------
func S3__store_gf_image(p_image_local_file_path_str string,
	p_image_thumbs       *Gf_image_thumbs,
	p_s3_bucket_name_str string,
	p_s3_info            *gf_core.Gf_s3_info,
	p_runtime_sys        *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_s3.S3__store_gf_image()")

	//--------------------
	//UPLOAD FULL_SIZE (ORIGINAL) IMAGE


	
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//FIX!! - target filename of the original image should not be its original file name (that might collide accross domains or other images),
	//        and instead should be the image ID with the file extension. 
	//        it also makes it more difficult to find the image on S3 that is represented by an Gf_img given 
	//        only the ID of that Gf_img
	s3_file_name_str := path.Base(p_image_local_file_path_str)
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!


	/*for files acquired by the Fetcher images are already uploaded 
	with their Gf_img ID as their filename. so here the p_image_local_file_path_str value is already 
	the image ID.
	
	
	ADD!! - have an explicit p_target_s3_file_name_str argument, and dont derive it
	        automatically from the the filename in p_image_local_file_path_str*/



	s3_response_str,gf_err := gf_core.S3__upload_file(p_image_local_file_path_str, //p_target_file__local_path_str
										s3_file_name_str,                          //p_target_file__s3_path_str
										p_s3_bucket_name_str,
										p_s3_info,
										p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	p_runtime_sys.Log_fun("INFO","s3_response_str - "+s3_response_str)
	//--------------------
	//UPLOAD THUMBS

	gf_err = S3__store_gf_image_thumbs(p_image_thumbs,
						p_s3_bucket_name_str,
						p_s3_info,
						p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//--------------------
	return nil
}
//---------------------------------------------------
func S3__store_gf_image_thumbs(p_image_thumbs *Gf_image_thumbs,
	p_s3_bucket_name_str string,
	p_s3_info            *gf_core.Gf_s3_info,
	p_runtime_sys        *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_s3.S3__store_gf_image_thumbs()")

	//IMPORTANT - for some image types (GIF) the system doesnt produce thumbs,
	//            and therefore p_image_thumbs is nil.
	if p_image_thumbs != nil {
		//--------------------
		//SMALL THUMB
		small_t_path_str         := p_image_thumbs.Small_local_file_path_str //thumbs_info_map["small__target_thumbnail_file_path_str"]
		small_t_s3_file_name_str := fmt.Sprintf("/thumbnails/%s",path.Base(small_t_path_str))

		s3_response_str,gf_err := gf_core.S3__upload_file(small_t_path_str, //p_target_file__local_path_str
											small_t_s3_file_name_str,       //p_target_file__s3_path_str
											p_s3_bucket_name_str,
											p_s3_info,
											p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
		p_runtime_sys.Log_fun("INFO","s3_response_str - "+s3_response_str)
		//--------------------
		//MEDIUM THUMB
		medium_t_path_str         := p_image_thumbs.Medium_local_file_path_str //thumbs_info_map["medium__target_thumbnail_file_path_str"]
		medium_t_s3_file_name_str := fmt.Sprintf("/thumbnails/%s",path.Base(medium_t_path_str))

		s3_response_str,gf_err = gf_core.S3__upload_file(medium_t_path_str, //p_target_file__local_path_str
												medium_t_s3_file_name_str,  //p_target_file__s3_path_str
												p_s3_bucket_name_str,
												p_s3_info,
												p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
		p_runtime_sys.Log_fun("INFO","s3_response_str - "+s3_response_str)
		//--------------------
		//LARGE THUMB
		large_t_path_str         := p_image_thumbs.Large_local_file_path_str //thumbs_info_map["large__target_thumbnail_file_path_str"]
		large_t_s3_file_name_str := fmt.Sprintf("/thumbnails/%s",path.Base(large_t_path_str))
		s3_response_str,gf_err    = gf_core.S3__upload_file(large_t_path_str, //p_target_file__local_path_str
												large_t_s3_file_name_str,     //p_target_file__s3_path_str
												p_s3_bucket_name_str,
												p_s3_info,
												p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
		p_runtime_sys.Log_fun("INFO","s3_response_str - "+s3_response_str)
		//--------------------
	}

	return nil
}
//------------------------------------------------
func S3__get_image_url(p_image_path_name_str string,
	p_s3_bucket_name_str string,
	p_runtime_sys        *gf_core.Runtime_sys) string {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_s3.S3__get_image_url()")

	//IMPORTANT!! - amazon URL escapes image file names when it makes them public in a bucket
	//              escaped_str := url.QueryEscape(*p_image_path_name_str)
	url_str := fmt.Sprintf("http://%s.s3-website-us-east-1.amazonaws.com/%s",
						p_s3_bucket_name_str,
						p_image_path_name_str) //escaped_str)
	return url_str
}

//---------------------------------------------------
//IMPORTANT!! - get the filepath of the gf_image's original 
//              image downloaded from the source site (or uploaded by user)

func S3__get_image_original_file_s3_filepath(p_image *Gf_image, p_runtime_sys *gf_core.Runtime_sys) string {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.S3__get_image_original_file_s3_filepath()")

	//when image is downloaded its renamed to its ID
	downloaded_image_filename_str := fmt.Sprintf("%s.%s",p_image.Id_str,p_image.Format_str)
	uploaded_s3_filepath_str      := downloaded_image_filename_str

	return uploaded_s3_filepath_str
}
//---------------------------------------------------
func S3__get_image_thumbs_s3_filepaths(p_image *Gf_image, p_runtime_sys *gf_core.Runtime_sys) (string,string,string) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_utils.S3__get_image_thumbs_s3_filepaths()")
	
	thumb_small__s3_filepath_str  := strings.Replace(p_image.Thumbnail_small_url_str, "/images/d","",1)
	thumb_medium__s3_filepath_str := strings.Replace(p_image.Thumbnail_medium_url_str,"/images/d","",1)
	thumb_large__s3_filepath_str  := strings.Replace(p_image.Thumbnail_large_url_str, "/images/d","",1)

	return thumb_small__s3_filepath_str,thumb_medium__s3_filepath_str,thumb_large__s3_filepath_str
}