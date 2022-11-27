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
	"strings"
	"strconv"
	"net/http"
	"encoding/hex"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func Http__get_arg__acc_address_hex(p_req *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	qs_map := p_req.URL.Query()

	var account_address_str string
	if acc_addr_lst, ok := qs_map["acc_addr"]; ok {
		account_address_hex__unverified_str := acc_addr_lst[0]


		// VALIDATE - parse hex address to validate it
		_ = eth_common.HexToAddress(account_address_hex__unverified_str)

		account_address_str = account_address_hex__unverified_str
	}

	return account_address_str, nil
}

//-------------------------------------------------
func Http__get_arg__block_range(p_resp http.ResponseWriter,
	p_req         *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (uint64, uint64, *gf_core.GFerror) {

	qs_map := p_req.URL.Query()

	var block_start_uint uint64
	var block_end_uint   uint64
	if block_range_lst, ok := qs_map["br"]; ok {

		block_range_str := block_range_lst[0]

		range_lst := strings.Split(block_range_str, "-")

		if len(range_lst) != 2 {
			gfErr := gf_core.ErrorCreate("suplied block range is not composed of start/end block number",
				"verify__invalid_value_error",
				map[string]interface{}{"block_range_str": block_range_str,},
				nil, "gf_eth_monitor_lib", pRuntimeSys)
			return 0, 0, gfErr
		}
		
		block_start_int, err := strconv.Atoi(range_lst[0])
		if err != nil {
			gfErr := gf_core.ErrorCreate("suplied block_start is not an integer",
				"int_parse_error",
				map[string]interface{}{"block_start_str": range_lst[0],},
				nil, "gf_eth_monitor_lib", pRuntimeSys)
			return 0, 0, gfErr
		}

		block_end_int, err := strconv.Atoi(range_lst[1])
		if err != nil {
			gfErr := gf_core.ErrorCreate("suplied block_end is not an integer",
				"int_parse_error",
				map[string]interface{}{"block_end_str": range_lst[1],},
				nil, "gf_eth_monitor_lib", pRuntimeSys)
			return 0, 0, gfErr
		}


		block_start_uint = uint64(block_start_int)
		block_end_uint   = uint64(block_end_int)
	}
	return block_start_uint, block_end_uint, nil
}

//-------------------------------------------------
func Http__get_arg__tx_id_hex(p_resp http.ResponseWriter,
	p_req         *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	qs_map := p_req.URL.Query()

	var tx_hex_str string
	if tx_lst, ok := qs_map["tx"]; ok {

		tx_hex_str = tx_lst[0]

		// IMPORTANT!! - TX hashes are 64 chars long, but have to contain
		//               the "0x" prefix, so the total is 66.
		if len(tx_hex_str) != 66 {
			gfErr := gf_core.ErrorCreate("supplied transaction ID is not 66 chars long",
				"verify__string_not_correct_length_error",
				map[string]interface{}{"tx_hex_str": tx_hex_str,},
				nil, "gf_eth_monitor_lib", pRuntimeSys)
			return "", gfErr
		}

		// check tax ID is prefixed with "0x"
		if !strings.HasPrefix(tx_hex_str, "0x") {
			gfErr := gf_core.ErrorCreate("supplied transaction ID is not prefixed with '0x'",
				"verify__invalid_value_error",
				map[string]interface{}{"tx_hex_str": tx_hex_str,},
				nil, "gf_eth_monitor_lib", pRuntimeSys)
			return "", gfErr
		}

		// check is Hex
		tx_hex_clean_str := strings.TrimPrefix(tx_hex_str, "0x")
		_, err := hex.DecodeString(tx_hex_clean_str)
		if err != nil {
			gfErr := gf_core.ErrorCreate("supplied transaction ID is not a hex string",
				"decode_hex",
				map[string]interface{}{"tx_hex_str": tx_hex_str,},
				nil, "gf_eth_monitor_lib", pRuntimeSys)
			return "", gfErr
		}

	}
	return tx_hex_str, nil
}

//-------------------------------------------------
func Http__get_arg__block_num(p_resp http.ResponseWriter,
	p_req         *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (uint64, *gf_core.GFerror) {
	
	qs_map := p_req.URL.Query()
	
	var block_num_int uint64
	if b_lst, ok := qs_map["b"]; ok {

		block_num_str := b_lst[0]

		i, err := strconv.Atoi(block_num_str)
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to parse input querystring arg as a an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"block_num": block_num_str,},
				err, "gf_eth_monitor_lib", pRuntimeSys)
			return 0, gfErr
		}
		block_num_int = uint64(i)
	}

	return block_num_int, nil
}

//-------------------------------------------------
func Http__get_arg__miner_addr(p_resp http.ResponseWriter,
	p_req         *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	miner_addr_str := ""
	
	return miner_addr_str, nil
}