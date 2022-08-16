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
	Image_id_str                      gf_images_core.GFimageID `json:"image_id_str"`
	Image_source_url_str              string                   `json:"image_source_url_str"`
	Thumbnail_small_relative_url_str  string                   `json:"thumbnail_small_relative_url_str"`
	Thumbnail_medium_relative_url_str string                   `json:"thumbnail_medium_relative_url_str"`
	Thumbnail_large_relative_url_str  string                   `json:"thumbnail_large_relative_url_str"`
}

//-------------------------------------------------
// CLIENT
//-------------------------------------------------
func RunLocalImgs(pClientTypeStr string,
	p_images_to_process_lst []gf_images_jobs_core.GF_image_local_to_process,
	p_flows_names_lst       []string,
	pJobsMngrCh             gf_images_jobs_core.JobsMngr,
	pRuntimeSys             *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, []*GF_job_expected_output, *gf_core.GFerror) {

	job_cmd_str    := "start_job_local_imgs"
	job_init_ch    := make(chan *gf_images_jobs_core.GFjobRunning)
	job_updates_ch := make(chan gf_images_jobs_core.JobUpdateMsg, 10)
	
	job_msg := gf_images_jobs_core.JobMsg{
		Client_type_str:             pClientTypeStr,
		Cmd_str:                     job_cmd_str,
		Job_init_ch:                 job_init_ch,
		Job_updates_ch:              job_updates_ch,
		Images_local_to_process_lst: p_images_to_process_lst,
		Flows_names_lst:             p_flows_names_lst,
	}

	// SEND_MSG
	pJobsMngrCh <- job_msg

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

	job_expected_outputs_lst, gf_err := getJobExpectedOutput(imgs_local_paths_lst, pRuntimeSys)
	if gf_err != nil {
		return nil, nil, gf_err
	}

	//-----------------

	return running_job, job_expected_outputs_lst, nil
}

//-------------------------------------------------
func RunUploadedImages(pClientTypeStr string,
	p_images_to_process_lst []gf_images_jobs_core.GF_image_uploaded_to_process,
	p_flows_names_lst       []string,
	pJobsMngrCh             gf_images_jobs_core.JobsMngr,
	pRuntimeSys             *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, *gf_core.GFerror) {

	job_cmd_str    := "start_job_uploaded_imgs"
	job_init_ch    := make(chan *gf_images_jobs_core.GFjobRunning)
	job_updates_ch := make(chan gf_images_jobs_core.JobUpdateMsg, 10)
	

	job_msg := gf_images_jobs_core.JobMsg{
		Client_type_str:                pClientTypeStr,
		Cmd_str:                        job_cmd_str,
		Job_init_ch:                    job_init_ch,
		Job_updates_ch:                 job_updates_ch,
		Images_uploaded_to_process_lst: p_images_to_process_lst,
		Flows_names_lst:                p_flows_names_lst,
	}

	// SEND_MSG
	pJobsMngrCh <- job_msg

	// RECEIVE_MSG - get running_job info back from jobs_mngr
	running_job := <- job_init_ch
	
	// spew.Dump(running_job)

	return running_job, nil
}

