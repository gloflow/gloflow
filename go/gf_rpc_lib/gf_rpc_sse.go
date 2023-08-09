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
type SSE_data_err_ch      chan gf_core.GFerror
type SSE_data_complete_ch chan bool
type handler_http_sse     func(context.Context, http.ResponseWriter, *http.Request) (SSE_data_update_ch, SSE_data_err_ch, SSE_data_complete_ch, *gf_core.GFerror)

//-------------------------------------------------

func CreateHandlerSSE(pPathStr string,
	pHandlerFun   handler_http_sse,
	pStoreRunBool bool,
	pRuntimeSys   *gf_core.RuntimeSys) {

	CreateHandlerHTTPwithMetrics(pPathStr,
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//-----------------------
				// HANDLER
				dataUpdatesCh, dataErrCh, dataCompleteCh, gfErr := pHandlerFun(pCtx, pResp, pReq)
				if gfErr != nil {
					return nil, gfErr
				}

				//-----------------------

				// don't close the connection, sending messages and flushing the response each time there is a new message to send along.
				// NOTE: we could loop endlessly; however, then you could not easily detect clients that dettach and the
				// server would continue to send them messages long after they're gone due to the "keep-alive" header.  One of
				// the nifty aspects of SSE is that clients automatically reconnect when they lose their connection.
				// a better way to do this is to use the CloseNotifier interface that will appear in future releases of
				// Go (this is written as of 1.0.3):
				// https://code.google.com/p/go/source/detail?name=3292433291b2

				flusher, ok := pResp.(http.Flusher)
				if !ok {
					errMsgStr := fmt.Sprintf("%s handler failed - SSE http streaming is not supported on the server", pPathStr)
					gfErr := gf_core.ErrorCreate(errMsgStr, 
						"http_server_flusher_not_supported_error",
						map[string]interface{}{"path_str": pPathStr,},
						nil, "gf_rpc_lib", pRuntimeSys)
					return nil, gfErr
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

				//-------------------------------------------------
				completeFun := func() {
					data_complete_lst, _ := json.Marshal(map[string]interface{}{"status_str": "complete"})
					fmt.Fprintf(pResp, "data: %s\n\n", data_complete_lst)
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
					case dataUpdate := <- dataUpdatesCh:

						sseEventUNIXtimeStr := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)
						sseEventIDstr       := sseEventUNIXtimeStr

						eventMap := map[string]interface{}{
							"data_map":   dataUpdate,
							"status_str": "ok",
						}
						sseEventLst, _ := json.Marshal(eventMap)

						// SSE_WRITE
						fmt.Fprintf(pResp, "id: %s\n", sseEventIDstr)
						fmt.Fprintf(pResp, "data: %s\n\n", sseEventLst)
						flusher.Flush()
					
					//-----------------------
					// ERROR_UPDATE
					case gfErr := <- dataErrCh:

						sseEventUNIXtimeStr := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)
						sseEventIDstr       := sseEventUNIXtimeStr
						
						eventMap := map[string]interface{}{
							"data_map":   gfErr,
							"status_str": "error",
						}
						sseEventLst, _ := json.Marshal(eventMap)

						// SSE_WRITE
						fmt.Fprintf(pResp, "id: %s\n", sseEventIDstr)
						fmt.Fprintf(pResp, "data: %s\n\n", sseEventLst)
						flusher.Flush()

						// an error occured, so complete SSE stream and exit
						completeFun()
						break

					//-----------------------
					// COMPLETE
					case completeBool := <- dataCompleteCh:

						if completeBool {

							completeFun()
							break
						}

					//-----------------------

					}	
				}
			}

			return nil, nil
		},
		nil,
		pStoreRunBool,
		pRuntimeSys)
}

//-------------------------------------------------

func clientParseResponseSSE(pBodyStr string,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	data_items_lst := []map[string]interface{}{}

	for _, line_str := range strings.Split(pBodyStr, `\n`) {

		pRuntimeSys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>")
		pRuntimeSys.LogFun("INFO", line_str)

		// filter out keep-alive new lines
		if line_str != "" && strings.HasPrefix(line_str, "data: ") {


			msg_str := strings.Replace(line_str, "data: ", "", 1)
			msg_map := map[string]interface{}{}
			err     := json.Unmarshal([]byte(msg_str), &msg_map)

			if err != nil {

				gf_err := gf_core.ErrorCreate("failed to parse JSON response line of the SSE stream (of even updates from a gf_images server)",
					"json_decode_error",
					map[string]interface{}{"line_str": line_str,},
					err, "gf_images_lib", pRuntimeSys)

				return nil, gf_err
			}

			//-------------------
			// STATUS
			if _,ok := msg_map["status_str"]; !ok {
				err_usr_msg := "sse message json doesnt container key status_str"
				gf_err      := gf_core.ErrorCreate(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", pRuntimeSys)
				return nil, gf_err
			}
			status_str := msg_map["status_str"].(string)

			if !(status_str == "ok" || status_str == "error") {

				err_usr_msg := "sse message json status_str key is not of value ok|error"
				gf_err      := gf_core.ErrorCreate(err_usr_msg,
					"verify__invalid_key_value_error",
					map[string]interface{}{
						"status_str": status_str,
						"msg_map":    msg_map,
					},
					nil, "gf_images_lib", pRuntimeSys)
				return nil, gf_err
			}

			//-------------------
			// DATA
			if _,ok := msg_map["data_map"]; !ok {
				err_usr_msg := "sse message json doesnt container key data_map"
				gf_err      := gf_core.ErrorCreate(err_usr_msg,
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_images_lib", pRuntimeSys)
				return nil, gf_err
			}
			
			data_map := msg_map["data_map"].(map[string]interface{})
			
			//-------------------
			data_items_lst = append(data_items_lst, data_map)
		}
	}
	return data_items_lst, nil
}