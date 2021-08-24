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

package gf_eth_blocks

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_core"
)

//-------------------------------------------------
// BLOCKS__DB__GET_COUNT
func DB__get_count(p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) (int64, *gf_core.GF_error) {

	coll_name_str := "gf_eth_blocks"
	coll := p_runtime.Runtime_sys.Mongo_db.Collection(coll_name_str)

	ctx := context.Background()
	
	count_int, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {

		// METRICS
		if p_metrics != nil {p_metrics.Errs_num__counter.Inc()}

		gf_err := gf_core.Mongo__handle_error("failed to DB count Blocks",
			"mongodb_count_error",
			map[string]interface{}{},
			err, "gf_eth_monitor_core", p_runtime.Runtime_sys)
		return 0, gf_err
	}

	return count_int, nil
}

//-------------------------------------------------
// BLOCKS__DB__WRITE_BULK
func DB__write_bulk(p_gf_blocks_lst []*GF_eth__block__int,
	p_ctx     context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) *gf_core.GF_error {

	ids_lst         := []string{}
	records_lst     := []interface{}{}
	blocks_nums_lst := []uint64{}
	for _, b := range p_gf_blocks_lst {
		ids_lst         = append(ids_lst, b.DB_id)
		records_lst     = append(records_lst, interface{}(b))
		blocks_nums_lst = append(blocks_nums_lst, b.Block_num_uint)
	}

	coll_name_str := "gf_eth_blocks"
	gf_err := gf_core.Mongo__insert_bulk(ids_lst, records_lst,
		coll_name_str,
		map[string]interface{}{
			"blocks_nums_lst":    blocks_nums_lst,
			"caller_err_msg_str": "failed to bulk insert Eth blocks (GF_eth__block__int) into DB",
		},
		p_ctx,
		p_runtime.Runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}