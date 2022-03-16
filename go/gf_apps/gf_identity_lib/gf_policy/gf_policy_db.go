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
// GET
func DBgetPolicies(pTargetResourceIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) ([]*GFpolicy, *gf_core.GF_error) {

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
		gfErr := gf_core.Mongo__handle_error("failed to get policy basic_info in the DB",
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
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GF_error {

	return nil
}