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

package gf_images_jobs_client

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// called "expected" because jobs are long-running processes, and they might fail at various stages
// of their processing. in that case some of these result values will be satisfied, others will not.
type GF_job_expected_output struct {
	Image_id_str                      gf_images_core.GF_image_id `json:"image_id_str"`
	Image_source_url_str              string                     `json:"image_source_url_str"`
	Thumbnail_small_relative_url_str  string                     `json:"thumbnail_small_relative_url_str"`
	Thumbnail_medium_relative_url_str string                     `json:"thumbnail_medium_relative_url_str"`
	Thumbnail_large_relative_url_str  string                     `json:"thumbnail_large_relative_url_str"`
}

//-------------------------------------------------
// CLIENT
//-------------------------------------------------
func Run_local_imgs(p_client_type_str string,
	p_images_to_process_lst []gf_images_jobs_core.GF_image_local_to_process,
	p_flows_names_lst       []string,
	p_jobs_mngr_ch          gf_images_jobs_core.Jobs_mngr,
	p_runtime_sys           *gf_core.Runtime_sys) (*gf_images_jobs_core.GF_job_running, []*GF_job_expected_output, *gf_core.GF_error) {

	job_cmd_str    := "start_job_local_imgs"
	job_init_ch    := make(chan *gf_images_jobs_core.GF_job_running)
	job_updates_ch := make(chan gf_images_jobs_core.Job_update_msg, 10)
	
	job_msg := gf_images_jobs_core.Job_msg{
		Client_type_str:             p_client_type_str,
		Cmd_str:                     job_cmd_str,
		Job_init_ch:                 job_init_ch,
		Job_updates_ch:              job_updates_ch,
		Images_local_to_process_lst: p_images_to_process_lst,
		Flows_names_lst:             p_flows_names_lst,
	}

	// SEND_MSG
	p_jobs_mngr_ch <- job_msg

	// RECEIVE_MSG - get running_job info back from jobs_mngr
	running_job := <- job_init_ch

	


	//-----------------
	// JOB_EXPECTED_OUTPUT - its "expected" because results are not available yet (and might not
	//                       be available for some time), and yet we still want to have some of the expected
	//                       values so that other parts of the system can initialize in parallel with the job 
	//                       completing.

	imgs_local_paths_lst := []string{}
	for _, image_to_process := range p_images_to_process_lst {
		imgs_local_paths_lst = append(imgs_local_paths_lst, image_to_process.Local_file_path_str)
	}

	job_expected_outputs_lst, gf_err := job__get_expected_output(imgs_local_paths_lst, p_runtime_sys)
	if gf_err != nil {
		return nil, nil, gf_err
	}

	//-----------------

	return running_job, job_expected_outputs_lst, nil
}

//-------------------------------------------------
func Run_uploaded_imgs(p_client_type_str string,
	p_images_to_process_lst []gf_images_jobs_core.GF_image_uploaded_to_process,
	p_flows_names_lst       []string,
	p_jobs_mngr_ch          gf_images_jobs_core.Jobs_mngr,
	p_runtime_sys           *gf_core.Runtime_sys) (*gf_images_jobs_core.GF_job_running, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_client.Run_uploaded_imgs()")

	job_cmd_str    := "start_job_uploaded_imgs"
	job_init_ch    := make(chan *gf_images_jobs_core.GF_job_running)
	job_updates_ch := make(chan gf_images_jobs_core.Job_update_msg, 10)
	

	job_msg := gf_images_jobs_core.Job_msg{
		Client_type_str:                p_client_type_str,
		Cmd_str:                        job_cmd_str,
		Job_init_ch:                    job_init_ch,
		Job_updates_ch:                 job_updates_ch,
		Images_uploaded_to_process_lst: p_images_to_process_lst,
		Flows_names_lst:                p_flows_names_lst,
	}

	// SEND_MSG
	p_jobs_mngr_ch <- job_msg

	// RECEIVE_MSG - get running_job info back from jobs_mngr
	running_job := <- job_init_ch
	
	// spew.Dump(running_job)

	return running_job, nil
}

