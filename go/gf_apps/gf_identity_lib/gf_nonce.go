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
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GF_user_nonce_val string
type GF_user_nonce struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`

	User_id_str     gf_core.GF_ID       `bson:"user_id_str"`
	Address_eth_str GF_user_address_eth `bson:"address_eth_str"`
	Val_str         GF_user_nonce_val   `bson:"val_str"`
}

//---------------------------------------------------
func nonce__create_and_persist(p_user_id_str gf_core.GF_ID,
	p_user_address_eth_str GF_user_address_eth,
	p_ctx                  context.Context,
	p_runtime_sys          *gf_core.Runtime_sys) (*GF_user_nonce, *gf_core.GF_error) {

	//------------------------
	// mark all existing nonces (if there are any) for this user_address_eth
	// as deleted
	gf_err := db__nonce__delete_all(p_user_address_eth_str, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------------
	nonce_val_str := gf_core.Str_random()

	// CREATE
	nonce, gf_err := nonce__create(GF_user_nonce_val(nonce_val_str),
		p_user_id_str,
		p_user_address_eth_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return nonce, nil
}

//---------------------------------------------------
func nonce__create(p_nonce_val_str GF_user_nonce_val,
	p_user_id_str          gf_core.GF_ID,
	p_user_address_eth_str GF_user_address_eth,
	p_ctx                  context.Context,
	p_runtime_sys          *gf_core.Runtime_sys) (*GF_user_nonce, *gf_core.GF_error) {

	creation_unix_time_f   := float64(time.Now().UnixNano())/1000000000.0
	unique_vals_for_id_lst := []string{string(p_nonce_val_str), }

	id_str := gf_core.ID__create(unique_vals_for_id_lst, creation_unix_time_f)
	
	nonce := &GF_user_nonce{
		V_str:                "0",
		Id_str:               id_str,
		Creation_unix_time_f: creation_unix_time_f,
		User_id_str:          p_user_id_str,
		Address_eth_str:      p_user_address_eth_str,
		Val_str:              p_nonce_val_str,
	}

	// DB
	gf_err := db__nonce__create(nonce, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return nonce, nil
}

//---------------------------------------------------
func db__nonce__delete_all(p_user_address_eth_str GF_user_address_eth,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	_, err := p_runtime_sys.Mongo_db.Collection("gf_users_nonces").UpdateMany(p_ctx, bson.M{
			"address_eth_str": p_user_address_eth_str,
			"deleted_bool":    false,
		},
		bson.M{"$set": bson.M{
			"deleted_bool": true,
		}})
		
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to mark all nonces for a user_address_eth as deleted",
			"mongodb_update_error",
			map[string]interface{}{
				"user_address_eth": p_user_address_eth_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}

//---------------------------------------------------
func db__nonce__create(p_nonce *GF_user_nonce,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	coll_name_str := "gf_users_nonces"
	gf_err := gf_core.Mongo__insert(p_nonce,
		coll_name_str,
		map[string]interface{}{
			"user_id_str":        p_nonce.User_id_str,
			"address_eth_str":    p_nonce.Address_eth_str,
			"caller_err_msg_str": "failed to insert GF_user_nonce into the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	return nil
}

//---------------------------------------------------
func db__nonce__get(p_user_address_eth_str GF_user_address_eth,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (GF_user_nonce_val, *gf_core.GF_error) {

	user_nonce := &GF_user_nonce{}
	err := p_runtime_sys.Mongo_db.Collection("gf_users_nonces").FindOne(p_ctx, bson.M{
			"address_eth_str": p_user_address_eth_str,
			"deleted_bool":    false,
		}).Decode(&user_nonce)
		
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to find user by address in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_address_eth_str": p_user_address_eth_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return GF_user_nonce_val(""), gf_err
	}

	user_nonce_val_str := user_nonce.Val_str
	
	return user_nonce_val_str, nil
}