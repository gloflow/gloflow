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
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFloginAttempt struct {
	Vstr              string             `bson:"v_str"` // schema_version
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	IDstr             gf_core.GF_ID      `bson:"id_str"`
	DeletedBool       bool               `bson:"deleted_bool"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`

	UserTypeStr        string     `bson:"user_type_str"` // "regular"|"admin"
	UserNameStr        GFuserName `bson:"user_name_str"`
	
	PassConfirmedBool  bool `bson:"pass_confirmed_bool"`
	EmailConfirmedBool bool `bson:"email_confirmed_bool"`
	MFAconfirmedBool   bool `bson:"mfa_confirmed_bool"`
}

//---------------------------------------------------

func LoginAttemptGetOrCreate(pUserNameStr GFuserName,
	pUserTypeStr string,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFloginAttempt, *gf_core.GFerror) {
	
	// GET
	// get a preexisting login_attempt if one exists and hasnt expired for this user.
	// if it has then a new one will have to be created.
	var loginAttempt *GFloginAttempt
	loginAttempt, gfErr := LoginAttemptGetIfValid(pUserNameStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	if loginAttempt == nil {

		//------------------------
		// CREATE_LOGIN_ATTEMPT

		loginAttempt, gfErr = loginAttempCreate(pUserNameStr, pUserTypeStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		//------------------------
	}

	return loginAttempt, nil
}

//---------------------------------------------------
// CREATE

func loginAttempCreate(pUserNameStr GFuserName,
	pUserTypeStr string,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFloginAttempt, *gf_core.GFerror) {

	userIdentifierStr := string(pUserNameStr)
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	loginAttemptIDstr := usersCreateID(userIdentifierStr, creationUNIXtimeF)

	loginAttempt := &GFloginAttempt{
		Vstr:              "0",
		IDstr:             loginAttemptIDstr,
		CreationUNIXtimeF: creationUNIXtimeF,
		UserTypeStr:       pUserTypeStr,
		UserNameStr:       pUserNameStr,
	}
	gfErr := dbLoginAttemptCreate(loginAttempt,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return loginAttempt, nil
}

//---------------------------------------------------
// GET_IF_VALID

func LoginAttemptGetIfValid(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFloginAttempt, *gf_core.GFerror) {

	loginAttemptMaxAgeSecondsF := 5*60.0
	
	var loginAttempt *GFloginAttempt
	loginAttempt, gfErr := dbLoginAttemptGetByUsername(pUserNameStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	if loginAttempt == nil {
		return nil, nil
	}


	// get the age of the login_attempt
	currentUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	ageF             := currentUNIXtimeF - loginAttempt.CreationUNIXtimeF


	// login_attempt has expired
	if ageF > loginAttemptMaxAgeSecondsF {

		// mark it as deleted
		expiredBool := true
		updateOp := &GFloginAttemptUpdateOp{DeletedBool: &expiredBool}
		gfErr = DBloginAttemptUpdate(&loginAttempt.IDstr,
			updateOp,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		return nil, nil
	} else {
		return loginAttempt, nil
	}

	return nil, nil
}