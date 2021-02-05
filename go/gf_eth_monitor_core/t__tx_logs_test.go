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
	// "os"
	"fmt"
	"testing"
	"context"
	"math/big"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func Test__get_tx_logs(p_test *testing.T) {

	fmt.Println("TEST__MAIN ==============================================")
	

	ctx := context.Background()

	//--------------------
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
		Mongodb_db_name_str: "gf_test",
	}
	runtime, err := Runtime__get(config, runtime_sys)
	if err != nil {
		p_test.Fail()
	}





	

	
	// METRICS
	metrics_port := 9110
	metrics, gf_err := Metrics__init(metrics_port)
	if gf_err != nil {
		p_test.Fail()
	}




	abi_map       := t__get_erc20_abi()
	coll_name_str := "gf_eth_meta__contracts_abi"
	gf_err = gf_core.Mongo__insert(abi_map, coll_name_str, &ctx, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
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
	decoded_logs_lst, gf_err := Eth_tx__enrich_logs(tx_logs,
		ctx,
		metrics,
		runtime)
	if gf_err != nil {
		p_test.Fail()
	}



	value_int := decoded_logs_lst[0]["value"].(*big.Int).Uint64()
	assert.EqualValues(p_test, value_int, 445550000, "the decoded event Eth log value should be equal to 445550000")
	





}

//---------------------------------------------------
func t__get_erc20_abi() map[string]interface{} {
	abi_json_str := `[
		{
			"constant": true,
			"inputs": [],
			"name": "name",
			"outputs": [
				{
					"name": "",
					"type": "string"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_spender",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "approve",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "totalSupply",
			"outputs": [
				{
					"name": "",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_from",
					"type": "address"
				},
				{
					"name": "_to",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "transferFrom",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "decimals",
			"outputs": [
				{
					"name": "",
					"type": "uint8"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [
				{
					"name": "_owner",
					"type": "address"
				}
			],
			"name": "balanceOf",
			"outputs": [
				{
					"name": "balance",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "symbol",
			"outputs": [
				{
					"name": "",
					"type": "string"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_to",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "transfer",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [
				{
					"name": "_owner",
					"type": "address"
				},
				{
					"name": "_spender",
					"type": "address"
				}
			],
			"name": "allowance",
			"outputs": [
				{
					"name": "",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"payable": true,
			"stateMutability": "payable",
			"type": "fallback"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "owner",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "spender",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "Approval",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "from",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "to",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "Transfer",
			"type": "event"
		}
	]`
	
	abi_lst := []map[string]interface{}{}
	json.Unmarshal([]byte(abi_json_str), &abi_lst)
	
	abi_map := map[string]interface{}{
		"type_str": "erc20",
		"def_lst":  abi_lst,
	}
	return abi_map
}