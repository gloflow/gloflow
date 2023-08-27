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
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//-------------------------------------------------

func DBmongoABIget(p_abi_type_str string,
	p_ctx     context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) ([]*GF_eth__abi, *gf_core.GFerror) {




	if !Is_type_valid(p_abi_type_str) {
		error_defs_map := gf_eth_core.ErrorGetDefs()
		gf_err := gf_core.ErrorCreateWithDefs("supplied Eth contract to get an ABI from DB for is not valid",
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