/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	"strings"
	"github.com/spf13/viper"
	gf_core "github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------
// CONFIG_INIT

/*
load config for gf_solo from all the sources.
Viper lib is used for config loading.
configs are loaded from:
- ENV vars
- gf_solo.yaml config file
- via CLI args
*/
func ConfigInit(pConfigDirPathStr string,
	pConfigFileNameStr   string,
	pConfigLoadPluginFun gf_core.GFpluginConfigLoadCallback,
	pRuntimeSys          *gf_core.RuntimeSys) (*gf_core.GFconfig, error) {

	configNameStr := strings.Split(pConfigFileNameStr, ".")[0] // viper expects just the file name, without extension

	// FILE
	viper.AddConfigPath(pConfigDirPathStr)
	viper.SetConfigName(configNameStr)

	//--------------------
	// ENV_VARS
	// all config members that have their mapstructure name for Viper config,
	// also have a corresponding ENV var name thats generated for them by
	// upper-casing their name.
	//--------------------
	// ENV_PREFIX - "GF" becomes "GF_" - prefix expected in all recognized ENV vars.

	viper.SetEnvPrefix("GF")

	// IMPORTANT!! - enable Viper parsing ENV vars.
	viper.AutomaticEnv()

	//--------------------

	//--------------------
	// ENV_VARS
	bindEnvVars()

	//--------------------

	// LOAD
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// CONFIG
	config := &gf_core.GFconfig{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	// IMPORTANT!! - Environment is ciritical param for all subsequent potential config loading,
	//               in pConfigLoadPluginFun, so setting it here right away.
	//               when config is loaded, either via ENV vars or from config file, the Env value is found out.
	pRuntimeSys.EnvStr = config.EnvStr

	// PLUGIN - allows users to specify a custom function to load additional configs
	//          external to GF core
	if pConfigLoadPluginFun != nil {
		newConfig, gfErr := pConfigLoadPluginFun(config, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr.Error
		}
		return newConfig, nil
	}

	return config, nil
}

//-------------------------------------------------------------

func bindEnvVars() {

	// ENV
	err := viper.BindEnv("env", "GF_ENV")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// PORT
	err = viper.BindEnv("port", "GF_PORT")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// PORT_ADMIN
	err = viper.BindEnv("port", "GF_PORT_ADMIN")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// PORT_METRICS
	err = viper.BindEnv("port_metrics", "GF_PORT_METRICS")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// MONGODB_HOST
	err = viper.BindEnv("mongodb_host", "GF_MONGODB_HOST")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// MONGODB_DB_NAME
	err = viper.BindEnv("mongodb_db_name", "GF_MONGODB_DB_NAME")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// SQL
	//--------------------
	// GF_SQL_USER_NAME
	err = viper.BindEnv("sql_user_name", "GF_SQL_USER_NAME")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// GF_SQL_PASS
	err = viper.BindEnv("sql_pass", "GF_SQL_PASS")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// GF_SQL_HOST
	err = viper.BindEnv("sql_host", "GF_SQL_HOST")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// GF_SQL_DB_NAME
	err = viper.BindEnv("sql_db_name", "GF_SQL_DB_NAME")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// GF_SENTRY_ENDPOINT
	err = viper.BindEnv("sentry_endpoint", "GF_SENTRY_ENDPOINT")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// GF_AUTH_SUBSYSTEM_TYPE
	err = viper.BindEnv("auth_subsystem_type", "GF_AUTH_SUBSYSTEM_TYPE")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// GF_ADMIN_MFA_SECRET_KEY_BASE32
	err = viper.BindEnv("admin_mfa_secret_key_base32", "GF_ADMIN_MFA_SECRET_KEY_BASE32")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// GF_ADMIN_EMAIL
	err = viper.BindEnv("admin_email", "GF_ADMIN_EMAIL")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	//--------------------
	// GF_DOMAIN_BASE
	err = viper.BindEnv("domain_base", "GF_DOMAIN_BASE")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// GF_DOMAIN_ADMIN_BASE
	err = viper.BindEnv("domain_admin_base", "GF_DOMAIN_ADMIN_BASE")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("alchemy_api_key", "GF_ALCHEMY_SERVICE_ACC__API_KEY")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("images_use_new_storage_engine", "GF_IMAGES_USE_NEW_STORAGE_ENGINE")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}

	// ENV
	err = viper.BindEnv("ipfs_node_host", "GF_IPFS__NODE_HOST")
	if err != nil {
		fmt.Println("failed to bind ENV var to Viper config")
		panic(err)
	}
}
