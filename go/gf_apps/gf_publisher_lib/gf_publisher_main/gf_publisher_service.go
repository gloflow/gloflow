/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package main

import (
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib"
)

//-------------------------------------------------
func main() {
	logFun, _ := gf_core.InitLogs()

	cli_args_map            := gf_publisher_lib.CLI__parse_args(logFun)
	run__start_service_bool := cli_args_map["run__start_service_bool"].(bool)
	port_str                := cli_args_map["port_str"].(string)
	mongodb_host_str        := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str     := cli_args_map["mongodb_db_name_str"].(string)
	gf_images_service_host_port_str := cli_args_map["gf_images_service_host_port_str"].(string)

	// START_SERVICE
	if run__start_service_bool {

		gf_images_runtime := &gf_publisher_lib.GF_images_extern_runtime_info{
			Jobs_mngr:             nil, //indicates not to send in-process messages to jobs_mngr goroutine, instead use HTTP REST API of gf_images
			Service_host_port_str: gf_images_service_host_port_str,
		}
		
		//init_done_ch := make(chan bool)
		gf_publisher_lib.Run_service(port_str,
			mongodb_host_str,
			mongodb_db_name_str,
			gf_images_runtime,
			nil, //init_done_ch,
			logFun)
		//<-init_done_ch
	}
}