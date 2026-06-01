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
)

//-------------------------------------------------
// COLLECT_GENERATED_IMAGE_IDS
//
// Collects image IDs as they are generated at the start of the processing pipeline.
// This happens synchronously before actual image processing begins.
// Maps source URIs to their generated IDs to maintain correct ordering.
//
// Parameters:
//   - pJobUpdatesCh: channel to receive job updates from the job manager
//   - pSourceURIsLst: list of source URIs/paths in the order they were submitted
//   - pRuntimeSys: runtime system for logging
//
// Returns:
//   - List of generated image IDs in the same order as pSourceURIsLst
func collectGeneratedImageIDs(pJobUpdatesCh chan gf_images_jobs_core.JobUpdateMsg,
	pSourceURIsLst []string,
	pRuntimeSys *gf_core.RuntimeSys) []gf_images_core.GFimageID {
	
	pTotalImages := len(pSourceURIsLst)
	
	// Map source URI -> generated ID
	uriToIDmap := make(map[string]gf_images_core.GFimageID)
	
	for len(uriToIDmap) < pTotalImages {
		update := <- pJobUpdatesCh
		
		if update.Type_str == gf_images_jobs_core.JOB_UPDATE_TYPE__ID_GENERATED {
			sourceURIstr := update.Image_source_url_str
			uriToIDmap[sourceURIstr] = update.ImageIDstr
			pRuntimeSys.LogFun("INFO", fmt.Sprintf("Image ID generated %d/%d: %s -> %s", 
				len(uriToIDmap), pTotalImages, sourceURIstr, update.ImageIDstr))
		}
		// Ignore other message types during ID collection phase
	}
	
	// Build ordered list matching input order
	generatedImageIDsLst := make([]gf_images_core.GFimageID, pTotalImages)
	for i, sourceURIstr := range pSourceURIsLst {
		if imageID, exists := uriToIDmap[sourceURIstr]; exists {
			generatedImageIDsLst[i] = imageID
		} else {
			pRuntimeSys.LogFun("ERROR", fmt.Sprintf("No generated ID found for source URI: %s", sourceURIstr))
			generatedImageIDsLst[i] = gf_images_core.GFimageID("") // shouldn't happen
		}
	}
	
	pRuntimeSys.LogFun("INFO", fmt.Sprintf("All %d image IDs collected and ordered", pTotalImages))
	return generatedImageIDsLst
}

//-------------------------------------------------
// WAIT_FOR_JOB_COMPLETION
//
// Waits for job completion by reading from the job updates channel
// until all images are processed (either completed or errored).
//
// Parameters:
//   - pJobUpdatesCh: channel to receive job updates from the job manager
//   - pTotalImages: total number of images expected to complete
//   - pJobTypeStr: type of job for logging (e.g., "Image", "Image classification")
//   - pRuntimeSys: runtime system for logging
//
// Returns:
//   - List of actual image IDs from completed processing (in order of completion)
func waitForJobCompletion(pJobUpdatesCh chan gf_images_jobs_core.JobUpdateMsg,
	pTotalImages int,
	pJobTypeStr string,
	pRuntimeSys *gf_core.RuntimeSys) []gf_images_core.GFimageID {
	
	completedCount := 0
	processedImageIDs := []gf_images_core.GFimageID{}
	
	for completedCount < pTotalImages {
		update := <- pJobUpdatesCh
		
		if update.Type_str == gf_images_jobs_core.JOB_UPDATE_TYPE__COMPLETED {
			completedCount++
			processedImageIDs = append(processedImageIDs, update.ImageIDstr)
			pRuntimeSys.LogFun("INFO", fmt.Sprintf("%s %d/%d completed: %s", pJobTypeStr, completedCount, pTotalImages, update.ImageIDstr))
		} else if update.Type_str == gf_images_jobs_core.JOB_UPDATE_TYPE__ERROR {
			pRuntimeSys.LogFun("ERROR", fmt.Sprintf("%s processing error: %s - %s", pJobTypeStr, update.ImageIDstr, update.Err_str))
		}
	}
	
	pRuntimeSys.LogFun("INFO", fmt.Sprintf("All %d %ss completed processing", pTotalImages, pJobTypeStr))
	return processedImageIDs
}

//-------------------------------------------------
// GET_JOB_UPDATE_CH

func GetJobUpdateCh(pJobIDstr string,
	pJobsMngrCh gf_images_jobs_core.JobsMngr,
	pRuntimeSys *gf_core.RuntimeSys) chan gf_images_jobs_core.JobUpdateMsg {

	responseCh := make(chan interface{})
	defer close(responseCh)

	jobCmdStr := "get_job_update_ch"
	jobMsg    := gf_images_jobs_core.JobMsg{
		JobIDstr:   pJobIDstr,
		CmdStr:     jobCmdStr,
		ResponseCh: responseCh,
	}

	pJobsMngrCh <- jobMsg

	response        := <- responseCh
	jobUpdatesCh, _ := response.(chan gf_images_jobs_core.JobUpdateMsg)

	return jobUpdatesCh
}

//-------------------------------------------------

func CleanupJob(pJobIDstr string,
	pJobsMngrCh gf_images_jobs_core.JobsMngr,
	pRuntimeSys *gf_core.RuntimeSys) {

	jobCmdStr := "cleanup_job"
	jobMsg    := gf_images_jobs_core.JobMsg{
		JobIDstr: pJobIDstr,
		CmdStr:    jobCmdStr,
	}

	pJobsMngrCh <- jobMsg
}

//-------------------------------------------------

// getJobExpectedOutput - Build expected output structures for processed images.
//
// IMPORTANT: pGeneratedImageIDsLst MUST be provided from collectGeneratedImageIDs().
// These are the actual IDs generated at pipeline start. Never regenerate IDs here
// as timestamp-based ID generation will produce different IDs.
func getJobExpectedOutput(pImagesSourceURIsLst []string,
	pGeneratedImageIDsLst []gf_images_core.GFimageID,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFjobExpectedOutput, *gf_core.GFerror) {

	if len(pGeneratedImageIDsLst) != len(pImagesSourceURIsLst) {
		gfErr := gf_core.ErrorCreate("image IDs count doesn't match source URIs count",
			"verify_error",
			map[string]interface{}{
				"ids_count":  len(pGeneratedImageIDsLst),
				"uris_count": len(pImagesSourceURIsLst),
			},
			nil, "gf_images_jobs_client", pRuntimeSys)
		return nil, gfErr
	}

	jobExpectedOutputsLst := []*GFjobExpectedOutput{}

	for i, imgSourceURLstr := range pImagesSourceURIsLst {

		//--------------
		// IMAGE_ID - use the ID generated at pipeline start
		imageIDstr := pGeneratedImageIDsLst[i]

		//--------------

		// all thumbs are stored as jpeg's
		thumbsExtStr := "jpeg"
		
		output := &GFjobExpectedOutput{
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
