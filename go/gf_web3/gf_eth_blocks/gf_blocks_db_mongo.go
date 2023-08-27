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
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//-------------------------------------------------
// BLOCKS__DB__GET_COUNT

func DBmongoGetCount(pMetrics *gf_eth_core.GF_metrics,
	pRuntime *gf_eth_core.GF_runtime) (int64, *gf_core.GFerror) {

	collNameStr := "gf_eth_blocks"
	coll := pRuntime.RuntimeSys.Mongo_db.Collection(collNameStr)

	ctx := context.Background()
	
	countInt, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {

		// METRICS
		if pMetrics != nil {pMetrics.Errs_num__counter.Inc()}

		gfErr := gf_core.MongoHandleError("failed to DB count Blocks",
			"mongodb_count_error",
			map[string]interface{}{},
			err, "gf_eth_monitor_core", pRuntime.RuntimeSys)
		return 0, gfErr
	}

	return countInt, nil
}

//-------------------------------------------------
// BLOCKS__DB__WRITE_BULK

func DBmongoWriteBulk(p_gf_blocks_lst []*GF_eth__block__int,
	pCtx     context.Context,
	pMetrics *gf_eth_core.GF_metrics,
	pRuntime *gf_eth_core.GF_runtime) *gf_core.GFerror {

	filterDocsByFieldsLst := []map[string]string{}
	recordsLst            := []interface{}{}
	blocksNumsLst         := []uint64{}

	for _, b := range p_gf_blocks_lst {

		filterDocsByFieldsLst = append(filterDocsByFieldsLst,
			map[string]string{"_id": b.DB_id,})

		recordsLst    = append(recordsLst, interface{}(b))
		blocksNumsLst = append(blocksNumsLst, b.Block_num_uint)
	}

	collNameStr := "gf_eth_blocks"
	_, gfErr := gf_core.MongoUpsertBulk(filterDocsByFieldsLst, recordsLst,
		collNameStr,
		map[string]interface{}{
			"blocks_nums_lst":    blocksNumsLst,
			"caller_err_msg_str": "failed to bulk insert Eth blocks (GF_eth__block__int) into DB",
		},
		pCtx,
		pRuntime.RuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}