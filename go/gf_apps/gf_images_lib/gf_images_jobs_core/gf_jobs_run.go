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
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
)

//-------------------------------------------------
func run_job__local_imgs(pImagesToProcessLst []GF_image_local_to_process,
	pFlowsNamesLst                        []string,
	p_images_store_local_dir_path_str     string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pS3info                               *gf_core.GFs3Info,
	pStorage                              *gf_images_storage.GFimageStorage,
	pJobRuntime                           *GFjobRuntime,
	pRuntimeSys                           *gf_core.Runtime_sys) []*gf_core.GFerror {
	



	gf_errors_lst := []*gf_core.GFerror{}
	for _, imageToProcess := range pImagesToProcessLst {
		

		fmt.Println(imageToProcess)



		gf_err := pipelineProcessLocalImage(pFlowsNamesLst,
			pS3info,
			pStorage,
			pRuntimeSys)
		if gf_err != nil {
			gf_errors_lst = append(gf_errors_lst, gf_err)
		}
	}

	return gf_errors_lst
}

//-------------------------------------------------
func run_job__uploaded_imgs(pImagesToProcessLst []GF_image_uploaded_to_process,
	pFlowsNamesLst                        []string,
	p_images_store_local_dir_path_str     string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	p_source_s3_bucket_name_str           string, // S3_bucket to which the image was uploaded to
	p_target_s3_bucket_name_str           string, // S3 bucket to which processed images are stored in after this pipeline processing
	pS3info                               *gf_core.GFs3Info,
	pStorage                              *gf_images_storage.GFimageStorage,
	pJobRuntime                           *GFjobRuntime,
	pRuntimeSys                           *gf_core.RuntimeSys) []*gf_core.GFerror {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_jobs_run.run_job__uploaded_imgs()")

	gfErrorsLst := []*gf_core.GF_error{}
	for _, imageToProcess := range pImagesToProcessLst {

		fmt.Println(imageToProcess)

		gf_image_id_str  := imageToProcess.GF_image_id_str
		s3_file_path_str := imageToProcess.S3_file_path_str
		meta_map         := imageToProcess.Meta_map

		gfErr := pipelineProcessUploadedImage(gf_image_id_str,
			s3_file_path_str,
			meta_map,
			p_images_store_local_dir_path_str,
			pImagesThumbnailsStoreLocalDirPathStr,
			pFlowsNamesLst,
			p_source_s3_bucket_name_str,
			p_target_s3_bucket_name_str,
			pS3info,
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
func run_job__extern_imgs(pImagesToProcessLst []GF_image_extern_to_process,
	pFlowsNamesLst                        []string,
	p_images_store_local_dir_path_str     string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pMediaDomainStr                       string,
	pS3bucketNameStr                      string,
	pS3info                               *gf_core.GFs3Info,
	pStorage                              *gf_images_storage.GFimageStorage,
	pJobRuntime                           *GFjobRuntime,
	pRuntimeSys                           *gf_core.RuntimeSys) []*gf_core.GFerror {

	ctx := context.Background()

	gf_errors_lst := []*gf_core.GF_error{}
	for _, imageToProcess := range pImagesToProcessLst {

		image_source_url_str      := imageToProcess.Source_url_str // FIX!! rename source_url_str to origin_url_str
		image_origin_page_url_str := imageToProcess.Origin_page_url_str

		//--------------
		// IMAGE_ID
		image_id_str, i_gf_err := gf_images_core.Image_ID__create_from_url(image_source_url_str, pRuntimeSys)

		if i_gf_err != nil {
			job_error_type_str := "create_image_id_error"
			_ = job_error__send(job_error_type_str, i_gf_err, image_source_url_str,
				image_id_str,
				pJobRuntime.job_id_str,
				pJobRuntime.job_updates_ch,
				pRuntimeSys)
			gf_errors_lst = append(gf_errors_lst, i_gf_err)
			continue
		}
		
		//--------------

		pRuntimeSys.Log_fun("INFO", "PROCESSING IMAGE - "+image_source_url_str)

		// IMPORTANT!! - 'ok' is '_' because Im already calling Get_image_ext_from_url()
		//               in Image__create_id_from_url()
		ext_str, ext_gf_err := gf_images_core.Get_image_ext_from_url(image_source_url_str, pRuntimeSys)
		
		if ext_gf_err != nil {
			job_error_type_str := "get_image_ext_error"
			_ = job_error__send(job_error_type_str, ext_gf_err, image_source_url_str, image_id_str,
				pJobRuntime.job_id_str,
				pJobRuntime.job_updates_ch,
				pRuntimeSys)
			gf_errors_lst = append(gf_errors_lst, ext_gf_err)
			continue
		}

		//--------------
		// GIF - gifs have their own processing pipeline

		// FIX!! - move GIF processing logic into gf_images_pipeline.go as well,
		//         it doesnt belong here in general images_job logic.

		if ext_str == "gif" {

			//-----------------
			// FLOWS_NAMES
			// check if "gifs" flow is already in the list
			b := false
			for _, s := range pFlowsNamesLst {
				if s == "gifs" {
					b = true
				}
			}
			
			var flows_names_lst []string
			if b {
				flows_names_lst = append([]string{"gifs"}, pFlowsNamesLst...)
			} else {
				flows_names_lst = pFlowsNamesLst
			}

			//-----------------

			_, gfErr := gf_gif_lib.ProcessAndUpload("", // p_gf_image_id_str
				image_source_url_str,
				image_origin_page_url_str,
				p_images_store_local_dir_path_str,
				pJobRuntime.job_client_type_str,
				flows_names_lst,
				true, // p_create_new_db_img_bool

				pMediaDomainStr,
				pS3bucketNameStr,
				pS3info,
				ctx,
				pRuntimeSys)

			if gfErr != nil {
				job_error_type_str := "gif_process_and_upload_error"
				_ = job_error__send(job_error_type_str, gfErr, image_source_url_str, image_id_str,
					pJobRuntime.job_id_str,
					pJobRuntime.job_updates_ch,
					pRuntimeSys)
				gf_errors_lst = append(gf_errors_lst, gfErr)
				continue
			}

			continue

		//-----------------------
		// STANDARD
		} else {
			gfErr := pipelineProcessExternImage(image_id_str,
				image_source_url_str,
				image_origin_page_url_str,
				p_images_store_local_dir_path_str,
				pImagesThumbnailsStoreLocalDirPathStr,
				pFlowsNamesLst,
				pS3bucketNameStr,
				pS3info,
				pJobRuntime,
				pRuntimeSys)

			if gfErr != nil {
				jobErrorTypeStr := "image_process_error"
				_ = job_error__send(jobErrorTypeStr, gfErr, image_source_url_str, image_id_str,
					pJobRuntime.job_id_str,
					pJobRuntime.job_updates_ch,
					pRuntimeSys)
				gf_errors_lst = append(gf_errors_lst, gfErr)
				continue
			}
		}

		//-----------------------
	}
	return gf_errors_lst
}