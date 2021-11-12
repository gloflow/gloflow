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
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func T__init() *gf_core.Runtime_sys {

	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"
	test__mongodb_url_str := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)


	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_identity_tests",
		Log_fun:          log_fun,
	}




	mongo_db, _, gf_err := gf_core.Mongo__connect_new(test__mongodb_url_str, test__mongodb_db_name_str, nil, runtime_sys)
	if gf_err != nil {
		panic(-1)
	}


	mongo_coll := mongo_db.Collection("data_symphony")
	runtime_sys.Mongo_db   = mongo_db
	runtime_sys.Mongo_coll = mongo_coll




	return runtime_sys
}