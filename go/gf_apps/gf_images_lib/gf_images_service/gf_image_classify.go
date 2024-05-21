/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

package gf_images_service

import (
	"fmt"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
)

//---------------------------------------------------

type GFimageClassifyInput struct {
	ClientTypeStr string
	ImagesIDsLst  []gf_images_core.GFimageID
}

//---------------------------------------------------

func ImageClassify(pInput *GFimageClassifyInput,
	pUserID      gf_core.GF_ID,
	pJobsMngrCh  chan gf_images_jobs_core.JobMsg,
	pServiceInfo *gf_images_core.GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {





	imagesToProcessLst := []gf_images_jobs_core.GFimageClassificationToProcess{}
	for _, imageID := range pInput.ImagesIDsLst {

		imageToProcess := gf_images_jobs_core.GFimageClassificationToProcess{
			GFimageIDstr: imageID,
		}

		imagesToProcessLst = append(imagesToProcessLst, imageToProcess)

	}

	classesLst, runningJob, gfErr := gf_images_jobs_client.RunClassifyImages(pInput.ClientTypeStr,
		imagesToProcessLst,
		pUserID,
		pJobsMngrCh,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
			



	jobIDstr := runningJob.IDstr

	fmt.Printf("job_id - %s\n", jobIDstr)



	


	//------------------
	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		eventMetaMap := map[string]interface{}{
			"user_id": pUserID,
		}
		gf_events.EmitApp(gf_images_core.GF_EVENT_APP__IMAGE_CLASSIFY,
			eventMetaMap,
			pUserID,
			pCtx,
			pRuntimeSys)
	}

	//------------------

	return classesLst, nil
}