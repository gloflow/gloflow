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
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// GET_BASIC_INFO
func db__user__get_basic_info_by_eth_addr(p_user_address_eth_str GF_user_address_eth,
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

// GET_BASIC_INFO
func db__user__get_basic_info_by_username(p_user_name_str GF_user_name,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (gf_core.GF_ID, *gf_core.GF_error) {

	user_id_str, gf_err := db__user__get_basic_info(bson.M{
			"user_name_str": p_user_name_str,
			"deleted_bool":  false,
		},
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
// GET_BY_ETH_ADDRESS
func db__user__get_by_eth_addr(p_user_address_eth_str GF_user_address_eth,
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

//---------------------------------------------------
// EXISTS_BY_ETH_ADDRESS
func db__user__exists_by_eth_addr(p_user_address_eth_str GF_user_address_eth,
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

//---------------------------------------------------
// CREATE
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
func db__user_creds__get_pass_hash(p_user_name_str GF_user_name,
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
// UPDATE
func db__user__update(p_user_address_eth_str GF_user_address_eth,
	p_update      *GF_user__update,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	//------------------------
	// FIELDS
	fields_targets := bson.M{}

	if string(p_update.User_name_str) != "" {
		fields_targets["username_str"] = p_update.User_name_str
	}

	if p_update.Description_str != "" {
		fields_targets["description_str"] = p_update.Description_str
	}
	
	//------------------------
	
	_, err := p_runtime_sys.Mongo_db.Collection("gf_users").UpdateMany(p_ctx, bson.M{
			"addresses_eth_lst": bson.M{"$in": bson.A{p_user_address_eth_str, }},
			"deleted_bool":      false,
		},
		bson.M{"$set": fields_targets})
		
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to to update user info",
			"mongodb_update_error",
			map[string]interface{}{
				"user_name_str":   p_update.User_name_str,
				"description_str": p_update.Description_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}

//---------------------------------------------------
// INVITE_LIST
func db__user__check_in_invitelist_by_username(p_user_email_str string,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {

	count_int, gf_err := gf_core.Mongo__count(bson.M{
		"user_email_str": p_user_email_str,
	},
	map[string]interface{}{
		"user_email_str": p_user_email_str,
		"caller_err_msg": "failed to check if the user_name is in the invite list",
	},
	p_runtime_sys.Mongo_db.Collection("gf_users_invite_list"),
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
// EXISTS
func db__user__exists_by_username(p_user_name_str GF_user_name,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {

	count_int, gf_err := gf_core.Mongo__count(bson.M{
			"user_name_str": p_user_name_str,
			"deleted_bool":  false,
		},
		map[string]interface{}{
			"user_name_str":  p_user_name_str,
			"caller_err_msg": "failed to check if there is a user in the DB with a given user_name",
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

//---------------------------------------------------
// ADD_TO_INVITE_LIST
func db__user__add_to_invite_list(p_user_email_str string,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	coll_name_str := "gf_users_invite_list"

	user_invite_map := map[string]interface{}{
		"user_email_str": p_user_email_str,
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