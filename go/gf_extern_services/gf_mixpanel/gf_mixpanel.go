/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_mixpanel

import (
	"fmt"
	// "encoding/json"
	// "encoding/base64"
	// "strings"
	// "github.com/parnurzeal/gorequest"
	"context"
	"github.com/mixpanel/mixpanel-go"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------

type GFmixpanelInfo struct {
	UsernameStr  string
	SecretStr    string
	ProjectIDstr string
	ProjectTokenStr string
}

//-------------------------------------------------------------

func EventSend(pEventTypeStr string,
	pEventMetaMap map[string]interface{},
	
	pInfo         *GFmixpanelInfo,
	pUserID       gf_core.GF_ID,
	pCtx          context.Context,
	pRuntimeSys   *gf_core.RuntimeSys) *gf_core.GFerror {

	mp := mixpanel.NewApiClient(pInfo.ProjectTokenStr)

	err := mp.Track(pCtx, []*mixpanel.Event{
		mp.NewEvent(pEventTypeStr, string(pUserID), pEventMetaMap),
	})
	
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to send event to mixpanel",
			"library_error",
			map[string]interface{}{"event_meta_map": pEventMetaMap,},
			err, "gf_mixpanel", pRuntimeSys)
		return gfErr
	}

	return nil
}

//-------------------------------------------------------------

/*
func EventSendHTTP(pEventTypeStr string,
	pEventMetaMap map[string]interface{},
	pInfo         *GFmixpanelInfo,
	pRuntimeSys   *gf_core.RuntimeSys) *gf_core.GFerror {

	request := gorequest.New()

	// AUTH - mixpanel uses basic http auth
	// request.Header.Add("user", fmt.Sprintf("%s:%s", pInfo.Username_str, pInfo.Secret_str))
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", pInfo.Username_str, pInfo.Secret_str)))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))

	dataMap := map[string]interface{}{
		"event":      pEventTypeStr,
		"properties": pEventMetaMap,
	}
	dataBytesLst, _ := json.Marshal(dataMap)


	urlStr := fmt.Sprintf("https://api.mixpanel.com/import?strict=1&project_id=%s", pInfo.Project_id_str)
	resp, body, errs := request.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	if errs != nil {
		err   := errs[0] // FIX!! - use all errors in some way, just in case
		gfErr := gf_core.ErrorCreate("failed to send event to mixpanel",
			"http_client_req_error",
			map[string]interface{}{"url_str": urlStr,},
			err, "gf_mixpanel", pRuntimeSys)
		return gfErr
	}



	if len(errs) > 0 {
		
		// log all errors
		errorMessages := make([]string, len(errs))
		for i, err := range errs {
			errorMessages[i] = err.Error()
		}
		gfErr := gf_core.ErrorCreate("failed to send event to mixpanel",
			"http_client_req_error",
			map[string]interface{}{"url_str": urlStr, "response": body},
			fmt.Errorf(strings.Join(errorMessages, ", ")), "gf_mixpanel", pRuntimeSys)
		return gfErr
	}



	// check the HTTP response status code for success
	if resp.StatusCode != 200 {
		gfErr := gf_core.ErrorCreate("unexpected status code from mixpanel",
			"http_client_req_error",
			map[string]interface{}{"url_str": urlStr, "response": body},
			fmt.Errorf("status code: %d", resp.StatusCode), "gf_mixpanel", pRuntimeSys)
		return gfErr
	}


	return nil
}
*/