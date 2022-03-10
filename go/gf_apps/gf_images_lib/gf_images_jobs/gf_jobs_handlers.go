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
func Jobs_mngr__init_handlers(p_mux *http.ServeMux,
	p_jobs_mngr_ch gf_images_jobs_core.Jobs_mngr,
	p_runtime_sys  *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_handlers.Jobs_mngr__init_handlers()")

	//---------------------
	// running_jobs_map := map[string]*Running_job{}

	// START_JOB
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/jobs/start",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//--------------------------
				// INPUT
				input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				
				p_runtime_sys.Log_fun("INFO", "input_map - "+fmt.Sprint(input_map))

				job_type_str                           := input_map["job_type_str"].(string)
				client_type_str                        := input_map["client_type_str"].(string)
				url_encoded_imgs_urls_str              := input_map["imgs_urls_str"].(string)
				url_encoded_imgs_origin_pages_urls_str := input_map["imgs_origin_pages_urls_str"].(string)

				// ADD!! - accept this flows_names argument from http arguments, not hardcoded as is here
				flows_names_lst := []string{"general",}

				p_runtime_sys.Log_fun("INFO", fmt.Sprintf("job_type_str    - %s", job_type_str))
				p_runtime_sys.Log_fun("INFO", fmt.Sprintf("client_type_str - %s", client_type_str))
				p_runtime_sys.Log_fun("INFO", fmt.Sprintf("flows_names_lst - %s", flows_names_lst))

				//-------------------
				// IMAGES_TO_PROCESS
				
				images_urls_str,_ := url.QueryUnescape(url_encoded_imgs_urls_str)
				images_urls_lst   := strings.Split(images_urls_str, ",")

				imgs_origin_pages_urls_str,_ := url.QueryUnescape(url_encoded_imgs_origin_pages_urls_str)
				imgs_origin_pages_urls_lst   := strings.Split(imgs_origin_pages_urls_str, ",")
				
				p_runtime_sys.Log_fun("INFO", "url_encoded_imgs_urls_str - "+url_encoded_imgs_urls_str)
				p_runtime_sys.Log_fun("INFO", "url_encoded_imgs_origin_pages_urls_str - "+url_encoded_imgs_origin_pages_urls_str)

				images_to_process_lst := []gf_images_jobs_core.GF_image_extern_to_process{}
				for i, image_url_str := range images_urls_lst {

					image_origin_page_url_str := imgs_origin_pages_urls_lst[i]

					img_to_process := gf_images_jobs_core.GF_image_extern_to_process{
						Source_url_str:      image_url_str,
						Origin_page_url_str: image_origin_page_url_str,
					}
					images_to_process_lst = append(images_to_process_lst, img_to_process)
				}

				//-------------------

				running_job, job_expected_outputs_lst, gf_err := gf_images_jobs_client.RunExternImgs(client_type_str,
					images_to_process_lst,
					flows_names_lst,
					p_jobs_mngr_ch,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				
				//------------------
				// OUTPUT
				output_map := map[string]interface{}{
					"running_job_id_str":       running_job.Id_str,
					"job_expected_outputs_lst": job_expected_outputs_lst,
				}
				return output_map, nil

				//------------------
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)

	//---------------------
	// GET_JOB_STATUS - SSE

	// DEPRECATED!! - use the new event streaming method general to GF,
	//                not this specific one for image jobs.
	gf_rpc_lib.CreateHandlerHTTPwithMux("/images/jobs/status",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {
				
				//-------------------------
				images_job_id_str := p_req.URL.Query().Get("images_job_id_str")
				p_runtime_sys.Log_fun("INFO", "images_job_id_str - "+images_job_id_str)

				/*if _,ok := running_jobs_map[images_job_id_str]; !ok {
					gf_rpc_lib.Error__in_handler("/images/jobs/status",
									nil, //err,
									"job with ID doesnt exist - "+images_job_id_str, //p_user_msg_str
									p_resp,
									p_mongodb_coll,
									p_log_fun)
					return
				}

				running_job := running_jobs_map[images_job_id_str]*/

				job_updates_ch := gf_images_jobs_client.Job__get_update_ch(images_job_id_str, p_jobs_mngr_ch, p_runtime_sys)

				//-------------------------

				// don't close the connection, sending messages and flushing the response each time there is a new message to send along.
				// NOTE: we could loop endlessly; however, then you could not easily detect clients that dettach and the
				// server would continue to send them messages long after they're gone due to the "keep-alive" header.  One of
				// the nifty aspects of SSE is that clients automatically reconnect when they lose their connection.
				// a better way to do this is to use the CloseNotifier interface that will appear in future releases of
				// Go (this is written as of 1.0.3):
				// https://code.google.com/p/go/source/detail?name=3292433291b2

				flusher, ok := p_resp.(http.Flusher)
				if !ok {
					err_msg_str := "/images/jobs/status handler failed - SSE http streaming is not supported on the server"
					gf_rpc_lib.Error__in_handler("/images/jobs/status", err_msg_str, nil, p_resp, p_runtime_sys)
					return nil, nil
				}

				notify := p_resp.(http.CloseNotifier).CloseNotify()
				go func() {
					<- notify
					p_runtime_sys.Log_fun("ERROR", "HTTP connection just closed")
				}()

				p_resp.Header().Set("Content-Type",                "text/event-stream")
				p_resp.Header().Set("Cache-Control",               "no-cache")
				p_resp.Header().Set("Connection",                  "keep-alive")
				p_resp.Header().Set("Access-Control-Allow-Origin", "*") // CORS

				for {

					// IMPORTANT!! - jobs_mngr, after processing a job and sending all status messages that it has to send, 
					//               does NOT close the channel. if the status messages client receiver, via this HTTP (SSE) handler
					//               is slow to connect/consume messages, then the jobs_mngr will complete the job (and close the channel)
					//               faster then the client here will consume them, and so a closed channel will be encoutered before all
					//               its messages were consumed.
					//               to avoid this a channel is not closed, and instead the last update message is waited for here, or an error,
					//               and cleanup only done after that.
					job_update := <- job_updates_ch
					
					sse_event__unix_time_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)
					sse_event_id_str         := sse_event__unix_time_str

					fmt.Fprintf(p_resp, "id: %s\n", sse_event_id_str)

					job_update_lst, _ := json.Marshal(job_update)
					fmt.Fprintf(p_resp, "data: %s\n\n", job_update_lst)

					flusher.Flush()

					// IMPORTANT!! - this is the last update message, so exit the loop
					if job_update.Type_str == gf_images_jobs_core.JOB_UPDATE_TYPE__ERROR || job_update.Type_str == gf_images_jobs_core.JOB_UPDATE_TYPE__COMPLETED {

						gf_images_jobs_client.Job__cleanup(images_job_id_str, p_jobs_mngr_ch, p_runtime_sys)
						break
					}
				}
				return nil, nil
			}
			return nil, nil
		},
		p_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)
	
	//---------------------
}