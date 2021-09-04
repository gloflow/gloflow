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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_contract"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_blocks"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_indexer chan(GF_indexer_cmd)
type GF_indexer_cmd struct {
	Block_start_uint uint64
	Block_end_uint   uint64
	// Ctx context.Context
}

//-------------------------------------------------
func Init(p_get_worker_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) (chan(GF_indexer_cmd), *gf_core.GF_error) {

	ctx := context.Background()
	
	// ABI_DEFS
	abis_defs_map, gf_err := gf_eth_contract.Eth_abi__get_defs(ctx, p_metrics, p_runtime)
	if gf_err != nil {
		return nil, gf_err
	}

	indexer_cmds_ch := make(chan GF_indexer_cmd, 10)
	go func() {

		// IMPORTANT!! - using a background context, and not a client supplied context
		//               (via cmd.Ctx) because clients just submit an index operation,
		//               and continue their work (or get response to their request). 
		//               the index op should complete independently of the client, in the future.
		ctx := context.Background()

		for {
			select {

			//----------------------------
			// INDEXER_COMMANDS
			case cmd := <- indexer_cmds_ch:

				// PERSIST_RANGE
				gf_errs_lst := index__range(cmd.Block_start_uint,
					cmd.Block_end_uint,
					p_get_worker_hosts_fn,
					abis_defs_map,
					ctx, // cmd.Ctx,
					p_metrics,
					p_runtime)
				
				if len(gf_errs_lst) > 0 {

				}
			}

			//----------------------------
		}
	}()
	
	return indexer_cmds_ch, nil
}

//-------------------------------------------------
func index__range(p_block_start_uint uint64,
	p_block_end_uint      uint64,
	p_get_worker_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string,
	p_abis_defs_map       map[string]*gf_eth_contract.GF_eth__abi,
	p_ctx                 context.Context,
	p_metrics             *gf_eth_core.GF_metrics,
	p_runtime             *gf_eth_core.GF_runtime) []*gf_core.GF_error {

	gf_errs_lst := []*gf_core.GF_error{}
	for b := p_block_start_uint; b <= p_block_end_uint; b++ {

		block_uint := b

		


		gf_err := gf_eth_blocks.Index__pipeline(block_uint,
			p_get_worker_hosts_fn,
			p_abis_defs_map,
			p_ctx,
			p_metrics,
			p_runtime)
		if gf_err != nil {
			gf_errs_lst = append(gf_errs_lst, gf_err)
			continue // continue processing subsequent blocks
		}

	}
	return gf_errs_lst
}