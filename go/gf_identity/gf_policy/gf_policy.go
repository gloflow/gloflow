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
	// "github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//-------------------------------------------------

type GFpolicy struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	ID                gf_core.GF_ID      `bson:"id_str"`
	DeletedBool       bool               `bson:"deleted_bool"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`
	
	// policy can be asssociated with multiple resources
	TargetResourceIDsLst  []gf_core.GF_ID `bson:"target_resource_ids_lst"`
	TargetResourceTypeStr string          `bson:"target_resource_type_str"`
	OwnerUserID           gf_core.GF_ID   `bson:"owner_user_id_str"`

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

type GFpolicyUpdateOutput struct {
	PolicyExistsBool bool
}

//-------------------------------------------------
// VERIFY

func Verify(pRequestedOpStr string,
	pTargetResourceID gf_core.GF_ID,
	pUserID           gf_core.GF_ID,
	pCtx              context.Context,
	pRuntimeSys       *gf_core.RuntimeSys) *gf_core.GFerror {

	// GET_POLICIES
	policiesLst, gfErr := DBsqlGetPolicies(pTargetResourceID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	// GET_DEFS
	policiesDefsMap := getDefs()

	
	// VALIDATE_POLICIES
	for _, policy := range policiesLst {

		verifiedBool := policySingleVerify(pRequestedOpStr,
			policy,
			pUserID,
			policiesDefsMap)
		if verifiedBool {

			// policy approved so dont raise an error
			return nil
		}
	}

	gfErr = gf_core.ErrorCreate("policy has failed to be validated",
		"policy__op_denied",
		map[string]interface{}{
			"user_id":            pUserID,
			"target_resource_id": pTargetResourceID,
		},
		nil, "gf_policy", pRuntimeSys)

	return gfErr
}

//-------------------------------------------------
// POLICY_SINGLE_VERIFY

func policySingleVerify(pRequestedOpStr string,
	pPolicy          *GFpolicy,
	pUserID          gf_core.GF_ID,
	pPoliciesDefsMap map[string][]string) bool {

	// if its the owner of the policy all operations are permitted
	if pUserID == pPolicy.OwnerUserID {
		return true
	}
	
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
				if pUserID == allowedUserIDstr {
					return true
				}
			}
			for _, allowedUserIDstr := range pPolicy.TaggersUserIDsLst {
				if pUserID == allowedUserIDstr {
					return true
				}
			}
			for _, allowedUserIDstr := range pPolicy.EditorsUserIDsLst {
				if pUserID == allowedUserIDstr {
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
				if pUserID == allowedUserIDstr {
					return true
				}
			}
			for _, allowedUserIDstr := range pPolicy.EditorsUserIDsLst {
				if pUserID == allowedUserIDstr {
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
				if pUserID == allowedUserIDstr {
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
	pPolicyID    gf_core.GF_ID,
	pOwnerUserID gf_core.GF_ID,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFpolicyUpdateOutput, *gf_core.GFerror) {

	output := &GFpolicyUpdateOutput{}

	//------------------------
	// EXISTS
	existsBool, gfErr := DBsqlExistsByID(pPolicyID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	if !existsBool {
		output.PolicyExistsBool = false
		return output, nil
	}

	//------------------------
	// DB - GET_POLICY_BY_ID
	publicViewBool := true
	updateOp := &GFpolicyUpdateOp{
		PublicViewBool: &publicViewBool,
	}
	gfErr = DBsqlUpdatePolicy(pPolicyID, updateOp, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------

	return output, nil
}

//-------------------------------------------------
// PIPELINE__CREATE

func PipelineCreate(pTargetResourceID gf_core.GF_ID,
	pTargetResourceTypeStr string,
	pOwnerUserID           gf_core.GF_ID,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (*GFpolicy, *gf_core.GFerror) {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	ID                := createID(pTargetResourceID, creationUNIXtimeF)

	policy := &GFpolicy{
		ID:                    ID,     
		CreationUNIXtimeF:     creationUNIXtimeF,
		TargetResourceIDsLst:  []gf_core.GF_ID{pTargetResourceID, },
		TargetResourceTypeStr: pTargetResourceTypeStr,
		OwnerUserID:           pOwnerUserID,
		PublicViewBool:        true,
	}

	// DB
	gfErr := DBsqlCreatePolicy(policy, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return policy, gfErr
}

//---------------------------------------------------

func createID(pTargetResourceID gf_core.GF_ID,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		string(pTargetResourceID),
		fmt.Sprintf("%f", pCreationUNIXtimeF),
	}
	gfID := gf_core.IDcreate(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfID
}