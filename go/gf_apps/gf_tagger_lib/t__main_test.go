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

package gf_tagger_lib

import (
	"fmt"
	"os"
	"context"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
var log_fun func(string,string)
var cli_args_map map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {
	log_fun = gf_core.Init_log_fun()
	cli_args_map = CLI__parse_args(log_fun)
	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------
func Test__main(p_test *testing.T) {

	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"
	test__mongodb_url_str := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)

	log_fun := gf_core.Init_log_fun()


	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_tagger_tests",
		Log_fun:          log_fun,
	}




	mongo_db, _, gf_err := gf_core.Mongo__connect_new(test__mongodb_url_str, test__mongodb_db_name_str, nil, runtime_sys)
	if gf_err != nil {
		panic(-1)
	}


	mongo_coll := mongo_db.Collection("data_symphony")
	runtime_sys.Mongo_db   = mongo_db
	runtime_sys.Mongo_coll = mongo_coll


	test_bookmarking(p_test, runtime_sys)
}

//-------------------------------------------------
func test_bookmarking(p_test *testing.T,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "t__main_test.test_bookmarking()")

	ctx := context.Background()
	validator := gf_core.Validate__init()

	test_user_id_str := gf_core.GF_ID("test_user")
	//------------------
	// CREATE
	input__create := &GF_bookmark__input_create{
		User_id_str: test_user_id_str,
		Url_str:     "https://gloflow.com",
		Description_str: "test bookmark",
		Tags_lst: []string{
			"test", "code", "art",
		},
	}
	gf_err := bookmarks__pipeline__create(input__create,
		validator,
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	//------------------


	// GET_ALL

	input__get_all := &GF_bookmark__input_get_all{
		User_id_str: test_user_id_str,
	}
	output, gf_err := bookmarks__pipeline__get_all(input__get_all,
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}



	spew.Dump(output.Bookmarks_lst)

	//------------------

}