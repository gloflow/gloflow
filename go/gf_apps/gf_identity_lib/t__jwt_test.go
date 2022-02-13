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
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Test__jwt(p_test *testing.T) {

	fmt.Println(" TEST__IDENTITY_JWT >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	runtime_sys := T__init()

	test_jwt_main(p_test, runtime_sys)
}

//-------------------------------------------------
func test_jwt_main(p_test *testing.T,
	p_runtime_sys *gf_core.Runtime_sys) {


	ctx := context.Background()

	test_user_address_eth := GF_user_address_eth("0xBA47Bef4ca9e8F86149D2f109478c6bd8A642C97")

	// JWT_GENERATE
	user_identifier_str := string(test_user_address_eth)
	jwt_val, gf_err := gf_session.JWT__pipeline__generate(user_identifier_str,
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}
	
	// JWT_VALIDATE
	valid_bool, user_identifier_str, gf_err := gf_session.JWT__validate(jwt_val,
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}



	assert.True(p_test, valid_bool == true, "test JWT token is not valid, when it should be")
	assert.True(p_test, user_identifier_str == string(test_user_address_eth),
		"test user_identifier extracted from JWT durring validation is the same as the input test eth address")
	
}