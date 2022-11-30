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

package gf_crawl_core

import (
	"fmt"
	"context"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
	// "github.com/davecgh/go-spew/spew"
)

//--------------------------------------------------

func images__stage__process_images(pCrawlerNameStr string,
	p_page_imgs__pipeline_infos_lst   []*gf_page_img__pipeline_info,
	p_images_store_local_dir_path_str string,
	p_origin_page_url_str             string,
	pMediaDomainStr                   string,
	pS3bucketNameStr                  string,
	pRuntime                          *GFcrawlerRuntime,
	pRuntimeSys                       *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - process_images")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _, page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		// IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		// IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if page_img__pinfo.exists_bool {
			continue
		}

		// IMPORTANT!! - check image is not flagged as a NSFV image
		if page_img__pinfo.nsfv_bool {
			continue
		}

		//----------------------------
		// IMAGE_PROCESS
		_, gf_image_thumbs, gfErr := imageProcess(page_img__pinfo.page_img,
			page_img__pinfo.gf_image_id_str, //pGFimageIDstr
			page_img__pinfo.local_file_path_str,
			p_images_store_local_dir_path_str,

			pMediaDomainStr,
			pS3bucketNameStr,
			pRuntime,
			pRuntimeSys)
		//----------------------------
		
		if gfErr != nil {
			t := "image_process__failed"
			m := "failed processing of image with img_url_str - "+page_img__pinfo.page_img.Url_str
			CreateErrorAndEvent(t,m,map[string]interface{}{"origin_page_url_str": p_origin_page_url_str,}, page_img__pinfo.page_img.Url_str, pCrawlerNameStr,
				gfErr, pRuntime, pRuntimeSys)

			page_img__pinfo.gf_error = gfErr
			continue // IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}

		// UPDATE__PAGE_IMG_PINFO
		page_img__pinfo.thumbs = gf_image_thumbs
	}
	return p_page_imgs__pipeline_infos_lst
}

//--------------------------------------------------

func imageProcess(pPageImg *GFcrawlerPageImage,
	pGFimageIDstr                     gf_images_core.GFimageID,
	p_local_image_file_path_str       string,
	p_images_store_local_dir_path_str string,
	pMediaDomainStr                   string,
	pS3bucketNameStr                  string,
	pRuntime                          *GFcrawlerRuntime,
	pRuntimeSys                       *gf_core.RuntimeSys) (*gf_images_core.GFimage, *gf_images_core.GFimageThumbs, *gf_core.GFerror) {

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	pRuntimeSys.LogFun("INFO", cyan("       >>>>>>>>>>>>> ----------------------------- ")+yellow("PROCESS_IMAGE"))

	//----------------------------
	// GIF
	if pPageImg.Img_ext_str == "gif" {

		image_client_type_str := "gf_crawl_images" 
		image_flows_names_lst := []string{"discovered", "gifs",}

		gif_download_and_frames__local_dir_path_str := p_images_store_local_dir_path_str

		ctx := context.Background()

		gf_gif, _, gfErr := gf_gif_lib.Process(pGFimageIDstr,
			pPageImg.Url_str,
			pPageImg.Origin_page_url_str,
			gif_download_and_frames__local_dir_path_str,
			image_client_type_str,
			image_flows_names_lst,
			true, // p_create_new_db_img_bool

			pMediaDomainStr,
			pS3bucketNameStr,
			pRuntime.S3info,
			ctx,
			pRuntimeSys)

		if gfErr != nil {
			return nil, nil, gfErr
		}													

		imageIDstr := gf_gif.GFimageIDstr
		gfErr       = image__db_update_after_process(pPageImg, imageIDstr, pRuntimeSys)
		if gfErr != nil {
			return nil, nil, gfErr
		}

		return nil, nil, nil

	//----------------------------
	// GENERAL
	} else {
	
		thumbnailsLocalDirPathStr := p_images_store_local_dir_path_str
		gfImage, gfImageThumbs, gfErr := imageProcessBitmap(pPageImg,
			pGFimageIDstr,
			p_local_image_file_path_str,
			thumbnailsLocalDirPathStr,
			pRuntime.PluginsPyDirPathStr,
			pRuntimeSys)

		if gfErr != nil {
			return nil, nil, gfErr
		}

		// IMPORTANT!! - if gf_image is nil and there is no error then imageProcessBitmap()
		//               determined that the image is in some way invalid and should not be further processesd
		//               (currently its nil if the image is smaller then the allowed dimension - the 
		//               image is some small icon or banner/etc.)
		if gfImage == nil {
			return nil, nil, nil
		}
		
		//spew.Dump(gf_image)

		imageIDstr := gfImage.IDstr
		gfErr       = image__db_update_after_process(pPageImg, imageIDstr, pRuntimeSys)
		if gfErr != nil {
			return nil, nil, gfErr
		}

		return gfImage, gfImageThumbs, nil
	}
	
	//----------------------------
	return nil, nil, nil
}

