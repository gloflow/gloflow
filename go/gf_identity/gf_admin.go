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

package gf_identity

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//---------------------------------------------------

type GFadminInputLogin struct {

	UserNameStr gf_identity_core.GFuserName `validate:"required,min=3,max=50"`

	// pass is not provided if email-login is used
	PassStr string `validate:"omitempty,min=8,max=50"`

	// admin email
	EmailStr string `validate:"omitempty,email"`
}

type GFadminOutputLogin struct {
	UserExistsBool     bool
	EmailConfirmedBool bool
	PassValidBool      bool
	UserIDstr          gf_core.GF_ID 
}

type GFadminOutputCreateAdmin struct {
	General *gf_identity_core.GFuserpassOutputCreate
}

type GFadminInputAddToInviteList struct {
	AdminUserID gf_core.GF_ID `validate:"required,min=3,max=50"`
	EmailStr    string        `validate:"required,email"`
}

type GFadminRemoveFromInviteListInput struct {
	AdminUserID gf_core.GF_ID `validate:"required,min=3,max=50"`
	EmailStr    string        `validate:"required,email"`
}

type GFadminUserViewOutput struct {
	ID                 gf_core.GF_ID                       `json:"id_str"`
	CreationUNIXtimeF  float64                             `json:"creation_unix_time_f"`
	UserNameStr        gf_identity_core.GFuserName         `json:"user_name_str"`
	ScreenNameStr      string                              `json:"screen_name_str"`
	AddressesETHlst    []gf_identity_core.GFuserAddressETH `json:"addresses_eth_lst"`
	EmailStr           string                              `json:"email_str"`
	EmailConfirmedBool bool                                `json:"email_confirmed_bool"`
	ProfileImageURLstr string                              `json:"profile_image_url_str"`
}

type GFadminResendConfirmEmailInput struct {
	UserID      gf_core.GF_ID               `validate:"required,min=3,max=50"`
	UserNameStr gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
	EmailStr    string                      `validate:"required,email"`
}

type GFadminUserDeleteInput struct {
	UserID      gf_core.GF_ID               `validate:"required,min=3,max=50"`
	UserNameStr gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
}

//------------------------------------------------
// PIPELINE_DELETE_USER

