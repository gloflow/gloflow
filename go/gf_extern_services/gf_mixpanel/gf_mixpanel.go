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
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------
type GFmixpanelInfo struct {
	Username_str   string
	Secret_str     string
	Project_id_str string
}

//-------------------------------------------------------------
func Event_send(pEventTypeStr string,
	pEventMetaMap map[string]interface{},
	pInfo         *GFmixpanelInfo,
	pRuntimeSys   *gf_core.RuntimeSys) *gf_core.GFerror {

	request := gorequest.New()

	// AUTH - mixpanel uses basic http auth
	request.Header.Add("user", fmt.Sprintf("%s:%s", pInfo.Username_str, pInfo.Secret_str))

	dataMap := map[string]interface{}{
		"event":      pEventTypeStr,
		"properties": pEventMetaMap,
	}
	dataBytesLst, _ := json.Marshal(dataMap)


	urlStr := fmt.Sprintf("https://api.mixpanel.com/import?strict=1&project_id=%s", pInfo.Project_id_str)
	_, _, errs := request.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	if errs != nil {
		err   := errs[0] // FIX!! - use all errors in some way, just in case
		gfErr := gf_core.Error__create("failed to send event to mixpanel",
			"http_client_req_error",
			map[string]interface{}{"url_str": urlStr,},
			err, "gf_mixpanel", pRuntimeSys)
		return gfErr
	}

	return nil
}