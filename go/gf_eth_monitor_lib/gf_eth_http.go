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
	"strconv"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func Http__get_arg__block_num(p_resp http.ResponseWriter,
	p_req         *http.Request,
	p_runtime_sys *gf_core.Runtime_sys) (uint64, error) {



	//------------------
	// INPUT
	qs_map := p_req.URL.Query()

	var block_num_int uint64
	if b_lst, ok := qs_map["b"]; ok {
		block_num_str := b_lst[0]

		i, err := strconv.Atoi(block_num_str)
		if err != nil {
			gf_rpc_lib.Error__in_handler("/gfethm_worker_inspect/v1/blocks",
				fmt.Sprintf("invalid input argument - %s", i),
				nil, p_resp, p_runtime_sys)
			return 0, err 
		}
		block_num_int = uint64(i)
	}
	
	//------------------

	return block_num_int, nil
}