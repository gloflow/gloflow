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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
// CLIENT__INDEX_BLOCK_RANGE
func Client__index_block_range(p_block_start_uint uint64,
	p_block_end_uint  uint64,
	p_indexer_cmds_ch chan(GF_indexer_cmd)) (GF_indexer_job_id, *gf_core.GFerror) {

	response_ch     := make(chan GF_indexer_job_id, 1)
	response_err_ch := make(chan gf_core.GFerror, 1)
	
	cmd := GF_indexer_cmd{
		Block_start_uint: p_block_start_uint,
		Block_end_uint:   p_block_end_uint,
		Response_ch:      response_ch,
		Response_err_ch:  response_err_ch,
	}

	p_indexer_cmds_ch <- cmd

	select {
	case job_id_str := <- response_ch:
		return job_id_str, nil
	case gfErr := <- response_err_ch:
		return GF_indexer_job_id(""), &gfErr
	}

	return GF_indexer_job_id(""), nil
}

//-------------------------------------------------
// CLIENT__NEW_CONSUMER
func Client__new_consumer(p_job_id_str GF_indexer_job_id,
	p_indexer_job_updates_new_consumer_ch GF_job_update_new_consumer_ch,
	p_ctx                                 context.Context) (GF_job_updates_ch, GF_job_err_ch, GF_job_complete_ch) {

	response_ch := make(chan GF_job_update_new_consumer_response)
	new_consumer := GF_job_update_new_consumer{
		Job_id_str:  p_job_id_str,
		Response_ch: response_ch,
		ctx:         p_ctx,
	}

	p_indexer_job_updates_new_consumer_ch <- new_consumer
	response := <- response_ch

	return response.Job_updates_ch, response.Job_err_ch, response.Job_complete_ch
}

//-------------------------------------------------
// HTTP
//-------------------------------------------------
func Client_http__index_block_range(p_block_start_uint uint64,
	p_block_end_uint uint64,
	p_host_port_str  string,
	p_ctx            context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (GF_indexer_job_id, *gf_core.GFerror) {

	url_str := fmt.Sprintf("http://%s/gfethm/v1/block/index?br=%d-%d",
		p_host_port_str,
		p_block_start_uint,
		p_block_end_uint)


	headers_map := map[string]string{}

	// GF_RPC_CLIENT
	data_map, gfErr := gf_rpc_lib.ClientRequest(url_str, headers_map, p_ctx, pRuntimeSys)
	if gfErr != nil {
		return GF_indexer_job_id(""), gfErr
	}

	job_id_str := GF_indexer_job_id(data_map["job_id_str"].(string))
	return job_id_str, nil
}

//-------------------------------------------------
func Client_http__index_job_updates(p_job_id_str GF_indexer_job_id,
	p_job_updates_ch chan(map[string]interface{}),
	p_host_port_str  string,
	p_ctx            context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {

	url_str := fmt.Sprintf("http://%s/gfethm/v1/block/index/job_updates?job_id=%s",
		p_host_port_str,
		p_job_id_str)
	
	headers_map := map[string]string{}

	// call SSE client that will block and process SSE events
	// going forward, but return this function right away so that 
	// the caller can consume messages from p_job_updates_ch.
	go func() {
		
		// p_job_updates_ch - channel to which to send SSE events as 
		//                    they're received over HTTP.
		gfErr := gf_rpc_lib.ClientRequestSSE(url_str,
			p_job_updates_ch,
			headers_map,
			p_ctx,
			pRuntimeSys)
		if gfErr != nil {

			// FIX!! - notify the caller of Client_http__index_job_updates() 
			//         via another error channel that SSE client failed.
			return
		}

	}()

	return nil
}