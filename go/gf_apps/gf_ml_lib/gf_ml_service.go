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

func InitService(pHTTPmux *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) {

	//-------------
	// HANDLERS
	gfErr := initHandlers(pHTTPmux, pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
}

//-------------------------------------------------

func RunService(pServiceInfo *GF_service_info,
	p_init_done_ch chan bool,
	pLogFun        func(string, string)) {

	//-------------
	// RUNTIME_SYS
	
	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_ml",
		LogFun:         pLogFun,
	}

	mongo_db, _, gf_err := gf_core.MongoConnectNew(pServiceInfo.Mongodb_host_str,
		pServiceInfo.Mongodb_db_name_str,
		nil,
		runtimeSys)
	if gf_err != nil {
		os.Exit(-1)
	}

	mongo_coll := mongo_db.Collection("gf_ml")

	runtimeSys.Mongo_db   = mongo_db
	runtimeSys.Mongo_coll = mongo_coll
	//-------------
	// INIT

	HTTPmux := http.NewServeMux()

	InitService(HTTPmux, runtimeSys)

	//-------------

	runtimeSys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	runtimeSys.LogFun("INFO", "STARTING HTTP SERVER - PORT - "+pServiceInfo.Port_str)
	runtimeSys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	http_err := http.ListenAndServe(":"+pServiceInfo.Port_str, nil)
	if http_err != nil {
		msg_str := fmt.Sprintf("cant start listening on port - ", pServiceInfo.Port_str)
		runtimeSys.LogFun("ERROR", msg_str)
		runtimeSys.LogFun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}