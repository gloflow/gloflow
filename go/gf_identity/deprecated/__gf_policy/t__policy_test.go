/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
	// "github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//-------------------------------------------------

func TestPolicy(pTest *testing.T) {

	fmt.Println(" TEST__POLICY >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	ctx := context.Background()
	runtimeSys := Tinit("gf_policy", cliArgsMap)

	gfErr := DBsqlCreateTables(runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	targetResourceID := gf_core.GF_ID("test_resource")
	ownerUserID      := gf_core.GF_ID("test_user")
	thirdpartyUserID := gf_core.GF_ID("other_user")

	//----------------------
	fmt.Println("create policy >>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	// CREATE
	policy, gfErr := CreateTestPolicy(targetResourceID, ownerUserID, ctx, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	gfErr = gf_core.DBsqlViewTableStructure("gf_policy", runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	// add some third-party user to a list of editors
	policy.EditorsUserIDsLst = append(policy.EditorsUserIDsLst, string(thirdpartyUserID))
	
	//----------------------
	fmt.Println("get policies >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	// GET_POLICIES
	policiesLst, gfErr := DBsqlGetPolicies(targetResourceID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}




	spew.Dump(policiesLst)


	//----------------------
	fmt.Println("validate policy >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")




	policiesDefsMap := getDefs()

	validBool := policySingleVerify(GF_POLICY_OP__FLOW_ADD_IMG,
		policy,
		thirdpartyUserID,
		policiesDefsMap)


	validDeleteBool := policySingleVerify(GF_POLICY_OP__FLOW_DELETE,
		policy,
		thirdpartyUserID,
		policiesDefsMap)

	fmt.Println("VALID", validBool, validDeleteBool)


	assert.True(pTest, validBool, "adding an image to a flow by a user who has the editor role should be allowed")
	assert.True(pTest, !validDeleteBool, "deleting a flow by a user who does not have the admin role should not be allowed")

	//----------------------
	fmt.Println("validate policy >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	policy.EditorsUserIDsLst = []string{}
	policy.TaggersUserIDsLst = append(policy.TaggersUserIDsLst, string(thirdpartyUserID))


	validTaggingBool := policySingleVerify(GF_POLICY_OP__FLOW_ADD_IMG,
		policy,
		thirdpartyUserID,
		policiesDefsMap)

	fmt.Println("VALID", validTaggingBool)


	assert.True(pTest, !validTaggingBool, "adding an image to a flow by a user who can only tag is not allowed")
}