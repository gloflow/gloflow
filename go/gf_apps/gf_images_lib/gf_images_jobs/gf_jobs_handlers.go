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

package gf_images_jobs

import (
	"fmt"
	"time"
	"context"
	"strings"
	"strconv"
	"net/url"
	"net/http"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
)

//-------------------------------------------------

func InitHandlers(pMux *http.ServeMux,
	pJobsMngrCh gf_images_jobs_core.JobsMngr,
	pRuntimeSys *gf_core.RuntimeSys) {

	//---------------------
	// START_JOB
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/jobs/start",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//--------------------------
				// INPUT
				inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				pRuntimeSys.LogFun("INFO", "input_map - "+fmt.Sprint(inputMap))

				job_type_str                       := inputMap["job_type_str"].(string)
				clientTypeStr                      := inputMap["client_type_str"].(string)
				urlEncodedImagesURLsStr            := inputMap["imgs_urls_str"].(string)
				urlEncodedImagesOriginPagesURLsStr := inputMap["imgs_origin_pages_urls_str"].(string)

				// ADD!! - accept this flows_names argument from http arguments, not hardcoded as is here
				flowsNamesLst := []string{"general",}

				pRuntimeSys.LogFun("INFO", fmt.Sprintf("job_type_str    - %s", job_type_str))
				pRuntimeSys.LogFun("INFO", fmt.Sprintf("client_type_str - %s", clientTypeStr))
				pRuntimeSys.LogFun("INFO", fmt.Sprintf("flows_names_lst - %s", flowsNamesLst))

				//-------------------
				// IMAGES_TO_PROCESS
				
				imagesURLsStr, _ := url.QueryUnescape(urlEncodedImagesURLsStr)
				imagesURLsLst    := strings.Split(imagesURLsStr, ",")

				imgs_origin_pages_urls_str, _ := url.QueryUnescape(urlEncodedImagesOriginPagesURLsStr)
				imgs_origin_pages_urls_lst    := strings.Split(imgs_origin_pages_urls_str, ",")
				
				pRuntimeSys.LogFun("INFO", "url_encoded_imgs_urls_str - "+urlEncodedImagesURLsStr)
				pRuntimeSys.LogFun("INFO", "url_encoded_imgs_origin_pages_urls_str - "+urlEncodedImagesOriginPagesURLsStr)

				imagesToProcessLst := []gf_images_jobs_core.GFimageExternToProcess{}
				for i, imageURLstr := range imagesURLsLst {

					imageOriginPageURLstr := imgs_origin_pages_urls_lst[i]

					img_to_process := gf_images_jobs_core.GFimageExternToProcess{
						SourceURLstr:     imageURLstr,
						OriginPageURLstr: imageOriginPageURLstr,
					}
					imagesToProcessLst = append(imagesToProcessLst, img_to_process)
				}

				//-------------------

				runningJob, jobExpectedOutputsLst, gfErr := gf_images_jobs_client.RunExternImages(clientTypeStr,
					imagesToProcessLst,
					flowsNamesLst,
					pJobsMngrCh,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				//------------------
				// OUTPUT
				outputMap := map[string]interface{}{
					"running_job_id_str":       runningJob.Id_str,
					"job_expected_outputs_lst": jobExpectedOutputsLst,
				}
				return outputMap, nil

				//------------------
			}
			return nil, nil
		},
		pMux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// GET_JOB_STATUS - SSE

	// DEPRECATED!! - use the new event streaming method general to GF,
	//                not this specific one for image jobs.
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/jobs/status",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {
				
				//-------------------------
				images_job_id_str := pReq.URL.Query().Get("images_job_id_str")
				pRuntimeSys.LogFun("INFO", "images_job_id_str - "+images_job_id_str)

				/*if _,ok := running_jobs_map[images_job_id_str]; !ok {
					gf_rpc_lib.Error__in_handler("/images/jobs/status",
									nil, //err,
									"job with ID doesnt exist - "+images_job_id_str, //p_user_msg_str
									pResp,
									p_mongodb_coll,
									pLogFun)
					return
				}

				running_job := running_jobs_map[images_job_id_str]*/

				jobUpdatesCh := gf_images_jobs_client.GetJobUpdateCh(images_job_id_str, pJobsMngrCh, pRuntimeSys)

				//-------------------------

				// don't close the connection, sending messages and flushing the response each time there is a new message to send along.
				// NOTE: we could loop endlessly; however, then you could not easily detect clients that dettach and the
				// server would continue to send them messages long after they're gone due to the "keep-alive" header.  One of
				// the nifty aspects of SSE is that clients automatically reconnect when they lose their connection.
				// a better way to do this is to use the CloseNotifier interface that will appear in future releases of
				// Go (this is written as of 1.0.3):
				// https://code.google.com/p/go/source/detail?name=3292433291b2

				flusher, ok := pResp.(http.Flusher)
				if !ok {
					err_msg_str := "/images/jobs/status handler failed - SSE http streaming is not supported on the server"
					gf_rpc_lib.ErrorInHandler("/images/jobs/status", err_msg_str, nil, pResp, pRuntimeSys)
					return nil, nil
				}

				notify := pResp.(http.CloseNotifier).CloseNotify()
				go func() {
					<- notify
					pRuntimeSys.LogFun("ERROR", "HTTP connection just closed")
				}()

				pResp.Header().Set("Content-Type",                "text/event-stream")
				pResp.Header().Set("Cache-Control",               "no-cache")
				pResp.Header().Set("Connection",                  "keep-alive")
				pResp.Header().Set("Access-Control-Allow-Origin", "*") // CORS

				for {

					// IMPORTANT!! - jobs_mngr, after processing a job and sending all status messages that it has to send, 
					//               does NOT close the channel. if the status messages client receiver, via this HTTP (SSE) handler
					//               is slow to connect/consume messages, then the jobs_mngr will complete the job (and close the channel)
					//               faster then the client here will consume them, and so a closed channel will be encoutered before all
					//               its messages were consumed.
					//               to avoid this a channel is not closed, and instead the last update message is waited for here, or an error,
					//               and cleanup only done after that.
					job_update := <- jobUpdatesCh
					
					sse_event__unix_time_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)
					sse_event_id_str         := sse_event__unix_time_str

					fmt.Fprintf(pResp, "id: %s\n", sse_event_id_str)

					job_update_lst, _ := json.Marshal(job_update)
					fmt.Fprintf(pResp, "data: %s\n\n", job_update_lst)

					flusher.Flush()

					// IMPORTANT!! - this is the last update message, so exit the loop
					if job_update.Type_str == gf_images_jobs_core.JOB_UPDATE_TYPE__ERROR || job_update.Type_str == gf_images_jobs_core.JOB_UPDATE_TYPE__COMPLETED {

						gf_images_jobs_client.CleanupJob(images_job_id_str, pJobsMngrCh, pRuntimeSys)
						break
					}
				}
				return nil, nil
			}
			return nil, nil
		},
		pMux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//---------------------
}