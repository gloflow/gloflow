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
	"bytes"
	"testing"
	"context"
	"encoding/json"
	"net/http"
	"time"
	"net/http/httptest"
	"io/ioutil"
	// "github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

var logFun func(string, string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//---------------------------------------------------

func TestMain(m *testing.M) {
	logFun, logNewFun = gf_core.LogsInit()
	cliArgsMap = gf_images_core.CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------

func TestBasicImageOps(pTest *testing.T) {

	serviceNameStr := "gf_images_ops_tests"

	/*
	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: serviceNameStr,
		LogFun:         logFun,
		LogNewFun:      logNewFun,
	}
	*/
	ctx := context.Background()

	// MONGODB
	test__mongodb_host_str := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	// test__mongodb_url_str     := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)
	// test__mongodb_db_name_str := "gf_tests"
	
	// SQL
	sqlHostStr := cliArgsMap["sql_host_str"].(string)


	/*
	mongodbDB, _, gfErr := gf_core.MongoConnectNew(test__mongodb_url_str, test__mongodb_db_name_str, nil, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()
	}
	mongodbColl := mongodbDB.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongodbDB
	runtimeSys.Mongo_coll = mongodbColl
	*/


	runtimeSys := gf_identity.Tinit(serviceNameStr, test__mongodb_host_str, sqlHostStr, logNewFun, logFun)


	gfErr := gf_images_core.DBsqlCreateTables(ctx, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()

	}

	//------------------
	// TEST_HTTP_PARSING
	metaMap := map[string]interface{}{
		"file_name_str":          "test_file",
		"file_name_original_str": "test_file_original",
	}
	metaBytesLst, _ := json.Marshal(metaMap)

	req := httptest.NewRequest(http.MethodPost, "/some_endpoint", 
        ioutil.NopCloser(bytes.NewBufferString(string(metaBytesLst))))

	iMap, gfErr :=  gf_core.HTTPgetInput(req, runtimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.Fail()
	}

	spew.Dump(iMap)
	
	//------------------
	// CREATE_TEST_IMAGES
	
	test_img_0 := &gf_images_core.GFimage{
		IDstr: gf_images_core.GFimageID(fmt.Sprintf("test_img_%d", time.Now().UnixNano())),
		T_str: "img",
		FlowsNamesLst: []string{"flow_0"},
	}


	// MONGO
	gfErr = gf_images_core.DBmongoPutImage(test_img_0, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}


	// SQL
	gfErr = gf_images_core.DBsqlPutImage(test_img_0, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//------------------
}