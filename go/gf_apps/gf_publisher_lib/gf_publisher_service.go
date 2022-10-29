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
	"fmt"
	// "os"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
type GF_images_extern_runtime_info struct {
	Jobs_mngr               gf_images_jobs_core.JobsMngr
	Service_host_port_str   string // "http://127.0.0.1:2060"
	Templates_dir_paths_map map[string]string
}

//-------------------------------------------------
func InitService(pHTTPmux *http.ServeMux,
	p_gf_images_runtime_info *GF_images_extern_runtime_info,
	pRuntimeSys              *gf_core.RuntimeSys) {

	//------------------------
	// STATIC FILES SERVING
	static_files__url_base_str := "/posts"
	local_dir_path_str         := "./static"
	gf_core.HTTPinitStaticServingWithMux(static_files__url_base_str,
		local_dir_path_str,
		pHTTPmux,
		pRuntimeSys)
	
	//------------------------

	err := init_handlers(p_gf_images_runtime_info,
		p_gf_images_runtime_info.Templates_dir_paths_map,
		pHTTPmux,
		pRuntimeSys)
	if err != nil {
		msg_str := "failed to initialize http handlers - "+fmt.Sprint(err)
		panic(msg_str)
	}
}

//-------------------------------------------------
func Run_service(p_port_str string,
	p_mongodb_host_str       string,
	p_mongodb_db_name_str    string,
	p_gf_images_runtime_info *GF_images_extern_runtime_info,
	p_init_done_ch           chan bool,
	pLogFun                func(string,string)) {
	pLogFun("FUN_ENTER","gf_publisher_service.Run_service()")

	pLogFun("INFO","")
	pLogFun("INFO"," >>>>>>>>>>> STARTING GF_PUBLISHER SERVICE")
	pLogFun("INFO","")
	logo_str := `
	                   #\   /##/      #
                    #   #/       #/
     ####\    /##\  #\__\#\     #/         /#
       \##\  /#  #\  ######|    #     /####/
         |#\_|___##| |#####|__ #/ _/######/
         \#########|_|##################/
           \###########/     \########/
            \#########|        \###|
        \##\ \########/   @@   |###| ___/#####
           #\ \######|    @@   |#########
            #\ //   \|         ||
            \##|     \\____ ####| /########
      _____  \#|_@@__|#####/....\##/ \#/  \#
     #######\ /######MMM#/ ......|#        \#
          /###/......\M/ ...... .\#######
         |#| .........|...........|###\
      ___|#|..........|......  .../| \##
     ########.. ......\........./##|   \#
         |#|.........../\.._____|##|     \##
         |##\  ...__.--|---#########
        /####\___/##/--|--|#######/ #
       /#    \######|-----|#/ \#    \##
     ##/     /|######\---/#/          \#
            ##/ |#########/            \#
               /########|               \##
              /#########|
             ## |#######|
                |#########\
                 |########|
                 |#########\
                 |##########\
                 |############\`
    pLogFun("INFO",logo_str)
	
	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_publisher",
		LogFun:           pLogFun,
	}

	mongo_db, _, gf_err   := gf_core.Mongo__connect_new(p_mongodb_host_str, p_mongodb_db_name_str, nil, runtime_sys)
	if gf_err != nil {
		panic(-1)
	}
	runtime_sys.Mongo_db   = mongo_db
	runtime_sys.Mongo_coll = mongo_db.Collection("data_symphony")

	//------------------------
	// INIT

	http_mux := http.NewServeMux()

	InitService(http_mux, p_gf_images_runtime_info, runtime_sys)

	//----------------------
	// IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}

	//----------------------
	runtime_sys.LogFun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	runtime_sys.LogFun("INFO","STARTING HTTP SERVER - PORT - "+p_port_str)
	runtime_sys.LogFun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	http_err := http.ListenAndServe(":"+p_port_str, nil)
	if http_err != nil {
		msg_str := "cant start listening on port - "+p_port_str
		runtime_sys.LogFun("ERROR", msg_str)
		runtime_sys.LogFun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}