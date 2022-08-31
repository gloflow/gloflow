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

package gf_identity_lib

import (
	"os"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func TestMain(m *testing.M) {

	logFun, _  = gf_core.InitLogs()
	cliArgsMap = CLI__parse_args(logFun)

	runtimeSys := T__init()

	portInt := 2000
	TestStartService(portInt, runtimeSys)

	/*test_port_int := 2000
	go func() {

		HTTPmux := http.NewServeMux()

		serviceInfo := &GF_service_info{

			// IMPORTANT!! - durring testing dont send emails
			Enable_email_bool: false,
		}
		InitService(HTTPmux, serviceInfo, runtime_sys)
		gf_rpc_lib.Server__init_with_mux(test_port_int, HTTPmux)
	}()
	time.Sleep(2*time.Second) // let server startup*/

	v := m.Run()
	os.Exit(v)
}