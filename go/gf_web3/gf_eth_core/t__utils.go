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

package gf_eth_core

import (
	"os"
	"fmt"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

func TestGenAddress(pTest *testing.T) {
	

	privateKeyHexStr, publicKeyHexStr, addressHexStr, _ := EthGenerateKeys()

	fmt.Printf("private key hex - %s\n", privateKeyHexStr)
	fmt.Printf("public key hex  - %s\n", publicKeyHexStr)
	fmt.Printf("address hex     - %s\n", addressHexStr)
}

//---------------------------------------------------

func TgetRuntime() (*GF_runtime, *GF_metrics, error) {
	
	//-----------------------
	// MONGODB_HOST
	var mongoDBhostPortStr string
	envMongoDBhostPortStr := os.Getenv("GF_TEST_MONGODB_HOST_PORT")
	if envMongoDBhostPortStr != "" {
		mongoDBhostPortStr = envMongoDBhostPortStr
	} else {
		mongoDBhostPortStr = "localhost:27017"
	}

	//-----------------------
	// SQL_HOST
	var sqlHostStr string
	envSQLhostStr := os.Getenv("GF_TEST_SQL_HOST_PORT")
	if envSQLhostStr != "" {
		sqlHostStr = envSQLhostStr
	} else {
		sqlHostStr = "localhost:5432"
	}

	//-----------------------
	// RUNTIME_SYS
	logFun, logNewFun  := gf_core.LogsInit()
	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_web3_monitor_test",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
		
		// SENTRY - enable it for error reporting
		ErrorsSendToSentryBool: true,
	}

	// VALIDATOR
	validator := gf_core.ValidateInit()
	runtimeSys.Validator = validator

	// HTTP_PROXY
	// if traffic should be routed through a http proxy
	var httpProxyURIstr string
	envHTTPproxyURIstr := os.Getenv("GF_TEST_HTTP_PROXY_URI")
	if envHTTPproxyURIstr != "" {
		httpProxyURIstr = envHTTPproxyURIstr
	}
	runtimeSys.HTTPproxyServerURIstr = httpProxyURIstr

	config := &GF_config{
		Mongodb_host_str:    mongoDBhostPortStr,
		Mongodb_db_name_str: "gf_web3_monitor_test",
	}

	// RUNTIME
	runtime, err := RuntimeGet(config, runtimeSys)
	if err != nil {
		return nil, nil, err
	}
	
	//--------------------
	// SQL

	dbNameStr := "gf_tests"
	dbUserStr := "gf"

	dbHostStr := sqlHostStr

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
	// SENTRY
	
	// FIX!! - load this from ENV var
	sentryEndpointURIstr := "https://702b507d193d45029674fbf98bcedaaf@o502595.ingest.sentry.io/5590469"

	SentryInit(sentryEndpointURIstr)

	//--------------------
	return runtime, nil, nil
}