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
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_gif_lib"
)
//-------------------------------------------------
func jobs_mngr__run_job(p_job_id_str string,
	p_job_client_type_str                        string,
	p_images_to_process_lst                      []Image_to_process,
	p_flows_names_lst                            []string,
	p_job_updates_ch                             chan *Job_update_msg,
	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_s3_bucket_name_str                         string,
	p_s3_info                                    *gf_core.Gf_s3_info,
	p_runtime_sys                                *gf_core.Runtime_sys) []*gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_job_run.jobs_mngr__run_job()")

	gf_errors_lst := []*gf_core.Gf_error{}
	for _,image_to_process := range p_images_to_process_lst {

		image_source_url_str      := image_to_process.Source_url_str
		image_origin_page_url_str := image_to_process.Origin_page_url_str

		//--------------
		//IMAGE_ID
		image_id_str,i_gf_err := gf_images_utils.Image__create_id_from_url(image_source_url_str,p_runtime_sys)

		if i_gf_err != nil {
			job_error_type_str := "create_image_id_error"
			_ = job_error__send(job_error_type_str,i_gf_err,image_source_url_str,image_id_str,p_job_id_str,p_job_updates_ch,p_runtime_sys)
			gf_errors_lst = append(gf_errors_lst,i_gf_err)
			continue
		}
		//--------------

		p_runtime_sys.Log_fun("INFO","PROCESSING IMAGE - "+image_source_url_str)

		//IMPORTANT!! - 'ok' is '_' because Im already calling Get_image_ext_from_url()
		//              in Image__create_id_from_url()
		ext_str,ext_gf_err := gf_images_utils.Get_image_ext_from_url(image_source_url_str,p_runtime_sys)
		
		if ext_gf_err != nil {
			job_error_type_str := "get_image_ext_error"
			_ = job_error__send(job_error_type_str,ext_gf_err,image_source_url_str,image_id_str,p_job_id_str,p_job_updates_ch,p_runtime_sys)
			gf_errors_lst = append(gf_errors_lst,ext_gf_err)
			continue
		}

		//--------------
		//GIF - gifs have their own processing pipeline

		//FIX!! - move GIF processing logic into gf_images_pipeline.go as well, it doesnt belong here in general images_job logic

		if ext_str == "gif" {

			//-----------------
			//FLOWS_NAMES

			//check if "gifs" flow is already in the list
			b := false
			for _,s := range p_flows_names_lst {
				if s == "gifs" {
					b = true
				}
			}
			
			var flows_names_lst []string
			if b {
				flows_names_lst = append([]string{"gifs"},p_flows_names_lst...)
			} else {
				flows_names_lst = p_flows_names_lst
			}
			//-----------------

			_,gf_err := gf_gif_lib.Process_and_upload(image_source_url_str,
											image_origin_page_url_str,
											p_images_store_local_dir_path_str,
											p_job_client_type_str,
											flows_names_lst,
											true, //p_create_new_db_img_bool
											p_s3_bucket_name_str,
											p_s3_info,
											p_runtime_sys)
			if gf_err != nil {
				job_error_type_str := "gif_process_and_upload_error"
				_ = job_error__send(job_error_type_str,gf_err,image_source_url_str,image_id_str,p_job_id_str,p_job_updates_ch,p_runtime_sys)
				gf_errors_lst = append(gf_errors_lst,gf_err)
				continue
			}

			continue
		//-----------------------
		} else {
			gf_err := pipeline__process_image(image_source_url_str,
							image_id_str,
							image_origin_page_url_str,
							p_images_store_local_dir_path_str,
							p_images_thumbnails_store_local_dir_path_str,
							p_flows_names_lst,
							p_job_id_str,
							p_job_client_type_str,
							p_job_updates_ch,
							p_s3_bucket_name_str,
							p_s3_info,
							job_error__send,
							p_runtime_sys)
			if gf_err != nil {
				job_error_type_str := "image_process_error"
				_ = job_error__send(job_error_type_str,gf_err,image_source_url_str,image_id_str,p_job_id_str,p_job_updates_ch,p_runtime_sys)
				gf_errors_lst = append(gf_errors_lst,gf_err)
				continue
			}
		}
		//-----------------------
	}

	return gf_errors_lst
}
//-------------------------------------------------
func job_error__send(p_job_error_type_str string,
	p_gf_err               *gf_core.Gf_error,
	p_image_source_url_str string,
	p_image_id_str         string,
	p_job_id_str           string,
	p_job_updates_ch       chan *Job_update_msg,
	p_runtime_sys          *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_job_run.job_error__send()")

	p_runtime_sys.Log_fun("ERROR",fmt.Sprintf("fetching image failed - %s - %s",
							p_image_source_url_str,
							p_gf_err.Error))

	error_str  := fmt.Sprint(p_gf_err.Error)
	pje_gf_err := job_error__persist(p_job_id_str,
					p_job_error_type_str,
					error_str,
					p_image_source_url_str,
					p_runtime_sys)
	if pje_gf_err != nil {
		return pje_gf_err
	}
	//------------
	update_msg := &Job_update_msg{
		Type_str            :p_gf_err.Type_str,
		Image_id_str        :p_image_id_str,
		Image_source_url_str:p_image_source_url_str,
		Err_str             :error_str,
	}
	p_job_updates_ch <- update_msg
	//------------
	return nil
}
//-------------------------------------------------
func job_error__persist(p_job_id_str string,
	p_error_type_str       string,
	p_error_str            string,
	p_image_source_url_str string,
	p_runtime_sys          *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_job_run.job_error__persist()")

	job_error := Job_Error{
		Type_str            :p_error_type_str,
		Error_str           :p_error_str,
		Image_source_url_str:p_image_source_url_str,
	}

	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t_str": "img_running_job",
			"id_str":p_job_id_str,
		},
		bson.M{
			"$push":bson.M{"errors_lst":job_error,},
		})

	if err != nil {
		gf_err := gf_core.Error__create("failed to update img_running_job type document in mongodb, to add a job error",
			"mongodb_update_error",
			&map[string]interface{}{
				"job_id_str":          p_job_id_str,
				"error_type_str":      p_error_type_str,
				"error_str":           p_error_str,
				"image_source_url_str":p_image_source_url_str,
			},
			err,"gf_images_jobs",p_runtime_sys)
		return gf_err
	}

	return nil
}