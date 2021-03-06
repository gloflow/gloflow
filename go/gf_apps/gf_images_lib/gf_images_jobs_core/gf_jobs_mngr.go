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
/*BosaC.Jan30.2020. <3 volim te zauvek*/

package gf_images_jobs_core

import (
	"fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type Jobs_mngr chan Job_msg

type GF_job_running struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	Id_str          string        `bson:"id_str"`
	T_str           string        `bson:"t"`
	Client_type_str string        `bson:"client_type_str"`
	Status_str      string        `bson:"status_str"` // "running"|"complete"
	Start_time_f    float64       `bson:"start_time_f"`
	End_time_f      float64       `bson:"end_time_f"`

	// LEGACY!! - update "images_to_process_lst" to "images_extern_to_process_lst" in the DB (bson)
	Images_extern_to_process_lst   []GF_image_extern_to_process   `bson:"images_to_process_lst"`
	Images_uploaded_to_process_lst []GF_image_uploaded_to_process `bson:"images_uploaded_to_process_lst"`
	
	Errors_lst     []Job_Error         `bson:"errors_lst"`
	job_updates_ch chan Job_update_msg `bson:"-"`
}

//------------------------
// IMAGES_TO_PROCESS
type GF_image_extern_to_process struct {
	Source_url_str      string `bson:"source_url_str"` // FIX!! - rename this to Origin_url_str to be consistent with other origin_url naming
	Origin_page_url_str string `bson:"origin_page_url_str"`
}

type GF_image_uploaded_to_process struct {
	Gf_image_id_str  gf_images_utils.Gf_image_id
	S3_file_path_str string // path to image in S3 in a bucket that it was originally uploaded to by client
}

type GF_image_local_to_process struct {
	Local_file_path_str string
}

//------------------------
// JOB_MSGS

type Job_msg struct {
	Job_id_str                     string // if its an existing job. for new jobs the mngr creates a new job ID
	Client_type_str                string 
	Cmd_str                        string // "start_job" | "get_running_job_ids"
	
	Job_init_ch                    chan *GF_job_running // used by clients for receiving outputs of job initialization by jobs_mngr
	Job_updates_ch                 chan Job_update_msg  // used by jobs_mngr to send job_updates to
	Msg_response_ch                chan interface{}     // DEPRECATED!! use a specific struct as a message format, interface{} too general.

	Images_extern_to_process_lst   []GF_image_extern_to_process
	Images_uploaded_to_process_lst []GF_image_uploaded_to_process
	Images_local_to_process_lst    []GF_image_local_to_process

	Flows_names_lst                []string
}

type Job_update_msg struct {
	Name_str             string                      `json:"name_str"`
	Type_str             job_update_type_val         `json:"type_str"`             // "ok" | "error" | "complete"
	Image_id_str         gf_images_utils.Gf_image_id `json:"image_id_str"`
	Image_source_url_str string                      `json:"image_source_url_str"`
	Err_str              string                      `json:"err_str,omitempty"`    // if the update indicates an error, this is its value
	Image_thumbs         *gf_images_utils.Gf_image_thumbs `json:"-"`
}

//------------------------
// JOBS_LIFECYCLE
type GF_jobs_lifecycle_callbacks struct {
	Job_type__transform_imgs__fun func() *gf_core.Gf_error
	Job_type__uploaded_imgs__fun  func() *gf_core.Gf_error
}

//------------------------
type job_status_val string
const JOB_STATUS__FAILED         job_status_val = "failed"
const JOB_STATUS__FAILED_PARTIAL job_status_val = "failed_partial"
const JOB_STATUS__COMPLETED      job_status_val = "completed"

type job_update_type_val string
const JOB_UPDATE_TYPE__OK        job_update_type_val = "ok"
const JOB_UPDATE_TYPE__ERROR     job_update_type_val = "error"
const JOB_UPDATE_TYPE__COMPLETED job_update_type_val = "completed"

