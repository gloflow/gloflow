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
	"os"
	"fmt"
	"path"
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_web3_lib"
)

//-------------------------------------------------
func main() {

	logFun, _ := gf_core.InitLogs()
	log.SetOutput(os.Stdout)

	cmd__base := cmds_init(logFun)
	err := cmd__base.Execute()
	if err != nil {
		panic(err)
	}
}

//-------------------------------------------------
func runtimeGet(p_config_path_str string,
	pLogFun func(string, string)) (*gf_eth_core.GF_runtime, error) {

	// CONFIG
	config_dir_path_str := path.Dir(p_config_path_str)  // "./../config/"
	config_name_str     := path.Base(p_config_path_str) // "gf_eth_monitor"
	
	config, err := config__init(config_dir_path_str, config_name_str)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to load config")
		return nil, err
	}

	// RUNTIME_SYS
	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_eth_monitor",
		LogFun:           pLogFun,

		// SENTRY - enable it for error reporting
		Errors_send_to_sentry_bool: true,	
	}

	runtime, err := gf_eth_core.RuntimeGet(config, runtime_sys)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

//-------------------------------------------------
func cmds_init(pLogFun func(string, string)) *cobra.Command {

	// BASE
	cmd__base := &cobra.Command{
		Use:   "gf_eth_monitor",
		Short: "gf_eth_monitor server",
		Long:  "",
		Run:   func(p_cmd *cobra.Command, p_args []string) {

		},
	}

	//--------------------
	// CLI_ARGUMENT - CONFIG
	cli_config_path__default_str := "./../config/gf_eth_monitor"
	var cli_config_path_str string
	cmd__base.PersistentFlags().StringVarP(&cli_config_path_str, "config", "", cli_config_path__default_str,
		"config file path on the local FS")

	//--------------------
	// CLI_ARGUMENT - PORT
	cmd__base.PersistentFlags().StringP("port", "p", "PORT NUMBER",
		"port on which to listen for HTTP traffic") // Cobra CLI argument
	err := viper.BindPFlag("port", cmd__base.PersistentFlags().Lookup("port")) // Bind Cobra CLI argument to a Viper configuration (for default value)
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}
	
	// ENV
	err = viper.BindEnv("port", "GF_PORT")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - PORT__METRICS
	cmd__base.PersistentFlags().StringP("port_metrics", "", "METRICS PORT NUMBER",
		"port on which to listen for metrics HTTP traffic")
	err = viper.BindPFlag("port_metrics", cmd__base.PersistentFlags().Lookup("port_metrics"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}
	
	// ENV
	err = viper.BindEnv("port_metrics", "GF_PORT_METRICS")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - MONGODB_HOST
	cmd__base.PersistentFlags().StringP("mongodb_host", "m", "MONGODB HOST", "mongodb host to which to connect") // Cobra CLI argument
	err = viper.BindPFlag("mongodb_host", cmd__base.PersistentFlags().Lookup("mongodb_host"))                    // Bind Cobra CLI argument to a Viper configuration (for default value)
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
	// CLI_ARGUMENT - INFLUXDB_HOST
	cmd__base.PersistentFlags().StringP("influxdb_host", "i", "INFLUXDB HOST", "influxdb host to which to connect") // Cobra CLI argument
	err = viper.BindPFlag("influxdb_host", cmd__base.PersistentFlags().Lookup("influxdb_host"))                     // Bind Cobra CLI argument to a Viper configuration (for default value)
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("influxdb_host", "GF_INFLUXDB_HOST")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - AWS_SQS_QUEUE
	cmd__base.PersistentFlags().StringP("aws_sqs_queue", "q", "AWS SQS QUEUE", "AWS SQS queue from which to consume events")
	err = viper.BindPFlag("aws_sqs_queue", cmd__base.PersistentFlags().Lookup("aws_sqs_queue"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("aws_sqs_queue", "GF_AWS_SQS_QUEUE")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - WORKERS_AWS_DISCOVERY
	cmd__base.PersistentFlags().StringP("workers_aws_discovery", "", "WORKERS HOSTS", "if AWS EC2 discovery should be enbaled to dynamicaly discover workers")
	err = viper.BindPFlag("workers_aws_discovery", cmd__base.PersistentFlags().Lookup("workers_aws_discovery"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("workers_aws_discovery", "GF_WORKERS_AWS_DISCOVERY")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - WORKERS_HOSTS
	cmd__base.PersistentFlags().StringP("workers_hosts", "", "WORKERS HOSTS", "list of all workers hosts, ',' separated")
	err = viper.BindPFlag("workers_hosts", cmd__base.PersistentFlags().Lookup("workers_hosts"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("workers_hosts", "GF_WORKERS_HOSTS")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - SENTRY_ENDPOINT
	cmd__base.PersistentFlags().StringP("sentry_endpoint", "", "SENTRY ENDPOINT", "Sentry endpoint to send errors to")
	err = viper.BindPFlag("sentry_endpoint", cmd__base.PersistentFlags().Lookup("sentry_endpoint"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("sentry_endpoint", "GF_SENTRY_ENDPOINT")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - EVENTS_CONSUME
	cmd__base.PersistentFlags().StringP("events_consume", "", "EVENTS CONSUME", "on/off consumption and processing of events from queue")
	err = viper.BindPFlag("events_consume", cmd__base.PersistentFlags().Lookup("events_consume"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("events_consume", "GF_EVENTS_CONSUME")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - PY_PLUGINS
	cmd__base.PersistentFlags().StringP("py_plugins_dir_path", "", "PY PLUGINS DIR PATH", "path to the directory holding Py plugin files")
	err = viper.BindPFlag("py_plugins_dir_path", cmd__base.PersistentFlags().Lookup("py_plugins_dir_path"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("py_plugins_dir_path", "GF_PY_PLUGINS_BASE_DIR_PATH")
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

			runtime, err := runtimeGet(cli_config_path_str, pLogFun)
			if err != nil {
				panic(err)
			}
			
			/*service_info := gf_web3_lib.GF_service_info{
				Port_str:           runtime.Config.Port_str,
				SQS_queue_name_str: runtime.Config.AWS_SQS_queue_str,
				Influxdb_client:    runtime.Influxdb_client,
			}*/

			// RUN_SERVICE
			gf_web3_lib.Run_service(runtime)			
		},
	}

	//--------------------
	// TEST
	cmd__test := &cobra.Command{
		Use:   "test",
		Short: "test some functionility",
		Run:   func(p_cmd *cobra.Command, p_args []string) {

		},
	}

	// TEST_WORKER_EVENT_PROCESS
	cmd__test_worker_event_process := &cobra.Command{
		Use:   "worker_event_process",
		Short: "test processing of worker events",
		Run:   func(p_cmd *cobra.Command, p_args []string) {




			runtime, err := runtimeGet(cli_config_path_str, pLogFun)
			if err != nil {
				panic(err)
			}


			SQS_queue_name_str := runtime.Config.AWS_SQS_queue_str
			queue_info, err    := gf_web3_lib.Event__init_queue(SQS_queue_name_str, nil)
			if err != nil {
				panic(err)
			}

			ctx := context.Background()

			// PROCESS_SINGLE_EVENT
			gf_web3_lib.Event__process_from_sqs(queue_info, ctx, nil, runtime)
		},
	}

	//--------------------
	cmd__start.AddCommand(cmd__start_service)
	cmd__test.AddCommand(cmd__test_worker_event_process)
	cmd__base.AddCommand(cmd__start)
	cmd__base.AddCommand(cmd__test)

	return cmd__base
}