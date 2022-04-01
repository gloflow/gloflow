/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_images_core

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GFserviceInfo struct {
	Port_str                                   string
	Mongodb_host_str                           string
	Mongodb_db_name_str                        string
	Images_store_local_dir_path_str            string
	Images_thumbnails_store_local_dir_path_str string
	Media_domain_str                           string
	Images_main_s3_bucket_name_str             string
	AWS_access_key_id_str                      string
	AWS_secret_access_key_str                  string
	AWS_token_str                              string
	Templates_paths_map                        map[string]string
	Config_file_path_str                       string

	// AUTH_LOGIN_URL - url of the login page to which the system should
	//                  redirect users when email is confirmed.
	AuthLoginURLstr string
}

//-------------------------------------------------
type GF_config struct {

	Store_local_dir_path_str            string `yaml:"store_local_dir_path"`
	Thumbnails_store_local_dir_path_str string `yaml:"thumbnails_store_local_dir_path"`

	Media_domain_str        string `yaml:"media_domain"`
	Main_s3_bucket_name_str string `yaml:"main_s3_bucket_name"`

	//------------------------
	// FUNCTIONS - buckets for particular functions in that system

	// UPLOADED_IMAGES - this is a special dedicated bucket, separate from buckets for all other flows.
	//                   Mainly because users are pushing data to it directly and so we want to possibly handle
	//                   it in a separate way from other buckets that only have internal GF systems
	//                   uploading data to it.
	Uploaded_images_s3_bucket_str string `yaml:"uploaded_images_s3_bucket"`

	// BOOKMARKS_IMAGES - dedicated bucket for screenshots of bookmarks
	Bookmark_images_s3_bucket_str string `yaml:"bookmark_images_s3_bucket"`

	//------------------------

	Images_flow_to_s3_bucket_default_str string            `yaml:"images_flow_to_s3_bucket_default"`
	Images_flow_to_s3_bucket_map         map[string]string `yaml:"images_flow_to_s3_bucket"`
}

//-------------------------------------------------
func Config__get_s3_bucket_for_flow(p_flow_name_str string,
	p_config *GF_config) string {

	var s3_bucket_name_final_str string
	if s3_bucket_str, ok := p_config.Images_flow_to_s3_bucket_map[p_flow_name_str]; !ok {
		s3_bucket_name_final_str = s3_bucket_str
	} else {
		s3_bucket_name_final_str = p_config.Images_flow_to_s3_bucket_default_str
	}
	return s3_bucket_name_final_str
}

//-------------------------------------------------
func Config__get(p_config_path_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_config, *gf_core.Gf_error) {

	config_str, err := ioutil.ReadFile(p_config_path_str) 
	if err != nil {
		
		gf_err := gf_core.Error__create("failed to read YAML config for gf_images",
			"file_read_error",
			map[string]interface{}{"config_path": p_config_path_str,},
			err, "gf_images_core", p_runtime_sys)
		return nil, gf_err
	}


	config := &GF_config{}
	err = yaml.Unmarshal([]byte(config_str), config)
	if err != nil {

		gf_err := gf_core.Error__create("failed to parse YAML config for gf_images",
			"yaml_decode_error",
			map[string]interface{}{"config_path": p_config_path_str,},
			err, "gf_images_core", p_runtime_sys)
		return nil, gf_err
	}

	return config, nil

	// flow_to_s3_bucket_map := flows__get_mapping_to_s3_buckets()
	//
	// config := &Config{
	// 	Uploaded_images_s3bucket_str: "gf--uploaded--img",
	// 	Images_flow_to_s3_bucket_map: flow_to_s3_bucket_map,
	// }
	// return config
}