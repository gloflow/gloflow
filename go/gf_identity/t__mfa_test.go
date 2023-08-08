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

package gf_identity

import (
	"fmt"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

func TestMFA(pTest *testing.T) {

	fmt.Println(" TEST__IDENTITY_MFA >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	serviceNameStr := "gf_identity_test"
	mongoHostStr   := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr     := cliArgsMap["sql_host_str"].(string)
	runtimeSys     := Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)

	testMFAmain(pTest, runtimeSys)
}

//-------------------------------------------------

func testMFAmain(pTest *testing.T,
	pRuntimeSys *gf_core.RuntimeSys) {

	// CODE THATS ENTERED INTO GOOGLE AUTH MANUALLY HAS TO BE 
	// BASE32 ENCODED
	secretKeyBase32str := "aabbccddeeffgghh"
	tokenStr, gfErr := TOTPgenerateValue(secretKeyBase32str, pRuntimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	fmt.Println("TOTP token - ", tokenStr)
}