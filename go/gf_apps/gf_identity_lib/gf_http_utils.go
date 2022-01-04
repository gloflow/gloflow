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

package gf_identity_lib

import (
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
type GF_user__http_input_update struct {
	Username_str         *string  `json:"username_str"    validate:"min=3,max=50"`
	Screenname_str       *string  `json:"screenname_str"  validate:"min=3,max=50"`
	Email_str            *string  `json:"email_str"       validate:"min=6,max=50"`
	Description_str      *string  `json:"description_str" validate:"min=1,max=2000"`

	Profile_image_url_str *string `json:"profile_image_url_str" validate:"min=1,max=100"` // FIX!! - validation
	Banner_image_url_str  *string `json:"banner_image_url_str"  validate:"min=1,max=100"` // FIX!! - validation
}

//---------------------------------------------------
func http__get_user_update(p_req *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__http_input_update, *gf_core.GF_error) {

	handler_url_path_str := p_req.URL.Path
	input             := GF_user__http_input_update{}
	body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
	err               := json.Unmarshal(body_bytes_lst, &input)
		
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse json http input for user update",
			"json_decode_error",
			map[string]interface{}{"handler_url_path_str": handler_url_path_str,},
			err, "gf_identity_lib", p_runtime_sys)
		return nil, gf_err
	}

	return &input, nil
}

//---------------------------------------------------
func http__get_user_address_eth(p_req *http.Request,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (GF_user_address_eth, *gf_core.GF_error) {

	query_args_map := p_req.URL.Query()
	if values_lst, ok := query_args_map["addr_eth"]; ok {
		return GF_user_address_eth(values_lst[0]), nil
	} else {
		gf_err := gf_core.Error__create("incoming http request is missing the addr_eth query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_lib", p_runtime_sys)
		return GF_user_address_eth(""), gf_err
	}
	return GF_user_address_eth(""), nil
}