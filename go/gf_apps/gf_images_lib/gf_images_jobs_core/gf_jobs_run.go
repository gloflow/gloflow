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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
)

//-------------------------------------------------
type GF_job_run__runtime struct {
	job_id_str          string
	job_client_type_str string
	job_updates_ch      chan Job_update_msg
	s3_info             *gf_core.GF_s3_info
}

//-------------------------------------------------
func run_job__local_imgs(p_images_to_process_lst []GF_image_local_to_process,

	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,

	p_target_s3_bucket_name_str string, // S3 bucket to which processed images are stored in after this pipeline processing
	p_job_run_runtime           *GF_job_run__runtime,
	p_runtime_sys               *gf_core.Runtime_sys) []*gf_core.GF_error {
	



	gf_errors_lst := []*gf_core.GF_error{}
	for _, image_to_process := range p_images_to_process_lst {
		

		fmt.Println(image_to_process)
	}

	return gf_errors_lst
}

//-------------------------------------------------
func run_job__uploaded_imgs(p_images_to_process_lst []GF_image_uploaded_to_process,
	p_flows_names_lst                            []string,
	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_source_s3_bucket_name_str                  string, // S3_bucket to which the image was uploaded to
	p_target_s3_bucket_name_str                  string, // S3 bucket to which processed images are stored in after this pipeline processing
	p_job_run_runtime                            *GF_job_run__runtime,
	p_runtime_sys                                *gf_core.Runtime_sys) []*gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_run.run_job__uploaded_imgs()")

	gf_errors_lst := []*gf_core.GF_error{}
	for _, image_to_process := range p_images_to_process_lst {

		fmt.Println(image_to_process)

		gf_image_id_str  := image_to_process.Gf_image_id_str
		s3_file_path_str := image_to_process.S3_file_path_str

		gf_err := job__pipeline__process_image_uploaded(gf_image_id_str,
			s3_file_path_str,
			p_images_store_local_dir_path_str,
			p_images_thumbnails_store_local_dir_path_str,
			p_flows_names_lst,
			p_job_run_runtime.job_id_str,
			p_job_run_runtime.job_client_type_str,
			p_job_run_runtime.job_updates_ch,
			p_source_s3_bucket_name_str,
			p_target_s3_bucket_name_str,
			p_job_run_runtime.s3_info,
			job_error__send,
			p_runtime_sys)
		if gf_err != nil {
			gf_errors_lst = append(gf_errors_lst, gf_err)
		}
	}

	return gf_errors_lst
}

//-------------------------------------------------
// RUN_JOB__EXTERN_IMAGES
func run_job__extern_imgs(p_images_to_process_lst []GF_image_extern_to_process,
	p_flows_names_lst                            []string,
	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_media_domain_str                           string,
	p_s3_bucket_name_str                         string,
	p_job_run_runtime                            *GF_job_run__runtime,
	p_runtime_sys                                *gf_core.Runtime_sys) []*gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_run.run_job__extern_imgs()")

	gf_errors_lst := []*gf_core.GF_error{}
	for _, image_to_process := range p_images_to_process_lst {

		image_source_url_str      := image_to_process.Source_url_str // FIX!! rename source_url_str to origin_url_str
		image_origin_page_url_str := image_to_process.Origin_page_url_str

		//--------------
		// IMAGE_ID
		image_id_str, i_gf_err := gf_images_utils.Image_ID__create_from_url(image_source_url_str, p_runtime_sys)

		if i_gf_err != nil {
			job_error_type_str := "create_image_id_error"
			_ = job_error__send(job_error_type_str, i_gf_err, image_source_url_str,
				image_id_str,
				p_job_run_runtime.job_id_str,
				p_job_run_runtime.job_updates_ch,
				p_runtime_sys)
			gf_errors_lst = append(gf_errors_lst, i_gf_err)
			continue
		}
		
		//--------------

		p_runtime_sys.Log_fun("INFO", "PROCESSING IMAGE - "+image_source_url_str)

		// IMPORTANT!! - 'ok' is '_' because Im already calling Get_image_ext_from_url()
		//               in Image__create_id_from_url()
		ext_str, ext_gf_err := gf_images_utils.Get_image_ext_from_url(image_source_url_str, p_runtime_sys)
		
		if ext_gf_err != nil {
			job_error_type_str := "get_image_ext_error"
			_ = job_error__send(job_error_type_str, ext_gf_err, image_source_url_str, image_id_str,
				p_job_run_runtime.job_id_str,
				p_job_run_runtime.job_updates_ch,
				p_runtime_sys)
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
			for _, s := range p_flows_names_lst {
				if s == "gifs" {
					b = true
				}
			}
			
			var flows_names_lst []string
			if b {
				flows_names_lst = append([]string{"gifs"}, p_flows_names_lst...)
			} else {
				flows_names_lst = p_flows_names_lst
			}

			//-----------------

			_, gf_err := gf_gif_lib.Process_and_upload("", // p_gf_image_id_str
				image_source_url_str,
				image_origin_page_url_str,
				p_images_store_local_dir_path_str,
				p_job_run_runtime.job_client_type_str,
				flows_names_lst,
				true, // p_create_new_db_img_bool

				p_media_domain_str,
				p_s3_bucket_name_str,
				p_job_run_runtime.s3_info,
				p_runtime_sys)

			if gf_err != nil {
				job_error_type_str := "gif_process_and_upload_error"
				_ = job_error__send(job_error_type_str, gf_err, image_source_url_str, image_id_str,
					p_job_run_runtime.job_id_str,
					p_job_run_runtime.job_updates_ch,
					p_runtime_sys)
				gf_errors_lst = append(gf_errors_lst, gf_err)
				continue
			}

			continue

		//-----------------------
		// STANDARD
		} else {
			gf_err := job__pipeline__process_image_extern(image_id_str,
				image_source_url_str,
				image_origin_page_url_str,
				p_images_store_local_dir_path_str,
				p_images_thumbnails_store_local_dir_path_str,
				p_flows_names_lst,
				p_job_run_runtime.job_id_str,
				p_job_run_runtime.job_client_type_str,
				p_job_run_runtime.job_updates_ch,
				p_s3_bucket_name_str,
				p_job_run_runtime.s3_info,
				job_error__send,
				p_runtime_sys)

			if gf_err != nil {
				job_error_type_str := "image_process_error"
				_ = job_error__send(job_error_type_str, gf_err, image_source_url_str, image_id_str,
					p_job_run_runtime.job_id_str,
					p_job_run_runtime.job_updates_ch,
					p_runtime_sys)
				gf_errors_lst = append(gf_errors_lst, gf_err)
				continue
			}
		}

		//-----------------------
	}
	return gf_errors_lst
}