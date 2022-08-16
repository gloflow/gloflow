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

package gf_publisher_lib

import (
	"os"
	"flag"
)

//-------------------------------------------------
func CLI__parse_args(pLogFun func(string,string)) map[string]interface{} {
	pLogFun("FUN_ENTER", "gf_publisher_cli.CLI__parse_args()")

	//-------------------
	run__start_service_bool         := flag.Bool("run__start_service",                true,                       "run the service daemon")
	port_str                        := flag.String("port",                            "2020",                     "port for the service to use")
	gf_images_service_host_port_str := flag.String("gf_images_service_host_port_str", "gf_images_service_1:3050", "gf_images service host")
	
	// MONGODB
	mongodb_host_str                := flag.String("mongodb_host",    "127.0.0.1", "host of mongodb to use")
	mongodb_db_name_str             := flag.String("mongodb_db_name", "prod_db",   "DB name to use")

	// MONGODB_ENV
	mongodb_host_env_str    := os.Getenv("GF_MONGODB_HOST")
	mongodb_db_name_env_str := os.Getenv("GF_MONGODB_DB_NAME")

	if mongodb_db_name_env_str != "" {
		*mongodb_db_name_str = mongodb_db_name_env_str
	}

	if mongodb_host_env_str != "" {
		*mongodb_host_str = mongodb_host_env_str
	}
	//-------------------

	flag.Parse()

	return map[string]interface{}{
		"run__start_service_bool":         *run__start_service_bool,
		"port_str":                        *port_str,
		"mongodb_host_str":                *mongodb_host_str,
		"mongodb_db_name_str":             *mongodb_db_name_str,
		"gf_images_service_host_port_str": *gf_images_service_host_port_str,
	}
}