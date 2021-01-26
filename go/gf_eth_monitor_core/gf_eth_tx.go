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
	eth_types "github.com/ethereum/go-ethereum/core/types"
)

//-------------------------------------------------
type GF_eth__tx struct {
	Hash_str              string           `json:"hash_str"       bson:"hash_str"`
	Index_int             uint64           `json:"index_int"      bson:"index_int"` // position of the transaction in the block
	From_addr_str         string           `json:"from_addr_str"  bson:"from_addr_str"`
	To_addr_str           string           `json:"to_addr_str"    bson:"to_addr_str"`
	Value_eth_f           float64          `json:"value_eth_f"    bson:"value_eth_f"`
	Data_bytes_lst        []byte           `json:"data_bytes_lst" bson:"data_bytes_lst"`
	Gas_used_int          uint64           `json:"gas_used_int"   bson:"gas_used_int"`
	Gas_price_int         uint64           `json:"gas_price_int"  bson:"gas_price_int"`
	Nonce_int             uint64           `json:"nonce_int"      bson:"nonce_int"`
	Size_f                float64          `json:"size_f"         bson:"size_f"`
	Cost_int              uint64           `json:"cost_int"       bson:"cost_int"`
	Contract_new   *GF_eth__contract_new `json:"contract_new_map" bson:"contract_new_map"`
	Logs_lst       []*eth_types.Log      `json:"logs_lst"         bson:"logs_lst"`
}

//-------------------------------------------------






