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

package gf_eth_monitor_core

import (
	"os"
	"testing"
	"context"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__contract_opcodes(p_test *testing.T) {

	ctx := context.Background()

	host_port_str := os.Getenv("GF_TEST_WORKER_INSPECTOR_HOST_PORT")

	runtime, _ := t__get_runtime(p_test)

	//--------------------
	tx_id_hex_str := "0x62974c8152c87e14880c54007260e0d5fe9d182c2cd22c58797735a9ae88370a"

	// GET_TRACE
	gf_tx_trace, gf_err := Eth_tx__get_trace__from_worker_inspector(tx_id_hex_str,
		host_port_str,
		ctx,
		runtime.Runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}









	spew.Dump(gf_tx_trace)



	
	//--------------------
	// PLOT

	plugins_info := &GF_py_plugins{
		Base_dir_path_str: "./../../py/plugins",
	}

	_, gf_err = py__run_plugin__plot_tx_trace(tx_id_hex_str,
		gf_tx_trace,
		plugins_info,
		runtime.Runtime_sys)
	if gf_err != nil {

		p_test.Fail()

	}




	//--------------------
}