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
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//--------------------------------------------------
type gf_page_img__pipeline_info struct {
	link                *gf_page_img_link
	page_img            *Gf_crawler_page_image
	page_img_ref        *Gf_crawler_page_image_ref
	local_file_path_str string
	thumbs              *gf_images_core.GFimageThumbs
	exists_bool         bool                   // has the page_img already been discovered in the past
	nsfv_bool           bool
	gf_error            *gf_core.GFerror       // if page_img processing failed at some stage

	// in some situations (or in tests) we wish to manually assign a gf_image_id,
	// instead of letting the gf_image processing/transformation
	// operations create those ID's themselves
	gf_image_id_str gf_images_core.GFimageID
}

type gf_page_img_link struct {
	img_src_str         string
	origin_page_url_str string
}

//--------------------------------------------------
func images_pipe__from_html(pURLfetch *Gf_crawler_url_fetch,
	pCycleRunIDstr         string,
	pCrawlerNameStr        string,
	pImagesLocalDirPathStr string,
	pMediaDomainStr        string,
	pS3bucketNameStr       string,
	pRuntime               *GFcrawlerRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_pipeline.images_pipe__from_html()")

	cyan := color.New(color.FgCyan).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	//yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))
	fmt.Println(">> IMAGES__GET_IN_PAGE - "+blue(pURLfetch.Url_str))
	fmt.Println(cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))

	origin_page_url_str := pURLfetch.Url_str

	//------------------
	// STAGE - pull all page image links

	page_imgs__pipeline_infos_lst := images__stage__pull_image_links(pURLfetch,
		pCrawlerNameStr,
		pCycleRunIDstr,
		pRuntime,
		pRuntimeSys)

	//------------------
	// STAGE - create gf_image/gf_image_refs structs

	page_imgs__pinfos_with_imgs_lst := images__stage__create_page_images(pCrawlerNameStr,
		pCycleRunIDstr,
		page_imgs__pipeline_infos_lst,
		pRuntime,
		pRuntimeSys)

	//------------------
	// STAGE - persist gf_image/gf_image_ref to DB
	
	page_imgs__pinfos_with_persists_lst := images__stage__page_images_persist(pCrawlerNameStr,
		page_imgs__pinfos_with_imgs_lst,
		pRuntime,
		pRuntimeSys)

	//------------------
	// STAGE - download gf_images from target URL
	
	page_imgs__pinfos_with_local_file_paths_lst := images__stage__download_images(pCrawlerNameStr,
		page_imgs__pinfos_with_persists_lst,
		pImagesLocalDirPathStr,
		origin_page_url_str,
		pRuntime,
		pRuntimeSys)

	//------------------
	// STAGES - process images

	page_imgs__pinfos_with_thumbs_lst := images__stages__process_images(pCrawlerNameStr,
		page_imgs__pinfos_with_local_file_paths_lst,
		pImagesLocalDirPathStr,
		origin_page_url_str,

		pMediaDomainStr,
		pS3bucketNameStr,
		pRuntime,
		pRuntimeSys)

	//------------------
	// STAGE - persist all images files (S3, etc.)

	page_imgs__pinfos_with_s3_lst := imagesS3stageStoreImages(pCrawlerNameStr,
		page_imgs__pinfos_with_thumbs_lst,
		origin_page_url_str,
		pS3bucketNameStr,
		pRuntime,
		pRuntimeSys)

	//------------------
	// STAGE - cleanup

	images__stages_cleanup(page_imgs__pinfos_with_s3_lst, pRuntime, pRuntimeSys)

	//------------------
}

//--------------------------------------------------
// SINGLE_IMAGE

func images_pipe__single_simple(pImage *Gf_crawler_page_image,
	pImagesStoreLocalDirPathStr   string,
	pMediaDomainStr               string,
	pCrawledImagesS3bucketNameStr string,
	pRuntime                      *GFcrawlerRuntime,
	pRuntimeSys                   *gf_core.RuntimeSys) (*gf_images_core.GFimage, *gf_images_core.GFimageThumbs, string, *gf_core.GFerror) {


	//------------------------
	// IMAGE_DOWNLOAD - download image from some external source
	localImageFilePathStr, gfErr := imageDownload(pImage, pImagesStoreLocalDirPathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, "", gfErr
	}

	//------------------------
	image, image_thumbs, gfErr := imageProcess(pImage,
		"", // p_gf_image_id_str
		localImageFilePathStr,
		pImagesStoreLocalDirPathStr,

		pMediaDomainStr,
		pCrawledImagesS3bucketNameStr,
		pRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, "", gfErr
	}

	//------------------------
	// S3_UPLOAD
	gfErr = imageS3upload(pImage,
		localImageFilePathStr,
		image_thumbs,
		pCrawledImagesS3bucketNameStr,
		pRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, "", gfErr
	}

	//------------------------

	return image, image_thumbs, localImageFilePathStr, nil
}

//--------------------------------------------------
// STAGES
//--------------------------------------------------
func images__stage__pull_image_links(pURLfetch *Gf_crawler_url_fetch,
	pCrawlerNameStr string,
	pCycleRunIDstr  string,
	pRuntime        *GFcrawlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_pipeline.images__stage__pull_image_links")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE - STAGE - pull_image_links")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	page_imgs__pipeline_infos_lst := []*gf_page_img__pipeline_info{}
	pURLfetch.goquery_doc.Find("img").Each(func(p_i int, p_elem *goquery.Selection) {

		img_src_str, _      := p_elem.Attr("src")
		origin_page_url_str := pURLfetch.Url_str
		
		// GF_PAGE_IMG__LINK
		page_img_link := &gf_page_img_link{
			img_src_str:         img_src_str,
			origin_page_url_str: origin_page_url_str,
		}

		// GF_PAGE_IMG__PIPELINE_INFO
		page_img__pipeline_info := &gf_page_img__pipeline_info{
			link: page_img_link,
		}

		page_imgs__pipeline_infos_lst = append(page_imgs__pipeline_infos_lst, page_img__pipeline_info)
	})

	return page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
