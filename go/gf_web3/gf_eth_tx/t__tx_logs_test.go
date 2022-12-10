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

package gf_eth_tx

import (
	// "os"
	"fmt"
	"testing"
	"context"
	"math/big"
	// "encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func Test__get_tx_logs(p_test *testing.T) {

	fmt.Println("TEST__TX_LOGS ==============================================")
	

	ctx := context.Background()

	//--------------------
	// RUNTIME_SYS
	logFun, _   := gf_core.LogsInit()
	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_web3_monitor_core__tests",
		LogFun:           logFun,
		
		// SENTRY - enable it for error reporting
		ErrorsSendToSentryBool: true,
	}

	config := &GF_config{
		Mongodb_host_str:    "localhost:27017",
		Mongodb_db_name_str: "gf_web3_monitor",
	}
	runtime, err := RuntimeGet(config, runtime_sys)
	if err != nil {
		p_test.Fail()
	}

	//--------------------
	// DB_INSERT
	abis_map       := t__get_abis()
	coll_name_str := "gf_web3_meta__contracts_abi"
	for _, gf_abi := range abis_map {
		gfErr := gf_core.MongoInsert(gf_abi, coll_name_str,
			map[string]interface{}{
				"caller_err_msg_str": "failed to insert contract ABI record into DB in gf_eth_monitor test t__tx_logs_test",
			},
			ctx, runtime_sys)
		if gfErr != nil {
			p_test.Fail()
		}
	}

	//--------------------

	// TETHER_ERC20_TX_LOG
	// https://etherscan.io/tx/0x0cc8148371a953793498823ecff6754d0a5f2ee648a1bdb696100a9dd8538d05
	tx_logs := []*GF_eth__log{
		{
			Address_str: "0xdac17f958d2ee523a2206206994597c13d831ec7", // Tether contract
			Topics_lst: []string{
				"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				"0x0000000000000000000000001235f58b63cbf96171f4454e0c9eb7dfba5597bd",
				"0x000000000000000000000000fe17e587923c8678fb2a93c167f121dce1a027ce",
			},
			Data_hex_str: "0x000000000000000000000000000000000000000000000000000000001a8e8db0",	
		},
	}

	decoded_logs_lst, gfErr := Eth_tx__enrich_logs(tx_logs,
		abis_map,
		ctx,
		nil,
		runtime)
	if gfErr != nil {
		p_test.Fail()
	}

	

	value_int := decoded_logs_lst[0]["value"].(*big.Int).Uint64()
	assert.EqualValues(p_test, value_int, 445550000, "the decoded event Eth log value should be equal to 445550000")
	



}