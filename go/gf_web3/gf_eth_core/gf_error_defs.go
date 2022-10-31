/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type ErrorDef struct {
	DescrStr string
}

//-------------------------------------------------
func ErrorGetDefs() map[string]gf_core.ErrorDef {

	errorDefs_map := map[string]gf_core.ErrorDef{

		//---------------
		// ETH_CONTRACT
		"eth_contract__not_supported_type": gf_core.ErrorDef{
			DescrStr: "eth contract type encountered is not supported",
		},
		"eth_contract__abi_not_loadable": gf_core.ErrorDef{
			DescrStr: "eth contract ABI cant be parsed from JSON to ABI struct",
		},
		"eth_tx_log__decode": gf_core.ErrorDef{
			DescrStr: "eth transaction log failed to be decoded with a given ABI",
		},
		"eth_contract__disassemble": gf_core.ErrorDef{
			DescrStr: "eth failed to diassemble contract bytecode",
		},

		//---------------
		// ETH_RPC
		"eth_rpc__dial": gf_core.ErrorDef{
			DescrStr: "failed to get Dial/Connect to Ethereum RPC-JSON API",
		},
		"eth_rpc__get_header": gf_core.ErrorDef{
			DescrStr: "failed to get Header via Ethereum RPC-JSON API",
		},
		"eth_rpc__get_block": gf_core.ErrorDef{
			DescrStr: "failed to get Block via Ethereum RPC-JSON API",
		},
		"eth_rpc__get_tx": gf_core.ErrorDef{
			DescrStr: "failed to get Transaction via Ethereum RPC-JSON API",
		},
		"eth_rpc__get_tx_receipt": gf_core.ErrorDef{
			DescrStr: "failed to get Transaction Receipt via Ethereum RPC-JSON API",
		},
		"eth_rpc__get_tx_sender": gf_core.ErrorDef{
			DescrStr: "failed to get Transaction Sender via Ethereum RPC-JSON API",
		},
		"eth_rpc__get_contract_code": gf_core.ErrorDef{
			DescrStr: "failed to Contract code via Ethereum RPC-JSON API",
		},

		//---------------
		
	}
	return errorDefs_map
}