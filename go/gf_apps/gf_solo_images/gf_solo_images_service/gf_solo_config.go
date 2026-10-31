/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

package gf_solo_images_service

import (
	"fmt"
	"strings"
	"github.com/spf13/viper"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------
// CONFIG_INIT

/*
load config for gf_solo_images from all the sources.
Viper lib is used for config loading.
the config struct is the shared gf_core.GFconfig (gf_images only uses a subset of its fields).
configs are loaded from:
- ENV vars
- config file (gf_solo_images.yaml)
- via CLI args
*/
func ConfigInit(pConfigDirPathStr string,
	pConfigFileNameStr string,
	pRuntimeSys        *gf_core.RuntimeSys) (*gf_core.GFconfig, error) {

	configNameStr := strings.Split(pConfigFileNameStr, ".")[0] // viper expects just the file name, without extension

	// FILE
	viper.AddConfigPath(pConfigDirPathStr)
	viper.SetConfigName(configNameStr)

	//--------------------
	// ENV_VARS
	// ENV_PREFIX - "GF" becomes "GF_" - prefix expected in all recognized ENV vars.
	viper.SetEnvPrefix("GF")

	// IMPORTANT!! - enable Viper parsing ENV vars.
	viper.AutomaticEnv()

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

	// IMPORTANT!! - env is needed by all subsequent loading/init, so set it right away.
	pRuntimeSys.EnvStr = config.EnvStr

	return config, nil
}

//-------------------------------------------------------------

func bindEnvVars() {

	// config key -> ENV var name
	envVarsMap := map[string]string{
		"env":                           "GF_ENV",
		"port":                          "GF_PORT",
		"port_metrics":                  "GF_PORT_METRICS",
		"mongodb_host":                  "GF_MONGODB_HOST",
		"mongodb_db_name":               "GF_MONGODB_DB_NAME",
		"sql_user_name":                 "GF_SQL_USER_NAME",
		"sql_pass":                      "GF_SQL_PASS",
		"sql_host":                      "GF_SQL_HOST",
		"sql_db_name":                   "GF_SQL_DB_NAME",
		"sentry_dsn":                    "SENTRY_DSN",
		"auth_subsystem_type":           "GF_AUTH_SUBSYSTEM_TYPE",
		"domain_base":                   "GF_DOMAIN_BASE",
		"domain_for_auth_cookies":       "GF_DOMAIN_FOR_AUTH_COOKIES",
		"images_use_new_storage_engine": "GF_IMAGES_USE_NEW_STORAGE_ENGINE",
		"ipfs_node_host":                "GF_IPFS__NODE_HOST",
	}
	for keyStr, envVarStr := range envVarsMap {
		err := viper.BindEnv(keyStr, envVarStr)
		if err != nil {
			fmt.Println("failed to bind ENV var to Viper config")
			panic(err)
		}
	}
}
