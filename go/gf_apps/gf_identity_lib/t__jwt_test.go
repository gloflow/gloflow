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
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func TestJWT(pTest *testing.T) {

	fmt.Println(" TEST__IDENTITY_JWT >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	runtimeSys := T__init()

	testJWTmain(pTest, runtimeSys)
}

//-------------------------------------------------

func testJWTmain(pTest *testing.T,
	pRuntimeSys *gf_core.RuntimeSys) {


	ctx := context.Background()

	testUserAddressETH := gf_identity_core.GFuserAddressETH("0xBA47Bef4ca9e8F86149D2f109478c6bd8A642C97")

	// JWT_GENERATE
	userIdentifierStr := string(testUserAddressETH)
	jwtVal, gfErr := gf_session.JWTpipelineGenerate(userIdentifierStr,
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}
	
	// JWT_VALIDATE
	valid_bool, userIdentifierStr, gfErr := gf_session.JWTvalidate(jwtVal,
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		pTest.Fail()
	}



	assert.True(pTest, valid_bool == true, "test JWT token is not valid, when it should be")
	assert.True(pTest, userIdentifierStr == string(testUserAddressETH),
		"test user_identifier extracted from JWT durring validation is the same as the input test eth address")
	
}