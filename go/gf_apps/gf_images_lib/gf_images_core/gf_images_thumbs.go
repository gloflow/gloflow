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
func CreateThumbnails(pImage image.Image,
	pImageIDstr                      GFimageID,
	pImageFormatStr                  string,
	pTargetThumbnailsLocalDirPathStr string,
	pSmallThumbSizePxInt             int,
	pMediumThumbSizePxInt            int,
	pLargeThumbSizePxInt             int,
	pRuntimeSys                      *gf_core.RuntimeSys) (*GFimageThumbs, *gf_core.GFerror) {

	// determine if the width or height is larger, and get its value
	largerDimensionPxInt, largerDimensionNameStr := GetImageLargerDimension(pImage, pRuntimeSys)

	//-----------------
	// SMALL THUMBS
	new_thumb_small_file_name_str         := fmt.Sprintf("%s_thumb_small.%s", pImageIDstr, pImageFormatStr)
	small__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", pTargetThumbnailsLocalDirPathStr, new_thumb_small_file_name_str)
	
	smallThumbWidthPxInt, smallThumbHeightPxInt := ThumbsGetSizeInPx(pSmallThumbSizePxInt,
		largerDimensionPxInt,
		largerDimensionNameStr)
	
	gfErr := resizeImage(pImage,
		small__target_thumbnail_file_path_str,
		smallThumbWidthPxInt,
		smallThumbHeightPxInt,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-----------------
	// MEDIUM THUMBS
	new_thumb_medium_file_name_str         := fmt.Sprintf("%s_thumb_medium.%s", pImageIDstr, pImageFormatStr)
	medium__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", pTargetThumbnailsLocalDirPathStr, new_thumb_medium_file_name_str)

	mediumThumbWidthPxInt, mediumThumbHeightPxInt := ThumbsGetSizeInPx(pMediumThumbSizePxInt, largerDimensionPxInt, largerDimensionNameStr)

	gfErr = resizeImage(pImage,
		medium__target_thumbnail_file_path_str,
		mediumThumbWidthPxInt,
		mediumThumbHeightPxInt,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-----------------
	// LARGE THUMBS
	new_thumb_large_file_name_str         := fmt.Sprintf("%s_thumb_large.%s", pImageIDstr, pImageFormatStr)
	large__target_thumbnail_file_path_str := fmt.Sprintf("%s/%s", pTargetThumbnailsLocalDirPathStr, new_thumb_large_file_name_str)

	largerThumbWidthPxInt, largerThumbHeightPxInt := ThumbsGetSizeInPx(pLargeThumbSizePxInt, largerDimensionPxInt, largerDimensionNameStr)

	gfErr = resizeImage(pImage,
		large__target_thumbnail_file_path_str,
		largerThumbWidthPxInt,
		largerThumbHeightPxInt,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-----------------

	thumb_small_relative_url_str  := fmt.Sprintf("/images/d/thumbnails/%s", new_thumb_small_file_name_str)
	thumb_medium_relative_url_str := fmt.Sprintf("/images/d/thumbnails/%s", new_thumb_medium_file_name_str)
	thumb_large_relative_url_str  := fmt.Sprintf("/images/d/thumbnails/%s", new_thumb_large_file_name_str)

	imageThumbs := &GFimageThumbs{
		Small_relative_url_str:     thumb_small_relative_url_str,
		Medium_relative_url_str:    thumb_medium_relative_url_str,
		Large_relative_url_str:     thumb_large_relative_url_str,
		Small_local_file_path_str:  small__target_thumbnail_file_path_str,
		Medium_local_file_path_str: medium__target_thumbnail_file_path_str,
		Large_local_file_path_str:  large__target_thumbnail_file_path_str,
	}

	return imageThumbs, nil
}

//---------------------------------------------------
func StoreThumbnails(pImageThumbs *GFimageThumbs,
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

//---------------------------------------------------
func ThumbsGetSizeInPx(pThumbSizeInPxInt int,
	pOriginalLargerDimensionInPxInt int,
	pOriginalLargerDimensionNameStr string) (int, int) {

		


	var thumbDimensionSizeInPxInt int
	if pOriginalLargerDimensionInPxInt > pThumbSizeInPxInt {

		// if the original image is larger than the desired
		// thumb size then use that desired thumb size as its
		// final size to use.
		thumbDimensionSizeInPxInt = pThumbSizeInPxInt
	} else {

		// if the original image is smaller than the desired
		// thumb size then dont upscale the image to fit with the 
		// desired thumbs size, and instead use the original size
		// as the final size to use.
		thumbDimensionSizeInPxInt = pOriginalLargerDimensionInPxInt	
	}

	var widthInt int
	var heightInt int 
	if pOriginalLargerDimensionNameStr == "width" {
		widthInt  = thumbDimensionSizeInPxInt

		// setting the dimension to value of 0 causes the resizer
		// to maintain the original aspect-ratio
		heightInt = 0

	} else {
		widthInt  = 0
		heightInt = thumbDimensionSizeInPxInt
	}

	return widthInt, heightInt
}