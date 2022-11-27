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
	"fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//---------------------------------------------------

type GFuserNonceVal string
type GFuserNonce struct {
	Vstr              string             `bson:"v_str"` // schema_version
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	IDstr             gf_core.GF_ID      `bson:"id_str"`
	DeletedBool       bool               `bson:"deleted_bool"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`

	UserIDstr     gf_core.GF_ID                     `bson:"user_id_str"`
	AddressETHstr gf_identity_core.GFuserAddressETH `bson:"address_eth_str"`
	ValStr        GFuserNonceVal                    `bson:"val_str"`
}

//---------------------------------------------------

func nonceCreateAndPersist(pUserIDstr gf_core.GF_ID,
	pUserAddressETHstr gf_identity_core.GFuserAddressETH,
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) (*GFuserNonce, *gf_core.GFerror) {

	//------------------------
	// mark all existing nonces (if there are any) for this user_address_eth
	// as deleted
	gfErr := dbNonceDeleteAll(pUserAddressETHstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	nonceValStr := fmt.Sprintf("gloflow:%s", gf_core.StrRandom())

	// CREATE
	nonce, gfErr := nonceCreate(GFuserNonceVal(nonceValStr),
		pUserIDstr,
		pUserAddressETHstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return nonce, nil
}

//---------------------------------------------------

func nonceCreate(pNonceValStr GFuserNonceVal,
	pUserIDstr         gf_core.GF_ID,
	pUserAddressETHstr gf_identity_core.GFuserAddressETH,
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) (*GFuserNonce, *gf_core.GFerror) {

	creationUNIXtimeF  := float64(time.Now().UnixNano())/1000000000.0
	uniqueValsForIDlst := []string{string(pNonceValStr), }

	idStr := gf_core.IDcreate(uniqueValsForIDlst, creationUNIXtimeF)
	
	nonce := &GFuserNonce{
		Vstr:              "0",
		IDstr:             idStr,
		CreationUNIXtimeF: creationUNIXtimeF,
		UserIDstr:         pUserIDstr,
		AddressETHstr:     pUserAddressETHstr,
		ValStr:            pNonceValStr,
	}

	// DB
	gfErr := dbNonceCreate(nonce, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return nonce, nil
}

//---------------------------------------------------

func dbNonceDeleteAll(pUserAddressETHstr gf_identity_core.GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	_, err := pRuntimeSys.Mongo_db.Collection("gf_users_nonces").UpdateMany(pCtx, bson.M{
			"address_eth_str": pUserAddressETHstr,
			"deleted_bool":    false,
		},
		bson.M{"$set": bson.M{
			"deleted_bool": true,
		}})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to mark all nonces for a user_address_eth as deleted",
			"mongodb_update_error",
			map[string]interface{}{
				"user_address_eth": pUserAddressETHstr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbNonceCreate(pNonce *GFuserNonce,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_users_nonces"
	gfErr := gf_core.MongoInsert(pNonce,
		collNameStr,
		map[string]interface{}{
			"user_id_str":        pNonce.UserIDstr,
			"address_eth_str":    pNonce.AddressETHstr,
			"caller_err_msg_str": "failed to insert GF_user_nonce into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------

func dbNonceGet(pUserAddressETHstr gf_identity_core.GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (GFuserNonceVal, bool, *gf_core.GFerror) {
	
	userNonce := &GFuserNonce{}
	err := pRuntimeSys.Mongo_db.Collection("gf_users_nonces").FindOne(pCtx, bson.M{
			"address_eth_str": pUserAddressETHstr,
			"deleted_bool":    false,
		}).Decode(&userNonce)
		
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return GFuserNonceVal(""), false, nil
		}

		gfErr := gf_core.MongoHandleError("failed to find user by address in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"user_address_eth_str": pUserAddressETHstr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return GFuserNonceVal(""), false, gfErr
	}

	userNonceValStr := userNonce.ValStr
	
	return userNonceValStr, true, nil
}