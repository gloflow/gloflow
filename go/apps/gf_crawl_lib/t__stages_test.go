package gf_crawl_lib

import (
	"testing"
	"gf_core"
	"apps/gf_crawl_lib/gf_crawl_core"
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
			p_runtime                         *gf_crawl_core.Crawler_runtime,
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