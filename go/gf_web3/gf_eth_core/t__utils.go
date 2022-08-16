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
	"github.com/gloflow/gloflow/go/gf_core"
)

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
	// RUNTIME_SYS
	logFun, _  := gf_core.InitLogs()
	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_web3_monitor_test",
		LogFun:           logFun,
		
		// SENTRY - enable it for error reporting
		Errors_send_to_sentry_bool: true,
	}

	// VALIDATOR
	validator := gf_core.Validate__init()
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
	

	// SENTRY
	sentry_endpoint_uri_str := "https://702b507d193d45029674fbf98bcedaaf@o502595.ingest.sentry.io/5590469"
	Sentry__init(sentry_endpoint_uri_str)


	return runtime, nil, nil
}