func AdminPipelineDeleteUser(pInput *GFadminUserDeleteInput,
	pCtx         context.Context,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	deletedBool := true
	updateOp := &gf_identity_core.GFuserUpdateOp{
		DeletedBool: &deletedBool,
	}

	// UPDATE_USER - mark user as email_confirmed
	gfErr := gf_identity_core.DBsqlUserUpdate(pInput.UserID,
		updateOp,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//------------------------------------------------

func AdminPipelineUserResendConfirmEmail(pInput *GFadminResendConfirmEmailInput,
	pCtx         context.Context,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	if pServiceInfo.EnableEmailBool {

		gfErr = gf_identity_core.UsersEmailPipelineVerify(pInput.EmailStr,
			pInput.UserNameStr,
			pInput.UserID,
			pServiceInfo.DomainBaseStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	return nil
}

//------------------------------------------------

func AdminPipelineGetAllUsers(pCtx context.Context,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) ([]*GFadminUserViewOutput, *gf_core.GFerror) {

	// DB
	usersLst, gfErr := gf_identity_core.DBsqlUserGetAll(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	outputLst := []*GFadminUserViewOutput{}
	for _, user := range usersLst {


		userView := &GFadminUserViewOutput{
			ID:                 user.ID,
			CreationUNIXtimeF:  user.CreationUNIXtimeF,
			UserNameStr:        user.UserNameStr,
			ScreenNameStr:      user.ScreenNameStr,
			AddressesETHlst:    user.AddressesETHlst,
			EmailStr:           user.EmailStr,
			EmailConfirmedBool: user.EmailConfirmedBool,
			ProfileImageURLstr: user.ProfileImageURLstr,
		}

		outputLst = append(outputLst, userView)
	}

	return outputLst, nil
}

//------------------------------------------------

func AdminPipelineGetAllInviteList(pCtx context.Context,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	// DB
	dbInviteListLst, gfErr := gf_identity_core.DBsqlUserGetAllInInviteList(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	inviteListLst := []map[string]interface{}{}
	for _, inviteMap := range dbInviteListLst {

		inviteListLst = append(inviteListLst, map[string]interface{}{
			"user_email_str":       inviteMap["user_email_str"],
			"creation_unix_time_f": inviteMap["creation_unix_time_f"],
		})
	}

	return inviteListLst, nil
}

//------------------------------------------------

func AdminPipelineUserAddToInviteList(pInput *GFadminInputAddToInviteList,
	pCtx         context.Context,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	gfErr = gf_identity_core.DBsqlUserAddToInviteList(pInput.EmailStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		
		adminUserNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(pInput.AdminUserID, pCtx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		eventMetaMap := map[string]interface{}{
			"user_id":                    pInput.AdminUserID,
			"user_name":                  adminUserNameStr,
			"email_added_to_invite_list": pInput.EmailStr,
		}
		gf_events.EmitApp(gf_identity_core.GF_EVENT_APP__ADMIN_ADDED_USER_TO_INVITE_LIST,
			eventMetaMap,
			pRuntimeSys)
	}

	return nil
}

//------------------------------------------------

func AdminPipelineUserRemoveFromInviteList(pInput *GFadminRemoveFromInviteListInput,
	pCtx         context.Context,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	gfErr = gf_identity_core.DBsqlUserRemoveFromInviteList(pInput.EmailStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		
		adminUserNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(pInput.AdminUserID, pCtx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		eventMetaMap := map[string]interface{}{
			"user_id":                    pInput.AdminUserID,
			"user_name":                  adminUserNameStr,
			"email_added_to_invite_list": pInput.EmailStr,
		}
		gf_events.EmitApp(gf_identity_core.GF_EVENT_APP__ADMIN_REMOVED_USER_FROM_INVITE_LIST,
			eventMetaMap,
			pRuntimeSys)
	}

	return nil
}

//---------------------------------------------------
// PIPELINE__LOGIN

// this function is entered mutliple times for complex logins where not only pass/eth_signature
// are verified, but where email/mfa have to be confirmed as well.
// for each of the login stages this function is entered, and the login_attempt record
// is used to keep track of which stages have completed.

func AdminPipelineLogin(pInput *GFadminInputLogin,
	pCtx         context.Context,
	pLocalHub    *sentry.Hub,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFadminOutputLogin, *gf_core.GFerror) {
	
	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	userNameStr := gf_identity_core.GFuserName(pInput.UserNameStr)
	output := &GFadminOutputLogin{}

	//------------------------
	// VERIFY

	userExistsBool, gfErr := gf_identity_core.DBsqlUserExistsByUsername(userNameStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	
	// BREADCRUMB
	gf_core.Breadcrumbs__add("auth", "admin user checked for existence",
		map[string]interface{}{"user_exists_bool": userExistsBool, "user_name_str": pInput.UserNameStr},
		pLocalHub)

	var userID gf_core.GF_ID
	
	// admin user doesnt exist
	if !userExistsBool {
		
		// so create it but only if its the root admin user.
		// other admin users have to be created explicitly
		if pInput.UserNameStr == "admin" {

			//------------------------	
			// PIPELINE__CREATE_ADMIN
			// if the admin user doesnt exist in the DB (most likely on first run of gloflow server),
			// create one in the DB

			inputCreate := &gf_identity_core.GFuserpassInputCreate{
				UserNameStr: userNameStr,
				PassStr:     pInput.PassStr,
				EmailStr:    pInput.EmailStr,
				UserTypeStr: "admin",
			}

			// BREADCRUMB
			gf_core.Breadcrumbs__add("auth", "creating new admin user",
				map[string]interface{}{"email_str": pInput.EmailStr, "user_name_str": pInput.UserNameStr},
				pLocalHub)
			
			// CREATE
			output_create, gfErr := adminPipelineCreateAdmin(inputCreate,
				pServiceInfo,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			//------------------------

			userID = output_create.General.UserID
			output.UserExistsBool = true
		} else {

			output.UserExistsBool = false
			return output, nil
		}
	
	} else {
		existingUserIDstr, gfErr := gf_identity_core.DBsqlGetBasicInfoByUsername(userNameStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		userID                = existingUserIDstr
		output.UserExistsBool = true
	}

	// BREADCRUMB
	gf_core.Breadcrumbs__add("auth", "got user_id for admin user",
		map[string]interface{}{"user_id_str": userID, "user_name_str": pInput.UserNameStr},
		pLocalHub)

	//------------------------
	// LOGIN_ATTEMPT

	userTypeStr := "admin"
	loginAttempt, gfErr := gf_identity_core.LoginAttemptGetOrCreate(userNameStr, userTypeStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//------------------------
	// VERIFY_PASSWORD

	// only verify password if the login_attempt didnt mark it yet as complete
	if !loginAttempt.PassConfirmedBool {

		passValidBool, gfErr := gf_identity_core.UserpassVerifyPass(userNameStr,
			pInput.PassStr,
			pServiceInfo,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		if !passValidBool {
			output.PassValidBool = false
			return output, nil
		} else {
			output.PassValidBool = true

			//------------------------
			// UPDATE_LOGIN_ATTEMPT
			// if password is valid then update the login_attempt 
			// to indicate that the password has been confirmed
			updateOp := &gf_identity_core.GFloginAttemptUpdateOp{PassConfirmedBool: &passValidBool}
			gfErr = gf_identity_core.DBsqlLoginAttemptUpdate(&loginAttempt.IDstr,
				updateOp,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			//------------------------

			// EVENT
			if pServiceInfo.EnableEventsAppBool {
				eventMeta := map[string]interface{}{
					"user_id":     userID,
					"user_name":   pInput.UserNameStr,
					"domain_base": pServiceInfo.DomainBaseStr,
				}
				gf_events.EmitApp(gf_identity_core.GF_EVENT_APP__ADMIN_LOGIN_PASS_CONFIRMED,
					eventMeta,
					pRuntimeSys)
			}
		}
	}

	//------------------------
	// EMAIL
	if pServiceInfo.EnableEmailBool {

		// go through the email verification pipeline if the email
		// has not yet been confirmed
		if !loginAttempt.EmailConfirmedBool {

			gfErr = gf_identity_core.UsersEmailPipelineVerify(pInput.EmailStr,
				userNameStr,
				userID,
				pServiceInfo.DomainBaseStr,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			// EVENT
			if pServiceInfo.EnableEventsAppBool {
				eventMeta := map[string]interface{}{
					"user_id":     userID,
					"user_name":   pInput.UserNameStr,
					"domain_base": pServiceInfo.DomainBaseStr,
				}
				gf_events.EmitApp(gf_identity_core.GF_EVENT_APP__ADMIN_LOGIN_EMAIL_VERIFICATION_SENT,
					eventMeta,
					pRuntimeSys)
			}

			//------------------------
		}
	}

	//------------------------
	
	return output, nil
}

//---------------------------------------------------

func adminPipelineCreateAdmin(pInput *gf_identity_core.GFuserpassInputCreate,
	pServiceInfo *gf_identity_core.GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFadminOutputCreateAdmin, *gf_core.GFerror) {

	//------------------------
	// PIPELINE
	output, gfErr := gf_identity_core.UserpassPipelineCreate(pInput,
		pServiceInfo,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		eventMeta := map[string]interface{}{
			"user_id":     output.UserID,
			"user_name":   pInput.UserNameStr,
			"user_type":   pInput.UserTypeStr,
			"domain_base": pServiceInfo.DomainBaseStr,
		}
		gf_events.EmitApp(gf_identity_core.GF_EVENT_APP__ADMIN_CREATE,
			eventMeta,
			pRuntimeSys)
	}

	//------------------------

	outputAdmin := &GFadminOutputCreateAdmin{
		General: output,
	}
	return outputAdmin, nil
}

//---------------------------------------------------

func AdminIs(pUserIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	userNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(pUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	if string(userNameStr) != "admin" {
		gfErr := gf_core.ErrorCreate("username thats not 'admin' is trying to login as admin",
			"verify__invalid_value_error",
			map[string]interface{}{
				"user_name_str": userNameStr,
			},
			nil, "gf_identity", pRuntimeSys)
		return gfErr
	}
	return nil
}