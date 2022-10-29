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

package gf_landing_page_lib

import (
	"fmt"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func InitService(p_templates_paths_map map[string]string,
	p_http_mux    *http.ServeMux,
	p_runtime_sys *gf_core.RuntimeSys) {

	//------------------------
	// STATIC FILES SERVING
	static_files__url_base_str := "/landing"
	local_dir_path_str         := "./static"
	gf_core.HTTPinitStaticServingWithMux(static_files__url_base_str,
		local_dir_path_str,
		p_http_mux,
		p_runtime_sys)

	//------------------------
	// HANDLERS
	gf_err := init_handlers(p_templates_paths_map, p_http_mux, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//------------------------
}

//-------------------------------------------------
func Run_service(p_port_str string,
	p_mongodb_host_str    string,
	p_mongodb_db_name_str string,
	p_init_done_ch        chan bool,
	pLogFun               func(string, string)) {
	pLogFun("FUN_ENTER", "gf_landing_page_service.Run_service()")

	pLogFun("INFO", "")
	pLogFun("INFO", " >>>>>>>>>>> STARTING GF_LANDING_PAGE SERVICE")
	pLogFun("INFO", "")

	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_landing_page",
		LogFun:           pLogFun,
	}

	mongo_db, _, gf_err := gf_core.Mongo__connect_new(p_mongodb_host_str, p_mongodb_db_name_str, nil, runtime_sys)
	if gf_err != nil {
		panic(-1)
	}

	runtime_sys.Mongo_db   = mongo_db 
	runtime_sys.Mongo_coll = mongo_db.Collection("data_symphony")
	
	//------------------------
	// INIT

	templates_dir_paths_map := map[string]string{
		"gf_landing_page": "./templates/gf_landing_page",
	}

	http_mux := http.NewServeMux()

	InitService(templates_dir_paths_map, http_mux, runtime_sys)

	//----------------------
	// IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}

	//----------------------

	runtime_sys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	runtime_sys.LogFun("INFO", "STARTING HTTP SERVER - PORT - "+p_port_str)
	runtime_sys.LogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	http_err := http.ListenAndServe(":"+p_port_str, nil)
	if http_err != nil {
		msg_str := "cant start listening on port - "+p_port_str
		runtime_sys.LogFun("ERROR", msg_str)
		runtime_sys.LogFun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}