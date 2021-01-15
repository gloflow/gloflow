/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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


package gf_rpc_lib

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Client__request(p_url_str string,
	p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {

	yellow   := color.New(color.FgYellow).SprintFunc()
	yellowBg := color.New(color.FgBlack, color.BgYellow).SprintFunc()

	fmt.Printf("%s - REQUEST SENT - %s\n", yellow("gf_rpc_client"), yellowBg(p_url_str))
	
	//-----------------------
	// FETCH_URL
	user_agent_str := "gf_rpc_client"
	gf_http_fetch, gf_err := gf_core.HTTP__fetch_url(p_url_str,
		user_agent_str,
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


	r_status_str := resp_map["status_str"].(string)

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