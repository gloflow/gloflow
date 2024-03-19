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
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_plugins"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type JobsMngr chan JobMsg

type GFjobRunning struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	IDstr           string        `bson:"id_str"`
	T_str           string        `bson:"t"`
	Client_type_str string        `bson:"client_type_str"`
	Status_str      string        `bson:"status_str"` // "running"|"complete"
	Start_time_f    float64       `bson:"start_time_f"`
	End_time_f      float64       `bson:"end_time_f"`

	// LEGACY!! - update "images_to_process_lst" to "images_extern_to_process_lst" in the DB (bson)
	Images_extern_to_process_lst   []GFimageExternToProcess   `bson:"images_to_process_lst"`
	Images_uploaded_to_process_lst []GFimageUploadedToProcess `bson:"images_uploaded_to_process_lst"`
	
	Errors_lst     []JobError        `bson:"errors_lst"`
	job_updates_ch chan JobUpdateMsg `bson:"-"`

	// user that submited the job
	UserID gf_core.GF_ID `bson:"user_id"`
}

type GFjobRuntime struct {
	job_id_str              string
	job_client_type_str     string
	job_updates_ch          chan JobUpdateMsg
	useNewStorageEngineBool bool
	metricsPlugins          *gf_images_plugins.GFmetrics
	userID                  gf_core.GF_ID
}

//------------------------
// IMAGES_TO_PROCESS

type GFimageExternToProcess struct {
	SourceURLstr     string `bson:"source_url_str"` // FIX!! - rename this to Origin_url_str to be consistent with other origin_url naming
	OriginPageURLstr string `bson:"origin_page_url_str"`
}

type GFimageUploadedToProcess struct {
	GFimageIDstr  gf_images_core.GFimageID
	S3filePathStr string                 // path to image in S3 in a bucket that it was originally uploaded to by client
	MetaMap       map[string]interface{} // metadata user might include for this image
}

type GFimageLocalToProcess struct {
	LocalFilePathStr string
}

type GFimageClassificationToProcess struct {
	GFimageIDstr gf_images_core.GFimageID
}

//------------------------
// JOB_MSGS

type JobMsg struct {
	Job_id_str                     string // if its an existing job. for new jobs the mngr creates a new job ID
	UserID                         gf_core.GF_ID
	Client_type_str                string 
	Cmd_str                        string // "start_job" | "get_running_job_ids"
	
	Job_init_ch                    chan *GFjobRunning   // used by clients for receiving outputs of job initialization by jobs_mngr
	Job_updates_ch                 chan JobUpdateMsg  // used by jobs_mngr to send job_updates to
	Msg_response_ch                chan interface{}     // DEPRECATED!! use a specific struct as a message format, interface{} too general.

	Images_extern_to_process_lst   []GFimageExternToProcess
	Images_uploaded_to_process_lst []GFimageUploadedToProcess
	Images_local_to_process_lst    []GFimageLocalToProcess
	ImagesToClassifyLst            []GFimageClassificationToProcess

	Flows_names_lst                []string
}

type JobUpdateMsg struct {
	Name_str             string                        `json:"name_str"`
	Type_str             job_update_type_val           `json:"type_str"`             // "ok" | "error" | "complete"
	ImageIDstr           gf_images_core.GFimageID      `json:"image_id_str"`
	Image_source_url_str string                        `json:"image_source_url_str"`
	Err_str              string                        `json:"err_str,omitempty"`    // if the update indicates an error, this is its value
	Image_thumbs         *gf_images_core.GFimageThumbs `json:"-"`
}

//------------------------
// JOBS_LIFECYCLE

