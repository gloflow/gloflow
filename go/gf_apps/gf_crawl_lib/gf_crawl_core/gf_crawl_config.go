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
type GFcrawlConfig struct {
	Crawlers_defs_lst []GFcrawlerDef `yaml:"crawlers-defs"`
}

type GFcrawlerDef struct {
	Name_str      string `yaml:"name"`
	Start_url_str string `yaml:"start-url"`
}

//--------------------------------------------------
func Get_all_crawlers(pCrawlConfigFilePathStr string,
	pRuntimeSys *gf_core.RuntimeSys) (map[string]GFcrawlerDef, *gf_core.GFerror) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_config.Get_all_crawlers()")

	// no config file found, so use hard-coded crawler definitions
	if _, err := os.Stat(pCrawlConfigFilePathStr); os.IsNotExist(err) {

		crawlersMap := map[string]GFcrawlerDef{
			"gloflow.com": GFcrawlerDef{
				Name_str:      "gloflow.com",
				Start_url_str: "http://gloflow.com/",
			},
		}
		return crawlersMap, nil
	} else {

		//-------------
		// OPEN_CONFIG_FILE
		configByteLst, gfErr := ioutil.ReadFile(pCrawlConfigFilePathStr)
		if gfErr != nil {
			gfErr := gf_core.ErrorCreate("failed to read a local file to load the image",
				"file_read_error",
				map[string]interface{}{"crawl_config_file_path_str": pCrawlConfigFilePathStr,},
				gfErr, "gf_crawl_lib", pRuntimeSys)
			return nil, gfErr
		}
		
		//-------------
		// PARSE_YAML
		crawlConfig := GFcrawlConfig{}
		err = yaml.Unmarshal(configByteLst, &crawlConfig)
		if err != nil {
			gfErr := gf_core.Mongo__handle_error("failed to parse gf_crawler YAML config file",
				"yaml_decode_error",
				map[string]interface{}{"crawl_config_file_path_str": pCrawlConfigFilePathStr,},
				err, "gf_crawl_core", pRuntimeSys)
			return nil, gfErr
		}

		// index crawler_defs by name
		crawlersMap := map[string]GFcrawlerDef{}
		for _, crawlerDef := range crawlConfig.Crawlers_defs_lst {
			crawlersMap[crawlerDef.Name_str] = crawlerDef
		}
		return crawlersMap, nil
	}
	return nil, nil
}