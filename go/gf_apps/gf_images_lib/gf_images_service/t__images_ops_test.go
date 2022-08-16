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

package gf_images_service

import (
	"os"
	"fmt"
	"testing"
	"context"
	// "github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	// "github.com/davecgh/go-spew/spew"
)


var logFun func(string, string)
var cli_args_map map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {
	logFun = gf_core.InitLogs()
	cli_args_map = gf_images_core.CLI__parse_args(logFun)
	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------
func Test__basic_image_ops(p_test *testing.T) {

	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_images_ops_tests",
		LogFun:           logFun,
	}

	// MONGODB
	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_url_str     := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)
	test__mongodb_db_name_str := "gf_tests"
	mongodb_db, _, gf_err := gf_core.Mongo__connect_new(test__mongodb_url_str, test__mongodb_db_name_str, nil, runtime_sys)
	if gf_err != nil {
		fmt.Println(gf_err.Error)
		p_test.Fail()
	}
	mongodb_coll := mongodb_db.Collection("data_symphony")
	runtime_sys.Mongo_coll = mongodb_coll
	
	//------------------
	ctx := context.Background()

	//------------------
	// CREATE_TEST_IMAGES
	test_img_0 := &gf_images_core.GF_image{
		Id_str: "test_img_0",
		T_str:  "img",
		Flows_names_lst: []string{"flow_0"},
	}
	gf_err = gf_images_core.DB__put_image(test_img_0, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	//------------------





}