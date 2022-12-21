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
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//-------------------------------------------------

type GFserviceInfo struct {
	Port_str                             string
	Mongodb_host_str                     string
	Mongodb_db_name_str                  string
	
	ImagesStoreLocalDirPathStr           string
	ImagesThumbnailsStoreLocalDirPathStr string
	VideoStoreLocalDirPathStr            string

	Media_domain_str                     string
	Images_main_s3_bucket_name_str       string
	AWS_access_key_id_str                string
	AWS_secret_access_key_str            string
	AWS_token_str                        string
	Templates_paths_map                  map[string]string
	Config_file_path_str                 string

	// AUTH_LOGIN_URL - url of the login page to which the system should
	//                  redirect users when email is confirmed.
	AuthLoginURLstr string
	KeyServer       *gf_identity_core.GFkeyServerInfo

	// NEW_STORAGE_ENGINE - flag indicating if the new image storage engine should be used
	UseNewStorageEngineBool bool

	// IPFS_NODE_HOST - host/gateway to use to connect to for IPFS operations
	IPFSnodeHostStr string
}

//-------------------------------------------------

type GFconfig struct {

	ImagesStoreLocalDirPathStr           string `yaml:"store_local_dir_path"`
	ImagesThumbnailsStoreLocalDirPathStr string `yaml:"thumbnails_store_local_dir_path"`
	VideoStoreLocalDirPathStr            string `yaml:"videos_store_local_dir_path"`

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

	ImagesFlowToS3bucketDefaultStr string            `yaml:"images_flow_to_s3_bucket_default"`
	ImagesFlowToS3bucketMap        map[string]string `yaml:"images_flow_to_s3_bucket"`

	UseNewStorageEngineBool bool `yaml:"use_new_storage_engine"`

	//------------------------
	// IPFS
	IPFSnodeHostStr string `yaml:"ipfs_node_host"`

	//------------------------
	// PLUGINS_PY
	// dir path from which to load Py plugins
	PluginsPyDirPathStr string `yaml:"plugins_py_dir_path"`

	//------------------------
}

//-------------------------------------------------

func ConfigGetS3bucketForFlow(pFlowNameStr string,
	pConfig *GFconfig) string {

	var s3bucketNameFinalStr string
	if s3bucketStr, ok := pConfig.ImagesFlowToS3bucketMap[pFlowNameStr]; !ok {
		s3bucketNameFinalStr = s3bucketStr
	} else {
		s3bucketNameFinalStr = pConfig.ImagesFlowToS3bucketDefaultStr
	}
	return s3bucketNameFinalStr
}

//-------------------------------------------------

func ConfigGet(pConfigPathStr string,
	pUseNewStorageEngineBool bool,
	pIPFSnodeHostStr         string,
	pRuntimeSys              *gf_core.RuntimeSys) (*GFconfig, *gf_core.GFerror) {

	configStr, err := ioutil.ReadFile(pConfigPathStr) 
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to read YAML config for gf_images",
			"file_read_error",
			map[string]interface{}{"config_path": pConfigPathStr,},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	config := &GFconfig{}
	
	// YAML - parse config file
	err = yaml.Unmarshal([]byte(configStr), config)
	if err != nil {

		gfErr := gf_core.ErrorCreate("failed to parse YAML config for gf_images",
			"yaml_decode_error",
			map[string]interface{}{"config_path": pConfigPathStr,},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	//------------------------
	// IPFS
	config.UseNewStorageEngineBool = pUseNewStorageEngineBool
	config.IPFSnodeHostStr         = pIPFSnodeHostStr

	//------------------------
	
	return config, nil
}