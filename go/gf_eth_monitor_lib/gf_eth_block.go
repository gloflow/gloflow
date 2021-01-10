/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_eth_monitor_lib

import (
	"fmt"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
	// eth_types "github.com/ethereum/go-ethereum/core/types"
)

//-------------------------------------------------
// BLOCK__INTERNAL - internal representation of the block, with fields
//                   that are not visible to the external public users.
type GF_eth__block__int struct {
	Block_num_int     uint64        `json:"block_num_int"`
	Gas_used_int      uint64        `json:"gas_used_int"`
	Gas_limit_int     uint64        `json:"gas_limit_int"`
	Coinbase_addr_str string        `json:"coinbase_addr_str"`
	Txs_lst           []*GF_eth__tx `json:"txs_lst"`
	Block             string        `json:"block"` // *eth_types.Block `json:"block"`
}

//-------------------------------------------------
func eth_block__get_block_pipeline(p_block_int uint64,
	p_runtime *GF_runtime) *gf_core.Gf_error {



	workers_inspectors_hosts_lst := strings.Split(p_runtime.Config.Workers_inspectors_hosts_str, ",")



	for _, host_str := range workers_inspectors_hosts_lst {



		gf_err := eth_block__worker_inspector__get_block(p_block_int, host_str, p_runtime)
		if gf_err != nil {
			return gf_err
		}
	}

	return nil
}

//-------------------------------------------------
func eth_block__worker_inspector__get_block(p_block_int uint64,
	p_host_str string,
	p_runtime  *GF_runtime) *gf_core.Gf_error {





	url_str := fmt.Sprintf("http://%s/gfethm_worker_inspect/v1/blocks?block=%s", p_host_str, p_block_int)

	gf_http_fetch, gf_err := gf_core.HTTP__fetch_url(url_str, p_runtime.Runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	fmt.Println(gf_http_fetch)





	return nil

}