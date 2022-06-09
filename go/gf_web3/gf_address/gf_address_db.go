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

package gf_address

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// GET_ALL
func DBgetAll(pAddressTypeStr string,
	pAddressChainNameStr string,
	pUserIDstr           gf_core.GF_ID,
	pCtx                 context.Context,
	pRuntimeSys          *gf_core.Runtime_sys) ([]*GFchainAddress, *gf_core.GFerror) {

	collNameStr := "gf_web3_addresses"

	findOpts := options.Find()
	cursor, gfErr := gf_core.MongoFind(bson.M{
			"owner_user_id_str": pUserIDstr,
			"type_str":          pAddressTypeStr,
			"chain_name_str":    pAddressChainNameStr,
			"deleted_bool":      false,
		},
		findOpts,
		map[string]interface{}{
			"owner_user_id_str":  pUserIDstr,
			"caller_err_msg_str": "failed to get home_viz record from the DB",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return nil, gfErr
	}

	var addressesLst []*GFchainAddress
	err := cursor.All(pCtx, &addressesLst)
	if err != nil {
		gfErr := gf_core.Mongo__handle_error("failed to get a home_viz record from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_address", pRuntimeSys)
		return nil, gfErr
	}

	return addressesLst, nil
}

//-------------------------------------------------
func DBadd(pAddress *GFchainAddress,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_web3_addresses"
	
	gfErr := gf_core.MongoInsert(pAddress,
		collNameStr,
		map[string]interface{}{
			"owner_user_id_str":  pAddress.OwnerUserIDstr,
			"caller_err_msg_str": "failed to insert GFhomeViz into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//-------------------------------------------------
// EXISTS
func DBexists(pAddressStr string,
	pChainStr   string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	collNameStr := "gf_web3_addresses"

	countInt, gfErr := gf_core.Mongo__count(bson.M{
			"address_str":    pAddressStr,
			"chain_name_str": pChainStr,
			"deleted_bool":   false,
		},
		map[string]interface{}{
			"address_str":    pAddressStr,
			"chain_name_str": pChainStr,
			"caller_err_msg": "failed to check if there is an address in the DB",
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

//-------------------------------------------------
// ADD_TAG
func DBaddTag(pTagsLst []string,
	pAddressStr string,
	pChainStr   string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	_, err := pRuntimeSys.Mongo_db.Collection("gf_home_viz").UpdateMany(pCtx, bson.M{
			"address_str":    pAddressStr,
			"chain_name_str": pChainStr,
			"deleted_bool":   false,
		},

		// push all tags into the "tags_lst" list.
		// "$each" - needed to push each element of pTagsLst as a separate element
		bson.M{"$push": bson.M{
			"tags_lst": bson.M{"$each": pTagsLst},
		}})

	if err != nil {
		gfErr := gf_core.Mongo__handle_error("failed to update adress with new tags",
			"mongodb_update_error",
			map[string]interface{}{
				"tags_lst":    pTagsLst,
				"address_str": pAddressStr,
				"chain_str":   pChainStr,
			},
			err, "gf_address", pRuntimeSys)
		return gfErr
	}

	

	return nil
}