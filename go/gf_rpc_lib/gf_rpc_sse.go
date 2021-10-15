/*
MIT License

Copyright (c) 2021 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_rpc_lib

import (
	"fmt"
	"context"
	"strconv"
	"time"
	"net/http"
	"strings"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type SSE_data_update_ch   chan interface{}
type SSE_data_err_ch      chan gf_core.GF_error
type SSE_data_complete_ch chan bool
type handler_http_sse     func(context.Context, http.ResponseWriter, *http.Request) (SSE_data_update_ch, SSE_data_err_ch, SSE_data_complete_ch, *gf_core.GF_error)

//-------------------------------------------------
func SSE_create_handler__http(p_path_str string,
	p_handler_fun    handler_http_sse,
	p_store_run_bool bool,
	p_runtime_sys    *gf_core.Runtime_sys) {

	Create_handler__http_with_metrics(p_path_str,
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//-----------------------
				// HANDLER
				data_updates_ch, data_err_ch, data_complete_ch, gf_err := p_handler_fun(p_ctx, p_resp, p_req)
				if gf_err != nil {
					return nil, gf_err
				}

				//-----------------------

				// don't close the connection, sending messages and flushing the response each time there is a new message to send along.
				// NOTE: we could loop endlessly; however, then you could not easily detect clients that dettach and the
				// server would continue to send them messages long after they're gone due to the "keep-alive" header.  One of
				// the nifty aspects of SSE is that clients automatically reconnect when they lose their connection.
				// a better way to do this is to use the CloseNotifier interface that will appear in future releases of
				// Go (this is written as of 1.0.3):
				// https://code.google.com/p/go/source/detail?name=3292433291b2

				flusher,ok := p_resp.(http.Flusher)
				if !ok {
					err_msg_str := fmt.Sprintf("%s handler failed - SSE http streaming is not supported on the server", p_path_str)
					gf_err := gf_core.Error__create(err_msg_str, 
						"http_server_flusher_not_supported_error",
						map[string]interface{}{"path_str": p_path_str,},
						nil, "gf_rpc_lib", p_runtime_sys)
					return nil, gf_err
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

				//-------------------------------------------------
				complete_fun := func() {
					data_complete_lst, _ := json.Marshal(map[string]interface{}{"status_str": "complete"})
					fmt.Fprintf(p_resp, "data: %s\n\n", data_complete_lst)
					flusher.Flush()
				}

				//-------------------------------------------------

				for {
					select {

					//-----------------------
					// DATA_UPDATE
					// IMPORTANT!! - jobs_mngr, after processing a job and sending all status messages that it has to send, 
					//               does NOT close the channel. if the status messages client receiver, via this HTTP (SSE) handler
					//               is slow to connect/consume messages, then the jobs_mngr will complete the job (and close the channel)
					//               faster then the client here will consume them, and so a closed channel will be encoutered before all
					//               its messages were consumed.
					//               to avoid this a channel is not closed, and instead the last update message is waited for here, or an error,
					//               and cleanup only done after that.
					case data_update := <- data_updates_ch:

						sse_event__unix_time_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)
						sse_event_id_str         := sse_event__unix_time_str

						event_map := map[string]interface{}{
							"data_map":   data_update,
							"status_str": "ok",
						}
						sse_event_lst, _ := json.Marshal(event_map)

						// SSE_WRITE
						fmt.Fprintf(p_resp, "id: %s\n", sse_event_id_str)
						fmt.Fprintf(p_resp, "data: %s\n\n", sse_event_lst)
						flusher.Flush()
					
					//-----------------------
					// ERROR_UPDATE
					case gf_err := <- data_err_ch:

						sse_event__unix_time_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)
						sse_event_id_str         := sse_event__unix_time_str
						
						event_map := map[string]interface{}{
							"data_map":   gf_err,
							"status_str": "error",
						}
						sse_event_lst, _ := json.Marshal(event_map)

						// SSE_WRITE
						fmt.Fprintf(p_resp, "id: %s\n", sse_event_id_str)
						fmt.Fprintf(p_resp, "data: %s\n\n", sse_event_lst)
						flusher.Flush()

						// an error occured, so complete SSE stream and exit
						complete_fun()
						break

					//-----------------------
					// COMPLETE
					case complete_bool := <- data_complete_ch:

						if complete_bool {

							complete_fun()
							break
						}

					//-----------------------

					}	
				}
			}

			return nil, nil
		},
		nil,
		p_store_run_bool,
		p_runtime_sys)
}

//-------------------------------------------------
func SSE_client__parse_response(p_body_str string,
	p_runtime_sys *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.GF_error) {
	// p_runtime_sys.Log_fun("FUN_ENTER", "gf_rpc_sse.SSE_client__parse_response()")

	data_items_lst := []map[string]interface{}{}

	for _, line_str := range strings.Split(p_body_str, `\n`) {

		p_runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>")
		p_runtime_sys.Log_fun("INFO", line_str)

		// filter out keep-alive new lines
		if line_str != "" && strings.HasPrefix(line_str, "data: ") {


			msg_str := strings.Replace(line_str, "data: ", "", 1)
			msg_map := map[string]interface{}{}
			err     := json.Unmarshal([]byte(msg_str), &msg_map)

			if err != nil {

				gf_err := gf_core.Error__create("failed to parse JSON response line of the SSE stream (of even updates from a gf_images server)",
					"json_unmarshal_error",
					map[string]interface{}{"line_str": line_str,},
					err, "gf_images_lib", p_runtime_sys)

				return nil, gf_err
			}

			//-------------------
			// STATUS
			if _,ok := msg_map["status_str"]; !ok {
				err_usr_msg := "sse message json doesnt container key status_str"
				gf_err      := gf_core.Error__create(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", p_runtime_sys)
				return nil, gf_err
			}
			status_str := msg_map["status_str"].(string)

			if !(status_str == "ok" || status_str == "error") {

				err_usr_msg := "sse message json status_str key is not of value ok|error"
				gf_err      := gf_core.Error__create(err_usr_msg,
					"verify__invalid_key_value_error",
					map[string]interface{}{
						"status_str": status_str,
						"msg_map":    msg_map,
					},
					nil, "gf_images_lib", p_runtime_sys)
				return nil, gf_err
			}

			//-------------------
			// DATA
			if _,ok := msg_map["data_map"]; !ok {
				err_usr_msg := "sse message json doesnt container key data_map"
				gf_err      := gf_core.Error__create(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", p_runtime_sys)
				return nil, gf_err
			}
			
			data_map := msg_map["data_map"].(map[string]interface{})
			
			//-------------------
			data_items_lst = append(data_items_lst, data_map)
		}
	}
	return data_items_lst, nil
}