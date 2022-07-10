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
	"path/filepath"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_storage"
)

//-------------------------------------------------
// PIPELINE__PROCESS_LOCAL_IMAGE
func pipelineProcessLocalImage(pFlowsNamesLst []string,
	pS3info     *gf_core.GFs3Info,
	pStorage    *gf_images_storage.GFimageStorage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {



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
	pSourceS3bucketNameStr       string, // S3_bucket to which the image was uploaded to
	pTargetS3bucketNameStr       string, // S3 bucket to which processed images are stored in after this pipeline processing
	pS3info                      *gf_core.GFs3Info,
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

		op := &gf_images_storage.GFgetOpDef{
			ImageLocalFilePathStr: imageLocalFilePathStr,
			TargetFilePathStr:     pS3filePathStr,
			S3bucketNameStr:       pSourceS3bucketNameStr,
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
			pSourceS3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gfErr != nil {
			error_type_str := "s3_download_for_processing_error"
			job_error__send(error_type_str, gfErr,
				"", // p_image_source_url_str,
				pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
			return gfErr
		}
	}

	//-----------------------
	// TRANSFORM_IMAGE
	
	imageThumbs, gfTerr := jobTransform(pImageIDstr,
		pFlowsNamesLst,
		"", // p_image_source_url_str,
		"", // p_image_origin_page_url_str,
		pMetaMap,
		imageLocalFilePathStr,
		pImagesThumbsLocalDirPathStr, 
		pJobRuntime,
		pRuntimeSys)
	if gfTerr != nil {
		return gfTerr
	}

	//-----------------------
	// SAVE_IMAGE_TO_FS
	
	// NEW_STORAGE
	if pJobRuntime.useNewStorageEngineBool {


		op := &gf_images_storage.GFcopyOpDef{
			SourceFilePathStr:         pS3filePathStr,
			SourceFileS3bucketNameStr: pSourceS3bucketNameStr,
			TargetFilePathStr:         pS3filePathStr,
			TargetFileS3bucketNameStr: pTargetS3bucketNameStr,
		}
		gfErr := gf_images_storage.FileCopy(op,
			pStorage,
			pRuntimeSys)
		if gfErr != nil {
			error_type_str := "s3_store_error"
			job_error__send(error_type_str, gfErr,
				"", // p_image_source_url_str,
				pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
				
			return gfErr
		}

	} else {
		// LEGACY
	
		// if the source and target S3 buckets are not the same for processing this image then
		// then copy this image from the source to the target bucket.
		// use the same image ID that is the name of the image.
		if pSourceS3bucketNameStr != pTargetS3bucketNameStr {

			// S3_FILE_COPY
			gfErr := gf_core.S3copyFile(pSourceS3bucketNameStr,
				pS3filePathStr,
				pTargetS3bucketNameStr,
				pS3filePathStr,
				pS3info,
				pRuntimeSys)
			if gfErr != nil {
				
				error_type_str := "s3_store_error"
				job_error__send(error_type_str, gfErr,
					"", // p_image_source_url_str,
					pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
				return gfErr
			}
		}

		// STORE__IMAGE_THUMBS
		gfErr := gf_images_core.S3storeThumbnails(imageThumbs,
			pTargetS3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gfErr != nil {
			error_type_str := "s3_store_error"
			job_error__send(error_type_str, gfErr,
				"", // p_image_source_url_str,
				pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
			return gfErr
		}
	}

	update_msg := JobUpdateMsg{
		Name_str:             "image_persist",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         pImageIDstr,
		Image_source_url_str: "", // p_image_source_url_str,
	}
	pJobRuntime.job_updates_ch <- update_msg

	//-----------------------
	// DONE
	update_msg = JobUpdateMsg{
		Name_str:             "image_done",
		Type_str:             JOB_UPDATE_TYPE__COMPLETED,
		Image_id_str:         pImageIDstr,
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
	p_image_source_url_str       string,
	p_image_origin_page_url_str  string,
	pImagesStoreLocalDirPathStr  string,
	pImagesThumbsLocalDirPathStr string,
	pFlowsNamesLst               []string,
	pS3bucketNameStr             string,
	pS3info                      *gf_core.GFs3Info,
	pJobRuntime                  *GFjobRuntime,
	pRuntimeSys                  *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_jobs_pipeline.job__pipeline__process_image_extern()")
	
	//-----------------------
	// FETCH_IMAGE
	image_local_file_path_str, _, gf_f_err := gf_images_core.Fetcher__get_extern_image(p_image_source_url_str,
		pImagesStoreLocalDirPathStr,
		false, // p_random_time_delay_bool
		pRuntimeSys)
	if gf_f_err != nil {
		error_type_str := "fetch_error"
		job_error__send(error_type_str, gf_f_err, p_image_source_url_str, pImageIDstr, 
			pJobRuntime.job_id_str,
			pJobRuntime.job_updates_ch, pRuntimeSys)
		return gf_f_err
	}

	updateMsg := JobUpdateMsg{
		Name_str:             "image_fetch",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         pImageIDstr,
		Image_source_url_str: p_image_source_url_str,
	}

	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	// TRANSFORM_IMAGE
	
	// FIX!! - this should be passed it from outside this function
	metaMap := map[string]interface{}{}

	gf_image_thumbs, gf_t_err := jobTransform(pImageIDstr,
		pFlowsNamesLst,
		p_image_source_url_str,
		p_image_origin_page_url_str,
		metaMap,
		image_local_file_path_str,
		pImagesThumbsLocalDirPathStr,
		pJobRuntime,
		pRuntimeSys)
	if gf_t_err != nil {
		return gf_t_err
	}

	//-----------------------
	// SAVE_IMAGE TO FS (S3)

	// NEW_STORAGE
	if pJobRuntime.useNewStorageEngineBool {

	} else {
		// LEGACY
	
		gf_s3_err := gf_images_core.S3storeImage(image_local_file_path_str,
			gf_image_thumbs,
			pS3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gf_s3_err != nil {
			error_type_str := "s3_store_error"
			job_error__send(error_type_str, gf_s3_err, p_image_source_url_str, pImageIDstr,
				pJobRuntime.job_id_str,
				pJobRuntime.job_updates_ch,
				pRuntimeSys)
			return gf_s3_err
		}
	}

	updateMsg = JobUpdateMsg{
		Name_str:             "image_persist",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         pImageIDstr,
		Image_source_url_str: p_image_source_url_str,
	}
	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	// DONE
	updateMsg = JobUpdateMsg{
		Name_str:             "image_done",
		Type_str:             JOB_UPDATE_TYPE__COMPLETED,
		Image_id_str:         pImageIDstr,
		Image_source_url_str: p_image_source_url_str,
		Image_thumbs:         gf_image_thumbs,
	}
	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	return nil
}

//-------------------------------------------------
func jobTransform(pImageIDstr gf_images_core.GFimageID,
	pFlowsNamesLst               []string,
	p_image_source_url_str       string,
	p_image_origin_page_url_str  string,
	p_meta_map                   map[string]interface{},
	pImageLocalFilePathStr       string,
	pImagesThumbsLocalDirPathStr string,
	pJobRuntime                  *GFjobRuntime,
	pRuntimeSys                  *gf_core.RuntimeSys) (*gf_images_core.GF_image_thumbs, *gf_core.GFerror) {

	// TRANSFORM
	ctx := context.Background()

	_, gfImageThumbs, gfTerr := gf_images_core.TransformImage(pImageIDstr,
		pJobRuntime.job_client_type_str,
		pFlowsNamesLst,
		p_image_source_url_str,
		p_image_origin_page_url_str,
		p_meta_map,
		pImageLocalFilePathStr,
		pImagesThumbsLocalDirPathStr,
		ctx,
		pRuntimeSys)

	if gfTerr != nil {
		error_type_str := "transform_error"
		job_error__send(error_type_str, gfTerr,
			p_image_source_url_str,
			pImageIDstr, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
		return nil, gfTerr
	}

	update_msg := JobUpdateMsg{
		Name_str:             "image_transform",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         pImageIDstr,
		Image_source_url_str: p_image_source_url_str,
	}
	pJobRuntime.job_updates_ch <- update_msg



	return gfImageThumbs, nil
}