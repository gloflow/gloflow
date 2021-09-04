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
)

//-------------------------------------------------
// INDEX_BLOCK_RANGE
func Client__index_block_range(p_block_start_uint uint64,
	p_block_end_uint  uint64,
	p_ctx             context.Context,
	p_indexer_cmds_ch chan(GF_indexer_cmd)) {






	cmd := GF_indexer_cmd{
		Block_start_uint: p_block_start_uint,
		Block_end_uint:   p_block_end_uint,
		// Ctx: p_ctx,
	}

	p_indexer_cmds_ch <- cmd
}