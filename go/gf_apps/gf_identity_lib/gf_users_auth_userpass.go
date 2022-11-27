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

package gf_identity_lib

import (
	"fmt"
	"time"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// io_login

type GFuserAuthUserpassInputLogin struct {

	// username is always required, with both pass and email login
	UserNameStr gf_identity_core.GFuserName `validate:"required,min=3,max=50"`

	// pass is not provided if email-login is used
	PassStr string `validate:"omitempty,min=8,max=50"`

	// for certain emails allow email-login
	EmailStr string `validate:"omitempty,email"`
}
type GFuserAuthUserpassOutputLogin struct {
	UserExistsBool     bool
	EmailConfirmedBool bool
	PassValidBool      bool
	UserIDstr          gf_core.GF_ID 
	JWTtokenVal        gf_session.GFjwtTokenVal
}

// io_login_finalize
type GFuserAuthUserpassInputLoginFinalize struct {
	UserNameStr gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
}
type GFuserAuthUserpassOutputLoginFinalize struct {
	EmailConfirmedBool bool
	UserIDstr          gf_core.GF_ID 
	JWTtokenVal        gf_session.GFjwtTokenVal
}

// io_create
type GFuserAuthUserpassInputCreate struct {
	UserNameStr gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
	PassStr     string                      `validate:"required,min=8,max=50"`
	EmailStr    string                      `validate:"required,email"`
	UserTypeStr string                      `validate:"required"` // "admin"|"standard"
}
type GFuserAuthUserpassOutputCreateRegular struct {
	UserExistsBool       bool
	UserInInviteListBool bool
	General              *GFuserAuthUserpassOutputCreate
}
type GFuserAuthUserpassOutputCreate struct {
	UserNameStr gf_identity_core.GFuserName
	UserIDstr   gf_core.GF_ID
}

//---------------------------------------------------
// PIPELINE__LOGIN

func usersAuthUserpassPipelineLogin(pInput *GFuserAuthUserpassInputLogin,
	pServiceInfo *GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFuserAuthUserpassOutputLogin, *gf_core.GFerror) {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	output := &GFuserAuthUserpassOutputLogin{}

	//------------------------
	// VERIFY

	userExistsBool, gfErr := dbUserExistsByUsername(gf_identity_core.GFuserName(pInput.UserNameStr),
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// user doesnt exists, so abort login
	if !userExistsBool {
		output.UserExistsBool = false
		return output, nil
	} else {
		output.UserExistsBool = true
	}

	//------------------------
	// VERIFY_PASSWORD
	passValidBool, gfErr := usersAuthUserpassVerifyPass(gf_identity_core.GFuserName(pInput.UserNameStr),
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
	}

	//------------------------
	// LOGIN_FINALIZE
	input := &GFuserAuthUserpassInputLoginFinalize{
		UserNameStr: gf_identity_core.GFuserName(pInput.UserNameStr),
	}
	loginFinalizeOutput, gfErr := usersAuthUserpassPipelineLoginFinalize(input,
		pServiceInfo,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	output.EmailConfirmedBool = loginFinalizeOutput.EmailConfirmedBool
	output.UserIDstr          = loginFinalizeOutput.UserIDstr
	output.JWTtokenVal        = loginFinalizeOutput.JWTtokenVal

	return output, nil
}

//---------------------------------------------------

func usersAuthUserpassPipelineLoginFinalize(pInput *GFuserAuthUserpassInputLoginFinalize,
	pServiceInfo *GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFuserAuthUserpassOutputLoginFinalize, *gf_core.GFerror) {

	output := &GFuserAuthUserpassOutputLoginFinalize{}
	userNameStr := gf_identity_core.GFuserName(pInput.UserNameStr)

	//------------------------
	// VERIFY_EMAIL_CONFIRMED
	// if this check is enabled, users that have not confirmed their email cant login.
	// this is the initial confirmation of an email on user creation, or user email update.
	if pServiceInfo.EnableEmailRequireConfirmForLoginBool {

		emailConfirmedBool, gfErr := dbUserGetEmailConfirmedByUsername(userNameStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		if !emailConfirmedBool {
			output.EmailConfirmedBool = false
			return output, nil
		} else {
			output.EmailConfirmedBool = true
		}
	}

	//------------------------
	// USER_ID
	
	userIDstr, gfErr := gf_identity_core.DBgetBasicInfoByUsername(userNameStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.UserIDstr = userIDstr

	//------------------------
	// JWT
	userIdentifierStr := string(userIDstr)
	JWTtokenVal, gfErr := gf_session.JWTpipelineGenerate(userIdentifierStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output.JWTtokenVal = JWTtokenVal

	//------------------------
	return output, nil
}

//---------------------------------------------------
// PIPELINE__CREATE_REGULAR

func usersAuthUserpassPipelineCreateRegular(pInput *GFuserAuthUserpassInputCreate,
	pServiceInfo *GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFuserAuthUserpassOutputCreateRegular, *gf_core.GFerror) {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	outputRegular := &GFuserAuthUserpassOutputCreateRegular{}

	//------------------------
	// VALIDATE

	userExistsBool, gfErr := dbUserExistsByUsername(pInput.UserNameStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// user already exists, so abort creation
	if userExistsBool {
		outputRegular.UserExistsBool = true
		return outputRegular, nil
	}

	// check if in invite list
	inInviteListBool, gfErr := dbUserCheckInInvitelistByEmail(pInput.EmailStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// user is not in the invite list, so abort the creation
	if inInviteListBool {
		outputRegular.UserInInviteListBool = true
	} else {
		outputRegular.UserInInviteListBool = false
		return outputRegular, nil
	}

	//------------------------
	// PIPELINE
	output, gfErr := usersAuthUserpassPipelineCreate(pInput,
		pServiceInfo,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	outputRegular.General = output

	//------------------------
	// EMAIL
	if pServiceInfo.EnableEmailBool {

		gfErr = usersEmailPipelineVerify(pInput.EmailStr,
			pInput.UserNameStr,
			output.UserIDstr,
			pServiceInfo.DomainBaseStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}
	
	//------------------------
	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		eventMeta := map[string]interface{}{
			"user_id_str":     output.UserIDstr,
			"user_name_str":   pInput.UserNameStr,
			"domain_base_str": pServiceInfo.DomainBaseStr,
		}
		gf_events.EmitApp(GF_EVENT_APP__USER_CREATE_REGULAR,
			eventMeta,
			pRuntimeSys)
	}

	//------------------------

	return outputRegular, nil
}

//---------------------------------------------------
// PIPELINE__CREATE

func usersAuthUserpassPipelineCreate(pInput *GFuserAuthUserpassInputCreate,
	pServiceInfo *GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFuserAuthUserpassOutputCreate, *gf_core.GFerror) {

	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	userTypeStr := pInput.UserTypeStr
	userNameStr := pInput.UserNameStr
	passStr     := pInput.PassStr
	emailStr    := pInput.EmailStr

	userIdentifierStr := string(userNameStr)
	userIDstr := usersCreateID(userIdentifierStr, creationUNIXtimeF)

	user := &GFuser{
		Vstr:              "0",
		IDstr:             userIDstr,
		CreationUNIXtimeF: creationUNIXtimeF,
		UserTypeStr:       userTypeStr,
		UserNameStr:       userNameStr,
		EmailStr:          emailStr,
	}

	
	passSaltStr := usersAuthUserpassGetPassSalt()
	passHashStr := usersAuthUserpassGetPassHash(passStr, passSaltStr)

	credsCreationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	userCredsIDstr         := usersCreateID(userIdentifierStr, credsCreationUNIXtimeF)

	userCreds := &GFuserCreds {
		Vstr:              "0",
		IDstr:             userCredsIDstr,
		CreationUNIXtimeF: credsCreationUNIXtimeF,
		UserIDstr:         userIDstr,
		UserNameStr:       userNameStr,
		PassSaltStr:       passSaltStr,
		PassHashStr:       passHashStr,
	}

	//------------------------
	// USER_PERSIST
	// DB__USER_CREATE
	gfErr = dbUserCreate(user, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// USER_CREDS_PERSIST

	// SECRETS_STORE
	if pServiceInfo.EnableUserCredsInSecretsStoreBool && 
		pRuntimeSys.ExternalPlugins.SecretCreateCallback != nil {

		secretNameStr := fmt.Sprintf("gf_user_creds@%s", userNameStr)
		secretDescriptionStr := fmt.Sprintf("user creds for a particular user")

		userCredsMap := map[string]interface{}{
			"user_creds_id_str":    userCredsIDstr, 
			"creation_unix_time_f": credsCreationUNIXtimeF,
			"user_id_str":          userIDstr,
			"user_name_str":        userNameStr,
			"pass_salt_str":        passSaltStr,
			"pass_hash_str":        passHashStr,
		}

		// SECRET_STORE__USER_CREDS_CREATE
		gfErr := pRuntimeSys.ExternalPlugins.SecretCreateCallback(secretNameStr,
			userCredsMap,
			secretDescriptionStr,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	} else {

		// DB__USER_CREDS_CREATE - otherwise use the regular DB
		gfErr = dbUserCredsCreate(userCreds, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}
	
	//------------------------
	// LOGIN_ATTEMPT
	// on user creation initiate a login process that completes after the user
	// confirms their email.
	_, gfErr = loginAttempCreate(userNameStr, userTypeStr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// EMAIL_VERIFY_ADDRESS
	if pServiceInfo.EnableEmailBool {

		// this SES email verification is done only once for a new email address,
		// so that SES allows sending to this email address.
		gfErr = gf_aws.AWS_SES__verify_address(emailStr,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}

	//------------------------

	output := &GFuserAuthUserpassOutputCreate{
		UserNameStr: userNameStr,
		UserIDstr:   userIDstr,
	}

	return output, nil
}

//---------------------------------------------------
// PASS
//---------------------------------------------------

func usersAuthUserpassVerifyPass(pUserNameStr gf_identity_core.GFuserName,
	pPassStr     string,
	pServiceInfo *GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	// GET_PASS_AND_SALT	
	var passSaltLoadedStr string
	var passHashLoadedStr string

	// SECRETS_STORE
	if pServiceInfo.EnableUserCredsInSecretsStoreBool && 
		pRuntimeSys.ExternalPlugins.SecretGetCallback != nil {

		secretNameStr := fmt.Sprintf("gf_user_creds@%s", pUserNameStr)

		// SECRET_GET
		secretMap, gfErr := pRuntimeSys.ExternalPlugins.SecretGetCallback(secretNameStr,
			pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}

		passSaltLoadedStr = secretMap["pass_salt_str"].(string)
		passHashLoadedStr = secretMap["pass_hash_str"].(string)
		
	} else {

		// DB
		dbPassSaltStr, dbPassHashStr, gfErr := dbUserCredsGetPassHash(pUserNameStr,
			pCtx, pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}

		passSaltLoadedStr = dbPassSaltStr
		passHashLoadedStr = dbPassHashStr
	}

	// GENERATE_PASS_HASH
	passHashExpectedStr := usersAuthUserpassGetPassHash(pPassStr, passSaltLoadedStr)


	if (passHashLoadedStr == passHashExpectedStr) {
		return true, nil
	} else {
		return false, nil
	}

	return false, nil
}

//---------------------------------------------------

func usersAuthUserpassGetPassHash(pPassStr string,
	pPassSaltStr string) string {

	saltedPassStr := fmt.Sprintf("%s:%s", pPassSaltStr, pPassStr)
	passHashStr   := gf_core.HashValSha256(saltedPassStr)
	return passHashStr
}

//---------------------------------------------------

func usersAuthUserpassGetPassSalt() string {
	randStr := gf_core.StrRandom()
	return randStr
}