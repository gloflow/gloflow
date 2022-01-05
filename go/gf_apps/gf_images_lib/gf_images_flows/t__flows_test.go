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

package gf_images_flows

import (
	"fmt"
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__get_all(p_test *testing.T) {

	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_images_flows_tests",
		Log_fun:          log_fun,
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
	

	ctx := context.Background()

	//------------------
	// CREATE_TEST_IMAGES
	test_img_0 := &gf_images_core.GF_image{
		Id_str: "test_img_0",
		T_str:  "img",
		Flows_names_lst: []string{"flow_0"},
	}
	test_img_1 := &gf_images_core.GF_image{
		Id_str: "test_img_1",
		T_str:  "img",
		Flows_names_lst: []string{"flow_0"},
	}
	test_img_2 := &gf_images_core.GF_image{
		Id_str: "test_img_2",
		T_str:  "img",
		Flows_names_lst: []string{"flow_0", "flow_1"},
	}
	test_img_3 := &gf_images_core.GF_image{
		Id_str: "test_img_3",
		T_str:  "img",
		Flows_names_lst: []string{"flow_1", "flow_2"},
	}
	gf_err = gf_images_core.DB__put_image(test_img_0, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}
	gf_err = gf_images_core.DB__put_image(test_img_1, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}
	gf_err = gf_images_core.DB__put_image(test_img_2, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}
	gf_err = gf_images_core.DB__put_image(test_img_3, ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}
 
	//------------------


	all_flows_names_lst, gf_err := flows__get_all__pipeline(ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	spew.Dump(all_flows_names_lst)

	assert.True(p_test, len(all_flows_names_lst) == 3, "3 flows in total should have been discovered")
	assert.True(p_test, all_flows_names_lst[0]["flow_name_str"].(string) == "flow_0", "first flow should be flow_0")
	assert.True(p_test, all_flows_names_lst[1]["flow_name_str"].(string) == "flow_1", "first flow should be flow_1")
	assert.True(p_test, all_flows_names_lst[2]["flow_name_str"].(string) == "flow_2", "first flow should be flow_2")

	assert.True(p_test, all_flows_names_lst[0]["flow_imgs_count_int"].(int32) == 3, "first flow should have a count of 3")
	assert.True(p_test, all_flows_names_lst[1]["flow_imgs_count_int"].(int32) == 2, "second flow should have a count of 2")
	assert.True(p_test, all_flows_names_lst[2]["flow_imgs_count_int"].(int32) == 1, "third flow should have a count of 1")
}