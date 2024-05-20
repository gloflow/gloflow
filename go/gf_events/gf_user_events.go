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

package gf_events

import (
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"context"
	"github.com/ianoshen/uaparser"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type GFuserEventInput struct {
	TypeStr       string                 `json:"type_str"`
	SourceTypeStr string                 `json:"source_type_str"` // "browser"
	DataMap       map[string]interface{} `json:"data_map"`
	ReqCtx        *GFuserEventReqCtx     `json:"-"`
}

type GFuserEventReqCtx struct {
	UserIPstr      string `json:"user_ip_str"`
	UserAgentStr   string `json:"user_agent_str"`
	BrowserNameStr string `json:"browser_name_str"`
	BrowserVerStr  string `json:"browser_ver_str"`
	OSnameStr      string `json:"os_name_str"`
	OSverStr       string `json:"os_ver_str"`
}

//-------------------------------------------------

func UserEventParseInput(pReq *http.Request,
	pResp        http.ResponseWriter,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserEventInput, *gf_core.GFerror) {

	//--------------------
	input           := GFuserEventInput{}
	bodyBytesLst, _ := ioutil.ReadAll(pReq.Body)
	err             := json.Unmarshal(bodyBytesLst, &input)

	//-----------------
	// BROWSER INFORMATION

	ipStr      := pReq.RemoteAddr
	cleanIPstr := strings.Split(ipStr,":")[0]

	userAgentStr := pReq.UserAgent()
	userAgent    := uaparser.Parse(userAgentStr)

	var browserNameStr string
	var browserVerStr  string
	if userAgent.Browser != nil {
		browserNameStr = userAgent.Browser.Name
		browserVerStr  = userAgent.Browser.Version
	}

	osNameStr    := userAgent.OS.Name
	osVersionStr := userAgent.OS.Version

	

	reqCtx := &GFuserEventReqCtx {
		UserIPstr:      cleanIPstr,
		UserAgentStr:   userAgentStr,
		BrowserNameStr: browserNameStr,
		BrowserVerStr:  browserVerStr,
		OSnameStr:      osNameStr,
		OSverStr:       osVersionStr,
	}

	input.ReqCtx = reqCtx

	//--------------------

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to parse json http input for user_event",
			"json_decode_error",
			nil, err, "gf_events", pRuntimeSys)
		return nil, gfErr
	}
	return &input, nil
}

//-------------------------------------------------

func UserEventFullType(pInput *GFuserEventInput) string {
	return fmt.Sprintf("%s:%s", pInput.SourceTypeStr, pInput.TypeStr)
}

//-------------------------------------------------

func UserEventCreate(pInput *GFuserEventInput,
	pUserID     gf_core.GF_ID,
	pCtx		context.Context,
	pRuntimeSys *gf_core.RuntimeSys) {

	//------------------
	// EVENT
	if pRuntimeSys.EnableEventsAppBool {
		
		eventFullTypeStr := UserEventFullType(pInput)
		eventMeta := pInput.DataMap

		EmitApp(eventFullTypeStr,
			eventMeta,
			pUserID,
			pCtx,
			pRuntimeSys)
	}

	//------------------
}

//-------------------------------------------------

/*
func session__get_id_cookie(p_req *http.Request,
	p_resp        http.ResponseWriter,
	pRuntimeSys *gf_core.RuntimeSys) string {

	cookie, _ := p_req.Cookie("gf")  
	if cookie == nil {
		session_id_str := session__create_id_cookie(p_req, p_resp, pRuntimeSys)
		return session_id_str
	} else {
		session_id_str := cookie.Value
		return session_id_str
	}
}

//-------------------------------------------------

func session__create_id_cookie(p_req *http.Request,
	p_resp        http.ResponseWriter,
	pRuntimeSys *gf_core.RuntimeSys) string {

	current_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
	ip_str               := p_req.RemoteAddr
	session_id_str       := fmt.Sprintf("%f_%s", current_time__unix_f, ip_str)

	pRuntimeSys.LogFun("INFO", "session_id_str - "+session_id_str)

	new_cookie := http.Cookie{Name:"gf", Value:session_id_str}
	http.SetCookie(p_resp, &new_cookie)

	return session_id_str
}
*/