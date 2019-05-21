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
	"os"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
)

//-------------------------------------------------
type Gf_images_extern_runtime_info struct {
	Jobs_mngr             gf_images_jobs.Jobs_mngr
	Service_host_port_str string //"http://127.0.0.1:2060"
}

//-------------------------------------------------
func Run_service(p_port_str string,
	p_mongodb_host_str       string,
	p_mongodb_db_name_str    string,
	p_gf_images_runtime_info *Gf_images_extern_runtime_info,
	p_init_done_ch           chan bool,
	p_log_fun                func(string,string)) {
	p_log_fun("FUN_ENTER","gf_publisher_service.Run_service()")

	p_log_fun("INFO","")
	p_log_fun("INFO"," >>>>>>>>>>> STARTING GF_PUBLISHER SERVICE")
	p_log_fun("INFO","")
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
    p_log_fun("INFO",logo_str)
    
	mongodb_db   := gf_core.Mongo__connect(p_mongodb_host_str, p_mongodb_db_name_str, p_log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_publisher",
		Log_fun:          p_log_fun,
		Mongodb_coll:     mongodb_coll,
	}

	//------------------------
	//STATIC FILES SERVING
	static_files__url_base_str := "/posts"
	gf_core.HTTP__init_static_serving(static_files__url_base_str, runtime_sys)
	//------------------------

	//TEMPLATES_DIR
	templates_dir_path_str := "./templates"
	if _, err := os.Stat(templates_dir_path_str); os.IsNotExist(err) {
		p_log_fun("ERROR", fmt.Sprintf("templates dir doesnt exist - %s", templates_dir_path_str))
		panic(1)
	}

	err := init_handlers(p_gf_images_runtime_info, templates_dir_path_str, runtime_sys)
	if err != nil {
		msg_str := "failed to initialize http handlers - "+fmt.Sprint(err)
		panic(msg_str)
	}

	//----------------------
	//IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}
	//----------------------
	runtime_sys.Log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	runtime_sys.Log_fun("INFO","STARTING HTTP SERVER - PORT - "+p_port_str)
	runtime_sys.Log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	http_err := http.ListenAndServe(":"+p_port_str, nil)
	if http_err != nil {
		msg_str := "cant start listening on port - "+p_port_str
		runtime_sys.Log_fun("ERROR", msg_str)
		runtime_sys.Log_fun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}