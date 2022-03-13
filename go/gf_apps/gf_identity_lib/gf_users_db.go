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
	// "fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//---------------------------------------------------
type GF_user__update_op struct {
	User_name_str        gf_identity_core.GFuserName
	Description_str      string
	Email_str            string
	Email_confirmed_bool bool
	MFA_confirm_bool     *bool // if nil dont update, else update to true/false
}

type GF_login_attempt__update_op struct {
	Pass_confirmed_bool  *bool
	Email_confirmed_bool *bool
	MFA_confirmed_bool   *bool
	Deleted_bool         *bool
}

//---------------------------------------------------
// USER
//---------------------------------------------------
// CREATE_USER
func db__user__create(p_user *GF_user,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	coll_name_str := "gf_users"

	gf_err := gf_core.Mongo__insert(p_user,
		coll_name_str,
		map[string]interface{}{
			"user_id_str":        p_user.Id_str,
			"user_name_str":      p_user.User_name_str,
			"description_str":    p_user.Description_str,
			"addresses_eth_lst":  p_user.Addresses_eth_lst, 
			"caller_err_msg_str": "failed to insert GF_user into the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	return nil
}

// UPDATE
func db__user__update(p_user_id_str gf_core.GF_ID, // p_user_address_eth_str GF_user_address_eth,
	p_update_op   *GF_user__update_op,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	//------------------------
	// FIELDS
	fields_targets := bson.M{}

	if string(p_update_op.User_name_str) != "" {
		fields_targets["username_str"] = p_update_op.User_name_str
	}

	if p_update_op.Description_str != "" {
		fields_targets["description_str"] = p_update_op.Description_str
	}

	if p_update_op.Email_str != "" {
		fields_targets["email_str"] = p_update_op.Email_str

		// IMPORTANT!! - if the email is changed then it needs to be confirmed
		//               again. this flag on the user can only be changed to false
		//               indirectly like here by upding the user email
		fields_targets["email_confirmed_bool"] = false
	}

	// email_confirmed_bool itself can only be updated to explicitly
	// if its being set to true
	if p_update_op.Email_confirmed_bool {
		fields_targets["email_confirmed_bool"] = true
	}
	
	// MFA_confirm_bool - is a pointer. if its not nil, then set
	//                    the MFA_confirm field to either true/false.
	if p_update_op.MFA_confirm_bool != nil {
		fields_targets["mfa_confirm_bool"] = p_update_op.MFA_confirm_bool
	}

	//------------------------
	
	_, err := p_runtime_sys.Mongo_db.Collection("gf_users").UpdateMany(p_ctx, bson.M{
			// "addresses_eth_lst": bson.M{"$in": bson.A{p_user_address_eth_str, }},
			"id_str":       p_user_id_str,
			"deleted_bool": false,
		},
		bson.M{"$set": fields_targets})
		
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to to update user info",
			"mongodb_update_error",
			map[string]interface{}{
				"user_name_str":   p_update_op.User_name_str,
				"description_str": p_update_op.Description_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}

// GET_BY_ID
func db__user__get_by_id(p_user_id_str gf_core.GF_ID,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user, *gf_core.GF_error) {

	find_opts := options.FindOne()
	
	user := GF_user{}
	err := p_runtime_sys.Mongo_db.Collection("gf_users").FindOne(p_ctx, bson.M{
			"id_str":       p_user_id_str,
			"deleted_bool": false,
		},
		find_opts).Decode(&user)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to find user by ID in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_id_str": p_user_id_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return nil, gf_err
	}

	return &user, nil
}

// GET_BY_USERNAME
func dbUserGetByUsername(pUserNameStr gf_identity_core.GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (*GF_user, *gf_core.GF_error) {

	find_opts := options.FindOne()
	
	user := GF_user{}
	err := pRuntimeSys.Mongo_db.Collection("gf_users").FindOne(pCtx, bson.M{
			"user_name_str": pUserNameStr,
			"deleted_bool":  false,
		},
		find_opts).Decode(&user)

	if err != nil {
		gfErr := gf_core.Mongo__handle_error("failed to find user by user_name in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": pUserNameStr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return nil, gfErr
	}

	return &user, nil
}

// GET_BASIC_INFO_BY_ETH_ADDR
func db__user__get_basic_info_by_eth_addr(p_user_address_eth_str gf_identity_core.GF_user_address_eth,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (gf_core.GF_ID, *gf_core.GF_error) {

	user_id_str, gf_err := db__user__get_basic_info(bson.M{
			"addresses_eth_lst": bson.M{"$in": bson.A{p_user_address_eth_str, }},
			"deleted_bool":      false,
		},
		map[string]interface{}{
			"user_address_eth_str": p_user_address_eth_str,
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_core.GF_ID(""), gf_err
	}

	return user_id_str, nil
}

// GET_BASIC_INFO_BY_USERNAME
func db__user__get_basic_info_by_username(p_user_name_str gf_identity_core.GFuserName,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (gf_core.GF_ID, *gf_core.GF_error) {

	user_id_str, gf_err := db__user__get_basic_info(bson.M{
			"user_name_str": p_user_name_str,
			"deleted_bool":  false,
		},
		// meta_map
		map[string]interface{}{
			"user_name_str": p_user_name_str,
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_core.GF_ID(""), gf_err
	}
	
	return user_id_str, nil
}

// GET_BASIC_INFO
func db__user__get_basic_info(p_query bson.M,
	p_meta_map    map[string]interface{}, // data describing the DB write op
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (gf_core.GF_ID, *gf_core.GF_error) {


	find_opts := options.FindOne()
	find_opts.Projection = map[string]interface{}{
		"id_str": 1,
	}
	
	user_basic_info_map := map[string]interface{}{}
	err := p_runtime_sys.Mongo_db.Collection("gf_users").FindOne(p_ctx,
		p_query,
		find_opts).Decode(&user_basic_info_map)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get user basic_info in the DB",
			"mongodb_find_error",
			p_meta_map,
			err, "gf_identity_lib", p_runtime_sys)
		return gf_core.GF_ID(""), gf_err
	}

	user_id_str := gf_core.GF_ID(user_basic_info_map["id_str"].(string))

	return user_id_str, nil
}

//---------------------------------------------------
// GET_BY_ETH_ADDR
func db__user__get_by_eth_addr(p_user_address_eth_str gf_identity_core.GF_user_address_eth,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_user, *gf_core.GF_error) {

	find_opts := options.FindOne()
	
	user := GF_user{}
	err := p_runtime_sys.Mongo_db.Collection("gf_users").FindOne(p_ctx, bson.M{
			"addresses_eth_lst": bson.M{"$in": bson.A{p_user_address_eth_str, }},
			"deleted_bool":      false,
		},
		find_opts).Decode(&user)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to find user by Eth address in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_address_eth_str": p_user_address_eth_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return nil, gf_err
	}

	return &user, nil
}

//---------------------------------------------------
// EXISTS_BY_USERNAME
func db__user__exists_by_username(p_user_name_str gf_identity_core.GFuserName,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {

	coll_name_str := "gf_users"

	count_int, gf_err := gf_core.Mongo__count(bson.M{
			"user_name_str": p_user_name_str,
			"deleted_bool":  false,
		},
		map[string]interface{}{
			"user_name_str":  p_user_name_str,
			"caller_err_msg": "failed to check if there is a user in the DB with a given user_name",
		},
		p_runtime_sys.Mongo_db.Collection(coll_name_str),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return false, gf_err
	}

	if count_int > 0 {
		return true, nil
	}
	return false, nil
}

// EXISTS_BY_ETH_ADDR
func db__user__exists_by_eth_addr(p_user_address_eth_str gf_identity_core.GF_user_address_eth,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {

	count_int, gf_err := gf_core.Mongo__count(bson.M{
			"addresses_eth_lst": bson.M{"$in": bson.A{p_user_address_eth_str, }},
			"deleted_bool":      false,
		},
		map[string]interface{}{
			"user_address_eth_str": p_user_address_eth_str,
			"caller_err_msg":       "failed to check if there is a user in the DB with a given address",
		},
		p_runtime_sys.Mongo_db.Collection("gf_users"),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return false, gf_err
	}

	if count_int > 0 {
		return true, nil
	}
	return false, nil
}

// EMAIL_IS_CONFIRMED
// for initial user creation only, checks if the if the user confirmed their email.
// this is done only once.
func db__user__email_is_confirmed(p_user_name_str gf_identity_core.GFuserName,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {

	find_opts := options.FindOne()
	find_opts.Projection = map[string]interface{}{
		"email_confirmed_bool": 1,
	}
	
	user_map := map[string]interface{}{}
	err := p_runtime_sys.Mongo_db.Collection("gf_users").FindOne(p_ctx,
		bson.M{
			"user_name_str": p_user_name_str,
			"deleted_bool":  false,
		},
		find_opts).Decode(&user_map)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get user email_confirmed from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": p_user_name_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return false, gf_err
	}

	email_confirmed_bool := user_map["email_confirmed_bool"].(bool)
	return email_confirmed_bool, nil
}

//---------------------------------------------------
// INVITE_LIST
//---------------------------------------------------
func db__user__get_all_in_invite_list(p_ctx context.Context,
	p_runtime_sys *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.GF_error) {

	coll_name_str := "gf_users_invite_list"

	find_opts := options.Find()
	cursor, gf_err := gf_core.Mongo__find(bson.M{
			"deleted_bool":  false,
		},
		find_opts,
		map[string]interface{}{
			"caller_err_msg_str": "failed to get all records in invite_list from the DB",
		},
		p_runtime_sys.Mongo_db.Collection(coll_name_str),
		p_ctx,
		p_runtime_sys)
	
	if gf_err != nil {
		return nil, gf_err
	}
	
	// no login_attempt found for user
	if cursor == nil {
		return nil, nil
	}
	
	invite_list_lst := []map[string]interface{}{}
	err := cursor.All(p_ctx, &invite_list_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get all records in invite_list from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_identity_lib", p_runtime_sys)
		return nil, gf_err
	}

	return invite_list_lst, nil
}

//---------------------------------------------------
// ADD_TO_INVITE_LIST
func db__user__add_to_invite_list(p_user_email_str string,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	coll_name_str   := "gf_users_invite_list"
	user_invite_map := map[string]interface{}{
		"user_email_str":       p_user_email_str,
		"creation_unix_time_f": creation_unix_time_f,
		"deleted_bool":         false,
	}
	gf_err := gf_core.Mongo__insert(user_invite_map,
		coll_name_str,
		map[string]interface{}{
			"user_email_str":     p_user_email_str,
			"caller_err_msg_str": "failed to add a user email to the invite_list in the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	return nil
}

//---------------------------------------------------
// CHECK_IN_INVITE_LIST_BY_USERNAME
func db__user__check_in_invitelist_by_username(p_user_email_str string,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {
	
	coll_name_str := "gf_users_invite_list"

	count_int, gf_err := gf_core.Mongo__count(bson.M{
		"user_email_str": p_user_email_str,
	},
	map[string]interface{}{
		"user_email_str": p_user_email_str,
		"caller_err_msg": "failed to check if the user_name is in the invite list",
	},
	p_runtime_sys.Mongo_db.Collection(coll_name_str),
	p_ctx,
	p_runtime_sys)
	if gf_err != nil {
		return false, gf_err
	}

	if count_int > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------
// USER_CREDS
//---------------------------------------------------
// CREATE_CREDS
func db__user_creds__create(p_user_creds *GF_user_creds,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	coll_name_str := "gf_users_creds"

	gf_err := gf_core.Mongo__insert(p_user_creds,
		coll_name_str,
		map[string]interface{}{
			"user_id_str":        p_user_creds.User_id_str,
			"user_name_str":      p_user_creds.User_name_str,
			"caller_err_msg_str": "failed to insert GF_user_creds into the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	return nil
}

//---------------------------------------------------
func db__user_creds__get_pass_hash(p_user_name_str gf_identity_core.GFuserName,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (string, string, *gf_core.GF_error) {

	coll_name_str := "gf_users_creds"
	
	find_opts := options.FindOne()
	find_opts.Projection = map[string]interface{}{
		"pass_salt_str": 1,
		"pass_hash_str": 1,
	}
	
	user_creds_info_map := map[string]interface{}{}
	err := p_runtime_sys.Mongo_db.Collection(coll_name_str).FindOne(p_ctx, bson.M{
			"user_name_str": string(p_user_name_str),
			"deleted_bool":  false,
		},
		find_opts).Decode(&user_creds_info_map)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to find user creds by user_name in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": p_user_name_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return "", "", gf_err
	}

	pass_salt_str := user_creds_info_map["pass_salt_str"].(string)
	pass_hash_str := user_creds_info_map["pass_hash_str"].(string)

	return pass_salt_str, pass_hash_str, nil
}

//---------------------------------------------------
// EMAIL
//---------------------------------------------------
// CREATE__EMAIL_CONFIRM
func db__user_email_confirm__create(p_user_name_str gf_identity_core.GFuserName,
	p_user_id_str      gf_core.GF_ID,
	p_confirm_code_str string,
	p_ctx              context.Context,
	p_runtime_sys      *gf_core.Runtime_sys) *gf_core.GF_error {

	coll_name_str        := "gf_users_email_confirm"
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	email_confirm_map := map[string]interface{}{
		"user_name_str":        p_user_name_str,
		"user_id_str":          p_user_id_str,
		"confirm_code_str":     p_confirm_code_str,
		"creation_unix_time_f": creation_unix_time_f,
	}

	gf_err := gf_core.Mongo__insert(email_confirm_map,
		coll_name_str,
		map[string]interface{}{
			"user_id_str":        p_user_id_str,
			"caller_err_msg_str": "failed to insert user email confirm_code into the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	return nil
}

//---------------------------------------------------
// GET__EMAIL_CONFIRM_CODE
func db__user_email_confirm__get_code(p_user_name_str gf_identity_core.GFuserName,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (string, float64, *gf_core.GF_error) {

	coll_name_str := "gf_users_email_confirm"

	find_opts := options.FindOne()
	find_opts.SetSort(map[string]interface{}{"creation_unix_time_f": -1})
	find_opts.Projection = map[string]interface{}{
		"confirm_code_str":     1,
		"creation_unix_time_f": 1,
	}
	
	email_confirm_map := map[string]interface{}{}
	err := p_runtime_sys.Mongo_db.Collection(coll_name_str).FindOne(p_ctx,
		bson.M{
			"user_name_str": string(p_user_name_str),
		},
		find_opts).Decode(&email_confirm_map)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get user email_confirm info from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": string(p_user_name_str),
			},
			err, "gf_identity_lib", p_runtime_sys)
		return "", 0.0, gf_err
	}

	confirm_code_str     := email_confirm_map["confirm_code_str"].(string)
	creation_unix_time_f := email_confirm_map["creation_unix_time_f"].(float64)

	return confirm_code_str, creation_unix_time_f, nil
}

//---------------------------------------------------
func db__user__get_email_confirmed_by_username(p_user_name_str gf_identity_core.GFuserName,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {


	coll_name_str := "gf_users"

	find_opts := options.FindOne()
	find_opts.Projection = map[string]interface{}{
		"email_confirmed_bool": 1,
	}
	
	user_map := map[string]interface{}{}
	err := p_runtime_sys.Mongo_db.Collection(coll_name_str).FindOne(p_ctx,
		bson.M{
			"user_name_str": string(p_user_name_str),
		},
		find_opts).Decode(&user_map)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get user email_confirm status of a user from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_name_str": string(p_user_name_str),
			},
			err, "gf_identity_lib", p_runtime_sys)
		return false, gf_err
	}

	email_confirmed_bool := user_map["email_confirmed_bool"].(bool)
	
	return email_confirmed_bool, nil
}

//---------------------------------------------------
// LOGIN_ATTEMPT
//---------------------------------------------------
func db__login_attempt__create(p_login_attempt *GF_login_attempt,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {


	coll_name_str := "gf_login_attempt"
	gf_err := gf_core.Mongo__insert(p_login_attempt,
		coll_name_str,
		map[string]interface{}{
			"login_attempt_id_str": p_login_attempt.Id_str,
			"user_type_str":        p_login_attempt.User_type_str,
			"user_name_str":        p_login_attempt.User_name_str,
			"caller_err_msg_str":   "failed to insert GF_login_attempt into the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}

//---------------------------------------------------
func db__login_attempt__get_by_username(p_user_name_str gf_identity_core.GFuserName,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_login_attempt, *gf_core.GF_error) {

	coll_name_str := "gf_login_attempt"

	find_opts := options.Find()
	cursor, gf_err := gf_core.Mongo__find(bson.M{
			"user_name_str": string(p_user_name_str),
			"deleted_bool":  false,
		},
		find_opts,
		map[string]interface{}{
			"user_name_str":      string(p_user_name_str),
			"caller_err_msg_str": "failed to get login_attempt by user_name from the DB",
		},
		p_runtime_sys.Mongo_db.Collection(coll_name_str),
		p_ctx,
		p_runtime_sys)
	
	if gf_err != nil {
		return nil, gf_err
	}
	
	// no login_attempt found for user
	if cursor == nil {
		return nil, nil
	}
	
	login_attempts_lst := []*GF_login_attempt{}
	err := cursor.All(p_ctx, &login_attempts_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get login_attempt from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{
				"user_name_str": string(p_user_name_str),
			},
			err, "gf_identity_lib", p_runtime_sys)
		return nil, gf_err
	}

	if len(login_attempts_lst) > 0 {
		login_attempt := login_attempts_lst[0]
		return login_attempt, nil
	}
	return nil, nil
}

//---------------------------------------------------
func db__login_attempt__update(p_login_attempt_id_str *gf_core.GF_ID,
	p_update_op   *GF_login_attempt__update_op,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {



	
	fields_targets := bson.M{}

	if p_update_op.Pass_confirmed_bool != nil {
		fields_targets["pass_confirmed_bool"] = *p_update_op.Pass_confirmed_bool
	}
	if p_update_op.Email_confirmed_bool != nil {
		fields_targets["email_confirmed_bool"] = *p_update_op.Email_confirmed_bool
	}
	if p_update_op.MFA_confirmed_bool != nil {
		fields_targets["mfa_confirmed_bool"] = *p_update_op.MFA_confirmed_bool
	}
	if p_update_op.Deleted_bool != nil {
		fields_targets["deleted_bool"] = *p_update_op.Deleted_bool
	}
	



	_, err := p_runtime_sys.Mongo_db.Collection("gf_login_attempt").UpdateMany(p_ctx, bson.M{
		"id_str":       p_login_attempt_id_str,
		"deleted_bool": false,
	},
	bson.M{"$set": fields_targets})
		
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to to update a login_attempt",
			"mongodb_update_error",
			map[string]interface{}{
				"login_attempt_id_str": string(*p_login_attempt_id_str),
			},
			err, "gf_identity_lib", p_runtime_sys)
		return gf_err
	}



	return nil
}