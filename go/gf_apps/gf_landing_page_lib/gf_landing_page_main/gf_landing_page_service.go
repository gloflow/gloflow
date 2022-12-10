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

package main

import (
	"flag"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_landing_page_lib"
)

//-------------------------------------------------

func main() {
	logFun, _ := gf_core.LogsInit()

	cli_args_map            := parse__cli_args(logFun)
	run__start_service_bool := cli_args_map["run__start_service_bool"].(bool)
	port_str                := cli_args_map["port_str"].(string)
	mongodb_host_str        := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str     := cli_args_map["mongodb_db_name_str"].(string)

	//START_SERVICE
	if run__start_service_bool {
		gf_landing_page_lib.Run_service(port_str,
			mongodb_host_str,
			mongodb_db_name_str,
			nil,
			logFun)
	}
}

//-------------------------------------------------

func parse__cli_args(pLogFun func(string, string)) map[string]interface{} {
	pLogFun("FUN_ENTER", "gf_landing_page_service.parse__cli_args()")

	//-------------------
	run__start_service_bool := flag.Bool("run__start_service", true,        "run the service daemon")
	port_str                := flag.String("port",             "2000",      "port for the service to use")
	mongodb_host_str        := flag.String("mongodb_host",     "127.0.0.1", "host of mongodb to use")
	mongodb_db_name_str     := flag.String("mongodb_db_name",  "prod_db",   "DB name to use")
	
	//-------------------
	flag.Parse()

	return map[string]interface{}{
		"run__start_service_bool": *run__start_service_bool,
		"port_str":                *port_str,
		"mongodb_host_str":        *mongodb_host_str,
		"mongodb_db_name_str":     *mongodb_db_name_str,
	}
}