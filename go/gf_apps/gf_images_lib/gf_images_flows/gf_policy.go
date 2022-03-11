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

package gf_images_flows

import (
	"fmt"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GFflowPolicy struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	IDstr             gf_core.GF_ID      `bson:"id_str"`
	DeletedBool       bool               `bson:"deleted_bool"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`
	
	FlowIDstr         string             `bson:"flow_id_str"`
	OwnerUserIDstr    gf_core.GF_ID      `bson:"owner_user_id_str"`

	// if the flow is fully public and all users (including anonymous) can view it
	PublicViewBool bool `bson:"public_view_bool"`

	// editors are users that are allowed by the owner to update/add/remove items to the flow
	EditorsUserIDsLst []gf_core.GF_ID `bson:"editors_user_ids_lst"`

	// viewers are users that can view if PublicViewBool is false
	ViewersUserIDsLst []gf_core.GF_ID `bson:"viewers_user_ids_lst"`
	
	// taggers are users that can attach tags/notes to flows
	TaggersUserIDsLst []gf_core.GF_ID `bson:"taggers_user_ids_lst"`
}

//-------------------------------------------------
// VERIFY
func policyVerify(pFlowsNamesLst []string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GF_error {


	for _, flowNameStr := range pFlowsNamesLst {

		policy, gfErr := DBgetPolicy(flowNameStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		fmt.Println(policy)
	}

	



	return nil
	
}

//-------------------------------------------------
// PIPELINE__CREATE
func policyPipelineCreate(pFlowIDstr string,
	pOwnerUserIDstr gf_core.GF_ID,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.Runtime_sys) *gf_core.GF_error {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	IDstr             := policyCreateID(pFlowIDstr, creationUNIXtimeF)

	policy := &GFflowPolicy{
		IDstr:             IDstr,     
		CreationUNIXtimeF: creationUNIXtimeF,
		FlowIDstr:         pFlowIDstr,
		OwnerUserIDstr:    pOwnerUserIDstr,
		PublicViewBool:    true,
	}

	gfErr := DBcreatePolicy(policy, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return gfErr
}

//---------------------------------------------------
func policyCreateID(pFlowIDstr string,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		pFlowIDstr,
		fmt.Sprintf("%f", pCreationUNIXtimeF),
	}
	gfIDstr := gf_core.ID__create(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfIDstr
}