//--------------------------------------------------

func imageProcessBitmap(p_page_img *GFcrawlerPageImage,
	pImageIDstr                     gf_images_core.GFimageID,
	p_local_image_file_path_str     string,
	p_thumbnails_local_dir_path_str string,
	pPluginsPyDirPathStr            string,
	pRuntimeSys                     *gf_core.RuntimeSys) (*gf_images_core.GFimage, *gf_images_core.GFimageThumbs, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_process.image__process_bitmap()")

	//----------------------
	// CONFIG
	image_client_type_str := "gf_crawl_images" 
	image_flows_names_lst := []string{"discovered",}

	//----------------------

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//-------------------
	imgWidthInt, imgHeightInt, gfErr := gf_images_core.GetImageDimensionsFromFilepath(p_local_image_file_path_str, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//-------------------

	// IMPORTANT!! - check that the image is too small, and is likely to be irrelevant 
	//               part of a particular page
	if imgWidthInt <= 130 || imgHeightInt <= 130 {
		pRuntimeSys.LogFun("INFO", yellow("IMG IS SMALLER THEN MINIMUM DIMENSIONS (width-"+cyan(fmt.Sprint(imgWidthInt))+"/height-"+cyan(fmt.Sprint(imgHeightInt))+")"))
		return nil, nil, nil
	} else {

		//--------------------------------
		// TRANSFORM DOWNLOADED IMAGE - CREATE THUMBS, SAVE TO DB, AND UPLOAD TO AWS_S3

		// IMPORTANT!! - a new gf_image ID is created if an external ID is not supplied
		var imageIDstr gf_images_core.GFimageID
		if pImageIDstr == "" {
			newImageIDstr, gfErr := gf_images_core.CreateIDfromURL(p_page_img.Url_str, pRuntimeSys)
			if gfErr != nil {
				return nil, nil, gfErr
			}
			imageIDstr = newImageIDstr
		} else {
			imageIDstr = pImageIDstr
		}

		imageOriginURLstr     := p_page_img.Url_str
		imageOriginPageURLstr := p_page_img.Origin_page_url_str
		meta_map := map[string]interface{}{}

		ctx := context.Background()


		// FINISH!! - properly create an instance of GFmetrics
		var imagesCoreMetrics *gf_images_core.GFmetrics

		// IMPORTANT!! - this creates a Gf_image object, and persists it in the DB ("t" == "img"),
		//               also creates gf_image thumbnails as local files.
		gf_image, gf_image_thumbs, gfErr := gf_images_core.TransformImage(imageIDstr,
			image_client_type_str,
			image_flows_names_lst,
			imageOriginURLstr,
			imageOriginPageURLstr,
			meta_map,
			p_local_image_file_path_str,
			p_thumbnails_local_dir_path_str,
			pPluginsPyDirPathStr,
			imagesCoreMetrics,
			ctx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, nil, gfErr
		}
		//--------------------------------

		return gf_image, gf_image_thumbs, nil
	}

	return nil, nil, nil
}