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

package gf_identity_core

import (
	"fmt"
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
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
	User_name_str    GFuserName `validate:"required,min=3,max=50"`
	Confirm_code_str string     `validate:"required,min=10,max=20"`
}

//---------------------------------------------------
// GET_USER_NAME_FROM_CTX
func GetUserIDfromCtx(pCtx context.Context) (gf_core.GF_ID, bool) {
	
	userID := pCtx.Value("gf_user_id")

	if userID != nil {
		userIDstr := gf_core.GF_ID(userID.(string))
		return userIDstr, true
	} else {
		return "", false
	}
	
	return "", false
}

//---------------------------------------------------
func GetSessionTTL() (int, int64) {
	sessionTTLhoursInt   := 24 * 30 // 1 month
	sessionTTLsecondsInt := int64(60*60*24*7)
	return sessionTTLhoursInt, sessionTTLsecondsInt
}

//---------------------------------------------------
// HTTP
//---------------------------------------------------
func HTTPgetUserStdInput(pCtx context.Context,
	p_req         *http.Request,
	p_resp        http.ResponseWriter,
	pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, gf_core.GF_ID, GF_user_address_eth, *gf_core.GF_error) {

	inputMap, gfErr := gf_core.HTTPgetInput(p_req, pRuntimeSys)
	if gfErr != nil {
		return nil, "", GF_user_address_eth(""), gfErr
	}
	
	// user-name is supplied if the traditional auth system is used, and not web3/eth
	var userIDstr gf_core.GF_ID
	if inputUserIDstr, ok := inputMap["user_id_str"].(gf_core.GF_ID); ok {
		userIDstr = inputUserIDstr
	} else {

		// logged in users are added to context by gf_rpc, not supplied explicitly
		// via http request input (as they are for unauthenticated requests).
		userIDfromCtxStr, ok := GetUserIDfromCtx(pCtx) // p_ctx.Value("gf_user_name").(string)
		if ok {
			userIDstr = userIDfromCtxStr
		}
	}

	fmt.Println("user ID:", userIDstr)

	// users eth address is used if the user picks that method instead of traditional
	var userAddressETHstr string;
	if input_user_address_eth_str, ok := inputMap["user_address_eth_str"].(string); ok {
		userAddressETHstr = input_user_address_eth_str
	}

	// one of the these values has to be supplied, they cant both be missing
	if userIDstr == "" && userAddressETHstr == "" {
		gfErr := gf_core.Mongo__handle_error("user_name_str or user_address_eth_str arguments are missing from request",
			"verify__input_data_missing_in_req_error",
			map[string]interface{}{},
			nil, "gf_identity_lib", pRuntimeSys)
		return nil, "", GF_user_address_eth(""), gfErr
	}

	return inputMap, userIDstr, GF_user_address_eth(userAddressETHstr), nil
}

//---------------------------------------------------
func Http__get_user_address_eth_input(p_req *http.Request,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) (GF_user_address_eth, *gf_core.GF_error) {

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
func Http__get_user_update_input(p_req *http.Request,
	p_runtime_sys *gf_core.RuntimeSys) (*GF_user__http_input_update, *gf_core.GF_error) {

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
func Http__get_email_confirm_input(p_req *http.Request,
	p_runtime_sys *gf_core.RuntimeSys) (*GF_user__http_input_email_confirm, *gf_core.GF_error) {

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
		User_name_str:    GFuserName(user_name_str),
		Confirm_code_str: confirmation_code_str,
	}

	return input, nil
}