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
	"strings"
	"github.com/spf13/viper"
)

//-------------------------------------------------------------

type GFconfig struct {

	// ENVIRONMENT
	EnvStr string `mapstructure:"env"`

	// DOMAIN - where this gf_solo instance is reachable on
	DomainBaseStr      string `mapstructure:"domain_base"`
	DomainAdminBaseStr string `mapstructure:"domain_admin_base"`
	
	// PORTS
	PortStr        string `mapstructure:"port"`
	PortAdminStr   string `mapstructure:"port_admin"`
	PortMetricsStr string `mapstructure:"port_metrics"`

	// MONGODB - this is the dedicated mongodb DB
	MongoHostStr   string `mapstructure:"mongodb_host"`
	MongoDBnameStr string `mapstructure:"mongodb_db_name"`

	// SQL
	SQLuserNameStr string `mapstructure:"sql_user_name"`
	SQLpassStr     string `mapstructure:"sql_pass"`
	SQLhostStr     string `mapstructure:"sql_host"`
	SQLdbNameStr   string `mapstructure:"sql_db_name"`

	// ELASTICSEARCH
	ElasticsearchHostStr string `mapstructure:"elasticsearch_host"`

	// SENTRY_ENDPOINT
	SentryEndpointStr string `mapstructure:"sentry_endpoint"`

	// TEMPLATES
	TemplatesPathsMap map[string]string `mapstructure:"templates_paths"`

	//--------------------
	// IDENTITY
	AuthSubsystemTypeStr       string `mapstructure:"auth_subsystem_type"`
	AdminMFAsecretKeyBase32str string `mapstructure:"admin_mfa_secret_key_base32"`
	AdminEmailStr              string `mapstructure:"admin_email"`

	//--------------------
	// GF_IMAGES
	ImagesConfigFilePathStr string `mapstructure:"images__config_file_path"`

	//--------------------
	// GF_ANALYTICS

	AnalyticsPyStatsDirsLst []string `mapstructure:"analytics__py_stats_dirs"`
	AnalyticsRunIndexerBool bool     `mapstructure:"analytics__run_indexer"`

	CrawlConfigFilePathStr     string `mapstructure:"crawl__config_file_path"`
	CrawlClusterNodeTypeStr    string `mapstructure:"crawl__cluster_node_type"`
	CrawlImagesLocalDirPathStr string `mapstructure:"crawl__images_local_dir_path"`

	//--------------------
	// ALCHEMY
	AlchemyAPIkeyStr string `mapstructure:"alchemy_api_key"`

	//--------------------
	// NEW_STORAGE_ENGINE - flag indicating if the new image storage engine should be used
	ImagesUseNewStorageEngineBool bool `mapstructure:"images_use_new_storage_engine"`

	// IPFS
	IPFSnodeHostStr string `mapstructure:"ipfs_node_host"`

	//--------------------
}

//-------------------------------------------------------------

func ConfigInit(pConfigDirPathStr string,
	pConfigFileNameStr string) (*GFconfig, error) {

	configNameStr := strings.Split(pConfigFileNameStr, ".")[0] // viper expects just the file name, without extension
	
	// FILE
	viper.AddConfigPath(pConfigDirPathStr)
	viper.SetConfigName(configNameStr)
	
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
	config := &GFconfig{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}