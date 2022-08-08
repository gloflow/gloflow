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
	"context"
	"image"
	"image/jpeg"
	// "image/png"
	"github.com/nfnt/resize"
	"github.com/h2non/bimg"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_plugins"
)

//-------------------------------------------------
// p_image_origin_page_url_str - urls of pages (html or some other resource) where the image image_url
//                               was found. this is valid for gf_chrome_ext image sources.
//                               its not relevant for direct image uploads from clients.

func TransformImage(pImageIDstr GFimageID,
	pImageClientTypeStr                   string,
	pImagesFlowsNamesLst                  []string,
	pImageOriginURLstr                    string,
	pImageOriginPageURLstr                string,
	pMetaMap                              map[string]interface{},
	pImageLocalFilePathStr                string,
	pImagesStoreThumbnailsLocalDirPathStr string,
	pCtx                                  context.Context,
	pRuntimeSys                           *gf_core.RuntimeSys) (*GFimage, *GFimageThumbs, *gf_core.GFerror) {

	normalizedExtStr, gfErr := GetImageExtFromURL(pImageLocalFilePathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	gfImage, gfImageThumbs, gfErr := TransformProcessImage(pImageIDstr,
		pImageClientTypeStr,
		pImagesFlowsNamesLst,
		pImageOriginURLstr,
		pImageOriginPageURLstr,
		pMetaMap,
		normalizedExtStr,
		pImageLocalFilePathStr,
		pImagesStoreThumbnailsLocalDirPathStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}


	//--------------------------
	// FINISH!! - this processing function uses "bimg", which uses
	//            and underlying C lib "libvips"
	gfErr = TransformProcessImageV2(pImageLocalFilePathStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//--------------------------
	// PY_PLUGINS
	// runs the Py VM in a new process, spawned via a new go-routine

	pluginsPyDirPathStr := "./gf_images/plugins"
	gf_images_plugins.RunPyImagePlugins(pImageLocalFilePathStr,
		pluginsPyDirPathStr,
		pCtx,
		pRuntimeSys)

	//--------------------------

	return gfImage, gfImageThumbs, nil
}

//---------------------------------------------------
// V2
func TransformProcessImageV2(pImageLocalFilePathStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	buffer, err := bimg.Read(pImageLocalFilePathStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newImage, err := bimg.NewImage(buffer).Convert(bimg.PNG)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if bimg.NewImage(newImage).Type() == "png" {
		fmt.Fprintln(os.Stderr, "The image was converted into png")
	}

	return nil
}

//---------------------------------------------------
func TransformProcessImage(pImageIDstr GFimageID,
	pImageClientTypeStr                string,
	pImagesFlowsNamesLst               []string,
	pImageOriginURLstr                 string,
	pImageOriginPageURLstr             string,
	pMetaMap                           map[string]interface{},
	pNormalizedExtStr                  string,
	pImageLocalFilePathStr             string,
	pImagesStoreThumbnailsLocalDirPathStr string,
	pCtx                                  context.Context,
	pRuntimeSys                           *gf_core.RuntimeSys) (*GF_image, *GFimageThumbs, *gf_core.GFerror) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_transformer.Trans__process_image()")

	//---------------------------------
	// LOAD_IMAGE

	img, gfErr := ImageLoadFile(pImageLocalFilePathStr, pNormalizedExtStr, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//---------------------------------
	// DIMENSIONS
	imageWidthInt, imageHeightInt := GetImageDimensionsFromImage(img, pRuntimeSys)

	//--------------------------
	// CREATE THUMBNAILS

	small_thumb_max_size_px_int  := 200
	medium_thumb_max_size_px_int := 400
	large_thumb_max_size_px_int  := 800

	gfImageThumbs, gfErr := CreateThumbnails(img,
		pImageIDstr,
		pNormalizedExtStr,
		pImagesStoreThumbnailsLocalDirPathStr,
		small_thumb_max_size_px_int,
		medium_thumb_max_size_px_int,
		large_thumb_max_size_px_int,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//--------------------------
	/* //DOMINANT COLOR DETERMINATION
	//it"s computed only for non-gif"s
	dominant_color_hex_str := gf_images_core_graphic.get_dominant_image_color(pImageLocalFilePathStr,p_log_fun)*/

	//--------------------------
	

	//--------------------------

	// SECURITY ISSUE!!
	// When you open a file, the file header is read to determine the file 
	// format and extract things like mode, size, and other properties 
	// required to decode the file, but the rest of the file is not processed until later.
	
	// someone can forge header information in an image
	imageTitleStr, gfErr := GetImageTitleFromURL(pImageOriginURLstr, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	imageInfo := &GFimageNewInfo{
		Id_str:                         pImageIDstr,
		Title_str:                      imageTitleStr,
		Flows_names_lst:                pImagesFlowsNamesLst,
		Image_client_type_str:          pImageClientTypeStr,
		Origin_url_str:                 pImageOriginURLstr,
		Origin_page_url_str:            pImageOriginPageURLstr,
		Original_file_internal_uri_str: pImageLocalFilePathStr,
		Thumbnail_small_url_str:        gfImageThumbs.Small_relative_url_str,
		Thumbnail_medium_url_str:       gfImageThumbs.Medium_relative_url_str,
		Thumbnail_large_url_str:        gfImageThumbs.Large_relative_url_str,
		Format_str:                     pNormalizedExtStr,
		Width_int:                      imageWidthInt,
		Height_int:                     imageHeightInt,

		Meta_map: pMetaMap,
	}

	//--------------------------
	// IMAGE_CREATE

	// IMPORTANT!! - creates a GF_Image struct and stores it in the DB
	gfImage, gfErr := ImageCreateNew(imageInfo, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//--------------------------

	return gfImage, gfImageThumbs, nil
}

//---------------------------------------------------
// OPS
//---------------------------------------------------
func resizeImage(pImg image.Image,
	pImageOutputPathStr string,
	pWidthPxInt         int,
	pHeightPxInt        int,
	pRuntimeSys         *gf_core.RuntimeSys) *gf_core.GFerror {
	
	m := resize.Resize(uint(pWidthPxInt), uint(pHeightPxInt), pImg, resize.Lanczos3)

	out, err := os.Create(pImageOutputPathStr)
	if err != nil {
		gfErr := gf_core.Error__create("OS failed to create a file to save a resized image to FS",
			"file_create_error",
			map[string]interface{}{"image_output_path_str": pImageOutputPathStr,},
			err, "gf_images_core", pRuntimeSys)
		return gfErr
	}
	defer out.Close()

	/*out_png, err := os.Create(fmt.Sprintf("%s.png", pImageOutputPathStr))
	if err != nil {
		gfErr := gf_core.Error__create("OS failed to create a file to save a resized image to FS",
			"file_create_error",
			map[string]interface{}{"image_output_path_str": pImageOutputPathStr,},
			err, "gf_images_core", pRuntimeSys)
		return gfErr
	}
	defer out.Close()*/

	//--------------------------
	// ADD!! - enable a way to set PNG format as a target.
	//         for very precise vector type images there is significant accuracy loss
	//         when recoding using JPEG
	// png.Encode(out_png, m)

	// IMPORTANT!! - using JPEG instead of PNG, because JPEG compression was made for photographic images,
	//               and so for these kinds of images it comes out with much smaller file size.	
	jpeg.Encode(out, m, nil)
	
	//--------------------------
	return nil
}