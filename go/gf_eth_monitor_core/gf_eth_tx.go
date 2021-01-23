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
	Hash_str      string           `json:"hash_str"`
	Index_int     uint             `json:"index_int"` // position of the transaction in the block
	Gas_used_int  uint64           `json:"gas_used_int"`
	Nonce_int     uint64           `json:"nonce_int"`
	Size_f        float64          `json:"size_f"`
	To_addr_str   string           `json:"to_addr_str"`
	Value_int     int64            `json:"value_int"`
	Gas_price_int int64            `json:"gas_price_int"`
	Cost_int      int64            `json:"cost_int"`
	Logs          []*eth_types.Log `json:"logs_lst"`
}

//-------------------------------------------------






