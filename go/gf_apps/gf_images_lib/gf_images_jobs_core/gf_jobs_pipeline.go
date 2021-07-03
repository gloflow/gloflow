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
	"path/filepath"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)


//-------------------------------------------------
// PIPELINE__PROCESS_IMAGE_LOCAL
func job__pipeline__process_image_local() {




	
}

//-------------------------------------------------
// PIPELINE__PROCESS_IMAGE_UPLOADED
func job__pipeline__process_image_uploaded(p_image_id_str gf_images_utils.Gf_image_id,
	p_s3_file_path_str                           string,
	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_flows_names_lst           []string,
	p_job_id_str                string,
	p_job_client_type_str       string,
	p_job_updates_ch            chan Job_update_msg,
	p_source_s3_bucket_name_str string, // S3_bucket to which the image was uploaded to
	p_target_s3_bucket_name_str string, // S3 bucket to which processed images are stored in after this pipeline processing
	p_s3_info                   *gf_core.GF_s3_info,
	p_send_error_fun            func(string, *gf_core.Gf_error, string, gf_images_utils.Gf_image_id, string, chan Job_update_msg, *gf_core.Runtime_sys) *gf_core.Gf_error,
	p_runtime_sys               *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_pipeline.job__pipeline__process_image_uploaded()")

	//-----------------------
	// S3_DOWNLOAD - of the uploaded user image

	// normalized_format_str  := "png"
	// image_s3_file_path_str := gf_images_utils.S3__get_image_s3_filepath(p_image_id_str,
	// 	normalized_format_str,
	// 	p_runtime_sys)

	image_local_file_path_str := fmt.Sprintf("%s/%s", p_images_store_local_dir_path_str, filepath.Base(p_s3_file_path_str))

	gf_err := gf_images_utils.S3__get_gf_image(p_s3_file_path_str,
		image_local_file_path_str,
		p_source_s3_bucket_name_str,
		p_s3_info,
		p_runtime_sys)
	if gf_err != nil {
		error_type_str := "s3_download_for_processing_error"
		p_send_error_fun(error_type_str, gf_err,
			"", // p_image_source_url_str,
			p_image_id_str, p_job_id_str, p_job_updates_ch, p_runtime_sys)
		return gf_err
	}

	//-----------------------
	// TRANSFORM_IMAGE
	
	image_client_type_str := p_job_client_type_str

	_, gf_image_thumbs, gf_t_err := gf_images_utils.Transform_image(p_image_id_str,
		image_client_type_str,
		p_flows_names_lst,
		"", // p_image_source_url_str,
		"", // p_image_origin_page_url_str,
		image_local_file_path_str,
		p_images_thumbnails_store_local_dir_path_str,
		p_runtime_sys)

	if gf_t_err != nil {
		error_type_str := "transform_error"
		p_send_error_fun(error_type_str, gf_t_err,
			"", // p_image_source_url_str
			p_image_id_str, p_job_id_str, p_job_updates_ch, p_runtime_sys)
		return gf_t_err
	}

	update_msg := Job_update_msg{
		Name_str:     "image_transform",
		Type_str:     JOB_UPDATE_TYPE__OK,
		Image_id_str: p_image_id_str,
	}
	p_job_updates_ch <- update_msg

	//-----------------------
	// SAVE_IMAGE TO FS (S3)

	// if the source and target S3 buckets are not the same for processing this image then
	// then copy this image from the source to the target bucket.
	// use the same image ID that is the name of the image.
	if p_source_s3_bucket_name_str != p_target_s3_bucket_name_str {

		// S3_FILE_COPY
		gf_err := gf_core.S3__copy_file(p_source_s3_bucket_name_str,
			p_s3_file_path_str,
			p_target_s3_bucket_name_str,
			p_s3_file_path_str,
			p_s3_info,
			p_runtime_sys)
		if gf_err != nil {
			error_type_str := "s3_store_error"
			p_send_error_fun(error_type_str, gf_err,
				"", // p_image_source_url_str,
				p_image_id_str, p_job_id_str, p_job_updates_ch, p_runtime_sys)
			return gf_err
		}
	}

	// STORE__IMAGE_THUMBS
	gf_err = gf_images_utils.S3__store_gf_image_thumbs(gf_image_thumbs,
		p_target_s3_bucket_name_str,
		p_s3_info,
		p_runtime_sys)
	if gf_err != nil {
		error_type_str := "s3_store_error"
		p_send_error_fun(error_type_str, gf_err,
			"", // p_image_source_url_str,
			p_image_id_str, p_job_id_str, p_job_updates_ch, p_runtime_sys)
		return gf_err
	}

	update_msg = Job_update_msg{
		Name_str:             "image_persist",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: "", // p_image_source_url_str,
	}
	p_job_updates_ch <- update_msg

	//-----------------------
	// DONE
	update_msg = Job_update_msg{
		Name_str:             "image_done",
		Type_str:             JOB_UPDATE_TYPE__COMPLETED,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: "", // p_image_source_url_str,
		Image_thumbs:         gf_image_thumbs,
	}
	p_job_updates_ch <- update_msg

	//-----------------------

	return nil
}

