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
	"github.com/nfnt/resize"
	"github.com/h2non/bimg"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// p_image_origin_page_url_str - urls of pages (html or some other resource) where the image image_url
//                               was found. this is valid for gf_chrome_ext image sources.
//                               its not relevant for direct image uploads from clients.

func TransformImage(p_image_id_str GFimageID,
	p_image_client_type_str               string,
	p_images_flows_names_lst              []string,
	p_image_origin_url_str                string,
	p_image_origin_page_url_str           string,
	p_meta_map                            map[string]interface{},
	pImageLocalFilePathStr                string,
	pImagesStoreThumbnailsLocalDirPathStr string,
	pCtx                                  context.Context,
	pRuntimeSys                           *gf_core.RuntimeSys) (*GFimage, *GF_image_thumbs, *gf_core.GFerror) {

	normalizedExtStr, gfErr := GetImageExtFromURL(pImageLocalFilePathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	gfImage, gfImageThumbs, gfErr := TransformProcessImage(p_image_id_str,
		p_image_client_type_str,
		p_images_flows_names_lst,
		p_image_origin_url_str,
		p_image_origin_page_url_str,
		p_meta_map,
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
	pRuntimeSys                           *gf_core.RuntimeSys) (*GF_image, *GF_image_thumbs, *gf_core.GFerror) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_transformer.Trans__process_image()")

	//---------------------------------
	// LOAD_IMAGE

	img, gfErr := ImageLoadFile(pImageLocalFilePathStr, pNormalizedExtStr, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}
	
	//--------------------------
	// CREATE THUMBNAILS

	small_thumb_max_size_px_int  := 200
	medium_thumb_max_size_px_int := 400
	large_thumb_max_size_px_int  := 600

	gf_image_thumbs, gfErr := CreateThumbnails(pImageIDstr,
		pNormalizedExtStr,
		pImageLocalFilePathStr,
		pImagesStoreThumbnailsLocalDirPathStr,
		small_thumb_max_size_px_int,
		medium_thumb_max_size_px_int,
		large_thumb_max_size_px_int,
		img,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//--------------------------
	/* //DOMINANT COLOR DETERMINATION
	//it"s computed only for non-gif"s
	dominant_color_hex_str := gf_images_core_graphic.get_dominant_image_color(pImageLocalFilePathStr,p_log_fun)*/

	//--------------------------
	image_width_int, image_height_int := Get_image_dimensions__from_image(img, pRuntimeSys)

	//--------------------------

	// SECURITY ISSUE!!
	// When you open a file, the file header is read to determine the file 
	// format and extract things like mode, size, and other properties 
	// required to decode the file, but the rest of the file is not processed until later.
	
	// someone can forge header information in an image
	imageTitleStr, gfErr := Get_image_title_from_url(pImageOriginURLstr, pRuntimeSys)
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
		Thumbnail_small_url_str:        gf_image_thumbs.Small_relative_url_str,
		Thumbnail_medium_url_str:       gf_image_thumbs.Medium_relative_url_str,
		Thumbnail_large_url_str:        gf_image_thumbs.Large_relative_url_str,
		Format_str:                     pNormalizedExtStr,
		Width_int:                      image_width_int,
		Height_int:                     image_height_int,

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

	return gfImage, gf_image_thumbs, nil
}

//---------------------------------------------------
func resizeImage(p_img image.Image,
	p_image_output_path_str string,
	p_image_format_str      string,
	p_size_px_int           int,
	pRuntimeSys           *gf_core.Runtime_sys) *gf_core.Gf_error {
	
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	
	m := resize.Resize(uint(p_size_px_int), 0, p_img, resize.Bilinear) // resize.Lanczos3)

	out, err := os.Create(p_image_output_path_str)
	if err != nil {
		gf_err := gf_core.Error__create("OS failed to create a file to save a resized image to FS",
			"file_create_error",
			map[string]interface{}{"image_output_path_str": p_image_output_path_str,},
			err, "gf_images_core", pRuntimeSys)
		return gf_err
	}
	defer out.Close()


	//--------------------------
	// ADD!! - enable a way to set PNG format as a target.
	//         for very precise vector type images there is significant accuracy loss
	//         when recoding using JPEG
	
	// IMPORTANT!! - using JPEG instead of PNG, because JPEG compression was made for photographic images,
	//               and so for these kinds of images it comes out with much smaller file size
	// write new image to file
	// jpeg.Encode(out, m, nil)
	jpeg.Encode(out, m, nil)

	//--------------------------
	return nil
}