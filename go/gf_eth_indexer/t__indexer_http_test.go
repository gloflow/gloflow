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

package gf_eth_indexer

import (
	"fmt"
	"os"
	"time"
	"testing"
	"context"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__indexer_http(p_test *testing.T) {

	fmt.Println("TEST__INDEXER_HTTP ==============================================")
	


	start_block_int     := uint64(13_414_401)
	end_block_int       := uint64(13_414_402)
	test_blocks_num_int := end_block_int - start_block_int
	worker__host_port_str := os.Getenv("GF_TEST_WORKER_INSPECTOR_HOST_PORT")
	ctx        := context.Background()
	runtime, _ := gf_eth_core.TgetRuntime(p_test)


	port_int := 3000
	host_port_str := fmt.Sprintf("localhost:%d", port_int)

	// INIT_TEST_SERVER
	init_done_ch := make(chan bool)
	go func() {

		indexer_cmds_ch, indexer_job_updates_new_consumer_ch, gf_err := Init(func(p_ctx context.Context, p_runtime *gf_eth_core.GF_runtime) []string {
				return []string{worker__host_port_str,}
			},
			nil,
			runtime)
		if gf_err != nil {
			p_test.Fail()
			return
		}

		Init_handlers(indexer_cmds_ch,
			indexer_job_updates_new_consumer_ch,
			nil,
			runtime)

		init_done_ch <- true
		gf_rpc_lib.Server__init(port_int)
		
	}()

	<- init_done_ch
	time.Sleep(1 * time.Second) // give server time to startup

	//---------------------
	// CLIENT_HTTP - START_INDEXING
	fmt.Println("TEST - START INDEXING...")

	job_id_str, gf_err := Client_http__index_block_range(start_block_int, end_block_int,
		host_port_str,
		ctx,
		runtime.Runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	//---------------------
	fmt.Println("TEST - START LISTENING FOR INDEXING UPDATES...")

	job_updates_ch := make(chan map[string]interface{}, 10)
	gf_err = Client_http__index_job_updates(job_id_str,
		job_updates_ch,
		host_port_str,
		ctx,
		runtime.Runtime_sys)

	if gf_err != nil {
		p_test.Fail()
	}

	received_job_updates_num_int := uint64(0)
	for {

		select {
		case update_map := <-job_updates_ch:
			spew.Dump(update_map)

			received_job_updates_num_int += 1
			if received_job_updates_num_int == test_blocks_num_int {

				// end test
				return

			} 
		}
	}

	//---------------------
}