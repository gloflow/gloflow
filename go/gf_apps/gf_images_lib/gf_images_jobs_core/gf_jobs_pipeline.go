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
/*BosaC.Jan30.2020. <3 volim te zauvek*/

package gf_images_jobs_core

import (
	"fmt"
	"context"
	"path"
	"path/filepath"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_video"
)

//-------------------------------------------------
// PIPELINE__PROCESS_EXTERN_VIDEO

func pipelineProcessExternVideo(pVideoIDstr gf_images_core.GFimageID,
	pVideoSourceURLstr           string,
	pOriginPageURLstr            string,
	pVideosStoreLocalDirPathStr  string,
	pImagesStoreLocalDirPathStr  string,
	pImagesThumbsLocalDirPathStr string,
	pFlowsNamesLst               []string,
	pPluginsPyDirPathStr         string,
	pStorage                     *gf_images_storage.GFimageStorage,
	pJobRuntime                  *GFjobRuntime,
	pRuntimeSys                  *gf_core.RuntimeSys) *gf_core.GFerror {

	
	

	
	//-----------------------
	// GET_VIDEO_FRAME_IMAGE
	imageFileNameStr      := fmt.Sprintf("%s.jpeg", pVideoIDstr)
	imageLocalFilePathStr := fmt.Sprintf("%s/%s", pImagesStoreLocalDirPathStr, imageFileNameStr)
	frameIndexInt := 1

	gfErr := gf_video.GetVideoFrameFromURL(pVideoSourceURLstr,
		imageLocalFilePathStr,
		frameIndexInt,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-----------------------
	// IMAGE_TRANSFORM	

	// FIX!! - this should be passed in from outside this function
	metaMap := map[string]interface{}{}

	imageThumbs, gfErr := jobTransform(pVideoIDstr,
		pFlowsNamesLst,
		pVideoSourceURLstr,
		pOriginPageURLstr,
		metaMap,
		imageLocalFilePathStr,
		pImagesThumbsLocalDirPathStr,
		pPluginsPyDirPathStr,
		pJobRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-----------------------
	// IMAGE_FULLSIZE_STORE

	op := &gf_images_storage.GFputFromLocalOpDef{
		ImageSourceLocalFilePathStr: imageLocalFilePathStr,
		ImageTargetFilePathStr:      imageFileNameStr,
	}
	if pStorage.TypeStr == "s3" {
		op.S3bucketNameStr = pStorage.S3.ExternImagesS3bucketNameStr
	}
	gfErr = gf_images_storage.FilePutFromLocal(op, pStorage, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//--------------------
	// THUMBS_STORE

	gfErr = gf_images_core.StoreThumbnails(imageThumbs,
		pStorage,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-----------------------
	return nil
}

//-------------------------------------------------
// PIPELINE__PROCESS_UPLOADED_IMAGE

func pipelineProcessUploadedImage(pImageIDstr gf_images_core.GFimageID,
	pS3filePathStr               string,
	pMetaMap                     map[string]interface{},
	pImagesStoreLocalDirPathStr  string,
	pImagesThumbsLocalDirPathStr string,
	pFlowsNamesLst               []string,
	// pSourceS3bucketNameStr       string, // S3_bucket to which the image was uploaded to
	// pTargetS3bucketNameStr       string, // S3 bucket to which processed images are stored in after this pipeline processing
	pS3info                      *gf_aws.GFs3Info,
	pPluginsPyDirPathStr         string,
	pStorage                     *gf_images_storage.GFimageStorage,
	pJobRuntime                  *GFjobRuntime,
	pRuntimeSys                  *gf_core.RuntimeSys) *gf_core.GFerror {

	//-----------------------
	// GET_IMAGE_FROM_FS

	imageLocalFilePathStr := fmt.Sprintf("%s/%s",
		pImagesStoreLocalDirPathStr,
		filepath.Base(pS3filePathStr))
	
	// NEW_STORAGE
	if pJobRuntime.useNewStorageEngineBool {

		imageFilePathStr := pS3filePathStr
		op := &gf_images_storage.GFgetOpDef{
			ImageSourceFilePathStr:      imageFilePathStr,
			ImageTargetLocalFilePathStr: imageLocalFilePathStr,
			
		}
		if pStorage.TypeStr == "s3" {
			op.S3bucketNameStr = pStorage.S3.UploadsSourceS3bucketNameStr // pSourceS3bucketNameStr
		}
		gfErr := gf_images_storage.FileGet(op,
			pStorage,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

	} else {
		// LEGACY

		// S3_DOWNLOAD - of the uploaded user image.
		//               the client uploads images to s3 directly for efficiency reasons, to avoid having
		//               all external image upload traffic going through GF servers.
		gfErr := gf_images_core.S3getImage(pS3filePathStr,
			imageLocalFilePathStr,
			pStorage.S3.UploadsSourceS3bucketNameStr, // pSourceS3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gfErr != nil {
			error_type_str := "s3_download_for_processing_error"
			jobErrorSend(error_type_str, gfErr,
				"", // p_image_source_url_str,
				pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
			return gfErr
		}
	}

	//-----------------------
	// TRANSFORM_IMAGE
	
	imageThumbs, gfErr := jobTransform(pImageIDstr,
		pFlowsNamesLst,
		"", // p_image_source_url_str,
		"", // p_image_origin_page_url_str,
		pMetaMap,
		imageLocalFilePathStr,
		pImagesThumbsLocalDirPathStr,
		pPluginsPyDirPathStr,
		pJobRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-----------------------
	// SAVE_IMAGE_TO_FS
		
	// if the source and target S3 buckets are not the same for processing this image then
	// then copy this image from the source to the target bucket.
	// use the same image ID that is the name of the image.
	if pStorage.S3.UploadsSourceS3bucketNameStr != pStorage.S3.UploadsTargetS3bucketNameStr {

		// NEW_STORAGE
		if pJobRuntime.useNewStorageEngineBool {

			imageFilePathStr := pS3filePathStr
			op := &gf_images_storage.GFcopyOpDef{
				ImageSourceFilePathStr: imageFilePathStr,
				ImageTargetFilePathStr: imageFilePathStr,
			}
			if pStorage.TypeStr == "s3" {
				op.SourceFileS3bucketNameStr = pStorage.S3.UploadsSourceS3bucketNameStr
				op.TargetFileS3bucketNameStr = pStorage.S3.UploadsTargetS3bucketNameStr
			}
			gfErr := gf_images_storage.FileCopy(op,
				pStorage,
				pRuntimeSys)
			if gfErr != nil {
				error_type_str := "s3_store_error"
				jobErrorSend(error_type_str, gfErr,
					"", // p_image_source_url_str,
					pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
					
				return gfErr
			}

		} else {
			// LEGACY

			// S3_FILE_COPY
			gfErr := gf_aws.S3copyFile(pStorage.S3.UploadsSourceS3bucketNameStr, // pSourceS3bucketNameStr,
				pS3filePathStr,
				pStorage.S3.UploadsTargetS3bucketNameStr, // pTargetS3bucketNameStr,
				pS3filePathStr,
				pS3info,
				pRuntimeSys)
			if gfErr != nil {
				
				error_type_str := "s3_store_error"
				jobErrorSend(error_type_str, gfErr,
					"", // p_image_source_url_str,
					pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
				return gfErr
			}
		}
		
	}

	// NEW_STORAGE
	if pJobRuntime.useNewStorageEngineBool {


		gfErr := gf_images_core.StoreThumbnails(imageThumbs, pStorage, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

	} else {
		// LEGACY

		// STORE__IMAGE_THUMBS
		gfErr := gf_images_core.S3storeThumbnails(imageThumbs,
			pStorage.S3.UploadsTargetS3bucketNameStr, // pTargetS3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gfErr != nil {
			error_type_str := "s3_store_error"
			jobErrorSend(error_type_str, gfErr,
				"", // p_image_source_url_str,
				pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
			return gfErr
		}
	}

	

	update_msg := JobUpdateMsg{
		Name_str:             "image_persist",
		Type_str:             JOB_UPDATE_TYPE__OK,
		ImageIDstr:           pImageIDstr,
		Image_source_url_str: "", // p_image_source_url_str,
	}
	pJobRuntime.job_updates_ch <- update_msg

	//-----------------------
	// DONE
	update_msg = JobUpdateMsg{
		Name_str:             "image_done",
		Type_str:             JOB_UPDATE_TYPE__COMPLETED,
		ImageIDstr:           pImageIDstr,
		Image_source_url_str: "", // p_image_source_url_str,
		Image_thumbs:         imageThumbs,
	}
	pJobRuntime.job_updates_ch <- update_msg

	//-----------------------

	return nil
}

//-------------------------------------------------
// PIPELINE__PROCESS_EXTERN_IMAGE

func pipelineProcessExternImage(pImageIDstr gf_images_core.GFimageID,
	pImageSourceURLstr           string,
	pOriginPageURLstr            string,
	pImagesStoreLocalDirPathStr  string,
	pImagesThumbsLocalDirPathStr string,
	pFlowsNamesLst               []string,
	pS3bucketNameStr             string,
	pS3info                      *gf_aws.GFs3Info,
	pPluginsPyDirPathStr         string,
	pStorage                     *gf_images_storage.GFimageStorage,
	pJobRuntime                  *GFjobRuntime,
	pRuntimeSys                  *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//-----------------------
	// FETCH_IMAGE
	imageLocalFilePathStr, _, gfErr := gf_images_core.FetcherGetExternImage(pImageSourceURLstr,
		pImagesStoreLocalDirPathStr,
		false, // p_random_time_delay_bool
		pRuntimeSys)
	if gfErr != nil {
		error_type_str := "fetch_error"
		jobErrorSend(error_type_str, gfErr, pImageSourceURLstr, pImageIDstr, 
			pJobRuntime.job_id_str,
			pJobRuntime.job_updates_ch, pRuntimeSys)
		return gfErr
	}

	updateMsg := JobUpdateMsg{
		Name_str:             "image_fetch",
		Type_str:             JOB_UPDATE_TYPE__OK,
		ImageIDstr:           pImageIDstr,
		Image_source_url_str: pImageSourceURLstr,
	}

	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	// TRANSFORM_IMAGE
	
	// FIX!! - this should be passed in from outside this function
	metaMap := map[string]interface{}{}

	imageThumbs, gfErr := jobTransform(pImageIDstr,
		pFlowsNamesLst,
		pImageSourceURLstr,
		pOriginPageURLstr,
		metaMap,
		imageLocalFilePathStr,
		pImagesThumbsLocalDirPathStr,
		pPluginsPyDirPathStr,
		pJobRuntime,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-----------------------
	// SAVE_IMAGE TO FS (S3)

	// NEW_STORAGE
	if pJobRuntime.useNewStorageEngineBool {

		//--------------------
		// IMAGE_FULLSIZE_STORE

		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// FIX!! - target filename of the original image should not be its original file name (that might collide accross domains or other images),
		//         and instead should be the image ID with the file extension. 
		//         it also makes it more difficult to find the image on S3 that is represented by an Gf_img given 
		//         only the ID of that Gf_img
		fileNameStr := path.Base(imageLocalFilePathStr)
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

		/* for files acquired by the Fetcher images are already uploaded 
		with their Gf_img ID as their filename. so here the pImageLocalFilePathStr value is already 
		the image ID.
		
		ADD!! - have an explicit p_target_fileNameStr argument, and dont derive it
				automatically from the the filename in pImageLocalFilePathStr */
		targetFilePathStr := fileNameStr

		op := &gf_images_storage.GFputFromLocalOpDef{
			ImageSourceLocalFilePathStr: imageLocalFilePathStr,
			ImageTargetFilePathStr:      targetFilePathStr,
		}
		if pStorage.TypeStr == "s3" {
			op.S3bucketNameStr = pStorage.S3.ExternImagesS3bucketNameStr
		}
		gfErr := gf_images_storage.FilePutFromLocal(op, pStorage, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		//--------------------
		// THUMBS_STORE

		gfErr = gf_images_core.StoreThumbnails(imageThumbs,
			pStorage,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		//--------------------
		
	} else {
		// LEGACY
	
		gfErr := gf_images_core.S3storeImage(imageLocalFilePathStr,
			imageThumbs,
			pS3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gfErr != nil {
			error_type_str := "s3_store_error"
			jobErrorSend(error_type_str, gfErr, pImageSourceURLstr, pImageIDstr,
				pJobRuntime.job_id_str,
				pJobRuntime.job_updates_ch,
				pRuntimeSys)
			return gfErr
		}
	}

	updateMsg = JobUpdateMsg{
		Name_str:             "image_persist",
		Type_str:             JOB_UPDATE_TYPE__OK,
		ImageIDstr:           pImageIDstr,
		Image_source_url_str: pImageSourceURLstr,
	}
	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	// DONE
	updateMsg = JobUpdateMsg{
		Name_str:             "image_done",
		Type_str:             JOB_UPDATE_TYPE__COMPLETED,
		ImageIDstr:           pImageIDstr,
		Image_source_url_str: pImageSourceURLstr,
		Image_thumbs:         imageThumbs,
	}
	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	return nil
}

//-------------------------------------------------
// PIPELINE__PROCESS_LOCAL_IMAGE

func pipelineProcessLocalImage(pFlowsNamesLst []string,
	pS3info              *gf_aws.GFs3Info,
	pPluginsPyDirPathStr string,
	pStorage             *gf_images_storage.GFimageStorage,
	pRuntimeSys          *gf_core.RuntimeSys) *gf_core.GFerror {



	return nil
}

//-------------------------------------------------

func jobTransform(pImageIDstr gf_images_core.GFimageID,
	pFlowsNamesLst               []string,
	pImageSourceURLstr           string,
	pImageOriginPageURLstr       string,
	pMetaMap                     map[string]interface{},
	pImageLocalFilePathStr       string,
	pImagesThumbsLocalDirPathStr string,
	pPluginsPyDirPathStr         string,
	pJobRuntime                  *GFjobRuntime,
	pRuntimeSys                  *gf_core.RuntimeSys) (*gf_images_core.GFimageThumbs, *gf_core.GFerror) {

	// TRANSFORM
	ctx := context.Background()

	_, imageThumbs, gfErr := gf_images_core.TransformImage(pImageIDstr,
		pJobRuntime.job_client_type_str,
		pFlowsNamesLst,
		pImageSourceURLstr,
		pImageOriginPageURLstr,
		pMetaMap,
		pImageLocalFilePathStr,
		pImagesThumbsLocalDirPathStr,
		pPluginsPyDirPathStr,
		pJobRuntime.metricsPlugins,
		ctx,
		pRuntimeSys)

	if gfErr != nil {
		error_type_str := "transform_error"
		jobErrorSend(error_type_str, gfErr,
			pImageSourceURLstr,
			pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
		return nil, gfErr
	}

	updateMsg := JobUpdateMsg{
		Name_str:             "image_transform",
		Type_str:             JOB_UPDATE_TYPE__OK,
		ImageIDstr:           pImageIDstr,
		Image_source_url_str: pImageSourceURLstr,
	}
	pJobRuntime.job_updates_ch <- updateMsg



	return imageThumbs, nil
}