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
	"strings"
	"strconv"
	"net/url"
	"net/http"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func Jobs_mngr__init_handlers(p_jobs_mngr_ch chan Job_msg,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_jobs_handlers.Jobs_mngr__init_handlers()")

	//---------------------
	//running_jobs_map := map[string]*Running_job{}

	//START_JOB
	http.HandleFunc("/images/jobs/start",func(p_resp http.ResponseWriter,
											p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/jobs/start ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			//INPUT
			input_map,gf_err := gf_rpc_lib.Get_http_input("/images/jobs/start",
											p_resp,
											p_req,
											p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/jobs/start",
								"failed parse input to start a job", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}
			//--------------------------
			p_runtime_sys.Log_fun("INFO","input_map - "+fmt.Sprint(input_map))

			job_type_str                           := input_map["job_type_str"].(string)
			client_type_str                        := input_map["client_type_str"].(string)
			url_encoded_imgs_urls_str              := input_map["imgs_urls_str"].(string)
			url_encoded_imgs_origin_pages_urls_str := input_map["imgs_origin_pages_urls_str"].(string)

			p_runtime_sys.Log_fun("INFO","job_type_str    - "+job_type_str)
			p_runtime_sys.Log_fun("INFO","client_type_str - "+client_type_str)
		
			//-------------------
			//IMAGES_TO_PROCESS
			
			images_urls_str,_ := url.QueryUnescape(url_encoded_imgs_urls_str)
			images_urls_lst   := strings.Split(images_urls_str,",")

			imgs_origin_pages_urls_str,_ := url.QueryUnescape(url_encoded_imgs_origin_pages_urls_str)
			imgs_origin_pages_urls_lst   := strings.Split(imgs_origin_pages_urls_str,",")
			
			p_runtime_sys.Log_fun("INFO","url_encoded_imgs_urls_str - "+url_encoded_imgs_urls_str)
			p_runtime_sys.Log_fun("INFO","url_encoded_imgs_origin_pages_urls_str - "+url_encoded_imgs_origin_pages_urls_str)

			images_to_process := []Image_to_process{}
			for i,image_url_str := range images_urls_lst {

				image_origin_page_url_str := imgs_origin_pages_urls_lst[i]
				img_to_process            := Image_to_process{
					Source_url_str     :image_url_str,
					Origin_page_url_str:image_origin_page_url_str,
				}
				images_to_process = append(images_to_process,img_to_process)
			}
			//-------------------

			flows_names_lst := []string{"general",}
			running_job,job_expected_outputs_lst,gf_err := Start_job(client_type_str,
															images_to_process, //p_images_to_process_lst
															flows_names_lst,
															p_jobs_mngr_ch,
															p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/images/jobs/start",
								"failed starting a job", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}
			
 			//------------------
			//OUTPUT
			data_map := map[string]interface{}{
				"running_job_id_str":      running_job.Id_str,
				"job_expected_outputs_lst":job_expected_outputs_lst,
			}
			gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
			//------------------
			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/jobs/start",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//---------------------
	//GET_JOB_STATUS - SSE
	http.HandleFunc("/images/jobs/status",func(p_resp http.ResponseWriter,
											p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST -- /images/jobs/status ----------")

		if p_req.Method == "GET" {
			
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//-------------------------
			images_job_id_str := p_req.URL.Query().Get("images_job_id_str")
			p_runtime_sys.Log_fun("INFO","images_job_id_str - "+images_job_id_str)

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

			job_updates_ch := get_running_job_update_ch(images_job_id_str,p_jobs_mngr_ch,p_runtime_sys)
			//-------------------------

			//don't close the connection,
			//sending messages and flushing the response each time
			//there is a new message to send along.
			//
			//NOTE: we could loop endlessly; however, then you
			//could not easily detect clients that dettach and the
			//server would continue to send them messages long after
			//they're gone due to the "keep-alive" header.  One of
			//the nifty aspects of SSE is that clients automatically
			//reconnect when they lose their connection.
			//
			//a better way to do this is to use the CloseNotifier
			//interface that will appear in future releases of
			//Go (this is written as of 1.0.3):
			//https://code.google.com/p/go/source/detail?name=3292433291b2

			flusher,ok := p_resp.(http.Flusher)
			if !ok {
				err_msg_str := "/images/jobs/status handler failed - SSE http streaming is not supported on the server"
				gf_rpc_lib.Error__in_handler("/images/jobs/status",
								err_msg_str,
								nil,p_resp,p_runtime_sys)
				return
			}

			notify := p_resp.(http.CloseNotifier).CloseNotify()
			go func() {
				<- notify
				p_runtime_sys.Log_fun("ERROR","HTTP connection just closed")
			}()

			p_resp.Header().Set("Content-Type"               ,"text/event-stream")
			p_resp.Header().Set("Cache-Control"              ,"no-cache")
			p_resp.Header().Set("Connection"                 ,"keep-alive")
			p_resp.Header().Set("Access-Control-Allow-Origin","*")

			for {
				job_update,more_bool := <- job_updates_ch //running_job.job_updates_ch
				if more_bool {
					sse_event__unix_time_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0,'f',10,64)
					sse_event_id_str         := sse_event__unix_time_str

					fmt.Fprintf(p_resp,"id: %s\n",sse_event_id_str)

					job_update_lst,_ := json.Marshal(job_update)
					fmt.Fprintf(p_resp,"data: %s\n\n",job_update_lst)

					flusher.Flush()
				} else {
					//channel has been closed, so exit the loop
					break
				}
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/images/jobs/status",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//---------------------
}