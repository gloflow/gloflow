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

package gf_eth_monitor_lib

import (
	"strconv"
	"net/http"
	"encoding/hex"
	// "github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func Http__get_arg__tx_id_hex(p_resp http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (string, *gf_core.Gf_error) {

	qs_map := p_req.URL.Query()

	var tx_hex_str string
	if tx_lst, ok := qs_map["tx"]; ok {

		tx_hex_str := tx_lst[0]

		// IMPORTANT!! - TX hashes are 64 chars long, but have to contain
		//               the "0x" prefix, so the total is 66.
		if len(tx_hex_str) != 66 {
			gf_err := gf_core.Error__create("supplied transaction ID is not 66 chars long",
				"verify__string_not_correct_length_error",
				map[string]interface{}{"tx_hex_str": tx_hex_str,},
				nil, "gf_eth_monitor_lib", p_runtime_sys)
			return "", gf_err
		}

		_, err := hex.DecodeString(tx_hex_str)
		if err != nil {
			gf_err := gf_core.Error__create("supplied transaction ID is not a hex string",
				"decode_hex",
				map[string]interface{}{"tx_hex_str": tx_hex_str,},
				nil, "gf_eth_monitor_lib", p_runtime_sys)
			return "", gf_err
		}

	}
	return tx_hex_str, nil
}

//-------------------------------------------------
func Http__get_arg__block_num(p_resp http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (uint64, *gf_core.Gf_error) {

	qs_map := p_req.URL.Query()

	var block_num_int uint64
	if b_lst, ok := qs_map["b"]; ok {
		block_num_str := b_lst[0]

		i, err := strconv.Atoi(block_num_str)
		if err != nil {
		
			gf_err := gf_core.Error__create("failed to parse input querystring arg as a an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"block_num": block_num_str,},
				err, "gf_eth_monitor_lib", p_runtime_sys)

			return 0, gf_err 
		}
		block_num_int = uint64(i)
	}

	return block_num_int, nil
}

//-------------------------------------------------
func Http__get_arg__miner_addr(p_resp http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (string, *gf_core.Gf_error) {





	miner_addr_str := ""
	return miner_addr_str, nil



}