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
	"fmt"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//-------------------------------------------------
type GFpolicy struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	IDstr             gf_core.GF_ID      `bson:"id_str"`
	DeletedBool       bool               `bson:"deleted_bool"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`
	
	// policy can be asssociated with multiple resources
	TargetResourceIDsLst  []gf_core.GF_ID `bson:"target_resource_ids_lst"`
	TargetResourceTypeStr string          `bson:"target_resource_type_str"`
	OwnerUserIDstr        gf_core.GF_ID   `bson:"owner_user_id_str"`

	// if the flow is fully public and all users (including anonymous) can view it
	PublicViewBool bool `bson:"public_view_bool"`

	//-----------------------
	// PRINCIPALS

	// viewers are users that can view if PublicViewBool is false
	ViewersUserIDsLst []gf_core.GF_ID `bson:"viewers_user_ids_lst"`

	// taggers are users that can attach tags/notes to flows
	TaggersUserIDsLst []gf_core.GF_ID `bson:"taggers_user_ids_lst"`

	// editors are users that are allowed by the owner to update/add/remove items to the flow
	EditorsUserIDsLst []gf_core.GF_ID `bson:"editors_user_ids_lst"`

	//-----------------------
}

//-------------------------------------------------
// VERIFY
func Verify(pRequestedOpStr string,
	pTargetResourceIDstr gf_core.GF_ID,
	pUserIDstr           gf_core.GF_ID,
	pCtx                 context.Context,
	pRuntimeSys          *gf_core.Runtime_sys) *gf_core.GF_error {

	// GET_POLICIES
	policiesLst, gfErr := DBgetPolicies(pTargetResourceIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	// GET_DEFS
	policiesDefsMap := getDefs()

	
	// VALIDATE_POLICIES
	for _, policy := range policiesLst {

		verifiedBool := policySingleVerify(pRequestedOpStr,
			policy,
			pUserIDstr,
			policiesDefsMap)
		if verifiedBool {

			// policy approved so dont raise an error
			return nil
		}
	}

	gfErr = gf_core.Mongo__handle_error("policy has failed to be validated",
		"policy__op_denied",
		map[string]interface{}{
			"target_resource_id_str": pTargetResourceIDstr,
		},
		nil, "gf_policy", pRuntimeSys)

	return gfErr
}

//-------------------------------------------------
// POLICY_SINGLE_VERIFY
func policySingleVerify(pRequestedOpStr string,
	pPolicy           *GFpolicy,
	pCurrentUserIDstr gf_core.GF_ID,
	pPoliciesDefsMap  map[string][]string) bool {

	// VIEWING
	// this is the lowest level set of permissions, so attempt to match that first
	for _, opStr := range pPoliciesDefsMap["viewing"] {

		if pRequestedOpStr == opStr {
			// for each allowed viwing user_id check if it equals to the
			// user_id requesting the operation permission.
			// IMPORTANT!! - view operations are lowest level,
			//               and any other policy can allow them.
			//               thats why not only viewer user_ids are checked
			//               but also tagging and editing ones.
			for _, allowedUserIDstr := range pPolicy.ViewersUserIDsLst {
				if pCurrentUserIDstr == allowedUserIDstr {
					return true
				}
			}
			for _, allowedUserIDstr := range pPolicy.TaggersUserIDsLst {
				if pCurrentUserIDstr == allowedUserIDstr {
					return true
				}
			}
			for _, allowedUserIDstr := range pPolicy.EditorsUserIDsLst {
				if pCurrentUserIDstr == allowedUserIDstr {
					return true
				}
			}
			return false
		}
	}

	// TAGGING
	for _, opStr := range pPoliciesDefsMap["tagging"] {

		if pRequestedOpStr == opStr {
			// for each allowed viwing user_id check if it equals to the
			// user_id requesting the operation permission.
			for _, allowedUserIDstr := range pPolicy.TaggersUserIDsLst {
				if pCurrentUserIDstr == allowedUserIDstr {
					return true
				}
			}
			for _, allowedUserIDstr := range pPolicy.EditorsUserIDsLst {
				if pCurrentUserIDstr == allowedUserIDstr {
					return true
				}
			}
			return false
		}
	}

	// EDITING
	// highest level permission set at the momement, so try last
	for _, opStr := range pPoliciesDefsMap["editing"] {

		if pRequestedOpStr == opStr {
			// for each allowed viwing user_id check if it equals to the
			// user_id requesting the operation permission.
			for _, allowedUserIDstr := range pPolicy.EditorsUserIDsLst {
				if pCurrentUserIDstr == allowedUserIDstr {
					return true
				}
			}
			return false
		}
	}

	return false
}

//-------------------------------------------------
func PipelineUpdate(pTargetResourceIDstr gf_core.GF_ID,
	pOwnerUserIDstr gf_core.GF_ID,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.Runtime_sys) *gf_core.GF_error {


	getDefs()

	return nil
}

//-------------------------------------------------
// PIPELINE__CREATE
func PipelineCreate(pTargetResourceIDstr gf_core.GF_ID,
	pOwnerUserIDstr gf_core.GF_ID,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.Runtime_sys) *gf_core.GF_error {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	IDstr             := createID(pTargetResourceIDstr, creationUNIXtimeF)

	policy := &GFpolicy{
		IDstr:                IDstr,     
		CreationUNIXtimeF:    creationUNIXtimeF,
		TargetResourceIDsLst: []gf_core.GF_ID{pTargetResourceIDstr, },
		OwnerUserIDstr:       pOwnerUserIDstr,
		PublicViewBool:       true,
	}

	// DB
	gfErr := DBcreatePolicy(policy, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return gfErr
}

//---------------------------------------------------
func createID(pTargetResourceIDstr gf_core.GF_ID,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		string(pTargetResourceIDstr),
		fmt.Sprintf("%f", pCreationUNIXtimeF),
	}
	gfIDstr := gf_core.ID__create(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfIDstr
}