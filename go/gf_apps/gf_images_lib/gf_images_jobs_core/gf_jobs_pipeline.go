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
// PIPELINE__PROCESS_IMAGE_LOCAL
func job__pipeline__process_image_local(pFlowsNamesLst []string,
	pS3info     *gf_core.GFs3Info,
	pStorage    *gf_images_storage.GFimageStorage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {



	return nil
	
}

//-------------------------------------------------
// PIPELINE__PROCESS_IMAGE_UPLOADED
func job__pipeline__process_image_uploaded(pImageIDstr gf_images_core.GFimageID,
	p_s3_file_path_str                 string,
	p_meta_map                         map[string]interface{},
	p_images_store_local_dir_path_str  string,
	p_images_thumbs_local_dir_path_str string,
	pFlowsNamesLst                     []string,
	pSourceS3bucketNameStr string, // S3_bucket to which the image was uploaded to
	pTargetS3bucketNameStr string, // S3 bucket to which processed images are stored in after this pipeline processing
	pS3info                *gf_core.GFs3Info,
	pStorage               *gf_images_storage.GFimageStorage,
	pJobRuntime            *GFjobRuntime,
	pRuntimeSys            *gf_core.RuntimeSys) *gf_core.GFerror {

	//-----------------------
	// GET_IMAGE_FROM_FS

	image_local_file_path_str := fmt.Sprintf("%s/%s",
		p_images_store_local_dir_path_str,
		filepath.Base(p_s3_file_path_str))
	
	// NEW_STORAGE
	if pJobRuntime.useNewStorageEngineBool {

		

	} else {
		// LEGACY

		// S3_DOWNLOAD - of the uploaded user image.
		//               the client uploads images to s3 directly for efficiency reasons, to avoid having
		//               all external image upload traffic going through GF servers.
		gfErr := gf_images_core.S3__get_gf_image(p_s3_file_path_str,
			image_local_file_path_str,
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
	
	gf_image_thumbs, gf_t_err := jobTransform(pImageIDstr,
		pFlowsNamesLst,
		"", // p_image_source_url_str,
		"", // p_image_origin_page_url_str,
		p_meta_map,
		image_local_file_path_str,
		p_images_thumbs_local_dir_path_str, 
		pJobRuntime,
		pRuntimeSys)
	if gf_t_err != nil {
		return gf_t_err
	}

	//-----------------------
	// SAVE_IMAGE_TO_FS
	
	// NEW_STORAGE
	if pJobRuntime.useNewStorageEngineBool {


		op := &gf_images_storage.GFcopyOpDef{
			SourceFilePathStr:         p_s3_file_path_str,
			SourceFileS3bucketNameStr: pSourceS3bucketNameStr,
			TargetFilePathStr:         p_s3_file_path_str,
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
				p_s3_file_path_str,
				pTargetS3bucketNameStr,
				p_s3_file_path_str,
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
		gfErr := gf_images_core.S3storeThumbnails(gf_image_thumbs,
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
		Image_thumbs:         gf_image_thumbs,
	}
	pJobRuntime.job_updates_ch <- update_msg

	//-----------------------

	return nil
}

//-------------------------------------------------
// PIPELINE__PROCESS_IMAGE_EXTERN
func job__pipeline__process_image_extern(p_image_id_str gf_images_core.GF_image_id,
	p_image_source_url_str             string,
	p_image_origin_page_url_str        string,
	p_images_store_local_dir_path_str  string,
	p_images_thumbs_local_dir_path_str string,
	pFlowsNamesLst                     []string,
	pS3bucketNameStr                   string,
	pS3info                            *gf_core.GFs3Info,
	pJobRuntime                        *GFjobRuntime,
	pRuntimeSys                        *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_jobs_pipeline.job__pipeline__process_image_extern()")
	
	//-----------------------
	// FETCH_IMAGE
	image_local_file_path_str, _, gf_f_err := gf_images_core.Fetcher__get_extern_image(p_image_source_url_str,
		p_images_store_local_dir_path_str,
		false, // p_random_time_delay_bool
		pRuntimeSys)
	if gf_f_err != nil {
		error_type_str := "fetch_error"
		job_error__send(error_type_str, gf_f_err, p_image_source_url_str, p_image_id_str, 
			pJobRuntime.job_id_str,
			pJobRuntime.job_updates_ch, pRuntimeSys)
		return gf_f_err
	}

	updateMsg := JobUpdateMsg{
		Name_str:             "image_fetch",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
	}

	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	// TRANSFORM_IMAGE
	
	// FIX!! - this should be passed it from outside this function
	metaMap := map[string]interface{}{}

	gf_image_thumbs, gf_t_err := jobTransform(p_image_id_str,
		pFlowsNamesLst,
		p_image_source_url_str,
		p_image_origin_page_url_str,
		metaMap,
		image_local_file_path_str,
		p_images_thumbs_local_dir_path_str,
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
			job_error__send(error_type_str, gf_s3_err, p_image_source_url_str, p_image_id_str,
				pJobRuntime.job_id_str,
				pJobRuntime.job_updates_ch,
				pRuntimeSys)
			return gf_s3_err
		}
	}

	updateMsg = JobUpdateMsg{
		Name_str:             "image_persist",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
	}
	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	// DONE
	updateMsg = JobUpdateMsg{
		Name_str:             "image_done",
		Type_str:             JOB_UPDATE_TYPE__COMPLETED,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
		Image_thumbs:         gf_image_thumbs,
	}
	pJobRuntime.job_updates_ch <- updateMsg

	//-----------------------
	return nil
}

//-------------------------------------------------
func jobTransform(p_image_id_str gf_images_core.GF_image_id,
	pFlowsNamesLst                     []string,
	p_image_source_url_str             string,
	p_image_origin_page_url_str        string,
	p_meta_map                         map[string]interface{},
	p_image_local_file_path_str        string,
	p_images_thumbs_local_dir_path_str string,
	pJobRuntime                        *GFjobRuntime,
	pRuntimeSys                        *gf_core.RuntimeSys) (*gf_images_core.GF_image_thumbs, *gf_core.GFerror) {

	// TRANSFORM
	ctx := context.Background()

	_, gfImageThumbs, gfTerr := gf_images_core.TransformImage(p_image_id_str,
		pJobRuntime.job_client_type_str,
		pFlowsNamesLst,
		p_image_source_url_str,
		p_image_origin_page_url_str,
		p_meta_map,
		p_image_local_file_path_str,
		p_images_thumbs_local_dir_path_str,
		ctx,
		pRuntimeSys)

	if gfTerr != nil {
		error_type_str := "transform_error"
		job_error__send(error_type_str, gfTerr,
			p_image_source_url_str,
			p_image_id_str, pJobRuntime.job_id_str, pJobRuntime.job_updates_ch, pRuntimeSys)
		return nil, gfTerr
	}

	update_msg := JobUpdateMsg{
		Name_str:             "image_transform",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
	}
	pJobRuntime.job_updates_ch <- update_msg



	return gfImageThumbs, nil
}