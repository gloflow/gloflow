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
                space of images in "flows". it copies images from gf_crawler storage to gf_images service storage*/

package gf_crawl_core

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_flows"
)

//--------------------------------------------------
// adds an image already crawled from an external source URL to some named list of flows in the gf_images app/system.
// to do this it adds the flow_name to the gf_image DB record, and then copies the discovered image file from
// gf_crawlers file_storage (S3/IPFS) to gf_images service file_storage (S3/IPFS).
// at the moment this is called directly in the gf_crawl HTTP handler.

func FlowsAddExternImage(pCrawlerPageImageIDstr Gf_crawler_page_image_id,
	pFlowsNamesLst                []string,
	pMediaDomainStr               string,
	pCrawledImagesS3bucketNameStr string,
	pImagesS3bucketNameStr        string,
	pRuntime                      *GFcrawlerRuntime,
	pRuntimeSys                   *gf_core.RuntimeSys) *gf_core.GFerror {

	green := color.New(color.BgGreen, color.FgBlack).SprintFunc()
	cyan := color.New(color.FgWhite, color.BgCyan).SprintFunc()

	// this is used temporarily to donwload images to, before upload to S3
	imagesStoreLocalDirPathStr := "."

	fmt.Printf("crawler_page_image_id_str - %s\n", pCrawlerPageImageIDstr)
	fmt.Printf("flows_names               - %s\n", fmt.Sprint(pFlowsNamesLst))

	// DB - get gf_crawler_page_image from the DB
	pageImage, gfErr := image__db_get(pCrawlerPageImageIDstr, pRuntime, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	imageIDstr := pageImage.Gf_image_id_str
	imagesS3bucketUploadCompleteBool := false

	//--------------------------
	// SPECIAL_CASE
	// IMPORTANT!! - some crawler_page_images dont have their imageIDstr set,
	//               which means that they dont have their corresponding gf_image.
	if imageIDstr == "" {

		pRuntimeSys.LogFun("INFO", "")
		pRuntimeSys.LogFun("INFO", "CRAWL_PAGE_IMAGE MISSING ITS GF_IMAGE --- STARTING_PROCESSING")
		pRuntimeSys.LogFun("INFO", "")


		gfImage, gf_image_thumbs, localImageFilePathStr, gfErr := images_pipe__single_simple(pageImage,
			imagesStoreLocalDirPathStr,

			pMediaDomainStr,
			pCrawledImagesS3bucketNameStr,
			pRuntime,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		imageIDstr = gfImage.IDstr

		//-------------------
		// S3_UPLOAD_TO_GF_IMAGES_BUCKET
		// IMPORTANT!! - crawler_page_image and its thumbs are uploaded to the crawled images S3 bucket,
		//               but gf_images service /images/d endpoint redirects users to the gf_images
		//               S3 bucket (gf--img).
		//               so we need to upload the new image to that gf_images S3 bucket as well.
		// FIX!! - too much uploading, very inefficient, figure out a better way!

		gfErr = gf_images_core.S3storeImage(localImageFilePathStr,
			gf_image_thumbs,
			pImagesS3bucketNameStr,
			pRuntime.S3_info,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		// IMPORTANT!! - gf_images service has its own dedicate S3 bucket, which is different from the gf_crawl bucket.
		//               gf_images_core.Trans__s3_store_image() uploads the image and its thumbs to S3, 
		//               to indicate that we dont need to upload it later again.
		imagesS3bucketUploadCompleteBool = true

		//-------------------
		//CLEANUP

		gfErr = image__cleanup(localImageFilePathStr, gf_image_thumbs, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
		//-------------------
	}

	//--------------------------
	// ADD_FLOWS_NAMES_TO_IMAGE_DB_RECORD

	// IMPORTANT!! - for each flow_name add that name to the target gf_image DB record.
	for _, flowNameStr := range pFlowsNamesLst {
		gfErr := gf_images_flows.Flows_db__add_flow_name_to_image(flowNameStr, imageIDstr, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	//--------------------------
	// S3_COPY_BETWEEN_BUCKETS - gf--discovered--img -> gf--img
	//                           only for gf_images that have not already been uploaded to the gf--img bucket
	//                           because they needed to be reprecossed and were downloaded from a URL onto
	//                           the local FS first.

	if !imagesS3bucketUploadCompleteBool {

		sourceCrawlS3bucketStr := pCrawledImagesS3bucketNameStr

		fmt.Printf("\n%s - %s -> %s\n\n", green("COPYING IMAGE between S3 BUCKETS"), cyan(sourceCrawlS3bucketStr), cyan(pImagesS3bucketNameStr))

		gfImage, gfErr := gf_images_core.DB__get_image(imageIDstr, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		/*S3__get_image_original_file_s3_filepath is wrong!! FIXX!!!
		the path of the originl_file that its returning is of a file named by its gf_img ID, which is wrong. 
		that filename should be of the file as it was original found in a html page or elsewhere.
		all images added via browser extension, or added by crawler, are named with original file name, not with ID.
		figure out if fixing this is going to break already added images (images added to a flow here from crawled images), 
		since they're all named by ID now (which is a bug)*/

		originalFileS3pathStr                              := gf_images_core.S3__get_image_original_file_s3_filepath(gfImage, pRuntimeSys)
		tSmallS3pathStr, tMediumS3pathStr, tLargeS3pathStr := gf_images_core.S3__get_image_thumbs_s3_filepaths(gfImage, pRuntimeSys)

		fmt.Printf("original_file_s3_path_str - %s\n", originalFileS3pathStr)
		fmt.Printf("t_small_s3_path_str       - %s\n", tSmallS3pathStr)
		fmt.Printf("t_medium_s3_path_str      - %s\n", tMediumS3pathStr)
		fmt.Printf("t_large_s3_path_str       - %s\n", tLargeS3pathStr)

		// ADD!! - copy t_small_s3_path_str first, and then copy originalFileS3pathStr and medium/large thumb in separate goroutines
		//         (in parallel and after the response returns back to the user). 
		//         this is critical to improve perceived user response time, since the small thumb is necessary to view an image in flows, 
		//         but the original_file and medium/large thumbs are not (and can take much longer to S3 copy without the user noticing).
		filesToCopyLst := []string{
			originalFileS3pathStr,
			tSmallS3pathStr, 
			tMediumS3pathStr,
			tLargeS3pathStr,
		}
		
		for _, S3pathStr := range filesToCopyLst {

			// IMPORTANT!! - the Crawler_page_img has alread been uploaded to S3, so we dont need 
			//               to download it from S3 and reupload to gf_images S3 bucket. Instead we do 
			//               a file copy operation within the S3 system without downloading here.

			// source_bucket_and_file__s3_path_str := filepath.Clean(fmt.Sprintf("/%s/%s", sourceCrawlS3bucketStr, s3_path_str))

			gfErr := gf_core.S3copyFile(sourceCrawlS3bucketStr, // p_source_file__s3_path_str
				S3pathStr,
				pImagesS3bucketNameStr, // p_target_bucket_name_str,
				S3pathStr,              // p_target_file__s3_path_str
				pRuntime.S3_info,
				pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
		}
	}
	//--------------------------
	
	return nil
}