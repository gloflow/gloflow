/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

package gf_policy

import (
	// "fmt"
	// "context"
	// "strings"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

var logFun func(p_g string, p_m string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//-------------------------------------------------

func Tinit(pServiceNameStr string,
	pCliArgsMap map[string]interface{}) *gf_core.RuntimeSys {

	/*
	testMongodbHostStr   := pMongoHostStr // cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	testMongodbDBnameStr := "gf_tests"
	testMongodbURLstr    := fmt.Sprintf("mongodb://%s", testMongodbHostStr)
	*/

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: pServiceNameStr, // "gf_identity_tests",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
		Validator:      gf_core.ValidateInit(),
	}

	/*
	mongoDB, _, gfErr := gf_core.MongoConnectNew(testMongodbURLstr,
		testMongodbDBnameStr,
		nil,
		runtimeSys)
	if gfErr != nil {
		panic(-1)
	}

	mongoColl := mongoDB.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongoDB
	runtimeSys.Mongo_coll = mongoColl
	*/

	//--------------------
	// SQL

	dbNameStr := "gf_tests"
	dbUserStr := "gf"

	dbHostStr := pCliArgsMap["sql_host_str"].(string)

	sqlDB, gfErr := gf_core.DBsqlConnect(dbNameStr,
		dbUserStr,
		"", // config.SQLpassStr,
		dbHostStr,
		runtimeSys)
	if gfErr != nil {
		panic(-1)
	}

	runtimeSys.SQLdb = sqlDB

	//--------------------
	
	return runtimeSys
}