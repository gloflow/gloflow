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
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func T__get_runtime(p_test *testing.T) (*GF_runtime, *GF_metrics) {

	// RUNTIME_SYS
	log_fun     := gf_core.Init_log_fun()
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_eth_monitor_core__tests",
		Log_fun:          log_fun,
		
		// SENTRY - enable it for error reporting
		Errors_send_to_sentry_bool: true,
	}

	config := &GF_config{
		Mongodb_host_str:    "localhost:27017",
		Mongodb_db_name_str: "gf_eth_monitor",
	}

	// RUNTIME
	runtime, err := RuntimeGet(config, runtime_sys)
	if err != nil {
		p_test.Fail()
	}
	

	// SENTRY
	sentry_endpoint_uri_str := "https://702b507d193d45029674fbf98bcedaaf@o502595.ingest.sentry.io/5590469"
	Sentry__init(sentry_endpoint_uri_str)


	return runtime, nil
}