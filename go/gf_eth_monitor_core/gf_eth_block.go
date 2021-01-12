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
	// eth_types "github.com/ethereum/go-ethereum/core/types"
)

//-------------------------------------------------
// BLOCK__INTERNAL - internal representation of the block, with fields
//                   that are not visible to the external public users.
type GF_eth__block__int struct {
	Block_num_int     uint64        `mapstructure:"block_num_int"     json:"block_num_int"`
	Gas_used_int      uint64        `mapstructure:"gas_used_int"      json:"gas_used_int"`
	Gas_limit_int     uint64        `mapstructure:"gas_limit_int"     json:"gas_limit_int"`
	Coinbase_addr_str string        `mapstructure:"coinbase_addr_str" json:"coinbase_addr_str"`
	Txs_lst           []*GF_eth__tx `mapstructure:"txs_lst"           json:"txs_lst"`
	Block             string        `mapstructure:"block"             json:"block"` // *eth_types.Block `json:"block"`
}

//-------------------------------------------------
func Eth_block__get_block_pipeline(p_block_int uint64,
	p_get_hosts_fn func() []string,
	p_ctx          context.Context,
	p_runtime      *GF_runtime) (map[string]*GF_eth__block__int, *gf_core.Gf_error) {

	

	worker_inspector__port_int := uint(2000)

	//---------------------

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




	block_from_workers_map := map[string]*GF_eth__block__int{}
	for _, host_str := range workers_inspectors_hosts_lst {

		
		// SPAN
		span_name_str    := fmt.Sprintf("worker_inspector__get_block:%s", host_str)
		span__get_blocks := sentry.StartSpan(p_ctx, span_name_str)

		// GET_BLOCK
		gf_block, gf_err := eth_block__worker_inspector__get_block(p_block_int,
			host_str,
			worker_inspector__port_int,
			p_runtime.Runtime_sys)

		span__get_blocks.Finish()

		if gf_err != nil {
			return nil, gf_err
		}



		block_from_workers_map[host_str] = gf_block


	}

	return block_from_workers_map, nil
}

//-------------------------------------------------
func eth_block__worker_inspector__get_block(p_block_int uint64,
	p_host_str    string,
	p_port_int    uint,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_eth__block__int, *gf_core.Gf_error) {


	


	url_str := fmt.Sprintf("http://%s:%d/gfethm_worker_inspect/v1/blocks?block=%s", p_host_str, p_port_int, p_block_int)

	data_map, gf_err := gf_rpc_lib.Client__request(url_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	fmt.Println(data_map)



	// DECODE_TO_STRUCT
	var gf_block GF_eth__block__int
	mapstructure.Decode(data_map, &gf_block)


	

	return &gf_block, nil

}