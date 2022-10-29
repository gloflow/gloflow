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

package gf_tagger_lib

import (
	"os"
	"flag"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
func main() {
	logFun, _ := gf_core.InitLogs()

	cli_args_map            := CLI__parse_args(logFun)
	run__start_service_bool := cli_args_map["run__start_service_bool"].(bool)
	port_str                := cli_args_map["port_str"].(string)
	mongodb_host_str        := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str     := cli_args_map["mongodb_db_name_str"].(string)

	// START_SERVICE
	if run__start_service_bool {

		// init_done_ch := make(chan bool)

		Run_service__in_process(port_str,
			mongodb_host_str,
			mongodb_db_name_str,
			nil, // init_done_ch,
			logFun)
		// <-init_done_ch
	}
}

//-------------------------------------------------
func CLI__parse_args(pLogFun func(string, string)) map[string]interface{} {
	pLogFun("FUN_ENTER", "gf_tagger_service.CLI__parse_args()")

	//-------------------
	run__start_service_bool := flag.Bool("run__start_service", true,   "run the service daemon")
	port_str                := flag.String("port",             "3000", "port for the service to use")

	// MONGODB
	mongodb_host_str        := flag.String("mongodb_host",    "127.0.0.1", "host of mongodb to use")
	mongodb_db_name_str     := flag.String("mongodb_db_name", "prod_db"  , "DB name to use")

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
		"run__start_service_bool": *run__start_service_bool,
		"port_str":                *port_str,
		"mongodb_host_str":        *mongodb_host_str,
		"mongodb_db_name_str":     *mongodb_db_name_str,
	}
}

//-------------------------------------------------
func InitService(p_templates_paths_map map[string]string,
	p_images_jobs_mngr gf_images_jobs_core.JobsMngr,
	p_http_mux         *http.ServeMux,
	pRuntimeSys        *gf_core.RuntimeSys) {
	
	//------------------------
	// DB_INDEXES
	gfErr := DBindexInit(pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}
	
	//------------------------
	// STATIC FILES SERVING
	url_base_str       := "/tags"
	local_dir_path_str := "./static"
	gf_core.HTTPinitStaticServingWithMux(url_base_str,
		local_dir_path_str,
		p_http_mux,
		pRuntimeSys)

	//------------------------
	
	gfErr = initHandlers(p_templates_paths_map,
		p_images_jobs_mngr,
		p_http_mux,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//------------------------
}

//-------------------------------------------------
func Run_service__in_process(p_port_str string,
	p_mongodb_host_str    string,
	p_mongodb_db_name_str string,
	p_init_done_ch        chan bool,
	pLogFun               func(string, string)) {
	pLogFun("FUN_ENTER", "gf_tagger_service.Run_service__in_process()")

	pLogFun("INFO", "")
	pLogFun("INFO", " >>>>>>>>>>> STARTING GF_TAGGER SERVICE")
	pLogFun("INFO", "")
	
	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_tagger",
		LogFun:           pLogFun,
	}

	mongo_db, _, gf_err := gf_core.Mongo__connect_new(p_mongodb_host_str, p_mongodb_db_name_str, nil, runtime_sys)
	if gf_err != nil {
		panic(-1)
	}
	runtime_sys.Mongo_db   = mongo_db 
	runtime_sys.Mongo_coll = mongo_db.Collection("data_symphony")

	//----------------------
	http_mux := http.NewServeMux()

	templates_dir_paths_map := map[string]string{
		"gf_tag_objects": "./templates/gf_tag_objects/gf_tag_objects.html",
	}

	// FIX!! - jobs_mngr shouldnt be used here. when gf_tagger service is run in a separate
	//         process from gf_images service, jobs_mngr can only be reaeched via HTTP or some other
	//         transport mechanism (not via Go messages as a goroutine).
	var jobs_mngr gf_images_jobs_core.JobsMngr

	InitService(templates_dir_paths_map,
		jobs_mngr,
		http_mux,
		runtime_sys)

	//----------------------
	// IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}

	//----------------------

	pLogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	pLogFun("INFO", "STARTING HTTP SERVER - PORT - "+p_port_str)
	pLogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	err := http.ListenAndServe(":"+p_port_str, nil)
	if err != nil {
		msg_str := "cant start listening on port - "+p_port_str
		pLogFun("ERROR", msg_str)
		panic(msg_str)
	}
}