type GF_jobs_lifecycle_callbacks struct {
	Job_type__transform_imgs__fun func() *gf_core.GFerror
	Job_type__uploaded_imgs__fun  func() *gf_core.GFerror
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
// INIT

func JobsMngrInit(pImagesStoreLocalDirPathStr string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pVideoStoreLocalDirPathStr            string,
	pMediaDomainStr                       string,
	pLifecycleCallbacks                   *GF_jobs_lifecycle_callbacks,
	pConfig                               *gf_images_core.GFconfig,
	pImageStorage                         *gf_images_storage.GFimageStorage,
	pS3info                               *gf_aws.GFs3Info,
	pMetricsCore                          *gf_images_core.GFmetrics,
	pRuntimeSys                           *gf_core.RuntimeSys) JobsMngr {

	jobsMngrCh := make(chan JobMsg, 100)


	
	



	// IMPORTANT!! - start jobs_mngr as an independent goroutine of the HTTP handlers at
	//               service initialization time
	go func() {
		
		//---------------------
		// JOBS_STATUS
		// IMPORTANT!! - version of streaming jobs status endpoint to move all clients to.
		//			version in gf_jobs_handlers.go is deprecated.

		SSEurlStr := "/v1/images/jobs/status"
		eventsCtx := gf_events.Init(SSEurlStr, pRuntimeSys)
		
		fmt.Println(eventsCtx)

		//---------------------
		
		//---------------------
		// METRICS 
		metricsJobs := MetricsCreate("gf_images_jobs")

		//---------------------
		
		runningJobsMap := map[string]chan JobUpdateMsg{}


		//-------------------------------------------------
		regularJobInitFun := func(pJobMsg JobMsg) (*GFjobRuntime, *GFjobRunning, *gf_core.GFerror) {

			// RUNNING_JOB
			runningJob, gfErr := JobsMngrCreateRunningJob(pJobMsg.Client_type_str,
				pJobMsg.Job_updates_ch,
				pJobMsg.UserID,
				pRuntimeSys)
			if gfErr != nil {
				return nil, nil, gfErr
			}

			// IMPORTANT!! - send sending running_job back to client, to avoid race conditions
			runningJobsMap[runningJob.IDstr] = pJobMsg.Job_updates_ch

			// SEND_MSG
			pJobMsg.Job_init_ch <- runningJob


			jobRuntime := &GFjobRuntime{
				job_id_str:          runningJob.IDstr,
				job_client_type_str: pJobMsg.Client_type_str,
				job_updates_ch:      pJobMsg.Job_updates_ch,
				useNewStorageEngineBool: pConfig.UseNewStorageEngineBool,
				userID:                  pJobMsg.UserID,
			}



			return jobRuntime, runningJob, nil
		}
		

		//-------------------------------------------------

		// listen to messages
		for {
			jobMsg := <- jobsMngrCh

			ctx := context.Background()

			// IMPORTANT!! - only one job is processed per jobs_mngr.
			//              Scaling is done with multiple jobs_mngr's (exp. per-core)           
			switch jobMsg.Cmd_str {


				//------------------------
				// START_JOB_CLASSIFY_IMAGES

				case "start_job_classify_imgs":


					//------------------------
					// METRICS
					metricsJobs.CmdStartJobClassifyImagesCount.Inc()

					jobRuntime, runningJob, gfErr := regularJobInitFun(jobMsg)
					if gfErr != nil {
						continue
					}

					//------------------------

					imagesLst := jobMsg.ImagesToClassifyLst

					gfErr = runJobClassifyImages(imagesLst,
						pConfig.ImagesClassifyPyDirPathStr,
						pImageStorage,
						jobRuntime,
						pMetricsCore,
						ctx,
						pRuntimeSys)
			

					//------------------------
					// JOB_STATUS
					var jobStatusStr job_status_val
					if gfErr != nil {
						jobStatusStr = JOB_STATUS__FAILED
					} else {
						jobStatusStr = JOB_STATUS__COMPLETED
					}
					_ = dbJobsMngrUpdateJobStatus(jobStatusStr, runningJob.IDstr, pRuntimeSys)

					//------------------------



				//------------------------
				// START_JOB_LOCAL_IMAGES
				case "start_job_local_imgs":

					
					// METRICS
					metricsJobs.Cmd__start_job_local_imgs__count.Inc()

					// RUNNING_JOB
					runningJob, gfErr := JobsMngrCreateRunningJob(jobMsg.Client_type_str,
						jobMsg.Job_updates_ch,
						jobMsg.UserID,
						pRuntimeSys)
					if gfErr != nil {
						continue
					}

					// IMPORTANT!! - send sending running_job back to client, to avoid race conditions
					runningJobsMap[runningJob.IDstr] = jobMsg.Job_updates_ch

					// SEND_MSG
					jobMsg.Job_init_ch <- runningJob

					//------------------------

					jobRuntime := &GFjobRuntime{
						job_id_str:              runningJob.IDstr,
						job_client_type_str:     jobMsg.Client_type_str,
						job_updates_ch:          jobMsg.Job_updates_ch,
						useNewStorageEngineBool: pConfig.UseNewStorageEngineBool,
						metricsPlugins:          metricsJobs.ImagesPluginsMetrics,
						userID:                  jobMsg.UserID,
					}

					runJobErrsLst := runJobLocalImages(jobMsg.Images_local_to_process_lst,
						jobMsg.Flows_names_lst,
						pImagesStoreLocalDirPathStr,
						pImagesThumbnailsStoreLocalDirPathStr,
						pS3info,
						pConfig.PluginsPyDirPathStr,
						pImageStorage,
						jobRuntime,
						pRuntimeSys)
					
					//------------------------
					// JOB_STATUS
					var jobStatusStr job_status_val
					if len(runJobErrsLst) == len(jobMsg.Images_local_to_process_lst) {
						jobStatusStr = JOB_STATUS__FAILED
					} else if len(runJobErrsLst) > 0 {
						jobStatusStr = JOB_STATUS__FAILED_PARTIAL
					} else {
						jobStatusStr = JOB_STATUS__COMPLETED
					}
					_ = dbJobsMngrUpdateJobStatus(jobStatusStr, runningJob.IDstr, pRuntimeSys)

					//------------------------

				//------------------------
				// START_JOB_TRANSFORM_IMAGES

				case "start_job_transform_imgs":

					// METRICS
					metricsJobs.Cmd__start_job_transform_imgs__count.Inc()

					/*
					// RUST
					// FIX!! - this just runs Rust job code for testing.
					//         pass in proper job_cmd argument.
					run_job_rust()
					*/

					//------------------------
					// LIFECYCLE_CALLBACK
					if pLifecycleCallbacks != nil {
						gfErr := pLifecycleCallbacks.Job_type__transform_imgs__fun()
						if gfErr != nil {
							continue
						}
					}

					//------------------------

				//------------------------
				// START_JOB_UPLOADED_IMAGES

				case "start_job_uploaded_imgs":
					
					// METRICS
					metricsJobs.Cmd__start_job_uploaded_imgs__count.Inc()

					// RUNNING_JOB
					runningJob, gfErr := JobsMngrCreateRunningJob(jobMsg.Client_type_str,
						jobMsg.Job_updates_ch,
						jobMsg.UserID,
						pRuntimeSys)
					if gfErr != nil {
						continue
					}

					// IMPORTANT!! - send sending running_job back to client, to avoid race conditions
					runningJobsMap[runningJob.IDstr] = jobMsg.Job_updates_ch

					// SEND_MSG
					jobMsg.Job_init_ch <- runningJob

					jobRuntime := &GFjobRuntime{
						job_id_str:          runningJob.IDstr,
						job_client_type_str: jobMsg.Client_type_str,
						job_updates_ch:      jobMsg.Job_updates_ch,
						useNewStorageEngineBool: pConfig.UseNewStorageEngineBool,
						metricsPlugins:          metricsJobs.ImagesPluginsMetrics,
						userID:                  jobMsg.UserID,
					}

					runJobErrsLst := runJobUploadedImages(jobMsg.Images_uploaded_to_process_lst,
						jobMsg.Flows_names_lst,
						pImagesStoreLocalDirPathStr,
						pImagesThumbnailsStoreLocalDirPathStr,
						pS3info,
						pConfig.PluginsPyDirPathStr,
						pImageStorage,
						jobRuntime,
						pRuntimeSys)
					
					//------------------------
					// JOB_STATUS
					var jobStatusStr job_status_val
					if len(runJobErrsLst) == len(jobMsg.Images_uploaded_to_process_lst) {
						jobStatusStr = JOB_STATUS__FAILED
					} else if len(runJobErrsLst) > 0 {
						jobStatusStr = JOB_STATUS__FAILED_PARTIAL
					} else {
						jobStatusStr = JOB_STATUS__COMPLETED
					}
					_ = dbJobsMngrUpdateJobStatus(jobStatusStr, runningJob.IDstr, pRuntimeSys)

					//------------------------
					// LIFECYCLE_CALLBACK
					if pLifecycleCallbacks != nil {
						gfErr = pLifecycleCallbacks.Job_type__uploaded_imgs__fun()
						if gfErr != nil {
							continue
						}
					}
				
				//------------------------
				// START_JOB_EXTERN_IMAGES
				// FIX!! - "start_job" needs to be "start_job_extern_imgs", update in all clients.

				case "start_job":

					// METRICS
					metricsJobs.Cmd__start_job_extern_imgs__count.Inc()

					// RUNNING_JOB
					runningJob, gfErr := JobsMngrCreateRunningJob(jobMsg.Client_type_str,
						jobMsg.Job_updates_ch,
						jobMsg.UserID,
						pRuntimeSys)
					if gfErr != nil {
						continue
					}

					// IMPORTANT!! - send sending running_job back to client, to avoid race conditions
					runningJobsMap[runningJob.IDstr] = jobMsg.Job_updates_ch

					// SEND_MSG
					jobMsg.Job_init_ch <- runningJob
					
					//------------------------
					// ADD!! - due to legacy reasons the "general" flow is still used as the main flow
					//         that images are added to when a external_images job is processed.
					//         this should be generalized so that images are added to dedicated flow S3 buckets
					//         if those flows have their S3_bucket mapping defined in Gf_config.Images_flow_to_s3_bucket_map
					s3bucketNameStr := pConfig.ImagesFlowToS3bucketMap["general"]

					//------------------------

					jobRuntime := &GFjobRuntime{
						job_id_str:          runningJob.IDstr,
						job_client_type_str: jobMsg.Client_type_str,
						job_updates_ch:      jobMsg.Job_updates_ch,
						useNewStorageEngineBool: pConfig.UseNewStorageEngineBool,
						metricsPlugins:          metricsJobs.ImagesPluginsMetrics,
						userID:                  jobMsg.UserID,
					}

					runJobErrsLst := runJobExternImages(jobMsg.Images_extern_to_process_lst,
						jobMsg.Flows_names_lst,
						pImagesStoreLocalDirPathStr,
						pImagesThumbnailsStoreLocalDirPathStr,
						pVideoStoreLocalDirPathStr,
						pMediaDomainStr,
						s3bucketNameStr,
						pS3info,
						pConfig.PluginsPyDirPathStr,
						pImageStorage,
						jobRuntime,
						ctx,
						pRuntimeSys)
					
					//------------------------
					// JOB_STATUS
					var jobStatusStr job_status_val
					if len(runJobErrsLst) == len(jobMsg.Images_extern_to_process_lst) {
						jobStatusStr = JOB_STATUS__FAILED
					} else if len(runJobErrsLst) > 0 {
						jobStatusStr = JOB_STATUS__FAILED_PARTIAL
					} else {
						jobStatusStr = JOB_STATUS__COMPLETED
					}
					_ = dbJobsMngrUpdateJobStatus(jobStatusStr, runningJob.IDstr, pRuntimeSys)

					//------------------------

				//------------------------
				// GET_JOB_UPDATE_CH

				case "get_job_update_ch":

					jobIDstr := jobMsg.Job_id_str

					if _, ok := runningJobsMap[jobIDstr]; ok {

						jobUpdatesCh := runningJobsMap[jobIDstr]
						jobMsg.Msg_response_ch <- jobUpdatesCh

					} else {
						jobMsg.Msg_response_ch <- nil
					}

				//------------------------
				// CLEANUP_JOB

				case "cleanup_job":
					jobIDstr := jobMsg.Job_id_str
					delete(runningJobsMap, jobIDstr) // remove running job from lookup, since its complete

				//------------------------
			} 
		}
	}()
	return jobsMngrCh
}

//-------------------------------------------------
// CREATE_RUNNING_JOB

func JobsMngrCreateRunningJob(pClientTypeStr string,
	p_job_updates_ch chan JobUpdateMsg,
	pUserID          gf_core.GF_ID,
	pRuntimeSys      *gf_core.RuntimeSys) (*GFjobRunning, *gf_core.GFerror) {

	if pUserID == "" {
		gfErr := gf_core.ErrorCreate("user_id for this job has not been specified (its empty)",
			"user_incorrect",
			map[string]interface{}{
				"client_type_str": pClientTypeStr,
			},
			nil, "gf_images_jobs_core", pRuntimeSys)
		return nil, gfErr
	}
	
	job_start_time_f := float64(time.Now().UnixNano())/1000000000.0
	jobIDstr         := fmt.Sprintf("job:%f", job_start_time_f)

	runningJob := &GFjobRunning{
		IDstr:           jobIDstr,
		T_str:           "img_running_job",
		Client_type_str: pClientTypeStr,
		Status_str:      "running",
		Start_time_f:    job_start_time_f,
		job_updates_ch:  p_job_updates_ch,
		UserID:          pUserID,
	}

	// DB
	gfErr := dbJobsMngrCreateRunningJob(runningJob, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return runningJob, nil
}

//-------------------------------------------------
// DB
//-------------------------------------------------

func dbJobsMngrCreateRunningJob(pRunningJob *GFjobRunning,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx           := context.Background()
	collNameStr := "gf_images__jobs_running" // pRuntimeSys.Mongo_coll.Name()
	gfErr        := gf_core.MongoInsert(pRunningJob,
		collNameStr,
		map[string]interface{}{
			"running_job_id_str": pRunningJob.IDstr,
			"client_type_str":    pRunningJob.Client_type_str,
			"caller_err_msg_str": "failed to create a Running_job record into the DB",
		},
		ctx,
		pRuntimeSys)
	
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//-------------------------------------------------

func dbJobsMngrUpdateJobStatus(p_status_str job_status_val,
	p_job_id_str  string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	if p_status_str != JOB_STATUS__COMPLETED && p_status_str != JOB_STATUS__FAILED {
		// status values are not generated at runtime, but are static, so its ok to panic here since
		// this should never be countered in production
		panic(fmt.Sprintf("job status value thats not allowed - %s", p_status_str))
	}

	ctx := context.Background()

	job_end_time_f := float64(time.Now().UnixNano())/1000000000.0
	coll           := pRuntimeSys.Mongo_db.Collection("gf_images__jobs_running")
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
		gfErr := gf_core.MongoHandleError("failed to update an img_running_job in the DB, as complete and its end_time",
			"mongodb_update_error",
			map[string]interface{}{
				"job_id_str":     p_job_id_str,
				"job_end_time_f": job_end_time_f,
			},
			err, "gf_jobs_mngr", pRuntimeSys)
		return gfErr
	}

	return nil

}