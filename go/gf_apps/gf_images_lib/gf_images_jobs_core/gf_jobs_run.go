// SPDX-License-Identifier: GPL-2.0
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

package gf_images_jobs_core

import (
	"fmt"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
)


//-------------------------------------------------

func runJobClassifyImages(pImagesToProcessLst []GFimageClassificationToProcess,
	pPyDirPathStr string,
	pStorage      *gf_images_storage.GFimageStorage,
	pJobRuntime   *GFjobRuntime,
	pMetricsCore  *gf_images_core.GFmetrics,
	pCtx          context.Context,
	pRuntimeSys   *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {

	
	// prepare python classify input
	imagesIDsLst := []gf_images_core.GFimageID{}
	for _, imageToProcess := range pImagesToProcessLst {
		imageIDstr := imageToProcess.GFimageIDstr
		imagesIDsLst = append(imagesIDsLst, gf_images_core.GFimageID(imageIDstr))
	}
		


	classesLst, gfErr := gf_images_core.RunClassify(imagesIDsLst,
		pPyDirPathStr,
		pMetricsCore,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}



	return classesLst, nil

}

//-------------------------------------------------

func runJobLocalImages(pImagesToProcessLst []GFimageLocalToProcess,
	pFlowsNamesLst                        []string,
	pImagesStoreLocalDirPathStr           string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pS3info                               *gf_aws.GFs3Info,
	pPluginsPyDirPathStr                  string,
	pStorage                              *gf_images_storage.GFimageStorage,
	pJobRuntime                           *GFjobRuntime,
	pRuntimeSys                           *gf_core.RuntimeSys) []*gf_core.GFerror {
	
	gfErrorsLst := []*gf_core.GFerror{}
	for _, imageToProcess := range pImagesToProcessLst {
		
		fmt.Println(imageToProcess)

		gfErr := pipelineProcessLocalImage(pFlowsNamesLst,
			pS3info,
			pPluginsPyDirPathStr,
			pStorage,
			pRuntimeSys)
		if gfErr != nil {
			gfErrorsLst = append(gfErrorsLst, gfErr)
		}
	}

	return gfErrorsLst
}

//-------------------------------------------------

func runJobUploadedImages(pImagesToProcessLst []GFimageUploadedToProcess,
	pFlowsNamesLst                        []string,
	pImagesStoreLocalDirPathStr           string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pS3info                               *gf_aws.GFs3Info,
	pPluginsPyDirPathStr                  string,
	pStorage                              *gf_images_storage.GFimageStorage,
	pJobRuntime                           *GFjobRuntime,
	pRuntimeSys                           *gf_core.RuntimeSys) []*gf_core.GFerror {

	gfErrorsLst := []*gf_core.GFerror{}
	for _, imageToProcess := range pImagesToProcessLst {

		fmt.Println(imageToProcess)

		imageIDstr    := imageToProcess.GFimageIDstr
		S3filePathStr := imageToProcess.S3filePathStr
		metaMap       := imageToProcess.MetaMap

		gfErr := pipelineProcessUploadedImage(imageIDstr,
			S3filePathStr,
			metaMap,
			pImagesStoreLocalDirPathStr,
			pImagesThumbnailsStoreLocalDirPathStr,
			pFlowsNamesLst,
			pS3info,
			pPluginsPyDirPathStr,
			pStorage,
			pJobRuntime,
			pRuntimeSys)
		if gfErr != nil {
			gfErrorsLst = append(gfErrorsLst, gfErr)
		}
	}

	return gfErrorsLst
}

//-------------------------------------------------
// RUN_JOB__EXTERN_IMAGES

func runJobExternImages(pImagesToProcessLst []GFimageExternToProcess,
	pFlowsNamesLst                        []string,
	pImagesStoreLocalDirPathStr           string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pVideoStoreLocalDirPathStr            string,
	pMediaDomainStr                       string,
	pS3bucketNameStr                      string,
	pS3info                               *gf_aws.GFs3Info,
	pPluginsPyDirPathStr                  string,
	pStorage                              *gf_images_storage.GFimageStorage,
	pJobRuntime                           *GFjobRuntime,
	pCtx                                  context.Context,
	pRuntimeSys                           *gf_core.RuntimeSys) []*gf_core.GFerror {

	ctx := context.Background()

	gfErrorsLst := []*gf_core.GFerror{}
	for _, imageToProcess := range pImagesToProcessLst {

		sourceURLstr     := imageToProcess.SourceURLstr // FIX!! rename source_url_str to origin_url_str
		originPageURLstr := imageToProcess.OriginPageURLstr

		//--------------
		// GET_MIME_CONTENT_TYPE
		// determening mime types from file headers is more robust/reliable 
		// than using file extensions for 
		headersMap, userAgentStr := gf_core.HTTPgetReqConfig()
		imageContentTypeStr, gfErr := gf_core.HTTPdetectMIMEtypeFromURL(sourceURLstr,
			headersMap,
			userAgentStr,
			pCtx,
			pRuntimeSys)
			
		if gfErr != nil {
			jobErrorTypeStr := "get_image_mime_content_type_error"
			_ = jobErrorSend(jobErrorTypeStr, gfErr, sourceURLstr,
				gf_images_core.GFimageID(""),
				pJobRuntime.job_id_str,
				pJobRuntime.job_updates_ch,
				pRuntimeSys)
			gfErrorsLst = append(gfErrorsLst, gfErr)
			continue
		}



		fmt.Printf("image content type - %s\n", imageContentTypeStr)

		
		//--------------
		// IMAGE_ID
		imageIDstr, gfErr := gf_images_core.CreateIDfromURL(sourceURLstr, pRuntimeSys)

		if gfErr != nil {
			jobErrorTypeStr := "create_image_id_error"
			_ = jobErrorSend(jobErrorTypeStr, gfErr, sourceURLstr,
				imageIDstr,
				pJobRuntime.job_id_str,
				pJobRuntime.job_updates_ch,
				pRuntimeSys)
			gfErrorsLst = append(gfErrorsLst, gfErr)
			continue
		}
		
		//--------------
		
		pRuntimeSys.LogFun("INFO", fmt.Sprintf("PROCESSING IMAGE - %s", sourceURLstr))

		//--------------
		// VIDEO

		if imageContentTypeStr == "video/mp4" || imageContentTypeStr == "video/webm" {

			videoIDstr := imageIDstr

			gfErr := pipelineProcessExternVideo(videoIDstr,
				sourceURLstr,
				originPageURLstr,
				pVideoStoreLocalDirPathStr,
				pImagesStoreLocalDirPathStr,
				pImagesThumbnailsStoreLocalDirPathStr,
				pFlowsNamesLst,
				pPluginsPyDirPathStr,
				pStorage,
				pJobRuntime,
				pRuntimeSys)
				
			if gfErr != nil {
				jobErrorTypeStr := "video_process_error"
				_ = jobErrorSend(jobErrorTypeStr, gfErr, sourceURLstr, videoIDstr,
					pJobRuntime.job_id_str,
					pJobRuntime.job_updates_ch,
					pRuntimeSys)
				gfErrorsLst = append(gfErrorsLst, gfErr)
				continue
			}
			
			continue

		//--------------
		// GIF - gifs have their own processing pipeline

		// FIX!! - move GIF processing logic into gf_images_pipeline.go as well,
		//         it doesnt belong here in general images_job logic.
		} else if imageContentTypeStr == "image/gif" {

			//-----------------
			// FLOWS_NAMES
			// check if "gifs" flow is already in the list
			b := false
			for _, s := range pFlowsNamesLst {
				if s == "gifs" {
					b = true
				}
			}
			
			var flowsNamesLst []string
			if b {
				flowsNamesLst = append([]string{"gifs"}, pFlowsNamesLst...)
			} else {
				flowsNamesLst = pFlowsNamesLst
			}

			//-----------------

			_, gfErr := gf_gif_lib.ProcessAndUpload("", // p_gf_imageIDstr
				sourceURLstr,
				originPageURLstr,
				pImagesStoreLocalDirPathStr,
				pJobRuntime.job_client_type_str,
				flowsNamesLst,
				true, // p_create_new_db_img_bool

				pMediaDomainStr,
				pS3bucketNameStr,
				pS3info,
				ctx,
				pRuntimeSys)

			if gfErr != nil {
				jobErrorTypeStr := "gif_process_and_upload_error"
				_ = jobErrorSend(jobErrorTypeStr, gfErr, sourceURLstr, imageIDstr,
					pJobRuntime.job_id_str,
					pJobRuntime.job_updates_ch,
					pRuntimeSys)
				gfErrorsLst = append(gfErrorsLst, gfErr)
				continue
			}

			continue

		//-----------------------
		// STANDARD
		} else {

			gfErr := pipelineProcessExternImage(imageIDstr,
				sourceURLstr,
				originPageURLstr,
				pImagesStoreLocalDirPathStr,
				pImagesThumbnailsStoreLocalDirPathStr,
				pFlowsNamesLst,
				pS3bucketNameStr,
				pS3info,
				pPluginsPyDirPathStr,
				pStorage,
				pJobRuntime,
				pRuntimeSys)

			if gfErr != nil {
				jobErrorTypeStr := "image_process_error"
				_ = jobErrorSend(jobErrorTypeStr, gfErr, sourceURLstr, imageIDstr,
					pJobRuntime.job_id_str,
					pJobRuntime.job_updates_ch,
					pRuntimeSys)
				gfErrorsLst = append(gfErrorsLst, gfErr)
				continue
			}
		}

		//-----------------------
	}
	return gfErrorsLst
}