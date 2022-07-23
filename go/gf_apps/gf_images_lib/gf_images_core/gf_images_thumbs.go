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
	"image"
	"path"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_storage"
)

//---------------------------------------------------
func CreateThumbnails(p_image_id_str GFimageID,
	p_image_format_str                     string,
	p_image_file_path_str                  string,
	p_local_target_thumbnails_dir_path_str string,
	p_small_thumb_max_size_px_int          int,
	p_medium_thumb_max_size_px_int         int,
	p_large_thumb_max_size_px_int          int,
	p_image                                image.Image,
	pRuntimeSys                            *gf_core.RuntimeSys) (*GF_image_thumbs, *gf_core.GFerror) {

	//-----------------
	// SMALL THUMBS
	new_thumb_small_file_name_str         := fmt.Sprintf("%s_thumb_small.%s", p_image_id_str, p_image_format_str)
	small__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", p_local_target_thumbnails_dir_path_str, new_thumb_small_file_name_str)

	gfErr := resizeImage(p_image, // p_image_file,
		small__target_thumbnail_file_path_str,
		p_image_format_str,
		p_small_thumb_max_size_px_int,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-----------------
	// MEDIUM THUMBS
	new_thumb_medium_file_name_str         := fmt.Sprintf("%s_thumb_medium.%s", p_image_id_str, p_image_format_str)
	medium__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", p_local_target_thumbnails_dir_path_str, new_thumb_medium_file_name_str)

	gfErr = resizeImage(p_image, // p_image_file,
		medium__target_thumbnail_file_path_str,
		p_image_format_str,
		p_medium_thumb_max_size_px_int,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-----------------
	// LARGE THUMBS
	new_thumb_large_file_name_str         := fmt.Sprintf("%s_thumb_large.%s", p_image_id_str, p_image_format_str)
	large__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", p_local_target_thumbnails_dir_path_str, new_thumb_large_file_name_str)

	gfErr = resizeImage(p_image, // p_image_file,
		large__target_thumbnail_file_path_str,
		p_image_format_str,
		p_large_thumb_max_size_px_int,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-----------------

	thumb_small_relative_url_str  := "/images/d/thumbnails/"+new_thumb_small_file_name_str
	thumb_medium_relative_url_str := "/images/d/thumbnails/"+new_thumb_medium_file_name_str
	thumb_large_relative_url_str  := "/images/d/thumbnails/"+new_thumb_large_file_name_str

	image_thumbs := &GF_image_thumbs{
		Small_relative_url_str:     thumb_small_relative_url_str,
		Medium_relative_url_str:    thumb_medium_relative_url_str,
		Large_relative_url_str:     thumb_large_relative_url_str,
		Small_local_file_path_str:  small__target_thumbnail_file_path_str,
		Medium_local_file_path_str: medium__target_thumbnail_file_path_str,
		Large_local_file_path_str:  large__target_thumbnail_file_path_str,
	}

	return image_thumbs, nil
}

//---------------------------------------------------
func StoreThumbnails(pImageThumbs *GF_image_thumbs,
	pStorage    *gf_images_storage.GFimageStorage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	// IMPORTANT - for some image types (GIF) the system doesnt produce thumbs,
	//             and therefore pImageThumbs is nil.
	if pImageThumbs != nil {

		//---------------------------------------------------
		filePutFunc := func(pLocalFilePathStr string,
			pTargetFilePathStr string) *gf_core.GFerror {

			op := &gf_images_storage.GFputFromLocalOpDef{
				SourceLocalFilePathStr: pLocalFilePathStr,
				TargetFilePathStr:      pTargetFilePathStr,
			}
			if pStorage.TypeStr == "s3" {
				op.S3bucketNameStr = pStorage.S3.ThumbsS3bucketNameStr
			}
			gfErr := gf_images_storage.FilePutFromLocal(op, pStorage, pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
			return nil
		}

		//---------------------------------------------------

		//--------------------
		// SMALL THUMB
		smallTpathStr           := pImageThumbs.Small_local_file_path_str // thumbs_info_map["small__target_thumbnail_file_path_str"]
		smallTtargetFilePathStr := fmt.Sprintf("/thumbnails/%s", path.Base(smallTpathStr))

		gfErr := filePutFunc(smallTpathStr, smallTtargetFilePathStr)
		if gfErr != nil {
			return gfErr
		}

		//--------------------
		// MEDIUM THUMB
		mediumTpathStr           := pImageThumbs.Medium_local_file_path_str // thumbs_info_map["medium__target_thumbnail_file_path_str"]
		mediumTtargetFilePathStr := fmt.Sprintf("/thumbnails/%s", path.Base(mediumTpathStr))

		gfErr = filePutFunc(mediumTpathStr, mediumTtargetFilePathStr)
		if gfErr != nil {
			return gfErr
		}

		//--------------------
		// LARGE THUMB
		largeTpathStr         := pImageThumbs.Large_local_file_path_str // thumbs_info_map["large__target_thumbnail_file_path_str"]
		largeTtargetFilePathStr := fmt.Sprintf("/thumbnails/%s",path.Base(largeTpathStr))

		gfErr = filePutFunc(largeTpathStr, largeTtargetFilePathStr)
		if gfErr != nil {
			return gfErr
		}

		//--------------------
	}

	return nil
}