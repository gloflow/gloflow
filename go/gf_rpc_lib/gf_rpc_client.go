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
	"encoding/json"
	"io/ioutil"
	"context"
	"strings"
	"bufio"
	log "github.com/sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Client__request_sse(p_url_str string,
	p_resp_data_ch chan(map[string]interface{}),
	p_headers_map  map[string]string,
	p_ctx          context.Context,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.GF_error {

	//-----------------------
	// FETCH_URL

	user_agent_str := "gf_rpc_client"
	gf_http_fetch, gf_err := gf_core.HTTP__fetch_url(p_url_str,
		p_headers_map,
		user_agent_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//-----------------------

	http_resp   := gf_http_fetch.Resp
	resp_reader := bufio.NewReader(http_resp.Body)

	for {

		line_lst, err := resp_reader.ReadBytes('\n')
		if err != nil {
			
			gf_err := gf_core.Error__create("failed to read a line in body reasponse stream for client http sse connection",
				"io_reader_error",
				map[string]interface{}{"url_str": p_url_str,},
				nil, "gf_rpc_lib", p_runtime_sys)
			return gf_err
		}

		line_str := string(line_lst)

		// filter out keep-alive new lines
		if line_str != "" && strings.HasPrefix(line_str, "data: ") {


			msg_str := strings.Replace(line_str, "data: ", "", 1)

			log.WithFields(log.Fields{"event": msg_str,}).Info("GF_RPC_CLIENT - SSE event")

			msg_map := map[string]interface{}{}
			err     := json.Unmarshal([]byte(msg_str), &msg_map)
			if err != nil {
				gf_err := gf_core.Error__create("failed to parse JSON response line of the SSE stream (of even updates from a gf_images server)",
					"json_unmarshal_error",
					map[string]interface{}{
						"url_str":  p_url_str,
						"line_str": line_str,
					},
					err, "gf_rpc_lib", p_runtime_sys)
				return gf_err
			}

			//-------------------
			// STATUS
			if _, ok := msg_map["status_str"]; !ok {
				gf_err := gf_core.Error__create("sse message json doesnt container key status_str",
					"verify__missing_key_error",
					map[string]interface{}{
						"url_str":  p_url_str,
						"line_str": line_str,
					},
					nil, "gf_rpc_lib", p_runtime_sys)
				return gf_err
			}
			status_str := msg_map["status_str"].(string)

			if !(status_str == "ok" || status_str == "error" || status_str == "complete") {
				gf_err := gf_core.Error__create("sse message json status_str key is not of value ok|complete|error",
					"verify__invalid_key_value_error",
					map[string]interface{}{
						"status_str": status_str,
						"url_str":    p_url_str,
						"line_str":   line_str,
					},
					nil, "gf_rpc_lib", p_runtime_sys)
				return gf_err
			}

			//-------------------
			// DATA
			if _, ok := msg_map["data_map"]; !ok {
				gf_err := gf_core.Error__create("sse message json doesnt container key data_map",
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msg_map,},
					nil, "gf_rpc_lib", p_runtime_sys)
				return gf_err
			}
			
			data_map := msg_map["data_map"].(map[string]interface{})
			
			//-------------------
			


			fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
			spew.Dump(data_map)

			p_resp_data_ch <- data_map
		}
	}
	
	return nil
}

//-------------------------------------------------
func Client__request(p_url_str string,
	p_headers_map map[string]string,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{}, *gf_core.GF_error) {

	yellow   := color.New(color.FgYellow).SprintFunc()
	yellowBg := color.New(color.FgBlack, color.BgYellow).SprintFunc()

	fmt.Printf("%s - REQUEST SENT - %s\n", yellow("gf_rpc_client"), yellowBg(p_url_str))
	

	//-----------------------
	// FETCH_URL
	user_agent_str := "gf_rpc_client"
	gf_http_fetch, gf_err := gf_core.HTTP__fetch_url(p_url_str,
		p_headers_map,
		user_agent_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//-----------------------
	// JSON_DECODE
	body_bytes_lst, _ := ioutil.ReadAll(gf_http_fetch.Resp.Body)

	var resp_map map[string]interface{}
	err := json.Unmarshal(body_bytes_lst, &resp_map)
	if err != nil {
		gf_err := gf_core.Error__create(fmt.Sprintf("failed to parse json response from gf_rpc_client"), 
			"json_decode_error",
			map[string]interface{}{"url_str": p_url_str,},
			err, "gf_rpc_lib", p_runtime_sys)
		return nil, gf_err
	}

	//-----------------------

	r_status_str := resp_map["status"].(string)

	if r_status_str == "OK" {
		data_map := resp_map["data"].(map[string]interface{})

		return data_map, nil
	} else {

		gf_err := gf_core.Error__create(fmt.Sprintf("received a non-OK response from GF HTTP REST API"), 
			"http_client_gf_status_error",
			map[string]interface{}{"url_str": p_url_str,},
			nil, "gf_rpc_lib", p_runtime_sys)
		return nil, gf_err
	}

	return nil, nil
}