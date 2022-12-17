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

package gf_policy

import (
	// "fmt"
	"context"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFpolicyUpdateOp struct {
	PublicViewBool *bool
}

//---------------------------------------------------
// GET_BY_ID

func DBgetPolicyByID(pPolicyIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFpolicy, *gf_core.GFerror) {

	collNameStr := "gf_policies"
	findOpts := options.FindOne()
	// findOpts.Projection = map[string]interface{}{}
	
	policy := &GFpolicy{}
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		bson.M{
			"id_str": pPolicyIDstr,
		},
		findOpts).Decode(policy)
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get policy by ID from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"policy_id_str": pPolicyIDstr,
			},
			err, "gf_policy", pRuntimeSys)
		return nil, gfErr
	}
	
	return policy, nil
}

//---------------------------------------------------
// GET

func DBgetPolicies(pTargetResourceIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFpolicy, *gf_core.GFerror) {

	collNameStr := "gf_policies"
	findOpts := options.FindOne()
	// findOpts.Projection = map[string]interface{}{}
	
	policiesLst := []*GFpolicy{}
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		bson.M{
			"target_resource_ids_lst": bson.M{"$in": []string{string(pTargetResourceIDstr)}},	
		},
		findOpts).Decode(&policiesLst)
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get policies by target_resource_id in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"target_resource_id_str": pTargetResourceIDstr,
			},
			err, "gf_policy", pRuntimeSys)
		return nil, gfErr
	}
	
	return policiesLst, nil
}

//---------------------------------------------------
// CREATE

func DBcreatePolicy(pPolicy *GFpolicy,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	return nil
}

//---------------------------------------------------
// UPDATE

func DBupdatePolicy(pPolicyIDstr gf_core.GF_ID,
	pUpdateOp   *GFpolicyUpdateOp,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	fieldsTargets := bson.M{}

	if pUpdateOp.PublicViewBool != nil {
		fieldsTargets["public_view_bool"] = *pUpdateOp.PublicViewBool
	}


	_, err := pRuntimeSys.Mongo_db.Collection("gf_policies").UpdateMany(pCtx, bson.M{
		"id_str":       pPolicyIDstr,
		"deleted_bool": false,
	},
	bson.M{"$set": fieldsTargets})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to to update policy in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"policy_id_str": string(pPolicyIDstr),
			},
			err, "gf_policy", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// EXISTS_BY_USERNAME

func DBexistsByID(pPolicyIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	collNameStr := "gf_policies"

	countInt, gfErr := gf_core.MongoCount(bson.M{
			"id_str":       pPolicyIDstr,
			"deleted_bool": false,
		},
		map[string]interface{}{
			"policy_id_str":  pPolicyIDstr,
			"caller_err_msg": "failed to check if there is a policy in the DB with a given ID",
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