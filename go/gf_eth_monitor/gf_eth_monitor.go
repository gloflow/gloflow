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

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_lib"
)

//-------------------------------------------------
type GF_eth_monitor_runtime struct {
	config      *GF_config
	runtime_sys *gf_core.Runtime_sys
}

//-------------------------------------------------
func main() {

	log_fun := gf_core.Init_log_fun()

	cmd__base := cmds_init(log_fun)
	err := cmd__base.Execute()
	if err != nil {
		panic(err)
	}
}

//-------------------------------------------------
func runtime__get(p_log_fun func(string, string)) (*GF_eth_monitor_runtime, error) {

	// CONFIG
	config_dir_path_str := "./../config/"
	config_name_str     := "gf_eth_monitor"
	
	config, err := config__init(config_dir_path_str, config_name_str)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to load config")
		return nil, err
	}

	// RUNTIME_SYS
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_eth_monitor",
		Log_fun:          p_log_fun,	
	}

	//--------------------
	// MONGODB
	mongodb_host_str := config.Mongodb_host_str
	fmt.Printf("mongodb_host - %s\n", mongodb_host_str)

	mongo_db, gf_err := gf_core.Mongo__connect_new(config.Mongodb_host_str,
		config.Mongodb_db_name_str,
		runtime_sys)
	if gf_err != nil {
		return nil, gf_err.Error
	}

	// mongo_coll := mongo_db.Collection(config.Mongodb_coll_name_str)
	runtime_sys.Mongo_db = mongo_db

	//--------------------
	// RUNTIME
	runtime := &GF_eth_monitor_runtime{
		config:      config,
		runtime_sys: runtime_sys,
	}

	return runtime, nil
}

//-------------------------------------------------
func cmds_init(p_log_fun func(string, string)) *cobra.Command {

	// BASE
	cmd__base := &cobra.Command{
		Use:   "gf_eth_monitor",
		Short: "gf_eth_monitor server",
		Long:  "",
		Run:   func(p_cmd *cobra.Command, p_args []string) {

		},
	}

	//--------------------
	// CLI_ARGUMENT - PORT
	cmd__base.PersistentFlags().StringP("port", "p", "PORT NUMBER", "port on which to listen for HTTP traffic") // Cobra CLI argument
	err := viper.BindPFlag("port", cmd__base.PersistentFlags().Lookup("port"))                                  // Bind Cobra CLI argument to a Viper configuration (for default value)
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - MONGODB_HOST
	cmd__base.PersistentFlags().StringP("mongodb_host", "m", "MONGODB HOST", "mongodb host to which to connect") // Cobra CLI argument
	err = viper.BindPFlag("mongodb_host", cmd__base.PersistentFlags().Lookup("mongodb_host"))                   // Bind Cobra CLI argument to a Viper configuration (for default value)
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("mongodb_host", "GF_MONGODB_HOST")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------

	// START
	cmd__start := &cobra.Command{
		Use:   "start",
		Short: "start some service or function",
		Run:   func(p_cmd *cobra.Command, p_args []string) {

		},
	}

	// START_SERVICE
	cmd__start_service := &cobra.Command{
		Use:   "service",
		Short: "start the gf_eth_monitor service",
		Long:  "start the gf_eth_monitor service in a target cluster",
		Run:   func(p_cmd *cobra.Command, p_args []string) {

			runtime, err := runtime__get(p_log_fun)
			if err != nil {
				return
			}
			
			service_info := gf_eth_monitor_lib.GF_service_info{
				Port_str: runtime.config.Port_str,
			}

			// RUN_SERVICE
			gf_eth_monitor_lib.Run_service(&service_info, runtime.runtime_sys)			
		},
	}

	cmd__start.AddCommand(cmd__start_service)
	cmd__base.AddCommand(cmd__start)

	return cmd__base
}