package main

import (
	"fmt"
	"flag"
	"net/http"
	"gf_core"
)
//-------------------------------------------------
func main() {
	log_fun := gf_core.Init_log_fun()

	cli_args_map            := parse__cli_args(log_fun)
	//in_container_bool       := cli_args_map["in_container_bool"].(bool)
	run__start_service_bool := cli_args_map["run__start_service_bool"].(bool)
	port_str                := cli_args_map["port_str"].(string)
	mongodb_host_str        := cli_args_map["mongodb_host_str"].(string)
	mongodb_db_name_str     := cli_args_map["mongodb_db_name_str"].(string)
	gf_images_service_host_port_str := cli_args_map["gf_images_service_host_port_str"].(string)

	gf_core.Init_exit_os_signals(log_fun)

	//START_SERVICE
	if run__start_service_bool {
		//init_done_ch := make(chan bool)
		Run_service__in_process(port_str,
			mongodb_host_str,
			mongodb_db_name_str,
			gf_images_service_host_port_str,
			nil, //init_done_ch,
			log_fun)
		//<-init_done_ch
	}
}
//-------------------------------------------------
func parse__cli_args(p_log_fun func(string,string)) map[string]interface{} {
	p_log_fun("FUN_ENTER","gf_publisher_service.parse__cli_args()")

	//-------------------
	//in_container_bool       := flag.Bool("in_container"      ,false      ,"is th service being run in a Docker container")
	run__start_service_bool         := flag.Bool("run__start_service",               true,                      "run the service daemon")
	port_str                        := flag.String("port",                           "2020",                    "port for the service to use")
	mongodb_host_str                := flag.String("mongodb_host",                   "127.0.0.1",               "host of mongodb to use")
	mongodb_db_name_str             := flag.String("mongodb_db_name",                "prod_db",                 "DB name to use")
	gf_images_service_host_port_str := flag.String("gf_images_service_host_port_str","gf_images_service_1:3050","gf_images service host")
	//-------------------
	flag.Parse()

	return map[string]interface{}{
		"run__start_service_bool":        *run__start_service_bool,
		"port_str":                       *port_str,
		"mongodb_host_str":               *mongodb_host_str,
		"mongodb_db_name_str":            *mongodb_db_name_str,
		"gf_images_service_host_port_str":*gf_images_service_host_port_str,
	}
}
//-------------------------------------------------
func Run_service__in_process(p_port_str string,
			p_mongodb_host_str                string,
			p_mongodb_db_name_str             string,
			p_gf_images_service_host_port_str string,
			p_init_done_ch                    chan bool,
			p_log_fun                         func(string,string)) {
	p_log_fun("FUN_ENTER","gf_publisher_service.Run_service__in_process()")

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
    
	mongo_db   := gf_core.Conn_to_mongodb(p_mongodb_host_str, p_mongodb_db_name_str, p_log_fun)
	mongo_coll := mongo_db.C("data_symphony")

	//------------------------
	//STATIC FILES SERVING
	static_files__url_base_str := "/posts"
	gf_core.HTTP__init_static_serving(&static_files__url_base_str,p_log_fun)
	//------------------------

	err := init_handlers(&p_gf_images_service_host_port_str, mongo_coll, p_log_fun)
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

	p_log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	p_log_fun("INFO","STARTING HTTP SERVER - PORT - "+p_port_str)
	p_log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	err = http.ListenAndServe(":"+p_port_str,nil)
	if err != nil {
		msg_str := "cant start listening on port - "+p_port_str
		p_log_fun("ERROR",msg_str)
		panic(msg_str)
	}
}