//-------------------------------------------------
// CREATE_RUNNING_JOB
func Jobs_mngr__create_running_job(p_client_type_str string,
	p_job_updates_ch chan Job_update_msg,
	p_runtime_sys    *gf_core.Runtime_sys) (*GF_job_running, *gf_core.Gf_error) {

	job_start_time_f := float64(time.Now().UnixNano())/1000000000.0
	job_id_str       := fmt.Sprintf("job:%f", job_start_time_f)

	running_job := &GF_job_running{
		Id_str:          job_id_str,
		T_str:           "img_running_job",
		Client_type_str: p_client_type_str,
		Status_str:      "running",
		Start_time_f:    job_start_time_f,
		job_updates_ch:  p_job_updates_ch,
	}

	// DB
	gf_err := db__jobs_mngr__create_running_job(running_job, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return running_job, nil
}

//-------------------------------------------------
// INIT
func Jobs_mngr__init(p_images_store_local_dir_path_str string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_media_domain_str                           string,
	p_lifecycle_callbacks                        *GF_jobs_lifecycle_callbacks,
	p_config                                     *gf_images_utils.GF_config,
	p_s3_info                                    *gf_core.GF_s3_info,
	p_runtime_sys                                *gf_core.Runtime_sys) Jobs_mngr {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_mngr.Jobs_mngr__init()")

	jobs_mngr_ch := make(chan Job_msg, 100)

	// IMPORTANT!! - start jobs_mngr as an independent goroutine of the HTTP handlers at
	//               service initialization time
	go func() {
		
		running_jobs_map := map[string]chan Job_update_msg{}

		// listen to messages
		for {
			job_msg := <- jobs_mngr_ch

			// IMPORTANT!! - only one job is processed per jobs_mngr.
			//              Scaling is done with multiple jobs_mngr's (exp. per-core)           
			switch job_msg.Cmd_str {


				//------------------------
				// START_JOB_LOCAL_IMAGES
				case "start_job_local_imgs":

					
					// RUNNING_JOB
					running_job, gf_err := Jobs_mngr__create_running_job(job_msg.Client_type_str,
						job_msg.Job_updates_ch,
						p_runtime_sys)
					if gf_err != nil {
						continue
					}

					//------------------------
					// S3_BUCKETS

					var target_s3_bucket_name_str string
					if len(job_msg.Flows_names_lst) > 0 {
						main_flow_str := job_msg.Flows_names_lst[0]

						// check if the specified flow has an associated s3 bucket
						var ok bool
						target_s3_bucket_name_str, ok = p_config.Images_flow_to_s3_bucket_map[main_flow_str]
						if !ok {
							target_s3_bucket_name_str = p_config.Images_flow_to_s3_bucket_default_str
						}
					} else {
						target_s3_bucket_name_str = p_config.Images_flow_to_s3_bucket_default_str
					}

					//------------------------

					job_run__runtime := &GF_job_run__runtime{
						job_id_str:          running_job.Id_str,
						job_client_type_str: job_msg.Client_type_str,
						job_updates_ch:      job_msg.Job_updates_ch,
						s3_info:             p_s3_info,
					}

					run_job_gf_errs_lst := run_job__local_imgs(job_msg.Images_local_to_process_lst,
						p_images_store_local_dir_path_str,
						p_images_thumbnails_store_local_dir_path_str,
						target_s3_bucket_name_str,
						job_run__runtime,
						p_runtime_sys)
					
					//------------------------
					// JOB_STATUS
					var job_status_str job_status_val
					if len(run_job_gf_errs_lst) == len(job_msg.Images_uploaded_to_process_lst) {
						job_status_str = JOB_STATUS__FAILED
					} else if len(run_job_gf_errs_lst) > 0 {
						job_status_str = JOB_STATUS__FAILED_PARTIAL
					} else {
						job_status_str = JOB_STATUS__COMPLETED
					}
					_ = db__jobs_mngr__update_job_status(job_status_str, running_job.Id_str, p_runtime_sys)

					//------------------------

				//------------------------
				// START_JOB_TRANSFORM_IMAGES
				case "start_job_transform_imgs":

					/*// RUST
					// FIX!! - this just runs Rust job code for testing.
					//         pass in proper job_cmd argument.
					run_job_rust()*/

					//------------------------
					// LIFECYCLE_CALLBACK
					gf_err := p_lifecycle_callbacks.Job_type__transform_imgs__fun()
					if gf_err != nil {
						continue
					}

					//------------------------

				//------------------------
				// START_JOB_UPLOADED_IMAGES
				case "start_job_uploaded_imgs":
					
					// RUNNING_JOB
					running_job, gf_err := Jobs_mngr__create_running_job(job_msg.Client_type_str,
						job_msg.Job_updates_ch,
						p_runtime_sys)
					if gf_err != nil {
						continue
					}

					// IMPORTANT!! - send sending running_job back to client, to avoid race conditions
					running_jobs_map[running_job.Id_str] = job_msg.Job_updates_ch

					// SEND_MSG
					job_msg.Job_init_ch <- running_job

					//------------------------
					// S3_BUCKETS
					source_s3_bucket_name_str := p_config.Uploaded_images_s3_bucket_str

					// ADD!! - due to legacy reasons the "general" flow is still used as the main flow
					//         that images are added to when a uploaded_images job is processed.
					//         this should be generalized so that images are added to dedicated flow S3 buckets
					//         if those flows have their S3_bucket mapping defined in Gf_config.Images_flow_to_s3_bucket_map
					target_s3_bucket_name_str := p_config.Images_flow_to_s3_bucket_map["general"]

					//------------------------

					job_run__runtime := &GF_job_run__runtime{
						job_id_str:          running_job.Id_str,
						job_client_type_str: job_msg.Client_type_str,
						job_updates_ch:      job_msg.Job_updates_ch,
						s3_info:             p_s3_info,
					}

					run_job_gf_errs_lst := run_job__uploaded_imgs(job_msg.Images_uploaded_to_process_lst,
						job_msg.Flows_names_lst,
						p_images_store_local_dir_path_str,
						p_images_thumbnails_store_local_dir_path_str,
						source_s3_bucket_name_str,
						target_s3_bucket_name_str,
						job_run__runtime,
						p_runtime_sys)
					
					//------------------------
					// JOB_STATUS
					var job_status_str job_status_val
					if len(run_job_gf_errs_lst) == len(job_msg.Images_uploaded_to_process_lst) {
						job_status_str = JOB_STATUS__FAILED
					} else if len(run_job_gf_errs_lst) > 0 {
						job_status_str = JOB_STATUS__FAILED_PARTIAL
					} else {
						job_status_str = JOB_STATUS__COMPLETED
					}
					_ = db__jobs_mngr__update_job_status(job_status_str, running_job.Id_str, p_runtime_sys)

					//------------------------
					// LIFECYCLE_CALLBACK
					gf_err = p_lifecycle_callbacks.Job_type__uploaded_imgs__fun()
					if gf_err != nil {
						continue
					}
				
				//------------------------
				// START_JOB_EXTERN_IMAGES
				// FIX!! - "start_job" needs to be "start_job_extern_imgs", update in all clients.
				case "start_job":

					// RUNNING_JOB
					running_job, gf_err := Jobs_mngr__create_running_job(job_msg.Client_type_str,
						job_msg.Job_updates_ch,
						p_runtime_sys)
					if gf_err != nil {
						continue
					}

					// IMPORTANT!! - send sending running_job back to client, to avoid race conditions
					running_jobs_map[running_job.Id_str] = job_msg.Job_updates_ch

					// SEND_MSG
					job_msg.Job_init_ch <- running_job
					
					//------------------------
					// ADD!! - due to legacy reasons the "general" flow is still used as the main flow
					//         that images are added to when a external_images job is processed.
					//         this should be generalized so that images are added to dedicated flow S3 buckets
					//         if those flows have their S3_bucket mapping defined in Gf_config.Images_flow_to_s3_bucket_map
					s3_bucket_name_str := p_config.Images_flow_to_s3_bucket_map["general"]

					//------------------------

					job_run__runtime := &GF_job_run__runtime{
						job_id_str:          running_job.Id_str,
						job_client_type_str: job_msg.Client_type_str,
						job_updates_ch:      job_msg.Job_updates_ch,
						s3_info:             p_s3_info,
					}

					run_job_gf_errs_lst := run_job__extern_imgs(job_msg.Images_extern_to_process_lst,
						job_msg.Flows_names_lst,
						p_images_store_local_dir_path_str,
						p_images_thumbnails_store_local_dir_path_str,

						p_media_domain_str,
						s3_bucket_name_str,
						job_run__runtime,
						p_runtime_sys)
					
					//------------------------
					// JOB_STATUS
					var job_status_str job_status_val
					if len(run_job_gf_errs_lst) == len(job_msg.Images_extern_to_process_lst) {
						job_status_str = JOB_STATUS__FAILED
					} else if len(run_job_gf_errs_lst) > 0 {
						job_status_str = JOB_STATUS__FAILED_PARTIAL
					} else {
						job_status_str = JOB_STATUS__COMPLETED
					}
					_ = db__jobs_mngr__update_job_status(job_status_str, running_job.Id_str, p_runtime_sys)

					//------------------------

				//------------------------
				// GET_JOB_UPDATE_CH
				case "get_job_update_ch":

					job_id_str := job_msg.Job_id_str

					if _, ok := running_jobs_map[job_id_str]; ok {

						job_updates_ch := running_jobs_map[job_id_str]
						job_msg.Msg_response_ch <- job_updates_ch

					} else {
						job_msg.Msg_response_ch <- nil
					}

				//------------------------
				// CLEANUP_JOB
				case "cleanup_job":
					job_id_str := job_msg.Job_id_str
					delete(running_jobs_map, job_id_str) // remove running job from lookup, since its complete

				//------------------------
			} 
		}
	}()
	return jobs_mngr_ch
}

//-------------------------------------------------
// DB
//-------------------------------------------------
func db__jobs_mngr__create_running_job(p_running_job *GF_job_running,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {


	ctx           := context.Background()
	coll_name_str := "gf_images__jobs_running" // p_runtime_sys.Mongo_coll.Name()
	gf_err        := gf_core.Mongo__insert(p_running_job,
		coll_name_str,
		map[string]interface{}{
			"running_job_id_str": p_running_job.Id_str,
			"client_type_str":    p_running_job.Client_type_str,
			"caller_err_msg_str": "failed to create a Running_job record into the DB",
		},
		ctx,
		p_runtime_sys)
	
	if gf_err != nil {
		return gf_err
	}

	return nil
}

//-------------------------------------------------
func db__jobs_mngr__update_job_status(p_status_str job_status_val,
	p_job_id_str  string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_mngr.db__jobs_mngr__update_job_status()")

	if p_status_str != JOB_STATUS__COMPLETED && p_status_str != JOB_STATUS__FAILED {
		// status values are not generated at runtime, but are static, so its ok to panic here since
		// this should never be countered in production
		panic(fmt.Sprintf("job status value thats not allowed - %s", p_status_str))
	}

	ctx := context.Background()

	job_end_time_f := float64(time.Now().UnixNano())/1000000000.0
	coll           := p_runtime_sys.Mongo_db.Collection("gf_images__jobs_running")
	_, err         := coll.UpdateMany(ctx, bson.M{
			"t":      "img_running_job",
			"id_str": p_job_id_str,
		},
		bson.M{
			"$set": bson.M{
				"status_str": p_status_str,
				"end_time_f": job_end_time_f,
			},
		},)
		
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update an img_running_job in the DB, as complete and its end_time",
			"mongodb_update_error",
			map[string]interface{}{
				"job_id_str":     p_job_id_str,
				"job_end_time_f": job_end_time_f,
			},
			err, "gf_jobs_mngr", p_runtime_sys)
		return gf_err
	}

	return nil

}