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
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func init_handlers(p_queue_info *GF_queue_info,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_eth_monitor_handlers.init_handlers()")

	//---------------------
	http.HandleFunc("/ethm/v1/register", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST - /ethm/v1/register ----------")


	})

	//---------------------

	return nil
}