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

package gf_core

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"bufio"
	"strings"
	"context"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"net/http"
	"mime"
)

//---------------------------------------------------
type GF_http_fetch = Gf_http_fetch
type Gf_http_fetch struct {
	Url_str          string            `bson:"url_str"`
	Status_code_int  int               `bson:"status_code_int"`
	Resp_headers_map map[string]string `bson:"resp_headers_map"`
	Req_time_f       float64           `bson:"req_time_f"`
	Resp_time_f      float64           `bson:"resp_time_f"`
	Resp             *http.Response    `bson:"-"`
}

//---------------------------------------------------
func HTTP__fetch_url(p_url_str string,
	p_headers_map    map[string]string,
	p_user_agent_str string,
	p_ctx            context.Context,
	p_runtime_sys    *Runtime_sys) (*GF_http_fetch, *GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_http_utils.HTTP__fetch_url()")


	// TIMEOUT
	timeout_sec := time.Second * 60

	client := &http.Client{
		Timeout: timeout_sec, // time.Second * 10, // to prevent requests taking too long to return

		/* IMPORTANT!! - golang http lib does not copy user-set headers on redirects, so a manual
		setting of these headers had to be added, via the CheckRedirect function
		that gets called on every redirect, which gives us a chance to to re-set
		user-agent headers again to the correct value*/
		/*CheckRedirect specifies the policy for handling redirects.
		If CheckRedirect is not nil, the client calls it before
		following an HTTP redirect. The arguments req and via are
		the upcoming request and the requests made already, oldest
		first. If CheckRedirect returns an error, the Client's Get
		method returns both the previous Response (with its Body
		closed) and CheckRedirect's error (wrapped in a url.Error)
		instead of issuing the Request req.
		As a special case, if CheckRedirect returns ErrUseLastResponse,
		then the most recent response is returned with its body
		unclosed, along with a nil error.
		If CheckRedirect is nil, the Client uses its default policy,
		which is to stop after 10 consecutive requests. */
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Del("User-Agent")
			req.Header.Set("User-Agent", p_user_agent_str)
			return nil
		},
	}

	req, err := http.NewRequest("GET", p_url_str, nil)
	if err != nil {
		gf_err := Error__create("image fetcher failed to create HTTP request to fetch a file",
			"http_client_req_error",
			map[string]interface{}{"url_str": p_url_str,},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}

	//-------------------------
	// HEADERS
	for k, v := range p_headers_map {
		req.Header.Set(k, v)
	}

	//-------------------------
	// USER_AGENT
	req.Header.Del("User-Agent")
	req.Header.Set("User-Agent", p_user_agent_str)

	//-------------------------
	// req_with_ctx := req.WithContext(p_ctx)

	// EXECUTE
	req_unix_time_f  := float64(time.Now().UnixNano())/1000000000.0
	resp, err        := client.Do(req)
	resp_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	if err != nil {
		gf_err := Error__create("http fetch failed to execute HTTP request to fetch a url",
			"http_client_req_error",
			map[string]interface{}{"url_str": p_url_str,},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}

	status_code_int := resp.StatusCode
	headers_map     := resp.Header



	fmt.Println(fmt.Sprintf("http response status_code - %d", status_code_int))


	resp_headers_map := map[string]string{}
	for k, v := range headers_map {
		resp_headers_map[k] = v[0]
	}

	gf_http_fetch := &GF_http_fetch{
		Url_str:          p_url_str, 
		Status_code_int:  status_code_int,
		Resp_headers_map: resp_headers_map,
		Req_time_f:       req_unix_time_f,
		Resp_time_f:      resp_unix_time_f,
		Resp:             resp,
	}

	return gf_http_fetch, nil
}

//---------------------------------------------------
// PUT_FILE
func HTTP__put_file(p_target_url_str string,
	p_file_path_str string,
	p_headers_map   map[string]string,
	p_runtime_sys   *Runtime_sys) (*http.Response, *Gf_error) {



	// FILE_OPEN
	f, err := os.Open(p_file_path_str)
	if err != nil {
		gf_err := Error__create("failed to open a file on the local FS that is to be sent to AWS S3",
			"file_open_error",
			map[string]interface{}{
				"target_url_str": p_target_url_str,
				"file_path_str":  p_file_path_str,
			},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}
	buffer := bufio.NewReader(f)



	req, err := http.NewRequest(http.MethodPut, p_target_url_str, buffer)
    if err != nil {
        gf_err := Error__create("failed to create a HTTP PUT request to upload file to S3",
			"http_client_req_error",
			map[string]interface{}{
				"target_url_str": p_target_url_str,
				"file_path_str":  p_file_path_str,
			},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}

	// golang http client sets "Transfer-Encoding": "chunked", 
	// which is rejected by some servers (AWS, etc.). so here we turn that off.
	req.TransferEncoding = []string{"identity"}



	// FILE_SIZE
	fi, err := os.Stat(p_file_path_str)
    if err != nil {
		gf_err := Error__create("failed to get file info via stat() to find out its size for uploading to S3 via HTTP PUT",
			"file_stat_error",
			map[string]interface{}{
				"target_url_str": p_target_url_str,
				"file_path_str":  p_file_path_str,
			},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
    }
	req.ContentLength = fi.Size()


	// HEADERS
	for k, v := range p_headers_map {
    	req.Header.Set(k, v)
	}

    client := http.Client{}

	p_runtime_sys.Log_fun("FUN_ENTER", fmt.Sprintf("ISSUING HTTP PUT REQUEST - %s", p_target_url_str))
    resp, err := client.Do(req)
    if err != nil {
		gf_err := Error__create("failed to execute a HTTP PUT request to upload file to S3",
			"http_client_req_error",
			map[string]interface{}{
				"target_url_str": p_target_url_str,
				"file_path_str":  p_file_path_str,
			},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
    }



	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))

	return resp, nil
}

//-------------------------------------------------
func HTTP__init_static_serving(p_url_base_str string,
	p_runtime_sys *Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_http_utils.HTTP__init_static_serving()")

	// IMPORTANT!! - trailing "/" in this url spec is important, since the desired urls that should
	//               match this are /*/static/some_further_text, and those will only match
	//               if the spec here ends with "/"
	url_str := p_url_base_str+"/static/"
	http.HandleFunc(url_str, func(p_resp http.ResponseWriter, p_req *http.Request) {
		// fmt.Println("FILE SERVE >>>>>>>>")

		if p_req.Method == "GET" {
			path_str := p_req.URL.Path

			//remove url_base
			file_path_str      := strings.Replace(path_str, url_str, "", 1) // "1" - just replace one occurance
			file_ext_str       := filepath.Ext(file_path_str)
			file_mime_type_str := mime.TypeByExtension(file_ext_str)
			local_path_str     := "./static/"+file_path_str

			p_resp.Header().Set("Content-Type", file_mime_type_str)

			p_runtime_sys.Log_fun("INFO","file_path_str  - "+file_path_str)
			p_runtime_sys.Log_fun("INFO","local_path_str - "+local_path_str)

		    http.ServeFile(p_resp, p_req, local_path_str)
		}
	})
}

//-------------------------------------------------
func HTTP__serialize_cookies(p_cookies_lst []*http.Cookie,
	p_runtime_sys *Runtime_sys) string {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_http_utils.HTTP__serialize_cookies()")

	buffer := bytes.NewBufferString("")
	for _, cookie := range p_cookies_lst {
		cookie_str := cookie.Raw
		buffer.WriteString("; "+cookie_str)
	}
	cookies_str := buffer.String()
	return cookies_str
}

//-------------------------------------------------
func HTTP__init_sse(p_resp http.ResponseWriter,
	p_runtime_sys *Runtime_sys) (http.Flusher, *Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_http_utils.HTTP__init_sse()")

	flusher,ok := p_resp.(http.Flusher)
	if !ok {
		err_msg_str := "GF - SSE streaming not supported by this server"
		http.Error(p_resp,err_msg_str,http.StatusInternalServerError)

		gf_err := Error__create(err_msg_str,
			"http_server_flusher_not_supported_error",
			nil, nil, "gf_core", p_runtime_sys)

		return nil, gf_err
	}

	// IMPORTANT!! - listening for the closing of the http connections
	notify := p_resp.(http.CloseNotifier).CloseNotify()
	go func() {
		<- notify
		p_runtime_sys.Log_fun("INFO", "HTTP SSE CONNECTION CLOSED")
	}()

	p_resp.Header().Set("Content-Type",                "text/event-stream")
	p_resp.Header().Set("Cache-Control",               "no-cache")
	p_resp.Header().Set("Connection",                  "keep-alive")
	p_resp.Header().Set("Access-Control-Allow-Origin", "*") // CORS

	flusher.Flush()

	return flusher, nil
}

//-------------------------------------------------
func HTTP__get_streaming_response(p_url_str string,
	p_runtime_sys *Runtime_sys) (*[]map[string]interface{}, *Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_http_utils.HTTP__get_streaming_response()")


	req,err := http.NewRequest("GET", p_url_str, nil)
    req.Header.Set("accept", "text/event-stream")

	client   := &http.Client{}
    resp,err := client.Do(req)
    if err != nil {
    	gf_err := Error__create("http get_streaming_response failed to execute HTTP request to fetch a url",
			"http_client_req_error",
			map[string]interface{}{"url_str": p_url_str,},
			err, "gf_core", p_runtime_sys)
    	return nil,gf_err
    }

	// resp, err := http.Get(p_url_str)
	// if err != nil {
	//	return nil,err
	// }

	data_lst := []map[string]interface{}{}
	reader   := bufio.NewReader(resp.Body)
	for {
	    line_lst,err := reader.ReadBytes('\n')
	    if err != nil {
	    	gf_err := Error__create("failed to read a line of SSE streaming response from a server url",
				"io_reader_error",
				map[string]interface{}{"url_str": p_url_str,},
				err, "gf_core", p_runtime_sys)
	    	return nil,gf_err
	    }
	    
	    line_str := string(line_lst)

	    if strings.HasPrefix(line_str,"data: ") {
	    	clean_line_str := strings.Replace(line_str,"data: ","",1)

	    	data_map := map[string]interface{}{}
	    	err      := json.Unmarshal([]byte(clean_line_str), &data_map)

	    	if err != nil {
	    		gf_err := Error__create("http get_streaming_response failed to parse JSON response",
					"json_decode_error",
					map[string]interface{}{"url_str": p_url_str,},
					err, "gf_core", p_runtime_sys)
	    		return nil, gf_err
	    	}

	    	data_lst = append(data_lst, data_map)
	    }
	}
	return &data_lst, nil
}