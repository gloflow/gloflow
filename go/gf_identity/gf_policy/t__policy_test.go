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
	// "github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func TestPolicy(pTest *testing.T) {

	fmt.Println(" TEST__POLICY >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")


	ctx := context.Background()
	runtimeSys := Tinit("gf_policy", cliArgsMap)


	gfErr := DBsqlCreateTables(runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}


	targetResourceID := gf_core.GF_ID("test_resource")
	ownerUserID := gf_core.GF_ID("test_user")

	// CREATE
	policy, gfErr := PipelineCreate(targetResourceID,
		ownerUserID,
		ctx,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	spew.Dump(policy)
	
}