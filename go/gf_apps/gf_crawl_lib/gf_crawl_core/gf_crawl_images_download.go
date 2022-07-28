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
	"time"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//--------------------------------------------------
func images__stage__download_images(p_crawler_name_str string,
	p_page_imgs__pipeline_infos_lst   []*gf_page_img__pipeline_info,
	p_images_store_local_dir_path_str string,
	p_origin_page_url_str             string,
	p_runtime                         *GFcrawlerRuntime,
	p_runtime_sys                     *gf_core.Runtime_sys) []*gf_page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_images_download.images__stage__download_images")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - download_images")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	//------------------
	// ADD!! - download images in batches. some pages have potentially 100's of images of various sizes. 
	//         browsers download these images in groups, because they limitations on the number of simultaneously 
	//         open TCP connections, so they might download 5-6 images at the same time. 
	//         so in this loop schedule image downlading to happen in separate goroutines, N at a time.
	//------------------

	for _, page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		// IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		// IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if page_img__pinfo.exists_bool {
			continue
		}

		start_time_f := float64(time.Now().UnixNano())/1000000000.0

		//------------------
		// DOWNLOAD
		// IMPORTANT!! - all images done as fast as possible (without sleeps/pauses)
		//               since when users view a page in their browser the browser issues all requests
		//               for all the images in the page immediatelly. 

		local_image_file_path_str, gf_err := image__download(page_img__pinfo.page_img,
			p_images_store_local_dir_path_str,
			p_runtime_sys)
		
		if gf_err != nil {
			t := "image_download__failed"
			m := "failed downloading of image with URL - "+page_img__pinfo.page_img.Url_str
			Create_error_and_event(t, m,
				map[string]interface{}{"origin_page_url_str": p_origin_page_url_str,}, 
				page_img__pinfo.page_img.Url_str,
				p_crawler_name_str,
				gf_err, p_runtime, p_runtime_sys)

			page_img__pinfo.gf_error = gf_err
			continue // IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
		}

		//------------------

		end_time_f := float64(time.Now().UnixNano())/1000000000.0

		page_img__pinfo.local_file_path_str = local_image_file_path_str

		//------------------
		// SEND_EVENT
		if p_runtime.Events_ctx != nil {
			events_id_str  := "crawler_events"
			event_type_str := "image_download__http_request__done"
			msg_str        := "completed downloading an image over HTTP"
			data_map       := map[string]interface{}{
				"img_url_str":  page_img__pinfo.page_img.Url_str,
				"start_time_f": start_time_f,
				"end_time_f":   end_time_f,
			}

			gf_events.Events__send_event(events_id_str,
				event_type_str, // p_type_str
				msg_str,        // p_msg_str
				data_map,
				p_runtime.Events_ctx,
				p_runtime_sys)
		}
		//------------------
	}

	return p_page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
func image__download(pImage *Gf_crawler_page_image,
	p_images_store_local_dir_path_str string,
	p_runtime_sys                     *gf_core.Runtime_sys) (string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_images_download.image__download()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	p_runtime_sys.Log_fun("INFO", cyan("       >>>>>>>>>>>>> ----------------------------- ")+yellow("DOWNLOAD_IMAGE"))

	//-------------------
	// DOWNLOAD
	// IMPORTANT!! - this creates a new gf_images ID, from the image URL
	localImageFilePathStr, imageIDstr, gfErr := gf_images_core.FetcherGetExternImage(pImage.Url_str,
		p_images_store_local_dir_path_str,

		// IMPORTANT!! - dont add any time delay, instead download images as fast as possible
		//               since they're all in the same page, and are expected to be downloaded 
		//               by the users browser in rapid succession, so no need to simulate user delay
		false, // p_random_time_delay_bool
		p_runtime_sys)
	if gfErr != nil {
		return "", gfErr
	}

	//-------------------
	// DB_UPDATE
	gfErr = image__db_mark_as_downloaded(pImage, p_runtime_sys)
	if gfErr != nil {
		return "", gfErr
	}

	gfErr = image__db_set_gf_image_id(imageIDstr, pImage, p_runtime_sys)
	if gfErr != nil {
		return "", gfErr
	}
	
	//-------------------

	pImage.Downloaded_bool = true
	pImage.Gf_image_id_str = imageIDstr

	return localImageFilePathStr, nil
}