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

package gf_eth_core

import (
	"context"
	"strings"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_eth__miner__int struct {
	Name_str        string `bson:"name_str" json:"name_str"`
	Address_hex_str string `bson:"addr_str" json:"address_hex_str"`
}

//-------------------------------------------------
// DB__GET_INFO
func Eth_miners__db__get_info(p_miner_address_str string,
	p_metrics *GF_metrics,
	p_ctx     context.Context,
	p_runtime *GF_runtime) (map[string]*GF_eth__miner__int, *gf_core.GF_error) {

	coll_name_str := "gf_eth_meta__miners"

	miner_address_lower_str := strings.ToLower(p_miner_address_str)
	q := bson.M{"addr_str": miner_address_lower_str, }

	cur, err := p_runtime.RuntimeSys.Mongo_db.Collection(coll_name_str).Find(p_ctx, q)
	if err != nil {

		// METRICS
		if p_metrics != nil {p_metrics.Errs_num__counter.Inc()}

		gf_err := gf_core.Mongo__handle_error("failed to find Miner with gives address in DB",
			"mongodb_find_error",
			map[string]interface{}{"miner_addr_str": miner_address_lower_str,},
			err, "gf_eth_core", p_runtime.RuntimeSys)
		return nil, gf_err
	}
	defer cur.Close(p_ctx)


	miners_map := map[string]*GF_eth__miner__int{}
	for cur.Next(p_ctx) {

		var miner GF_eth__miner__int
		err := cur.Decode(&miner)
		if err != nil {
			gf_err := gf_core.Mongo__handle_error("failed to decode mongodb result of query to get Miners",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_eth_core", p_runtime.RuntimeSys)
				
			return nil, gf_err
		}
	
		miners_map[miner.Name_str] = &miner
	}

	return miners_map, nil
}