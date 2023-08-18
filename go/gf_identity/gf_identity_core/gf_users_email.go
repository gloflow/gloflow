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

package gf_identity_core

import (
	"fmt"
	"context"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
)

//---------------------------------------------------

func UsersEmailPipelineVerify(pEmailAddressStr string,
	pUserNameStr   GFuserName,
	pUserIDstr     gf_core.GF_ID,
	pDomainBaseStr string,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//------------------------
	// EMAIL_CONFIRM

	confirmCodeStr := usersEmailGenerateConfirmationCode()

	// DB
	gfErr := dbSQLuserEmailConfirmCreate(pUserNameStr,
		pUserIDstr,
		confirmCodeStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	msgSubjectStr, msgBodyHTMLstr, msgBodyTextStr := usersEmailGetConfirmMsgInfo(pUserNameStr,
		confirmCodeStr,
		pDomainBaseStr)

	// sender address
	senderAddressStr := fmt.Sprintf("gf-email-confirm@%s", pDomainBaseStr)

	gfErr = gf_aws.SESsendMessage(pEmailAddressStr,
		senderAddressStr,
		msgSubjectStr,
		msgBodyHTMLstr,
		msgBodyTextStr,
		pRuntimeSys)
	
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	return nil
}

//---------------------------------------------------

func UsersEmailPipelineConfirm(pInput *GFuserHTTPinputEmailConfirm,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {

	dbConfirmCodeStr, expiredBool, gfErr := usersEmailGetConfirmationCode(pInput.UserNameStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return false, "", gfErr
	}
	
	if expiredBool {
		return false, "email confirmation code has expired", nil
	}

	// confirm_code is correct
	if pInput.ConfirmCodeStr == dbConfirmCodeStr {
		
		// GET_USER_ID
		userIDstr, gfErr := DBsqlGetBasicInfoByUsername(pInput.UserNameStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		//------------------------
		// initial user email confirmation. only for new users.
		// user confirmed their email as valid.
		userEmailConfirmedBool, gfErr := DBsqlUserEmailIsConfirmed(pInput.UserNameStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		// if users email is not already marked as confirmed in the DB, updated
		// it to mark it as confirmed.
		if !userEmailConfirmedBool {
			
			emailIsConfirmedBool := true
			updateOp := &GFuserUpdateOp{
				EmailConfirmedBool: &emailIsConfirmedBool,
			}
	
			// UPDATE_USER - mark user as email_confirmed
			gfErr = DBsqlUserUpdate(userIDstr,
				updateOp,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return false, "", gfErr
			}
		}

		//------------------------

		//------------------------
		// UPDATE_LOGIN_ATTEMPT
		// if email is confirmed then update the login_attempt

		// get a preexisting login_attempt if one exists and hasnt expired for this user.
		// if it has then a new one will have to be created.
		var loginAttempt *GFloginAttempt
		loginAttempt, gfErr = LoginAttemptGetIfValid(GFuserName(pInput.UserNameStr),
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		if loginAttempt == nil {
			return false, "login_attempt for this user has not been found in the DB, to mark its email_confirmed flag as true", nil
		}

		loginEmailConfirmedBool := true
		updateOp := &GFloginAttemptUpdateOp{EmailConfirmedBool: &loginEmailConfirmedBool}
		gfErr = DBsqlLoginAttemptUpdate(&loginAttempt.IDstr,
			updateOp,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, "", gfErr
		}

		//------------------------

		return true, "", nil

	} else {
		return false, "received confirm code and DB confirm code are not the same", nil
	}
	return false, "", nil
}

//---------------------------------------------------

func usersEmailGetConfirmationCode(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, bool, *gf_core.GFerror) {

	expiredBool := false

	confirmCodeStr, confirmCodeCreationTimeF, gfErr := dbSQLuserEmailConfirmGetCode(pUserNameStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", expiredBool, gfErr
	}

	//------------------------
	// check confirm_code didnt expire
	currentUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	confirmCodeAgeTimeF := currentUNIXtimeF - confirmCodeCreationTimeF

	// check if older than 5min
	if (5.0 < confirmCodeAgeTimeF / 60) {
		expiredBool = true
		return "", expiredBool, nil
	}

	//------------------------

	return confirmCodeStr, expiredBool, nil
}

//---------------------------------------------------

func usersEmailGenerateConfirmationCode() string {
	cStr := fmt.Sprintf("%s:%s", gf_core.StrRandom(), gf_core.StrRandom())
	return cStr
}

//---------------------------------------------------

func usersEmailGetConfirmMsgInfo(pUserNameStr GFuserName,
	pConfirmCodeStr string,
	pDomainStr      string) (string, string, string) {

	subjectStr := fmt.Sprintf("%s - confirm your email", pDomainStr)
	
	welcomeMsgStr := fmt.Sprintf(`
		<div id="welcome_message" style="
			margin-left: 10px;
			padding-top: 9px;
			margin-top:  33px;">
			Hello <span style="font-weight: bold;">%s</span>.
			Welcome to %s!</div>
		<div>`,
		pUserNameStr,
		pDomainStr)

	confirmBtnStr := fmt.Sprintf(`
		<div id="confirm_btn" style="
			width: 100%%;
			background-color: #e17d44;
			text-align: center;
			padding-top: 12px;
			padding-bottom: 14px;
			cursor: pointer;">
			<a style="color: white; cursor: pointer;text-decoration: none;" href="https://%s/v1/identity/email_confirm?u=%s&c=%s">confirm email</a>
		</div>`,
		pDomainStr,
		pUserNameStr,
		pConfirmCodeStr)

	htmlStr := fmt.Sprintf(`
		<div>
			<style>
				body {
					margin:      0px;
					font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;

					/*turn off horizontal scroll*/
					overflow-x: hidden;

					background-color: #d7d7d7;
				}
				
			</style>
			<div id='gf_logo' style="margin-top: 75px;">
				<img src="https://gloflow.com/images/d/gf_logo_0.3.png"></img>
			</div>
			<div>
				%s
			</div>
			<div id="confirm_email" style="
				background-color: rgb(234, 206, 196);
				margin-top: 50px;
				padding: 10px;
				width: 360px;
				margin-bottom: 50px;">
				
				<div style="
					font-size:     12px;
					margin-bottom: 8px;">Please click on the bellow link to confirm your email address.
				</div>

				%s
				
			</div>
			<div>
				<div id="message" style="
					margin-top: 5px;
					margin-bottom: 5px;
					padding-left: 11px;">
					"There is no spoon ...it is only yourself."
				</div>
				<img src="https://gloflow.com/images/d/thumbnails/b2373f98d61208c60155fce191399f9f_thumb_large.png"></img>
			</div>
			<div style="font-size: 10px; padding: 3px; padding-left: 7px; margin-top: 140px;">
				don't reply to this email
			</div>
		</div>`,
		welcomeMsgStr,
		confirmBtnStr)

	textStr := fmt.Sprintf(`
		Welcome to %s!
		There is no spoon. ...it is only yourself.

		Please open the following link in your browser to confirm your email address.
		
		https://%s/v1/identity/email_confirm?c=%s`,
		pDomainStr,
		pDomainStr,
		pConfirmCodeStr)

	return subjectStr, htmlStr, textStr
}