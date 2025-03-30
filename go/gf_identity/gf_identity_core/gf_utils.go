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
	// "fmt"
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------

type GFuserHTTPinputUpdate struct {
	ScreenNameStr       *string  `json:"screen_name_str" validate:"min=3,max=50"`
	EmailStr            *string  `json:"email_str"       validate:"min=6,max=50"`
	DescriptionStr      *string  `json:"description_str" validate:"min=1,max=2000"`

	ProfileImageURLstr *string `json:"profile_image_url_str" validate:"min=1,max=100"` // FIX!! - validation
	BannerImageURLstr  *string `json:"banner_image_url_str"  validate:"min=1,max=100"` // FIX!! - validation
}

type GFuserHTTPinputEmailConfirm struct {
	UserNameStr    GFuserName `validate:"required,min=3,max=50"`
	ConfirmCodeStr string     `validate:"required,min=10,max=20"`
}

type GFuserHTTPinputEmailLogin struct {
	EmailStr *string  `json:"email_str" validate:"min=6,max=50"`
}

type GFuserHTTPinputEmailLoginConfirm struct {
	UserNameStr    GFuserName `validate:"required,min=3,max=50"`
	ConfirmCodeStr string     `validate:"required,min=10,max=20"`
}

//---------------------------------------------------

func GetUserIDfromReq(pReq *http.Request,
	pAuthSubsystemTypeStr  string,
	pCtx		context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, gf_core.GF_ID, *gf_core.GFerror) {


	sessionID, sessionIDfoundBool := GetSessionID(pReq, pRuntimeSys)
	if !sessionIDfoundBool {
		return false, gf_core.GF_ID(""), nil
	}

	// AUTH0
	if pAuthSubsystemTypeStr == GF_AUTH_SUBSYSTEM_TYPE__AUTH0 {
		auth0session, gfErr := DBsqlAuth0getSession(sessionID, pCtx, pRuntimeSys)
		if gfErr != nil {
			return false, gf_core.GF_ID(""), gfErr
		}

		userID := auth0session.UserID

		return true, userID, nil
	}

	return false, gf_core.GF_ID(""), nil
}

//---------------------------------------------------

func ResolveUserName(pUserID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) GFuserName {

	var userNameStr GFuserName

	/*
	LEGACY!! - old images dont have a user_id associated with them.
		before the user system was fully integrated into gf_images, images were added anonimously
		and did not have a user ID associated with them.
		for those images it is not possible to associate them with user_names. 
	*/
	if pUserID != "" {

		userID := pUserID

		resolvedUserNameStr, gfErr := DBsqlGetUserNameByID(userID, pCtx, pRuntimeSys)
		if gfErr != nil {

			// failing to resolve user_name should still return a user_name
			userNameStr = GFuserName("?")
			return userNameStr
		}
		userNameStr = resolvedUserNameStr
		
	} else {

		// IMPORTANT!! - pre-auth-system images are marked as owned by anonymous users.
		userNameStr = GFuserName("anon")
	}

	return userNameStr
}

//---------------------------------------------------
// CREATE_ID

func usersCreateID(pUserIdentifierStr string,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		pUserIdentifierStr,
	}
	gfIDstr := gf_core.IDcreate(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfIDstr
}

//---------------------------------------------------
// GET_USER_NAME_FROM_CTX

func GetUserIDfromCtx(pCtx context.Context) (gf_core.GF_ID, bool) {
	
	userID := pCtx.Value("gf_user_id")

	if userID != nil {
		userID := gf_core.GF_ID(userID.(string))
		return userID, true
	}
	
	return gf_core.GF_ID(""), false
}

//---------------------------------------------------
// GET_SESSION_ID_FROM_CTX

func GetSessionIDfromCtx(pCtx context.Context) (gf_core.GF_ID, bool) {
	
	sessionIDval := pCtx.Value("gf_session_id")

	if sessionIDval != nil {
		sessionID := gf_core.GF_ID(sessionIDval.(string))
		return sessionID, true
	}
	
	return gf_core.GF_ID(""), false
}

//---------------------------------------------------
// HTTP
//---------------------------------------------------

func HTTPgetEmailLoginInput(pReq *http.Request) (*GFuserHTTPinputEmailLogin, *gf_core.GFerror) {
	input := &GFuserHTTPinputEmailLogin{}
	return input, nil
}

//---------------------------------------------------
// GET_AUTH_SUBSYSTEM_TYPE

func HTTPgetAuthSubsystemType(pReq *http.Request) string {
	authSubsystemTypeStr := pReq.Header.Get("gf_auth_type")
	return authSubsystemTypeStr
}

//---------------------------------------------------

