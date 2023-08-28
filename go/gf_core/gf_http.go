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
	"io"
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
// COOKIES
//---------------------------------------------------

func HTTPgetCookieFromReq(pCookieNameStr string,
	pReq        *http.Request,
	pRuntimeSys *RuntimeSys) (bool, string) {

	pRuntimeSys.LogNewFun("DEBUG", `cookies in request...`,
		HTTPgetAllCookies(pReq))

	for _, cookie := range pReq.Cookies() {
		if (cookie.Name == pCookieNameStr) {
			dataStr := cookie.Value
			return true, dataStr
		}
	}
	return false, ""
}

//---------------------------------------------------

func HTTPgetAllCookies(pReq *http.Request) map[string]interface{} {

	cookies := pReq.Cookies()
	cookiesMap := map[string]interface{}{}
	for _, cookie := range cookies {
		cookiesMap[cookie.Name] = cookie.Value
	}
	return cookiesMap
}

//---------------------------------------------------

func HTTPsetCookieOnReq(pCookieNameStr string,
	pDataStr     string,
	pResp        http.ResponseWriter,
	pTTLhoursInt int) {

	ttl    := time.Duration(pTTLhoursInt) * time.Hour
	expire := time.Now().Add(ttl)
	
	cookie := http.Cookie{
		Name:  pCookieNameStr,
		Value: pDataStr,

		// if not set the cookie would be a session cookie and would be
		// deleted on restart of the browser.
		Expires: expire,

		// IMPORTANT!! - session cookie should be set for all paths
		//               on the same domain, not just the /v1/identity/...
		//               paths, because session is verified on all of them.
		//               otherwise the cookie will only be set requests
		//               that are on some subset of urls relative to the root.
		// IMPORTANT!! - In Go (and in HTTP cookies in general), if the Path for a cookie is
		//               not explicitly set the cookie's path will default to the path of 
		//               the URL where the Set-Cookie HTTP response header was received from.
		//               This means that the cookie will be sent only for requests to this path and its subpaths.
		Path: "/", 
		
		// ADD!! - ability to specify multiple domains that the session is
		//         set for in case the GF services and API endpoints are spread
		//         across multiple domains.
		// Domain: "", 

		// force cookie to only be accessible over HTTPS connections
		Secure: true,

		// IMPORTANT!! - make cookie http_only, disabling browser js context
		//               from being able to read its value
		HttpOnly: true,

		// SameSite allows a server to define a cookie attribute making it impossible for
		// the browser to send this cookie along with cross-site requests. The main
		// goal is to mitigate the risk of cross-origin information leakage, and provide
		// some protection against cross-site request forgery attacks.
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(pResp, &cookie)
}

//---------------------------------------------------

func HTTPdeleteCookieOnResp(pCookieNameStr string,
	pResp http.ResponseWriter) {

	expiredCookie := http.Cookie{
		Name:     pCookieNameStr,
		Value:    "", // empty value
		Expires:  time.Unix(0, 0), // set expiration in the past
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	
	http.SetCookie(pResp, &expiredCookie)
}

//---------------------------------------------------
// VAR
//---------------------------------------------------

func HTTPdetectMIMEtypeFromURL(pURLstr string,
	pHeadersMap   map[string]string,
	pUserAgentStr string,
	pCtx          context.Context,
	pRuntimeSys   *RuntimeSys) (string, *GFerror) {

	pRuntimeSys.LogNewFun("DEBUG", "detecting remote URL MIME type (fetching initial 512 bytes)...", map[string]interface{}{
		"url_str": pURLstr,
	})

	// fetch the first 512 bytes, which is all we need
	// for determening MIME type.
	//
	// "Range" - HTTP request header indicates the part of a document 
	//           that the server should return.
	pHeadersMap["Range"] = "bytes=0-511"

	// HTTP REQUEST
	HTTPfetch, gfErr := HTTPfetchURL(pURLstr, pHeadersMap, pUserAgentStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}
	defer HTTPfetch.Resp.Body.Close()
	

	bodyBytesLst, _ := ioutil.ReadAll(HTTPfetch.Resp.Body)

	// REMOVE!! - for debugging only
	// fmt.Println("body bytes - ", string(bodyBytesLst))

	contentTypeStr := http.DetectContentType(bodyBytesLst)

	return contentTypeStr, nil
}

//---------------------------------------------------

func HTTPgetInput(pReq *http.Request,
	pRuntimeSys *RuntimeSys) (map[string]interface{}, *GFerror) {

	bodyBytesLst, _ := ioutil.ReadAll(pReq.Body)

	// parse body bytes only if they're larger than 0
	if len(bodyBytesLst) > 0 {

		i, gfErr := ParseJSONfromByteList(bodyBytesLst, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		return i, nil
	}

	return nil, nil
}

//---------------------------------------------------

func HTTPfetchURL(pURLstr string,
	pHeadersMap   map[string]string,
	pUserAgentStr string,
	pCtx          context.Context,
	pRuntimeSys   *RuntimeSys) (*GF_http_fetch, *GFerror) {

	pRuntimeSys.LogNewFun("INFO", "fetching URL", map[string]interface{}{
		"url_str": pURLstr,
	})

	// TIMEOUT
	timeoutSec := time.Second * 60

	client := &http.Client{
		Timeout: timeoutSec, // time.Second * 10, // to prevent requests taking too long to return

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
			req.Header.Set("User-Agent", pUserAgentStr)
			return nil
		},
	}

	req, err := http.NewRequest("GET", pURLstr, nil)
	if err != nil {
		gfErr := ErrorCreate("image fetcher failed to create HTTP request to fetch a file",
			"http_client_req_error",
			map[string]interface{}{"url_str": pURLstr,},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}

	//-------------------------
	// HEADERS
	for k, v := range pHeadersMap {
		req.Header.Set(k, v)
	}

	//-------------------------
	// USER_AGENT
	req.Header.Del("User-Agent")
	req.Header.Set("User-Agent", pUserAgentStr)

	//-------------------------
	// req_with_ctx := req.WithContext(pCtx)

	
	reqUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

	// EXECUTE
	resp, err := client.Do(req)
	
	respUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

	if err != nil {
		gfErr := ErrorCreate("http fetch failed to execute HTTP request to fetch a url",
			"http_client_req_error",
			map[string]interface{}{"url_str": pURLstr,},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}

	statusCodeInt := resp.StatusCode
	headersMap    := resp.Header

	pRuntimeSys.LogNewFun("DEBUG", "http response received ", map[string]interface{}{
		"status_code_int": statusCodeInt,
	})

	respHeadersMap := map[string]string{}
	for k, v := range headersMap {
		respHeadersMap[k] = v[0]
	}

	HTTPfetch := &GF_http_fetch{
		Url_str:          pURLstr, 
		Status_code_int:  statusCodeInt,
		Resp_headers_map: respHeadersMap,
		Req_time_f:       reqUNIXtimeF,
		Resp_time_f:      respUNIXtimeF,
		Resp:             resp,
	}

	return HTTPfetch, nil
}

//---------------------------------------------------

func HTTPgetFile(pTargetURLstr string,
	pFileLocalPathStr string,
	pCtx              context.Context,
	pRuntimeSys       *RuntimeSys) *GFerror {
	

	//--------------
	headersMap, userAgentStr := HTTPgetReqConfig()

	HTTPfetch, gfErr := HTTPfetchURL(pTargetURLstr, headersMap, userAgentStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	defer HTTPfetch.Resp.Body.Close()

	//--------------
	// WRITE TO FILE
	fmt.Printf("local file path - %s\n", pFileLocalPathStr)

	out, err := os.Create(pFileLocalPathStr)
	defer out.Close()

	if err != nil {
		gfErr := ErrorCreate("failed to create local file for fetched file",
			"file_create_error",
			map[string]interface{}{"file_local_path_str": pFileLocalPathStr,},
			err, "gf_core", pRuntimeSys)
		return gfErr
	}

	_, err = io.Copy(out, HTTPfetch.Resp.Body)
	if err != nil {
		gfErr := ErrorCreate("failed to copy HTTP GET response Body buffer to a file",
			"file_buffer_copy_error",
			map[string]interface{}{
				"file_local_path_str": pFileLocalPathStr,
				"target_url_str":      pTargetURLstr,
			},
			err, "gf_core", pRuntimeSys)
		return gfErr
	}

	//--------------
	return nil
}

//---------------------------------------------------
// PUT_FILE

func HTTPputFile(pTargetURLstr string,
	pFilePathStr string,
	pHeadersMap  map[string]string,
	pRuntimeSys  *RuntimeSys) (*http.Response, *GFerror) {



	// FILE_OPEN
	f, err := os.Open(pFilePathStr)
	if err != nil {
		gfErr := ErrorCreate("failed to open a file on the local FS that is to be sent to AWS S3",
			"file_open_error",
			map[string]interface{}{
				"target_url_str": pTargetURLstr,
				"file_path_str":  pFilePathStr,
			},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}
	buffer := bufio.NewReader(f)



	req, err := http.NewRequest(http.MethodPut, pTargetURLstr, buffer)
    if err != nil {
        gfErr := ErrorCreate("failed to create a HTTP PUT request to upload file to S3",
			"http_client_req_error",
			map[string]interface{}{
				"target_url_str": pTargetURLstr,
				"file_path_str":  pFilePathStr,
			},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}

	// golang http client sets "Transfer-Encoding": "chunked", 
	// which is rejected by some servers (AWS, etc.). so here we turn that off.
	req.TransferEncoding = []string{"identity"}



	// FILE_SIZE
	fi, err := os.Stat(pFilePathStr)
    if err != nil {
		gfErr := ErrorCreate("failed to get file info via stat() to find out its size for uploading to S3 via HTTP PUT",
			"file_stat_error",
			map[string]interface{}{
				"target_url_str": pTargetURLstr,
				"file_path_str":  pFilePathStr,
			},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
    }
	req.ContentLength = fi.Size()


	// HEADERS
	for k, v := range pHeadersMap {
    	req.Header.Set(k, v)
	}

    client := http.Client{}

	pRuntimeSys.LogNewFun("INFO", "issuing HTTP PUT request", map[string]interface{}{
		"target_url_str": pTargetURLstr,
	})

    resp, err := client.Do(req)
    if err != nil {
		gfErr := ErrorCreate("failed to execute a HTTP PUT request to upload file to S3",
			"http_client_req_error",
			map[string]interface{}{
				"target_url_str": pTargetURLstr,
				"file_path_str":  pFilePathStr,
			},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
    }



	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))

	return resp, nil
}

//-------------------------------------------------

func HTTPinitStaticServingWithMux(pURLbaseStr string,
	p_local_dir_path_str string,
	p_mux                *http.ServeMux,
	pRuntimeSys        *RuntimeSys) {
	
	// IMPORTANT!! - trailing "/" in this url spec is important, since the desired urls that should
	//               match this are /*/static/some_further_text, and those will only match
	//               if the spec here ends with "/"
	urlStr := fmt.Sprintf("%s/static/", pURLbaseStr)
	p_mux.HandleFunc(urlStr, func(pResp http.ResponseWriter, pReq *http.Request) {

		HTTPserveFile(p_local_dir_path_str,
			urlStr,
			pReq, pResp, pRuntimeSys)
	})
}

//-------------------------------------------------

func HTTPinitStaticServing(pURLbaseStr string,
	pRuntimeSys *RuntimeSys) {
	
	local_dir_str := fmt.Sprintf("./static")

	// IMPORTANT!! - trailing "/" in this url spec is important, since the desired urls that should
	//               match this are /*/static/some_further_text, and those will only match
	//               if the spec here ends with "/"
	urlStr := fmt.Sprintf("%s/static/", pURLbaseStr)
	http.HandleFunc(urlStr, func(pResp http.ResponseWriter, pReq *http.Request) {

		HTTPserveFile(local_dir_str,
			urlStr,
			pReq, pResp, pRuntimeSys)

		/*if pReq.Method == "GET" {
			path_str := pReq.URL.Path

			//remove url_base
			file_path_str      := strings.Replace(path_str, url_str, "", 1) // "1" - just replace one occurance
			file_ext_str       := filepath.Ext(file_path_str)
			file_mime_type_str := mime.TypeByExtension(file_ext_str)
			local_path_str     := fmt.Sprintf("./static/%s", file_path_str)

			pResp.Header().Set("Content-Type", file_mime_type_str)

			pRuntimeSys.LogFun("INFO", "file_path_str  - "+file_path_str)
			pRuntimeSys.LogFun("INFO", "local_path_str - "+local_path_str)

		    http.ServeFile(pResp, pReq, local_path_str)
		}*/
	})
}

//-------------------------------------------------

func HTTPserveFile(pLocalDirStr string,
	pURLstr     string,
	pReq        *http.Request,
	pResp       http.ResponseWriter,
	pRuntimeSys *RuntimeSys) {

	if pReq.Method == "GET" {
		path_str := pReq.URL.Path

		// remove url_base
		filePathStr     := strings.Replace(path_str, pURLstr, "", 1) // "1" - just replace one occurance
		fileExtStr      := filepath.Ext(filePathStr)
		fileMimeTypeStr := mime.TypeByExtension(fileExtStr)
		localPathStr    := fmt.Sprintf("%s/%s", pLocalDirStr, filePathStr)

		pResp.Header().Set("Content-Type", fileMimeTypeStr)

		pRuntimeSys.LogNewFun("DEBUG", "serving static file", map[string]interface{}{
			"file_url_path_str":   filePathStr,
			"file_local_path_str": localPathStr,
		})

		http.ServeFile(pResp, pReq, localPathStr)
	}
}

//-------------------------------------------------

func HTTPserializeCookies(pCookiesLst []*http.Cookie,
	pRuntimeSys *RuntimeSys) string {

	buffer := bytes.NewBufferString("")
	for _, cookie := range pCookiesLst {
		cookieStr := cookie.Raw
		buffer.WriteString("; "+cookieStr)
	}
	cookiesStr := buffer.String()
	return cookiesStr
}

//-------------------------------------------------

func HTTPinitSSE(pResp http.ResponseWriter,
	pRuntimeSys *RuntimeSys) (http.Flusher, *GFerror) {

	flusher, ok := pResp.(http.Flusher)
	if !ok {
		errMsgStr := "GF - SSE streaming not supported by this server"
		http.Error(pResp, errMsgStr, http.StatusInternalServerError)

		gfErr := ErrorCreate(errMsgStr,
			"http_server_flusher_not_supported_error",
			nil, nil, "gf_core", pRuntimeSys)

		return nil, gfErr
	}

	// IMPORTANT!! - listening for the closing of the http connections
	notify := pResp.(http.CloseNotifier).CloseNotify()
	go func() {
		<- notify
		pRuntimeSys.LogNewFun("DEBUG", "HTTP SSE connection closed", nil)
	}()

	pResp.Header().Set("Content-Type",                "text/event-stream")
	pResp.Header().Set("Cache-Control",               "no-cache")
	pResp.Header().Set("Connection",                  "keep-alive")
	pResp.Header().Set("Access-Control-Allow-Origin", "*") // CORS

	flusher.Flush()

	return flusher, nil
}

//-------------------------------------------------

func HTTPgetStreamingResponse(pURLstr string,
	pRuntimeSys *RuntimeSys) (*[]map[string]interface{}, *GFerror) {


	req,err := http.NewRequest("GET", pURLstr, nil)
    req.Header.Set("accept", "text/event-stream")

	client   := &http.Client{}
    resp,err := client.Do(req)
    if err != nil {
    	gfErr := ErrorCreate("http get_streaming_response failed to execute HTTP request to fetch a url",
			"http_client_req_error",
			map[string]interface{}{"url_str": pURLstr,},
			err, "gf_core", pRuntimeSys)
    	return nil, gfErr
    }

	dataLst := []map[string]interface{}{}
	reader  := bufio.NewReader(resp.Body)
	for {
	    lineLst, err := reader.ReadBytes('\n')
	    if err != nil {
	    	gfErr := ErrorCreate("failed to read a line of SSE streaming response from a server url",
				"io_reader_error",
				map[string]interface{}{"url_str": pURLstr,},
				err, "gf_core", pRuntimeSys)
	    	return nil, gfErr
	    }
	    
	    line_str := string(lineLst)

	    if strings.HasPrefix(line_str,"data: ") {
	    	clean_line_str := strings.Replace(line_str, "data: ", "", 1)

	    	dataMap := map[string]interface{}{}
	    	err     := json.Unmarshal([]byte(clean_line_str), &dataMap)

	    	if err != nil {
	    		gfErr := ErrorCreate("http get_streaming_response failed to parse JSON response",
					"json_decode_error",
					map[string]interface{}{"url_str": pURLstr,},
					err, "gf_core", pRuntimeSys)
	    		return nil, gfErr
	    	}

	    	dataLst = append(dataLst, dataMap)
	    }
	}
	return &dataLst, nil
}

//-------------------------------------------------

func HTTPgetReqConfig() (map[string]string, string) {
	headersMap   := map[string]string{}
	userAgentStr := "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1"
	return headersMap, userAgentStr
}

//-------------------------------------------------

func HTTPdisableCachingOfResponse(pResp http.ResponseWriter) {
	pResp.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    pResp.Header().Set("Pragma", "no-cache")
    pResp.Header().Set("Expires", "0")
}