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
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)

//--------------------------------------------------
type gf_page_img__pipeline_info struct {
	link                *gf_page_img_link
	page_img            *Gf_crawler_page_img
	page_img_ref        *Gf_crawler_page_img_ref
	exists_bool         bool                   //has the page_img already been discovered in the past
	local_file_path_str string
	nsfv_bool           bool
	thumbs              *gf_images_utils.Gf_image_thumbs
	gf_error            *gf_core.Gf_error      //if page_img processing failed at some stage
}

type gf_page_img_link struct {
	img_src_str         string
	origin_page_url_str string
}

//--------------------------------------------------
func images_pipe__from_html(p_url_fetch *Gf_crawler_url_fetch,
	p_cycle_run_id_str          string,
	p_crawler_name_str          string,
	p_images_local_dir_path_str string,
	p_s3_bucket_name_str        string,
	p_runtime                   *Gf_crawler_runtime,
	p_runtime_sys               *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_images_pipeline.images_pipe__from_html()")

	cyan := color.New(color.FgCyan).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	//yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))
	fmt.Println(">> IMAGES__GET_IN_PAGE - "+blue(p_url_fetch.Url_str))
	fmt.Println(cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))

	origin_page_url_str := p_url_fetch.Url_str

	//------------------
	//STAGE - pull all page image links

	page_imgs__pinfos_lst := images__stage__pull_image_links(p_url_fetch,
		p_crawler_name_str,
		p_cycle_run_id_str,
		p_runtime,
		p_runtime_sys)
	//------------------
	//STAGE - create gf_image/gf_image_refs structs

	page_imgs__pinfos_with_imgs_lst := images__stage__create_page_images(p_crawler_name_str,
		p_cycle_run_id_str,
		page_imgs__pinfos_lst,
		p_runtime,
		p_runtime_sys)
	//------------------
	//STAGE - persist gf_image/gf_image_ref to DB
	
	page_imgs__pinfos_with_persists_lst := images__stage__page_images_persist(p_crawler_name_str,
		page_imgs__pinfos_with_imgs_lst,
		p_runtime,
		p_runtime_sys)
	//------------------
	//STAGE - download gf_images
	
	page_imgs__pinfos_with_local_file_paths_lst := images__stage__download_images(p_crawler_name_str,
		page_imgs__pinfos_with_persists_lst,
		p_images_local_dir_path_str,
		origin_page_url_str,
		p_runtime,
		p_runtime_sys)
	//------------------
	//STAGES - process images

	page_imgs__pinfos_with_thumbs_lst := images__stages__process_images(p_crawler_name_str,
		page_imgs__pinfos_with_local_file_paths_lst,
		p_images_local_dir_path_str,
		origin_page_url_str,
		p_s3_bucket_name_str,
		p_runtime,
		p_runtime_sys)
	//------------------
	//STAGE - S3 store all images

	page_imgs__pinfos_with_s3_lst := images_s3__stage__store_images(p_crawler_name_str,
		page_imgs__pinfos_with_thumbs_lst,
		origin_page_url_str,
		p_s3_bucket_name_str,
		p_runtime,
		p_runtime_sys)
	//------------------
	//STAGE - cleanup

	images__stages_cleanup(page_imgs__pinfos_with_s3_lst, p_runtime, p_runtime_sys)
	//------------------
}

//--------------------------------------------------
//SINGLE_IMAGE

func images_pipe__single_simple(p_image *Gf_crawler_page_img,
	p_images_store_local_dir_path_str   string,
	p_crawled_images_s3_bucket_name_str string,
	p_runtime                           *Gf_crawler_runtime,
	p_runtime_sys                       *gf_core.Runtime_sys) (*gf_images_utils.Gf_image, *gf_images_utils.Gf_image_thumbs, string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_process.images_pipe__single_simple")


	//------------------------
	local_image_file_path_str, gf_err := image__download(p_image, p_images_store_local_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return nil, nil, "", gf_err
	}
	//------------------------
	image, image_thumbs, gf_err := image__process(p_image,
		local_image_file_path_str,
		p_images_store_local_dir_path_str,
		p_crawled_images_s3_bucket_name_str,
		p_runtime,
		p_runtime_sys)
	if gf_err != nil {
		return nil, nil, "", gf_err
	}
	//------------------------
	gf_err = image_s3__upload(p_image,
		local_image_file_path_str,
		image_thumbs,
		p_crawled_images_s3_bucket_name_str,
		p_runtime,
		p_runtime_sys)
	if gf_err != nil {
		return nil, nil, "", gf_err
	}
	//------------------------

	return image, image_thumbs, local_image_file_path_str, nil
}

//--------------------------------------------------
//STAGES
//--------------------------------------------------
func images__stage__pull_image_links(p_url_fetch *Gf_crawler_url_fetch,
	p_crawler_name_str string,
	p_cycle_run_id_str string,
	p_runtime          *Gf_crawler_runtime,
	p_runtime_sys      *gf_core.Runtime_sys) []*gf_page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_pipeline.images__stage__pull_image_links")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE - STAGE - pull_image_links")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	page_imgs__pipeline_infos_lst := []*gf_page_img__pipeline_info{}
	p_url_fetch.goquery_doc.Find("img").Each(func(p_i int, p_elem *goquery.Selection) {

		img_src_str,_       := p_elem.Attr("src")
		origin_page_url_str := p_url_fetch.Url_str
		
		page_img_link := &gf_page_img_link{
			img_src_str:        img_src_str,
			origin_page_url_str:origin_page_url_str,
		}

		page_img__pipeline_info := &gf_page_img__pipeline_info{
			link:page_img_link,
		}

		page_imgs__pipeline_infos_lst = append(page_imgs__pipeline_infos_lst, page_img__pipeline_info)
	})

	return page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
