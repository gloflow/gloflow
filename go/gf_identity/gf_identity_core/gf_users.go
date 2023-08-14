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
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFuserName       string
type GFauthSignature  string
type GFuserAddressETH string

type GFuser struct {
	Vstr              string             `bson:"v_str"` // schema_version
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	ID                gf_core.GF_ID      `bson:"id_str"`
	DeletedBool       bool               `bson:"deleted_bool"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`

	UserTypeStr     string     `bson:"user_type_str"`   // "admin" | "standard"
	UserNameStr     GFuserName `bson:"user_name_str"`   // set once at the creation of the user
	ScreenNameStr   string     `bson:"screen_name_str"` // changable durring the lifetime of the user
	
	DescriptionStr  string             `bson:"description_str"`
	AddressesETHlst []GFuserAddressETH `bson:"addresses_eth_lst"`

	EmailStr           string `bson:"email_str"`
	EmailConfirmedBool bool   `bson:"email_confirmed_bool"` // one-time confirmation on user-creation to validate user
	
	// IMAGES
	ProfileImageURLstr string `bson:"profile_image_url_str"`
	BannerImageURLstr  string `bson:"banner_image_url_str"`
}

// ADD!! - provide logic/plugin for storing this record in some alternative store
//         separate from the main DB
type GFuserCreds struct {
	Vstr              string             `bson:"v_str"` // schema_version
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	ID                gf_core.GF_ID      `bson:"id_str"`
	DeletedBool       bool               `bson:"deleted_bool"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`

	UserID      gf_core.GF_ID            `bson:"user_id_str"`
	UserNameStr GFuserName               `bson:"user_name_str"`
	PassSaltStr string                   `bson:"pass_salt_str"`
	PassHashStr string                   `bson:"pass_hash_str"`
}

// io_update
type GFuserInputUpdate struct {
	UserID            gf_core.GF_ID    `validate:"required"`                 // required - not updated, but for lookup
	UserAddressETHstr GFuserAddressETH `validate:"omitempty,eth_addr"`       // optional - add an Eth address to the user
	ScreenNameStr     *string          `validate:"omitempty,min=3,max=50"`   // optional
	EmailStr          *string          `validate:"omitempty,email"`          // optional
	DescriptionStr    *string          `validate:"omitempty,min=1,max=2000"` // optional

	ProfileImageURLstr *string `validate:"omitempty,min=1,max=100"` // optional // FIX!! - validation
	BannerImageURLstr  *string `validate:"omitempty,min=1,max=100"` // optional // FIX!! - validation
}
type GFuserOutputUpdate struct {
	
}

type GFuserInputGet struct {
	UserID gf_core.GF_ID
}

type GFuserOutputGet struct {
	UserNameStr        GFuserName
	EmailStr           string
	DescriptionStr     string
	ProfileImageURLstr string
	BannerImageURLstr  string
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
// PIPELINE__UPDATE

func UsersPipelineUpdate(pInput *GFuserInputUpdate,
	pServiceInfo *GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFuserOutputUpdate, *gf_core.GFerror) {
	
	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// USER_NAME
	userNameStr, gfErr := DBsqlGetUserNameByID(pInput.UserID, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// EMAIL
	if pServiceInfo.EnableEmailBool {
		if *pInput.EmailStr != "" {

			gfErr = UsersEmailPipelineVerify(*pInput.EmailStr,
				userNameStr,
				pInput.UserID,
				pServiceInfo.DomainBaseStr,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}
		}
	}

	//------------------------

	output := &GFuserOutputUpdate{}

	return output, nil
}

//---------------------------------------------------
// PIPELINE__GET

func UsersPipelineGet(pInput *GFuserInputGet,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuserOutputGet, *gf_core.GFerror) {

	//------------------------
	// VALIDATE
	gfErr := gf_core.ValidateStruct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	
	user, gfErr := DBsqlUserGetByID(pInput.UserID,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output := &GFuserOutputGet{
		UserNameStr:        user.UserNameStr,
		EmailStr:           user.EmailStr,
		DescriptionStr:     user.DescriptionStr,
		ProfileImageURLstr: user.ProfileImageURLstr,
		BannerImageURLstr:  user.BannerImageURLstr,
	}

	return output, nil
}