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
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_tagger_core"
)

//---------------------------------------------------

func main() {

	logFun, _ := gf_core.LogsInit()

	cliArgsMap              := gf_tagger_core.CLIparseArgs(logFun)
	run__start_service_bool := cliArgsMap["run__start_service_bool"].(bool)
	port_str                := cliArgsMap["port_str"].(string)
	mongodb_host_str        := cliArgsMap["mongodb_host_str"].(string)
	mongodb_db_name_str     := cliArgsMap["mongodb_db_name_str"].(string)

	// START_SERVICE
	if run__start_service_bool {

		// init_done_ch := make(chan bool)

		gf_tagger_core.RunService(port_str,
			mongodb_host_str,
			mongodb_db_name_str,
			nil, // init_done_ch,
			logFun)
		// <-init_done_ch
	}
}