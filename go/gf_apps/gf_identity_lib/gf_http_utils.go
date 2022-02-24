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
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//------------------------------------------------
type GF_user__http_input_update struct {
	Screen_name_str       *string  `json:"screen_name_str" validate:"min=3,max=50"`
	Email_str             *string  `json:"email_str"       validate:"min=6,max=50"`
	Description_str       *string  `json:"description_str" validate:"min=1,max=2000"`

	Profile_image_url_str *string `json:"profile_image_url_str" validate:"min=1,max=100"` // FIX!! - validation
	Banner_image_url_str  *string `json:"banner_image_url_str"  validate:"min=1,max=100"` // FIX!! - validation
}

type GF_user__http_input_email_confirm struct {
	User_name_str    GF_user_name `validate:"required,min=3,max=50"`
	Confirm_code_str string       `validate:"required,min=10,max=20"`
}

//---------------------------------------------------
func Http__get_user_std_input(p_req *http.Request,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) (map[string]interface{}, string, GF_user_address_eth, *gf_core.GF_error) {

	input_map, gf_err := gf_rpc_lib.Get_http_input(p_resp, p_req, p_runtime_sys)
	if gf_err != nil {
		return nil, "", GF_user_address_eth(""), gf_err
	}

	// user-name is supplied if the traditional auth system is used, and not web3/eth
	var user_name_str string;
	if input_user_name_str, ok := input_map["user_name_str"].(string); ok {
		user_name_str = input_user_name_str
	}

	// users eth address is used if the user picks that method instead of traditional
	var user_address_eth_str string;
	if input_user_address_eth_str, ok := input_map["user_address_eth_str"].(string); ok {
		user_address_eth_str = input_user_address_eth_str
	}

	// one of the these values has to be supplied, they cant both be missing
	if user_name_str == "" && user_address_eth_str == "" {
		gf_err := gf_core.Mongo__handle_error("user_name_str or user_address_eth_str arguments are missing from request",
			"verify__input_data_missing_in_req_error",
			map[string]interface{}{},
			nil, "gf_identity_lib", p_runtime_sys)
		return nil, "", GF_user_address_eth(""), gf_err
	}

	return input_map, user_name_str, GF_user_address_eth(user_address_eth_str), nil
}

//---------------------------------------------------
func http__get_user_address_eth_input(p_req *http.Request,
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

//---------------------------------------------------
func http__get_user_update_input(p_req *http.Request,
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
func http__get_email_confirm_input(p_req *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user__http_input_email_confirm, *gf_core.GF_error) {

	var user_name_str         string
	var confirmation_code_str string

	query_args_map := p_req.URL.Query()
	
	if values_lst, ok := query_args_map["u"]; ok {
		user_name_str = values_lst[0]
	} else {
		gf_err := gf_core.Error__create("incoming http request is missing the email user_name query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_lib", p_runtime_sys)
		return nil, gf_err
	}

	if values_lst, ok := query_args_map["c"]; ok {
		confirmation_code_str = values_lst[0]
	} else {
		gf_err := gf_core.Error__create("incoming http request is missing the email confirmation_code query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_lib", p_runtime_sys)
		return nil, gf_err
	}

	input := &GF_user__http_input_email_confirm{
		User_name_str:    GF_user_name(user_name_str),
		Confirm_code_str: confirmation_code_str,
	}

	return input, nil
}