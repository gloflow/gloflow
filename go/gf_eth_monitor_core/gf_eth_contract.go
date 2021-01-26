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
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum/ethclient"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_eth__contract_new struct {
	Addr_str       string  `json:"addr_str"`
	Code_bytes_lst []byte  `json:"code_bytes_lst"`
	Block_num_int   uint64 `json:"block_num_int"`
}

//-------------------------------------------------
func Eth_contract__get_code(p_contract_addr_str string,
	p_block_num_int  uint64,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	p_runtime_sys    *gf_core.Runtime_sys) ([]byte, *gf_core.Gf_error) {





	contract_addr := eth_common.HexToAddress(p_contract_addr_str)
	code_bytes_lst, err := p_eth_rpc_client.CodeAt(p_ctx,
		contract_addr,
		big.NewInt(0).SetUint64(p_block_num_int))
		
	if err != nil {
		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get code at particular account address in target block",
			"eth_rpc__get_contract_code",
			map[string]interface{}{"contract_addr_str": p_contract_addr_str, "block_num_int": p_block_num_int,},
			err, "gf_eth_monitor_core", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}



	return code_bytes_lst, nil



	


}