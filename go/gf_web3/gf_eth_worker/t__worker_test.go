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

package gf_eth_worker

import (
	"os"
	"fmt"
	"testing"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__worker(p_test *testing.T) {

	fmt.Println("TEST__WORKER ==============================================")
	
	ctx := context.Background()

	block_int     := 4634748
	host_port_str := os.Getenv("GF_TEST_WORKER_INSPECTOR_HOST_PORT")


	//--------------------
	// RUNTIME_SYS
	logFun, _   := gf_core.LogsInit()
	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_eth_monitor_core__tests",
		LogFun:           logFun,
		
		// SENTRY - enable it for error reporting
		ErrorsSendToSentryBool: true,
	}

	config := &GF_config{
		Mongodb_host_str:    "localhost:27017",
		Mongodb_db_name_str: "gf_eth_monitor",
	}
	runtime, err := RuntimeGet(config, runtimeSys)
	if err != nil {
		p_test.Fail()
	}

	//--------------------
	// GET_BLOCK__FROM_WORKER_INSPECTOR
	gf_block, gfErr := eth_blocks__get_block__from_worker_inspector(uint64(block_int),
		host_port_str,
		ctx,
		runtimeSys)

	if gfErr != nil {
		p_test.Fail()
	}




	spew.Dump(gf_block)


	abi_type_str := "erc20"
	abis_lst, gfErr := Eth_contract__db__get_abi(abi_type_str, ctx, nil, runtime)
	if gfErr != nil {
		p_test.Fail()
	}
	abis_map := map[string]*GF_eth__abi{
		"erc20": abis_lst[0],
	}

	gfErr = eth_tx__enrich_from_block(gf_block,
		abis_map,
		ctx,
		nil,
		runtime)
	if gfErr != nil {
		p_test.Fail()
	}




}