func images__stage__create_page_images(pCrawlerNameStr string,
	pCycleRunIDstr                  string,
	p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	pRuntime                        *GFcrawlerRuntime,
	pRuntimeSys                     *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_pipeline.images__stage__create_page_images")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE - STAGE - create_page_images")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _, page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		//------------------
		// CRAWLER_PAGE_IMG

		gf_img, gf_err := images_adt__prepare_and_create(pCrawlerNameStr,
			pCycleRunIDstr,                       // pCycleRunIDstr
			page_img__pinfo.link.img_src_str,         // p_img_src_url_str
			page_img__pinfo.link.origin_page_url_str, // p_origin_page_url_str
			pRuntime,
			pRuntimeSys)
		if gf_err != nil {
			page_img__pinfo.gf_error = gf_err
			continue // IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}
		//------------------
		// CRAWLER_PAGE_IMG_REF

		gf_img_ref := images_adt__ref_create(pCrawlerNameStr,
			pCycleRunIDstr,
			gf_img.Url_str,                           // p_image_url_str
			gf_img.Domain_str,                        // p_image_url_domain_str
			page_img__pinfo.link.origin_page_url_str, // p_origin_page_url_str
			gf_img.Origin_page_url_domain_str,        // p_origin_page_url_domain_str
			pRuntimeSys)

		//------------------
		// GIF
		if  gf_img.Img_ext_str == "gif" {

			//IMPORTANT!! - all GIF images are valid_for_usage, regardless of size
			gf_img.Valid_for_usage_bool = true
		}
		//------------------

		page_img__pinfo.page_img     = gf_img
		page_img__pinfo.page_img_ref = gf_img_ref
	}

	return p_page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
func images__stage__page_images_persist(pCrawlerNameStr string,
	p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	pRuntime                        *GFcrawlerRuntime,
	pRuntimeSys                     *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_pipeline.images__stage__page_images_persist")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - page_images_persist")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _, page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		// IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		page_img := page_img__pinfo.page_img

		//------------------
		img_exists_bool, gfErr := Image__db_create(page_img__pinfo.page_img, pRuntime, pRuntimeSys)
		if gfErr != nil {
			t := "image_db_create__failed"
			m := "failed db creation of image with img_url_str - "+page_img.Url_str
			Create_error_and_event(t,m,map[string]interface{}{"origin_page_url_str": page_img__pinfo.link.origin_page_url_str,}, page_img.Url_str, pCrawlerNameStr,
				gfErr, pRuntime, pRuntimeSys)

			page_img__pinfo.gf_error = gfErr
			continue // IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}

		page_img__pinfo.exists_bool = img_exists_bool

		//------------------
		gfErr = Image__db_create_ref(page_img__pinfo.page_img_ref, pRuntime, pRuntimeSys)
		if gfErr != nil {
			t := "image_ref_db_create__failed"
			m := "failed db creation of image_ref with img_url_str - "+page_img.Url_str
			Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str":page_img__pinfo.link.origin_page_url_str,}, page_img.Url_str, pCrawlerNameStr,
				gfErr, pRuntime, pRuntimeSys)

			page_img__pinfo.gf_error = gfErr
			continue // IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}
		
		//------------------
	}
	return p_page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
func images__stages__process_images(pCrawlerNameStr string,
	p_page_imgs__pipeline_infos_lst   []*gf_page_img__pipeline_info,
	p_images_store_local_dir_path_str string,
	p_origin_page_url_str             string,

	pMediaDomainStr                   string,
	pS3bucketNameStr                  string,
	pRuntime                          *GFcrawlerRuntime,
	pRuntimeSys                       *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images.images__stages__process_images")

	//------------------
	// // STAGE - determine if image is NSFV (contains nudity)
	// // FIX!! - check if the processing cost of large images is not lower then determening NSFV first on large images,
	// //         and then processing (which is whats done now). perhaps processing all images and then taking the 
	
	// page_imgs__pinfos_with_nsfv_lst := images__stage__determine_are_nsfv(pCrawlerNameStr,
	// 	p_page_imgs__pipeline_infos_lst,
	// 	p_origin_page_url_str,
	// 	pRuntime,
	// 	pRuntimeSys)
	//------------------
	// STAGE - process images - resize for all thumbnail sizes

	page_imgs__pinfos_with_thumbs_lst := images__stage__process_images(pCrawlerNameStr,
		p_page_imgs__pipeline_infos_lst, // page_imgs__pinfos_with_nsfv_lst,
		p_images_store_local_dir_path_str,
		p_origin_page_url_str,

		pMediaDomainStr,
		pS3bucketNameStr,
		pRuntime,
		pRuntimeSys)

	//------------------
	return page_imgs__pinfos_with_thumbs_lst
}

//--------------------------------------------------
func images__stages_cleanup(p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	pRuntime    *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_pipeline.images__stages_cleanup")

	// IMPORTANT!! - delete local tmp transformed image, since the files
	//               have just been uploaded to S3 so no need for them localy anymore
	//               crawling servers are not meant to hold their own image files,
	//               and service runs in Docker with temporary 
	for _, page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		// IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		// IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if page_img__pinfo.exists_bool {
			continue
		}

		gfErr := image__cleanup(page_img__pinfo.local_file_path_str, page_img__pinfo.thumbs, pRuntimeSys)
		if gfErr != nil {
			page_img__pinfo.gf_error = gfErr
			continue // IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}
	}

	return p_page_imgs__pipeline_infos_lst
}