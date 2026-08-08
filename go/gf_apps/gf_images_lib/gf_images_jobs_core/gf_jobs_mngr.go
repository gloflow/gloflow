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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	// "github.com/davecgh/go-spew/spew"
)

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


	
	

	pRuntimeSys.LogNewFun("DEBUG", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", nil)
	pRuntimeSys.LogNewFun("DEBUG", "image storage config...", map[string]interface{}{
		"type": pImageStorage.TypeStr,
	})
	



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
			switch jobMsg.CmdStr {


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

					classesLst, gfErr := runJobClassifyImages(imagesLst,
						pConfig.ImagesClassifyPyDirPathStr,
						pImageStorage,
						jobRuntime,
						pMetricsCore,
						ctx,
						pRuntimeSys)
					if gfErr != nil {
						continue
					}
					
					//------------------------
					// OUTPUT
					// synchronously send the response back to the client (unlike most other job types
					// that dont send any sync output back to clients).

					jobMsg.ResponseCh <- classesLst

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
						jobMsg.JPEGqualityInt,
						jobMsg.ImageOutputFormat,
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

					jobIDstr := jobMsg.JobIDstr

					if _, ok := runningJobsMap[jobIDstr]; ok {

						jobUpdatesCh := runningJobsMap[jobIDstr]
						jobMsg.ResponseCh <- jobUpdatesCh

					} else {
						jobMsg.ResponseCh <- nil
					}

				//------------------------
				// CLEANUP_JOB

				case "cleanup_job":
					jobIDstr := jobMsg.JobIDstr
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