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
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)
//-------------------------------------------------
func Test__stages(p_test *testing.T) {

	test__images_s3_bucket_name_str         := "gf--test"
	test__crawler_images_local_dir_path_str := "./test_data/crawled_images"

	runtime_sys,crawler_runtime := gf_crawl_core.T__init()

	test__stages(test__crawler_images_local_dir_path_str,
		test__images_s3_bucket_name_str,
		crawler_runtime,
		runtime_sys)
}
//-------------------------------------------------
func test__stages(p_test__crawler_images_local_dir_path_str string,
	p_test__images_s3_bucket_name_str string,
	p_runtime                         *gf_crawl_core.Gf_crawler_runtime,
	p_runtime_sys                     *gf_core.Runtime_sys) {




	/*fetch_url(p_url_str string,
		p_link             *Crawler_page_outgoing_link,
		p_cycle_run_id_str string,
		p_crawler_name_str string,
		p_runtime          *Crawler_runtime,
		p_runtime_sys      *gf_core.Runtime_sys)


	
	crawled_images_lst,crawled_images_refs_lst := images__stage__pull_image_links(p_url_fetch,
																			p_crawler_name_str,
																			p_cycle_run_id_str,
																			p_runtime,
																			p_runtime_sys)*/
}