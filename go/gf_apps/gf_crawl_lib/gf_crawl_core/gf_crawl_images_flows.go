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

/*IMPORTANT!! - functions in this file are responsible for bridging the gf_crawler space of images, with the gf_images service
                space of images in "flows"*/

package gf_crawl_core

import (
	"fmt"
	"path/filepath"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)

//--------------------------------------------------
/*ads an image already crawled from an external source URL to some named list of flows.
to do this it ads the flow_name to the gf_image DB record, and then copies the discovered image file from
gf_crawlers file_storage (S3) to gf_images service file_storage (S3).*/

func Flows__add_extern_image(p_crawler_page__gf_image_id_str string,
	p_flows_names_lst                   []string,
	p_crawled_images_s3_bucket_name_str string,
	p_gf_images_s3_bucket_name_str      string,
	p_runtime                           *Gf_crawler_runtime,
	p_runtime_sys                       *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_images_flows.Flows__add_extern_image()")

	green := color.New(color.BgGreen, color.FgBlack).SprintFunc()
	cyan := color.New(color.FgWhite, color.BgCyan).SprintFunc()

	//this is used temporarily to donwload images to, before upload to S3
	images_store_local_dir_path_str := "."

	fmt.Printf("crawler_page__gf_image_id - %s\n", p_crawler_page__gf_image_id_str)
	fmt.Printf("flows_names               - %s\n", fmt.Sprint(p_flows_names_lst))

	gf_page_img, gf_err := image__db_get(p_crawler_page__gf_image_id_str, p_runtime, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	gf_image_id_str := gf_page_img.Image_id_str
	gf_images_s3_bucket__upload_complete__bool := false

	//--------------------------
	//SPECIAL_CASE
	//IMPORTANT!! - some crawler_page_images dont have their gf_image_id_str set,
	//              which means that they dont have their corresponding gf_image.
	if gf_image_id_str == "" {

		p_runtime_sys.Log_fun("INFO", "")
		p_runtime_sys.Log_fun("INFO", "CRAWL_PAGE_IMAGE MISSING ITS GF_IMAGE --- STARTING_PROCESSING")
		p_runtime_sys.Log_fun("INFO", "")

		//S3_UPLOAD - images__process_crawler_page_image() uploads image and its thumbs to S3 
		//            after it finishes processing it.
		gf_image, gf_image_thumbs, local_image_file_path_str, gf_err := images_pipe__single_simple(gf_page_img,
			images_store_local_dir_path_str,
			p_crawled_images_s3_bucket_name_str,
			p_runtime,
			p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		gf_image_id_str = gf_image.Id_str
		//-------------------
		//S3_UPLOAD_TO_GF_IMAGES_BUCKET
		//IMPORTANT!! - crawler_page_image and its thumbs are uploaded to the crawled images S3 bucket,
		//              but gf_images service /images/d endpoint redirects users to the gf_images
		//              S3 bucket (gf--img).
		//              so we need to upload the new image to that gf_images S3 bucket as well.
		//FIX!! - too much uploading, very inefficient, figure out a better way!

		gf_err = gf_images_utils.S3__store_gf_image(local_image_file_path_str,
			gf_image_thumbs,
			p_gf_images_s3_bucket_name_str,
			p_runtime.S3_info,
			p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		//IMPORTANT!! - gf_images service has its own dedicate S3 bucket, which is different from the gf_crawl bucket.
		//              gf_images_utils.Trans__s3_store_image() uploads the image and its thumbs to S3, 
		//              to indicate that we dont need to upload it later again.
		gf_images_s3_bucket__upload_complete__bool = true
		//-------------------
		//CLEANUP

		gf_err = image__cleanup(local_image_file_path_str, gf_image_thumbs, p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
		//-------------------
	}
	//--------------------------
	//ADD_FLOWS_NAMES_TO_IMAGE_DB_RECORD

	//IMPORTANT!! - for each flow_name add that name to the target gf_image DB record.
	for _, flow_name_str := range p_flows_names_lst {
		gf_err := gf_images_lib.Flows_db__add_flow_name_to_image(flow_name_str, gf_image_id_str, p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
	}
	//--------------------------
	//S3_COPY_BETWEEN_BUCKETS - gf--discovered--img -> gf--img
	//                          only for gf_images that have not already been uploaded to the gf--img bucket
	//                          because they needed to be reprecossed and were downloaded from a URL onto
	//                          the local FS first.

	if !gf_images_s3_bucket__upload_complete__bool {

		source_gf_crawl_s3_bucket_str := p_crawled_images_s3_bucket_name_str

		fmt.Printf("\n%s - %s -> %s\n\n", green("COPYING IMAGE between S3 BUCKETS"), cyan(source_gf_crawl_s3_bucket_str), cyan(p_gf_images_s3_bucket_name_str))

		gf_image, gf_err := gf_images_utils.DB__get_image(gf_image_id_str, p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		/*S3__get_image_original_file_s3_filepath is wrong!! FIXX!!!
		the path of the originl_file that its returning is of a file named by its gf_img ID, which is wrong. 
		that filename should be of the file as it was original found in a html page or elsewhere.
		all images added via browser extension, or added by crawler, are named with original file name, not with ID.
		figure out if fixing this is going to break already added images (images added to a flow here from crawled images), 
		since they're all named by ID now (which is a bug)*/

		original_file_s3_path_str                                      := gf_images_utils.S3__get_image_original_file_s3_filepath(gf_image, p_runtime_sys)
		t_small_s3_path_str, t_medium_s3_path_str, t_large_s3_path_str := gf_images_utils.S3__get_image_thumbs_s3_filepaths(gf_image, p_runtime_sys)



		fmt.Printf("original_file_s3_path_str - %s\n", original_file_s3_path_str)
		fmt.Printf("t_small_s3_path_str       - %s\n", t_small_s3_path_str)
		fmt.Printf("t_medium_s3_path_str      - %s\n", t_medium_s3_path_str)
		fmt.Printf("t_large_s3_path_str       - %s\n", t_large_s3_path_str)




		//ADD!! - copy t_small_s3_path_str first, and then copy original_file_s3_path_str and medium/large thumb in separate goroutines
		//        (in parallel and after the response returns back to the user). 
		//        this is critical to improve perceived user response time, since the small thumb is necessary to view an image in flows, 
		//        but the original_file and medium/large thumbs are not (and can take much longer to S3 copy without the user noticing).
		files_to_copy_lst := []string{
			original_file_s3_path_str,
			t_small_s3_path_str, 
			t_medium_s3_path_str,
			t_large_s3_path_str,
		}
		


		
		for _, s3_path_str := range files_to_copy_lst {

			//IMPORTANT!! - the Crawler_page_img has alread been uploaded to S3, so we dont need 
			//              to download it from S3 and reupload to gf_images S3 bucket. Instead we do 
			//              a file copy operation within the S3 system without downloading here.

			source_bucket_and_file__s3_path_str := filepath.Clean(fmt.Sprintf("/%s/%s", source_gf_crawl_s3_bucket_str, s3_path_str))


			//DEBUGGING
			fmt.Println("")
			fmt.Println("==========================")
			fmt.Println(p_gf_images_s3_bucket_name_str)
			fmt.Println(source_bucket_and_file__s3_path_str)
			fmt.Println(s3_path_str)
			fmt.Println("")

			gf_err := gf_core.S3__copy_file(p_gf_images_s3_bucket_name_str, //p_target_bucket_name_str,
				source_bucket_and_file__s3_path_str, //p_source_file__s3_path_str
				s3_path_str,                         //p_target_file__s3_path_str
				p_runtime.S3_info,
				p_runtime_sys)
			if gf_err != nil {
				return gf_err
			}
		}
	}
	//--------------------------
	
	return nil
}