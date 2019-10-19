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

package gf_crawl_lib

import (
	"os"
	"testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)

var runtime_sys *gf_core.Runtime_sys
var crawler_runtime *gf_crawl_core.Gf_crawler_runtime

//---------------------------------------------------
func TestMain(m *testing.M) {
	runtime_sys, crawler_runtime = gf_crawl_core.T__init()
	if runtime_sys == nil || crawler_runtime == nil {
		return
	}
	v := m.Run()
	os.Exit(v)
}

//-------------------------------------------------
func Test__run_crawl_cycle(p_test *testing.T) {

	test__crawled_images_s3_bucket_name_str := "gf--test--discovered--img"
	test__crawler_images_local_dir_path_str := "./test_data/crawled_images"
	test__crawl_config_file_path_str        := "./test_data/config/test_crawl_config.yaml"

	test__run_crawl_cycle(p_test,
		test__crawler_images_local_dir_path_str,
		test__crawled_images_s3_bucket_name_str,
		test__crawl_config_file_path_str,
		crawler_runtime,
		runtime_sys)
}

//---------------------------------------------------
func test__run_crawl_cycle(p_test *testing.T,
	p_test__crawler_images_local_dir_path_str string,
	p_test__crawled_images_s3_bucket_name_str string,
	p_test__crawl_config_file_path_str        string,
	p_runtime                                 *gf_crawl_core.Gf_crawler_runtime,
	p_runtime_sys                             *gf_core.Runtime_sys) {

	crawlers_map, gf_err := gf_crawl_core.Get_all_crawlers(p_test__crawl_config_file_path_str, p_runtime_sys)
	if gf_err != nil {
		p_test.Errorf("failed to get all crawler definitions from config file [%s]", p_test__crawl_config_file_path_str)
		return
	}
	
	spew.Dump(crawlers_map)

	crawler := crawlers_map["gloflow"]

	gf_err = Run_crawler_cycle(crawler,
		p_test__crawler_images_local_dir_path_str,
		p_test__crawled_images_s3_bucket_name_str,
		p_runtime,
		p_runtime_sys)
	
	if gf_err != nil {
		p_test.Errorf("failed to run a crawler_cycle [%s], images_local_dir [%s] and s3 bucket [%s]",
			crawler.Name_str,
			p_test__crawler_images_local_dir_path_str,
			p_test__crawled_images_s3_bucket_name_str)
		return
	}
}