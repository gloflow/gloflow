/*
GloFlow application and media management/publishing platform
Copyright (C) 2026 Ivan Trajkovic

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

package gf_solo_service

import (
	"fmt"
	"path"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func RuntimeGet(pConfigPathStr string,
	pBootHooks     *gf_core.ExternalBootHook,
	pExternalHooks *gf_core.ExternalHooks,
	pLogFun        func(string, string),
	pLogNewFun     gf_core.GFlogFun) (*gf_core.RuntimeSys, *gf_core.GFconfig, error) {

	//--------------------
	// RUNTIME_SYS
	runtimeSys := &gf_core.RuntimeSys{
		AppNameStr:     "gf_solo",
		ServiceNameStr: "gf_solo",
		LogFun:         pLogFun,
		LogNewFun:      pLogNewFun,

		// VALIDATOR
		Validator: gf_core.ValidateInit(),

		// EXTERNAL_HOOKS
		ExternalHooks: pExternalHooks,
	}

	//--------------------
	// CONFIG
	configDirPathStr := path.Dir(pConfigPathStr)  // "./../config/"
	configNameStr    := path.Base(pConfigPathStr) // "gf_solo"

	// HOOK
	var hookConfigLoadCallbackFun gf_core.GFhookConfigLoadCallback
	if pBootHooks != nil {
		hookConfigLoadCallbackFun = pBootHooks.ConfigLoadCallback
	}

	config, err := ConfigInit(configDirPathStr,
		configNameStr,
		hookConfigLoadCallbackFun,
		runtimeSys)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to load config")
		return nil, nil, err
	}

	runtimeSys.Config = config
	runtimeSys.EnvStr = config.EnvStr

	//--------------------
	// GF_IDENTITY
	runtimeSys.IdentitySubsystemTypeStr = config.AuthSubsystemTypeStr

	//--------------------
	// SENTRY - ERROR_REPORTING
	if config.SentryDSNstr != "" {

		fmt.Println("Initializing Sentry error reporting...")

		sentryDSNstr := config.SentryDSNstr
		runtimeSys.SentryDSNstr = sentryDSNstr
		runtimeSys.ErrorsSendToMongodbBool = true // enable it for error reporting

		sentrySampleRateDefaultF := 1.0
		sentryTracingRateForHandlersMap := map[string]float64{

		}
		err := gf_core.ErrorInitSentry(sentryDSNstr,
			sentryTracingRateForHandlersMap,
			sentrySampleRateDefaultF)
		if err != nil {
			panic(err)
		}

		defer sentry.Flush(2 * time.Second)
	}

	//--------------------
	// SQL

	sqlDB, sqlDSNstr, gfErr := gf_core.DBsqlConnect(config.SQLdbNameStr,
		config.SQLuserNameStr,
		config.SQLpassStr,
		config.SQLhostStr,
		"require", // SSL mode - required for PostgreSQL 18
		runtimeSys)

	runtimeSys.SQLdb = sqlDB
	runtimeSys.SQLdsnStr = sqlDSNstr

	//--------------------
	// MONGODB
	mongodbHostStr := config.MongoHostStr
	mongodbURLstr  := fmt.Sprintf("mongodb://%s", mongodbHostStr)
	fmt.Printf("mongodb_host    - %s\n", mongodbHostStr)
	fmt.Printf("mongodb_db_name - %s\n", config.MongoDBnameStr)

	mongodbDB, _, gfErr := gf_core.MongoConnectNew(mongodbURLstr,
		config.MongoDBnameStr,
		nil,
		runtimeSys)
	if gfErr != nil {
		return nil, nil, gfErr.Error
	}

	runtimeSys.Mongo_db   = mongodbDB
	runtimeSys.Mongo_coll = mongodbDB.Collection("data_symphony")
	fmt.Printf("mongodb connected...\n")

	//--------------------
	return runtimeSys, config, nil
}
