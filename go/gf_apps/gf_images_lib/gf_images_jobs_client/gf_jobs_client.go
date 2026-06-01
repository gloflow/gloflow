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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func RunLocalImgs(pClientTypeStr string,
	pImagesToProcessLst []gf_images_jobs_core.GFimageLocalToProcess,
	pFlowsNamesLst      []string,
	pWaitForCompletionBool bool,
	pUserID             gf_core.GF_ID,
	pJobsMngrCh         gf_images_jobs_core.JobsMngr,
	pRuntimeSys         *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, []*GFjobExpectedOutput, *gf_core.GFerror) {

	jobCmdStr    := "start_job_local_imgs"
	jobInitCh    := make(chan *gf_images_jobs_core.GFjobRunning)
	jobUpdatesCh := make(chan gf_images_jobs_core.JobUpdateMsg, 10)
	
	jobMsg := gf_images_jobs_core.JobMsg{
		UserID:                      pUserID,
		Client_type_str:             pClientTypeStr,
		CmdStr:                     jobCmdStr,
		Job_init_ch:                 jobInitCh,
		Job_updates_ch:              jobUpdatesCh,
		Images_local_to_process_lst: pImagesToProcessLst,
		Flows_names_lst:             pFlowsNamesLst,
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

	imagesLocalPathsLst := []string{}
	for _, imageToProcess := range pImagesToProcessLst {
		imagesLocalPathsLst = append(imagesLocalPathsLst, imageToProcess.LocalFilePathStr)
	}

	//-----------------
	// COLLECT_GENERATED_IDS
	// Collect image IDs as they are generated (happens immediately before processing)
	generatedImageIDsLst := collectGeneratedImageIDs(jobUpdatesCh, imagesLocalPathsLst, pRuntimeSys)

	//-----------------
	// GET_JOB_EXPECTED_OUTPUT
	// Use the generated IDs that were created at pipeline start
	jobExpectedOutputsLst, gfErr := getJobExpectedOutput(imagesLocalPathsLst, generatedImageIDsLst, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//-----------------
	// WAIT_FOR_COMPLETION
	if pWaitForCompletionBool {
		_ = waitForJobCompletion(jobUpdatesCh, len(pImagesToProcessLst), "Image", pRuntimeSys)
	}

	//-----------------

	return runningJob, jobExpectedOutputsLst, nil
}

//-------------------------------------------------

func RunUploadedImages(pClientTypeStr string,
	pImagesToProcessLst []gf_images_jobs_core.GFimageUploadedToProcess,
	pFlowsNamesLst      []string,
	pWaitForCompletionBool bool,
	pUserID             gf_core.GF_ID,
	pJobsMngrCh         gf_images_jobs_core.JobsMngr,
	pRuntimeSys         *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, *gf_core.GFerror) {

	job_cmd_str    := "start_job_uploaded_imgs"
	job_init_ch    := make(chan *gf_images_jobs_core.GFjobRunning)
	job_updates_ch := make(chan gf_images_jobs_core.JobUpdateMsg, 10)

	job_msg := gf_images_jobs_core.JobMsg{
		UserID:                         pUserID,
		Client_type_str:                pClientTypeStr,
		CmdStr:                        job_cmd_str,
		Job_init_ch:                    job_init_ch,
		Job_updates_ch:                 job_updates_ch,
		Images_uploaded_to_process_lst: pImagesToProcessLst,
		Flows_names_lst:                pFlowsNamesLst,
	}

	// SEND_MSG
	pJobsMngrCh <- job_msg

	// RECEIVE_MSG - get running_job info back from jobs_mngr
	runningJob := <- job_init_ch

	//-----------------
	// WAIT_FOR_COMPLETION
	if pWaitForCompletionBool {
		_ = waitForJobCompletion(job_updates_ch, len(pImagesToProcessLst), "Image", pRuntimeSys)
	}

	//-----------------

	return runningJob, nil
}

//-------------------------------------------------
// RUN

func RunExternImages(pClientTypeStr string,
	pImagesExternToProcessLst []gf_images_jobs_core.GFimageExternToProcess,
	pFlowsNamesLst            []string,
	pWaitForCompletionBool    bool,
	pUserID                   gf_core.GF_ID,
	pJobsMngrCh               gf_images_jobs_core.JobsMngr,
	pRuntimeSys               *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, []*GFjobExpectedOutput, *gf_core.GFerror) {

	//-----------------
	// SEND_MSG_TO_JOBS_MNGR
	jobCmdStr    := "start_job"
	jobInitCh    := make(chan *gf_images_jobs_core.GFjobRunning)
	jobUpdatesCh := make(chan gf_images_jobs_core.JobUpdateMsg, 10) // ADD!! - channel buffer size should be larger for large jobs (with a lot of images)

	jobMsg := gf_images_jobs_core.JobMsg{
		UserID:                       pUserID,
		Client_type_str:              pClientTypeStr,
		CmdStr:                       jobCmdStr,
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
	for _, imageToProcess := range pImagesExternToProcessLst {
		imagesSourceURLsLst = append(imagesSourceURLsLst, imageToProcess.SourceURLstr)
	}

	//-----------------
	// COLLECT_GENERATED_IDS
	// Collect image IDs as they are generated (happens immediately before processing)
	generatedImageIDsLst := collectGeneratedImageIDs(jobUpdatesCh, imagesSourceURLsLst, pRuntimeSys)

	//-----------------
	// GET_JOB_EXPECTED_OUTPUT
	// Use the generated IDs that were created at pipeline start
	jobExpectedOutputsLst, gfErr := getJobExpectedOutput(imagesSourceURLsLst, generatedImageIDsLst, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//-----------------
	// WAIT_FOR_COMPLETION
	if pWaitForCompletionBool {
		_ = waitForJobCompletion(jobUpdatesCh, len(pImagesExternToProcessLst), "Image", pRuntimeSys)
	}

	//-----------------
	return runningJob, jobExpectedOutputsLst, nil
}

//-------------------------------------------------

func RunClassifyImages(pClientTypeStr string,
	pImagesToProcessLst []gf_images_jobs_core.GFimageClassificationToProcess,
	pWaitForCompletionBool bool,
	pUserID             gf_core.GF_ID,
	pJobsMngrCh         gf_images_jobs_core.JobsMngr,
	pRuntimeSys         *gf_core.RuntimeSys) ([]string, *gf_images_jobs_core.GFjobRunning, *gf_core.GFerror) {

	jobCmdStr    := "start_job_classify_imgs"
	jobInitCh    := make(chan *gf_images_jobs_core.GFjobRunning)
	jobUpdatesCh := make(chan gf_images_jobs_core.JobUpdateMsg, 10)
	responseCh := make(chan interface{})

	jobMsg := gf_images_jobs_core.JobMsg{
		UserID:              pUserID,
		Client_type_str:     pClientTypeStr,
		CmdStr:             jobCmdStr,
		Job_init_ch:         jobInitCh,
		Job_updates_ch:      jobUpdatesCh,
		ResponseCh:          responseCh,
		ImagesToClassifyLst: pImagesToProcessLst,
	}

	// SEND_MSG
	pJobsMngrCh <- jobMsg

	// RECEIVE_MSG - get running_job info back from jobs_mngr
	runningJob := <- jobInitCh

	classes := <- responseCh
	classesStrLst, _ := classes.([]string)

	//-----------------
	// WAIT_FOR_COMPLETION
	// Note: RunClassifyImages already waits for the response channel, 
	// but we add explicit waiting logic for consistency
	if pWaitForCompletionBool {
		_ = waitForJobCompletion(jobUpdatesCh, len(pImagesToProcessLst), "Image classification", pRuntimeSys)
	}

	//-----------------

	return classesStrLst, runningJob, nil
}

//-------------------------------------------------
