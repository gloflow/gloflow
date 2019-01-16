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

package gf_crawl_core

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_gif_lib"
	//"github.com/davecgh/go-spew/spew"
)
//--------------------------------------------------
func images__stage__process_images(p_crawler_name_str string,
	p_page_imgs__pipeline_infos_lst   []*gf__page_img__pipeline_info,
	p_images_store_local_dir_path_str string,
	p_origin_page_url_str             string,
	p_s3_bucket_name_str              string,
	p_runtime                         *Crawler_runtime,
	p_runtime_sys                     *gf_core.Runtime_sys) []*gf__page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_process.images__stage__process_images")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - process_images")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _,page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		//IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		//IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if page_img__pinfo.exists_bool {
			continue
		}

		//IMPORTANT!! - check image is not flagged as a NSFV image
		if page_img__pinfo.nsfv_bool {
			continue
		}

		//----------------------------
		//IMAGE_PROCESS
		_, gf_image_thumbs, gf_err := image__process(page_img__pinfo.page_img,
			page_img__pinfo.local_file_path_str,
			p_images_store_local_dir_path_str,
			p_s3_bucket_name_str,
			p_runtime,
			p_runtime_sys)
		//----------------------------
		
		if gf_err != nil {
			t:="image_process__failed"
			m:="failed processing of image with img_url_str - "+page_img__pinfo.page_img.Url_str
			Create_error_and_event(t,m,map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,}, page_img__pinfo.page_img.Url_str, p_crawler_name_str,
				gf_err, p_runtime, p_runtime_sys)

			page_img__pinfo.gf_error = gf_err
			continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}

		//UPDATE__PAGE_IMG_PINFO
		page_img__pinfo.thumbs = gf_image_thumbs
	}
	return p_page_imgs__pipeline_infos_lst
}
//--------------------------------------------------
func image__process(p_page_img *Crawler_page_img,
	p_local_image_file_path_str       string,
	p_images_store_local_dir_path_str string,
	p_s3_bucket_name_str              string,
	p_runtime                         *Crawler_runtime,
	p_runtime_sys                     *gf_core.Runtime_sys) (*gf_images_utils.Gf_image, *gf_images_utils.Gf_image_thumbs, *gf_core.Gf_error) {
	//p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_process.image__process()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	p_runtime_sys.Log_fun("INFO",cyan("       >>>>>>>>>>>>> ----------------------------- ")+yellow("PROCESS_IMAGE"))

	//----------------------------
	//GIF
	if p_page_img.Img_ext_str == "gif" {

		image_client_type_str := "gf_crawl_images" 
		image_flows_names_lst := []string{"discovered","gifs",}

		gif_download_and_frames__local_dir_path_str := p_images_store_local_dir_path_str
		gf_gif, _, gf_err := gf_gif_lib.Process(p_page_img.Url_str, //p_image_source_url_str *string,
			p_page_img.Origin_page_url_str,
			gif_download_and_frames__local_dir_path_str,
			image_client_type_str,
			image_flows_names_lst,
			true, //p_create_new_db_img_bool
			p_s3_bucket_name_str,
			p_runtime.S3_info,
			p_runtime_sys)

		if gf_err != nil {
			return nil, nil, gf_err
		}													

		gf_image_id_str := gf_gif.Gf_image_id_str
		gf_err           = image__update_after_process(p_page_img, gf_image_id_str, p_runtime_sys)
		if gf_err != nil {
			return nil, nil, gf_err
		}

		return nil, nil, nil
	//----------------------------
	//GENERAL
	} else {
	
		thumbnails_local_dir_path_str     := p_images_store_local_dir_path_str
		gf_image, gf_image_thumbs, gf_err := image__process_bitmap(p_page_img,
			p_local_image_file_path_str,
			thumbnails_local_dir_path_str,
			p_runtime_sys)
		if gf_err != nil {
			return nil, nil, gf_err
		}


		//spew.Dump(gf_image)


		gf_image_id_str := gf_image.Id_str
		gf_err           = image__update_after_process(p_page_img, gf_image_id_str, p_runtime_sys)
		if gf_err != nil {
			return nil, nil, gf_err
		}

		return gf_image, gf_image_thumbs, nil
	}
	//----------------------------
	return nil, nil, nil
}
//--------------------------------------------------
func image__process_bitmap(p_page_img *Crawler_page_img,
	p_local_image_file_path_str     string,
	p_thumbnails_local_dir_path_str string,
	p_runtime_sys                   *gf_core.Runtime_sys) (*gf_images_utils.Gf_image, *gf_images_utils.Gf_image_thumbs, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_process.image__process_bitmap()")

	//----------------------
	//CONFIG
	image_client_type_str := "gf_crawl_images" 
	image_flows_names_lst := []string{"discovered",}
	//----------------------

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//-------------------
	img_width_int,img_height_int,gf_err := gf_images_utils.Get_image_dimensions__from_filepath(p_local_image_file_path_str, p_runtime_sys)
	if gf_err != nil {
		return nil, nil, gf_err
	}
	//-------------------

	//IMPORTANT!! - check that the image is too small, and is likely to be irrelevant 
	//              part of a particular page
	if img_width_int <= 130 || img_height_int <= 130 {
		p_runtime_sys.Log_fun("INFO",yellow("IMG IS SMALLER THEN MINIMUM DIMENSIONS (width-"+cyan(fmt.Sprint(img_width_int))+"/height-"+cyan(fmt.Sprint(img_height_int))+")"))
		return nil, nil, nil
	} else {

		//--------------------------------
		//TRANSFORM DOWNLOADED IMAGE - CREATE THUMBS, SAVE TO DB, AND UPLOAD TO AWS_S3

		gf_image_id_str,gf_err := gf_images_utils.Image__create_id_from_url(p_page_img.Url_str, p_runtime_sys)
		if gf_err != nil {
			return nil,nil,gf_err
		}
		image_origin_url_str      := p_page_img.Url_str
		image_origin_page_url_str := p_page_img.Origin_page_url_str
		
		//IMPORTANT!! - this creates a Gf_image object, and persists it in the DB ("t" == "img"),
		//              also creates gf_image thumbnails as local files.


		gf_image, gf_image_thumbs, gf_err := gf_images_utils.Transform_image(gf_image_id_str,
			image_client_type_str,
			image_flows_names_lst,
			image_origin_url_str,
			image_origin_page_url_str,
			p_local_image_file_path_str,
			p_thumbnails_local_dir_path_str,
			p_runtime_sys)
		if gf_err != nil {
			return nil, nil, gf_err
		}
		//--------------------------------

		return gf_image, gf_image_thumbs, nil
	}

	return nil, nil, nil
}