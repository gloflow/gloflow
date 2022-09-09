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
	"context"
	"strings"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//-------------------------------------------------
type GF_eth__abi struct {
	Type_str string                   `bson:"type_str"`
	Def_lst  []map[string]interface{} `bson:"def_lst"`
}

//-------------------------------------------------
func Eth_abi__get_defs(p_ctx context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) (map[string]*GF_eth__abi, *gf_core.GF_error) {



	abis_defs_map := map[string]*GF_eth__abi{}


	// DB_GET
	abi_type_str := "erc20"
	abis_lst, gf_err := Eth_abi__db__get(abi_type_str, p_ctx, p_metrics, p_runtime)
	if gf_err != nil {
		return nil, gf_err
	}

	if len(abis_lst) > 0 {
		abis_defs_map["erc20"] = abis_lst[0]
	}

	return abis_defs_map, nil
}

//-------------------------------------------------
func Eth_abi__get(p_gf_abi *GF_eth__abi,
	p_ctx     context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) (*abi.ABI, *gf_core.GF_error) {

	//---------------------
	/*// VALIDATE
	if p_abi_type_str != "erc20" {
		error_defs_map := Error__get_defs()
		gf_err := gf_core.ErrorCreate_with_defs("ABI type is not supported",
			"eth_contract__abi_not_loadable",
			map[string]interface{}{
				"abi_type_str": p_abi_type_str,
			},
			nil, "gf_eth_monitor_core", error_defs_map, p_runtime.RuntimeSys)
		return nil, gf_err
	}*/

	//---------------------
	// LOAD
	
	abi_def_lst := p_gf_abi.Def_lst // abis_lst[0].Def_lst
	abi_def_str, _ := json.Marshal(abi_def_lst)

	abi, err := abi.JSON(strings.NewReader(string(abi_def_str)))
	if err != nil {
		error_defs_map := gf_eth_core.Error__get_defs()
		gf_err := gf_core.ErrorCreate_with_defs("cant load ABI JSON whos definition was loaded from DB",
			"eth_contract__abi_not_loadable",
			map[string]interface{}{
				"abi_type_str": p_gf_abi.Type_str,
				"abi_def_str":  abi_def_str,
			},
			err, "gf_eth_monitor_core", error_defs_map, 1, p_runtime.RuntimeSys)
		return nil, gf_err
	}

	//---------------------

	return &abi, nil
}

//-------------------------------------------------
func Eth_abi__db__get(p_abi_type_str string,
	p_ctx     context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) ([]*GF_eth__abi, *gf_core.GF_error) {




	if !Is_type_valid(p_abi_type_str) {
		error_defs_map := gf_eth_core.Error__get_defs()
		gf_err := gf_core.ErrorCreate_with_defs("supplied Eth contract to get an ABI from DB for is not valid",
			"eth_contract__not_supported_type",
			map[string]interface{}{"type_str": p_abi_type_str,},
			nil, "gf_eth_monitor_core", error_defs_map, 1, p_runtime.RuntimeSys)
		return nil, gf_err
	}




	coll_name_str := "gf_eth_meta__contracts_abi"

	q := bson.M{"type_str": p_abi_type_str, }

	cur, err := p_runtime.RuntimeSys.Mongo_db.Collection(coll_name_str).Find(p_ctx, q)
	if err != nil {

		// METRICS
		if p_metrics != nil {
			p_metrics.Errs_num__counter.Inc()
		}

		gf_err := gf_core.MongoHandleError("failed to find Contract ABI with given type in DB",
			"mongodb_find_error",
			map[string]interface{}{"type_str": p_abi_type_str,},
			err, "gf_eth_monitor_core", p_runtime.RuntimeSys)
		return nil, gf_err
	}
	defer cur.Close(p_ctx)


	abis_lst := []*GF_eth__abi{}
	for cur.Next(p_ctx) {

		var gf_abi GF_eth__abi
		err := cur.Decode(&gf_abi)
		if err != nil {
			gf_err := gf_core.MongoHandleError("failed to decode mongodb result of query to get contract ABIs",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_eth_monitor_core", p_runtime.RuntimeSys)

			return nil, gf_err
		}
	
		abis_lst = append(abis_lst, &gf_abi)
	}

	return abis_lst, nil
}