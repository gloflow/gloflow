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
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFuserUpdateOp struct {
	DeletedBool        *bool // if nil dont update, else update to true/false
	UserNameStr        GFuserName
	DescriptionStr     string
	EmailStr           string
	EmailConfirmedBool bool
	MFAconfirmBool     *bool // if nil dont update, else update to true/false
}

type GFloginAttemptUpdateOp struct {
	PassConfirmedBool  *bool
	EmailConfirmedBool *bool
	MFAconfirmedBool   *bool
	DeletedBool        *bool
}

//---------------------------------------------------
// AUTH0
//---------------------------------------------------

func dbAuth0createNewSession(pAuth0session *GFauth0session,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	collNameStr := "gf_auth0_session"
	gfErr := gf_core.MongoInsert(pAuth0session,
		collNameStr,
		map[string]interface{}{
			"caller_err_msg_str": "failed to insert GFauth0session into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	
	return nil
}

//---------------------------------------------------


func dbAuth0GetSession(pGFsessionIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFauth0session, *gf_core.GFerror) {

	findOpts := options.FindOne()
	
	session := GFauth0session{}
	collNameStr := "gf_auth0_session"

	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx, bson.M{
			"id_str":       pGFsessionIDstr,
			"deleted_bool": false,
		},
		findOpts).Decode(&session)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to find Auth0 session by ID in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"auth0_session_id_str": pGFsessionIDstr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	return &session, nil
}

//---------------------------------------------------

func dbAuth0UpdateSession(pGFsessionIDstr gf_core.GF_ID,
	pLoginCompleteBool bool,
	pAccessTokenStr    string,
	pAuth0profileMap   map[string]interface{},
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// FIELDS
	fieldsTargets := bson.M{}
	fieldsTargets["login_complete_bool"] = pLoginCompleteBool
	fieldsTargets["access_token_str"] = pAccessTokenStr
	fieldsTargets["profile_map"]      = pAuth0profileMap

	//------------------------
	collNameStr := "gf_auth0_session"
	_, err := pRuntimeSys.Mongo_db.Collection(collNameStr).UpdateMany(pCtx, bson.M{
			"id_str":       pGFsessionIDstr,
			"deleted_bool": false,
		},
		bson.M{"$set": fieldsTargets})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to to update Auth0 session in the DB",
			"mongodb_update_error",
			map[string]interface{}{
				"id_str": pGFsessionIDstr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// USER
//---------------------------------------------------

func DBuserGetAll(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFuser, *gf_core.GFerror) {

	collNameStr := "gf_users"

	findOpts := options.Find()
	cursor, gfErr := gf_core.MongoFind(bson.M{
			"deleted_bool":  false,
		},
		findOpts,
		map[string]interface{}{
			"caller_err_msg_str": "failed to get all users records from the DB",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return nil, gfErr
	}
	
	// no login_attempt found for user
	if cursor == nil {
		return nil, nil
	}
	
	usersLst := []*GFuser{}
	err := cursor.All(pCtx, &usersLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get all users records from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	return usersLst, nil
}

//---------------------------------------------------
// CREATE_USER

func dbUserCreate(pUser *GFuser,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_users"

	gfErr := gf_core.MongoInsert(pUser,
		collNameStr,
		map[string]interface{}{
			"user_id_str":        pUser.IDstr,
			"user_name_str":      pUser.UserNameStr,
			"description_str":    pUser.DescriptionStr,
			"addresses_eth_lst":  pUser.AddressesETHlst, 
			"caller_err_msg_str": "failed to insert GF_user into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------
// UPDATE

func DBuserUpdate(pUserIDstr gf_core.GF_ID, // p_user_address_eth_str GF_user_address_eth,
	pUpdateOp   *GFuserUpdateOp,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// FIELDS
	fieldsTargets := bson.M{}

	// DeletedBool - is a pointer. if its not nil, then set
	//               the deleted_bool field to either true/false.
	if pUpdateOp.DeletedBool != nil {
		fieldsTargets["deleted_bool"] = pUpdateOp.DeletedBool
	}

	if string(pUpdateOp.UserNameStr) != "" {
		fieldsTargets["username_str"] = pUpdateOp.UserNameStr
	}

	if pUpdateOp.DescriptionStr != "" {
		fieldsTargets["description_str"] = pUpdateOp.DescriptionStr
	}

	if pUpdateOp.EmailStr != "" {
		fieldsTargets["email_str"] = pUpdateOp.EmailStr

		// IMPORTANT!! - if the email is changed then it needs to be confirmed
		//               again. this flag on the user can only be changed to false
		//               indirectly like here by upding the user email
		fieldsTargets["email_confirmed_bool"] = false
	}

	// email_confirmed_bool itself can only be updated to explicitly
	// if its being set to true
	if pUpdateOp.EmailConfirmedBool {
		fieldsTargets["email_confirmed_bool"] = true
	}
	
	// MFA_confirm_bool - is a pointer. if its not nil, then set
	//                    the MFA_confirm field to either true/false.
	if pUpdateOp.MFAconfirmBool != nil {
		fieldsTargets["mfa_confirm_bool"] = pUpdateOp.MFAconfirmBool
	}

	//------------------------
	
	_, err := pRuntimeSys.Mongo_db.Collection("gf_users").UpdateMany(pCtx, bson.M{
			// "addresses_eth_lst": bson.M{"$in": bson.A{p_user_address_eth_str, }},
			"id_str":       pUserIDstr,
			"deleted_bool": false,
		},
		bson.M{"$set": fieldsTargets})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to to update user info",
			"mongodb_update_error",
			map[string]interface{}{
				"user_name_str":   pUpdateOp.UserNameStr,
				"description_str": pUpdateOp.DescriptionStr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// GET_BY_ID

func dbUserGetByID(pUserIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuser, *gf_core.GFerror) {

	findOpts := options.FindOne()
	
	user := GFuser{}
	collNameStr := "gf_users"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx, bson.M{
			"id_str":       pUserIDstr,
			"deleted_bool": false,
		},
		findOpts).Decode(&user)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to find user by ID in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_id_str": pUserIDstr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	return &user, nil
}

//---------------------------------------------------
// GET_BY_USERNAME

func dbUserGetByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuser, *gf_core.GFerror) {

	findOpts := options.FindOne()
	
	user := GFuser{}
	collNameStr := "gf_users"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx, bson.M{
			"user_name_str": pUserNameStr,
			"deleted_bool":  false,
		},
		findOpts).Decode(&user)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to find user by user_name in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": pUserNameStr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	return &user, nil
}

//---------------------------------------------------
// GET_BY_ETH_ADDR

func dbUserGetByETHaddr(pUserAddressETHstr GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuser, *gf_core.GFerror) {

	findOpts := options.FindOne()
	
	user := GFuser{}
	collNameStr := "gf_users"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx, bson.M{
			"addresses_eth_lst": bson.M{"$in": bson.A{pUserAddressETHstr, }},
			"deleted_bool":      false,
		},
		findOpts).Decode(&user)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to find user by Eth address in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_address_eth_str": pUserAddressETHstr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	return &user, nil
}

//---------------------------------------------------
// EXISTS_BY_USERNAME

func DBuserExistsByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	collNameStr := "gf_users"
	countInt, gfErr := gf_core.MongoCount(bson.M{
			"user_name_str": pUserNameStr,
			"deleted_bool":  false,
		},
		map[string]interface{}{
			"user_name_str":  pUserNameStr,
			"caller_err_msg": "failed to check if there is a user in the DB with a given user_name",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------
// EXISTS_BY_ETH_ADDR

func dbUserExistsByETHaddr(pUserAddressETHstr GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	collNameStr := "gf_users"
	countInt, gfErr := gf_core.MongoCount(bson.M{
			"addresses_eth_lst": bson.M{"$in": bson.A{pUserAddressETHstr, }},
			"deleted_bool":      false,
		},
		map[string]interface{}{
			"user_address_eth_str": pUserAddressETHstr,
			"caller_err_msg":       "failed to check if there is a user in the DB with a given address",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------
// EMAIL_IS_CONFIRMED

// for initial user creation only, checks if the if the user confirmed their email.
// this is done only once.
func dbUserEmailIsConfirmed(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	
	findOpts := options.FindOne()
	findOpts.Projection = map[string]interface{}{
		"email_confirmed_bool": 1,
	}
	
	userMap := map[string]interface{}{}
	collNameStr := "gf_users"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		bson.M{
			"user_name_str": pUserNameStr,
			"deleted_bool":  false,
		},
		findOpts).Decode(&userMap)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get user email_confirmed from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": pUserNameStr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return false, gfErr
	}

	emailConfirmedBool := userMap["email_confirmed_bool"].(bool)
	return emailConfirmedBool, nil
}


//---------------------------------------------------

func DBgetUserNameByID(pUserIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (GFuserName, *gf_core.GFerror) {

	findOpts := options.FindOne()
	findOpts.Projection = map[string]interface{}{
		"user_name_str": 1,
	}
	
	userBasicInfoMap := map[string]interface{}{}
	collNameStr := "gf_users"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		bson.M{
			"id_str":       string(pUserIDstr),
			"deleted_bool": false,
		},
		findOpts).Decode(&userBasicInfoMap)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get user basic_info in the DB",
			"mongodb_find_error",
			map[string]interface{}{"user_id_str": pUserIDstr,},
			err, "gf_identity_core", pRuntimeSys)
		return GFuserName(""), gfErr
	}

	userNameStr := GFuserName(userBasicInfoMap["user_name_str"].(string))

	return userNameStr, nil
}

//---------------------------------------------------
// GET_BASIC_INFO_BY_ETH_ADDR

func DBgetBasicInfoByETHaddr(pUserAddressETHstr GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	userIDstr, gfErr := DBgetUserID(bson.M{
			"addresses_eth_lst": bson.M{"$in": bson.A{pUserAddressETHstr, }},
			"deleted_bool":      false,
		},
		map[string]interface{}{
			"user_address_eth_str": pUserAddressETHstr,
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gf_core.GF_ID(""), gfErr
	}

	return userIDstr, nil
}

//---------------------------------------------------
// GET_BASIC_INFO_BY_USERNAME

func DBgetBasicInfoByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	userIDstr, gfErr := DBgetUserID(bson.M{
			"user_name_str": pUserNameStr,
			"deleted_bool":  false,
		},
		// meta_map
		map[string]interface{}{
			"user_name_str": pUserNameStr,
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gf_core.GF_ID(""), gfErr
	}
	
	return userIDstr, nil
}

//---------------------------------------------------
// DB_GET_USER_ID

func DBgetUserID(pQuery bson.M,
	pMetaMap    map[string]interface{}, // data describing the DB write op
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {


	findOpts := options.FindOne()
	findOpts.Projection = map[string]interface{}{
		"id_str": 1,
	}
	
	userBasicInfoMap := map[string]interface{}{}
	collNameStr := "gf_users"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		pQuery,
		findOpts).Decode(&userBasicInfoMap)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get user basic_info in the DB",
			"mongodb_find_error",
			pMetaMap,
			err, "gf_identity_core", pRuntimeSys)
		return gf_core.GF_ID(""), gfErr
	}

	userIDstr := gf_core.GF_ID(userBasicInfoMap["id_str"].(string))

	return userIDstr, nil
}

//---------------------------------------------------
// INVITE_LIST
//---------------------------------------------------

func DBuserGetAllInInviteList(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	collNameStr := "gf_users_invite_list"

	findOpts := options.Find()
	cursor, gfErr := gf_core.MongoFind(bson.M{
			"deleted_bool":  false,
		},
		findOpts,
		map[string]interface{}{
			"caller_err_msg_str": "failed to get all records in invite_list from the DB",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return nil, gfErr
	}
	
	// no login_attempt found for user
	if cursor == nil {
		return nil, nil
	}
	
	inviteListLst := []map[string]interface{}{}
	err := cursor.All(pCtx, &inviteListLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get all records in invite_list from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	return inviteListLst, nil
}

//---------------------------------------------------
// ADD_TO_INVITE_LIST

func DBuserAddToInviteList(pUserEmailStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

	collNameStr   := "gf_users_invite_list"
	userInviteMap := map[string]interface{}{
		"user_email_str":       pUserEmailStr,
		"creation_unix_time_f": creationUNIXtimeF,
		"deleted_bool":         false,
	}
	gfErr := gf_core.MongoInsert(userInviteMap,
		collNameStr,
		map[string]interface{}{
			"user_email_str":     pUserEmailStr,
			"caller_err_msg_str": "failed to add a user to the invite_list in the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------

func DBuserRemoveFromInviteList(pUserEmailStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	collNameStr := "gf_users_invite_list"
	fieldsTargets := bson.M{
		"deleted_bool": true,
	}
	
	_, err := pRuntimeSys.Mongo_db.Collection(collNameStr).UpdateMany(pCtx, bson.M{
		"user_email_str": pUserEmailStr,
		"deleted_bool":   false,
	},
	bson.M{"$set": fieldsTargets})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to remove user from invite list",
			"mongodb_update_error",
			map[string]interface{}{
				"user_email_str": pUserEmailStr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return gfErr
	}
	return nil
}

//---------------------------------------------------
// CHECK_IN_INVITE_LIST_BY_EMAIL

func dbUserCheckInInvitelistByEmail(pUserEmailStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {
	
	collNameStr := "gf_users_invite_list"
	countInt, gfErr := gf_core.MongoCount(bson.M{
			"user_email_str": pUserEmailStr,
			"deleted_bool":   false,
		},
		map[string]interface{}{
			"user_email_str": pUserEmailStr,
			"caller_err_msg": "failed to check if the user_name is in the invite list",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------
// USER_CREDS
//---------------------------------------------------
// CREATE_CREDS

func dbUserCredsCreate(pUserCreds *GFuserCreds,
	pCtx         context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_users_creds"

	gfErr := gf_core.MongoInsert(pUserCreds,
		collNameStr,
		map[string]interface{}{
			"user_id_str":        pUserCreds.UserIDstr,
			"user_name_str":      pUserCreds.UserNameStr,
			"caller_err_msg_str": "failed to insert gf_user_creds into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------

func dbUserCredsGetPassHash(pUserNameStr GFuserName,
	pCtx         context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, string, *gf_core.GFerror) {

	collNameStr := "gf_users_creds"
	
	findOpts := options.FindOne()
	findOpts.Projection = map[string]interface{}{
		"pass_salt_str": 1,
		"pass_hash_str": 1,
	}
	
	userCredsInfoMap := map[string]interface{}{}
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx, bson.M{
			"user_name_str": string(pUserNameStr),
			"deleted_bool":  false,
		},
		findOpts).Decode(&userCredsInfoMap)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to find user creds by user_name in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": pUserNameStr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return "", "", gfErr
	}

	passSaltStr := userCredsInfoMap["pass_salt_str"].(string)
	passHashStr := userCredsInfoMap["pass_hash_str"].(string)

	return passSaltStr, passHashStr, nil
}

//---------------------------------------------------
// EMAIL
//---------------------------------------------------
// CREATE__EMAIL_CONFIRM

func dbUserEmailConfirmCreate(pUserNameStr GFuserName,
	pUserIDstr      gf_core.GF_ID,
	pConfirmCodeStr string,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr       := "gf_users_email_confirm"
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

	email_confirm_map := map[string]interface{}{
		"user_name_str":        pUserNameStr,
		"user_id_str":          pUserIDstr,
		"confirm_code_str":     pConfirmCodeStr,
		"creation_unix_time_f": creationUNIXtimeF,
	}

	gfErr := gf_core.MongoInsert(email_confirm_map,
		collNameStr,
		map[string]interface{}{
			"user_id_str":        pUserIDstr,
			"caller_err_msg_str": "failed to insert user email confirm_code into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------
// GET__EMAIL_CONFIRM_CODE

func dbUserEmailConfirmGetCode(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, float64, *gf_core.GFerror) {

	collNameStr := "gf_users_email_confirm"

	findOpts := options.FindOne()
	findOpts.SetSort(map[string]interface{}{"creation_unix_time_f": -1})
	findOpts.Projection = map[string]interface{}{
		"confirm_code_str":     1,
		"creation_unix_time_f": 1,
	}
	
	emailConfirmMap := map[string]interface{}{}
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		bson.M{
			"user_name_str": string(pUserNameStr),
		},
		findOpts).Decode(&emailConfirmMap)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get user email_confirm info from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": string(pUserNameStr),
			},
			err, "gf_identity_lib", pRuntimeSys)
		return "", 0.0, gfErr
	}

	confirmCodeStr    := emailConfirmMap["confirm_code_str"].(string)
	creationUNIXtimeF := emailConfirmMap["creation_unix_time_f"].(float64)

	return confirmCodeStr, creationUNIXtimeF, nil
}

//---------------------------------------------------

func dbUserGetEmailConfirmedByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	collNameStr := "gf_users"

	findOpts := options.FindOne()
	findOpts.Projection = map[string]interface{}{
		"email_confirmed_bool": 1,
	}
	
	userMap := map[string]interface{}{}
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		bson.M{
			"user_name_str": string(pUserNameStr),
		},
		findOpts).Decode(&userMap)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get user email_confirm status of a user from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": string(pUserNameStr),
			},
			err, "gf_identity_lib", pRuntimeSys)
		return false, gfErr
	}

	emailConfirmedBool := userMap["email_confirmed_bool"].(bool)
	
	return emailConfirmedBool, nil
}

//---------------------------------------------------
// LOGIN_ATTEMPT
//---------------------------------------------------

func dbLoginAttemptCreate(pLoginAttempt *GFloginAttempt,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_login_attempt"
	gfErr := gf_core.MongoInsert(pLoginAttempt,
		collNameStr,
		map[string]interface{}{
			"login_attempt_id_str": pLoginAttempt.IDstr,
			"user_type_str":        pLoginAttempt.UserTypeStr,
			"user_name_str":        pLoginAttempt.UserNameStr,
			"caller_err_msg_str":   "failed to insert GFloginAttempt into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbLoginAttemptGetByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFloginAttempt, *gf_core.GFerror) {

	collNameStr := "gf_login_attempt"

	findOpts := options.Find()
	cursor, gfErr := gf_core.MongoFind(bson.M{
			"user_name_str": string(pUserNameStr),
			"deleted_bool":  false,
		},
		findOpts,
		map[string]interface{}{
			"user_name_str":      string(pUserNameStr),
			"caller_err_msg_str": "failed to get login_attempt by user_name from the DB",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return nil, gfErr
	}
	
	// no login_attempt found for user
	if cursor == nil {
		return nil, nil
	}
	
	loginAttemptsLst := []*GFloginAttempt{}
	err := cursor.All(pCtx, &loginAttemptsLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get login_attempt from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{
				"user_name_str": string(pUserNameStr),
			},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	if len(loginAttemptsLst) > 0 {
		login_attempt := loginAttemptsLst[0]
		return login_attempt, nil
	}
	return nil, nil
}

//---------------------------------------------------

func DBloginAttemptUpdate(pLoginAttemptIDstr *gf_core.GF_ID,
	pUpdateOp   *GFloginAttemptUpdateOp,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	fieldsTargets := bson.M{}

	if pUpdateOp.PassConfirmedBool != nil {
		fieldsTargets["pass_confirmed_bool"] = *pUpdateOp.PassConfirmedBool
	}
	if pUpdateOp.EmailConfirmedBool != nil {
		fieldsTargets["email_confirmed_bool"] = *pUpdateOp.EmailConfirmedBool
	}
	if pUpdateOp.MFAconfirmedBool != nil {
		fieldsTargets["mfa_confirmed_bool"] = *pUpdateOp.MFAconfirmedBool
	}
	if pUpdateOp.DeletedBool != nil {
		fieldsTargets["deleted_bool"] = *pUpdateOp.DeletedBool
	}
	
	_, err := pRuntimeSys.Mongo_db.Collection("gf_login_attempt").UpdateMany(pCtx, bson.M{
		"id_str":       pLoginAttemptIDstr,
		"deleted_bool": false,
	},
	bson.M{"$set": fieldsTargets})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to to update a login_attempt",
			"mongodb_update_error",
			map[string]interface{}{
				"login_attempt_id_str": string(*pLoginAttemptIDstr),
			},
			err, "gf_identity_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}