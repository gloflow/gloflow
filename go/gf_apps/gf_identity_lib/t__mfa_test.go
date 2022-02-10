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

package gf_identity_lib

import (
	"fmt"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func Test__mfa(p_test *testing.T) {

	fmt.Println(" TEST__IDENTITY_MFA >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	runtime_sys := T__init()

	test_mfa_main(p_test, runtime_sys)
}

//-------------------------------------------------
func test_mfa_main(p_test *testing.T,
	p_runtime_sys *gf_core.Runtime_sys) {




	// CODE THATS ENTERED INTO GOOGLE AUTH MANUALLY HAS TO BE 
	// BASE32 ENCODED
	secret_key_base32_str := "aabbccddeeffgghh"
	token_str, gf_err := TOTP_generate_value(secret_key_base32_str, p_runtime_sys)
	if gf_err != nil {
		p_test.FailNow()
	}

	fmt.Println("TOTP token - ", token_str)
}