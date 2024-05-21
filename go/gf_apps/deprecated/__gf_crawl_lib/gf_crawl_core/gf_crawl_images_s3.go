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
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//--------------------------------------------------
// STAGE

func imagesS3stageStoreImages(pCrawlerNameStr string,
	pPageImagesPipelineInfosLst []*gf_page_img__pipeline_info,
	pOriginPageURLstr           string,
	pS3bucketNameStr            string,
	pRuntime                    *GFcrawlerRuntime,
	pRuntimeSys                 *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - s3_store_images")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _, pageImagePipelineInfo := range pPageImagesPipelineInfosLst {

		// IMPORTANT!! - skip failed images
		if pageImagePipelineInfo.gf_error != nil {
			continue
		}

		// IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if pageImagePipelineInfo.exists_bool {
			continue
		}

		// IMPORTANT!! - check image is not flagged as a NSFV image
		if pageImagePipelineInfo.nsfv_bool {
			continue
		}

		//------------------
		// IMPORTANT!! - only store/persist if they are valid (of the right dimensions) or
		//               if they're a GIF (all GIF's are stored/persisted,
		//               even if they determined to be NSFV for some reason).

		if pageImagePipelineInfo.page_img.Img_ext_str == "gif" || pageImagePipelineInfo.page_img.Valid_for_usage_bool {

			gfErr := imageS3upload(pageImagePipelineInfo.page_img,
				pageImagePipelineInfo.local_file_path_str,
				pageImagePipelineInfo.thumbs,
				pS3bucketNameStr,
				pRuntime,
				pRuntimeSys)

			if gfErr != nil {
				t := "image_s3_upload__failed"
				m := "failed s3 uploading of image with img_url_str - "+pageImagePipelineInfo.page_img.Url_str
				CreateErrorAndEvent(t, m, 
					map[string]interface{}{"origin_page_url_str": pOriginPageURLstr,},
					pageImagePipelineInfo.page_img.Url_str,
					pCrawlerNameStr,
					gfErr, pRuntime, pRuntimeSys)
				pageImagePipelineInfo.gf_error = gfErr
				continue // IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
			}
		}

		//------------------
	}

	return pPageImagesPipelineInfosLst
}

//--------------------------------------------------

func imageS3upload(pImage *GFcrawlerPageImage,
	pLocalImageFilePathStr string,
	pImageThumbs           *gf_images_core.GFimageThumbs,
	pS3bucketNameStr       string,
	pRuntime               *GFcrawlerRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) *gf_core.GFerror {
	
	cyan   := color.New(color.FgCyan, color.BgWhite).SprintFunc()
	yellow := color.New(color.FgYellow, color.BgBlack).SprintFunc()
	fmt.Printf("\n%s GF_CRAWL_PAGE_IMG TO S3 - id[%s] - local_file[%s]\n\n", cyan("UPLOADING"),
		yellow(pImage.IDstr),
		yellow(pLocalImageFilePathStr))

	gfErr := gf_images_core.S3storeImage(pLocalImageFilePathStr,
		pImageThumbs,
		pS3bucketNameStr,
		pRuntime.S3info,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	imageS3dbFlagAsUploaded(pImage, pRuntimeSys)

	return nil
}

//--------------------------------------------------
// UPDATE_DB

// flag crawler page image as persisted on s3
func imageS3dbFlagAsUploaded(pImage *GFcrawlerPageImage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()

	pImage.S3_stored_bool = true
	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"t":        "crawler_page_img",
			"hash_str": pImage.Hash_str,
		},
		bson.M{
			"$set": bson.M{"s3_stored_bool": true},
		})
		
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update an crawler_page_img s3_stored flag by its hash",
			"mongodb_update_error",
			map[string]interface{}{"image_hash_str": pImage.Hash_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}
	return nil
}