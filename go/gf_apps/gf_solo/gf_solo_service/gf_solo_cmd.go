/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_solo_service

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//-------------------------------------------------
func Cmds_init(p_log_fun func(string, string)) *cobra.Command {

	// BASE
	cmd__base := &cobra.Command{
		Use:   "gf_solo",
		Short: "gf_solo",
		Long:  "",
		Run:   func(p_cmd *cobra.Command, p_args []string) {

		},
	}

	//--------------------
	// CLI_ARGUMENT - CONFIG
	cli_config_path__default_str := "./config/gf_solo"
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
	// CLI_ARGUMENT - MONGODB_DB_NAME
	cmd__base.PersistentFlags().StringP("mongodb_db_name", "", "MONGODB HOST", "mongodb db name to which to connect") // Cobra CLI argument
	err = viper.BindPFlag("mongodb_db_name", cmd__base.PersistentFlags().Lookup("mongodb_db_name"))                 // Bind Cobra CLI argument to a Viper configuration (for default value)
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("mongodb_db_name", "GF_MONGODB_DB_NAME")
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
	// CLI_ARGUMENT - AWS_ACCESS_KEY_ID
	cmd__base.PersistentFlags().StringP("aws_access_key_id", "", "AWS ACCESS_KEY_ID", "AWS access_key_id")
	err = viper.BindPFlag("aws_access_key_id", cmd__base.PersistentFlags().Lookup("aws_access_key_id"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("aws_access_key_id", "AWS_ACCESS_KEY_ID")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - AWS_SECRET_ACCESS_KEY
	cmd__base.PersistentFlags().StringP("aws_secret_access_key", "", "AWS SECRET_ACCESS_KEY", "AWS secret_access_key")
	err = viper.BindPFlag("aws_secret_access_key", cmd__base.PersistentFlags().Lookup("aws_secret_access_key"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("aws_secret_access_key", "AWS_SECRET_ACCESS_KEY")
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
		Short: "start the gf_solo service",
		Long:  "start the gf_solo service",
		Run:   func(p_cmd *cobra.Command, p_args []string) {

			




			runtime_sys, config, err := Runtime__get(cli_config_path_str, p_log_fun)
			if err != nil {
				panic(err)
			}

			// RUN_SERVICE
			Run(config, runtime_sys)			
		},
	}

	//--------------------
	cmd__start.AddCommand(cmd__start_service)
	cmd__base.AddCommand(cmd__start)

	return cmd__base
}