//-------------------------------------------------
// PIPELINE__PROCESS_IMAGE_EXTERN
func job__pipeline__process_image_extern(p_image_id_str gf_images_utils.Gf_image_id,
	p_image_source_url_str                       string,
	p_image_origin_page_url_str                  string,
	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_flows_names_lst                            []string,
	p_job_id_str                                 string,
	p_job_client_type_str                        string,
	p_job_updates_ch                  chan Job_update_msg,
	p_s3_bucket_name_str              string,
	p_s3_info                         *gf_core.GF_s3_info,
	p_send_error_fun                  func(string, *gf_core.Gf_error, string, gf_images_utils.Gf_image_id, string, chan Job_update_msg, *gf_core.Runtime_sys) *gf_core.Gf_error,
	p_runtime_sys                     *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_pipeline.job__pipeline__process_image_extern()")

	//-----------------------
	// FETCH_IMAGE
	local_image_file_path_str, _, gf_f_err := gf_images_utils.Fetcher__get_extern_image(p_image_source_url_str,
		p_images_store_local_dir_path_str,
		false, // p_random_time_delay_bool
		p_runtime_sys)
	if gf_f_err != nil {
		error_type_str := "fetch_error"
		p_send_error_fun(error_type_str, gf_f_err, p_image_source_url_str, p_image_id_str, p_job_id_str, p_job_updates_ch, p_runtime_sys)
		return gf_f_err
	}

	update_msg := Job_update_msg{
		Name_str:             "image_fetch",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
	}

	p_job_updates_ch <- update_msg

	//-----------------------
	// TRANSFORM_IMAGE
	
	image_client_type_str        := p_job_client_type_str
	_, gf_image_thumbs, gf_t_err := gf_images_utils.Transform_image(p_image_id_str,
		image_client_type_str,
		p_flows_names_lst,
		p_image_source_url_str,
		p_image_origin_page_url_str,
		local_image_file_path_str,
		p_images_thumbnails_store_local_dir_path_str,
		p_runtime_sys)

	if gf_t_err != nil {
		error_type_str := "transform_error"
		p_send_error_fun(error_type_str, gf_t_err, p_image_source_url_str, p_image_id_str, p_job_id_str, p_job_updates_ch, p_runtime_sys)
		return gf_t_err
	}

	update_msg = Job_update_msg{
		Name_str:             "image_transform",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
	}
	p_job_updates_ch <- update_msg

	//-----------------------
	// SAVE_IMAGE TO FS (S3)

	gf_s3_err := gf_images_utils.S3__store_gf_image(local_image_file_path_str,
		gf_image_thumbs,
		p_s3_bucket_name_str,
		p_s3_info,
		p_runtime_sys)
	if gf_s3_err != nil {
		error_type_str := "s3_store_error"
		p_send_error_fun(error_type_str, gf_s3_err, p_image_source_url_str, p_image_id_str, p_job_id_str, p_job_updates_ch, p_runtime_sys)
		return gf_s3_err
	}

	update_msg = Job_update_msg{
		Name_str:             "image_persist",
		Type_str:             JOB_UPDATE_TYPE__OK,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
	}
	p_job_updates_ch <- update_msg

	//-----------------------
	// DONE
	update_msg = Job_update_msg{
		Name_str:             "image_done",
		Type_str:             JOB_UPDATE_TYPE__COMPLETED,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
		Image_thumbs:         gf_image_thumbs,
	}
	p_job_updates_ch <- update_msg

	//-----------------------
	return nil
}