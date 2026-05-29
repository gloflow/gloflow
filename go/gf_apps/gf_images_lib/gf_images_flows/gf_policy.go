/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"	
)

//-------------------------------------------------
// VERIFY_POLICY

func VerifyPolicy(pOpStr string,
	pFlowsNamesLst []string,
	pUserID gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	// POLICY_VERIFY
	flowsIDsLst, gfErr := DBsqlGetFlowsIDs(pFlowsNamesLst, pCtx, pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	//---------------------
	// HOOK
	if pRuntimeSys != nil && pRuntimeSys.PolicyHooks != nil &&
		pRuntimeSys.PolicyHooks.VerifyCallback != nil {

		for _, flowIDstr := range flowsIDsLst {
			// Use hook for verification
			allowed, gfErr := pRuntimeSys.PolicyHooks.VerifyCallback(
				pUserID,
				flowIDstr,
				"flow",
				pOpStr,
				pCtx,
				pRuntimeSys)

			if gfErr != nil {
				return false, gfErr
			}

			if !allowed {
				return false, nil
			}
		}
		return true, nil
	}

	//---------------------

	// Fallback: No policy hooks available - only allow actions on flows owned by the user
	// This protects against misconfiguration where policies are accidentally disabled
	for _, flowIDstr := range flowsIDsLst {
		isOwned, gfErr := DBsqlIsFlowOwnedByUser(flowIDstr, pUserID, pCtx, pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}
		if !isOwned {
			return false, nil
		}
	}

	return true, nil
}
