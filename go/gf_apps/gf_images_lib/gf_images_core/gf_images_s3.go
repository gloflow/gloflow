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
	"path"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
)

//---------------------------------------------------

func S3getImage(pImageS3filePathStr string,
	pTargetFileLocalPathStr string,
	pS3bucketNameStr string,
	pS3info          *gf_aws.GFs3Info,
	pRuntimeSys      *gf_core.RuntimeSys) *gf_core.GFerror {

	gfErr := gf_aws.S3getFile(pImageS3filePathStr,
		pTargetFileLocalPathStr,
		pS3bucketNameStr,
		pS3info,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// DEPRECATED!!

func S3storeImage(pImageLocalFilePathStr string,
	pImageThumbs     *GFimageThumbs,
	pS3bucketNameStr string,
	pS3info          *gf_aws.GFs3Info,
	pRuntimeSys      *gf_core.RuntimeSys) *gf_core.GFerror {

	//--------------------
	// UPLOAD FULL_SIZE (ORIGINAL) IMAGE

	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// FIX!! - target filename of the original image should not be its original file name (that might collide accross domains or other images),
	//         and instead should be the image ID with the file extension. 
	//         it also makes it more difficult to find the image on S3 that is represented by an Gf_img given 
	//         only the ID of that Gf_img
	s3FileNameStr := path.Base(pImageLocalFilePathStr)
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	/* for files acquired by the Fetcher images are already uploaded 
	with their Gf_img ID as their filename. so here the pImageLocalFilePathStr value is already 
	the image ID.
	
	ADD!! - have an explicit p_target_s3FileNameStr argument, and dont derive it
	        automatically from the the filename in pImageLocalFilePathStr */

	targetFileLocalPathStr := pImageLocalFilePathStr
	targetFileS3pathStr    := s3FileNameStr
	_, gfErr := gf_aws.S3putFile(targetFileLocalPathStr,
		targetFileS3pathStr,
		pS3bucketNameStr,
		pS3info,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}
	
	//--------------------
	// UPLOAD THUMBS

	gfErr = S3storeThumbnails(pImageThumbs,
		pS3bucketNameStr,
		pS3info,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//--------------------
	return nil
}

//---------------------------------------------------
// DEPRECATED!!

func S3storeThumbnails(pImageThumbs *GFimageThumbs,
	pS3bucketNameStr string,
	pS3info          *gf_aws.GFs3Info,
	pRuntimeSys      *gf_core.RuntimeSys) *gf_core.GFerror {

	// IMPORTANT - for some image types (GIF) the system doesnt produce thumbs,
	//             and therefore pImageThumbs is nil.
	if pImageThumbs != nil {

		//--------------------
		// SMALL THUMB
		small_t_path_str         := pImageThumbs.Small_local_file_path_str // thumbs_info_map["small__target_thumbnail_file_path_str"]
		small_t_s3_file_name_str := fmt.Sprintf("/thumbnails/%s", path.Base(small_t_path_str))
		s3ResponseStr, gfErr  := gf_aws.S3putFile(small_t_path_str, small_t_s3_file_name_str, pS3bucketNameStr, pS3info, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
		pRuntimeSys.LogFun("INFO","s3ResponseStr - "+s3ResponseStr)

		//--------------------
		// MEDIUM THUMB
		medium_t_path_str         := pImageThumbs.Medium_local_file_path_str // thumbs_info_map["medium__target_thumbnail_file_path_str"]
		medium_t_s3_file_name_str := fmt.Sprintf("/thumbnails/%s", path.Base(medium_t_path_str))

		s3ResponseStr, gfErr = gf_aws.S3putFile(medium_t_path_str, medium_t_s3_file_name_str, pS3bucketNameStr, pS3info, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
		pRuntimeSys.LogFun("INFO","s3ResponseStr - "+s3ResponseStr)

		//--------------------
		// LARGE THUMB
		large_t_path_str         := pImageThumbs.Large_local_file_path_str // thumbs_info_map["large__target_thumbnail_file_path_str"]
		large_t_s3_file_name_str := fmt.Sprintf("/thumbnails/%s",path.Base(large_t_path_str))
		s3ResponseStr, gfErr   = gf_aws.S3putFile(large_t_path_str, large_t_s3_file_name_str, pS3bucketNameStr, pS3info, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
		pRuntimeSys.LogFun("INFO","s3ResponseStr - "+s3ResponseStr)

		//--------------------
	}
	return nil
}

//---------------------------------------------------
// S3__get_image_original_file_s3_filepath returns the S3 filepath of a gf_image's original image.
// Original image is the full-size file that was initially acquired, whether fetched from an external source
// or uploaded via API by other programs or by users via UI's).
// As input it requires a Gf_image struct.

func S3getImageOriginalFileS3filepath(p_image *GFimage,
	pRuntimeSys *gf_core.RuntimeSys) string {
	
	// when image is downloaded its renamed to its ID
	downloaded_image_filename_str := fmt.Sprintf("%s.%s", p_image.IDstr, p_image.Format_str)
	s3_filepath_str               := downloaded_image_filename_str

	return s3_filepath_str
}

//---------------------------------------------------

func S3getImageFilepath(pImageIDstr GFimageID,
	pImageFormatStr string,
	pRuntimeSys     *gf_core.RuntimeSys) string {

	imageFileNameStr := ImageGetFilepathFromID(pImageIDstr, pImageFormatStr) // fmt.Sprintf("%s.%s", p_gf_image_id_str, pImageFormatStr)
	s3FilepathStr    := imageFileNameStr
	return s3FilepathStr
}

//---------------------------------------------------

func S3getImageThumbsS3filepaths(p_image *GFimage,
	pRuntimeSys *gf_core.RuntimeSys) (string, string, string) {
	
	thumb_small__s3_filepath_str  := strings.Replace(p_image.Thumbnail_small_url_str,  "/images/d", "", 1)
	thumb_medium__s3_filepath_str := strings.Replace(p_image.Thumbnail_medium_url_str, "/images/d", "", 1)
	thumb_large__s3_filepath_str  := strings.Replace(p_image.Thumbnail_large_url_str,  "/images/d", "", 1)

	return thumb_small__s3_filepath_str, thumb_medium__s3_filepath_str, thumb_large__s3_filepath_str
}