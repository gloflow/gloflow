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
	"os"
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
var log_fun func(p_g string, p_m string)
var cli_args_map map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {

	log_fun      = gf_core.Init_log_fun()
	cli_args_map = CLI__parse_args(log_fun)

	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------
func Test__jwt(p_test *testing.T) {

	fmt.Println(" TEST__IDENTITY_JWT >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	runtime_sys := T__init()

	test_jwt_main(p_test, runtime_sys)
}

//-------------------------------------------------
func test_jwt_main(p_test *testing.T,
	p_runtime_sys *gf_core.Runtime_sys) {



	test_user_address_eth := GF_user_address_eth("")
	test_signing_key_str  := GF_jwt_secret_key_val("fdsfsdf")
	creation_unix_time_f  := float64(time.Now().UnixNano())/1000000000.0

	// JWT_GENERATE
	jwt_val, gf_err := jwt__generate(test_user_address_eth,
		test_signing_key_str,
		creation_unix_time_f,
		p_runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}
	
	// JWT_VALIDATE
	valid_bool, gf_err := jwt__validate(jwt_val,
		test_signing_key_str,
		p_runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}



	assert.True(p_test, valid_bool == true, "test JWT token is not valid, when it should be")
	
}