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
	AppStr        string                 `json:"app_str"`
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
		
		eventAppStr      := pInput.AppStr
		eventFullTypeStr := UserEventFullType(pInput)
		eventMeta := pInput.DataMap

		EmitApp(eventFullTypeStr,
			eventMeta,
			eventAppStr,
			pUserID,
			pCtx,
			pRuntimeSys)
	}

	//------------------
}

//-------------------------------------------------