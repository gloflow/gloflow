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

package gf_images_lib



import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
)

//---------------------------------------------------
func T__test_image_job__updates(p_job_id_str string,
	p_jobs_mngr   gf_images_jobs.Jobs_mngr,
	p_runtime_sys *gf_core.Runtime_sys) {

	//-------------
	//TEST_JOB_UPDATES
	job_updates_ch := gf_images_jobs.Job__get_update_ch(p_job_id_str, p_jobs_mngr, p_runtime_sys)

	for ;; {

		fmt.Println("\n\n------------------------- TESTING - GET_JOB_UPDATE -----")
		job_update := <-job_updates_ch

		spew.Dump(job_update)

		job_update_type_str := job_update.Type_str
		if job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__ERROR {
			panic("job encountered an error while processing")
		}

		if !(job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__OK || job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__COMPLETED) {
			panic(fmt.Sprintf("job_update is expected to be of type 'ok' but instead is - %s", job_update_type_str))
		}
		
		//test complete
		if job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__COMPLETED {
			break
		}
	}
	//-------------
}