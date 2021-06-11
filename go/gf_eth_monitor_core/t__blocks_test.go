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
	"os"
	"fmt"
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	// "github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_stats/gf_stats_lib"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__blocks__get_and_persist_bulk__pipeline(p_test *testing.T) {

	fmt.Println("TEST__BLOCKS ==============================================")
	


	ctx := context.Background()

	host_port_str := os.Getenv("GF_TEST_WORKER_INSPECTOR_HOST_PORT")

	runtime, metrics := t__get_runtime(p_test)




	block_start_uint := uint64(2_000_000)
	block_end_uint   := uint64(2_001_000)

	get_worker_hosts_fn := func(p_ctx context.Context, p_runtime *GF_runtime) []string {
		return []string{host_port_str, }
	}
	gf_errs_lst := Eth_blocks__get_and_persist_bulk__pipeline(block_start_uint,
		block_end_uint,
		get_worker_hosts_fn,
		ctx,
		metrics,
		runtime)

	if len(gf_errs_lst) > 0 {
		p_test.Fail()
	}


	
	txs__to_test_map := map[string]uint64{
		"0xc55e2b90168af6972193c1f86fa4d7d7b31a29c156665d15b9cd48618b5177ef": 2_000_000,
		"0x0a819ec79aa1ce1cb1d408f69c6ac6b4af187ac8cd7f094532e278d0848ddba3": 2_000_002,
		"0x00af51187daefca9e0a7ff9eee7fff2fde30eb1a449c7682288a524f36df3f01": 2_000_002,
	}


	for tx_hash_str, block_num_int := range txs__to_test_map {
		gf_tx, gf_err := eth_tx__db__get(tx_hash_str, ctx, metrics, runtime)
		if gf_err != nil {
			p_test.Fail()
		}


		spew.Dump(gf_tx)
		

		assert.EqualValues(p_test, gf_tx.Block_num_int, block_num_int,
			"test TX fetched from DB doesnt have the same block number is the specified test block that contains it")
		
	}




	fmt.Println("+++++++++++++++++++++++++++++++")
	db_coll_stats, gf_err := gf_stats_lib.Db_stats__coll("gf_eth_txs_traces", ctx, runtime.Runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}
	spew.Dump(db_coll_stats)


}