package gf_crawl_core

import (
	"github.com/globalsign/mgo/bson"
	"gf_core"
)
//-------------------------------------------------
func t__cleanup__test_page_imgs(p_test__crawler_name_str string,
						p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","t__utils.t__cleanup__test_page_imgs()")

	_,err := p_runtime_sys.Mongodb_coll.RemoveAll(bson.M{
							"t":               bson.M{"$in":[]string{"crawler_page_img","crawler_page_img_ref",},},
							"crawler_name_str":p_test__crawler_name_str,
						})
	if err != nil {
		panic(err)
	}


}
//-------------------------------------------------
func T__init() (*gf_core.Runtime_sys,*Crawler_runtime) {

	test__mongodb_host_str      := "127.0.0.1"
	test__mongodb_db_name_str   := "test_db"
	test__cluster_node_type_str := "master"
	
	log_fun  := gf_core.Init_log_fun()
	mongo_db := gf_core.Mongo__connect(test__mongodb_host_str,
							test__mongodb_db_name_str,
							log_fun)
	mongodb_coll := mongo_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str:"gf_crawl_tests",
		Log_fun:         log_fun,
		Mongodb_coll:    mongodb_coll,
	}


	esearch_client,gf_err := gf_core.Elastic__get_client(runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}


	//-------------
	//S3
	s3_info,gf_err := gf_core.S3__init(runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	//-------------

	crawler_runtime := &Crawler_runtime{
		Events_ctx:           nil,
		Esearch_client:       esearch_client,
		S3_info:              s3_info,
		Cluster_node_type_str:test__cluster_node_type_str,
	}

	return runtime_sys,crawler_runtime
}