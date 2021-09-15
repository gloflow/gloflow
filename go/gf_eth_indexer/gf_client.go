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
	p_indexer_cmds_ch chan(GF_indexer_cmd)) GF_indexer_job_id {

	response_ch := make(chan GF_indexer_job_id) 
	cmd := GF_indexer_cmd{
		Block_start_uint: p_block_start_uint,
		Block_end_uint:   p_block_end_uint,
		Response_ch:      response_ch,
	}

	p_indexer_cmds_ch <- cmd
	job_id_str := <- response_ch
	
	return job_id_str
}

//-------------------------------------------------
func Client__new_consumer(p_job_id_str GF_indexer_job_id,
	p_indexer_job_updates_new_consumer_ch GF_job_update_new_consumer_ch,
	p_ctx                                 context.Context,) (GF_job_updates_ch, GF_job_complete_ch) {



	response_ch := make(chan GF_job_update_new_consumer_response)
	new_consumer := GF_job_update_new_consumer{
		Job_id_str:  p_job_id_str,
		Response_ch: response_ch,
		ctx:         p_ctx,
	}

	p_indexer_job_updates_new_consumer_ch <- new_consumer
	response := <- response_ch

	return response.Job_updates_ch, response.Job_complete_ch
}