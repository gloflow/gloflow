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

package gf_home_lib

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
var logFun func(p_g string, p_m string)
var cliArgsMap map[string]interface{}

//-------------------------------------------------
func Tinit() *gf_core.RuntimeSys {

	test__mongodb_host_str    := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"
	test__mongodb_url_str := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)


	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_home_tests",
		LogFun:           logFun,
		Validator:        gf_core.Validate__init(),
	}




	mongo_db, _, gf_err := gf_core.Mongo__connect_new(test__mongodb_url_str,
		test__mongodb_db_name_str,
		nil,
		runtimeSys)
	if gf_err != nil {
		panic(-1)
	}


	mongo_coll := mongo_db.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongo_db
	runtimeSys.Mongo_coll = mongo_coll




	return runtimeSys
}

//-------------------------------------------------