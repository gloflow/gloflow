/*
GloFlow media management/publishing system
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

package gf_images_jobs

import (
	"gf_core"
	"apps/gf_images_lib/gf_images_utils"
)
//-------------------------------------------------
func pipeline__process_image(p_image_source_url_str string,
	p_image_id_str                               string,
	p_image_origin_page_url_str                  string,
	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_flows_names_lst                            []string,
	p_job_id_str                                 string,
	p_job_client_type_str                        string,
	p_job_updates_ch                  chan *Job_update_msg,
	p_s3_bucket_name_str              string,
	p_s3_info                         *gf_core.Gf_s3_info,
	p_send_error_fun                  func(string,*gf_core.Gf_error,string,string,string,chan *Job_update_msg,*gf_core.Runtime_sys) *gf_core.Gf_error,
	p_runtime_sys                     *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_pipeline.pipeline__process_image()")

	//-----------------------
	//FETCH_IMAGE
	local_image_file_path_str,f_err := gf_images_utils.Fetch_image(p_image_source_url_str,
															p_images_store_local_dir_path_str,
															p_runtime_sys)
	if f_err != nil {
		error_type_str := "fetch_error"
		p_send_error_fun(error_type_str,
					f_err,
					p_image_source_url_str,
					p_image_id_str,
					p_job_id_str,
					p_job_updates_ch,
					p_runtime_sys)
		return f_err
	}

	update_msg := &Job_update_msg{
		Type_str            :"fetch_ok",
		Image_id_str        :p_image_id_str,
		Image_source_url_str:p_image_source_url_str,
	}
	p_job_updates_ch <- update_msg
	//-----------------------
	//TRANSFORM_IMAGE
	
	image_client_type_str   := p_job_client_type_str
	_,gf_image_thumbs,t_err := gf_images_utils.Transform_image(p_image_id_str,
													image_client_type_str,
													p_flows_names_lst,
													p_image_source_url_str,
													p_image_origin_page_url_str,
													local_image_file_path_str,
													p_images_thumbnails_store_local_dir_path_str,
													p_runtime_sys)
	if t_err != nil {
		error_type_str := "transform_error"
		p_send_error_fun(error_type_str,
					t_err,
					p_image_source_url_str,
					p_image_id_str,
					p_job_id_str,
					p_job_updates_ch,
					p_runtime_sys)
		return t_err
	}

	update_msg = &Job_update_msg{
		Type_str:            "transform_ok",
		Image_id_str:        p_image_id_str,
		Image_source_url_str:p_image_source_url_str,
	}
	p_job_updates_ch <- update_msg
	//-----------------------
	//SAVE_IMAGE TO FS (S3)

	s3_err := gf_images_utils.S3__store_gf_image(local_image_file_path_str,
											gf_image_thumbs,
											p_s3_bucket_name_str,
											p_s3_info,
											p_runtime_sys)
	if s3_err != nil {
		error_type_str := "s3_store_error"
		p_send_error_fun(error_type_str,
					s3_err,
					p_image_source_url_str,
					p_image_id_str,
					p_job_id_str,
					p_job_updates_ch,
					p_runtime_sys)
		return s3_err
	}
	//-----------------------
	update_msg = &Job_update_msg{
		Type_str:            "image_done",
		Image_id_str:        p_image_id_str,
		Image_source_url_str:p_image_source_url_str,
		Image_thumbs:        gf_image_thumbs,
	}
	p_job_updates_ch <- update_msg
	//-----------------------
	return nil
}