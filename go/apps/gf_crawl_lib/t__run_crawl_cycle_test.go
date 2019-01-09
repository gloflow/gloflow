package gf_crawl_lib

import (
	"testing"
	"gf_core"
	"apps/gf_crawl_lib/gf_crawl_core"
)
//-------------------------------------------------
func Test__run_crawl_cycle(p_test *testing.T) {

	test__crawled_images_s3_bucket_name_str := "gf--test--discovered--img"
	test__crawler_images_local_dir_path_str := "./test_data/crawled_images"



	runtime_sys,crawler_runtime := gf_crawl_core.T__init()


	test__run_crawl_cycle(test__crawler_images_local_dir_path_str,
			test__crawled_images_s3_bucket_name_str,
			crawler_runtime,
			runtime_sys)
}
//---------------------------------------------------
func test__run_crawl_cycle(p_test__crawler_images_local_dir_path_str string,
			p_test__crawled_images_s3_bucket_name_str string,
			p_runtime                                 *gf_crawl_core.Crawler_runtime,
			p_runtime_sys                             *gf_core.Runtime_sys) {



	crawlers_map := Get_all_crawlers()
	crawler      := crawlers_map["r2-r.tumblr"]


	gf_err := Run_crawler_cycle(crawler,
						p_test__crawler_images_local_dir_path_str,
						p_test__crawled_images_s3_bucket_name_str,
						p_runtime,
						p_runtime_sys)
	
	if gf_err != nil {
		panic(gf_err.Error)
	}



}