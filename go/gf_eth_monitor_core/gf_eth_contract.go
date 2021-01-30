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
	"fmt"
	"context"
	"math/big"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/ethereum/go-ethereum/ethclient"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_eth__contract_new struct {
	Addr_str       string `json:"addr_str"`
	Code_bytes_lst []byte `json:"-"` // in json serialization []byte is not included, just the base64 encoding
	Code_b64_str   string `json:"code_b64_str"`
	Block_num_int  uint64 `json:"block_num_int"`
}

type GF_eth__abi struct {
	Type_str string                   `bson:"type_str"`
	Def_lst  []map[string]interface{} `bson:"def_lst"`
}

//-------------------------------------------------
func Eth_contract__enrich(p_ctx context.Context,
	p_metrics *GF_metrics,
	p_runtime *GF_runtime) *gf_core.Gf_error {


	abi_type_str := ""
	abis_lst, gf_err := Eth_contract__db__get_abi(abi_type_str, p_ctx, p_metrics, p_runtime)
	if gf_err != nil {
		


		return gf_err
	}



	fmt.Println(abis_lst)


	return nil

}



//-------------------------------------------------

func Eth_contract__db__get_abi(p_abi_type_str string,
	p_ctx     context.Context,
	p_metrics *GF_metrics,
	p_runtime *GF_runtime) ([]*GF_eth__abi, *gf_core.Gf_error) {




	if !Eth_contract__is_type_valid(p_abi_type_str) {
		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("supplied Eth contract to get an ABI from DB for is not valid",
			"eth_contract__not_supported_type",
			map[string]interface{}{"type_str": p_abi_type_str,},
			nil, "gf_eth_monitor_core", error_defs_map, p_runtime.Runtime_sys)
		return nil, gf_err
	}




	coll_name_str := "gf_eth_meta__contracts_abi"

	q := bson.M{"type_str": p_abi_type_str, }

	cur, err := p_runtime.Mongodb_db.Collection(coll_name_str).Find(p_ctx, q)
	if err != nil {




		// METRICS
		if p_metrics != nil {
			p_metrics.Counter__errs_num.Inc()
		}

		gf_err := gf_core.Mongo__handle_error("failed to find Miner with gives address in DB",
			"mongodb_find_error",
			map[string]interface{}{"type_str": p_abi_type_str,},
			err, "gf_eth_monitor_core", p_runtime.Runtime_sys)
		return nil, gf_err
	}
	defer cur.Close(p_ctx)


	abis_lst := []*GF_eth__abi{}
	for cur.Next(p_ctx) {

		var gf_abi GF_eth__abi
		err := cur.Decode(&gf_abi)
		if err != nil {
			gf_err := gf_core.Mongo__handle_error("failed to decode mongodb result of query to get contract ABIs",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_eth_monitor_core", p_runtime.Runtime_sys)



			return nil, gf_err
		}
	
		abis_lst = append(abis_lst, &gf_abi)
	}

	return abis_lst, nil
}

//-------------------------------------------------
func Eth_contract__get_via_rpc(p_contract_addr_str string,
	p_block_num_int  uint64,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	p_runtime_sys    *gf_core.Runtime_sys) (*GF_eth__contract_new, *gf_core.Gf_error) {

	code_bytes_lst, gf_err := Eth_contract__get_code(p_contract_addr_str,
		p_block_num_int,
		p_ctx,
		p_eth_rpc_client,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}


	// base64
	code_b64_str := base64.StdEncoding.EncodeToString(code_bytes_lst)


	contract__new := &GF_eth__contract_new{
		Addr_str:       p_contract_addr_str,
		Code_bytes_lst: code_bytes_lst,
		Code_b64_str:   code_b64_str,
		Block_num_int:  p_block_num_int,
	}

	return contract__new, gf_err
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

//-------------------------------------------------

func Eth_contract__is_type_valid(p_type_str string) bool {
	types_map := map[string]bool{
		"erc20": true,
	}
	if _, ok := types_map[p_type_str]; ok {
		return true
	}
	return false
}