func images__stage__create_page_images(p_crawler_name_str string,
	p_cycle_run_id_str              string,
	p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	p_runtime                       *Gf_crawler_runtime,
	p_runtime_sys                   *gf_core.Runtime_sys) []*gf_page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_pipeline.images__stage__create_page_images")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE - STAGE - create_page_images")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _,page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		//------------------
		//CRAWLER_PAGE_IMG

		gf_img,gf_err := images__prepare_and_create(p_crawler_name_str,
			p_cycle_run_id_str,                       //p_cycle_run_id_str
			page_img__pinfo.link.img_src_str,         //p_img_src_url_str
			page_img__pinfo.link.origin_page_url_str, //p_origin_page_url_str
			p_runtime,
			p_runtime_sys)
		if gf_err != nil {
			page_img__pinfo.gf_error = gf_err
			continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}
		//------------------
		//CRAWLER_PAGE_IMG_REF

		gf_img_ref := images__ref_create(p_crawler_name_str,
			p_cycle_run_id_str,
			gf_img.Url_str,                           //p_image_url_str
			gf_img.Domain_str,                        //p_image_url_domain_str
			page_img__pinfo.link.origin_page_url_str, //p_origin_page_url_str
			gf_img.Origin_page_url_domain_str,        //p_origin_page_url_domain_str
			p_runtime_sys)
		//------------------
		//GIF
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
func images__stage__page_images_persist(p_crawler_name_str string,
	p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	p_runtime                       *Gf_crawler_runtime,
	p_runtime_sys                   *gf_core.Runtime_sys) []*gf_page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_pipeline.images__stage__page_images_persist")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - page_images_persist")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _,page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		//IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		page_img := page_img__pinfo.page_img
		//------------------
		img_exists_bool, gf_err := Image__db_create(page_img__pinfo.page_img, p_runtime, p_runtime_sys)
		if gf_err != nil {
			t:="image_db_create__failed"
			m:="failed db creation of image with img_url_str - "+page_img.Url_str
			Create_error_and_event(t,m,map[string]interface{}{"origin_page_url_str":page_img__pinfo.link.origin_page_url_str,}, page_img.Url_str, p_crawler_name_str,
				gf_err, p_runtime, p_runtime_sys)

			page_img__pinfo.gf_error = gf_err
			continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}

		page_img__pinfo.exists_bool = img_exists_bool
		//------------------
		gf_err = Image__db_create_ref(page_img__pinfo.page_img_ref, p_runtime, p_runtime_sys)
		if gf_err != nil {
			t:="image_ref_db_create__failed"
			m:="failed db creation of image_ref with img_url_str - "+page_img.Url_str
			Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str":page_img__pinfo.link.origin_page_url_str,}, page_img.Url_str, p_crawler_name_str,
				gf_err, p_runtime, p_runtime_sys)

			page_img__pinfo.gf_error = gf_err
			continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}
		//------------------
	}
	return p_page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
func images__stages__process_images(p_crawler_name_str string,
	p_page_imgs__pipeline_infos_lst   []*gf_page_img__pipeline_info,
	p_images_store_local_dir_path_str string,
	p_origin_page_url_str             string,
	p_s3_bucket_name_str              string,
	p_runtime                         *Gf_crawler_runtime,
	p_runtime_sys                     *gf_core.Runtime_sys) []*gf_page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images.images__stages__process_images")

	//------------------
	//STAGE - determine if image is NSFV (contains nudity)
	//FIX!! - check if the processing cost of large images is not lower then determening NSFV first on large images,
	//        and then processing (which is whats done now). perhaps processing all images and then taking the 
	
	page_imgs__pinfos_with_nsfv_lst := images__stage__determine_are_nsfv(p_crawler_name_str,
		p_page_imgs__pipeline_infos_lst,
		p_origin_page_url_str,
		p_runtime,
		p_runtime_sys)
	//------------------
	//STAGE - process images - resize for all thumbnail sizes

	page_imgs__pinfos_with_thumbs_lst := images__stage__process_images(p_crawler_name_str,
		page_imgs__pinfos_with_nsfv_lst,
		p_images_store_local_dir_path_str,
		p_origin_page_url_str,
		p_s3_bucket_name_str,
		p_runtime,
		p_runtime_sys)
	//------------------
	return page_imgs__pinfos_with_thumbs_lst
}

//--------------------------------------------------
func images__stages_cleanup(p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) []*gf_page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_pipeline.images__stages_cleanup")

	//IMPORTANT!! - delete local tmp transformed image, since the files
	//              have just been uploaded to S3 so no need for them localy anymore
	//              crawling servers are not meant to hold their own image files,
	//              and service runs in Docker with temporary 
	for _,page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		//IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		//IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if page_img__pinfo.exists_bool {
			continue
		}

		gf_err := image__cleanup(page_img__pinfo.local_file_path_str, page_img__pinfo.thumbs, p_runtime_sys)
		if gf_err != nil {
			page_img__pinfo.gf_error = gf_err
			continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}
	}

	return p_page_imgs__pipeline_infos_lst
}