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
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_utils"
)
//-------------------------------------------------
type Jobs_mngr chan Job_msg

type Running_job struct {
	Id                    bson.ObjectId       `bson:"_id,omitempty"`
	Id_str                string              `bson:"id_str"`
	T_str                 string              `bson:"t"`
	Client_type_str       string              `bson:"client_type_str"`
	Status_str            string              `bson:"status_str"` //"running"|"complete"
	Start_time_f          float64             `bson:"start_time_f"`
	End_time_f            float64             `bson:"end_time_f"`
	Images_to_process_lst []Image_to_process  `bson:"images_to_process_lst"`
	Errors_lst            []Job_Error         `bson:"errors_lst"`
	job_updates_ch        chan Job_update_msg `bson:"-"`
}

type Image_to_process struct {
	Source_url_str      string `bson:"source_url_str"` //FIX!! - rename this to Origin_url_str to be consistent with other origin_url naming
	Origin_page_url_str string `bson:"origin_page_url_str"`
}

type Job_msg struct {
	job_id_str            string
	client_type_str       string 
	cmd_str               string //"start_job"|"get_running_job_ids"
	msg_response_ch       chan interface{}
	job_updates_ch        chan Job_update_msg
	images_to_process_lst []Image_to_process
	flows_names_lst       []string
}

type Job_update_msg struct {
	Name_str             string              `json:"name_str"`
	Type_str             job_update_type_val `json:"type_str"`             //"ok"|"error"|"complete"
	Image_id_str         string              `json:"image_id_str"`
	Image_source_url_str string              `json:"image_source_url_str"`
	Err_str              string              `json:"err_str,omitempty"`    //if the update indicates an error, this is its value
	Image_thumbs         *gf_images_utils.Gf_image_thumbs `json:"-"`
}

type job_status_val string
const JOB_STATUS__FAILED    job_status_val = "failed"
const JOB_STATUS__COMPLETED job_status_val = "completed"

type job_update_type_val string
const JOB_UPDATE_TYPE__OK        job_update_type_val = "ok"
const JOB_UPDATE_TYPE__ERROR     job_update_type_val = "error"
const JOB_UPDATE_TYPE__COMPLETED job_update_type_val = "completed"
//-------------------------------------------------
//SERVER
//-------------------------------------------------
func Jobs_mngr__init(p_images_store_local_dir_path_str string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_s3_bucket_name_str                         string,
	p_s3_info                                    *gf_core.Gf_s3_info,
	p_runtime_sys                                *gf_core.Runtime_sys) Jobs_mngr {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_jobs_mngr.Jobs_mngr__init()")

	jobs_mngr_ch := make(chan Job_msg, 100)

	//IMPORTANT!! - start jobs_mngr as an independent goroutine of the HTTP handlers at
	//              service initialization time
	go func() {
		
		running_jobs_map := map[string]chan Job_update_msg{}

		//listen to messages
		for {
			job_msg := <- jobs_mngr_ch

			//IMPORTANT!! - only one job is processed per jobs_mngr.
			//              Scaling is done with multiple jobs_mngr's (exp. per-core)           
			switch job_msg.cmd_str {

				//------------------------
				//START_JOB
				case "start_job":

					job_id_str                  := job_msg.job_id_str
					running_jobs_map[job_id_str] = job_msg.job_updates_ch

					run_job_gf_err := run_job(job_id_str,
						job_msg.client_type_str,
						job_msg.images_to_process_lst,
						job_msg.flows_names_lst,
						job_msg.job_updates_ch,
						p_images_store_local_dir_path_str,
						p_images_thumbnails_store_local_dir_path_str,
						p_s3_bucket_name_str,
						p_s3_info,
						p_runtime_sys)

					if run_job_gf_err != nil {
						_ = db__jobs_mngr__update_job_status(JOB_STATUS__FAILED, job_id_str, p_runtime_sys)
					} else {
						_ = db__jobs_mngr__update_job_status(JOB_STATUS__COMPLETED, job_id_str, p_runtime_sys)
					}
				//------------------------
				//GET_JOB_UPDATE_CH
				case "get_job_update_ch":

					job_id_str := job_msg.job_id_str

					if _,ok := running_jobs_map[job_id_str]; ok {

						job_updates_ch := running_jobs_map[job_id_str]
						job_msg.msg_response_ch <- job_updates_ch

					} else {
						job_msg.msg_response_ch <- nil
					}
				//------------------------
				//CLEANUP_JOB
				case "cleanup_job":
					job_id_str := job_msg.job_id_str
					delete(running_jobs_map, job_id_str) //remove running job from lookup, since its complete
				//------------------------
			} 
		}
	}()
	return jobs_mngr_ch
}
//-------------------------------------------------
func db__jobs_mngr__update_job_status(p_status_str job_status_val,
	p_job_id_str  string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_jobs_mngr.db__jobs_mngr__update_job_status()")

	if p_status_str != JOB_STATUS__COMPLETED && p_status_str != JOB_STATUS__FAILED {
		//status values are not generated at runtime, but are static, so its ok to panic here since
		//this should never be countered in production
		panic(fmt.Sprintf("job status value thats not allowed - %s", p_status_str))
	}

	job_end_time_f := float64(time.Now().UnixNano())/1000000000.0
	err            := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t":      "img_running_job",
			"id_str": p_job_id_str,
		},
		bson.M{
			"$set":bson.M{
				"status_str":p_status_str,
				"end_time_f":job_end_time_f,
			},
		},)
		
	if err != nil {
		gf_err := gf_core.Error__create("failed to update an img_running_job in the DB, as complete and its end_time",
			"mongodb_update_error",
			&map[string]interface{}{
				"job_id_str":    p_job_id_str,
				"job_end_time_f":job_end_time_f,
			},
			err, "gf_images_jobs", p_runtime_sys)
		return gf_err
	}

	return nil

}