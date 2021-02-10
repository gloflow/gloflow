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
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func TestMain(m *testing.M) {
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func t__get_runtime(p_test *testing.T) (*GF_runtime, *GF_metrics) {

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

	// RUNTIME
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


	return runtime, metrics
}

//---------------------------------------------------
func t__get_abis() map[string]*GF_eth__abi {
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
	
	abis_map := map[string]*GF_eth__abi{
		"erc20": &GF_eth__abi{
			Type_str: "erc20",
			Def_lst:  abi_lst,
		},
	}
	return abis_map
}