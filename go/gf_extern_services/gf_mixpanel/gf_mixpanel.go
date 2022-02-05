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
type GF_mixpanel_info struct {
	Username_str   string
	Secret_str     string
	Project_id_str string
}

//-------------------------------------------------------------
func Event_send(p_event_type_str string,
	p_event_meta_map map[string]interface{},
	p_info           *GF_mixpanel_info,
	p_runtime_sys    *gf_core.Runtime_sys) *gf_core.GF_error {

	request := gorequest.New()

	// AUTH - mixpanel uses basic http auth
	request.Header.Add("user", fmt.Sprintf("%s:%s", p_info.Username_str, p_info.Secret_str))


	data_map := map[string]interface{}{
		"event":      p_event_type_str,
		"properties": p_event_meta_map,
	}
	data_bytes_lst, _ := json.Marshal(data_map)


	url_str := fmt.Sprintf("https://api.mixpanel.com/import?strict=1&project_id=%s", p_info.Project_id_str)
	_, _, errs := request.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	if errs != nil {
		err    := errs[0] // FIX!! - use all errors in some way, just in case
		gf_err := gf_core.Error__create("failed to send event to mixpanel",
			"http_client_req_error",
			map[string]interface{}{"url_str": url_str,},
			err, "gf_mixpanel", p_runtime_sys)
		return gf_err
	}


	return nil
}