//-------------------------------------------------
// START
func Run_extern_imgs(p_client_type_str string,
	p_images_extern_to_process_lst []gf_images_jobs_core.GF_image_extern_to_process,
	p_flows_names_lst              []string,
	p_jobs_mngr_ch                 gf_images_jobs_core.Jobs_mngr,
	p_runtime_sys                  *gf_core.Runtime_sys) (*gf_images_jobs_core.GF_job_running, []*GF_job_expected_output, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_client.Run_extern_imgs()")
	p_runtime_sys.Log_fun("INFO",      "images_extern_to_process - "+fmt.Sprint(p_images_extern_to_process_lst))

	//-----------------
	// SEND_MSG_TO_JOBS_MNGR
	job_cmd_str    := "start_job"
	job_init_ch    := make(chan *gf_images_jobs_core.GF_job_running)
	job_updates_ch := make(chan gf_images_jobs_core.Job_update_msg, 10) // ADD!! - channel buffer size should be larger for large jobs (with a lot of images)

	job_msg := gf_images_jobs_core.Job_msg{
		Client_type_str:              p_client_type_str,
		Cmd_str:                      job_cmd_str,
		Job_init_ch:                  job_init_ch,
		Job_updates_ch:               job_updates_ch,
		Images_extern_to_process_lst: p_images_extern_to_process_lst,
		Flows_names_lst:              p_flows_names_lst,
	}

	// SEND_MSG
	p_jobs_mngr_ch <- job_msg

	// RECEIVE_MSG - get running_job info back from jobs_mngr
	running_job := <- job_init_ch

	


	//-----------------
	// JOB_EXPECTED_OUTPUT - its "expected" because results are not available yet (and might not
	//                       be available for some time), and yet we still want to have some of the expected
	//                       values so that other parts of the system can initialize in parallel with the job 
	//                       completing.

	imgs_source_urls_lst := []string{}
	for _, image_to_process := range p_images_extern_to_process_lst {
		imgs_source_urls_lst = append(imgs_source_urls_lst, image_to_process.Source_url_str)
	}

	job_expected_outputs_lst, gf_err := job__get_expected_output(imgs_source_urls_lst, p_runtime_sys)
	if gf_err != nil {
		return nil, nil, gf_err
	}

	//-----------------
	return running_job, job_expected_outputs_lst, nil
}

//-------------------------------------------------
func Job__get_update_ch(p_job_id_str string,
	p_jobs_mngr_ch gf_images_jobs_core.Jobs_mngr,
	p_runtime_sys  *gf_core.Runtime_sys) chan gf_images_jobs_core.Job_update_msg {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_client.Job__get_update_ch()")

	msg_response_ch := make(chan interface{})
	defer close(msg_response_ch)

	job_cmd_str := "get_job_update_ch"
	job_msg     := gf_images_jobs_core.Job_msg{
		Job_id_str:      p_job_id_str,
		Cmd_str:         job_cmd_str,
		Msg_response_ch: msg_response_ch,
	}

	p_jobs_mngr_ch <- job_msg

	response          := <-msg_response_ch
	job_updates_ch, _ := response.(chan gf_images_jobs_core.Job_update_msg)

	return job_updates_ch
}

//-------------------------------------------------
func Job__cleanup(p_job_id_str string,
	p_jobs_mngr_ch gf_images_jobs_core.Jobs_mngr,
	p_runtime_sys  *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_client.Job__cleanup()")

	job_cmd_str := "cleanup_job"
	job_msg     := gf_images_jobs_core.Job_msg{
		Job_id_str: p_job_id_str,
		Cmd_str:    job_cmd_str,
	}

	p_jobs_mngr_ch <- job_msg
}

//-------------------------------------------------
// VAR
//-------------------------------------------------
func job__get_expected_output(p_images_source_URIs_lst []string,
	p_runtime_sys *gf_core.Runtime_sys) ([]*GF_job_expected_output, *gf_core.GF_error) {

	

	job_expected_outputs_lst := []*GF_job_expected_output{}

	for _, img_source_url_str := range p_images_source_URIs_lst {

		//--------------
		// IMAGE_ID
		image_id_str, i_gf_err := gf_images_core.Image_ID__create_from_url(img_source_url_str, p_runtime_sys)
		if i_gf_err != nil {
			return nil, i_gf_err
		}

		//--------------
		// GET FILE_FORMAT
		normalized_ext_str, gf_err := gf_images_core.Get_image_ext_from_url(img_source_url_str, p_runtime_sys)
		
		// FIX!! - it should not fail the whole job if one image is invalid,
		//         it should continue and just mark that image with an error.
		if gf_err != nil {
			return nil, gf_err
		}

		//--------------

		output := &GF_job_expected_output{
			Image_id_str:                      image_id_str,
			Image_source_url_str:              img_source_url_str,
			Thumbnail_small_relative_url_str : fmt.Sprintf("/images/d/thumbnails/%s_thumb_small.%s",  image_id_str, normalized_ext_str),
			Thumbnail_medium_relative_url_str: fmt.Sprintf("/images/d/thumbnails/%s_thumb_medium.%s", image_id_str, normalized_ext_str),
			Thumbnail_large_relative_url_str:  fmt.Sprintf("/images/d/thumbnails/%s_thumb_large.%s",  image_id_str, normalized_ext_str),
		}
		job_expected_outputs_lst = append(job_expected_outputs_lst, output)
	}

	//-----------------

	return job_expected_outputs_lst, nil
}