//-------------------------------------------------
// START
func RunExternImages(pClientTypeStr string,
	pImagesExternToProcessLst []gf_images_jobs_core.GF_image_extern_to_process,
	pFlowsNamesLst            []string,
	pJobsMngrCh               gf_images_jobs_core.JobsMngr,
	pRuntimeSys               *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, []*GF_job_expected_output, *gf_core.GFerror) {
	pRuntimeSys.LogFun("INFO", fmt.Sprintf("images_extern_to_process - %s", fmt.Sprint(pImagesExternToProcessLst)))

	//-----------------
	// SEND_MSG_TO_JOBS_MNGR
	jobCmdStr    := "start_job"
	jobInitCh    := make(chan *gf_images_jobs_core.GFjobRunning)
	jobUpdatesCh := make(chan gf_images_jobs_core.JobUpdateMsg, 10) // ADD!! - channel buffer size should be larger for large jobs (with a lot of images)

	jobMsg := gf_images_jobs_core.JobMsg{
		Client_type_str:              pClientTypeStr,
		Cmd_str:                      jobCmdStr,
		Job_init_ch:                  jobInitCh,
		Job_updates_ch:               jobUpdatesCh,
		Images_extern_to_process_lst: pImagesExternToProcessLst,
		Flows_names_lst:              pFlowsNamesLst,
	}

	// SEND_MSG
	pJobsMngrCh <- jobMsg

	// RECEIVE_MSG - get running_job info back from jobs_mngr
	runningJob := <- jobInitCh

	//-----------------
	// JOB_EXPECTED_OUTPUT - its "expected" because results are not available yet (and might not
	//                       be available for some time), and yet we still want to have some of the expected
	//                       values so that other parts of the system can initialize in parallel with the job 
	//                       completing.

	imagesSourceURLsLst := []string{}
	for _, image_to_process := range pImagesExternToProcessLst {
		imagesSourceURLsLst = append(imagesSourceURLsLst, image_to_process.Source_url_str)
	}

	jobExpectedOutputsLst, gfErr := getJobExpectedOutput(imagesSourceURLsLst, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//-----------------
	return runningJob, jobExpectedOutputsLst, nil
}

//-------------------------------------------------
func GetJobUpdateCh(pJobIDstr string,
	pJobsMngrCh gf_images_jobs_core.JobsMngr,
	pRuntimeSys *gf_core.RuntimeSys) chan gf_images_jobs_core.JobUpdateMsg {

	msg_response_ch := make(chan interface{})
	defer close(msg_response_ch)

	job_cmd_str := "get_job_update_ch"
	job_msg     := gf_images_jobs_core.JobMsg{
		Job_id_str:      pJobIDstr,
		Cmd_str:         job_cmd_str,
		Msg_response_ch: msg_response_ch,
	}

	pJobsMngrCh <- job_msg

	response          := <-msg_response_ch
	job_updates_ch, _ := response.(chan gf_images_jobs_core.JobUpdateMsg)

	return job_updates_ch
}

//-------------------------------------------------
func CleanupJob(pJobIDstr string,
	pJobsMngrCh gf_images_jobs_core.JobsMngr,
	pRuntimeSys *gf_core.RuntimeSys) {

	jobCmdStr := "cleanup_job"
	jobMsg    := gf_images_jobs_core.JobMsg{
		Job_id_str: pJobIDstr,
		Cmd_str:    jobCmdStr,
	}

	pJobsMngrCh <- jobMsg
}

//-------------------------------------------------
// VAR
//-------------------------------------------------
func getJobExpectedOutput(pImagesSourceURIsLst []string,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GF_job_expected_output, *gf_core.GFerror) {

	jobExpectedOutputsLst := []*GF_job_expected_output{}

	for _, imgSourceURLstr := range pImagesSourceURIsLst {

		//--------------
		// IMAGE_ID
		imageIDstr, gfErr := gf_images_core.CreateIDfromURL(imgSourceURLstr, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		//--------------

		// all thumbs are stored as jpeg's
		thumbsExtStr := "jpeg"
		
		output := &GF_job_expected_output{
			Image_id_str:                      imageIDstr,
			Image_source_url_str:              imgSourceURLstr,
			Thumbnail_small_relative_url_str : fmt.Sprintf("/images/d/thumbnails/%s_thumb_small.%s",  imageIDstr, thumbsExtStr),
			Thumbnail_medium_relative_url_str: fmt.Sprintf("/images/d/thumbnails/%s_thumb_medium.%s", imageIDstr, thumbsExtStr),
			Thumbnail_large_relative_url_str:  fmt.Sprintf("/images/d/thumbnails/%s_thumb_large.%s",  imageIDstr, thumbsExtStr),
		}
		jobExpectedOutputsLst = append(jobExpectedOutputsLst, output)
	}

	//-----------------

	return jobExpectedOutputsLst, nil
}