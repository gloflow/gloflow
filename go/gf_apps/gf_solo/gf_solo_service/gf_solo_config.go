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
type GF_config struct {

	// DOMAIN - where this gf_solo instance is reachable on
	Domain_base_str string `mapstructure:"domain_base"`

	// PORTS
	Port_str         string `mapstructure:"port"`
	Port_admin_str   string `mapstructure:"port_admin"`
	Port_metrics_str string `mapstructure:"port_metrics"`

	// MONGODB - this is the dedicated mongodb DB
	Mongodb_host_str    string `mapstructure:"mongodb_host"`
	Mongodb_db_name_str string `mapstructure:"mongodb_db_name"`

	// ELASTICSEARCH
	Elasticsearch_host_str string `mapstructure:"elasticsearch_host"`

	// SENTRY_ENDPOINT
	Sentry_endpoint_str string `mapstructure:"sentry_endpoint"`

	// TEMPLATES
	Templates_paths_map map[string]string `mapstructure:"templates_paths"`

	//--------------------
	// GF_IMAGES
	Images__config_file_path_str string `mapstructure:"images__config_file_path"`
	/*Images__store_local_dir_path_str            string `mapstructure:"images__store_local_dir_path"`
	Images__thumbnails_store_local_dir_path_str string `mapstructure:"images__thumbnails_store_local_dir_path"`
	Images__main_s3_bucket_name_str             string `mapstructure:"images__main_s3_bucket_name"`

	Images__uploaded_s3_bucket_str        string            `mapstructure:"images__uploaded_s3_bucket"`
	Images__flow_to_s3_bucket_default_str string            `mapstructure:"images__flow_to_s3_bucket_default"`
	Images__flow_to_s3_bucket_map         map[string]string `mapstructure:"images__flow_to_s3_bucket"`*/

	//--------------------
	// GF_ANALYTICS

	Analytics__py_stats_dirs_lst []string `mapstructure:"analytics__py_stats_dirs"`
	Analytics__run_indexer_bool  bool     `mapstructure:"analytics__run_indexer"`

	Crawl__config_file_path_str      string `mapstructure:"crawl__config_file_path"`
	Crawl__cluster_node_type_str     string `mapstructure:"crawl__cluster_node_type"`
	Crawl__images_local_dir_path_str string `mapstructure:"crawl__images_local_dir_path"`

	//--------------------
	// AWS
	AWS_access_key_id_str     string `mapstructure:"aws_access_key_id"`
	AWS_secret_access_key_str string `mapstructure:"aws_secret_access_key"`
	AWS_token_str             string `mapstructure:"aws_token"`

	//--------------------
	// ADMIN_EMAIL
	Admin_email_str string `mapstructure:"admin_email"`

	//--------------------
}

//-------------------------------------------------------------
func Config__init(p_config_dir_path_str string,
	p_config_file_name_str string) (*GF_config, error) {


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
	config := &GF_config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}