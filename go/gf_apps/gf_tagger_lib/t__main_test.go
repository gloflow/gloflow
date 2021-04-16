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
	"os"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
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

	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) //"127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"

	log_fun      := gf_core.Init_log_fun()
	mongodb_db   := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_tagger_tests",
		Log_fun:          log_fun,
		Mongodb_coll:     mongodb_coll,
	}

	test_posts_tagging(runtime_sys)
}

//-------------------------------------------------
func test_posts_tagging(p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","t__main_test.test_posts_tagging()")


}