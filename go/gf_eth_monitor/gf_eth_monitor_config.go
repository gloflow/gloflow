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
	"strings"
	"github.com/spf13/viper"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_core"
)

//-------------------------------------------------------------
func config__init(p_config_dir_path_str string,
	p_config_file_name_str string) (*gf_eth_monitor_core.GF_config, error) {


	config_name_str := strings.Split(p_config_file_name_str, ".")[0] // viper expects just the file name, without extension
	
	// FILE
	viper.AddConfigPath(p_config_dir_path_str)
	viper.SetConfigName(config_name_str)
	
	//--------------------
	// ENV_PREFIX - "GF" becomes "GF_" - prefix expected in all recognized ENV vars.
	viper.SetEnvPrefix("GF")

	// ENV_VARS
	// IMPORTANT!! - enable Viper parsing ENV vars.
	//               all config members that have their mapstructure name for Viper config, 
	//               also have a corresponding ENV var name thats generated for them by
	//               upper-casing their name.
	viper.AutomaticEnv()
	
	//--------------------

	// LOAD
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// CONFIG
	config := &gf_eth_monitor_core.GF_config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}