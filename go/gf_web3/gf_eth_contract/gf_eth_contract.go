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

package gf_eth_contract

import (
	"fmt"
	"context"
	"math/big"
	// "encoding/base64"
	"github.com/ethereum/go-ethereum/ethclient"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//-------------------------------------------------
type GF_eth__contract_new struct {
	Addr_str       string `json:"addr_str"`
	Code_bytes_lst []byte `json:"-"` // in json serialization []byte is not included, just the base64 encoding
	Code_hex_str   string `json:"code_hex_str"`
	// Code_b64_str   string `json:"code_b64_str"`
	Block_num_int  uint64 `json:"block_num_int"`
}

//-------------------------------------------------
func Enrich(p_gf_abi *GF_eth__abi,
	p_ctx     context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) *gf_core.GFerror {

	// sometimes no ABI's are fetched from the DB
	if p_gf_abi != nil {

		// abi_type_str := "erc20"
		abi, gf_err := Eth_abi__get(p_gf_abi, p_ctx, p_metrics, p_runtime)
		if gf_err != nil {
			return gf_err
		}

		fmt.Println(abi)
	}

	return nil
}

//-------------------------------------------------
func Get_via_rpc(p_contract_addr_str string,
	p_block_num_int  uint64,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	pRuntimeSys    *gf_core.RuntimeSys) (*GF_eth__contract_new, *gf_core.GFerror) {

	code_bytes_lst, gf_err := Get_code(p_contract_addr_str,
		p_block_num_int,
		p_ctx,
		p_eth_rpc_client,
		pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}


	// // base64
	// code_b64_str := base64.StdEncoding.EncodeToString(code_bytes_lst)


	code_hex_str := eth_common.BytesToHash(code_bytes_lst).Hex()

	contract__new := &GF_eth__contract_new{
		Addr_str:       p_contract_addr_str,
		Code_bytes_lst: code_bytes_lst,
		Code_hex_str:   code_hex_str,
		// Code_b64_str:   code_b64_str,
		Block_num_int:  p_block_num_int,
	}

	return contract__new, gf_err
}

//-------------------------------------------------
func Get_code(p_contract_addr_str string,
	p_block_num_int  uint64,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	pRuntimeSys    *gf_core.RuntimeSys) ([]byte, *gf_core.GFerror) {

	contract_addr := eth_common.HexToAddress(p_contract_addr_str)

	// CODE_AT
	code_bytes_lst, err := p_eth_rpc_client.CodeAt(p_ctx,
		contract_addr,
		big.NewInt(0).SetUint64(p_block_num_int))
		
	if err != nil {
		error_defs_map := gf_eth_core.ErrorGetDefs()
		gf_err := gf_core.ErrorCreateWithDefs("failed to get code at particular account address in target block",
			"eth_rpc__get_contract_code",
			map[string]interface{}{"contract_addr_str": p_contract_addr_str, "block_num_int": p_block_num_int,},
			err, "gf_eth_monitor_core", error_defs_map, 1, pRuntimeSys)
		return nil, gf_err
	}

	return code_bytes_lst, nil
}

//-------------------------------------------------
func Is_type_valid(p_type_str string) bool {
	types_map := map[string]bool{
		"erc20": true,
	}
	if _, ok := types_map[p_type_str]; ok {
		return true
	}
	return false
}