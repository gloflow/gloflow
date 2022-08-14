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

package gf_images_lib

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
)

//---------------------------------------------------
func T__test_image_job__updates(pJobIDstr string,
	pJobsMngr   gf_images_jobs_core.JobsMngr,
	pRuntimeSys *gf_core.RuntimeSys) {

	//-------------
	// TEST_JOB_UPDATES
	jobUpdatesCh := gf_images_jobs_client.GetJobUpdateCh(pJobIDstr, pJobsMngr, pRuntimeSys)

	for ;; {

		fmt.Println("\n\n------------------------- TESTING - GET_JOB_UPDATE -----")
		jobUpdate := <- jobUpdatesCh

		spew.Dump(jobUpdate)

		jobUpdateTypeStr := jobUpdate.Type_str
		if jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__ERROR {
			panic("job encountered an error while processing")
		}

		if !(jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__OK || jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__COMPLETED) {
			panic(fmt.Sprintf("job_update is expected to be of type 'ok' but instead is - %s", jobUpdateTypeStr))
		}
		
		// test complete
		if jobUpdateTypeStr == gf_images_jobs_core.JOB_UPDATE_TYPE__COMPLETED {
			break
		}
	}

	//-------------
}