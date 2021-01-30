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
	"strings"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/mitchellh/mapstructure"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// BLOCK__INTERNAL - internal representation of the block, with fields
//                   that are not visible to the external public users.
type GF_eth__block__int struct {
	Hash_str          string        `mapstructure:"hash_str"          json:"hash_str"`
	Parent_hash_str   string        `mapstructure:"parent_hash_str"   json:"parent_hash_str"`
	Block_num_int     uint64        `mapstructure:"block_num_int"     json:"block_num_int"`
	Gas_used_int      uint64        `mapstructure:"gas_used_int"      json:"gas_used_int"`
	Gas_limit_int     uint64        `mapstructure:"gas_limit_int"     json:"gas_limit_int"`
	Coinbase_addr_str string        `mapstructure:"coinbase_addr_str" json:"coinbase_addr_str"`
	Txs_lst           []*GF_eth__tx `mapstructure:"txs_lst"           json:"txs_lst"`
	Time              uint64        `mapstructure:"time_int"          json:"time_int"`
	Block             string        `mapstructure:"block"             json:"block"` // *eth_types.Block `json:"block"`
}

//-------------------------------------------------
func Eth_blocks__get_block_pipeline(p_block_int uint64,
	p_get_hosts_fn func() []string,
	
	p_ctx     context.Context,
	p_metrics *GF_metrics,
	p_runtime *GF_runtime) (map[string]*GF_eth__block__int, map[string]*GF_eth__miner__int, *gf_core.Gf_error) {

	worker_inspector__port_int := uint(2000)

	//---------------------
	// GET_WORKER_HOSTS
	span__get_worker_hosts := sentry.StartSpan(p_ctx, "get_worker_hosts")
	span__get_worker_hosts.SetTag("workers_aws_discovery", fmt.Sprint(p_runtime.Config.Workers_aws_discovery_bool))

	var workers_inspectors_hosts_lst []string
	if p_runtime.Config.Workers_aws_discovery_bool {
		workers_inspectors_hosts_lst = p_get_hosts_fn()
	} else {
		workers_inspectors_hosts_str := p_runtime.Config.Workers_hosts_str
		workers_inspectors_hosts_lst = strings.Split(workers_inspectors_hosts_str, ",")
	}



	span__get_worker_hosts.Finish()

	//---------------------
	// GET_BLOCKS__FROM_WORKERS_INSPECTORS__ALL

	span := sentry.StartSpan(p_ctx, "get_blocks__workers_inspectors__all")
	defer span.Finish()

	block_from_workers_map   := map[string]*GF_eth__block__int{}
	gf_errs_from_workers_map := map[string]*gf_core.Gf_error{}

	for _, host_str := range workers_inspectors_hosts_lst {

		ctx := span.Context()

		// GET_BLOCK__FROM_WORKER
		gf_block, gf_err := eth_blocks__worker_inspector__get_block(p_block_int,
			host_str,
			worker_inspector__port_int,
			ctx,
			p_runtime.Runtime_sys)

		if gf_err != nil {
			gf_errs_from_workers_map[host_str] = gf_err
			
			// mark a block coming from this worker_inspector host as nil,
			// and continue processing other hosts. 
			// a particular host may fail to return a particular block for various reasons,
			// it might not have synced to that block. 
			block_from_workers_map[host_str] = nil
			continue
		}





		for _, tx := range gf_block.Txs_lst {


			// this is a new_contract transaction
			if tx.Contract_new != nil {


				gf_err := Eth_contract__enrich(ctx, p_metrics, p_runtime)
				if gf_err != nil {
					return nil, nil, gf_err
				}



			}


		}






		block_from_workers_map[host_str] = gf_block




	}

	span.Finish()

	//---------------------
	// GET_MINERS - that own this address, potentially multiple records for the same address


	// get coinbase address from the block comming from the first worker_inspector
	var block_miner_addr_hex_str string
	for _, gf_block := range block_from_workers_map {
		
		// if worker failed to return a block, it will be set to nil, so go to the 
		// next one from which a coinbase could be acquired.
		if gf_block != nil {
			block_miner_addr_hex_str = gf_block.Coinbase_addr_str
			break
		}

	}

	miners_map, gf_err := Eth_miners__db__get_info(block_miner_addr_hex_str,
		p_metrics,
		p_ctx,
		p_runtime)
	if gf_err != nil {
		return nil, nil, gf_err
	}

	//---------------------

	return block_from_workers_map, miners_map, nil
}

//-------------------------------------------------
func eth_blocks__worker_inspector__get_block(p_block_int uint64,
	p_host_str    string,
	p_port_int    uint,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_eth__block__int, *gf_core.Gf_error) {


	


	url_str := fmt.Sprintf("http://%s:%d/gfethm_worker_inspect/v1/blocks?b=%d", p_host_str, p_port_int, p_block_int)

	//-----------------------
	


	// SPAN
	span_name_str    := fmt.Sprintf("worker_inspector__get_block:%s", p_host_str)
	span__get_blocks := sentry.StartSpan(p_ctx, span_name_str)
	
	// adding tracing ID as a header, to allow for distributed tracing, correlating transactions
	// across services.
	sentry_trace_id_str := span__get_blocks.ToSentryTrace()
	headers_map         := map[string]string{"sentry-trace": sentry_trace_id_str,}
		
	// GF_RPC_CLIENT
	data_map, gf_err := gf_rpc_lib.Client__request(url_str, headers_map, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	span__get_blocks.Finish()

	//-----------------------

	block_map := data_map["block_map"].(map[string]interface{})


	// DECODE_TO_STRUCT
	var gf_block GF_eth__block__int
	err := mapstructure.Decode(block_map, &gf_block)
	if err != nil {

		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to load response block_map into a GF_eth__block__int struct",
			"mapstruct__decode",
			map[string]interface{}{
				"url_str":   url_str,
				"block_map": block_map,
			},
			err, "gf_eth_monitor_core", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}


	return &gf_block, nil

}