func HTTPgetUserStdInput(pCtx context.Context,
	pReq        *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, gf_core.GF_ID, GFuserAddressETH, *gf_core.GFerror) {

	inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
	if gfErr != nil {
		return nil, "", GFuserAddressETH(""), gfErr
	}
	
	// user-name is supplied if the traditional auth system is used, and not web3/eth
	var userIDstr gf_core.GF_ID
	if inputUserIDstr, ok := inputMap["user_id_str"].(gf_core.GF_ID); ok {
		userIDstr = inputUserIDstr
	} else {

		// logged in users are added to context by gf_rpc, not supplied explicitly
		// via http request input (as they are for unauthenticated requests).
		userIDfromCtxStr, ok := GetUserIDfromCtx(pCtx) // pCtx.Value("gf_user_name").(string)
		if ok {
			userIDstr = userIDfromCtxStr
		}
	}

	pRuntimeSys.LogNewFun("DEBUG", "getting HTTP user std input", map[string]interface{}{"user_id_str": userIDstr,})

	// users eth address is used if the user picks that method instead of traditional
	var userAddressETHstr string;
	if inputUserAddressETHstr, ok := inputMap["user_address_eth_str"].(string); ok {
		userAddressETHstr = inputUserAddressETHstr
	}

	// one of the these values has to be supplied, they cant both be missing
	if userIDstr == "" && userAddressETHstr == "" {
		gfErr := gf_core.MongoHandleError("user_name_str or user_address_eth_str arguments are missing from request",
			"verify__input_data_missing_in_req_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return nil, "", GFuserAddressETH(""), gfErr
	}

	return inputMap, userIDstr, GFuserAddressETH(userAddressETHstr), nil
}

//---------------------------------------------------

func HTTPgetUserAddressETHinput(pReq *http.Request,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (GFuserAddressETH, *gf_core.GFerror) {

	queryArgsMap := pReq.URL.Query()
	if valuesLst, ok := queryArgsMap["addr_eth"]; ok {
		return GFuserAddressETH(valuesLst[0]), nil
	} else {
		gfErr := gf_core.ErrorCreate("incoming http request is missing the addr_eth query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return GFuserAddressETH(""), gfErr
	}
	return GFuserAddressETH(""), nil
}

//---------------------------------------------------

func HTTPgetUserUpdateInput(pReq *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserHTTPinputUpdate, *gf_core.GFerror) {

	handlerURLpathStr := pReq.URL.Path
	input             := GFuserHTTPinputUpdate{}
	bodyBytesLst, _   := ioutil.ReadAll(pReq.Body)
	err               := json.Unmarshal(bodyBytesLst, &input)
		
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to parse json http input for user update",
			"json_decode_error",
			map[string]interface{}{"handler_url_path_str": handlerURLpathStr,},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	return &input, nil
}

//---------------------------------------------------
// GET_EMAIL_CONFIRM_INPUT

func HTTPgetEmailConfirmInput(pReq *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserHTTPinputEmailConfirm, *gf_core.GFerror) {

	var userNameStr         string
	var confirmationCodeStr string

	queryArgsMap := pReq.URL.Query()
		
	//----------------------------
	// USER_NAME
	if valuesLst, ok := queryArgsMap["u"]; ok {
		userNameStr = valuesLst[0]
	} else {
		gfErr := gf_core.ErrorCreate("incoming http request is missing the email user_name query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	//----------------------------
	// CONFIRMATION_CODE
	if valuesLst, ok := queryArgsMap["c"]; ok {
		confirmationCodeStr = valuesLst[0]
	} else {
		gfErr := gf_core.ErrorCreate("incoming http request is missing the email confirmation_code query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	//----------------------------
	input := &GFuserHTTPinputEmailConfirm{
		UserNameStr:    GFuserName(userNameStr),
		ConfirmCodeStr: confirmationCodeStr,
	}

	return input, nil
}

//---------------------------------------------------
// GET_EMAIL_LOGIN_CONFIRM_INPUT

func HTTPgetEmailLoginConfirmInput(pReq *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserHTTPinputEmailLoginConfirm, *gf_core.GFerror) {

	var userNameStr         string
	var confirmationCodeStr string

	queryArgsMap := pReq.URL.Query()
	
	if valuesLst, ok := queryArgsMap["u"]; ok {
		userNameStr = valuesLst[0]
	} else {
		gfErr := gf_core.ErrorCreate("incoming http request is missing the email user_name query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	if valuesLst, ok := queryArgsMap["c"]; ok {
		confirmationCodeStr = valuesLst[0]
	} else {
		gfErr := gf_core.ErrorCreate("incoming http request is missing the email confirmation_code query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	input := &GFuserHTTPinputEmailLoginConfirm{
		UserNameStr:    GFuserName(userNameStr),
		ConfirmCodeStr: confirmationCodeStr,
	}

	return input, nil
}

//---------------------------------------------------