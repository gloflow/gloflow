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

package gf_ml_lib

import (
	"os"
	"fmt"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_service_info struct {
	Port_str                string
	Mongodb_host_str        string
	Mongodb_db_name_str     string
	Templates_dir_paths_map map[string]interface{}
	Config_file_path_str    string
}

//-------------------------------------------------
func Init_service(p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_ml_service.Init_service()")

	//-------------
	// HANDLERS
	gf_err := init_handlers(p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
}

//-------------------------------------------------
func Run_service(p_service_info *GF_service_info,
	p_init_done_ch chan bool,
	p_log_fun      func(string, string)) {
	p_log_fun("FUN_ENTER", "gf_ml_service.Run_service()")

	//-------------
	// RUNTIME_SYS
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_ml",
		Log_fun:          p_log_fun,
	}

	mongo_db, _, gf_err := gf_core.Mongo__connect_new(p_service_info.Mongodb_host_str,
		p_service_info.Mongodb_db_name_str,
		nil,
		runtime_sys)
	if gf_err != nil {
		os.Exit(-1)
	}

	mongo_coll := mongo_db.Collection("g_ml")

	runtime_sys.Mongo_db   = mongo_db
	runtime_sys.Mongo_coll = mongo_coll
	//-------------
	// INIT
	Init_service(runtime_sys)

	//-------------

	runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	runtime_sys.Log_fun("INFO", "STARTING HTTP SERVER - PORT - "+p_service_info.Port_str)
	runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	http_err := http.ListenAndServe(":"+p_service_info.Port_str, nil)
	if http_err != nil {
		msg_str := fmt.Sprintf("cant start listening on port - ", p_service_info.Port_str)
		runtime_sys.Log_fun("ERROR", msg_str)
		runtime_sys.Log_fun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}