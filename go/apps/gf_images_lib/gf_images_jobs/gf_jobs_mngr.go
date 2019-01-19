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
	Type_str             string        `json:"type_str"`
	Image_id_str         string        `json:"image_id_str"`
	Image_source_url_str string        `json:"image_source_url_str"`
	Err_str              string        `json:"err_str,omitempty"`  //if the update indicates an error, this is its value
	Image_thumbs         *gf_images_utils.Gf_image_thumbs `json:"-"`
}
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
						continue
					}
					//--------------------
					//MARK_JOB_AS_COMPLETE

					job_end_time_f := float64(time.Now().UnixNano())/1000000000.0
					update_err     := p_runtime_sys.Mongodb_coll.Update(bson.M{
							"t":"img_running_job","id_str":job_id_str,
						},
						bson.M{
							"$set":bson.M{
								"status_str":"complete",
								"end_time_f":job_end_time_f,
							},
						},)
					if update_err != nil {
						_ = gf_core.Error__create("failed to update an img_running_job in the DB, as complete and its end_time",
							"mongodb_update_error",
							&map[string]interface{}{
								"job_id_str":    job_id_str,
								"job_end_time_f":job_end_time_f,
							},
							update_err, "gf_images_jobs", p_runtime_sys)
					}
					//--------------------

					delete(running_jobs_map,job_id_str) //remove running job from lookup, since its complete
					close(job_msg.job_updates_ch)
				//------------------------
				case "get_running_job_update_ch":

					job_id_str := job_msg.job_id_str

					if _,ok := running_jobs_map[job_id_str]; ok {

						job_updates_ch := running_jobs_map[job_id_str]
						job_msg.msg_response_ch <- job_updates_ch

					} else {
						job_msg.msg_response_ch <- nil
					}
				//------------------------
			} 
		}
	}()
	return jobs_mngr_ch
}