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

package gf_eth_tx

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_core"
)

//-------------------------------------------------
// DB__GET_COUNT
func DB__get_count(p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) (uint64, uint64, *gf_core.GF_error) {

	//-------------------------------------------------
	count_fn := func(p_coll_name_str string) (uint64, *gf_core.GF_error) {
		coll := p_runtime.Runtime_sys.Mongo_db.Collection(p_coll_name_str)

		ctx := context.Background()
		
		count_int, err := coll.CountDocuments(ctx, bson.M{})
		if err != nil {

			// METRICS
			if p_metrics != nil {p_metrics.Errs_num__counter.Inc()}

			gf_err := gf_core.Mongo__handle_error("failed to DB count Transactions",
				"mongodb_count_error",
				map[string]interface{}{},
				err, "gf_eth_monitor_core", p_runtime.Runtime_sys)
			return 0, gf_err
		}
		return uint64(count_int), nil
	}

	//-------------------------------------------------
	txs_count_int, gf_err := count_fn("gf_eth_txs")
	if gf_err != nil {
		return 0, 0, gf_err
	}
	txs_traces_count_int, gf_err := count_fn("gf_eth_txs_traces")
	if gf_err != nil {
		return 0, 0, gf_err
	}

	return txs_count_int, txs_traces_count_int, nil
}

//-------------------------------------------------
// DB__GET
func DB__get(p_tx_hash_str string,
	p_ctx     context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) (*GF_eth__tx, *gf_core.GF_error) {

	coll_name_str := "gf_eth_txs"


	q := bson.M{"hash_str": p_tx_hash_str, }

	var gf_tx GF_eth__tx
	err := p_runtime.Runtime_sys.Mongo_db.Collection(coll_name_str).FindOne(p_ctx, q).Decode(&gf_tx)
	if err != nil {


		// METRICS
		if p_metrics != nil {
			p_metrics.Errs_num__counter.Inc()
		}

		gf_err := gf_core.Mongo__handle_error("failed to find Transaction with gives hash in DB",
			"mongodb_find_error",
			map[string]interface{}{"tx_hash_str": p_tx_hash_str,},
			err, "gf_eth_monitor_core", p_runtime.Runtime_sys)
		return nil, gf_err
	}

	return &gf_tx, nil
}

//-------------------------------------------------
// DB__WRITE_BULK
func DB__write_bulk(p_txs_lst []*GF_eth__tx,
	p_ctx     context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) *gf_core.GF_error {

	coll_name_str := "gf_eth_txs"

	ids_lst        := []string{}
	records_lst    := []interface{}{}
	txs_hashes_lst := []string{}
	for _, tx := range p_txs_lst {
		ids_lst        = append(ids_lst, tx.DB_id)
		records_lst    = append(records_lst, interface{}(tx))
		txs_hashes_lst = append(txs_hashes_lst, tx.Hash_str)
	}

	gf_err := gf_core.Mongo__insert_bulk(ids_lst, records_lst,
		coll_name_str,
		map[string]interface{}{
			"txs_hashes_lst":     txs_hashes_lst,
			"caller_err_msg_str": "failed to bulk insert Eth txs (GF_eth__tx) into DB",
		},
		p_ctx, p_runtime.Runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}