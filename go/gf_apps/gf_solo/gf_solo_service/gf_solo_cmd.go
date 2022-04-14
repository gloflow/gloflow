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
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func Cmds_init(p_external_plugins *gf_core.External_plugins,
	p_log_fun func(string, string)) *cobra.Command {

	// BASE
	cmdBase := &cobra.Command{
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
	cmdBase.PersistentFlags().StringVarP(&cli_config_path_str, "config", "", cli_config_path__default_str,
		"config file path on the local FS")

	//--------------------
	// CLI_ARGUMENT - PORT

	// Cobra CLI argument
	cmdBase.PersistentFlags().StringP("port", "p", "port for the main service",
		"port on which to listen for HTTP traffic")

	// Bind Cobra CLI argument to a Viper configuration (for default value)
	err := viper.BindPFlag("port", cmdBase.PersistentFlags().Lookup("port"))
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
	// CLI_ARGUMENT - PORT_ADMIN

	// Cobra CLI argument
	cmdBase.PersistentFlags().StringP("port_admin", "a", "port for the admin service",
		"port on which to listen for HTTP admin traffic")
	
	// Bind Cobra CLI argument to a Viper configuration (for default value)
	err = viper.BindPFlag("port_admin", cmdBase.PersistentFlags().Lookup("port_admin"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}
	
	// ENV
	err = viper.BindEnv("port", "GF_PORT_ADMIN")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - PORT__METRICS
	cmdBase.PersistentFlags().StringP("port_metrics", "", "METRICS PORT NUMBER",
		"port on which to listen for metrics HTTP traffic")
	err = viper.BindPFlag("port_metrics", cmdBase.PersistentFlags().Lookup("port_metrics"))
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
	cmdBase.PersistentFlags().StringP("mongodb_host", "m", "MONGODB HOST", "mongodb host to which to connect") // Cobra CLI argument
	err = viper.BindPFlag("mongodb_host", cmdBase.PersistentFlags().Lookup("mongodb_host"))                    // Bind Cobra CLI argument to a Viper configuration (for default value)
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
	cmdBase.PersistentFlags().StringP("mongodb_db_name", "", "MONGODB HOST", "mongodb db name to which to connect") // Cobra CLI argument
	err = viper.BindPFlag("mongodb_db_name", cmdBase.PersistentFlags().Lookup("mongodb_db_name"))                 // Bind Cobra CLI argument to a Viper configuration (for default value)
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
	cmdBase.PersistentFlags().StringP("sentry_endpoint", "", "SENTRY ENDPOINT", "Sentry endpoint to send errors to")
	err = viper.BindPFlag("sentry_endpoint", cmdBase.PersistentFlags().Lookup("sentry_endpoint"))
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
	cmdBase.PersistentFlags().StringP("aws_access_key_id", "", "AWS ACCESS_KEY_ID", "AWS access_key_id")
	err = viper.BindPFlag("aws_access_key_id", cmdBase.PersistentFlags().Lookup("aws_access_key_id"))
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
	cmdBase.PersistentFlags().StringP("aws_secret_access_key", "", "AWS SECRET_ACCESS_KEY", "AWS secret_access_key")
	err = viper.BindPFlag("aws_secret_access_key", cmdBase.PersistentFlags().Lookup("aws_secret_access_key"))
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
	// CLI_ARGUMENT - ADMIN_MFA_SECRET_KEY_BASE32
	admin_mfa_secret_key_base32_str := "aabbccddeeffgghh"
	cmdBase.PersistentFlags().StringP("admin_mfa_secret_key_base32", "", admin_mfa_secret_key_base32_str,
		"secret key used to verify mfa (totp/hotp), base32 encoded")
	err = viper.BindPFlag("admin_mfa_secret_key_base32", cmdBase.PersistentFlags().Lookup("admin_mfa_secret_key_base32"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("admin_mfa_secret_key_base32", "GF_ADMIN_MFA_SECRET_KEY_BASE32")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// CLI_ARGUMENT - ADMIN_EMAIL
	admin_email__default_str := ""
	cmdBase.PersistentFlags().StringP("admin_email", "", admin_email__default_str, "default email of the administrator to use")
	err = viper.BindPFlag("admin_email", cmdBase.PersistentFlags().Lookup("admin_email"))
	if err != nil {
		fmt.Println("failed to bind CLI arg to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("admin_email", "GF_ADMIN_EMAIL")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------

	// ENV
	err = viper.BindEnv("domain_base", "GF_DOMAIN_BASE")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("domain_admin_base", "GF_DOMAIN_ADMIN_BASE")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// START
	cmdStart := &cobra.Command{
		Use:   "start",
		Short: "start some service or function",
		Run:   func(pCmd *cobra.Command, p_args []string) {

		},
	}

	// START_SERVICE
	cmdStartService := &cobra.Command{
		Use:   "service",
		Short: "start the gf_solo service",
		Long:  "start the gf_solo service",
		Run:   func(pCmd *cobra.Command, p_args []string) {

			




			runtimeSys, config, err := Runtime__get(cli_config_path_str, p_external_plugins, p_log_fun)
			if err != nil {
				panic(err)
			}

			// RUN_SERVICE
			Run(config, runtimeSys)			
		},
	}

	//--------------------
	cmdStart.AddCommand(cmdStartService)
	cmdBase.AddCommand(cmdStart)

	return cmdBase
}