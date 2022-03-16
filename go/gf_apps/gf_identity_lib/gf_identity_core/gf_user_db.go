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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// GET_BASIC_INFO_BY_ETH_ADDR
func DBgetBasicInfoByETHaddr(pUserAddressETHstr GF_user_address_eth,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (gf_core.GF_ID, *gf_core.GF_error) {

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
	pRuntimeSys *gf_core.Runtime_sys) (gf_core.GF_ID, *gf_core.GF_error) {

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
func DBgetUserID(p_query bson.M,
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