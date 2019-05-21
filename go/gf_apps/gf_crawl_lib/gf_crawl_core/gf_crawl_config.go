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

package gf_crawl_core

import (
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
type Gf_crawl_config struct {
	Crawlers_defs_lst []Gf_crawler_def `yaml:"crawlers-defs"`
}

type Gf_crawler_def struct {
	Name_str      string `yaml:"name"`
	Start_url_str string `yaml:"start-url"`
}

//--------------------------------------------------
func Get_all_crawlers(p_crawl_config_file_path_str string, p_runtime_sys *gf_core.Runtime_sys) (map[string]Gf_crawler_def, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_config.Get_all_crawlers()")

	//no config file found, so use hard-coded crawler definitions
	if _, err := os.Stat(p_crawl_config_file_path_str); os.IsNotExist(err) {

		crawlers_map := map[string]Gf_crawler_def{
			"gloflow.com": Gf_crawler_def{
				Name_str:      "gloflow.com",
				Start_url_str: "http://gloflow.com/",
			},
		}
		return crawlers_map, nil
	} else {

		//-------------
		//OPEN_CONFIG_FILE
		config_byte_lst, fs_err := ioutil.ReadFile(p_crawl_config_file_path_str)
		if fs_err != nil {
			gf_err := gf_core.Error__create("failed to read a local file to load the image",
				"file_read_error",
				map[string]interface{}{"crawl_config_file_path_str": p_crawl_config_file_path_str,},
				fs_err, "gf_crawl_lib", p_runtime_sys)
			return nil, gf_err
		}
		//-------------
		//PARSE_YAML
		crawl_config := Gf_crawl_config{}
		err = yaml.Unmarshal(config_byte_lst, &crawl_config)
		if err != nil {
			gf_err := gf_core.Mongo__handle_error("failed to parse gf_crawler YAML config file",
				"mongodb_update_error",
				map[string]interface{}{"crawl_config_file_path_str": p_crawl_config_file_path_str,},
				err, "gf_crawl_core", p_runtime_sys)
			return nil, gf_err
		}

		//index crawler_defs by name
		crawlers_map := map[string]Gf_crawler_def{}
		for _, crawler_def := range crawl_config.Crawlers_defs_lst {
			crawlers_map[crawler_def.Name_str] = crawler_def
		}
		return crawlers_map, nil
	}
	return nil, nil
}