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
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// CLIENT_REQUEST_SSE

func ClientRequestSSE(pURLstr string,
	pRespDataCh chan(map[string]interface{}),
	pHeadersMap map[string]string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//-----------------------
	// FETCH_URL

	userAgentStr := "gf_rpc_client"
	gf_http_fetch, gfErr := gf_core.HTTPfetchURL(pURLstr,
		pHeadersMap,
		userAgentStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//-----------------------

	httpResp   := gf_http_fetch.Resp
	respReader := bufio.NewReader(httpResp.Body)

	for {

		lineLst, err := respReader.ReadBytes('\n')
		if err != nil {
			
			gfErr := gf_core.ErrorCreate("failed to read a line in body reasponse stream for client http sse connection",
				"io_reader_error",
				map[string]interface{}{"url_str": pURLstr,},
				nil, "gf_rpc_lib", pRuntimeSys)
			return gfErr
		}

		lineStr := string(lineLst)

		// filter out keep-alive new lines
		if lineStr != "" && strings.HasPrefix(lineStr, "data: ") {


			msg_str := strings.Replace(lineStr, "data: ", "", 1)

			log.WithFields(log.Fields{"event": msg_str,}).Info("GF_RPC_CLIENT - SSE event")

			msgMap := map[string]interface{}{}
			err     := json.Unmarshal([]byte(msg_str), &msgMap)
			if err != nil {
				gfErr := gf_core.ErrorCreate("failed to parse JSON response line of the SSE stream (of even updates from a gf_images server)",
					"json_decode_error",
					map[string]interface{}{
						"url_str":  pURLstr,
						"line_str": lineStr,
					},
					err, "gf_rpc_lib", pRuntimeSys)
				return gfErr
			}

			//-------------------
			// STATUS
			if _, ok := msgMap["status_str"]; !ok {
				gfErr := gf_core.ErrorCreate("sse message json doesnt container key status_str",
					"verify__missing_key_error",
					map[string]interface{}{
						"url_str":  pURLstr,
						"line_str": lineStr,
					},
					nil, "gf_rpc_lib", pRuntimeSys)
				return gfErr
			}
			statusStr := msgMap["status_str"].(string)

			if !(statusStr == "ok" || statusStr == "error" || statusStr == "complete") {
				gfErr := gf_core.ErrorCreate("sse message json status_str key is not of value ok|complete|error",
					"verify__invalid_key_value_error",
					map[string]interface{}{
						"status_str": statusStr,
						"url_str":    pURLstr,
						"line_str":   lineStr,
					},
					nil, "gf_rpc_lib", pRuntimeSys)
				return gfErr
			}

			//-------------------
			// DATA
			if _, ok := msgMap["data_map"]; !ok {
				gfErr := gf_core.ErrorCreate("sse message json doesnt container key data_map",
					"verify__missing_key_error",
					map[string]interface{}{"msg_map": msgMap,},
					nil, "gf_rpc_lib", pRuntimeSys)
				return gfErr
			}
			
			dataMap := msgMap["data_map"].(map[string]interface{})
			
			//-------------------

			pRespDataCh <- dataMap
		}
	}
	
	return nil
}

//-------------------------------------------------
// CLIENT_REQUEST

func ClientRequest(pURLstr string,
	pHeadersMap map[string]string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {

	yellow   := color.New(color.FgYellow).SprintFunc()
	yellowBg := color.New(color.FgBlack, color.BgYellow).SprintFunc()

	fmt.Printf("%s - REQUEST SENT - %s\n", yellow("gf_rpc_client"), yellowBg(pURLstr))
	

	//-----------------------
	// FETCH_URL
	userAgentStr := "gf_rpc_client"
	HTTPfetch, gfErr := gf_core.HTTPfetchURL(pURLstr,
		pHeadersMap,
		userAgentStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-----------------------
	// JSON_DECODE
	bodyBytesLst, _ := ioutil.ReadAll(HTTPfetch.Resp.Body)

	var respMap map[string]interface{}
	err := json.Unmarshal(bodyBytesLst, &respMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from gf_rpc_client"), 
			"json_decode_error",
			map[string]interface{}{"url_str": pURLstr,},
			err, "gf_rpc_lib", pRuntimeSys)
		return nil, gfErr
	}

	//-----------------------

	rStatusStr := respMap["status"].(string)

	if rStatusStr == "OK" {
		dataMap := respMap["data"].(map[string]interface{})

		return dataMap, nil
	} else {

		gfErr := gf_core.ErrorCreate(fmt.Sprintf("received a non-OK response from GF HTTP REST API"), 
			"http_client_gf_status_error",
			map[string]interface{}{"url_str": pURLstr,},
			nil, "gf_rpc_lib", pRuntimeSys)
		return nil, gfErr
	}

	return nil, nil
}