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
	logFun, _ = gf_core.InitLogs()
	cli_args_map = gf_images_core.CLI__parse_args(logFun)
	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------

func Test__basic_image_ops(p_test *testing.T) {

	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_images_ops_tests",
		LogFun:           logFun,
	}

	// MONGODB
	test__mongodb_host_str    := cli_args_map["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_url_str     := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)
	test__mongodb_db_name_str := "gf_tests"
	
	mongodbDB, _, gfErr := gf_core.MongoConnectNew(test__mongodb_url_str, test__mongodb_db_name_str, nil, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		p_test.Fail()
	}
	mongodbColl := mongodbDB.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongodbDB
	runtimeSys.Mongo_coll = mongodbColl
	
	//------------------
	ctx := context.Background()

	//------------------
	// CREATE_TEST_IMAGES
	test_img_0 := &gf_images_core.GFimage{
		IDstr: "test_img_0",
		T_str:  "img",
		Flows_names_lst: []string{"flow_0"},
	}
	gfErr = gf_images_core.DBputImage(test_img_0, ctx, runtimeSys)
	if gfErr != nil {
		p_test.Fail()
	}

	//------------------





}